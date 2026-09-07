package termutil

import (
	"os"
	"testing"
	"time"
	"unicode/utf8"
)

// Escape sequences arrive from whatever the child process decides to print, so
// the parser is an untrusted-input boundary. This target asserts the two
// properties that matter there: it never panics, and it never grows the line
// buffer past the configured scrollback.
//
//	go test ./internal/app/darktile/termutil -run FuzzParser -fuzz FuzzParser
func FuzzParser(f *testing.F) {
	seeds := []string{
		"hello world",
		"\x1b[31mred\x1b[0m",
		"\x1b[1;1H\x1b[2J",
		"\x1b[3;8r\x1b[5;1Habc",
		"\x1b]0;title\x07",
		"\x1b]4;0;rgb:f0f0/f0f0/f0f0\x1b\\",
		"\x1b(0qqqq\x1b(B",
		"\x1b[?1049h\x1b[?1049l",
		"\x1b[99M\x1b[99L\x1b[99P\x1b[99@",
		"\x1b[1;;1;1;1;3*y",
		"\x1bPq#0;2;0;0;0#0~~~\x1b\\",
		"\x1b_apc\x1b\\\x1b^pm\x1b\\",
		"héllo → ✓ 😀",
		"\x1b[38;2;18;52;86m\x1b[48;5;196m",
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	// Replies would otherwise fill a pipe buffer and block the parser mid-run.
	discard, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		f.Fatalf("failed to open %s: %v", os.DevNull, err)
	}
	f.Cleanup(func() { _ = discard.Close() })

	theme := testTheme()

	f.Fuzz(func(t *testing.T, input string) {
		term := New(WithTheme(theme))
		term.SetWindowManipulator(&testWindow{})
		term.pty = discard
		for _, buf := range term.buffers {
			buf.resizeView(80, 24)
		}

		// A sequence cut off by the end of the input leaves the parser waiting
		// for the rest of it, so the input is followed by bytes that terminate
		// whichever kind of sequence is open. Without this every truncated
		// input would look like a hang.
		for _, r := range input + fuzzSequenceTerminators {
			select {
			case term.processChan <- MeasuredRune{Rune: r, Width: utf8.RuneLen(r)}:
			default:
				return // input too large to queue; nothing to learn from it
			}
		}

		done := make(chan struct{})
		go func() {
			defer close(done)
			for len(term.processChan) > 0 {
				term.processSequence(<-term.processChan)
			}
		}()

		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Fatalf("parser did not finish; input %q holds it open", input)
		}

		for _, buf := range term.buffers {
			if height := uint64(buf.Height()); height > buf.GetMaxLines() {
				t.Fatalf("buffer grew to %d lines, past the %d line limit", height, buf.GetMaxLines())
			}
		}
	})
}

// BEL closes an OSC, ESC \ closes a string sequence, and the backslash doubles
// as a CSI final byte. 'A' feeds any handler still waiting on a single rune.
const fuzzSequenceTerminators = "\x07\x1b\\A\x07\x1b\\A"
