package stringsutil

import (
	"strings"
	"unicode"

	"github.com/gertd/go-pluralize"
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

var pluralizeClient = pluralize.NewClient()

func Plural(word string) string {
	return pluralizeClient.Plural(word)
}

func Singular(word string) string {
	return pluralizeClient.Singular(word)
}
