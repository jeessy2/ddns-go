package util

import "testing"

func TestWriteString(t *testing.T) {
	tests := []struct {
		input    []string
		expected string
	}{
		{[]string{"hello", "world"}, "helloworld"},
		{[]string{"", "test"}, "test"},
		{[]string{"hello", " ", "world"}, "hello world"},
		{[]string{""}, ""},
	}

	for _, tt := range tests {
		result := WriteString(tt.input...)
		if result != tt.expected {
			t.Errorf("Expected %s, but got %s", tt.expected, result)
		}
	}
}

func TestToHostname(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"With https scheme", "https://www.example.com", "www.example.com"},
		{"With path", "www.example.com/path", "www.example.com"},
		{"With https scheme and path", "https://www.example.com/path", "www.example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toHostname(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, but got %s", tt.expected, result)
			}
		})
	}
}
