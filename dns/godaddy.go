package dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
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
	throttle  util.Throttle
	client    *http.Client
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
	g.throttle, _ = util.GetThrottle(55)
	g.client = util.CreateHTTPClient()
}

func (g *GoDaddyDNS) updateDomainRecord(rType string, data string, domains []*config.Domain) {
	for _, domain := range domains {
		domain.UpdateStatus = config.UpdatedFailed
		if _, err := g.sendReq(http.MethodPut, rType, domain, &godaddyRecords{godaddyRecord{
			Data: data,
			Name: domain.SubDomain,
			TTL:  g.ttl,
			Type: rType,
		}}); err == nil {
			domain.UpdateStatus = config.UpdatedSuccess
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

func (g *GoDaddyDNS) sendReq(method string, rType string, domain *config.Domain, data any) (*godaddyRecords, error) {
	for !g.throttle.Try() {
		runtime.Gosched()
	}
	var body *bytes.Buffer
	if data != nil {
		if buffer, err := json.Marshal(data); err != nil {
			return nil, err
		} else {
			body = bytes.NewBuffer(buffer)
		}
	}
	path := fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/%s/%s",
		domain.DomainName, rType, domain.SubDomain)
	log.Printf("向godaddy发送请求，请求地址为%s", path)
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	req.Header = g.header
	resp, err := g.client.Do(req)
	res := &godaddyRecords{}
	err = util.GetHTTPResponse(resp, path, err, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
