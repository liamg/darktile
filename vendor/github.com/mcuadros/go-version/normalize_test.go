package version

import (
	"testing"
)

var versions = map[string]string{
	"1-stable":            "1.0.0.0",
	"1.0.0":               "1.0.0.0",
	"1.2.3.4":             "1.2.3.4",
	"1.0.0RC1dev":         "1.0.0.0-RC1-dev",
	"1.0.0-rC15-dev":      "1.0.0.0-RC15-dev",
	"1.0.0.RC.15-dev":     "1.0.0.0-RC15-dev",
	"1.0.0-rc1":           "1.0.0.0-RC1",
	"1.0.0.pl3-dev":       "1.0.0.0-patch3-dev",
	"1.0-dev":             "1.0.0.0-dev",
	"0":                   "0.0.0.0",
	"10.4.13-beta":        "10.4.13.0-beta",
	"10.4.13-b":           "10.4.13.0-beta",
	"10.4.13-b5":          "10.4.13.0-beta5",
	"v1.0.0":              "1.0.0.0",
	"v20100102":           "20100102",
	"2010.01":             "2010-01",
	"2010.01.02":          "2010-01-02",
	"2010-01-02":          "2010-01-02",
	"2010-01-02.5":        "2010-01-02-5",
	"20100102-203040":     "20100102-203040",
	"20100102203040-10":   "20100102203040-10",
	"20100102-203040-p1":  "20100102-203040-patch1",
	"dev-master":          "9999999-dev",
	"dev-trunk":           "9999999-dev",
	"1.x-dev":             "1.9999999.9999999.9999999-dev",
	"dev-feature-foo":     "dev-feature-foo",
	"DEV-FOOBAR":          "dev-FOOBAR",
	"dev-feature/foo":     "dev-feature/foo",
	"dev-master as 1.0.0": "9999999-dev",
}

func TestNormalize(t *testing.T) {
	for in, out := range versions {
		if x := Normalize(in); x != out {
			t.Errorf("FAIL: Normalize(%v) = %v: want %v", in, x, out)
		}
	}
}

var branches = map[string]string{
	"v1.x":      "1.9999999.9999999.9999999-dev",
	"v1.*":      "1.9999999.9999999.9999999-dev",
	"v1.0":      "1.0.9999999.9999999-dev",
	"2.0":       "2.0.9999999.9999999-dev",
	"v1.0.x":    "1.0.9999999.9999999-dev",
	"v1.0.3.*":  "1.0.3.9999999-dev",
	"v2.4.0":    "2.4.0.9999999-dev",
	"2.4.4":     "2.4.4.9999999-dev",
	"master":    "9999999-dev",
	"trunk":     "9999999-dev",
	"feature-a": "dev-feature-a",
	"FOOBAR":    "dev-FOOBAR",
}

func TestNormalizeBranch(t *testing.T) {
	for in, out := range branches {
		if x := normalizeBranch(in); x != out {
			t.Errorf("FAIL: normalizeBranch(%v) = %v: want %v", in, x, out)
		}
	}
}
