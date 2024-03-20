package exec

import (
	"strings"

	"gopkg.in/hedzr/errors.v3"
)

// Call executes the command line via system (OS).
//
// DO NOT QUOTE: in 'cmd', A command line shouldn't has quoted parts.
// These are bad:
//
//	cmd := "ls '/usr/bin'"
//	cmd := `tar "c:/My Documents/"`
//
// Uses CallSlice if your args includes space (like 'c:/My Documents/')
func Call(cmd string, fn func(retCode int, stdoutText string)) (err error) {
	a := strings.Split(cmd, " ")
	err = internalCallImpl(a, fn, true)
	return
}

// CallQuiet executes the command line via system (OS) without error printing.
//
// DO NOT QUOTE: in 'cmd', A command line shouldn't has quoted parts.
// These are bad:
//
//	cmd := "ls '/usr/bin'"
//	cmd := `tar "c:/My Documents/"`
//
// Uses CallSliceQuiet if your args includes space (like 'c:/My Documents/')
func CallQuiet(cmd string, fn func(retCode int, stdoutText string)) (err error) {
	a := strings.Split(cmd, " ")
	err = internalCallImpl(a, fn, false)
	return
}

// CallSlice executes the command line via system (OS).
func CallSlice(cmd []string, fn func(retCode int, stdoutText string)) (err error) {
	err = internalCallImpl(cmd, fn, true)
	return
}

// CallSliceQuiet executes the command line via system (OS) without error printing.
func CallSliceQuiet(cmd []string, fn func(retCode int, stdoutText string)) (err error) {
	err = internalCallImpl(cmd, fn, false)
	return
}

// internalCallImpl executes the command line via system (OS) without error printing.
func internalCallImpl(cmd []string, fn func(retCode int, stdoutText string), autoErrReport bool) (err error) {
	var (
		str string
		rc  int
	)

	_, str, err = RunWithOutput(cmd[0], cmd[1:]...)
	if err != nil {
		if autoErrReport {
			err = errors.New("Error on launching '%v': %v", cmd, err)
		}
		return
	}
	fn(rc, str)
	return
}
