package sixel

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
	"strings"
)

type Sixel struct {
	px     map[uint]map[uint]colour
	width  uint
	height uint
}

type colour [3]uint8

func decompress(data string) string {

	output := ""

	inMarker := false
	countStr := ""

	for _, r := range data {

		if !inMarker {
			if r == '!' {
				inMarker = true
				countStr = ""
			} else {
				output += string(r)
			}
			continue
		}

		if r >= 0x30 && r <= 0x39 {
			countStr = fmt.Sprintf("%s%c", countStr, r)
		} else {
			count, _ := strconv.Atoi(countStr)
			output += strings.Repeat(string(r), count)
			inMarker = false
		}
	}

	return output
}

// pass in everything after ESC+P and before ST
func ParseString(data string) (*Sixel, error) {

	data = decompress(data)

	inHeader := true
	inColour := false

	six := Sixel{}
	var x, y uint

	colourStr := ""

	colourMap := map[string]colour{}
	var selectedColour colour

	headerStr := ""

	remainMode := false

	var ratio uint

	// read p1 p2 p3
	for i, r := range data {
		switch true {
		case inHeader:
			// todo read p1 p2 p3
			if r == 'q' {
				headers := strings.Split(headerStr, ";")
				switch headers[0] {
				case "0", "1":
					ratio = 5
				case "2":
					ratio = 3
				case "3", "4", "5", "6":
					ratio = 2
				case "7", "8", "9", "":
					ratio = 1
				}
				if len(headers) > 1 {
					remainMode = headers[1] == "1"
				}
				inHeader = false
			} else {
				headerStr = fmt.Sprintf("%s%c", headerStr, r)
			}
		case inColour:
			colourStr = fmt.Sprintf("%s%c", colourStr, r)
			if i+1 >= len(data) || data[i+1] < 0x30 || data[i+1] > 0x3b {
				// process colour string
				inColour = false
				parts := strings.Split(colourStr, ";")

				// select colour
				if len(parts) == 1 {
					c, ok := colourMap[parts[0]]
					if ok {
						selectedColour = c
					}
				} else if len(parts) == 5 {
					switch parts[1] {
					case "1":
						// HSL
						return nil, fmt.Errorf("HSL colours are not yet supported")
					case "2":
						// RGB
						r, _ := strconv.Atoi(parts[2])
						g, _ := strconv.Atoi(parts[3])
						b, _ := strconv.Atoi(parts[4])
						colourMap[parts[0]] = colour([3]uint8{
							uint8(r & 0xff),
							uint8(g & 0xff),
							uint8(b & 0xff),
						})
					default:
						return nil, fmt.Errorf("Unknown colour definition type: %s", parts[1])
					}
				} else {
					return nil, fmt.Errorf("Invalid colour directive: #%s", colourStr)
				}

				colourStr = ""
			}

		default:
			switch r {
			case '-':
				y += 6
				x = 0
			case '$':
				x = 0
			case '#':
				inColour = true
			default:
				if r < 63 || r > 126 {
					continue
				}
				b := (r & 0xff) - 0x3f
				var bit int
				for bit = 5; bit >= 0; bit-- {
					if b&(1<<uint(bit)) > 0 {
						six.setPixel(x, y+uint(bit), selectedColour, ratio)
					} else if !remainMode {
						// @todo use background colour here
						//six.setPixel(x, y+uint(bit), selectedColour)
					}
				}
				x++
			}
		}
	}
	return &six, nil
}

func (six *Sixel) setPixel(x, y uint, c colour, vhRatio uint) {

	if six.px == nil {
		six.px = map[uint]map[uint]colour{}
	}

	if _, exists := six.px[x]; !exists {
		six.px[x] = map[uint]colour{}
	}

	if x+1 > six.width {
		six.width = x
	}

	ay := vhRatio * y

	var i uint
	for i = 0; i < vhRatio; i++ {
		if ay+i+1 > six.height {
			six.height = ay + i + 1
		}
		six.px[x][ay+i] = c
	}

}

func (six *Sixel) RGBA() *image.RGBA {
	rgba := image.NewRGBA(image.Rect(0, 0, int(six.width), int(six.height)))

	for x, r := range six.px {
		for y, colour := range r {
			rgba.Set(int(x), int(six.height)-int(y), color.RGBA{
				R: colour[0],
				G: colour[1],
				B: colour[2],
				A: 255,
			})
		}
	}

	return rgba
}
