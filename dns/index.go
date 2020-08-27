package dns

import (
	"ddns-go/config"
	"log"
	"strings"
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
	return d.SubDomain + "." + d.DomainName
}

// ParseDomain 解析域名
func ParseDomain(domainArr []string) (domains []*Domain) {
	for _, domainStr := range domainArr {
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
	return
}
