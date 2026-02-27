package util

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"
)

var dialer = &net.Dialer{
	Timeout:   30 * time.Second,
	KeepAlive: 30 * time.Second,
}

var defaultTransport = &http.Transport{
	// from http.DefaultTransport
	Proxy: http.ProxyFromEnvironment,
	DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
		return dialer.DialContext(ctx, network, address)
	},
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

// insecureSkipVerify 全局TLS验证跳过标志
var insecureSkipVerify bool

// CreateHTTPClient Create Default HTTP Client
func CreateHTTPClient() *http.Client {
	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: defaultTransport,
	}
}

// GetLocalAddrFromInterface 根据网卡名称获取本地IP地址
func GetLocalAddrFromInterface(ifaceName string) (string, error) {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return "", fmt.Errorf("找不到网卡 %s: %v", ifaceName, err)
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return "", fmt.Errorf("获取网卡 %s 地址失败: %v", ifaceName, err)
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && ipNet.IP.IsGlobalUnicast() {
			return ipNet.IP.String(), nil
		}
	}
	return "", fmt.Errorf("网卡 %s 没有可用的单播地址", ifaceName)
}

// CreateHTTPClientWithInterface 创建绑定指定网卡的HTTP客户端
func CreateHTTPClientWithInterface(ifaceName string) *http.Client {
	if ifaceName == "" {
		return CreateHTTPClient()
	}
	localIP, err := GetLocalAddrFromInterface(ifaceName)
	if err != nil {
		Log("绑定网卡失败, 将使用默认网卡. 网卡: %s, 错误: %s", ifaceName, err)
		return CreateHTTPClient()
	}
	localAddr := &net.TCPAddr{IP: net.ParseIP(localIP)}
	boundDialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		LocalAddr: localAddr,
	}
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			return boundDialer.DialContext(ctx, network, address)
		},
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	if insecureSkipVerify {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}
}

var noProxyTcp4Transport = &http.Transport{
	// no proxy
	// DisableKeepAlives
	DisableKeepAlives: true,
	// tcp4
	DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
		return dialer.DialContext(ctx, "tcp4", address)
	},
	// from http.DefaultTransport
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

var noProxyTcp6Transport = &http.Transport{
	// no proxy
	// DisableKeepAlives
	DisableKeepAlives: true,
	// tcp6
	DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
		return dialer.DialContext(ctx, "tcp6", address)
	},
	// from http.DefaultTransport
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

// CreateNoProxyHTTPClient Create NoProxy HTTP Client
func CreateNoProxyHTTPClient(network string) *http.Client {
	if network == "tcp6" {
		return &http.Client{
			Timeout:   30 * time.Second,
			Transport: noProxyTcp6Transport,
		}
	}

	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: noProxyTcp4Transport,
	}
}

// SetInsecureSkipVerify 将所有 http.Transport 的 InsecureSkipVerify 设置为 true
func SetInsecureSkipVerify() {
	insecureSkipVerify = true
	transports := []*http.Transport{defaultTransport, noProxyTcp4Transport, noProxyTcp6Transport}

	for _, transport := range transports {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
}
