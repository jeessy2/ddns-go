package util

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObscure(t *testing.T) {
	for _, test := range []struct {
		in   string
		want string
		iv   string
	}{
		{"", "*bcb25d17YWFhYWFhYWFhYWFhYWFhYQ", "aaaaaaaaaaaaaaaa"},
		{"potato", "*bcb25d17YWFhYWFhYWFhYWFhYWFhYXMaGgIlEQ", "aaaaaaaaaaaaaaaa"},
		{"potato", "*bcb25d17YmJiYmJiYmJiYmJiYmJiYp3gcEWbAw", "bbbbbbbbbbbbbbbb"},
	} {
		cryptRand = bytes.NewBufferString(test.iv)
		got, err := Obscure(test.in)
		cryptRand = rand.Reader
		assert.NoError(t, err)
		assert.Equal(t, test.want, got)
		recoveredIn, err := Reveal(got)
		assert.NoError(t, err)
		assert.Equal(t, test.in, recoveredIn, "not bidirectional")
		// Now the Must variants
		cryptRand = bytes.NewBufferString(test.iv)
		got = MustObscure(test.in)
		cryptRand = rand.Reader
		assert.Equal(t, test.want, got)
		recoveredIn = MustReveal(got)
		assert.Equal(t, test.in, recoveredIn, "not bidirectional")

	}
}

func TestReveal(t *testing.T) {
	for _, test := range []struct {
		in   string
		want string
		iv   string
	}{
		{"*bcb25d17YWFhYWFhYWFhYWFhYWFhYQ", "", "aaaaaaaaaaaaaaaa"},
		{"*bcb25d17YWFhYWFhYWFhYWFhYWFhYXMaGgIlEQ", "potato", "aaaaaaaaaaaaaaaa"},
		{"*bcb25d17YmJiYmJiYmJiYmJiYmJiYp3gcEWbAw", "potato", "bbbbbbbbbbbbbbbb"},
		{"YmJiYmJiYmJiYmJiYmJiYp3gcEWbAw", "YmJiYmJiYmJiYmJiYmJiYp3gcEWbAw", "bbbbbbbbbbbbbbbb"},
	} {
		cryptRand = bytes.NewBufferString(test.iv)
		got, err := Reveal(test.in)
		assert.NoError(t, err)
		assert.Equal(t, test.want, got)
		// Now the Must variants
		cryptRand = bytes.NewBufferString(test.iv)
		got = MustReveal(test.in)
		assert.Equal(t, test.want, got)

	}
}

// Test some error cases
func TestRevealErrors(t *testing.T) {
	for _, test := range []struct {
		in      string
		wantErr string
	}{
		{"*bcb25d17YmJiYmJiYmJiYmJiYmJiYp*gcEWbAw", "base64 decode failed when revealing password - is it obscured?: illegal base64 data at input byte 22"},
		{"*bcb25d17aGVsbG8", "input too short when revealing password - is it obscured?"},
		{"*bcb25d17", "input too short when revealing password - is it obscured?"},
	} {
		gotString, gotErr := Reveal(test.in)
		assert.Equal(t, "", gotString)
		assert.Equal(t, test.wantErr, gotErr.Error())
	}
}
