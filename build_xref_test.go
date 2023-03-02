// Copyright © 2020 Hedzr Yeh.

package cmdr_test

import (
	"os"
	"strings"
	"testing"

	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/conf"
	"github.com/hedzr/cmdr/tool"
	"github.com/hedzr/log/dir"
	"github.com/hedzr/logex"
	"gopkg.in/hedzr/errors.v3"
)

// TestAddonPlugin for pluggable addons/plugins
func TestAddonPlugin(t *testing.T) {
	conf.AppName = "fluent"
	t.Logf("    - pwd=%q", dir.GetCurrentDir())

	defer logex.CaptureLog(t).Release()
	if tool.SavedOsArgs == nil {
		tool.SavedOsArgs = os.Args
	}
	defer func() {
		os.Args = tool.SavedOsArgs
	}()

	cmdr.ResetOptions()
	cmdr.InternalResetWorkerForTest()

	var err error
	// v1, v2 := 11, 0
	// var cmd *Command
	rootCmdX := &cmdr.RootCommand{
		Command: cmdr.Command{
			BaseOpt: cmdr.BaseOpt{
				Name: "flags",
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
	commands := []struct { //nolint:govet //just a testcase
		line      string
		validator func(t *testing.T, err error) error
	}{
		{"flags dx", func(t *testing.T, err error) error {
			if err != nil {
				t.Logf("flags dx ~ -------- has warned error: %v", err)
			} else {
				t.Logf("flags dx ~ -------- no errors")
			}
			return nil
		}},
	}
	for _, cc := range commands {
		os.Args = strings.Split(cc.line, " ")
		cmdr.SetInternalOutputStreams(nil, nil)
		cmdr.ResetOptions()
		if err = cmdr.Exec(rootCmdX, cmdr.WithNoWatchConfigFiles(true)); // cmdr.WithUnhandledErrorHandler(onUnhandleErrorHandler),
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

// TestPE for pluggable extensions
func TestPE(t *testing.T) {
	conf.AppName = "flags"

	defer logex.CaptureLog(t).Release()
	if tool.SavedOsArgs == nil {
		tool.SavedOsArgs = os.Args
	}
	defer func() {
		os.Args = tool.SavedOsArgs
	}()

	cmdr.ResetOptions()
	cmdr.InternalResetWorkerForTest()

	var err error
	// v1, v2 := 11, 0
	// var cmd *Command
	rootCmdX := &cmdr.RootCommand{
		Command: cmdr.Command{
			BaseOpt: cmdr.BaseOpt{
				Name: "flags",
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
	commands := []struct { //nolint:govet //just a testcase
		line      string
		validator func(t *testing.T, err error) error
	}{
		{"flags -cc 3.14159-2.56i", func(t *testing.T, err error) error {
			if cmdr.GetComplex128("app.complex") != 3.14159-2.56i {
				return errors.New("something wrong complex. |expected %v|got %v|", 3.14159-2.56i, cmdr.GetComplex128("app.complex"))
			}
			t.Logf("flags -cc PI ~ -------- no errors")
			return nil
		}},
		{"flags cpu", func(t *testing.T, err error) error {
			if err != nil {
				t.Logf("flags cpu ~ -------- has warned error: %v", err)
			} else {
				t.Logf("flags cpu ~ -------- no errors")
			}
			return nil
		}},
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
