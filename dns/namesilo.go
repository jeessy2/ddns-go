package dns

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jeessy2/ddns-go/v5/config"
	"github.com/jeessy2/ddns-go/v5/util"
)

var (
	nameSiloListRecordEndpoint   = "https://www.namesilo.com/api/dnsListRecords?version=1&type=xml&key=#{password}&domain=#{domain}"
	nameSiloAddRecordEndpoint    = "https://www.namesilo.com/api/dnsAddRecord?version=1&type=xml&key=#{password}&domain=#{domain}&rrtype=#{recordType}&rrvalue=#{ip}&rrttl=#{ttl}"
	nameSiloUpdateRecordEndpoint = "https://www.namesilo.com/api/dnsUpdateRecord?version=1&type=xml&key=#{password}&domain=#{domain}&rrid=#{recordID}&rrvalue=#{ip}&rrttl=#{ttl}"
)

const (
	minTTL = 3600
	maxTTL = 2592000
)

// NameSilo Domain
type NameSilo struct {
	DNS      config.DNS
	Domains  config.Domains
	lastIpv4 string
	lastIpv6 string
	TTL      int
}

// NameSiloResp 修改域名解析结果
type NameSiloResp struct {
	XMLName xml.Name      `xml:"namesilo"`
	Request Request       `xml:"request"`
	Reply   ReplyResponse `xml:"reply"`
}

type ReplyResponse struct {
	Code     int    `xml:"code"`
	Detail   string `xml:"detail"`
	RecordID string `xml:"record_id"`
}

type NameSiloDNSListRecordResp struct {
	XMLName xml.Name `xml:"namesilo"`
	Request Request  `xml:"request"`
	Reply   Reply    `xml:"reply"`
}

type Request struct {
	Operation string `xml:"operation"`
	IP        string `xml:"ip"`
}

type Reply struct {
	Code          int              `xml:"code"`
	Detail        string           `xml:"detail"`
	ResourceItems []ResourceRecord `xml:"resource_record"`
}

type ResourceRecord struct {
	RecordID string `xml:"record_id"`
	Type     string `xml:"type"`
	Host     string `xml:"host"`
	Value    string `xml:"value"`
	TTL      int    `xml:"ttl"`
	Distance int    `xml:"distance"`
}

// Init 初始化
func (ns *NameSilo) Init(dnsConf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache) {
	ns.Domains.Ipv4Cache = ipv4cache
	ns.Domains.Ipv6Cache = ipv6cache
	ns.lastIpv4 = ipv4cache.Addr
	ns.lastIpv6 = ipv6cache.Addr

	ns.DNS = dnsConf.DNS
	ns.Domains.GetNewIp(dnsConf)
	// 默认3600s 官方支持 最小3600s 最大2592000s
	ns.TTL = minTTL
	if dnsConf.TTL != "" {
		ttl, err := strconv.Atoi(dnsConf.TTL)
		if err == nil && (ttl > minTTL || ttl <= maxTTL) {
			ns.TTL = ttl
		}
	}
}

// AddUpdateDomainRecords 添加或更新IPv4/IPv6记录
func (ns *NameSilo) AddUpdateDomainRecords() config.Domains {
	ns.addUpdateDomainRecords("A")
	ns.addUpdateDomainRecords("AAAA")
	return ns.Domains
}

func (ns *NameSilo) addUpdateDomainRecords(recordType string) {
	ipAddr, domains := ns.Domains.GetNewIpResult(recordType)

	if ipAddr == "" {
		return
	}

	for _, domain := range domains {
		// 拿到DNS记录列表，从列表中去取对应域名的id，有id进行修改，没ID进行新增
		records, err := ns.listRecords(domain)
		if err != nil {
			log.Printf("获取域名列表 %s 失败！", domain)
			domain.UpdateStatus = config.UpdatedFailed
			return
		}
		items := records.Reply.ResourceItems
		record := findResourceRecord(items, recordType, domain.String())
		var isAdd bool
		var recordID string
		if record == nil {
			isAdd = true
		} else {
			recordID = record.RecordID
		}
		ns.modify(domain, recordID, recordType, ipAddr, isAdd)
	}
}

// 修改
func (ns *NameSilo) modify(domain *config.Domain, recordID, recordType, ipAddr string, isAdd bool) {
	var err error
	var result string
	var requestType string
	if isAdd {
		requestType = "新增"
		if domain.GetSubDomain() != "@" {
			nameSiloAddRecordEndpoint += "&rrhost=#{host}"
		}
		result, err = ns.request(ipAddr, domain, "", recordType, nameSiloAddRecordEndpoint)
	} else {
		requestType = "修改"
		if domain.GetSubDomain() != "@" {
			nameSiloUpdateRecordEndpoint += "&rrhost=#{host}"
		}
		result, err = ns.request(ipAddr, domain, recordID, "", nameSiloUpdateRecordEndpoint)
	}
	if err != nil {
		log.Printf("修改域名解析 %s 失败！", domain)
		domain.UpdateStatus = config.UpdatedFailed
		return
	}
	var resp NameSiloResp
	xml.Unmarshal([]byte(result), &resp)
	if resp.Reply.Code == 300 {
		log.Printf("%s 域名解析 %s 成功！IP: %s\n", requestType, domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	} else {
		log.Printf("%s 域名解析 %s 失败！Deatil: %s\n", requestType, domain, resp.Reply.Detail)
		domain.UpdateStatus = config.UpdatedFailed
	}
}

func (ns *NameSilo) listRecords(domain *config.Domain) (resp NameSiloDNSListRecordResp, err error) {
	result, err := ns.request("", domain, "", "", nameSiloListRecordEndpoint)
	err = xml.Unmarshal([]byte(result), &resp)
	return
}

// request 统一请求接口
func (ns *NameSilo) request(ipAddr string, domain *config.Domain, recordID, recordType, url string) (result string, err error) {
	url = strings.ReplaceAll(url, "#{host}", domain.GetSubDomain())
	url = strings.ReplaceAll(url, "#{domain}", domain.DomainName)
	url = strings.ReplaceAll(url, "#{password}", ns.DNS.Secret)
	url = strings.ReplaceAll(url, "#{recordID}", recordID)
	url = strings.ReplaceAll(url, "#{recordType}", recordType)
	url = strings.ReplaceAll(url, "#{ip}", ipAddr)
	url = strings.ReplaceAll(url, "#{ttl}", strconv.Itoa(ns.TTL))
	req, err := http.NewRequest(
		http.MethodGet,
		url,
		http.NoBody,
	)

	if err != nil {
		log.Println("http.NewRequest失败. Error: ", err)
		return
	}

	client := util.CreateHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		log.Println("client.Do失败. Error: ", err)
		return
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	result = string(data)
	return
}

func findResourceRecord(data []ResourceRecord, recordType, domain string) *ResourceRecord {
	for i := 0; i < len(data); i++ {
		if data[i].Host == domain && data[i].Type == recordType {
			return &data[i]
		}
	}
	return nil
}
