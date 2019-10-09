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

// GetExecutableDir returns the executable file directory
func GetExecutableDir() string {
	// _ = ioutil.WriteFile("/tmp/11", []byte(strings.Join(os.Args,",")), 0644)
	// fmt.Printf("os.Args[0] = %v\n", os.Args[0])

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
// it should be equal with os.Getenv("PWD")
func GetCurrentDir() string {
	dir, _ := os.Getwd()
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	// fmt.Println(dir)
	return dir
}

// IsDirectory tests whether `path` is a directory or not
func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

// IsRegularFile tests whether `path` is a normal regular file or not
func IsRegularFile(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.Mode().IsRegular(), err
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
	if len(s) == 0 {
		return s
	}

	s = os.Expand(s, os.Getenv)
	if s[0] == '/' {
		return s
	} else if strings.HasPrefix(s, "./") {
		return path.Join(GetCurrentDir(), s)
	} else if strings.HasPrefix(s, "../") {
		return path.Dir(path.Join(GetCurrentDir(), s))
	} else if strings.HasPrefix(s, "~/") {
		return path.Join(os.Getenv("HOME"), s[2:])
	} else {
		return s
	}
}

// GetPredefinedLocations return the searching locations for loading config files.
func GetPredefinedLocations() []string {
	return uniqueWorker.predefinedLocations
}

// SetPredefinedLocations to customize the searching locations for loading config files.
//
// It MUST be invoked before `cmdr.Exec`. Such as:
// ```go
//     SetPredefinedLocations([]string{"./config", "~/.config/cmdr/", "$GOPATH/running-configs/cmdr"})
// ```
func SetPredefinedLocations(locations []string) {
	uniqueWorker.predefinedLocations = locations
}
