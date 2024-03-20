// Copyright © 2020 Hedzr Yeh.

package dir

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/cmdr/v2/pkg/exec"
)

// GetExecutableDir returns the executable file directory
func GetExecutableDir() string {
	// _ = ioutil.WriteFile("/tmp/11", []byte(strings.Join(os.Args,",")), 0644)
	// fmt.Printf("os.Args[0] = %v\n", os.Args[0])

	p, _ := os.Executable()
	p, _ = filepath.Abs(p)
	d, _ := filepath.Abs(filepath.Dir(p))
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	// fmt.Println(d)
	return d
}

// GetExecutablePath returns the executable file path
func GetExecutablePath() string {
	p, _ := os.Executable()
	p, _ = filepath.Abs(p)
	return p
}

// GetCurrentDir returns the current workingFlag directory
// it should be equal with os.Getenv("PWD")
func GetCurrentDir() string {
	d, _ := os.Getwd()
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	// fmt.Println(d)
	return d
}

// IsDirectory tests whether `path` is a directory or not
func IsDirectory(filePath string) (bool, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

// IsRegularFile tests whether `path` is a normal regular file or not
func IsRegularFile(filePath string) (bool, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}
	return fileInfo.Mode().IsRegular(), err
}

func timeSpecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(ts.Sec, ts.Nsec)
}

// FileModeIs tests the mode of 'filepath' with 'tester'. Examples:
//
//	var yes = exec.FileModeIs("/etc/passwd", exec.IsModeExecAny)
//	var yes = exec.FileModeIs("/etc/passwd", exec.IsModeDirectory)
func FileModeIs(filePath string, tester func(mode os.FileMode) bool) (ret bool) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return
	}
	ret = tester(fileInfo.Mode())
	return
}

// IsModeRegular give the result of whether a file is a regular file
func IsModeRegular(mode os.FileMode) bool { return mode.IsRegular() }

// IsModeDirectory give the result of whether a file is a directory
func IsModeDirectory(mode os.FileMode) bool { return mode&os.ModeDir != 0 }

// IsModeSymbolicLink give the result of whether a file is a symbolic link
func IsModeSymbolicLink(mode os.FileMode) bool { return mode&os.ModeSymlink != 0 }

// IsModeDevice give the result of whether a file is a device
func IsModeDevice(mode os.FileMode) bool { return mode&os.ModeDevice != 0 }

// IsModeNamedPipe give the result of whether a file is a named pipe
func IsModeNamedPipe(mode os.FileMode) bool { return mode&os.ModeNamedPipe != 0 }

// IsModeSocket give the result of whether a file is a socket file
func IsModeSocket(mode os.FileMode) bool { return mode&os.ModeSocket != 0 }

// IsModeSetuid give the result of whether a file has the setuid bit
func IsModeSetuid(mode os.FileMode) bool { return mode&os.ModeSetuid != 0 }

// IsModeSetgid give the result of whether a file has the setgid bit
func IsModeSetgid(mode os.FileMode) bool { return mode&os.ModeSetgid != 0 }

// IsModeCharDevice give the result of whether a file is a character device
func IsModeCharDevice(mode os.FileMode) bool { return mode&os.ModeCharDevice != 0 }

// IsModeSticky give the result of whether a file is a sticky file
func IsModeSticky(mode os.FileMode) bool { return mode&os.ModeSticky != 0 }

// IsModeIrregular give the result of whether a file is a non-regular file; nothing else is known about this file
func IsModeIrregular(mode os.FileMode) bool { return mode&os.ModeIrregular != 0 }

//

// IsModeExecOwner give the result of whether a file can be invoked by its unix-owner
func IsModeExecOwner(mode os.FileMode) bool { return mode&0100 != 0 }

// IsModeExecGroup give the result of whether a file can be invoked by its unix-group
func IsModeExecGroup(mode os.FileMode) bool { return mode&0010 != 0 }

// IsModeExecOther give the result of whether a file can be invoked by its unix-all
func IsModeExecOther(mode os.FileMode) bool { return mode&0001 != 0 }

// IsModeExecAny give the result of whether a file can be invoked by anyone
func IsModeExecAny(mode os.FileMode) bool { return mode&0111 != 0 }

// IsModeExecAll give the result of whether a file can be invoked by all users
func IsModeExecAll(mode os.FileMode) bool { return mode&0111 == 0111 }

//

// IsModeWriteOwner give the result of whether a file can be written by its unix-owner
func IsModeWriteOwner(mode os.FileMode) bool { return mode&0200 != 0 }

// IsModeWriteGroup give the result of whether a file can be written by its unix-group
func IsModeWriteGroup(mode os.FileMode) bool { return mode&0020 != 0 }

// IsModeWriteOther give the result of whether a file can be written by its unix-all
func IsModeWriteOther(mode os.FileMode) bool { return mode&0002 != 0 }

// IsModeWriteAny give the result of whether a file can be written by anyone
func IsModeWriteAny(mode os.FileMode) bool { return mode&0222 != 0 }

// IsModeWriteAll give the result of whether a file can be written by all users
func IsModeWriteAll(mode os.FileMode) bool { return mode&0222 == 0222 }

//

// IsModeReadOwner give the result of whether a file can be read by its unix-owner
func IsModeReadOwner(mode os.FileMode) bool { return mode&0400 != 0 }

// IsModeReadGroup give the result of whether a file can be read by its unix-group
func IsModeReadGroup(mode os.FileMode) bool { return mode&0040 != 0 }

// IsModeReadOther give the result of whether a file can be read by its unix-all
func IsModeReadOther(mode os.FileMode) bool { return mode&0004 != 0 }

// IsModeReadAny give the result of whether a file can be read by anyone
func IsModeReadAny(mode os.FileMode) bool { return mode&0444 != 0 }

// IsModeReadAll give the result of whether a file can be read by all users
func IsModeReadAll(mode os.FileMode) bool { return mode&0444 == 0444 }

//

// Exists returns the existence of an directory or file.
// See the short version FileExists.
func Exists(filePath string) (fi os.FileInfo, exists bool, err error) {
	if fi, err = os.Stat(os.ExpandEnv(filePath)); err != nil {
		if os.IsNotExist(err) {
			return
		}
	}
	exists = true
	return
}

// FileExists returns the existence of an directory or file
func FileExists(filePath string) bool {
	if _, err := os.Stat(os.ExpandEnv(filePath)); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// EnsureDir checks and creates the directory.
func EnsureDir(d string) (err error) {
	if d == "" {
		return errors.NewLite("empty directory")
	}
	if !FileExists(d) {
		err = os.MkdirAll(d, 0755)
	}
	return
}

// EnsureDirEnh checks and creates the directory, via sudo if necessary.
func EnsureDirEnh(d string) (err error) {
	if d == "" {
		return errors.NewLite("empty directory")
	}
	if !FileExists(d) {
		err = os.MkdirAll(d, 0755)
		err = checkAndSudoMkdir(d, err)
	}
	return
}

func checkAndSudoMkdir(d string, err error) error {
	if e, ok := err.(*os.PathError); ok && e.Err == syscall.EACCES {
		// var u *user.User
		u, err1 := user.Current()
		if err1 == nil {
			if _, _, err1 = exec.Sudo("mkdir", "-p", d); err == nil {
				_, _, err1 = exec.Sudo("chown", u.Username+":", d)
			}

			// if _, _, err = exec.Sudo("mkdir", "-p", d); err != nil {
			//	logrus.Warnf("Failed to create directory %q, using default stderr. error is: %v", d, err)
			// } else if _, _, err = exec.Sudo("chown", u.Username+":", d); err != nil {
			//	logrus.Warnf("Failed to create directory %q, using default stderr. error is: %v", d, err)
			// }
			return err1
		}
	}
	return nil
}

// RemoveDirRecursive removes a directory and any children it contains.
func RemoveDirRecursive(d string) (err error) {
	// RemoveContentsInDir(d)
	err = os.RemoveAll(d)
	return
}

// // RemoveContentsInDir removes all file and sub-directory in a directory
// func RemoveContentsInDir(dir string) error {
// 	d, err := os.Open(dir)
// 	if err != nil {
// 		return err
// 	}
// 	defer d.Close()
// 	names, err := d.Readdirnames(-1)
// 	if err != nil {
// 		return err
// 	}
// 	for _, name := range names {
// 		err = os.RemoveAll(filepath.Join(dir, name))
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// NormalizeDir make dir name normalized
func NormalizeDir(s string) string {
	return normalizeDirImpl(s)
}

func normalizeDirImpl(s string) string {
	p := normalizeDirBasic(s)
	p = filepath.Clean(p)
	return p
}

func normalizeDirBasic(s string) string {
	if s == "" {
		return s
	}

	s1 := os.Expand(s, os.Getenv)
	if s1[0] == '/' {
		return s1
	}
	if strings.HasPrefix(s1, "./") {
		return path.Join(GetCurrentDir(), s1)
	}
	if strings.HasPrefix(s1, "../") {
		return path.Dir(path.Join(GetCurrentDir(), s1))
	}
	if strings.HasPrefix(s1, "~/") {
		return path.Join(os.Getenv("HOME"), s1[2:])
	}

	return s1
}

// AbsPath returns a clean, normalized and absolute path string
// for the given pathname.
func AbsPath(pathname string) string {
	return absPathImpl(pathname)
}

// AbsPathL returns a clean, normalized, symbolic links resolved,
// absolute path string for the given pathname.
func AbsPathL(pathname string) string {
	abs := absPathImpl(pathname)
	if t, err := filepath.EvalSymlinks(abs); err == nil {
		return t
	}
	return abs
}

func absPathImpl(pathname string) (abs string) {
	abs = normalizePathImpl(pathname)
	if s, err := filepath.Abs(abs); err == nil {
		abs = s
	}
	return
}

// FollowSymLink resolves the symbolic links for a given pathname.
func FollowSymLink(pathname string) string {
	if t, err := filepath.EvalSymlinks(pathname); err == nil {
		return t
	}
	return pathname
}

// NormalizePath cleans up the given pathname
func NormalizePath(pathname string) string {
	return normalizePathImpl(pathname)
}

func normalizePathImpl(pathname string) string {
	p := normalizePathBasic(pathname)
	p = filepath.Clean(p)
	return p
}

func normalizePathBasic(pathname string) string {
	if pathname == "" {
		return pathname
	}

	name := os.Expand(pathname, os.Getenv)
	if name[0] == '/' {
		return name
	}
	if strings.HasPrefix(name, "~/") {
		return path.Join(os.Getenv("HOME"), name[2:])
	}
	return name
}

// ForDir walks on `root` directory and its children
func ForDir(
	root string,
	cb func(depth int, dirName string, fi os.DirEntry) (stop bool, err error),
	excludes ...string,
) (err error) {
	err = ForDirMax(root, 0, -1, cb, excludes...)
	return
}

// ForDirMax walks on `root` directory and its children with nested levels up to `maxLength`.
//
// Example - discover folder just one level
//
//	     _ = ForDirMax(dir, 0, 1, func(depth int, dirname string, fi os.FileInfo) (stop bool, err error) {
//				if fi.IsDir() {
//					return
//				}
//	         // ... doing something for a file,
//				return
//			})
//
// maxDepth = -1: no limit.
// initialDepth: 0 if no idea.
func ForDirMax(
	root string,
	initialDepth int,
	maxDepth int,
	cb func(depth int, dirName string, fi os.DirEntry) (stop bool, err error),
	excludes ...string,
) (err error) {
	if maxDepth > 0 && initialDepth >= maxDepth {
		return
	}

	// rootDir := os.ExpandEnv(root)
	rootDir := path.Clean(NormalizeDir(root))

	return forDirMaxR(rootDir, initialDepth, maxDepth, cb, excludes...)
}

func forDirMaxR(
	rootDir string,
	initialDepth int,
	maxDepth int,
	cb func(depth int, dirName string, fi os.DirEntry) (stop bool, err error),
	excludes ...string,
) (err error) {
	var dirs []os.DirEntry
	dirs, err = os.ReadDir(rootDir)
	if err != nil {
		// Logger.Fatalf("error in ForDirMax(): %v", err)
		return
	}

	// var stop bool
	//
	// // // files, err :=os.ReadDir(rootDir)
	// // var fi os.FileInfo
	// // fi, err = os.Stat(rootDir)
	// // if err != nil {
	// // 	return
	// // }
	// // if stop, err = cb(initialDepth, rootDir, fi); stop { //nolint:ineffassign
	// // 	return
	// // }

	_, err = forDirMaxLoops(dirs, rootDir, initialDepth, maxDepth, cb, excludes...) //nolint:ineffassign,staticcheck
	return
}

func forDirMaxLoops( //nolint:revive
	dirs []os.DirEntry,
	rootDir string,
	initialDepth int,
	maxDepth int,
	cb func(depth int, dirName string, fi os.DirEntry) (stop bool, err error),
	excludes ...string,
) (stop bool, err error) {
	ec := errors.New(`forDirMaxLoops have errors`)
	defer ec.Defer(&err)

	for _, f := range dirs {
		// Logger.Printf("  - %v", f.Name())
		if err != nil {
			continue
		}

		if (maxDepth <= 0 || (maxDepth > 0 && initialDepth+1 < maxDepth)) && f.IsDir() {
			d := path.Join(rootDir, f.Name())
			if forFileMatched(d, excludes...) {
				continue
			}

			if stop, err = cb(initialDepth, d, f); stop {
				return
			}
			if err = ForDirMax(d, initialDepth+1, maxDepth, cb); err != nil {
				ec.Attach(err)
			}
		}
	}

	return
}

// ForFile walks on `root` directory and its children
func ForFile(
	root string,
	cb func(depth int, dirName string, fi os.DirEntry) (stop bool, err error),
	excludes ...string,
) (err error) {
	err = ForFileMax(root, 0, -1, cb, excludes...)
	return
}

// ForFileMax walks on `root` directory and its children with nested levels up to `maxLength`.
//
// Example - discover folder just one level
//
//	     _ = ForFileMax(dir, 0, 1, func(depth int, dirName string, fi os.FileInfo) (stop bool, err error) {
//				if fi.IsDir() {
//					return
//				}
//	         // ... doing something for a file,
//				return
//			})
//
// maxDepth = -1: no limit.
// initialDepth: 0 if no idea.
//
// Known issue:
// can't walk at ~/.local/share/NuGet/v3-cache/1ca707a4d90792ce8e42453d4e350886a0fdaa4d:_api.nuget.org_v3_index.json.
// workaround: use filepath.Walk
func ForFileMax(
	root string,
	initialDepth, maxDepth int,
	cb func(depth int, dirName string, fi os.DirEntry) (stop bool, err error),
	excludes ...string,
) (err error) {
	if maxDepth > 0 && initialDepth >= maxDepth {
		return
	}

	// rootDir := os.ExpandEnv(root)
	rootDir := path.Clean(NormalizeDir(root))

	return forFileMaxR(rootDir, initialDepth, maxDepth, cb, excludes...)
}

func forFileMaxR( //nolint:revive
	rootDir string,
	initialDepth, maxDepth int,
	cb func(depth int, dirName string, fi os.DirEntry) (stop bool, err error),
	excludes ...string,
) (err error) {
	ec := errors.New(`forFileMax have errors`)
	defer ec.Defer(&err)

	var dirs []os.DirEntry
	dirs, err = os.ReadDir(rootDir)
	// var dirs []os.DirEntry
	// dirs, err = os.ReadDir(rootDir)
	if err != nil {
		// Logger.Fatalf("error in ForFileMax(): %v", err)
		ec.Attach(err)
		return
	}

	var stop bool
	for _, f := range dirs {
		// Logger.Printf("  - %v", f.Name())
		if err != nil {
			continue
		}

		if (maxDepth <= 0 || (maxDepth > 0 && initialDepth < maxDepth)) && f.IsDir() {
			d := path.Join(rootDir, f.Name())
			if forFileMatched(d, excludes...) {
				continue
			}

			if err = ForFileMax(d, initialDepth+1, maxDepth, cb, excludes...); err != nil {
				ec.Attach(err)
			}

			continue
		}

		if !f.IsDir() {
			d := path.Join(rootDir, f.Name())
			if forFileMatched(d, excludes...) {
				continue
			}

			// log.Infof(" - %s", f.Name())
			// fi, _ := f.Info()
			if stop, err = cb(initialDepth, rootDir, f); stop {
				return
			}
		}
	}

	return
}

// Basename returns filename without pathname and suffix.
// For a file "docs/tech-notes/cpp-dev/ringbuf/RingBuffer.md", Basename returns "RingBuffer".
func Basename(file string) string {
	base, ext := path.Base(file), path.Ext(file)
	if l := len(ext); l > 0 {
		base = base[:len(base)-l]
	}
	return base
}

// Base returns filename without pathname.
// For a file "docs/tech-notes/cpp-dev/ringbuf/RingBuffer.md", Basename returns "RingBuffer.md".
//
// Base is synonym to path.Base
func Base(file string) string { return path.Base(file) }
func Ext(file string) string  { return path.Ext(file) } // return Ext part
func Name(file string) string { return path.Dir(file) } // return DirName part

// RelName returns relative pathname from root
func RelName(file, root string) string {
	if str, yes := strings.CutPrefix(file, root); yes {
		if str[0] == '/' {
			return str[1:]
		}
		return str
	}
	return file
}

// RelNameForce returns the relative pathname even if 'file' and 'root' is not on a same volume.
func RelNameForce(file, root string) string {
	if str, yes := strings.CutPrefix(file, root); yes {
		if str[0] == '/' {
			return str[1:]
		}
		return str
	}
	return strings.Repeat("../", strings.Count(root, "/")) + file
}

// func forDirMatched(f os.DirEntry, root string, excludes ...string) (matched bool) {
//	fullName := path.Join(root, f.Name())
//	for _, ptn := range excludes {
//		if IsWildMatch(fullName, ptn) {
//			matched = true
//			break
//		}
//	}
//	return
// }

// func forFileMatched(f os.FileInfo, root string, excludes ...string) (matched bool) {
//	fullName := path.Join(root, f.Name())
//	matched = inExcludes(fullName, excludes...)
//	//if matched, _ = filepath.Match(ptn, fullName); matched {
//	//	break
//	//}
//	return
// }

func forFileMatched(name string, excludes ...string) (yes bool) {
	for _, ptn := range excludes {
		if yes = IsWildMatch(name, ptn); yes {
			break
		}
	}
	return
}

// PushDir provides a shortcut to enter a folder and restore at
// the end of your current function scope.
//
// PushDir returns a functor and assumes you will DEFER call it.
// For example:
//
//	func TestSth() {
//	    defer dir.PushDir("/your/working/dir")()
//	    // do sth under '/your/working/dir' ...
//	}
//
// BEWARE DON'T miss the ending brackets for defer call.
// NOTE that current directory would not be changed if chdir(dirName) failed,
func PushDir(dirName string) (closer func()) {
	savedDir := GetCurrentDir()
	var err error
	if err = os.Chdir(dirName); err != nil {
		// err = nil //ignore path err
		return func() {}
	}
	return func() {
		if err == nil {
			_ = os.Chdir(savedDir)
		}
	}
}

// PushDirEx provides a shortcut to enter a folder and restore at
// the end of your current function scope.
func PushDirEx(dirName string) (closer func(), err error) {
	savedDir := GetCurrentDir()
	if err = os.Chdir(dirName); err != nil {
		// err = nil //ignore path err
		return
	}
	closer = func() {
		if err == nil {
			_ = os.Chdir(savedDir)
		}
	}
	return
}

// DeleteFile deletes a file if exists
func DeleteFile(dst string) (err error) {
	dst = os.ExpandEnv(dst) //nolint:revive
	if FileExists(dst) {
		err = os.Remove(dst)
	}
	return
}

// CopyFileByLinkFirst copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherwise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFileByLinkFirst(src, dst string) (err error) {
	return copyFileByLinkFirst(src, dst, true)
}

// CopyFile will make a content clone of src.
func CopyFile(src, dst string) (err error) {
	return copyFileByLinkFirst(src, dst, false)
}

func copyFileByLinkFirst(src, dst string, linkAtFirst bool) (err error) { //nolint:revive
	src = os.ExpandEnv(src) //nolint:revive
	dst = os.ExpandEnv(dst) //nolint:revive
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("[CopyFile]: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("[CopyFile]: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if linkAtFirst {
		if err = os.Link(src, dst); err == nil {
			return
		}
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
