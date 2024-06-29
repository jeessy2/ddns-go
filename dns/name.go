package dns

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	var items NameComItemList
	err := cloud.request("GET", api, nil, &items)
	return items, err
}

// request 统一请求接口
func (cloud *Name) request(method string, url string, data interface{}, result interface{}) (err error) {
	body, _ := json.Marshal(data)
	request, _ := http.NewRequest(method, url, bytes.NewReader(body))
	request.Header = cloud.header
	request.SetBasicAuth(cloud.DNS.ID, cloud.DNS.Secret)
	response, err := cloud.client.Do(request)
	if err != nil {
		return err
	}
	body, err = util.GetHTTPResponseOrg(response, err)
	if err != nil {
		return err
	}
	if result != nil {
		err = json.Unmarshal(body, result)
	}
	return nil
}
func (cloud *Name) addDomainRecord(record NameComItem, domain *config.Domain, ipAddr string) {
	api := fmt.Sprintf("%s/v4/domains/%s/records", nameAPI, record.Domain)
	item := NameComItem{
		Host: record.Host,
		Type: record.Type,
		IP:   record.IP,
		TTL:  cloud.ttl,
	}
	err := cloud.request("POST", api, item, nil)
	if err != nil {
		util.Log("添加域名解析 %s 失败! 异常信息: %v", domain, err)
		domain.UpdateStatus = config.UpdatedFailed
	} else {
		util.Log("添加域名解析 %s 成功! IP: %s", domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	}
}
func (cloud *Name) updateDomainRecord(record NameComItem, domain *config.Domain, ipAddr string) {
	api := fmt.Sprintf("%s/v4/domains/%s/records/%d", nameAPI, record.Domain, record.ID)
	item := NameComItem{
		Host: record.Host,
		Type: record.Type,
		IP:   record.IP,
		TTL:  cloud.ttl,
	}
	err := cloud.request("PUT", api, item, nil)
	if err != nil {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, err)
		domain.UpdateStatus = config.UpdatedFailed
	} else {
		util.Log("更新域名解析 %s 成功! IP: %s", domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	}
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
		if len(records.Records) == 0 {
			cloud.addDomainRecord(NameComItem{
				Domain: domain.DomainName,
				Host:   domain.SubDomain,
				IP:     ipAddr,
				Type:   recordType,
			}, domain, ipAddr)
			continue
		}
		hasMatch := false
		for _, record := range records.Records {
			if record.Host == domain.SubDomain {
				hasMatch = true
				// ip 没有变化，不需要重新解析
				if record.IP == cloud.lastIpv4 {
					util.Log("你的IPv4未变化, 未触发 %s 请求", "name.com")
					break
				}
				cloud.updateDomainRecord(NameComItem{
					Domain: domain.DomainName,
					Host:   domain.SubDomain,
					IP:     ipAddr,
					Type:   recordType,
					ID:     record.ID,
				}, domain, ipAddr)
				break
			}
		}
		if !hasMatch {
			cloud.addDomainRecord(NameComItem{
				Domain: domain.DomainName,
				Host:   domain.SubDomain,
				IP:     ipAddr,
				Type:   recordType,
			}, domain, ipAddr)
		}
	}

}
