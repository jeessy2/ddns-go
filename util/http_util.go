package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	proxyCheckUrl = "http://www.gstatic.com/generate_204"
)

// GetHTTPResponse 处理HTTP结果，返回序列化的json
func GetHTTPResponse(resp *http.Response, url string, err error, result interface{}) error {
	body, err := GetHTTPResponseOrg(resp, url, err)

	if err == nil {
		// log.Println(string(body))
		if len(body) != 0 {
			err = json.Unmarshal(body, &result)
			if err != nil {
				log.Printf("请求接口%s解析json结果失败! ERROR: %s\n", url, err)
			}
		}
	}

	return err

}

// GetHTTPResponseOrg 处理HTTP结果，返回byte
func GetHTTPResponseOrg(resp *http.Response, url string, err error) ([]byte, error) {
	if err != nil {
		log.Printf("请求接口%s失败! ERROR: %s\n", url, err)
		Ipv4Cache.ForceCompare = true
		Ipv6Cache.ForceCompare = true
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		Ipv4Cache.ForceCompare = true
		Ipv6Cache.ForceCompare = true
		log.Printf("请求接口%s失败! ERROR: %s\n", url, err)
	}

	// 300及以上状态码都算异常
	if resp.StatusCode >= 300 {
		errMsg := fmt.Sprintf("请求接口 %s 失败! 返回内容: %s ,返回状态码: %d\n", url, string(body), resp.StatusCode)
		log.Println(errMsg)
		err = fmt.Errorf(errMsg)
	}

	return body, err
}

// CreateHTTPClient CreateHTTPClient
func CreateHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			IdleConnTimeout:     10 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}
}

func CreateHTTPClientWithProxy(proxyUrl string) *http.Client {

	// Proxy URL check: format check
	client := CreateHTTPClient()
	_, err := url.ParseRequestURI(proxyUrl)
	if err != nil {
		log.Println("Proxy parse error, disable the proxy")
		return client
	}

	// Set proxy url
	client.Transport.(*http.Transport).Proxy = func(_ *http.Request) (*url.URL, error) {
		return url.Parse(proxyUrl)
	}

	// Check if proxy is alive
	req, err := http.NewRequest(http.MethodHead, proxyCheckUrl, nil)
	if err != nil {
		log.Println("Proxy test request create error, disable the proxy")
		client.Transport.(*http.Transport).Proxy = nil
		return client
	}
	resp, err := client.Do(req)
	if err != nil || resp == nil || resp.StatusCode >= 300 {
		// Not alive, proxy will not be set
		log.Println("Proxy test failed (Cannot access", proxyCheckUrl, "using proxy", proxyUrl, "), disable the proxy")
		client.Transport.(*http.Transport).Proxy = nil
		return client
	}

	// Return HTTP client with proxy
	return client
}
