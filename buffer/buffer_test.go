package buffer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestOffsets(t *testing.T) {
	b := NewBuffer(10, 3, CellAttributes{})
	b.Write([]rune("hello\r\n")...)
	b.Write([]rune("hello\r\n")...)
	b.Write([]rune("hello\r\n")...)
	b.Write([]rune("hello\r\n")...)
	b.Write([]rune("hello")...)
	assert.Equal(t, uint16(10), b.ViewWidth())
	assert.Equal(t, uint16(10), b.Width())
	assert.Equal(t, uint16(3), b.ViewHeight())
	assert.Equal(t, 5, b.Height())
}

func TestBufferCreation(t *testing.T) {
	b := NewBuffer(10, 20, CellAttributes{})
	assert.Equal(t, uint16(10), b.Width())
	assert.Equal(t, uint16(20), b.ViewHeight())
	assert.Equal(t, uint16(0), b.CursorColumn())
	assert.Equal(t, uint16(0), b.CursorLine())
	assert.NotNil(t, b.lines)
}

func TestBufferWriteIncrementsCursorCorrectly(t *testing.T) {

	b := NewBuffer(5, 4, CellAttributes{})

	/*01234
	 |-----
	0|xxxxx
	1|
	2|
	3|
	 |-----
	*/

	b.Write('x')
	require.Equal(t, uint16(1), b.CursorColumn())
	require.Equal(t, uint16(0), b.CursorLine())

	b.Write('x')
	require.Equal(t, uint16(2), b.CursorColumn())
	require.Equal(t, uint16(0), b.CursorLine())

	b.Write('x')
	require.Equal(t, uint16(3), b.CursorColumn())
	require.Equal(t, uint16(0), b.CursorLine())

	b.Write('x')
	require.Equal(t, uint16(4), b.CursorColumn())
	require.Equal(t, uint16(0), b.CursorLine())

	b.Write('x')
	require.Equal(t, uint16(5), b.CursorColumn())
	require.Equal(t, uint16(0), b.CursorLine())

	b.Write('x')
	require.Equal(t, uint16(1), b.CursorColumn())
	require.Equal(t, uint16(1), b.CursorLine())

	b.Write('x')
	require.Equal(t, uint16(2), b.CursorColumn())
	require.Equal(t, uint16(1), b.CursorLine())

	lines := b.GetVisibleLines()
	require.Equal(t, 2, len(lines))
	assert.Equal(t, "xxxxx", lines[0].String())
	assert.Equal(t, "xx", lines[1].String())

}

func TestWritingNewLineAsFirstRuneOnWrappedLine(t *testing.T) {
	b := NewBuffer(3, 20, CellAttributes{})
	b.Write('a', 'b', 'c')
	assert.Equal(t, uint16(3), b.cursorX)
	assert.Equal(t, uint16(0), b.cursorY)
	b.Write(0x0a)
	assert.Equal(t, uint16(0), b.cursorX)
	assert.Equal(t, uint16(1), b.cursorY)

	b.Write('d', 'e', 'f')
	assert.Equal(t, uint16(3), b.cursorX)
	assert.Equal(t, uint16(1), b.cursorY)
	b.Write(0x0a)

	assert.Equal(t, uint16(0), b.cursorX)
	assert.Equal(t, uint16(2), b.cursorY)

	require.Equal(t, 2, len(b.lines))
	assert.Equal(t, "abc", b.lines[0].String())
	assert.Equal(t, "def", b.lines[1].String())

}

func TestWritingNewLineAsSecondRuneOnWrappedLine(t *testing.T) {
	b := NewBuffer(3, 20, CellAttributes{})
	/*
		|abc
		|d
		|ef
		|
		|
		|z
	*/

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

	b := NewBuffer(120, 80, CellAttributes{})
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
	b := NewBuffer(120, 80, CellAttributes{})
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

	b := NewBuffer(80, 10, CellAttributes{})
	b.Write([]rune("hello 1\r\n")...)
	b.Write([]rune("hello 2\r\n")...)
	b.Write([]rune("hello 3\r\n")...)
	b.Write([]rune("hello 4\r\n")...)
	b.Write([]rune("hello 5\r\n")...)
	b.Write([]rune("hello 6\r\n")...)
	b.Write([]rune("hello 7\r\n")...)
	b.Write([]rune("hello 8\r\n")...)
	b.Write([]rune("hello 9\r\n")...)
	b.Write([]rune("hello 10\r\n")...)
	b.Write([]rune("hello 11\r\n")...)
	b.Write([]rune("hello 12\r\n")...)
	b.Write([]rune("hello 13\r\n")...)
	b.Write([]rune("hello 14")...)

	lines := b.GetVisibleLines()
	require.Equal(t, 10, len(lines))
	assert.Equal(t, "hello 5", lines[0].String())
	assert.Equal(t, "hello 14", lines[9].String())

}

func TestClearWithoutFullView(t *testing.T) {
	b := NewBuffer(80, 10, CellAttributes{})
	b.Write([]rune("hello 1\r\n")...)
	b.Write([]rune("hello 2\r\n")...)
	b.Write([]rune("hello 3")...)
	b.Clear()
	lines := b.GetVisibleLines()
	for _, line := range lines {
		assert.Equal(t, "", line.String())
	}
}

func TestClearWithFullView(t *testing.T) {
	b := NewBuffer(80, 5, CellAttributes{})
	b.Write([]rune("hello 1\r\n")...)
	b.Write([]rune("hello 2\r\n")...)
	b.Write([]rune("hello 3\r\n")...)
	b.Write([]rune("hello 4\r\n")...)
	b.Write([]rune("hello 5\r\n")...)
	b.Write([]rune("hello 6\r\n")...)
	b.Write([]rune("hello 7\r\n")...)
	b.Write([]rune("hello 8\r\n")...)
	b.Clear()
	lines := b.GetVisibleLines()
	for _, line := range lines {
		assert.Equal(t, "", line.String())
	}
}

func TestCarriageReturn(t *testing.T) {
	b := NewBuffer(80, 20, CellAttributes{})
	b.Write([]rune("hello!")...)
	b.CarriageReturn()
	b.Write([]rune("secret")...)
	lines := b.GetVisibleLines()
	assert.Equal(t, "secret", lines[0].String())
}

func TestCarriageReturnOnFullLine(t *testing.T) {
	b := NewBuffer(20, 20, CellAttributes{})
	b.Write([]rune("abcdeabcdeabcdeabcde")...)
	b.CarriageReturn()
	b.Write([]rune("xxxxxxxxxxxxxxxxxxxx")...)
	lines := b.GetVisibleLines()
	assert.Equal(t, "xxxxxxxxxxxxxxxxxxxx", lines[0].String())
}

func TestCarriageReturnOnFullLastLine(t *testing.T) {
	b := NewBuffer(20, 2, CellAttributes{})
	b.Write([]rune("\nabcdeabcdeabcdeabcde\rxxxxxxxxxxxxxxxxxxxx")...)
	lines := b.GetVisibleLines()
	assert.Equal(t, "", lines[0].String())
	assert.Equal(t, "xxxxxxxxxxxxxxxxxxxx", lines[1].String())
}

func TestCarriageReturnOnWrappedLine(t *testing.T) {
	b := NewBuffer(80, 6, CellAttributes{})
	b.Write([]rune("hello!\rsecret")...)
	lines := b.GetVisibleLines()
	assert.Equal(t, "secret", lines[0].String())
}

func TestCarriageReturnOnLineThatDoesntExist(t *testing.T) {
	b := NewBuffer(6, 10, CellAttributes{})
	b.cursorY = 3
	b.Write('\r')
	assert.Equal(t, uint16(0), b.cursorX)
	assert.Equal(t, uint16(3), b.cursorY)
}

func TestGetCell(t *testing.T) {
	b := NewBuffer(80, 20, CellAttributes{})
	b.Write([]rune("Hello\r\nthere\r\nsomething...")...)
	cell := b.GetCell(8, 2)
	require.NotNil(t, cell)
	assert.Equal(t, 'g', cell.Rune())
}

func TestGetCellWithHistory(t *testing.T) {
	b := NewBuffer(80, 2, CellAttributes{})
	b.Write([]rune("Hello\r\nthere\r\nsomething...")...)
	cell := b.GetCell(8, 1)
	require.NotNil(t, cell)
	assert.Equal(t, 'g', cell.Rune())
}

func TestGetCellWithBadCursor(t *testing.T) {
	b := NewBuffer(80, 2, CellAttributes{})
	b.Write([]rune("Hello\r\nthere\r\nsomething...")...)
	require.Nil(t, b.GetCell(8, 3))
	require.Nil(t, b.GetCell(90, 0))

}

func TestCursorAttr(t *testing.T) {
	b := NewBuffer(80, 2, CellAttributes{})
	assert.Equal(t, &b.cursorAttr, b.CursorAttr())
}

func TestCursorPositionQuerying(t *testing.T) {
	b := NewBuffer(80, 20, CellAttributes{})
	b.cursorX = 17
	b.cursorY = 9
	assert.Equal(t, b.cursorX, b.CursorColumn())
	assert.Equal(t, b.cursorY, b.CursorLine())
}

func TestRawPositionQuerying(t *testing.T) {
	b := NewBuffer(80, 5, CellAttributes{})
	b.Write([]rune("a\r\na\r\na\r\na\r\na\r\na\r\na\r\na\r\na\r\na")...)
	b.cursorX = 3
	b.cursorY = 4
	assert.Equal(t, uint64(9), b.RawLine())
}

// CSI 2 K
func TestEraseLine(t *testing.T) {
	b := NewBuffer(80, 5, CellAttributes{})
	b.Write([]rune("hello, this is a test\r\nthis line should be deleted")...)
	b.EraseLine()
	assert.Equal(t, "hello, this is a test", b.lines[0].String())
	assert.Equal(t, "", b.lines[1].String())
}

// CSI 1 K
func TestEraseLineToCursor(t *testing.T) {
	b := NewBuffer(80, 5, CellAttributes{})
	b.Write([]rune("hello, this is a test\r\ndeleted")...)
	b.MovePosition(-3, 0)
	b.EraseLineToCursor()
	assert.Equal(t, "hello, this is a test", b.lines[0].String())
	assert.Equal(t, "\x00\x00\x00\x00\x00ed", b.lines[1].String())
}

// CSI 0 K
func TestEraseLineAfterCursor(t *testing.T) {
	b := NewBuffer(80, 5, CellAttributes{})
	b.Write([]rune("hello, this is a test\r\ndeleted")...)
	b.MovePosition(-3, 0)
	b.EraseLineFromCursor()
	assert.Equal(t, "hello, this is a test", b.lines[0].String())
	assert.Equal(t, "dele", b.lines[1].String())
}
func TestEraseDisplay(t *testing.T) {
	b := NewBuffer(80, 5, CellAttributes{})
	b.Write([]rune("hello\r\nasdasd\r\nthing")...)
	b.MovePosition(2, 1)
	b.EraseDisplay()
	lines := b.GetVisibleLines()
	for _, line := range lines {
		assert.Equal(t, "", line.String())
	}
}
func TestEraseDisplayToCursor(t *testing.T) {
	b := NewBuffer(80, 5, CellAttributes{})
	b.Write([]rune("hello\r\nasdasd\r\nthing")...)
	b.MovePosition(-2, 0)
	b.EraseDisplayToCursor()
	lines := b.GetVisibleLines()
	assert.Equal(t, "", lines[0].String())
	assert.Equal(t, "", lines[1].String())
	assert.Equal(t, "\x00\x00\x00ng", lines[2].String())

}

func TestEraseDisplayFromCursor(t *testing.T) {
	b := NewBuffer(80, 5, CellAttributes{})
	b.Write([]rune("hello\r\nasdasd\r\nthings")...)
	b.MovePosition(-3, -1)
	b.EraseDisplayFromCursor()
	lines := b.GetVisibleLines()
	assert.Equal(t, "hello", lines[0].String())
	assert.Equal(t, "asd", lines[1].String())
	assert.Equal(t, "", lines[2].String())
}
func TestBackspace(t *testing.T) {
	b := NewBuffer(80, 5, CellAttributes{})
	b.Write([]rune("hello")...)
	b.Backspace()
	b.Backspace()
	b.Write([]rune("p")...)
	lines := b.GetVisibleLines()
	assert.Equal(t, "helpo", lines[0].String())
}

func TestHorizontalResizeView(t *testing.T) {
	b := NewBuffer(80, 10, CellAttributes{})

	// 60 characters
	b.Write([]rune(
		`hellohellohellohellohellohellohellohellohellohellohellohello
goodbyegoodbye`)...)

	require.Equal(t, uint16(14), b.cursorX)
	require.Equal(t, uint16(1), b.cursorY)

	b.ResizeView(40, 10)

	expected := `hellohellohellohellohellohellohellohello
hellohellohellohello
goodbyegoodbye`

	require.Equal(t, uint16(14), b.cursorX)
	require.Equal(t, uint16(2), b.cursorY)

	lines := b.GetVisibleLines()
	strs := []string{}
	for _, l := range lines {
		strs = append(strs, l.String())
	}
	require.Equal(t, expected, strings.Join(strs, "\n"))

	b.ResizeView(20, 10)

	expected = `hellohellohellohello
hellohellohellohello
hellohellohellohello
goodbyegoodbye`

	lines = b.GetVisibleLines()
	strs = []string{}
	for _, l := range lines {
		strs = append(strs, l.String())
	}
	require.Equal(t, expected, strings.Join(strs, "\n"))

	b.ResizeView(10, 10)

	expected = `hellohello
hellohello
hellohello
hellohello
hellohello
hellohello
goodbyegoo
dbye`

	lines = b.GetVisibleLines()
	strs = []string{}
	for _, l := range lines {
		strs = append(strs, l.String())
	}
	require.Equal(t, expected, strings.Join(strs, "\n"))

	b.ResizeView(80, 20)

	expected = `hellohellohellohellohellohellohellohellohellohellohellohello
goodbyegoodbye`

	lines = b.GetVisibleLines()
	strs = []string{}
	for _, l := range lines {
		strs = append(strs, l.String())
	}
	require.Equal(t, expected, strings.Join(strs, "\n"))

	require.Equal(t, uint16(1), b.cursorY)
	require.Equal(t, uint16(14), b.cursorX)
}

/*
hellohellohellohellohellohellohellohellohellohellohellohello
goodbyegoo
dbye
*/
