package util

import (
	"log"
	"net/http"
	"net/url"
	"os"
)

const HTTPProxyEnv = "DDNS_GO_HTTP_PROXY"

// getHTTPProxy 获取 HTTPProxyEnv 的值。如果值为 URL 则使用值的 HTTP 代理，否则使用变量的 HTTP 代理。
func getHTTPProxy() func(*http.Request) (*url.URL, error) {
	proxy := os.Getenv(HTTPProxyEnv)
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			log.Println("解析 HTTP 代理失败：", err)
			return http.ProxyFromEnvironment
		}
		return http.ProxyURL(proxyURL)
	}
	return http.ProxyFromEnvironment
}
