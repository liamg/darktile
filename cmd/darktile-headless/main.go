// Command darktile-headless runs the terminal core against a real pty with no
// window, renderer or GPU. It exists so external conformance suites such as
// esctest can drive the escape sequence handling directly:
//
//	darktile-headless --command 'python3 esctest.py --expected-terminal=xterm'
//
// Escape sequences written by the child reach the same parser the GUI uses, and
// the terminal's replies go back down the pty, which is what the suites read.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/liamg/darktile/internal/app/darktile/config"
	"github.com/liamg/darktile/internal/app/darktile/termutil"
)

func main() {
	var (
		cols    = flag.Uint("cols", 80, "terminal width in characters")
		rows    = flag.Uint("rows", 24, "terminal height in characters")
		command = flag.String("command", "", "command to run in the terminal (required)")
		shell   = flag.String("shell", "/bin/sh", "shell used to launch the command")
		logFile = flag.String("log-file", "", "write a parser trace here, or '-' for stdout")
	)
	flag.Parse()

	if *command == "" {
		fmt.Fprintln(os.Stderr, "--command is required")
		os.Exit(2)
	}

	theme, err := config.DefaultTheme(config.DefaultConfig())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build theme: %v\n", err)
		os.Exit(1)
	}

	window := &headlessWindow{
		cols: uint16(*cols),
		rows: uint16(*rows),
	}

	opts := []termutil.Option{
		termutil.WithTheme(theme),
		termutil.WithShell(*shell),
		termutil.WithInitialCommand(*command),
		termutil.WithWindowManipulator(window),
	}
	if *logFile != "" {
		opts = append(opts, termutil.WithLogFile(*logFile))
	}

	terminal := termutil.New(opts...)
	window.terminal = terminal

	// Nothing draws, but the terminal still signals when it wants a redraw and
	// drops the signal if nobody is listening.
	updateChan := make(chan struct{}, 1)
	go func() {
		for range updateChan {
		}
	}()

	if err := terminal.Run(updateChan, window.rows, window.cols); err != nil {
		fmt.Fprintf(os.Stderr, "terminal exited: %v\n", err)
		os.Exit(1)
	}
}

// headlessWindow answers the window queries a conformance suite makes without
// there being a window. Sizes are reported from the configured grid using a
// fixed notional cell, so pixel-based reports stay self-consistent.
type headlessWindow struct {
	terminal   *termutil.Terminal
	cols, rows uint16
	title      string
	titleStack []string
	fullscreen bool
}

const (
	cellWidthPx  = 10
	cellHeightPx = 20
	screenCols   = 200
	screenRows   = 60
)

func (w *headlessWindow) State() termutil.WindowState { return termutil.StateNormal }
func (w *headlessWindow) Minimise()                   {}
func (w *headlessWindow) Maximise()                   {}
func (w *headlessWindow) Restore()                    {}
func (w *headlessWindow) SetTitle(title string)       { w.title = title }
func (w *headlessWindow) GetTitle() string            { return w.title }
func (w *headlessWindow) Position() (int, int)        { return 0, 0 }
func (w *headlessWindow) IsFullscreen() bool          { return w.fullscreen }
func (w *headlessWindow) SetFullscreen(enabled bool)  { w.fullscreen = enabled }
func (w *headlessWindow) Move(x, y int)               {}

func (w *headlessWindow) SizeInPixels() (int, int) {
	return int(w.cols) * cellWidthPx, int(w.rows) * cellHeightPx
}

func (w *headlessWindow) CellSizeInPixels() (int, int) { return cellWidthPx, cellHeightPx }
func (w *headlessWindow) SizeInChars() (int, int)      { return int(w.cols), int(w.rows) }

func (w *headlessWindow) ScreenSizeInPixels() (int, int) {
	return screenCols * cellWidthPx, screenRows * cellHeightPx
}

func (w *headlessWindow) ScreenSizeInChars() (int, int) { return screenCols, screenRows }

func (w *headlessWindow) ResizeInPixels(width, height int) {
	w.ResizeInChars(width/cellWidthPx, height/cellHeightPx)
}

func (w *headlessWindow) ResizeInChars(width, height int) {
	if width <= 0 || height <= 0 {
		return
	}
	w.cols, w.rows = uint16(width), uint16(height)
	_ = w.terminal.SetSize(w.rows, w.cols)
}

func (w *headlessWindow) SaveTitleToStack() {
	w.titleStack = append(w.titleStack, w.title)
}

func (w *headlessWindow) RestoreTitleFromStack() {
	if len(w.titleStack) == 0 {
		return
	}
	w.title = w.titleStack[len(w.titleStack)-1]
	w.titleStack = w.titleStack[:len(w.titleStack)-1]
}

func (w *headlessWindow) ReportError(err error) {
	fmt.Fprintf(os.Stderr, "terminal error: %v\n", err)
}
