package domainprovider

import (
	"log"
	"net/url"
	"strings"

	"github.com/jeessy2/ddns-go/v4/config"
)

func init() {
	RegisterDomainProvider(&DomainConfigProvider{})
}

type DomainConfigProvider struct {
	config *config.Config
}

func (d *DomainConfigProvider) Code() string {
	return "config"
}

func (d *DomainConfigProvider) Init(conf *config.Config) error {
	d.config = conf
	return nil
}

func (d *DomainConfigProvider) GetDomains() []*config.Domain {
	var domains []*config.Domain
	domains = append(domains, d.getIPV4Domains()...)
	domains = append(domains, d.getIPV6Domains()...)
	return domains
}

func (d *DomainConfigProvider) getIPV6Domains() []*config.Domain {
	// return config.GetConfigCache().Domains.IPV4
	if !d.config.Ipv6.Enable {
		return []*config.Domain{}
	}
	return checkParseDomains(d.config.Ipv4.Domains)
}

func (d *DomainConfigProvider) getIPV4Domains() []*config.Domain {
	// return config.GetConfigCache().Domains.IPV4
	if !d.config.Ipv4.Enable {
		return []*config.Domain{}
	}
	return checkParseDomains(d.config.Ipv4.Domains)
}

// checkParseDomains 校验并解析用户输入的域名
func checkParseDomains(domainArr []string) (domains []*config.Domain) {
	for _, domainStr := range domainArr {
		domainStr = strings.TrimSpace(domainStr)
		if domainStr != "" {
			domain := &config.Domain{}
			dp := strings.Split(domainStr, ":")
			dplen := len(dp)
			if dplen == 1 { // 自动识别域名
				sp := strings.Split(domainStr, ".")
				length := len(sp)
				if length <= 1 {
					log.Println(domainStr, "域名不正确")
					continue
				}
				// 处理域名
				domain.DomainName = sp[length-2] + "." + sp[length-1]
				// 如包含在org.cn等顶级域名下，后三个才为用户主域名
				staticMainDomains := config.GetStaticMainDomains()
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

			} else if dplen == 2 { // 主机记录:域名 格式
				sp := strings.Split(dp[1], ".")
				length := len(sp)
				if length <= 1 {
					log.Println(domainStr, "域名不正确")
					continue
				}
				domain.DomainName = dp[1]
				domain.SubDomain = dp[0]
			} else {
				log.Println(domainStr, "域名不正确")
				continue
			}

			// 参数条件
			if strings.Contains(domain.DomainName, "?") {
				u, err := url.Parse("http://" + domain.DomainName)
				if err != nil {
					log.Println(domainStr, "域名解析失败")
					continue
				}
				domain.DomainName = u.Host
				domain.CustomParams = u.Query().Encode()
			}
			domains = append(domains, domain)
		}
	}
	return
}
