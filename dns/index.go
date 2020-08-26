package dns

import "ddns-go/config"

// DNS interface
type DNS interface {
	AddRecord(conf *config.Config) (ipv4 bool, ipv6 bool)
}
