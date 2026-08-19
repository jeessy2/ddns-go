package dns

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

func TestDesecAddUpdateDomainRecords(t *testing.T) {
	tests := []struct {
		name         string
		subDomain    string
		recordType   string
		recordIP     string
		ipAddr       string
		wantRequests int // 预期请求数(含查询)
		wantMethod   string
	}{
		{
			name:         "ip unchanged, no update request",
			subDomain:    "www",
			recordType:   "A",
			recordIP:     "192.0.2.1",
			ipAddr:       "192.0.2.1",
			wantRequests: 1,
		},
		{
			name:         "ip changed, patch existing record",
			subDomain:    "www",
			recordType:   "A",
			recordIP:     "192.0.2.1",
			ipAddr:       "192.0.2.2",
			wantRequests: 2,
			wantMethod:   "PATCH",
		},
		{
			name:         "record not found, create new record",
			subDomain:    "www",
			recordType:   "A",
			recordIP:     "",
			ipAddr:       "192.0.2.2",
			wantRequests: 2,
			wantMethod:   "POST",
		},
		{
			name:         "apex domain uses empty subname",
			subDomain:    "",
			recordType:   "AAAA",
			recordIP:     "2001:db8::1",
			ipAddr:       "2001:db8::2",
			wantRequests: 2,
			wantMethod:   "PATCH",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestCount := 0
			var lastMethod string
			var lastBody []DeSECRRSet
			client := &http.Client{
				Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
					requestCount++
					if request.Method == "GET" {
						// 校验过滤参数
						if got := request.URL.Query().Get("subname"); got != tt.subDomain {
							t.Errorf("subname filter = %q, want %q", got, tt.subDomain)
						}
						if got := request.URL.Query().Get("type"); got != tt.recordType {
							t.Errorf("type filter = %q, want %q", got, tt.recordType)
						}
						records := "[]"
						if tt.recordIP != "" {
							records = `[{"subdomain":"` + tt.subDomain + `","type":"` + tt.recordType + `","records":["` + tt.recordIP + `"],"ttl":3600}]`
						}
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(strings.NewReader(records)),
							Header:     make(http.Header),
						}, nil
					}
					lastMethod = request.Method
					// POST为单个rrset对象，PATCH为rrset数组，统一包裹为数组处理
					body, readErr := io.ReadAll(request.Body)
					if readErr != nil {
						t.Fatalf("read request body: %v", readErr)
					}
					if lastMethod == "POST" {
						var single DeSECRRSet
						if err := json.Unmarshal(body, &single); err != nil {
							t.Fatalf("decode request body: %v", err)
						}
						lastBody = []DeSECRRSet{single}
					} else if err := json.Unmarshal(body, &lastBody); err != nil {
						t.Fatalf("decode request body: %v", err)
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`[]`)),
						Header:     make(http.Header),
					}, nil
				}),
			}

			desec := DeSEC{TTL: 3600, httpClient: client}
			domain := &config.Domain{
				DomainName: "example.com",
				SubDomain:  tt.subDomain,
			}

			ipv4Cache := &util.IpCache{}
			ipv6Cache := &util.IpCache{}
			desec.Domains.Ipv4Cache = ipv4Cache
			desec.Domains.Ipv6Cache = ipv6Cache

			if tt.recordType == "A" {
				desec.Domains.Ipv4Addr = tt.ipAddr
				desec.Domains.Ipv4Domains = []*config.Domain{domain}
				// 使缓存检查通过
				ipv4Cache.Times = 1
			} else {
				desec.Domains.Ipv6Addr = tt.ipAddr
				desec.Domains.Ipv6Domains = []*config.Domain{domain}
				ipv6Cache.Times = 1
			}

			desec.addUpdateDomainRecords(tt.recordType)

			if requestCount != tt.wantRequests {
				t.Fatalf("request count = %d, want %d", requestCount, tt.wantRequests)
			}
			if tt.wantRequests == 2 {
				if lastMethod != tt.wantMethod {
					t.Errorf("method = %s, want %s", lastMethod, tt.wantMethod)
				}
				if len(lastBody) == 0 {
					t.Fatalf("request body is empty")
				}
				rrset := lastBody[0]
				if rrset.SubDomain != tt.subDomain {
					t.Errorf("subname = %q, want %q", rrset.SubDomain, tt.subDomain)
				}
				if rrset.Type != tt.recordType {
					t.Errorf("type = %s, want %s", rrset.Type, tt.recordType)
				}
				if len(rrset.Records) != 1 || rrset.Records[0] != tt.ipAddr {
					t.Errorf("records = %v, want [%s]", rrset.Records, tt.ipAddr)
				}
				if rrset.TTL != 3600 {
					t.Errorf("ttl = %d, want 3600", rrset.TTL)
				}
			}
		})
	}
}

func TestDesecZoneNotFound(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(strings.NewReader(`{"detail":"Not found."}`)),
				Header:     make(http.Header),
			}, nil
		}),
	}

	desec := DeSEC{TTL: 3600, httpClient: client}
	domain := &config.Domain{
		DomainName: "nonexistent.com",
		SubDomain:  "www",
	}
	desec.Domains.Ipv4Cache = &util.IpCache{}
	desec.Domains.Ipv6Cache = &util.IpCache{}
	desec.Domains.Ipv4Addr = "192.0.2.1"
	desec.Domains.Ipv4Domains = []*config.Domain{domain}
	desec.Domains.Ipv4Cache.Times = 1

	desec.addUpdateDomainRecords("A")

	if domain.UpdateStatus != config.UpdatedFailed {
		t.Errorf("update status = %v, want UpdatedFailed", domain.UpdateStatus)
	}
}

func TestDesecInitTTL(t *testing.T) {
	tests := []struct {
		name      string
		ttl       string
		wantValue int
	}{
		{name: "default ttl", ttl: "", wantValue: 3600},
		{name: "custom ttl", ttl: "7200", wantValue: 7200},
		{name: "below minimum clamped", ttl: "300", wantValue: 3600},
		{name: "invalid ttl falls back", ttl: "abc", wantValue: 3600},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := config.DnsConfig{}
			conf.DNS.Name = "desec"
			conf.TTL = tt.ttl
			desec := DeSEC{}
			desec.Init(&conf, &util.IpCache{}, &util.IpCache{})
			if desec.TTL != tt.wantValue {
				t.Errorf("ttl = %d, want %d", desec.TTL, tt.wantValue)
			}
		})
	}
}
