package stringsutil

import (
	"strings"
	"unicode"
)

func ToUpperCamelCase(in string) string {
	s := strings.Split(in, "_")
	var buf strings.Builder
	for _, v := range s {
		buf.WriteRune(unicode.ToUpper(rune(v[0])))
		buf.WriteString(v[1:])
	}
	return buf.String()
}

func ToLowerCamelCase(in string) string {
	s := ToUpperCamelCase(in)

	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}
