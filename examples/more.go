package examples

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/hedzr/is"
	"github.com/hedzr/is/term"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/pkg/dir"
	"github.com/hedzr/cmdr/v2/pkg/logz"
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

	cb := parent.Cmd("start", "s", "st", "run", "startup").
		Description("startup this system service/daemon", ``).
		OnAction(serverStartup)
	cb.Flg("foreground", "f", "fg", "fore").
		Default(false).
		Description("run foreground", ``).
		Build()
	cb.Build()

	cb = parent.Cmd("run1", "nf", "startup1").
		Description("startup this system service/daemon", ``).
		OnAction(serverStartup)
	cb.Cmd("run1", "nf", "startup1").
		Description("startup this system service/daemon", ``).
		OnAction(serverStartup).
		Build()
	cb.Build()

	parent.Cmd("stop", "t", "pause", "resume", "stp", "hlt", "halt").
		Description("stop this system service/daemon", ``).
		OnAction(serverStop).
		Build()
	parent.Cmd("restart", "r", "reload", "live-reload").
		Description("restart/reload/live-reload this system service/daemon", ``).
		OnAction(serverRestart).
		Build()
	parent.Cmd("status", "st").
		Description("display its running status as a system service/daemon", ``).
		OnAction(serverStatus).
		Build()
	parent.Cmd("install", "i", "inst", "setup").
		Description("install as a system service/daemon", ``).
		OnAction(serverInstall).
		Build()
	parent.Cmd("uninstall", "u", "remove", "rm", "delete").
		Description("uninstall a system service/daemon", ``).
		OnAction(serverUninstall).
		Build()
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

	cb := parent.Cmd("backup", "b", "bf", "bkp").
		Description("Dump Consul's KV database to a JSON/YAML file", ``).
		OnAction(kvBackup)
	cb.Flg("output", "o", "out").
		Default("consul-backup.json").
		Description("Write output to a file (*.json / *.yml)", ``).
		PlaceHolder("FILE").
		CompActionStr(`*.(json|yml|yaml)`). //  \*.\(ps\|eps\)
		// ':postscript file:_files -g \*.\(ps\|eps\)'
		Build()
	cb.Build()

	cb = parent.Cmd("restore", "r").
		Description("restore to Consul's KV store, from a a JSON/YAML backup file", ``).
		OnAction(kvRestore)
	cb.Flg("input", "i", "in").
		Default("consul-backup.json").
		Description("Read the input file (*.json / *.yml)", ``).
		PlaceHolder("FILE").
		CompActionStr(`*.(json|yml|yaml)`). //  \*.\(ps\|eps\)
		// ':postscript file:_files -g \*.\(ps\|eps\)'
		Build()
	cb.Build()
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
		OnAction(msList).
		Build()

	cb := parent.Cmd("tags", "t").
		Description("tags operations of a micro-service", "").
		Group("")
	AttachConsulConnectFlags(cb)
	tagsCommands(cb)
	cb.Build()

	// tags := parent.Cmd("tags", "t").
	// 	Description("tags operations of a micro-service", "").
	// 	Group("")
	// AttachConsulConnectFlags(tags)
	// tags.Build()
}

func tagsCommands(parent cli.CommandBuilder) { //nolint:revive
	// ms tags ls

	cb := parent.Cmd("list", "ls", "l", "lst", "dir").
		Description("list tags for ms tags cmd").
		Group("2333.List").
		OnAction(msTagsList)
	cb.Build()

	// ms tags add

	cb = parent.Cmd("add", "a", "new", "create").
		Description("add tags").
		Deprecated("0.2.1").
		Group("")

	cb.Flg("list", "ls", "l", "lst", "dir").
		Default([]string{}).
		Description("tags add: a comma list to be added").
		Group("").
		PlaceHolder("LIST").
		Build()

	cb.Cmd("check", "c", "chk").
		Description("[sub] check").
		Group("").
		Build()

	cbc := cb.Cmd("check-point", "pt", "chk-pt").
		Description("[sub][sub] checkpoint").
		Group("")

	cbc.Flg("add", "a", "add-list").
		Default([]string{}).
		Description("checkpoint: a comma list to be added.").
		PlaceHolder("LIST").
		Group("List").
		Build()

	cbc.Flg("remove", "r", "rm-list", "rm", "del", "delete").
		Default([]string{}).
		Description("checkpoint: a comma list to be removed.", ``).
		PlaceHolder("LIST").
		Group("List").
		Build()
	cbc.Build()

	cbc = cb.Cmd("check-in", "in", "chk-in").
		Description("[sub][sub] check-in").
		Group("")

	cbc.Flg("n", "name").
		Default("").
		Description("check-in name: a string to be added.").
		Group("").
		Build()

	cbc.Cmd("demo-1", "d1").
		Description("[sub][sub] check-in sub, d1").
		Group("").
		Build()

	cbc.Cmd("demo-2", "d2").
		Description("[sub][sub] check-in sub, d2").
		Group("").
		Build()

	cbc.Cmd("demo-3", "d3").
		Description("[sub][sub] check-in sub, d3").
		Group("").
		Build()

	cbc.Cmd("check-out", "out", "chk-out").
		Description("[sub][sub] check-out").
		Group("").
		Build()
	cbc.Build()

	cb.Build()

	// ms tags rm

	cb = parent.Cmd("rm", "r", "remove", "delete", "del", "erase").
		Description("remove tags").
		Group("")
	cb.Flg("list", "ls", "l", "lst", "dir").
		Default([]string{}).
		Description("tags rm: a comma list to be added").
		Group("").
		PlaceHolder("LIST").
		Build()
	cb.Build()

	// ms tags modify

	cb = parent.Cmd("modify", "m", "mod", "modi", "update", "change").
		Description("modify tags of a service.").
		Group("").
		OnAction(msTagsModify)

	cb.Flg("delim", "d", "delimiter").
		Default("=").
		Description("delimiter char in `non-plain` mode.").
		PlaceHolder("").
		CompJustOnce(true).
		Build()

	AttachModifyFlags(cb)

	cb.Flg("add", "a", "add-list").
		Default([]string{}).
		Description("tags modify: a comma list to be added.").
		PlaceHolder("LIST").
		Group("List").
		Build()
	cb.Flg("remove", "r", "rm-list", "rm", "del", "delete").
		Default([]string{}).
		Description("tags modify: a comma list to be removed.").
		PlaceHolder("LIST").
		Group("List").
		Build()

	cb.Build()

	// ms tags toggle

	cb = parent.Cmd("toggle", "t", "tog", "switch").
		Description("toggle tags").
		Group("").
		OnAction(msTagsToggle)

	AttachModifyFlags(cb)

	cb.Flg("set", "s").
		Default([]string{}).
		Description("tags toggle: a comma list to be set").
		Group("").
		PlaceHolder("LIST").
		Build()

	cb.Flg("unset", "un").
		Default([]string{}).
		Description("tags toggle: a comma list to be unset").
		Group("").
		PlaceHolder("LIST").
		Build()

	cb.Flg("address", "a", "addr").
		Default("").
		Description("tags toggle: the address of the service (by id or name)").
		PlaceHolder("HOST:PORT").
		Build()

	cb.Build()
}

func AttachMoreCommandsForTest(parent cli.CommandBuilder, moreAndMore bool) { //nolint:revive
	// test/debug build, many multilevel subcommands here

	parent.Flg("enable-ueh", "ueh").
		Default(false).
		Description("Enables the unhandled exception handler?").
		Build()

	parent.Titles("more", "m").
		Description("More commands").
		With(func(b cli.CommandBuilder) {
			AddXyPrintCommand(b)
			AddKilobytesCommand(b)
			AddSoundexCommand(b)
			AddTtySizeTestCommand(b)

			AddToggleGroupCommand(b)
			AddMxCommand(b)

			AddPanicTestCommand(b)
			AddExternalEditorFlag(b)
			AddTypedFlags(b)

			if moreAndMore {
				AddMultiLevelTestCommand(b)
				cmdrManyCommandsTest(b)
			}

		})

	// pprof.AttachToCmdr(more.RootCmdOpt())
}

func AddMxCommand(parent cli.CommandBuilder) { //nolint:revive
	cb := parent.Cmd("mx-test", "mx").
		Description("mx test new features", "mx test new features,\nverbose long descriptions here.").
		Group("Test").
		OnAction(mxTest)

	cb.Flg("test", "t").
		Default("").
		Description("the test text.", "").
		EnvVars("COOLT", "TEST").
		Group("").
		Build()

	cb.Flg("password", "pp").
		Default("").
		Description("the password requesting.", "").
		Group("").
		PlaceHolder("PASSWORD").
		ExternalEditor(cli.ExternalToolPasswordInput).
		Build()

	cb.Flg("message", "m", "msg").
		Default("").
		Description("the message requesting.", "").
		Group("").
		PlaceHolder("MESG").
		ExternalEditor(cli.ExternalToolEditor).
		Build()

	cb.Flg("fruit", "fr").
		Default("").
		Description("the message.", "").
		Group("").
		PlaceHolder("FRUIT").
		ValidArgs("apple", "banana", "orange").
		Build()

	cb.Flg("head", "hd").
		Default(1).
		Description("the head lines.", "").
		Group("").
		PlaceHolder("LINES").
		HeadLike(true, 1, 3000).
		Build()

	cb.Flg("stdin", "c").
		Default(false).
		Description("read file content from stdin.", "").
		Group("").
		Build()

	cb.Build()
}

func mxTest(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
	_, _ = cmd, args
	// cmdr.Set("test.1", 8)
	cmd.Set().Set("test.deep.branch.1", "test")
	z := cmd.Set().MustString("app.test.deep.branch.1")
	fmt.Printf("*** Got app.test.deep.branch.1: %s\n", z)
	if z != "test" {
		logz.Fatal("[cmdr] err, expect 'test', but got z", "z", z)
	}

	cmd.Set().Remove("app.test.deep.branch.1")
	if cmd.Set().Has("app.test.deep.branch.1") {
		logz.Fatal("[cmdr] FAILED, expect key not found, but found a value associated with: ", "value", cmd.Set().MustR("app.test.deep.branch.1"))
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
			logz.ErrorContext(ctx, "[cmdr] error:", "err", err)
			return
		}
		fmt.Println("> The input contents are:")
		fmt.Print(string(data))
		fmt.Println()
	}
	return
}

func AddXyPrintCommand(parent cli.CommandBuilder) {
	// xy-print

	parent.Cmd("xy-print", "xy").
		Description("test terminal control sequences", "xy-print test terminal control sequences,\nverbose long descriptions here.").
		Group("Test").
		OnAction(xyPrint).
		Build()
}

func xyPrint(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
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
}

func AddPanicTestCommand(parent cli.CommandBuilder) {
	// panic test

	cb := parent.Cmd("panic-test", "pa").
		Description("test panic inside cmdr actions", "").
		Group("Test")

	val := 9
	zeroVal := zero

	cb.Cmd("division-by-zero", "dz").
		Description("").
		Group("Test").
		OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
			_, _ = cmd, args
			fmt.Println(val / zeroVal)
			return
		}).
		Build()
	cb.Cmd("panic", "pa").
		Description("").
		Group("Test").
		OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
			_, _ = cmd, args
			panic(9)
		}).
		Build()

	cb.Build()
}

func AddTtySizeTestCommand(parent cli.CommandBuilder) {
	parent.Cmd("cols", "rows", "tty-size").
		Description("detected tty size").
		Group("Test").
		OnAction(ttySize).
		Build()
}

func ttySize(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
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

		cb := parent.Cmd(t, fmt.Sprintf("sc%v", i)).
			Description(fmt.Sprintf("subcommands %v.sc%v test", ttl, i)).
			Group("Test")
		AddToggleGroupFlags(cb)
		AddValidArgsFlag(cb)
		cb.Build()
	}
}

const zero = 0
