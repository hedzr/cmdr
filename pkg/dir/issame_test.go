package dir

import (
	"testing"
)

func TestIsSameVolume(t *testing.T) {
	tests := []struct {
		f1, f2 string
		want   bool
	}{
		{"/Volumes/VolHack/work/00.md", "/Volumes/VolHack/work/00.md", true},
		{"/Volumes/VolHack/work/00.md", "/Users/hz/Downloads/00.md", false},
		{"/Volumes/VolHack/work", "/Users/hz/Downloads", false},

		// for a split /home in linux,
		// {"/usr", "/Users/hz", false},
	}

	for _, tt := range tests {
		if !FileExists(tt.f1) || !FileExists(tt.f2) {
			continue
		}

		if got := IsSameVolume(tt.f1, tt.f2); got != tt.want {
			t.Errorf("IsSameVolume() = %v, want %v", got, tt.want)
		}
	}
}
