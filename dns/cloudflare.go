package dns

import (
	"bytes"
	"ddns-go/config"
	"ddns-go/util"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	zonesAPI string = "https://api.cloudflare.com/client/v4/zones"
)

// Cloudflare Cloudflare实现
type Cloudflare struct {
	DNSConfig config.DNSConfig
	Domains
}

// CloudflareZonesResp cloudflare zones返回结果
type CloudflareZonesResp struct {
	CloudflareStatus
	Result []struct {
		ID     string
		Name   string
		Status string
		Paused bool
	}
}

// CloudflareRecordsResp records
type CloudflareRecordsResp struct {
	CloudflareStatus
	Result []CloudflareRecord
}

// CloudflareRecord 记录实体
type CloudflareRecord struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Proxied bool   `json:"proxied"`
	TTL     int    `json:"ttl"`
}

// CloudflareStatus 公共状态
type CloudflareStatus struct {
	Success  bool
	Messages []string
}

// Init 初始化
func (cf *Cloudflare) Init(conf *config.Config) {
	cf.DNSConfig = conf.DNS
	cf.Domains.ParseDomain(conf)
}

// AddUpdateIpv4DomainRecords 添加或更新IPV4记录
func (cf *Cloudflare) AddUpdateIpv4DomainRecords() {
	cf.addUpdateDomainRecords("A")
}

// AddUpdateIpv6DomainRecords 添加或更新IPV6记录
func (cf *Cloudflare) AddUpdateIpv6DomainRecords() {
	cf.addUpdateDomainRecords("AAAA")
}

func (cf *Cloudflare) addUpdateDomainRecords(recordType string) {
	ipAddr := cf.Ipv4Addr
	domains := cf.Ipv4Domains
	if recordType == "AAAA" {
		ipAddr = cf.Ipv6Addr
		domains = cf.Ipv6Domains
	}

	if ipAddr == "" {
		return
	}

	for _, domain := range domains {
		// get zone
		result, err := cf.getZones(domain)
		if err != nil || len(result.Result) != 1 {
			return
		}
		zoneID := result.Result[0].ID

		var records CloudflareRecordsResp
		// getDomains 最多更新前50条
		cf.request(
			"GET",
			fmt.Sprintf(zonesAPI+"/%s/dns_records?type=%s&name=%s&per_page=50", zoneID, recordType, domain),
			nil,
			&records,
		)

		if records.Success && len(records.Result) > 0 {
			// 更新
			cf.modify(records, zoneID, domain, recordType, ipAddr)
		} else {
			// 新增
			cf.create(zoneID, domain, recordType, ipAddr)
		}
	}
}

// 创建
func (cf *Cloudflare) create(zoneID string, domain *Domain, recordType string, ipAddr string) {
	record := &CloudflareRecord{
		Type:    recordType,
		Name:    domain.String(),
		Content: ipAddr,
		Proxied: false,
		// auto ttl
		TTL: 1,
	}
	var status CloudflareStatus
	err := cf.request(
		"POST",
		fmt.Sprintf(zonesAPI+"/%s/dns_records", zoneID),
		record,
		&status,
	)
	if err == nil && status.Success {
		log.Printf("新增域名解析 %s 成功！IP: %s", domain, ipAddr)
	} else {
		log.Printf("新增域名解析 %s 失败！Messages: %s", domain, status.Messages)
	}
}

// 修改
func (cf *Cloudflare) modify(result CloudflareRecordsResp, zoneID string, domain *Domain, recordType string, ipAddr string) {

	for _, record := range result.Result {
		// 相同不修改
		if record.Content == ipAddr {
			log.Printf("你的IP %s 没有变化, 未有任何操作。域名 %s", ipAddr, domain)
			continue
		}
		var status CloudflareStatus
		record.Content = ipAddr

		err := cf.request(
			"PUT",
			fmt.Sprintf(zonesAPI+"/%s/dns_records/%s", zoneID, record.ID),
			record,
			&status,
		)

		if err == nil && status.Success {
			log.Printf("更新域名解析 %s 成功！IP: %s", domain, ipAddr)
		} else {
			log.Printf("更新域名解析 %s 失败！Messages: %s", domain, status.Messages)
		}
	}
}

// 获得域名记录列表
func (cf *Cloudflare) getZones(domain *Domain) (result CloudflareZonesResp, err error) {
	err = cf.request(
		"GET",
		fmt.Sprintf(zonesAPI+"?name=%s&status=%s&per_page=%s", domain.DomainName, "active", "50"),
		nil,
		&result,
	)

	return
}

// request 统一请求接口
func (cf *Cloudflare) request(method string, url string, data interface{}, result interface{}) (err error) {
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
		log.Println("http.NewRequest失败. Error: ", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+cf.DNSConfig.Secret)
	req.Header.Set("Content-Type", "application/json")

	clt := http.Client{}
	clt.Timeout = 1 * time.Minute
	resp, err := clt.Do(req)
	err = util.GetHTTPResponse(resp, url, err, result)

	return
}
