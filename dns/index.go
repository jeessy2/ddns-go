package dns

import (
	"log"
	"time"

	"github.com/jeessy2/ddns-go/v4/config"
	"github.com/jeessy2/ddns-go/v4/util"
)

// DNS interface
type DNS interface {
	Init(conf *config.Config)
	// 添加或更新IPv4/IPv6记录
	AddUpdateDomainRecords() (domains config.Domains)
}

var domainCache map[string]DNS

func DNSInit() {
	domainCache = map[string]DNS{
		"alidns":       &Alidns{},
		"dnspod":       &Dnspod{},
		"cloudflare":   &Cloudflare{},
		"huaweicloud":  &Huaweicloud{},
		"callback":     &Callback{},
		"baiducloud":   &BaiduCloud{},
		"porkbun":      &Porkbun{},
		"godaddy":      &GoDaddyDNS{},
		"googledomain": &GoogleDomain{},
	}
}

// RunTimer 定时运行
func RunTimer(firstDelay time.Duration, delay time.Duration) {
	log.Printf("第一次运行将等待 %d 秒后运行 (等待网络)", int(firstDelay.Seconds()))
	DNSInit()
	time.Sleep(firstDelay)
	for {
		RunOnce()
		time.Sleep(delay)
	}
}

// RunOnce RunOnce
func RunOnce() {
	cglobal, err := config.GetConfigGlobal()
	if err != nil {
		return
	}
	cmap := config.GetConfigMap()

	for name, conf := range cmap {
		domainCache[name].Init(&conf)
		domains := domainCache[name].AddUpdateDomainRecords()
		config.ExecWebhook(&domains, &cglobal)
	}
	util.ForceCompare = false
}
