package util

import "strings"

// WriteString creates a new string using [strings.Builder].
func WriteString(strs ...string) string {
	var b strings.Builder
	for _, str := range strs {
		b.WriteString(str)
	}

	return b.String()
}

// toHostname normalizes a URL with a https scheme to just its hostname.
//
// See also:
//
//   - https://github.com/moby/moby/blob/v25.0.3/registry/auth.go#L132
func toHostname(url string) string {
	stripped := url
	stripped = strings.TrimPrefix(stripped, "https://")

	return strings.Split(stripped, "/")[0]
}

// SplitLines splits a string into lines by '\r\n' or '\n'.
func SplitLines(s string) []string {
	if strings.Contains(s, "\r\n") {
		return strings.Split(s, "\r\n")
	}

	return strings.Split(s, "\n")
}
