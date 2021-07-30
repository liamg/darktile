package hinters

import (
	"testing"

	"github.com/liamg/darktile/internal/app/darktile/termutil"
	"github.com/stretchr/testify/assert"
)

func Test_b64_hinter_resolves_from_base64_correctly_activated(t *testing.T) {

	hinter := &Base64Hinter{}
	api := &TestAPI{}
	text := "This is the result SGVsbG8gTGlhbQ=="

	match, offset, length := hinter.Match(text, 28)

	assert.Equal(t, true, match)
	hinter.Activate(api, text[offset:offset+length], termutil.Position{}, termutil.Position{})

	assert.Equal(t, "Base64 decodes to:\nHello Liam", api.highlighted)
}

func Test_b64_hinter_resolves_from_base64_correctly_activated_then_cleared(t *testing.T) {

	hinter := &Base64Hinter{}
	api := &TestAPI{}
	text := "This is the result SGVsbG8gTGlhbQ=="

	match, offset, length := hinter.Match(text, 28)

	assert.Equal(t, true, match)
	hinter.Activate(api, text[offset:offset+length], termutil.Position{}, termutil.Position{})

	assert.Equal(t, "Base64 decodes to:\nHello Liam", api.highlighted)
	hinter.Deactivate(api)
	assert.Equal(t, "", api.highlighted)
}

func Test_b64_hinter_doesnt_match_random_junk(t *testing.T) {

	hinter := &Base64Hinter{}
	text := "This is the result SGVsbG8eTGlhbQ=="

	match, _, _ := hinter.Match(text, 10)

	assert.Equal(t, false, match)
}
