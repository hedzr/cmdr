/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr_test

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/tool"
	"github.com/hedzr/logex"
	"gopkg.in/hedzr/errors.v3"
)

// // TestIsDirectory tests more
// //
// // usage:
// //   go test ./... -v -test.run '^TestIsDirectory$'
// //
// func TestIsDirectory(t *testing.T) {
//	t.Logf("osargs[0] = %v", os.Args[0])
//	t.Logf("InTesting: %v", cmdr.InTesting())
//	t.Logf("InDebugging: %v", cmdr.InDebugging())
//
//	cmdr.NormalizeDir("")
//
//	if yes, err := cmdr.IsDirectory("./conf.d1"); yes {
//		t.Fatal(err)
//	}
//	if yes, err := cmdr.IsDirectory("./ci"); !yes {
//		t.Fatal(err)
//	}
//	if yes, err := cmdr.IsRegularFile("./doc1.golang"); yes {
//		t.Fatal(err)
//	}
//	if yes, err := cmdr.IsRegularFile("./doc.go"); !yes {
//		t.Fatal(err)
//	}
// }

func TestDumpers(t *testing.T) {
	if cmdr.DumpAsString() == "" {
		t.Fatal("fatal DumpAsString")
	}
}

func TestMatchPreQ(t *testing.T) {
	if len(strings.Split("server start ", " ")) != 3 {
		t.Fatal("expect 3")
	}
	if len(strings.Split("server start  ", " ")) != 4 {
		t.Fatal("expect 4")
	}

	t.Logf("%q", strings.Split("server start ", " "))
	t.Logf("%q", strings.Split("server start  ", " "))
}

//nolint:funlen //for test
func TestMatch(t *testing.T) {
	defer logex.CaptureLog(t).Release()

	cmdr.ResetOptions()
	cmdr.InternalResetWorkerForTest()

	var err error
	var cmd *cmdr.Command
	outX := bytes.NewBufferString("")
	errX := bytes.NewBufferString("")
	outBuf := bufio.NewWriterSize(outX, 16384)
	errBuf := bufio.NewWriterSize(errX, 16384)
	cmdr.SetInternalOutputStreams(outBuf, errBuf)
	// cmdr.SetCustomShowVersion(nil)
	// cmdr.SetCustomShowBuildInfo(nil)

	rc := rootCmdForTesting()
	rc.AppendPostActions(func(cmd *cmdr.Command, args []string) {
	})
	rc.AppendPreActions(func(cmd *cmdr.Command, args []string) (err error) {
		return
	})
	if !rc.IsRoot() {
		t.Fatal("IsRoot() test failed")
	}
	if rc.GetRoot() != nil {
		t.Log("GetRoot() has a non-nil value from last testing")
	}
	t.Log("root.Name")

	defer func() {
		postWorks(t)

		if errX.Len() > 0 {
			t.Log("--------- stderr")
			// t.Fatalf("Error!! %v", errX.String())
			t.Errorf("Error for testing (it might not be failed)!! %v", errX.String())
		}

		resetOsArgs()

		// x := outX.String()
		// t.Logf("--------- stdout // %v // %v\n%v", cmdr.GetExecutableDir(), cmdr.GetExcutablePath(), x)
	}()

	onUnhandledErrorHandler := func(err interface{}) {
		t.Fatal(errors.DumpStacksAsString(false))
	}

	t.Log("xxx: -------- loops for execTestingsMatch")
	for sss, verifier := range execTestingsMatch {
		cmdr.InternalResetWorkerForTest()
		// resetFlagsAndLog(t)
		cmdr.ResetOptions()

		// cmdr.ShouldIgnoreWrongEnumValue = true

		println("xxx: ***: ", sss)
		w := cmdr.Worker3(rc)
		w.AddOnAfterXrefBuilt(func(root *cmdr.RootCommand, args []string) {})
		w.AddOnBeforeXrefBuilding(func(root *cmdr.RootCommand, args []string) {})
		if cmd, err = cmdr.MatchForTest(sss,
			cmdr.WithUnhandledErrorHandler(onUnhandledErrorHandler),
			cmdr.WithNoCommandAction(true),
			cmdr.WithOnSwitchCharHit(func(parsed *cmdr.Command, switchChar string, args []string) (err error) {
				return
			}),
		); err != nil {
			if e, ok := err.(*cmdr.ErrorForCmdr); !ok || !e.Ignorable { //nolint:errorlint //like it
				t.Fatal(err)
			}
		}

		// if sss == "consul-tags kv unknown" {
		// 	errX = bytes.NewBufferString("")
		// }

		t.Logf("xxx: matched: cmd = %v, err = %v", cmd, err)
		if err = verifier(t, cmd, err); err != nil {
			t.Fatal(err)
		}
	}
}

func TestHeadLike1(t *testing.T) {
	testFramework(t, rootCmdForTesting, execTestingsHeadLike)
}

func TestWrongEnum(t *testing.T) {
	testFramework(t, rootCmdForTesting, execTestingsWronEnum)
}

// func TestHeadLike(t *testing.T) {
//
//	cmdr.ResetOptions()
//	cmdr.InternalResetWorkerForTest()
//
//	var err error
//	var outX = bytes.NewBufferString("")
//	var errX = bytes.NewBufferString("")
//	var outBuf = bufio.NewWriterSize(outX, 16384)
//	var errBuf = bufio.NewWriterSize(errX, 16384)
//	cmdr.SetInternalOutputStreams(outBuf, errBuf)
//	// cmdr.SetCustomShowVersion(nil)
//	// cmdr.SetCustomShowBuildInfo(nil)
//
//	copyRootCmd = rootCmdForTesting
//
//	defer func() {
//
//		postWorks(t)
//
//		if errX.Len() > 0 {
//			t.Log("--------- stderr")
//			// t.Fatalf("Error!! %v", errX.String())
//			t.Errorf("Error for testing (it might not be failed)!! %v", errX.String())
//		}
//
//		resetOsArgs()
//
//		// x := outX.String()
//		// t.Logf("--------- stdout // %v // %v\n%v", cmdr.GetExecutableDir(), cmdr.GetExcutablePath(), x)
//	}()
//
//	t.Log("xxx: -------- loops for execTestings")
//	for sss, verifier := range execTestingsHeadLike {
//		cmdr.InternalResetWorkerForTest()
//		//resetFlagsAndLog(t)
//		cmdr.ResetOptions()
//
//		// cmdr.ShouldIgnoreWrongEnumValue = true
//
//		println("xxx: ***: ", sss)
//		// w := cmdr.Worker2(sss != "consul-tags server -e oil")
//		w := cmdr.Worker2(true)
//		w.AddOnAfterXrefBuilt(func(root *cmdr.RootCommand, args []string) {
//		})
//		w.AddOnBeforeXrefBuilding(func(root *cmdr.RootCommand, args []string) {
//		})
//
//		if sss == "consul-tags server -e oil" {
//			fmt.Println()
//		}
//		var cc *cmdr.Command
//		if cc, err = w.InternalExecFor(rootCmdForTesting, strings.Split(sss, " ")); err != nil {
//			var perr *cmdr.ErrorForCmdr
//			if errors.As(err, &perr) && !perr.Ignorable {
//				t.Fatal(err)
//			}
//		}
//
//		// if sss == "consul-tags kv unknown" {
//		// 	errX = bytes.NewBufferString("")
//		// }
//
//		if err = verifier(t, cc, err); err != nil {
//			var perr *cmdr.ErrorForCmdr
//			if errors.As(err, &perr) && !perr.Ignorable {
//				t.Fatal(err)
//			} else {
//				t.Logf("[Warn] error occurs: %v", err)
//			}
//		}
//	}
//
// }

func TestComplexOpt1(t *testing.T) {
	rootCmdX := func() *cmdr.RootCommand {
		return &cmdr.RootCommand{
			Command: cmdr.Command{
				BaseOpt: cmdr.BaseOpt{
					Name: "consul-tags",
				},
				Flags: []*cmdr.Flag{
					{
						BaseOpt: cmdr.BaseOpt{
							Short: "cc", Full: "complex",
						},
						DefaultValue: complex128(0),
					},
				},
			},
		}
	}
	commands := testCases{
		"consul-tags -cc 3.14159-2.56i": func(t *testing.T, cmd *cmdr.Command, err error) error {
			if cmdr.GetComplex128("app.complex") != 3.14159-2.56i {
				return errors.New("something wrong complex. |expected %v|got %v|", 3.14159-2.56i, cmdr.GetComplex128("app.complex"))
			}
			fmt.Println("consul-tags kv b ~ -------- no errors")
			return nil
		},
		// {"consul-tags pa", func(t *testing.T, err error) error { return nil }},
	}
	testFramework(t, rootCmdX, commands)
}

// func TestComplexOpt(t *testing.T) {
//	defer logex.CaptureLog(t).Release()
//	if tool.SavedOsArgs == nil {
//		tool.SavedOsArgs = os.Args
//	}
//	defer func() {
//		os.Args = tool.SavedOsArgs
//	}()
//
//	cmdr.ResetOptions()
//	cmdr.InternalResetWorkerForTest()
//
//	var err error
//	// v1, v2 := 11, 0
//	// var cmd *Command
//	var rootCmdX = &cmdr.RootCommand{
//		Command: cmdr.Command{
//			BaseOpt: cmdr.BaseOpt{
//				Name: "consul-tags",
//			},
//			Flags: []*cmdr.Flag{
//				{
//					BaseOpt: cmdr.BaseOpt{
//						Short: "cc", Full: "complex",
//					},
//					DefaultValue: complex128(0),
//				},
//			},
//		},
//	}
//
//	// cmd = &rootCmdX.Command
//	var commands = []struct {
//		line      string
//		validator func(t *testing.T, err error) error
//	}{
//		{"consul-tags -cc 3.14159-2.56i", func(t *testing.T, err error) error {
//			if cmdr.GetComplex128("app.complex") != 3.14159-2.56i {
//				return errors.New("something wrong complex. |expected %v|got %v|", 3.14159-2.56i, cmdr.GetComplex128("app.complex"))
//			}
//			fmt.Println("consul-tags kv b ~ -------- no errors")
//			return nil
//		}},
//		// {"consul-tags pa", func(t *testing.T, err error) error { return nil }},
//	}
//	for _, cc := range commands {
//		os.Args = strings.Split(cc.line, " ")
//		cmdr.SetInternalOutputStreams(nil, nil)
//		cmdr.ResetOptions()
//		if err = cmdr.Exec(rootCmdX); // cmdr.WithUnhandledErrorHandler(onUnhandleErrorHandler),
//		// cmdr.WithOnSwitchCharHit(func(parsed *cmdr.Command, switchChar string, args []string) (err error) {
//		// 	return
//		// }),
//		err != nil {
//			t.Log(err) // hi, here is not real error occurs
//		}
//		if cc.validator != nil {
//			err = cc.validator(t, err)
//			if err != nil {
//				t.Fatal(err)
//			}
//		}
//	}
// }

func TestTildeOptionsAndToggleGroupBranch1(t *testing.T) {
	rootCmdX := func() *cmdr.RootCommand {
		return &cmdr.RootCommand{
			Command: cmdr.Command{
				BaseOpt: cmdr.BaseOpt{
					Name: "consul-tags",
				},
				Flags: []*cmdr.Flag{
					{
						BaseOpt: cmdr.BaseOpt{
							Short: "ff", Full: "float",
						},
						DefaultValue: float32(0),
					},
					{
						BaseOpt: cmdr.BaseOpt{
							Short: "tg1", Full: "tg1",
						},
						ToggleGroup: "ToggleGroup",
					},
					{
						BaseOpt: cmdr.BaseOpt{
							Short: "tg2", Full: "tg2",
						},
						ToggleGroup: "ToggleGroup",
					},
					{
						BaseOpt: cmdr.BaseOpt{
							Short: "tg3", Full: "tg3",
						},
						DefaultValue: true,
						ToggleGroup:  "ToggleGroup",
					},
				},
			},
		}
	}

	// cmd = &rootCmdX.Command
	commands := testCases{
		"consul-tags -":                     nil,
		"consul-tags --":                    nil,
		"consul-tags --help ~~debug":        nil,
		"consul-tags --help ~~debug ~~env":  nil,
		"consul-tags --help ~~debug ~~raw":  nil,
		"consul-tags --help ~~debug ~~more": nil,
		// {"consul-tags pa", func(t *testing.T, err error) error { return nil }},
	}
	testFramework(t, rootCmdX, commands)
}

// func TestTildeOptionsAndToggleGroupBranch(t *testing.T) {
//	defer logex.CaptureLog(t).Release()
//	if tool.SavedOsArgs == nil {
//		tool.SavedOsArgs = os.Args
//	}
//	defer func() {
//		os.Args = tool.SavedOsArgs
//	}()
//
//	cmdr.InternalResetWorkerForTest()
//	cmdr.SetInternalOutputStreams(nil, nil)
//
//	var err error
//	var rootCmdX = &cmdr.RootCommand{
//		Command: cmdr.Command{
//			BaseOpt: cmdr.BaseOpt{
//				Name: "consul-tags",
//			},
//			Flags: []*cmdr.Flag{
//				{
//					BaseOpt: cmdr.BaseOpt{
//						Short: "ff", Full: "float",
//					},
//					DefaultValue: float32(0),
//				},
//				{
//					BaseOpt: cmdr.BaseOpt{
//						Short: "tg1", Full: "tg1",
//					},
//					ToggleGroup: "ToggleGroup",
//				},
//				{
//					BaseOpt: cmdr.BaseOpt{
//						Short: "tg2", Full: "tg2",
//					},
//					ToggleGroup: "ToggleGroup",
//				},
//				{
//					BaseOpt: cmdr.BaseOpt{
//						Short: "tg3", Full: "tg3",
//					},
//					DefaultValue: true,
//					ToggleGroup:  "ToggleGroup",
//				},
//			},
//		},
//	}
//
//	// cmd = &rootCmdX.Command
//	var commands = []struct {
//		line      string
//		validator func(t *testing.T, err error) error
//	}{
//		{"consul-tags -", nil},
//		{"consul-tags --", nil},
//		{"consul-tags --help ~~debug", nil},
//		{"consul-tags --help ~~debug ~~env", nil},
//		{"consul-tags --help ~~debug ~~raw", nil},
//		{"consul-tags --help ~~debug ~~more", nil},
//		// {"consul-tags pa", func(t *testing.T, err error) error { return nil }},
//	}
//	for _, cc := range commands {
//		os.Args = strings.Split(cc.line, " ")
//		cmdr.ResetOptions()
//		if err = cmdr.Exec(rootCmdX,
//			cmdr.WithOnSwitchCharHit(func(parsed *cmdr.Command, switchChar string, args []string) (err error) {
//				return
//			}),
//		); err != nil {
//			t.Log(err) // hi, here is not real error occurs
//		}
//		if cc.validator != nil {
//			err = cc.validator(t, err)
//			if err != nil {
//				t.Fatal(err)
//			}
//		}
//	}
// }

func TestHandlerPassThru1(t *testing.T) {
	rootCmdX := func() *cmdr.RootCommand {
		return &cmdr.RootCommand{
			Command: cmdr.Command{
				BaseOpt: cmdr.BaseOpt{
					Name: "consul-tags",
				},
				Flags: []*cmdr.Flag{
					{
						BaseOpt: cmdr.BaseOpt{
							Short: "ff", Full: "float",
						},
						DefaultValue: float32(0),
					},
				},
				SubCommands: []*cmdr.Command{
					{
						BaseOpt: cmdr.BaseOpt{
							Full: "c1", Short: "c1",
						},
						SubCommands: []*cmdr.Command{
							{
								BaseOpt: cmdr.BaseOpt{
									Full: "c1", Short: "c1",
								},
							},
							{
								BaseOpt: cmdr.BaseOpt{
									Full: "c2", Short: "c2",
								},
							},
						},
					},
					{
						BaseOpt: cmdr.BaseOpt{
							Full: "c2", Short: "c2",
						},
					},
				},
			},
		}
	}

	// cmd = &rootCmdX.Command
	commands := testCases{
		"consul-tags --help -- ~~debug": nil,
		"consul-tags c1 --help":         nil,
		"consul-tags c1 c2 --help":      nil,
		"consul-tags c1 c3 --help":      nil,
		"consul-tags c2 --help":         nil,
		// {"consul-tags --help ~~debug ~~env", nil},
		// {"consul-tags --help ~~debug ~~raw", nil},
		// {"consul-tags --help ~~debug ~~more", nil},
		// {"consul-tags pa", func(t *testing.T, err error) error { return nil }},
	}
	testFramework(t, rootCmdX, commands)
}

// func TestHandlerPassThru(t *testing.T) {
//	defer logex.CaptureLog(t).Release()
//	if tool.SavedOsArgs == nil {
//		tool.SavedOsArgs = os.Args
//	}
//	defer func() {
//		os.Args = tool.SavedOsArgs
//	}()
//
//	cmdr.InternalResetWorkerForTest()
//	cmdr.SetInternalOutputStreams(nil, nil)
//
//	var err error
//	var rootCmdX = &cmdr.RootCommand{
//		Command: cmdr.Command{
//			BaseOpt: cmdr.BaseOpt{
//				Name: "consul-tags",
//			},
//			Flags: []*cmdr.Flag{
//				{
//					BaseOpt: cmdr.BaseOpt{
//						Short: "ff", Full: "float",
//					},
//					DefaultValue: float32(0),
//				},
//			},
//			SubCommands: []*cmdr.Command{
//				{
//					BaseOpt: cmdr.BaseOpt{
//						Full: "c1", Short: "c1",
//					},
//					SubCommands: []*cmdr.Command{
//						{
//							BaseOpt: cmdr.BaseOpt{
//								Full: "c1", Short: "c1",
//							},
//						},
//						{
//							BaseOpt: cmdr.BaseOpt{
//								Full: "c2", Short: "c2",
//							},
//						},
//					},
//				},
//				{
//					BaseOpt: cmdr.BaseOpt{
//						Full: "c2", Short: "c2",
//					},
//				},
//			},
//		},
//	}
//
//	// cmd = &rootCmdX.Command
//	var commands = []struct {
//		line      string
//		validator func(t *testing.T, err error) error
//	}{
//		{"consul-tags --help -- ~~debug", nil},
//		{"consul-tags c1 --help", nil},
//		{"consul-tags c1 c2 --help", nil},
//		{"consul-tags c1 c3 --help", nil},
//		{"consul-tags c2 --help", nil},
//		// {"consul-tags --help ~~debug ~~env", nil},
//		// {"consul-tags --help ~~debug ~~raw", nil},
//		// {"consul-tags --help ~~debug ~~more", nil},
//		// {"consul-tags pa", func(t *testing.T, err error) error { return nil }},
//	}
//	for _, cc := range commands {
//		os.Args = strings.Split(cc.line, " ")
//		cmdr.ResetOptions()
//		if err = cmdr.Exec(rootCmdX,
//			cmdr.WithOnPassThruCharHit(func(parsed *cmdr.Command, switchChar string, args []string) (err error) {
//				fmt.Println(parsed.GetDottedNamePath())
//				return
//			}),
//		); err != nil {
//			t.Log(err) // hi, here is not real error occurs
//		}
//
//		if lastCommand, e := cmdr.MatchForTest(cc.line,
//			cmdr.WithAfterArgsParsed(func(cmd *cmdr.Command, args []string) (err error) {
//				return
//			}),
//		); lastCommand != nil || e != nil {
//			t.Logf("lastCommand = %v, e = %v", lastCommand, e)
//		}
//
//		if cc.validator != nil {
//			err = cc.validator(t, err)
//			if err != nil {
//				t.Fatal(err)
//			}
//		}
//	}
// }

//nolint:funlen //for test
func testFramework(t *testing.T, rootCommand func() *cmdr.RootCommand, cases testCases, opts ...cmdr.ExecOption) {
	defer logex.CaptureLog(t).Release()

	deferFn := prepareConfD(t)
	defer deferFn()

	if tool.SavedOsArgs == nil {
		tool.SavedOsArgs = os.Args
	}
	defer func() { os.Args = tool.SavedOsArgs }()

	// cmdr.ResetOptions()
	// cmdr.InternalResetWorkerForTest()

	var err error
	var cmd *cmdr.Command
	outX := bytes.NewBufferString("")
	errX := bytes.NewBufferString("")
	outBuf := bufio.NewWriterSize(outX, 32768)
	errBuf := bufio.NewWriterSize(errX, 16384)
	cmdr.SetInternalOutputStreams(outBuf, errBuf)
	// cmdr.SetCustomShowVersion(nil)
	// cmdr.SetCustomShowBuildInfo(nil)

	var copyOfRootCommand *cmdr.RootCommand
	resetRootCommand := func() { // not yet
		// copyOfRootCommand = new(cmdr.RootCommand)
		// if err := cmdr.CloneViaGob(copyOfRootCommand, rootCommand); err != nil {
		//	t.Fatal(err)
		// }
		copyOfRootCommand = rootCommand()
	}
	t.Log("root.Name")

	defer func() {
		postWorks(t)

		if errBuf.Buffered() > 0 {
			t.Log("--------- stderr")
			// t.Fatalf("Error!! %v", errX.String())
			_ = errBuf.Flush()
			t.Errorf("Error for testing (it might not be failed)!! %v", errX.String())
		}

		resetOsArgs()

		// x := outX.String()
		// t.Logf("--------- stdout // %v // %v\n%v", cmdr.GetExecutableDir(), cmdr.GetExcutablePath(), x)
	}()

	onUnhandledErrorHandler := func(err interface{}) { t.Fatal(errors.DumpStacksAsString(false)) }

	t.Log("xxx: -------- loops for TestGenerateShellCommand")
	for sss, verifier := range cases {
		cmdr.InternalResetWorkerForTest()
		cmdr.ResetOptions()
		resetRootCommand()

		// cmdr.ShouldIgnoreWrongEnumValue = true

		println("xxx: ***: ", sss)

		w := cmdr.Worker3(copyOfRootCommand)
		// w.AddOnAfterXrefBuilt(func(root *cmdr.RootCommand, args []string) {})
		// w.AddOnBeforeXrefBuilding(func(root *cmdr.RootCommand, args []string) {})
		cmdr.WithUnhandledErrorHandler(onUnhandledErrorHandler)(w)
		cmdr.WithInternalOutputStreams(outBuf, errBuf)(w)
		cmdr.WithNoWatchConfigFiles(true)(w)
		// cmdr.WithIgnoreWrongEnumValue(true)(w)
		for _, opt := range opts {
			opt(w)
		}

		args := strings.Split(sss, " ")

		cmd, err = w.InternalExecFor(copyOfRootCommand, args)

		// if sss == "consul-tags kv unknown" {
		// 	errX = bytes.NewBufferString("")
		// }

		t.Logf("xxx: matched: cmd = %v, err = %v", cmd, err)
		if verifier != nil {
			if err = verifier(t, cmd, err); !cmdr.IsIgnorableError(err) {
				t.Fatal(err)
			}
		}
	}
}

type testCases map[string]func(t *testing.T, c *cmdr.Command, e error) (err error)

var (
	execTestingsMatch = testCases{
		"server star": func(t *testing.T, c *cmdr.Command, e error) error {
			return e
		},
		"server start": func(t *testing.T, c *cmdr.Command, e error) error {
			return e
		},
		"server start ": func(t *testing.T, c *cmdr.Command, e error) error {
			return e
		},
		"server start -": func(t *testing.T, c *cmdr.Command, e error) error {
			return e
		},
	}

	// testing args
	execTestingsWronEnum = testCases{
		// enum test
		"consul-tags server -e oil": func(t *testing.T, c *cmdr.Command, e error) (err error) {
			if e == nil {
				return errors.New("expecting an 'unexpected enumerable value' exception threw for command-line: 'consul-tags server -e oil'")
			}
			if strings.Contains(e.Error(), "unexpected enumerable value") {
				println("unexpected enumerable value found. This is a test, not an error.")
				return nil
			}
			return e
		},
	}

	execTestingsHeadLike = testCases{
		"consul-tags server -e orange": func(t *testing.T, c *cmdr.Command, e error) (err error) {
			if cmdr.GetStringR("server.enum") != "orange" {
				println("unexpected enumerable value found. This is an error")
				return errors.New("unexpected enumerable value '%v' found. This is an error",
					cmdr.GetStringR("server.enum"))
			}
			return e
		},

		"consul-tags server -1 -tt 3": func(t *testing.T, c *cmdr.Command, e error) (err error) {
			if cmdr.GetIntR("server.retry") != 3 || cmdr.GetIntR("server.tail") != 1 {
				return fmt.Errorf("wrong: server.retry=%v(expect: %v), server.tail=%v(expect: %v)",
					cmdr.GetIntR("server.retry"), 3, cmdr.GetIntR("server.tail"), 1)
			}
			return nil
		},
		"consul-tags server -2 -t 5": func(t *testing.T, c *cmdr.Command, e error) (err error) {
			if cmdr.GetIntR("retry") != 5 || cmdr.GetIntR("server.tail") != 2 {
				return errors.New("wrong: retry=%v(expect: %v), server.tail=%v(expect: %v)",
					cmdr.GetIntR("retry"), 5, cmdr.GetIntR("server.tail"), 2)
			}
			return nil
		},
		"consul-tags server -5": func(t *testing.T, c *cmdr.Command, e error) (err error) {
			if cmdr.GetIntR("server.tail") != 5 {
				return fmt.Errorf("wrong: server.tail=%v(expect: %v)",
					cmdr.GetIntR("server.tail"), 5)
			}
			return nil
		},
		"consul-tags server -1973": func(t *testing.T, c *cmdr.Command, e error) (err error) {
			if cmdr.GetIntR("server.tail") != 1973 {
				return errors.New("wrong: server.tail=%v (expect: %v)",
					cmdr.GetIntR("server.tail"), 1973)
			}
			return nil
		},
	}
)
