package hinters

import (
	"encoding/base64"
	"regexp"

	"github.com/liamg/darktile/internal/app/darktile/termutil"
)

func init() {
	register(&Base64Hinter{}, PriorityVeryLow)
}

type Base64Hinter struct {
	target string
}

var base64Matcher = regexp.MustCompile("(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{4})")

func (h *Base64Hinter) Match(text string, cursorIndex int) (matched bool, offset int, length int) {
	matches := base64Matcher.FindAllStringIndex(text, -1)
	for _, match := range matches {
		if match[0] <= cursorIndex && match[1] > cursorIndex {
			result := text[match[0]:match[1]]
			if len(result) > 4 && isReadable(result) {
				return true, match[0], match[1] - match[0]
			}
		}
	}
	return
}

func isReadable(result string) bool {
	parts, err := base64.StdEncoding.DecodeString(result)
	if err != nil {
		return false
	}
	for i := range parts {
		if (parts[i] > 0x7e || parts[i] < 0x20) && parts[i] != 0x0a && parts[i] != 0x0d {
			return false
		}
	}
	return true
}

func (h *Base64Hinter) Activate(api HintAPI, match string, start termutil.Position, end termutil.Position) error {
	h.target = match
	result, err := base64.StdEncoding.DecodeString(match)
	if err != nil {
		return err
	}
	api.Highlight(start, end, "Base64 decodes to:\n"+string(result), nil)
	return nil
}

func (h *Base64Hinter) Deactivate(api HintAPI) error {
	api.ClearHighlight()
	return nil
}

func (h *Base64Hinter) Click(api HintAPI) error {
	return nil
}
