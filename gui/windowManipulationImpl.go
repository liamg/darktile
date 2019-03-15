package gui

import (
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/liamg/aminal/terminal"
)

//
// Implementation of the terminal.WindowManipulationInterface
//

func (gui *GUI) RestoreWindow(term *terminal.Terminal) error {
	return gui.executeInMainThread(func() error {
		gui.window.Restore()
		return nil
	})
}

func (gui *GUI) IconifyWindow(term *terminal.Terminal) error {
	return gui.executeInMainThread(func() error {
		return gui.window.Iconify()
	})
}

func (gui *GUI) MoveWindow(term *terminal.Terminal, pixelX int, pixelY int) error {
	return gui.executeInMainThread(func() error {
		gui.window.SetPos(pixelX, pixelY)
		return nil
	})
}

func (gui *GUI) ResizeWindowByPixels(term *terminal.Terminal, pixelsHeight int, pixelsWidth int) error {
	return gui.executeInMainThread(func() error {
		term.Unlock()
		gui.window.SetSize(pixelsWidth, pixelsHeight)
		term.Lock()
		return nil
	})
}

func (gui *GUI) BringWindowToFront(term *terminal.Terminal) error {
	var err error
	if gui.window.GetAttrib(glfw.Iconified) != 0 {
		err = gui.executeInMainThread(func() error {
			return gui.window.Restore()
		})
	}

	if err != nil {
		err = gui.window.Focus()
	}

	return err
}

func (gui *GUI) ResizeWindowByChars(term *terminal.Terminal, charsHeight int, charsWidth int) error {
	return gui.executeInMainThread(func() error {
		return term.SetSize(uint(charsWidth), uint(charsHeight))
	})
}

func (gui *GUI) MaximizeWindow(term *terminal.Terminal) error {
	return gui.executeInMainThread(func() error {
		term.Lock()
		err := gui.window.Maximize()
		term.Unlock()
		return err
	})
}

func (gui *GUI) ReportWindowState(term *terminal.Terminal) error {
	// Report xterm window state. If the xterm window is open (non-iconified), it returns CSI 1 t .
	// If the xterm window is iconified, it returns CSI 2 t .
	if gui.window.GetAttrib(glfw.Iconified) != 0 {
		_ = term.Write([]byte("\x1b[2t"))
	} else {
		_ = term.Write([]byte("\x1b[1t"))
	}

	return nil
}

func (gui *GUI) ReportWindowPosition(term *terminal.Terminal) error {
	// Report xterm window position as CSI 3 ; x; yt
	x, y := gui.window.GetPos()

	_ = term.Write([]byte(fmt.Sprintf("\x1b[3;%d;%dt", x, y)))

	return nil
}

func (gui *GUI) ReportWindowSizeInPixels(term *terminal.Terminal) error {
	// Report xterm window in pixels as CSI 4 ; height ; width t
	_ = term.Write([]byte(fmt.Sprintf("\x1b[4;%d;%dt", gui.height, gui.width)))

	return nil
}

func (gui *GUI) ReportWindowSizeInChars(term *terminal.Terminal) error {
	// Report the size of the text area in characters as CSI 8 ; height ; width t
	charsWidth, charsHeight := gui.renderer.GetTermSize()

	_ = term.Write([]byte(fmt.Sprintf("\x1b[8;%d;%dt", charsHeight, charsWidth)))

	return nil
}
