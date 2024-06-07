package dns

import (
	"io"
	"net/http"
	"strings"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

const (
	nameCheapEndpoint string = "https://dynamicdns.park-your-domain.com/update?host=#{host}&domain=#{domain}&password=#{password}&ip=#{ip}"
)

// NameCheap Domain
type NameCheap struct {
	DNS      config.DNS
	Domains  config.Domains
	lastIpv4 string
	lastIpv6 string
}

// NameCheap 修改域名解析结果
type NameCheapResp struct {
	Status string
	Errors []string
}

// Init 初始化
func (nc *NameCheap) Init(dnsConf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache) {
	nc.Domains.Ipv4Cache = ipv4cache
	nc.Domains.Ipv6Cache = ipv6cache
	nc.lastIpv4 = ipv4cache.Addr
	nc.lastIpv6 = ipv6cache.Addr

	nc.DNS = dnsConf.DNS
	nc.Domains.GetNewIp(dnsConf)
}

// AddUpdateDomainRecords 添加或更新IPv4/IPv6记录
func (nc *NameCheap) AddUpdateDomainRecords() config.Domains {
	nc.addUpdateDomainRecords("A")
	nc.addUpdateDomainRecords("AAAA")
	return nc.Domains
}

func (nc *NameCheap) addUpdateDomainRecords(recordType string) {
	ipAddr, domains := nc.Domains.GetNewIpResult(recordType)

	if ipAddr == "" {
		return
	}

	// 防止多次发送Webhook通知
	if recordType == "A" {
		if nc.lastIpv4 == ipAddr {
			util.Log("你的IPv4未变化, 未触发 %s 请求", "NameCheap")
			return
		}
	} else {
		// https://www.namecheap.com/support/knowledgebase/article.aspx/29/11/how-to-dynamically-update-the-hosts-ip-with-an-http-request/
		util.Log("Namecheap 不支持更新 IPv6")
		return
	}

	for _, domain := range domains {
		nc.modify(domain, ipAddr)
	}
}

// 修改
func (nc *NameCheap) modify(domain *config.Domain, ipAddr string) {
	var result NameCheapResp
	err := nc.request(&result, ipAddr, domain)

	if err != nil {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, err)
		domain.UpdateStatus = config.UpdatedFailed
		return
	}

	switch result.Status {
	case "Success":
		util.Log("更新域名解析 %s 成功! IP: %s", domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	default:
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, result.Status)
		domain.UpdateStatus = config.UpdatedFailed
	}
}

// request 统一请求接口
func (nc *NameCheap) request(result *NameCheapResp, ipAddr string, domain *config.Domain) (err error) {
	url := strings.NewReplacer(
		"#{host}", domain.GetSubDomain(),
		"#{domain}", domain.DomainName,
		"#{password}", nc.DNS.Secret,
		"#{ip}", ipAddr,
	).Replace(nameCheapEndpoint)

	req, err := http.NewRequest(
		http.MethodGet,
		url,
		http.NoBody,
	)

	if err != nil {
		return
	}

	client := util.CreateHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	status := string(data)

	if strings.Contains(status, "<ErrCount>0</ErrCount>") {
		result.Status = "Success"
	} else {
		result.Status = status
	}

	return
}
