package dns

import (
	"net/url"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

const (
	recordListAPI   string = "https://dnsapi.cn/Record.List"
	recordModifyURL string = "https://dnsapi.cn/Record.Modify"
	recordCreateAPI string = "https://dnsapi.cn/Record.Create"
)

// https://cloud.tencent.com/document/api/302/8516
// Dnspod 腾讯云dns实现
type Dnspod struct {
	DNS     config.DNS
	Domains config.Domains
	TTL     string
}

// DnspodRecord DnspodRecord
type DnspodRecord struct {
	ID      string
	Name    string
	Type    string
	Value   string
	Enabled string
}

// DnspodRecordListResp recordListAPI结果
type DnspodRecordListResp struct {
	DnspodStatus
	Records []DnspodRecord
}

// DnspodStatus DnspodStatus
type DnspodStatus struct {
	Status struct {
		Code    string
		Message string
	}
}

// Init 初始化
func (dnspod *Dnspod) Init(dnsConf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache) {
	dnspod.Domains.Ipv4Cache = ipv4cache
	dnspod.Domains.Ipv6Cache = ipv6cache
	dnspod.DNS = dnsConf.DNS
	dnspod.Domains.GetNewIp(dnsConf)
	if dnsConf.TTL == "" {
		// 默认600s
		dnspod.TTL = "600"
	} else {
		dnspod.TTL = dnsConf.TTL
	}
}

// AddUpdateDomainRecords 添加或更新IPv4/IPv6记录
func (dnspod *Dnspod) AddUpdateDomainRecords() config.Domains {
	dnspod.addUpdateDomainRecords("A")
	dnspod.addUpdateDomainRecords("AAAA")
	return dnspod.Domains
}

func (dnspod *Dnspod) addUpdateDomainRecords(recordType string) {
	ipAddr, domains := dnspod.Domains.GetNewIpResult(recordType)

	if ipAddr == "" {
		return
	}

	for _, domain := range domains {
		result, err := dnspod.getRecordList(domain, recordType)
		if err != nil {
			util.Log("查询域名信息发生异常! %s", err)
			domain.UpdateStatus = config.UpdatedFailed
			return
		}

		if len(result.Records) > 0 {
			// 默认第一个
			recordSelected := result.Records[0]
			params := domain.GetCustomParams()
			if params.Has("record_id") {
				for i := 0; i < len(result.Records); i++ {
					if result.Records[i].ID == params.Get("record_id") {
						recordSelected = result.Records[i]
					}
				}
			}
			// 更新
			dnspod.modify(recordSelected, domain, recordType, ipAddr)
		} else {
			// 新增
			dnspod.create(domain, recordType, ipAddr)
		}
	}
}

// 创建
func (dnspod *Dnspod) create(domain *config.Domain, recordType string, ipAddr string) {
	params := domain.GetCustomParams()
	params.Set("login_token", dnspod.DNS.ID+","+dnspod.DNS.Secret)
	params.Set("domain", domain.DomainName)
	params.Set("sub_domain", domain.GetSubDomain())
	params.Set("record_type", recordType)
	params.Set("value", ipAddr)
	params.Set("ttl", dnspod.TTL)
	params.Set("format", "json")

	if !params.Has("record_line") {
		params.Set("record_line", "默认")
	}

	status, err := dnspod.commonRequest(recordCreateAPI, params, domain)

	if err != nil {
		util.Log("新增域名解析 %s 失败! 异常信息: %s", domain, err)
		domain.UpdateStatus = config.UpdatedFailed
		return
	}

	if err == nil && status.Status.Code == "1" {
		util.Log("新增域名解析 %s 成功! IP: %s", domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	} else {
		util.Log("新增域名解析 %s 失败! 异常信息: %s", domain, status.Status.Message)
		domain.UpdateStatus = config.UpdatedFailed
	}
}

// 修改
func (dnspod *Dnspod) modify(record DnspodRecord, domain *config.Domain, recordType string, ipAddr string) {

	// 相同不修改
	if record.Value == ipAddr {
		util.Log("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
		return
	}

	params := domain.GetCustomParams()
	params.Set("login_token", dnspod.DNS.ID+","+dnspod.DNS.Secret)
	params.Set("domain", domain.DomainName)
	params.Set("sub_domain", domain.GetSubDomain())
	params.Set("record_type", recordType)
	params.Set("value", ipAddr)
	params.Set("ttl", dnspod.TTL)
	params.Set("format", "json")
	params.Set("record_id", record.ID)

	if !params.Has("record_line") {
		params.Set("record_line", "默认")
	}

	status, err := dnspod.commonRequest(recordModifyURL, params, domain)

	if err != nil {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, err)
		domain.UpdateStatus = config.UpdatedFailed
		return
	}

	if status.Status.Code == "1" {
		util.Log("更新域名解析 %s 成功! IP: %s", domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	} else {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, status.Status.Message)
		domain.UpdateStatus = config.UpdatedFailed
	}
}

// 公共
func (dnspod *Dnspod) commonRequest(apiAddr string, values url.Values, domain *config.Domain) (status DnspodStatus, err error) {
	client := util.CreateHTTPClient()
	resp, err := client.PostForm(
		apiAddr,
		values,
	)

	err = util.GetHTTPResponse(resp, err, &status)

	return
}

// 获得域名记录列表
func (dnspod *Dnspod) getRecordList(domain *config.Domain, typ string) (result DnspodRecordListResp, err error) {

	params := domain.GetCustomParams()
	params.Set("login_token", dnspod.DNS.ID+","+dnspod.DNS.Secret)
	params.Set("domain", domain.DomainName)
	params.Set("record_type", typ)
	params.Set("sub_domain", domain.GetSubDomain())
	params.Set("format", "json")

	client := util.CreateHTTPClient()
	resp, err := client.PostForm(
		recordListAPI,
		params,
	)

	err = util.GetHTTPResponse(resp, err, &result)

	return
}
