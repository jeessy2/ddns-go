package dns

import (
	"ddns-go/config"
	"ddns-go/util"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// Webhook Webhook
type Webhook struct {
	DNSConfig config.DNSConfig
	Domains
}

// Init 初始化
func (wh *Webhook) Init(conf *config.Config) {
	wh.DNSConfig = conf.DNS
	wh.Domains.ParseDomain(conf)
}

var beforeDomains Domains

// AddUpdateDomainRecords 添加或更新IPV4/IPV6记录
func (wh *Webhook) AddUpdateDomainRecords() {
	if beforeDomains.Ipv4Addr == "" && beforeDomains.Ipv6Addr == "" {
		// 以前为空，现在有什么也不做
		log.Println("暂时没有以前的IP记录, 不会调用Webhook, 等待下次")
	} else {
		// 有一处变化，webhook...
		if beforeDomains.Ipv4Addr != wh.Domains.Ipv4Addr || beforeDomains.Ipv6Addr != wh.Domains.Ipv6Addr {
			method := "GET"
			postPara := ""
			if wh.DNSConfig.Secret != "" {
				method = "POST"
				postPara = wh.replacePara(wh.DNSConfig.Secret)
			}

			req, err := http.NewRequest(method, wh.replacePara(wh.DNSConfig.ID), strings.NewReader(postPara))

			clt := http.Client{}
			clt.Timeout = 30 * time.Second
			resp, err := clt.Do(req)
			body, err := util.GetHTTPResponseOrg(resp, wh.replacePara(wh.DNSConfig.ID), err)
			if err == nil {
				if wh.Domains.Ipv4Addr != "" {
					log.Println(fmt.Sprintf("Webhook调用成功, 新IPV4: %s", wh.Domains.Ipv4Addr))
				}
				if wh.Domains.Ipv6Addr != "" {
					log.Println(fmt.Sprintf("Webhook调用成功, 新IPV6: %s", wh.Domains.Ipv6Addr))
				}
				log.Printf("返回数据: %s", string(body))
			} else {
				log.Println("Webhook调用失败, 下次重新调用")
				// 不改变旧值，下次重新调用
				return
			}
		} else {
			if wh.Domains.Ipv4Addr != "" {
				log.Println(fmt.Sprintf("你的IP: %s 没有变化, 未调用Webhook", wh.Domains.Ipv4Addr))
			}
			if wh.Domains.Ipv6Addr != "" {
				log.Println(fmt.Sprintf("你的IP: %s 没有变化, 未调用Webhook", wh.Domains.Ipv6Addr))
			}
		}
	}

	// required
	beforeDomains = wh.Domains
}

// replacePara 替换参数
func (wh *Webhook) replacePara(orgPara string) (newPara string) {
	orgPara = strings.ReplaceAll(orgPara, "#{ipv4New}", wh.Domains.Ipv4Addr)
	orgPara = strings.ReplaceAll(orgPara, "#{ipv4Old}", beforeDomains.Ipv4Addr)
	orgPara = strings.ReplaceAll(orgPara, "#{ipv6New}", wh.Domains.Ipv6Addr)
	orgPara = strings.ReplaceAll(orgPara, "#{ipv6Old}", beforeDomains.Ipv6Addr)
	orgPara = strings.ReplaceAll(orgPara, "#{ipv4Domains}", getDomainsStr(wh.Domains.Ipv4Domains))
	orgPara = strings.ReplaceAll(orgPara, "#{ipv6Domains}", getDomainsStr(wh.Domains.Ipv6Domains))

	return orgPara
}

// getDomainsStr 用逗号分割域名
func getDomainsStr(domains []*Domain) string {
	str := ""
	for i, v46 := range domains {
		str += v46.String()
		if i != len(domains)-1 {
			str += ","
		}
	}

	return str
}
