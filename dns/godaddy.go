package dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jeessy2/ddns-go/v4/config"
	"github.com/jeessy2/ddns-go/v4/util"
	"log"
	"net/http"
	"runtime"
	"sync"
)

type godaddyRecord struct {
	Data string `json:"data"`
	Name string `json:"name"`
	TTL  int    `json:"ttl"`
	Type string `json:"type"`
}

type godaddyRecords []godaddyRecord

type recordFactory struct {
	Pool *sync.Pool
}

func (r *recordFactory) getRecords() *godaddyRecords {
	return r.Pool.Get().(*godaddyRecords)
}

func (r *recordFactory) putRecords(record *godaddyRecords) {
	r.Pool.Put(record)
}

func getRecordFactory() *recordFactory {
	return &recordFactory{Pool: &sync.Pool{New: func() any {
		return &godaddyRecords{}
	}}}
}

type GoDaddyDNS struct {
	dnsConfig config.DNSConfig
	domains   config.Domains
	ttl       string
	header    http.Header
	factory   *recordFactory
	throttle  util.Throttle
	client    *http.Client
}

func (g *GoDaddyDNS) Init(conf *config.Config) {
	g.dnsConfig = conf.DNS
	//g.domains.GetNewIp(conf)
	if conf.TTL == "" {
		// 默认600s
		g.ttl = "600"
	} else {
		g.ttl = conf.TTL
	}
	g.header = map[string][]string{
		"Authorization": {fmt.Sprintf("sso-key %s:%s", g.dnsConfig.ID, g.dnsConfig.Secret)},
		"Content-Type":  {"application/json"},
	}
	g.throttle, _ = util.GetThrottle(55)
	g.factory = getRecordFactory()
	g.client = util.CreateHTTPClient()
	log.Println("godaddy dns plugin init successful")
}

func (g *GoDaddyDNS) AddUpdateDomainRecords() (domains config.Domains) {
	panic("implements me")
}

func (g *GoDaddyDNS) sendReq(method string, rType string, domain *config.Domain, data any) (*godaddyRecords, error) {
	for !g.throttle.Try() {
		runtime.Gosched()
	}
	var body *bytes.Buffer
	if data != nil {
		if buffer, err := json.Marshal(data); err != nil {
			return nil, err
		} else {
			body = bytes.NewBuffer(buffer)
		}
	}
	path := fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/%s/%s",
		domain.DomainName, rType, domain.SubDomain)
	log.Printf("向godaddy发送请求，请求地址为%s", path)
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	req.Header = g.header
	resp, err := g.client.Do(req)
	res := g.factory.getRecords()
	err = util.GetHTTPResponse(resp, path, err, res)
	if err != nil {
		g.factory.putRecords(res)
		return nil, err
	}
	return res, nil
}
