package dns

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
	"golang.org/x/net/http/httpguts"
)

type Callback struct {
	DNS        config.DNS
	Domains    config.Domains
	TTL        string
	lastIpv4   string
	lastIpv6   string
	httpClient *http.Client
	ipv4Enable bool
	ipv6Enable bool
}

// Init 初始化
func (cb *Callback) Init(dnsConf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache) {
	cb.Domains.Ipv4Cache = ipv4cache
	cb.Domains.Ipv6Cache = ipv6cache
	cb.lastIpv4 = ipv4cache.Addr
	cb.lastIpv6 = ipv6cache.Addr
	cb.ipv4Enable = dnsConf.Ipv4.Enable
	cb.ipv6Enable = dnsConf.Ipv6.Enable

	cb.DNS = dnsConf.DNS
	cb.Domains.GetNewIp(dnsConf)
	if dnsConf.TTL == "" {
		// 默认600
		cb.TTL = "600"
	} else {
		cb.TTL = dnsConf.TTL
	}
	cb.httpClient = dnsConf.GetHTTPClient()
}

// AddUpdateDomainRecords 添加或更新IPv4/IPv6记录
func (cb *Callback) AddUpdateDomainRecords() config.Domains {
	if cb.ipv4Enable {
		cb.addUpdateDomainRecords("A")
	}
	if cb.ipv6Enable {
		cb.addUpdateDomainRecords("AAAA")
	}
	return cb.Domains
}

func (cb *Callback) addUpdateDomainRecords(recordType string) {
	ipAddr, domains := cb.Domains.GetNewIpResult(recordType)

	if ipAddr == "" {
		return
	}

	// 防止多次发送Webhook通知
	if recordType == "A" {
		if cb.lastIpv4 == ipAddr {
			util.Log("你的IPv4未变化, 未触发 %s 请求", "Callback")
			return
		}
	} else {
		if cb.lastIpv6 == ipAddr {
			util.Log("你的IPv6未变化, 未触发 %s 请求", "Callback")
			return
		}
	}

	for _, domain := range domains {
		method := "GET"
		postPara := ""
		contentType := "application/x-www-form-urlencoded"
		if cb.DNS.Secret != "" {
			method = "POST"
			postPara = cb.replacePara(cb.DNS.Secret, ipAddr, domain, recordType, cb.TTL)
			if json.Valid([]byte(postPara)) {
				contentType = "application/json"
			}
		}
		requestURL := cb.replacePara(cb.DNS.ID, ipAddr, domain, recordType, cb.TTL)
		u, err := url.Parse(requestURL)
		if err != nil {
			util.Log("Callback的URL不正确")
			return
		}
		req, err := http.NewRequest(method, u.String(), strings.NewReader(postPara))
		if err != nil {
			util.Log("异常信息: %v", err)
			domain.UpdateStatus = config.UpdatedFailed
			return
		}
		// 添加自定义请求头（支持与 URL/RequestBody 相同的变量替换），为空则仅保留默认请求头
		headers := extractCallbackHeaders(cb.replacePara(cb.DNS.Headers, ipAddr, domain, recordType, cb.TTL))
		for key, value := range headers {
			req.Header.Add(key, value)
		}
		// 仅在用户未自定义 content-type 时设置默认值，避免出现多个 Content-Type 头
		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("content-type", contentType)
		}

		resp, err := cb.httpClient.Do(req)
		body, err := util.GetHTTPResponseOrg(resp, err)
		if err == nil {
			util.Log("Callback调用成功, 域名: %s, IP: %s, 返回数据: %s", domain, ipAddr, string(body))
			domain.UpdateStatus = config.UpdatedSuccess
		} else {
			util.Log("Callback调用失败, 异常信息: %v", err)
			domain.UpdateStatus = config.UpdatedFailed
		}
	}
}

// extractCallbackHeaders 将"每行 Header: value"格式的字符串转换为 map。
// 与 webhook 的 extractHeaders 实现隔离，避免影响 webhook 的日志语义。
func extractCallbackHeaders(s string) map[string]string {
	lines := util.SplitLines(s)
	headers := make(map[string]string, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			util.Log("Callback Header格式不正确: %s", line)
			continue
		}
		k, v := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		// 校验 header 名值合法性，避免非法字段导致整个请求被 net/http 拒绝
		if !httpguts.ValidHeaderFieldName(k) || !httpguts.ValidHeaderFieldValue(v) {
			util.Log("Callback Header格式不正确: %s", line)
			continue
		}
		headers[k] = v
	}
	return headers
}

// replacePara 替换参数
func (cb *Callback) replacePara(orgPara, ipAddr string, domain *config.Domain, recordType string, ttl string) string {
	// params 使用 map 以便添加更多参数
	params := map[string]string{
		"ip":         ipAddr,
		"domain":     domain.String(),
		"recordType": recordType,
		"ttl":        ttl,
		"ipv4Addr":   cb.Domains.Ipv4Addr,
		"ipv6Addr":   cb.Domains.Ipv6Addr,
		"timestamp":  strconv.FormatInt(time.Now().UTC().Unix(), 10),
	}

	// 也替换域名的自定义参数
	for k, v := range domain.GetCustomParams() {
		if len(v) == 1 {
			params[k] = v[0]
		}
	}

	// 将 map 转换为 [NewReplacer] 所需的参数
	// map 中的每个元素占用 2 个位置（kv），因此需要预留 2 倍的空间
	oldnew := make([]string, 0, len(params)*2)
	for k, v := range params {
		k = fmt.Sprintf("#{%s}", k)
		oldnew = append(oldnew, k, v)
	}

	return strings.NewReplacer(oldnew...).Replace(orgPara)
}
