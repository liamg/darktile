package buffer

import (
	"fmt"
	"strings"

	"github.com/liamg/aminal/hints"
)

func (buffer *Buffer) GetHintAtPosition(col uint16, row uint16) *hints.Hint {

	cell := buffer.GetCell(col, row)
	if cell == nil || cell.Rune() == 0x00 {
		return nil
	}

	candidate := ""

	for i := col; i >= 0; i-- {
		cell := buffer.GetCell(i, row)
		if cell == nil {
			break
		}
		if isRuneWordSelectionMarker(cell.Rune()) {
			break
		}
		candidate = fmt.Sprintf("%c%s", cell.Rune(), candidate)
	}

	trimmed := strings.TrimLeft(candidate, " ")
	sx := col - uint16(len(trimmed)-1)

	for i := col + 1; i < buffer.viewWidth; i++ {
		cell := buffer.GetCell(i, row)
		if cell == nil {
			break
		}
		if isRuneWordSelectionMarker(cell.Rune()) {
			break
		}

		candidate = fmt.Sprintf("%s%c", candidate, cell.Rune())
	}

	line := buffer.lines[buffer.convertViewLineToRawLine(row)]

	return hints.Get(strings.Trim(candidate, " "), line.String(), sx, row)

}
