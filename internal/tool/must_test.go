package tool

import (
	"testing"
)

func TestMust(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"123", 123},
	}
	for _, tt := range tests {
		got := Must[int64](tt.input)
		// if err != nil {
		// 	t.Fatal(err)
		// }
		if got != tt.want {
			t.Errorf("want %v, got %v", tt.want, got)
		}
	}

	tests1 := []struct {
		input string
		want  uint64
	}{
		{"123", uint64(123)},
	}
	for _, tt := range tests1 {
		got := Must[uint64](tt.input)
		// if err != nil {
		// 	t.Fatal(err)
		// }
		if got != tt.want {
			t.Errorf("want %v, got %v", tt.want, got)
		}
	}

	tests2 := []struct {
		input string
		want  float64
	}{
		{"123", 123},
		{"123.4567890123", 123.4567890123},
	}
	for _, tt := range tests2 {
		got := Must[float64](tt.input)
		// if err != nil {
		// 	t.Fatal(err)
		// }
		if got != tt.want {
			t.Errorf("want %v, got %v", tt.want, got)
		}
	}
}
