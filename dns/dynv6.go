package dns

import (
	"bytes"
	"encoding/json"
	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
	"net/http"
	"strconv"
	"strings"
)

const (
	dynv6Endpoint = "https://dynv6.com"
)

type Dynv6 struct {
	DNS     config.DNS
	Domains config.Domains
	TTL     string
}

type Dynv6Zone struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Ipv4 string `json:"ipv4address"`
	Ipv6 string `json:"ipv6prefix"`
}

type Dynv6Record struct {
	ID     uint   `json:"id"`
	ZoneID uint   `json:"zoneID"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Data   string `json:"data"`
}

// Init 初始化
func (dynv6 *Dynv6) Init(dnsConf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache) {
	dynv6.Domains.Ipv4Cache = ipv4cache
	dynv6.Domains.Ipv6Cache = ipv6cache
	dynv6.DNS = dnsConf.DNS
	dynv6.Domains.GetNewIp(dnsConf)
	if dnsConf.TTL == "" {
		// 默认600s
		dynv6.TTL = "600"
	} else {
		dynv6.TTL = dnsConf.TTL
	}
}

// AddUpdateDomainRecords 添加或更新IPv4/IPv6记录
func (dynv6 *Dynv6) AddUpdateDomainRecords() config.Domains {
	dynv6.addUpdateDomainRecords("A")
	dynv6.addUpdateDomainRecords("AAAA")
	return dynv6.Domains
}

func (dynv6 *Dynv6) addUpdateDomainRecords(recordType string) {
	ipAddr, domains := dynv6.Domains.GetNewIpResult(recordType)

	if ipAddr == "" {
		return
	}

	for _, domain := range domains {
		isFindZone, findZone, isMain, err := dynv6.findZone(domain)

		if err != nil {
			util.Log("查询域名信息发生异常! %s", err)
			domain.UpdateStatus = config.UpdatedFailed
			return
		}

		if !isFindZone {
			util.Log("在DNS服务商中未找到根域名: %s", domain)
			domain.UpdateStatus = config.UpdatedFailed
			continue
		}

		zoneId := strconv.FormatUint(uint64(findZone.ID), 10)

		if isMain {
			// 如果使用的域名是主域名，对比DNS记录确定是否调用更新接口
			if (recordType == "A" && findZone.Ipv4 == ipAddr) || (recordType == "AAAA" && findZone.Ipv6 == ipAddr) {
				// ip与dns服务器一致，不执行更新
				util.Log("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
				domain.UpdateStatus = config.UpdatedNothing
			} else {
				dynv6.modifyMain(domain, zoneId, recordType, ipAddr)
			}
		} else {
			// 如果是子域名，检查是否有该子域名记录，有就更新记录，没有就创建

			// 处理subDomain
			processSubDomainOk := dynv6.processSubDomain(domain, findZone)

			if !processSubDomainOk {
				util.Log("域名: %s 不正确", domain)
				domain.UpdateStatus = config.UpdatedFailed
				continue
			}

			isFindRecord, findRecord, err := dynv6.findRecord(domain, zoneId, recordType)

			if err != nil {
				util.Log("查询域名信息发生异常! %s", err)
				domain.UpdateStatus = config.UpdatedFailed
				return
			}

			if isFindRecord {
				// 判断是否需要更新
				if findRecord.Type == recordType && findRecord.Data == ipAddr {
					// ip与dns服务器一致，不执行更新
					util.Log("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
					domain.UpdateStatus = config.UpdatedNothing
				} else {
					dynv6.modify(domain, zoneId, findRecord, recordType, ipAddr)
				}
			} else {
				// 创建记录
				dynv6.create(domain, zoneId, recordType, ipAddr)
			}
		}
	}
}

func (dynv6 *Dynv6) processSubDomain(domain *config.Domain, zone Dynv6Zone) bool {
	// 确定subDomain
	subDomainLen := len(domain.String()) - len(zone.Name) - 1
	if subDomainLen <= 0 {
		return false
	}
	subDomain := domain.String()[:subDomainLen]

	domain.DomainName = zone.Name
	domain.SubDomain = subDomain
	return true
}

// 根据domain获取zone
func (dynv6 *Dynv6) findZone(domain *config.Domain) (isFind bool, zone Dynv6Zone, isMain bool, err error) {
	var zones []Dynv6Zone
	isFind = false
	isMain = false

	// 获取所有zone
	err = dynv6.request("GET", dynv6Endpoint+"/api/v2/zones", nil, &zones)

	if err != nil {
		return
	}

	// 遍历token权限下所有zone，确定当前域名属于哪个zone，并判断当前域名是主域名还是子域名
	for _, z := range zones {
		if strings.HasSuffix(domain.String(), z.Name) {
			isFind = true
			zone = z
			if domain.String() == z.Name {
				isMain = true
			}
			break
		}
	}

	return
}

// 根据domain获取record
func (dynv6 *Dynv6) findRecord(domain *config.Domain, zoneId string, recordType string) (isFind bool, record Dynv6Record, err error) {
	var records []Dynv6Record
	isFind = false

	err = dynv6.request("GET", dynv6Endpoint+"/api/v2/zones/"+zoneId+"/records", nil, &records)
	if err != nil {
		return
	}

	// 遍历zone下所有record，判断是更新还是创建
	for _, r := range records {
		if r.Name == domain.SubDomain && r.Type == recordType {
			isFind = true
			record = r
			break
		}
	}

	return
}

// modify 更新根域名
func (dynv6 *Dynv6) modifyMain(domain *config.Domain, zoneId string, recordType string, ipAddr string) {
	var zoneUpdateReq = Dynv6Zone{}
	if recordType == "A" {
		zoneUpdateReq.Ipv4 = ipAddr
	} else {
		zoneUpdateReq.Ipv6 = ipAddr
	}

	err := dynv6.request("PATCH", dynv6Endpoint+"/api/v2/zones/"+zoneId, zoneUpdateReq, &Dynv6Zone{})

	if err != nil {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, err)
		domain.UpdateStatus = config.UpdatedFailed
	} else {
		util.Log("更新域名解析 %s 成功! IP: %s", domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	}
}

// create 创建新的解析
func (dynv6 *Dynv6) create(domain *config.Domain, zoneId string, recordType string, ipAddr string) {
	recordUpdateReq := Dynv6Record{
		Name: domain.SubDomain,
		Type: recordType,
		Data: ipAddr,
	}

	err := dynv6.request("POST", dynv6Endpoint+"/api/v2/zones/"+zoneId+"/records", recordUpdateReq, &Dynv6Record{})

	if err != nil {
		util.Log("新增域名解析 %s 失败! 异常信息: %s", domain, err)
		domain.UpdateStatus = config.UpdatedFailed
	} else {
		util.Log("新增域名解析 %s 成功! IP: %s", domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	}
}

// modify 更新解析
func (dynv6 *Dynv6) modify(domain *config.Domain, zoneId string, record Dynv6Record, recordType string, ipAddr string) {
	record.Type = recordType
	record.Data = ipAddr

	recordId := strconv.FormatUint(uint64(record.ID), 10)

	err := dynv6.request("PATCH", dynv6Endpoint+"/api/v2/zones/"+zoneId+"/records/"+recordId, record, &Dynv6Record{})

	if err != nil {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, err)
		domain.UpdateStatus = config.UpdatedFailed
	} else {
		util.Log("更新域名解析 %s 成功! IP: %s", domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	}
}

// request 统一请求接口
func (dynv6 *Dynv6) request(method string, url string, data interface{}, result interface{}) (err error) {
	jsonStr := make([]byte, 0)
	if data != nil {
		jsonStr, _ = json.Marshal(data)
	}

	req, err := http.NewRequest(
		method,
		url,
		bytes.NewBuffer(jsonStr),
	)

	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+dynv6.DNS.Secret)
	req.Header.Set("Content-Type", "application/json")

	client := util.CreateHTTPClient()
	resp, err := client.Do(req)
	err = util.GetHTTPResponse(resp, err, result)
	return err
}
