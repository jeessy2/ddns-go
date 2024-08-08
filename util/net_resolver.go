package util

import (
	"context"
	"net"
	"net/url"
	"strings"

	"golang.org/x/text/language"
)

// BackupDNS will be used if DNS error occurs.
var BackupDNS = []string{"1.1.1.1", "8.8.8.8", "9.9.9.9", "223.5.5.5"}

func InitBackupDNS(customDNS, lang string) {
	if customDNS != "" {
		BackupDNS = []string{customDNS}
		return
	}

	if lang == language.Chinese.String() {
		BackupDNS = []string{"223.5.5.5", "114.114.114.114", "119.29.29.29"}
	}

}

// SetDNS sets the dialer.Resolver to use the given DNS server.
func SetDNS(dns string) {

	if !strings.Contains(dns, "://") {
		dns = "udp://" + dns
	}
	svrParse, _ := url.Parse(dns)

	var network string
	switch strings.ToLower(svrParse.Scheme) {
	case "tcp":
		network = "tcp"
	default:
		network = "udp"
	}

	if svrParse.Port() == "" {
		dns = net.JoinHostPort(svrParse.Host, "53")
	} else {
		dns = svrParse.Host
	}

	dialer.Resolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, _, address string) (net.Conn, error) {
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
