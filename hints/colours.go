package hints

import (
	"regexp"
	"strconv"
)

func init() {
	hinters = append(hinters, hintColours)
}

func hintColours(word string, context string, wordX uint16, wordY uint16) *Hint {
	item := NewHint(word, context, wordX, wordY)

	if isColour(word) {

		r, err := strconv.ParseInt(word[1:3], 16, 64)
		if err != nil {
			return nil
		}
		g, err := strconv.ParseInt(word[3:5], 16, 64)
		if err != nil {
			return nil
		}
		b, err := strconv.ParseInt(word[5:7], 16, 64)
		if err != nil {
			return nil
		}

		item.Description = word
		item.BackgroundColour = [3]float32{float32(r) / 255, float32(g) / 255, float32(b) / 255}
		if (r+g+b)/3 < 128 {
			item.ForegroundColour = [3]float32{1, 1, 1}
		} else {
			item.ForegroundColour = [3]float32{0, 0, 0}
		}
		return item
	}

	return nil
}

func isColour(s string) bool {
	re := regexp.MustCompile("#[0-9A-Fa-f]{6}")
	return re.MatchString(s)
}
