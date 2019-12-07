/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"bytes"
	"github.com/hedzr/errors"
	"github.com/hedzr/logex"
	"github.com/sirupsen/logrus"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

// Worker returns unexported worker for testing
func Worker() *ExecWorker {
	return internalGetWorker()
}

// Worker2 + shouldIgnoreWrongEnumValue
func Worker2(b bool) *ExecWorker {
	internalGetWorker().shouldIgnoreWrongEnumValue = b
	return internalGetWorker()
}

// ResetWorker function
func ResetWorker() {
	InternalResetWorker()
}

// ResetRootInWorker function
func ResetRootInWorker() {
	internalGetWorker().rootCommand = nil
}

func TestTplApply(t *testing.T) {
	tplApply("{{ .dkl }}", &struct{ sth bool }{false})
}

func TestFlag(t *testing.T) {
	t.Log(IsDebuggerAttached())
	t.Log(InTesting())
	t.Log(StripPrefix("8.yes", "8."))
	t.Log(StripQuotes("'8.yes'"))
	t.Log(IsDigitHeavy("not-digit"))
	t.Log(IsDigitHeavy("8-is-not-digit"))

	in := bytes.NewBufferString("\n")
	PressEnterToContinue(in, "ok...")
	in = bytes.NewBufferString("\n")
	PressEnterToContinue(in)

	in = bytes.NewBufferString("\n")
	t.Log(PressAnyKeyToContinue(in, "ok..."))
	in = bytes.NewBufferString("\n")
	t.Log(PressAnyKeyToContinue(in))

	x := SavedOsArgs
	defer func() {
		SavedOsArgs = x
	}()
	SavedOsArgs = []string{"xx.test"}
	t.Log(InTesting())
	SavedOsArgs = []string{"xx.runtime"}
	t.Log(InTesting())
	SavedOsArgs = []string{"xx.runtime", "-test.v"}
	t.Log(InTesting())

	var rootCmdX = &RootCommand{
		Command: Command{
			BaseOpt: BaseOpt{
				Name: "consul-tags",
			},
			SubCommands: []*Command{
				{
					BaseOpt: BaseOpt{
						Name: "consul-tags",
					},
				},
				{
					BaseOpt: BaseOpt{
						Name: "consul-tags",
					},
				},
			},
		},
	}
	_ = walkFromCommand(&rootCmdX.Command, 0, func(cmd *Command, index int) (err error) {
		if index > 0 {
			return ErrBadArg
		}
		return nil
	})
}

func TestUnknownXXX(t *testing.T) {
	defer logex.CaptureLog(t).Release()

	// // RaiseInterrupt(t, 16)
	// go func() {
	// 	time.Sleep(16 * time.Second)
	// 	SignalTermSignal()
	// }()

	if SavedOsArgs == nil {
		SavedOsArgs = os.Args
	}
	defer func() {
		os.Args = SavedOsArgs
	}()

	var pkg *ptpkg
	var cmd *Command
	var args []string

	var rootCmdX = &RootCommand{
		Command: Command{
			BaseOpt: BaseOpt{
				Name: "consul-tags",
			},
		},
	}
	cmd = &rootCmdX.Command
	var commands = []string{
		"consul-tags --help -q",
	}
	for _, cc := range commands {
		os.Args = strings.Split(cc, " ")
		SetInternalOutputStreams(nil, nil)
		ResetOptions()
		if err := Exec(rootCmdX); err != nil {
			t.Fatal(err)
		}
	}

	pkg = &ptpkg{}
	unknownCommand(pkg, cmd, args)
	unknownFlagDetector(pkg, cmd, args)
}

// TestSliceConverters functions
func TestSliceConverters(t *testing.T) {
	stringSliceToInt64Slice([]string{"x"})
	intSliceToUint64Slice([]int{1})
	int64SliceToUint64Slice([]int64{1})
	uint64SliceToInt64Slice([]uint64{1})

	w := internalGetWorker()

	Set("x", []string{"1"})
	GetIntSliceR("x")
	w.rxxtOptions.GetInt64Slice("app.x")
	w.rxxtOptions.GetUint64Slice("app.x")
	Set("x", "1,2,3")
	GetIntSliceR("x")
	w.rxxtOptions.GetInt64Slice("app.x")
	w.rxxtOptions.GetUint64Slice("app.x")
	Set("x", []int{1, 2})
	GetIntSliceR("x")
	w.rxxtOptions.GetInt64Slice("app.x")
	w.rxxtOptions.GetUint64Slice("app.x")
	Set("x", []int64{1, 2})
	GetIntSliceR("x")
	w.rxxtOptions.GetInt64Slice("app.x")
	w.rxxtOptions.GetUint64Slice("app.x")
	Set("x", []uint64{1, 2})
	GetIntSliceR("x")
	w.rxxtOptions.GetInt64Slice("app.x")
	w.rxxtOptions.GetUint64Slice("app.x")
	Set("x", []byte{1, 2})
	GetIntSliceR("x")
	w.rxxtOptions.GetInt64Slice("app.x")
	w.rxxtOptions.GetUint64Slice("app.x")
	Set("x", 57)
	GetIntSliceR("x")
	w.rxxtOptions.GetInt64Slice("app.x")
	w.rxxtOptions.GetUint64Slice("app.x")

	mxIx("", "")
}

func (pkg *ptpkg) setOwner(cmd *Command) {
	if pkg.flg != nil {
		pkg.flg.owner = cmd
	}
}

func TestLog(t *testing.T) {
	var rootCmdX = &RootCommand{
		Command: Command{
			BaseOpt: BaseOpt{
				Name: "consul-tags",
			},
		},
	}

	for _, x := range []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL", "PANIC", ""} {
		Set("logger.level", x)
		_ = internalGetWorker().getWithLogexInitializor(logrus.DebugLevel)(&rootCmdX.Command, []string{})
	}

	Set("logger.target", "journal")
	Set("logger.format", "json")
	_ = internalGetWorker().getWithLogexInitializor(logrus.DebugLevel)(&rootCmdX.Command, []string{})
}

// TestPtpkgToggleGroup functions
func TestPtpkgToggleGroup(t *testing.T) {
	pkg := &ptpkg{flg: &Flag{
		ToggleGroup: "XX",
	}}
	pkg.setOwner(&Command{
		Flags: []*Flag{
			{
				ToggleGroup: "XX",
			},
			{
				ToggleGroup: "XX",
			},
		},
	})

	pkg.toggleGroup()

	pkg = &ptpkg{flg: &Flag{
		DefaultValue: time.Second,
	}}
	_ = pkg.tryExtractingOthers([]string{}, reflect.Chan)
	_ = pkg.tryExtractingOthers([]string{"sss"}, reflect.Int)
	_ = pkg.processExternalTool()
}

// ExecWith is main entry of `cmdr`.
// for testing
func ExecWith(rootCmd *RootCommand, beforeXrefBuildingX, afterXrefBuiltX HookFunc) (err error) {
	w := internalGetWorker()

	if beforeXrefBuildingX != nil {
		w.beforeXrefBuilding = append(w.beforeXrefBuilding, beforeXrefBuildingX)
	}
	if afterXrefBuiltX != nil {
		w.afterXrefBuilt = append(w.afterXrefBuilt, afterXrefBuiltX)
	}

	err = w.InternalExecFor(rootCmd, os.Args)
	return
}

// SetInternalOutputStreams sets the internal output streams for debugging
// for testing
func SetInternalOutputStreams(out, err *bufio.Writer) {
	w := internalGetWorker()
	w.defaultStdout = out
	w.defaultStderr = err

	if w.defaultStdout == nil {
		w.defaultStdout = bufio.NewWriterSize(os.Stdout, 16384)
	}
	if w.defaultStderr == nil {
		w.defaultStderr = bufio.NewWriterSize(os.Stderr, 16384)
	}
}

// SetPredefinedLocationsForTesting
// for testing
func SetPredefinedLocationsForTesting(locations ...string) {
	internalGetWorker().predefinedLocations = locations
}

func TestNewError(t *testing.T) {

	errWrongEnumValue := newErrTmpl("unexpect enumerable value '%s' for option '%s', under command '%s'")

	err := newError(false, errWrongEnumValue, "ds", "head", "server")
	println(err)

	err = newError(true, newErr("unexpect enumerable value"))
	println(err.Error())

	err = newErrorWithMsg("Holo", errors.New("unexpect enumerable value"))
	println(err.Error())

	errWrongEnumValue = newErrTmpl("unexpect enumerable value '%s' for option '%s', under command '%s'")
	_ = errWrongEnumValue.Template("x").Format().Msg("x %v", 1).Nest(err)
}
