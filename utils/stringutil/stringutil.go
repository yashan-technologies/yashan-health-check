// The stringutil package encapsulates functions related to strings.
package stringutil

import (
	"encoding/base64"
	"regexp"
)

const (
	STR_EMPTY         = ""
	STR_BLANK_SPACE   = " "
	STR_NEWLINE       = "\n"
	STR_COMMA         = ","
	STR_DOT           = "."
	STR_HYPHEN        = "-"
	STR_BAR           = "|"
	STR_FORWARD_SLASH = "/"
	STR_UNDER_SCORE   = "_"
	STR_HASH          = "#"
	STR_HTML_BR       = "<br>"
	STR_QUESTION_MARK = "?"
	STR_EQUAL_SIGN    = "="
	STR_SINGEL_QUOTE  = "\""
	STR_COLON         = ":"
)

// IsEmpty checks whether a string is empty.
func IsEmpty(str string) bool {
	return len(str) == 0
}

func RemoveExtraSpaces(str string) string {
	regex := regexp.MustCompile(`\s+`)
	return regex.ReplaceAllString(str, STR_BLANK_SPACE)
}

func IsBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}
