package dns

import (
	"ddns-go/config"
	"log"

	alidnssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

// Alidns 阿里云dns实现
type Alidns struct {
	client *alidnssdk.Client
	Domains
}

// Init 初始化
func (ali *Alidns) Init(conf *config.Config) {
	client, err := alidnssdk.NewClientWithAccessKey("cn-hangzhou", conf.DNS.ID, conf.DNS.Secret)
	if err != nil {
		log.Println("Alidns链接失败")
	}
	ali.client = client

	ali.Domains.ParseDomain(conf)

}

// AddUpdateIpv4DomainRecords 添加或更新IPV4记录
func (ali *Alidns) AddUpdateIpv4DomainRecords() {
	ali.addUpdateDomainRecords("A")
}

// AddUpdateIpv6DomainRecords 添加或更新IPV6记录
func (ali *Alidns) AddUpdateIpv6DomainRecords() {
	ali.addUpdateDomainRecords("AAAA")
}

func (ali *Alidns) addUpdateDomainRecords(recordType string) {
	ipAddr := ali.Ipv4Addr
	domains := ali.Ipv4Domains
	if recordType == "AAAA" {
		ipAddr = ali.Ipv6Addr
		domains = ali.Ipv6Domains
	}

	if ipAddr == "" {
		return
	}

	existReq := alidnssdk.CreateDescribeSubDomainRecordsRequest()
	existReq.Type = recordType

	for _, domain := range domains {
		existReq.SubDomain = domain.SubDomain + "." + domain.DomainName
		rep, err := ali.client.DescribeSubDomainRecords(existReq)
		if err != nil {
			log.Println(err)
		}
		if rep.TotalCount > 0 {
			// Update
			for _, record := range rep.DomainRecords.Record {
				if record.Value == ipAddr {
					log.Printf("你的IP %s 没有变化, 未有任何操作。域名 %s", ipAddr, domain)
					continue
				}
				request := alidnssdk.CreateUpdateDomainRecordRequest()
				request.Scheme = "https"
				request.Value = ipAddr
				request.Type = recordType
				request.RR = domain.SubDomain
				request.RecordId = record.RecordId

				updateResp, err := ali.client.UpdateDomainRecord(request)
				if err != nil || !updateResp.BaseResponse.IsSuccess() {
					log.Printf("更新域名解析 %s 失败！IP: %s, Error: %s, Response: %s", domain, ipAddr, err, updateResp.GetHttpContentString())
				} else {
					log.Printf("更新域名解析 %s 成功！IP: %s", domain, ipAddr)
				}
			}
		} else {
			// Add
			request := alidnssdk.CreateAddDomainRecordRequest()
			request.Scheme = "https"
			request.Value = ipAddr
			request.Type = recordType
			request.RR = domain.SubDomain
			request.DomainName = domain.DomainName

			createResp, err := ali.client.AddDomainRecord(request)
			if err != nil || !createResp.BaseResponse.IsSuccess() {
				log.Printf("新增域名解析 %s 失败！IP: %s, Error: %s, Response: %s", domain, ipAddr, err, createResp.GetHttpContentString())
			} else {
				log.Printf("新增域名解析 %s 成功！IP: %s", domain, ipAddr)
			}
		}
	}
}
