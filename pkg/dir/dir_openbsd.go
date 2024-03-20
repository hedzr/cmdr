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
	ts := fileInfo.Sys().(*syscall.Stat_t)
	tm = timeSpecToTime(ts.Ctimespec)
	return
}

// FileAccessedTime return the creation time of a file
func FileAccessedTime(fileInfo os.FileInfo) (tm time.Time) {
	ts := fileInfo.Sys().(*syscall.Stat_t)
	tm = timeSpecToTime(ts.Atimespec)
	return
}

// FileModifiedTime return the creation time of a file
func FileModifiedTime(fileInfo os.FileInfo) (tm time.Time) {
	ts := fileInfo.Sys().(*syscall.Stat_t)
	tm = timeSpecToTime(ts.Mtimespec)
	return
}
