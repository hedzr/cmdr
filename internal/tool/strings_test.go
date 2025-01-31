package tool

import (
	"reflect"
	"testing"
)

func TestSplitCommandString(t *testing.T) {
	tests := []struct {
		input string
		out   []string
	}{
		{`bash -c "echo hello world!"`, []string{"bash", "-c", "echo hello world!"}},
		{`bash -c 'echo hello world!'`, []string{"bash", "-c", "echo hello world!"}},
	}

	for _, tt := range tests {
		got := SplitCommandString(tt.input, '\'')
		if !reflect.DeepEqual(got, tt.out) {
			t.Errorf("SplitCommandString(%v) = %v, want %v", tt.input, got, tt.out)
			// spew.SDump(got)
		}
	}
}

func TestStripQuotes(t *testing.T) {
	tests := []struct {
		input string
		out   string
	}{
		{`"echo hello world!"`, "echo hello world!"},
		{`'echo hello world!'`, "echo hello world!"},
	}

	for _, tt := range tests {
		got := StripQuotes(tt.input)
		got1 := TrimQuotes(tt.input)
		if !reflect.DeepEqual(got, tt.out) || got != got1 {
			t.Errorf("StripQuotes(%v) = %v, want %v", tt.input, got, tt.out)
			// spew.SDump(got)
		}
	}
}
