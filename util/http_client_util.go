package util

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"time"
)

const SkipVerifyENV = "DDNS_SKIP_VERIFY"

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

// CreateHTTPClient Create Default HTTP Client
func CreateHTTPClient() *http.Client {
	dialer.Resolver = CustomDNSResolver()
	// SkipVerfiry
	defaultTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: os.Getenv(SkipVerifyENV) == "true"}
	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: defaultTransport,
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
	dialer.Resolver = CustomDNSResolver()
	if network == "tcp6" {
		// SkipVerfiry
		noProxyTcp6Transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: os.Getenv(SkipVerifyENV) == "true"}
		return &http.Client{
			Timeout:   30 * time.Second,
			Transport: noProxyTcp6Transport,
		}
	}

	// SkipVerfiry
	noProxyTcp4Transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: os.Getenv(SkipVerifyENV) == "true"}
	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: noProxyTcp4Transport,
	}
}
