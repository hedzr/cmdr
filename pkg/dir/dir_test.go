/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dir_test

import (
	"os"
	"path"
	"syscall"
	"testing"

	"gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/cmdr/v2/pkg/dir"
)

// TestIsDirectory tests more
//
// usage:
//
//	go test ./... -v -test.run '^TestIsDirectory$'
func TestIsDirectory(t *testing.T) {
	t.Logf("osargs[0] = %v", os.Args[0])
	// t.Logf("InTesting: %v", cmdr.InTesting())
	// t.Logf("InDebugging: %v", cmdr.InDebugging())

	dir.NormalizeDir("")

	if yes, err := dir.IsDirectory("./conf.d1"); yes {
		t.Fatal(err)
	}
	if yes, err := dir.IsDirectory("../dir"); !yes {
		t.Fatal(err)
	}
	if yes, err := dir.IsRegularFile("./doc1.golang"); yes {
		t.Fatal(err)
	}
	if yes, err := dir.IsRegularFile("./dir.go"); !yes {
		t.Fatal(err)
	}
}

func TestForDir(t *testing.T) {
	// defer logex.CaptureLog(t).Release()

	dirName := "$HOME/.local/share"
	if !dir.FileExists(dirName) {
		dirName = "$HOME/.config/containers"
	}
	if !dir.FileExists(dirName) {
		dirName = "$HOME/.config"
	}

	err := dir.ForDir(dirName, func(depth int, dirName string, fi os.DirEntry) (stop bool, err error) {
		if fi.IsDir() {
			t.Logf("  - dir: %v - [%v]", dirName, fi.Name())
		} else {
			t.Logf("  - file: %v - %v", dirName, fi.Name())
		}
		return
	})

	if err != nil && !errors.TypeIs(err, &os.PathError{Err: syscall.ENOENT}) {
		if err != nil && !errors.TypeIs(err, syscall.ENOENT) {
			t.Errorf("wrong for ForDir(): %v", err)
		}
	}
}

func TestForDirMax(t *testing.T) {
	// defer logex.CaptureLog(t).Release()

	dirName := "$HOME/.local"
	if !dir.FileExists(dirName) {
		dirName = "$HOME/.config"
	}

	err := dir.ForDirMax(dirName, 0, 2, func(depth int, dirName string, fi os.DirEntry) (stop bool, err error) {
		if fi.IsDir() {
			t.Logf("  - dir: %v - [%v]", dirName, fi.Name())
		} else {
			t.Logf("  - file: %v - %v", dirName, fi.Name())
		}
		return
	})

	if err != nil {
		t.Errorf("wrong for ForDir(): %v", err)
	}
}

// func TestWalk(t *testing.T) {
//	dirName := "$HOME/.local"
//	if !dir.FileExists(dirName) {
//		dirName = "$HOME/.config"
//	}
//	err := filepath.Walk(os.ExpandEnv(dirName), func(path string, info fs.FileInfo, err error) error {
//		if info == nil {
//			t.Logf("  - file: %v - ERR: %v", path, err)
//		} else {
//			t.Logf("  - file: %v - %v - %v", path, info.Name(), err)
//		}
//		return err
//	})
//	if err != nil {
//		t.Error(err)
//	}
// }

func TestForFileMax(t *testing.T) {
	// defer logex.CaptureLog(t).Release()

	// dirName := "$HOME/.local"
	// if !dir.FileExists(dirName) {
	//	dirName = "$HOME/.config"
	// }
	dirName := "$HOME/.config"
	if !dir.FileExists(dirName) {
		dirName = "$HOME"
	}

	err := dir.ForFileMax(dirName, 0, 6, func(depth int, dirName string, fi os.DirEntry) (stop bool, err error) {
		if fi.IsDir() {
			t.Logf("  - dir: %v - [%v]", dirName, fi.Name())
		} else {
			t.Logf("  - file: %v - %v", dirName, fi.Name())
		}
		return
	}, "*/1.x/*", "*/1.x", "*/2.c", "*/node_modules", "*/.git", "*/usr*", "*/share")

	if err != nil {
		t.Errorf("wrong for ForDir(): %v", err)
	}
}

func TestGetExecutableDir(t *testing.T) {
	t.Logf("GetExecutablePath = %v", dir.GetExecutablePath())
	t.Logf("GetExecutableDir = %v", dir.GetExecutableDir())
	t.Logf("GetCurrentDir = %v", dir.GetCurrentDir())

	fn := path.Join(dir.GetCurrentDir(), "dir.go")
	if ok, err := dir.IsRegularFile(fn); err != nil || !ok {
		t.Fatal("expecting regular file detected.")
	}
	if !dir.FileExists(fn) {
		t.Fatal("expecting regular file existed.")
	}

	fileInfo, err := os.Stat(fn)
	if err != nil {
		t.Fatal(err)
	}

	_ = dir.FileModeIs(fn, dir.IsModeIrregular)
	_ = dir.FileModeIs(fn, dir.IsModeRegular)
	_ = dir.FileModeIs(fn, dir.IsModeDirectory)
	_ = dir.FileModeIs("/etc", dir.IsModeDirectory)
	_ = dir.FileModeIs("/etc", dir.IsModeIrregular)
	_ = dir.FileModeIs("/etc/not-existence", dir.IsModeIrregular)

	_ = dir.IsModeExecOwner(fileInfo.Mode())
	_ = dir.IsModeExecGroup(fileInfo.Mode())
	_ = dir.IsModeExecOther(fileInfo.Mode())
	_ = dir.IsModeExecAny(fileInfo.Mode())
	_ = dir.IsModeExecAll(fileInfo.Mode())

	_ = dir.IsModeWriteOwner(fileInfo.Mode())
	_ = dir.IsModeWriteGroup(fileInfo.Mode())
	_ = dir.IsModeWriteOther(fileInfo.Mode())
	_ = dir.IsModeWriteAny(fileInfo.Mode())
	_ = dir.IsModeWriteAll(fileInfo.Mode())

	_ = dir.IsModeReadOwner(fileInfo.Mode())
	_ = dir.IsModeReadGroup(fileInfo.Mode())
	_ = dir.IsModeReadOther(fileInfo.Mode())
	_ = dir.IsModeReadAny(fileInfo.Mode())
	_ = dir.IsModeReadAll(fileInfo.Mode())

	_ = dir.IsModeDirectory(fileInfo.Mode())
	_ = dir.IsModeSymbolicLink(fileInfo.Mode())
	_ = dir.IsModeDevice(fileInfo.Mode())
	_ = dir.IsModeNamedPipe(fileInfo.Mode())
	_ = dir.IsModeSocket(fileInfo.Mode())
	_ = dir.IsModeSetuid(fileInfo.Mode())
	_ = dir.IsModeSetgid(fileInfo.Mode())
	_ = dir.IsModeCharDevice(fileInfo.Mode())
	_ = dir.IsModeSticky(fileInfo.Mode())
	_ = dir.IsModeIrregular(fileInfo.Mode())
}

func TestEnsureDir(t *testing.T) { //nolint:revive
	//

	if err := dir.EnsureDir(""); err == nil {
		t.Fatal("expecting an error.")
	}

	if err := dir.EnsureDirEnh(""); err == nil {
		t.Fatal("expecting an error.")
	}

	//

	dn := path.Join(dir.GetCurrentDir(), ".tmp.1")
	if err := dir.EnsureDir(dn); err != nil {
		t.Fatal(err)
	}
	if err := dir.RemoveDirRecursive(dn); err != nil {
		t.Fatal(err)
	}

	//

	dn = path.Join(dir.GetCurrentDir(), ".github")
	if err := dir.EnsureDir(dn); err != nil {
		t.Fatal(err)
	}
	if err := dir.EnsureDirEnh(dn); err != nil {
		t.Fatal(err)
	}

	dn = path.Join(dir.GetCurrentDir(), ".tmp1")
	if err := dir.EnsureDirEnh(dn); err != nil {
		t.Fatal(err)
	}
	if err := dir.RemoveDirRecursive(dn); err != nil {
		t.Fatal(err)
	}

	dn = path.Join(dn, ".tmp2")
	if err := dir.EnsureDirEnh(dn); err != nil {
		t.Fatal(err)
	}
	if err := dir.RemoveDirRecursive(dn); err != nil {
		t.Fatal(err)
	}
}

func TestNormalizeDir(t *testing.T) { //nolint:revive
	dir.NormalizeDir("")
	dir.NormalizeDir(".")
	dir.NormalizeDir("./ad/./c")
	dir.NormalizeDir("./ad/../c")
	dir.NormalizeDir("/ad/./c")
	dir.NormalizeDir("../ad/./c")
	dir.NormalizeDir("~/ad/./c")
}

func TestDirTimestamps(t *testing.T) {
	fileInfo, err := os.Stat("/tmp")
	if err != nil {
		return
	}
	t.Logf("create time: %v", dir.FileCreatedTime(fileInfo))
	t.Logf("access time: %v", dir.FileAccessedTime(fileInfo))
	t.Logf("modified time: %v", dir.FileModifiedTime(fileInfo))
}
