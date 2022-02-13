/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr_test

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/log/dir"
	"gopkg.in/yaml.v3"
	"strings"
	"testing"
	"time"
)

func TestCommandMethods(t *testing.T) {
	root := cmdr.Root("aa", "1.0.1").
		Header("sds")

	msCmd := cmdr.NewSubCmd().
		Titles("microservice", "ms").Name("ms").
		Description("", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		}).
		AttachTo(root)
	cmdr.NewSubCmd().
		Titles("list", "ls", "l", "lst", "dir").
		Description("list tags", "").
		Group("2333.List").
		Hidden(true).
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		}).
		AttachTo(msCmd)
	tagsSC := cmdr.NewSubCmd().
		Titles("tags", "t").
		Description("tags operations of a micro-service", "").
		Group("").
		AttachTo(msCmd)

	if cc := msCmd.ToCommand(); cc != nil {
		assertBool(len(cc.SubCommands) == 2, t, "want len(cc.SubCommands) == 2")
	}

	cmdr.NewSubCmd().
		Titles("list", "ls", "l", "lst", "dir").
		Description("list tags", "").
		Group("2333.List").
		Hidden(true).
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		}).
		AttachTo(tagsSC)

	assertBool(len(tagsSC.ToCommand().SubCommands) == 1, t, "want len(cc.SubCommands) == 2")

	xy := cmdr.NewSubCmd().
		Titles("xy-print", "xy").
		Description("test terminal control sequences", "test terminal control sequences,\nverbose long descriptions here.").
		Group("Test").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			fmt.Println("\x1b[2J") // clear screen

			for i, s := range args {
				fmt.Printf("\x1b[s\x1b[%d;%dH%s\x1b[u", 15+i, 30, s)
			}

			return
		}).
		AttachTo(root)
	cmdr.NewSubCmd().
		Titles("mx-test", "mx").
		Description("test new features", "test new features,\nverbose long descriptions here.").
		Group("001.Test").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			fmt.Printf("*** Got pp: %s\n", cmdr.GetString("app.mx-test.password"))
			fmt.Printf("*** Got msg: %s\n", cmdr.GetString("app.mx-test.message"))
			return
		}).
		AttachTo(root)

	cmd := xy.RootCommand().SubCommands[1]
	if cmd.GetHitStr() != "" {
		t.Failed()
	}
	if cmd.GetQuotedGroupName() != "Test" {
		t.Failed()
	}
	if cmd.GetSubCommandNamesBy(",") != "" {
		t.Failed()
	}

	cmd = xy.RootCommand().SubCommands[2]
	if cmd.GetQuotedGroupName() != "Test" {
		t.Failed()
	}

	cmd = xy.RootCommand().SubCommands[0]
	if cmd.GetSubCommandNamesBy(",") != "tags" {
		t.Failed()
	}
}

func assertBool(cond bool, t *testing.T, msgFailed ...string) {
	if cond == false {
		for _, msg := range msgFailed {
			t.Fatalf(msg)
		}
		t.Fatalf("cond NOT TRUE!")
	}
}

func TestFluentAPIDefault(t *testing.T) {
	root := cmdr.Root("aa", "1.0.1").
		Header("aa - test for cmdr - no version - hedzr").
		Copyright("s", "x")
	rootCmd1 := root.RootCommand()
	t.Log(rootCmd1)
	t.Log(root.ToCommand())

	root.AddOptFlag(cmdr.NewBool(false))
	root.AddFlag(cmdr.NewBool(false).ToFlag())
	root.AddOptCmd(cmdr.NewCmd())
	root.AddCommand(cmdr.NewCmd().ToCommand())

	optcmd := cmdr.NewCmd()

	optcmd.AttachToRoot(rootCmd1)
	optcmd.AttachTo(root)
	optcmd.AttachToCommand(&rootCmd1.Command)

	cmdr.NewBool(true).AttachToCommand(&rootCmd1.Command)
	cmdr.NewBool(true).AttachToRoot(rootCmd1)
	cmdr.NewBool(true).AttachTo(root)

	cmdr.NewCmdFrom(&rootCmd1.Command)
	cmdr.NewCmd()
	cmdr.NewSubCmd()
	cmdr.NewBool(false).Name("e")
	cmdr.NewDuration(0)
	cmdr.NewInt(0)
	cmdr.NewInt64(0)
	cmdr.NewIntSlice()
	cmdr.NewUintSlice()
	cmdr.NewString("")
	cmdr.NewStringSlice()
	cmdr.NewUint(0)
	cmdr.NewUint64(0)
	cmdr.NewFloat32(0)
	cmdr.NewFloat64(0)
	cmdr.NewComplex64(0)
	cmdr.NewComplex128(0)

}

func TestAsXXX(t *testing.T) {
	cmdr.AsYaml()
	cmdr.AsJSON()
	_, _ = cmdr.AsJSONExt(true)
	cmdr.AsToml()
	cmdr.GetHierarchyList()
	if _, err := cmdr.AsTomlExt(); err != nil {
		t.Logf("AsTomlExt error: %v", err)
	}
}

func TestKiloBytes(t *testing.T) {
	opts := cmdr.CurrentOptions()
	for _, str := range []string{
		"8K", "8M", "8G", "8T", "8P", "8E",
	} {
		//r:=str[1]
		n := opts.FromKilobytes(str)
		n = opts.FromKibiBytes(str)
		t.Log(n)
	}
}

func createRootOld() (rootOpt *cmdr.RootCmdOpt) {
	root := cmdr.Root("aa", "1.0.1").
		AddGlobalPreAction(func(cmd *cmdr.Command, args []string) (err error) {
			return
		}).
		AddGlobalPostAction(func(cmd *cmdr.Command, args []string) {
		}).
		RunAsSubCommand("generate.shell").
		Header("aa - test for cmdr - no version - hedzr").
		Copyright("s", "x")

	root.AppendPreActions(func(cmd *cmdr.Command, args []string) (err error) {
		return
	})
	root.AppendPostActions(func(cmd *cmdr.Command, args []string) {
	})

	// ms

	co1 := cmdr.NewSubCmd().
		Titles("micro-service-1", "ms-1").
		Short("ms-1").Long("micro-service-1").Aliases("goms-1").
		Examples(``).Hidden(false).Deprecated("").
		PreAction(nil).PostAction(nil).Action(nil).
		TailPlaceholder("").
		Description("", "")
	co1.AttachTo(root)

	ff2 := cmdr.NewInt(3).
		Titles("retry1x", "t1x").
		Description("1(2)3`FILE` usage", "").
		Group("").
		DefaultValue(false, "RETRY")

	ff2.AttachTo(co1)
	ff2.SetOwner(co1)
	ff2.SetOwner(nil)

	co := cmdr.NewSubCmd().
		Titles("micro-service", "ms").
		Short("ms").Long("micro-service").Aliases("goms").
		Examples(``).Hidden(false).Deprecated("").
		PreAction(nil).PostAction(nil).Action(nil).
		TailPlaceholder("").
		Description("", "").
		Group("").
		VendorHidden(false).
		AttachTo(root)

	co.OwnerCommand()
	co.SetOwner(root)
	co.RootCmdOpt()

	cmdr.NewUint().
		Titles("retry", "t").
		Short("tt").Long("retry-tt").Aliases("go-tt").
		Examples(``).Hidden(false).Deprecated("").
		Action(nil).
		ExternalTool(cmdr.ExternalToolEditor).
		ExternalTool(cmdr.ExternalToolPasswordInput).
		Description("", "").
		Group("").
		CompletionActionStr("").CompletionMutualExclusiveFlags("").
		CompletionPrerequisitesFlags("").CompletionJustOnce(false).
		CompletionCircuitBreak(false).DoubleTildeOnly(false).
		DefaultValue(uint(3), "RETRY").
		AttachTo(co).
		SetOwner(root)

	ff1 := cmdr.NewBool().
		Titles("retry1", "t1").
		Description("1(2)3", "").
		Group("").
		DefaultValue(false, "RETRY").
		AttachTo(co).
		OwnerCommand()
	ff1.SetOwner(co)
	ff1.SetOwner(nil)

	cmdr.NewInt().
		Titles("retry2", "t2").
		Description("", "").
		Group("").ToggleGroup("").
		VendorHidden(false).
		DefaultValue(3, "RETRY").
		AttachTo(co).
		RootCommand()

	cmdr.NewUint().
		Titles("retry3", "t3").
		Description("", "").
		Group("").
		DefaultValue(uint64(3), "RETRY").
		AttachTo(co)

	cmdr.NewInt64().
		Titles("retry4", "t4").
		Description("", "").
		Group("").
		DefaultValue(int64(3), "RETRY").
		AttachTo(co)

	cmdr.NewStringSlice().
		Titles("retry5", "t5").
		Description("", "").
		Group("").
		DefaultValue([]string{"a", "b"}, "RETRY").
		AttachTo(co)

	cmdr.NewIntSlice().
		Titles("retry6", "t6").
		Description("", "").
		Group("").
		DefaultValue([]int{1, 2, 3}, "RETRY").
		AttachTo(co)

	cmdr.NewDuration().
		Titles("retry7", "t7").
		Description("", "").
		Group("").
		DefaultValue(3, "RETRY").
		AttachTo(co)

	f0 := cmdr.NewFloat32().
		Titles("retry8", "t8").
		Description("", "").
		Group("").
		DefaultValue(3.14, "PI").
		AttachTo(co)
	f0.ToFlag().GetTitleZshFlagNamesArray()

	cmdr.NewFloat64().
		Titles("retry9", "t9").
		Description("", "").
		Group("").
		DefaultValue(3.14159265358979323846264338327950288419716939937510582097494459230781640628620899, "PI").
		AttachTo(co)

	cmdr.NewComplex64().
		Titles("retry10", "t10").
		Description("", "").
		Group("").
		DefaultValue(3.14, "PI").
		AttachTo(co)

	cmdr.NewComplex128().
		Titles("retry11", "t11").
		Description("", "").
		Group("").
		DefaultValue(3.14159265358979323846264338327950288419716939937510582097494459230781640628620899, "PI").
		AttachTo(co)

	cmdr.NewInt().
		Titles("head", "h").
		Description("", "").
		Group("").
		DefaultValue(1, "").
		HeadLike(true, 1, 8000).
		AttachTo(co)

	f1 := cmdr.NewString().
		Titles("ienum", "i").
		Description("", "").
		Group("").
		DefaultValue("", "").
		ValidArgs("apple", "banana", "orange").
		AttachTo(co)
	f2 := f1.ToFlag()
	f2.GetDottedNamePath()
	f2.Delete()
	f2.GetTitleZshFlagNamesArray()
	f2.GetTitleZshFlagShortName()
	f2.GetTitleZshNamesBy(",", true, true)
	(&cmdr.Flag{}).Delete()

	// ms tags

	cTags := cmdr.NewSubCmd().
		Titles("tags", "t").
		Description("", "").
		Group("").
		AttachTo(co)

	cmdr.NewString().
		Titles("addr", "a").
		Description("", "").
		Group("").
		DefaultValue("consul.ops.local", "ADDR").
		AttachTo(cTags)

	// ms tags ls

	cmdr.NewSubCmd().
		Titles("list", "ls").
		Description("", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		}).
		AttachTo(cTags)

	fn := func() {
		defer func() {
			if e := recover(); e != nil {
				print(e)
			}
		}()

		c1 := cmdr.NewSubCmd().
			Titles("", "").
			AttachTo(cTags)
		c1.ToCommand().GetName()
	}
	fn()

	c8 := cmdr.NewSubCmd().
		Titles("add1", "a1").
		Description("", "").
		Group("").
		AttachTo(cTags)

	c9 := cmdr.NewSubCmd().
		Titles("add", "").
		Description("", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		}).
		AttachTo(cTags)
	f91 := cmdr.NewString().
		Titles("ienum", "").Aliases("ie91").
		Description("", "").
		Group("").
		DefaultValue("", "").
		ValidArgs("apple", "banana", "orange").
		AttachTo(c9)
	f91.ToFlag().GetTitleZshFlagShortName()

	c91 := c9.ToCommand()
	c91.Match("")
	c91.Match("ie91")
	c91.GetTriggeredTimes()
	//c9.SetOwner(nil)
	c91.Delete()

	c8.SetOwner(nil)
	c8.ToCommand().Delete()

	cTags.ToCommand().Delete()

	return root
}

func createRoot() (rootOpt *cmdr.RootCmdOpt) {
	root := cmdr.Root("aa", "1.0.1").
		Header("aa - test for cmdr - no version - hedzr").
		Copyright("s", "x")

	// ms

	co := cmdr.NewSubCmd().
		Titles("micro-service", "ms").
		Short("ms").Long("micro-service").Aliases("goms").
		Examples(``).Hidden(false).Deprecated("").
		PreAction(nil).PostAction(nil).Action(nil).
		TailPlaceholder("").
		Description("", "").
		Group("").
		AttachTo(root)

	co.OwnerCommand()
	co.SetOwner(root)

	cmdr.NewUint(3).
		Titles("retry", "t").
		Short("tt").Long("retry-tt").Aliases("go-tt").
		Examples(``).Hidden(false).Deprecated("").
		Action(nil).
		ExternalTool(cmdr.ExternalToolEditor).
		ExternalTool(cmdr.ExternalToolPasswordInput).
		Description("", "").
		Group("").
		Placeholder("RETRY").
		AttachTo(co).
		SetOwner(root)

	cmdr.NewBool().
		Titles("retry1", "t1").
		Description("", "").
		Group("").
		Placeholder("RETRY").
		AttachTo(co).
		OwnerCommand()

	cmdr.NewInt(5).
		Titles("retry2", "t2").
		Description("", "").
		Group("").ToggleGroup("").
		AttachTo(co).
		RootCommand()

	cmdr.NewUint64(uint64(3)).
		Titles("retry3", "t3").
		Description("", "").
		Group("").
		AttachTo(co)

	cmdr.NewInt64(int64(3)).
		Titles("retry4", "t4").
		Description("", "").
		Group("").
		AttachTo(co)

	cmdr.NewStringSlice("a", "b").
		Titles("retry5", "t5").
		Description("", "").
		AttachTo(co)

	cmdr.NewIntSlice(1, 2, 3).
		Titles("retry6", "t6").
		Description("", "").
		AttachTo(co)

	cmdr.NewUintSlice(1, 2, 3).
		Titles("retry61", "t61").
		Description("", "").
		AttachTo(co)

	cmdr.NewDuration(time.Second).
		Titles("retry7", "t7").
		AttachTo(co)

	cmdr.NewFloat32(float32(3.14)).
		Titles("retry8", "t8").
		Description("", "").
		Group("").
		Placeholder("PI").
		AttachTo(co)

	cmdr.NewFloat64(3.14159265358979323846264338327950288419716939937510582097494459230781640628620899).
		Titles("retry9", "t9").
		Description("", "").
		Group("").
		Placeholder("PI").
		AttachTo(co)

	cmdr.NewComplex64(complex64(3.14+9i)).
		Titles("retry10", "t10").
		Description("", "").
		Group("").
		AttachTo(co)

	cmdr.NewComplex128(3.14+9i).
		Titles("retry11", "t11").
		Description("", "").
		Group("").
		AttachTo(co)

	cmdr.NewInt(1).
		Titles("head", "h").
		Description("", "").
		Group("").
		HeadLike(true, 1, 8000).
		EnvKeys("AVCX").
		Required().
		Required(false, false, true, false).
		AttachTo(co)

	cmdr.NewString("").
		Titles("ienum", "i").
		Description("", "").
		Group("").
		ValidArgs("apple", "banana", "orange").
		AttachTo(co)

	// ms tags

	cTags := cmdr.NewSubCmd().
		Titles("tags", "t").
		Description("", "").
		Group("").
		AttachTo(co)

	cmdr.NewString("consul.ops.local").
		Titles("addr", "a").
		Description("", "").
		Group("").
		Placeholder("ADDR").
		AttachTo(cTags)

	// ms tags ls

	cmdr.NewSubCmd().
		Titles("list", "ls").
		Description("", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		}).
		AttachTo(cTags)

	cmdr.NewSubCmd().
		Titles("add", "a").
		Description("", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		}).
		AttachTo(cTags)

	return root
}

func TestFluentAPINew(t *testing.T) {
	cmdr.ResetOptions()
	cmdr.InternalResetWorkerForTest()

	root := createRoot()
	rootCmd1 := root.RootCommand()
	t.Log(rootCmd1)

	if s, err := cmdr.AsYamlExt(); err != nil {
		t.Fatalf("AsYamlExt error: %v", err)
	} else {
		t.Log(s)
	}
	if s, err := cmdr.AsTomlExt(); err != nil {
		t.Fatalf("AsTomlExt error: %v", err)
	} else {
		t.Log(s)
	}
}

func TestFluentAPIOld(t *testing.T) {
	cmdr.ResetOptions()
	cmdr.InternalResetWorkerForTest()

	root := createRootOld()
	rootCmd1 := root.RootCommand()
	t.Log(rootCmd1)
}

func TestMergeWith(t *testing.T) {
	cmdr.Set("test.1", 8)
	cmdr.Set("test.deep.branch.1", "test")

	if cmdr.GetString("app.test.deep.branch.1") != "test" {
		t.Fatalf("err, expect 'test', but got '%v'", cmdr.GetString("app.test.deep.branch.1"))
	}

	var m = make(map[string]interface{})
	err := yaml.Unmarshal([]byte(`
app:
  test:
    1: 9
    deep:
      branch:
        1: test-ok
`), &m)
	if err != nil {
		t.Fatal(err)
	}

	err = cmdr.MergeWith(m)
	if err != nil {
		t.Fatal(err)
	}

	if cmdr.GetInt("app.test.1") != 9 {
		t.Fatalf("err, expect 9, but got %v", cmdr.GetInt("app.test.1"))
	}
	if cmdr.GetString("app.test.deep.branch.1") != "test-ok" {
		t.Fatalf("err, expect 'test-ok', but got '%v'", cmdr.GetString("app.test.deep.branch.1"))
	}
}

func TestDelete(t *testing.T) {
	cmdr.Set("test.1", 8)
	cmdr.Set("test.deep.branch.1", "test")

	if cmdr.GetString("app.test.deep.branch.1") != "test" {
		t.Fatalf("err, expect 'test', but got '%v'", cmdr.GetString("app.test.deep.branch.1"))
	}

	cmdr.DeleteKey("app.test.branch.1")

	if cmdr.HasKey("app.test.branch.1") {
		t.Fatalf("FAILED, expect key not found, but found: %v", cmdr.Get("app.test.branch.1"))
	}
}

func addDupFlags(root *cmdr.RootCmdOpt) {
	co := root.RootCommand().SubCommands[0]

	co.Flags = append(co.Flags, &cmdr.Flag{
		BaseOpt: cmdr.BaseOpt{
			Short:       "tt",
			Full:        "retry-tt",
			Aliases:     []string{"retry-tt"},
			Group:       "",
			Description: "",
		},
		DefaultValue: false,
	})
	co.Flags = append(co.Flags, &cmdr.Flag{
		BaseOpt: cmdr.BaseOpt{
			Name:        "retry-tt",
			Group:       "",
			Description: "",
		},
		DefaultValue: false,
	})
	co.Flags = append(co.Flags, &cmdr.Flag{
		BaseOpt: cmdr.BaseOpt{
			Name:        "retry-tt-not-dup",
			Group:       "",
			Description: "",
		},
		DefaultValue: false,
	})

	r := root.RootCommand()
	r.SubCommands = append(r.SubCommands, &cmdr.Command{
		BaseOpt: cmdr.BaseOpt{
			Short:           "ms",
			Full:            "micro-service",
			Aliases:         []string{"micro-service"},
			Group:           "",
			Description:     "",
			LongDescription: "",
		},
	})
	r.SubCommands = append(r.SubCommands, &cmdr.Command{
		BaseOpt: cmdr.BaseOpt{
			Name:            "micro-service",
			Short:           "",
			Full:            "",
			Group:           "",
			Description:     "",
			LongDescription: "",
		},
	})
}

func TestAlreadyUsed(t *testing.T) {
	cmdr.ResetOptions()
	cmdr.InternalResetWorkerForTest()

	root := createRoot()

	var err error
	deferFn := prepareConfD(t)
	outX, errX := prepareStreams()
	defer func() {

		x := outX.String()
		t.Logf("--------- stdout // %v // %v\n%v", dir.GetExecutableDir(), dir.GetExecutablePath(), x)

		if errX.Len() > 0 {
			t.Log("--------- stderr")
			t.Logf("Warn for normal err-info!! %v", errX.String())
		}

		resetOsArgs()
		deferFn()

	}()

	addDupFlags(root)

	t.Log("xxx: -------- loops for alreadyUsedTestings")
	for sss, verifier := range alreadyUsedTestings {
		resetFlagsAndLog(t)
		cmdr.ResetRootInWorkerForTest()

		t.Log("xxx: ***: ", sss)

		if _, err = cmdr.Worker().InternalExecFor(root.RootCommand(), strings.Split(sss, " ")); err != nil {
			t.Fatal(err, fmt.Sprintf("rootCmd = %p", root.RootCommand()))
		}
		if err = verifier(t); err != nil {
			t.Fatal(err)
		}
	}
}

var (
	// testing args
	alreadyUsedTestings = map[string]func(t *testing.T) error{
		// "consul-tags -qq": func(t *testing.T) error {
		// 	return nil
		// },
		"consul-tags --help": func(t *testing.T) error {
			return nil
		},
		// "consul-tags --help ~~debug": func(t *testing.T) error {
		// 	return nil
		// },
	}
)
