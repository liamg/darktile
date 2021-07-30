package fontinfo

import "strings"

type matcher func(m *fontMetadata) bool

// MatchFamily is a matcher which matches fonts with the specified font family (case insensitive)
func MatchFamily(family string) matcher {
	return func(m *fontMetadata) bool {
		return strings.EqualFold(m.FontFamily, family)
	}
}

// MatchStyle is a matcher which matches fonts with the specified font family (case insensitive)
func MatchStyle(style string) matcher {
	return func(m *fontMetadata) bool {
		return strings.EqualFold(m.FontStyle, style)
	}
}
