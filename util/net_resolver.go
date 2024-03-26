package util

import (
	"context"
	"net"

	"golang.org/x/text/language"
)

// BackupDNS will be used if DNS error occurs.
var BackupDNS = []string{}

func InitDefaultDNS(customDNS, lang string) {
	if customDNS != "" {
		BackupDNS = []string{customDNS}
		return
	}

	if lang == language.Chinese.String() {
		BackupDNS = []string{"223.5.5.5", "114.114.114.114"}
		return
	}

	BackupDNS = []string{"1.1.1.1", "8.8.8.8"}
}

// SetDNS sets the dialer.Resolver to use the given DNS server.
func SetDNS(dns string) {
	// Error means that the given DNS doesn't have a port. Add it.
	if _, _, err := net.SplitHostPort(dns); err != nil {
		dns = net.JoinHostPort(dns, "53")
	}

	dialer.Resolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return net.Dial(network, dns)
		},
	}
}

// LookupHost looks up the host based on the given URL using the dialer.Resolver.
// A wrapper for [net.Resolver.LookupHost].
func LookupHost(url string) error {
	name := toHostname(url)

	_, err := dialer.Resolver.LookupHost(context.Background(), name)
	return err
}
