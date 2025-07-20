package dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

const gcoreAPIEndpoint = "https://api.gcore.com/dns/v2"

// Gcore Gcore DNS实现
type Gcore struct {
	DNS     config.DNS
	Domains config.Domains
	TTL     int
}

// GcoreZoneResponse zones返回结果
type GcoreZoneResponse struct {
	Zones       []GcoreZone `json:"zones"`
	TotalAmount int         `json:"total_amount"`
}

// GcoreZone 域名信息
type GcoreZone struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GcoreRRSetListResponse RRSet列表返回结果
type GcoreRRSetListResponse struct {
	RRSets      []GcoreRRSet `json:"rrsets"`
	TotalAmount int          `json:"total_amount"`
}

// GcoreRRSet RRSet记录实体
type GcoreRRSet struct {
	Name            string                 `json:"name"`
	Type            string                 `json:"type"`
	TTL             int                    `json:"ttl"`
	ResourceRecords []GcoreResourceRecord  `json:"resource_records"`
	Meta            map[string]interface{} `json:"meta,omitempty"`
}

// GcoreResourceRecord 资源记录
type GcoreResourceRecord struct {
	Content []interface{}          `json:"content"`
	Enabled bool                   `json:"enabled"`
	ID      int                    `json:"id,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

// GcoreInputRRSet 输入的RRSet
type GcoreInputRRSet struct {
	TTL             int                        `json:"ttl"`
	ResourceRecords []GcoreInputResourceRecord `json:"resource_records"`
	Meta            map[string]interface{}     `json:"meta,omitempty"`
}

// GcoreInputResourceRecord 输入的资源记录
type GcoreInputResourceRecord struct {
	Content []interface{}          `json:"content"`
	Enabled bool                   `json:"enabled"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

// Init 初始化
func (gc *Gcore) Init(dnsConf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache) {
	gc.Domains.Ipv4Cache = ipv4cache
	gc.Domains.Ipv6Cache = ipv6cache
	gc.DNS = dnsConf.DNS
	gc.Domains.GetNewIp(dnsConf)
	if dnsConf.TTL == "" {
		// 默认 120 秒（免费版最低值）
		gc.TTL = 120
	} else {
		ttl, err := strconv.Atoi(dnsConf.TTL)
		if err != nil {
			gc.TTL = 120
		} else {
			gc.TTL = ttl
		}
	}
}

// AddUpdateDomainRecords 添加或更新 IPv4 / IPv6 记录
func (gc *Gcore) AddUpdateDomainRecords() config.Domains {
	gc.addUpdateDomainRecords("A")
	gc.addUpdateDomainRecords("AAAA")
	return gc.Domains
}

func (gc *Gcore) addUpdateDomainRecords(recordType string) {
	ipAddr, domains := gc.Domains.GetNewIpResult(recordType)

	if ipAddr == "" {
		return
	}

	for _, domain := range domains {
		// get zone
		zoneInfo, err := gc.getZoneByDomain(domain)
		if err != nil {
			util.Log("查询域名信息发生异常! %s", err)
			domain.UpdateStatus = config.UpdatedFailed
			continue
		}

		if zoneInfo == nil {
			util.Log("在DNS服务商中未找到根域名: %s", domain.DomainName)
			domain.UpdateStatus = config.UpdatedFailed
			continue
		}

		// 查询现有记录
		existingRecord, err := gc.getRRSet(zoneInfo.Name, domain.GetSubDomain(), recordType)
		if err != nil {
			util.Log("查询域名信息发生异常! %s", err)
			domain.UpdateStatus = config.UpdatedFailed
			continue
		}

		if existingRecord != nil {
			// 更新现有记录
			gc.updateRecord(zoneInfo.Name, domain, recordType, ipAddr, existingRecord)
		} else {
			// 创建新记录
			gc.createRecord(zoneInfo.Name, domain, recordType, ipAddr)
		}
	}
}

// 获取域名对应的Zone信息
func (gc *Gcore) getZoneByDomain(domain *config.Domain) (*GcoreZone, error) {
	var result GcoreZoneResponse
	params := url.Values{}
	params.Set("name", domain.DomainName)

	err := gc.request(
		"GET",
		fmt.Sprintf("%s/zones?%s", gcoreAPIEndpoint, params.Encode()),
		nil,
		&result,
	)

	if err != nil {
		return nil, err
	}

	if len(result.Zones) > 0 {
		return &result.Zones[0], nil
	}

	return nil, nil
}

// 获取指定的RRSet记录
func (gc *Gcore) getRRSet(zoneName, recordName, recordType string) (*GcoreRRSet, error) {
	var result GcoreRRSetListResponse

	err := gc.request(
		"GET",
		fmt.Sprintf("%s/zones/%s/rrsets", gcoreAPIEndpoint, zoneName),
		nil,
		&result,
	)

	if err != nil {
		return nil, err
	}

	// 查找匹配的记录
	fullRecordName := recordName
	if recordName != "" && recordName != "@" {
		fullRecordName = recordName + "." + zoneName
	} else {
		fullRecordName = zoneName
	}

	for _, rrset := range result.RRSets {
		if rrset.Name == fullRecordName && rrset.Type == recordType {
			return &rrset, nil
		}
	}

	return nil, nil
}

// 创建新记录
func (gc *Gcore) createRecord(zoneName string, domain *config.Domain, recordType string, ipAddr string) {
	recordName := domain.GetSubDomain()
	if recordName == "" || recordName == "@" {
		recordName = zoneName
	} else {
		recordName = recordName + "." + zoneName
	}

	inputRRSet := GcoreInputRRSet{
		TTL: gc.TTL,
		ResourceRecords: []GcoreInputResourceRecord{
			{
				Content: []interface{}{ipAddr},
				Enabled: true,
			},
		},
	}

	var result interface{}
	err := gc.request(
		"POST",
		fmt.Sprintf("%s/zones/%s/%s/%s", gcoreAPIEndpoint, zoneName, recordName, recordType),
		inputRRSet,
		&result,
	)

	if err != nil {
		util.Log("新增域名解析 %s 失败! 异常信息: %s", domain, err)
		domain.UpdateStatus = config.UpdatedFailed
		return
	}

	util.Log("新增域名解析 %s 成功! IP: %s", domain, ipAddr)
	domain.UpdateStatus = config.UpdatedSuccess
}

// 更新现有记录
func (gc *Gcore) updateRecord(zoneName string, domain *config.Domain, recordType string, ipAddr string, existingRecord *GcoreRRSet) {
	// 检查IP是否相同
	if len(existingRecord.ResourceRecords) > 0 && len(existingRecord.ResourceRecords[0].Content) > 0 {
		if existingRecord.ResourceRecords[0].Content[0] == ipAddr {
			util.Log("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
			return
		}
	}

	recordName := domain.GetSubDomain()
	if recordName == "" || recordName == "@" {
		recordName = zoneName
	} else {
		recordName = recordName + "." + zoneName
	}

	inputRRSet := GcoreInputRRSet{
		TTL: gc.TTL,
		ResourceRecords: []GcoreInputResourceRecord{
			{
				Content: []interface{}{ipAddr},
				Enabled: true,
			},
		},
	}

	var result interface{}
	err := gc.request(
		"PUT",
		fmt.Sprintf("%s/zones/%s/%s/%s", gcoreAPIEndpoint, zoneName, recordName, recordType),
		inputRRSet,
		&result,
	)

	if err != nil {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, err)
		domain.UpdateStatus = config.UpdatedFailed
		return
	}

	util.Log("更新域名解析 %s 成功! IP: %s", domain, ipAddr)
	domain.UpdateStatus = config.UpdatedSuccess
}

// request 统一请求接口
func (gc *Gcore) request(method string, url string, data interface{}, result interface{}) (err error) {
	jsonStr := make([]byte, 0)
	if data != nil {
		jsonStr, _ = json.Marshal(data)
	}

	req, err := http.NewRequest(
		method,
		url,
		bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", "APIKey "+gc.DNS.Secret)
	req.Header.Set("Content-Type", "application/json")

	client := util.CreateHTTPClient()
	resp, err := client.Do(req)
	err = util.GetHTTPResponse(resp, err, result)

	return
}
