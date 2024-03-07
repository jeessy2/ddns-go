package util

import "testing"

const (
	testDNS = "1.1.1.1"
	testURL = "https://cloudflare.com"
)

func TestSetDNS(t *testing.T) {
	SetDNS(testDNS)

	if dialer.Resolver == nil {
		t.Error("Failed to set dialer.Resolver")
	}
}

func TestLookupHost(t *testing.T) {
	t.Run("Valid URL", func(t *testing.T) {
		if err := LookupHost(testURL); err != nil {
			t.Errorf("Expected nil error, got %v", err)
		}
	})

	t.Run("Invalid URL", func(t *testing.T) {
		if err := LookupHost("invalidurl"); err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("After SetDNS", func(t *testing.T) {
		SetDNS(testDNS)

		if err := LookupHost(testURL); err != nil {
			t.Errorf("Expected nil error, got %v", err)
		}
	})
}
