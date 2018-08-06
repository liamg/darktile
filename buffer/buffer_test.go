package buffer

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestOffsets(t *testing.T) {
	b := NewBuffer(10, 8)
	test := "hellothere\nhellothere\nhellothere\nhellothere\nhellothere\nhellothere\nhellothere\nhellothere\nhellothere\nhellothere\nhellothere\nhellothere\n?"
	b.Write([]rune(test)...)
	assert.Equal(t, uint16(10), b.ViewWidth())
	assert.Equal(t, uint16(10), b.Width())
	assert.Equal(t, uint16(8), b.ViewHeight())
	assert.Equal(t, 13, b.Height())
}

func TestBufferCreation(t *testing.T) {
	b := NewBuffer(10, 20)
	assert.Equal(t, uint16(10), b.Width())
	assert.Equal(t, uint16(20), b.ViewHeight())
	assert.Equal(t, uint16(0), b.CursorColumn())
	assert.Equal(t, uint16(0), b.CursorLine())
	assert.NotNil(t, b.lines)
}

func TestBufferCursorIncrement(t *testing.T) {

	b := NewBuffer(5, 4)
	b.incrementCursorPosition()
	require.Equal(t, uint16(1), b.CursorColumn())
	require.Equal(t, uint16(0), b.CursorLine())

	b.incrementCursorPosition()
	require.Equal(t, uint16(2), b.CursorColumn())
	require.Equal(t, uint16(0), b.CursorLine())

	b.incrementCursorPosition()
	require.Equal(t, uint16(3), b.CursorColumn())
	require.Equal(t, uint16(0), b.CursorLine())

	b.incrementCursorPosition()
	require.Equal(t, uint16(4), b.CursorColumn())
	require.Equal(t, uint16(0), b.CursorLine())

	b.incrementCursorPosition()
	require.Equal(t, uint16(0), b.CursorColumn())
	require.Equal(t, uint16(1), b.CursorLine())

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

	require.Equal(t, uint16(0), b.CursorColumn())
	require.Equal(t, uint16(3), b.CursorLine())

	b.Write([]rune("hello\n")...)
	b.Write([]rune("hello\n")...)
	b.Write([]rune("hello\n")...)
	b.Write([]rune("hello\n")...)
	b.Write([]rune("hello\n")...)
	b.Write([]rune("hello")...)
	b.SetPosition(0, 2)
	b.incrementCursorPosition()

}
func TestBufferWrite(t *testing.T) {

	b := NewBuffer(5, 20)

	assert.Equal(t, uint16(0), b.CursorColumn())
	assert.Equal(t, uint16(0), b.CursorLine())

	b.Write('a')
	assert.Equal(t, uint16(1), b.CursorColumn())
	assert.Equal(t, uint16(0), b.CursorLine())

	b.Write('b')
	assert.Equal(t, uint16(2), b.CursorColumn())
	assert.Equal(t, uint16(0), b.CursorLine())

	b.Write('c')
	assert.Equal(t, uint16(3), b.CursorColumn())
	assert.Equal(t, uint16(0), b.CursorLine())

	b.Write('d')
	assert.Equal(t, uint16(4), b.CursorColumn())
	assert.Equal(t, uint16(0), b.CursorLine())

	b.Write('e')
	assert.Equal(t, uint16(0), b.CursorColumn())
	assert.Equal(t, uint16(1), b.CursorLine())

	b.Write('f')
	assert.Equal(t, uint16(1), b.CursorColumn())
	assert.Equal(t, uint16(1), b.CursorLine())

	//b.lines[0].cells[]

}

func TestWritingNewLineAsFirstRuneOnWrappedLine(t *testing.T) {
	b := NewBuffer(3, 20)
	b.Write('a', 'b', 'c')
	assert.Equal(t, uint16(0), b.cursorX)
	b.Write(0x0a)
	b.Write('d', 'e', 'f')
	b.Write(0x0a)

	assert.Equal(t, "abc", b.lines[0].String())
	assert.Equal(t, "def", b.lines[1].String())
	assert.Equal(t, "", b.lines[2].String())

}

func TestWritingNewLineAsSecondRuneOnWrappedLine(t *testing.T) {
	b := NewBuffer(3, 20)
	b.Write('a', 'b', 'c', 'd')
	b.Write(0x0a)
	b.Write('e', 'f')
	b.Write(0x0a)
	b.Write(0x0a)
	b.Write(0x0a)
	b.Write('z')

	assert.Equal(t, "abc", b.lines[0].String())
	assert.Equal(t, "d", b.lines[1].String())
	assert.Equal(t, "ef", b.lines[2].String())
	assert.Equal(t, "", b.lines[3].String())
	assert.Equal(t, "", b.lines[4].String())
	assert.Equal(t, "z", b.lines[5].String())
}

func TestSetPosition(t *testing.T) {

	b := NewBuffer(120, 80)
	assert.Equal(t, 0, int(b.CursorColumn()))
	assert.Equal(t, 0, int(b.CursorLine()))
	b.SetPosition(60, 10)
	assert.Equal(t, 60, int(b.CursorColumn()))
	assert.Equal(t, 10, int(b.CursorLine()))
	b.SetPosition(0, 0)
	assert.Equal(t, 0, int(b.CursorColumn()))
	assert.Equal(t, 0, int(b.CursorLine()))
	b.SetPosition(120, 90)
	assert.Equal(t, 119, int(b.CursorColumn()))
	assert.Equal(t, 79, int(b.CursorLine()))

}

func TestMovePosition(t *testing.T) {
	b := NewBuffer(120, 80)
	assert.Equal(t, 0, int(b.CursorColumn()))
	assert.Equal(t, 0, int(b.CursorLine()))
	b.MovePosition(-1, -1)
	assert.Equal(t, 0, int(b.CursorColumn()))
	assert.Equal(t, 0, int(b.CursorLine()))
	b.MovePosition(30, 20)
	assert.Equal(t, 30, int(b.CursorColumn()))
	assert.Equal(t, 20, int(b.CursorLine()))
	b.MovePosition(30, 20)
	assert.Equal(t, 60, int(b.CursorColumn()))
	assert.Equal(t, 40, int(b.CursorLine()))
	b.MovePosition(-1, -1)
	assert.Equal(t, 59, int(b.CursorColumn()))
	assert.Equal(t, 39, int(b.CursorLine()))
	b.MovePosition(100, 100)
	assert.Equal(t, 119, int(b.CursorColumn()))
	assert.Equal(t, 79, int(b.CursorLine()))
}

func TestVisibleLines(t *testing.T) {

	b := NewBuffer(80, 10)
	b.Write([]rune("hello 1\n")...)
	b.Write([]rune("hello 2\n")...)
	b.Write([]rune("hello 3\n")...)
	b.Write([]rune("hello 4\n")...)
	b.Write([]rune("hello 5\n")...)
	b.Write([]rune("hello 6\n")...)
	b.Write([]rune("hello 7\n")...)
	b.Write([]rune("hello 8\n")...)
	b.Write([]rune("hello 9\n")...)
	b.Write([]rune("hello 10\n")...)
	b.Write([]rune("hello 11\n")...)
	b.Write([]rune("hello 12\n")...)
	b.Write([]rune("hello 13\n")...)
	b.Write([]rune("hello 14")...)

	lines := b.GetVisibleLines()
	require.Equal(t, 10, len(lines))
	assert.Equal(t, "hello 5", lines[0].String())
	assert.Equal(t, "hello 14", lines[9].String())

}

func TestClearWithoutFullView(t *testing.T) {
	b := NewBuffer(80, 10)
	b.Write([]rune("hello 1\n")...)
	b.Write([]rune("hello 2\n")...)
	b.Write([]rune("hello 3")...)
	b.Clear()
	lines := b.GetVisibleLines()
	for _, line := range lines {
		assert.Equal(t, "", line.String())
	}
}

func TestClearWithFullView(t *testing.T) {
	b := NewBuffer(80, 5)
	b.Write([]rune("hello 1\n")...)
	b.Write([]rune("hello 2\n")...)
	b.Write([]rune("hello 3\n")...)
	b.Write([]rune("hello 4\n")...)
	b.Write([]rune("hello 5\n")...)
	b.Write([]rune("hello 6\n")...)
	b.Write([]rune("hello 7\n")...)
	b.Write([]rune("hello 8\n")...)
	b.Clear()
	lines := b.GetVisibleLines()
	for _, line := range lines {
		assert.Equal(t, "", line.String())
	}
}

func TestResizeView(t *testing.T) {
	b := NewBuffer(80, 20)
	b.ResizeView(40, 10)
}
