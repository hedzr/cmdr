package tool

import (
	"testing"
)

func TestParseT(t *testing.T) {
	tests := []struct {
		input string
		want  any
	}{
		{"123", 123},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var got int
			var err error
			got, err = ParseT[int](tt.input)
			if err != nil {
				t.Failed()
			}
			if got != tt.want {
				t.Errorf("ParseT() = %v, want %v", got, tt.want)
			}
		})
	}
}
