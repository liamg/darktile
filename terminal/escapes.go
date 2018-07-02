package terminal

import (
	"strconv"
)

func (terminal *Terminal) processInput(buffer chan rune) {

	// https://en.wikipedia.org/wiki/ANSI_escape_code

	for {
		b := <-buffer

		if b == 0x1b { // if the byte is an escape character, read the next byte to determine which one
			b = <-buffer
			switch b {
			case 0x5b: // CSI: Control Sequence Introducer ]
				var final rune
				params := []string{}
			CSI:
				for {
					b = <-buffer
					param := ""
					switch true {
					case b >= 0x30 && b <= 0x3F:
						param = param + string(b)
					case b >= 0x20 && b <= 0x2F:
						params = append(params, param)
						param = ""
					case b >= 0x40 && b <= 0x7e:
						params = append(params, param)
						final = b
						break CSI
					}
				}

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

					terminal.position.Line += distance
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
						if err != nil {
							distance = 1
						}
					}

					terminal.position.Col = distance - 1 // 1 based to 0 based

				case 'H', 'f':

					x, y := 1, 1
					if len(params) == 2 {
						var err error
						if params[0] != "" {
							x, err = strconv.Atoi(string(params[0]))
							if err != nil {
								x = 1
							}
						}
						if params[1] != "" {
							y, err = strconv.Atoi(string(params[y]))
							if err != nil {
								y = 1
							}
						}
						terminal.position.Col = x - 1
						terminal.position.Line = y - 1
					}
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
							for i := 0; i < terminal.position.Col; i++ {
								line.Cells[i].r = 0
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
						terminal.position.Col = 0
						terminal.position.Line = 0
					case "3":
						terminal.lines = []Line{}
						terminal.position.Col = 0
						terminal.position.Line = 0
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
						terminal.ClearToEndOfLine()
					case "1":
						line := terminal.getBufferedLine(terminal.position.Line)
						if line != nil {
							for i := 0; i < terminal.position.Col; i++ {
								line.Cells[i].r = 0
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
