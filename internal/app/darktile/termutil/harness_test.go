package termutil

import (
	"fmt"
	"image/color"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode/utf8"
)

// testWindow is a WindowManipulator that records what the terminal asks of it
// instead of touching a real window.
type testWindow struct {
	mu         sync.Mutex
	title      string
	titleStack []string
	fullscreen bool
	errs       []error
}

func (w *testWindow) State() WindowState { return StateNormal }
func (w *testWindow) Minimise()          {}
func (w *testWindow) Maximise()          {}
func (w *testWindow) Restore()           {}

func (w *testWindow) SetTitle(title string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.title = title
}

func (w *testWindow) GetTitle() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.title
}

func (w *testWindow) SaveTitleToStack() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.titleStack = append(w.titleStack, w.title)
}

func (w *testWindow) RestoreTitleFromStack() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.titleStack) == 0 {
		return
	}
	w.title = w.titleStack[len(w.titleStack)-1]
	w.titleStack = w.titleStack[:len(w.titleStack)-1]
}

func (w *testWindow) Position() (int, int)           { return 0, 0 }
func (w *testWindow) SizeInPixels() (int, int)       { return 800, 600 }
func (w *testWindow) CellSizeInPixels() (int, int)   { return 10, 20 }
func (w *testWindow) SizeInChars() (int, int)        { return 80, 30 }
func (w *testWindow) ResizeInPixels(int, int)        {}
func (w *testWindow) ResizeInChars(int, int)         {}
func (w *testWindow) ScreenSizeInPixels() (int, int) { return 1920, 1080 }
func (w *testWindow) ScreenSizeInChars() (int, int)  { return 192, 54 }
func (w *testWindow) Move(x, y int)                  {}
func (w *testWindow) IsFullscreen() bool             { return w.fullscreen }
func (w *testWindow) SetFullscreen(enabled bool)     { w.fullscreen = enabled }

func (w *testWindow) ReportError(err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.errs = append(w.errs, err)
}

// testTheme mirrors the shipped default theme so tests assert against the
// colours users actually see. It cannot be imported from the config package,
// which depends on termutil.
func testTheme() *Theme {
	rgb := func(v uint32) color.Color {
		return color.RGBA{uint8(v >> 16), uint8(v >> 8), uint8(v), 0xff}
	}
	return NewThemeFactory().
		WithColour(ColourBlack, rgb(0x1d1f21)).
		WithColour(ColourRed, rgb(0xcc6666)).
		WithColour(ColourGreen, rgb(0xb5bd68)).
		WithColour(ColourYellow, rgb(0xf0c674)).
		WithColour(ColourBlue, rgb(0x81a2be)).
		WithColour(ColourMagenta, rgb(0xb294bb)).
		WithColour(ColourCyan, rgb(0x8abeb7)).
		WithColour(ColourWhite, rgb(0xc5c8c6)).
		WithColour(ColourBrightBlack, rgb(0x666666)).
		WithColour(ColourBrightRed, rgb(0xd54e53)).
		WithColour(ColourBrightGreen, rgb(0xb9ca4a)).
		WithColour(ColourBrightYellow, rgb(0xe7c547)).
		WithColour(ColourBrightBlue, rgb(0x7aa6da)).
		WithColour(ColourBrightMagenta, rgb(0xc397d8)).
		WithColour(ColourBrightCyan, rgb(0x70c0b1)).
		WithColour(ColourBrightWhite, rgb(0xeaeaea)).
		WithColour(ColourBackground, rgb(0x000000)).
		WithColour(ColourForeground, rgb(0xc5c8c6)).
		WithColour(ColourSelectionBackground, rgb(0x33aa33)).
		WithColour(ColourSelectionForeground, rgb(0xffffff)).
		WithColour(ColourCursorForeground, rgb(0x1d1f21)).
		WithColour(ColourCursorBackground, rgb(0xc5c8c6)).
		Build()
}

// testTerm drives a Terminal without a shell, a pty or a GPU.
type testTerm struct {
	*Terminal
	t       *testing.T
	window  *testWindow
	replies *os.File
}

// newTestTerm builds a terminal whose replies (DSR, DA, window reports) land in
// a pipe rather than a pty, so they can be asserted on.
func newTestTerm(t *testing.T, cols, rows uint16) *testTerm {
	t.Helper()

	replyR, replyW, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create reply pipe: %v", err)
	}
	t.Cleanup(func() {
		_ = replyR.Close()
		_ = replyW.Close()
	})

	window := &testWindow{}
	term := New(WithTheme(testTheme()))
	term.SetWindowManipulator(window)
	term.pty = replyW

	for _, buf := range term.buffers {
		buf.resizeView(cols, rows)
	}

	return &testTerm{Terminal: term, t: t, window: window, replies: replyR}
}

// feed pushes input through the real parser and returns once the terminal has
// settled. The parser reads from processChan while decoding a sequence, so the
// queue is filled first and then drained on a single goroutine — that keeps
// ordering identical to production while staying synchronous for assertions.
func (tt *testTerm) feed(input string) {
	tt.t.Helper()

	for _, r := range input {
		select {
		// Terminal.Write derives Width from bufio.ReadRune's byte count, so
		// match that rather than display width, or the harness diverges from
		// production.
		case tt.processChan <- MeasuredRune{Rune: r, Width: utf8.RuneLen(r)}:
		default:
			tt.t.Fatalf("process channel full after %d runes", len(tt.processChan))
		}
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for len(tt.processChan) > 0 {
			tt.processSequence(<-tt.processChan)
		}
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		// An unterminated sequence leaves handleANSI blocked on the channel.
		tt.t.Fatalf("parser stalled on input %q — sequence is incomplete or the handler never returns", input)
	}
}

// screen renders the visible grid as text, with unwritten cells as spaces and
// trailing blanks stripped.
func (tt *testTerm) screen() string {
	lines := tt.GetActiveBuffer().GetVisibleLines()
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		out = append(out, strings.TrimRight(strings.ReplaceAll(line.String(), "\x00", " "), " "))
	}
	for len(out) > 0 && out[len(out)-1] == "" {
		out = out[:len(out)-1]
	}
	return strings.Join(out, "\n")
}

// reply returns everything the terminal has written back to the host, or "" if
// it wrote nothing.
func (tt *testTerm) reply() string {
	tt.t.Helper()
	if err := tt.replies.SetReadDeadline(time.Now().Add(100 * time.Millisecond)); err != nil {
		tt.t.Fatalf("failed to set read deadline: %v", err)
	}
	buf := make([]byte, 4096)
	n, err := tt.replies.Read(buf)
	if err != nil && n == 0 {
		return ""
	}
	return string(buf[:n])
}

func (tt *testTerm) cursor() (col, row uint16) {
	buf := tt.GetActiveBuffer()
	return buf.CursorColumn(), buf.CursorLine()
}

func (tt *testTerm) cellAt(col, row uint16) *Cell {
	tt.t.Helper()
	cell := tt.GetActiveBuffer().GetCell(col, row)
	if cell == nil {
		tt.t.Fatalf("no cell at %d,%d", col, row)
	}
	return cell
}

func hexColour(c color.Color) string {
	if c == nil {
		return "nil"
	}
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
}
