package dns

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

// https://www.todaynic.com/docApi/
// Nowcn nowcn DNS实现
type Nowcn struct {
	DNS     config.DNS
	Domains config.Domains
	TTL     string
}

// NowcnRecord DNS记录结构
type NowcnRecord struct {
	ID     int `json:"id"`
	Domain string
	Host   string
	Type   string
	Value  string
	State  int
	// Name    string
	// Enabled string
}

// NowcnRecordListResp 记录列表响应
type NowcnRecordListResp struct {
	NowcnBaseResult
	Data []NowcnRecord
}

// NowcnStatus API响应状态
type NowcnBaseResult struct {
	RequestId string `json:"RequestId"`
	Id        int    `json:"Id"`
	Error     string `json:"error"`
}

// Init 初始化
func (nowcn *Nowcn) Init(dnsConf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache) {
	nowcn.Domains.Ipv4Cache = ipv4cache
	nowcn.Domains.Ipv6Cache = ipv6cache
	nowcn.DNS = dnsConf.DNS
	nowcn.Domains.GetNewIp(dnsConf)
	if dnsConf.TTL == "" {
		// 默认600s
		nowcn.TTL = "600"
	} else {
		nowcn.TTL = dnsConf.TTL
	}
}

// AddUpdateDomainRecords 添加或更新IPv4/IPv6记录
func (nowcn *Nowcn) AddUpdateDomainRecords() config.Domains {
	nowcn.addUpdateDomainRecords("A")
	nowcn.addUpdateDomainRecords("AAAA")
	return nowcn.Domains
}

func (nowcn *Nowcn) addUpdateDomainRecords(recordType string) {
	ipAddr, domains := nowcn.Domains.GetNewIpResult(recordType)

	if ipAddr == "" {
		return
	}

	for _, domain := range domains {
		result, err := nowcn.getRecordList(domain, recordType)
		if err != nil {
			util.Log("查询域名信息发生异常! %s", err)
			domain.UpdateStatus = config.UpdatedFailed
			return
		}

		if len(result.Data) > 0 {
			// 默认第一个
			recordSelected := result.Data[0]
			params := domain.GetCustomParams()
			if params.Has("Id") {
				for i := 0; i < len(result.Data); i++ {
					if strconv.Itoa(result.Data[i].ID) == params.Get("Id") {
						recordSelected = result.Data[i]
					}
				}
			}
			// 更新
			nowcn.modify(recordSelected, domain, recordType, ipAddr)
		} else {
			// 新增
			nowcn.create(domain, recordType, ipAddr)
		}
	}
}

// create 创建DNS记录
func (nowcn *Nowcn) create(domain *config.Domain, recordType string, ipAddr string) {
	param := map[string]string{
		"Domain": domain.DomainName,
		"Host":   domain.GetSubDomain(),
		"Type":   recordType,
		"Value":  ipAddr,
		"Ttl":    nowcn.TTL,
	}
	res, err := nowcn.request("/api/Dns/AddDomainRecord", param, "GET")
	if err != nil {
		util.Log("新增域名解析 %s 失败! 异常信息: %s", domain, err.Error())
		domain.UpdateStatus = config.UpdatedFailed
	}
	var result NowcnBaseResult
	err = json.Unmarshal(res, &result)
	if err != nil {
		util.Log("新增域名解析 %s 失败! 异常信息: %s", domain, err.Error())
		domain.UpdateStatus = config.UpdatedFailed
	}
	if result.Error != "" {
		util.Log("新增域名解析 %s 失败! 异常信息: %s", domain, result.Error)
		domain.UpdateStatus = config.UpdatedFailed
	} else {
		domain.UpdateStatus = config.UpdatedSuccess
	}
}

// modify 修改DNS记录
func (nowcn *Nowcn) modify(record NowcnRecord, domain *config.Domain, recordType string, ipAddr string) {
	// 相同不修改
	if record.Value == ipAddr {
		util.Log("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
		return
	}
	param := map[string]string{
		"Id":     strconv.Itoa(record.ID),
		"Domain": domain.DomainName,
		"Host":   domain.GetSubDomain(),
		"Type":   recordType,
		"Value":  ipAddr,
		"Ttl":    nowcn.TTL,
	}
	res, err := nowcn.request("/api/Dns/UpdateDomainRecord", param, "GET")
	if err != nil {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, err.Error())
		domain.UpdateStatus = config.UpdatedFailed
	}
	var result NowcnBaseResult
	err = json.Unmarshal(res, &result)
	if err != nil {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, err.Error())
		domain.UpdateStatus = config.UpdatedFailed
	}
	if result.Error != "" {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, result.Error)
		domain.UpdateStatus = config.UpdatedFailed
	} else {
		util.Log("更新域名解析 %s 成功! IP: %s", domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	}
}

// getRecordList 获取域名记录列表
func (nowcn *Nowcn) getRecordList(domain *config.Domain, typ string) (result NowcnRecordListResp, err error) {
	param := map[string]string{
		"Domain": domain.DomainName,
		"Type":   typ,
		"Host":   domain.GetSubDomain(),
	}
	res, err := nowcn.request("/api/Dns/DescribeRecordIndex", param, "GET")
	err = json.Unmarshal(res, &result)
	return
}

func (t *Nowcn) sign(params map[string]string, method string) (string, error) {
	// 添加公共参数
	params["AccessKeyID"] = t.DNS.ID
	params["SignatureMethod"] = "HMAC-SHA1"
	params["SignatureNonce"] = fmt.Sprintf("%d", time.Now().UnixNano())
	params["Timestamp"] = time.Now().UTC().Format("2006-01-02T15:04:05Z")

	// 1. 排序参数(按首字母顺序)
	var keys []string
	for k := range params {
		if k != "Signature" { // 排除Signature参数
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// 2. 构造规范化请求字符串
	var canonicalizedQuery []string
	for _, k := range keys {
		// URL编码参数名和参数值
		encodedKey := util.PercentEncode(k)
		encodedValue := util.PercentEncode(params[k])
		canonicalizedQuery = append(canonicalizedQuery, encodedKey+"="+encodedValue)
	}
	canonicalizedQueryString := strings.Join(canonicalizedQuery, "&")

	// 3. 构造待签名字符串
	stringToSign := method + "&" + util.PercentEncode("/") + "&" + util.PercentEncode(canonicalizedQueryString)

	// 4. 计算HMAC-SHA1签名
	key := t.DNS.Secret + "&"
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// 5. 添加签名到参数中
	params["Signature"] = signature

	// 6. 重新构造最终的查询字符串(包含签名)
	keys = append(keys, "Signature")
	sort.Strings(keys)
	var finalQuery []string
	for _, k := range keys {
		encodedKey := util.PercentEncode(k)
		encodedValue := util.PercentEncode(params[k])
		finalQuery = append(finalQuery, encodedKey+"="+encodedValue)
	}

	return strings.Join(finalQuery, "&"), nil
}

func (t *Nowcn) request(apiPath string, params map[string]string, method string) ([]byte, error) {
	// 生成签名
	queryString, err := t.sign(params, method)
	if err != nil {
		return nil, fmt.Errorf("生成签名失败: %v", err)
	}

	// 构造完整URL
	baseURL := "https://api.now.cn"
	fullURL := baseURL + apiPath + "?" + queryString

	// 创建HTTP请求
	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Accept", "application/json")

	// 发送请求
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
