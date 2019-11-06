/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"github.com/hedzr/logex"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func internalGetWorker() (w *ExecWorker) {
	w = uniqueWorker
	return
}

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

func RaiseInterrupt(t *testing.T, timeout int) {
	go func() {
		time.Sleep(time.Duration(timeout) * time.Second)
		p, err := os.FindProcess(os.Getpid())
		if err != nil {
			t.Fatal(err)
		}
		err = p.Signal(os.Interrupt)
		if err != nil {
			t.Fatal(err)
		}
	}()
}

func TestTrapSignals(t *testing.T) {

	if SavedOsArgs == nil {
		SavedOsArgs = os.Args
	}
	defer func() {
		os.Args = SavedOsArgs
	}()

	RaiseInterrupt(t, 6)
	TrapSignals(func(s os.Signal) {
		//
	})

	_ = RemoveDirRecursive("docs")

	// testTypes(t)
}

func TestUnknownXXX(t *testing.T) {
	defer logex.CaptureLog(t).Release()
	RaiseInterrupt(t, 16)

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

	Set("x", []string{"1"})
	GetIntSliceR("x")
	uniqueWorker.rxxtOptions.GetInt64Slice("app.x")
	uniqueWorker.rxxtOptions.GetUint64Slice("app.x")
	Set("x", "1,2,3")
	GetIntSliceR("x")
	uniqueWorker.rxxtOptions.GetInt64Slice("app.x")
	uniqueWorker.rxxtOptions.GetUint64Slice("app.x")
	Set("x", []int{1, 2})
	GetIntSliceR("x")
	uniqueWorker.rxxtOptions.GetInt64Slice("app.x")
	uniqueWorker.rxxtOptions.GetUint64Slice("app.x")
	Set("x", []int64{1, 2})
	GetIntSliceR("x")
	uniqueWorker.rxxtOptions.GetInt64Slice("app.x")
	uniqueWorker.rxxtOptions.GetUint64Slice("app.x")
	Set("x", []uint64{1, 2})
	GetIntSliceR("x")
	uniqueWorker.rxxtOptions.GetInt64Slice("app.x")
	uniqueWorker.rxxtOptions.GetUint64Slice("app.x")
	Set("x", []byte{1, 2})
	GetIntSliceR("x")
	uniqueWorker.rxxtOptions.GetInt64Slice("app.x")
	uniqueWorker.rxxtOptions.GetUint64Slice("app.x")
	Set("x", 57)
	GetIntSliceR("x")
	uniqueWorker.rxxtOptions.GetInt64Slice("app.x")
	uniqueWorker.rxxtOptions.GetUint64Slice("app.x")

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
		_ = uniqueWorker.getWithLogexInitializor(&rootCmdX.Command, []string{})
	}

	Set("logger.target", "journal")
	Set("logger.format", "json")
	_ = uniqueWorker.getWithLogexInitializor(&rootCmdX.Command, []string{})
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
	w := uniqueWorker

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
	uniqueWorker.defaultStdout = out
	uniqueWorker.defaultStderr = err

	if uniqueWorker.defaultStdout == nil {
		uniqueWorker.defaultStdout = bufio.NewWriterSize(os.Stdout, 16384)
	}
	if uniqueWorker.defaultStderr == nil {
		uniqueWorker.defaultStderr = bufio.NewWriterSize(os.Stderr, 16384)
	}
}

// SetPredefinedLocationsForTesting
// for testing
func SetPredefinedLocationsForTesting(locations ...string) {
	uniqueWorker.predefinedLocations = locations
}
