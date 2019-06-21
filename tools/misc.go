package tools

import "strings"

var (
	Verbose = false
)

func IsEmpty(s string) bool {
	return len(s) == 0
}

func FilterEmptyLines(s string) string {
	r := []string{}

	for _, st := range strings.Split(s, "\n") {
		if !IsEmpty(st) {
			r = append(r, st)
		}
	}
	return strings.Join(r, "\n")
}
