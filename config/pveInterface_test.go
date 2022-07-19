package config

import (
	"testing"
)

func TestGetPveInterface(t *testing.T) {
	ipv4PveInterfaces, ipv6PveInterfaces, err := GetPveInterface("100")
	if err != nil {
		t.Error(err)
	}
	t.Log(ipv4PveInterfaces, ipv6PveInterfaces)
}
