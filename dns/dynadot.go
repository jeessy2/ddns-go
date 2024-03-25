package dns

import (
	"bytes"
	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// https://www.dynadot.com/set_ddns
const (
	dynadotEndpoint string = "https://www.dynadot.com/set_ddns"
)

// Dynadot Dynadot
type Dynadot struct {
	DNS      config.DNS
	Domains  config.Domains
	TTL      string
	LastIpv4 string
	LastIpv6 string
}

// DynadotRecord record
type DynadotRecord struct {
	DomainName     string
	SubDomainNames []string
	CustomParams   url.Values
	Domains        []*config.Domain
	ContainRoot    bool
}

// DynadotResp 修改/添加返回结果
type DynadotResp struct {
	Status    string   `json:"status"`
	ErrorCode int      `json:"error_code"`
	Content   []string `json:"content"`
}

// Init 初始化
func (dynadot *Dynadot) Init(dnsConf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache) {
	dynadot.Domains.Ipv4Cache = ipv4cache
	dynadot.Domains.Ipv6Cache = ipv6cache
	dynadot.LastIpv4 = ipv4cache.Addr
	dynadot.LastIpv6 = ipv6cache.Addr
	dynadot.DNS = dnsConf.DNS
	dynadot.Domains.GetNewIp(dnsConf)
	if dnsConf.TTL == "" {
		// 默认600s
		dynadot.TTL = "600"
	} else {
		dynadot.TTL = dnsConf.TTL
	}
}

// AddUpdateDomainRecords 添加或更新IPv4/IPv6记录
func (dynadot *Dynadot) AddUpdateDomainRecords() config.Domains {
	dynadot.addOrUpdateDomainRecords("A")
	dynadot.addOrUpdateDomainRecords("AAAA")
	return dynadot.Domains
}

// addOrUpdateDomainRecords 添加或更新记录
func (dynadot *Dynadot) addOrUpdateDomainRecords(recordType string) {
	ipAddr, domains := dynadot.Domains.GetNewIpResult(recordType)

	if len(ipAddr) == 0 {
		return
	}

	// 防止多次发送Webhook通知
	if recordType == "A" {
		if dynadot.LastIpv4 == ipAddr {
			util.Log("你的IPv4未变化, 未触发 %s 请求", "dynadot")
			return
		}
	} else {
		if dynadot.LastIpv6 == ipAddr {
			util.Log("你的IPv6未变化, 未触发 %s 请求", "dynadot")
			return
		}
	}

	records := mergeDomains(domains)
	// dynadot 仅支持一个域名对应一个dynamic password
	if len(records) != 1 {
		util.Log("dynadot仅支持单域名配置，多个域名请添加更多配置")
		return
	}
	for _, record := range records {
		// 创建或更新
		dynadot.createOrModify(record, recordType, ipAddr)
	}
}

// 合并域名的子域名
func mergeDomains(domains []*config.Domain) (records []*DynadotRecord) {
	records = make([]*DynadotRecord, 0)
	for _, domain := range domains {
		var record *DynadotRecord
		for _, r := range records {
			if r.DomainName == domain.DomainName {
				record = r
				params := domain.GetCustomParams()
				for key := range params {
					record.CustomParams.Add(key, params.Get(key))
				}
				record.Domains = append(record.Domains, domain)
				record.SubDomainNames = append(record.SubDomainNames, domain.GetSubDomain())
				break
			}
		}
		if record == nil {
			record = &DynadotRecord{
				DomainName:     domain.DomainName,
				CustomParams:   domain.GetCustomParams(),
				Domains:        []*config.Domain{domain},
				SubDomainNames: []string{domain.GetSubDomain()},
			}
			records = append(records, record)
		}
		if len(domain.SubDomain) == 0 {
			// 包含根域名
			record.ContainRoot = true
		}
	}
	return records
}

// 创建或变更记录
func (dynadot *Dynadot) createOrModify(record *DynadotRecord, recordType string, ipAddr string) {
	params := record.CustomParams
	params.Set("domain", record.DomainName)
	params.Set("subDomain", strings.Join(record.SubDomainNames, ","))
	params.Set("type", recordType)
	params.Set("ip", ipAddr)
	params.Set("pwd", dynadot.DNS.Secret)
	params.Set("ttl", dynadot.TTL)
	params.Set("containRoot", strconv.FormatBool(record.ContainRoot))

	var result DynadotResp
	err := dynadot.request(params, &result)

	domains := record.Domains
	for _, domain := range domains {

		if err != nil {
			util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, err)
			domain.UpdateStatus = config.UpdatedFailed
			return
		}

		if result.ErrorCode != -1 {
			util.Log("更新域名解析 %s 成功! IP: %s", domain, ipAddr)
			domain.UpdateStatus = config.UpdatedSuccess
		} else {
			util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, strings.Join(result.Content, ","))
			domain.UpdateStatus = config.UpdatedFailed
		}
	}

}

// request 统一请求接口
func (dynadot *Dynadot) request(params url.Values, result interface{}) (err error) {

	req, err := http.NewRequest(
		"GET",
		dynadotEndpoint,
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
