package dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

// 文档 https://www.name.com/zh-cn/api-docs

const (
	nameAPI string = "https://api.name.com"
)

type NameComItemList struct {
	Records []NameComItem `json:"records"`
}

type NameComError struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

type Name struct {
	DNS      config.DNS
	Domains  config.Domains
	lastIpv4 string
	lastIpv6 string
	ttl      int
	header   http.Header
	client   *http.Client
}

func (cloud *Name) Init(dnsConf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache) {
	cloud.Domains.Ipv4Cache = ipv4cache
	cloud.Domains.Ipv6Cache = ipv6cache
	cloud.lastIpv4 = ipv4cache.Addr
	cloud.lastIpv6 = ipv6cache.Addr
	cloud.DNS = dnsConf.DNS
	cloud.Domains.GetNewIp(dnsConf)
	cloud.ttl = 600
	if val, err := strconv.Atoi(dnsConf.TTL); err == nil {
		cloud.ttl = val
	}
	cloud.header = map[string][]string{
		"Content-Type": {"application/json"},
	}
	cloud.client = util.CreateHTTPClient()
}

// 添加或更新IPv4/IPv6记录
func (cloud *Name) AddUpdateDomainRecords() (domains config.Domains) {
	if ipv4Addr, ipv4Domains := cloud.Domains.GetNewIpResult("A"); ipv4Addr != "" {
		cloud.addOrUpdateDomain("A", ipv4Addr, ipv4Domains)
	}
	if ipv6Addr, ipv6Domains := cloud.Domains.GetNewIpResult("AAAA"); ipv6Addr != "" {
		cloud.addOrUpdateDomain("AAAA", ipv6Addr, ipv6Domains)
	}
	return cloud.Domains
}

type NameComItem struct {
	ID     int64  `json:"id,omitempty"`
	Domain string `json:"domainName,omitempty"`
	Host   string `json:"host,omitempty"`
	Type   string `json:"type,omitempty"`
	IP     string `json:"answer,omitempty"`
	TTL    int    `json:"ttl,omitempty"`
}

func (cloud *Name) getDomainRecords(domain string) (NameComItemList, error) {
	api := fmt.Sprintf("%s/v4/domains/%s/records", nameAPI, domain)
	request, _ := http.NewRequest("GET", api, nil)
	request.SetBasicAuth(cloud.DNS.ID, cloud.DNS.Secret)
	response, err := cloud.client.Do(request)
	if err != nil {
		return NameComItemList{}, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return NameComItemList{}, err
	}
	if response.StatusCode != 200 {
		return NameComItemList{}, fmt.Errorf(string(body))
	}
	var items NameComItemList
	err = json.Unmarshal(body, &items)
	return items, err
}

func (cloud *Name) addDomainRecord(domain NameComItem) error {
	api := fmt.Sprintf("%s/v4/domains/%s/records", nameAPI, domain.Domain)
	item := NameComItem{
		Host: domain.Host,
		Type: domain.Type,
		IP:   domain.IP,
		TTL:  cloud.ttl,
	}
	body, _ := json.Marshal(item)
	request, _ := http.NewRequest("POST", api, bytes.NewReader(body))
	request.Header = cloud.header
	request.SetBasicAuth(cloud.DNS.ID, cloud.DNS.Secret)
	response, err := cloud.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf(string(responseBody))
	}
	return nil
}
func (cloud *Name) updateDomainRecord(id int64, domain NameComItem) error {
	api := fmt.Sprintf("%s/v4/domains/%s/records/%d", nameAPI, domain.Domain, id)
	item := NameComItem{
		Host: domain.Host,
		Type: domain.Type,
		IP:   domain.IP,
		TTL:  cloud.ttl,
	}
	body, _ := json.Marshal(item)
	request, _ := http.NewRequest("PUT", api, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(cloud.DNS.ID, cloud.DNS.Secret)
	response, err := cloud.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf(string(responseBody))
	}
	return nil
}

func (cloud *Name) addOrUpdateDomain(recordType string, ipAddr string, domains []*config.Domain) {

	if ipAddr == "" {
		return
	}
	for _, domain := range domains {
		records, apiErr := cloud.getDomainRecords(domain.DomainName)
		if apiErr != nil {
			util.Log("查询域名信息发生异常! %s", apiErr)
			domain.UpdateStatus = config.UpdatedFailed
			continue
		}
		var err error
		if len(records.Records) == 0 {
			err = cloud.addDomainRecord(NameComItem{
				Domain: domain.DomainName,
				Host:   domain.SubDomain,
				IP:     ipAddr,
				Type:   recordType,
			})
		} else {
			hasMatch := false
			for _, record := range records.Records {
				if record.Host == domain.SubDomain {
					// ip 没有变化，不需要重新解析
					if record.IP == cloud.lastIpv4 {
						util.Log("你的IPv4未变化, 未触发 %s 请求", "name.com")
						hasMatch = true
						break
					}
					err = cloud.updateDomainRecord(record.ID, NameComItem{
						Domain: domain.DomainName,
						Host:   domain.SubDomain,
						IP:     ipAddr,
						Type:   recordType,
					})
					hasMatch = true
					break
				}
			}
			if !hasMatch {
				err = cloud.addDomainRecord(NameComItem{
					Domain: domain.DomainName,
					Host:   domain.SubDomain,
					IP:     ipAddr,
					Type:   recordType,
				})
			}
		}
		if err == nil {
			util.Log("更新域名解析 %s 成功! IP: %s", domain, ipAddr)
			domain.UpdateStatus = config.UpdatedSuccess
		} else {
			util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, err)
			domain.UpdateStatus = config.UpdatedFailed
		}
	}

}
