/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr_test

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/hedzr/cmdr"
	"os"
	"strings"
	"testing"
	"time"
)

func TestSingleCommandLine1(t *testing.T) {
	copyRootCmd = rootCmd
	var err error
	defer func() {
		_ = os.Remove(".tmp.1.json")
		_ = os.Remove(".tmp.1.yaml")
		_ = os.Remove(".tmp.1.toml")
	}()

	os.Args = []string{"consul-tags", "kv", "b"}

	cmdr.SetInternalOutputStreams(nil, nil)
	cmdr.SetHelpTabStop(70)
	cmdr.SetUnknownOptionHandler(nil)

	_ = cmdr.Exec(rootCmd,
		cmdr.WithXrefBuildingHooks(func(root *cmdr.RootCommand, args []string) {}, func(root *cmdr.RootCommand, args []string) {}),
		cmdr.WithAutomaticEnvHooks(func(root *cmdr.RootCommand, opts *cmdr.Options) {}),
		cmdr.WithEnvPrefix("APP_"),
		cmdr.WithOptionsPrefix("app"),
		cmdr.WithRxxtPrefix("app"),
		cmdr.WithPredefinedLocations(),
		cmdr.WithIgnoreWrongEnumValue(true),
		cmdr.WithBuiltinCommands(true, true, true, true, true),
		cmdr.WithInternalOutputStreams(nil, nil),
		cmdr.WithCustomShowVersion(func() {}),
		cmdr.WithCustomShowBuildInfo(func() {}),
		cmdr.WithNoLoadConfigFiles(false),
		cmdr.WithHelpPainter(nil),
		cmdr.WithConfigLoadedListener(nil),
		cmdr.WithHelpTabStop(70),
		cmdr.WithUnknownOptionHandler(func(isFlag bool, title string, cmd *cmdr.Command, args []string) (fallback bool) {
			return
		}),
		cmdr.WithSimilarThreshold(0.7),
	)

	cmdr.ResetWorker()
	_ = cmdr.ExecWith(rootCmd, nil, nil)

	_ = cmdr.SaveAsYaml(".tmp.1.yaml")
	_ = cmdr.SaveAsJSON(".tmp.1.json")
	if err = cmdr.SaveAsToml(".tmp.1.toml"); err != nil {
		// t.Fatal("dump toml failed", err)
	}
	// _ = os.Remove(".tmp.json")

	cmdr.AddOnAfterXrefBuilt(func(root *cmdr.RootCommand, args []string) {
		return
	})
	cmdr.AddOnBeforeXrefBuilding(func(root *cmdr.RootCommand, args []string) {
		return
	})
	cmdr.AddOnConfigLoadedListener(&cfgLoaded{})
	_ = cmdr.ExecWith(rootCmd, func(root *cmdr.RootCommand, args []string) {
		return
	}, func(root *cmdr.RootCommand, args []string) {
		return
	})

	resetOsArgs()
}

func TestUnknownHandler(t *testing.T) {
	copyRootCmd = rootCmd
	cmdr.ResetOptions()

	defer prepareConfD(t)()

	cmdr.SetUnknownOptionHandler(func(isFlag bool, title string, cmd *cmdr.Command, args []string) (fallback bool) {
		t.Logf("isFlag: %v, title: %v, cmd: %v, args: %v", isFlag, title, cmd, args)
		return
	})

	os.Args = []string{"consul-tags", "--confug", "./conf.d"}
	cmdr.SetInternalOutputStreams(nil, nil)
	if err := cmdr.Exec(rootCmd); err != nil {
		t.Fatal(err)
	}
	resetOsArgs()
	cmdr.ResetOptions()

	os.Args = []string{"consul-tags", "tigs"}
	cmdr.SetInternalOutputStreams(nil, nil)
	if err := cmdr.Exec(rootCmd); err != nil {
		t.Fatal(err)
	}
	resetOsArgs()
	cmdr.ResetOptions()

	cmdr.SetUnknownOptionHandler(nil)

	os.Args = []string{"consul-tags", "--confug", "./conf.d"}
	cmdr.SetInternalOutputStreams(nil, nil)
	if err := cmdr.Exec(rootCmd); err != nil {
		t.Fatal(err)
	}
	resetOsArgs()
	cmdr.ResetOptions()

	os.Args = []string{"consul-tags", "tegs"}
	cmdr.SetInternalOutputStreams(nil, nil)
	if err := cmdr.Exec(rootCmd); err != nil {
		t.Fatal(err)
	}
	resetOsArgs()
	cmdr.ResetOptions()

}

func TestConfigOption(t *testing.T) {
	copyRootCmd = rootCmd
	cmdr.ResetOptions()

	defer prepareConfD(t)()

	os.Args = []string{"consul-tags", "--config", "./conf.d"}
	cmdr.SetInternalOutputStreams(nil, nil)
	if err := cmdr.Exec(rootCmd); err != nil {
		t.Fatal(err)
	}
	resetOsArgs()
	cmdr.ResetOptions()

	os.Args = []string{"consul-tags", "--config=./conf.d"}
	cmdr.SetInternalOutputStreams(nil, nil)
	if err := cmdr.Exec(rootCmd); err != nil {
		t.Fatal(err)
	}
	resetOsArgs()
	cmdr.ResetOptions()

	os.Args = []string{"consul-tags", "--config./conf.d/tmp.yaml"}
	cmdr.SetInternalOutputStreams(nil, nil)
	if err := cmdr.Exec(rootCmd); err != nil {
		t.Fatal(err)
	}
	resetOsArgs()
	cmdr.ResetOptions()
}

func TestStrictMode(t *testing.T) {
	copyRootCmd = rootCmd
	cmdr.ResetOptions()
	os.Args = []string{"consul-tags", "ms", "tags", "add", "--strict-mode"}
	cmdr.SetInternalOutputStreams(nil, nil)
	if err := cmdr.Exec(rootCmd); err != nil {
		t.Fatal(err)
	}
	cmdr.ResetOptions()

	os.Args = []string{"consul-tags", "server", "start", "~f", "--strict-mode"}
	cmdr.SetInternalOutputStreams(nil, nil)
	if err := cmdr.Exec(rootCmd); err != nil {
		t.Fatal(err)
	}
	cmdr.ResetOptions()

	os.Args = []string{"consul-tags", "server", "nf", "nf", "--strict-mode"}
	cmdr.SetInternalOutputStreams(nil, nil)
	if err := cmdr.Exec(rootCmd); err != nil {
		t.Fatal(err)
	}

	resetOsArgs()
	cmdr.ResetOptions()
}

func TestTreeDump(t *testing.T) {
	copyRootCmd = rootCmd
	cmdr.ResetOptions()
	for _, args := range [][]string{
		{"consul-tags", "--tree"},
		{"consul-tags", "--no-color", "--tree"},
	} {
		os.Args = args
		cmdr.SetInternalOutputStreams(nil, nil)
		if err := cmdr.Exec(rootCmd); err != nil {
			t.Fatal(err)
		}
		resetOsArgs()
		cmdr.ResetOptions()
	}
}

func TestVersionCommand(t *testing.T) {
	copyRootCmd = rootCmd
	cmdr.ResetOptions()

	defer func() {
		_ = os.Remove(".tmp.1.json")
		_ = os.Remove(".tmp.1.yaml")
		_ = os.Remove(".tmp.1.toml")
	}()

	os.Args = []string{"consul-tags", "version"}
	cmdr.SetInternalOutputStreams(nil, nil)
	if err := cmdr.Exec(rootCmd); err != nil {
		t.Fatal(err)
	}
	cmdr.ResetOptions()

	os.Args = []string{"consul-tags", "ver"}
	cmdr.SetInternalOutputStreams(nil, nil)
	if err := cmdr.Exec(rootCmd); err != nil {
		t.Fatal(err)
	}
	cmdr.ResetOptions()

	os.Args = []string{"consul-tags", "--version"}
	cmdr.SetInternalOutputStreams(nil, nil)
	if err := cmdr.Exec(rootCmd); err != nil {
		t.Fatal(err)
	}
	cmdr.ResetOptions()

	os.Args = []string{"consul-tags", "-#"}
	cmdr.SetInternalOutputStreams(nil, nil)
	if err := cmdr.Exec(rootCmd); err != nil {
		t.Fatal(err)
	}
	cmdr.ResetOptions()

	resetOsArgs()
}

func TestGlobalShow(t *testing.T) {
	copyRootCmd = rootCmd
	cmdr.SetInternalOutputStreams(nil, nil)

	os.Args = []string{"consul-tags", "--version"}
	if err := cmdr.Exec(rootCmd); err != nil {
		t.Fatal(err)
	}

	cmdr.SetCustomShowVersion(func() {
		//
	})
	if err := cmdr.Exec(rootCmd); err != nil {
		t.Fatal(err)
	}

	os.Args = []string{"consul-tags", "--build-info"}
	if err := cmdr.Exec(rootCmd); err != nil {
		t.Fatal(err)
	}

	cmdr.SetCustomShowBuildInfo(func() {
		//
	})
	if err := cmdr.Exec(rootCmd); err != nil {
		t.Fatal(err)
	}

	resetOsArgs()
}

func TestPP(t *testing.T) {
	copyRootCmd = rootCmd
	cmdr.ResetOptions()
	os.Args = []string{"consul-tags", "-pp"}
	cmdr.SetInternalOutputStreams(nil, nil)
	if err := cmdr.Exec(rootCmd); err != nil {
		t.Fatal(err)
	}
	resetOsArgs()

	// os.Args = []string{"consul-tags", "-qq"}
	// cmdr.SetInternalOutputStreams(nil, nil)
	// if err := cmdr.Exec(rootCmd); err != nil {
	// 	t.Fatal(err)
	// }
	// resetOsArgs()

	cmdr.ResetOptions()
}

func TestForGenerateCommands(t *testing.T) {
	copyRootCmd = rootCmd

	defer func() {
		_ = os.Remove(".tmp.1.json")
		_ = os.Remove(".tmp.1.yaml")
		_ = os.Remove(".tmp.1.toml")
	}()

	var commands = []string{
		"consul-tags gen shell --auto",
		"consul-tags gen shell --auto --force-bash",
		"consul-tags gen doc",
		"consul-tags gen pdf",
		"consul-tags gen docx",
		"consul-tags gen tex",
		"consul-tags gen markdown",
		"consul-tags gen d",
		"consul-tags gen doc --pdf",
		"consul-tags gen doc --markdown",
		"consul-tags gen doc --tex",
		"consul-tags gen doc --doc",
		"consul-tags gen doc --docx",
		"consul-tags gen shell --bash",
		"consul-tags gen shell --zsh",
		"consul-tags gen shell",
	}
	for _, cc := range commands {
		cmdr.Set("generate.shell.zsh", false)
		cmdr.Set("generate.shell.bash", false)
		cmdr.Set("generate.shell.auto", false)
		cmdr.Set("generate.shell.force-bash", false)
		cmdr.Set("generate.doc.pdf", false)
		cmdr.Set("generate.doc.markdown", false)
		cmdr.Set("generate.doc.tex", false)
		cmdr.Set("generate.doc.doc", false)
		cmdr.Set("generate.doc.docx", false)

		os.Args = strings.Split(cc, " ")
		cmdr.SetInternalOutputStreams(nil, nil)
		if err := cmdr.Exec(rootCmd); err != nil {
			t.Fatal(err)
		}
	}

	resetOsArgs()
	cmdr.ResetOptions()
}

func TestForGenerateDoc(t *testing.T) {
	copyRootCmd = rootCmd
	var commands = []string{
		"consul-tags gen doc",
	}
	for _, cc := range commands {
		os.Args = strings.Split(cc, " ")
		cmdr.SetInternalOutputStreams(nil, nil)
		if err := cmdr.Exec(rootCmd); err != nil {
			t.Fatal(err)
		}
	}

	resetOsArgs()
	cmdr.ResetOptions()
}

func TestForGenerateMan(t *testing.T) {
	copyRootCmd = rootCmd
	var commands = []string{
		"consul-tags gen man",
	}
	for _, cc := range commands {
		os.Args = strings.Split(cc, " ")
		cmdr.SetInternalOutputStreams(nil, nil)
		if err := cmdr.Exec(rootCmd); err != nil {
			t.Fatal(err)
		}
	}

	resetOsArgs()
	cmdr.ResetOptions()
}

func TestReflectOfSlice(t *testing.T) {
	xs := doubleSlice([]string{"foo", "bar"}).([]string)
	fmt.Println("data =", xs, "len =", len(xs), "cap =", cap(xs))

	ys := doubleSlice([]int{3, 1, 4}).([]int)
	fmt.Println("data =", ys, "len =", len(ys), "cap =", cap(ys))
}

func TestGetSectionFrom(t *testing.T) {

	cmdr.Set("debug", true)
	ac := new(testStruct)
	_ = cmdr.GetSectionFrom("debug", &ac) // for app.debug
	if ac.Debug == false {
		t.Fatal(ac.Debug)
	}

	cmdr.Set("server.head", 3)
	cmdr.Set("server.tail", 5)
	cmdr.Set("server.retry", 7)
	cmdr.Set("server.enum", "bug")
	sc := new(testServerStruct)
	_ = cmdr.GetSectionFrom("server", &sc) // for app.server
	if sc.Enum != "bug" {
		t.Fatal(sc.Enum)
	}

	resetFlagsAndLog(t)
}

func TestTightFlag(t *testing.T) {
	copyRootCmd = rootCmd
	var commands = []string{
		"consul-tags -t3 -s 5 kv b --help-zsh 2 ~~",
		"consul-tags -? -vD --vv kv backup --prefix'' -h ~~debug",
	}
	for _, cc := range commands {
		t.Logf("--- command-line: %v", cc)
		os.Args = strings.Split(cc, " ")
		cmdr.SetInternalOutputStreams(nil, nil)
		if err := cmdr.Exec(rootCmd); err != nil {
			t.Fatal(err)
		}
	}

	resetOsArgs()
	cmdr.ResetOptions()
}

func TestCmdrClone(t *testing.T) {
	cmdr.ResetOptions()

	if rootCmd.SubCommands[1].SubCommands[0].Flags[0] == rootCmd.SubCommands[2].Flags[0] {
		t.Log(rootCmd.SubCommands[1].SubCommands[0].Flags)
		t.Log(rootCmd.SubCommands[2].Flags)
		t.Fatal("should not equal.")
	}

	flags := *cmdr.Clone(&consulConnectFlags, &[]*cmdr.Flag{}).(*[]*cmdr.Flag)
	t.Log(flags)
}

func TestExec(t *testing.T) {
	cmdr.ResetOptions()

	defer prepareConfD(t)()

	var err error
	var outX = bytes.NewBufferString("")
	var errX = bytes.NewBufferString("")
	var outBuf = bufio.NewWriterSize(outX, 16384)
	var errBuf = bufio.NewWriterSize(errX, 16384)
	cmdr.SetInternalOutputStreams(outBuf, errBuf)
	cmdr.SetCustomShowVersion(nil)
	cmdr.SetCustomShowBuildInfo(nil)

	defer func() {
		x := outX.String()
		t.Logf("--------- stdout // %v // %v\n%v", cmdr.GetExcutableDir(), cmdr.GetExcutablePath(), x)

		_ = cmdr.EnsureDir("ci")
		if err = cmdr.EnsureDir(""); err == nil {
			t.Failed()
		}
		if err = cmdr.EnsureDir(".tmp"); err == nil {
			_ = os.Remove(".tmp")
		}

		cmdr.SetPredefinedLocations([]string{})
		if len(cmdr.GetPredefinedLocations()) != 0 {
			t.Failed()
		}
		cmdr.SetNoLoadConfigFiles(false)
		cmdr.SetCurrentHelpPainter(nil)
	}()

	copyRootCmd = rootCmd
	_ = cmdr.RootFrom(rootCmd)
	flg := cmdr.FindFlag("ddduration", &copyRootCmd.Command)
	flgOpt := cmdr.NewDurationFrom(flg)
	flgOpt.OnSet(func(keyPath string, value interface{}) {})

	t.Log("xxx: -------- loops for execTestings")
	for sss, verifier := range execTestings {
		resetFlagsAndLog(t)

		t.Log("xxx: ***: ", sss)

		if sss == "consul-tags -qq" {
			fmt.Println("xxx: ***: ", sss)
			rootCmd.Header = "fhsdjkfhdsfhskfhsdjhfksdjfkhsjhfds"
			cmdr.SetCustomShowVersion(func() {
			})
			cmdr.SetCustomShowBuildInfo(func() {
			})
		}
		if sss == "consul-tags ms tags modify -h ~~debug --port8509 --prefix/" {
			fmt.Println("xx*: ***: ", sss)
		}

		if err = cmdr.Worker().InternalExecFor(rootCmd, strings.Split(sss, " ")); err != nil {
			t.Fatal(err)
		}
		if sss == "consul-tags kv unknown" {
			errX = bytes.NewBufferString("")
		}
		if err = verifier(t); err != nil {
			t.Fatal(err)
		}

		if cmdr.GetStrictMode() == false && cmdr.GetQuietMode() == false {
			rootCmd.Header = ""
		}
	}

	postWorks(t)

	if errX.Len() > 0 {
		t.Log("--------- stderr")
		t.Fatalf("Error!! %v", errX.String())
	}

	resetOsArgs()
}

var (
	// testing args
	execTestings = map[string]func(t *testing.T) error{
		// "consul-tags -qq": func(t *testing.T) error {
		// 	return nil
		// },
		"consul-tags --help --help-zsh 1": func(t *testing.T) error {
			return nil
		},
		"consul-tags --help --help-bash": func(t *testing.T) error {
			return nil
		},
		"consul-tags ms dr --help": func(t *testing.T) error {
			return nil
		},
		"consul-tags ms dz --help": func(t *testing.T) error {
			fmt.Println("~ consul-tags ms dz --help")
			return nil
		},
		"consul-tags ms dz dz --help": func(t *testing.T) error {
			fmt.Println("~ consul-tags ms dz dz --help")
			return nil
		},
		"consul-tags ms ls --help": func(t *testing.T) error {
			return nil
		},
		"consul-tags --no-color --help": func(t *testing.T) error {
			return nil
		},
		"consul-tags --version-sim 3.3.3": func(t *testing.T) error {
			return nil
		},
		"consul-tags -pp": func(t *testing.T) error {
			return nil
		},
		"consul-tags -dd 1h": func(t *testing.T) error {
			return nil
		},
		"consul-tags ~dd 1h": func(t *testing.T) error {
			return nil
		},
		"consul-tags ms tags --help --no-color": func(t *testing.T) error {
			return nil
		},
		"consul-tags kv b -K- -K+ --": func(t *testing.T) error {
			// gocov Command.PrintXXX
			fmt.Println("consul-tags kv b -------- no errors")
			return nil
		},
		"consul-tags -t3 -s 5 kv b --help-zsh 2 ~~": func(t *testing.T) error {
			// gocov Command.PrintXXX
			fmt.Println("consul-tags -t3 -s5 -pp kv b ~~ -------- no errors")
			return nil
		},
		"consul-tags server --help": func(t *testing.T) error {
			fmt.Println("consul-tags server --help -------- no errors")
			return nil
		},
		"consul-tags kv b ~": func(t *testing.T) error {
			// gocov Command.PrintXXX
			fmt.Println("consul-tags kv b ~ -------- no errors")
			return nil
		},
		"consul-tags kv unknown": func(t *testing.T) error {
			return nil
		},
		"consul-tags ms tags modify -h ~~debug --port8509 --prefix/": func(t *testing.T) error {
			if cmdr.GetInt("app.ms.tags.port") != 8509 || cmdr.GetString("app.ms.tags.prefix") != "/" ||
				!cmdr.GetBool("app.help") || !cmdr.GetBool("debug") {
				return fmt.Errorf("something wrong 1. |%v|%v|%v|%v",
					cmdr.GetInt("app.ms.tags.port"), cmdr.GetString("app.ms.tags.prefix"),
					cmdr.GetBool("app.help"), cmdr.GetBool("debug"))
			}
			return nil
		},
		"consul-tags -? -vD --vv kv backup --prefix'' -h ~~debug": func(t *testing.T) error {
			if cmdr.GetInt("app.kv.port") != 8500 || cmdr.GetString("app.kv.prefix") != "" ||
				!cmdr.GetBool("app.help") || !cmdr.GetBool("debug") ||
				!cmdr.GetVerboseMode() || !cmdr.GetDebugMode() {
				return fmt.Errorf("something wrong 2. |%v|%v|%v|%v|%v|%v",
					cmdr.GetInt("app.kv.port"), cmdr.GetString("app.kv.prefix"),
					cmdr.GetBool("app.help"), cmdr.GetBool("debug"),
					cmdr.GetVerboseMode(), cmdr.GetDebugMode())
			}
			return nil
		},
		"consul-tags -vD --vv ms tags modify --prefix'' --help ~~debug --prefix\"\" --prefix'cmdr' --prefix\"app\" --prefix=/ --prefix/ --prefix /": func(t *testing.T) error {
			if cmdr.GetInt("app.ms.tags.port") != 8500 || cmdr.GetString("app.ms.tags.prefix") != "/" ||
				!cmdr.GetBool("app.help") || !cmdr.GetBool("debug") ||
				!cmdr.GetVerboseMode() || !cmdr.GetDebugMode() {
				return fmt.Errorf("something wrong 3. |%v|%v|%v|%v|%v|%v",
					cmdr.GetInt("app.ms.tags.port"), cmdr.GetString("app.ms.tags.prefix"),
					cmdr.GetBool("app.help"), cmdr.GetBool("debug"),
					cmdr.GetVerboseMode(), cmdr.GetDebugMode())
			}
			return nil
		},
		"consul-tags -vD ms tags -K- modify --prefix'' -a a,b,v -z 1,2,3 -x '-1,-2' -? ~~debug --port8509 -p8507 -p=8506 -p 8503": func(t *testing.T) error {
			fmt.Println(cmdr.GetStringSlice("app.ms.tags.modify.add"))
			fmt.Println(cmdr.GetIntSlice("app.ms.tags.modify.zed"))
			fmt.Println(cmdr.GetStringSlice("app.ms.tags.modify.xed"))
			fmt.Println(cmdr.GetIntSlice("app.ms.tags.modify.xed"))
			if cmdr.GetInt("app.ms.tags.port") != 8503 || cmdr.GetString("app.ms.tags.prefix") != "" ||
				!cmdr.GetBool("app.help") || !cmdr.GetBool("debug") ||
				!cmdr.GetVerboseMode() || !cmdr.GetDebugMode() {
				return fmt.Errorf("something wrong 4. |%v|%v|%v|%v|%v|%v",
					cmdr.GetInt("app.ms.tags.port"), cmdr.GetString("app.ms.tags.prefix"),
					cmdr.GetBool("app.help"), cmdr.GetBool("debug"),
					cmdr.GetVerboseMode(), cmdr.GetDebugMode())
			}
			return nil
		},
	}

	// testing rootCmd

	copyRootCmd *cmdr.RootCommand

	rootCmd = &cmdr.RootCommand{
		Command: cmdr.Command{
			BaseOpt: cmdr.BaseOpt{
				Name: "consul-tags",
			},
			Flags: []*cmdr.Flag{
				// global options here.
				{
					BaseOpt: cmdr.BaseOpt{
						Short:       "t",
						Full:        "retry",
						Description: "ss",
						Examples:    `random examples`,
						Deprecated:  "1.2.3",
					},
					DefaultValue:            1,
					DefaultValuePlaceholder: "RETRY",
				},
				{
					BaseOpt: cmdr.BaseOpt{
						Short:       "s",
						Full:        "",
						Description: "",
						Action: func(cmd *cmdr.Command, args []string) (err error) {
							flg := cmd.Flags[0]
							if flg.GetDescZsh() != "ss" {
								err = errors.New("err `t`.GetDescZsh()")
							}
							if flg.GetTitleZshFlagNames(",") == "" {
								err = errors.New("err ss.GetTitleZshFlagNames()")
							}
							if len(flg.GetTitleZshFlagNamesArray()) != 2 {
								err = errors.New("err ss.GetTitleZshFlagNamesArray()")
							}
							flg = cmd.Flags[1]
							if len(flg.GetDescZsh()) == 0 {
								err = errors.New("err sss.GetDescZsh()")
							}
							if flg.GetTitleZshFlagNames(",") == "" {
								err = errors.New("err ss.GetTitleZshFlagNames()")
							}
							if len(flg.GetTitleZshFlagNamesArray()) != 2 {
								err = errors.New("err ss.GetTitleZshFlagNamesArray()")
							}
							flg = cmd.Flags[2]
							if len(flg.GetDescZsh()) == 0 {
								err = errors.New("err ssss.GetDescZsh()")
							}
							if flg.GetTitleZshFlagNames(",") == "" {
								err = errors.New("err ss.GetTitleZshFlagNames()")
							}
							if len(flg.GetTitleZshFlagNamesArray()) != 2 {
								err = errors.New("err ss.GetTitleZshFlagNamesArray()")
							}
							return
						},
					},
					DefaultValue: uint(1),
				},
				{
					BaseOpt: cmdr.BaseOpt{
						Short:       "pp",
						Full:        "spasswd",
						Description: "",
						Action: func(cmd *cmdr.Command, args []string) (err error) {
							cmd.PrintVersion()
							cmdr.PrintBuildInfo()
							cmd.PrintBuildInfo()

							cmdr.SetCustomShowVersion(nil)
							cmdr.SetCustomShowBuildInfo(nil)
							return
						},
					},
					DefaultValue: "",
					ExternalTool: cmdr.ExternalToolPasswordInput,
				},
				{
					BaseOpt: cmdr.BaseOpt{
						Short:       "qq",
						Full:        "qqpasswd",
						Description: "",
					},
					DefaultValue: "567",
					ExternalTool: cmdr.ExternalToolEditor,
				},
				{
					BaseOpt: cmdr.BaseOpt{
						Short:       "dd",
						Full:        "ddduration",
						Description: "",
					},
					DefaultValue: time.Second,
				},
			},
			PreAction: func(cmd *cmdr.Command, args []string) (err error) {
				return
			},
			PostAction: func(cmd *cmdr.Command, args []string) {
				return
			},
			SubCommands: []*cmdr.Command{
				// dnsCommands,
				// playCommand,
				// generatorCommands,
				serverCommands,
				msCommands,
				kvCommands,
			},
		},

		AppName: "consul-tags",
		Version: "0.0.1",
		Header:  `dsjlfsdjflsdfjlsdjflksjdfdsfsd`,
		// Version:    consul_tags.Version,
		// VersionInt: consul_tags.VersionInt,
		Copyright: "consul-tags is an effective devops tool",
		Author:    "Hedzr Yeh <hedzrz@gmail.com>",
	}

	serverCommands = &cmdr.Command{
		BaseOpt: cmdr.BaseOpt{
			// Name:        "server",
			Short:       "s",
			Full:        "server",
			Aliases:     []string{"serve", "svr"},
			Description: "server ops: for linux service/daemon.",
			Deprecated:  "1.0",
			Examples:    `random examples`,
		},
		Flags: []*cmdr.Flag{
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "h",
					Full:        "head",
					Description: "head -1 like",
				},
				DefaultValue: 0,
				HeadLike:     true,
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "l",
					Full:        "tail",
					Description: "tail -1 like",
				},
				DefaultValue: 0,
				HeadLike:     true,
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "e",
					Full:        "enum",
					Description: "enum tests",
				},
				DefaultValue: "apple",
				ValidArgs:    []string{"apple", "banana", "orange"},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "tt",
					Full:        "retry",
					Description: "ss",
				},
				DefaultValue:            1,
				DefaultValuePlaceholder: "RETRY",
			},
		},
		SubCommands: []*cmdr.Command{
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "s",
					Full:        "start",
					Aliases:     []string{"run", "startup"},
					Description: "startup this system service/daemon.",
					// Action:impl.ServerStart,
				},
				// PreAction: impl.ServerStartPre,
				// PostAction: impl.ServerStartPost,
				Flags: []*cmdr.Flag{
					// {
					// 	BaseOpt: cmdr.BaseOpt{
					// 		Short:       "t",
					// 		Full:        "retry",
					// 		Description: "ss",
					// 	},
					// 	DefaultValue:            1,
					// 	DefaultValuePlaceholder: "RETRY",
					// },
					// {
					// 	BaseOpt: cmdr.BaseOpt{
					// 		Short:       "t",
					// 		Full:        "retry",
					// 		Description: "ss: dup test",
					// 	},
					// 	DefaultValue:            1,
					// 	DefaultValuePlaceholder: "RETRY",
					// },
					// {
					// 	BaseOpt: cmdr.BaseOpt{
					// 		Name:        "retry",
					// 		Description: "ss: dup test",
					// 	},
					// 	DefaultValue:            1,
					// 	DefaultValuePlaceholder: "RETRY",
					// },
				},
			},
			// {
			// 	BaseOpt: cmdr.BaseOpt{
			// 		Short:       "s",
			// 		Full:        "start",
			// 		Aliases:     []string{"run", "startup"},
			// 		Description: "dup test: startup this system service/daemon.",
			// 		// Action:impl.ServerStart,
			// 	},
			// 	// PreAction: impl.ServerStartPre,
			// 	// PostAction: impl.ServerStartPost,
			// },
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "nf", // parent no Full
					Aliases:     []string{"run1", "startup1"},
					Description: "dup test: startup this system service/daemon.",
				},
				PreAction: func(cmd *cmdr.Command, args []string) (err error) {
					fmt.Println(cmd.GetParentName(), cmd.GetName(), cmd.GetQuotedGroupName(), cmd.GetExpandableNames())
					fmt.Println(cmd.GetRoot().GetParentName())
					return
				},
				SubCommands: []*cmdr.Command{
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "nf", // parent no Full
							Aliases:     []string{"run", "startup"},
							Description: "dup test: startup this system service/daemon.",
						},
						PreAction: func(cmd *cmdr.Command, args []string) (err error) {
							fmt.Println(cmd.GetParentName(), cmd.GetName(), cmd.GetQuotedGroupName(), cmd.GetExpandableNames())
							return
						},
					},
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "t",
					Full:        "stop",
					Aliases:     []string{"stp", "halt", "pause"},
					Description: "stop this system service/daemon.",
					// Action:impl.ServerStop,
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "r",
					Full:        "restart",
					Aliases:     []string{"reload"},
					Description: "restart this system service/daemon.",
					// Action:impl.ServerRestart,
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Full:        "status",
					Aliases:     []string{"st"},
					Description: "display its running status as a system service/daemon.",
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "i",
					Full:        "install",
					Aliases:     []string{"setup"},
					Description: "install as a system service/daemon.",
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "u",
					Full:        "uninstall",
					Aliases:     []string{"remove"},
					Description: "remove from a system service/daemon.",
				},
			},
		},
	}

	kvCommands = &cmdr.Command{
		BaseOpt: cmdr.BaseOpt{
			Name:        "kvstore",
			Full:        "kv",
			Aliases:     []string{"kvstore"},
			Description: "consul kv store operations...",
		},
		Flags: *cmdr.Clone(&consulConnectFlags, &[]*cmdr.Flag{}).(*[]*cmdr.Flag),
		SubCommands: []*cmdr.Command{
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "b",
					Full:        "backup",
					Aliases:     []string{"bk", "bf", "bkp"},
					Description: "Dump Consul's KV database to a JSON/YAML file",
					Group:       "bbb",
					Action: func(cmd *cmdr.Command, args []string) (err error) {
						// for gocov
						cmd.PrintHelp(false)
						cmd.PrintVersion()

						if cmd.GetRoot() != copyRootCmd {
							return errors.New("failed: root is wrong")
						}
						if copyRootCmd.IsRoot() != true {
							return errors.New("failed: root test is wrong")
						}
						if cmd.GetHitStr() != "b" {
							return errors.New("failed: GetHitStr() is wrong")
						}
						if cmd.GetName() != "backup" {
							return errors.New("failed: GetName() is wrong")
						}
						if cmd.GetExpandableNames() != "{backup,b}" {
							return errors.New("failed: GetExpandableNames() is wrong")
						}
						if cmd.GetQuotedGroupName() != "[bbb]" {
							return errors.New("failed: GetQuotedGroupName() is wrong")
						}

						if cmd.GetParentName() != "kv" {
							return errors.New("failed: GetParentName() is wrong")
						}
						if cmd.GetOwner().GetSubCommandNamesBy(",") != "b,backup,bk,bf,bkp,r,restore" {
							return errors.New(fmt.Sprintf("failed: GetSubCommandNamesBy() is wrong: '%s'", cmd.GetOwner().GetSubCommandNamesBy(",")))
						}

						cmd.PrintHelp(true)
						return
					},
				},
				Flags: []*cmdr.Flag{
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "o",
							Full:        "output",
							Description: "Write output to a file (*.json / *.yml)",
							Deprecated:  "2.0",
						},
						DefaultValue:            "consul-backup.json",
						DefaultValuePlaceholder: "FILE",
					},
				},
				PreAction: func(cmd *cmdr.Command, args []string) (err error) {
					return
				},
				PostAction: func(cmd *cmdr.Command, args []string) {
					return
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "r",
					Full:        "restore",
					Description: "restore to Consul's KV store, from a a JSON/YAML backup file",
					// Action:      kvRestore,
				},
				Flags: []*cmdr.Flag{
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "i",
							Full:        "input",
							Description: "read the input file (*.json / *.yml)",
						},
						DefaultValue:            "consul-backup.json",
						DefaultValuePlaceholder: "FILE",
					},
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "hh",
					Full:        "hidden-cmd",
					Description: "restore to Consul's KV store, from a a JSON/YAML backup file",
					Hidden:      true,
				},
			},
		},
	}

	msCommands = &cmdr.Command{
		BaseOpt: cmdr.BaseOpt{
			Name:        "microservices",
			Full:        "ms",
			Aliases:     []string{"microservice", "micro-service"},
			Description: "micro-service operations...",
		},
		Flags: []*cmdr.Flag{
			{
				BaseOpt: cmdr.BaseOpt{
					Short:           "n",
					Full:            "name",
					Description:     "name of the service",
					LongDescription: `fdhsjsfhdsjk`,
					Examples:        `fsdhjkfsdhk`,
				},
				DefaultValue:            "",
				DefaultValuePlaceholder: "NAME",
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "i",
					Full:        "id",
					Description: "unique id of the service",
				},
				DefaultValue:            "",
				DefaultValuePlaceholder: "ID",
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "a",
					Full:        "all",
					Description: "all services",
				},
				DefaultValue: false,
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "cc",
					Full:        "",
					Description: "unique id of the service",
				},
				DefaultValue:            "",
				DefaultValuePlaceholder: "ID",
			},
		},
		SubCommands: []*cmdr.Command{
			tagsCommands,
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "l",
					Full:        "list",
					Aliases:     []string{"ls", "lst"},
					Description: "list services.",
					// Action:      msList,
					Group: " ",
				},
				PreAction: func(cmd *cmdr.Command, args []string) (err error) {
					fmt.Println(cmd.GetParentName(), cmd.GetName(), cmd.GetQuotedGroupName(), cmd.GetExpandableNames())
					return
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					// Short: "",
					// Full:        "",
					// Aliases:     []string{"ls", "lst", "dir"},
					Description: "3 empty - list services.",
					Group:       "56.vvvvvv",
				},
				PreAction: func(cmd *cmdr.Command, args []string) (err error) {
					fmt.Println(cmd.GetParentName(), cmd.GetName(), cmd.GetQuotedGroupName(), cmd.GetExpandableNames())
					return
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short: "dr",
					// Full:        "list",
					// Aliases:     []string{"ls", "lst", "dir"},
					Description: "list services.",
					Group:       "56.vvvvvv",
				},
				PreAction: func(cmd *cmdr.Command, args []string) (err error) {
					fmt.Println(cmd.GetParentName(), cmd.GetName(), cmd.GetQuotedGroupName(), cmd.GetExpandableNames())
					return
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Name: "dz",
					// Full:        "list",
					// Aliases:     []string{"ls", "lst", "dir"},
					Description: "list services.",
					Group:       "56.vvvvvv",
				},
				PreAction: func(cmd *cmdr.Command, args []string) (err error) {
					fmt.Println(cmd, "'s owner is", cmd.GetOwner())
					fmt.Println(cmd.GetParentName(), cmd.GetName(), cmd.GetQuotedGroupName(), cmd.GetExpandableNames())
					return
				},
				SubCommands: []*cmdr.Command{
					{
						BaseOpt: cmdr.BaseOpt{
							Name: "dz",
							// Full:        "list",
							// Aliases:     []string{"ls", "lst", "dir"},
							Description: "list services.",
							Group:       "56.vvvvvv",
						},
						PreAction: func(cmd *cmdr.Command, args []string) (err error) {
							fmt.Println(cmd, "'s owner is", cmd.GetOwner())
							fmt.Println(cmd.GetParentName(), cmd.GetName(), cmd.GetQuotedGroupName(), cmd.GetExpandableNames())
							return
						},
					},
				},
			},
		},
	}

	tagsCommands = &cmdr.Command{
		BaseOpt: cmdr.BaseOpt{
			// Short:       "t",
			Full:        "tags",
			Aliases:     []string{},
			Description: "tags op.",
		},
		Flags: *cmdr.Clone(&consulConnectFlags, &[]*cmdr.Flag{}).(*[]*cmdr.Flag),
		SubCommands: []*cmdr.Command{
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "ls",
					Full:        "list",
					Aliases:     []string{"l", "lst", "dir"},
					Description: "list tags.",
					// Action:      msTagsList,
					Group: "2323.List",
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "a",
					Full:        "add",
					Aliases:     []string{"create", "new"},
					Description: "add tags.",
					// Action:      msTagsAdd,
					Group: "311Z.Add",
				},
				Flags: []*cmdr.Flag{
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "ls",
							Full:        "list",
							Aliases:     []string{"l", "lst", "dir"},
							Description: "a comma list to be added",
						},
						DefaultValue:            []string{},
						DefaultValuePlaceholder: "LIST",
					},
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "r",
					Full:        "rm",
					Aliases:     []string{"remove", "erase", "delete", "del"},
					Description: "remove tags.",
					// Action:      msTagsRemove,
				},
				Flags: []*cmdr.Flag{
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "ls",
							Full:        "list",
							Aliases:     []string{"l", "lst", "dir"},
							Description: "a comma list to be added.",
						},
						DefaultValue:            []string{},
						DefaultValuePlaceholder: "LIST",
					},
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "m",
					Full:        "modify",
					Aliases:     []string{"mod", "update", "change"},
					Description: "modify tags.",
					// Action:      msTagsModify,
				},
				Flags: []*cmdr.Flag{
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "a",
							Full:        "add",
							Description: "a comma list to be added.",
						},
						DefaultValue:            []string{},
						DefaultValuePlaceholder: "LIST",
					},
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "r",
							Full:        "rm",
							Aliases:     []string{"remove", "erase", "del"},
							Description: "a comma list to be removed.",
						},
						DefaultValue:            []string{},
						DefaultValuePlaceholder: "LIST",
					},
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "u",
							Full:        "ued",
							Description: "a comma list to be removed.",
						},
						DefaultValue:            "7,99",
						DefaultValuePlaceholder: "LIST",
					},
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "w",
							Full:        "wed",
							Description: "a comma list to be removed.",
						},
						DefaultValue:            []string{"2", "3"},
						DefaultValuePlaceholder: "LIST",
					},
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "z",
							Full:        "zed",
							Description: "a comma list to be removed.",
						},
						DefaultValue:            []uint{2, 3},
						DefaultValuePlaceholder: "LIST",
					},
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "x",
							Full:        "xed",
							Description: "a comma list to be removed.",
						},
						DefaultValue:            []int{4, 5},
						DefaultValuePlaceholder: "LIST",
					},
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "v",
							Full:        "ved",
							Description: "a comma list to be removed.",
						},
						DefaultValue: 2 * time.Second,
					},
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "t",
					Full:        "toggle",
					Aliases:     []string{"tog", "switch"},
					Description: "toggle tags for ms.",
					// Action:      msTagsToggle,
				},
				Flags: []*cmdr.Flag{
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "x",
							Full:        "address",
							Description: "the address of the service (by id or name)",
						},
						DefaultValue:            "",
						DefaultValuePlaceholder: "HOST:PORT",
					},
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "s",
							Full:        "set",
							Description: "set to `tag` which service specified by --address",
						},
						DefaultValue:            []string{},
						DefaultValuePlaceholder: "LIST",
					},
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "u",
							Full:        "unset",
							Aliases:     []string{"reset"},
							Description: "and reset the others service nodes to `tag`",
						},
						DefaultValue:            []string{},
						DefaultValuePlaceholder: "LIST",
					},
				},
			},
		},
	}

	consulConnectFlags = []*cmdr.Flag{
		{
			BaseOpt: cmdr.BaseOpt{
				Short:       "a",
				Full:        "addr",
				Description: "Consul ip/host and port: HOST[:PORT] (No leading 'http(s)://')",
			},
			DefaultValue:            "consul.ops.local",
			DefaultValuePlaceholder: "HOST[:PORT]",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:       "p",
				Full:        "port",
				Description: "Consul port",
			},
			DefaultValue:            8500,
			DefaultValuePlaceholder: "PORT",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:       "K",
				Full:        "insecure",
				Description: "Skip TLS host verification",
			},
			DefaultValue:            true,
			DefaultValuePlaceholder: "PORT",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:       "",
				Full:        "prefix",
				Description: "Root key prefix",
			},
			DefaultValue:            "/",
			DefaultValuePlaceholder: "ROOT",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:       "",
				Full:        "cacert",
				Description: "Client CA cert",
			},
			DefaultValue:            "",
			DefaultValuePlaceholder: "FILE",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:       "",
				Full:        "cert",
				Description: "Client cert",
			},
			DefaultValue:            "",
			DefaultValuePlaceholder: "FILE",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:       "",
				Full:        "scheme",
				Description: "Consul connection scheme (HTTP or HTTPS)",
			},
			DefaultValue:            "",
			DefaultValuePlaceholder: "SCHEME",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:       "u",
				Full:        "username",
				Description: "HTTP Basic auth user",
			},
			DefaultValue:            "",
			DefaultValuePlaceholder: "USERNAME",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:       "pw",
				Full:        "password",
				Aliases:     []string{"passwd", "pwd"},
				Description: "HTTP Basic auth password",
			},
			DefaultValue:            "",
			DefaultValuePlaceholder: "PASSWORD",
		},
	}
)
