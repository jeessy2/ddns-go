package dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

const (
	rainyunEndpoint = "https://api.v2.rainyun.com"
)

// https://s.apifox.cn/a4595cc8-44c5-4678-a2a3-eed7738dab03/api-153559362
// Rainyun Rainyun
type Rainyun struct {
	DNS        config.DNS
	Domains    config.Domains
	TTL        int
	httpClient *http.Client
}

// RainyunRecord 雨云DNS记录
type RainyunRecord struct {
	RecordID int64  `json:"record_id"`
	Host     string `json:"host"`
	Type     string `json:"type"`
	Value    string `json:"value"`
	Line     string `json:"line"`
	TTL      int    `json:"ttl"`
	Level    int    `json:"level"`
}

// RainyunResp 雨云API通用响应
type RainyunResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// Init 初始化
func (rainyun *Rainyun) Init(dnsConf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache) {
	rainyun.Domains.Ipv4Cache = ipv4cache
	rainyun.Domains.Ipv6Cache = ipv6cache
	rainyun.DNS = dnsConf.DNS
	rainyun.Domains.GetNewIp(dnsConf)
	if dnsConf.TTL == "" {
		// 默认600s
		rainyun.TTL = 600
	} else {
		ttlInt, _ := strconv.Atoi(dnsConf.TTL)
		rainyun.TTL = ttlInt
	}
	rainyun.httpClient = dnsConf.GetHTTPClient()
}

// AddUpdateDomainRecords 添加或更新IPv4/IPv6记录
func (rainyun *Rainyun) AddUpdateDomainRecords() (domains config.Domains) {
	rainyun.addUpdateDomainRecords("A")
	rainyun.addUpdateDomainRecords("AAAA")
	return rainyun.Domains
}

func (rainyun *Rainyun) addUpdateDomainRecords(recordType string) {
	ipAddr, domains := rainyun.Domains.GetNewIpResult(recordType)
	if ipAddr == "" {
		return
	}

	for _, domain := range domains {
		// 获取Domain ID
		domainID := rainyun.DNS.ID

		// 获取记录列表
		records, err := rainyun.getRecordList(domainID)
		if err != nil {
			util.Log("查询域名信息发生异常! %s", err)
			domain.UpdateStatus = config.UpdatedFailed
			continue
		}

		// 查找匹配的记录
		var recordSelected *RainyunRecord
		for i := range records {
			if strings.EqualFold(records[i].Host, domain.GetSubDomain()) &&
				strings.EqualFold(records[i].Type, recordType) {
				recordSelected = &records[i]
				break
			}
		}

		if recordSelected != nil {
			// 更新记录
			rainyun.modify(domainID, recordSelected, domain, ipAddr)
		} else {
			// 新增记录
			rainyun.create(domainID, domain, recordType, ipAddr)
		}
	}
}

// getRecordList 获取域名记录列表
func (rainyun *Rainyun) getRecordList(domainID string) ([]RainyunRecord, error) {
	query := url.Values{}
	query.Set("limit", "100")
	query.Set("page_no", "1")

	var result struct {
		TotalRecords int             `json:"TotalRecords"`
		Records      []RainyunRecord `json:"Records"`
	}
	err := rainyun.request(
		http.MethodGet,
		fmt.Sprintf("/product/domain/%s/dns/", url.PathEscape(domainID)),
		query,
		nil,
		&result,
	)
	if err != nil {
		return nil, err
	}
	return result.Records, nil
}

// create 创建DNS记录
func (rainyun *Rainyun) create(domainID string, domain *config.Domain, recordType string, ipAddr string) {
	record := &RainyunRecord{
		Host:  domain.GetSubDomain(),
		Type:  recordType,
		Value: ipAddr,
		Line:  "DEFAULT",
		TTL:   rainyun.TTL,
		Level: 10,
	}

	err := rainyun.createRecord(domainID, record)
	if err != nil {
		util.Log("新增域名解析 %s 失败! 异常信息: %s", domain, err)
		domain.UpdateStatus = config.UpdatedFailed
		return
	}

	util.Log("新增域名解析 %s 成功! IP: %s", domain, ipAddr)
	domain.UpdateStatus = config.UpdatedSuccess
}

// createRecord 发送POST请求创建记录
func (rainyun *Rainyun) createRecord(domainID string, record *RainyunRecord) error {
	payload := map[string]any{
		"host":      record.Host,
		"line":      record.Line,
		"level":     record.Level,
		"ttl":       record.TTL,
		"type":      record.Type,
		"value":     record.Value,
		"record_id": 0,
	}

	byt, _ := json.Marshal(payload)
	return rainyun.request(
		http.MethodPost,
		fmt.Sprintf("/product/domain/%s/dns", url.PathEscape(domainID)),
		nil,
		byt,
		nil,
	)
}

// modify 修改DNS记录
func (rainyun *Rainyun) modify(domainID string, record *RainyunRecord, domain *config.Domain, ipAddr string) {
	if record.Value == ipAddr {
		util.Log("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
		return
	}

	record.Value = ipAddr
	record.TTL = rainyun.TTL

	err := rainyun.patchRecord(domainID, record)
	if err != nil {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, err)
		domain.UpdateStatus = config.UpdatedFailed
		return
	}

	util.Log("更新域名解析 %s 成功! IP: %s", domain, ipAddr)
	domain.UpdateStatus = config.UpdatedSuccess
}

// patchRecord 发送PATCH请求更新记录
func (rainyun *Rainyun) patchRecord(domainID string, record *RainyunRecord) error {
	payload := map[string]any{
		"host":      record.Host,
		"line":      record.Line,
		"level":     record.Level,
		"ttl":       record.TTL,
		"type":      record.Type,
		"value":     record.Value,
		"record_id": record.RecordID,
	}

	byt, _ := json.Marshal(payload)
	return rainyun.request(
		http.MethodPatch,
		fmt.Sprintf("/product/domain/%s/dns", url.PathEscape(domainID)),
		nil,
		byt,
		nil,
	)
}

// request 统一请求接口
func (rainyun *Rainyun) request(method string, path string, query url.Values, body []byte, result any) error {
	u, err := url.Parse(rainyunEndpoint)
	if err != nil {
		return err
	}
	u.Path = path
	if query != nil {
		u.RawQuery = query.Encode()
	}

	var reader *bytes.Reader
	if body == nil {
		reader = bytes.NewReader(nil)
	} else {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, u.String(), reader)
	if err != nil {
		return err
	}
	// 认证
	req.Header.Set("x-api-key", rainyun.DNS.Secret)
	if method == http.MethodPost || method == http.MethodPatch || method == http.MethodPut {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := rainyun.httpClient.Do(req)
	if err != nil {
		return err
	}

	var apiResp RainyunResp
	err = util.GetHTTPResponse(resp, err, &apiResp)
	if err != nil {
		return err
	}
	if apiResp.Code != 200 {
		if apiResp.Message != "" {
			return fmt.Errorf("%s", apiResp.Message)
		}
		return fmt.Errorf("Rainyun API error, code=%d", apiResp.Code)
	}

	if result == nil {
		return nil
	}

	dataBytes, err := json.Marshal(apiResp.Data)
	if err != nil {
		return err
	}
	return json.Unmarshal(dataBytes, result)
}
