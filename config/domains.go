package config

import (
	"log"
	"strings"
)

// 固定的主域名
var staticMainDomains = []string{"com.cn", "org.cn", "net.cn", "ac.cn"}

// Domains Ipv4/Ipv6 domains
type Domains struct {
	Ipv4Addr    string
	Ipv4Domains []*Domain
	Ipv6Addr    string
	Ipv6Domains []*Domain
}

// Domain 域名实体
type Domain struct {
	DomainName   string
	SubDomain    string
	UpdateStatus updateStatusType // 更新状态
}

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
// 阿里云，dnspod需要
func (d Domain) GetSubDomain() string {
	if d.SubDomain != "" {
		return d.SubDomain
	}
	return "@"
}

// ParseDomain 接口获得ip并校验用户输入的域名
func (domains *Domains) ParseDomain(conf *Config) {
	// IPV4
	ipv4Addr := conf.GetIpv4Addr()
	if ipv4Addr != "" {
		domains.Ipv4Addr = ipv4Addr
		domains.Ipv4Domains = parseDomainArr(conf.Ipv4.Domains)
	}
	// IPV6
	ipv6Addr := conf.GetIpv6Addr()
	if ipv6Addr != "" {
		domains.Ipv6Addr = ipv6Addr
		domains.Ipv6Domains = parseDomainArr(conf.Ipv6.Domains)
	}
}

// parseDomainArr 校验用户输入的域名
func parseDomainArr(domainArr []string) (domains []*Domain) {
	for _, domainStr := range domainArr {
		domainStr = strings.TrimSpace(domainStr)
		if domainStr != "" {
			domain := &Domain{}
			sp := strings.Split(domainStr, ".")
			length := len(sp)
			if length <= 1 {
				log.Println(domainStr, "域名不正确")
				continue
			}
			// 处理域名
			domain.DomainName = sp[length-2] + "." + sp[length-1]
			// 如包含在org.cn等顶级域名下，后三个才为用户主域名
			for _, staticMainDomain := range staticMainDomains {
				if staticMainDomain == domain.DomainName {
					domain.DomainName = sp[length-3] + "." + domain.DomainName
					break
				}
			}

			domainLen := len(domainStr) - len(domain.DomainName)
			if domainLen > 0 {
				domain.SubDomain = domainStr[:domainLen-1]
			} else {
				domain.SubDomain = domainStr[:domainLen]
			}

			domains = append(domains, domain)
		}
	}
	return
}

// ParseDomainResult 获得ParseDomain结果
func (domains *Domains) ParseDomainResult(recordType string) (ipAddr string, retDomains []*Domain) {
	if recordType == "AAAA" {
		return domains.Ipv6Addr, domains.Ipv6Domains
	}
	return domains.Ipv4Addr, domains.Ipv4Domains

}
