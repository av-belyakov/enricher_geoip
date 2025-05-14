package supportingfunctions

import "strings"

// ReplaceCommaCharacter заменяет двойную кавычку одинарной
func ReplaceCommaCharacter(v string) string {
	v = strings.ReplaceAll(v, "\"", "'")
	v = strings.ReplaceAll(v, "\n", "")
	v = strings.ReplaceAll(v, "\\", "\\\\")
	return v
}
