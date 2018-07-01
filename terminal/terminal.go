package terminal

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"go.uber.org/zap"
)

type Terminal struct {
	lines           []Line   // lines, where 0 is earliest, n is latest
	position        Position // line and col
	lock            sync.Mutex
	pty             *os.File
	logger          *zap.SugaredLogger
	title           string
	onUpdate        []func()
	size            Winsize
	colourScheme    ColourScheme
	cellAttr        CellAttributes
	defaultCellAttr CellAttributes
}

type Line struct {
	Cells   []Cell
	wrapped bool
}

func NewLine() Line {
	return Line{
		Cells: []Cell{},
	}
}

func (line *Line) String() string {
	s := ""
	for _, c := range line.Cells {
		s += string(c.r)
	}
	return s
}

func (line *Line) CutCellsAfter(n int) []Cell {
	cut := line.Cells[n:]
	line.Cells = line.Cells[:n]
	return cut
}

func (line *Line) CutCellsFromBeginning(n int) []Cell {
	if n > len(line.Cells) {
		n = len(line.Cells)
	}
	cut := line.Cells[:n]
	line.Cells = line.Cells[n:]
	return cut
}

func (line *Line) CutCellsFromEnd(n int) []Cell {
	cut := line.Cells[len(line.Cells)-n:]
	line.Cells = line.Cells[:len(line.Cells)-n]
	return cut
}

func (line *Line) GetRenderedLength() int {
	l := 0
	for x, c := range line.Cells {
		if c.r > 0 {
			l = x
		}
	}
	return l
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

func New(pty *os.File, logger *zap.SugaredLogger, colourScheme ColourScheme) *Terminal {

	defaultCellAttr := CellAttributes{
		FgColour: colourScheme.DefaultFg,
		BgColour: colourScheme.DefaultBg,
	}

	return &Terminal{
		lines: []Line{
			NewLine(),
		},
		pty:             pty,
		logger:          logger,
		onUpdate:        []func(){},
		cellAttr:        defaultCellAttr,
		defaultCellAttr: defaultCellAttr,
		colourScheme:    colourScheme,
	}
}

func (terminal *Terminal) OnUpdate(handler func()) {
	terminal.onUpdate = append(terminal.onUpdate, handler)
}

func (terminal *Terminal) triggerOnUpdate() {
	for _, handler := range terminal.onUpdate {
		go handler()
	}
}

func (terminal *Terminal) getPosition() Position {
	return terminal.position
}

func (terminal *Terminal) incrementPosition() {
	position := terminal.getPosition()
	if position.Col+1 >= int(terminal.size.Width) {
		position.Line++
		_, h := terminal.GetSize()
		if position.Line >= h {
			position.Line--
		}
		position.Col = 0
	} else {
		position.Col++
	}
	terminal.SetPosition(position)
}

func (terminal *Terminal) SetPosition(position Position) {
	terminal.position = position
}

func (terminal *Terminal) GetPosition() Position {
	return terminal.position
}

func (terminal *Terminal) GetTitle() string {
	return terminal.title
}

// Write sends data, i.e. locally typed keystrokes to the pty
func (terminal *Terminal) Write(data []byte) error {
	_, err := terminal.pty.Write(data)
	return err
}

func (terminal *Terminal) ClearToEndOfLine() {

	position := terminal.getPosition()

	line := terminal.getBufferedLine(position.Line)
	if line != nil {
		if position.Col < len(line.Cells) {
			line.Cells = line.Cells[:position.Col]
		}
	}

}

// we have thousands of lines of output. if the terminal is X lines high, we just want to lookat the most recent X lines to render (unless scroll etc)
func (terminal *Terminal) getBufferedLine(line int) *Line {

	if len(terminal.lines) >= int(terminal.size.Height) {
		line = len(terminal.lines) - int(terminal.size.Height) + line
	}

	if line >= len(terminal.lines) {
		return nil
	}

	return &terminal.lines[line]
}

func (terminal *Terminal) processInput(buffer chan rune) {

	// https://en.wikipedia.org/wiki/ANSI_escape_code

	for {
		b := <-buffer

		if b == 0x1b { // if the byte is an escape character, read the next byte to determine which one
			b = <-buffer
			switch b {
			case 0x5b: // CSI: Control Sequence Introducer ]
				var final rune
				params := []rune{}
				intermediate := []rune{}
			CSI:
				for {
					b = <-buffer
					switch true {
					case b >= 0x30 && b <= 0x3F:
						params = append(params, b)
					case b >= 0x20 && b <= 0x2F:
						intermediate = append(intermediate, b)
					case b >= 0x40 && b <= 0x7e:
						final = b
						break CSI
					}
				}

				switch final {
				case rune('A'):
					distance := 1
					if len(params) > 0 {
						var err error
						distance, err = strconv.Atoi(string(params[0]))
						if err != nil {
							distance = 1
						}
					}
					if terminal.position.Line-distance >= 0 {
						terminal.position.Line -= distance
					}
				case rune('B'):
					distance := 1
					if len(params) > 0 {
						var err error
						distance, err = strconv.Atoi(string(params[0]))
						if err != nil {
							distance = 1
						}
					}

					terminal.position.Line += distance

				case 0x4b: // K - EOL - Erase to end of line
					if len(params) == 0 || params[0] == rune('0') {
						terminal.ClearToEndOfLine()
					} else {
						terminal.logger.Errorf("Unsupported EL")
					}
				case rune('m'):
					// SGR: colour and shit
					sgr := string(params)
					sgrParams := strings.Split(sgr, ";")
					for i := range sgrParams {
						param := sgrParams[i]
						switch param {
						case "0":
							terminal.cellAttr = terminal.defaultCellAttr
						case "1":
							terminal.cellAttr.Bold = true
						case "2":
							terminal.cellAttr.Dim = true
						case "4":
							terminal.cellAttr.Underline = true
						case "5":
							terminal.cellAttr.Blink = true
						case "7":
							terminal.cellAttr.Reverse = true
						case "8":
							terminal.cellAttr.Hidden = true
						case "21":
							terminal.cellAttr.Bold = false
						case "22":
							terminal.cellAttr.Dim = false
						case "24":
							terminal.cellAttr.Underline = false
						case "25":
							terminal.cellAttr.Blink = false
						case "27":
							terminal.cellAttr.Reverse = false
						case "28":
							terminal.cellAttr.Hidden = false
						case "39":
							terminal.cellAttr.FgColour = terminal.colourScheme.DefaultFg
						case "30":
							terminal.cellAttr.FgColour = terminal.colourScheme.BlackFg
						case "31":
							terminal.cellAttr.FgColour = terminal.colourScheme.RedFg
						case "32":
							terminal.cellAttr.FgColour = terminal.colourScheme.GreenFg
						case "33":
							terminal.cellAttr.FgColour = terminal.colourScheme.YellowFg
						case "34":
							terminal.cellAttr.FgColour = terminal.colourScheme.BlueFg
						case "35":
							terminal.cellAttr.FgColour = terminal.colourScheme.MagentaFg
						case "36":
							terminal.cellAttr.FgColour = terminal.colourScheme.CyanFg
						case "37":
							terminal.cellAttr.FgColour = terminal.colourScheme.LightGreyFg
						case "90":
							terminal.cellAttr.FgColour = terminal.colourScheme.DarkGreyFg
						case "91":
							terminal.cellAttr.FgColour = terminal.colourScheme.LightRedFg
						case "92":
							terminal.cellAttr.FgColour = terminal.colourScheme.LightGreenFg
						case "93":
							terminal.cellAttr.FgColour = terminal.colourScheme.LightYellowFg
						case "94":
							terminal.cellAttr.FgColour = terminal.colourScheme.LightBlueFg
						case "95":
							terminal.cellAttr.FgColour = terminal.colourScheme.LightMagentaFg
						case "96":
							terminal.cellAttr.FgColour = terminal.colourScheme.LightCyanFg
						case "97":
							terminal.cellAttr.FgColour = terminal.colourScheme.WhiteFg
						case "49":
							terminal.cellAttr.BgColour = terminal.colourScheme.DefaultBg
						case "40":
							terminal.cellAttr.BgColour = terminal.colourScheme.BlackBg
						case "41":
							terminal.cellAttr.BgColour = terminal.colourScheme.RedBg
						case "42":
							terminal.cellAttr.BgColour = terminal.colourScheme.GreenBg
						case "43":
							terminal.cellAttr.BgColour = terminal.colourScheme.YellowBg
						case "44":
							terminal.cellAttr.BgColour = terminal.colourScheme.BlueBg
						case "45":
							terminal.cellAttr.BgColour = terminal.colourScheme.MagentaBg
						case "46":
							terminal.cellAttr.BgColour = terminal.colourScheme.CyanBg
						case "47":
							terminal.cellAttr.BgColour = terminal.colourScheme.LightGreenBg
						case "100":
							terminal.cellAttr.BgColour = terminal.colourScheme.DarkGreyBg
						case "101":
							terminal.cellAttr.BgColour = terminal.colourScheme.LightRedBg
						case "102":
							terminal.cellAttr.BgColour = terminal.colourScheme.LightGreenBg
						case "103":
							terminal.cellAttr.BgColour = terminal.colourScheme.LightYellowBg
						case "104":
							terminal.cellAttr.BgColour = terminal.colourScheme.LightBlueBg
						case "105":
							terminal.cellAttr.BgColour = terminal.colourScheme.LightMagentaBg
						case "106":
							terminal.cellAttr.BgColour = terminal.colourScheme.LightCyanBg
						case "107":
							terminal.cellAttr.BgColour = terminal.colourScheme.WhiteBg

						}
					}

				default:
					b = <-buffer
					terminal.logger.Errorf("Unknown CSI control sequence: 0x%02X (%s)", final, string(final))
				}
			case 0x5d: // OSC: Operating System Command
				b = <-buffer
				switch b {
				case rune('0'):
					b = <-buffer
					if b == rune(';') {
						title := []rune{}
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
						terminal.logger.Errorf("Invalid OSC 0 control sequence: 0x%02X", b)
					}
				default:
					terminal.logger.Errorf("Unknown OSC control sequence: 0x%02X", b)
				}
			case rune('c'):
				terminal.logger.Errorf("RIS not yet supported")
			case rune(')'), rune('('):
				b = <-buffer
				terminal.logger.Debugf("Ignoring character set control code )%s", string(b))
			default:
				terminal.logger.Errorf("Unknown control sequence: 0x%02X [%s]", b, string(b))
			}
		} else {

			switch b {
			case 0x0a:
				terminal.position.Line++
				_, h := terminal.GetSize()
				if terminal.position.Line >= h {
					terminal.position.Line--
				}
				terminal.lines = append(terminal.lines, NewLine())
			case 0x0d:
				terminal.position.Col = 0
			case 0x08:
				// backspace
				terminal.position.Col--
			case 0x07:
				// @todo ring bell
			default:
				// render character at current location
				//		fmt.Printf("%s\n", string([]byte{b}))
				terminal.writeRune(b)
			}

		}
		terminal.triggerOnUpdate()
	}
}

// Read needs to be run on a goroutine, as it continually reads output to set on the terminal
func (terminal *Terminal) Read() error {

	buffer := make(chan rune, 0xffff)

	reader := bufio.NewReader(terminal.pty)

	go terminal.processInput(buffer)
	for {
		r, size, err := reader.ReadRune()
		if err != nil {
			return err
		} else if size > 0 {
			buffer <- r
		}
	}
}

func (terminal *Terminal) writeRune(r rune) {
	terminal.setRuneAtPos(terminal.position, r)
	terminal.incrementPosition()

}

func (terminal *Terminal) Clear() {
	// @todo actually should just add a bunch of newlines?
	for i := 0; i < int(terminal.size.Height); i++ {
		terminal.lines = append(terminal.lines, NewLine())
	}
	terminal.SetPosition(Position{Line: 0, Col: 0})
}

func (terminal *Terminal) GetCellAtPos(pos Position) (*Cell, error) {

	if int(terminal.size.Height) <= pos.Line {
		terminal.logger.Errorf("Line %d does not exist", pos.Line)
		return nil, fmt.Errorf("Line %d does not exist", pos.Line)
	}

	if int(terminal.size.Width) <= pos.Col {
		terminal.logger.Errorf("Col %d does not exist", pos.Col)
		return nil, fmt.Errorf("Col %d does not exist", pos.Col)
	}

	line := terminal.getBufferedLine(pos.Line)
	if line == nil {
		return nil, fmt.Errorf("Line missing")
	}
	for pos.Col >= len(line.Cells) {
		line.Cells = append(line.Cells, terminal.NewCell())
	}
	return &line.Cells[pos.Col], nil
}

func (terminal *Terminal) setRuneAtPos(pos Position, r rune) error {

	if int(terminal.size.Width) <= pos.Col {
		terminal.logger.Errorf("Col %d does not exist", pos.Col)
		return fmt.Errorf("Col %d does not exist", pos.Col)
	}

	if pos.Line == 0 && pos.Col == 0 {
		fmt.Printf("\n\nSetting %d %d to %q\n\n\n", pos.Line, pos.Col, string(r))
	}

	for terminal.position.Line >= len(terminal.lines) {
		terminal.lines = append(terminal.lines, NewLine())
	}

	line := terminal.getBufferedLine(pos.Line)
	if line == nil {
		return fmt.Errorf("Impossible?")
	}

	for pos.Col >= len(line.Cells) {
		line.Cells = append(line.Cells, terminal.NewCell())
	}

	line.Cells[pos.Col].attr = terminal.cellAttr
	line.Cells[pos.Col].r = r
	return nil
}

func (terminal *Terminal) GetSize() (int, int) {
	return int(terminal.size.Width), int(terminal.size.Height)
}

func (terminal *Terminal) SetSize(newCols int, newLines int) error {
	terminal.lock.Lock()
	defer terminal.lock.Unlock()

	oldCols := int(terminal.size.Width)
	oldLines := int(terminal.size.Height)

	if oldLines > 0 && oldCols > 0 { // only bother resizing content if there is some
		if newCols < oldCols { // if the width decreased, we need to do some line trimming

			for l := range terminal.lines {
				if terminal.lines[l].GetRenderedLength() > newCols {
					cells := terminal.lines[l].CutCellsAfter(newCols)
					line := Line{
						Cells:   cells,
						wrapped: true,
					}
					terminal.lines = append(terminal.lines[:l+1], append([]Line{line}, terminal.lines[l+1:]...)...)
					if terminal.getPosition().Line > l {
						terminal.position.Line++
					} else if terminal.getPosition().Line == l {
						if terminal.getPosition().Col >= newCols {
							terminal.position.Line++
						}
					}
				}
			}

		} else if newCols > oldCols { // if width increased, we need to potentially unwrap some lines
			for l := 0; l < len(terminal.lines); l++ {
				if terminal.lines[l].GetRenderedLength() < newCols { // there is space here to unwrap a line if needed
					if l+1 < len(terminal.lines) {
						if terminal.lines[l+1].wrapped {
							wrapSize := newCols - terminal.lines[l].GetRenderedLength()
							cells := terminal.lines[l+1].CutCellsFromBeginning(wrapSize)
							terminal.lines[l].Cells = append(terminal.lines[l].Cells, cells...)
							if terminal.lines[l+1].GetRenderedLength() == 0 {
								// remove line
								terminal.lines = append(terminal.lines[:l+1], terminal.lines[l+2:]...)
								if terminal.getPosition().Line >= l+1 {
									terminal.position.Line--
								}
							}
						}
					}
				}
			}

		}
	}

	terminal.size.Width = uint16(newCols)
	terminal.size.Height = uint16(newLines)

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
