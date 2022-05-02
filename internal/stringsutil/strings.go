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
		buf.WriteString(strings.ToLower(v[1:]))
	}
	return buf.String()
}

func ToLowerCamelCase(in string) string {
	s := ToUpperCamelCase(in)

	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}

func ToUpperSnakeCase(in string) string {
	var buf strings.Builder
	for i, v := range in {
		if i != 0 && unicode.IsUpper(v) {
			buf.WriteRune('_')
		}
		buf.WriteRune(unicode.ToUpper(v))
	}
	return buf.String()
}

func ToLowerSnakeCase(in string) string {
	var buf strings.Builder
	for i, v := range in {
		if i != 0 && unicode.IsUpper(v) {
			buf.WriteRune('_')
		}
		buf.WriteRune(unicode.ToLower(v))
	}
	return buf.String()
}

var pluralizeClient = pluralize.NewClient()

func Plural(word string) string {
	return pluralizeClient.Plural(word)
}

func Singular(word string) string {
	return pluralizeClient.Singular(word)
}
