package yasqlutil

import "strings"

const (
	DOUBLE_QUOTE = `"`
)

func HasQuoted(s string) bool {
	return strings.HasPrefix(s, DOUBLE_QUOTE) &&
		strings.HasSuffix(s, DOUBLE_QUOTE)
}

func Quote(s string, force ...bool) string {
	if len(force) == 0 && HasQuoted(s) {
		return s
	}
	return DOUBLE_QUOTE + s + DOUBLE_QUOTE
}
