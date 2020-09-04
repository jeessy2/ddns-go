package dns

import (
	"ddns-go/config"
	"ddns-go/util"
	"log"
	"net/http"
	"net/url"
)

const (
	recordListAPI   string = "https://dnsapi.cn/Record.List"
	recordModifyURL string = "https://dnsapi.cn/Record.Modify"
	recordCreateAPI string = "https://dnsapi.cn/Record.Create"
)

// Dnspod 腾讯云dns实现
type Dnspod struct {
	DNSConfig config.DNSConfig
	Domains
}

// DnspodRecordListResp recordListAPI结果
type DnspodRecordListResp struct {
	DnspodStatus
	Records []struct {
		ID      string
		Name    string
		Type    string
		Value   string
		Enabled string
	}
}

// DnspodStatus DnspodStatus
type DnspodStatus struct {
	Status struct {
		Code    string
		Message string
	}
}

// Init 初始化
func (dnspod *Dnspod) Init(conf *config.Config) {
	dnspod.DNSConfig = conf.DNS
	dnspod.Domains.ParseDomain(conf)
}

// AddUpdateIpv4DomainRecords 添加或更新IPV4记录
func (dnspod *Dnspod) AddUpdateIpv4DomainRecords() {
	dnspod.addUpdateDomainRecords("A")
}

// AddUpdateIpv6DomainRecords 添加或更新IPV6记录
func (dnspod *Dnspod) AddUpdateIpv6DomainRecords() {
	dnspod.addUpdateDomainRecords("AAAA")
}

func (dnspod *Dnspod) addUpdateDomainRecords(recordType string) {
	ipAddr := dnspod.Ipv4Addr
	domains := dnspod.Ipv4Domains
	if recordType == "AAAA" {
		ipAddr = dnspod.Ipv6Addr
		domains = dnspod.Ipv6Domains
	}

	if ipAddr == "" {
		return
	}

	for _, domain := range domains {
		result, err := dnspod.getRecordList(domain, recordType)
		if err != nil {
			return
		}

		if len(result.Records) > 0 {
			// 更新
			dnspod.modify(result, domain, recordType, ipAddr)
		} else {
			// 新增
			dnspod.create(result, domain, recordType, ipAddr)
		}
	}
}

// 创建
func (dnspod *Dnspod) create(result DnspodRecordListResp, domain *Domain, recordType string, ipAddr string) {
	status, err := dnspod.commonRequest(
		recordCreateAPI,
		url.Values{
			"login_token": {dnspod.DNSConfig.ID + "," + dnspod.DNSConfig.Secret},
			"domain":      {domain.DomainName},
			"sub_domain":  {domain.GetSubDomain()},
			"record_type": {recordType},
			"record_line": {"默认"},
			"value":       {ipAddr},
			"format":      {"json"},
		},
		domain,
	)
	if err == nil && status.Status.Code == "1" {
		log.Printf("新增域名解析 %s 成功！IP: %s", domain, ipAddr)
	} else {
		log.Printf("新增域名解析 %s 失败！Code: %s, Message: %s", domain, status.Status.Code, status.Status.Message)
	}
}

// 修改
func (dnspod *Dnspod) modify(result DnspodRecordListResp, domain *Domain, recordType string, ipAddr string) {
	for _, record := range result.Records {
		// 相同不修改
		if record.Value == ipAddr {
			log.Printf("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
			continue
		}
		status, err := dnspod.commonRequest(
			recordModifyURL,
			url.Values{
				"login_token": {dnspod.DNSConfig.ID + "," + dnspod.DNSConfig.Secret},
				"domain":      {domain.DomainName},
				"sub_domain":  {domain.GetSubDomain()},
				"record_type": {recordType},
				"record_line": {"默认"},
				"record_id":   {record.ID},
				"value":       {ipAddr},
				"format":      {"json"},
			},
			domain,
		)
		if err == nil && status.Status.Code == "1" {
			log.Printf("更新域名解析 %s 成功！IP: %s", domain, ipAddr)
		} else {
			log.Printf("更新域名解析 %s 失败！Code: %s, Message: %s", domain, status.Status.Code, status.Status.Message)
		}
	}
}

// 公共
func (dnspod *Dnspod) commonRequest(apiAddr string, values url.Values, domain *Domain) (status DnspodStatus, err error) {
	resp, err := http.PostForm(
		apiAddr,
		values,
	)

	err = util.GetHTTPResponse(resp, apiAddr, err, &status)

	return
}

// 获得域名记录列表
func (dnspod *Dnspod) getRecordList(domain *Domain, typ string) (result DnspodRecordListResp, err error) {
	values := url.Values{
		"login_token": {dnspod.DNSConfig.ID + "," + dnspod.DNSConfig.Secret},
		"domain":      {domain.DomainName},
		"record_type": {typ},
		"sub_domain":  {domain.GetSubDomain()},
		"format":      {"json"},
	}

	resp, err := http.PostForm(
		recordListAPI,
		values,
	)

	err = util.GetHTTPResponse(resp, recordListAPI, err, &result)

	return
}
