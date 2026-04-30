package dns

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

const (
	cloudnsEndpoint string = "https://api.cloudns.net/dns/"
)

// ClouDNS ClouDNS
type ClouDNS struct {
	DNS        config.DNS
	Domains    config.Domains
	TTL        string
	httpClient *http.Client
}

// ClouDNSRecord record
type ClouDNSRecord struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Host string `json:"host"`
	Value string `json:"record"`
}

// ClouDNSResp generic response
type ClouDNSResp struct {
	Status            string `json:"status"`
	StatusDescription string `json:"statusDescription"`
}

// Init 初始化
func (cl *ClouDNS) Init(dnsConf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache) {
	cl.Domains.Ipv4Cache = ipv4cache
	cl.Domains.Ipv6Cache = ipv6cache
	cl.DNS = dnsConf.DNS
	cl.Domains.GetNewIp(dnsConf)
	if dnsConf.TTL == "" {
		// Default 3600 (ClouDNS minimum for some plans)
		cl.TTL = "3600"
	} else {
		cl.TTL = dnsConf.TTL
	}
	cl.httpClient = dnsConf.GetHTTPClient()
}

// AddUpdateDomainRecords 添加或更新IPv4/IPv6记录
func (cl *ClouDNS) AddUpdateDomainRecords() config.Domains {
	cl.addUpdateDomainRecords("A")
	cl.addUpdateDomainRecords("AAAA")
	return cl.Domains
}

func (cl *ClouDNS) addUpdateDomainRecords(recordType string) {
	ipAddr, domains := cl.Domains.GetNewIpResult(recordType)

	if ipAddr == "" {
		return
	}

	for _, domain := range domains {
		var records map[string]ClouDNSRecord
		// Get current record information
		params := url.Values{}
		params.Set("auth-id", cl.DNS.ID)
		params.Set("auth-password", cl.DNS.Secret)
		params.Set("domain-name", domain.DomainName)
		params.Set("host", domain.GetSubDomain())
		params.Set("type", recordType)

		err := cl.request("list-records.json", params, &records)
		if err != nil {
			util.Log("查询域名 %s 信息发生异常! %s", domain, err)
			domain.UpdateStatus = config.UpdatedFailed
			return
		}

		var recordSelected *ClouDNSRecord
		if len(records) > 0 {
			// Find the first record of the matching type and host
			for _, r := range records {
				if r.Type == recordType && r.Host == domain.GetSubDomain() {
					recordSelected = &r
					break
				}
			}
		}

		if recordSelected != nil {
			// Exist, modify
			cl.modify(recordSelected, domain, recordType, ipAddr)
		} else {
			// Not exist, create
			cl.create(domain, recordType, ipAddr)
		}
	}
}

// create
func (cl *ClouDNS) create(domain *config.Domain, recordType string, ipAddr string) {
	params := url.Values{}
	params.Set("auth-id", cl.DNS.ID)
	params.Set("auth-password", cl.DNS.Secret)
	params.Set("domain-name", domain.DomainName)
	params.Set("host", domain.GetSubDomain())
	params.Set("type", recordType)
	params.Set("record", ipAddr)
	params.Set("ttl", cl.TTL)

	var result ClouDNSResp
	err := cl.request("add-record.json", params, &result)

	if err != nil {
		util.Log("新增域名解析 %s 失败! 异常信息: %s", domain, err)
		domain.UpdateStatus = config.UpdatedFailed
		return
	}

	if result.Status == "Success" {
		util.Log("新增域名解析 %s 成功! IP: %s", domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	} else {
		util.Log("新增域名解析 %s 失败! 异常信息: %s", domain, result.StatusDescription)
		domain.UpdateStatus = config.UpdatedFailed
	}
}

// modify
func (cl *ClouDNS) modify(recordSelected *ClouDNSRecord, domain *config.Domain, recordType string, ipAddr string) {
	// Same, no change
	if recordSelected.Value == ipAddr {
		util.Log("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
		return
	}

	params := url.Values{}
	params.Set("auth-id", cl.DNS.ID)
	params.Set("auth-password", cl.DNS.Secret)
	params.Set("domain-name", domain.DomainName)
	params.Set("record-id", recordSelected.ID)
	params.Set("host", domain.GetSubDomain())
	params.Set("record", ipAddr)
	params.Set("ttl", cl.TTL)

	var result ClouDNSResp
	err := cl.request("modify-record.json", params, &result)

	if err != nil {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, err)
		domain.UpdateStatus = config.UpdatedFailed
		return
	}

	if result.Status == "Success" {
		util.Log("更新域名解析 %s 成功! IP: %s", domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	} else {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, result.StatusDescription)
		domain.UpdateStatus = config.UpdatedFailed
	}
}

// request
func (cl *ClouDNS) request(action string, params url.Values, result interface{}) (err error) {
	resp, err := cl.httpClient.PostForm(cloudnsEndpoint+action, params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(result)
	return err
}
