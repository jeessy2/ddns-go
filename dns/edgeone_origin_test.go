package dns

import (
	"testing"

	"github.com/jeessy2/ddns-go/v6/config"
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
