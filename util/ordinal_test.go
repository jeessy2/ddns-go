package util

import "testing"

func TestOrdinal(t *testing.T) {
	lang := "en"

	tests := []struct {
		name string
		got  string
		want string
	}{
		{"0", Ordinal(0, lang), "0th"},
		{"1", Ordinal(1, lang), "1st"},
		{"2", Ordinal(2, lang), "2nd"},
		{"3", Ordinal(3, lang), "3rd"},
		{"4", Ordinal(4, lang), "4th"},
		{"10", Ordinal(10, lang), "10th"},
		{"11", Ordinal(11, lang), "11th"},
		{"12", Ordinal(12, lang), "12th"},
		{"13", Ordinal(13, lang), "13th"},
		{"21", Ordinal(21, lang), "21st"},
		{"32", Ordinal(32, lang), "32nd"},
		{"43", Ordinal(43, lang), "43rd"},
		{"101", Ordinal(101, lang), "101st"},
		{"102", Ordinal(102, lang), "102nd"},
		{"103", Ordinal(103, lang), "103rd"},
		{"211", Ordinal(211, lang), "211th"},
		{"212", Ordinal(212, lang), "212th"},
		{"213", Ordinal(213, lang), "213th"},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("On %s, Expected %s, but got %s", tt.name, tt.want, tt.got)
		}
	}
}
