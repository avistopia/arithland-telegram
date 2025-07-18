package clean

import (
	"github.com/avistopia/arithland-telegram/internal/pkg/texts"
	"unicode/utf8"
)

// UserDisplayName validates and cleans the user display name field. Returns cleand value, and an error if not valid.
func UserDisplayName(displayName string) (string, string) {
	displayName = texts.TrimSpaces(displayName)

	if utf8.RuneCountInString(displayName) < 2 {
		return "", texts.DisplayNameIsTooShort
	}

	if utf8.RuneCountInString(displayName) > 30 {
		return "", texts.DisplayNameIsTooLong
	}

	return displayName, ""
}
