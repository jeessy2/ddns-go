// Based on https://github.com/creativeprojects/go-selfupdate/blob/v1.1.1/decompress_test.go

package update

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

var buf = []byte{'a', 'b', 'c'}

func TestCompressionNotRequired(t *testing.T) {
	want := bytes.NewReader(buf)
	r, err := decompressCommand(want, "https://github.com/foo/bar/releases/download/v1.2.3/foo", "foo")
	if err != nil {
		t.Fatal(err)
	}

	have, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(buf, have) {
		t.Errorf("expected %v, got %v", buf, have)
	}
}

func TestMatchExecutableName(t *testing.T) {
	testData := []struct {
		cmd    string
		target string
		found  bool
	}{
		{"gostuff", "gostuff", true},
		{"gostuff", "gostuff_linux_x86_64", false},
		{"gostuff", "gostuff_darwin_amd64", false},
		{"gostuff", "gostuff.exe", true},
		{"gostuff", "gostuff_windows_amd64.exe", false},
	}

	for _, testItem := range testData {
		t.Run(testItem.target, func(t *testing.T) {
			if matchExecutableName(testItem.cmd, testItem.target) != testItem.found {
				t.Errorf("Expected '%t' but got '%t'", testItem.found, matchExecutableName(testItem.cmd, testItem.target))
			}
		})
	}
}

func TestErrorFromReader(t *testing.T) {
	extensions := []string{
		"zip",
		"tar.gz",
	}

	for _, extension := range extensions {
		t.Run(extension, func(t *testing.T) {
			reader, err := decompressCommand(bytes.NewReader(buf), "foo."+extension, "foo."+extension)
			if err != nil {
				if !strings.Contains(err.Error(), errCannotDecompressFile.Error()) {
					t.Fatalf("Expected error: EOF, got: %v", err)
				}
			} else {
				_, err = io.ReadAll(reader)
				if err == nil {
					t.Fatalf("An error is expected but got nil.")
				}
			}
		})
	}
}
