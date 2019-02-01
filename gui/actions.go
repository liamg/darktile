package gui

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/liamg/aminal/config"
)

var actionMap = map[config.UserAction]func(gui *GUI){
	config.ActionCopy:        actionCopy,
	config.ActionPaste:       actionPaste,
	config.ActionToggleDebug: actionToggleDebug,
	config.ActionSearch:      actionSearchSelection,
	config.ActionToggleSlomo: actionToggleSlomo,
	config.ActionReportBug:   actionReportBug,
}

func actionCopy(gui *GUI) {
	selectedText := gui.terminal.ActiveBuffer().GetSelectedText()

	if selectedText != "" {
		gui.window.SetClipboardString(selectedText)
	}
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

func actionSearchSelection(gui *GUI) {
	keywords := gui.terminal.ActiveBuffer().GetSelectedText()
	if keywords != "" && gui.config.SearchURL != "" && strings.Contains(gui.config.SearchURL, "$QUERY") {
		gui.launchTarget(fmt.Sprintf(strings.Replace(gui.config.SearchURL, "$QUERY", "%s", 1), url.QueryEscape(keywords)))
	}
}

func actionToggleSlomo(gui *GUI) {
	gui.config.Slomo = !gui.config.Slomo
}

func actionReportBug(gui *GUI) {
	gui.launchTarget("https://github.com/liamg/aminal/issues/new/choose")
}
