/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

// FindSubCommand find sub-command with `longName` from `cmd`
func FindSubCommand(longName string, cmd *Command) (res *Command) {
	for _, cx := range cmd.SubCommands {
		if longName == cx.Full {
			res = cx
			return
		}
	}
	return
}

// FindFlag find flag with `longName` from `cmd`
func FindFlag(longName string, cmd *Command) (res *Flag) {
	for _, cx := range cmd.Flags {
		if longName == cx.Full {
			res = cx
			return
		}
	}
	return
}

// FindSubCommandRecursive find sub-command with `longName` from `cmd` recursively
func FindSubCommandRecursive(longName string, cmd *Command) (res *Command) {
	for _, cx := range cmd.SubCommands {
		if longName == cx.Full {
			res = cx
			return
		}
	}
	for _, cx := range cmd.SubCommands {
		if len(cx.SubCommands) > 0 {
			if res = FindSubCommandRecursive(longName, cx); res != nil {
				return
			}
		}
	}
	return
}

// FindFlagRecursive find flag with `longName` from `cmd` recursively
func FindFlagRecursive(longName string, cmd *Command) (res *Flag) {
	for _, cx := range cmd.Flags {
		if longName == cx.Full {
			res = cx
			return
		}
	}
	for _, cx := range cmd.SubCommands {
		// if len(cx.SubCommands) > 0 {
		if res = FindFlagRecursive(longName, cx); res != nil {
			return
		}
		// }
	}
	return
}

func manBr(s string) string {
	var lines []string
	for _, l := range strings.Split(s, "\n") {
		lines = append(lines, l+"\n.br")
	}
	return strings.Join(lines, "\n")
}

func manWs(fmtStr string, args ...interface{}) string {
	str := fmt.Sprintf(fmtStr, args...)
	str = strings.ReplaceAll(strings.TrimSpace(str), " ", `\ `)
	return str
}

func manExamples(s string, data interface{}) string {
	var (
		sources  = strings.Split(s, "\n")
		lines    []string
		lastLine string
	)
	for _, l := range sources {
		if strings.HasPrefix(l, "$ {{.AppName}}") {
			lines = append(lines, `.TP \w'{{.AppName}}\ 'u
.BI {{.AppName}} \ `+manWs(l[14:]))
		} else {
			if len(lastLine) == 0 {
				lastLine = strings.TrimSpace(l)
				// ignore multiple empty lines, compat them as one line.
				if len(lastLine) != 0 {
					lines = append(lines, lastLine+"\n.br")
				}
			} else {
				lastLine = strings.TrimSpace(l)
				lines = append(lines, lastLine+"\n.br")
			}
		}
	}
	return tplApply(strings.Join(lines, "\n"), data)
}

func tplApply(tmpl string, data interface{}) string {
	var w = new(bytes.Buffer)
	var tpl = template.Must(template.New("x").Parse(tmpl))
	if err := tpl.Execute(w, data); err != nil {
		logrus.Errorf("tpl execute error: %v", err)
	}
	return w.String()
}

//
// external
//

// Launch executes a command setting both standard input, output and error.
func Launch(cmd string, args ...string) (err error) {
	c := exec.Command(cmd, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	err = c.Run()

	if err != nil {
		if _, isExitError := err.(*exec.ExitError); !isExitError {
			return err
		}
	}

	return nil
}

// // LaunchSudo executes a command under "sudo".
// func LaunchSudo(cmd string, args ...string) error {
// 	return Launch("sudo", append([]string{cmd}, args...)...)
// }

//
// editor
//

// func getEditor() (string, error) {
// 	if GetEditor != nil {
// 		return GetEditor()
// 	}
// 	return exec.LookPath(DefaultEditor)
// }

// SavedOsArgs is a copy of os.Args, just for testing
var SavedOsArgs []string

func init() {
	// bug: can't copt slice to slice: _ = StandardCopier.Copy(&SavedOsArgs, &os.Args)
	for _, s := range os.Args {
		SavedOsArgs = append(SavedOsArgs, s)
	}
}

// InTesting detects whether is running under go test mode
func InTesting() bool {
	return strings.HasSuffix(SavedOsArgs[0], ".test") ||
		strings.Contains(SavedOsArgs[0], "/T/___Test") ||
		strings.Contains(SavedOsArgs[0], "/T/go-build")
}

func randomFilename() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return os.Getenv("HOME") + ".CMDR_EDIT_FILE"
	}
	return fmt.Sprintf("%v/.CMDR_%x", os.Getenv("HOME"), buf)
}

// LaunchEditor launches the specified editor
func LaunchEditor(editor string) (content []byte, err error) {
	return launchEditorWith(editor, randomFilename())
}

// LaunchEditor launches the specified editor with a filename
func LaunchEditorWith(editor string, filename string) (content []byte, err error) {
	return launchEditorWith(editor, filename)
}

func launchEditorWith(editor, filename string) (content []byte, err error) {
	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		if _, isExitError := err.(*exec.ExitError); !isExitError {
			return
		}
	}

	content, err = ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, nil
	}
	return
}
