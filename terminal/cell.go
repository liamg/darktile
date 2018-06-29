package terminal

import "github.com/go-gl/mathgl/mgl32"

type Cell struct {
	r         rune
	wrapper   bool
	isWrapped bool
}

func (cell *Cell) GetRune() rune {
	return cell.r
}

func (cell *Cell) IsHidden() bool {
	return cell.r == 0
}

func (cell *Cell) GetColour() (r float32, g float32, b float32) {

	if cell.wrapper {
		return 0, 1, 0
	}

	if cell.isWrapped {
		return 1, 1, 0
	}

	if cell.IsHidden() {
		return 0, 0, 1
	}

	return 1, 1, 1

}

func (cell *Cell) GetColourVec() mgl32.Vec3 {
	r, g, b := cell.GetColour()
	return mgl32.Vec3{r, g, b}
}
