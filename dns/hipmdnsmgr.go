package dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

const (
	hipmDnsMgrEndpoint = "https://dnsmgr.example.com"
)

// HiPMDnsMgr HiPM DNSMgr 提供商
// 参考实现: c:\Users\HINS\Documents\Trae\DNSMgr-1\server\src\lib\dns\providers\dnsmgr.ts
type HiPMDnsMgr struct {
	DNS        config.DNS
	Domains    config.Domains
	TTL        string
	lastIpv4   string
	lastIpv6   string
	httpClient *http.Client
}

// DnsMgrApiResponse DNSMgr API 响应结构
// 对应 dnsmgr.ts 中的 DnsMgrApiResponse<T>
type DnsMgrApiResponse struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
	Msg  string          `json:"msg"`
}

// DnsMgrDomain DNSMgr 域名结构
// 对应 dnsmgr.ts 中的 DnsMgrDomain
type DnsMgrDomain struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	AccountID   int    `json:"account_id"`
	ThirdID     string `json:"third_id"`
	RecordCount int    `json:"record_count"`
}

// DnsMgrRecord DNSMgr 记录结构
// 对应 dnsmgr.ts 中的 DnsMgrRecord
type DnsMgrRecord struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	Line      string `json:"line"`
	TTL       int    `json:"ttl"`
	MX        int    `json:"mx"`
	Weight    int    `json:"weight"`
	Status    int    `json:"status"`
	Remark    string `json:"remark"`
	UpdatedAt string `json:"updated_at"`
	Proxiable bool   `json:"proxiable"`
	Cloudflare *struct {
		Proxied   bool `json:"proxied"`
		Proxiable bool `json:"proxiable"`
	} `json:"cloudflare"`
}

// DnsMgrRecordList DNSMgr 记录列表响应
// 对应 dnsmgr.ts 中的 PageResult<DnsMgrRecord>
type DnsMgrRecordList struct {
	Total int            `json:"total"`
	List  []DnsMgrRecord `json:"list"`
}

// Init 初始化
func (h *HiPMDnsMgr) Init(dnsConf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache) {
	h.Domains.Ipv4Cache = ipv4cache
	h.Domains.Ipv6Cache = ipv6cache
	h.lastIpv4 = ipv4cache.Addr
	h.lastIpv6 = ipv6cache.Addr

	h.DNS = dnsConf.DNS
	h.Domains.GetNewIp(dnsConf)
	if dnsConf.TTL == "" {
		h.TTL = "600"
	} else {
		h.TTL = dnsConf.TTL
	}
	h.httpClient = dnsConf.GetHTTPClient()
}

// AddUpdateDomainRecords 添加或更新IPv4/IPv6记录
func (h *HiPMDnsMgr) AddUpdateDomainRecords() config.Domains {
	h.addUpdateDomainRecords("A")
	h.addUpdateDomainRecords("AAAA")
	return h.Domains
}

func (h *HiPMDnsMgr) addUpdateDomainRecords(recordType string) {
	ipAddr, domains := h.Domains.GetNewIpResult(recordType)

	if ipAddr == "" {
		return
	}

	// 防止多次发送请求
	if recordType == "A" {
		if h.lastIpv4 == ipAddr {
			util.Log("你的IPv4未变化, 未触发 %s 请求", "HiPMDnsMgr")
			return
		}
	} else {
		if h.lastIpv6 == ipAddr {
			util.Log("你的IPv6未变化, 未触发 %s 请求", "HiPMDnsMgr")
			return
		}
	}

	for _, domain := range domains {
		err := h.updateRecord(domain, ipAddr, recordType)
		if err != nil {
			util.Log("HiPMDnsMgr更新记录失败, 域名: %s, IP: %s, 错误: %s", domain, ipAddr, err)
			domain.UpdateStatus = config.UpdatedFailed
		} else {
			util.Log("HiPMDnsMgr更新记录成功, 域名: %s, IP: %s", domain, ipAddr)
			domain.UpdateStatus = config.UpdatedSuccess
		}
	}
}

// updateRecord 更新或创建 DNS 记录
func (h *HiPMDnsMgr) updateRecord(domain *config.Domain, ipAddr string, recordType string) error {
	baseURL := h.DNS.ID
	if baseURL == "" {
		baseURL = hipmDnsMgrEndpoint
	}
	apiToken := h.DNS.Secret
	if apiToken == "" {
		return fmt.Errorf("API Token 不能为空")
	}

	// 获取域名ID
	domainID, err := h.getDomainID(baseURL, apiToken, domain.DomainName)
	if err != nil {
		return fmt.Errorf("获取域名ID失败: %w", err)
	}

	// 获取现有记录
	record, err := h.getRecord(baseURL, apiToken, domainID, domain.SubDomain, recordType)
	if err != nil {
		return fmt.Errorf("获取记录失败: %w", err)
	}

	ttl, _ := strconv.Atoi(h.TTL)
	if ttl == 0 {
		ttl = 600
	}

	if record != nil {
		// 更新现有记录
		return h.updateExistingRecord(baseURL, apiToken, domainID, record.ID, domain.SubDomain, recordType, ipAddr, ttl)
	}
	// 创建新记录
	return h.createRecord(baseURL, apiToken, domainID, domain.SubDomain, recordType, ipAddr, ttl)
}

// getHeaders 获取请求头
// 参考 dnsmgr.ts 中的 getHeaders() 方法
func (h *HiPMDnsMgr) getHeaders(apiToken string) map[string]string {
	return map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + apiToken,
	}
}

// request 发送 HTTP 请求
// 参考 dnsmgr.ts 中的 request<T>() 方法
func (h *HiPMDnsMgr) request(baseURL, apiToken, method, path string, body interface{}) (*DnsMgrApiResponse, error) {
	// 参考 dnsmgr.ts 中的 URL 处理方式
	// Ensure baseUrl doesn't end with /api and path starts with /
	base := strings.TrimSuffix(baseURL, "/")
	base = strings.TrimSuffix(base, "/api")
	
	normalizedPath := path
	if !strings.HasPrefix(normalizedPath, "/") {
		normalizedPath = "/" + normalizedPath
	}
	
	url := base + "/api" + normalizedPath
	
	util.Log("HiPMDnsMgr请求: %s %s", method, url)
	
	var bodyReader *bytes.Buffer
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体失败: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	} else {
		bodyReader = bytes.NewBuffer(nil)
	}
	
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	
	// 设置请求头
	headers := h.getHeaders(apiToken)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	
	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var apiResp DnsMgrApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}
	
	util.Log("HiPMDnsMgr响应: code=%d, msg=%s", apiResp.Code, apiResp.Msg)
	
	return &apiResp, nil
}

// getDomainID 获取域名ID
// 参考 dnsmgr.ts 中的 getDomainList() 方法
func (h *HiPMDnsMgr) getDomainID(baseURL, apiToken, domainName string) (int, error) {
	// DnsMgr API returns array directly in data, not { total, list } format
	apiResp, err := h.request(baseURL, apiToken, "GET", "/domains?page=1&pageSize=100", nil)
	if err != nil {
		return 0, err
	}
	
	if apiResp.Code != 0 {
		return 0, fmt.Errorf("API错误: %s", apiResp.Msg)
	}
	
	var domains []DnsMgrDomain
	if err := json.Unmarshal(apiResp.Data, &domains); err != nil {
		return 0, err
	}
	
	// 参考 dnsmgr.ts 中的域名过滤逻辑
	for _, d := range domains {
		if d.Name == domainName {
			return d.ID, nil
		}
	}
	
	return 0, fmt.Errorf("域名 %s 未找到", domainName)
}

// getRecord 获取 DNS 记录
// 参考 dnsmgr.ts 中的 getDomainRecords() 方法
func (h *HiPMDnsMgr) getRecord(baseURL, apiToken string, domainID int, subDomain, recordType string) (*DnsMgrRecord, error) {
	path := fmt.Sprintf("/domains/%d/records?page=1&pageSize=100&subdomain=%s&type=%s", 
		domainID, subDomain, recordType)
	
	apiResp, err := h.request(baseURL, apiToken, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	
	if apiResp.Code != 0 {
		return nil, fmt.Errorf("API错误: %s", apiResp.Msg)
	}
	
	var recordList DnsMgrRecordList
	if err := json.Unmarshal(apiResp.Data, &recordList); err != nil {
		return nil, err
	}
	
	// 查找匹配的记录
	for _, r := range recordList.List {
		if r.Name == subDomain && r.Type == recordType {
			return &r, nil
		}
	}
	
	return nil, nil
}

// createRecord 创建新记录
// 参考 dnsmgr.ts 中的 addDomainRecord() 方法
func (h *HiPMDnsMgr) createRecord(baseURL, apiToken string, domainID int, name, recordType, value string, ttl int) error {
	path := fmt.Sprintf("/domains/%d/records", domainID)
	
	// 参考 dnsmgr.ts 中的请求体构造
	body := map[string]interface{}{
		"name":  name,
		"type":  recordType,
		"value": value,
		"ttl":   ttl,
		"line":  "0",
	}
	
	if recordType == "MX" {
		body["mx"] = 10
	}
	
	apiResp, err := h.request(baseURL, apiToken, "POST", path, body)
	if err != nil {
		return err
	}
	
	if apiResp.Code != 0 {
		return fmt.Errorf("API错误: %s", apiResp.Msg)
	}
	
	return nil
}

// updateExistingRecord 更新现有记录
// 参考 dnsmgr.ts 中的 updateDomainRecord() 方法
func (h *HiPMDnsMgr) updateExistingRecord(baseURL, apiToken string, domainID int, recordID, name, recordType, value string, ttl int) error {
	path := fmt.Sprintf("/domains/%d/records/%s", domainID, recordID)
	
	// 参考 dnsmgr.ts 中的请求体构造
	body := map[string]interface{}{
		"name":  name,
		"type":  recordType,
		"value": value,
		"ttl":   ttl,
		"line":  "0",
	}
	
	if recordType == "MX" {
		body["mx"] = 10
	}
	
	apiResp, err := h.request(baseURL, apiToken, "PUT", path, body)
	if err != nil {
		return err
	}
	
	if apiResp.Code != 0 {
		return fmt.Errorf("API错误: %s", apiResp.Msg)
	}
	
	return nil
}
