package stringsutil

import (
	"strings"
	"unicode"

	"github.com/gertd/go-pluralize"
)

var word = []string{"UUID", "UID", "WWIDs", "IQN", "ISCSI", "API", "ACME", "CHAP", "CIDRs", "CIDR", "PID", "ID", "DNS",
	"IPC", "IPs", "IP", "PUT", "POST", "QOS", "OS", "NFS", "FS", "FC", "RBD", "TCP", "UDP", "SCTP", "GET", "URI", "URL",
	"TLS", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH", "HTTPS", "HTTP", "SELinux", "FQDN", "TTY", "WWNs", "GCE",
	"AWS", "IO", "CSI", "GRPC", "SSL", "GMSA", "HMAC", "HEAD"}

var wordDic map[string]struct{}

func init() {
	wordDic = make(map[string]struct{})
	for _, v := range word {
		wordDic[v] = struct{}{}
	}
}

func ToUpperCamelCase(in string) string {
	var s []string
	var start, cur int
	for cur < len(in) {
		switch in[cur] {
		case '_', '-', '.':
			s = append(s, in[start:cur])
			start = cur + 1
		}
		cur++
	}
	s = append(s, in[start:])

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
	s := splitString(in)
	for i := range s {
		s[i] = strings.ToUpper(s[i])
	}
	return strings.Join(s, "_")
}

func ToLowerSnakeCase(in string) string {
	s := splitString(in)
	for i := range s {
		s[i] = strings.ToLower(s[i])
	}
	return strings.Join(s, "_")
}

func splitString(in string) []string {
	s := []string{in}
	i := 0
Loop:
	for i < len(s) {
		str := s[i]

		for _, v := range word {
			if v == str {
				i++
				continue Loop
			}

			if idx := strings.Index(str, v); idx != -1 {
				if idx == 0 {
					s[i] = str[idx : idx+len(v)]
					i++
				} else {
					s = insert(s, i+1, str[idx:idx+len(v)])
					s[i] = str[:idx]
				}
				if idx+len(v) < len(str) {
					s = append(s, str[idx+len(v):])
				}
				continue Loop
			}
		}

		for k := 0; k < len(s[i]); k++ {
			if k >= 1 && unicode.IsUpper(rune(s[i][k])) {
				s = insert(s, i, s[i][:k])
				s[i+1] = s[i+1][k:]
				i++
				continue Loop
			}

			switch s[i][k] {
			case '-', '_':
				if k == 0 {
					if len(s[i]) == 1 {
						s = remove(s, i)
					} else {
						s[i] = s[i][k+1:]
					}
				} else {
					s = insert(s, i, s[i][:k])
					s[i+1] = s[i+1][k+1:]
					i++
				}
				continue Loop
			}
		}
		i++
	}

	return s
}

func insert(s []string, i int, v string) []string {
	if len(s) == i {
		return append(s, v)
	}
	s = append(s[:i+1], s[i:]...)
	s[i] = v
	return s
}

func remove(s []string, i int) []string {
	return append(s[:i], s[i+1:]...)
}

var pluralizeClient = pluralize.NewClient()

func Plural(word string) string {
	return pluralizeClient.Plural(word)
}

func Singular(word string) string {
	return pluralizeClient.Singular(word)
}
