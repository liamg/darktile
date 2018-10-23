package terminal

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"syscall"
	"unsafe"

	"github.com/liamg/aminal/buffer"
	"github.com/liamg/aminal/config"
	"go.uber.org/zap"
)

const (
	MainBuffer uint8 = 0
	AltBuffer  uint8 = 1
)

type MouseMode uint

const (
	MouseModeNone MouseMode = iota
	MouseModeX10
	MouseModeVT200
	MouseModeVT200Highlight
	MouseModeButtonEvent
	MouseModeAnyEvent
)

type Terminal struct {
	buffers            []*buffer.Buffer
	activeBufferIndex  uint8
	lock               sync.Mutex
	pty                *os.File
	logger             *zap.SugaredLogger
	title              string
	size               Winsize
	config             *config.Config
	titleHandlers      []chan bool
	pauseChan          chan bool
	resumeChan         chan bool
	modes              Modes
	mouseMode          MouseMode
	bracketedPasteMode bool
	isDirty            bool
}

type Modes struct {
	ShowCursor            bool
	ApplicationCursorKeys bool
	BlinkingCursor        bool
}

type Winsize struct {
	Height uint16
	Width  uint16
	x      uint16 //ignored, but necessary for ioctl calls
	y      uint16 //ignored, but necessary for ioctl calls
}

type Position struct {
	Line int
	Col  int
}

func New(pty *os.File, logger *zap.SugaredLogger, config *config.Config) *Terminal {

	return &Terminal{
		buffers: []*buffer.Buffer{
			buffer.NewBuffer(1, 1, buffer.CellAttributes{
				FgColour: config.ColourScheme.Foreground,
				BgColour: config.ColourScheme.Background,
			}),
			buffer.NewBuffer(1, 1, buffer.CellAttributes{
				FgColour: config.ColourScheme.Foreground,
				BgColour: config.ColourScheme.Background,
			}),
		},
		pty:           pty,
		logger:        logger,
		config:        config,
		titleHandlers: []chan bool{},
		pauseChan:     make(chan bool, 1),
		resumeChan:    make(chan bool, 1),
		modes: Modes{
			ShowCursor: true,
		},
	}

}

func (terminal *Terminal) SetBracketedPasteMode(enabled bool) {
	terminal.bracketedPasteMode = enabled
}

func (terminal *Terminal) CheckDirty() bool {
	d := terminal.isDirty
	terminal.isDirty = false
	return d || terminal.ActiveBuffer().IsDirty()
}

func (terminal *Terminal) SetDirty() {
	terminal.isDirty = true
}

func (terminal *Terminal) IsApplicationCursorKeysModeEnabled() bool {
	return terminal.modes.ApplicationCursorKeys
}

func (terminal *Terminal) SetMouseMode(mode MouseMode) {
	terminal.mouseMode = mode
}

func (terminal *Terminal) GetMouseMode() MouseMode {
	return terminal.mouseMode
}

func (terminal *Terminal) UseMainBuffer() {
	terminal.activeBufferIndex = MainBuffer
	terminal.SetSize(uint(terminal.size.Width), uint(terminal.size.Height))
}

func (terminal *Terminal) UseAltBuffer() {
	terminal.activeBufferIndex = AltBuffer
	terminal.SetSize(uint(terminal.size.Width), uint(terminal.size.Height))
}

func (terminal *Terminal) ActiveBuffer() *buffer.Buffer {
	return terminal.buffers[terminal.activeBufferIndex]
}

func (terminal *Terminal) GetScrollOffset() uint {
	return terminal.ActiveBuffer().GetScrollOffset()
}

func (terminal *Terminal) ScrollDown(lines uint16) {
	terminal.logger.Infof("Scrolling down %d", lines)
	terminal.ActiveBuffer().ScrollDown(lines)

}

func (terminal *Terminal) ScrollUp(lines uint16) {
	terminal.logger.Infof("Scrolling up %d", lines)
	terminal.ActiveBuffer().ScrollUp(lines)
}

func (terminal *Terminal) ScrollPageDown() {
	terminal.ActiveBuffer().ScrollPageDown()
}
func (terminal *Terminal) ScrollPageUp() {
	terminal.ActiveBuffer().ScrollPageUp()
}
func (terminal *Terminal) ScrollToEnd() {
	terminal.ActiveBuffer().ScrollToEnd()
}

func (terminal *Terminal) GetVisibleLines() []buffer.Line {
	return terminal.ActiveBuffer().GetVisibleLines()
}

func (terminal *Terminal) GetCell(col uint16, row uint16) *buffer.Cell {
	return terminal.ActiveBuffer().GetCell(col, row)
}

func (terminal *Terminal) AttachTitleChangeHandler(handler chan bool) {
	terminal.titleHandlers = append(terminal.titleHandlers, handler)
}

func (terminal *Terminal) Modes() Modes {
	return terminal.modes
}

func (terminal *Terminal) emitTitleChange() {
	for _, h := range terminal.titleHandlers {
		go func(c chan bool) {
			c <- true
		}(h)
	}
}

func (terminal *Terminal) GetLogicalCursorX() uint16 {
	if terminal.ActiveBuffer().CursorColumn() >= terminal.ActiveBuffer().Width() {
		return 0
	}

	return terminal.ActiveBuffer().CursorColumn()
}

func (terminal *Terminal) GetLogicalCursorY() uint16 {
	if terminal.ActiveBuffer().CursorColumn() >= terminal.ActiveBuffer().Width() {
		return terminal.ActiveBuffer().CursorLine() + 1
	}

	return terminal.ActiveBuffer().CursorLine()
}

func (terminal *Terminal) GetTitle() string {
	return terminal.title
}

func (terminal *Terminal) SetTitle(title string) {
	terminal.title = title
	terminal.emitTitleChange()
}

// Write sends data, i.e. locally typed keystrokes to the pty
func (terminal *Terminal) Write(data []byte) error {
	_, err := terminal.pty.Write(data)
	return err
}

func (terminal *Terminal) Paste(data []byte) error {

	if terminal.bracketedPasteMode {
		data = []byte(fmt.Sprintf("\x1b[200~%s\x1b[201~", string(data)))
	}
	_, err := terminal.pty.Write(data)
	return err
}

// Read needs to be run on a goroutine, as it continually reads output to set on the terminal
func (terminal *Terminal) Read() error {

	buffer := make(chan rune, 0xffff)

	reader := bufio.NewReader(terminal.pty)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go terminal.processInput(ctx, buffer)
	for {
		r, size, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		} else if size > 0 {
			buffer <- r
		}
	}

	//clean exit
	return nil
}

func (terminal *Terminal) Clear() {
	terminal.ActiveBuffer().Clear()
}

func (terminal *Terminal) GetSize() (int, int) {
	return int(terminal.size.Width), int(terminal.size.Height)
}

func (terminal *Terminal) SetSize(newCols uint, newLines uint) error {
	terminal.lock.Lock()
	defer terminal.lock.Unlock()

	terminal.size.Width = uint16(newCols)
	terminal.size.Height = uint16(newLines)

	terminal.ActiveBuffer().ResizeView(terminal.size.Width, terminal.size.Height)

	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(terminal.pty.Fd()),
		uintptr(syscall.TIOCSWINSZ), uintptr(unsafe.Pointer(&terminal.size)))
	if err != 0 {
		return fmt.Errorf("Failed to set terminal size vai ioctl: Error no %d", err)
	}

	return nil
}
