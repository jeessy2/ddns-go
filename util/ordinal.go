package util

import (
	"strconv"

	"golang.org/x/text/language"
)

// Ordinal returns the ordinal format of the given number.
//
// See also: https://github.com/dustin/go-humanize/blob/master/ordinals.go
func Ordinal(x int, lang string) string {
	s := strconv.Itoa(x)

	// Chinese doesn't require an ordinal
	if lang == language.Chinese.String() {
		return s
	}

	suffix := "th"
	switch x % 10 {
	case 1:
		if x%100 != 11 {
			suffix = "st"
		}
	case 2:
		if x%100 != 12 {
			suffix = "nd"
		}
	case 3:
		if x%100 != 13 {
			suffix = "rd"
		}
	}
	return s + suffix
}
