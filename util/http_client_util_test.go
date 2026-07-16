package util

import (
	"net"
	"testing"
)

func TestNewHTTPDialerInheritsCustomDNSResolver(t *testing.T) {
	originalResolver := dialer.Resolver
	t.Cleanup(func() {
		dialer.Resolver = originalResolver
	})

	customResolver := &net.Resolver{PreferGo: true}
	dialer.Resolver = customResolver

	boundDialer := newHTTPDialer(&net.TCPAddr{IP: net.ParseIP("192.0.2.1")})
	if boundDialer.Resolver != customResolver {
		t.Fatal("bound HTTP dialer did not inherit the custom DNS resolver")
	}
}
