package config

import (
	"reflect"
	"testing"
)

// TestExtractHeaders 测试 parseHeaderArr
func TestExtractHeaders(t *testing.T) {
	input := `
a: foo
b: bar
c: foo:bar`
	expected := map[string]string{
		"a": "foo",
		"b": "bar",
		"c": "foo:bar",
	}

	parsedHeaders := extractHeaders(input)
	if !reflect.DeepEqual(parsedHeaders, expected) {
		t.Errorf("Expected %v, got %v", expected, parsedHeaders)
	}
}

// TestIpv4ToUint32 测试IPv4地址转数字格式
func TestIpv4ToUint32(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"219.239.110.138", "3689901706"},
		{"0.0.0.0", "0"},
		{"255.255.255.255", "4294967295"},
		{"1.2.3.4", "16909060"},
		{"", ""},
		{"not-an-ip", ""},
		{"2001:db8::1", ""},
	}

	for _, tt := range tests {
		if got := ipv4ToUint32(tt.input); got != tt.expected {
			t.Errorf("ipv4ToUint32(%q) = %q, expected %q", tt.input, got, tt.expected)
		}
	}
}

// TestReplacePara 测试 replacePara 中 ipv4AddrNum 的替换
func TestReplacePara(t *testing.T) {
	domains := &Domains{Ipv4Addr: "219.239.110.138", Ipv6Addr: "2001:db8::1"}

	got := replacePara(domains, "#{ipv4AddrNum}|#{ipv4Addr}|#{ipv4AddrNum}", UpdatedSuccess, UpdatedNothing, "123")
	expected := "3689901706|219.239.110.138|3689901706"
	if got != expected {
		t.Errorf("replacePara got %q, expected %q", got, expected)
	}
}
