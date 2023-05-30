package config

import (
	"testing"
)

// TestParseDomainArr 测试 parseDomainArr
func TestParseDomainArr(t *testing.T) {
	domains := []string{"mydomain.com", "test.mydomain.com", "test2.test.mydomain.com", "mydomain.com.cn",
		"test.mydomain.com.cn", "test:mydomain.com.cn",
		"test.mydomain.com?Line=oversea&RecordId=123", "test.mydomain.com.cn?Line=oversea&RecordId=123",
		"test2:test.mydomain.com?Line=oversea&RecordId=123"}
	result := []Domain{
		{DomainName: "mydomain.com", SubDomain: ""},
		{DomainName: "mydomain.com", SubDomain: "test"},
		{DomainName: "mydomain.com", SubDomain: "test2.test"},
		{DomainName: "mydomain.com.cn", SubDomain: ""},
		{DomainName: "mydomain.com.cn", SubDomain: "test"},
		{DomainName: "mydomain.com.cn", SubDomain: "test"},
		{DomainName: "mydomain.com", SubDomain: "test", CustomParams: "Line=oversea&RecordId=123"},
		{DomainName: "mydomain.com.cn", SubDomain: "test", CustomParams: "Line=oversea&RecordId=123"},
		{DomainName: "test.mydomain.com", SubDomain: "test2", CustomParams: "Line=oversea&RecordId=123"},
	}

	parsedDomains := checkParseDomains(domains)
	for i := 0; i < len(parsedDomains); i++ {
		if parsedDomains[i].DomainName != result[i].DomainName ||
			parsedDomains[i].SubDomain != result[i].SubDomain ||
			parsedDomains[i].CustomParams != result[i].CustomParams {
			t.Errorf("解析 %s 失败：\n期待 DomainName：%s，得到 DomainName：%s\n期待 SubDomain：%s，得到 SubDomain：%s\n期待 CustomParams：%s，得到 CustomParams：%s",
				parsedDomains[i].String(),
				result[i].DomainName, parsedDomains[i].DomainName,
				result[i].SubDomain, parsedDomains[i].SubDomain,
				result[i].CustomParams, parsedDomains[i].CustomParams)
		}
	}

}
