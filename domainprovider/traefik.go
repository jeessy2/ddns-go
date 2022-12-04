package domainprovider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	dac "github.com/xinsnake/go-http-digest-auth-client"

	"github.com/jeessy2/ddns-go/v4/config"
)

func init() {
	RegisterDomainProvider(&TraefikDomainProvider{})
}

type TraefikDomainProvider struct {
	config *config.Config
}

func (t *TraefikDomainProvider) Code() string {
	return "traefik"
}

func (t *TraefikDomainProvider) Init(conf *config.Config) error {
	t.config = conf
	return nil
}

func (t *TraefikDomainProvider) GetDomains() []*config.Domain {
	var domains []*config.Domain
	if !t.config.Traefik.Enable {
		return domains
	}
	url := fmt.Sprintf("%s://%s/api/http/routers", t.config.Traefik.Schema, t.config.Traefik.Host)
	var routers []struct {
		Rule string `json:"rule"`
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("获取traefik路由失败: %s", err)
		return domains
	}
	var resp *http.Response
	if t.config.Traefik.BasicAuth {
		t := dac.NewTransport(t.config.Traefik.Username, t.config.Traefik.Password)
		resp, err = t.RoundTrip(req)
		if err != nil {
			log.Fatalln(err)
			return domains
		}
	} else {
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalln(err)
			return domains
		}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("获取traefik路由失败: %s", err)
		return domains
	}
	log.Printf("获取traefik路由成功,%s", string(body))
	if err := json.Unmarshal(body, &routers); err != nil {
		log.Printf("解析traefik路由失败: %s", err)
		return domains
	}
	var domainArr []string
	for _, router := range routers {
		if strings.HasPrefix(router.Rule, "Host(`") {
			domainStr := strings.TrimPrefix(router.Rule, "Host(`")
			domainStr = strings.TrimSuffix(domainStr, "`)")
			domainArr = append(domainArr, domainStr)
		}
	}
	domains = append(domains, checkParseDomains(domainArr)...)
	return domains
}
