package dns

import (
	"log"
	"time"

	"github.com/jeessy2/ddns-go/v4/config"
	"github.com/jeessy2/ddns-go/v4/domainprovider"
)

var _dnsRunnerMap = map[string]DNS{}

func RegisterDNS(dns DNS) {
	if dns == nil {
		panic("dns is nil")
	}
	_dnsRunnerMap[dns.Code()] = dns
	log.Printf("注册DNS[%s],Name:%s", dns.Code(), dns.Name())
}

// DNS interface
type DNS interface {
	// DNS的编码
	Code() string
	// DNS的名称
	Name() string
	// 初始化配置
	Init(conf *config.Config)
	// 添加或更新IPv4/IPv6记录
	AddUpdateDomainRecords() (domains config.Domains)
	// 添加或更新IPv4/IPv6记录
	AddUpdateDomainRecordsFromDomains(domains []*config.Domain) config.Domains
}

// RunTimer 定时运行
func RunTimer(firstDelay time.Duration, delay time.Duration) {
	log.Printf("第一次运行将等待 %d 秒后运行 (等待网络)", int(firstDelay.Seconds()))
	time.Sleep(firstDelay)
	for {
		RunOnce()
		time.Sleep(delay)
	}
}

// RunOnce RunOnce
func RunOnce() {
	conf, err := config.GetConfigCache()
	if err != nil {
		return
	}
	// var dnsSelected DNS
	dnsSelected, ok := _dnsRunnerMap[conf.DNS.Name]
	if !ok {
		log.Printf("未找到 %s DNS", conf.DNS.Name)
		return
	}
	dnsSelected.Init(&conf)
	sourceDomains := domainprovider.GetDomains()
	domains := dnsSelected.AddUpdateDomainRecordsFromDomains(sourceDomains)
	config.ExecWebhook(&domains, &conf)
}
