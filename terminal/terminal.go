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

	"gitlab.com/liamg/raft/buffer"
	"gitlab.com/liamg/raft/config"
	"go.uber.org/zap"
)

type Terminal struct {
	buffer        *buffer.Buffer
	lock          sync.Mutex
	pty           *os.File
	logger        *zap.SugaredLogger
	title         string
	size          Winsize
	config        config.Config
	titleHandlers []chan bool
	pauseChan     chan bool
	resumeChan    chan bool
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

func New(pty *os.File, logger *zap.SugaredLogger, config config.Config) *Terminal {

	return &Terminal{
		buffer: buffer.NewBuffer(0, 0, buffer.CellAttributes{
			FgColour: config.ColourScheme.DefaultFg,
			BgColour: config.ColourScheme.DefaultBg,
		}),
		pty:           pty,
		logger:        logger,
		config:        config,
		titleHandlers: []chan bool{},
		pauseChan:     make(chan bool, 1),
		resumeChan:    make(chan bool, 1),
	}
}

func (terminal *Terminal) GetCell(col int, row int) *buffer.Cell {
	return terminal.buffer.GetCell(col, row)
}

func (terminal *Terminal) AttachDisplayChangeHandler(handler chan bool) {
	terminal.buffer.AttachDisplayChangeHandler(handler)
}

func (terminal *Terminal) AttachTitleChangeHandler(handler chan bool) {
	terminal.titleHandlers = append(terminal.titleHandlers, handler)
}

func (terminal *Terminal) emitTitleChange() {
	for _, h := range terminal.titleHandlers {
		go func(c chan bool) {
			c <- true
		}(h)
	}
}

func (terminal *Terminal) GetLogicalCursorX() uint16 {
	if terminal.buffer.CursorColumn() >= terminal.buffer.Width() {
		return 0
	}

	return terminal.buffer.CursorColumn()
}

func (terminal *Terminal) GetLogicalCursorY() uint16 {
	if terminal.buffer.CursorColumn() >= terminal.buffer.Width() {
		return terminal.buffer.CursorLine() + 1
	}

	return terminal.buffer.CursorLine()
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
	terminal.buffer.Clear()
}

func (terminal *Terminal) GetSize() (int, int) {
	return int(terminal.size.Width), int(terminal.size.Height)
}

func (terminal *Terminal) SetSize(newCols int, newLines int) error {
	terminal.lock.Lock()
	defer terminal.lock.Unlock()

	terminal.size.Width = uint16(newCols)
	terminal.size.Height = uint16(newLines)

	terminal.buffer.ResizeView(terminal.size.Width, terminal.size.Height)

	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(terminal.pty.Fd()),
		uintptr(syscall.TIOCSWINSZ), uintptr(unsafe.Pointer(&terminal.size)))
	if err != 0 {
		return fmt.Errorf("Failed to set terminal size vai ioctl: Error no %d", err)
	}

	return nil
}

/*
------------------ ->
ssssssssssssssssss
ssssPPPPPPPPPPPPPP
xxxxxxxxx
xxxxxxxxxxxxxxxxxx
--------------------------
ssssssssssssssssss
SsssPPPPPPPPPPPPPP
xxxxxxxxx
xxxxxxxxxxxxxxxxxx




*/
