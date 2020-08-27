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

	// IPV4
	ipv4Addr, err := conf.GetIpv4Addr()
	if err == nil {
		ali.Ipv4Addr = ipv4Addr
		ali.Ipv4Domains = ParseDomain(conf.Ipv4.Domains)
	}

	// IPV6
	ipv6Addr, err := conf.GetIpv6Addr()
	if err != nil {
		ali.Ipv6Addr = ipv6Addr
		ali.Ipv6Domains = ParseDomain(conf.Ipv6.Domains)
	}

}

// AddUpdateIpv4DomainRecords 添加或更新IPV4记录
func (ali *Alidns) AddUpdateIpv4DomainRecords() {
	ali.addUpdateDomainRecords("A")
}

// AddUpdateIpv6DomainRecords 添加或更新IPV4记录
func (ali *Alidns) AddUpdateIpv6DomainRecords() {
	ali.addUpdateDomainRecords("AAAA")
}

func (ali *Alidns) addUpdateDomainRecords(typ string) {
	typeName := "ipv4"
	ipAddr := ali.Ipv4Addr
	domains := ali.Ipv4Domains
	if typ == "AAAA" {
		typeName = "ipv6"
		ipAddr = ali.Ipv6Addr
		domains = ali.Ipv6Domains
	}

	if ipAddr == "" {
		return
	}

	existReq := alidnssdk.CreateDescribeSubDomainRecordsRequest()
	existReq.Type = typ

	for _, dom := range domains {
		existReq.SubDomain = dom.SubDomain + "." + dom.DomainName
		rep, err := ali.client.DescribeSubDomainRecords(existReq)
		if err != nil {
			log.Println(err)
		}
		if rep.TotalCount > 0 {
			// Update
			if rep.DomainRecords.Record[0].Value != ipAddr {
				request := alidnssdk.CreateUpdateDomainRecordRequest()
				request.Scheme = "https"
				request.Value = ipAddr
				request.Type = typ
				request.RR = dom.SubDomain
				request.RecordId = rep.DomainRecords.Record[0].RecordId

				_, err = ali.client.UpdateDomainRecord(request)
				if err != nil {
					log.Println("Update ipv4 error! Domain:", dom, " IP:", ipAddr, " ERROR: ", err.Error())
				} else {
					log.Println("Update ipv4 success! Domain:", dom, " IP:", ipAddr)
				}
				if rep.TotalCount > 1 {
					log.Println(dom, "records more than 2, We just update the first!")
				}
			} else {
				log.Println(typeName, "address is the same！Domain:", dom, " IP:", ipAddr)
			}
		} else {
			// Add
			request := alidnssdk.CreateAddDomainRecordRequest()
			request.Scheme = "https"
			request.Value = ipAddr
			request.Type = typ
			request.RR = dom.SubDomain
			request.DomainName = dom.DomainName

			_, err = ali.client.AddDomainRecord(request)
			if err != nil {
				log.Println("Add ", typeName, " error! Domain: ", dom, " IP: ", ipAddr, " ERROR: ", err.Error())
			} else {
				log.Println("Add ", typeName, " success! Domain: ", dom, " IP: ", ipAddr)
			}
		}
	}
}
