package dns

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/jeessy2/ddns-go/v4/config"
	"github.com/jeessy2/ddns-go/v4/util"
)

type Callback struct {
	DNSConfig config.DNSConfig
	Domains   config.Domains
	TTL       string
}

// Init 初始化
func (cb *Callback) Init(conf *config.Config) {
	cb.DNSConfig = conf.DNS
	cb.Domains.GetNewIp(conf)
	if conf.TTL == "" {
		// 默认600
		cb.TTL = "600"
	} else {
		cb.TTL = conf.TTL
	}
}

// AddUpdateDomainRecords 添加或更新IPv4/IPv6记录
func (cb *Callback) AddUpdateDomainRecords() config.Domains {
	cb.addUpdateDomainRecords("A")
	cb.addUpdateDomainRecords("AAAA")
	return cb.Domains
}

var lastIpv4 string
var lastIpv6 string

func (cb *Callback) addUpdateDomainRecords(recordType string) {
	ipAddr, domains := cb.Domains.GetNewIpResult(recordType)

	if ipAddr == "" {
		return
	}

	if recordType == "A" {
		if lastIpv4 == ipAddr {
			log.Println("你的IPv4未变化, 未触发Callback")
			return
		}
		lastIpv4 = ipAddr
	} else {
		if lastIpv6 == ipAddr {
			log.Println("你的IPv6未变化, 未触发Callback")
			return
		}
		lastIpv6 = ipAddr
	}

	for _, domain := range domains {
		method := "GET"
		postPara := ""
		contentType := "application/x-www-form-urlencoded"
		if cb.DNSConfig.Secret != "" {
			method = "POST"
			postPara = replacePara(cb.DNSConfig.Secret, ipAddr, domain, recordType, cb.TTL)
			if json.Valid([]byte(postPara)) {
				contentType = "application/json"
			}
		}
		requestURL := replacePara(cb.DNSConfig.ID, ipAddr, domain, recordType, cb.TTL)
		u, err := url.Parse(requestURL)
		if err != nil {
			log.Println("Callback的URL不正确")
			return
		}
		req, err := http.NewRequest(method, u.String(), strings.NewReader(postPara))
		if err != nil {
			log.Println("创建Callback请求异常, Err:", err)
			return
		}
		req.Header.Add("content-type", contentType)

		clt := util.CreateHTTPClient()
		resp, err := clt.Do(req)
		body, err := util.GetHTTPResponseOrg(resp, requestURL, err)
		if err == nil {
			log.Printf("Callback调用成功, 返回数据: %s\n", string(body))
			domain.UpdateStatus = config.UpdatedSuccess
		} else {
			log.Printf("Callback调用失败，Err：%s\n", err)
			domain.UpdateStatus = config.UpdatedFailed
		}
	}
}

// replacePara 替换参数
func replacePara(orgPara, ipAddr string, domain *config.Domain, recordType string, ttl string) (newPara string) {
	orgPara = strings.ReplaceAll(orgPara, "#{ip}", ipAddr)
	orgPara = strings.ReplaceAll(orgPara, "#{domain}", domain.String())
	orgPara = strings.ReplaceAll(orgPara, "#{recordType}", recordType)
	orgPara = strings.ReplaceAll(orgPara, "#{ttl}", ttl)

	for k, v := range domain.GetCustomParams() {
		if len(v) == 1 {
			orgPara = strings.ReplaceAll(orgPara, "#{"+k+"}", v[0])
		}
	}

	return orgPara
}
