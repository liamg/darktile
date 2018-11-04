package buffer

import (
	"fmt"
	"net/url"
	"time"
)

type Buffer struct {
	lines                 []Line
	cursorX               uint16
	cursorY               uint16
	viewHeight            uint16
	viewWidth             uint16
	cursorAttr            CellAttributes
	displayChangeHandlers []chan bool
	savedX                uint16
	savedY                uint16
	scrollLinesFromBottom uint
	topMargin             uint // see DECSTBM docs - this is for scrollable regions
	bottomMargin          uint // see DECSTBM docs - this is for scrollable regions
	replaceMode           bool // overwrite character at cursor or insert new
	autoWrap              bool
	dirty                 bool
	selectionStart        *Position
	selectionEnd          *Position
	selectionComplete     bool // whether the selected text can update or whether it is final
	selectionExpanded     bool // whether the selection to word expansion has already run on this point
	selectionClickTime    time.Time
}

type Position struct {
	Line int
	Col  int
}

// NewBuffer creates a new terminal buffer
func NewBuffer(viewCols uint16, viewLines uint16, attr CellAttributes) *Buffer {
	b := &Buffer{
		cursorX:    0,
		cursorY:    0,
		lines:      []Line{},
		cursorAttr: attr,
		autoWrap:   true,
	}
	b.SetVerticalMargins(0, uint(viewLines-1))
	b.ResizeView(viewCols, viewLines)
	return b
}

func (buffer *Buffer) GetURLAtPosition(col uint16, row uint16) string {

	cell := buffer.GetCell(col, row)
	if cell == nil || cell.Rune() == 0x00 {
		return ""
	}

	candidate := ""

	for i := col; i >= 0; i-- {
		cell := buffer.GetCell(i, row)
		if cell == nil {
			break
		}
		if isRuneURLSelectionMarker(cell.Rune()) {
			break
		}
		candidate = fmt.Sprintf("%c%s", cell.Rune(), candidate)
	}

	for i := col + 1; i < buffer.viewWidth; i++ {
		cell := buffer.GetCell(i, row)
		if cell == nil {
			break
		}
		if isRuneURLSelectionMarker(cell.Rune()) {
			break
		}
		candidate = fmt.Sprintf("%s%c", candidate, cell.Rune())
	}

	// check if url
	_, err := url.ParseRequestURI(candidate)
	if err != nil {
		return ""
	}
	return candidate
}

func (buffer *Buffer) SelectWordAtPosition(col uint16, row uint16) {

	cell := buffer.GetCell(col, row)
	if cell == nil || cell.Rune() == 0x00 {
		return
	}

	start := col
	end := col

	for i := col; i >= 0; i-- {
		cell := buffer.GetCell(i, row)
		if cell == nil {
			break
		}
		if isRuneWordSelectionMarker(cell.Rune()) {
			break
		}
		start = i
	}

	for i := col; i < buffer.viewWidth; i++ {
		cell := buffer.GetCell(i, row)
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
		Line: int(buffer.convertViewLineToRawLine(row)),
	}
	buffer.selectionEnd = &Position{
		Col:  int(end),
		Line: int(buffer.convertViewLineToRawLine(row)),
	}
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

	for row := buffer.selectionStart.Line; row <= buffer.selectionEnd.Line; row++ {

		if row >= len(buffer.lines) {
			break
		}

		line := buffer.lines[row]

		minX := 0
		maxX := int(buffer.viewWidth) - 1
		if row == buffer.selectionStart.Line {
			minX = buffer.selectionStart.Col
		} else if !line.wrapped {
			text += "\n"
		}
		if row == buffer.selectionEnd.Line {
			maxX = buffer.selectionEnd.Col
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

func (buffer *Buffer) StartSelection(col uint16, row uint16) {
	if buffer.selectionComplete {
		buffer.selectionEnd = nil

		if buffer.selectionStart != nil && time.Since(buffer.selectionClickTime) < time.Millisecond*500 {
			if buffer.selectionExpanded {
				//select whole line!
				buffer.selectionStart = &Position{
					Col:  0,
					Line: int(buffer.convertViewLineToRawLine(row)),
				}
				buffer.selectionEnd = &Position{
					Col:  int(buffer.ViewWidth() - 1),
					Line: int(buffer.convertViewLineToRawLine(row)),
				}
				buffer.emitDisplayChange()
			} else {
				buffer.SelectWordAtPosition(col, row)
				buffer.selectionExpanded = true
			}
			return
		}

		buffer.selectionExpanded = false
	}

	buffer.selectionComplete = false
	buffer.selectionStart = &Position{
		Col:  int(col),
		Line: int(buffer.convertViewLineToRawLine(row)),
	}
	buffer.selectionClickTime = time.Now()
}

func (buffer *Buffer) EndSelection(col uint16, row uint16, complete bool) {

	if buffer.selectionComplete {
		return
	}

	buffer.selectionComplete = complete

	defer buffer.emitDisplayChange()

	if buffer.selectionStart == nil {
		buffer.selectionEnd = nil
		return
	}

	if int(col) == buffer.selectionStart.Col && int(buffer.convertViewLineToRawLine(row)) == int(buffer.selectionStart.Line) && complete {
		return
	}

	buffer.selectionEnd = &Position{
		Col:  int(col),
		Line: int(buffer.convertViewLineToRawLine(row)),
	}
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

	rawY := int(buffer.convertViewLineToRawLine(row))
	return (rawY > y1 || (rawY == y1 && int(col) >= x1)) && (rawY < y2 || (rawY == y2 && int(col) <= x2))
}

func (buffer *Buffer) IsDirty() bool {
	if !buffer.dirty {
		return false
	}
	buffer.dirty = false
	return true
}

func (buffer *Buffer) SetAutoWrap(enabled bool) {
	buffer.autoWrap = enabled
}

func (buffer *Buffer) SetInsertMode() {
	buffer.replaceMode = false
}

func (buffer *Buffer) SetReplaceMode() {
	buffer.replaceMode = true
}

func (buffer *Buffer) SetVerticalMargins(top uint, bottom uint) {
	buffer.topMargin = top
	buffer.bottomMargin = bottom
}

func (buffer *Buffer) GetScrollOffset() uint {
	return buffer.scrollLinesFromBottom
}

func (buffer *Buffer) HasScrollableRegion() bool {
	return buffer.topMargin > 0 || buffer.bottomMargin < uint(buffer.ViewHeight())-1
}

func (buffer *Buffer) InScrollableRegion() bool {
	return buffer.HasScrollableRegion() && uint(buffer.cursorY) >= buffer.topMargin && uint(buffer.cursorY) <= buffer.bottomMargin
}

func (buffer *Buffer) ScrollDown(lines uint16) {

	defer buffer.emitDisplayChange()

	if buffer.Height() < int(buffer.ViewHeight()) {
		return
	}

	if uint(lines) > buffer.scrollLinesFromBottom {
		lines = uint16(buffer.scrollLinesFromBottom)
	}
	buffer.scrollLinesFromBottom -= uint(lines)
}

func (buffer *Buffer) ScrollUp(lines uint16) {

	defer buffer.emitDisplayChange()

	if buffer.Height() < int(buffer.ViewHeight()) {
		return
	}

	if uint(lines)+buffer.scrollLinesFromBottom >= (uint(buffer.Height()) - uint(buffer.ViewHeight())) {
		buffer.scrollLinesFromBottom = uint(buffer.Height()) - uint(buffer.ViewHeight())
	} else {
		buffer.scrollLinesFromBottom += uint(lines)
	}
}

func (buffer *Buffer) ScrollPageDown() {
	buffer.ScrollDown(buffer.viewHeight)
}
func (buffer *Buffer) ScrollPageUp() {
	buffer.ScrollUp(buffer.viewHeight)
}
func (buffer *Buffer) ScrollToEnd() {
	defer buffer.emitDisplayChange()
	buffer.scrollLinesFromBottom = 0
}

func (buffer *Buffer) SaveCursor() {
	buffer.savedX = buffer.cursorX
	buffer.savedY = buffer.cursorY
}

func (buffer *Buffer) RestoreCursor() {
	buffer.cursorX = buffer.savedX
	buffer.cursorY = buffer.savedY
}

func (buffer *Buffer) CursorAttr() *CellAttributes {
	return &buffer.cursorAttr
}

func (buffer *Buffer) GetCell(viewCol uint16, viewRow uint16) *Cell {

	rawLine := buffer.convertViewLineToRawLine(viewRow)

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
	return buffer.cursorX
}

// Line returns cursor line
func (buffer *Buffer) CursorLine() uint16 {
	return buffer.cursorY
}

func (buffer *Buffer) TopMargin() uint {
	return buffer.topMargin
}

func (buffer *Buffer) BottomMargin() uint {
	return buffer.bottomMargin
}

// translates the cursor line to the raw buffer line
func (buffer *Buffer) RawLine() uint64 {
	return buffer.convertViewLineToRawLine(buffer.cursorY)
}

func (buffer *Buffer) convertViewLineToRawLine(viewLine uint16) uint64 {
	rawHeight := buffer.Height()
	if int(buffer.viewHeight) > rawHeight {
		return uint64(viewLine)
	}
	return uint64(int(viewLine) + (rawHeight - int(buffer.viewHeight)))
}

func (buffer *Buffer) convertRawLineToViewLine(rawLine uint64) uint16 {
	rawHeight := buffer.Height()
	if int(buffer.viewHeight) > rawHeight {
		return uint16(rawLine)
	}
	return uint16(int(rawLine) - (rawHeight - int(buffer.viewHeight)))
}

// Width returns the width of the buffer in columns
func (buffer *Buffer) Width() uint16 {
	return buffer.viewWidth
}

func (buffer *Buffer) ViewWidth() uint16 {
	return buffer.viewWidth
}

func (buffer *Buffer) Height() int {
	return len(buffer.lines)
}

func (buffer *Buffer) ViewHeight() uint16 {
	return buffer.viewHeight
}

func (buffer *Buffer) insertLine() {

	defer buffer.emitDisplayChange()

	if !buffer.InScrollableRegion() {
		pos := buffer.RawLine()
		out := make([]Line, len(buffer.lines)+1)
		copy(out[:pos], buffer.lines[:pos])
		out[pos] = newLine()
		copy(out[pos+1:], buffer.lines[pos:])
		buffer.lines = out
	} else {
		topIndex := buffer.convertViewLineToRawLine(uint16(buffer.topMargin))
		bottomIndex := buffer.convertViewLineToRawLine(uint16(buffer.bottomMargin))
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

func (buffer *Buffer) InsertLines(count int) {

	if buffer.HasScrollableRegion() && !buffer.InScrollableRegion() {
		// should have no effect outside of scrollable region
		return
	}

	buffer.cursorX = 0

	for i := 0; i < count; i++ {
		buffer.insertLine()
	}

}

func (buffer *Buffer) Index() {

	// This sequence causes the active position to move downward one line without changing the column position.
	// If the active position is at the bottom margin, a scroll up is performed."

	defer buffer.emitDisplayChange()

	if buffer.InScrollableRegion() {

		if uint(buffer.cursorY) < buffer.bottomMargin {
			buffer.cursorY++
		} else {

			topIndex := buffer.convertViewLineToRawLine(uint16(buffer.topMargin))
			bottomIndex := buffer.convertViewLineToRawLine(uint16(buffer.bottomMargin))

			for i := topIndex; i < bottomIndex; i++ {
				buffer.lines[i] = buffer.lines[i+1]
			}

			buffer.lines[bottomIndex] = newLine()
		}

		return
	}

	if buffer.cursorY >= buffer.ViewHeight()-1 {
		buffer.lines = append(buffer.lines, newLine())
	} else {
		buffer.cursorY++
	}
}

func (buffer *Buffer) ReverseIndex() {

	defer buffer.emitDisplayChange()

	if buffer.InScrollableRegion() {

		if uint(buffer.cursorY) > buffer.topMargin {
			buffer.cursorY--
		} else {

			topIndex := buffer.convertViewLineToRawLine(uint16(buffer.topMargin))
			bottomIndex := buffer.convertViewLineToRawLine(uint16(buffer.bottomMargin))

			for i := bottomIndex; i > topIndex; i-- {
				buffer.lines[i] = buffer.lines[i-1]
			}

			buffer.lines[topIndex] = newLine()
		}
		return
	}

	if buffer.cursorY > 0 {
		buffer.cursorY--
	}
}

// Write will write a rune to the terminal at the position of the cursor, and increment the cursor position
func (buffer *Buffer) Write(runes ...rune) {

	// scroll to bottom on input
	inc := true
	buffer.scrollLinesFromBottom = 0

	for _, r := range runes {
		if r == 0x0a {
			buffer.NewLine()
			continue
		} else if r == 0x0d {
			buffer.CarriageReturn()
			continue
		} else if r == 0x9 {
			buffer.Tab()
			continue
		}
		line := buffer.getCurrentLine()

		if buffer.replaceMode {
			for int(buffer.CursorColumn()) >= len(line.cells) {
				line.cells = append(line.cells, NewBackgroundCell(buffer.cursorAttr.BgColour))
			}
			line.cells[buffer.cursorX].attr = buffer.cursorAttr
			line.cells[buffer.cursorX].setRune(r)
			buffer.incrementCursorPosition()
			continue
		}

		if buffer.CursorColumn() >= buffer.Width() { // if we're after the line, move to next

			if buffer.autoWrap {

				buffer.NewLine()

				newLine := buffer.getCurrentLine()
				newLine.setWrapped(true)
				if len(newLine.cells) == 0 {
					newLine.cells = []Cell{Cell{}}
				}
				cell := &newLine.cells[buffer.CursorColumn()]
				cell.setRune(r)
				cell.attr = buffer.cursorAttr

			} else {
				buffer.cursorX = buffer.Width() - 1
				inc = false
			}

			// @todo if next line is wrapped then prepend to it and shuffle characters along line, wrapping to next if necessary
		} else {

			for int(buffer.CursorColumn()) >= len(line.cells) {
				line.cells = append(line.cells, NewBackgroundCell(buffer.cursorAttr.BgColour))
			}

			cell := &line.cells[buffer.CursorColumn()]
			cell.setRune(r)
			cell.attr = buffer.cursorAttr

		}

		if inc {
			buffer.incrementCursorPosition()
		}
	}
}

func (buffer *Buffer) incrementCursorPosition() {

	defer buffer.emitDisplayChange()

	// we can increment one column past the end of the line.
	// this is effectively the beginning of the next line, except when we \r etc.
	if buffer.CursorColumn() < buffer.Width() { // if not at end of line

		buffer.cursorX++

	}
}

func (buffer *Buffer) Backspace() {

	if buffer.cursorX == 0 {
		line := buffer.getCurrentLine()
		if line.wrapped {
			buffer.MovePosition(int16(buffer.Width()-1), -1)
		} else {
			//@todo ring bell or whatever - actually i think the pty will trigger this
		}
	} else {
		buffer.MovePosition(-1, 0)
	}
}

func (buffer *Buffer) CarriageReturn() {
	defer buffer.emitDisplayChange()
	buffer.cursorX = 0
}

func (buffer *Buffer) Tab() {
	defer buffer.emitDisplayChange()
	tabSize := 4
	shift := int(buffer.cursorX-1) % tabSize
	if shift == 0 {
		shift = tabSize
	}
	for i := 0; i < shift; i++ {
		buffer.Write(' ')
	}
}

func (buffer *Buffer) NewLine() {
	defer buffer.emitDisplayChange()

	buffer.cursorX = 0
	buffer.Index()
}

func (buffer *Buffer) MovePosition(x int16, y int16) {

	var toX uint16
	var toY uint16

	if int16(buffer.cursorX)+x < 0 {
		toX = 0
	} else {
		toX = uint16(int16(buffer.cursorX) + x)
	}

	if int16(buffer.cursorY)+y < 0 {
		toY = 0
	} else {
		toY = uint16(int16(buffer.cursorY) + y)
	}

	buffer.SetPosition(toX, toY)
}

func (buffer *Buffer) SetPosition(col uint16, line uint16) {
	defer buffer.emitDisplayChange()

	if col >= buffer.ViewWidth() {
		col = buffer.ViewWidth() - 1
		//logrus.Errorf("Cannot set cursor position: column %d is outside of the current view width (%d columns)", col, buffer.ViewWidth())
	}
	if line >= buffer.ViewHeight() {
		line = buffer.ViewHeight() - 1
		//logrus.Errorf("Cannot set cursor position: line %d is outside of the current view height (%d lines)", line, buffer.ViewHeight())
	}

	buffer.cursorX = col
	buffer.cursorY = line
}

func (buffer *Buffer) GetVisibleLines() []Line {
	lines := []Line{}

	for i := buffer.Height() - int(buffer.ViewHeight()); i < buffer.Height(); i++ {
		y := i - int(buffer.scrollLinesFromBottom)
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
	return buffer.getViewLine(buffer.cursorY)
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
	for i := 0; i <= int(buffer.cursorX); i++ {
		if i < len(line.cells) {
			line.cells[i].erase()
		}
	}
}

func (buffer *Buffer) EraseLineFromCursor() {
	defer buffer.emitDisplayChange()
	line := buffer.getCurrentLine()

	if len(line.cells) > 0 {
		cx := buffer.cursorX
		if int(cx) < len(line.cells) {
			line.cells = line.cells[:buffer.cursorX]
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
	if int(buffer.cursorX) >= len(line.cells) {
		return
	}
	before := line.cells[:buffer.cursorX]
	if int(buffer.cursorX)+n >= len(line.cells) {
		n = len(line.cells) - int(buffer.cursorX)
	}
	after := line.cells[int(buffer.cursorX)+n:]
	line.cells = append(before, after...)
}

func (buffer *Buffer) EraseCharacters(n int) {
	defer buffer.emitDisplayChange()

	line := buffer.getCurrentLine()

	max := int(buffer.cursorX) + n
	if max > len(line.cells) {
		max = len(line.cells)
	}

	for i := int(buffer.cursorX); i < max; i++ {
		line.cells[i].erase()
	}
}

func (buffer *Buffer) EraseDisplayFromCursor() {
	defer buffer.emitDisplayChange()
	line := buffer.getCurrentLine()

	max := int(buffer.cursorX)
	if max > len(line.cells) {
		max = len(line.cells)
	}

	line.cells = line.cells[:max]
	for i := buffer.cursorY + 1; i < buffer.ViewHeight(); i++ {
		rawLine := buffer.convertViewLineToRawLine(i)
		if int(rawLine) < len(buffer.lines) {
			buffer.lines[int(rawLine)].cells = []Cell{}
		}
	}
}

func (buffer *Buffer) EraseDisplayToCursor() {
	defer buffer.emitDisplayChange()
	line := buffer.getCurrentLine()

	for i := 0; i < int(buffer.cursorX); i++ {
		line.cells[i].erase()
	}
	for i := uint16(0); i < buffer.cursorY; i++ {
		rawLine := buffer.convertViewLineToRawLine(i)
		if int(rawLine) < len(buffer.lines) {
			buffer.lines[int(rawLine)].cells = []Cell{}
		}
	}
}

func (buffer *Buffer) ResizeView(width uint16, height uint16) {

	defer buffer.emitDisplayChange()

	if buffer.viewHeight == 0 {
		buffer.viewWidth = width
		buffer.viewHeight = height
		return
	}

	// @todo scroll to bottom on resize
	line := buffer.getCurrentLine()
	cXFromEndOfLine := len(line.cells) - int(buffer.cursorX+1)

	cursorYMovement := 0

	if width < buffer.viewWidth { // wrap lines if we're shrinking
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

				if i+1 <= int(buffer.cursorY) {
					cursorYMovement++
				}

				newLine := newLine()
				newLine.setWrapped(true)
				newLine.cells = sillyCells
				after := append([]Line{newLine}, buffer.lines[i+1:]...)
				buffer.lines = append(buffer.lines[:i+1], after...)

			}
		}
	} else if width > buffer.viewWidth { // unwrap lines if we're growing
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

					if i+offset <= int(buffer.cursorY) {
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

	buffer.viewWidth = width
	buffer.viewHeight = height

	cY := uint16(len(buffer.lines) - 1)
	if cY >= buffer.viewHeight {
		cY = buffer.viewHeight - 1
	}
	buffer.cursorY = cY

	// position cursorX
	line = buffer.getCurrentLine()
	buffer.cursorX = uint16((len(line.cells) - cXFromEndOfLine) - 1)

	buffer.SetVerticalMargins(0, uint(buffer.viewHeight-1))
}
