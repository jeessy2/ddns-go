package main

import (
	"ddns-go/config"
	"ddns-go/dns"
)

const (
	ipv4Addr = "https://api-ipv4.ip.sb/ip"
	ipv6Addr = "https://api-ipv6.ip.sb/ip"
)

func main() {
	conf := &config.Config{}
	conf.GetConfigFromFile()

	var dnsSelected dns.DNS
	switch conf.DNS.Name {
	case "alidns":
		dnsSelected = &dns.Alidns{}
	}
	dnsSelected.AddRecord(conf)

}
