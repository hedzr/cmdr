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

	msCmd := root.NewSubCommand().
		Titles("microservice", "ms").Name("ms").
		Description("", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})
	msCmd.NewSubCommand().
		Titles("list", "ls", "l", "lst", "dir").
		Description("list tags", "").
		Group("2333.List").
		Hidden(true).
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})
	msCmd.NewSubCommand().
		Titles("tags", "t").
		Description("tags operations of a micro-service", "").
		Group("")

	xy := root.NewSubCommand().
		Titles("xy-print", "xy").
		Description("test terminal control sequences", "test terminal control sequences,\nverbose long descriptions here.").
		Group("Test").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			fmt.Println("\x1b[2J") // clear screen

			for i, s := range args {
				fmt.Printf("\x1b[s\x1b[%d;%dH%s\x1b[u", 15+i, 30, s)
			}

			return
		})
	root.NewSubCommand().
		Titles("mx-test", "mx").
		Description("test new features", "test new features,\nverbose long descriptions here.").
		Group("001.Test").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			fmt.Printf("*** Got pp: %s\n", cmdr.GetString("app.mx-test.password"))
			fmt.Printf("*** Got msg: %s\n", cmdr.GetString("app.mx-test.message"))
			return
		})

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
	cmdr.AsToml()
	cmdr.GetHierarchyList()
	if _, err := cmdr.AsTomlExt(); err != nil {
		t.Logf("AsTomlExt error: %v", err)
	}
}

func createRootOld() (rootOpt *cmdr.RootCmdOpt) {
	root := cmdr.Root("aa", "1.0.1").
		AddGlobalPreAction(func(cmd *cmdr.Command, args []string) (err error) {
			return
		}).AddGlobalPostAction(func(cmd *cmdr.Command, args []string) {
	}).
		Header("aa - test for cmdr - no version - hedzr").
		Copyright("s", "x")

	// ms

	co := root.NewSubCommand().
		Titles("micro-service", "ms").
		Short("ms").Long("micro-service").Aliases("goms").
		Examples(``).Hidden(false).Deprecated("").
		PreAction(nil).PostAction(nil).Action(nil).
		TailPlaceholder("").
		Description("", "").
		Group("").
		VendorHidden(false)

	co.OwnerCommand()
	co.SetOwner(root)
	co.RootCmdOpt()

	co.NewFlag(cmdr.OptFlagTypeUint).
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
		DefaultValue(uint(3), "RETRY").SetOwner(root)

	co.NewFlag(cmdr.OptFlagTypeBool).
		Titles("retry1", "t1").
		Description("", "").
		Group("").
		DefaultValue(false, "RETRY").OwnerCommand()

	co.NewFlag(cmdr.OptFlagTypeInt).
		Titles("retry2", "t2").
		Description("", "").
		Group("").ToggleGroup("").
		VendorHidden(false).
		DefaultValue(3, "RETRY").RootCommand()

	co.NewFlag(cmdr.OptFlagTypeUint64).
		Titles("retry3", "t3").
		Description("", "").
		Group("").
		DefaultValue(uint64(3), "RETRY")

	co.NewFlag(cmdr.OptFlagTypeInt64).
		Titles("retry4", "t4").
		Description("", "").
		Group("").
		DefaultValue(int64(3), "RETRY")

	co.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("retry5", "t5").
		Description("", "").
		Group("").
		DefaultValue([]string{"a", "b"}, "RETRY")

	co.NewFlag(cmdr.OptFlagTypeIntSlice).
		Titles("retry6", "t6").
		Description("", "").
		Group("").
		DefaultValue([]int{1, 2, 3}, "RETRY")

	co.NewFlag(cmdr.OptFlagTypeDuration).
		Titles("retry7", "t7").
		Description("", "").
		Group("").
		DefaultValue(3, "RETRY")

	f0 := co.NewFlag(cmdr.OptFlagTypeFloat32).
		Titles("retry8", "t8").
		Description("", "").
		Group("").
		DefaultValue(3.14, "PI")
	f0.ToFlag().GetTitleZshFlagNamesArray()

	co.NewFlag(cmdr.OptFlagTypeFloat64).
		Titles("retry9", "t9").
		Description("", "").
		Group("").
		DefaultValue(3.14159265358979323846264338327950288419716939937510582097494459230781640628620899, "PI")

	co.NewFlag(cmdr.OptFlagTypeComplex64).
		Titles("retry10", "t10").
		Description("", "").
		Group("").
		DefaultValue(3.14, "PI")

	co.NewFlag(cmdr.OptFlagTypeComplex128).
		Titles("retry11", "t11").
		Description("", "").
		Group("").
		DefaultValue(3.14159265358979323846264338327950288419716939937510582097494459230781640628620899, "PI")

	co.NewFlag(cmdr.OptFlagTypeInt).
		Titles("head", "h").
		Description("", "").
		Group("").
		DefaultValue(1, "").
		HeadLike(true, 1, 8000)

	f1 := co.NewFlag(cmdr.OptFlagTypeString).
		Titles("ienum", "i").
		Description("", "").
		Group("").
		DefaultValue("", "").
		ValidArgs("apple", "banana", "orange")
	f2 := f1.ToFlag()
	f2.GetDottedNamePath()
	f2.Delete()
	f2.GetTitleZshFlagNamesArray()
	f2.GetTitleZshFlagShortName()
	f2.GetTitleZshNamesBy(",", true, true)
	(&cmdr.Flag{}).Delete()

	// ms tags

	cTags := co.NewSubCommand().
		Titles("tags", "t").
		Description("", "").
		Group("")

	cTags.NewFlag(cmdr.OptFlagTypeString).
		Titles("addr", "a").
		Description("", "").
		Group("").
		DefaultValue("consul.ops.local", "ADDR")

	// ms tags ls

	cTags.NewSubCommand().
		Titles("list", "ls").
		Description("", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	fn := func() {
		defer func() {
			if e := recover(); e != nil {
				print(e)
			}
		}()

		c1 := cTags.NewSubCommand().
			Titles("", "")
		c1.ToCommand().GetName()
	}
	fn()

	c8 := cTags.NewSubCommand().
		Titles("add1", "a1").
		Description("", "").
		Group("")

	c9 := cTags.NewSubCommand().
		Titles("add", "").
		Description("", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})
	f91 := c9.NewFlag(cmdr.OptFlagTypeString).
		Titles("ienum", "").Aliases("ie91").
		Description("", "").
		Group("").
		DefaultValue("", "").
		ValidArgs("apple", "banana", "orange")
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

	co := root.NewSubCommand().
		Titles("micro-service", "ms").
		Short("ms").Long("micro-service").Aliases("goms").
		Examples(``).Hidden(false).Deprecated("").
		PreAction(nil).PostAction(nil).Action(nil).
		TailPlaceholder("").
		Description("", "").
		Group("")

	co.OwnerCommand()
	co.SetOwner(root)

	co.NewFlagV(uint(3)).
		Titles("retry", "t").
		Short("tt").Long("retry-tt").Aliases("go-tt").
		Examples(``).Hidden(false).Deprecated("").
		Action(nil).
		ExternalTool(cmdr.ExternalToolEditor).
		ExternalTool(cmdr.ExternalToolPasswordInput).
		Description("", "").
		Group("").
		Placeholder("RETRY").SetOwner(root)

	co.NewFlagV(true).
		Titles("retry1", "t1").
		Description("", "").
		Group("").
		Placeholder("RETRY").OwnerCommand()

	co.NewFlagV(3).
		Titles("retry2", "t2").
		Description("", "").
		Group("").ToggleGroup("").
		RootCommand()

	co.NewFlagV(uint64(3)).
		Titles("retry3", "t3").
		Description("", "").
		Group("")

	co.NewFlagV(int64(3)).
		Titles("retry4", "t4").
		Description("", "").
		Group("")

	co.NewFlagV([]string{"a", "b"}).
		Titles("retry5", "t5").
		Description("", "")

	co.NewFlagV([]int{1, 2, 3}).
		Titles("retry6", "t6").
		Description("", "")

	co.NewFlagV([]uint{1, 2, 3}).
		Titles("retry61", "t61").
		Description("", "")

	co.NewFlagV(time.Second).
		Titles("retry7", "t7")

	co.NewFlagV(float32(3.14)).
		Titles("retry8", "t8").
		Description("", "").
		Group("").
		Placeholder("PI")

	co.NewFlagV(3.14159265358979323846264338327950288419716939937510582097494459230781640628620899).
		Titles("retry9", "t9").
		Description("", "").
		Group("").
		Placeholder("PI")

	co.NewFlagV(complex64(3.14+9i)).
		Titles("retry10", "t10").
		Description("", "").
		Group("")

	co.NewFlagV(3.14+9i).
		Titles("retry11", "t11").
		Description("", "").
		Group("")

	co.NewFlagV(1).
		Titles("head", "h").
		Description("", "").
		Group("").
		HeadLike(true, 1, 8000).
		EnvKeys("AVCX").
		Required().
		Required(false, false, true, false)

	co.NewFlagV("").
		Titles("ienum", "i").
		Description("", "").
		Group("").
		ValidArgs("apple", "banana", "orange")

	// ms tags

	cTags := co.NewSubCommand().
		Titles("tags", "t").
		Description("", "").
		Group("")

	cTags.NewFlagV("consul.ops.local").
		Titles("addr", "a").
		Description("", "").
		Group("").
		Placeholder("ADDR")

	// ms tags ls

	cTags.NewSubCommand().
		Titles("list", "ls").
		Description("", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	cTags.NewSubCommand().
		Titles("add", "a").
		Description("", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	return root
}

func TestFluentAPINew(t *testing.T) {
	cmdr.ResetOptions()
	cmdr.InternalResetWorker()

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
	cmdr.InternalResetWorker()

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
	cmdr.InternalResetWorker()

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
		cmdr.ResetRootInWorker()

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
