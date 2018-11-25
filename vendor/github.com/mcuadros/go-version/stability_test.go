package version

import (
	"testing"
)

var stabilityValues = map[string]int{
	"1.0":       Stable,
	"1.0-dev":   Development,
	"1.0-alpha": Alpha,
	"1.0b1":     Beta,
	"1.0rc1":    RC,
}

func TestGetStability(t *testing.T) {
	for in, out := range stabilityValues {
		if x := GetStability(in); x != out {
			t.Errorf("FAIL: GetStability(%v) = %v: want %v", in, x, out)
		}
	}
}
