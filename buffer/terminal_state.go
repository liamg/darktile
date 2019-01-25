package buffer

type TerminalState struct {
	scrollLinesFromBottom uint
	cursorX               uint16
	cursorY               uint16
	cursorAttr            CellAttributes
	defaultCell           Cell
	viewHeight            uint16
	viewWidth             uint16
	topMargin             uint // see DECSTBM docs - this is for scrollable regions
	bottomMargin          uint // see DECSTBM docs - this is for scrollable regions
	ReplaceMode           bool // overwrite character at cursor or insert new
	OriginMode            bool // see DECOM docs - whether cursor is positioned within the margins or not
	LineFeedMode          bool
	AutoWrap              bool
	maxLines              uint64
}

// NewTerminalMode creates a new terminal state
func NewTerminalState(viewCols uint16, viewLines uint16, attr CellAttributes, maxLines uint64) *TerminalState {
	b := &TerminalState{
		cursorX:      0,
		cursorY:      0,
		cursorAttr:   attr,
		AutoWrap:     true,
		defaultCell:  Cell{attr: attr},
		maxLines:     maxLines,
		viewWidth:    viewCols,
		viewHeight:   viewLines,
		topMargin:    0,
		bottomMargin: uint(viewLines - 1),
	}
	return b
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
