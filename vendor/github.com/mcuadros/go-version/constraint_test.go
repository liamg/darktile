package version

import (
	"testing"
)

func TestGetOperator(t *testing.T) {
	constraint := NewConstrain("=", "1.0.0")
	out := "="

	if x := constraint.GetOperator(); x != "=" {
		t.Errorf("FAIL: GetOperator() = {%s}: want {%s}", x, out)
	}
}

func TestGetVersion(t *testing.T) {
	constraint := NewConstrain("=", "1.0.0")
	out := "1.0.0"

	if x := constraint.GetVersion(); x != "1.0.0" {
		t.Errorf("FAIL: GetVersion() = {%s}: want {%s}", x, out)
	}
}

func TestString(t *testing.T) {
	constraint := NewConstrain("=", "1.0.0")
	out := "= 1.0.0"

	if x := constraint.String(); x != out {
		t.Errorf("FAIL: String() = {%s}: want {%s}", x, out)
	}
}

func TestMatchSuccess(t *testing.T) {
	constraint := NewConstrain("=", "1.0.0")
	out := true

	if x := constraint.Match("1.0"); x != out {
		t.Errorf("FAIL: Match() = {%v}: want {%v}", x, out)
	}
}

func TestMatchFail(t *testing.T) {
	constraint := NewConstrain("=", "1.0.0")
	out := false

	if x := constraint.Match("2.0"); x != out {
		t.Errorf("FAIL: Match() = {%v}: want {%v}", x, out)
	}
}
