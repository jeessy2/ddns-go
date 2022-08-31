package web

import (
	"log"
	"net/http"
	"strings"

	"github.com/jeessy2/ddns-go/v4/config"
)

// WebhookTest 测试webhook
func WebhookTest(writer http.ResponseWriter, request *http.Request) {
	url := strings.TrimSpace(request.FormValue("URL"))
	requestBody := strings.TrimSpace(request.FormValue("RequestBody"))

	var domains = make([]*config.Domain, 1)
	domains[0] = &config.Domain{}
	domains[0].DomainName = "example.com"
	domains[0].SubDomain = "test"
	domains[0].UpdateStatus = config.UpdatedSuccess

	fakeDomains := &config.Domains{
		Ipv4Addr:    "127.0.0.1",
		Ipv4Domains: domains,
		Ipv6Addr:    "::1",
		Ipv6Domains: domains,
	}

	fakeConfig := &config.Config{
		Webhook: config.Webhook{
			WebhookURL:         url,
			WebhookRequestBody: requestBody,
		},
	}

	if url != "" {
		config.ExecWebhook(fakeDomains, fakeConfig)
	} else {
		log.Println("请输入Webhook的URL")
	}
}
