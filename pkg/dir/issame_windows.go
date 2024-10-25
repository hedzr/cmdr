package dir

import (
	"os"
	"strings"
)

// IsSameVolume detects if two files are in same volume (Unix like)
// or partition (Windows like).
//
// For windows, c:\1.txt & d:\2.txt shall be in different volumes.
// For Linux/Unix/macOS, suppose you have a split `/home` than `/`,
// then /usr and /home/hz should be in different volumes.
//
// On same a volume, moving one file as another one is quick and
// without file content copying; but on the contrary, it needs
// to copy source to destine and deleting source file to take
// affects.
//
// TODO:For Windows, this function is not tested yet.
func IsSameVolume(f1, f2 string) bool {
	stat1, _ := os.Stat(f1)
	stat2, _ := os.Stat(f2)

	// todo

	n1 := stat1.Name()
	n2 := stat2.Name()
	if strings.Contains(n1, ":") {
		pos := strings.Index(n1, ":")
		n1 = n1[:pos+1]
	}
	if strings.Contains(n2, ":") {
		pos := strings.Index(n2, ":")
		n2 = n2[:pos+1]
	}
	return n1 == n2
}
