package dns

import (
	"ddns-go/config"
	"log"
	"strings"
	"time"
)

// DNS interface
type DNS interface {
	Init(conf *config.Config)
	// 添加或更新IPV4记录
	AddUpdateIpv4DomainRecords()
	// 添加或更新IPV6记录
	AddUpdateIpv6DomainRecords()
}

// Domains Ipv4/Ipv6 domains
type Domains struct {
	Ipv4Addr    string
	Ipv4Domains []*Domain
	Ipv6Addr    string
	Ipv6Domains []*Domain
}

// Domain 域名实体
type Domain struct {
	DomainName string
	SubDomain  string
	Exist      bool
}

func (d Domain) String() string {
	if d.SubDomain != "" {
		return d.SubDomain + "." + d.DomainName
	}
	return d.DomainName
}

// RunTimer 定时运行
func RunTimer() {
	for {
		RunOnce()
		time.Sleep(time.Minute * time.Duration(5))
	}
}

// RunOnce RunOnce
func RunOnce() {
	conf := &config.Config{}
	err := conf.InitConfigFromFile()
	if err != nil {
		return
	}

	var dnsSelected DNS
	switch conf.DNS.Name {
	case "alidns":
		dnsSelected = &Alidns{}
	case "dnspod":
		dnsSelected = &Dnspod{}
	case "cloudflare":
		dnsSelected = &Cloudflare{}
	default:
		dnsSelected = &Alidns{}
	}
	dnsSelected.Init(conf)
	dnsSelected.AddUpdateIpv4DomainRecords()
	dnsSelected.AddUpdateIpv6DomainRecords()
}

// ParseDomain 解析域名
func (domains *Domains) ParseDomain(conf *config.Config) {
	// IPV4
	ipv4Addr := conf.GetIpv4Addr()
	if ipv4Addr != "" {
		domains.Ipv4Addr = ipv4Addr
		domains.Ipv4Domains = parseDomainInner(conf.Ipv4.Domains)
	}

	// IPV6
	ipv6Addr := conf.GetIpv6Addr()
	if ipv6Addr != "" {
		domains.Ipv6Addr = ipv6Addr
		domains.Ipv6Domains = parseDomainInner(conf.Ipv6.Domains)
	}
}

// parseDomainInner 解析域名inner
func parseDomainInner(domainArr []string) (domains []*Domain) {
	for _, domainStr := range domainArr {
		domainStr = strings.Trim(domainStr, " ")
		if domainStr != "" {
			domain := &Domain{}
			sp := strings.Split(domainStr, ".")
			length := len(sp)
			if length <= 1 {
				log.Println(domainStr, "域名不正确")
				continue
			} else if length == 2 {
				domain.DomainName = domainStr
			} else {
				// >=3
				domain.DomainName = sp[length-2] + "." + sp[length-1]
				domain.SubDomain = domainStr[:len(domainStr)-len(domain.DomainName)-1]
			}
			domains = append(domains, domain)
		}
	}
	return
}
