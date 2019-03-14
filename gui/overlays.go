package gui

type overlay interface {
	render(gui *GUI)
}

func (gui *GUI) setOverlay(m overlay) {
	gui.overlay = m
	gui.terminal.NotifyDirty()
}

func (gui *GUI) renderOverlay() {
	if gui.overlay == nil {
		return
	}

	gui.overlay.render(gui)
}
