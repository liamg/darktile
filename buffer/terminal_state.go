package buffer

type TerminalState struct {
	scrollLinesFromBottom uint
	cursorX               uint16
	cursorY               uint16
	CursorAttr            CellAttributes
	viewHeight            uint16
	viewWidth             uint16
	topMargin             uint // see DECSTBM docs - this is for scrollable regions
	bottomMargin          uint // see DECSTBM docs - this is for scrollable regions
	ReplaceMode           bool // overwrite character at cursor or insert new
	OriginMode            bool // see DECOM docs - whether cursor is positioned within the margins or not
	LineFeedMode          bool
	ScreenMode            bool // DECSCNM (black on white background)
	AutoWrap              bool
	maxLines              uint64
	tabStops              map[uint16]struct{}
	Charsets              []*map[rune]rune // array of 2 charsets, nil means ASCII (no conversion)
	CurrentCharset        int              // active charset index in Charsets array, valid values are 0 or 1
}

// NewTerminalState creates a new terminal state
func NewTerminalState(viewCols uint16, viewLines uint16, attr CellAttributes, maxLines uint64) *TerminalState {
	b := &TerminalState{
		cursorX:      0,
		cursorY:      0,
		CursorAttr:   attr,
		AutoWrap:     true,
		maxLines:     maxLines,
		viewWidth:    viewCols,
		viewHeight:   viewLines,
		topMargin:    0,
		bottomMargin: uint(viewLines - 1),
		Charsets:     []*map[rune]rune{nil, nil},
		LineFeedMode: true,
	}
	b.TabReset()
	return b
}

func (terminalState *TerminalState) DefaultCell(applyEffects bool) Cell {
	attr := terminalState.CursorAttr
	if !applyEffects {
		attr.Blink = false
		attr.Bold = false
		attr.Dim = false
		attr.Inverse = false
		attr.Underline = false
		attr.Dim = false
	}
	return Cell{attr: attr}
}

func (terminalState *TerminalState) SetVerticalMargins(top uint, bottom uint) {
	terminalState.topMargin = top
	terminalState.bottomMargin = bottom
}

// ResetVerticalMargins resets margins to extreme positions
func (terminalState *TerminalState) ResetVerticalMargins() {
	terminalState.SetVerticalMargins(0, uint(terminalState.viewHeight-1))
}

func (terminalState *TerminalState) IsNewLineMode() bool {
	return terminalState.LineFeedMode == false
}

func (terminalState *TerminalState) TabZonk() {
	terminalState.tabStops = make(map[uint16]struct{})
}

func (terminalState *TerminalState) TabSet(index uint16) {
	terminalState.tabStops[index] = struct{}{}
}

func (terminalState *TerminalState) TabClear(index uint16) {
	delete(terminalState.tabStops, index)
}

func (terminalState *TerminalState) getTabIndexFromCursor() uint16 {
	index := terminalState.cursorX
	if index == terminalState.viewWidth {
		index = 0
	}
	return index
}

func (terminalState *TerminalState) IsTabSetAtCursor() bool {
	index := terminalState.getTabIndexFromCursor()
	_, ok := terminalState.tabStops[index]
	return ok
}

func (terminalState *TerminalState) TabClearAtCursor() {
	terminalState.TabClear(terminalState.getTabIndexFromCursor())
}

func (terminalState *TerminalState) TabSetAtCursor() {
	terminalState.TabSet(terminalState.getTabIndexFromCursor())
}

func (terminalState *TerminalState) TabReset() {
	terminalState.TabZonk()
	const MaxTabs uint16 = 1024
	const TabStep = 4
	var i uint16
	for i < MaxTabs {
		terminalState.TabSet(i)
		i += TabStep
	}
}

func (terminalState *TerminalState) ViewHeight() uint16 {
	return terminalState.viewHeight
}

func (terminalState *TerminalState) GetScrollOffset() uint {
	return terminalState.scrollLinesFromBottom
}

func (terminalState *TerminalState) SetScrollOffset(offset uint) {
	terminalState.scrollLinesFromBottom = offset
}
