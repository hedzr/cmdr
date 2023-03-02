/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hedzr/evendeep"

	cmdrbase "github.com/hedzr/cmdr-base"
	"github.com/hedzr/cmdr/tool"
	"github.com/hedzr/log/detects"
	"github.com/hedzr/log/dir"

	"gopkg.in/hedzr/errors.v3"
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

// Worker3 + shouldIgnoreWrongEnumValue
func Worker3(root *RootCommand) *ExecWorker {
	w := internalGetWorker()
	w.shouldIgnoreWrongEnumValue = false
	w.rootCommand = root
	return w
}

// BacktraceCmdNamesForTest _
func BacktraceCmdNamesForTest(cmd *Command, v bool) string {
	return backtraceCmdNames(cmd, v)
}

// // ResetWorker function
// func ResetWorker() {
// 	InternalResetWorkerForTest()
// }

func StopExitingChannelForFsWatcherAlways() {
	stopExitingChannelForFsWatcherAlways()
}

// InternalResetWorkerForTest is an internal helper, esp for debugging
func InternalResetWorkerForTest() (w *ExecWorker) {
	uniqueWorkerLock.Lock()
	w = internalResetWorkerNoLock()
	noResetWorker = false
	uniqueWorkerLock.Unlock()
	return
}

// InternalResetWorkerNoLockForTest is an internal helper, esp for debugging
func InternalResetWorkerNoLockForTest() (w *ExecWorker) {
	w = internalResetWorkerNoLock()
	return
}

// ResetRootInWorkerForTest function
func ResetRootInWorkerForTest() {
	uniqueWorkerLock.Lock()
	w := internalResetWorkerNoLock()
	w.rootCommand = nil
	uniqueWorkerLock.Unlock()
}

// GetTextPiecesForTest _
func GetTextPiecesForTest(str string, start, want int) string {
	s, _ := getTextPiece(str, start, want)
	return s
}

// Cpt _
func Cpt() ColorTranslator {
	return &cpt
}

// CptNC _
func CptNC() ColorTranslator {
	return &cptNC
}

// NewOptionsForTest _
func NewOptionsForTest() *Options { return newOptions() }

// IsTypeUint _
func IsTypeUint(kind reflect.Kind) bool {
	return isTypeUint(kind)
}

// IsTypeSInt _
func IsTypeSInt(kind reflect.Kind) bool {
	return isTypeSInt(kind)
}

// IsTypeFloat _
func IsTypeFloat(kind reflect.Kind) bool {
	return isTypeFloat(kind)
}

// IsTypeComplex _
func IsTypeComplex(kind reflect.Kind) bool {
	return isTypeComplex(kind)
}

func TestEmptyUnknownOptionHandler(t *testing.T) {
	emptyUnknownOptionHandler(false, "", nil, nil)
}

func TestTplApply(t *testing.T) {
	tplApply("{{ .dkl }}", &struct{ sth bool }{false})
}

func tLog(a ...interface{}) {}

//nolint:funlen //for test
func TestFlag(t *testing.T) {
	ResetOptions()
	ResetRootInWorkerForTest()
	internalGetWorker().rxxtPrefixes = []string{}
	t.Log(wrapWithRxxtPrefix("x"))
	internalGetWorker().rxxtPrefixes = []string{"app"}
	InternalResetWorkerForTest()

	noResetWorker = false
	tLog(GetStringR("version"))
	noResetWorker = true
	tLog(GetStringR("version"))

	t.Log(detects.IsDebuggerAttached())
	t.Log(InTesting())
	t.Log(InDevelopingTime())
	SetDebugMode(false)
	t.Log(GetDebugMode())
	t.Log(InDevelopingTime())
	t.Log(detects.InDebugging())
	SetDebugMode(true)
	t.Log(GetDebugMode())
	t.Log(InDevelopingTime())
	SetTraceMode(false)
	t.Log(GetTraceMode())
	SetTraceMode(true)
	t.Log(GetTraceMode())
	t.Log(detects.InDockerEnvSimple())
	t.Log(tool.StripPrefix("8.yes", "8."))
	t.Log(tool.IsDigitHeavy("not-digit"))
	t.Log(tool.IsDigitHeavy("8-is-not-digit"))

	in := bytes.NewBufferString("\n")
	tool.PressEnterToContinue(in, "ok...")
	in = bytes.NewBufferString("\n")
	tool.PressEnterToContinue(in)

	in = bytes.NewBufferString("\n")
	t.Log(tool.PressAnyKeyToContinue(in, "ok..."))
	in = bytes.NewBufferString("\n")
	t.Log(tool.PressAnyKeyToContinue(in))

	_ = isTypeFloat(reflect.TypeOf(8).Kind())
	_ = isTypeFloat(reflect.TypeOf(8.9).Kind())

	_ = isTypeComplex(reflect.TypeOf(8).Kind())
	_ = isTypeComplex(reflect.TypeOf(8.9).Kind())
	_ = isTypeComplex(reflect.TypeOf(8.9 + 0i).Kind())
	_ = isTypeComplex(reflect.TypeOf(8.9 - 2i).Kind())

	x := tool.SavedOsArgs
	defer func() {
		tool.SavedOsArgs = x
	}()
	tool.SavedOsArgs = []string{"xx.test"}
	t.Log(InTesting())
	tool.SavedOsArgs = []string{"xx.runtime"}
	t.Log(InTesting())
	tool.SavedOsArgs = []string{"xx.runtime", "-test.v"}
	t.Log(InTesting())

	rootCmdX := &RootCommand{
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
	_ = walkFromCommand(&rootCmdX.Command, 0, 0, func(cmd *Command, index, level int) (err error) {
		if index > 0 {
			return ErrBadArg
		}
		return nil
	})
}

func dumpStacks() { //nolint:deadcode,unused //keep it
	fmt.Printf("=== BEGIN goroutine stack dump ===\n%s\n=== END goroutine stack dump ===\n", errors.DumpStacksAsString(true))
}

func TestHandlePanic(t *testing.T) {
	// defer logex.CaptureLog(t).Release()
	if tool.SavedOsArgs == nil {
		tool.SavedOsArgs = os.Args
	}
	defer func() {
		os.Args = tool.SavedOsArgs
	}()

	ResetOptions()
	InternalResetWorkerForTest()

	onUnhandledErrorHandler1 := func(err interface{}) {
		// debug.PrintStack()
		// pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
		// dumpStacks()
		fmt.Println("some error handled: ", err)
	}

	v1, v2 := 11, 0
	// var cmd *Command
	rootCmdX := &RootCommand{
		Command: Command{
			BaseOpt: BaseOpt{
				Name: "consul-tags",
			},
			SubCommands: []*Command{
				{
					BaseOpt: BaseOpt{
						Short: "dz", Full: "division-by-zero",
						Action: func(cmd *Command, args []string) (err error) {
							fmt.Println(v1 / v2)
							return
						},
					},
				},
				{
					BaseOpt: BaseOpt{
						Short: "pa", Full: "panic",
						Action: func(cmd *Command, args []string) (err error) {
							panic(8.1)
							// return
						},
					},
				},
			},
		},
	}

	// cmd = &rootCmdX.Command
	commands := []string{
		"consul-tags dz",
		"consul-tags pa",
	}
	for _, cc := range commands {
		os.Args = strings.Split(cc, " ")
		SetInternalOutputStreams(nil, nil)
		ResetOptions()
		if err := Exec(rootCmdX,
			WithUnhandledErrorHandler(onUnhandledErrorHandler1),
			WithOnSwitchCharHit(func(parsed *Command, switchChar string, args []string) (err error) {
				return
			}),
		); err == nil {
			t.Error("BAD !! / ERROR !! / expecting an error returned without unexpected program terminated") // hi, here is not real error occurs
		}
	}

	t.Log(GetPredefinedLocations())
}

func TestW_parsePredefinedLocation(t *testing.T) {
	w := Worker()
	os.Args = []string{"", "--config=/tmp/1"}
	_ = w.parsePredefinedLocation()
	os.Args = []string{"", "--config/tmp/1"}
	_ = w.parsePredefinedLocation()
	os.Args = []string{"", "--config", "/tmp"}
	_ = w.parsePredefinedLocation()

	GetSecondaryLocations()
	GetPredefinedAlterLocations()
	setPredefinedLocations("")
	setSecondaryLocations("")
	setAlterLocations("")
}

func TestNewOptions(t *testing.T) {
	newOptions()
	newOptionsWith(nil)
}

func TestUnknownXXX(t *testing.T) {
	// defer logex.CaptureLog(t).Release()

	// // RaiseInterrupt(t, 16)
	// go func() {
	// 	time.Sleep(16 * time.Second)
	// 	SignalTermSignal()
	// }()

	if tool.SavedOsArgs == nil {
		tool.SavedOsArgs = os.Args
	}
	defer func() {
		os.Args = tool.SavedOsArgs
	}()

	var pkg *ptpkg
	var cmd *Command
	var args []string

	rootCmdX := &RootCommand{
		Command: Command{
			BaseOpt: BaseOpt{
				Name: "consul-tags",
			},
		},
	}
	cmd = &rootCmdX.Command
	commands := []string{
		"consul-tags --help -q",
	}
	for _, cc := range commands {
		os.Args = strings.Split(cc, " ")
		SetInternalOutputStreams(nil, nil)
		ResetOptions()
		if err := Exec(rootCmdX, WithNoWatchConfigFiles(true)); err != nil {
			t.Fatal(err)
		}
	}

	pkg = &ptpkg{}
	unknownCommand(pkg, cmd, args)
	unknownFlagDetector(pkg, cmd, args)
}

// TestBuildXref _
func TestBuildXrefNilBranch(t *testing.T) {
	w := internalGetWorker()
	w.beforeConfigFileLoading = append(w.beforeConfigFileLoading, func(root *RootCommand, args []string) {})
	w.afterConfigFileLoading = append(w.afterConfigFileLoading, func(root *RootCommand, args []string) {})
	_ = w.buildXref(nil, nil)
}

func rootCmdForAliasesTest() *RootCommand {
	var cc *Command
	root := &RootCommand{
		Command: Command{
			BaseOpt: BaseOpt{
				Full:   "Z",
				Action: func(cmd *Command, args []string) (err error) { return },
			},
			plainCmds:   make(map[string]*Command),
			allCmds:     make(map[string]map[string]*Command),
			Invoke:      "cc arg1",
			InvokeShell: "ls -l",
			InvokeProc:  "ls -a",
			Shell:       "/usr/bin/env bash",
		},
	}

	cc = &Command{
		BaseOpt: BaseOpt{
			Full:    "cc",
			Aliases: []string{"ccc"},
			Group:   "G",
			Action:  func(cmd *Command, args []string) (err error) { return },
		},
		plainCmds: make(map[string]*Command),
		allCmds:   make(map[string]map[string]*Command),
		SubCommands: []*Command{
			{
				BaseOpt: BaseOpt{
					Full:    "cc1",
					Aliases: []string{"ccc1"},
					Action:  func(cmd *Command, args []string) (err error) { return },
				},
				plainCmds: make(map[string]*Command),
				allCmds:   make(map[string]map[string]*Command),
				SubCommands: []*Command{
					{
						BaseOpt: BaseOpt{
							Full:    "cc1c1",
							Aliases: []string{"ccc1c1"},
							Action:  func(cmd *Command, args []string) (err error) { return },
						},
					},
				},
			},
		},
	}

	root.SubCommands = uniAddCmd(root.SubCommands, cc)

	BuildXref(root)

	return root
}

func BuildXref(root *RootCommand) {
	w := internalGetWorker()
	w.doNotLoadingConfigFiles = true
	err := w.buildXref(root, nil)
	if err != nil {
		return
	}
}

// TestBuildAliasesCrossRefsErrBranch _
func TestBuildAliasesCrossRefsErrBranch(t *testing.T) {
	w := internalGetWorker()
	w.buildAliasesCrossRefs(nil)

	w.rootCommand = rootCmdForAliasesTest()
	root := &w.rootCommand.Command
	cc := root.SubCommands[0]

	_ = w._toolAddCmd(root, "G1", cc)
	_, _ = w.locateCommand("cc", nil)
	_, _ = w.locateCommand("", nil)

	var h Handler
	h = w.getInvokeAction(root)
	_ = h(root, []string{"1", "2"})
	h = w.getInvokeProcAction(root)
	_ = h(root, []string{"1", "2"})
	h = w.getInvokeShellAction(root)
	_ = h(root, []string{"1", "2"})
	root.Shell = ""
	_ = h(root, []string{"1", "2"})

	// _ = InvokeCommand("cc")
}

func TestAliasActions(t *testing.T) {
	w := internalGetWorker()
	w.buildAliasesCrossRefs(nil)

	w.rootCommand = rootCmdForAliasesTest()
	root := &w.rootCommand.Command
	cc := root.SubCommands[0]

	_ = w._toolAddCmd(root, "G1", cc)

	_ = InvokeCommand("cc")

	_ = cc.Match("ccc1")
	_ = cc.Match("ccc")
}

func TestColorPrintTool(t *testing.T) {
	for _, s := range []string{
		"black", "red", "green", "yellow", "blue", "magenta", "cyan",
		"lightgray", "light-gray", "darkgray", "dark-gray", "lightred", "light-red",
		"lightgreen", "light-green", "lightyellow", "light-yellow", "lightblue", "light-blue",
		"lightmagenta", "light-magenta", "lightcyan", "light-cyan", "white",
		"??",
	} {
		_ = cpt.toColorInt(s)
	}
	_ = cpt._sz("")
	_ = cpt._ss("\x1b[2m")
}

// TestSliceConverters _
//
//nolint:funlen //for test
func TestSliceConverters(t *testing.T) {
	stringSliceToInt64Slice([]string{"x"})
	intSliceToUint64Slice([]int{1})
	int64SliceToUint64Slice([]int64{1})
	uint64SliceToInt64Slice([]uint64{1})

	w := internalGetWorker()

	var val interface{} = "1,2,3"
	valary := []int{1, 2, 3}

	w.rxxtOptions.setToReplaceMode()

	Set("x", []string{"1"})
	Set("x", val)                                                                 // val = "1,2,3"
	if v1, v2 := GetIntSliceR("x"), []int{1, 2, 3}; !evendeep.DeepEqual(v1, v2) { // equalSlice(reflect.ValueOf(v1), reflect.ValueOf(v2)) {
		t.Errorf("want %v, but got %v", v2, v1)
	}

	Set("x", []string{"5", "7", "19"})
	Set("x", val)                                                                 // val = "1,2,3"
	if v1, v2 := GetIntSliceR("x"), []int{1, 2, 3}; !evendeep.DeepEqual(v1, v2) { // !equalSlice(reflect.ValueOf(v1), reflect.ValueOf(v2)) {
		t.Errorf("want %v, but got %v", v2, v1)
	}

	Set("x", []string{"1"})

	//
	//
	//

	w.rxxtOptions.setToAppendMode()

	SetOverwrite("x", []string{"5", "7", "19"})
	SetOverwrite("x", valary)                                                     // val = "1,2,3"
	if v1, v2 := GetIntSliceR("x"), []int{1, 2, 3}; !evendeep.DeepEqual(v1, v2) { // !equalSlice(reflect.ValueOf(v1), reflect.ValueOf(v2)) {
		t.Errorf("want %v, but got %v", v2, v1)
	}

	SetOverwrite("x", []string{"1"})
	if v1, v2 := GetIntSliceR("x"), []int{1}; !evendeep.DeepEqual(v1, v2) { // !equalSlice(reflect.ValueOf(v1), reflect.ValueOf(v2)) {
		t.Errorf("want %v, but got %v", v2, v1)
	}
	w.rxxtOptions.GetInt64Slice("app.x")
	w.rxxtOptions.GetUint64Slice("app.x")

	SetOverwrite("x", valary)                                                     // val = "1,2,3"
	if v1, v2 := GetIntSliceR("x"), []int{1, 2, 3}; !evendeep.DeepEqual(v1, v2) { // !equalSlice(reflect.ValueOf(v1), reflect.ValueOf(v2)) {
		t.Errorf("want %v, but got %v", v2, v1)
	}
	w.rxxtOptions.GetInt64Slice("app.x")
	w.rxxtOptions.GetUint64Slice("app.x")

	Set("x", []int{1, 2, 4})                                                         // old val type = string, so new value replace it
	if v1, v2 := GetIntSliceR("x"), []int{1, 2, 3, 4}; !evendeep.DeepEqual(v1, v2) { // !equalSlice(reflect.ValueOf(v1), reflect.ValueOf(v2)) {
		t.Errorf("want %v, but got %v", v2, v1)
	}
	w.rxxtOptions.GetInt64Slice("app.x")
	w.rxxtOptions.GetUint64Slice("app.x")

	Set("x", []int64{5, 2})                                                                 // slices will be merged, but without dup elem check of course
	if v1, v2 := GetInt64SliceR("x"), []int64{1, 2, 3, 4, 5}; !evendeep.DeepEqual(v1, v2) { // !equalSlice(reflect.ValueOf(v1), reflect.ValueOf(v2)) {
		t.Errorf("want %v, but got %v", v2, v1)
	}
	w.rxxtOptions.GetIntSlice("app.x")
	w.rxxtOptions.GetUint64Slice("app.x")

	Set("x", []uint64{9, 5})
	if v1, v2 := GetUint64SliceR("x"), []uint64{1, 2, 3, 4, 5, 9}; !evendeep.DeepEqual(v1, v2) { // !equalSlice(reflect.ValueOf(v1), reflect.ValueOf(v2)) {
		t.Errorf("want %v, but got %v", v2, v1)
	}
	w.rxxtOptions.GetInt64Slice("app.x")
	w.rxxtOptions.GetIntSlice("app.x")

	Set("x", []byte{11, 13})
	GetIntSliceR("x")
	w.rxxtOptions.GetInt64Slice("app.x")
	w.rxxtOptions.GetUint64Slice("app.x")

	Set("x", 57)
	GetIntSliceR("x")
	w.rxxtOptions.GetInt64Slice("app.x")
	w.rxxtOptions.GetUint64Slice("app.x")

	mxIx("", "")
}

// TestWithShellCompletionPrivates _
func TestWithShellCompletionPrivates(t *testing.T) {
	withShellCompletionPartialMatch(true)(Worker())
}

// // GetWithLogexInitializer _
// func (w *ExecWorker) GetWithLogexInitializer(lvl Level, opts ...logex.Option) Handler {
//	return w.getWithLogexInitializer(lvl, opts...)
// }

func (pkg *ptpkg) setOwnerForTest(cmd *Command) {
	if pkg.flg != nil {
		pkg.flg.owner = cmd
	}
}

// TestPtpkgToggleGroup functions
func TestPtpkgToggleGroup(t *testing.T) {
	pkg := &ptpkg{flg: &Flag{
		ToggleGroup: "XX",
	}}
	pkg.setOwnerForTest(&Command{
		Flags: []*Flag{
			{
				ToggleGroup: "XX",
			},
			{
				ToggleGroup: "XX",
			},
		},
	})

	pkg.tryToggleGroup()

	pkg = &ptpkg{flg: &Flag{
		DefaultValue: time.Second,
	}}
	_ = pkg.tryExtractingOthers([]string{}, reflect.Chan)
	_ = pkg.tryExtractingOthers([]string{"sss"}, reflect.Int)
	_ = pkg.processExternalTool()
}

// ExecWithForTest is main entry of `cmdr`.
// for testing
func ExecWithForTest(rootCmd *RootCommand, beforeXrefBuildingX, afterXrefBuiltX HookFunc) (err error) {
	w := internalGetWorker()

	if beforeXrefBuildingX != nil {
		w.beforeXrefBuilding = append(w.beforeXrefBuilding, beforeXrefBuildingX)
	}
	if afterXrefBuiltX != nil {
		w.afterXrefBuilt = append(w.afterXrefBuilt, afterXrefBuiltX)
	}

	_, err = w.InternalExecFor(rootCmd, os.Args)
	return
}

// SetInternalOutputStreams sets the internal output streams for debugging
// for testing
func SetInternalOutputStreams(out, err *bufio.Writer) {
	w := internalGetWorker()

	w.defaultStdout = out
	w.defaultStderr = err
	w.bufferedStdio = true

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

// MatchForTest try parsing the input command-line, the result is the last hit *Command.
func MatchForTest(inputCommandlineWithoutArg0 string, opts ...ExecOption) (last *Command, err error) {
	return matchForTest(inputCommandlineWithoutArg0, opts...)
}

func TestNewError(t *testing.T) {
	errWrongEnumValue := newErrTmpl("unexpected enumerable value '%s' for option '%s', under command '%s'")

	err := newError(false, errWrongEnumValue, "ds", "head", "server")
	println(err)

	err = newError(true, newErr("unexpected enumerable value"))
	println(err.Error())

	err = newErrorWithMsg("Holo", errors.New("unexpected enumerable value"))
	println(err.Error())

	var perr *os.PathError
	err = newErrorWithMsg("hooloo", &os.PathError{Err: io.EOF, Op: "find", Path: "/"})
	if errors.As(err, &perr) {
		t.Logf("As() ok: %+v", *perr)
	} else {
		t.Fatal("As() failed: expect it is a os.PathError{}")
	}

	if !err.(*errors.WithStackInfo).As(&perr) { //nolint:errorlint //for test only
		t.Fatal("As() failed: expect it is a os.PathError{}")
	}

	if !err.(*errors.WithStackInfo).Is(perr) { //nolint:errorlint //for test only
		t.Fatal("As() failed: expect it is a os.PathError{}")
	}

	// errWrongEnumValue = newErrTmpl("unexpected enumerable value '%s' for option '%s', under command '%s'")
	// _ = errWrongEnumValue.Template("x").Format().Msg("x %v", 1).Nest(err)
}

func TestGenPowerShell1(t *testing.T) {
	InternalResetWorkerForTest()
	ResetOptions()

	// copyRootCmd = rootCmdForTesting
	rootCmdX := &RootCommand{
		Command: Command{
			BaseOpt: BaseOpt{
				Name: "consul-tags",
			},
		},
	}

	commands := []string{
		"consul-tags --help -q",
	}

	for _, cc := range commands {
		os.Args = strings.Split(cc, " ")
		SetInternalOutputStreams(nil, nil)

		if err := Exec(rootCmdX); err != nil {
			t.Fatal(err)
		}
	}

	w := Worker()

	_ = w.genShellPowershell(os.Stdout, "", &rootCmdX.Command, nil)
}

// GenManualForCommandForTest _
func GenManualForCommandForTest(cmd *Command) (fn string, err error) {
	return genManualForCommand(cmd)
}

func TestWorkerAddIt(t *testing.T) {
	InternalResetWorkerForTest()
	ResetOptions()

	// copyRootCmd = rootCmdForTesting
	rootCmdX := rootCmdForAliasesTest()

	w := Worker()

	f, err := dir.TempFile("", "example")
	// f, err := os.CreateTemp("", "example") // go 1.17+ only
	if err != nil {
		t.Errorf("err: %v", err)
		return
	}
	defer os.Remove(f.Name())

	fi, err2 := f.Stat()
	if err2 != nil {
		t.Errorf("err: %v", err2)
		return
	}

	root := &rootCmdX.Command

	_ = w._addOnAddIt(fi, rootCmdX, ".")

	_ = w._addonAddCmd(root, "v", "v", &a1{}, &p1{})

	_ = w._addIt(fi, rootCmdX, ".")

	w._intFlgAdd(rootCmdX, "full", "full", "full", func(ff *Flag) {})
	w._intFlgAdd(rootCmdX, "full1", "full1", "", func(ff *Flag) {})
	w._intFlgAdd(rootCmdX, "full1", "full1", "", func(ff *Flag) {})
}

type p1 struct{}

func (p *p1) Name() string                      { return "vvv" }
func (p *p1) ShortName() string                 { return "v" }
func (p *p1) Aliases() []string                 { return []string{"v1"} }
func (p *p1) Description() string               { return "v" }
func (p *p1) SubCommands() []cmdrbase.PluginCmd { return []cmdrbase.PluginCmd{&p2{}} }
func (p *p1) Flags() []cmdrbase.PluginFlag      { return []cmdrbase.PluginFlag{&f1{}} }
func (p *p1) Action(args []string) (err error)  { return nil }

type p2 struct{}

func (p *p2) Name() string                      { return "vvv2" }
func (p *p2) ShortName() string                 { return "v2" }
func (p *p2) Aliases() []string                 { return []string{"v2"} }
func (p *p2) Description() string               { return "v2" }
func (p *p2) SubCommands() []cmdrbase.PluginCmd { return []cmdrbase.PluginCmd{} }
func (p *p2) Flags() []cmdrbase.PluginFlag      { return []cmdrbase.PluginFlag{&f1{}} }
func (p *p2) Action(args []string) (err error)  { return nil }

type f1 struct{}

func (f *f1) Name() string              { return "flag1" }
func (f *f1) ShortName() string         { return "f1" }
func (f *f1) Aliases() []string         { return []string{"f1f1"} }
func (f *f1) Description() string       { return "f1" }
func (f *f1) DefaultValue() interface{} { return true }
func (f *f1) PlaceHolder() string       { return "" }
func (f *f1) Action() (err error)       { return nil }

type a1 struct{}

func (a *a1) Name() string                      { return "addon1" } //nolint:goconst //keep it
func (a *a1) ShortName() string                 { return "a1" }
func (a *a1) Aliases() []string                 { return []string{"a1a1"} }
func (a *a1) Description() string               { return "addon1" }
func (a *a1) SubCommands() []cmdrbase.PluginCmd { return nil }
func (a *a1) Flags() []cmdrbase.PluginFlag      { return nil }
func (a *a1) Action(args []string) (err error)  { return nil }
func (a *a1) AddonTitle() string                { return "addon1" }
func (a *a1) AddonDescription() string          { return "addon1" }
func (a *a1) AddonCopyright() string            { return "addon1" }
func (a *a1) AddonVersion() string              { return "addon1" }

func TestWorkerHelpSystemPrint(t *testing.T) {
	InternalResetWorkerForTest()
	ResetOptions()

	// copyRootCmd = rootCmdForTesting
	rootCmdX := rootCmdForAliasesTest()

	commands := []string{
		"consul-tags --help -q",
	}

	for _, cc := range commands {
		os.Args = strings.Split(cc, " ")
		SetInternalOutputStreams(nil, nil)

		if err := Exec(rootCmdX); err != nil {
			t.Fatal(err)
		}
	}

	w := Worker()

	_ = w.helpSystemAction(&rootCmdX.Command, []string{})
	_ = w.helpSystemAction(&rootCmdX.Command, []string{"generate", "shell"})
	_ = w.helpSystemAction(&rootCmdX.Command, []string{"generate", "sh", "--zsh"})

	root := &rootCmdX.Command

	_ = w.helpSystemAction(root, []string{"generate", "r"})
	_ = w.helpSystemAction(root, []string{""})
	_ = w.helpSystemAction(root, []string{"generate", ""})
	_ = w.helpSystemAction(root, []string{"generate", "s"})
	_ = w.helpSystemAction(root, []string{"generate", "shell", ""})
	_ = w.helpSystemAction(root, []string{"generate", "sh", "-"})
	_ = w.helpSystemAction(root, []string{"generate", "shell", "--z"})
	_ = w.helpSystemAction(root, []string{"generate", "shell", "--zsh"})

	_ = DottedPathToCommand("generate.shell", nil)
	_ = DottedPathToCommand("version", root)
	_ = DottedPathToCommand("cc.cc1.cc1c1", root)

	_ = DottedPathToCommand("versions", root)

	_, _ = DottedPathToCommandOrFlag("generate.shell", nil)
	_, _ = DottedPathToCommandOrFlag("version", root)
	_, _ = DottedPathToCommandOrFlag("cc.cc1.cc1c1", root)
	_, _ = DottedPathToCommandOrFlag("generate.shell.zsh", nil)

	w.rootCommand.RunAsSubCommand = "generate.shell"
	w.preparePtPkg(&ptpkg{})
	w._setSwChars("windows")
}
