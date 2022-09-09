package dns

import (
	"github.com/jeessy2/ddns-go/v4/config"
	"net/http"
	"testing"
)

func TestGetRecords(t *testing.T) {
	dns := GoDaddyDNS{}
	dns.Init(&config.Config{
		DNS: config.DNSConfig{
			Name:   "",
			ID:     "gGpkXxwq7kfa_QMqMtjLRfzmV8dFE4kgRWe",
			Secret: "Jq7A9oL3ifDWbZAuNKpxmF",
		},
		User:              config.User{},
		Webhook:           config.Webhook{},
		NotAllowWanAccess: false,
		TTL:               "",
	})
	dns.sendReq(http.MethodPut, "A", &config.Domain{SubDomain: "@", DomainName: "furryfandom.xyz"},
		[]godaddyRecord{godaddyRecord{
			Data: "1.1.1.1",
			Name: "@",
			TTL:  600,
			Type: "A",
		}})

}
