package dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

// 获取配置目录（兼容Windows/Linux/macOS）
func getConfigDir() string {
	// 优先读取环境变量
	if dir := os.Getenv("DDNS_GO_CONFIG"); dir != "" {
		return dir
	}

	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./" // 降级到当前目录
	}

	// 根据系统返回默认配置目录
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(homeDir, "AppData", "Roaming", "ddns-go")
	case "darwin":
		return filepath.Join(homeDir, "Library", "Application Support", "ddns-go")
	default: // Linux/FreeBSD等
		return filepath.Join(homeDir, ".ddns_go")
	}
}

// 固定DNSHE API基础地址
const dnsheAPIBase = "https://api005.dnshe.com/index.php?m=domain_hub"

// DNSHE DNSHE服务商接口实现
type DNSHE struct {
	DNS     config.DNS
	Domains config.Domains
	TTL     int
	// 移除apiLogger字段，删除API日志器相关定义
}

// Init 初始化DNSHE客户端
func (d *DNSHE) Init(dnsConf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache) {
	d.Domains.Ipv4Cache = ipv4cache
	d.Domains.Ipv6Cache = ipv6cache
	d.DNS = dnsConf.DNS
	d.Domains.GetNewIp(dnsConf)

	// 初始化TTL，默认600
	if dnsConf.TTL == "" {
		d.TTL = 600
	} else {
		ttl, err := strconv.Atoi(dnsConf.TTL)
		d.TTL = 600
		if err == nil {
			d.TTL = ttl
		}
	}

	// 移除API日志器初始化代码，不再创建日志器实例
}

// AddUpdateDomainRecords 新增或更新域名解析记录
func (d *DNSHE) AddUpdateDomainRecords() config.Domains {
	d.addUpdateDomainRecords("A")
	d.addUpdateDomainRecords("AAAA")
	return d.Domains
}

// addUpdateDomainRecords 处理指定类型的域名解析记录
func (d *DNSHE) addUpdateDomainRecords(recordType string) {
	ipAddr, domains := d.Domains.GetNewIpResult(recordType)
	if ipAddr == "" {
		return
	}

	for _, domain := range domains {
		fullDomain := domain.GetFullDomain()
		rootDomain, firstPrefix, multiPrefix := splitDomainToMultiLevels(fullDomain)
		if rootDomain == "" {
			util.Log("域名格式非法: %s", fullDomain)
			domain.UpdateStatus = config.UpdatedFailed
			continue
		}

		// 查询一级子域（不再自动注册）
		firstSubDomain := fmt.Sprintf("%s.%s", firstPrefix, rootDomain)
		if firstPrefix == "" {
			firstSubDomain = rootDomain
		}
		subID, err := d.getFirstSubdomain(firstPrefix, rootDomain)
		if err != nil || subID <= 0 {
			util.Log("一级子域%s不存在或查询失败: %s", firstSubDomain, err)
			domain.UpdateStatus = config.UpdatedFailed
			continue
		}

		// 确定记录名称（多级前缀）
		recordName := multiPrefix
		targetFullName := fullDomain
		if multiPrefix == "" {
			targetFullName = firstSubDomain
			recordName = ""
		}

		// 查询现有DNS记录
		existRec, err := d.findRecordByFullName(subID, targetFullName, recordType)
		if err != nil {
			util.Log("查询DNS记录异常: %s", err)
			domain.UpdateStatus = config.UpdatedFailed
			continue
		}

		// 更新或创建记录
		if existRec != nil {
			if existRec.Content == ipAddr {
				util.Log("IP未变化: %s -> %s", ipAddr, fullDomain)
				continue
			}
			if err := d.updateRecord(existRec.ID, ipAddr); err != nil {
				util.Log("更新解析%s失败: %s", fullDomain, err)
				domain.UpdateStatus = config.UpdatedFailed
				continue
			}
			util.Log("更新解析%s成功: %s", fullDomain, ipAddr)
			domain.UpdateStatus = config.UpdatedSuccess
		} else {
			if err := d.createRecordWithMultiPrefix(subID, recordName, recordType, ipAddr); err != nil {
				util.Log("新增解析%s失败: %s", fullDomain, err)
				domain.UpdateStatus = config.UpdatedFailed
				continue
			}
			util.Log("新增解析%s成功: %s", fullDomain, ipAddr)
			domain.UpdateStatus = config.UpdatedSuccess
		}
	}
}

// splitDomainToMultiLevels 拆分域名为根域、一级前缀、多级前缀
// 规则: 根域为最后两段, 一级前缀为倒数第三段, 多级前缀为前面所有段
func splitDomainToMultiLevels(fullDomain string) (rootDomain, firstPrefix, multiPrefix string) {
	fullDomain = strings.TrimSuffix(fullDomain, ".")
	parts := strings.Split(fullDomain, ".")
	if len(parts) < 2 {
		return "", "", ""
	}

	rootDomain = strings.Join(parts[len(parts)-2:], ".")
	if len(parts) == 2 {
		return rootDomain, "", ""
	}
	if len(parts) == 3 {
		return rootDomain, parts[0], ""
	}

	firstPrefix = parts[len(parts)-3]
	multiPrefix = strings.Join(parts[:len(parts)-3], ".")
	return rootDomain, firstPrefix, multiPrefix
}

// convertToInt 将interface{}类型值转为int，兼容string/int/float64类型
func convertToInt(v interface{}) (int, error) {
	switch val := v.(type) {
	case int:
		return val, nil
	case string:
		return strconv.Atoi(val)
	case float64:
		return int(val), nil
	default:
		return 0, fmt.Errorf("不支持的类型转换: %T", v)
	}
}

// getFirstSubdomain 仅查询一级子域，不自动注册
func (d *DNSHE) getFirstSubdomain(prefix, root string) (int, error) {
	// 1. 查询现有子域名列表
	var listResp dnsheListSubdomainsResp
	u := fmt.Sprintf("%s&endpoint=subdomains&action=list", dnsheAPIBase)
	if err := d.request("GET", u, nil, &listResp); err != nil {
		return 0, fmt.Errorf("查询子域名列表失败: %s", err)
	}

	// 2. 匹配目标一级子域
	targetFullDomain := fmt.Sprintf("%s.%s", prefix, root)
	if prefix == "" {
		targetFullDomain = root
	}
	if listResp.Success {
		for _, s := range listResp.Subdomains {
			if strings.EqualFold(s.FullDomain, targetFullDomain) {
				return s.ID, nil
			}
		}
	}

	// 3. 未找到子域时返回错误
	return 0, fmt.Errorf("子域%s不存在，请先手动注册", targetFullDomain)
}

// findRecordByFullName 按完整域名查询DNS记录
func (d *DNSHE) findRecordByFullName(subID int, fullName, recordType string) (*dnsheRecord, error) {
	var resp dnsheListRecordsResp
	qs := url.Values{}
	qs.Set("subdomain_id", strconv.Itoa(subID))
	u := fmt.Sprintf("%s&endpoint=dns_records&action=list&%s", dnsheAPIBase, qs.Encode())

	if err := d.request("GET", u, nil, &resp); err != nil {
		return nil, fmt.Errorf("请求失败: %s", err)
	}
	if !resp.Success {
		errMsg := "查询无结果"
		if resp.Error != "" {
			errMsg = resp.Error
		}
		return nil, fmt.Errorf("查询DNS记录异常: %s", errMsg)
	}

	for _, r := range resp.Records {
		if strings.EqualFold(r.Type, recordType) && strings.EqualFold(r.Name, fullName) {
			return &r, nil
		}
	}
	return nil, nil
}

// createRecordWithMultiPrefix 创建带多级前缀的DNS记录
func (d *DNSHE) createRecordWithMultiPrefix(subID int, multiPrefix, recordType, ip string) error {
	req := dnsheCreateRecordReq{
		SubdomainID: subID,
		Type:        recordType,
		Content:     ip,
		Name:        multiPrefix,
		TTL:         d.TTL,
	}
	var resp dnsheCreateRecordResp
	u := fmt.Sprintf("%s&endpoint=dns_records&action=create", dnsheAPIBase)

	if err := d.request("POST", u, req, &resp); err != nil {
		return fmt.Errorf("请求失败: %s", err)
	}
	if !resp.Success {
		errMsg := "创建无响应"
		if resp.Error != "" {
			errMsg = resp.Error
		}
		return fmt.Errorf("创建DNS记录异常: %s", errMsg)
	}

	// 打印record_id日志，不显示原始响应
	switch v := resp.RecordID.(type) {
	case int:
		util.Log("创建记录成功，record_id (int): %d", v)
	case string:
		util.Log("创建记录成功，record_id (string): %s", v)
	default:
		util.Log("创建记录成功，record_id (未知类型): %v", v)
	}
	return nil
}

// updateRecord 更新DNS记录
func (d *DNSHE) updateRecord(recordID int, ip string) error {
	req := dnsheUpdateRecordReq{RecordID: recordID, Content: ip, TTL: d.TTL}
	var resp dnsheUpdateRecordResp
	u := fmt.Sprintf("%s&endpoint=dns_records&action=update", dnsheAPIBase)

	if err := d.request("POST", u, req, &resp); err != nil {
		return fmt.Errorf("请求失败: %s", err)
	}
	if !resp.Success {
		errMsg := "更新无响应"
		if resp.Error != "" {
			errMsg = resp.Error
		}
		return fmt.Errorf("更新DNS记录异常: %s", errMsg)
	}
	return nil
}

// request 通用HTTP请求方法，处理API通信
func (d *DNSHE) request(method, urlStr string, data interface{}, result interface{}) (err error) {
	var reqBody bytes.Buffer
	if method != "GET" && data != nil {
		jsonBytes, marshalErr := json.Marshal(data)
		if marshalErr != nil {
			return fmt.Errorf("序列化失败: %s", marshalErr)
		}
		reqBody = *bytes.NewBuffer(jsonBytes)
	}

	// 创建HTTP请求
	req, err := http.NewRequest(method, urlStr, &reqBody)
	if err != nil {
		return fmt.Errorf("创建请求失败: %s", err)
	}
	req.Header.Set("X-API-Key", d.DNS.ID)
	req.Header.Set("X-API-Secret", d.DNS.Secret)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := util.CreateHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %s", err)
	}
	defer resp.Body.Close()

	// 读取原始响应（仅用于反序列化，不再处理日志）
	rawBody, _ := io.ReadAll(resp.Body)



	// 反序列化响应
	bodyReader := bytes.NewReader(rawBody)
	if err = json.NewDecoder(bodyReader).Decode(result); err != nil {
		util.Log("JSON反序列化失败，但API可能已执行成功")
		return nil
	}
	return nil
}

// --- 适配原有接口 ---
func (d *DNSHE) findRecordByType(subID int, domain *config.Domain, recordType string) (*dnsheRecord, error) {
	return d.findRecordByFullName(subID, domain.GetFullDomain(), recordType)
}

func (d *DNSHE) createRecord(subID int, recordType, ip string) error {
	return d.createRecordWithMultiPrefix(subID, "", recordType, ip)
}

// --- 扩展功能接口 ---
// GetSubdomainDetail 获取子域名详情
func (d *DNSHE) GetSubdomainDetail(subdomainID int) (*dnsheSubdomainDetailResp, error) {
	var resp dnsheSubdomainDetailResp
	qs := url.Values{}
	qs.Set("subdomain_id", strconv.Itoa(subdomainID))
	u := fmt.Sprintf("%s&endpoint=subdomains&action=get&%s", dnsheAPIBase, qs.Encode())
	if err := d.request("GET", u, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetQuota 查询账户配额
func (d *DNSHE) GetQuota() (*dnsheQuota, error) {
	var resp dnsheQuotaResp
	u := fmt.Sprintf("%s&endpoint=quota", dnsheAPIBase)
	if err := d.request("GET", u, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Quota, nil
}
