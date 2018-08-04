package terminal

import (
	"fmt"
)

func (terminal *Terminal) delete(n int) error {
	if len(terminal.lines) <= terminal.position.Line {
		return fmt.Errorf("Cannot delete character at current position - line does not exist")
	}
	line := &terminal.lines[terminal.position.Line]

	if terminal.position.Col >= len(line.Cells) {
		return fmt.Errorf("Line not long enough to delete anything")
	}

	for terminal.position.Col+n > len(line.Cells) {
		n--
	}
	after := line.Cells[terminal.position.Col+n:]
	before := line.Cells[:terminal.position.Col]

	line.Cells = append(before, after...)

	// @todo rewrap lines here
	// so if line overflows and then we delete characters from beginnign of the line

	return nil
}

// @todo remove this debug func
func (terminal *Terminal) GetLineString() string {
	if len(terminal.lines) <= terminal.position.Line {
		return ""
	}
	line := &terminal.lines[terminal.position.Line]

	return line.String()
}
