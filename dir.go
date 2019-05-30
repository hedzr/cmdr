/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// GetExcutableDir returns the executable file directory
func GetExcutableDir() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	// fmt.Println(dir)
	return dir
}

// GetExcutablePath returns the executable file path
func GetExcutablePath() string {
	p, _ := filepath.Abs(os.Args[0])
	return p
}

// GetCurrentDir returns the current workingFlag directory
func GetCurrentDir() string {
	dir, _ := os.Getwd()
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	// fmt.Println(dir)
	return dir
}

// FileExists returns the existence of an directory or file
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// EnsureDir checks and creates the directory.
func EnsureDir(dir string) (err error) {
	if len(dir) == 0 {
		return fmt.Errorf("empty directory")
	}
	if !FileExists(dir) {
		err = os.MkdirAll(dir, 0755)
	}
	return
}

// NormalizeDir make dir name normalized
func NormalizeDir(s string) string {
	return normalizeDir(s)
}

func normalizeDir(s string) string {
	s = os.Expand(s, os.Getenv)
	if s[0] == '/' {
		return s
	} else if strings.HasPrefix(s, "./") {
		return path.Join(GetCurrentDir(), s)
	} else if strings.HasPrefix(s, "../") {
		return path.Join(GetCurrentDir(), s)
	} else if strings.HasPrefix(s, "~/") {
		return path.Join(os.Getenv("HOME"), s[2:])
	} else {
		return s
	}
}

// getExpandedPredefinedLocations for internal using
func getExpandedPredefinedLocations() (locations []string) {
	for _, d := range predefinedLocations {
		locations = append(locations, normalizeDir(d))
	}
	return
}

// GetPredefinedLocations return the searching locations for loading config files.
func GetPredefinedLocations() []string {
	return predefinedLocations
}

// SetPredefinedLocations to customize the searching locations for loading config files.
// It MUST be invoked before `cmdr.Exec`. Such as:
// ```go
//     SetPredefinedLocations([]string{"./config", "~/.config/cmdr/", "$GOPATH/running-configs/cmdr"})
// ```
func SetPredefinedLocations(locations []string) {
	predefinedLocations = locations
}
