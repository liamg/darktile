package gui

type RenderState struct {
	cells map[[2]uint]RenderedCell
}

func NewRenderState() *RenderState {
	return &RenderState{
		cells: map[[2]uint]RenderedCell{},
	}
}

type RenderedCell struct {
	bg       [3]float32
	fg       [3]float32
	contents rune
	dirty    bool
}

func (rs *RenderState) Reset() {
	rs.cells = map[[2]uint]RenderedCell{}
}

func (rs *RenderState) SetDirty(x uint, y uint) {
	rs.cells[[2]uint{x, y}] = RenderedCell{
		bg:       [3]float32{0, 0, 0},
		fg:       [3]float32{0, 0, 0},
		contents: 0,
		dirty:    true,
	}
}

func (rs *RenderState) RequiresRender(x uint, y uint, bg [3]float32, fg [3]float32, contents rune, empty bool) bool {

	state, found := rs.cells[[2]uint{x, y}]
	if !found {
		if empty {
			//return false
		}
		rs.cells[[2]uint{x, y}] = RenderedCell{
			bg:       bg,
			fg:       fg,
			contents: contents,
		}
		return true
	}

	if state.bg != bg || state.fg != fg || state.contents != contents || state.dirty {
		rs.cells[[2]uint{x, y}] = RenderedCell{
			bg:       bg,
			fg:       fg,
			contents: contents,
		}
		return true
	}

	return false
}
