package buffer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestTabbing(t *testing.T) {
	b := NewBuffer(NewTerminalState(30, 3, CellAttributes{}, 1000))
	b.Write([]rune("hello")...)
	b.Tab()
	b.Write([]rune("x")...)
	b.Tab()
	b.Write([]rune("goodbye")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hell")...)
	b.Tab()
	b.Write([]rune("xxx")...)
	b.Tab()
	b.Write([]rune("good")...)
	b.CarriageReturn()
	b.NewLine()
	expected := `
hello   x   goodbye
hell    xxx good
`

	lines := b.GetVisibleLines()
	strs := []string{}
	for _, l := range lines {
		strs = append(strs, l.String())
	}
	require.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(strings.Join(strs, "\n")))
}

func TestOffsets(t *testing.T) {
	b := NewBuffer(NewTerminalState(10, 3, CellAttributes{}, 1000))
	b.Write([]rune("hello")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello")...)
	assert.Equal(t, uint16(10), b.ViewWidth())
	assert.Equal(t, uint16(10), b.Width())
	assert.Equal(t, uint16(3), b.ViewHeight())
	assert.Equal(t, 5, b.Height())
}

func TestBufferCreation(t *testing.T) {
	b := NewBuffer(NewTerminalState(10, 20, CellAttributes{}, 1000))
	assert.Equal(t, uint16(10), b.Width())
	assert.Equal(t, uint16(20), b.ViewHeight())
	assert.Equal(t, uint16(0), b.CursorColumn())
	assert.Equal(t, uint16(0), b.CursorLine())
	assert.NotNil(t, b.lines)
}

func TestBufferWriteIncrementsCursorCorrectly(t *testing.T) {

	b := NewBuffer(NewTerminalState(5, 4, CellAttributes{}, 1000))

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
	b := NewBuffer(NewTerminalState(3, 20, CellAttributes{}, 1000))
	b.terminalState.LineFeedMode = false

	b.Write('a', 'b', 'c')
	assert.Equal(t, uint16(3), b.terminalState.cursorX)
	assert.Equal(t, uint16(0), b.terminalState.cursorY)
	b.NewLine()
	assert.Equal(t, uint16(0), b.terminalState.cursorX)
	assert.Equal(t, uint16(1), b.terminalState.cursorY)

	b.Write('d', 'e', 'f')
	assert.Equal(t, uint16(3), b.terminalState.cursorX)
	assert.Equal(t, uint16(1), b.terminalState.cursorY)
	b.NewLine()

	assert.Equal(t, uint16(0), b.terminalState.cursorX)
	assert.Equal(t, uint16(2), b.terminalState.cursorY)

	require.Equal(t, 3, len(b.lines))
	assert.Equal(t, "abc", b.lines[0].String())
	assert.Equal(t, "def", b.lines[1].String())

}

func TestWritingNewLineAsSecondRuneOnWrappedLine(t *testing.T) {
	b := NewBuffer(NewTerminalState(3, 20, CellAttributes{}, 1000))
	b.terminalState.LineFeedMode = false
	/*
		|abc
		|d
		|ef
		|
		|
		|z
	*/

	b.Write('a', 'b', 'c', 'd')
	b.NewLine()
	b.Write('e', 'f')
	b.NewLine()
	b.NewLine()
	b.NewLine()
	b.Write('z')

	assert.Equal(t, "abc", b.lines[0].String())
	assert.Equal(t, "d", b.lines[1].String())
	assert.Equal(t, "ef", b.lines[2].String())
	assert.Equal(t, "", b.lines[3].String())
	assert.Equal(t, "", b.lines[4].String())
	assert.Equal(t, "z", b.lines[5].String())
}

func TestSetPosition(t *testing.T) {

	b := NewBuffer(NewTerminalState(120, 80, CellAttributes{}, 1000))
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
	b := NewBuffer(NewTerminalState(120, 80, CellAttributes{}, 1000))
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

	b := NewBuffer(NewTerminalState(80, 10, CellAttributes{}, 1000))
	b.Write([]rune("hello 1")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello 2")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello 3")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello 4")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello 5")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello 6")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello 7")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello 8")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello 9")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello 10")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello 11")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello 12")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello 13")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello 14")...)

	lines := b.GetVisibleLines()
	require.Equal(t, 10, len(lines))
	assert.Equal(t, "hello 5", lines[0].String())
	assert.Equal(t, "hello 14", lines[9].String())

}

func TestClearWithoutFullView(t *testing.T) {
	b := NewBuffer(NewTerminalState(80, 10, CellAttributes{}, 1000))
	b.Write([]rune("hello 1")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello 1")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello 1")...)
	b.Clear()
	lines := b.GetVisibleLines()
	for _, line := range lines {
		assert.Equal(t, "", line.String())
	}
}

func TestClearWithFullView(t *testing.T) {
	b := NewBuffer(NewTerminalState(80, 5, CellAttributes{}, 1000))
	b.Write([]rune("hello 1")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello 1")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello 1")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello 1")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello 1")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello 1")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello 1")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("hello 1")...)
	b.Clear()
	lines := b.GetVisibleLines()
	for _, line := range lines {
		assert.Equal(t, "", line.String())
	}
}

func TestCarriageReturn(t *testing.T) {
	b := NewBuffer(NewTerminalState(80, 20, CellAttributes{}, 1000))
	b.Write([]rune("hello!")...)
	b.CarriageReturn()
	b.Write([]rune("secret")...)
	lines := b.GetVisibleLines()
	assert.Equal(t, "secret", lines[0].String())
}

func TestCarriageReturnOnFullLine(t *testing.T) {
	b := NewBuffer(NewTerminalState(20, 20, CellAttributes{}, 1000))
	b.Write([]rune("abcdeabcdeabcdeabcde")...)
	b.CarriageReturn()
	b.Write([]rune("xxxxxxxxxxxxxxxxxxxx")...)
	lines := b.GetVisibleLines()
	assert.Equal(t, "xxxxxxxxxxxxxxxxxxxx", lines[0].String())
}

func TestCarriageReturnOnFullLastLine(t *testing.T) {
	b := NewBuffer(NewTerminalState(20, 2, CellAttributes{}, 1000))
	b.NewLine()
	b.Write([]rune("abcdeabcdeabcdeabcde")...)
	b.CarriageReturn()
	b.Write([]rune("xxxxxxxxxxxxxxxxxxxx")...)
	lines := b.GetVisibleLines()
	assert.Equal(t, "", lines[0].String())
	assert.Equal(t, "xxxxxxxxxxxxxxxxxxxx", lines[1].String())
}

func TestCarriageReturnOnWrappedLine(t *testing.T) {
	b := NewBuffer(NewTerminalState(80, 6, CellAttributes{}, 1000))
	b.Write([]rune("hello!")...)
	b.CarriageReturn()
	b.Write([]rune("secret")...)

	lines := b.GetVisibleLines()
	assert.Equal(t, "secret", lines[0].String())
}

func TestCarriageReturnOnLineThatDoesntExist(t *testing.T) {
	b := NewBuffer(NewTerminalState(6, 10, CellAttributes{}, 1000))
	b.terminalState.cursorY = 3
	b.CarriageReturn()
	assert.Equal(t, uint16(0), b.terminalState.cursorX)
	assert.Equal(t, uint16(3), b.terminalState.cursorY)
}

func TestGetCell(t *testing.T) {
	b := NewBuffer(NewTerminalState(80, 20, CellAttributes{}, 1000))
	b.Write([]rune("Hello")...)
	b.CarriageReturn()
	b.NewLine()

	b.Write([]rune("there")...)
	b.CarriageReturn()
	b.NewLine()

	b.Write([]rune("something...")...)
	cell := b.GetCell(8, 2)
	require.NotNil(t, cell)
	assert.Equal(t, 'g', cell.Rune())
}

func TestGetCellWithHistory(t *testing.T) {
	b := NewBuffer(NewTerminalState(80, 2, CellAttributes{}, 1000))

	b.Write([]rune("Hello")...)
	b.CarriageReturn()
	b.NewLine()

	b.Write([]rune("there")...)
	b.CarriageReturn()
	b.NewLine()

	b.Write([]rune("something...")...)

	cell := b.GetCell(8, 1)
	require.NotNil(t, cell)
	assert.Equal(t, 'g', cell.Rune())
}

func TestGetCellWithBadCursor(t *testing.T) {
	b := NewBuffer(NewTerminalState(80, 2, CellAttributes{}, 1000))
	b.Write([]rune("Hello\r\nthere\r\nsomething...")...)
	require.Nil(t, b.GetCell(8, 3))
	require.Nil(t, b.GetCell(90, 0))

}

func TestCursorAttr(t *testing.T) {
	b := NewBuffer(NewTerminalState(80, 2, CellAttributes{}, 1000))
	assert.Equal(t, &b.terminalState.CursorAttr, b.CursorAttr())
}

func TestCursorPositionQuerying(t *testing.T) {
	b := NewBuffer(NewTerminalState(80, 20, CellAttributes{}, 1000))
	b.terminalState.cursorX = 17
	b.terminalState.cursorY = 9
	assert.Equal(t, b.terminalState.cursorX, b.CursorColumn())
	assert.Equal(t, b.terminalState.cursorY, b.CursorLine())
}

func TestRawPositionQuerying(t *testing.T) {
	b := NewBuffer(NewTerminalState(80, 5, CellAttributes{}, 1000))
	b.Write([]rune("a")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("a")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("a")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("a")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("a")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("a")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("a")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("a")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("a")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("a")...)

	b.terminalState.cursorX = 3
	b.terminalState.cursorY = 4
	assert.Equal(t, uint64(9), b.RawLine())
}

// CSI 2 K
func TestEraseLine(t *testing.T) {
	b := NewBuffer(NewTerminalState(80, 5, CellAttributes{}, 1000))
	b.Write([]rune("hello, this is a test")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("this line should be deleted")...)
	b.EraseLine()
	assert.Equal(t, "hello, this is a test", b.lines[0].String())
	assert.Equal(t, "", b.lines[1].String())
}

// CSI 1 K
func TestEraseLineToCursor(t *testing.T) {
	b := NewBuffer(NewTerminalState(80, 5, CellAttributes{}, 1000))
	b.Write([]rune("hello, this is a test")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("deleted")...)

	b.MovePosition(-3, 0)
	b.EraseLineToCursor()
	assert.Equal(t, "hello, this is a test", b.lines[0].String())
	assert.Equal(t, "\x00\x00\x00\x00\x00ed", b.lines[1].String())
}

// CSI 0 K
func TestEraseLineAfterCursor(t *testing.T) {
	b := NewBuffer(NewTerminalState(80, 5, CellAttributes{}, 1000))
	b.Write([]rune("hello, this is a test")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("deleted")...)
	b.MovePosition(-3, 0)
	b.EraseLineFromCursor()
	assert.Equal(t, "hello, this is a test", b.lines[0].String())
	assert.Equal(t, "dele", b.lines[1].String())
}
func TestEraseDisplay(t *testing.T) {
	b := NewBuffer(NewTerminalState(80, 5, CellAttributes{}, 1000))
	b.Write([]rune("hello")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("asdasd")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("thing")...)
	b.MovePosition(2, 1)
	b.EraseDisplay()
	lines := b.GetVisibleLines()
	for _, line := range lines {
		assert.Equal(t, "", line.String())
	}
}
func TestEraseDisplayToCursor(t *testing.T) {
	b := NewBuffer(NewTerminalState(80, 5, CellAttributes{}, 1000))
	b.Write([]rune("hello")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("asdasd")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("thing")...)
	b.MovePosition(-2, 0)
	b.EraseDisplayToCursor()
	lines := b.GetVisibleLines()
	assert.Equal(t, "", lines[0].String())
	assert.Equal(t, "", lines[1].String())
	assert.Equal(t, "\x00\x00\x00\x00g", lines[2].String())

}

func TestEraseDisplayFromCursor(t *testing.T) {
	b := NewBuffer(NewTerminalState(80, 5, CellAttributes{}, 1000))
	b.Write([]rune("hello")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("asdasd")...)
	b.CarriageReturn()
	b.NewLine()
	b.Write([]rune("things")...)
	b.MovePosition(-3, -1)
	b.EraseDisplayFromCursor()
	lines := b.GetVisibleLines()
	assert.Equal(t, "hello", lines[0].String())
	assert.Equal(t, "asd", lines[1].String())
	assert.Equal(t, "", lines[2].String())
}
func TestBackspace(t *testing.T) {
	b := NewBuffer(NewTerminalState(80, 5, CellAttributes{}, 1000))
	b.Write([]rune("hello")...)
	b.Backspace()
	b.Backspace()
	b.Write([]rune("p")...)
	lines := b.GetVisibleLines()
	assert.Equal(t, "helpo", lines[0].String())
}

func TestHorizontalResizeView(t *testing.T) {
	b := NewBuffer(NewTerminalState(80, 10, CellAttributes{}, 1000))

	// 60 characters
	b.Write([]rune(`hellohellohellohellohellohellohellohellohellohellohellohello`)...)

	b.CarriageReturn()
	b.NewLine()

	b.Write([]rune(`goodbyegoodbye`)...)

	require.Equal(t, uint16(14), b.terminalState.cursorX)
	require.Equal(t, uint16(1), b.terminalState.cursorY)

	b.ResizeView(40, 10)

	expected := `hellohellohellohellohellohellohellohello
hellohellohellohello
goodbyegoodbye`

	require.Equal(t, uint16(14), b.terminalState.cursorX)
	require.Equal(t, uint16(2), b.terminalState.cursorY)

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

	require.Equal(t, uint16(1), b.terminalState.cursorY)
	require.Equal(t, uint16(14), b.terminalState.cursorX)
}

/*
hellohellohellohellohellohellohellohellohellohellohellohello
goodbyegoo
dbye
*/

func TestBufferMaxLines(t *testing.T) {
	b := NewBuffer(NewTerminalState(80, 2, CellAttributes{}, 2))
	b.terminalState.LineFeedMode = false

	b.Write([]rune("hello")...)
	b.NewLine()
	b.Write([]rune("funny")...)
	b.NewLine()
	b.Write([]rune("world")...)

	assert.Equal(t, 2, len(b.lines))
	assert.Equal(t, "funny", b.lines[0].String())
	assert.Equal(t, "world", b.lines[1].String())
}

func makeBufferForTestingSelection() *Buffer {
	b := NewBuffer(NewTerminalState(80, 10, CellAttributes{}, 10))
	b.terminalState.LineFeedMode = false

	b.Write([]rune("The quick brown")...)
	b.NewLine()
	b.Write([]rune("fox jumps over")...)
	b.NewLine()
	b.Write([]rune("the lazy dog")...)

	return b
}

func TestSelectingChars(t *testing.T) {
	b := makeBufferForTestingSelection()

	b.StartSelection(2, 0, SelectionChar)
	b.ExtendSelection(4, 1, true)

	assert.Equal(t, "e quick brown\nfox j", b.GetSelectedText())
}

func TestSelectingWordsDown(t *testing.T) {
	b := makeBufferForTestingSelection()

	b.StartSelection(6, 1, SelectionWord)
	b.ExtendSelection(5, 2, true)

	assert.Equal(t, "jumps over\nthe lazy", b.GetSelectedText())
}

func TestSelectingWordsUp(t *testing.T) {
	b := makeBufferForTestingSelection()

	b.StartSelection(5, 2, SelectionWord)
	b.ExtendSelection(6, 1, true)

	assert.Equal(t, "jumps over\nthe lazy", b.GetSelectedText())
}

func TestSelectingLinesDown(t *testing.T) {
	b := makeBufferForTestingSelection()

	b.StartSelection(6, 1, SelectionLine)
	b.ExtendSelection(4, 2, true)

	assert.Equal(t, "fox jumps over\nthe lazy dog", b.GetSelectedText())
}

func TestSelectingLineUp(t *testing.T) {
	b := makeBufferForTestingSelection()

	b.StartSelection(8, 2, SelectionLine)
	b.ExtendSelection(3, 1, true)

	assert.Equal(t, "fox jumps over\nthe lazy dog", b.GetSelectedText())
}

func TestSelectingAfterText(t *testing.T) {
	b := makeBufferForTestingSelection()

	b.StartSelection(6, 3, SelectionChar)
	b.ExtendSelection(6, 3, true)

	start, end := b.getActualSelection()

	assert.Equal(t, start.Col, 0)
	assert.Equal(t, start.Line, 3)
	assert.Equal(t, end.Col, 79)
	assert.Equal(t, end.Line, 3)
}
