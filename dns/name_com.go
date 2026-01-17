package dns

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

const (
	listRecords  = "https://api.name.com/core/v1/domains/%s/records"
	createRecord = "https://api.name.com/core/v1/domains/%s/records"
	updateRecord = "https://api.name.com/core/v1/domains/%s/records/%d"
)

type NameCom struct {
	DNS     config.DNS
	Domains config.Domains
	TTL     string
}

type NameComRecord struct {
	TTL    int    `json:"ttl"`
	Type   string `json:"type"`
	Answer string `json:"answer"`
	Host   string `json:"host"`
}

type NameComRecordResp struct {
	TTL        int    `json:"ttl"`
	Type       string `json:"type"`
	Answer     string `json:"answer"`
	DomainName string `json:"domainName"`
	Fqdn       string `json:"fqdn"`
	Host       string `json:"host"`
	Id         int    `json:"id"`
	Priority   int    `json:"priority"`
}

type NameComRecordListResp struct {
	TotalCount int                 `json:"totalCount"`
	From       int                 `json:"from"`
	To         int                 `json:"to"`
	Records    []NameComRecordResp `json:"records"`
	LastPage   int                 `json:"lastPage"`
	NextPage   int                 `json:"nextPage"`
}

func (n *NameCom) Init(dnsConf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache) {
	n.Domains.Ipv4Cache = ipv4cache
	n.Domains.Ipv6Cache = ipv6cache
	n.DNS = dnsConf.DNS
	n.Domains.GetNewIp(dnsConf)
	if dnsConf.TTL == "" {
		n.TTL = "300"
	} else {
		n.TTL = dnsConf.TTL
	}
}

func (n *NameCom) AddUpdateDomainRecords() (domains config.Domains) {
	n.addUpdateDomainRecords("A")
	n.addUpdateDomainRecords("AAAA")
	domains = n.Domains
	return
}

func (n *NameCom) addUpdateDomainRecords(recordType string) {
	ipAddr, domains := n.Domains.GetNewIpResult(recordType)
	if ipAddr == "" {
		return
	}

	for _, domain := range domains {
		resp, err := n.getRecordList(domain)
		if err != nil {
			util.Log("查询域名信息发生异常! %s", err)
			domain.UpdateStatus = config.UpdatedFailed
			return
		}
		resp4TypeRecords := make([]NameComRecordResp, 0, resp.TotalCount)
		if resp.TotalCount > 0 {
			for _, r := range resp.Records {
				if r.Type == recordType && r.Host == domain.SubDomain {
					resp4TypeRecords = append(resp4TypeRecords, r)
				}
			}
		}
		if len(resp4TypeRecords) > 0 {
			for _, r := range resp4TypeRecords {
				err := n.update(r, domain, ipAddr, recordType)
				if err != nil {
					domain.UpdateStatus = config.UpdatedFailed
					return
				}
			}
		} else {
			_, err := n.create(domain, recordType, ipAddr)
			if err != nil {
				domain.UpdateStatus = config.UpdatedFailed
				return
			}
		}
	}
}

func (n *NameCom) getRecordList(domain *config.Domain) (resp *NameComRecordListResp, err error) {
	url := fmt.Sprintf(listRecords, domain.DomainName)
	err = n.request("GET", url, nil, &resp)
	return
}

func (n *NameCom) create(domain *config.Domain, recordType string, ipAddr string) (resp *NameComRecord, err error) {
	i, err := strconv.Atoi(n.TTL)
	if err != nil {
		return
	}

	resq := &NameComRecord{
		TTL:    i,
		Answer: ipAddr,
		Host:   domain.SubDomain,
		Type:   recordType,
	}
	url := fmt.Sprintf(createRecord, domain.DomainName)
	err = n.request("POST", url, resq, resp)
	if err != nil {
		util.Log("新增域名解析 %s 失败! 异常信息: %s", domain, err)
		return
	}
	util.Log("新增域名解析 %s 成功! IP: %s", domain, ipAddr)
	return
}

func (n *NameCom) update(record NameComRecordResp, domain *config.Domain, ipAddr, recordType string) (err error) {
	if record.Answer == ipAddr {
		util.Log("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
		return
	}
	record.Answer = ipAddr
	record.Type = recordType
	url := fmt.Sprintf(updateRecord, domain.DomainName, record.Id)
	err = n.request("PUT", url, record, nil)
	if err != nil {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, err)
		return
	}
	util.Log("更新域名解析 %s 成功! IP: %s", domain, ipAddr)
	return
}

func (n *NameCom) request(action string, url string, data any, result any) (err error) {
	jsonStr := make([]byte, 0)
	if data != nil {
		jsonStr, err = json.Marshal(data)
		if err != nil {
			return
		}
	}
	req, err := http.NewRequest(
		action,
		url,
		bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(n.DNS.ID+":"+n.DNS.Secret)))
	if strings.EqualFold(action, "POST") || strings.EqualFold(action, "PUT") {
		req.Header.Add("Content-Type", "application/json")
	}

	client := util.CreateHTTPClient()
	resp, err := client.Do(req)
	err = util.GetHTTPResponse(resp, err, result)

	return
}
