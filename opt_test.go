/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr_test

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"gopkg.in/yaml.v2"
	"testing"
	"time"
)

func TestCommandMethods(t *testing.T) {
	root := cmdr.Root("aa", "1.0.1").
		Header("sds")

	msCmd := root.NewSubCommand().
		Titles("ms", "microservice").
		Description("", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})
	msCmd.NewSubCommand().
		Titles("ls", "list", "l", "lst", "dir").
		Description("list tags", "").
		Group("2333.List").
		Hidden(true).
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})
	msCmd.NewSubCommand().
		Titles("t", "tags").
		Description("tags operations of a micro-service", "").
		Group("")

	xy := root.NewSubCommand().
		Titles("xy", "xy-print").
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
		Titles("mx", "mx-test").
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
	cmdr.NewBool(false)
	cmdr.NewDuration(0)
	cmdr.NewInt(0)
	cmdr.NewInt64(0)
	cmdr.NewIntSlice([]int{})
	cmdr.NewString("")
	cmdr.NewStringSlice([]string{})
	cmdr.NewUint(0)
	cmdr.NewUint64(0)
	cmdr.NewFloat32(0)
	cmdr.NewFloat64(0)

	cmdr.NewOptions()
	cmdr.NewOptionsWith(nil)
}

func TestAsXXX(t *testing.T) {
	cmdr.AsYaml()
	cmdr.AsJSON()
	cmdr.AsToml()
	cmdr.GetHierarchyList()
}

func TestFluentAPINew(t *testing.T) {
	root := cmdr.Root("aa", "1.0.1").
		Header("aa - test for cmdr - no version - hedzr").
		Copyright("s", "x")
	rootCmd1 := root.RootCommand()
	t.Log(rootCmd1)

	// ms

	co := root.NewSubCommand().
		Titles("ms", "micro-service").
		Short("ms").Long("micro-service").Aliases("goms").
		Examples(``).Hidden(false).Deprecated("").
		PreAction(nil).PostAction(nil).Action(nil).
		TailPlaceholder("").
		Description("", "").
		Group("")

	co.OwnerCommand()
	co.SetOwner(root)

	co.NewFlagV(uint(3)).
		Titles("t", "retry").
		Short("tt").Long("retry-tt").Aliases("go-tt").
		Examples(``).Hidden(false).Deprecated("").
		Action(nil).
		ExternalTool(cmdr.ExternalToolEditor).
		ExternalTool(cmdr.ExternalToolPasswordInput).
		Description("", "").
		Group("").
		Placeholder("RETRY").SetOwner(root)

	co.NewFlagV(true).
		Titles("t1", "retry1").
		Description("", "").
		Group("").
		Placeholder("RETRY").OwnerCommand()

	co.NewFlagV(3).
		Titles("t2", "retry2").
		Description("", "").
		Group("").ToggleGroup("").
		RootCommand()

	co.NewFlagV(uint64(3)).
		Titles("t3", "retry3").
		Description("", "").
		Group("")

	co.NewFlagV(int64(3)).
		Titles("t4", "retry4").
		Description("", "").
		Group("")

	co.NewFlagV([]string{"a", "b"}).
		Titles("t5", "retry5").
		Description("", "")

	co.NewFlagV([]int{1, 2, 3}).
		Titles("t6", "retry6").
		Description("", "")

	co.NewFlagV([]uint{1, 2, 3}).
		Titles("t61", "retry61").
		Description("", "")

	co.NewFlagV(time.Second).
		Titles("t7", "retry7")

	co.NewFlagV(float32(3.14)).
		Titles("t8", "retry8").
		Description("", "").
		Group("").
		Placeholder("PI")

	co.NewFlagV(3.14159265358979323846264338327950288419716939937510582097494459230781640628620899).
		Titles("t9", "retry9").
		Description("", "").
		Group("").
		Placeholder("PI")

	co.NewFlagV(1).
		Titles("h", "head").
		Description("", "").
		Group("").
		HeadLike(true, 1, 8000).EnvKeys("AVCX")

	co.NewFlagV("").
		Titles("i", "ienum").
		Description("", "").
		Group("").
		ValidArgs("apple", "banana", "orange")

	// ms tags

	cTags := co.NewSubCommand().
		Titles("t", "tags").
		Description("", "").
		Group("")

	cTags.NewFlagV("consul.ops.local").
		Titles("a", "addr").
		Description("", "").
		Group("").
		Placeholder("ADDR")

	// ms tags ls

	cTags.NewSubCommand().
		Titles("ls", "list").
		Description("", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	cTags.NewSubCommand().
		Titles("a", "add").
		Description("", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

}

func TestFluentAPI(t *testing.T) {

	root := cmdr.Root("aa", "1.0.1").
		Header("aa - test for cmdr - no version - hedzr").
		Copyright("s", "x")
	rootCmd1 := root.RootCommand()
	t.Log(rootCmd1)

	// ms

	co := root.NewSubCommand().
		Titles("ms", "micro-service").
		Short("ms").Long("micro-service").Aliases("goms").
		Examples(``).Hidden(false).Deprecated("").
		PreAction(nil).PostAction(nil).Action(nil).
		TailPlaceholder("").
		Description("", "").
		Group("")

	co.OwnerCommand()
	co.SetOwner(root)

	co.NewFlag(cmdr.OptFlagTypeUint).
		Titles("t", "retry").
		Short("tt").Long("retry-tt").Aliases("go-tt").
		Examples(``).Hidden(false).Deprecated("").
		Action(nil).
		ExternalTool(cmdr.ExternalToolEditor).
		ExternalTool(cmdr.ExternalToolPasswordInput).
		Description("", "").
		Group("").
		DefaultValue(uint(3), "RETRY").SetOwner(root)

	co.NewFlag(cmdr.OptFlagTypeBool).
		Titles("t1", "retry1").
		Description("", "").
		Group("").
		DefaultValue(false, "RETRY").OwnerCommand()

	co.NewFlag(cmdr.OptFlagTypeInt).
		Titles("t2", "retry2").
		Description("", "").
		Group("").ToggleGroup("").
		DefaultValue(3, "RETRY").RootCommand()

	co.NewFlag(cmdr.OptFlagTypeUint64).
		Titles("t3", "retry3").
		Description("", "").
		Group("").
		DefaultValue(uint64(3), "RETRY")

	co.NewFlag(cmdr.OptFlagTypeInt64).
		Titles("t4", "retry4").
		Description("", "").
		Group("").
		DefaultValue(int64(3), "RETRY")

	co.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("t5", "retry5").
		Description("", "").
		Group("").
		DefaultValue([]string{"a", "b"}, "RETRY")

	co.NewFlag(cmdr.OptFlagTypeIntSlice).
		Titles("t6", "retry6").
		Description("", "").
		Group("").
		DefaultValue([]int{1, 2, 3}, "RETRY")

	co.NewFlag(cmdr.OptFlagTypeDuration).
		Titles("t7", "retry7").
		Description("", "").
		Group("").
		DefaultValue(3, "RETRY")

	co.NewFlag(cmdr.OptFlagTypeFloat32).
		Titles("t8", "retry8").
		Description("", "").
		Group("").
		DefaultValue(3.14, "PI")

	co.NewFlag(cmdr.OptFlagTypeFloat64).
		Titles("t9", "retry9").
		Description("", "").
		Group("").
		DefaultValue(3.14159265358979323846264338327950288419716939937510582097494459230781640628620899, "PI")

	co.NewFlag(cmdr.OptFlagTypeInt).
		Titles("h", "head").
		Description("", "").
		Group("").
		DefaultValue(1, "").
		HeadLike(true, 1, 8000)

	co.NewFlag(cmdr.OptFlagTypeString).
		Titles("i", "ienum").
		Description("", "").
		Group("").
		DefaultValue("", "").
		ValidArgs("apple", "banana", "orange")

	// ms tags

	cTags := co.NewSubCommand().
		Titles("t", "tags").
		Description("", "").
		Group("")

	cTags.NewFlag(cmdr.OptFlagTypeString).
		Titles("a", "addr").
		Description("", "").
		Group("").
		DefaultValue("consul.ops.local", "ADDR")

	// ms tags ls

	cTags.NewSubCommand().
		Titles("ls", "list").
		Description("", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	cTags.NewSubCommand().
		Titles("a", "add").
		Description("", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

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
