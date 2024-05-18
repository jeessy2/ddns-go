package dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

const (
	porkbunEndpoint string = "https://porkbun.com/api/json/v3/dns"
)

type Porkbun struct {
	DNSConfig config.DNS
	Domains   config.Domains
	TTL       string
}
type PorkbunDomainRecord struct {
	Name    *string `json:"name"`    // subdomain
	Type    *string `json:"type"`    // record type, e.g. A AAAA CNAME
	Content *string `json:"content"` // value
	Ttl     *string `json:"ttl"`     // default 300
}

type PorkbunResponse struct {
	Status string `json:"status"`
}

type PorkbunDomainQueryResponse struct {
	*PorkbunResponse
	Records []PorkbunDomainRecord `json:"records"`
}

type PorkbunApiKey struct {
	AccessKey string `json:"apikey"`
	SecretKey string `json:"secretapikey"`
}

type PorkbunDomainCreateOrUpdateVO struct {
	*PorkbunApiKey
	*PorkbunDomainRecord
}

// Init 初始化
func (pb *Porkbun) Init(conf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache) {
	pb.Domains.Ipv4Cache = ipv4cache
	pb.Domains.Ipv6Cache = ipv6cache
	pb.DNSConfig = conf.DNS
	pb.Domains.GetNewIp(conf)
	if conf.TTL == "" {
		// 默认600s
		pb.TTL = "600"
	} else {
		pb.TTL = conf.TTL
	}
}

// AddUpdateDomainRecords 添加或更新IPv4/IPv6记录
func (pb *Porkbun) AddUpdateDomainRecords() config.Domains {
	pb.addUpdateDomainRecords("A")
	pb.addUpdateDomainRecords("AAAA")
	return pb.Domains
}

func (pb *Porkbun) addUpdateDomainRecords(recordType string) {
	ipAddr, domains := pb.Domains.GetNewIpResult(recordType)

	if ipAddr == "" {
		return
	}

	for _, domain := range domains {
		var record PorkbunDomainQueryResponse
		// 获取当前域名信息
		err := pb.request(
			porkbunEndpoint+fmt.Sprintf("/retrieveByNameType/%s/%s/%s", domain.DomainName, recordType, domain.SubDomain),
			&PorkbunApiKey{
				AccessKey: pb.DNSConfig.ID,
				SecretKey: pb.DNSConfig.Secret,
			},
			&record,
		)

		if err != nil {
			util.Log("查询域名信息发生异常! %s", err)
			domain.UpdateStatus = config.UpdatedFailed
			return
		}
		if record.Status == "SUCCESS" {
			if len(record.Records) > 0 {
				// 存在，更新
				pb.modify(&record, domain, recordType, ipAddr)
			} else {
				// 不存在，创建
				pb.create(domain, recordType, ipAddr)
			}
		} else {
			util.Log("在DNS服务商中未找到根域名: %s", domain.DomainName)
			domain.UpdateStatus = config.UpdatedFailed
		}
	}
}

// 创建
func (pb *Porkbun) create(domain *config.Domain, recordType string, ipAddr string) {
	var response PorkbunResponse

	err := pb.request(
		porkbunEndpoint+fmt.Sprintf("/create/%s", domain.DomainName),
		&PorkbunDomainCreateOrUpdateVO{
			PorkbunApiKey: &PorkbunApiKey{
				AccessKey: pb.DNSConfig.ID,
				SecretKey: pb.DNSConfig.Secret,
			},
			PorkbunDomainRecord: &PorkbunDomainRecord{
				Name:    &domain.SubDomain,
				Type:    &recordType,
				Content: &ipAddr,
				Ttl:     &pb.TTL,
			},
		},
		&response,
	)

	if err != nil {
		util.Log("新增域名解析 %s 失败! 异常信息: %s", domain, err)
		domain.UpdateStatus = config.UpdatedFailed
		return
	}

	if response.Status == "SUCCESS" {
		util.Log("新增域名解析 %s 成功! IP: %s", domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	} else {
		util.Log("新增域名解析 %s 失败! 异常信息: %s", domain, response.Status)
		domain.UpdateStatus = config.UpdatedFailed
	}
}

// 修改
func (pb *Porkbun) modify(record *PorkbunDomainQueryResponse, domain *config.Domain, recordType string, ipAddr string) {

	// 相同不修改
	if len(record.Records) > 0 && *record.Records[0].Content == ipAddr {
		util.Log("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
		return
	}

	var response PorkbunResponse

	err := pb.request(
		porkbunEndpoint+fmt.Sprintf("/editByNameType/%s/%s/%s", domain.DomainName, recordType, domain.SubDomain),
		&PorkbunDomainCreateOrUpdateVO{
			PorkbunApiKey: &PorkbunApiKey{
				AccessKey: pb.DNSConfig.ID,
				SecretKey: pb.DNSConfig.Secret,
			},
			PorkbunDomainRecord: &PorkbunDomainRecord{
				Content: &ipAddr,
				Ttl:     &pb.TTL,
			},
		},
		&response,
	)

	if err != nil {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, err)
		domain.UpdateStatus = config.UpdatedFailed
		return
	}

	if response.Status == "SUCCESS" {
		util.Log("更新域名解析 %s 成功! IP: %s", domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	} else {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, response.Status)
		domain.UpdateStatus = config.UpdatedFailed
	}
}

// request 统一请求接口
func (pb *Porkbun) request(url string, data interface{}, result interface{}) (err error) {
	jsonStr := make([]byte, 0)
	if data != nil {
		jsonStr, _ = json.Marshal(data)
	}
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := util.CreateHTTPClient()
	resp, err := client.Do(req)
	err = util.GetHTTPResponse(resp, err, result)

	return
}
