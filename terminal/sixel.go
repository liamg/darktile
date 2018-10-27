package terminal

import (
	"fmt"
	"math"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/liamg/aminal/sixel"
)

func sixelHandler(pty chan rune, terminal *Terminal) error {

	data := []rune{}

	for {
		b := <-pty
		if b == 0x1b { // terminated by ESC bell or ESC \
			_ = <-pty // swallow \ or bell
			break
		}
		data = append(data, b)
	}

	six, err := sixel.ParseString(string(data))
	if err != nil {
		return fmt.Errorf("Failed to parse sixel data: %s", err)
	}

	x, y := terminal.ActiveBuffer().CursorColumn(), terminal.ActiveBuffer().CursorLine()
	terminal.ActiveBuffer().Write(' ')
	cell := terminal.ActiveBuffer().GetCell(x, y)
	if cell == nil {
		return fmt.Errorf("Missing cell for sixel")
	}

	gl.UseProgram(terminal.program)
	cell.SetImage(six.RGBA())

	imageHeight := float64(cell.Image().Bounds().Size().Y)
	lines := int(math.Ceil(imageHeight / float64(terminal.charHeight)))
	for l := 0; l <= int(lines+1); l++ {
		terminal.ActiveBuffer().NewLine()
	}

	return nil
}
