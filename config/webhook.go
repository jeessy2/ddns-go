package config

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/jeessy2/ddns-go/v4/util"
)

// Webhook Webhook
type Webhook struct {
	WebhookURL         string
	WebhookRequestBody string
}

// updateStatusType 更新状态
type updateStatusType string

const (
	// UpdatedNothing 未改变
	UpdatedNothing updateStatusType = "未改变"
	// UpdatedFailed 更新失败
	UpdatedFailed = "失败"
	// UpdatedSuccess 更新成功
	UpdatedSuccess = "成功"
)

// ExecWebhook 添加或更新IPv4/IPv6记录
func ExecWebhook(domains *Domains, conf *Config) {
	v4Status := getDomainsStatus(domains.Ipv4Domains)
	v6Status := getDomainsStatus(domains.Ipv6Domains)

	if conf.WebhookURL != "" && (v4Status != UpdatedNothing || v6Status != UpdatedNothing) {
		// 成功和失败都要触发webhook
		method := "GET"
		postPara := ""
		contentType := "application/x-www-form-urlencoded"
		if conf.WebhookRequestBody != "" {
			method = "POST"
			postPara = replacePara(domains, conf.WebhookRequestBody, v4Status, v6Status)
			if json.Valid([]byte(postPara)) {
				contentType = "application/json"
			}
		}
		requestURL := replacePara(domains, conf.WebhookURL, v4Status, v6Status)
		u, err := url.Parse(requestURL)
		if err != nil {
			log.Println("Webhook配置中的URL不正确")
			return
		}
		req, err := http.NewRequest(method, fmt.Sprintf("%s://%s%s?%s", u.Scheme, u.Host, u.Path, u.Query().Encode()), strings.NewReader(postPara))
		if err != nil {
			log.Println("创建Webhook请求异常, Err:", err)
			return
		}
		req.Header.Add("content-type", contentType)

		clt := util.CreateHTTPClient()
		resp, err := clt.Do(req)
		body, err := util.GetHTTPResponseOrg(resp, requestURL, err)
		if err == nil {
			log.Println(fmt.Sprintf("Webhook调用成功, 返回数据: %s", string(body)))
		} else {
			log.Println(fmt.Sprintf("Webhook调用失败，Err：%s", err))
		}
	}
}

// getDomainsStr 用逗号分割域名
func getDomainsStatus(domains []*Domain) updateStatusType {
	successNum := 0
	for _, v46 := range domains {
		switch v46.UpdateStatus {
		case UpdatedFailed:
			// 一个失败，全部失败
			return UpdatedFailed
		case UpdatedSuccess:
			successNum++
		}
	}

	if successNum > 0 {
		// 迭代完成后一个成功，就成功
		return UpdatedSuccess
	}
	return UpdatedNothing
}

// replacePara 替换参数
func replacePara(domains *Domains, orgPara string, ipv4Result updateStatusType, ipv6Result updateStatusType) (newPara string) {
	orgPara = strings.ReplaceAll(orgPara, "#{ipv4Addr}", domains.Ipv4Addr)
	orgPara = strings.ReplaceAll(orgPara, "#{ipv4Result}", string(ipv4Result))
	orgPara = strings.ReplaceAll(orgPara, "#{ipv4Domains}", getDomainsStr(domains.Ipv4Domains))

	orgPara = strings.ReplaceAll(orgPara, "#{ipv6Addr}", domains.Ipv6Addr)
	orgPara = strings.ReplaceAll(orgPara, "#{ipv6Result}", string(ipv6Result))
	orgPara = strings.ReplaceAll(orgPara, "#{ipv6Domains}", getDomainsStr(domains.Ipv6Domains))

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
