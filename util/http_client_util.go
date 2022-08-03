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

var dialer = &net.Dialer{
	Timeout:   30 * time.Second,
	KeepAlive: 30 * time.Second,
}

// CreateNoProxyHTTPClient CreateNoProxyHTTPClient
func CreateNoProxyHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			// no proxy
			// from http.DefaultTransport
			DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
				return dialer.DialContext(ctx, network, address)
			},
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}
