package dns

// DNS interface
type DNS interface {
	addRecord() (ipv4 bool, ipv6 bool)
}
