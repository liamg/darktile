package termutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

// Workloads follow vtebench: each one isolates a cost the parser pays, so a
// regression points at the part of the pipeline that caused it rather than at
// an aggregate number.
//
//	go test ./internal/app/darktile/termutil -run '^$' -bench Parser -benchmem
//
// Capture a baseline before touching dependencies or the buffer layer, then
// compare with benchstat:
//
//	go test ... -bench Parser -count=10 > /tmp/before.txt
//	benchstat /tmp/before.txt /tmp/after.txt
var parserWorkloads = []struct {
	name  string
	cols  uint16
	rows  uint16
	build func() string
	file  string
}{
	{
		name: "plain-ascii", cols: 80, rows: 24,
		// The cheapest possible path: no escape sequences at all.
		build: func() string {
			return strings.Repeat("the quick brown fox jumps over the lazy dog\r\n", 1200)
		},
	},
	{
		name: "scrolling", cols: 80, rows: 24,
		// Every line pushes one out of the view, so this measures the buffer
		// churn rather than the parser.
		build: func() string {
			var b strings.Builder
			for i := 0; i < 700; i++ {
				fmt.Fprintf(&b, "line %d: %s\r\n", i, strings.Repeat("x", 60))
			}
			return b.String()
		},
	},
	{
		name: "sgr-churn", cols: 80, rows: 24,
		// Colour changes per cell, the pattern syntax highlighting produces.
		build: func() string {
			var b strings.Builder
			for i := 0; i < 120; i++ {
				for col := 0; col < 40; col++ {
					fmt.Fprintf(&b, "\x1b[38;5;%dm#", col%256)
				}
				b.WriteString("\x1b[0m\r\n")
			}
			return b.String()
		},
	},
	{
		name: "cursor-motion", cols: 80, rows: 24,
		// Absolute addressing with no text, as a full-screen TUI redraw does.
		build: func() string {
			var b strings.Builder
			for i := 0; i < 2000; i++ {
				fmt.Fprintf(&b, "\x1b[%d;%dH*", (i%24)+1, (i%80)+1)
			}
			return b.String()
		},
	},
	{
		name: "unicode", cols: 80, rows: 24,
		build: func() string {
			return strings.Repeat("héllo wörld → ✓ 日本語テキスト\r\n", 900)
		},
	},
	{
		name: "alt-screen", cols: 80, rows: 24,
		// Alternate buffer redraws, which is where full-screen editors live.
		build: func() string {
			var b strings.Builder
			b.WriteString("\x1b[?1049h")
			for i := 0; i < 45; i++ {
				b.WriteString("\x1b[2J\x1b[H")
				for row := 0; row < 24; row++ {
					fmt.Fprintf(&b, "row %02d %s\r\n", row, strings.Repeat("=", 40))
				}
			}
			b.WriteString("\x1b[?1049l")
			return b.String()
		},
	},
	{name: "256-colour-capture", cols: 80, rows: 24, file: "256-colour.vt"},
	{name: "true-colour-capture", cols: 80, rows: 24, file: "true-colour.vt"},
}

func BenchmarkParser(b *testing.B) {
	discard, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		b.Fatalf("failed to open %s: %v", os.DevNull, err)
	}
	b.Cleanup(func() { _ = discard.Close() })

	theme := testTheme()

	for _, workload := range parserWorkloads {
		b.Run(workload.name, func(b *testing.B) {
			input := ""
			if workload.file != "" {
				raw, err := os.ReadFile(filepath.Join("testdata", "sessions", workload.file))
				if err != nil {
					b.Fatalf("missing capture: %v", err)
				}
				input = string(raw)
			} else {
				input = workload.build()
			}

			runes := []rune(input)
			// The parser reads its input from a buffered channel, so a single
			// batch has to fit in it.
			if len(runes) >= processChanCapacity {
				b.Fatalf("workload is %d runes, over the %d the channel holds", len(runes), processChanCapacity)
			}

			b.SetBytes(int64(len(input)))
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				// Each iteration starts from an empty terminal. Reusing one
				// would let scrollback pile up, making the result depend on
				// how many iterations the framework chose to run.
				b.StopTimer()
				term := New(WithTheme(theme))
				term.SetWindowManipulator(&testWindow{})
				term.pty = discard
				for _, buf := range term.buffers {
					buf.resizeView(workload.cols, workload.rows)
				}
				b.StartTimer()

				benchFeed(term, runes)
			}
		})
	}
}

// benchFeed mirrors how bytes reach the parser in production: they are queued
// and consumed by a separate goroutine, which matters because a handler blocks
// on the queue while it decodes.
func benchFeed(term *Terminal, runes []rune) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for len(term.processChan) > 0 {
			term.processSequence(<-term.processChan)
		}
	}()

	for _, r := range runes {
		term.processChan <- MeasuredRune{Rune: r, Width: utf8.RuneLen(r)}
	}
	<-done
}
