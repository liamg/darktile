package gui

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"time"
)

func (g *GUI) RequestScreenshot(filename string) {
	g.screenshotRequested = true
	if filename == "" {
		filename = fmt.Sprintf("darktile-screenshot-%d.png", time.Now().UnixNano())
		targetdir, err := os.UserHomeDir()
		if err != nil {
			targetdir = "/tmp"
		}
		filename = filepath.Join(targetdir, filename)
	}
	g.screenshotFilename = filename
}

func (g *GUI) takeScreenshot(screen image.Image) {
	g.screenshotRequested = false

	file, err := os.Create(g.screenshotFilename)
	if err != nil {
		g.ShowError(fmt.Sprintf("Screenshot failed: %s", err))
		return
	}
	defer file.Close()

	if err := png.Encode(file, screen); err != nil {
		g.ShowError(fmt.Sprintf("Screenshot failed: %s", err))
		return
	}

	g.ShowMessage(fmt.Sprintf("Screenshot saved: %s", g.screenshotFilename))
}
