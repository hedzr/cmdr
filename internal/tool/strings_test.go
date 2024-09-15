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

func TestStringToBool(t *testing.T) {
	tests := []struct {
		input string
		out   bool
	}{
		{`0`, false},
		{`F`, false},
		{`f`, false},
		{`false`, false},
		{`FALSE`, false},
		{`OFF`, false},
		{`Off`, false},
		{`off`, false},
		{`N`, false},
		{`n`, false},
		{`NO`, false},
		{`No`, false},
		{`no`, false},

		{`1`, true},
		{`t`, true},
		{`true`, true},
		{`T`, true},
		{`TRUE`, true},
		{`True`, true},
		{`Y`, true},
		{`y`, true},
		{`YES`, true},
		{`Yes`, true},
		{`yes`, true},
		{`ON`, true},
		{`On`, true},
		{`on`, true},
	}

	for _, tt := range tests {
		got := StringToBool(tt.input)
		if !reflect.DeepEqual(got, tt.out) {
			t.Errorf("StringToBool(%v) = %v, want %v", tt.input, got, tt.out)
			// spew.SDump(got)
		}
	}
}

func TestToBool(t *testing.T) {
	tests := []struct {
		input any
		out   bool
	}{
		{int8(0), false},
		{int16(0), false},
		{int32(0), false},
		{int64(0), false},
		{int(0), false},
		{uint8(0), false},
		{uint16(0), false},
		{uint32(0), false},
		{uint64(0), false},
		{uint(0), false},
		{int8(1), true},
		{int16(1), true},
		{int32(1), true},
		{int64(1), true},
		{int(1), true},
		{uint8(1), true},
		{uint16(1), true},
		{uint32(1), true},
		{uint64(1), true},
		{uint(1), true},

		{false, false},
		{true, true},

		{`0`, false},
		{`F`, false},
		{`f`, false},
		{`false`, false},
		{`FALSE`, false},
		{`OFF`, false},
		{`Off`, false},
		{`off`, false},
		{`N`, false},
		{`n`, false},
		{`NO`, false},
		{`No`, false},
		{`no`, false},

		{`1`, true},
		{`t`, true},
		{`true`, true},
		{`T`, true},
		{`TRUE`, true},
		{`True`, true},
		{`Y`, true},
		{`y`, true},
		{`YES`, true},
		{`Yes`, true},
		{`yes`, true},
		{`ON`, true},
		{`On`, true},
		{`on`, true},
	}

	for _, tt := range tests {
		got := ToBool(tt.input)
		if !reflect.DeepEqual(got, tt.out) {
			t.Errorf("ToBool(%v) = %v, want %v", tt.input, got, tt.out)
			// spew.SDump(got)
		}
	}
}
