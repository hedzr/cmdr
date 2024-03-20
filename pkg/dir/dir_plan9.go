/*
 * Copyright © 2021 Hedzr Yeh.
 */

package dir

import (
	"os"
	"syscall"
	"time"
)

// FileCreatedTime return the creation time of a file
func FileCreatedTime(fileInfo os.FileInfo) (tm time.Time) {
	tm = time.Unix(int64(fileInfo.Sys().(*syscall.Dir).Mtime), 0)
	return
}

// FileAccessedTime return the creation time of a file
func FileAccessedTime(fileInfo os.FileInfo) (tm time.Time) {
	tm = time.Unix(int64(fileInfo.Sys().(*syscall.Dir).Atime), 0)
	return
}

// FileModifiedTime return the creation time of a file
func FileModifiedTime(fileInfo os.FileInfo) (tm time.Time) {
	tm = time.Unix(int64(fileInfo.Sys().(*syscall.Dir).Mtime), 0)
	return
}
