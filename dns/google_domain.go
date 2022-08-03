package dns

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/jeessy2/ddns-go/v4/config"
	"github.com/jeessy2/ddns-go/v4/util"
)

const (
	googleDomainEndpoint string = "https://domains.google.com/nic/update"
)

// https://support.google.com/domains/answer/6147083?hl=zh-Hans#zippy=%2C使用-api-更新您的动态-dns-记录
// GoogleDomain Google Domain
type GoogleDomain struct {
	DNSConfig config.DNSConfig
	Domains   config.Domains
}

// GoogleDomainResp 修改域名解析结果
type GoogleDomainResp struct {
	Status  string
	SetedIP string
}

// Init 初始化
func (gd *GoogleDomain) Init(conf *config.Config) {
	gd.DNSConfig = conf.DNS
	gd.Domains.GetNewIp(conf)
}

// AddUpdateDomainRecords 添加或更新IPv4/IPv6记录
func (gd *GoogleDomain) AddUpdateDomainRecords() config.Domains {
	gd.addUpdateDomainRecords("A")
	gd.addUpdateDomainRecords("AAAA")
	return gd.Domains
}

func (gd *GoogleDomain) addUpdateDomainRecords(recordType string) {
	ipAddr, domains := gd.Domains.GetNewIpResult(recordType)

	if ipAddr == "" {
		return
	}

	for _, domain := range domains {

		gd.modify(domain, recordType, ipAddr)
	}
}

// 修改
func (gd *GoogleDomain) modify(domain *config.Domain, recordType string, ipAddr string) {
	params := domain.GetCustomParams()
	params.Set("hostname", domain.GetFullDomain())
	params.Set("myip", ipAddr)

	var result GoogleDomainResp
	err := gd.request(params, &result)

	if err != nil {
		log.Printf("新增域名解析 %s 失败！", domain)
		domain.UpdateStatus = config.UpdatedFailed
		return
	}

	switch result.Status {
	case "nochg":
		log.Printf("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
	case "good":
		log.Printf("新增域名解析 %s 成功！IP: %s", domain, ipAddr)
		domain.UpdateStatus = config.UpdatedSuccess
	default:
		log.Printf("新增域名解析 %s 失败！Status: %s", domain, result.Status)
	}
}

// request 统一请求接口
func (gd *GoogleDomain) request(params url.Values, result *GoogleDomainResp) (err error) {

	req, err := http.NewRequest(
		http.MethodPost,
		googleDomainEndpoint,
		http.NoBody,
	)

	if err != nil {
		log.Println("http.NewRequest失败. Error: ", err)
		return
	}

	req.URL.RawQuery = params.Encode()
	req.SetBasicAuth(gd.DNSConfig.ID, gd.DNSConfig.Secret)

	client := util.CreateHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		log.Println("client.Do失败. Error: ", err)
		return
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	status := string(data)

	if s := strings.Split(status, " "); s[0] == "good" || s[0] == "nochg" { // Success status
		result.Status = s[0]
		if len(s) > 1 {
			result.SetedIP = s[1]
		}
	} else {
		result.Status = status
	}
	return
}
