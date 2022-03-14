package stringtool

import "strings"

// ConvertFirstUpper method
func ConvertFirstUpper(str string) string {
	return strings.ToUpper(str[0:1]) + str[1:]
}

// ConvertFirstLower method
func ConvertFirstLower(str string) string {
	return strings.ToLower(str[0:1]) + str[1:]
}

// StringJoin method
func StringJoin(str ...string) string {
	s := make([]string, 0)
	s = append(s, str...)
	return strings.Join(s, "")
}
