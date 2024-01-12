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

const (
	huaweicloudEndpoint string = "https://dns.myhuaweicloud.com"
)

// https://support.huaweicloud.com/api-dns/dns_api_64001.html
// Huaweicloud Huaweicloud
type Huaweicloud struct {
	DNS     config.DNS
	Domains config.Domains
	TTL     int
}

// HuaweicloudZonesResp zones response
type HuaweicloudZonesResp struct {
	Zones []struct {
		ID         string
		Name       string
		Recordsets []HuaweicloudRecordsets
	}
}

// HuaweicloudRecordsResp 记录返回结果
type HuaweicloudRecordsResp struct {
	Recordsets []HuaweicloudRecordsets
}

// HuaweicloudRecordsets 记录
type HuaweicloudRecordsets struct {
	ID      string
	Name    string `json:"name"`
	ZoneID  string `json:"zone_id"`
	Status  string
	Type    string   `json:"type"`
	TTL     int      `json:"ttl"`
	Records []string `json:"records"`
}

// Init 初始化
func (hw *Huaweicloud) Init(dnsConf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache) {
	hw.Domains.Ipv4Cache = ipv4cache
	hw.Domains.Ipv6Cache = ipv6cache
	hw.DNS = dnsConf.DNS
	hw.Domains.GetNewIp(dnsConf)
	if dnsConf.TTL == "" {
		// 默认300s
		hw.TTL = 300
	} else {
		ttl, err := strconv.Atoi(dnsConf.TTL)
		if err != nil {
			hw.TTL = 300
		} else {
			hw.TTL = ttl
		}
	}
}

// AddUpdateDomainRecords 添加或更新IPv4/IPv6记录
func (hw *Huaweicloud) AddUpdateDomainRecords() config.Domains {
	hw.addUpdateDomainRecords("A")
	hw.addUpdateDomainRecords("AAAA")
	return hw.Domains
}

func (hw *Huaweicloud) addUpdateDomainRecords(recordType string) {
	ipAddr, domains := hw.Domains.GetNewIpResult(recordType)

	if ipAddr == "" {
		return
	}

	for _, domain := range domains {

		var records HuaweicloudRecordsResp

		err := hw.request(
			"GET",
			fmt.Sprintf(huaweicloudEndpoint+"/v2/recordsets?type=%s&name=%s", recordType, domain),
			nil,
			&records,
		)

		if err != nil {
			util.Log("查询域名信息发生异常! %s", err)
			domain.UpdateStatus = config.UpdatedFailed
			return
		}

		find := false
		for _, record := range records.Recordsets {
			// 名称相同才更新。华为云默认是模糊搜索
			if record.Name == domain.String()+"." {
				// 更新
				hw.modify(record, domain, recordType, ipAddr)
				find = true
				break
			}
		}

		if !find {
			// 新增
			hw.create(domain, recordType, ipAddr)
		}

	}
}

// 创建
func (hw *Huaweicloud) create(domain *config.Domain, recordType string, ipAddr string) {
	zone, err := hw.getZones(domain)
	if err != nil {
		util.Log("查询域名信息发生异常! %s", err)
		domain.UpdateStatus = config.UpdatedFailed
		return
	}

	if len(zone.Zones) == 0 {
		util.Log("在DNS服务商中未找到域名: %s", domain.String())
		domain.UpdateStatus = config.UpdatedFailed
		return
	}

	zoneID := zone.Zones[0].ID
	for _, z := range zone.Zones {
		if z.Name == domain.DomainName+"." {
			zoneID = z.ID
			break
		}
	}

	record := &HuaweicloudRecordsets{
		Type:    recordType,
		Name:    domain.String() + ".",
		Records: []string{ipAddr},
		TTL:     hw.TTL,
	}
	var result HuaweicloudRecordsets
	err = hw.request(
		"POST",
		fmt.Sprintf(huaweicloudEndpoint+"/v2/zones/%s/recordsets", zoneID),
		record,
		&result,
	)

	if err != nil {
		util.Log("新增域名解析 %s 失败! 异常信息: %s", domain, err)
		domain.UpdateStatus = config.UpdatedFailed
		return
	}

	if len(result.Records) > 0 && result.Records[0] == ipAddr {
		util.Log("新增域名解析 %s 成功! IP: %s", domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	} else {
		util.Log("新增域名解析 %s 失败! 异常信息: %s", domain, result.Status)
		domain.UpdateStatus = config.UpdatedFailed
	}
}

// 修改
func (hw *Huaweicloud) modify(record HuaweicloudRecordsets, domain *config.Domain, recordType string, ipAddr string) {

	// 相同不修改
	if len(record.Records) > 0 && record.Records[0] == ipAddr {
		util.Log("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
		return
	}

	var request map[string]interface{} = make(map[string]interface{})
	request["records"] = []string{ipAddr}
	request["ttl"] = hw.TTL

	var result HuaweicloudRecordsets

	err := hw.request(
		"PUT",
		fmt.Sprintf(huaweicloudEndpoint+"/v2/zones/%s/recordsets/%s", record.ZoneID, record.ID),
		&request,
		&result,
	)

	if err != nil {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, err)
		domain.UpdateStatus = config.UpdatedFailed
		return
	}

	if len(result.Records) > 0 && result.Records[0] == ipAddr {
		util.Log("更新域名解析 %s 成功! IP: %s", domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	} else {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, result.Status)
		domain.UpdateStatus = config.UpdatedFailed
	}
}

// 获得域名记录列表
func (hw *Huaweicloud) getZones(domain *config.Domain) (result HuaweicloudZonesResp, err error) {
	err = hw.request(
		"GET",
		fmt.Sprintf(huaweicloudEndpoint+"/v2/zones?name=%s", domain.DomainName),
		nil,
		&result,
	)

	return
}

// request 统一请求接口
func (hw *Huaweicloud) request(method string, url string, data interface{}, result interface{}) (err error) {
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

	s := util.Signer{
		Key:    hw.DNS.ID,
		Secret: hw.DNS.Secret,
	}
	s.Sign(req)

	req.Header.Add("content-type", "application/json")

	client := util.CreateHTTPClient()
	resp, err := client.Do(req)
	err = util.GetHTTPResponse(resp, err, result)

	return
}
