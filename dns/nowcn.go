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
	nowcnRecordListAPI   string = "https://todapi.now.cn:2443/api/dns/describe-record-index.json"
	nowcnRecordModifyURL string = "https://todapi.now.cn:2443/api/dns/update-domain-record.json"
	nowcnRecordCreateAPI string = "https://todapi.now.cn:2443/api/dns/add-domain-record.json"
)

// https://www.todaynic.com/partner/mode_Http_Api_detail.php?target_id=d15d8028-7c4f-4a5c-9d15-3a4481c4178e
// Nowcn nowcn DNS实现
type Nowcn struct {
	DNS     config.DNS
	Domains config.Domains
	TTL     string
}

// NowcnRecord DNS记录结构
type NowcnRecord struct {
	ID     int `json:"id"`
	Domain string
	Host   string
	Type   string
	Value  string
	State  int
	// Name    string
	// Enabled string
}

// NowcnRecordListResp 记录列表响应
type NowcnRecordListResp struct {
	NowcnStatus
	Data []NowcnRecord
}

// NowcnStatus API响应状态
type NowcnStatus struct {
	RequestId string `json:"RequestId"`
	Id        int    `json:"Id"`
	Error     string `json:"error"`
}

// Init 初始化
func (nowcn *Nowcn) Init(dnsConf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache) {
	nowcn.Domains.Ipv4Cache = ipv4cache
	nowcn.Domains.Ipv6Cache = ipv6cache
	nowcn.DNS = dnsConf.DNS
	nowcn.Domains.GetNewIp(dnsConf)
	if dnsConf.TTL == "" {
		// 默认600s
		nowcn.TTL = "600"
	} else {
		nowcn.TTL = dnsConf.TTL
	}
}

// AddUpdateDomainRecords 添加或更新IPv4/IPv6记录
func (nowcn *Nowcn) AddUpdateDomainRecords() config.Domains {
	nowcn.addUpdateDomainRecords("A")
	nowcn.addUpdateDomainRecords("AAAA")
	return nowcn.Domains
}

func (nowcn *Nowcn) addUpdateDomainRecords(recordType string) {
	ipAddr, domains := nowcn.Domains.GetNewIpResult(recordType)

	if ipAddr == "" {
		return
	}

	for _, domain := range domains {
		result, err := nowcn.getRecordList(domain, recordType)
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
			nowcn.modify(recordSelected, domain, recordType, ipAddr)
		} else {
			// 新增
			nowcn.create(domain, recordType, ipAddr)
		}
	}
}

// create 创建DNS记录
func (nowcn *Nowcn) create(domain *config.Domain, recordType string, ipAddr string) {
	param := map[string]any{
		"Domain": domain.DomainName,
		"Host":   domain.GetSubDomain(),
		"Type":   recordType,
		"Value":  ipAddr,
		"Ttl":    nowcn.TTL,
	}
	res, err := nowcn.request(nowcnRecordCreateAPI, param)
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
func (nowcn *Nowcn) modify(record NowcnRecord, domain *config.Domain, recordType string, ipAddr string) {
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
		"Ttl":    nowcn.TTL,
	}
	res, err := nowcn.request(nowcnRecordModifyURL, param)
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
func (nowcn *Nowcn) request(apiAddr string, param map[string]any) (status NowcnStatus, err error) {
	param["auth-userid"] = nowcn.DNS.ID
	param["api-key"] = nowcn.DNS.Secret

	fullURL := apiAddr + "?" + nowcn.queryParams(param)
	client := util.CreateHTTPClient()
	resp, err := client.Get(fullURL)

	// 处理响应
	err = util.GetHTTPResponse(resp, err, &status)

	return
}

// getRecordList 获取域名记录列表
func (nowcn *Nowcn) getRecordList(domain *config.Domain, typ string) (result NowcnRecordListResp, err error) {
	param := map[string]any{
		"Domain":      domain.DomainName,
		"auth-userid": nowcn.DNS.ID,
		"api-key":     nowcn.DNS.Secret,
	}
	fullURL := nowcnRecordListAPI + "?" + nowcn.queryParams(param)
	client := util.CreateHTTPClient()
	resp, err := client.Get(fullURL)
	var response NowcnRecordListResp
	result = NowcnRecordListResp{
		Data: make([]NowcnRecord, 0),
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

func (nowcn *Nowcn) queryParams(param map[string]any) string {
	var queryParams []string
	for key, value := range param {
		// 只对键进行URL编码，值保持原样（特别是@符号）
		encodedKey := url.QueryEscape(key)
		valueStr := fmt.Sprintf("%v", value)
		// 对值进行选择性编码，保留@符号
		encodedValue := strings.ReplaceAll(url.QueryEscape(valueStr), "%40", "@")
		queryParams = append(queryParams, encodedKey+"="+encodedValue)
	}
	return strings.Join(queryParams, "&")
}
