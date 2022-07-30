package config

import (
	"testing"
)

// TestParseDomainArr 测试 parseDomainArr
func TestParseDomainArr(t *testing.T) {
	domains := []string{"mydomain.com", "test.mydomain.com", "test2.test.mydomain.com", "mydomain.com.cn",
		"test.mydomain.com.cn", "test:mydomain.com.cn",
		"test.mydomain.com?Line=oversea&RecordId=123", "test2:test.mydomain.com?Line=oversea&RecordId=123"}
	result := []Domain{
		{DomainName: "mydomain.com", SubDomain: ""},
		{DomainName: "mydomain.com", SubDomain: "test"},
		{DomainName: "mydomain.com", SubDomain: "test2.test"},
		{DomainName: "mydomain.com.cn", SubDomain: ""},
		{DomainName: "mydomain.com.cn", SubDomain: "test"},
		{DomainName: "mydomain.com.cn", SubDomain: "test"},
		{DomainName: "mydomain.com", SubDomain: "test", CustomParams: "Line=oversea&RecordId=123"},
		{DomainName: "test.mydomain.com", SubDomain: "test2", CustomParams: "Line=oversea&RecordId=123"},
	}

	parsedDomains := checkParseDomains(domains)
	for i := 0; i < len(parsedDomains); i++ {
		if parsedDomains[i].DomainName != result[i].DomainName ||
			parsedDomains[i].SubDomain != result[i].SubDomain ||
			parsedDomains[i].CustomParams != result[i].CustomParams {
			t.Error(parsedDomains[i].String() + "解析失败")
		}
	}

}
