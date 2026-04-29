package dns

import (
	"testing"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

func TestEdgeOneIsOriginGroupDomain(t *testing.T) {
	eo := &EdgeOne{}
	if !eo.isOriginGroupDomain(&config.Domain{
		DomainName:   "example.com",
		CustomParams: "GroupId=origin-123",
	}) {
		t.Fatal("expected GroupId to enable origin group mode")
	}
	if !eo.isOriginGroupDomain(&config.Domain{
		DomainName:   "example.com",
		CustomParams: "OriginGroupName=my-group",
	}) {
		t.Fatal("expected OriginGroupName to enable origin group mode")
	}
	if eo.isOriginGroupDomain(&config.Domain{DomainName: "example.com"}) {
		t.Fatal("did not expect plain DNS record to enter origin group mode")
	}
}

func TestSameEdgeOneOriginRecords(t *testing.T) {
	current := []EdgeOneOriginRecord{
		{Record: "2001:db8::1", Type: edgeoneOriginRecordType, Weight: 100},
		{Record: "1.1.1.1", Type: edgeoneOriginRecordType, Weight: 100},
	}
	desired := []EdgeOneOriginRecord{
		{Record: "1.1.1.1", Type: edgeoneOriginRecordType, Weight: 100},
		{Record: "2001:db8::1", Type: edgeoneOriginRecordType, Weight: 100},
	}

	if !sameEdgeOneOriginRecords(current, desired) {
		t.Fatal("expected records with different order to be treated as identical")
	}

	desired[1].Weight = 90
	if sameEdgeOneOriginRecords(current, desired) {
		t.Fatal("expected weight difference to be treated as a change")
	}
}

func TestBuildEdgeOneDomainTuplesKeepsCombinedAddresses(t *testing.T) {
	domain := &config.Domain{
		DomainName:   "example.com",
		CustomParams: "GroupId=origin-123",
	}
	results := buildEdgeOneDomainTuples(
		config.Domains{
			Ipv4Addr: "1.1.1.1",
			Ipv6Addr: "2001:db8::1",
		},
		"1.1.1.1",
		[]*config.Domain{domain},
		"",
		nil,
	)
	if len(results) != 1 {
		t.Fatalf("expected 1 tuple, got %d", len(results))
	}

	tuple := results[domain.String()]
	if tuple == nil {
		t.Fatal("expected tuple to exist")
	}
	if tuple.RecordType != "A" {
		t.Fatalf("expected record type A, got %s", tuple.RecordType)
	}
	if tuple.Ipv4Addr != "1.1.1.1" || tuple.Ipv6Addr != "2001:db8::1" {
		t.Fatalf("expected both current addresses to be preserved, got ipv4=%s ipv6=%s", tuple.Ipv4Addr, tuple.Ipv6Addr)
	}
}

func TestIpCacheSingleRoundReuse(t *testing.T) {
	cache := util.IpCache{}
	if !cache.Check("1.1.1.1") {
		t.Fatal("expected first check to trigger compare")
	}
	if cache.Check("1.1.1.1") {
		t.Fatal("expected second check in same state to be skipped")
	}
}
