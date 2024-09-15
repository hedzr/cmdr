package dir

import (
	"os"
	"syscall"
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
func IsSameVolume(f1, f2 string) bool {
	stat1, _ := os.Stat(f1)
	stat2, _ := os.Stat(f2)

	// // returns *syscall.Stat_t
	// fmt.Println(reflect.TypeOf(stat1.Sys()))
	//
	// fmt.Println(stat1.Sys().(*syscall.Stat_t).Dev)
	// fmt.Println(stat2.Sys().(*syscall.Stat_t).Dev)

	return stat1.Sys().(*syscall.Stat_t).Dev == stat2.Sys().(*syscall.Stat_t).Dev
}
