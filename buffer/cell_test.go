package buffer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetRune(t *testing.T) {
	cell := newCell()
	assert.False(t, cell.hasContent)
	cell.setRune('X')
	assert.True(t, cell.hasContent)
	assert.Equal(t, 'X', cell.r)
	cell.setRune('Y')
	assert.True(t, cell.hasContent)
	assert.Equal(t, 'Y', cell.r)
}
