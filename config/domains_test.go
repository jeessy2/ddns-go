package config

import (
	"testing"
)

// TestParseDomainArr 测试 parseDomainArr
func TestParseDomainArr(t *testing.T) {
	domains := []string{"mydomain.com", "test.mydomain.com", "test2.test.mydomain.com", "mydomain.com.cn", "test.mydomain.com.cn"}
	result := []Domain{
		{DomainName: "mydomain.com", SubDomain: ""},
		{DomainName: "mydomain.com", SubDomain: "test"},
		{DomainName: "mydomain.com", SubDomain: "test2.test"},
		{DomainName: "mydomain.com.cn", SubDomain: ""},
		{DomainName: "mydomain.com.cn", SubDomain: "test"},
	}

	parsedDomains := parseDomainArr(domains)
	for i := 0; i < len(parsedDomains); i++ {
		if parsedDomains[i].DomainName != result[i].DomainName || parsedDomains[i].SubDomain != result[i].SubDomain {
			t.Error(parsedDomains[i].String() + "解析失败")
		}
	}

}
