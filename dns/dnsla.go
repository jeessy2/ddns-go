package dns

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
	"io"
	"net/http"
	"strconv"
)

const (
	recordList   string = "http://api.dns.la/api/recordList"
	recordModify string = "http://api.dns.la/api/record"
	recordCreate string = "http://api.dns.la/api/record"
)

// https://www.dns.la/docs/ApiDoc
// dnsla dnsla实现
type Dnsla struct {
	DNS     config.DNS
	Domains config.Domains
	TTL     int
}

// DnslaRecord
type DnslaRecord struct {
	ID   string `json:"id"`
	Host string `json:"host"`
	Type int    `json:"type"`
	Data string `json:"data"`
}

// DnslaRecordListResp recordListAPI结果
type DnslaRecordListResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Total   int           `json:"total"`
		Results []DnslaRecord `json:"results"`
	} `json:"data"`
}

// DnslaStatus DnslaStatus
type DnslaStatus struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Id string `json:"id"`
	} `json:"data"`
}

// Init 初始化
func (dnsla *Dnsla) Init(dnsConf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache) {
	dnsla.Domains.Ipv4Cache = ipv4cache
	dnsla.Domains.Ipv6Cache = ipv6cache
	dnsla.DNS = dnsConf.DNS
	dnsla.Domains.GetNewIp(dnsConf)
	if dnsConf.TTL == "" {
		// 默认600s
		dnsla.TTL = 600
	} else {
		ttlInt, _ := strconv.Atoi(dnsConf.TTL)
		dnsla.TTL = ttlInt
	}
}

// AddUpdateDomainRecords 添加或更新IPv4/IPv6记录
func (dnsla *Dnsla) AddUpdateDomainRecords() config.Domains {
	dnsla.addUpdateDomainRecords("A")
	dnsla.addUpdateDomainRecords("AAAA")
	return dnsla.Domains
}

func (dnsla *Dnsla) addUpdateDomainRecords(recordType string) {
	ipAddr, domains := dnsla.Domains.GetNewIpResult(recordType)
	if ipAddr == "" {
		return
	}
	for _, domain := range domains {
		resultByte, err := dnsla.getRecordList(domain, recordType)
		if err != nil {
			util.Log("查询域名信息发生异常! %s", err)
			domain.UpdateStatus = config.UpdatedFailed
			return
		}
		var jsonResult DnslaRecordListResp
		errU := json.Unmarshal(resultByte, &jsonResult)
		if errU != nil {
			util.Log(errU.Error())
			return
		}
		if jsonResult.Data.Total > 0 { // 默认第一个
			recordSelected := jsonResult.Data.Results[0]
			params := domain.GetCustomParams()
			if params.Has("id") {
				for i := 0; i < len(jsonResult.Data.Results); i++ {
					if jsonResult.Data.Results[i].ID == params.Get("id") {
						recordSelected = jsonResult.Data.Results[i]
					}
				}
			}
			// 更新
			dnsla.modify(recordSelected, domain, recordType, ipAddr)
		} else {
			// 新增
			dnsla.create(domain, recordType, ipAddr)
		}
	}
}

// 创建
func (dnsla *Dnsla) create(domain *config.Domain, recordType string, ipAddr string) {
	recordTypeInt := 1
	if recordType == "AAAA" {
		recordTypeInt = 28
	}
	type CreateParams struct {
		Domain string `json:"Domain"`
		Host   string `json:"Host"`
		Type   int    `json:"Type"`
		Data   string `json:"Data"`
		TTL    int    `json:"TTL"`
	}
	createParams := CreateParams{
		Domain: domain.DomainName,
		Host:   domain.GetSubDomain(),
		Type:   recordTypeInt,
		Data:   ipAddr,
		TTL:    dnsla.TTL,
	}
	jsonData, _ := json.Marshal(createParams)
	resultByte, err := dnsla.request("POST", recordCreate, jsonData)
	if err != nil {
		util.Log("新增域名解析 %s 失败! 异常信息: %s", domain, err)
		domain.UpdateStatus = config.UpdatedFailed
		return
	}
	var jsonResult DnslaStatus
	errU := json.Unmarshal(resultByte, &jsonResult)
	if errU != nil {
		util.Log(errU.Error())
		return
	}
	if jsonResult.Code == 200 {
		util.Log("新增域名解析 %s 成功! IP: %s", domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	} else {
		util.Log("新增域名解析 %s 失败! 异常信息: %s", domain, jsonResult.Msg)
		domain.UpdateStatus = config.UpdatedFailed
	}
}

// 修改
func (dnsla *Dnsla) modify(record DnslaRecord, domain *config.Domain, recordType string, ipAddr string) {
	// 相同不修改
	if record.Data == ipAddr {
		util.Log("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
		return
	}
	recordTypeInt := 1
	if recordType == "AAAA" {
		recordTypeInt = 28
	}
	type ModifyParams struct {
		ID   string `json:"Id"`
		Host string `json:"Host"`
		Type int    `json:"Type"`
		Data string `json:"Data"`
		TTL  int    `json:"TTL"`
	}
	modifyParams := ModifyParams{
		ID:   record.ID,
		Host: domain.GetSubDomain(),
		Type: recordTypeInt,
		Data: ipAddr,
		TTL:  dnsla.TTL,
	}
	jsonData, _ := json.Marshal(modifyParams)
	resultByte, err := dnsla.request("PUT", recordModify, jsonData)

	if err != nil {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, err)
		domain.UpdateStatus = config.UpdatedFailed
		return
	}

	var jsonResult DnslaStatus
	errU := json.Unmarshal(resultByte, &jsonResult)
	if errU != nil {
		util.Log(errU.Error())
		return
	}
	if jsonResult.Code == 200 {
		util.Log("更新域名解析 %s 成功! IP: %s", domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	} else {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, jsonResult.Msg)
		domain.UpdateStatus = config.UpdatedFailed
	}
}

// request sends a POST request to the given API with the given values.
func (dnsla *Dnsla) request(method, apiAddr string, values []byte) (body []byte, err error) {
	req, err := http.NewRequest(
		method,
		apiAddr,
		bytes.NewReader(values),
	)
	if err != nil {
		panic(err)
	}
	// 设置自定义 Headers
	byteBuff := []byte(dnsla.DNS.ID + ":" + dnsla.DNS.Secret)
	token := "Basic " + base64.StdEncoding.EncodeToString(byteBuff)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	// 4. 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	return
}

// 获得域名记录列表
func (dnsla *Dnsla) getRecordList(domain *config.Domain, typ string) (result []byte, err error) {
	recordTypeInt := "1"
	if typ == "AAAA" {
		recordTypeInt = "28"
	}
	params := domain.GetCustomParams()
	params.Set("domain", domain.DomainName)
	params.Set("host", domain.GetSubDomain())
	params.Set("type", recordTypeInt)
	params.Set("pageIndex", "1")
	params.Set("pageSize", "999")

	url := recordList + "?" + params.Encode()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	byteBuff := []byte(dnsla.DNS.ID + ":" + dnsla.DNS.Secret)
	token := "Basic " + base64.StdEncoding.EncodeToString(byteBuff)
	// 设置 Headers
	req.Header.Set("Authorization", token)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 读取响应
	result, errR := io.ReadAll(resp.Body)
	if errR != nil {
		util.Log(errR.Error())
		return
	}
	return
}
