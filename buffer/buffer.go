package buffer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
)

type Buffer struct {
	lines                 []Line
	displayChangeHandlers []chan bool
	savedX                uint16
	savedY                uint16
	dirty                 bool
	selectionStart        *Position
	selectionEnd          *Position
	selectionComplete     bool // whether the selected text can update or whether it is final
	terminalState         *TerminalState
}

type Position struct {
	Line int
	Col  int
}

// NewBuffer creates a new terminal buffer
func NewBuffer(terminalState *TerminalState) *Buffer {
	b := &Buffer{
		lines:         []Line{},
		terminalState: terminalState,
	}
	return b
}

func (buffer *Buffer) GetURLAtPosition(col uint16, viewRow uint16) string {

	row := buffer.convertViewLineToRawLine((viewRow)) - uint64(buffer.terminalState.scrollLinesFromBottom)

	cell := buffer.GetRawCell(col, row)
	if cell == nil || cell.Rune() == 0x00 {
		return ""
	}

	candidate := ""

	for i := col; i >= uint16(0); i-- {
		cell := buffer.GetRawCell(i, row)
		if cell == nil {
			break
		}
		if isRuneURLSelectionMarker(cell.Rune()) {
			break
		}
		candidate = fmt.Sprintf("%c%s", cell.Rune(), candidate)
	}

	for i := col + 1; i < buffer.terminalState.viewWidth; i++ {
		cell := buffer.GetRawCell(i, row)
		if cell == nil {
			break
		}
		if isRuneURLSelectionMarker(cell.Rune()) {
			break
		}
		candidate = fmt.Sprintf("%s%c", candidate, cell.Rune())
	}

	if candidate == "" || candidate[0] == '/' {
		return ""
	}

	// check if url
	_, err := url.ParseRequestURI(candidate)
	if err != nil {
		return ""
	}
	return candidate
}

func (buffer *Buffer) IsSelectionComplete() bool {
	return buffer.selectionComplete
}

func (buffer *Buffer) SelectLineAtPosition(col uint16, viewRow uint16) {
	row := buffer.convertViewLineToRawLine(viewRow) - uint64(buffer.terminalState.scrollLinesFromBottom)

	buffer.selectionStart = &Position {
		Col: 0,
		Line: int(row),
	}
	buffer.selectionEnd = &Position {
		Col: int(buffer.ViewWidth() - 1),
		Line: int(row),
	}

	buffer.selectionComplete = true
	buffer.emitDisplayChange()
}

func (buffer *Buffer) SelectWordAtPosition(col uint16, viewRow uint16) {

	row := buffer.convertViewLineToRawLine(viewRow) - uint64(buffer.terminalState.scrollLinesFromBottom)

	cell := buffer.GetRawCell(col, row)
	if cell == nil || cell.Rune() == 0x00 {
		return
	}

	start := col
	end := col

	for i := col; i >= uint16(0); i-- {
		cell := buffer.GetRawCell(i, row)
		if cell == nil {
			break
		}
		if isRuneWordSelectionMarker(cell.Rune()) {
			break
		}
		start = i
	}

	for i := col; i < buffer.terminalState.viewWidth; i++ {
		cell := buffer.GetRawCell(i, row)
		if cell == nil {
			break
		}
		if isRuneWordSelectionMarker(cell.Rune()) {
			break
		}
		end = i
	}

	buffer.selectionStart = &Position{
		Col:  int(start),
		Line: int(row),
	}
	buffer.selectionEnd = &Position{
		Col:  int(end),
		Line: int(row),
	}

	buffer.selectionComplete = true
	buffer.emitDisplayChange()
}

// bounds for word selection
func isRuneWordSelectionMarker(r rune) bool {
	switch r {
	case ',', ' ', ':', ';', 0, '\'', '"', '[', ']', '(', ')', '{', '}':
		return true
	}

	return false
}

func isRuneURLSelectionMarker(r rune) bool {
	switch r {
	case ' ', 0, '\'', '"', '{', '}':
		return true
	}

	return false
}

func (buffer *Buffer) GetSelectedText() string {
	if buffer.selectionStart == nil || buffer.selectionEnd == nil {
		return ""
	}

	text := ""

	var x1, x2, y1, y2 int

	if buffer.selectionStart.Line > buffer.selectionEnd.Line || (buffer.selectionStart.Line == buffer.selectionEnd.Line && buffer.selectionStart.Col > buffer.selectionEnd.Col) {
		y2 = buffer.selectionStart.Line
		y1 = buffer.selectionEnd.Line
		x2 = buffer.selectionStart.Col
		x1 = buffer.selectionEnd.Col
	} else {
		y1 = buffer.selectionStart.Line
		y2 = buffer.selectionEnd.Line
		x1 = buffer.selectionStart.Col
		x2 = buffer.selectionEnd.Col
	}

	for row := y1; row <= y2; row++ {

		if row >= len(buffer.lines) {
			break
		}

		line := buffer.lines[row]

		minX := 0
		maxX := int(buffer.terminalState.viewWidth) - 1
		if row == y1 {
			minX = x1
		} else if !line.wrapped {
			text += "\n"
		}
		if row == y2 {
			maxX = x2
		}

		for col := minX; col <= maxX; col++ {
			if col >= len(line.cells) {
				break
			}
			cell := line.cells[col]
			text += string(cell.Rune())
		}

	}

	return text
}

func (buffer *Buffer) StartSelection(col uint16, viewRow uint16) {
	row := buffer.convertViewLineToRawLine(viewRow) - uint64(buffer.terminalState.scrollLinesFromBottom)
	buffer.selectionComplete = false

	buffer.selectionStart = &Position {
		Col:  int(col),
		Line: int(row),
	}

	buffer.selectionEnd = nil
}

func (buffer *Buffer) EndSelection(col uint16, viewRow uint16, complete bool) {

	if buffer.selectionComplete {
		return
	}

	buffer.selectionComplete = complete

	defer buffer.emitDisplayChange()

	if buffer.selectionStart == nil {
		buffer.selectionEnd = nil
		return
	}

	row := buffer.convertViewLineToRawLine(viewRow) - uint64(buffer.terminalState.scrollLinesFromBottom)

	if int(col) == buffer.selectionStart.Col && int(row) == int(buffer.selectionStart.Line) && complete {
		return
	}

	buffer.selectionEnd = &Position{
		Col:  int(col),
		Line: int(row),
	}
}

func (buffer *Buffer) ClearSelection() {
	buffer.selectionStart = nil
	buffer.selectionEnd = nil
	buffer.selectionComplete = true

	buffer.emitDisplayChange()
}

func (buffer *Buffer) InSelection(col uint16, row uint16) bool {

	if buffer.selectionStart == nil || buffer.selectionEnd == nil {
		return false
	}

	var x1, x2, y1, y2 int

	// first, let's put the selection points in the correct order, earliest first
	if buffer.selectionStart.Line > buffer.selectionEnd.Line || (buffer.selectionStart.Line == buffer.selectionEnd.Line && buffer.selectionStart.Col > buffer.selectionEnd.Col) {
		y2 = buffer.selectionStart.Line
		y1 = buffer.selectionEnd.Line
		x2 = buffer.selectionStart.Col
		x1 = buffer.selectionEnd.Col
	} else {
		y1 = buffer.selectionStart.Line
		y2 = buffer.selectionEnd.Line
		x1 = buffer.selectionStart.Col
		x2 = buffer.selectionEnd.Col
	}

	rawY := int(buffer.convertViewLineToRawLine(row) - uint64(buffer.terminalState.scrollLinesFromBottom))
	return (rawY > y1 || (rawY == y1 && int(col) >= x1)) && (rawY < y2 || (rawY == y2 && int(col) <= x2))
}

func (buffer *Buffer) IsDirty() bool {
	if !buffer.dirty {
		return false
	}
	buffer.dirty = false
	return true
}

func (buffer *Buffer) GetScrollOffset() uint {
	return buffer.terminalState.scrollLinesFromBottom
}

func (buffer *Buffer) HasScrollableRegion() bool {
	return buffer.terminalState.topMargin > 0 || buffer.terminalState.bottomMargin < uint(buffer.ViewHeight())-1
}

func (buffer *Buffer) InScrollableRegion() bool {
	return buffer.HasScrollableRegion() && uint(buffer.terminalState.cursorY) >= buffer.terminalState.topMargin && uint(buffer.terminalState.cursorY) <= buffer.terminalState.bottomMargin
}

func (buffer *Buffer) ScrollDown(lines uint16) {

	defer buffer.emitDisplayChange()

	if buffer.Height() < int(buffer.ViewHeight()) {
		return
	}

	if uint(lines) > buffer.terminalState.scrollLinesFromBottom {
		lines = uint16(buffer.terminalState.scrollLinesFromBottom)
	}
	buffer.terminalState.scrollLinesFromBottom -= uint(lines)
}

func (buffer *Buffer) ScrollUp(lines uint16) {

	defer buffer.emitDisplayChange()

	if buffer.Height() < int(buffer.ViewHeight()) {
		return
	}

	if uint(lines)+buffer.terminalState.scrollLinesFromBottom >= (uint(buffer.Height()) - uint(buffer.ViewHeight())) {
		buffer.terminalState.scrollLinesFromBottom = uint(buffer.Height()) - uint(buffer.ViewHeight())
	} else {
		buffer.terminalState.scrollLinesFromBottom += uint(lines)
	}
}

func (buffer *Buffer) ScrollPageDown() {
	buffer.ScrollDown(buffer.terminalState.viewHeight)
}
func (buffer *Buffer) ScrollPageUp() {
	buffer.ScrollUp(buffer.terminalState.viewHeight)
}
func (buffer *Buffer) ScrollToEnd() {
	defer buffer.emitDisplayChange()
	buffer.terminalState.scrollLinesFromBottom = 0
}

func (buffer *Buffer) SaveCursor() {
	buffer.savedX = buffer.terminalState.cursorX
	buffer.savedY = buffer.terminalState.cursorY
}

func (buffer *Buffer) RestoreCursor() {
	buffer.terminalState.cursorX = buffer.savedX
	buffer.terminalState.cursorY = buffer.savedY
}

func (buffer *Buffer) CursorAttr() *CellAttributes {
	return &buffer.terminalState.cursorAttr
}

func (buffer *Buffer) GetCell(viewCol uint16, viewRow uint16) *Cell {
	rawLine := buffer.convertViewLineToRawLine(viewRow)
	return buffer.GetRawCell(viewCol, rawLine)
}

func (buffer *Buffer) GetRawCell(viewCol uint16, rawLine uint64) *Cell {

	if viewCol < 0 || rawLine < 0 || int(rawLine) >= len(buffer.lines) {
		return nil
	}
	line := &buffer.lines[rawLine]
	if int(viewCol) >= len(line.cells) {
		return nil
	}
	return &line.cells[viewCol]
}

func (buffer *Buffer) emitDisplayChange() {
	buffer.dirty = true
}

// Column returns cursor column
func (buffer *Buffer) CursorColumn() uint16 {
	// @todo originMode and left margin
	return buffer.terminalState.cursorX
}

// Line returns cursor line
func (buffer *Buffer) CursorLine() uint16 {
	if buffer.terminalState.OriginMode {
		result := buffer.terminalState.cursorY - uint16(buffer.terminalState.topMargin)
		if result < 0 {
			result = 0
		}
		return result
	}
	return buffer.terminalState.cursorY
}

func (buffer *Buffer) TopMargin() uint {
	return buffer.terminalState.topMargin
}

func (buffer *Buffer) BottomMargin() uint {
	return buffer.terminalState.bottomMargin
}

// translates the cursor line to the raw buffer line
func (buffer *Buffer) RawLine() uint64 {
	return buffer.convertViewLineToRawLine(buffer.terminalState.cursorY)
}

func (buffer *Buffer) convertViewLineToRawLine(viewLine uint16) uint64 {
	rawHeight := buffer.Height()
	if int(buffer.terminalState.viewHeight) > rawHeight {
		return uint64(viewLine)
	}
	return uint64(int(viewLine) + (rawHeight - int(buffer.terminalState.viewHeight)))
}

func (buffer *Buffer) convertRawLineToViewLine(rawLine uint64) uint16 {
	rawHeight := buffer.Height()
	if int(buffer.terminalState.viewHeight) > rawHeight {
		return uint16(rawLine)
	}
	return uint16(int(rawLine) - (rawHeight - int(buffer.terminalState.viewHeight)))
}

// Width returns the width of the buffer in columns
func (buffer *Buffer) Width() uint16 {
	return buffer.terminalState.viewWidth
}

func (buffer *Buffer) ViewWidth() uint16 {
	return buffer.terminalState.viewWidth
}

func (buffer *Buffer) Height() int {
	return len(buffer.lines)
}

func (buffer *Buffer) ViewHeight() uint16 {
	return buffer.terminalState.viewHeight
}

func (buffer *Buffer) deleteLine() {
	index := int(buffer.RawLine())
	buffer.lines = buffer.lines[:index+copy(buffer.lines[index:], buffer.lines[index+1:])]
}

func (buffer *Buffer) insertLine() {

	defer buffer.emitDisplayChange()

	if !buffer.InScrollableRegion() {
		pos := buffer.RawLine()
		maxLines := buffer.getMaxLines()
		newLineCount := uint64(len(buffer.lines) + 1)
		if newLineCount > maxLines {
			newLineCount = maxLines
		}

		out := make([]Line, newLineCount)
		copy(
			out[:pos-(uint64(len(buffer.lines))+1-newLineCount)],
			buffer.lines[uint64(len(buffer.lines))+1-newLineCount:pos])
		out[pos] = newLine()
		copy(out[pos+1:], buffer.lines[pos:])
		buffer.lines = out
	} else {
		topIndex := buffer.convertViewLineToRawLine(uint16(buffer.terminalState.topMargin))
		bottomIndex := buffer.convertViewLineToRawLine(uint16(buffer.terminalState.bottomMargin))
		before := buffer.lines[:topIndex]
		after := buffer.lines[bottomIndex+1:]
		out := make([]Line, len(buffer.lines))
		copy(out[0:], before)

		pos := buffer.RawLine()
		for i := topIndex; i < bottomIndex; i++ {
			if i < pos {
				out[i] = buffer.lines[i]
			} else {
				out[i+1] = buffer.lines[i]
			}
		}

		copy(out[bottomIndex+1:], after)

		out[pos] = newLine()
		buffer.lines = out
	}
}

func (buffer *Buffer) InsertBlankCharacters(count int) {

	index := int(buffer.RawLine())
	for i := 0; i < count; i++ {
		cells := buffer.lines[index].cells
		buffer.lines[index].cells = append(cells[:buffer.terminalState.cursorX], append([]Cell{buffer.terminalState.defaultCell}, cells[buffer.terminalState.cursorX:]...)...)
	}
}

func (buffer *Buffer) InsertLines(count int) {

	if buffer.HasScrollableRegion() && !buffer.InScrollableRegion() {
		// should have no effect outside of scrollable region
		return
	}

	buffer.terminalState.cursorX = 0

	for i := 0; i < count; i++ {
		buffer.insertLine()
	}

}

func (buffer *Buffer) DeleteLines(count int) {

	if buffer.HasScrollableRegion() && !buffer.InScrollableRegion() {
		// should have no effect outside of scrollable region
		return
	}

	buffer.terminalState.cursorX = 0

	for i := 0; i < count; i++ {
		buffer.deleteLine()
	}

}

func (buffer *Buffer) Index() {

	// This sequence causes the active position to move downward one line without changing the column position.
	// If the active position is at the bottom margin, a scroll up is performed."

	defer buffer.emitDisplayChange()

	if buffer.InScrollableRegion() {

		if uint(buffer.terminalState.cursorY) < buffer.terminalState.bottomMargin {
			buffer.terminalState.cursorY++
		} else {

			topIndex := buffer.convertViewLineToRawLine(uint16(buffer.terminalState.topMargin))
			bottomIndex := buffer.convertViewLineToRawLine(uint16(buffer.terminalState.bottomMargin))

			for i := topIndex; i < bottomIndex; i++ {
				buffer.lines[i] = buffer.lines[i+1]
			}

			buffer.lines[bottomIndex] = newLine()
		}

		return
	}

	if buffer.terminalState.cursorY >= buffer.ViewHeight()-1 {
		buffer.lines = append(buffer.lines, newLine())
		maxLines := buffer.getMaxLines()
		if uint64(len(buffer.lines)) > maxLines {
			copy(buffer.lines, buffer.lines[uint64(len(buffer.lines))-maxLines:])
			buffer.lines = buffer.lines[:maxLines]
		}
	} else {
		buffer.terminalState.cursorY++
	}
}

func (buffer *Buffer) ReverseIndex() {

	defer buffer.emitDisplayChange()

	if buffer.InScrollableRegion() {

		if uint(buffer.terminalState.cursorY) > buffer.terminalState.topMargin {
			buffer.terminalState.cursorY--
		} else {

			topIndex := buffer.convertViewLineToRawLine(uint16(buffer.terminalState.topMargin))
			bottomIndex := buffer.convertViewLineToRawLine(uint16(buffer.terminalState.bottomMargin))

			for i := bottomIndex; i > topIndex; i-- {
				buffer.lines[i] = buffer.lines[i-1]
			}

			buffer.lines[topIndex] = newLine()
		}
		return
	}

	if buffer.terminalState.cursorY > 0 {
		buffer.terminalState.cursorY--
	}
}

// Write will write a rune to the terminal at the position of the cursor, and increment the cursor position
func (buffer *Buffer) Write(runes ...rune) {

	// scroll to bottom on input
	buffer.terminalState.scrollLinesFromBottom = 0

	for _, r := range runes {

		line := buffer.getCurrentLine()

		if buffer.terminalState.ReplaceMode {

			if buffer.CursorColumn() >= buffer.Width() {
				// @todo replace rune at position 0 on next line down
				return
			}

			for int(buffer.CursorColumn()) >= len(line.cells) {
				line.cells = append(line.cells, buffer.terminalState.defaultCell)
			}
			line.cells[buffer.terminalState.cursorX].attr = buffer.terminalState.cursorAttr
			line.cells[buffer.terminalState.cursorX].setRune(r)
			buffer.incrementCursorPosition()
			continue
		}

		if buffer.CursorColumn() >= buffer.Width() { // if we're after the line, move to next

			if buffer.terminalState.AutoWrap {

				buffer.NewLineEx(true)

				newLine := buffer.getCurrentLine()
				if len(newLine.cells) == 0 {
					newLine.cells = append(newLine.cells, buffer.terminalState.defaultCell)
				}
				cell := &newLine.cells[0]
				cell.setRune(r)
				cell.attr = buffer.terminalState.cursorAttr

			} else {
				// no more room on line and wrapping is disabled
				return
			}

			// @todo if next line is wrapped then prepend to it and shuffle characters along line, wrapping to next if necessary
		} else {

			for int(buffer.CursorColumn()) >= len(line.cells) {
				line.cells = append(line.cells, buffer.terminalState.defaultCell)
			}

			cell := &line.cells[buffer.CursorColumn()]
			cell.setRune(r)
			cell.attr = buffer.terminalState.cursorAttr
		}

		buffer.incrementCursorPosition()
	}
}

func (buffer *Buffer) incrementCursorPosition() {
	// we can increment one column past the end of the line.
	// this is effectively the beginning of the next line, except when we \r etc.
	if buffer.CursorColumn() < buffer.Width() {
		buffer.terminalState.cursorX++
	}
}

func (buffer *Buffer) inDoWrap() bool {
	// xterm uses 'do_wrap' flag for this special terminal state
	// we use the cursor position right after the boundary
	// let's see how it works out
	return buffer.terminalState.cursorX == buffer.terminalState.viewWidth // @todo rightMargin
}

func (buffer *Buffer) Backspace() {

	if buffer.terminalState.cursorX == 0 {
		line := buffer.getCurrentLine()
		if line.wrapped {
			buffer.MovePosition(int16(buffer.Width()-1), -1)
		} else {
			//@todo ring bell or whatever - actually i think the pty will trigger this
		}
	} else if buffer.inDoWrap() {
		// the "do_wrap" implementation
		buffer.MovePosition(-2, 0)
	} else {
		buffer.MovePosition(-1, 0)
	}
}

func (buffer *Buffer) CarriageReturn() {

	for {
		line := buffer.getCurrentLine()
		if line == nil {
			break
		}
		if line.wrapped && buffer.terminalState.cursorY > 0 {
			buffer.terminalState.cursorY--
		} else {
			break
		}
	}

	buffer.terminalState.cursorX = 0
}

func (buffer *Buffer) Tab() {
	tabSize := 4
	max := tabSize

	// @todo rightMargin
	if buffer.terminalState.cursorX < buffer.terminalState.viewWidth {
		max = int(buffer.terminalState.viewWidth - buffer.terminalState.cursorX - 1)
	}

	shift := tabSize - (int(buffer.terminalState.cursorX+1) % tabSize)

	if shift > max {
		shift = max
	}

	for i := 0; i < shift; i++ {
		buffer.Write(' ')
	}
}

func (buffer *Buffer) NewLine() {
	buffer.NewLineEx(false)
}

func (buffer *Buffer) NewLineEx(forceCursorToMargin bool) {

	if buffer.terminalState.IsNewLineMode() || forceCursorToMargin {
		buffer.terminalState.cursorX = 0
	}
	buffer.Index()

	for {
		line := buffer.getCurrentLine()
		if !line.wrapped {
			break
		}
		buffer.Index()
	}
}

func (buffer *Buffer) IsNewLineMode() bool {
	return buffer.terminalState.LineFeedMode == false
}

func (buffer *Buffer) MovePosition(x int16, y int16) {

	var toX uint16
	var toY uint16

	if int16(buffer.CursorColumn())+x < 0 {
		toX = 0
	} else {
		toX = uint16(int16(buffer.CursorColumn()) + x)
	}

	// should either use CursorLine() and SetPosition() or use absolutes, mind Origin Mode (DECOM)
	if int16(buffer.CursorLine())+y < 0 {
		toY = 0
	} else {
		toY = uint16(int16(buffer.CursorLine()) + y)
	}

	buffer.SetPosition(toX, toY)
}

func (buffer *Buffer) SetPosition(col uint16, line uint16) {
	defer buffer.emitDisplayChange()

	useCol := col
	useLine := line
	maxLine := buffer.ViewHeight() - 1

	if buffer.terminalState.OriginMode {
		useLine += uint16(buffer.terminalState.topMargin)
		maxLine = uint16(buffer.terminalState.bottomMargin)
		// @todo left and right margins
	}
	if useLine > maxLine {
		useLine = maxLine
	}

	if useCol >= buffer.ViewWidth() {
		useCol = buffer.ViewWidth() - 1
		//logrus.Errorf("Cannot set cursor position: column %d is outside of the current view width (%d columns)", col, buffer.ViewWidth())
	}

	buffer.terminalState.cursorX = useCol
	buffer.terminalState.cursorY = useLine
}

func (buffer *Buffer) GetVisibleLines() []Line {
	lines := []Line{}

	for i := buffer.Height() - int(buffer.ViewHeight()); i < buffer.Height(); i++ {
		y := i - int(buffer.terminalState.scrollLinesFromBottom)
		if y >= 0 && y < len(buffer.lines) {
			lines = append(lines, buffer.lines[y])
		}
	}
	return lines
}

// tested to here

func (buffer *Buffer) Clear() {
	defer buffer.emitDisplayChange()
	for i := 0; i < int(buffer.ViewHeight()); i++ {
		buffer.lines = append(buffer.lines, newLine())
	}
	buffer.SetPosition(0, 0) // do we need to set position?
}

// creates if necessary
func (buffer *Buffer) getCurrentLine() *Line {
	return buffer.getViewLine(buffer.terminalState.cursorY)
}

func (buffer *Buffer) getViewLine(index uint16) *Line {

	if index >= buffer.ViewHeight() { // @todo is this okay?#
		return &buffer.lines[len(buffer.lines)-1]
	}

	if len(buffer.lines) < int(buffer.ViewHeight()) {
		for int(index) >= len(buffer.lines) {
			buffer.lines = append(buffer.lines, newLine())
		}
		return &buffer.lines[int(index)]
	}

	if int(buffer.convertViewLineToRawLine(index)) < len(buffer.lines) {
		return &buffer.lines[buffer.convertViewLineToRawLine(index)]
	}

	panic(fmt.Sprintf("Failed to retrieve line for %d", index))
}

func (buffer *Buffer) EraseLine() {
	defer buffer.emitDisplayChange()
	line := buffer.getCurrentLine()
	line.cells = []Cell{}
}

func (buffer *Buffer) EraseLineToCursor() {
	defer buffer.emitDisplayChange()
	line := buffer.getCurrentLine()
	for i := 0; i <= int(buffer.terminalState.cursorX); i++ {
		if i < len(line.cells) {
			line.cells[i].erase(buffer.terminalState.defaultCell.attr.BgColour)
		}
	}
}

func (buffer *Buffer) EraseLineFromCursor() {
	defer buffer.emitDisplayChange()
	line := buffer.getCurrentLine()

	if len(line.cells) > 0 {
		cx := buffer.terminalState.cursorX
		if int(cx) < len(line.cells) {
			line.cells = line.cells[:buffer.terminalState.cursorX]
		}
	}

	max := int(buffer.ViewWidth()) - len(line.cells)

	buffer.SaveCursor()
	for i := 0; i < max; i++ {
		buffer.Write(0)
	}
	buffer.RestoreCursor()
}

func (buffer *Buffer) EraseDisplay() {
	defer buffer.emitDisplayChange()
	for i := uint16(0); i < (buffer.ViewHeight()); i++ {
		rawLine := buffer.convertViewLineToRawLine(i)
		if int(rawLine) < len(buffer.lines) {
			buffer.lines[int(rawLine)].cells = []Cell{}
		}
	}
}

func (buffer *Buffer) DeleteChars(n int) {
	defer buffer.emitDisplayChange()

	line := buffer.getCurrentLine()
	if int(buffer.terminalState.cursorX) >= len(line.cells) {
		return
	}
	before := line.cells[:buffer.terminalState.cursorX]
	if int(buffer.terminalState.cursorX)+n >= len(line.cells) {
		n = len(line.cells) - int(buffer.terminalState.cursorX)
	}
	after := line.cells[int(buffer.terminalState.cursorX)+n:]
	line.cells = append(before, after...)
}

func (buffer *Buffer) EraseCharacters(n int) {
	defer buffer.emitDisplayChange()

	line := buffer.getCurrentLine()

	max := int(buffer.terminalState.cursorX) + n
	if max > len(line.cells) {
		max = len(line.cells)
	}

	for i := int(buffer.terminalState.cursorX); i < max; i++ {
		line.cells[i].erase(buffer.terminalState.defaultCell.attr.BgColour)
	}
}

func (buffer *Buffer) EraseDisplayFromCursor() {
	defer buffer.emitDisplayChange()
	line := buffer.getCurrentLine()

	max := int(buffer.terminalState.cursorX)
	if max > len(line.cells) {
		max = len(line.cells)
	}

	line.cells = line.cells[:max]
	for i := buffer.terminalState.cursorY + 1; i < buffer.ViewHeight(); i++ {
		rawLine := buffer.convertViewLineToRawLine(i)
		if int(rawLine) < len(buffer.lines) {
			buffer.lines[int(rawLine)].cells = []Cell{}
		}
	}
}

func (buffer *Buffer) EraseDisplayToCursor() {
	defer buffer.emitDisplayChange()
	line := buffer.getCurrentLine()

	for i := 0; i <= int(buffer.terminalState.cursorX); i++ {
		if i >= len(line.cells) {
			break
		}
		line.cells[i].erase(buffer.terminalState.defaultCell.attr.BgColour)
	}
	for i := uint16(0); i < buffer.terminalState.cursorY; i++ {
		rawLine := buffer.convertViewLineToRawLine(i)
		if int(rawLine) < len(buffer.lines) {
			buffer.lines[int(rawLine)].cells = []Cell{}
		}
	}
}

func (buffer *Buffer) ResizeView(width uint16, height uint16) {

	defer buffer.emitDisplayChange()

	if buffer.terminalState.viewHeight == 0 {
		buffer.terminalState.viewWidth = width
		buffer.terminalState.viewHeight = height
		return
	}

	// @todo scroll to bottom on resize
	line := buffer.getCurrentLine()
	cXFromEndOfLine := len(line.cells) - int(buffer.terminalState.cursorX+1)

	cursorYMovement := 0

	if width < buffer.terminalState.viewWidth { // wrap lines if we're shrinking
		for i := 0; i < len(buffer.lines); i++ {
			line := &buffer.lines[i]
			//line.Cleanse()
			if len(line.cells) > int(width) { // only try wrapping a line if it's too long
				sillyCells := line.cells[width:] // grab the cells we need to wrap
				line.cells = line.cells[:width]

				// we need to move cut cells to the next line
				// if the next line is wrapped anyway, we can push them onto the beginning of that line
				// otherwise, we need add a new wrapped line
				if i+1 < len(buffer.lines) {
					nextLine := &buffer.lines[i+1]
					if nextLine.wrapped {

						nextLine.cells = append(sillyCells, nextLine.cells...)
						continue
					}
				}

				if i+1 <= int(buffer.terminalState.cursorY) {
					cursorYMovement++
				}

				newLine := newLine()
				newLine.setWrapped(true)
				newLine.cells = sillyCells
				after := append([]Line{newLine}, buffer.lines[i+1:]...)
				buffer.lines = append(buffer.lines[:i+1], after...)

			}
		}
	} else if width > buffer.terminalState.viewWidth { // unwrap lines if we're growing
		for i := 0; i < len(buffer.lines)-1; i++ {
			line := &buffer.lines[i]
			//line.Cleanse()
			for offset := 1; i+offset < len(buffer.lines); offset++ {
				nextLine := &buffer.lines[i+offset]
				//nextLine.Cleanse()
				if !nextLine.wrapped { // if the next line wasn't wrapped, we don't need to move characters back to this line
					break
				}
				spaceOnLine := int(width) - len(line.cells)
				if spaceOnLine <= 0 { // no more space to unwrap
					break
				}
				moveCount := spaceOnLine
				if moveCount > len(nextLine.cells) {
					moveCount = len(nextLine.cells)
				}
				line.cells = append(line.cells, nextLine.cells[:moveCount]...)
				if moveCount == len(nextLine.cells) {

					if i+offset <= int(buffer.terminalState.cursorY) {
						cursorYMovement--
					}

					// if we unwrapped all cells off the next line, delete it
					buffer.lines = append(buffer.lines[:i+offset], buffer.lines[i+offset+1:]...)

					offset--

				} else {
					// otherwise just remove the characters we moved up a line
					nextLine.cells = nextLine.cells[moveCount:]
				}
			}

		}
	}

	buffer.terminalState.viewWidth = width
	buffer.terminalState.viewHeight = height

	cY := uint16(len(buffer.lines) - 1)
	if cY >= buffer.terminalState.viewHeight {
		cY = buffer.terminalState.viewHeight - 1
	}
	buffer.terminalState.cursorY = cY

	// position cursorX
	line = buffer.getCurrentLine()
	buffer.terminalState.cursorX = uint16((len(line.cells) - cXFromEndOfLine) - 1)

	buffer.terminalState.ResetVerticalMargins()
}

func (buffer *Buffer) getMaxLines() uint64 {
	result := buffer.terminalState.maxLines
	if result < uint64(buffer.terminalState.viewHeight) {
		result = uint64(buffer.terminalState.viewHeight)
	}

	return result
}

func (buffer *Buffer) Save(path string) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for _, line := range buffer.lines {
		f.WriteString(line.String())
	}
}

func (buffer *Buffer) Compare(path string) bool {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	bufferContent := []byte{}
	for _, line := range buffer.lines {
		lineBytes := []byte(line.String())
		bufferContent = append(bufferContent, lineBytes...)
	}
	return bytes.Equal(f, bufferContent)
}
