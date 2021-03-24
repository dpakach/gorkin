package utils

import (
	"strings"
)

// AreArrayEqual returns if given two array of strings are equal
func AreArrayEqual(a, b []string) bool {
	if len(a) == 0 {
		if len(b) == 0 {
			return true
		}
		return false
	}
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// AreMapEqual returns if given two maps of strings are equal
func AreMapEqual(a, b map[string]string) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ReplaceNth Replace the nth occurrence of old in s by new.
func ReplaceNth(s, old, new string, n int) string {
	i := 0
	for m := 1; m <= n; m++ {
		x := strings.Index(s[i:], old)
		if x < 0 {
			break
		}
		i += x
		if m == n {
			return s[:i] + new + s[i+len(old):]
		}
		i += len(old)
	}
	return s
}
