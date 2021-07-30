package hinters

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/liamg/darktile/internal/app/darktile/termutil"
)

func init() {
	register(&PermsHinter{}, PriorityLow)
}

type PermsHinter struct {
}

var permsMatcher = regexp.MustCompile("^[d|-][rwx-]{9}")

func (h *PermsHinter) Match(text string, cursorIndex int) (matched bool, offset int, length int) {
	matches := permsMatcher.FindAllStringIndex(text, -1)
	for _, match := range matches {
		if match[0] <= cursorIndex && match[1] > cursorIndex {
			return true, match[0], match[1] - match[0]
		}
	}
	return
}

func (h *PermsHinter) Activate(api HintAPI, match string, start termutil.Position, end termutil.Position) error {
	result, err := getPermsFromString(match)
	if err != nil {
		return err
	}
	api.Highlight(start, end, result, nil)
	return nil
}

func getPermsFromString(match string) (string, error) {
	match = strings.NewReplacer(
		"r", "1",
		"w", "1",
		"x", "1",
		"-", "0").
		Replace(match)

	return strings.Join([]string{
		"0",
		readPermPart(match, 1, 4),
		readPermPart(match, 4, 7),
		readPermPart(match, 7, 10),
	}, ""), nil
}

func readPermPart(match string, from, to int) string {
	permPart := match[from:to]

	i, err := strconv.ParseInt(permPart, 2, 0)
	if err != nil {
		return "0"
	}
	return fmt.Sprint(strconv.FormatInt(i, 8))
}

func (h *PermsHinter) Deactivate(api HintAPI) error {
	api.ClearHighlight()
	return nil
}

func (h *PermsHinter) Click(api HintAPI) error {
	return nil
}
