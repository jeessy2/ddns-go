package config

import (
	"reflect"
	"testing"
)

// TestExtractHeaders 测试 parseHeaderArr
func TestExtractHeaders(t *testing.T) {
	input := `
a: foo
b: bar
c: foo:bar`
	expected := map[string]string{
		"a": "foo",
		"b": "bar",
		"c": "foo:bar",
	}

	parsedHeaders := extractHeaders(input)
	if !reflect.DeepEqual(parsedHeaders, expected) {
		t.Errorf("Expected %v, got %v", expected, parsedHeaders)
	}
}
