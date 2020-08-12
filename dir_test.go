/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr_test

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/tool"
	"github.com/hedzr/logex"
	"gopkg.in/hedzr/errors.v2"
	"os"
	"strings"
	"testing"
)

// TestIsDirectory tests more
//
// usage:
//   go test ./... -v -test.run '^TestIsDirectory$'
//
func TestIsDirectory(t *testing.T) {
	t.Logf("osargs[0] = %v", os.Args[0])
	t.Logf("InTesting: %v", tool.InTesting())
	t.Logf("InDebugging: %v", tool.InDebugging())

	cmdr.NormalizeDir("")

	if yes, err := cmdr.IsDirectory("./conf.d1"); yes {
		t.Fatal(err)
	}
	if yes, err := cmdr.IsDirectory("./ci"); !yes {
		t.Fatal(err)
	}
	if yes, err := cmdr.IsRegularFile("./doc1.golang"); yes {
		t.Fatal(err)
	}
	if yes, err := cmdr.IsRegularFile("./doc.go"); !yes {
		t.Fatal(err)
	}
}

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

func TestMatch(t *testing.T) {

	cmdr.ResetOptions()
	cmdr.InternalResetWorker()

	var err error
	var cmd *cmdr.Command
	var outX = bytes.NewBufferString("")
	var errX = bytes.NewBufferString("")
	var outBuf = bufio.NewWriterSize(outX, 16384)
	var errBuf = bufio.NewWriterSize(errX, 16384)
	cmdr.SetInternalOutputStreams(outBuf, errBuf)
	// cmdr.SetCustomShowVersion(nil)
	// cmdr.SetCustomShowBuildInfo(nil)

	copyRootCmd = rootCmdForTesting
	copyRootCmd.AppendPostActions(func(cmd *cmdr.Command, args []string) {
		//
	})
	if !copyRootCmd.IsRoot() {
		t.Fatal("IsRoot() test failed")
	}
	if copyRootCmd.GetRoot() != nil {
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
		cmdr.InternalResetWorker()
		resetFlagsAndLog(t)

		// cmdr.ShouldIgnoreWrongEnumValue = true

		println("xxx: ***: ", sss)
		w := cmdr.Worker3(rootCmdForTesting)
		w.AddOnAfterXrefBuilt(func(root *cmdr.RootCommand, args []string) {})
		w.AddOnBeforeXrefBuilding(func(root *cmdr.RootCommand, args []string) {})
		if cmd, err = cmdr.Match(sss,
			cmdr.WithUnhandledErrorHandler(onUnhandledErrorHandler),
			cmdr.WithNoCommandAction(true),
			cmdr.WithOnSwitchCharHit(func(parsed *cmdr.Command, switchChar string, args []string) (err error) {
				return
			}),
		); err != nil {
			if e, ok := err.(*cmdr.ErrorForCmdr); !ok || !e.Ignorable {
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

func TestHeadLike(t *testing.T) {

	cmdr.ResetOptions()
	cmdr.InternalResetWorker()

	var err error
	var outX = bytes.NewBufferString("")
	var errX = bytes.NewBufferString("")
	var outBuf = bufio.NewWriterSize(outX, 16384)
	var errBuf = bufio.NewWriterSize(errX, 16384)
	cmdr.SetInternalOutputStreams(outBuf, errBuf)
	// cmdr.SetCustomShowVersion(nil)
	// cmdr.SetCustomShowBuildInfo(nil)

	copyRootCmd = rootCmdForTesting

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

	t.Log("xxx: -------- loops for execTestings")
	for sss, verifier := range execTestingsHeadLike {
		cmdr.InternalResetWorker()
		resetFlagsAndLog(t)

		// cmdr.ShouldIgnoreWrongEnumValue = true

		println("xxx: ***: ", sss)
		w := cmdr.Worker2(true)
		w.AddOnAfterXrefBuilt(func(root *cmdr.RootCommand, args []string) {
		})
		w.AddOnBeforeXrefBuilding(func(root *cmdr.RootCommand, args []string) {
		})
		if _, err = w.InternalExecFor(rootCmdForTesting, strings.Split(sss, " ")); err != nil {
			var perr *cmdr.ErrorForCmdr
			if errors.As(err, &perr) && !perr.Ignorable {
				t.Fatal(err)
			}
		}

		// if sss == "consul-tags kv unknown" {
		// 	errX = bytes.NewBufferString("")
		// }

		if err = verifier(t, err); err != nil {
			t.Fatal(err)
		}
	}

}

func TestComplexOpt(t *testing.T) {
	defer logex.CaptureLog(t).Release()
	if tool.SavedOsArgs == nil {
		tool.SavedOsArgs = os.Args
	}
	defer func() {
		os.Args = tool.SavedOsArgs
	}()

	cmdr.ResetOptions()
	cmdr.InternalResetWorker()

	var err error
	// v1, v2 := 11, 0
	// var cmd *Command
	var rootCmdX = &cmdr.RootCommand{
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

	// cmd = &rootCmdX.Command
	var commands = []struct {
		line      string
		validator func(t *testing.T, err error) error
	}{
		{"consul-tags -cc 3.14159-2.56i", func(t *testing.T, err error) error {
			if cmdr.GetComplex128("app.complex") != 3.14159-2.56i {
				return errors.New("something wrong complex. |expected %v|got %v|", 3.14159-2.56i, cmdr.GetComplex128("app.complex"))
			}
			fmt.Println("consul-tags kv b ~ -------- no errors")
			return nil
		}},
		// {"consul-tags pa", func(t *testing.T, err error) error { return nil }},
	}
	for _, cc := range commands {
		os.Args = strings.Split(cc.line, " ")
		cmdr.SetInternalOutputStreams(nil, nil)
		cmdr.ResetOptions()
		if err = cmdr.Exec(rootCmdX); // cmdr.WithUnhandledErrorHandler(onUnhandleErrorHandler),
		// cmdr.WithOnSwitchCharHit(func(parsed *cmdr.Command, switchChar string, args []string) (err error) {
		// 	return
		// }),
		err != nil {
			t.Log(err) // hi, here is not real error occurs
		}
		if cc.validator != nil {
			err = cc.validator(t, err)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestTlideOptionsAndToggleGroupBranch(t *testing.T) {
	defer logex.CaptureLog(t).Release()
	if tool.SavedOsArgs == nil {
		tool.SavedOsArgs = os.Args
	}
	defer func() {
		os.Args = tool.SavedOsArgs
	}()

	cmdr.InternalResetWorker()
	cmdr.SetInternalOutputStreams(nil, nil)

	var err error
	var rootCmdX = &cmdr.RootCommand{
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

	// cmd = &rootCmdX.Command
	var commands = []struct {
		line      string
		validator func(t *testing.T, err error) error
	}{
		{"consul-tags -", nil},
		{"consul-tags --", nil},
		{"consul-tags --help ~~debug", nil},
		{"consul-tags --help ~~debug ~~env", nil},
		{"consul-tags --help ~~debug ~~raw", nil},
		{"consul-tags --help ~~debug ~~more", nil},
		// {"consul-tags pa", func(t *testing.T, err error) error { return nil }},
	}
	for _, cc := range commands {
		os.Args = strings.Split(cc.line, " ")
		cmdr.ResetOptions()
		if err = cmdr.Exec(rootCmdX,
			cmdr.WithOnSwitchCharHit(func(parsed *cmdr.Command, switchChar string, args []string) (err error) {
				return
			}),
		); err != nil {
			t.Log(err) // hi, here is not real error occurs
		}
		if cc.validator != nil {
			err = cc.validator(t, err)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestHandlerPassThru(t *testing.T) {
	defer logex.CaptureLog(t).Release()
	if tool.SavedOsArgs == nil {
		tool.SavedOsArgs = os.Args
	}
	defer func() {
		os.Args = tool.SavedOsArgs
	}()

	cmdr.InternalResetWorker()
	cmdr.SetInternalOutputStreams(nil, nil)

	var err error
	var rootCmdX = &cmdr.RootCommand{
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

	// cmd = &rootCmdX.Command
	var commands = []struct {
		line      string
		validator func(t *testing.T, err error) error
	}{
		{"consul-tags --help -- ~~debug", nil},
		{"consul-tags c1 --help", nil},
		{"consul-tags c1 c2 --help", nil},
		{"consul-tags c1 c3 --help", nil},
		{"consul-tags c2 --help", nil},
		// {"consul-tags --help ~~debug ~~env", nil},
		// {"consul-tags --help ~~debug ~~raw", nil},
		// {"consul-tags --help ~~debug ~~more", nil},
		// {"consul-tags pa", func(t *testing.T, err error) error { return nil }},
	}
	for _, cc := range commands {
		os.Args = strings.Split(cc.line, " ")
		cmdr.ResetOptions()
		if err = cmdr.Exec(rootCmdX,
			cmdr.WithOnPassThruCharHit(func(parsed *cmdr.Command, switchChar string, args []string) (err error) {
				fmt.Println(parsed.GetDottedNamePath())
				return
			}),
		); err != nil {
			t.Log(err) // hi, here is not real error occurs
		}

		if lastCommand, e := cmdr.Match(cc.line,
			cmdr.WithAfterArgsParsed(func(cmd *cmdr.Command, args []string) (err error) {
				return
			}),
		); lastCommand != nil || e != nil {
			t.Logf("lastCommand = %v, e = %v", lastCommand, e)
		}

		if cc.validator != nil {
			err = cc.validator(t, err)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

var (
	execTestingsMatch = map[string]func(t *testing.T, c *cmdr.Command, e error) error{
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
	execTestingsHeadLike = map[string]func(t *testing.T, e error) error{
		// enum test
		"consul-tags server -e oil": func(t *testing.T, e error) error {
			if strings.Index(e.Error(), "unexpected enumerable value") >= 0 {
				println("unexpected enumerable value found. This is a test, not an error.")
				return nil
			}
			return e
		},
		"consul-tags server -e orange": func(t *testing.T, e error) error {
			if cmdr.GetStringR("server.enum") != "orange" {
				println("unexpected enumerable value found. This is an error")
				return errors.New("unexpected enumerable value '%v' found. This is an error.",
					cmdr.GetStringR("server.enum"))
			}
			return e
		},

		"consul-tags server -1 -tt 3": func(t *testing.T, e error) error {
			if cmdr.GetIntR("server.retry") != 3 || cmdr.GetIntR("server.tail") != 1 {
				return fmt.Errorf("wrong: server.retry=%v(expect: %v), server.tail=%v(expect: %v)",
					cmdr.GetIntR("server.retry"), 3, cmdr.GetIntR("server.tail"), 1)
			}
			return nil
		},
		"consul-tags server -2 -t 5": func(t *testing.T, e error) error {
			if cmdr.GetIntR("retry") != 5 || cmdr.GetIntR("server.tail") != 2 {
				return errors.New("wrong: retry=%v(expect: %v), server.tail=%v(expect: %v)",
					cmdr.GetIntR("retry"), 5, cmdr.GetIntR("server.tail"), 2)
			}
			return nil
		},
		"consul-tags server -5": func(t *testing.T, e error) error {
			if cmdr.GetIntR("server.tail") != 5 {
				return fmt.Errorf("wrong: server.tail=%v(expect: %v)",
					cmdr.GetIntR("server.tail"), 5)
			}
			return nil
		},
		"consul-tags server -1973": func(t *testing.T, e error) error {
			if cmdr.GetIntR("server.tail") != 1973 {
				return errors.New("wrong: server.tail=%v (expect: %v)",
					cmdr.GetIntR("server.tail"), 1973)
			}
			return nil
		},
	}
)
