package dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jeessy2/ddns-go/v4/config"
	"github.com/jeessy2/ddns-go/v4/util"
)

type godaddyRecord struct {
	Data string `json:"data"`
	Name string `json:"name"`
	TTL  int    `json:"ttl"`
	Type string `json:"type"`
}

type godaddyRecords []godaddyRecord

type GoDaddyDNS struct {
	dnsConfig config.DNSConfig
	domains   config.Domains
	ttl       int
	header    http.Header
	client    *http.Client
	lastIpv4  string
	lastIpv6  string
}

func (g *GoDaddyDNS) Init(conf *config.Config) {
	g.dnsConfig = conf.DNS
	g.domains.GetNewIp(conf)
	g.ttl = 600
	if val, err := strconv.Atoi(conf.TTL); err == nil {
		g.ttl = val
	}
	g.header = map[string][]string{
		"Authorization": {fmt.Sprintf("sso-key %s:%s", g.dnsConfig.ID, g.dnsConfig.Secret)},
		"Content-Type":  {"application/json"},
	}

	g.client = util.CreateHTTPClient()
}

func (g *GoDaddyDNS) updateDomainRecord(recordType string, ipAddr string, domains []*config.Domain) {
	if ipAddr == "" {
		return
	}

	if recordType == "A" {
		if g.lastIpv4 == ipAddr {
			log.Println("你的IPv4未变化, 未触发请求")
			return
		}
		g.lastIpv4 = ipAddr
	} else {
		if g.lastIpv6 == ipAddr {
			log.Println("你的IPv6未变化, 未触发请求")
			return
		}
		g.lastIpv6 = ipAddr
	}

	for _, domain := range domains {
		err := g.sendReq(http.MethodPut, recordType, domain, &godaddyRecords{godaddyRecord{
			Data: ipAddr,
			Name: domain.SubDomain,
			TTL:  g.ttl,
			Type: recordType,
		}})
		if err == nil {
			log.Printf("更新域名解析 %s 成功! IP: %s", domain, ipAddr)
			domain.UpdateStatus = config.UpdatedSuccess
		} else {
			log.Printf("更新域名解析 %s 失败！", domain)
			domain.UpdateStatus = config.UpdatedFailed
		}
	}
}

func (g *GoDaddyDNS) AddUpdateDomainRecords() config.Domains {
	if ipv4Addr, ipv4Domains := g.domains.GetNewIpResult("A"); ipv4Addr != "" {
		g.updateDomainRecord("A", ipv4Addr, ipv4Domains)
	}
	if ipv6Addr, ipv6Domains := g.domains.GetNewIpResult("AAAA"); ipv6Addr != "" {
		g.updateDomainRecord("AAAA", ipv6Addr, ipv6Domains)
	}
	return g.domains
}

func (g *GoDaddyDNS) sendReq(method string, rType string, domain *config.Domain, data *godaddyRecords) error {

	var body *bytes.Buffer
	if data != nil {
		if buffer, err := json.Marshal(data); err != nil {
			return err
		} else {
			body = bytes.NewBuffer(buffer)
		}
	}
	path := fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/%s/%s",
		domain.DomainName, rType, domain.SubDomain)

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return err
	}
	req.Header = g.header
	resp, err := g.client.Do(req)
	_, err = util.GetHTTPResponseOrg(resp, path, err)
	return err
}
