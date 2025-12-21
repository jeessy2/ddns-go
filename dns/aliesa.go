package dns

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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

	siteCache   map[string]AliesaSite
	domainCache config.DomainTuples
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
	RecordId   int64
	RecordName string
	Data       struct {
		Value string
	}
}

// AliesaResp 修改/添加返回结果
type AliesaResp struct {
	OriginPoolId int64 `json:"Id"`
	RecordID     int64
	RequestID    string
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
	ali.siteCache = make(map[string]AliesaSite)
	ali.domainCache = ali.Domains.GetAllNewIpResult("A/AAAA")
	ali.addUpdateDomainRecords("A")
	ali.addUpdateDomainRecords("AAAA")
	ali.addUpdateDomainRecords("A/AAAA")
	return ali.Domains
}

func (ali *Aliesa) addUpdateDomainRecords(recordType string) {
	for _, domain := range ali.domainCache {
		if domain.RecordType != recordType {
			continue
		}

		// 获取站点
		siteSelected, err := ali.getSite(domain)
		if err != nil {
			util.Log("查询域名信息发生异常! %s", err)
			domain.SetUpdateStatus(config.UpdatedFailed)
			return
		}
		if siteSelected.SiteId == 0 {
			util.Log("在DNS服务商中未找到根域名: %s", domain.Primary.DomainName)
			domain.SetUpdateStatus(config.UpdatedFailed)
			return
		}

		// 处理源地址池
		poolId, origins, err := ali.getOriginPool(siteSelected, domain)
		if err != nil {
			util.Log("查询域名信息发生异常! %s", err)
			domain.SetUpdateStatus(config.UpdatedFailed)
			return
		}
		// TODO：不允许相同ip
		if len(origins) != 0 {
			ali.updateOriginPool(siteSelected, domain, poolId, origins)
			return
		}

		// 获取记录
		recordSelected, err := ali.getRecord(siteSelected, domain, "A/AAAA")
		if err != nil {
			util.Log("查询域名信息发生异常! %s", err)
			domain.SetUpdateStatus(config.UpdatedFailed)
			return
		}
		if recordSelected.RecordId != 0 {
			// 存在，更新
			ali.modify(recordSelected, domain, "A/AAAA")
		} else {
			// 不存在，创建
			ali.create(siteSelected, domain, "A/AAAA")
		}
	}
}

// 创建
// https://help.aliyun.com/zh/edge-security-acceleration/esa/api-esa-2024-09-10-createrecord
func (ali *Aliesa) create(site AliesaSite, domainTuple *config.DomainTuple, recordType string) {
	domain := domainTuple.Primary
	ipAddr := domainTuple.GetIpAddrPool(",")

	params := domain.GetCustomParams()
	params.Set("Action", "CreateRecord")
	params.Set("SiteId", strconv.FormatInt(site.SiteId, 10))
	params.Set("RecordName", domain.String())

	params.Set("Type", recordType)
	params.Set("Data", `{"Value":"`+ipAddr+`"}`)
	params.Set("Ttl", ali.TTL)

	// 兼容 CNAME 接入方式
	if site.AccessType == "CNAME" && !params.Has("Proxied") {
		params.Set("Proxied", "true")
	}
	if params.Has("Proxied") && !params.Has("BizName") {
		params.Set("BizName", "web")
	}

	var result AliesaResp
	err := ali.request(http.MethodPost, params, &result)

	if err != nil {
		util.Log("新增域名解析 %s 失败! 异常信息: %s", domain, err)
		domainTuple.SetUpdateStatus(config.UpdatedFailed)
		return
	}

	if result.RecordID != 0 {
		util.Log("新增域名解析 %s 成功! IP: %s", domain, ipAddr)
		domainTuple.SetUpdateStatus(config.UpdatedSuccess)
	} else {
		util.Log("新增域名解析 %s 失败! 异常信息: %s", domain, "返回RecordId为空")
		domainTuple.SetUpdateStatus(config.UpdatedFailed)
	}
}

// 修改
// https://help.aliyun.com/zh/edge-security-acceleration/esa/api-esa-2024-09-10-updaterecord
func (ali *Aliesa) modify(record AliesaRecord, domainTuple *config.DomainTuple, recordType string) {
	domain := domainTuple.Primary
	ipAddr := domainTuple.GetIpAddrPool(",")
	// 相同不修改
	if record.Data.Value == ipAddr {
		util.Log("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
		return
	}

	params := domain.GetCustomParams()
	params.Set("Action", "UpdateRecord")
	params.Set("RecordId", strconv.FormatInt(record.RecordId, 10))

	params.Set("Type", recordType)
	params.Set("Data", `{"Value":"`+ipAddr+`"}`)
	params.Set("Ttl", ali.TTL)

	var result AliesaResp
	err := ali.request(http.MethodPost, params, &result)

	if err != nil {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, err)
		domainTuple.SetUpdateStatus(config.UpdatedFailed)
		return
	}

	// 不检查 result.RecordID ，更新成功也会返回 0
	util.Log("更新域名解析 %s 成功! IP: %s", domain, ipAddr)
	domainTuple.SetUpdateStatus(config.UpdatedSuccess)
}

// 获取当前域名信息
// https://help.aliyun.com/zh/edge-security-acceleration/esa/api-esa-2024-09-10-listrecords
func (ali *Aliesa) getRecord(site AliesaSite, domainTuple *config.DomainTuple, recordType string) (result AliesaRecord, err error) {
	domain := domainTuple.Primary
	var recordResp AliesaRecordResp

	params := url.Values{}
	params.Set("Action", "ListRecords")
	params.Set("SiteId", strconv.FormatInt(site.SiteId, 10))
	params.Set("RecordName", domain.String())
	params.Set("Type", recordType)
	err = ali.request(http.MethodGet, params, &recordResp)

	// recordResp.TotalCount == 0
	if len(recordResp.Records) == 0 {
		return
	}

	// 指定 RecordId
	recordId := domain.GetCustomParams().Get("RecordId")
	if recordId != "" {
		for i := 0; i < len(recordResp.Records); i++ {
			if strconv.FormatInt(recordResp.Records[i].RecordId, 10) == recordId {
				return recordResp.Records[i], nil
			}
		}
	}
	return recordResp.Records[0], nil
}

// 获取域名的站点信息
// https://help.aliyun.com/zh/edge-security-acceleration/esa/api-esa-2024-09-10-listsites
func (ali *Aliesa) getSite(domainTuple *config.DomainTuple) (result AliesaSite, err error) {
	domain := domainTuple.Primary
	if site, ok := ali.siteCache[domain.DomainName]; ok {
		return site, nil
	}

	// 解析自定义参数 SiteId，但不使用 api GetSite 查询
	siteIdStr := domain.GetCustomParams().Get("SiteId")
	if siteId, _ := strconv.ParseInt(siteIdStr, 10, 64); siteId != 0 {
		// 兼容 CNAME 接入方式
		result.AccessType = "CNAME"
		result.SiteName = domain.DomainName
		result.SiteId = siteId
		return
	}

	var siteResp AliesaSiteResp
	params := url.Values{}
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
	ali.siteCache[domain.DomainName] = result
	return
}

// getOriginPool 获取源地址池
// https://help.aliyun.com/zh/edge-security-acceleration/esa/api-esa-2024-09-10-listoriginpools
func (ali *Aliesa) getOriginPool(site AliesaSite, domainTuple *config.DomainTuple) (id int64, origins []map[string]interface{}, err error) {
	name, found := strings.CutSuffix(domainTuple.Primary.SubDomain, ".origin-pool")
	if !found {
		return
	}

	params := url.Values{}
	params.Set("Action", "ListOriginPools")
	params.Set("SiteId", strconv.FormatInt(site.SiteId, 10))
	params.Set("Name", name)
	params.Set("MatchType", "exact")

	result := struct {
		TotalCount  int
		OriginPools []struct {
			Id      int64
			Origins []map[string]interface{}
		}
	}{}

	err = ali.request(http.MethodGet, params, &result)
	if err == nil && len(result.OriginPools) > 0 {
		pool := result.OriginPools[0]
		id = pool.Id
		origins = pool.Origins
	}
	return
}

// updateOriginPool 更新源地址池
// https://help.aliyun.com/zh/edge-security-acceleration/esa/api-esa-2024-09-10-updateoriginpool
func (ali *Aliesa) updateOriginPool(site AliesaSite, domainTuple *config.DomainTuple, id int64, origins []map[string]interface{}) {
	needUpdate := false
	count := len(domainTuple.Domains)
	for _, origin := range origins {
		// 源地址池不能有多个相同地址，因此 Domain 更少放内层
		for i, d := range domainTuple.Domains {
			name := d.GetCustomParams().Get("Name")
			if origin["Name"] != name {
				continue
			}
			// 相同不修改
			address := domainTuple.IpAddrs[i]
			if origin["Address"] != address {
				origin["Address"] = address
				needUpdate = true
			}
			count--
			break
		}
	}

	domain := domainTuple.Primary
	ipAddr := domainTuple.GetIpAddrPool(",")
	if count > 0 {
		// 有新增的源地址
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, "不支持新增源地址")
		domainTuple.SetUpdateStatus(config.UpdatedFailed)
		return
	}
	if !needUpdate {
		util.Log("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
		return
	}

	originsData, _ := json.Marshal(origins)
	params := url.Values{}
	params.Set("Action", "UpdateOriginPool")
	params.Set("SiteId", strconv.FormatInt(site.SiteId, 10))
	params.Set("Id", strconv.FormatInt(id, 10))
	params.Set("Origins", string(originsData))

	result := AliesaResp{}
	err := ali.request(http.MethodPost, params, &result)

	if err != nil {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, err)
		domainTuple.SetUpdateStatus(config.UpdatedFailed)
		return
	}

	if result.OriginPoolId != 0 {
		util.Log("更新域名解析 %s 成功! IP: %s", domain, ipAddr)
		domainTuple.SetUpdateStatus(config.UpdatedSuccess)
	} else {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, "返回 OriginPool Id为空")
		domainTuple.SetUpdateStatus(config.UpdatedFailed)
	}
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
