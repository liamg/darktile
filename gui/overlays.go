package gui

type overlay interface {
	render(gui *GUI)
}

func (gui *GUI) setOverlay(m overlay) {
	defer gui.terminal.SetDirty()
	gui.overlay = m
}

func (gui *GUI) renderOverlay() {
	if gui.overlay == nil || !gui.terminal.UsingMainBuffer() {
		return
	}

	gui.overlay.render(gui)
}
