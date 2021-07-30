package hinters

import (
	"testing"

	"github.com/liamg/darktile/internal/app/darktile/termutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_converter(t *testing.T) {

	tests := map[string]string{
		"dr-xr-xr-x": "0555",
		"-r-xr-xr-x": "0555",
		"drwxr-xr-x": "0755",
		"-rwxr-xr-x": "0755",
		"-r--------": "0400",
		"-rw-------": "0600",
	}

	for perm, octet := range tests {
		perm, err := getPermsFromString(perm)
		require.NoError(t, err)
		assert.Equal(t, octet, perm)
	}
}

func Test_perm_hinter_resolves_from_string_and_correctly_activates(t *testing.T) {

	hinter := &PermsHinter{}
	api := &TestAPI{}

	input := "dr-xr-xr-x"
	match, offset, length := hinter.Match(input, 3)

	assert.Equal(t, true, match)
	hinter.Activate(api, input[offset:offset+length], termutil.Position{}, termutil.Position{})
	assert.Equal(t, "0555", api.highlighted)
}

func Test_perm_hinter_resolves_from_string_and_correctly_activates_then_cleared(t *testing.T) {

	hinter := &PermsHinter{}
	api := &TestAPI{}

	input := "drwxr-xr-x"
	match, offset, length := hinter.Match(input, 3)

	assert.Equal(t, true, match)
	hinter.Activate(api, input[offset:offset+length], termutil.Position{}, termutil.Position{})
	assert.Equal(t, "0755", api.highlighted)
	hinter.Deactivate(api)
	assert.Equal(t, "", api.highlighted)
}
