package main

import (
	"ddns-go/dns"
)

const (
	ipv4Addr = "https://api-ipv4.ip.sb/ip"
	ipv6Addr = "https://api-ipv6.ip.sb/ip"
)

func main() {
	conf := &Config{}
	conf.getConfigFromFile()

	ipv4, errIpv4 := conf.getIpv4Addr()
	ipv6, errIpv6 := conf.getIpv4Addr()

	var dnsSelected dns.DNS
	switch conf.DNS.Name {
	case "alidns":
		dnsSelected = &dns.Alidns{}
	}
	dnsSelected.addRecord()

}
