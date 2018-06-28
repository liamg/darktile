package terminal

import (
	"fmt"
	"os"
	"sync"
	"syscall"
	"unsafe"

	"go.uber.org/zap"
)

type Terminal struct {
	cells    [][]Cell // y, x
	lock     sync.Mutex
	pty      *os.File
	logger   *zap.SugaredLogger
	title    string
	position Position
}

type Winsize struct {
	Height uint16
	Width  uint16
	x      uint16 //ignored, but necessary for ioctl calls
	y      uint16 //ignored, but necessary for ioctl calls
}

type Position struct {
	Col int
	Row int
}

func New(pty *os.File, logger *zap.SugaredLogger) *Terminal {
	return &Terminal{
		cells:  [][]Cell{},
		pty:    pty,
		logger: logger,
	}
}

func (terminal *Terminal) GetTitle() string {
	return terminal.title
}

// Write sends data, i.e. locally typed keystrokes to the pty
func (terminal *Terminal) Write(data []byte) error {
	_, err := terminal.pty.Write(data)
	return err
}

// Read needs to be run on a goroutine, as it continually reads output to set on the terminal
func (terminal *Terminal) Read() error {

	buffer := make(chan byte, 0xffff)

	go func() {

		// https://en.wikipedia.org/wiki/ANSI_escape_code

		for {
			b := <-buffer

			if b == 0x1b { // if the byte is an escape character, read the next byte to determine which one
				b = <-buffer
				terminal.logger.Debugf("Escape: 0x%X", b)
				switch b {
				case 0x5b: // CSI: Control Sequence Introducer
					b = <-buffer
					switch b {
					default:
						terminal.logger.Debugf("Unknown CSI control sequence: 0x%X", b)
					}

				case 0x5d: // OSC: Operating System Command
					b = <-buffer
					switch b {
					case byte('0'):
						if <-buffer == byte(';') {
							title := []byte{}
							for {
								b = <-buffer
								if b == 0x07 {
									break
								}
								title = append(title, b)
							}
							terminal.logger.Debugf("Terminal title set to: %s", string(title))
							terminal.title = string(title)
						} else {
							terminal.logger.Debugf("Invalid OSC 0 control sequence")
						}
					default:
						terminal.logger.Debugf("Unknown OSC control sequence: 0x%X", b)
					}
				default:
					terminal.logger.Debugf("Unknown control sequence: 0x%X", b)
				}
			} else {
				// render character at current location
				terminal.writeRune([]rune(string([]byte{b}))[0])
			}

		}
	}()

	for {
		readBytes := make([]byte, 1024)
		n, err := terminal.pty.Read(readBytes)
		if err != nil {
			terminal.logger.Errorf("Failed to read from pty: %s", err)
		}
		if len(readBytes) > 0 {
			readBytes = readBytes[:n]
			for _, x := range readBytes {
				buffer <- x
			}
		}
	}
}

func (terminal *Terminal) writeRune(r rune) error {
	fmt.Println(string(r))
	err := terminal.setRuneAtPos(terminal.position, r)
	if err != nil {
		return err
	}
	w, h := terminal.GetSize()
	if terminal.position.Col < w-1 {
		terminal.position.Col++
	} else {
		terminal.position.Col = 0
		if terminal.position.Row <= h-1 {
			terminal.position.Row++
		} else {
			panic(fmt.Errorf("Not implemented - need to shuffle all rows up one"))
		}
	}
	return nil
}

func (terminal *Terminal) setRuneAtPos(pos Position, r rune) error {

	if len(terminal.cells) <= pos.Col {
		return fmt.Errorf("Col %d does not exist", pos.Col)
	}

	if len(terminal.cells) < 1 || len(terminal.cells[0]) <= pos.Row {
		return fmt.Errorf("Row %d does not exist", pos.Row)
	}

	terminal.cells[pos.Row][pos.Col].rune = r
	return nil
}

func (terminal *Terminal) GetSize() (int, int) {
	terminal.lock.Lock()
	defer terminal.lock.Unlock()
	if len(terminal.cells) == 0 {
		return 0, 0
	}
	return len(terminal.cells[0]), len(terminal.cells)
}

func (terminal *Terminal) SetSize(cols int, rows int) error {
	terminal.lock.Lock()
	defer terminal.lock.Unlock()
	cells := make([][]Cell, rows)
	for i := range cells {
		cells[i] = make([]Cell, cols)
	}
	terminal.cells = cells

	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(terminal.pty.Fd()),
		uintptr(syscall.TIOCSWINSZ), uintptr(unsafe.Pointer(&Winsize{Width: uint16(cols), Height: uint16(rows)})))
	if err != 0 {
		return fmt.Errorf("Failed to set terminal size vai ioctl: Error no %d", err)
	}
	return nil
}
