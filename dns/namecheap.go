package dns

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/jeessy2/ddns-go/v4/config"
	"github.com/jeessy2/ddns-go/v4/util"
)

const (
	nameCheapEndpoint string = "https://dynamicdns.park-your-domain.com/update?host=#{host}&domain=#{domain}&password=#{serectKey}&ip=#{ip}"
)

// NameCheap Domain
type NameCheap struct {
	DNSConfig config.DNSConfig
	Domains   config.Domains
}

// NameCheap 修改域名解析结果
type NameCheapResp struct {
	Status string
	Errors []string
}

// Init 初始化
func (nc *NameCheap) Init(conf *config.Config) {
	nc.DNSConfig = conf.DNS
	nc.Domains.GetNewIp(conf)
}

// AddUpdateDomainRecords 添加或更新IPv4/IPv6记录
func (nc *NameCheap) AddUpdateDomainRecords() config.Domains {
	nc.addUpdateDomainRecords("A")
	nc.addUpdateDomainRecords("AAAA")
	return nc.Domains
}

func (nc *NameCheap) addUpdateDomainRecords(recordType string) {
	ipAddr, domains := nc.Domains.GetNewIpResult(recordType)

	if ipAddr == "" {
		return
	}

	for _, domain := range domains {

		nc.modify(domain, recordType, ipAddr)
	}
}

// 修改
func (nc *NameCheap) modify(domain *config.Domain, recordType string, ipAddr string) {
	var result NameCheapResp
	err := nc.request(&result, ipAddr, domain)

	if err != nil {
		log.Printf("新增域名解析 %s 失败！", domain)
		domain.UpdateStatus = config.UpdatedFailed
		return
	}

	switch result.Status {
	case "Success":
		log.Printf("新增域名解析 %s 成功！IP: %s", domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	default:
		log.Printf("新增域名解析 %s 失败！Status: %s", domain, result.Status)
	}
}

// request 统一请求接口
func (nc *NameCheap) request(result *NameCheapResp, ipAddr string, domain *config.Domain) (err error) {
	var url string = nameCheapEndpoint
	url = strings.ReplaceAll(url, "#{host}", domain.SubDomain)
	url = strings.ReplaceAll(url, "#{domain}", domain.DomainName)
	url = strings.ReplaceAll(url, "#{serectKey}", nc.DNSConfig.Secret)
	url = strings.ReplaceAll(url, "#{ip}", ipAddr)

	log.Println("Start to request url: ", url)

	req, err := http.NewRequest(
		http.MethodGet,
		url,
		http.NoBody,
	)

	if err != nil {
		log.Println("http.NewRequest失败. Error: ", err)
		return
	}

	client := util.CreateHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		log.Println("client.Do失败. Error: ", err)
		return
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("请求namecheap失败")
		return err
	}

	status := string(data)

	if strings.Contains(status, "<ErrCount>0</ErrCount>") {
		result.Status = "Success"
	} else {
		result.Status = status
	}

	return
}
