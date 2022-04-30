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

func TestStripTty(t *testing.T) {
	src := `[0m[90mscan folder and save [3mresult[0m[90m to [51;1mbgo.yml[0m[90m, as [7mproject settings[0m[90m[0m`
	tgt := `scan folder and save result to bgo.yml, as project settings`
	tt := StripEscapes(src)
	if tt != tgt {
		t.Fatal("wrong")
	}
}
