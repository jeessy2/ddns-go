package dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

const spaceshipAPI = "https://spaceship.dev/api/v1/dns/records"
const maxRecords = 500

type Spaceship struct {
	domains config.Domains
	header  http.Header
	ttl     int
}

func (s *Spaceship) Init(dnsConf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache) {
	s.domains.Ipv4Cache = ipv4cache
	s.domains.Ipv6Cache = ipv6cache
	s.domains.GetNewIp(dnsConf)

	s.ttl = 600
	if val, err := strconv.Atoi(dnsConf.TTL); err == nil {
		s.ttl = val
	}
	s.header = http.Header{
		"X-API-Key":    {dnsConf.DNS.ID},
		"X-API-Secret": {dnsConf.DNS.Secret},
		"Content-Type": {"application/json"},
	}
}

func (s *Spaceship) AddUpdateDomainRecords() (domains config.Domains) {
	for _, recordType := range []string{"A", "AAAA"} {
		ip, domains := s.domains.GetNewIpResult(recordType)
		if ip == "" {
			continue
		}
		for _, domain := range domains {
			hasUpdated, err := s.updateRecord(recordType, ip, domain)
			if err != nil {
				util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, err)
				domain.UpdateStatus = config.UpdatedFailed
				continue
			}
			if !hasUpdated {
				util.Log("你的IP %s 没有变化, 域名 %s", ip, domain)
			} else {
				util.Log("更新域名解析 %s 成功! IP: %s", domain, ip)
				domain.UpdateStatus = config.UpdatedSuccess
			}
		}
	}
	return s.domains
}

func (s *Spaceship) request(domain *config.Domain, method string, query url.Values, payload []byte) (response []byte, err error) {
	url := fmt.Sprintf("%s/%s", spaceshipAPI, domain.DomainName)
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return
	}
	req.Header = s.header
	req.URL.RawQuery = query.Encode()

	cli := util.CreateHTTPClient()
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	response, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	type DataItem struct {
		Field   string `json:"field"`
		Details string `json:"details"`
	}

	type ErrorResponse struct {
		Detail string      `json:"detail"`
		Data   *[]DataItem `json:"data,omitempty"`
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		var e ErrorResponse
		err = json.Unmarshal(response, &e)
		if err != nil {
			return
		}
		err = fmt.Errorf("request error: %s", e.Detail)
		return
	}

	return
}

func (s *Spaceship) createRecord(recordType string, ip string, domain *config.Domain) (err error) {
	type Item struct {
		Type    string `json:"type"`
		Address string `json:"address"`
		Name    string `json:"name"`
		TTL     int    `json:"ttl"`
	}

	type Payload struct {
		Force bool   `json:"force"`
		Items []Item `json:"items"`
	}

	payload := Payload{
		Force: true,
		Items: []Item{
			{
				Type:    recordType,
				Address: ip,
				Name:    domain.SubDomain,
				TTL:     s.ttl,
			},
		},
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	_, err = s.request(domain, "PUT", url.Values{}, data)
	return
}

func (s *Spaceship) getRecords(recordType string, domain *config.Domain) (ips []string, err error) {
	type Group struct {
		Type string `json:"type"`
	}

	type Item struct {
		Type    string `json:"type"`
		Address string `json:"address"`
		Name    string `json:"name"`
		TTL     int    `json:"ttl"`
		Group   Group  `json:"group"`
	}

	type Response struct {
		Items []Item `json:"items"`
		Total int    `json:"total"`
	}

	resp, err := s.request(domain, "GET", url.Values{"take": {strconv.Itoa(maxRecords)}, "skip": {"0"}}, []byte{})
	if err != nil {
		return
	}

	var response Response
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return
	}

	if response.Total > maxRecords {
		err = fmt.Errorf("could not fetch all %d records in a one request", response.Total)
		return
	}

	for _, item := range response.Items {
		if item.Type == recordType && item.Name == domain.SubDomain {
			ips = append(ips, item.Address)
		}
	}
	return
}

func (s *Spaceship) deleteRecords(recordType string, domain *config.Domain, ips []string) (err error) {
	if len(ips) == 0 {
		return
	}

	if len(ips) > maxRecords {
		err = fmt.Errorf("could not delete all %d records in a one request", len(ips))
		return
	}

	type Item struct {
		Type    string `json:"type"`
		Address string `json:"address"`
		Name    string `json:"name"`
	}
	var payload []Item
	for _, ip := range ips {
		payload = append(payload, Item{
			Type:    recordType,
			Address: ip,
			Name:    domain.SubDomain,
		})
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	_, err = s.request(domain, "DELETE", url.Values{}, data)
	return
}

func (s *Spaceship) updateRecord(recordType string, ip string, domain *config.Domain) (hasUpdated bool, err error) {
	ips, err := s.getRecords(recordType, domain)
	if err != nil {
		return
	}
	if len(ips) == 1 && ips[0] == ip {
		return
	}
	err = s.deleteRecords(recordType, domain, ips)
	if err != nil {
		return
	}
	err = s.createRecord(recordType, ip, domain)
	hasUpdated = true
	return
}
