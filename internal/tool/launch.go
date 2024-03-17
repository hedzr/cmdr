package tool

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// LaunchEditor launches the specified editor
func LaunchEditor(editor string, fnamegetter func() string) (content []byte, err error) {
	getter := fnamegetter
	if getter == nil {
		getter = shellEditorRandomFilename
	}
	return launchEditorImpl(editor, getter())
}

func shellEditorRandomFilename() (fn string) {
	buf := make([]byte, 16)
	fn = os.Getenv("HOME") + ".CMDR_EDIT_FILE"
	if _, err := rand.Read(buf); err == nil {
		fn = fmt.Sprintf("%v/.CMDR_%x", os.Getenv("HOME"), buf)
	}
	return
}

// LaunchEditorWith launches the specified editor with a filename
func LaunchEditorWith(editor, filename string) (content []byte, err error) {
	return launchEditorImpl(editor, filename)
}

func launchEditorImpl(editor, filename string) (content []byte, err error) {
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

	content, err = ReadFile(filename)
	defer func() { _ = DeleteFile(filename) }()
	if err != nil {
		return []byte{}, nil
	}
	return
}

// ReadFile reads the file named by filename and returns the contents.
// A successful call returns err == nil, not err == EOF. Because ReadFile
// reads the whole file, it does not treat an EOF from Read as an error
// to be reported.
//
// As of Go 1.16, this function simply calls os.ReadFile.
func ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// DeleteFile deletes a file if exists
func DeleteFile(dst string) (err error) {
	str := os.ExpandEnv(dst)
	if FileExists(str) {
		err = os.Remove(str)
	}
	return
}

// FileExists returns the existence of an directory or file
func FileExists(filepath string) bool {
	if _, err := os.Stat(os.ExpandEnv(filepath)); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
