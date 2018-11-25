package config

import (
	"testing"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyCombinations(t *testing.T) {

	combi, err := parseKeyCombination("ctrl + alt + a")
	require.Nil(t, err)
	require.NotNil(t, combi)

	assert.Equal(t, glfw.KeyA, combi.key)
	assert.Equal(t, glfw.ModControl+glfw.ModAlt, combi.mods)

	assert.True(t, combi.Match(glfw.ModControl^glfw.ModAlt, glfw.KeyA))
	assert.False(t, combi.Match(glfw.ModControl^glfw.ModAlt, glfw.KeyB))
	assert.False(t, combi.Match(glfw.ModControl, glfw.KeyA))
	assert.False(t, combi.Match(glfw.ModAlt, glfw.KeyA))
	assert.False(t, combi.Match(0, glfw.KeyA))
	assert.False(t, combi.Match(glfw.ModControl^glfw.ModAlt^glfw.ModShift, glfw.KeyA))

}
