/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr_test

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/errors"
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
	t.Logf("InTesting: %v", cmdr.InTesting())
	t.Logf("InDebugging: %v", cmdr.InDebugging())

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

	t.Log("xxx: -------- loops for execTestingsMatch")
	for sss, verifier := range execTestingsMatch {
		cmdr.InternalResetWorker()
		resetFlagsAndLog(t)

		// cmdr.ShouldIgnoreWrongEnumValue = true

		println("xxx: ***: ", sss)
		w := cmdr.Worker3(rootCmdForTesting)
		w.AddOnAfterXrefBuilt(func(root *cmdr.RootCommand, args []string) {})
		w.AddOnBeforeXrefBuilding(func(root *cmdr.RootCommand, args []string) {})
		if cmd, err = cmdr.Match(sss); err != nil {
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
			if e, ok := err.(*cmdr.ErrorForCmdr); !ok || !e.Ignorable {
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
