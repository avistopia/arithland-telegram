package texts

import (
	"strings"
	"unicode"
)

func TrimSpaces(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		// Check for zero-width characters and whitespace
		return unicode.IsSpace(r) || r == '\u200b' || r == '\u200c' || r == '\u200d'
	})
}
