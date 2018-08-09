package config

import (
	"bytes"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestColourTomlEncoding(t *testing.T) {
	target := struct {
		Orange Colour `toml:"colour"`
	}{
		Orange: Colour([3]float32{1, 0.5, 0}),
	}
	var buf bytes.Buffer
	e := toml.NewEncoder(&buf)
	err := e.Encode(target)
	require.Nil(t, err)
	assert.Equal(t, `colour = "#ff7f00"
`, buf.String())

}
func TestColourTomlUnmarshalling(t *testing.T) {
	target := struct {
		Purple Colour `toml:"colour"`
	}{}
	err := toml.Unmarshal([]byte(`colour = "#7f00ff"`), &target)
	require.Nil(t, err)
	assert.InDelta(t, 0.5, target.Purple[0], 0.01)
	assert.InDelta(t, 0.0, target.Purple[1], 0.01)
	assert.InDelta(t, 1.0, target.Purple[2], 0.01)
}
