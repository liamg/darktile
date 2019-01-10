package matrix

// AutoMatrix -- automatically growing matrix
type AutoMatrix struct {
	lines [][]rune
}

// NewAutoMatrix creates a new auto-matrix
func NewAutoMatrix() *AutoMatrix {
	m := &AutoMatrix{
		lines: make([][]rune, 0),
	}
	return m
}

// ExtractFrom extracts from (x1, y1) until the end
func (matrix *AutoMatrix) ExtractFrom(x1, y1 int) []rune {
	result := make([]rune, 0)
	y := y1
	for y < len(matrix.lines) {
		if matrix.lines[y] != nil {
			xMin := 0
			if y == y1 {
				xMin = x1
			}
			result = append(result, matrix.lines[y][xMin:]...)
		}
		y++
	}
	return result
}

// Extract from (x1, y1) until (x2, y2)
func (matrix *AutoMatrix) Extract(x1, y1, x2, y2 int) []rune {
	result := make([]rune, 0)
	y := y1
	for y <= y2 && y < len(matrix.lines) {
		xMin := 0
		if y == y1 {
			xMin = x1
		}
		if matrix.lines[y] != nil {
			xMax := x2
			if y != y2 {
				xMax = len(matrix.lines[y]) - 1
			}
			result = append(result, matrix.lines[y][xMin:xMax]...)
		}
		y++
	}
	return result
}

func (matrix *AutoMatrix) setAtLine(value rune, x int, line []rune) []rune {
	if x >= len(line) {
		line = append(line, make([]rune, x-len(line)+1)...)
	}
	line[x] = value
	return line
}

// SetAt sets a value at (x, y) growing the matrix as necessary
func (matrix *AutoMatrix) SetAt(value rune, x int, y int) {

	if y >= len(matrix.lines) {
		matrix.lines = append(matrix.lines, make([][]rune, y-len(matrix.lines)+1)...)
	}
	matrix.lines[y] = matrix.setAtLine(value, x, matrix.lines[y])
}
