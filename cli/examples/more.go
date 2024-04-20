package examples

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/hedzr/is"
	"github.com/hedzr/is/term"
	logz "github.com/hedzr/logg/slog" //nolint:gci

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/pkg/dir"
	"github.com/hedzr/cmdr/v2/pkg/text"
)

func AttachServerCommand(parent cli.CommandBuilder) { //nolint:revive
	// kv

	defer parent.Build()

	parent.Titles("server", "s", "svr", "serve").
		Description("server ops: for linux service/daemon", ``)

	parent.Flg("head", "h").
		Default(1).
		Description("head -1 like", ``).
		HeadLike(true, 1, 65536).
		// CompActionStr(`*.(json|yml|yaml)`). //  \*.\(ps\|eps\)
		// ':postscript file:_files -g \*.\(ps\|eps\)'
		Build()
	parent.Flg("tail", "t").
		Default(1).
		Description("todo (dup): tail -1 like", ``).
		// HeadLike(true, 1, 65536).
		// CompActionStr(`*.(json|yml|yaml)`). //  \*.\(ps\|eps\)
		// ':postscript file:_files -g \*.\(ps\|eps\)'
		Build()
	parent.Flg("enum", "e").
		Default("none").
		Description("enum test", ``).
		ValidArgs("none", "apple", "banana", "mongo", "orange", "zig").
		// CompActionStr(`*.(json|yml|yaml)`). //  \*.\(ps\|eps\)
		// ':postscript file:_files -g \*.\(ps\|eps\)'
		Build()
	parent.Flg("retry", "tt").
		Default(1).
		Description("retry tt", ``).
		PlaceHolder("RETRY").
		// CompActionStr(`*.(json|yml|yaml)`). //  \*.\(ps\|eps\)
		// ':postscript file:_files -g \*.\(ps\|eps\)'
		Build()

	parent.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("start", "s", "st", "run", "startup").
			Description("startup this system service/daemon", ``).
			OnAction(serverStartup)
		b.Flg("foreground", "f", "fg", "fore").
			Default(false).
			Description("run foreground", ``).
			Build()
	})

	parent.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("run1", "nf", "startup1").
			Description("startup this system service/daemon", ``).
			OnAction(serverStartup)
		b.AddCmd(func(b cli.CommandBuilder) {
			b.Titles("run1", "nf", "startup1").
				Description("startup this system service/daemon", ``).
				OnAction(serverStartup)
		})
	})

	parent.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("stop", "t", "pause", "resume", "stp", "hlt", "halt").
			Description("stop this system service/daemon", ``).
			OnAction(serverStop)
	})

	parent.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("restart", "r", "reload", "live-reload").
			Description("restart/reload/live-reload this system service/daemon", ``).
			OnAction(serverRestart)
	})

	parent.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("status", "st").
			Description("display its running status as a system service/daemon", ``).
			OnAction(serverStatus)
	})

	parent.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("install", "i", "inst", "setup").
			Description("install as a system service/daemon", ``).
			OnAction(serverInstall)
	})

	parent.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("uninstall", "u", "remove", "rm", "delete").
			Description("uninstall a system service/daemon", ``).
			OnAction(serverUninstall)
	})
}

// AttachKvCommand adds 'kv' command to a builder (eg: app.App)
//
// Example:
//
//	app := cli.New().Info(...)
//	app.AddCmd(func(b cli.CommandBuilder) {
//	  examples.AttachKvCommand(b)
//	})
//	// Or:
//	examples.AttachKvCommand(app.NewCommandBuilder())
func AttachKvCommand(parent cli.CommandBuilder) {
	// kv

	defer parent.Build()

	// parent.AddCmd(func(b cli.CommandBuilder) {
	parent.Titles("kv-store", "kv", "kvstore").
		Description("consul kv store operations...", ``)
	AttachConsulConnectFlags(parent)

	parent.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("backup", "b", "bf", "bkp").
			Description("Dump Consul's KV database to a JSON/YAML file", ``).
			OnAction(kvBackup)
		b.Flg("output", "o", "out").
			Default("consul-backup.json").
			Description("Write output to a file (*.json / *.yml)", ``).
			PlaceHolder("FILE").
			CompActionStr(`*.(json|yml|yaml)`). //  \*.\(ps\|eps\)
			// ':postscript file:_files -g \*.\(ps\|eps\)'
			Build()
	})

	parent.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("restore", "r").
			Description("restore to Consul's KV store, from a a JSON/YAML backup file", ``).
			OnAction(kvRestore)
		b.Flg("input", "i", "in").
			Default("consul-backup.json").
			Description("Read the input file (*.json / *.yml)", ``).
			PlaceHolder("FILE").
			CompActionStr(`*.(json|yml|yaml)`). //  \*.\(ps\|eps\)
			// ':postscript file:_files -g \*.\(ps\|eps\)'
			Build()
	})
	// })
}

func AttachMsCommand(parent cli.CommandBuilder) { //nolint:funlen //a test
	// ms

	defer parent.Build()

	// parent.AddCmd(func(b cli.CommandBuilder) {
	parent.Titles("micro-service", "ms", "microservice").
		Description("micro-service operations...", "").
		Group("")

	parent.Flg("money", "mm").
		Default(false).
		Description("A placeholder flag - money.", "").
		Group("").
		PlaceHolder("").
		Build()

	parent.Flg("name", "n").
		Default("").
		Description("name of the service", ``).
		PlaceHolder("NAME").
		Build()
	parent.Flg("id", "i", "ID").
		Default("").
		Description("unique id of the service", ``).
		PlaceHolder("ID").
		Build()
	parent.Flg("all", "a").
		Default(false).
		Description("all services", ``).
		PlaceHolder("").
		Build()
	parent.Flg("retry", "t").
		Default(3).
		Description("retry times for ms cmd", "").
		Group("").
		PlaceHolder("RETRY").
		Build()

	parent.Cmd("list", "ls", "l", "lst", "dir").
		Description("list tags for ms cmd", "").
		Group("2333.List").
		OnAction(func(cmd *cli.Command, args []string) (err error) {
			_, _ = cmd, args
			return
		}).
		Build()

	parent.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("tags", "t").
			Description("tags operations of a micro-service", "").
			Group("")
		AttachConsulConnectFlags(b)
		tagsCommands(b)
	})

	// tags := parent.Cmd("tags", "t").
	// 	Description("tags operations of a micro-service", "").
	// 	Group("")
	// AttachConsulConnectFlags(tags)
	// tags.Build()
}

func tagsCommands(parent cli.CommandBuilder) { //nolint:revive
	// ms tags ls

	parent.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("list", "ls", "l", "lst", "dir").
			Description("list tags for ms tags cmd").
			Group("2333.List").
			OnAction(func(cmd *cli.Command, args []string) (err error) {
				_, _ = cmd, args
				return
			})
	})

	// ms tags add

	parent.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("add", "a", "new", "create").
			Description("add tags").
			Deprecated("0.2.1").
			Group("")

		b.Flg("list", "ls", "l", "lst", "dir").
			Default([]string{}).
			Description("tags add: a comma list to be added").
			Group("").
			PlaceHolder("LIST").
			Build()

		b.AddCmd(func(b1 cli.CommandBuilder) {
			b1.Titles("check", "c", "chk").
				Description("[sub] check").
				Group("")
		})

		b.AddCmd(func(b1 cli.CommandBuilder) {
			b1.Titles("check-point", "pt", "chk-pt").
				Description("[sub][sub] checkpoint").
				Group("")

			b1.Flg("add", "a", "add-list").
				Default([]string{}).
				Description("checkpoint: a comma list to be added.").
				PlaceHolder("LIST").
				Group("List").
				Build()

			b1.Flg("remove", "r", "rm-list", "rm", "del", "delete").
				Default([]string{}).
				Description("checkpoint: a comma list to be removed.", ``).
				PlaceHolder("LIST").
				Group("List").
				Build()
		})

		b.AddCmd(func(b1 cli.CommandBuilder) {
			b1.Titles("check-in", "in", "chk-in").
				Description("[sub][sub] check-in").
				Group("")

			b1.AddFlg(func(b2 cli.FlagBuilder) {
				b2.Titles("n", "name").
					Default("").
					Description("check-in name: a string to be added.").
					Group("")
			})

			b1.AddCmd(func(b2 cli.CommandBuilder) {
				b2.Titles("demo-1", "d1").
					Description("[sub][sub] check-in sub, d1").
					Group("")
			})

			b1.AddCmd(func(b2 cli.CommandBuilder) {
				b2.Titles("demo-2", "d2").
					Description("[sub][sub] check-in sub, d2").
					Group("")
			})

			b1.AddCmd(func(b2 cli.CommandBuilder) {
				b2.Titles("demo-3", "d3").
					Description("[sub][sub] check-in sub, d3").
					Group("")
			})

			b1.AddCmd(func(b2 cli.CommandBuilder) {
				b2.Titles("check-out", "out", "chk-out").
					Description("[sub][sub] check-out").
					Group("")
			})
		})
	})

	// ms tags rm

	parent.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("rm", "r", "remove", "delete", "del", "erase").
			Description("remove tags").
			Group("")
		b.Flg("list", "ls", "l", "lst", "dir").
			Default([]string{}).
			Description("tags rm: a comma list to be added").
			Group("").
			PlaceHolder("LIST").
			Build()
	})

	// ms tags modify

	parent.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("modify", "m", "mod", "modi", "update", "change").
			Description("modify tags of a service.").
			Group("").
			OnAction(msTagsModify)

		AttachModifyFlags(b)

		b.Flg("add", "a", "add-list").
			Default([]string{}).
			Description("tags modify: a comma list to be added.").
			PlaceHolder("LIST").
			Group("List").
			Build()
		b.Flg("remove", "r", "rm-list", "rm", "del", "delete").
			Default([]string{}).
			Description("tags modify: a comma list to be removed.").
			PlaceHolder("LIST").
			Group("List").
			Build()
	})

	// ms tags toggle

	parent.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("toggle", "t", "tog", "switch").
			Description("toggle tags").
			Group("").
			OnAction(func(cmd *cli.Command, args []string) (err error) {
				_, _ = cmd, args
				return
			})

		AttachModifyFlags(b)

		b.Flg("set", "s").
			Default([]string{}).
			Description("tags toggle: a comma list to be set").
			Group("").
			PlaceHolder("LIST").
			Build()

		b.Flg("unset", "un").
			Default([]string{}).
			Description("tags toggle: a comma list to be unset").
			Group("").
			PlaceHolder("LIST").
			Build()

		b.Flg("address", "a", "addr").
			Default("").
			Description("tags toggle: the address of the service (by id or name)").
			PlaceHolder("HOST:PORT").
			Build()
	})
}

func AttachMoreCommandsForTest(parent cli.CommandBuilder, moreAndMore bool) { //nolint:revive
	// test/debug build, many multilevel subcommands here

	defer parent.Build()

	more := parent.Titles("more", "m").
		Description("More commands")
	defer more.Build()

	cmdrXyPrint(more)
	cmdrKbPrint(more)
	cmdrSoundex(more)
	cmdrTtySize(more)

	tgCommand(more)
	mxCommand(more)

	parent.AddFlg(func(b cli.FlagBuilder) {
		b.Default(false).
			Titles("enable-ueh", "ueh").
			Description("Enables the unhandled exception handler?")
	})

	cmdrPanic(more)

	if moreAndMore {
		cmdrMultiLevelTest(more)
		cmdrManyCommandsTest(more)
	}

	// pprof.AttachToCmdr(more.RootCmdOpt())
}

func tgCommand(parent cli.CommandBuilder) { //nolint:revive
	// toggle-group-test - without a default choice

	parent.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("tg-test", "tg", "toggle-group-test").
			Description("tg test new features", "tg test new features,\nverbose long descriptions here.").
			Group("Test").
			OnAction(func(cmd *cli.Command, args []string) (err error) {
				fmt.Printf("*** Got fruit (toggle group): %v\n", cmd.Set().MustString("app.tg-test.fruit"))

				fmt.Printf("> STDIN MODE: %v \n", cmd.Set().MustBool("mx-test.stdin"))
				fmt.Println()

				// logrus.Debug("debug")
				// logrus.Info("debug")
				// logrus.Warning("debug")
				// logrus.WithField(logex.SKIP, 1).Warningf("dsdsdsds")
				_, _ = cmd, args

				return
			}).
			Build()
		b.AddFlg(func(b cli.FlagBuilder) {
			b.Default(false).
				Titles("apple", "").
				Description("the test text.", "").
				ToggleGroup("fruit").
				Build()
		})
		b.AddFlg(func(b cli.FlagBuilder) {
			b.Default(false).Titles("banana", "").
				Description("the test text.", "").
				ToggleGroup("fruit").
				Build()
		})
		b.AddFlg(func(b cli.FlagBuilder) {
			b.Default(true).Titles("orange", "").
				Description("the test text.", "").
				ToggleGroup("fruit").
				Build()
		})
		b.Build()
	})

	// tg2 - with a default choice

	parent.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("tg-test2", "tg2", "toggle-group-test2").
			Description("tg2 test new features", "tg2 test new features,\nverbose long descriptions here.").
			Group("Test").
			OnAction(func(cmd *cli.Command, args []string) (err error) {
				fmt.Printf("*** Got fruit (toggle group): %v\n", cmd.Set().MustString("app.tg-test2.fruit"))
				_, _ = cmd, args

				fmt.Printf("> STDIN MODE: %v \n", cmd.Set().MustBool("mx-test.stdin"))
				fmt.Println()
				return
			}).
			Build()

		b.AddFlg(func(b cli.FlagBuilder) {
			b.Default(false).Titles("apple", "").
				Description("the test text.", "").
				ToggleGroup("fruit").
				Build()
		})
		b.AddFlg(func(b cli.FlagBuilder) {
			b.Default(false).Titles("banana", "").
				Description("the test text.", "").
				ToggleGroup("fruit").
				Build()
		})
		b.AddFlg(func(b cli.FlagBuilder) {
			b.Default(false).Titles("orange", "").
				Description("the test text.", "").
				ToggleGroup("fruit").
				Build()
		})
	})
}

func mxCommand(parent cli.CommandBuilder) { //nolint:revive
	// mx-test

	parent.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("mx-test", "mx").
			Description("mx test new features", "mx test new features,\nverbose long descriptions here.").
			Group("Test").
			OnAction(func(cmd *cli.Command, args []string) (err error) {
				_, _ = cmd, args
				// cmdr.Set("test.1", 8)
				cmd.Set().Set("test.deep.branch.1", "test")
				z := cmd.Set().MustString("app.test.deep.branch.1")
				fmt.Printf("*** Got app.test.deep.branch.1: %s\n", z)
				if z != "test" {
					logz.Fatal("err, expect 'test', but got z", "z", z)
				}

				cmd.Set().Remove("app.test.deep.branch.1")
				if cmd.Set().Has("app.test.deep.branch.1") {
					logz.Fatal("FAILED, expect key not found, but found a value associated with: ", "value", cmd.Set().MustR("app.test.deep.branch.1"))
				}
				fmt.Printf("*** Got app.test.deep.branch.1 (after deleted): %s\n", cmd.Set().MustString("app.test.deep.branch.1"))

				fmt.Printf("*** Got pp: %s\n", cmd.Set().MustString("app.mx-test.password"))
				fmt.Printf("*** Got msg: %s\n", cmd.Set().MustString("app.mx-test.message"))
				fmt.Printf("*** Got fruit (valid args): %v\n", cmd.Set().MustString("app.mx-test.fruit"))
				fmt.Printf("*** Got head (head-like): %v\n", cmd.Set().MustInt("app.mx-test.head"))
				fmt.Println()
				fmt.Printf("*** test text: %s\n", cmd.Set().MustString("mx-test.test"))
				fmt.Println()
				// fmt.Printf("> InTesting: args[0]=%v \n", tool.SavedOsArgs[0])
				// fmt.Println()
				// fmt.Printf("> Used config file: %v\n", cmd.Set().GetUsedConfigFile())
				// fmt.Printf("> Used config files: %v\n", cmd.Set().GetUsingConfigFiles())
				// fmt.Printf("> Used config sub-dir: %v\n", cmd.Set().GetUsedConfigSubDir())

				fmt.Printf("> STDIN MODE: %v \n", cmd.Set().MustBool("mx-test.stdin"))
				fmt.Println()

				// logrus.Debug("debug")
				// logrus.Info("debug")
				// logrus.Warning("debug")
				// logrus.WithField(logex.SKIP, 1).Warningf("dsdsdsds")

				if cmd.Set().MustBool("mx-test.stdin") {
					fmt.Println("> Type your contents here, press Ctrl-D to end it:")
					var data []byte
					data, err = dir.ReadAll(os.Stdin)
					if err != nil {
						logz.Error("error:", "err", err)
						return
					}
					fmt.Println("> The input contents are:")
					fmt.Print(string(data))
					fmt.Println()
				}
				return
			}).
			Build()

		b.AddFlg(func(b cli.FlagBuilder) {
			b.Default("").
				Titles("test", "t").
				Description("the test text.", "").
				EnvVars("COOLT", "TEST").
				Group("").
				Build()
		})

		b.AddFlg(func(b cli.FlagBuilder) {
			b.Default("").
				Titles("password", "pp").
				Description("the password requesting.", "").
				Group("").
				PlaceHolder("PASSWORD").
				ExternalEditor(cli.ExternalToolPasswordInput).
				Build()
		})

		b.AddFlg(func(b cli.FlagBuilder) {
			b.Default("").
				Titles("message", "m", "msg").
				Description("the message requesting.", "").
				Group("").
				PlaceHolder("MESG").
				ExternalEditor(cli.ExternalToolEditor).
				Build()
		})

		b.AddFlg(func(b cli.FlagBuilder) {
			b.Default("").
				Titles("fruit", "fr").
				Description("the message.", "").
				Group("").
				PlaceHolder("FRUIT").
				ValidArgs("apple", "banana", "orange").
				Build()
		})

		b.AddFlg(func(b cli.FlagBuilder) {
			b.Default(1).
				Titles("head", "hd").
				Description("the head lines.", "").
				Group("").
				PlaceHolder("LINES").
				HeadLike(true, 1, 3000).
				Build()
		})

		b.AddFlg(func(b cli.FlagBuilder) {
			b.Default(false).
				Titles("stdin", "c").
				Description("read file content from stdin.", "").
				Group("").
				Build()
		})
	})
}

func cmdrXyPrint(parent cli.CommandBuilder) {
	// xy-print

	parent.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("xy-print", "xy").
			Description("test terminal control sequences", "xy-print test terminal control sequences,\nverbose long descriptions here.").
			Group("Test").
			OnAction(func(cmd *cli.Command, args []string) (err error) {
				//
				// https://en.wikipedia.org/wiki/ANSI_escape_code
				// https://zh.wikipedia.org/wiki/ANSI%E8%BD%AC%E4%B9%89%E5%BA%8F%E5%88%97
				// https://en.wikipedia.org/wiki/POSIX_terminal_interface
				//

				_, _ = cmd, args

				fmt.Println("\x1b[2J") // clear screen

				for i, s := range args {
					fmt.Printf("\x1b[s\x1b[%d;%dH%s\x1b[u", 15+i, 30, s)
				}

				return
			}).
			Build()
	})
}

func cmdrKbPrint(parent cli.CommandBuilder) {
	// kb-print

	parent.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("kb-print", "kb").
			Description("kilobytes test", "test kibibytes' input,\nverbose long descriptions here.").
			Group("Test").
			Examples(`
$ {{.AppName}} kb --size 5kb
  5kb = 5,120 bytes
$ {{.AppName}} kb --size 8T
  8TB = 8,796,093,022,208 bytes
$ {{.AppName}} kb --size 1g
  1GB = 1,073,741,824 bytes
		`).
			OnAction(func(cmd *cli.Command, args []string) (err error) {
				// fmt.Printf("Got size: %v (literal: %v)\n\n", cmdr.GetKibibytesR("kb-print.size"), cmdr.GetStringR("kb-print.size"))
				_, _ = cmd, args
				return
			}).
			Build()
		b.AddFlg(func(b cli.FlagBuilder) {
			b.Default("1k").Titles("size", "s").
				Description("max message size. Valid formats: 2k, 2kb, 2kB, 2KB. Suffixes: k, m, g, t, p, e.", "").
				Group("").
				Build()
		})
	})
}

func cmdrPanic(parent cli.CommandBuilder) {
	// panic test
	parent.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("panic-test", "pa").
			Description("test panic inside cmdr actions", "").
			Group("Test").
			Build()

		val := 9
		zeroVal := zero

		b.AddCmd(func(b cli.CommandBuilder) {
			b.Titles("division-by-zero", "dz").
				Description("").
				Group("Test").
				OnAction(func(cmd *cli.Command, args []string) (err error) {
					_, _ = cmd, args
					fmt.Println(val / zeroVal)
					return
				}).
				Build()
		})
		b.AddCmd(func(b cli.CommandBuilder) {
			b.Titles("panic", "pa").
				Description("").
				Group("Test").
				OnAction(func(cmd *cli.Command, args []string) (err error) {
					_, _ = cmd, args
					panic(9)
				}).
				Build()
		})
	})
}

func cmdrSoundex(parent cli.CommandBuilder) {
	parent.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("soundex", "snd", "sndx", "sound").
			Description("soundex test").
			Group("Test").
			TailPlaceHolders("[text1, text2, ...]").
			OnAction(func(cmd *cli.Command, args []string) (err error) {
				_, _ = cmd, args
				for ix, s := range args {
					fmt.Printf("%5d. %s => %s\n", ix, s, text.Soundex(s))
				}
				return
			}).
			Build()
	})
}

func cmdrTtySize(parent cli.CommandBuilder) {
	parent.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("cols", "rows", "tty-size").
			Description("detected tty size").
			Group("Test").
			OnAction(func(cmd *cli.Command, args []string) (err error) {
				_, _ = cmd, args

				cols, rows := term.GetTtySize()
				fmt.Printf(" 1. cols = %v, rows = %v\n\n", cols, rows)

				cols, rows, err = term.GetSize(int(os.Stdout.Fd()))
				fmt.Printf(" 2. cols = %v, rows = %v | in-docker: %v\n\n", cols, rows, is.InDocker())
				if err != nil {
					log.Printf("    err: %v", err)
				}

				var out []byte
				cc := exec.Command("stty", "size")
				cc.Stdin = os.Stdin
				out, err = cc.Output()
				fmt.Printf(" 3. out: %v", string(out))
				fmt.Printf("    err: %v\n", err)

				if is.InDocker() {
					log.Printf(" 4  run-in-docker: %v", true)
				}
				return
			}).
			Build()
	})
}

func cmdrManyCommandsTest(parent cli.CommandBuilder) {
	for i := 1; i <= 23; i++ {
		t := fmt.Sprintf("subcmd-%v", i)
		// var ttls []string
		// for o, l := parent, obj.CommandBuilder(nil); o != nil && o != l; {
		// 	ttls = append(ttls, o.ToCommand().GetTitleName())
		// 	l, o = o, o.OwnerCommand()
		// }
		// ttl := strings2.Join(strings.ReverseStringSlice(ttls), ".")
		ttl := ""

		parent.AddCmd(func(b cli.CommandBuilder) {
			b.Titles(t, fmt.Sprintf("sc%v", i)).
				Description(fmt.Sprintf("subcommands %v.sc%v test", ttl, i)).
				Group("Test").
				Build()
			cmdrAddFlags(b)
		})
	}
}

func cmdrMultiLevelTest(parent cli.CommandBuilder) {
	parent.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("mls", "mls").
			Description("multi-level subcommands test").
			Group("Test").
			Build()
		// Sets(func(cmd obj.CommandBuilder) {
		//	cmdrAddFlags(cmd)
		// })

		// cmd := root.NewSubCommand("mls", "mls").
		//	Description("multi-level subcommands test").
		//	Group("Test")
		cmdrAddFlags(b)
		cmdrMultiLevel(b, 1)
	})
}

func cmdrMultiLevel(parent cli.CommandBuilder, depth int) {
	if depth > 3 {
		return
	}

	for i := 1; i < 4; i++ {
		t := fmt.Sprintf("subcmd-%v", i)
		// var ttls []string
		// for o, l := parent, obj.CommandBuilder(nil); o != nil && o != l; {
		// 	ttls = append(ttls, o.ToCommand().GetTitleName())
		// 	l, o = o, o.OwnerCommand()
		// }
		// ttl := strings.Join(tool.ReverseStringSlice(ttls), ".")
		ttl := ""

		parent.AddCmd(func(b cli.CommandBuilder) {
			b.Titles(t, fmt.Sprintf("sc%v", i)).
				// cc := parent.NewSubCommand(t, fmt.Sprintf("sc%v", i)).
				Description(fmt.Sprintf("subcommands %v.sc%v test", ttl, i)).
				Group("Test").
				Build()
			cmdrAddFlags(b)
			cmdrMultiLevel(b, depth+1)
		})
	}
}

func cmdrAddFlags(c cli.CommandBuilder) { //nolint:revive
	c.AddFlg(func(b cli.FlagBuilder) {
		b.Default(false).Titles("apple", "").
			Description("the test text.", "").
			ToggleGroup("fruit").
			Build()
	})
	c.AddFlg(func(b cli.FlagBuilder) {
		b.Default(false).Titles("banana", "").
			Description("the test text.", "").
			ToggleGroup("fruit").
			Build()
	})
	c.AddFlg(func(b cli.FlagBuilder) {
		b.Default(false).Titles("orange", "").
			Description("the test text.", "").
			ToggleGroup("fruit").
			Build()
	})

	c.AddFlg(func(b cli.FlagBuilder) {
		b.Default("").Titles("message", "m", "msg").
			Description("the message requesting.", "").
			Group("").
			PlaceHolder("MESG").
			ExternalEditor(cli.ExternalToolEditor).
			Build()
	})
	c.AddFlg(func(b cli.FlagBuilder) {
		b.Default("").Titles("fruit", "fr").
			Description("the message.", "").
			Group("").
			PlaceHolder("FRUIT").
			ValidArgs("apple", "banana", "orange").
			Build()
	})
	c.AddFlg(func(b cli.FlagBuilder) {
		b.Default(false).
			Titles("bool", "b").
			Description("A bool flag", "").
			Group("").
			Build()
	})
	c.AddFlg(func(b cli.FlagBuilder) {
		b.Default(1).
			Titles("int", "i").
			Description("A int flag", "").
			Group("1000.Integer").
			Build()
	})
	c.AddFlg(func(b cli.FlagBuilder) {
		b.Default(int64(2)).
			Titles("int64", "i64").
			Description("A int64 flag", "").
			Group("1000.Integer").
			Build()
	})
	c.AddFlg(func(b cli.FlagBuilder) {
		b.Default(uint(3)).
			Titles("uint", "u").
			Description("A uint flag", "").
			Group("1000.Integer").
			Build()
	})
	c.AddFlg(func(b cli.FlagBuilder) {
		b.Default(uint64(4)).
			Titles("uint64", "u64").
			Description("A uint64 flag", "").
			Group("1000.Integer").
			Build()
	})
	c.AddFlg(func(b cli.FlagBuilder) {
		b.Default(float32(2.71828)).
			Titles("float32", "f", "float").
			Description("A float32 flag with 'e' value", "").
			Group("2000.Float").
			Build()
	})
	c.AddFlg(func(b cli.FlagBuilder) {
		b.Default(3.14159265358979323846264338327950288419716939937510582097494459230781640628620899).
			Titles("float64", "f64").
			Description("A float64 flag with a `PI` value", "").
			Group("2000.Float").
			Build()
	})
	c.AddFlg(func(b cli.FlagBuilder) {
		b.Default(complex64(3.14+9i)).
			Titles("complex64", "c64").
			Description("A complex64 flag", "").
			Group("2010.Complex").
			Build()
	})
	c.AddFlg(func(b cli.FlagBuilder) {
		b.Default(complex128(3.14+9i)).
			Titles("complex128", "c128").
			Description("A complex128 flag", "").
			Group("2010.Complex").
			Build()
	})

	// a set of booleans

	c.AddFlg(func(b cli.FlagBuilder) {
		b.Default(false).
			Titles("single", "s").
			Description("A bool flag: single", "").
			Group("Boolean").
			EnvVars("").
			Build()
	})
	c.AddFlg(func(b cli.FlagBuilder) {
		b.Default(false).
			Titles("double", "d").
			Description("A bool flag: double", "").
			Group("Boolean").
			EnvVars("").
			Build()
	})
	c.AddFlg(func(b cli.FlagBuilder) {
		b.Default(false).
			Titles("norway", "n", "nw").
			Description("A bool flag: norway", "").
			Group("Boolean").
			Build()
	})
	c.AddFlg(func(b cli.FlagBuilder) {
		b.Default(false).
			Titles("mongo", "mongo").
			Description("A bool flag: mongo", "").
			Group("Boolean").
			Build()
	})
}

const zero = 0
