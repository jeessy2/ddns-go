package dns

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

const (
	eranetRecordListAPI   string = "http://api.eranet.com:2080/api/dns/describe-record-index.json"
	eranetRecordModifyURL string = "http://api.eranet.com:2080/api/dns/update-domain-record.json"
	eranetRecordCreateAPI string = "http://api.eranet.com:2080/api/dns/add-domain-record.json"
)

// https://partner.tnet.hk/adminCN/mode_Http_Api_detail.php
// Eranet DNS实现
type Eranet struct {
	DNS     config.DNS
	Domains config.Domains
	TTL     string
}

type EranetRecord struct {
	ID     int `json:"id"`
	Domain string
	Host   string
	Type   string
	Value  string
	State  int
	// Name    string
	// Enabled string
}

type EranetRecordListResp struct {
	EranetStatus
	Data []EranetRecord
}

type EranetStatus struct {
	RequestId string `json:"RequestId"`
	Id        int    `json:"Id"`
	Error     string `json:"error"`
}

// Init 初始化
func (eranet *Eranet) Init(dnsConf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache) {
	eranet.Domains.Ipv4Cache = ipv4cache
	eranet.Domains.Ipv6Cache = ipv6cache
	eranet.DNS = dnsConf.DNS
	eranet.Domains.GetNewIp(dnsConf)
	if dnsConf.TTL == "" {
		// 默认600s
		eranet.TTL = "600"
	} else {
		eranet.TTL = dnsConf.TTL
	}
}

// AddUpdateDomainRecords 添加或更新IPv4/IPv6记录
func (eranet *Eranet) AddUpdateDomainRecords() config.Domains {
	eranet.addUpdateDomainRecords("A")
	eranet.addUpdateDomainRecords("AAAA")
	return eranet.Domains
}

func (eranet *Eranet) addUpdateDomainRecords(recordType string) {
	ipAddr, domains := eranet.Domains.GetNewIpResult(recordType)

	if ipAddr == "" {
		return
	}

	for _, domain := range domains {
		result, err := eranet.getRecordList(domain, recordType)
		if err != nil {
			util.Log("查询域名信息发生异常! %s", err)
			domain.UpdateStatus = config.UpdatedFailed
			return
		}

		if len(result.Data) > 0 {
			// 默认第一个
			recordSelected := result.Data[0]
			params := domain.GetCustomParams()
			if params.Has("Id") {
				for i := 0; i < len(result.Data); i++ {
					if strconv.Itoa(result.Data[i].ID) == params.Get("Id") {
						recordSelected = result.Data[i]
					}
				}
			}
			// 更新
			eranet.modify(recordSelected, domain, recordType, ipAddr)
		} else {
			// 新增
			eranet.create(domain, recordType, ipAddr)
		}
	}
}

// create 创建DNS记录
func (eranet *Eranet) create(domain *config.Domain, recordType string, ipAddr string) {
	param := map[string]any{
		"Domain": domain.DomainName,
		"Host":   domain.GetSubDomain(),
		"Type":   recordType,
		"Value":  ipAddr,
		"Ttl":    eranet.TTL,
	}
	res, err := eranet.request(eranetRecordCreateAPI, param)
	if err != nil {
		util.Log("新增域名解析 %s 失败! 异常信息: %s", domain, err.Error())
		domain.UpdateStatus = config.UpdatedFailed
	} else if res.Error != "" {
		util.Log("新增域名解析 %s 失败! 异常信息: %s", domain, res.Error)
		domain.UpdateStatus = config.UpdatedFailed
	} else {
		util.Log("新增域名解析 %s 成功! IP: %s", domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	}
}

// modify 修改DNS记录
func (eranet *Eranet) modify(record EranetRecord, domain *config.Domain, recordType string, ipAddr string) {
	// 相同不修改
	if record.Value == ipAddr {
		util.Log("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
		return
	}
	param := map[string]any{
		"Id":     record.ID,
		"Domain": domain.DomainName,
		"Host":   domain.GetSubDomain(),
		"Type":   recordType,
		"Value":  ipAddr,
		"Ttl":    eranet.TTL,
	}
	res, err := eranet.request(eranetRecordModifyURL, param)
	if err != nil {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, err.Error())
		domain.UpdateStatus = config.UpdatedFailed
	} else if res.Error != "" {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, res.Error)
		domain.UpdateStatus = config.UpdatedFailed
	} else {
		util.Log("更新域名解析 %s 成功! IP: %s", domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	}
}

// request 发送HTTP请求
func (eranet *Eranet) request(apiAddr string, param map[string]any) (status EranetStatus, err error) {
	param["auth-userid"] = eranet.DNS.ID
	param["api-key"] = eranet.DNS.Secret

	fullURL := apiAddr + "?" + eranet.queryParams(param)
	client := util.CreateHTTPClient()
	resp, err := client.Get(fullURL)

	// 处理响应
	err = util.GetHTTPResponse(resp, err, &status)

	return
}

// getRecordList 获取域名记录列表
func (eranet *Eranet) getRecordList(domain *config.Domain, typ string) (result EranetRecordListResp, err error) {
	param := map[string]any{
		"Domain":      domain.DomainName,
		"auth-userid": eranet.DNS.ID,
		"api-key":     eranet.DNS.Secret,
	}
	fullURL := eranetRecordListAPI + "?" + eranet.queryParams(param)
	client := util.CreateHTTPClient()
	resp, err := client.Get(fullURL)
	var response EranetRecordListResp
	result = EranetRecordListResp{
		Data: make([]EranetRecord, 0),
	}
	err = util.GetHTTPResponse(resp, err, &response)
	for _, v := range response.Data {
		if v.Host == domain.GetSubDomain() {
			result.Data = append(result.Data, v)
			break
		}

	}
	return
}

func (eranet *Eranet) queryParams(param map[string]any) string {
	var queryParams []string
	for key, value := range param {
		// 只对键进行URL编码，值保持原样（特别是@符号）
		encodedKey := url.QueryEscape(key)
		valueStr := fmt.Sprintf("%v", value)
		// 对值进行选择性编码，保留@符号
		encodedValue := strings.ReplaceAll(url.QueryEscape(valueStr), "%40", "@")
		encodedValue = strings.ReplaceAll(encodedValue, "%3A", ":")
		queryParams = append(queryParams, encodedKey+"="+encodedValue)
	}
	return strings.Join(queryParams, "&")
}
