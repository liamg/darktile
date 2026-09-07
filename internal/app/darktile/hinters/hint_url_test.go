package hinters

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// The hinter only underlines what the cursor is actually over, so a match is a
// function of both the line and the column. xurls decides where a URL ends,
// which is the part that moves when it is upgraded.
func Test_url_hinter_matches(t *testing.T) {

	tests := []struct {
		name string
		text string
		want string
	}{
		{"https", "see https://example.com for more", "https://example.com"},
		{"http", "see http://example.com for more", "http://example.com"},
		{"with path and query", "go to https://example.com/a/b?c=1&d=2 now", "https://example.com/a/b?c=1&d=2"},
		{"with port", "serving on http://localhost:8080/ok", "http://localhost:8080/ok"},
		{"trailing full stop is excluded", "read https://example.com/page.", "https://example.com/page"},
		{"inside parentheses", "(https://example.com)", "https://example.com"},
		{"mailto", "mail me at mailto:a@example.com ok", "mailto:a@example.com"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hinter := &URLHinter{}
			index := strings.Index(test.text, test.want)
			assert.GreaterOrEqual(t, index, 0, "test text must contain the expected match")

			matched, offset, length := hinter.Match(test.text, index)

			assert.True(t, matched)
			assert.Equal(t, test.want, test.text[offset:offset+length])
		})
	}
}

// Strict mode is what keeps a hostname in ordinary output from turning the
// whole line into a link.
func Test_url_hinter_ignores_text_without_a_scheme(t *testing.T) {

	tests := []struct{ name, text string }{
		{"bare domain", "example.com is a domain"},
		{"prose", "this is just a sentence"},
		{"file path", "/usr/share/doc/readme"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hinter := &URLHinter{}
			matched, _, _ := hinter.Match(test.text, 2)
			assert.False(t, matched)
		})
	}
}

// A URL elsewhere on the line is not the one under the cursor.
func Test_url_hinter_does_not_match_away_from_the_cursor(t *testing.T) {

	hinter := &URLHinter{}
	text := "https://example.com and then some trailing prose"

	matched, _, _ := hinter.Match(text, len(text)-4)

	assert.False(t, matched)
}
