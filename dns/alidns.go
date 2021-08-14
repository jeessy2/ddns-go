package dns

import (
	"bytes"
	"ddns-go/config"
	"ddns-go/util"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	alidnsEndpoint string = "https://alidns.aliyuncs.com/"
)

// https://help.aliyun.com/document_detail/29776.html?spm=a2c4g.11186623.6.672.715a45caji9dMA
// Alidns Alidns
type Alidns struct {
	DNSConfig config.DNSConfig
	Domains   config.Domains
	TTL       string
}

// AlidnsSubDomainRecords 记录
type AlidnsSubDomainRecords struct {
	TotalCount    int
	DomainRecords struct {
		Record []struct {
			DomainName string
			RecordID   string
			Value      string
		}
	}
}

// AlidnsResp 修改/添加返回结果
type AlidnsResp struct {
	RecordID  string
	RequestID string
}

// Init 初始化
func (ali *Alidns) Init(conf *config.Config) {
	ali.DNSConfig = conf.DNS
	ali.Domains.ParseDomain(conf)
	if conf.TTL == "" {
		// 默认600s
		ali.TTL = "600"
	} else {
		ali.TTL = conf.TTL
	}
}

// AddUpdateDomainRecords 添加或更新IPv4/IPv6记录
func (ali *Alidns) AddUpdateDomainRecords() config.Domains {
	ali.addUpdateDomainRecords("A")
	ali.addUpdateDomainRecords("AAAA")
	return ali.Domains
}

func (ali *Alidns) addUpdateDomainRecords(recordType string) {
	ipAddr, domains := ali.Domains.ParseDomainResult(recordType)

	if ipAddr == "" {
		return
	}

	for _, domain := range domains {
		var record AlidnsSubDomainRecords
		// 获取当前域名信息
		params := url.Values{}
		params.Set("Action", "DescribeSubDomainRecords")
		params.Set("SubDomain", domain.GetFullDomain())
		params.Set("Type", recordType)
		err := ali.request(params, &record)

		if err != nil {
			return
		}

		if record.TotalCount > 0 {
			// 存在，更新
			ali.modify(record, domain, recordType, ipAddr)
		} else {
			// 不存在，创建
			ali.create(domain, recordType, ipAddr)
		}

	}
}

// 创建
func (ali *Alidns) create(domain *config.Domain, recordType string, ipAddr string) {
	params := url.Values{}
	params.Set("Action", "AddDomainRecord")
	params.Set("DomainName", domain.DomainName)
	params.Set("RR", domain.GetSubDomain())
	params.Set("Type", recordType)
	params.Set("Value", ipAddr)
	params.Set("TTL", ali.TTL)

	var result AlidnsResp
	err := ali.request(params, &result)

	if err == nil && result.RecordID != "" {
		log.Printf("新增域名解析 %s 成功！IP: %s", domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	} else {
		log.Printf("新增域名解析 %s 失败！", domain)
		domain.UpdateStatus = config.UpdatedFailed
	}
}

// 修改
func (ali *Alidns) modify(record AlidnsSubDomainRecords, domain *config.Domain, recordType string, ipAddr string) {

	// 相同不修改
	if len(record.DomainRecords.Record) > 0 && record.DomainRecords.Record[0].Value == ipAddr {
		log.Printf("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
		return
	}

	params := url.Values{}
	params.Set("Action", "UpdateDomainRecord")
	params.Set("RR", domain.GetSubDomain())
	params.Set("RecordId", record.DomainRecords.Record[0].RecordID)
	params.Set("Type", recordType)
	params.Set("Value", ipAddr)
	params.Set("TTL", ali.TTL)

	var result AlidnsResp
	err := ali.request(params, &result)

	if err == nil && result.RecordID != "" {
		log.Printf("更新域名解析 %s 成功！IP: %s", domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	} else {
		log.Printf("更新域名解析 %s 失败！", domain)
		domain.UpdateStatus = config.UpdatedFailed
	}
}

// request 统一请求接口
func (ali *Alidns) request(params url.Values, result interface{}) (err error) {

	util.AliyunSigner(ali.DNSConfig.ID, ali.DNSConfig.Secret, &params)

	req, err := http.NewRequest(
		"GET",
		alidnsEndpoint,
		bytes.NewBuffer(nil),
	)
	req.URL.RawQuery = params.Encode()

	if err != nil {
		log.Println("http.NewRequest失败. Error: ", err)
		return
	}

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	err = util.GetHTTPResponse(resp, alidnsEndpoint, err, result)

	return
}
