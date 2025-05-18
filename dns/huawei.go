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
	Weight  int      `json:"weight"`
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
		customParams := domain.GetCustomParams()
		params := url.Values{}
		params.Set("name", domain.String())
		params.Set("type", recordType)

		// 如果有精准匹配
		// 详见 查询记录集 https://support.huaweicloud.com/api-dns/dns_api_64002.html
		if customParams.Has("zone_id") && customParams.Has("recordset_id") {
			var record HuaweicloudRecordsets
			err := hw.request(
				"GET",
				fmt.Sprintf(huaweicloudEndpoint+"/v2.1/zones/%s/recordsets/%s", customParams.Get("zone_id"), customParams.Get("recordset_id")),
				params,
				&record,
			)

			if err != nil {
				util.Log("查询域名信息发生异常！ %s", err)
				domain.UpdateStatus = config.UpdatedFailed
				return
			}

			// 更新
			hw.modify(record, domain, ipAddr)

		} else { // 没有精准匹配，则支持更多的查询参数。详见 查询租户记录集列表 https://support.huaweicloud.com/api-dns/dns_api_64003.html
			// 复制所有自定义参数
			util.CopyUrlParams(customParams, params, nil)
			// 参数名修正
			if params.Has("recordset_id") {
				params.Set("id", params.Get("recordset_id"))
				params.Del("recordset_id")
			}

			var records HuaweicloudRecordsResp
			err := hw.request(
				"GET",
				huaweicloudEndpoint+"/v2.1/recordsets",
				params,
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
					hw.modify(record, domain, ipAddr)
					find = true
					break
				}
			}

			if !find {
				thIdParamName := ""
				if customParams.Has("id") {
					thIdParamName = "id"
				} else if customParams.Has("recordset_id") {
					thIdParamName = "recordset_id"
				}

				if thIdParamName != "" {
					util.Log("域名 %s 解析未找到，且因添加了参数 %s=%s 导致无法创建。本次更新已被忽略", domain, thIdParamName, customParams.Get(thIdParamName))
				} else {
					// 新增
					hw.create(domain, recordType, ipAddr)
				}
			}
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
		util.Log("在DNS服务商中未找到根域名: %s", domain.DomainName)
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
		Weight:  1,
	}
	var result HuaweicloudRecordsets
	err = hw.request(
		"POST",
		fmt.Sprintf(huaweicloudEndpoint+"/v2.1/zones/%s/recordsets", zoneID),
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
func (hw *Huaweicloud) modify(record HuaweicloudRecordsets, domain *config.Domain, ipAddr string) {
	originDnsRecordCount := 0
	if record.Type == "A" {
		// 只多不少，记录最大值。防止误删
		if len(record.Records) > hw.Domains.Ipv4DnsRecordCount {
			hw.Domains.Ipv4DnsRecordCount = len(record.Records)
		}
		originDnsRecordCount = hw.Domains.Ipv4DnsRecordCount
	} else if record.Type == "AAAA" && hw.Domains.Ipv6DnsRecordCount == 0 {
		if len(record.Records) > hw.Domains.Ipv6DnsRecordCount {
			hw.Domains.Ipv6DnsRecordCount = len(record.Records)
		}
		originDnsRecordCount = hw.Domains.Ipv6DnsRecordCount
	}

	// 相同不修改{
	if len(record.Records) > 0 {
		found := false
		for index, item := range record.Records {
			if item == ipAddr {
				found = true
				if index == 0 {
					util.Log("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
					return
				} else {
					// 元素不在首位则挪动并覆盖首位（元素个数减一）
					// 主要为了保留用户预设的固定的内网IP地址（计算机离线可用的那种，本地 Wireguard 网络之类的）
					// 同时可以覆盖以下场景
					// 1. 公网IP地址时有时无
					// 2. DHCP分配的IP地址时有时无
					if len(record.Records) == originDnsRecordCount {
						record.Records = append([]string{ipAddr}, append(record.Records[1:index], record.Records[index+1:]...)...)
					}
					break
				}
			}
		}

		if !found {
			if len(record.Records)+1 == originDnsRecordCount {
				// 元素个数不够，系以往发现重复而被删了的。直接插入到首位，给它补回去
				record.Records = append([]string{ipAddr}, record.Records[:]...)
			} else { // 其它情况，不管3721 就只修改首位元素
				record.Records[0] = ipAddr
			}
		}
	} else {
		record.Records = []string{ipAddr}
	}

	var request = make(map[string]interface{})
	request["name"] = record.Name
	request["type"] = record.Type
	request["records"] = record.Records
	request["ttl"] = hw.TTL

	var result HuaweicloudRecordsets

	err := hw.request(
		"PUT",
		fmt.Sprintf(huaweicloudEndpoint+"/v2.1/zones/%s/recordsets/%s", record.ZoneID, record.ID),
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
		huaweicloudEndpoint+"/v2/zones",
		url.Values{"name": []string{domain.DomainName}},
		&result,
	)

	return
}

// request 统一请求接口
func (hw *Huaweicloud) request(method string, urlString string, data interface{}, result interface{}) (err error) {
	var (
		req *http.Request
	)

	if method == "GET" {
		req, err = http.NewRequest(
			method,
			urlString,
			bytes.NewBuffer(nil),
		)

		req.URL.RawQuery = data.(url.Values).Encode()
	} else {
		jsonStr := make([]byte, 0)
		if data != nil {
			jsonStr, _ = json.Marshal(data)
		}

		req, err = http.NewRequest(
			method,
			urlString,
			bytes.NewBuffer(jsonStr),
		)
	}

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
