package terminal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/liamg/aminal/buffer"
	"github.com/liamg/aminal/config"
)

func sgrSequenceHandler(params []string, terminal *Terminal) error {

	if len(params) == 0 {
		params = []string{"0"}
	}

	for i := range params {

		p := strings.Replace(strings.Replace(params[i], "[", "", -1), "]", "", -1)

		switch p {
		case "00", "0", "":
			attr := terminal.ActiveBuffer().CursorAttr()
			*attr = buffer.CellAttributes{
				FgColour: terminal.config.ColourScheme.Foreground,
				BgColour: terminal.config.ColourScheme.Background,
			}
		case "1", "01":
			terminal.ActiveBuffer().CursorAttr().Bold = true
		case "2", "02":
			terminal.ActiveBuffer().CursorAttr().Dim = true
		case "4", "04":
			terminal.ActiveBuffer().CursorAttr().Underline = true
		case "5", "05":
			terminal.ActiveBuffer().CursorAttr().Blink = true
		case "7", "07":
			terminal.ActiveBuffer().CursorAttr().Inverse = true
		case "8", "08":
			terminal.ActiveBuffer().CursorAttr().Hidden = true
		case "21":
			terminal.ActiveBuffer().CursorAttr().Bold = false
		case "22":
			terminal.ActiveBuffer().CursorAttr().Dim = false
		case "23":
			// not italic
		case "24":
			terminal.ActiveBuffer().CursorAttr().Underline = false
		case "25":
			terminal.ActiveBuffer().CursorAttr().Blink = false
		case "27":
			terminal.ActiveBuffer().CursorAttr().Inverse = false
		case "28":
			terminal.ActiveBuffer().CursorAttr().Hidden = false
		case "29":
			// not strikethrough
		case "39":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.Foreground
		case "30":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.Black
		case "31":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.Red
		case "32":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.Green
		case "33":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.Yellow
		case "34":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.Blue
		case "35":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.Magenta
		case "36":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.Cyan
		case "37":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.White
		case "90":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.DarkGrey
		case "91":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.LightRed
		case "92":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.LightGreen
		case "93":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.LightYellow
		case "94":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.LightBlue
		case "95":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.LightMagenta
		case "96":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.LightCyan
		case "97":
			terminal.ActiveBuffer().CursorAttr().FgColour = terminal.config.ColourScheme.White
		case "49":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.Background
		case "40":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.Black
		case "41":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.Red
		case "42":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.Green
		case "43":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.Yellow
		case "44":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.Blue
		case "45":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.Magenta
		case "46":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.Cyan
		case "47":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.White
		case "100":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.DarkGrey
		case "101":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.LightRed
		case "102":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.LightGreen
		case "103":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.LightYellow
		case "104":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.LightBlue
		case "105":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.LightMagenta
		case "106":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.LightCyan
		case "107":
			terminal.ActiveBuffer().CursorAttr().BgColour = terminal.config.ColourScheme.White
		case "38": // set foreground
			c, err := terminal.getANSIColour(params[i:])
			if err != nil {
				return err
			}
			terminal.ActiveBuffer().CursorAttr().FgColour = c
			return nil
		case "48": // set background
			c, err := terminal.getANSIColour(params[i:])
			if err != nil {
				return err
			}
			terminal.ActiveBuffer().CursorAttr().BgColour = c
			return nil
		default:
			return fmt.Errorf("Unknown SGR control sequence: (ESC[%sm)", params[i:])
		}
	}

	//terminal.logger.Debugf("SGR control sequence: (ESC[%s%sm)", param, intermediate)

	return nil
}

func (terminal *Terminal) getANSIColour(params []string) (config.Colour, error) {

	if len(params) > 2 {
		switch params[1] {
		case "5":
			// 8 bit colour
			colNum, err := strconv.Atoi(params[2])

			if err != nil || colNum >= 256 || colNum < 0 {
				return [3]float32{0, 0, 0}, fmt.Errorf("Invalid 8-bit colour specifier")
			}
			return terminal.get8BitSGRColour(uint8(colNum)), nil

		case "2":
			if len(params) < 4 {
				return [3]float32{0, 0, 0}, fmt.Errorf("Invalid true colour specifier")
			}
			// 24 bit colour
			if len(params) == 5 { // standard true colour

				r, err := strconv.Atoi(params[2])
				if err != nil {
					return [3]float32{0, 0, 0}, fmt.Errorf("Invalid true colour specifier")
				}
				g, err := strconv.Atoi(params[3])
				if err != nil {
					return [3]float32{0, 0, 0}, fmt.Errorf("Invalid true colour specifier")
				}
				b, err := strconv.Atoi(params[4])
				if err != nil {
					return [3]float32{0, 0, 0}, fmt.Errorf("Invalid true colour specifier")
				}
				return [3]float32{
					float32(r) / 0xff,
					float32(g) / 0xff,
					float32(b) / 0xff,
				}, nil
			} else if len(params) > 5 { // ISO/IEC International Standard 8613-6
				r, err := strconv.Atoi(params[3])
				if err != nil {
					return [3]float32{0, 0, 0}, fmt.Errorf("Invalid true colour specifier")
				}
				g, err := strconv.Atoi(params[4])
				if err != nil {
					return [3]float32{0, 0, 0}, fmt.Errorf("Invalid true colour specifier")
				}
				b, err := strconv.Atoi(params[5])
				if err != nil {
					return [3]float32{0, 0, 0}, fmt.Errorf("Invalid true colour specifier")
				}
				return [3]float32{
					float32(r) / 0xff,
					float32(g) / 0xff,
					float32(b) / 0xff,
				}, nil
			}
		}
	}

	return [3]float32{}, fmt.Errorf("Unknown ANSI colour format identifier")

}

func (terminal *Terminal) get8BitSGRColour(colNum uint8) [3]float32 {

	// https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit

	switch colNum {
	case 0:
		return terminal.config.ColourScheme.Black
	case 1:
		return terminal.config.ColourScheme.Red
	case 2:
		return terminal.config.ColourScheme.Green
	case 3:
		return terminal.config.ColourScheme.Yellow
	case 4:
		return terminal.config.ColourScheme.Blue
	case 5:
		return terminal.config.ColourScheme.Magenta
	case 6:
		return terminal.config.ColourScheme.Cyan
	case 7:
		return terminal.config.ColourScheme.White
	case 8:
		return terminal.config.ColourScheme.DarkGrey
	case 9:
		return terminal.config.ColourScheme.LightRed
	case 10:
		return terminal.config.ColourScheme.LightGreen
	case 11:
		return terminal.config.ColourScheme.LightYellow
	case 12:
		return terminal.config.ColourScheme.LightBlue
	case 13:
		return terminal.config.ColourScheme.LightMagenta
	case 14:
		return terminal.config.ColourScheme.LightCyan
	case 15:
		return terminal.config.ColourScheme.White
	}

	if colNum < 232 {

		r := 0
		g := 0
		b := 0

		index := int(colNum - 16) // 0-216

		for i := 0; i < index; i++ {
			if b == 0 {
				b = 95
			} else if b < 255 {
				b += 40
			} else {
				b = 0
				if g == 0 {
					g = 95
				} else if g < 255 {
					g += 40
				} else {
					g = 0
					if r == 0 {
						r = 95
					} else if r < 255 {
						r += 40
					} else {
						break
					}
				}
			}
		}

		return [3]float32{float32(r) / 0xff, float32(g) / 0xff, float32(b) / 0xff}
	}

	c := float32(colNum-232) / 0x18
	return [3]float32{c, c, c}
}
