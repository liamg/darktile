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

	assert.Equal(t, 'a', combi.char)
	assert.Equal(t, glfw.ModControl+glfw.ModAlt, combi.mods)

	assert.True(t, combi.Match(glfw.ModControl^glfw.ModAlt, 'a'))
	assert.False(t, combi.Match(glfw.ModControl^glfw.ModAlt, 'b'))
	assert.False(t, combi.Match(glfw.ModControl, 'b'))
	assert.False(t, combi.Match(glfw.ModAlt, 'd'))
	assert.False(t, combi.Match(0, 'e'))
	assert.False(t, combi.Match(glfw.ModControl^glfw.ModAlt^glfw.ModShift, 'f'))

}
