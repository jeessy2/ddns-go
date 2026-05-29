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

	// Prevent duplicate requests
	if recordType == "A" {
		if h.lastIpv4 == ipAddr {
			util.Log("IPv4 address unchanged, skipping %s request", "HiPMDnsMgr")
			return
		}
	} else {
		if h.lastIpv6 == ipAddr {
			util.Log("IPv6 address unchanged, skipping %s request", "HiPMDnsMgr")
			return
		}
	}

	for _, domain := range domains {
		err := h.updateRecord(domain, ipAddr, recordType)
		if err != nil {
			util.Log("HiPMDnsMgr failed to update record, domain: %s, IP: %s, error: %s", domain, ipAddr, err)
			domain.UpdateStatus = config.UpdatedFailed
		} else {
			util.Log("HiPMDnsMgr updated record successfully, domain: %s, IP: %s", domain, ipAddr)
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
		return fmt.Errorf("API token cannot be empty")
	}

	// Get domain ID
	domainID, err := h.getDomainID(baseURL, apiToken, domain.DomainName)
	if err != nil {
		return fmt.Errorf("failed to get domain ID: %w", err)
	}

	// Get existing record
	record, err := h.getRecord(baseURL, apiToken, domainID, domain.SubDomain, recordType)
	if err != nil {
		return fmt.Errorf("failed to get record: %w", err)
	}

	ttl, _ := strconv.Atoi(h.TTL)
	if ttl == 0 {
		ttl = 600
	}

	if record != nil {
		// Update existing record
		return h.updateExistingRecord(baseURL, apiToken, domainID, record.ID, domain.SubDomain, recordType, ipAddr, ttl)
	}
	// Create new record
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
	
	util.Log("HiPMDnsMgr request: %s %s", method, url)
	
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
	
	util.Log("HiPMDnsMgr response: code=%d, msg=%s", apiResp.Code, apiResp.Msg)
	
	// Log detailed error information when code != 0
	if apiResp.Code != 0 {
		util.Log("HiPMDnsMgr error details: url=%s, method=%s, code=%d, msg=%s", url, method, apiResp.Code, apiResp.Msg)
	}
	
	return &apiResp, nil
}

// getDomainID Get domain ID
// Prefer using keyword parameter for direct query, with list matching as fallback
func (h *HiPMDnsMgr) getDomainID(baseURL, apiToken, domainName string) (int, error) {
	// Method 1: Use keyword parameter for direct query (efficient)
	util.Log("HiPMDnsMgr querying domain by keyword: %s", domainName)
	path := fmt.Sprintf("/domains?page=1&pageSize=1&keyword=%s", domainName)
	
	apiResp, err := h.request(baseURL, apiToken, "GET", path, nil)
	if err != nil {
		return 0, err
	}
	
	if apiResp.Code != 0 {
		return 0, fmt.Errorf("API error: %s", apiResp.Msg)
	}
	
	var domains []DnsMgrDomain
	
	// Smart detection: support both array and object formats
	var rawData interface{}
	if err := json.Unmarshal(apiResp.Data, &rawData); err != nil {
		return 0, fmt.Errorf("failed to parse response data: %w", err)
	}
	
	switch v := rawData.(type) {
	case []interface{}:
		jsonData, _ := json.Marshal(v)
		if err := json.Unmarshal(jsonData, &domains); err != nil {
			return 0, fmt.Errorf("failed to parse domain list: %w", err)
		}
	case map[string]interface{}:
		if listData, ok := v["list"]; ok {
			jsonData, _ := json.Marshal(listData)
			if err := json.Unmarshal(jsonData, &domains); err != nil {
				return 0, fmt.Errorf("failed to parse domain list: %w", err)
			}
		} else {
			return 0, fmt.Errorf("invalid response format: missing list field")
		}
	default:
		return 0, fmt.Errorf("unknown response data format: %T", rawData)
	}
	
	// Check if exact match is found
	for _, d := range domains {
		if d.Name == domainName {
			util.Log("HiPMDnsMgr keyword query successful: domain=%s, id=%d", domainName, d.ID)
			return d.ID, nil
		}
	}
	
	// Method 2: If keyword query not found, use list matching as fallback (compatible with old API)
	// Paginate through all domains to find the target
	util.Log("HiPMDnsMgr keyword query not found, trying paginated list matching: %s", domainName)
	
	const pageSize = 100
	currentPage := 1
	
	for {
		path := fmt.Sprintf("/domains?page=%d&pageSize=%d", currentPage, pageSize)
		apiResp, err := h.request(baseURL, apiToken, "GET", path, nil)
		if err != nil {
			return 0, fmt.Errorf("paginated query failed at page %d: %w", currentPage, err)
		}
		
		if apiResp.Code != 0 {
			return 0, fmt.Errorf("paginated query API error at page %d: %s", currentPage, apiResp.Msg)
		}
		
		// Parse response with smart format detection
		var pageDomains []DnsMgrDomain
		var total int
		
		var rawData interface{}
		if err := json.Unmarshal(apiResp.Data, &rawData); err == nil {
			switch v := rawData.(type) {
			case []interface{}:
				jsonData, _ := json.Marshal(v)
				json.Unmarshal(jsonData, &pageDomains)
			case map[string]interface{}:
				if listData, ok := v["list"]; ok {
					jsonData, _ := json.Marshal(listData)
					json.Unmarshal(jsonData, &pageDomains)
				}
				if totalData, ok := v["total"]; ok {
					if t, ok := totalData.(float64); ok {
						total = int(t)
					}
				}
			}
		}
		
		// Search in current page
		for _, d := range pageDomains {
			if d.Name == domainName {
				util.Log("HiPMDnsMgr paginated list matching successful: domain=%s, id=%d, page=%d", domainName, d.ID, currentPage)
				return d.ID, nil
			}
		}
		
		// Check if we've reached the end
		if len(pageDomains) < pageSize || (total > 0 && currentPage*pageSize >= total) {
			break
		}
		
		currentPage++
		
		// Safety limit: stop after 10 pages (1000 domains)
		if currentPage > 10 {
			util.Log("HiPMDnsMgr warning: searched 10 pages (1000 domains), stopping to avoid excessive queries")
			break
		}
	}
	
	return 0, fmt.Errorf("domain %s not found", domainName)
}

// getRecord Get DNS record
// Paginate through all records to find the target
func (h *HiPMDnsMgr) getRecord(baseURL, apiToken string, domainID int, subDomain, recordType string) (*DnsMgrRecord, error) {
	const pageSize = 100
	currentPage := 1
	
	util.Log("HiPMDnsMgr querying record with pagination: domainID=%d, subdomain=%s, type=%s", domainID, subDomain, recordType)
	
	for {
		path := fmt.Sprintf("/domains/%d/records?page=%d&pageSize=%d&subdomain=%s&type=%s", 
			domainID, currentPage, pageSize, subDomain, recordType)
		
		apiResp, err := h.request(baseURL, apiToken, "GET", path, nil)
		if err != nil {
			return nil, fmt.Errorf("paginated record query failed at page %d: %w", currentPage, err)
		}
		
		if apiResp.Code != 0 {
			return nil, fmt.Errorf("paginated record query API error at page %d: %s", currentPage, apiResp.Msg)
		}
		
		var recordList DnsMgrRecordList
		if err := json.Unmarshal(apiResp.Data, &recordList); err != nil {
			return nil, fmt.Errorf("failed to parse record list: %w", err)
		}
		
		util.Log("HiPMDnsMgr page %d: found %d records (total: %d)", currentPage, len(recordList.List), recordList.Total)
		
		// Find matching record in current page
		for _, r := range recordList.List {
			if r.Name == subDomain && r.Type == recordType {
				util.Log("HiPMDnsMgr matched existing record: id=%s, value=%s, page=%d", r.ID, r.Value, currentPage)
				return &r, nil
			}
		}
		
		// Check if we've reached the end
		if len(recordList.List) < pageSize || (recordList.Total > 0 && currentPage*pageSize >= recordList.Total) {
			break
		}
		
		currentPage++
		
		// Safety limit: stop after 10 pages (1000 records)
		if currentPage > 10 {
			util.Log("HiPMDnsMgr warning: searched 10 pages (1000 records), stopping to avoid excessive queries")
			break
		}
	}
	
	util.Log("HiPMDnsMgr no matching record found after searching %d pages, will create new record", currentPage)
	return nil, nil
}

// createRecord Create new record
// Reference: addDomainRecord() method in dnsmgr.ts
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
		return fmt.Errorf("API error: %s", apiResp.Msg)
	}
	
	return nil
}

// updateExistingRecord Update existing record
// Reference: updateDomainRecord() method in dnsmgr.ts
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
		return fmt.Errorf("API error: %s", apiResp.Msg)
	}
	
	return nil
}
