package dns

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

func TestGoDaddyInit(t *testing.T) {
	conf := config.DnsConfig{}
	conf.DNS.Name = "godaddy"
	conf.DNS.Secret = "gd_testpat123"
	conf.TTL = "600"
	g := GoDaddyDNS{}
	g.Init(&conf, &util.IpCache{}, &util.IpCache{})
	if got := g.header.Get("Authorization"); got != "Bearer gd_testpat123" {
		t.Errorf("Authorization = %q, want %q", got, "Bearer gd_testpat123")
	}
	if g.ttl != 600 {
		t.Errorf("ttl = %d, want 600", g.ttl)
	}
}

func TestGoDaddyUpdateDomainRecord(t *testing.T) {
	tests := []struct {
		name         string
		subDomain    string
		recordType   string
		recordIP     string // "" means GET returns an empty list
		getStatus    int
		getBareArray bool // GET returns a bare [...] instead of {"items": [...]}
		deleteStatus int
		ipAddr       string
		wantRequests int
		wantMethods  []string
		wantSuccess  bool
	}{
		{
			name:         "ip unchanged, no write",
			subDomain:    "home",
			recordType:   "A",
			recordIP:     "192.0.2.1",
			getStatus:    http.StatusOK,
			deleteStatus: http.StatusNoContent,
			ipAddr:       "192.0.2.1",
			wantRequests: 1,
			wantMethods:  []string{"GET"},
			wantSuccess:  true,
		},
		{
			name:         "ip changed, delete then create",
			subDomain:    "home",
			recordType:   "A",
			recordIP:     "192.0.2.1",
			getStatus:    http.StatusOK,
			deleteStatus: http.StatusNoContent,
			ipAddr:       "192.0.2.2",
			wantRequests: 3,
			wantMethods:  []string{"GET", "DELETE", "POST"},
			wantSuccess:  true,
		},
		{
			name:         "record not found, create only",
			subDomain:    "home",
			recordType:   "A",
			recordIP:     "",
			getStatus:    http.StatusNotFound,
			deleteStatus: http.StatusNoContent,
			ipAddr:       "192.0.2.2",
			wantRequests: 2,
			wantMethods:  []string{"GET", "POST"},
			wantSuccess:  true,
		},
		{
			name:         "delete returns 404, still creates",
			subDomain:    "home",
			recordType:   "A",
			recordIP:     "192.0.2.1",
			getStatus:    http.StatusOK,
			deleteStatus: http.StatusNotFound,
			ipAddr:       "192.0.2.2",
			wantRequests: 3,
			wantMethods:  []string{"GET", "DELETE", "POST"},
			wantSuccess:  true,
		},
		{
			name:         "apex domain uses @ name",
			subDomain:    "",
			recordType:   "AAAA",
			recordIP:     "2001:db8::1",
			getStatus:    http.StatusOK,
			deleteStatus: http.StatusNoContent,
			ipAddr:       "2001:db8::2",
			wantRequests: 3,
			wantMethods:  []string{"GET", "DELETE", "POST"},
			wantSuccess:  true,
		},
		{
			name:         "get returns bare array, still works",
			subDomain:    "home",
			recordType:   "A",
			recordIP:     "192.0.2.1",
			getStatus:    http.StatusOK,
			getBareArray: true,
			deleteStatus: http.StatusNoContent,
			ipAddr:       "192.0.2.2",
			wantRequests: 3,
			wantMethods:  []string{"GET", "DELETE", "POST"},
			wantSuccess:  true,
		},
		{
			name:         "get returns 401, marks failed",
			subDomain:    "home",
			recordType:   "A",
			recordIP:     "",
			getStatus:    http.StatusUnauthorized,
			deleteStatus: http.StatusNoContent,
			ipAddr:       "192.0.2.2",
			wantRequests: 1,
			wantMethods:  []string{"GET"},
			wantSuccess:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var methods []string
			var paths []string
			var queries []string
			var postBody []byte
			client := &http.Client{
				Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
					// Every request must hit the v3 API with a Bearer token, never v1/sso-key.
					if got := request.URL.Scheme + "://" + request.URL.Host; got != "https://api.godaddy.com" {
						t.Errorf("request host = %q, want https://api.godaddy.com", got)
					}
					if !strings.HasPrefix(request.URL.Path, "/v3/domains/zones/example.com/dns-records") {
						t.Errorf("request path = %q, want v3 zone path", request.URL.Path)
					}
					if got := request.Header.Get("Authorization"); got != "Bearer gd_testpat" {
						t.Errorf("Authorization = %q, want Bearer gd_testpat", got)
					}
					methods = append(methods, request.Method)
					paths = append(paths, request.URL.Path)
					queries = append(queries, request.URL.RawQuery)

					switch request.Method {
					case "GET":
						if tt.getStatus != http.StatusOK {
							return &http.Response{
								StatusCode: tt.getStatus,
								Body:       io.NopCloser(strings.NewReader(`{"code":"ERROR"}`)),
								Header:     make(http.Header),
							}, nil
						}
						name := tt.subDomain
						if name == "" {
							name = "@"
						}
						if tt.recordIP == "" {
							return &http.Response{
								StatusCode: http.StatusOK,
								Body:       io.NopCloser(strings.NewReader(`{"items":[]}`)),
								Header:     make(http.Header),
							}, nil
						}
						item := `{"recordId":"rec-123","data":"` + tt.recordIP + `","name":"` + name + `","ttl":600,"type":"` + tt.recordType + `"}`
						body := `{"items":[` + item + `]}`
						if tt.getBareArray {
							body = `[` + item + `]`
						}
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(strings.NewReader(body)),
							Header:     make(http.Header),
						}, nil
					case "DELETE":
						return &http.Response{
							StatusCode: tt.deleteStatus,
							Body:       io.NopCloser(strings.NewReader("")),
							Header:     make(http.Header),
						}, nil
					case "POST":
						body, readErr := io.ReadAll(request.Body)
						if readErr != nil {
							t.Fatalf("read request body: %v", readErr)
						}
						postBody = body
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(strings.NewReader(`{"recordId":"rec-456"}`)),
							Header:     make(http.Header),
						}, nil
					default:
						t.Fatalf("unexpected method %s (a v3 PUT must never be sent)", request.Method)
						return nil, nil
					}
				}),
			}

			g := GoDaddyDNS{
				ttl: 600,
				header: http.Header{
					"Authorization": {"Bearer gd_testpat"},
					"Content-Type":  {"application/json"},
				},
				client: client,
			}
			domain := &config.Domain{
				DomainName: "example.com",
				SubDomain:  tt.subDomain,
			}
			g.updateDomainRecord(tt.recordType, tt.ipAddr, []*config.Domain{domain})

			if len(methods) != tt.wantRequests {
				t.Fatalf("request count = %d (%v), want %d (%v)", len(methods), methods, tt.wantRequests, tt.wantMethods)
			}
			for i, want := range tt.wantMethods {
				if methods[i] != want {
					t.Errorf("request %d method = %s, want %s", i, methods[i], want)
				}
			}
			// The record name in the query is the subdomain, or @ for the apex.
			wantName := tt.subDomain
			if wantName == "" {
				wantName = "@"
			}
			// GET lists by query parameters on the collection URL.
			if len(paths) > 0 {
				wantPath := "/v3/domains/zones/example.com/dns-records"
				if paths[0] != wantPath {
					t.Errorf("GET path = %q, want %q", paths[0], wantPath)
				}
				if len(queries) > 0 {
					q := queries[0]
					if !strings.Contains(q, "type="+tt.recordType) {
						t.Errorf("GET query = %q, want type=%s", q, tt.recordType)
					}
					if !strings.Contains(q, "name="+url.QueryEscape(wantName)) {
						t.Errorf("GET query = %q, want name=%s", q, wantName)
					}
				}
			}
			// DELETE targets the record by its server-assigned recordId.
			for i, m := range methods {
				if m == "DELETE" {
					wantDelete := "/v3/domains/zones/example.com/dns-records/rec-123"
					if paths[i] != wantDelete {
						t.Errorf("DELETE path = %q, want %q", paths[i], wantDelete)
					}
				}
			}
			if tt.wantSuccess {
				if domain.UpdateStatus != config.UpdatedSuccess {
					t.Errorf("update status = %v, want UpdatedSuccess", domain.UpdateStatus)
				}
			} else {
				if domain.UpdateStatus != config.UpdatedFailed {
					t.Errorf("update status = %v, want UpdatedFailed", domain.UpdateStatus)
				}
			}
			// v3 POST creates ONE record: the body is a single object,
			// not an array, and carries no recordId.
			if len(tt.wantMethods) > 0 && tt.wantMethods[len(tt.wantMethods)-1] == "POST" {
				var rec godaddyRecord
				if err := json.Unmarshal(postBody, &rec); err != nil {
					t.Fatalf("decode POST body: %v", err)
				}
				if strings.HasPrefix(strings.TrimSpace(string(postBody)), "[") {
					t.Errorf("POST body is an array %s, want a single object", string(postBody))
				}
				if rec.Data != tt.ipAddr {
					t.Errorf("POST data = %q, want %q", rec.Data, tt.ipAddr)
				}
				if rec.Name != wantName {
					t.Errorf("POST name = %q, want %q", rec.Name, wantName)
				}
				if rec.Type != tt.recordType {
					t.Errorf("POST type = %q, want %q", rec.Type, tt.recordType)
				}
				if rec.TTL != 600 {
					t.Errorf("POST ttl = %d, want 600", rec.TTL)
				}
				if rec.RecordID != "" {
					t.Errorf("POST recordId = %q, want empty (server-assigned)", rec.RecordID)
				}
			}
		})
	}
}
