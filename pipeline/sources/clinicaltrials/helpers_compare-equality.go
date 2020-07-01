package main

import "strings"

// isEqual will determine if to entries are equal
func isEqual(value1, value2 string) (fullEqual, partialEqual bool) {

	if value1 == value2 {
		fullEqual = true
	}

	if clean(value1) == clean(value2) {
		partialEqual = true
	}

	return fullEqual, partialEqual

}

func clean(value string) string {

	s := []byte(value)

	j := 0
	for _, b := range s {
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') {
			s[j] = b
			j++
		}
	}
	return strings.ToLower(string(s[:j]))
}
