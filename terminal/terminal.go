package terminal

import (
	"bufio"
	"fmt"
	"io"
	"sync"

	"github.com/liamg/aminal/buffer"
	"github.com/liamg/aminal/config"
	"github.com/liamg/aminal/platform"
	"go.uber.org/zap"
)

const (
	MainBuffer     uint8 = 0
	AltBuffer      uint8 = 1
	InternalBuffer uint8 = 2
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
	program                   uint32
	buffers                   []*buffer.Buffer
	activeBuffer              *buffer.Buffer
	lock                      sync.Mutex
	pty                       platform.Pty
	logger                    *zap.SugaredLogger
	title                     string
	size                      Winsize
	config                    *config.Config
	titleHandlers             []chan bool
	modes                     Modes
	mouseMode                 MouseMode
	bracketedPasteMode        bool
	isDirty                   bool
	charWidth                 float32
	charHeight                float32
	lastBuffer                uint8
	platformDependentSettings platform.PlatformDependentSettings
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

func New(pty platform.Pty, logger *zap.SugaredLogger, config *config.Config) *Terminal {
	t := &Terminal{
		buffers: []*buffer.Buffer{
			buffer.NewBuffer(1, 1, buffer.CellAttributes{
				FgColour: config.ColourScheme.Foreground,
				BgColour: config.ColourScheme.Background,
			}, config.MaxLines),
			buffer.NewBuffer(1, 1, buffer.CellAttributes{
				FgColour: config.ColourScheme.Foreground,
				BgColour: config.ColourScheme.Background,
			}, config.MaxLines),
			buffer.NewBuffer(1, 1, buffer.CellAttributes{
				FgColour: config.ColourScheme.Foreground,
				BgColour: config.ColourScheme.Background,
			}, config.MaxLines),
		},
		pty:           pty,
		logger:        logger,
		config:        config,
		titleHandlers: []chan bool{},
		modes: Modes{
			ShowCursor: true,
		},
		platformDependentSettings: pty.GetPlatformDependentSettings(),
	}
	t.activeBuffer = t.buffers[0]
	return t

}

func (terminal *Terminal) SetProgram(program uint32) {
	terminal.program = program
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

func (terminal *Terminal) IsOSCTerminator(char rune) bool {
	_, ok := terminal.platformDependentSettings.OSCTerminators[char]
	return ok
}

func (terminal *Terminal) UseMainBuffer() {
	terminal.activeBuffer = terminal.buffers[MainBuffer]
	terminal.SetSize(uint(terminal.size.Width), uint(terminal.size.Height))
}

func (terminal *Terminal) UseAltBuffer() {
	terminal.activeBuffer = terminal.buffers[AltBuffer]
	terminal.SetSize(uint(terminal.size.Width), uint(terminal.size.Height))
}

func (terminal *Terminal) UseInternalBuffer() {
	terminal.activeBuffer = terminal.buffers[InternalBuffer]
	terminal.SetSize(uint(terminal.size.Width), uint(terminal.size.Height))
}

func (terminal *Terminal) ExitInternalBuffer() {
	terminal.activeBuffer = terminal.buffers[terminal.lastBuffer]
}

func (terminal *Terminal) ActiveBuffer() *buffer.Buffer {
	return terminal.activeBuffer
}

func (terminal *Terminal) UsingMainBuffer() bool {
	return terminal.activeBuffer == terminal.buffers[MainBuffer]
}

func (terminal *Terminal) GetScrollOffset() uint {
	return terminal.ActiveBuffer().GetScrollOffset()
}

func (terminal *Terminal) ScrollDown(lines uint16) {
	terminal.ActiveBuffer().ScrollDown(lines)

}

func (terminal *Terminal) SetCharSize(w float32, h float32) {
	terminal.charWidth = w
	terminal.charHeight = h
}

func (terminal *Terminal) ScrollUp(lines uint16) {
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

	go terminal.processInput(buffer)
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		buffer <- r
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

	err := terminal.pty.Resize(int(newCols), int(newLines))
	if err != nil {
		return fmt.Errorf("Failed to set terminal size vai ioctl: Error no %d", err)
	}

	terminal.ActiveBuffer().ResizeView(terminal.size.Width, terminal.size.Height)
	return nil
}
