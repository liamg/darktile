package gui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liamg/darktile/internal/app/darktile/termutil"
)

type WindowManipulator struct {
	g          *GUI
	title      string
	titleStack []string
}

func NewManipulator(g *GUI) *WindowManipulator {
	return &WindowManipulator{
		g: g,
	}
}

func (m *WindowManipulator) ReportError(err error) {
	m.g.ShowError(err.Error())
}

func (m *WindowManipulator) CellSizeInPixels() (int, int) {
	size := m.g.fontManager.CharSize()
	return size.X, size.Y
}

func (m *WindowManipulator) Position() (int, int) {
	return ebiten.WindowPosition()
}

func (m *WindowManipulator) GetTitle() string {
	return m.title
}

func (m *WindowManipulator) SetTitle(title string) {
	m.title = title
	ebiten.SetWindowTitle(m.title)
}

func (m *WindowManipulator) SaveTitleToStack() {
	m.titleStack = append(m.titleStack, m.title)
}

func (m *WindowManipulator) RestoreTitleFromStack() {
	if len(m.titleStack) == 0 {
		m.SetTitle("")
	}

	title := m.titleStack[len(m.titleStack)-1]
	m.titleStack = m.titleStack[:len(m.titleStack)-1]
	m.SetTitle(title)
}

func (m *WindowManipulator) State() termutil.WindowState {
	if ebiten.IsWindowMinimized() {
		return termutil.StateMinimised
	}

	if ebiten.IsWindowMaximized() {
		return termutil.StateMaximised
	}

	return termutil.StateNormal
}

func (m *WindowManipulator) Minimise() {
	ebiten.MinimizeWindow()
}

func (m *WindowManipulator) Maximise() {
	ebiten.MaximizeWindow()
}

func (m *WindowManipulator) Restore() {
	ebiten.RestoreWindow()
}

func (m *WindowManipulator) SizeInPixels() (int, int) {
	return m.g.size.X, m.g.size.Y
}

func (m *WindowManipulator) SizeInChars() (int, int) {
	return int(m.g.terminal.GetActiveBuffer().ViewWidth()), int(m.g.terminal.GetActiveBuffer().ViewHeight())
}

func (m *WindowManipulator) ResizeInPixels(x int, y int) {
	ebiten.SetWindowSize(x, y)
}

func (m *WindowManipulator) ResizeInChars(cols int, rows int) {
	x := cols * m.g.fontManager.CharSize().X
	y := rows * m.g.fontManager.CharSize().Y
	ebiten.SetWindowSize(x, y)
}

func (m *WindowManipulator) ScreenSizeInPixels() (int, int) {
	return ebiten.WindowSize()
}

func (m *WindowManipulator) ScreenSizeInChars() (int, int) {
	w, h := ebiten.WindowSize()
	return w / m.g.fontManager.CharSize().X, h / m.g.fontManager.CharSize().Y
}

func (m *WindowManipulator) Move(x, y int) {
	ebiten.SetWindowPosition(x, y)
}

func (m *WindowManipulator) SetFullscreen(enabled bool) {
	ebiten.SetFullscreen(enabled)
}

func (m *WindowManipulator) IsFullscreen() bool {
	return ebiten.IsFullscreen()
}
