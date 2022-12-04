package domainprovider

import (
	"log"

	"github.com/jeessy2/ddns-go/v4/config"
)

type DomainProviderInterface interface {
	Code() string
	Init(conf *config.Config) error
	GetDomains() []*config.Domain
}

var _domainProviderMap = make(map[string]DomainProviderInterface)

func RegisterDomainProvider(provider DomainProviderInterface) {
	_domainProviderMap[provider.Code()] = provider
	log.Printf("注册Domain's Provider[%s]", provider.Code())
}

func GetDomains() []*config.Domain {
	var domains []*config.Domain
	for _, provider := range _domainProviderMap {
		conf, err := config.GetConfigCache()
		if err != nil {
			log.Printf("获取配置失败: %s", err)
			continue
		}
		provider.Init(&conf)
		domains = append(domains, provider.GetDomains()...)
	}
	return domains
}
