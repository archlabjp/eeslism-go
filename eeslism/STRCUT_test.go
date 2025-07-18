
package eeslism

import "testing"

func TestSTRCUT(t *testing.T) {
	testCases := []struct {
		name     string
		data     string
		a        string
		expected string
	}{
		{"simple case", "abcde", "c", "ab"},
		{"substring at the end", "abcde", "e", "abcd"},
		{"substring at the beginning", "abcde", "a", ""},
		{"substring not found", "abcde", "f", ""},
		{"multiple occurrences", "abacada", "a", "abacad"},
		{"empty data", "", "a", ""},
		{"empty substring", "abcde", "", "abcde"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := STRCUT(tc.data, tc.a)
			if result != tc.expected {
				t.Errorf("Expected '%s', but got '%s'", tc.expected, result)
			}
		})
	}
}
