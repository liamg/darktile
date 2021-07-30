package hinters

import (
	"github.com/liamg/darktile/internal/app/darktile/termutil"
	"github.com/skratchdot/open-golang/open"
	"mvdan.cc/xurls"
)

func init() {
	register(&URLHinter{}, PriorityHigh)
}

type URLHinter struct {
	target string
}

func (h *URLHinter) Match(text string, cursorIndex int) (matched bool, offset int, length int) {
	matches := xurls.Strict.FindAllStringIndex(text, -1)
	for _, match := range matches {
		if match[0] <= cursorIndex && match[1] > cursorIndex {
			return true, match[0], match[1] - match[0]
		}
	}
	return
}

func (h *URLHinter) Activate(api HintAPI, match string, start termutil.Position, end termutil.Position) error {
	h.target = match
	api.Highlight(start, end, "CTRL + CLICK: Open in browser", nil)
	api.SetCursorToPointer()
	return nil
}

func (h *URLHinter) Deactivate(api HintAPI) error {
	api.ClearHighlight()
	api.ResetCursor()
	return nil
}

func (h *URLHinter) Click(api HintAPI) error {
	api.ShowMessage("Launching URL in your browser...")
	return open.Run(h.target)
}
