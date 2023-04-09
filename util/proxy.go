package util

import (
	"log"
	"net/http"
	"net/url"
	"os"
)

var proxyEnvs = []string{"HTTP_PROXY", "HTTPS_PROXY", "http_proxy", "https_proxy"}

// getHTTPProxy 获取 HTTP 代理变量的值。如果值为 URL 则使用值，否则使用 ProxyFromEnvironment。
func getHTTPProxy() func(*http.Request) (*url.URL, error) {
	for _, key := range proxyEnvs {
		proxy := os.Getenv(key)
		if proxy != "" {
			proxyURL, err := url.Parse(proxy)
			if err != nil {
				log.Println("解析 HTTP 代理失败：", err)
				return http.ProxyFromEnvironment
			}
			return http.ProxyURL(proxyURL)
		}
	}
	return http.ProxyFromEnvironment
}
