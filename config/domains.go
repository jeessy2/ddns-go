package config

import (
	"net/url"
	"strings"

	"github.com/jeessy2/ddns-go/v6/util"
	"golang.org/x/net/idna"
	"golang.org/x/net/publicsuffix"
)

// Domains Ipv4/Ipv6 domains
type Domains struct {
	Ipv4Addr    string
	Ipv4Cache   *util.IpCache
	Ipv4Domains []*Domain
	Ipv6Addr    string
	Ipv6Cache   *util.IpCache
	Ipv6Domains []*Domain
}

// Domain 域名实体
type Domain struct {
	// DomainName 根域名
	DomainName string
	// SubDomain 子域名
	SubDomain    string
	CustomParams string
	UpdateStatus updateStatusType // 更新状态
}

// DomainTuples 域名元组映射 key: Domain.String()
type DomainTuples map[string]*DomainTuple

// DomainTuple 域名元组
type DomainTuple struct {
	RecordType string
	// Primary 首要域名 Domains[-1] = Primary
	Primary  *Domain
	Domains  []*Domain
	IpAddrs  []string
	Ipv4Addr string
	Ipv6Addr string
}

// nontransitionalLookup implements the nontransitional processing as specified in
// Unicode Technical Standard 46 with almost all checkings off to maximize user freedom.
//
// Copied from: https://github.com/cloudflare/cloudflare-go/blob/v0.97.0/dns.go#L95
var nontransitionalLookup = idna.New(
	idna.MapForLookup(),
	idna.StrictDomainName(false),
	idna.ValidateLabels(false),
)

func (d Domain) String() string {
	if d.SubDomain != "" {
		return d.SubDomain + "." + d.DomainName
	}
	return d.DomainName
}

// GetFullDomain 获得全部的，子域名
func (d Domain) GetFullDomain() string {
	if d.SubDomain != "" {
		return d.SubDomain + "." + d.DomainName
	}
	return "@" + "." + d.DomainName
}

// GetSubDomain 获得子域名，为空返回@
// 阿里云/腾讯云/dnspod/GoDaddy/namecheap 需要
func (d Domain) GetSubDomain() string {
	if d.SubDomain != "" {
		return d.SubDomain
	}
	return "@"
}

// GetCustomParams not be nil
func (d Domain) GetCustomParams() url.Values {
	if d.CustomParams != "" {
		q, err := url.ParseQuery(d.CustomParams)
		if err == nil {
			return q
		}
	}
	return url.Values{}
}

// ToASCII converts [Domain] to its ASCII form,
// using non-transitional process specified in UTS 46.
//
// Note: conversion errors are silently discarded and partial conversion
// results are used.
func (d Domain) ToASCII() string {
	name, _ := nontransitionalLookup.ToASCII(d.String())
	return name
}

// GetNewIp 接口/网卡/命令获得 ip 并校验用户输入的域名
func (domains *Domains) GetNewIp(dnsConf *DnsConfig) {
	domains.Ipv4Domains = checkParseDomains(dnsConf.Ipv4.Domains)
	domains.Ipv6Domains = checkParseDomains(dnsConf.Ipv6.Domains)

	// IPv4
	if dnsConf.Ipv4.Enable && len(domains.Ipv4Domains) > 0 {
		ipv4Addr := dnsConf.GetIpv4Addr()
		if ipv4Addr != "" {
			domains.Ipv4Addr = ipv4Addr
			domains.Ipv4Cache.TimesFailedIP = 0
		} else {
			// 启用IPv4 & 未获取到IP & 填写了域名 & 失败刚好3次，防止偶尔的网络连接失败，并且只发一次
			domains.Ipv4Cache.TimesFailedIP++
			if domains.Ipv4Cache.TimesFailedIP == 3 {
				domains.Ipv4Domains[0].UpdateStatus = UpdatedFailed
			}
			util.Log("未能获取IPv4地址, 将不会更新")
		}
	}

	// IPv6
	if dnsConf.Ipv6.Enable && len(domains.Ipv6Domains) > 0 {
		ipv6Addr := dnsConf.GetIpv6Addr()
		if ipv6Addr != "" {
			domains.Ipv6Addr = ipv6Addr
			domains.Ipv6Cache.TimesFailedIP = 0
		} else {
			// 启用IPv6 & 未获取到IP & 填写了域名 & 失败刚好3次，防止偶尔的网络连接失败，并且只发一次
			domains.Ipv6Cache.TimesFailedIP++
			if domains.Ipv6Cache.TimesFailedIP == 3 {
				domains.Ipv6Domains[0].UpdateStatus = UpdatedFailed
			}
			util.Log("未能获取IPv6地址, 将不会更新")
		}
	}

}

// checkParseDomains 校验并解析用户输入的域名
func checkParseDomains(domainArr []string) (domains []*Domain) {
	for _, domainStr := range domainArr {
		domainStr = strings.TrimSpace(domainStr)
		if domainStr == "" {
			continue
		}

		domain := &Domain{}

		// qp(queryParts) 从域名中提取自定义参数，如 baidu.com?q=1 => [baidu.com, q=1]
		qp := strings.Split(domainStr, "?")
		domainStr = qp[0]

		// dp(domainParts) 将域名（qp[0]）分割为子域名与根域名，如 www:example.cn.eu.org => [www, example.cn.eu.org]
		dp := strings.Split(domainStr, ":")

		switch len(dp) {
		case 1: // 不使用冒号分割，自动识别域名
			domainName, err := publicsuffix.EffectiveTLDPlusOne(domainStr)
			if err != nil {
				util.Log("域名: %s 不正确", domainStr)
				util.Log("异常信息: %s", err)
				continue
			}
			domain.DomainName = domainName

			domainLen := len(domainStr) - len(domainName) - 1
			if domainLen > 0 {
				domain.SubDomain = domainStr[:domainLen]
			}
		case 2: // 使用冒号分隔，为 子域名:根域名 格式
			sp := strings.Split(dp[1], ".")
			if len(sp) <= 1 {
				util.Log("域名: %s 不正确", domainStr)
				continue
			}
			domain.DomainName = dp[1]
			domain.SubDomain = dp[0]
		default:
			util.Log("域名: %s 不正确", domainStr)
			continue
		}

		// 参数条件
		if len(qp) == 2 {
			u, err := url.Parse("https://baidu.com?" + qp[1])
			if err != nil {
				util.Log("域名: %s 解析失败", domainStr)
				continue
			}
			domain.CustomParams = u.Query().Encode()
		}
		domains = append(domains, domain)
	}
	return
}

// GetNewIpResult 获得GetNewIp结果
func (domains *Domains) GetNewIpResult(recordType string) (ipAddr string, retDomains []*Domain) {
	if recordType == "AAAA" {
		if domains.Ipv6Cache.Check(domains.Ipv6Addr) {
			return domains.Ipv6Addr, domains.Ipv6Domains
		} else {
			util.Log("IPv6未改变, 将等待 %d 次后与DNS服务商进行比对", domains.Ipv6Cache.Times)
			return "", domains.Ipv6Domains
		}
	}
	// IPv4
	if domains.Ipv4Cache.Check(domains.Ipv4Addr) {
		return domains.Ipv4Addr, domains.Ipv4Domains
	} else {
		util.Log("IPv4未改变, 将等待 %d 次后与DNS服务商进行比对", domains.Ipv4Cache.Times)
		return "", domains.Ipv4Domains
	}
}

// GetAllNewIpResult 获得getNewIp结果
func (domains *Domains) GetAllNewIpResult(multiRecordType string) (results DomainTuples) {
	ipv4Addr, ipv4Domains := domains.GetNewIpResult("A")
	ipv6Addr, ipv6Domains := domains.GetNewIpResult("AAAA")
	if ipv4Addr == "" && ipv6Addr == "" {
		return
	}
	cap := 0
	if ipv4Addr != "" {
		cap += len(ipv4Domains)
	}
	if ipv6Addr != "" {
		cap += len(ipv6Domains)
	}

	results = make(DomainTuples, cap)
	results.append(ipv4Addr, ipv4Domains, multiRecordType, DomainTuple{RecordType: "A", Ipv4Addr: domains.Ipv4Addr, Ipv6Addr: domains.Ipv6Addr})
	results.append(ipv6Addr, ipv6Domains, multiRecordType, DomainTuple{RecordType: "AAAA", Ipv4Addr: domains.Ipv4Addr, Ipv6Addr: domains.Ipv6Addr})
	return
}

// append 添加域名到域名元组映射
func (domains DomainTuples) append(ipAddr string, retDomains []*Domain, multiRecordType string, template DomainTuple) {
	if ipAddr == "" {
		return
	}

	for _, domain := range retDomains {
		domainStr := domain.String()
		if tuple, ok := domains[domainStr]; ok {
			if tuple.RecordType != template.RecordType {
				tuple.RecordType = multiRecordType
			}
			tuple.Primary = domain
			tuple.Domains = append(tuple.Domains, domain)
			tuple.IpAddrs = append(tuple.IpAddrs, ipAddr)
		} else {
			tuple := template
			domains[domainStr] = &tuple
			tuple.Primary = domain
			tuple.Domains = []*Domain{domain}
			tuple.IpAddrs = []string{ipAddr}
		}
	}
}

// SetUpdateStatus 设置更新状态
func (d *DomainTuple) SetUpdateStatus(status updateStatusType) {
	if d.Primary.UpdateStatus == status {
		return
	}

	for _, domain := range d.Domains {
		domain.UpdateStatus = status
	}
}

// GetIpAddrPool 设置更新状态
func (d *DomainTuple) GetIpAddrPool(separator string) (result string) {
	s := d.Primary.GetCustomParams().Get("IpAddrPool")
	if len(s) != 0 {
		return strings.NewReplacer(
			"{ipv4Addr}", d.Ipv4Addr,
			"{ipv6Addr}", d.Ipv6Addr,
		).Replace(s)
	}
	switch d.RecordType {
	case "A":
		return d.Ipv4Addr
	case "AAAA":
		return d.Ipv6Addr
	default:
		return d.Ipv4Addr + separator + d.Ipv6Addr
	}
}
