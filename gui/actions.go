package gui

import (
	"fmt"
	"net/url"

	"github.com/liamg/aminal/config"
)

var actionMap = map[config.UserAction]func(gui *GUI){
	config.ActionCopy:        actionCopy,
	config.ActionPaste:       actionPaste,
	config.ActionToggleDebug: actionToggleDebug,
	config.ActionGoogle:      actionGoogleSelection,
	config.ActionToggleSlomo: actionToggleSlomo,
	config.ActionReportBug:   actionReportBug,
}

func actionCopy(gui *GUI) {
	gui.window.SetClipboardString(gui.terminal.ActiveBuffer().GetSelectedText())
}

func actionPaste(gui *GUI) {
	if s, err := gui.window.GetClipboardString(); err == nil {
		_ = gui.terminal.Paste([]byte(s))
	}
}

func actionToggleDebug(gui *GUI) {
	gui.showDebugInfo = !gui.showDebugInfo
	gui.terminal.SetDirty()
}

func actionGoogleSelection(gui *GUI) {
	keywords := gui.terminal.ActiveBuffer().GetSelectedText()
	if keywords != "" {
		gui.launchTarget(fmt.Sprintf("https://www.google.com/search?q=%s", url.QueryEscape(keywords)))
	}
}

func actionToggleSlomo(gui *GUI) {
	gui.config.Slomo = !gui.config.Slomo
}

func actionReportBug(gui *GUI) {
	gui.launchTarget("https://github.com/liamg/aminal/issues/new/choose")
}
