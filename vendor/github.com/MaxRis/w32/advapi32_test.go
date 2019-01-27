package w32

import (
	"testing"
)

func TestInitializeSecurityDescriptor(t *testing.T) {
	sd, err := InitializeSecurityDescriptor(1)
	if err != nil {
		t.Errorf("Failed: %v", err)
	}
	t.Logf("SD:\n%#v\n", *sd)
}

func TestSetSecurityDescriptorDacl(t *testing.T) {

	sd, err := InitializeSecurityDescriptor(1)
	if err != nil {
		t.Errorf("Failed to initialize: %v", err)
	}
	err = SetSecurityDescriptorDacl(sd, nil)
	if err != nil {
		t.Errorf("Failed to set NULL DACL: %v", err)
	}
	t.Logf("[OK] Set NULL DACL")

	empty := &ACL{
		AclRevision: 4,
		Sbz1:        0,
		AclSize:     4,
		AceCount:    0,
		Sbz2:        0,
	}
	err = SetSecurityDescriptorDacl(sd, empty)
	if err != nil {
		t.Errorf("Failed to set empty DACL: %v", err)
	}
	t.Logf("[OK] Set empty DACL")
	t.Logf("SD:\n%#v\n", *sd)

}
