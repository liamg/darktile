package termutil

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var updateGolden = flag.Bool("update", false, "rewrite golden files from current behaviour")

// goldenSessions lists the streams replayed through the parser. The view size
// is part of the fixture: the same bytes produce a different grid at a
// different width.
//
// Synthetic cases carry their input inline so it stays readable. Cases with a
// file name replay a real session captured from a program via TestRecordSession.
var goldenSessions = []struct {
	name  string
	cols  uint16
	rows  uint16
	input string
	file  string
}{
	{name: "256-colour", cols: 80, rows: 24, file: "256-colour.vt"},
	{name: "true-colour", cols: 80, rows: 24, file: "true-colour.vt"},
	{
		name: "sgr-attributes", cols: 40, rows: 10,
		input: "plain\r\n" +
			"\x1b[1mbold\x1b[0m \x1b[2mdim\x1b[0m \x1b[3mitalic\x1b[0m\r\n" +
			"\x1b[4munderline\x1b[0m \x1b[9mstrike\x1b[0m\r\n" +
			"\x1b[7minverse\x1b[0m \x1b[5mblink\x1b[0m\r\n" +
			"\x1b[31mred\x1b[32mgreen\x1b[34mblue\x1b[39m\r\n" +
			"\x1b[41mon-red\x1b[49m \x1b[91mbright\x1b[39m\r\n" +
			"\x1b[38;5;196m8bit\x1b[38;2;18;52;86m24bit\x1b[0m\r\n" +
			"\x1b[1;4;31mcombined\x1b[0m",
	},
	{
		name: "box-drawing", cols: 40, rows: 8,
		input: "\x1b(0" +
			"lqqqqqqqqk\r\n" +
			"x\x1b(Bplain\x1b(0    x\r\n" +
			"mqqqqqqqqj" +
			"\x1b(B",
	},
	{
		name: "scroll-region", cols: 20, rows: 10,
		input: "\x1b[2Jheader\r\n" +
			"\x1b[3;8r\x1b[3;1H" +
			"a\r\nb\r\nc\r\nd\r\ne\r\nf\r\ng\r\nh" +
			"\x1b[r\x1b[10;1Hfooter",
	},
	{
		name: "alt-screen", cols: 20, rows: 6,
		input: "main line one\r\nmain line two" +
			"\x1b[?1049h\x1b]0;alt title\x07" +
			"\x1b[2J\x1b[1;1Hfullscreen app" +
			"\x1b[?1049l",
	},
	{
		name: "wrapping", cols: 12, rows: 8,
		input: "the quick brown fox jumps\r\n" +
			"exactly12chr\r\n" +
			"short",
	},
	{
		name: "erase-ops", cols: 20, rows: 6,
		input: "abcdefghij\r\nklmnopqrst\r\nuvwxyz\r\n" +
			"\x1b[1;5H\x1b[0K" +
			"\x1b[2;5H\x1b[1K" +
			"\x1b[3;3H\x1b[2P" +
			"\x1b[3;1H\x1b[3@",
	},
}

// TestGoldenSessions replays each recorded stream and compares a full snapshot
// of the resulting terminal state against a checked-in golden file. This is the
// regression net for dependency bumps and refactors: any change in parsing,
// buffer handling or attribute tracking shows up as a diff.
func TestGoldenSessions(t *testing.T) {
	for _, session := range goldenSessions {
		t.Run(session.name, func(t *testing.T) {
			goldenPath := filepath.Join("testdata", "sessions", session.name+".golden")

			input := session.input
			if session.file != "" {
				raw, err := os.ReadFile(filepath.Join("testdata", "sessions", session.file))
				require.NoError(t, err, "missing capture — record it with TestRecordSession")
				input = string(raw)
			}

			term := newTestTerm(t, session.cols, session.rows)
			term.feed(input)
			got := term.snapshot()

			if *updateGolden {
				require.NoError(t, os.WriteFile(goldenPath, []byte(got), 0o644))
				return
			}

			want, err := os.ReadFile(goldenPath)
			require.NoError(t, err, "missing golden — regenerate with 'go test -run TestGoldenSessions -update'")
			require.Equal(t, string(want), got)
		})
	}
}

// snapshot renders the full terminal state in a diff-friendly form: the grid as
// text, followed by run-length encoded attributes for any cell that is not at
// its defaults.
func (tt *testTerm) snapshot() string {
	buf := tt.GetActiveBuffer()

	var out strings.Builder
	fmt.Fprintf(&out, "size: %dx%d\n", buf.ViewWidth(), buf.ViewHeight())
	fmt.Fprintf(&out, "buffer: %s\n", tt.activeBufferName())
	fmt.Fprintf(&out, "title: %q\n", tt.GetTitle())
	fmt.Fprintf(&out, "cursor: col=%d row=%d visible=%t\n", buf.CursorColumn(), buf.CursorLine(), buf.IsCursorVisible())
	if buf.HasScrollableRegion() {
		fmt.Fprintf(&out, "scroll-region: %d-%d\n", buf.TopMargin(), buf.BottomMargin())
	} else {
		out.WriteString("scroll-region: none\n")
	}

	out.WriteString("\n--- text ---\n")
	for row := uint16(0); row < buf.ViewHeight(); row++ {
		fmt.Fprintf(&out, "%3d|%s|\n", row, tt.rowText(row))
	}

	out.WriteString("\n--- attributes ---\n")
	attrs := tt.attributeRuns()
	if attrs == "" {
		out.WriteString("(all cells at defaults)\n")
	} else {
		out.WriteString(attrs)
	}

	return out.String()
}

func (tt *testTerm) activeBufferName() string {
	for i, buf := range tt.buffers {
		if buf == tt.activeBuffer {
			switch uint8(i) {
			case MainBuffer:
				return "main"
			case AltBuffer:
				return "alt"
			case InternalBuffer:
				return "internal"
			}
		}
	}
	return "unknown"
}

func (tt *testTerm) rowText(row uint16) string {
	buf := tt.GetActiveBuffer()
	var line strings.Builder
	for col := uint16(0); col < buf.ViewWidth(); col++ {
		cell := buf.GetCell(col, row)
		if cell == nil || cell.Rune().Rune == 0 {
			line.WriteRune(' ')
			continue
		}
		line.WriteRune(cell.Rune().Rune)
	}
	return strings.TrimRight(line.String(), " ")
}

// attributeRuns collapses each row into runs of identical styling, emitting
// only the runs that differ from a freshly written default cell.
func (tt *testTerm) attributeRuns() string {
	buf := tt.GetActiveBuffer()

	var out strings.Builder
	for row := uint16(0); row < buf.ViewHeight(); row++ {
		var runs []string
		runStart := uint16(0)
		runSig := ""

		flush := func(end uint16) {
			switch runSig {
			case "", defaultCellSignature, unsetCellSignature:
				return
			}
			runs = append(runs, fmt.Sprintf("[%d-%d]%s", runStart, end-1, runSig))
		}

		for col := uint16(0); col < buf.ViewWidth(); col++ {
			sig := cellSignature(buf.GetCell(col, row))
			if sig != runSig {
				flush(col)
				runStart = col
				runSig = sig
			}
		}
		flush(buf.ViewWidth())

		if len(runs) > 0 {
			fmt.Fprintf(&out, "row %d: %s\n", row, strings.Join(runs, " "))
		}
	}
	return out.String()
}

const unsetCellSignature = "unset"

// Cells at these colours carry no styling worth recording, so they are left out
// of the snapshot. Derived from the theme so the two cannot drift apart.
var defaultCellSignature = fmt.Sprintf(
	"fg=%s bg=%s",
	hexColour(testTheme().DefaultForeground()),
	hexColour(testTheme().DefaultBackground()),
)

func cellSignature(cell *Cell) string {
	if cell == nil {
		return unsetCellSignature
	}

	sig := fmt.Sprintf("fg=%s bg=%s", hexColour(cell.Fg()), hexColour(cell.Bg()))
	for _, flag := range []struct {
		name string
		on   bool
	}{
		{"bold", cell.Bold()},
		{"dim", cell.Dim()},
		{"italic", cell.Italic()},
		{"underline", cell.Underline()},
		{"strikethrough", cell.Strikethrough()},
		{"blink", cell.attr.blink},
		{"inverse", cell.attr.inverse},
		{"hidden", cell.attr.hidden},
	} {
		if flag.on {
			sig += "+" + flag.name
		}
	}
	return sig
}
