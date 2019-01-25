package buffer

import (
	"fmt"
	"strings"

	"github.com/liamg/aminal/hints"
)

func (buffer *Buffer) GetHintAtPosition(col uint16, viewRow uint16) *hints.Hint {

	row := buffer.convertViewLineToRawLine(viewRow) - uint64(buffer.terminalState.scrollLinesFromBottom)

	cell := buffer.GetRawCell(col, row)
	if cell == nil || cell.Rune() == 0x00 {
		return nil
	}

	candidate := ""

	for i := int(col); i >= 0; i-- {
		cell := buffer.GetRawCell(uint16(i), row)
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

	for i := col + 1; i < buffer.terminalState.viewWidth; i++ {
		cell := buffer.GetRawCell(i, row)
		if cell == nil {
			break
		}
		if isRuneWordSelectionMarker(cell.Rune()) {
			break
		}

		candidate = fmt.Sprintf("%s%c", candidate, cell.Rune())
	}

	line := buffer.lines[row]

	return hints.Get(strings.Trim(candidate, " "), line.String(), sx, viewRow)

}
