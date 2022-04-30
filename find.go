// Copyright © 2020 Hedzr Yeh.

package cmdr

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/hedzr/cmdr/conf"
	"github.com/hedzr/cmdr/tool"
	log2 "log"
	"os"
	"strings"
	"text/template"
	"time"
)

// FindSubCommand find sub-command with `longName` from `cmd`
// if cmd == nil: finding from root command
func FindSubCommand(longName string, cmd *Command) (res *Command) {
	if cmd == nil {
		cmd = &internalGetWorker().rootCommand.Command
	}
	res = cmd.FindSubCommand(longName)
	return
}

// FindFlag find flag with `longName` from `cmd`
// if cmd == nil: finding from root command
func FindFlag(longName string, cmd *Command) (res *Flag) {
	if cmd == nil {
		cmd = &internalGetWorker().rootCommand.Command
	}
	res = cmd.FindFlag(longName)
	return
}

// FindSubCommandRecursive find sub-command with `longName` from `cmd` recursively
// if cmd == nil: finding from root command
func FindSubCommandRecursive(longName string, cmd *Command) (res *Command) {
	if cmd == nil {
		cmd = &internalGetWorker().rootCommand.Command
	}
	res = cmd.FindSubCommandRecursive(longName)
	return
}

// FindFlagRecursive find flag with `longName` from `cmd` recursively
// if cmd == nil: finding from root command
func FindFlagRecursive(longName string, cmd *Command) (res *Flag) {
	if cmd == nil {
		cmd = &internalGetWorker().rootCommand.Command
	}
	res = cmd.FindFlagRecursive(longName)
	return
}

//
//

func manBr(s string) string {
	var lines []string //nolint:prealloc //can't prealloc
	for _, l := range strings.Split(s, "\n") {
		lines = append(lines, l+"\n.br")
	}
	return strings.Join(lines, "\n")
}

func manWs(fmtStr string, args ...interface{}) string {
	str := fmt.Sprintf(fmtStr, args...)
	str = replaceAll(strings.TrimSpace(str), " ", `\ `)
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
			if lastLine == "" {
				lastLine = strings.TrimSpace(l)
				// ignore multiple empty lines, compat them as one line.
				if lastLine != "" {
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
		log2.Printf("tpl execute error: %v", err)
		return ""
	}
	return w.String()
}

//
//

func (w *ExecWorker) setupRootCommand(rootCmd *RootCommand) {
	uniqueWorkerLock.Lock()

	w.rootCommand = rootCmd

	w.ow()

	w.rootCommand.PreActions = append(w.rootCommand.PreActions, w.preActions...)
	w.rootCommand.PostActions = append(w.rootCommand.PostActions, w.postActions...)

	uniqueWorkerLock.Unlock()

	w.syncRootCommand()
}

func (w *ExecWorker) ow() {
	if w.bufferedStdio {
		if w.defaultStdout == nil {
			w.defaultStdout = bufio.NewWriterSize(os.Stdout, 16384)
		}
		if w.defaultStderr == nil {
			w.defaultStderr = bufio.NewWriterSize(os.Stderr, 16384)
		}

		w.rootCommand.ow = w.defaultStdout
		w.rootCommand.oerr = w.defaultStderr
	} else {
		w.rootCommand.ow = bufio.NewWriter(os.Stdout)
		w.rootCommand.oerr = bufio.NewWriter(os.Stderr)
	}
}

func (w *ExecWorker) syncRootCommand() {
	if conf.AppName == "" {
		conf.AppName = w.rootCommand.AppName
		conf.Version = w.rootCommand.Version
	}
	_ = os.Setenv("APPNAME", conf.AppName)
	if conf.Buildstamp == "" {
		conf.Buildstamp = time.Now().Format(time.RFC1123)
	}

	w.rootCommand.Version = conf.Version

	SetRaw("cmdr.Version", Version)
}

func (w *ExecWorker) preparePtPkg(pkg *ptpkg) {
	if w.rootCommand.RunAsSubCommand != "" {
		pkg.aliasCommand = dottedPathToCommand(w.rootCommand.RunAsSubCommand, nil)
	}
}

func (w *ExecWorker) getPrefix() string {
	return strings.Join(w.rxxtPrefixes, ".")
}

func (w *ExecWorker) tmpGetRemainArgs(pkg *ptpkg, args []string) []string {
	if pkg.remainArgs == nil && pkg.i < len(args) {
		return args[pkg.i:]
	}
	return pkg.remainArgs
}

func (w *ExecWorker) getRemainArgs(pkg *ptpkg, args []string) []string {
	if pkg.remainArgs == nil && pkg.i < len(args) {
		pkg.remainArgs = args[pkg.i:]
	}
	return pkg.remainArgs
}

// GetRemainArgs returns the remain arguments after command line parsed
func GetRemainArgs() []string {
	w := internalGetWorker()
	return w.getRemainArgs(w.lastPkg, tool.SavedOsArgs)
}
