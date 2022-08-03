package dns

import (
	"log"
	"time"

	"github.com/jeessy2/ddns-go/v4/config"
)

// DNS interface
type DNS interface {
	Init(conf *config.Config)
	// 添加或更新IPv4/IPv6记录
	AddUpdateDomainRecords() (domains config.Domains)
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

	var dnsSelected DNS
	switch conf.DNS.Name {
	case "alidns":
		dnsSelected = &Alidns{}
	case "dnspod":
		dnsSelected = &Dnspod{}
	case "cloudflare":
		dnsSelected = &Cloudflare{}
	case "huaweicloud":
		dnsSelected = &Huaweicloud{}
	case "callback":
		dnsSelected = &Callback{}
	case "baiducloud":
		dnsSelected = &BaiduCloud{}
	case "porkbun":
		dnsSelected = &Porkbun{}
	case "godaddy":
		dnsSelected = &GoDaddyDNS{}
	case "googledomain":
		dnsSelected = &GoogleDomain{}
	default:
		dnsSelected = &Alidns{}
	}
	dnsSelected.Init(&conf)

	domains := dnsSelected.AddUpdateDomainRecords()
	config.ExecWebhook(&domains, &conf)
}
