package dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

const nsoneAPIEndpoint = "https://api.nsone.net/v1/zones"

type NSOne struct {
	DNS     config.DNS
	Domains config.Domains
	TTL     int
}

type NSOneZone struct {
	AssignedNameservers []string `json:"assigned_nameservers"`
	DNSServers          []string `json:"dns_servers"`
	Expiry              int      `json:"expiry"`
	Name                string   `json:"name"`
	Link                string   `json:"link"`
	PrimaryMaster       string   `json:"primary_master"`
	Hostmaster          string   `json:"hostmaster"`
	ID                  string   `json:"id"`
	Meta                struct {
		Asn           []string `json:"asn"`
		CaProvince    []string `json:"ca_province"`
		Connections   int      `json:"connections"`
		Country       []string `json:"country"`
		Georegion     []string `json:"georegion"`
		HighWatermark float64  `json:"high_watermark"`
		IPPrefixes    []string `json:"ip_prefixes"`
		Latitude      float64  `json:"latitude"`
		LoadAvg       float64  `json:"loadAvg"`
		Longitude     float64  `json:"longitude"`
		LowWatermark  float64  `json:"low_watermark"`
		Note          string   `json:"note"`
		Priority      int      `json:"priority"`
		Pulsar        string   `json:"pulsar"`
		Requests      int      `json:"requests"`
		Up            bool     `json:"up"`
		UsState       []string `json:"us_state"`
		Weight        float64  `json:"weight"`
	} `json:"meta"`
	NetworkPools []string `json:"network_pools"`
	Networks     []int    `json:"networks"`
	NxTTL        int      `json:"nx_ttl"`
	Serial       int      `json:"serial"`
	Primary      struct {
		Enabled     bool `json:"enabled"`
		Secondaries []struct {
			IP      string `json:"ip"`
			Network int    `json:"network"`
			Notify  bool   `json:"notify"`
			Port    int    `json:"port"`
			Tsig    struct {
				Enabled bool   `json:"enabled"`
				Hash    string `json:"hash"`
				Name    string `json:"name"`
				Key     string `json:"key"`
			} `json:"tsig"`
		} `json:"secondaries"`
	} `json:"primary"`
	Refresh   int `json:"refresh"`
	Retry     int `json:"retry"`
	Secondary struct {
		Status         string `json:"status"`
		Error          string `json:"error"`
		LastXfr        int    `json:"last_xfr"`
		LastTry        int    `json:"last_try"`
		Enabled        bool   `json:"enabled"`
		Expired        bool   `json:"expired"`
		PrimaryIP      string `json:"primary_ip"`
		PrimaryPort    int    `json:"primary_port"`
		PrimaryNetwork int    `json:"primary_network"`
		Tsig           struct {
			Enabled        bool   `json:"enabled"`
			Hash           string `json:"hash"`
			Name           string `json:"name"`
			Key            string `json:"key"`
			SignedNotifies bool   `json:"signed_notifies"`
		} `json:"tsig"`
		OtherPorts      []int    `json:"other_ports"`
		OtherIps        []string `json:"other_ips"`
		OtherNetworks   []int    `json:"other_networks"`
		OtherNotifyOnly []bool   `json:"other_notify_only"`
	} `json:"secondary"`
	TTL       int      `json:"ttl"`
	Zone      string   `json:"zone"`
	Views     []string `json:"views"`
	LocalTags []string `json:"local_tags"`
	Tags      struct {
		ID int64 `json:"id"`
	} `json:"tags"`
	CreatedAt  int  `json:"created_at"`
	UpdatedAt  int  `json:"updated_at"`
	Dnssec     bool `json:"dnssec"`
	Signatures []struct {
		Answer []string `json:"answer"`
	} `json:"signatures"`
	Presigned     bool `json:"presigned"`
	IDVersion     int  `json:"id_version"`
	ActiveVersion bool `json:"active_version"`
}

type NSOneRecordAnswer struct {
	Answer []string `json:"answer"`
	ID     string   `json:"id,omitempty"`
	Meta   struct {
		ID int64 `json:"id,omitempty"`
	} `json:"meta,omitempty"`
	Region string `json:"region,omitempty"`
	Feeds  []struct {
		Source string `json:"source,omitempty"`
		Feed   string `json:"feed,omitempty"`
	} `json:"feeds,omitempty"`
}

type NSOneRecordResponse struct {
	Answers []NSOneRecordAnswer `json:"answers"`
	Domain  string              `json:"domain"`
	Filters []struct {
		Config struct {
			Eliminate bool `json:"eliminate"`
		} `json:"config"`
	} `json:"filters"`
	Link string `json:"link"`
	Meta struct {
		ID int64 `json:"id"`
	} `json:"meta"`
	Networks []int `json:"networks"`
	Regions  struct {
		ID int64 `json:"id"`
	} `json:"regions"`
	Tier            int      `json:"tier"`
	TTL             int      `json:"ttl"`
	OverrideTTL     bool     `json:"override_ttl"`
	Type            string   `json:"type"`
	UseClientSubnet bool     `json:"use_client_subnet"`
	Zone            string   `json:"zone"`
	ZoneName        string   `json:"zone_name"`
	BlockedTags     []string `json:"blocked_tags"`
	LocalTags       []string `json:"local_tags"`
	Tags            struct {
		ID int64 `json:"id"`
	} `json:"tags"`
	OverrideAddressRecords bool `json:"override_address_records"`
	Signatures             []struct {
		Answer []string `json:"answer"`
	} `json:"signatures"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
	ID        string `json:"id"`
	Customer  int    `json:"customer"`
	Feeds     []struct {
		Source string `json:"source"`
		Feed   string `json:"feed"`
	} `json:"feeds"`
}

type NSOneRecordRequest struct {
	Answers []NSOneRecordAnswer `json:"answers"`
	Domain  string              `json:"domain"`
	TTL     int                 `json:"ttl"`
	Type    string              `json:"type"`
	Zone    string              `json:"zone"`
}

func (nsone *NSOne) Init(dnsConf *config.DnsConfig, ipv4cache *util.IpCache, ipv6cache *util.IpCache) {
	nsone.Domains.Ipv4Cache = ipv4cache
	nsone.Domains.Ipv6Cache = ipv6cache
	nsone.DNS = dnsConf.DNS
	nsone.Domains.GetNewIp(dnsConf)
	if dnsConf.TTL == "" {
		nsone.TTL = 60
	} else {
		ttl, err := strconv.Atoi(dnsConf.TTL)
		if err != nil {
			// Default TTL in documentation is 1 hour
			nsone.TTL = 3600
		} else {
			nsone.TTL = ttl
		}
	}
}

func (nsone *NSOne) AddUpdateDomainRecords() config.Domains {
	nsone.addUpdateDomainRecords("A")
	nsone.addUpdateDomainRecords("AAAA")
	return nsone.Domains
}

func (nsone *NSOne) addUpdateDomainRecords(recordType string) {
	ipAddr, domains := nsone.Domains.GetNewIpResult(recordType)

	if ipAddr == "" {
		return
	}

	for _, domain := range domains {
		zoneInfo, err := nsone.getZone(domain)
		if err != nil {
			util.Log("查询域名信息发生异常! %s", err)
			domain.UpdateStatus = config.UpdatedFailed
			continue
		}

		if zoneInfo == nil {
			util.Log("在DNS服务商中未找到根域名: %s", domain.DomainName)
			domain.UpdateStatus = config.UpdatedFailed
			continue
		}

		existingRecord, err := nsone.getRecord(domain, recordType)
		if err != nil {
			util.Log("查询域名信息发生异常! %s", err)
			domain.UpdateStatus = config.UpdatedFailed
			continue
		}

		if existingRecord != nil {
			nsone.updateRecord(domain, recordType, ipAddr, existingRecord)
		} else {
			nsone.createRecord(domain, recordType, ipAddr)
		}
	}
}

func (nsone *NSOne) getZone(domain *config.Domain) (*NSOneZone, error) {
	var result NSOneZone
	params := url.Values{}
	params.Set("records", "false")

	err := nsone.request(
		"GET",
		fmt.Sprintf("%s/%s?%s", nsoneAPIEndpoint, domain.DomainName, params.Encode()),
		nil,
		&result,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (nsone *NSOne) getRecord(domain *config.Domain, recordType string) (*NSOneRecordResponse, error) {
	var result NSOneRecordResponse
	params := url.Values{}
	params.Set("records", "false")

	err := nsone.request(
		"GET",
		fmt.Sprintf("%s/%s/%s/%s?%s", nsoneAPIEndpoint, domain.DomainName, domain.GetFullDomain(), recordType, params.Encode()),
		nil,
		&result,
	)

	if err == nil && len(result.Answers) > 0 {
		return &result, nil
	}

	return nil, nil
}

func (nsone *NSOne) createRecord(domain *config.Domain, recordType string, ipAddr string) {
	recordName := domain.GetFullDomain()
	request := NSOneRecordRequest{
		Answers: []NSOneRecordAnswer{
			{
				Answer: []string{
					ipAddr,
				},
			},
		},
		Domain: recordName,
		TTL:    nsone.TTL,
		Type:   recordType,
		Zone:   domain.DomainName,
	}

	var response NSOneRecordResponse
	err := nsone.request(
		"PUT",
		fmt.Sprintf("%s/%s/%s/%s", nsoneAPIEndpoint, domain.DomainName, recordName, recordType),
		request,
		&response,
	)

	if err != nil {
		util.Log("新增域名解析 %s 失败! 异常信息: %s", domain, err)
		domain.UpdateStatus = config.UpdatedFailed
		return
	}

	util.Log("新增域名解析 %s 成功! IP: %s", domain, ipAddr)
	domain.UpdateStatus = config.UpdatedSuccess
}

func (nsone *NSOne) updateRecord(domain *config.Domain, recordType string, ipAddr string, existingRecord *NSOneRecordResponse) {
	if len(existingRecord.Answers) > 0 && len(existingRecord.Answers[0].Answer) > 0 {
		if existingRecord.Answers[0].Answer[0] == ipAddr {
			util.Log("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
			return
		}
	}

	recordName := domain.GetFullDomain()
	request := NSOneRecordRequest{
		Answers: []NSOneRecordAnswer{
			{
				Answer: []string{
					ipAddr,
				},
			},
		},
		Domain: recordName,
		TTL:    nsone.TTL,
		Type:   recordType,
		Zone:   domain.DomainName,
	}

	var response NSOneRecordResponse
	err := nsone.request(
		"POST",
		fmt.Sprintf("%s/%s/%s/%s", nsoneAPIEndpoint, domain.DomainName, recordName, recordType),
		request,
		&response,
	)

	if err != nil {
		util.Log("更新域名解析 %s 失败! 异常信息: %s", domain, err)
		domain.UpdateStatus = config.UpdatedFailed
		return
	}

	util.Log("更新域名解析 %s 成功! IP: %s", domain, ipAddr)
	domain.UpdateStatus = config.UpdatedSuccess
}

func (nsone *NSOne) request(method string, url string, data interface{}, result interface{}) (err error) {
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
		return
	}

	req.Header.Set("X-NSONE-Key", nsone.DNS.Secret)
	req.Header.Set("Content-Type", "application/json")

	client := util.CreateHTTPClient()
	resp, err := client.Do(req)
	err = util.GetHTTPResponse(resp, err, result)

	return
}
