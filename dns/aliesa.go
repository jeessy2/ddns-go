package dns

import (
	"bytes"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

const (
	aliesaEndpoint string = "https://esa.cn-hangzhou.aliyuncs.com/"
)

// Aliesa Aliesa
type Aliesa struct {
	DNS     config.DNS
	Domains config.Domains
	TTL     string
}

// AliesaSiteResp 站点返回结果
type AliesaSiteResp struct {
	TotalCount int
	Sites      []AliesaSite
}

// AliesaSites 站点
type AliesaSite struct {
	SiteId     int64
	SiteName   string
	AccessType string
}

// AliesaRecordResp 记录返回结果
type AliesaRecordResp struct {
	TotalCount int
	Records    []AliesaRecord
}

// AliesaRecord 记录
type AliesaRecord struct {
	RecordName string
	RecordId   int64
	Data       struct {
		Value string
	}
}

// AliesaResp 修改/添加返回结果
type AliesaResp struct {
	RecordID  int64
	RequestID string
}

// Init 初始化
func (ali *Aliesa) Init(dnsConf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache) {
	ali.Domains.Ipv4Cache = ipv4cache
	ali.Domains.Ipv6Cache = ipv6cache
	ali.DNS = dnsConf.DNS
	ali.Domains.GetNewIp(dnsConf)
	if dnsConf.TTL == "" {
		// 默认600s
		ali.TTL = "600"
	} else {
		ali.TTL = dnsConf.TTL
	}
}

// AddUpdateDomainRecords 添加或更新IPv4/IPv6记录
func (ali *Aliesa) AddUpdateDomainRecords() config.Domains {
	siteCache := make(map[string]AliesaSite)
	ali.addUpdateDomainRecords("A", siteCache)
	ali.addUpdateDomainRecords("AAAA", siteCache)
	return ali.Domains
}

func (ali *Aliesa) addUpdateDomainRecords(recordType string, siteCache map[string]AliesaSite) {
	ipAddr, domains := ali.Domains.GetNewIpResult(recordType)

	if ipAddr == "" {
		return
	}

	for _, domain := range domains {
		siteSelected, err := ali.getSite(domain, siteCache)

		if err != nil {
			util.Log("查询域名信息发生异常! %s", err)
			domain.UpdateStatus = config.UpdatedFailed
			return
		}

		if siteSelected.SiteId == 0 {
			util.Log("在DNS服务商中未找到根域名: %s", domain.DomainName)
			domain.UpdateStatus = config.UpdatedFailed
			return
		}

		// 获取当前域名信息
		// https://help.aliyun.com/zh/edge-security-acceleration/esa/api-esa-2024-09-10-listrecords
		var recordResp AliesaRecordResp
		params := domain.GetCustomParams()
		params.Set("Action", "ListRecords")
		params.Set("SiteId", strconv.FormatInt(siteSelected.SiteId, 10))
		params.Set("RecordName", domain.String())
		params.Set("Type", "A/AAAA")
		err = ali.request(http.MethodGet, params, &recordResp)

		if err != nil {
			util.Log("查询域名信息发生异常! %s", err)
			domain.UpdateStatus = config.UpdatedFailed
			return
		}

		if recordSelected, ok := getFrom(recordResp, recordType, params.Get("RecordId")); ok {
			// 存在，更新
			ali.modify(recordSelected, domain, ipAddr)
		} else {
			// 不存在，创建
			ali.create(siteSelected, domain, ipAddr)
		}
	}
}

// 创建
// https://help.aliyun.com/zh/edge-security-acceleration/esa/api-esa-2024-09-10-createrecord
func (ali *Aliesa) create(site AliesaSite, domain *config.Domain, ipAddr string) {
	params := domain.GetCustomParams()
	params.Set("Action", "CreateRecord")
	params.Set("SiteId", strconv.FormatInt(site.SiteId, 10))
	params.Set("RecordName", domain.String())

	params.Set("Type", "A/AAAA")
	params.Set("Data", `{"Value":"`+ipAddr+`"}`)
	params.Set("Ttl", ali.TTL)

	// 兼容 CNAME 接入方式
	if site.AccessType == "CNAME" {
		if !params.Has("Proxied") {
			params.Set("Proxied", "true")
		}
		if !params.Has("BizName") {
			params.Set("BizName", "web")
		}
	}

	var result AliesaResp
	err := ali.request(http.MethodPost, params, &result)

	if err != nil {
		util.Log("新增域名解析 %s 失败! 异常信息: %s", domain, err)
		domain.UpdateStatus = config.UpdatedFailed
		return
	}

	if result.RecordID != 0 {
		util.Log("新增域名解析 %s 成功! IP: %s", domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	} else {
		util.Log("新增域名解析 %s 失败! 异常信息: %s", domain, "返回RecordId为空")
		domain.UpdateStatus = config.UpdatedFailed
	}
}

// 修改
// https://help.aliyun.com/zh/edge-security-acceleration/esa/api-esa-2024-09-10-updaterecord
func (ali *Aliesa) modify(recordSelected AliesaRecord, domain *config.Domain, ipAddr string) {
	// 相同不修改
	if recordSelected.Data.Value == ipAddr {
		util.Log("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
		return
	}

	params := domain.GetCustomParams()
	params.Set("Action", "UpdateRecord")
	params.Set("RecordId", strconv.FormatInt(recordSelected.RecordId, 10))

	params.Set("Type", "A/AAAA")
	params.Set("Data", `{"Value":"`+ipAddr+`"}`)
	params.Set("Ttl", ali.TTL)

	var result AliesaResp
	err := ali.request(http.MethodPost, params, &result)

	if err != nil {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, err)
		domain.UpdateStatus = config.UpdatedFailed
		return
	}

	// 不检查 result.RecordID ，更新成功也会返回 0
	util.Log("更新域名解析 %s 成功! IP: %s", domain, ipAddr)
	domain.UpdateStatus = config.UpdatedSuccess
}

// 获取域名的站点信息
// https://help.aliyun.com/zh/edge-security-acceleration/esa/api-esa-2024-09-10-listsites
func (ali *Aliesa) getSite(domain *config.Domain, siteCache map[string]AliesaSite) (result AliesaSite, err error) {
	if site, ok := siteCache[domain.DomainName]; ok {
		return site, nil
	}

	params := domain.GetCustomParams()

	// 解析自定义参数 SiteId，但不使用 GetSite api
	if siteId, _ := strconv.ParseInt(params.Get("SiteId"), 10, 64); siteId != 0 {
		// 兼容 CNAME 接入方式
		result.AccessType = "CNAME"
		result.SiteName = domain.DomainName
		result.SiteId = siteId
		return
	}

	var siteResp AliesaSiteResp
	params.Set("Action", "ListSites")
	params.Set("SiteName", domain.DomainName)
	err = ali.request(http.MethodGet, params, &siteResp)

	if err != nil {
		return
	}

	// siteResp.TotalCount == 0
	if len(siteResp.Sites) == 0 {
		return
	}

	result = siteResp.Sites[0]
	siteCache[domain.DomainName] = result
	return
}

func getFrom(recordResp AliesaRecordResp, recordType string, recordId string) (result AliesaRecord, ok bool) {
	if recordResp.TotalCount == 0 {
		return
	}

	// 指定 RecordId
	if recordId != "" {
		for i := 0; i < len(recordResp.Records); i++ {
			if strconv.FormatInt(recordResp.Records[i].RecordId, 10) == recordId {
				return recordResp.Records[i], true
			}
		}
	}

	// Alidns 的 recordType 不区分 A/AAAA
	if recordType == "AAAA" {
		for i := 0; i < len(recordResp.Records); i++ {
			ip := net.ParseIP(recordResp.Records[i].Data.Value)
			// ipv4.To16() 不为 nil
			if ip.To4() == nil {
				return recordResp.Records[i], true
			}
		}
	}

	if recordType == "A" {
		for i := 0; i < len(recordResp.Records); i++ {
			ip := net.ParseIP(recordResp.Records[i].Data.Value)
			if ip.To4() != nil {
				return recordResp.Records[i], true
			}
		}
	}

	return
}

// request 统一请求接口
func (ali *Aliesa) request(method string, params url.Values, result interface{}) (err error) {
	util.AliyunSigner(ali.DNS.ID, ali.DNS.Secret, &params, method, "2024-09-10")

	req, err := http.NewRequest(
		method,
		aliesaEndpoint,
		bytes.NewBuffer(nil),
	)
	req.URL.RawQuery = params.Encode()

	if err != nil {
		return
	}

	client := util.CreateHTTPClient()
	resp, err := client.Do(req)
	err = util.GetHTTPResponse(resp, err, result)

	return
}
