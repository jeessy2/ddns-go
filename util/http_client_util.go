package util

import (
	"context"
	"net"
	"net/http"
	"time"
)

// CreateHTTPClient Create Default HTTP Client
func CreateHTTPClient() *http.Client {
	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: http.DefaultTransport,
	}
}

var noProxyDialer = &net.Dialer{
	Timeout:   30 * time.Second,
	KeepAlive: 30 * time.Second,
}

var noProxyTcp4Transport = &http.Transport{
	// no proxy
	// DisableKeepAlives
	DisableKeepAlives: true,
	// tcp4
	DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
		return noProxyDialer.DialContext(ctx, "tcp4", address)
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
		return noProxyDialer.DialContext(ctx, "tcp6", address)
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
