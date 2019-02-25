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
type MouseExtMode uint

const (
	MouseModeNone MouseMode = iota
	MouseModeX10
	MouseModeVT200
	MouseModeVT200Highlight
	MouseModeButtonEvent
	MouseModeAnyEvent
	MouseExtNone MouseExtMode = iota
	MouseExtUTF
	MouseExtSGR
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
	resizeHandlers            []chan bool
	reverseHandlers           []chan bool
	modes                     Modes
	mouseMode                 MouseMode
	mouseExtMode              MouseExtMode
	bracketedPasteMode        bool
	isDirty                   bool
	charWidth                 float32
	charHeight                float32
	lastBuffer                uint8
	terminalState             *buffer.TerminalState
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
		terminalState: buffer.NewTerminalState(1, 1, buffer.CellAttributes{
			FgColour: config.ColourScheme.Foreground,
			BgColour: config.ColourScheme.Background,
		}, config.MaxLines),
		pty:           pty,
		logger:        logger,
		config:        config,
		titleHandlers: []chan bool{},
		modes: Modes{
			ShowCursor: true,
		},
		platformDependentSettings: pty.GetPlatformDependentSettings(),
	}
	t.buffers = []*buffer.Buffer{
		buffer.NewBuffer(t.terminalState),
		buffer.NewBuffer(t.terminalState),
		buffer.NewBuffer(t.terminalState),
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

func (terminal *Terminal) SetMouseExtMode(mode MouseExtMode) {
	terminal.mouseExtMode = mode
}

func (terminal *Terminal) GetMouseExtMode() MouseExtMode {
	return terminal.mouseExtMode
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
	return terminal.terminalState.GetScrollOffset()
}

func (terminal *Terminal) ScreenScrollDown(lines uint16) {
	defer terminal.SetDirty()
	buffer := terminal.ActiveBuffer()

	if buffer.Height() < int(buffer.ViewHeight()) {
		return
	}

	offset := terminal.terminalState.GetScrollOffset()
	if uint(lines) > offset {
		lines = uint16(offset)
	}
	terminal.terminalState.SetScrollOffset(offset - uint(lines))
}

func (terminal *Terminal) SetCharSize(w float32, h float32) {
	terminal.charWidth = w
	terminal.charHeight = h
}

func (terminal *Terminal) AreaScrollUp(lines uint16) {
	terminal.ActiveBuffer().AreaScrollUp(lines)
}

func (terminal *Terminal) AreaScrollDown(lines uint16) {
	terminal.ActiveBuffer().AreaScrollDown(lines)
}

func (terminal *Terminal) ScreenScrollUp(lines uint16) {
	defer terminal.SetDirty()
	buffer := terminal.ActiveBuffer()

	if buffer.Height() < int(buffer.ViewHeight()) {
		return
	}

	offset := terminal.terminalState.GetScrollOffset()

	if uint(lines)+offset >= (uint(buffer.Height()) - uint(buffer.ViewHeight())) {
		terminal.terminalState.SetScrollOffset(uint(buffer.Height()) - uint(buffer.ViewHeight()))
	} else {
		terminal.terminalState.SetScrollOffset(offset + uint(lines))
	}
}

func (terminal *Terminal) ScrollPageDown() {
	terminal.ScreenScrollDown(terminal.terminalState.ViewHeight())
}
func (terminal *Terminal) ScrollPageUp() {
	terminal.ScreenScrollUp(terminal.terminalState.ViewHeight())
}

func (terminal *Terminal) ScrollToEnd() {
	defer terminal.SetDirty()
	terminal.terminalState.SetScrollOffset(0)
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

func (terminal *Terminal) AttachResizeHandler(handler chan bool) {
	terminal.resizeHandlers = append(terminal.resizeHandlers, handler)
}

func (terminal *Terminal) AttachReverseHandler(handler chan bool) {
	terminal.reverseHandlers = append(terminal.reverseHandlers, handler)
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

func (terminal *Terminal) emitResize() {
	for _, h := range terminal.resizeHandlers {
		go func(c chan bool) {
			c <- true
		}(h)
	}
}

func (terminal *Terminal) emitReverse(reverse bool) {
	for _, h := range terminal.reverseHandlers {
		go func(c chan bool) {
			c <- reverse
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
		return terminal.ActiveBuffer().CursorLineAbsolute() + 1
	}

	return terminal.ActiveBuffer().CursorLineAbsolute()
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

func (terminal *Terminal) WriteReturn() error {
	if terminal.terminalState.IsNewLineMode() {
		return terminal.Write([]byte{0x0d, 0x0a})
	} else {
		return terminal.Write([]byte{0x0d})
	}
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

	if terminal.size.Width == uint16(newCols) && terminal.size.Height == uint16(newLines) {
		return nil
	}

	err := terminal.pty.Resize(int(newCols), int(newLines))
	if err != nil {
		return fmt.Errorf("Failed to set terminal size vai ioctl: Error no %d", err)
	}

	terminal.size.Width = uint16(newCols)
	terminal.size.Height = uint16(newLines)

	terminal.ActiveBuffer().ResizeView(terminal.size.Width, terminal.size.Height)

	terminal.emitResize()
	return nil
}

func (terminal *Terminal) SetAutoWrap(enabled bool) {
	terminal.terminalState.AutoWrap = enabled
}

func (terminal *Terminal) IsAutoWrap() bool {
	return terminal.terminalState.AutoWrap
}

func (terminal *Terminal) SetOriginMode(enabled bool) {
	terminal.terminalState.OriginMode = enabled
	terminal.ActiveBuffer().SetPosition(0, 0)
}

func (terminal *Terminal) SetInsertMode() {
	terminal.terminalState.ReplaceMode = false
}

func (terminal *Terminal) SetReplaceMode() {
	terminal.terminalState.ReplaceMode = true
}

func (terminal *Terminal) SetNewLineMode() {
	terminal.terminalState.LineFeedMode = false
}

func (terminal *Terminal) SetLineFeedMode() {
	terminal.terminalState.LineFeedMode = true
}

func (terminal *Terminal) ResetVerticalMargins() {
	terminal.terminalState.ResetVerticalMargins()
}

func (terminal *Terminal) SetScreenMode(enabled bool) {
	if terminal.terminalState.ScreenMode == enabled {
		return
	}
	terminal.terminalState.ScreenMode = enabled
	terminal.terminalState.CursorAttr.ReverseVideo()
	for _, buffer := range terminal.buffers {
		buffer.ReverseVideo()
	}
	terminal.emitReverse(enabled)
}
