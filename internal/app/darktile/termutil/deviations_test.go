package termutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The tests below pin down behaviour that is known to be wrong. They assert
// what the terminal does today, not what it should do, so that a refactor or a
// dependency bump cannot change it silently. When one of these is fixed, the
// test fails and should be moved into parser_test.go with the correct
// expectation.

// SGR 0 resets the whole attribute struct, which leaves the foreground and
// background as nil rather than the theme defaults that SGR 39/49 apply.
// Nothing breaks today only because the renderer substitutes a default when it
// sees a nil colour (gui/render/row.go). A renderer without that guard would
// mis-draw every cell written after a reset, which is almost all of them.
func TestDeviationSGRResetClearsColoursToNil(t *testing.T) {
	term := newTestTerm(t, 20, 6)
	term.feed("\x1b[0mx")

	cell := term.cellAt(0, 0)
	assert.Nil(t, cell.Fg(), "SGR 0 currently nils the foreground")
	assert.Nil(t, cell.Bg(), "SGR 0 currently nils the background")

	// A cell written with no SGR at all does get the theme defaults, which is
	// what makes the reset path inconsistent.
	untouched := newTestTerm(t, 20, 6)
	untouched.feed("x")
	assert.Equal(t, "#c5c8c6", hexColour(untouched.cellAt(0, 0).Fg()))
	assert.Equal(t, "#000000", hexColour(untouched.cellAt(0, 0).Bg()))
}

// OSC strings are terminated by BEL or by ST (the two-byte ESC \). The handler
// only looks for the bare backslash, so the ESC that precedes it has already
// been appended to the parameter by the time the loop breaks.
func TestDeviationOSCStringTerminatorLeavesEscapeInTitle(t *testing.T) {
	term := newTestTerm(t, 20, 6)
	term.feed("\x1b]0;st title\x1b\\")

	assert.Equal(t, "st title\x1b", term.GetTitle(), "ST-terminated OSC keeps a trailing ESC")

	// BEL-terminated titles are unaffected.
	bel := newTestTerm(t, 20, 6)
	bel.feed("\x1b]0;bel title\x07")
	assert.Equal(t, "bel title", bel.GetTitle())
}

// A consequence of the same bug: a literal backslash inside an OSC string is
// treated as a terminator, so titles containing one are truncated.
func TestDeviationOSCTitleTruncatedAtBackslash(t *testing.T) {
	term := newTestTerm(t, 20, 6)
	term.feed("\x1b]0;C:\\Users\\me\x07")

	assert.Equal(t, "C:", term.GetTitle(), "backslash currently terminates the OSC string")
}

// Reverse index at the top margin scrolls the visible area down. That works
// once the buffer holds more than one line, but with a single line the scroll
// range is a single row, so the content is blanked instead of pushed down.
func TestDeviationReverseIndexLosesContentOnSingleLineBuffer(t *testing.T) {
	single := newTestTerm(t, 10, 4)
	single.feed("one")
	single.feed("\x1b[1;1H\x1bM")
	assert.Equal(t, "", single.screen(), "single-line buffer loses its content")

	// With two or more lines the same sequence behaves correctly.
	two := newTestTerm(t, 10, 4)
	two.feed("one\r\ntwo")
	two.feed("\x1b[1;1H\x1bM")
	require.Equal(t, "\none", two.screen())
}

// An escape sequence that is cut off mid-way leaves the parser blocked waiting
// for the rest of it. That is tolerable for a streaming parser fed by a pty,
// but it means parser state is not a pure function of its input, which the
// fuzzing and golden layers have to work around.
func TestDeviationTruncatedSequenceBlocksParser(t *testing.T) {
	term := newTestTerm(t, 20, 6)

	for _, r := range "ok\x1b[" {
		term.processChan <- MeasuredRune{Rune: r, Width: 1}
	}

	settled := make(chan struct{})
	go func() {
		defer close(settled)
		for len(term.processChan) > 0 {
			term.processSequence(<-term.processChan)
		}
	}()

	select {
	case <-settled:
		t.Fatal("parser no longer blocks on a truncated sequence — good news, update this test")
	case <-time.After(250 * time.Millisecond):
	}
}

// Insert Line outside a scrollable region grows the line buffer instead of
// shifting content down and discarding what falls off the bottom. The cursor's
// raw line then sits above the view, and converting it back underflows uint16,
// so the terminal reports a nonsense position to the program that asked.
//
// xterm keeps the buffer height fixed here. Fixing this means changing how IL
// interacts with scrollback, so it is pinned rather than corrected.
func TestDeviationInsertLineGrowsBufferAndUnderflowsCursorReport(t *testing.T) {
	term := newTestTerm(t, 20, 6)
	term.feed("one\r\ntwo\r\nthree")

	before := term.GetActiveBuffer().Height()
	term.feed("\x1b[1;1H\x1b[99L")
	after := term.GetActiveBuffer().Height()

	assert.Greater(t, after, before, "buffer grew by the inserted lines")

	term.feed("\x1b[6n")
	assert.Equal(t, "\x1b[65534;1R", term.reply(), "cursor row underflowed")
}

// Every rune advances the cursor by exactly one column, whatever its display
// width. Wide characters (CJK, most emoji) should take two cells and combining
// marks none, so CJK text overlaps the cell to its right and an accent pushes
// the rest of the line along.
//
// MeasuredRune.Width carries the rune's UTF-8 byte length, not its display
// width, and the buffer ignores it. ucs-detect stops at its entry probe because
// of this (scripts/ucs-detect.sh), so this test is the regression guard until
// widths are implemented against wcwidth.
func TestDeviationEveryRuneOccupiesOneColumn(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  uint16
	}{
		{"ascii", "a", 1},
		{"latin-1 precomposed", "é", 1},
		{"CJK ideograph", "一", 2},
		{"emoji", "⌚", 2},
		{"combining mark", "é", 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			term := newTestTerm(t, 20, 6)
			term.feed(test.input)

			col, _ := term.cursor()
			assert.Equal(t, uint16(len([]rune(test.input))), col,
				"cursor advanced one column per rune")
			if col != test.want {
				t.Logf("correct width would put the cursor at column %d", test.want)
			}
		})
	}
}
