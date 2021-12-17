package hinters

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/liamg/darktile/internal/app/darktile/termutil"
)

func init() {
	register(&DmesgTimestampHinter{}, PriorityVeryLow)
	setSysStartTime()
}

var sysStart time.Time

type DmesgTimestampHinter struct{}

var dmsegTsMatcher = regexp.MustCompile(`^\[\s*\d+.\d{6}\]`)

func (h *DmesgTimestampHinter) Match(text string, cursorIndex int) (matched bool, offset int, length int) {
	matches := dmsegTsMatcher.FindAllStringIndex(text, -1)
	for _, match := range matches {
		if match[0] <= cursorIndex && match[1] > cursorIndex {
			return true, match[0], match[1] - match[0]
		}
	}
	return
}

func (h *DmesgTimestampHinter) Activate(api HintAPI, match string, start termutil.Position, end termutil.Position) error {
	match = strings.Split(strings.Trim(match, "[] "), ".")[0]
	seconds, err := strconv.ParseFloat(match, 32)
	if err != nil {
		return err
	}
	result := sysStart.Add(time.Duration(seconds) * time.Second).Format(time.ANSIC)
	api.Highlight(start, end, result, nil)
	return nil
}

func (h *DmesgTimestampHinter) Deactivate(api HintAPI) error {
	api.ClearHighlight()
	return nil
}

func (h *DmesgTimestampHinter) Click(api HintAPI) error {
	return nil
}

func setSysStartTime() {
	sysStart = time.Now().Local().Add(time.Duration(int(getUptime()*-1)) * time.Second)
}
