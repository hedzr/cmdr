// Copyright © 2020 Hedzr Yeh.

package tool

import "testing"

func TestStripQuotes(t *testing.T) {

	for _, tc := range []struct {
		src      string
		expected string
	}{
		{"'8.yes'", "8.yes"},
		{"'8.yes", "8.yes"},
		{"8.yes'", "8.yes"},
		{"'8.y'es'", "8.y'es"},
		{"\"8.yes\"", "8.yes"},
		{"\"8.yes", "8.yes"},
		{"8.yes\"", "8.yes"},
		{"\"8.y'es\"", "8.y'es"},
	} {
		tgt := StripQuotes(tc.src)
		if tgt != tc.expected {
			t.Fatalf("StripQuotes(%q): expect %q, but got %q", tc.src, tc.expected, tgt)
		}
	}

}
