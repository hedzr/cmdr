/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"os"
	"testing"
)

// Worker returns unexported worker for testing
func Worker() *ExecWorker {
	return uniqueWorker
}

// Worker2 + shouldIgnoreWrongEnumValue
func Worker2(b bool) *ExecWorker {
	uniqueWorker.shouldIgnoreWrongEnumValue = b
	return uniqueWorker
}

// ResetWorker function
func ResetWorker() {
	uniqueWorker = &ExecWorker{
		envPrefixes:  []string{"CMDR"},
		rxxtPrefixes: []string{"app"},

		predefinedLocations: []string{
			"./ci/etc/%s/%s.yml",
			"/etc/%s/%s.yml",
			"/usr/local/etc/%s/%s.yml",
			"$HOME/.%s/%s.yml",
			"$HOME/.config/%s/%s.yml",
		},

		shouldIgnoreWrongEnumValue: true,

		enableVersionCommands:  true,
		enableHelpCommands:     true,
		enableVerboseCommands:  true,
		enableCmdrCommands:     true,
		enableGenerateCommands: true,

		doNotLoadingConfigFiles: false,

		currentHelpPainter: new(helpPainter),

		defaultStdout: bufio.NewWriterSize(os.Stdout, 16384),
		defaultStderr: bufio.NewWriterSize(os.Stderr, 16384),

		rxxtOptions: NewOptions(),
	}
}

func TestTrapSignals(t *testing.T) {
	TrapSignals(func(s os.Signal) {
		//
	})

	// testTypes(t)
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

// ExecWith is main entry of `cmdr`.
//
// Deprecated: from v1.5.0
//
// for Testing
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
//
// Deprecated: from v1.5.0
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
