package render

func (r *Render) drawContent() {
	// draw base content for each row
	defBg := r.theme.DefaultBackground()
	defFg := r.theme.DefaultForeground()
	for viewY := int(r.buffer.ViewHeight() - 1); viewY >= 0; viewY-- {
		r.drawRow(viewY, defBg, defFg)
	}
}
