package buffer

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestBufferCreation(t *testing.T) {
	b := NewBuffer(10)
	assert.Equal(t, 10, b.Width())
	assert.Equal(t, 0, b.Column())
	assert.Equal(t, 0, b.Line())
	assert.NotNil(t, b.lines)
}

func TestBufferCursorIncrement(t *testing.T) {

	b := NewBuffer(5)
	b.incrementCursorPosition()
	require.Equal(t, 1, b.Column())
	require.Equal(t, 0, b.Line())

	b.incrementCursorPosition()
	require.Equal(t, 2, b.Column())
	require.Equal(t, 0, b.Line())

	b.incrementCursorPosition()
	require.Equal(t, 3, b.Column())
	require.Equal(t, 0, b.Line())

	b.incrementCursorPosition()
	require.Equal(t, 4, b.Column())
	require.Equal(t, 0, b.Line())

	b.incrementCursorPosition()
	require.Equal(t, 0, b.Column())
	require.Equal(t, 1, b.Line())

	b.incrementCursorPosition()
	b.incrementCursorPosition()
	b.incrementCursorPosition()
	b.incrementCursorPosition()
	b.incrementCursorPosition()
	b.incrementCursorPosition()
	b.incrementCursorPosition()
	b.incrementCursorPosition()
	b.incrementCursorPosition()
	b.incrementCursorPosition()

	require.Equal(t, 0, b.Column())
	require.Equal(t, 3, b.Line())

}

func TestBufferWrite(t *testing.T) {

}
