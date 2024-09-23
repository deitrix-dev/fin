package stringsx

import "strings"

// UpperFirst returns the string with the first character uppercased.
func UpperFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}
