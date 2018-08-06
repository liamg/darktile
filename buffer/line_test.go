package buffer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLine(t *testing.T) {

	line := newLine()
	line.cells = []Cell{
		{r: 'h'},
		{r: 'e'},
		{r: 'l'},
		{r: 'l'},
		{r: 'o'},
	}

	assert.Equal(t, "hello", line.String())
	assert.False(t, line.wrapped)

	line.setWrapped(true)
	assert.True(t, line.wrapped)

	line.setWrapped(false)
	assert.False(t, line.wrapped)

}
