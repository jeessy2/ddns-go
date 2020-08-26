package dns

import (
	"ddns-go/config"
	"fmt"
	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

// Alidns 阿里云dns实现
type Alidns struct{}

// AddRecord 添加记录
func (ali *Alidns) AddRecord(conf *config.Config) (ipv4 bool, ipv6 bool) {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", conf.DNS.ID, conf.DNS.Secret)
	if err != nil {
		log.Println("Alidns链接失败")
		return false, false
	}
	ipv4Success := addIpv4Record(client, conf)
	ipv6Success := addIpv6Record(client, conf)

	return ipv4Success, ipv6Success
}

func addIpv4Record(client *alidns.Client, conf *config.Config) bool {
	ipv4Addr, err := conf.GetIpv4Addr()
	if err != nil {
		log.Println("获得IPV4失败")
	}
	for _, domain := range conf.Ipv4.Domains {
		subDomain := strings.Split(domain, ".")
		if len(subDomain) >= 2 {
			reqExist := alidns.CreateDescribeDomainsRequest()
			reqExist.Domain = domain[len(subDomain[0])+1:]
			reqExist.PageSize = "500"
			reqExist.PageNumber = "1"

			rep, err := client.DescribeDomains(reqExist)
			fmt.Println(rep.Domains)

			request := alidns.CreateAddDomainRecordRequest()
			request.Scheme = "https"
			request.Value = ipv4Addr
			request.Type = "A"
			request.RR = subDomain[0]
			request.DomainName = domain[len(subDomain[0])+1:]

			_, err = client.AddDomainRecord(request)
			if err != nil {
				fmt.Print(err.Error())
				return false
			}
			log.Println("成功添加域名：", domain)
		} else {
			log.Println(domain, "不是一个域名")
			return false
		}
	}
	return true
}

func addIpv6Record(client *alidns.Client, conf *config.Config) bool {
	ipv6Addr, err := conf.GetIpv6Addr()
	if err != nil {
		log.Println("获得IPV6失败")
	}
	for _, domain := range conf.Ipv4.Domains {
		subDomain := strings.Split(domain, ".")
		if len(subDomain) >= 2 {
			request := alidns.CreateAddDomainRecordRequest()
			request.Scheme = "https"
			request.Value = ipv6Addr
			request.Type = "AAAA"
			request.RR = subDomain[0]
			request.DomainName = domain[len(subDomain[0])+1:]

			_, err := client.AddDomainRecord(request)
			if err != nil {
				fmt.Print(err.Error())
				return false
			}
			log.Println("成功添加域名：", domain)
		} else {
			log.Println(domain, "不是一个域名")
			return false
		}
	}

	return true

}
