package dns

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

type Callback struct {
	DNS      config.DNS
	Domains  config.Domains
	TTL      string
	lastIpv4 string
	lastIpv6 string
}

// Init 初始化
func (cb *Callback) Init(dnsConf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache) {
	cb.Domains.Ipv4Cache = ipv4cache
	cb.Domains.Ipv6Cache = ipv6cache
	cb.lastIpv4 = ipv4cache.Addr
	cb.lastIpv6 = ipv6cache.Addr

	cb.DNS = dnsConf.DNS
	cb.Domains.GetNewIp(dnsConf)
	if dnsConf.TTL == "" {
		// 默认600
		cb.TTL = "600"
	} else {
		cb.TTL = dnsConf.TTL
	}
}

// AddUpdateDomainRecords 添加或更新IPv4/IPv6记录
func (cb *Callback) AddUpdateDomainRecords() config.Domains {
	cb.addUpdateDomainRecords("A")
	cb.addUpdateDomainRecords("AAAA")
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
			postPara = replacePara(cb.DNS.Secret, ipAddr, domain, recordType, cb.TTL)
			if json.Valid([]byte(postPara)) {
				contentType = "application/json"
			}
		}
		requestURL := replacePara(cb.DNS.ID, ipAddr, domain, recordType, cb.TTL)
		u, err := url.Parse(requestURL)
		if err != nil {
			util.Log("Callback的URL不正确")
			return
		}
		req, err := http.NewRequest(method, u.String(), strings.NewReader(postPara))
		if err != nil {
			util.Log("异常信息: %s", err)
			domain.UpdateStatus = config.UpdatedFailed
			return
		}
		req.Header.Add("content-type", contentType)

		clt := util.CreateHTTPClient()
		resp, err := clt.Do(req)
		body, err := util.GetHTTPResponseOrg(resp, err)
		if err == nil {
			util.Log("Callback调用成功, 域名: %s, IP: %s, 返回数据: %s", domain, ipAddr, string(body))
			domain.UpdateStatus = config.UpdatedSuccess
		} else {
			util.Log("Callback调用失败, 异常信息: %s", err)
			domain.UpdateStatus = config.UpdatedFailed
		}
	}
}

// replacePara 替换参数
func replacePara(orgPara, ipAddr string, domain *config.Domain, recordType string, ttl string) (newPara string) {
	orgPara = strings.ReplaceAll(orgPara, "#{ip}", ipAddr)
	orgPara = strings.ReplaceAll(orgPara, "#{domain}", domain.String())
	orgPara = strings.ReplaceAll(orgPara, "#{recordType}", recordType)
	orgPara = strings.ReplaceAll(orgPara, "#{ttl}", ttl)

	for k, v := range domain.GetCustomParams() {
		if len(v) == 1 {
			orgPara = strings.ReplaceAll(orgPara, "#{"+k+"}", v[0])
		}
	}

	return orgPara
}
