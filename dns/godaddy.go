package dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

// godaddyRecord is one DNS record in the GoDaddy v3 API. RecordID is
// server-assigned; it is only present on records read back from the API
// and is omitted when creating a record.
type godaddyRecord struct {
	RecordID string `json:"recordId,omitempty"`
	Data     string `json:"data"`
	Name     string `json:"name"`
	TTL      int    `json:"ttl"`
	Type     string `json:"type"`
}

type godaddyRecords []godaddyRecord

// godaddyListResponse is the v3 list shape: a paginated object whose
// records live under "items". (Some list calls return a bare array;
// getRecords accepts both.)
type godaddyListResponse struct {
	Items godaddyRecords `json:"items"`
}

// godaddyV3Base is the GoDaddy Domains API v3 zone endpoint.
// v1 (api.godaddy.com/v1) and its sso-key auth are deprecated by GoDaddy;
// v3 authenticates with a Personal Access Token (Bearer).
const godaddyV3Base = "https://api.godaddy.com/v3/domains/zones"

type GoDaddyDNS struct {
	dns      config.DNS
	domains  config.Domains
	ttl      int
	header   http.Header
	client   *http.Client
	lastIpv4 string
	lastIpv6 string
}

func (g *GoDaddyDNS) Init(dnsConf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache) {
	g.domains.Ipv4Cache = ipv4cache
	g.domains.Ipv6Cache = ipv6cache
	g.lastIpv4 = ipv4cache.Addr
	g.lastIpv6 = ipv6cache.Addr

	g.dns = dnsConf.DNS
	g.domains.GetNewIp(dnsConf)
	g.ttl = 600
	if val, err := strconv.Atoi(dnsConf.TTL); err == nil {
		g.ttl = val
	}
	// v3 takes a Personal Access Token in the Secret field. The v1 Key
	// (ID field) is no longer used. A token that does not start with
	// "gd_" is almost certainly a v1 key, so warn instead of failing
	// with a bare 401 from the API.
	if g.dns.Secret != "" && !strings.HasPrefix(g.dns.Secret, "gd_") {
		util.Log("GoDaddy Secret 看起来不是 v3 Personal Access Token (应以 gd_ 开头)。请在 https://developer.godaddy.com/keys 创建 Token 并填入 Secret 字段")
	}
	g.header = map[string][]string{
		"Authorization": {"Bearer " + g.dns.Secret},
		"Content-Type":  {"application/json"},
	}

	g.client = dnsConf.GetHTTPClient()
}

func (g *GoDaddyDNS) updateDomainRecord(recordType string, ipAddr string, domains []*config.Domain) {
	if ipAddr == "" {
		return
	}

	// 防止多次发送Webhook通知
	if recordType == "A" {
		if g.lastIpv4 == ipAddr {
			util.Log("你的IPv4未变化, 未触发 %s 请求", "godaddy")
			return
		}
	} else {
		if g.lastIpv6 == ipAddr {
			util.Log("你的IPv6未变化, 未触发 %s 请求", "godaddy")
			return
		}
	}

	for _, domain := range domains {
		err := g.replaceRecord(recordType, ipAddr, domain)
		if err == nil {
			util.Log("更新域名解析 %s 成功! IP: %s", domain, ipAddr)
			domain.UpdateStatus = config.UpdatedSuccess
		} else {
			util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, err)
			domain.UpdateStatus = config.UpdatedFailed
		}
	}
}

func (g *GoDaddyDNS) AddUpdateDomainRecords() config.Domains {
	if ipv4Addr, ipv4Domains := g.domains.GetNewIpResult("A"); ipv4Addr != "" {
		g.updateDomainRecord("A", ipv4Addr, ipv4Domains)
	}
	if ipv6Addr, ipv6Domains := g.domains.GetNewIpResult("AAAA"); ipv6Addr != "" {
		g.updateDomainRecord("AAAA", ipv6Addr, ipv6Domains)
	}
	return g.domains
}

// replaceRecord brings one record to ipAddr via the v3 API. It first reads
// the current records and only writes when they differ, so a restart (or a
// record changed elsewhere) does not cause a spurious write.
//
// A v3 PUT is never used here: PUT replaces a whole record set, and a PUT
// to a wrong or missing endpoint has wiped zones before
// (godaddy/cli#136). Instead we DELETE the old records by their
// server-assigned recordId and POST the new one, one record at a time.
func (g *GoDaddyDNS) replaceRecord(recordType string, ipAddr string, domain *config.Domain) error {
	current, err := g.getRecords(recordType, domain)
	if err != nil {
		return err
	}
	if len(current) == 1 && current[0].Data == ipAddr {
		util.Log("域名解析 %s 已是最新 IP: %s, 无需更新", domain, ipAddr)
		return nil
	}
	for _, rec := range current {
		if rec.RecordID == "" {
			return fmt.Errorf("GoDaddy record %s has no recordId, cannot delete it", rec.Name)
		}
		if err := g.deleteRecord(domain, rec.RecordID); err != nil {
			return err
		}
	}
	return g.createRecord(recordType, ipAddr, domain)
}

// listPath returns the v3 collection URL for one zone.
func (g *GoDaddyDNS) listPath(domain *config.Domain) string {
	return fmt.Sprintf("%s/%s/dns-records", godaddyV3Base, domain.DomainName)
}

// do sends one v3 request and returns the body and HTTP status code.
// Transport errors come back as err; HTTP error statuses do not, so the
// caller can treat a 404 as "not found" instead of a failure.
func (g *GoDaddyDNS) do(req *http.Request) ([]byte, int, error) {
	req.Header = g.header
	resp, err := g.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024000))
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return body, resp.StatusCode, nil
}

func (g *GoDaddyDNS) getRecords(recordType string, domain *config.Domain) (godaddyRecords, error) {
	// v3 lists by query parameters, not by a /{type}/{name} path, and
	// returns a paginated object {"items": [...]}.
	u := g.listPath(domain) + "?type=" + url.QueryEscape(recordType) +
		"&name=" + url.QueryEscape(domain.GetSubDomain()) +
		"&page=1&pageSize=100"
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	body, status, err := g.do(req)
	if err != nil {
		return nil, err
	}
	// 404 means the zone does not exist yet: no current record.
	if status == http.StatusNotFound {
		return godaddyRecords{}, nil
	}
	if status >= 300 {
		return nil, fmt.Errorf("%s", util.LogStr("返回内容: %s ,返回状态码: %d", string(body), status))
	}
	if len(body) == 0 {
		return godaddyRecords{}, nil
	}
	// Preferred shape: {"items": [...]}.
	var listed godaddyListResponse
	if err := json.Unmarshal(body, &listed); err == nil && listed.Items != nil {
		return listed.Items, nil
	}
	// Fallback: a bare array.
	var records godaddyRecords
	if err := json.Unmarshal(body, &records); err != nil {
		return nil, err
	}
	return records, nil
}

func (g *GoDaddyDNS) deleteRecord(domain *config.Domain, recordID string) error {
	// v3 deletes one record by its server-assigned recordId.
	req, err := http.NewRequest(http.MethodDelete,
		fmt.Sprintf("%s/%s", g.listPath(domain), recordID), nil)
	if err != nil {
		return err
	}
	body, status, err := g.do(req)
	if err != nil {
		return err
	}
	// 404 means there was nothing to delete.
	if status == http.StatusNotFound {
		return nil
	}
	if status >= 300 {
		return fmt.Errorf("%s", util.LogStr("返回内容: %s ,返回状态码: %d", string(body), status))
	}
	return nil
}

func (g *GoDaddyDNS) createRecord(recordType string, ipAddr string, domain *config.Domain) error {
	// v3 POST creates ONE record and takes a single object, not an array.
	record := godaddyRecord{
		Data: ipAddr,
		Name: domain.GetSubDomain(),
		TTL:  g.ttl,
		Type: recordType,
	}
	buffer, err := json.Marshal(record)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, g.listPath(domain), bytes.NewBuffer(buffer))
	if err != nil {
		return err
	}
	body, status, err := g.do(req)
	if err != nil {
		return err
	}
	if status >= 300 {
		return fmt.Errorf("%s", util.LogStr("返回内容: %s ,返回状态码: %d", string(body), status))
	}
	return nil
}
