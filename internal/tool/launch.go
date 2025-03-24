package tool

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/hedzr/is/dir"
	"github.com/hedzr/is/stringtool"
)

// LaunchEditor launches the specified editor
func LaunchEditor(editor string) (content []byte, err error) {
	return LaunchEditorWithGetter(editor, nil, false)
}

func LaunchEditorWithGetter(editor string, filenamegetter func() string, simulate bool) (content []byte, err error) {
	getter := filenamegetter
	if getter == nil {
		getter = shellEditorRandomFilename
	}
	return launchEditorImpl(editor, getter(), simulate)
}

func shellEditorRandomFilename() (fn string) {
	buf := make([]byte, 5)
	fn = os.Getenv("HOME") + ".CMDR_edit_file_.tmp"
	if _, err := rand.Read(buf); err == nil {
		// generate a random name like '.CMDR_986ec1b553.tmp'
		fn = fmt.Sprintf("%v/.CMDR_%x.tmp", os.Getenv("HOME"), buf)
	}
	return
}

// LaunchEditorWith launches the specified editor with a filename
func LaunchEditorWith(editor, filename string) (content []byte, err error) {
	return launchEditorImpl(editor, filename, false)
}

func launchEditorImpl(editor, filename string, simulate bool) (content []byte, err error) {
	if simulate {
		content = []byte(stringtool.RandomStringPure(10))
		return
	}

	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		var _t0 *exec.ExitError
		if isExitError := errors.Is(err, _t0); !isExitError {
			return
		}
	}

	var exists bool
	if _, exists, err = dir.Exists(filename); err == nil && exists {
		defer func() { _ = dir.DeleteFile(filename) }()
	}

	content, err = dir.ReadFile(filename)
	if err != nil {
		return []byte{}, nil
	}
	return
}

// // ReadFile reads the file named by filename and returns the contents.
// // A successful call returns err == nil, not err == EOF. Because ReadFile
// // reads the whole file, it does not treat an EOF from Read as an error
// // to be reported.
// //
// // As of Go 1.16, this function simply calls os.ReadFile.
// func ReadFile(filename string) ([]byte, error) {
// 	return os.ReadFile(filename)
// }

// // DeleteFile deletes a file if exists
// func DeleteFile(dst string) (err error) {
// 	str := os.ExpandEnv(dst)
// 	if FileExists(str) {
// 		err = os.Remove(str)
// 	}
// 	return
// }

// // FileExists returns the existence of an directory or file
// func FileExists(filepath string) bool {
// 	if _, err := os.Stat(os.ExpandEnv(filepath)); err != nil {
// 		if os.IsNotExist(err) {
// 			return false
// 		}
// 	}
// 	return true
// }
