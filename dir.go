/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"os"
	"path/filepath"
)

func GetExcutableDir() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	// fmt.Println(dir)
	return dir
}

func GetCurrentDir() string {
	dir, _ := os.Getwd()
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	// fmt.Println(dir)
	return dir
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func EnsureDir(dir string) (err error) {
	if !FileExists(dir) {
		err = os.MkdirAll(dir, 0755)
	}
	return
}
