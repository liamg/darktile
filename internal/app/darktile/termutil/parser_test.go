package termutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParserScreenOutput drives whole escape sequences through the parser and
// asserts on the resulting grid. These cover the CSI/OSC layer, which the
// buffer-level tests bypass entirely.
func TestParserScreenOutput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "plain text",
			input: "hello",
			want:  "hello",
		},
		{
			name:  "carriage return and line feed",
			input: "one\r\ntwo",
			want:  "one\ntwo",
		},
		{
			name:  "cursor position CUP",
			input: "\x1b[3;5Hmoved",
			want:  "\n\n    moved",
		},
		{
			name:  "cursor position defaults to home",
			input: "xxxx\x1b[Hy",
			want:  "yxxx",
		},
		{
			name:  "cursor forward CUF",
			input: "a\x1b[3Cb",
			want:  "a   b",
		},
		{
			name:  "cursor backward CUB overwrites",
			input: "abcdef\x1b[3DX",
			want:  "abcXef",
		},
		{
			name:  "cursor up and down",
			input: "l1\r\nl2\r\nl3\x1b[2Ax\x1b[2Bz",
			want:  "l1x\nl2\nl3 z",
		},
		{
			name:  "erase in line to end EL0",
			input: "abcdef\x1b[3G\x1b[0K",
			want:  "ab",
		},
		{
			name:  "erase in line to start EL1",
			input: "abcdef\x1b[3G\x1b[1K",
			want:  "   def",
		},
		{
			name:  "erase whole line EL2",
			input: "abcdef\x1b[2K",
			want:  "",
		},
		{
			name:  "erase in display ED2",
			input: "one\r\ntwo\r\nthree\x1b[2J",
			want:  "",
		},
		{
			name:  "erase characters ECH",
			input: "abcdef\x1b[1;2H\x1b[3X",
			want:  "a   ef",
		},
		{
			name:  "delete characters DCH",
			input: "abcdef\x1b[1;2H\x1b[2P",
			want:  "adef",
		},
		{
			name:  "insert blank characters ICH",
			input: "abcdef\x1b[1;2H\x1b[2@",
			want:  "a  bcdef",
		},
		{
			name:  "insert line IL",
			input: "one\r\ntwo\x1b[1;1H\x1b[1L",
			want:  "\none\ntwo",
		},
		{
			name:  "delete line DL",
			input: "one\r\ntwo\r\nthree\x1b[1;1H\x1b[1M",
			want:  "two\nthree",
		},
		{
			name:  "save and restore cursor",
			input: "\x1b[2;3H\x1b7\x1b[5;5Hx\x1b8y",
			want:  "\n  y\n\n\n    x",
		},
		{
			name:  "backspace moves without erasing",
			input: "abc\bX",
			want:  "abX",
		},
		{
			name:  "tab advances to next stop",
			input: "a\tb",
			want:  "a       b",
		},
		{
			name:  "reverse index scrolls content down at top margin",
			input: "l1\r\nl2\r\nl3\r\nl4\r\nl5\r\nl6\x1b[1;1H\x1bMzero",
			want:  "zero\nl1\nl2\nl3\nl4\nl5",
		},
		{
			name:  "cursor character absolute CHA",
			input: "abcdef\x1b[3GX",
			want:  "abXdef",
		},
		{
			name:  "line position absolute VPA",
			input: "\x1b[3dx",
			want:  "\n\nx",
		},
		{
			name:  "next line CNL",
			input: "one\x1b[1Etwo",
			want:  "one\ntwo",
		},
		{
			name:  "unknown sequence is ignored",
			input: "a\x1b[999zb",
			want:  "ab",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			term := newTestTerm(t, 20, 6)
			term.feed(test.input)
			assert.Equal(t, test.want, term.screen())
		})
	}
}

func TestParserCursorPosition(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantCol  uint16
		wantLine uint16
	}{
		{name: "home", input: "\x1b[H", wantCol: 0, wantLine: 0},
		{name: "explicit position", input: "\x1b[4;7H", wantCol: 6, wantLine: 3},
		{name: "position is clamped to view", input: "\x1b[99;99H", wantCol: 19, wantLine: 5},
		{name: "cursor up stops at top", input: "\x1b[10A", wantCol: 0, wantLine: 0},
		{name: "cursor back stops at column zero", input: "\x1b[10D", wantCol: 0, wantLine: 0},
		{name: "writing advances the cursor", input: "abc", wantCol: 3, wantLine: 0},
		{name: "carriage return resets column", input: "abc\r", wantCol: 0, wantLine: 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			term := newTestTerm(t, 20, 6)
			term.feed(test.input)
			col, line := term.cursor()
			assert.Equal(t, test.wantCol, col, "column")
			assert.Equal(t, test.wantLine, line, "line")
		})
	}
}

func TestParserSGRAttributes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(t *testing.T, cell *Cell)
	}{
		{
			name:  "bold",
			input: "\x1b[1mx",
			check: func(t *testing.T, cell *Cell) { assert.True(t, cell.Bold()) },
		},
		{
			name:  "italic",
			input: "\x1b[3mx",
			check: func(t *testing.T, cell *Cell) { assert.True(t, cell.Italic()) },
		},
		{
			name:  "underline",
			input: "\x1b[4mx",
			check: func(t *testing.T, cell *Cell) { assert.True(t, cell.Underline()) },
		},
		{
			name:  "strikethrough",
			input: "\x1b[9mx",
			check: func(t *testing.T, cell *Cell) { assert.True(t, cell.Strikethrough()) },
		},
		{
			name:  "dim",
			input: "\x1b[2mx",
			check: func(t *testing.T, cell *Cell) { assert.True(t, cell.Dim()) },
		},
		{
			name:  "reset clears attributes",
			input: "\x1b[1;3;4m\x1b[0mx",
			check: func(t *testing.T, cell *Cell) {
				assert.False(t, cell.Bold())
				assert.False(t, cell.Italic())
				assert.False(t, cell.Underline())
			},
		},
		{
			name:  "bold off via 22",
			input: "\x1b[1m\x1b[22mx",
			check: func(t *testing.T, cell *Cell) { assert.False(t, cell.Bold()) },
		},
		{
			name:  "4-bit foreground",
			input: "\x1b[31mx",
			check: func(t *testing.T, cell *Cell) { assert.Equal(t, "#cc6666", hexColour(cell.Fg())) },
		},
		{
			name:  "4-bit background",
			input: "\x1b[42mx",
			check: func(t *testing.T, cell *Cell) { assert.Equal(t, "#b5bd68", hexColour(cell.Bg())) },
		},
		{
			name:  "bright foreground",
			input: "\x1b[91mx",
			check: func(t *testing.T, cell *Cell) { assert.Equal(t, "#d54e53", hexColour(cell.Fg())) },
		},
		{
			name:  "24-bit foreground",
			input: "\x1b[38;2;18;52;86mx",
			check: func(t *testing.T, cell *Cell) { assert.Equal(t, "#123456", hexColour(cell.Fg())) },
		},
		{
			name:  "8-bit foreground from the 216-colour cube",
			input: "\x1b[38;5;196mx",
			check: func(t *testing.T, cell *Cell) { assert.Equal(t, "#ff0000", hexColour(cell.Fg())) },
		},
		{
			name:  "inverse swaps foreground and background",
			input: "\x1b[31;42m\x1b[7mx",
			check: func(t *testing.T, cell *Cell) {
				assert.Equal(t, "#b5bd68", hexColour(cell.Fg()))
				assert.Equal(t, "#cc6666", hexColour(cell.Bg()))
			},
		},
		{
			name:  "default foreground",
			input: "\x1b[31m\x1b[39mx",
			check: func(t *testing.T, cell *Cell) { assert.Equal(t, "#c5c8c6", hexColour(cell.Fg())) },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			term := newTestTerm(t, 20, 6)
			term.feed(test.input)
			test.check(t, term.cellAt(0, 0))
		})
	}
}

func TestParserSGRAppliesOnlyToFollowingCells(t *testing.T) {
	term := newTestTerm(t, 20, 6)
	term.feed("a\x1b[31mb\x1b[39mc")

	assert.Equal(t, "#c5c8c6", hexColour(term.cellAt(0, 0).Fg()), "cell before SGR")
	assert.Equal(t, "#cc6666", hexColour(term.cellAt(1, 0).Fg()), "cell inside SGR")
	assert.Equal(t, "#c5c8c6", hexColour(term.cellAt(2, 0).Fg()), "cell after default-fg")
}

func TestParserOSCTitle(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "OSC 0 sets title, BEL terminated", input: "\x1b]0;my title\x07", want: "my title"},
		{name: "OSC 2 sets title, BEL terminated", input: "\x1b]2;other title\x07", want: "other title"},
		{name: "empty title", input: "\x1b]0;\x07", want: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			term := newTestTerm(t, 20, 6)
			term.feed(test.input)
			assert.Equal(t, test.want, term.GetTitle())
		})
	}
}

func TestParserAltBuffer(t *testing.T) {
	term := newTestTerm(t, 20, 6)

	term.feed("main content")
	require.Equal(t, "main content", term.screen())

	term.feed("\x1b[?1049h")
	assert.Equal(t, "", term.screen(), "alt buffer starts empty")

	term.feed("alt content")
	assert.Equal(t, "alt content", term.screen())

	term.feed("\x1b[?1049l")
	assert.Equal(t, "main content", term.screen(), "main buffer survives the round trip")
}

func TestParserCursorVisibility(t *testing.T) {
	term := newTestTerm(t, 20, 6)
	require.True(t, term.GetActiveBuffer().IsCursorVisible())

	term.feed("\x1b[?25l")
	assert.False(t, term.GetActiveBuffer().IsCursorVisible())

	term.feed("\x1b[?25h")
	assert.True(t, term.GetActiveBuffer().IsCursorVisible())
}

func TestParserScrollingRegion(t *testing.T) {
	term := newTestTerm(t, 20, 6)
	term.feed("\x1b[2;4r")

	buf := term.GetActiveBuffer()
	assert.True(t, buf.HasScrollableRegion())
	assert.Equal(t, uint(1), buf.TopMargin())
	assert.Equal(t, uint(3), buf.BottomMargin())

	term.feed("\x1b[r")
	assert.False(t, term.GetActiveBuffer().HasScrollableRegion())
}

func TestParserDECALNFillsScreen(t *testing.T) {
	term := newTestTerm(t, 4, 3)
	term.feed("\x1b#8")

	assert.Equal(t, "EEEE\nEEEE\nEEEE", term.screen())

	col, line := term.cursor()
	assert.Equal(t, uint16(0), col)
	assert.Equal(t, uint16(0), line)
}

// TestParserReplies covers the sequences where the terminal answers the host.
// These are what esctest-style compliance suites rely on.
func TestParserReplies(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "device status report is OK",
			input: "\x1b[5n",
			want:  "\x1b[0n",
		},
		{
			name:  "cursor position report is 1-based",
			input: "\x1b[3;7H\x1b[6n",
			want:  "\x1b[3;7R",
		},
		{
			name:  "cursor position report at home",
			input: "\x1b[6n",
			want:  "\x1b[1;1R",
		},
		{
			name:  "primary device attributes",
			input: "\x1b[c",
			want:  "\x1b[?1;2c",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			term := newTestTerm(t, 20, 6)
			term.feed(test.input)
			assert.Equal(t, test.want, term.reply())
		})
	}
}

func TestParserCharsetSelection(t *testing.T) {
	term := newTestTerm(t, 20, 6)

	// Select the DEC special graphics set into G0 and shift it in: 'q' becomes
	// a horizontal line, which is how box drawing reaches the screen.
	term.feed("\x1b(0qqq")
	assert.Equal(t, "───", term.screen())

	term.feed("\r\x1b(Bqqq")
	assert.Equal(t, "qqq", term.screen(), "ASCII charset restored")
}

func TestParserLineWrap(t *testing.T) {
	term := newTestTerm(t, 5, 4)
	term.feed("abcdefgh")
	assert.Equal(t, "abcde\nfgh", term.screen())
}

func TestParserUnicode(t *testing.T) {
	term := newTestTerm(t, 20, 4)
	term.feed("héllo → ✓")
	assert.Equal(t, "héllo → ✓", term.screen())
}

// A CSI sequence carrying an intermediate byte that no handler claims used to
// deadlock the parser goroutine: handleANSI holds the terminal mutex and the
// unknown-sequence fallback tried to take it again. Any program probing for
// DECRQCRA, DECRQM or DECSCA would freeze the terminal for good.
func TestParserUnknownCSIWithIntermediateIsIgnored(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"DECRQCRA", "\x1b[1;;1;1;1;3*y"},
		{"DECRQM", "\x1b[?1049$p"},
		{"DECSCA", "\x1b[1\"q"},
		{"unassigned intermediate", "\x1b[5#z"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			term := newTestTerm(t, 20, 6)
			term.feed(test.input + "ok")
			assert.Equal(t, "ok", term.screen(), "sequence must not reach the screen")

			term.feed("\x1b[6n")
			assert.Equal(t, "\x1b[1;3R", term.reply(), "terminal must still respond")
		})
	}
}

// Line editing sequences take a repeat count straight from the program, so the
// count routinely exceeds the number of lines that actually exist. Deleting
// past the end used to slice out of range and take the terminal down.
func TestParserLineEditingCountsBeyondBuffer(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"delete lines", "\x1b[99M"},
		{"insert lines", "\x1b[99L"},
		{"delete characters", "\x1b[99P"},
		{"insert characters", "\x1b[99@"},
		{"erase characters", "\x1b[99X"},
		{"scroll up", "\x1b[99S"},
		{"scroll down", "\x1b[99T"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			term := newTestTerm(t, 20, 6)
			term.feed("one\r\ntwo\r\nthree")
			term.feed("\x1b[1;1H" + test.input)

			term.feed("\x1b[6n")
			assert.Regexp(t, `^\x1b\[\d+;\d+R$`, term.reply(), "terminal must survive the sequence")
		})
	}
}
