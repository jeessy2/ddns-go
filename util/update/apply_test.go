// Based on https://github.com/inconshreveable/go-update/blob/7a872911e5b39953310f0a04161f0d50c7e63755/apply_test.go

package update

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

var (
	oldFile = []byte{0xDE, 0xAD, 0xBE, 0xEF}
	newFile = []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
)

func cleanup(path string) {
	os.Remove(path)
	os.Remove(fmt.Sprintf("%s.new", path))
}

// we write with a separate name for each test so that we can run them in parallel
func writeOldFile(path string, t *testing.T) {
	if err := os.WriteFile(path, oldFile, 0777); err != nil {
		t.Fatalf("Failed to write file for testing preparation: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("Failed to stat file for testing preparation: %v", err)
	}
}

func validateUpdate(path string, err error, t *testing.T) {
	if err != nil {
		t.Fatalf("Failed to update: %v", err)
	}

	buf, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file post-update: %v", err)
	}

	if !bytes.Equal(buf, newFile) {
		t.Fatalf("File was not updated! Bytes read: %v, Bytes expected: %v", buf, newFile)
	}
}

func TestApply(t *testing.T) {
	t.Parallel()

	fName := "TestApply"
	defer cleanup(fName)
	writeOldFile(fName, t)

	err := apply(bytes.NewReader(newFile), fName)
	validateUpdate(fName, err, t)
}
