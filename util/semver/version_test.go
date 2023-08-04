// Based on https://github.com/Masterminds/semver/blob/v3.2.1/version_test.go

package semver

import "testing"

func TestNewVersion(t *testing.T) {
	tests := []struct {
		version string
		err     bool
	}{
		{"1.2.3", false},
		{"1.2.3+test.01", false},
		{"1.2.3-alpha.-1", false},
		{"v1.2.3", false},
		{"1.0", false},
		{"v1.0", false},
		{"1", false},
		{"v1", false},
		{"1.2.beta", true},
		{"v1.2.beta", true},
		{"foo", true},
		{"1.2-5", false},
		{"v1.2-5", false},
		{"1.2-beta.5", false},
		{"v1.2-beta.5", false},
		{"\n1.2", true},
		{"\nv1.2", true},
		{"1.2.0-x.Y.0+metadata", false},
		{"v1.2.0-x.Y.0+metadata", false},
		{"1.2.0-x.Y.0+metadata-width-hyphen", false},
		{"v1.2.0-x.Y.0+metadata-width-hyphen", false},
		{"1.2.3-rc1-with-hyphen", false},
		{"v1.2.3-rc1-with-hyphen", false},
		{"1.2.3.4", true},
		{"v1.2.3.4", true},
		{"1.2.2147483648", false},
		{"1.2147483648.3", false},
		{"2147483648.3.0", false},

		// Due to having 4 parts these should produce an error. See
		// https://github.com/Masterminds/semver/issues/185 for the reason for
		// these tests.
		{"12.3.4.1234", true},
		{"12.23.4.1234", true},
		{"12.3.34.1234", true},

		// The SemVer spec in a pre-release expects to allow [0-9A-Za-z-].
		{"20221209-update-renovatejson-v4", false},
	}

	for _, tc := range tests {
		_, err := NewVersion(tc.version)
		if tc.err && err == nil {
			t.Fatalf("expected error for version: %s", tc.version)
		} else if !tc.err && err != nil {
			t.Fatalf("error for version %s: %s", tc.version, err)
		}
	}
}

func TestParts(t *testing.T) {
	v, err := NewVersion("1.2.3")
	if err != nil {
		t.Error("Error parsing version 1.2.3")
	}

	if v.major != 1 {
		t.Error("major returning wrong value")
	}
	if v.minor != 2 {
		t.Error("minor returning wrong value")
	}
	if v.patch != 3 {
		t.Error("patch returning wrong value")
	}
}

func TestCoerceString(t *testing.T) {
	tests := []struct {
		version  string
		expected string
	}{
		{"1.2.3", "1.2.3"},
		{"v1.2.3", "1.2.3"},
		{"1.0", "1.0.0"},
		{"v1.0", "1.0.0"},
		{"1", "1.0.0"},
		{"v1", "1.0.0"},
	}

	for _, tc := range tests {
		v, err := NewVersion(tc.version)
		if err != nil {
			t.Errorf("Error parsing version %s", tc)
		}

		s := v.String()
		if s != tc.expected {
			t.Errorf("Error generating string. Expected '%s' but got '%s'", tc.expected, s)
		}
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.2.3", "1.5.1", -1},
		{"2.2.3", "1.5.1", 1},
		{"2.2.3", "2.2.2", 1},
	}

	for _, tc := range tests {
		v1, err := NewVersion(tc.v1)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		v2, err := NewVersion(tc.v2)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		a := v1.compare(v2)
		e := tc.expected
		if a != e {
			t.Errorf(
				"Comparison of '%s' and '%s' failed. Expected '%d', got '%d'",
				tc.v1, tc.v2, e, a,
			)
		}
	}
}

func TestGreaterThan(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected bool
	}{
		{"1.2.3", "1.5.1", false},
		{"2.2.3", "1.5.1", true},
		{"3.2-beta", "3.2-beta", false},
		{"3.2.0-beta.1", "3.2.0-beta.5", false},
		{"7.43.0-SNAPSHOT.99", "7.43.0-SNAPSHOT.103", false},
		{"7.43.0-SNAPSHOT.99", "7.43.0-SNAPSHOT.BAR", false},
	}

	for _, tc := range tests {
		v1, err := NewVersion(tc.v1)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		v2, err := NewVersion(tc.v2)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		a := v1.GreaterThan(v2)
		e := tc.expected
		if a != e {
			t.Errorf(
				"Comparison of '%s' and '%s' failed. Expected '%t', got '%t'",
				tc.v1, tc.v2, e, a,
			)
		}
	}
}

func TestGreaterThanOrEqual(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected bool
	}{
		{"1.2.3", "1.5.1", false},
		{"2.2.3", "1.5.1", true},
		{"3.2-beta", "3.2-beta", true},
		{"3.2-beta.4", "3.2-beta.2", true},
		{"7.43.0-SNAPSHOT.FOO", "7.43.0-SNAPSHOT.103", true},
	}

	for _, tc := range tests {
		v1, err := NewVersion(tc.v1)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		v2, err := NewVersion(tc.v2)
		if err != nil {
			t.Errorf("Error parsing version: %s", err)
		}

		a := v1.GreaterThanOrEqual(v2)
		e := tc.expected
		if a != e {
			t.Errorf(
				"Comparison of '%s' and '%s' failed. Expected '%t', got '%t'",
				tc.v1, tc.v2, e, a,
			)
		}
	}
}
