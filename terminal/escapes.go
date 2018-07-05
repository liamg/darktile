package terminal

import (
	"strconv"
	"strings"
)

// Wish list here: http://invisible-island.net/xterm/ctlseqs/ctlseqs.html

type TerminalCharSet int

func (terminal *Terminal) processInput(buffer chan rune) {

	// https://en.wikipedia.org/wiki/ANSI_escape_code

	lineOverflow := false

	for {

		b := <-buffer

		if b == 0x1b { // if the byte is an escape character, read the next byte to determine which one
			b = <-buffer
			switch b {
			case '[': // CSI: Control Sequence Introducer [
				var final rune
				param := ""
				intermediate := ""
			CSI:
				for {
					b = <-buffer
					switch true {
					case b >= 0x30 && b <= 0x3F:
						param = param + string(b)
					case b >= 0x20 && b <= 0x2F:
						//intermediate? useful?
						intermediate += string(b)
					case b >= 0x40 && b <= 0x7e:
						final = b
						break CSI
					}
				}

				params := strings.Split(param, ";")

				switch final {
				case 'A':
					distance := 1
					if len(params) > 0 {
						var err error
						distance, err = strconv.Atoi(params[0])
						if err != nil {
							distance = 1
						}
					}
					if terminal.position.Line-distance >= 0 {
						terminal.position.Line -= distance
					} else {
						terminal.position.Line = 0
					}
				case 'B':
					distance := 1
					if len(params) > 0 {
						var err error
						distance, err = strconv.Atoi(params[0])
						if err != nil {
							distance = 1
						}
					}

					_, h := terminal.GetSize()
					if terminal.position.Line+distance >= h {
						terminal.position.Line = h - 1
					} else {
						terminal.position.Line += distance
					}
				case 'C':

					distance := 1
					if len(params) > 0 {
						var err error
						distance, err = strconv.Atoi(params[0])
						if err != nil {
							distance = 1
						}
					}

					terminal.position.Col += distance
					w, _ := terminal.GetSize()
					if terminal.position.Col >= w {
						terminal.position.Col = w - 1
					}

				case 'D':

					distance := 1
					if len(params) > 0 {
						var err error
						distance, err = strconv.Atoi(params[0])
						if err != nil {
							distance = 1
						}
					}

					terminal.position.Col -= distance
					if terminal.position.Col < 0 {
						terminal.position.Col = 0
					}

				case 'E':
					distance := 1
					if len(params) > 0 {
						var err error
						distance, err = strconv.Atoi(params[0])
						if err != nil {
							distance = 1
						}
					}

					terminal.position.Line += distance
					terminal.position.Col = 0

				case 'F':

					distance := 1
					if len(params) > 0 {
						var err error
						distance, err = strconv.Atoi(params[0])
						if err != nil {
							distance = 1
						}
					}
					if terminal.position.Line-distance >= 0 {
						terminal.position.Line -= distance
					} else {
						terminal.position.Line = 0
					}
					terminal.position.Col = 0

				case 'G':

					distance := 1
					if len(params) > 0 {
						var err error
						distance, err = strconv.Atoi(params[0])
						if err != nil || params[0] == "" {
							distance = 1
						}
					}

					terminal.position.Col = distance - 1 // 1 based to 0 based

				case 'H', 'f':

					x, y := 1, 1
					if len(params) == 2 {
						var err error
						if params[0] != "" {
							y, err = strconv.Atoi(string(params[0]))
							if err != nil {
								y = 1
							}
						}
						if params[1] != "" {
							x, err = strconv.Atoi(string(params[1]))
							if err != nil {
								x = 1
							}
						}
					}
					terminal.position.Col = x - 1
					terminal.position.Line = y - 1

				case 'J':

					n := "0"
					if len(params) > 0 {
						n = params[0]
					}

					switch n {

					case "0", "":
						line := terminal.getBufferedLine(terminal.position.Line)
						if line != nil {
							line.Cells = line.Cells[:terminal.position.Col]
						}
						_, h := terminal.GetSize()
						for i := terminal.position.Line + 1; i < h; i++ {
							line := terminal.getBufferedLine(i)
							if line != nil {
								line.Cells = []Cell{}
							}
						}
					case "1":
						line := terminal.getBufferedLine(terminal.position.Line)
						if line != nil {
							for i := 0; i <= terminal.position.Col; i++ {
								if i < len(line.Cells) {
									line.Cells[i].r = 0
								}
							}
						}
						for i := 0; i < terminal.position.Line; i++ {
							line := terminal.getBufferedLine(i)
							if line != nil {
								line.Cells = []Cell{}
							}
						}

					case "2":
						_, h := terminal.GetSize()
						for i := 0; i < h; i++ {
							line := terminal.getBufferedLine(i)
							if line != nil {
								line.Cells = []Cell{}
							}
						}
					case "3":
						terminal.lines = []Line{}

					default:
						terminal.logger.Errorf("Unknown CSI ED sequence: %s", n)
					}

				case 'K': // K - EOL - Erase to end of line
					n := "0"
					if len(params) > 0 {
						n = params[0]
					}

					switch n {
					case "0", "":
						line := terminal.getBufferedLine(terminal.position.Line)
						if line != nil {
							line.Cells = line.Cells[:terminal.position.Col]
						}
					case "1":
						line := terminal.getBufferedLine(terminal.position.Line)
						if line != nil {
							for i := 0; i <= terminal.position.Col; i++ {
								if i < len(line.Cells) {
									line.Cells[i].r = 0
								}
							}
						}
					case "2":
						line := terminal.getBufferedLine(terminal.position.Line)
						if line != nil {
							line.Cells = []Cell{}
						}
					default:
						terminal.logger.Errorf("Unsupported EL: %s", n)
					}

				case 'm':
					// SGR: colour and shit
					for i := range params {
						param := params[i]
						switch param {
						case "00", "0", "":
							terminal.cellAttr = terminal.defaultCellAttr
						case "1", "01":
							terminal.cellAttr.Bold = true
						case "2", "02":
							terminal.cellAttr.Dim = true
						case "4", "04":
							terminal.cellAttr.Underline = true
						case "5", "05":
							terminal.cellAttr.Blink = true
						case "7", "07":
							terminal.cellAttr.Reverse = true
						case "8", "08":
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
							terminal.cellAttr.FgColour = terminal.colourScheme.WhiteFg
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
							terminal.cellAttr.BgColour = terminal.colourScheme.WhiteBg
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
						default:

							terminal.logger.Errorf("Unknown SGR control sequence: (ESC[%s%s%s)", param, intermediate, string(final))
						}

						//terminal.logger.Debugf("SGR control sequence: (ESC[%s%s%s)", param, intermediate, string(final))
					}

				default:
					switch param + intermediate + string(final) {
					case "?25h":
						terminal.showCursor()
					case "?25l":
						terminal.hideCursor()
					case "?12h":
						// todo enable cursor blink
					case "?12l":
						// todo disable cursor blink
					default:
						terminal.logger.Errorf("Unknown CSI control sequence: 0x%02X (ESC[%s%s%s)", final, param, intermediate, string(final))
					}

				}
				terminal.logger.Debugf("Received CSI control sequence: 0x%02X (ESC[%s%s%s)", final, param, intermediate, string(final))
			case 0x5d: // OSC: Operating System Command
				b = <-buffer
				switch b {
				case rune('0'):
					b = <-buffer
					if b == rune(';') {
						title := []rune{}
						for {
							b = <-buffer
							if b == 0x07 || b == 0x5c { // 0x07 -> BELL, 0x5c -> ST (\)
								break
							}
							title = append(title, b)
						}
						terminal.title = string(title)
					} else {
						terminal.logger.Errorf("Invalid OSC 0 control sequence: 0x%02X", b)
					}
				default:
					terminal.logger.Errorf("Unknown OSC control sequence: 0x%02X", b)
				}
			case 'c':
				terminal.logger.Errorf("RIS not yet supported")
			case '(':
				b = <-buffer
				switch b {
				case 'A': //uk @todo handle these?
					//terminal.charSet = C0
				case 'B': //us
					//terminal.charSet = C0
				}
			case ')':
				b = <-buffer
				switch b {
				case 'A': //uk @todo handle these?
					//terminal.charSet = C1
				case 'B': //us
					//terminal.charSet = C1
				}
			case '*':
				b = <-buffer
				switch b {
				case 'A': //uk @todo handle these?
					//terminal.charSet = C2
				case 'B': //us
					//terminal.charSet = C2
				}
			case '+':
				b = <-buffer
				switch b {
				case 'A': //uk @todo handle these?
					//terminal.charSet = C3
				case 'B': //us
					//terminal.charSet = C3
				}
			case '>':
				// numeric char selection @todo
			case '=':
				//alternate char selection @todo
			case '?':
				pm := ""
				for {
					b = <-buffer
					switch b {
					case 'h':
						switch pm {
						default:
							terminal.logger.Errorf("Unknown private code ESC?%sh", pm)
						}
					case 'l':
						switch pm {
						default:
							terminal.logger.Errorf("Unknown private code ESC?%sl", pm)
						}
					default:
						pm += string(b)
					}
				}
			default:
				terminal.logger.Errorf("Unknown control sequence: 0x%02X [%s]", b, string(b))
			}
		} else {

			//fmt.Printf("%s", string(b))

			if b != 0x0d {
				lineOverflow = false
			}

			switch b {
			case 0x0a:

				_, h := terminal.GetSize()
				if terminal.position.Line+1 >= h {
					terminal.lines = append(terminal.lines, NewLine())
				} else {
					terminal.position.Line++
				}

			case 0x0d:
				if terminal.position.Col == 0 && terminal.position.Line > 0 && lineOverflow {
					terminal.position.Line--
					terminal.logger.Debugf("Swallowing forced new line for CR")
					lineOverflow = false
				}
				terminal.position.Col = 0

			case 0x08:
				// backspace
				terminal.position.Col--
				if terminal.position.Col < 0 {
					terminal.position.Col = 0
				}
			case 0x07:
				// @todo ring bell
			default:
				// render character at current location
				//		fmt.Printf("%s\n", string([]byte{b}))
				if b >= 0x20 {
					terminal.writeRune(b)
					lineOverflow = terminal.position.Col == 0
				} else {
					terminal.logger.Error("Non-readable rune received: 0x%X", b)
				}
			}

		}
		terminal.triggerOnUpdate()
	}
}
