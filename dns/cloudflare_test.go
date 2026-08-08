package dns

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/jeessy2/ddns-go/v6/config"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

// RoundTrip implements http.RoundTripper for Cloudflare request tests.
func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestCloudflareModifyProxied(t *testing.T) {
	tests := []struct {
		name          string
		customParams  string
		recordIP      string
		ipAddr        string
		recordProxied bool
		wantRequest   bool
		wantProxied   bool
	}{
		{
			name:          "preserve proxy status without custom parameter",
			recordIP:      "192.0.2.1",
			ipAddr:        "192.0.2.1",
			recordProxied: true,
			wantRequest:   false,
		},
		{
			name:          "preserve proxy status when IP changes",
			recordIP:      "192.0.2.1",
			ipAddr:        "192.0.2.2",
			recordProxied: true,
			wantRequest:   true,
			wantProxied:   true,
		},
		{
			name:          "sync explicitly disabled proxy status without IP changes",
			customParams:  "proxied=false",
			recordIP:      "192.0.2.1",
			ipAddr:        "192.0.2.1",
			recordProxied: true,
			wantRequest:   true,
			wantProxied:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestCount := 0
			var updatedRecord CloudflareRecord
			client := &http.Client{
				Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
					requestCount++
					if err := json.NewDecoder(request.Body).Decode(&updatedRecord); err != nil {
						t.Fatalf("decode request body: %v", err)
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"success":true}`)),
						Header:     make(http.Header),
					}, nil
				}),
			}

			cf := Cloudflare{TTL: 1, httpClient: client}
			domain := &config.Domain{
				DomainName:   "example.com",
				SubDomain:    "test",
				CustomParams: tt.customParams,
			}
			cf.modify(CloudflareRecordsResp{
				Result: []CloudflareRecord{{
					ID:      "record-id",
					Content: tt.recordIP,
					Proxied: tt.recordProxied,
				}},
			}, "zone-id", domain, tt.ipAddr)

			if got := requestCount > 0; got != tt.wantRequest {
				t.Fatalf("request sent = %v, want %v", got, tt.wantRequest)
			}
			if tt.wantRequest && updatedRecord.Proxied != tt.wantProxied {
				t.Errorf("updated proxied = %v, want %v", updatedRecord.Proxied, tt.wantProxied)
			}
		})
	}
}
