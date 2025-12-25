package dns

// --- 子域名管理相关结构体 ---
// dnsheSubdomain 子域名信息
type dnsheSubdomain struct {
	ID         int    `json:"id"`
	Subdomain  string `json:"subdomain"`
	Rootdomain string `json:"rootdomain"`
	FullDomain string `json:"full_domain"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}

// dnsheListSubdomainsResp 列出子域名响应
type dnsheListSubdomainsResp struct {
	Success    bool             `json:"success"`
	Count      int              `json:"count"`
	Subdomains []dnsheSubdomain `json:"subdomains"`
	Error      string           `json:"error,omitempty"`
}

// dnsheRegisterReq 注册子域名请求参数
type dnsheRegisterReq struct {
	Subdomain  string `json:"subdomain"`
	Rootdomain string `json:"rootdomain"`
}

// dnsheRegisterResp 注册子域名响应
type dnsheRegisterResp struct {
	Success     bool        `json:"success"`
	Message     string      `json:"message"`
	SubdomainID interface{} `json:"subdomain_id"`
	FullDomain  string      `json:"full_domain"`
	Error       string      `json:"error,omitempty"`
}

// dnsheSubdomainDetailResp 获取子域名详情响应
type dnsheSubdomainDetailResp struct {
	Success    bool           `json:"success"`
	Subdomain  dnsheSubdomain `json:"subdomain"`
	DNSRecords []dnsheRecord  `json:"dns_records"`
	DNSCount   int            `json:"dns_count"`
	Error      string         `json:"error,omitempty"`
}

// dnsheDeleteSubdomainReq 删除子域名请求参数
type dnsheDeleteSubdomainReq struct {
	SubdomainID int `json:"subdomain_id"`
}

// dnsheDeleteSubdomainResp 删除子域名响应
type dnsheDeleteSubdomainResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// dnsheRenewSubdomainReq 续期子域名请求参数
type dnsheRenewSubdomainReq struct {
	SubdomainID int `json:"subdomain_id"`
}

// dnsheRenewSubdomainResp 续期子域名响应
type dnsheRenewSubdomainResp struct {
	Success         bool   `json:"success"`
	Message         string `json:"message"`
	SubdomainID     int    `json:"subdomain_id"`
	Subdomain       string `json:"subdomain"`
	PreviousExpires string `json:"previous_expires_at"`
	NewExpires      string `json:"new_expires_at"`
	RenewedAt       string `json:"renewed_at"`
	NeverExpires    int    `json:"never_expires"`
	Status          string `json:"status"`
	RemainingDays   int    `json:"remaining_days"`
	Error           string `json:"error,omitempty"`
}

// --- DNS 记录管理相关结构体 ---
// dnsheRecord DNS记录信息
type dnsheRecord struct {
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	Type      string      `json:"type"`
	Content   string      `json:"content"`
	TTL       int         `json:"ttl"`
	Priority  interface{} `json:"priority"`
	Proxied   bool        `json:"proxied,omitempty"`
	Status    string      `json:"status"`
	CreatedAt string      `json:"created_at,omitempty"`
}

// dnsheListRecordsResp 列出DNS记录响应
type dnsheListRecordsResp struct {
	Success bool           `json:"success"`
	Count   int            `json:"count"`
	Records []dnsheRecord  `json:"records"`
	Error   string         `json:"error,omitempty"`
}

// dnsheCreateRecordReq 创建DNS记录请求参数
type dnsheCreateRecordReq struct {
	SubdomainID int         `json:"subdomain_id"`
	Type        string      `json:"type"`
	Content     string      `json:"content"`
	Name        string      `json:"name,omitempty"`
	TTL         int         `json:"ttl,omitempty"`
	Priority    interface{} `json:"priority,omitempty"`
}

// dnsheCreateRecordResp 创建DNS记录响应
type dnsheCreateRecordResp struct {
	Success  bool        `json:"success"`
	Message  string      `json:"message"`
	RecordID interface{} `json:"record_id"`
	Error    string      `json:"error,omitempty"`
}

// dnsheUpdateRecordReq 更新DNS记录请求参数
type dnsheUpdateRecordReq struct {
	RecordID  int    `json:"record_id"`
	Content   string `json:"content,omitempty"`
	TTL       int    `json:"ttl,omitempty"`
	Priority  int    `json:"priority,omitempty"`
}

// dnsheUpdateRecordResp 更新DNS记录响应
type dnsheUpdateRecordResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// dnsheDeleteRecordReq 删除DNS记录请求参数
type dnsheDeleteRecordReq struct {
	RecordID int `json:"record_id"`
}

// dnsheDeleteRecordResp 删除DNS记录响应
type dnsheDeleteRecordResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// --- API密钥管理相关结构体 ---
// dnsheAPIKey API密钥信息
type dnsheAPIKey struct {
	ID           int    `json:"id"`
	KeyName      string `json:"key_name"`
	APIKey       string `json:"api_key"`
	Status       string `json:"status"`
	RequestCount int    `json:"request_count"`
	LastUsedAt   string `json:"last_used_at"`
	CreatedAt    string `json:"created_at"`
}

// dnsheListAPIKeysResp 列出API密钥响应
type dnsheListAPIKeysResp struct {
	Success bool           `json:"success"`
	Count   int            `json:"count"`
	Keys    []dnsheAPIKey  `json:"keys"`
	Error   string         `json:"error,omitempty"`
}

// dnsheCreateAPIKeyReq 创建API密钥请求参数
type dnsheCreateAPIKeyReq struct {
	KeyName     string `json:"key_name"`
	IPWhitelist string `json:"ip_whitelist,omitempty"`
}

// dnsheCreateAPIKeyResp 创建API密钥响应
type dnsheCreateAPIKeyResp struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
	Warning   string `json:"warning"`
	Error     string `json:"error,omitempty"`
}

// dnsheDeleteAPIKeyReq 删除API密钥请求参数
type dnsheDeleteAPIKeyReq struct {
	KeyID int `json:"key_id"`
}

// dnsheDeleteAPIKeyResp 删除API密钥响应
type dnsheDeleteAPIKeyResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// dnsheRegenAPIKeyReq 重新生成API密钥请求参数
type dnsheRegenAPIKeyReq struct {
	KeyID int `json:"key_id"`
}

// dnsheRegenAPIKeyResp 重新生成API密钥响应
type dnsheRegenAPIKeyResp struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
	Warning   string `json:"warning"`
	Error     string `json:"error,omitempty"`
}

// --- 配额查询相关结构体 ---
// dnsheQuota 配额信息
type dnsheQuota struct {
	Used        int `json:"used"`
	Base        int `json:"base"`
	InviteBonus int `json:"invite_bonus"`
	Total       int `json:"total"`
	Available   int `json:"available"`
}

// dnsheQuotaResp 配额查询响应
type dnsheQuotaResp struct {
	Success bool        `json:"success"`
	Quota   dnsheQuota  `json:"quota"`
	Error   string      `json:"error,omitempty"`
}

// --- 通用错误响应结构体 ---
type dnsheErrorResp struct {
	Error     string `json:"error"`
	Limit     int    `json:"limit,omitempty"`
	Remaining int    `json:"remaining,omitempty"`
	ResetAt   string `json:"reset_at,omitempty"`
}