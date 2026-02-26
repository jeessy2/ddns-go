package config

import (
	"testing"
)

func TestGetNetInterface(t *testing.T) {
	ipv4NetInterfaces, ipv6NetInterfaces, err := GetNetInterface()
	if err != nil {
		t.Error(err)
	}
	t.Log(ipv4NetInterfaces, ipv6NetInterfaces)
}
