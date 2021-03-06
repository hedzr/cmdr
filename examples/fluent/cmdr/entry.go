package cmdr

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/tool"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"os/exec"
)

// Entry is real main entry for this app
func Entry() {
	// logrus.SetLevel(logrus.DebugLevel)
	// logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})

	// defer func() {
	// 	fmt.Println("defer caller")
	// 	if err := recover(); err != nil {
	// 		fmt.Printf("recover success. error: %v", err)
	// 	}
	// }()

	if err := cmdr.Exec(buildRootCmd(),

		// To disable internal commands and flags, uncomment the following codes
		// cmdr.WithBuiltinCommands(false, false, false, false, false),

		//cmdr.WithHelpTabStop(41),
		// daemon.WithDaemon(svr.NewDaemon(), nil, nil, nil),

		// integrate with logex library
		//cmdr.WithLogex(cmdr.DebugLevel),
		//cmdr.WithLogexPrefix("logger"),
		//cmdr.WithLogx(build.New(cmdr.NewLoggerConfigWith(true, "logrus", "debug"))),
		cmdr.WithLogxShort(true, "logrus", "debug"),

		cmdr.WithWatchMainConfigFileToo(true),
		// cmdr.WithNoWatchConfigFiles(false),

		cmdr.WithOptionMergeModifying(onOptionMergeModifying),
		cmdr.WithUnknownOptionHandler(onUnknownOptionHandler),
		cmdr.WithUnhandledErrorHandler(onUnhandledErrorHandler),
		cmdr.WithOnSwitchCharHit(onSwitchCharHit),
		cmdr.WithOnPassThruCharHit(onPassThruCharHit),

		optAddTraceOption,
		optAddServerExtOption,
	); err != nil {
		cmdr.Logger.Fatalf("error: %v", err)
	}
}

func buildRootCmd() (rootCmd *cmdr.RootCommand) {

	// var cmd *Command

	// cmdr.Root("aa", "1.0.1").
	// 	Header("sds").
	// 	NewSubCommand().
	// 	Titles("ms", "microservice").
	// 	Description("", "").
	// 	Group("").
	// 	Action(func(cmd *cmdr.Command, args []string) (err error) {
	// 		return
	// 	})

	// root

	root := cmdr.Root(appName, cmdr.Version).
		Header("fluent - test for cmdr - no version - hedzr").
		Description(desc, longDesc).
		Examples(examples)
	rootCmd = root.RootCommand()

	root.NewSubCommand("cols", "rows", "tty-size").
		Description("detected tty size").
		Group("Test").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			cols, rows := tool.GetTtySize()
			fmt.Printf(" 1. cols = %v, rows = %v\n\n", cols, rows)

			cols, rows, err = terminal.GetSize(int(os.Stdout.Fd()))
			fmt.Printf(" 2. cols = %v, rows = %v | in-docker: %v\n\n", cols, rows, cmdr.InDockerEnv())

			var out []byte
			cc := exec.Command("stty", "size")
			cc.Stdin = os.Stdin
			out, err = cc.Output()
			fmt.Printf(" 3. out: %v", string(out))
			fmt.Printf("    err: %v\n", err)

			if cmdr.InDockerEnv() {
				//
			}
			return
		})

	// soundex

	root.NewSubCommand("soundex", "snd", "sndx", "sound").
		Description("soundex test").
		Group("Test").
		TailPlaceholder("[text1, text2, ...]").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			for ix, s := range args {
				fmt.Printf("%5d. %s => %s\n", ix, s, tool.Soundex(s))
			}
			return
		})

	// panic test

	pa := root.NewSubCommand("panic-test", "pa").
		Description("test panic inside cmdr actions", "").
		Group("Test")

	val := 9
	zeroVal := zero

	pa.NewSubCommand("division-by-zero", "dz").
		Description("").
		Group("Test").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			fmt.Println(val / zeroVal)
			return
		})

	pa.NewSubCommand("panic", "pa").
		Description("").
		Group("Test").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			panic(9)
			return
		})

	// kb-print

	kb := root.NewSubCommand("kb-print", "kb").
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
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			fmt.Printf("Got size: %v (literal: %v)\n\n", cmdr.GetKibibytesR("kb-print.size"), cmdr.GetStringR("kb-print.size"))
			return
		})

	kb.NewFlagV("1k", "size", "s").
		Description("max message size. Valid formats: 2k, 2kb, 2kB, 2KB. Suffixes: k, m, g, t, p, e.", "").
		Group("")

	// xy-print

	root.NewSubCommand("xy-print", "xy").
		Description("test terminal control sequences", "test terminal control sequences,\nverbose long descriptions here.").
		Group("Test").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			//
			// https://en.wikipedia.org/wiki/ANSI_escape_code
			// https://zh.wikipedia.org/wiki/ANSI%E8%BD%AC%E4%B9%89%E5%BA%8F%E5%88%97
			// https://en.wikipedia.org/wiki/POSIX_terminal_interface
			//

			fmt.Println("\x1b[2J") // clear screen

			for i, s := range args {
				fmt.Printf("\x1b[s\x1b[%d;%dH%s\x1b[u", 15+i, 30, s)
			}

			return
		})

	// mx-test

	mx := root.NewSubCommand("mx-test", "mx").
		Description("test new features", "test new features,\nverbose long descriptions here.").
		Group("Test").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			// cmdr.Set("test.1", 8)
			cmdr.Set("test.deep.branch.1", "test")
			z := cmdr.GetString("app.test.deep.branch.1")
			fmt.Printf("*** Got app.test.deep.branch.1: %s\n", z)
			if z != "test" {
				cmdr.Logger.Fatalf("err, expect 'test', but got '%v'", z)
			}

			cmdr.DeleteKey("app.test.deep.branch.1")
			if cmdr.HasKey("app.test.deep.branch.1") {
				cmdr.Logger.Fatalf("FAILED, expect key not found, but found: %v", cmdr.Get("app.test.deep.branch.1"))
			}
			fmt.Printf("*** Got app.test.deep.branch.1 (after deleted): %s\n", cmdr.GetString("app.test.deep.branch.1"))

			fmt.Printf("*** Got pp: %s\n", cmdr.GetString("app.mx-test.password"))
			fmt.Printf("*** Got msg: %s\n", cmdr.GetString("app.mx-test.message"))
			fmt.Printf("*** Got fruit (toggle group): %v\n", cmdr.GetString("app.mx-test.fruit"))
			fmt.Printf("*** Got head (head-like): %v\n", cmdr.GetInt("app.mx-test.head"))
			fmt.Println()
			fmt.Printf("*** test text: %s\n", cmdr.GetStringR("mx-test.test"))
			fmt.Println()
			fmt.Printf("> InTesting: args[0]=%v \n", tool.SavedOsArgs[0])
			fmt.Println()
			fmt.Printf("> Used config file: %v\n", cmdr.GetUsedConfigFile())
			fmt.Printf("> Used config files: %v\n", cmdr.GetUsingConfigFiles())
			fmt.Printf("> Used config sub-dir: %v\n", cmdr.GetUsedConfigSubDir())

			fmt.Printf("> STDIN MODE: %v \n", cmdr.GetBoolR("mx-test.stdin"))
			fmt.Println()

			//logrus.Debug("debug")
			//logrus.Info("debug")
			//logrus.Warning("debug")
			//logrus.WithField(logex.SKIP, 1).Warningf("dsdsdsds")

			if cmdr.GetBoolR("mx-test.stdin") {
				fmt.Println("> Type your contents here, press Ctrl-D to end it:")
				var data []byte
				data, err = ioutil.ReadAll(os.Stdin)
				if err != nil {
					cmdr.Logger.Printf("error: %+v", err)
					return
				}
				fmt.Println("> The input contents are:")
				fmt.Print(string(data))
				fmt.Println()
			}
			return
		})
	mx.NewFlagV("", "test", "t").
		Description("the test text.", "").
		EnvKeys("COOLT", "TEST").
		Group("")
	mx.NewFlagV("", "password", "pp").
		Description("the password requesting.", "").
		Group("").
		Placeholder("PASSWORD").
		ExternalTool(cmdr.ExternalToolPasswordInput)
	mx.NewFlagV("", "message", "m", "msg").
		Description("the message requesting.", "").
		Group("").
		Placeholder("MESG").
		ExternalTool(cmdr.ExternalToolEditor)
	mx.NewFlagV("", "fruit", "fr").
		Description("the message.", "").
		Group("").
		Placeholder("FRUIT").
		ValidArgs("apple", "banana", "orange")
	mx.NewFlagV(1, "head", "hd").
		Description("the head lines.", "").
		Group("").
		Placeholder("LINES").
		HeadLike(true, 1, 3000)
	mx.NewFlagV(false, "stdin", "c").
		Description("read file content from stdin.", "").
		Group("")

	// kv

	kvCmd := root.NewSubCommand("kvstore", "kv").
		Description("consul kv store operations...", ``)

	attachConsulConnectFlags(kvCmd)

	kvBackupCmd := kvCmd.NewSubCommand("backup", "b", "bf", "bkp").
		Description("Dump Consul's KV database to a JSON/YAML file", ``).
		Action(kvBackup)
	kvBackupCmd.NewFlagV("consul-backup.json", "output", "o").
		Description("Write output to a file (*.json / *.yml)", ``).
		Placeholder("FILE")

	kvRestoreCmd := kvCmd.NewSubCommand("restore", "r").
		Description("restore to Consul's KV store, from a a JSON/YAML backup file", ``).
		Action(kvRestore)
	kvRestoreCmd.NewFlagV("consul-backup.json", "input", "i").
		Description("Read the input file (*.json / *.yml)", ``).
		Placeholder("FILE")

	// ms

	msCmd := root.NewSubCommand("micro-service", "ms", "microservice").
		Description("micro-service operations...", "").
		Group("")

	msCmd.NewFlagV(false, "money", "mm").
		Description("A placeholder flag.", "").
		Group("").
		Placeholder("")

	msCmd.NewFlagV("", "name", "n").
		Description("name of the service", ``).
		Placeholder("NAME")
	msCmd.NewFlagV("", "id", "i", "ID").
		Description("unique id of the service", ``).
		Placeholder("ID")
	msCmd.NewFlagV(false, "all", "a").
		Description("all services", ``).
		Placeholder("")

	msCmd.NewFlagV(3, "retry", "t").
		Description("", "").
		Group("").
		Placeholder("RETRY")

	// ms ls

	msCmd.NewSubCommand("list", "ls", "l", "lst", "dir").
		Description("list tags", "").
		Group("2333.List").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	// ms tags

	msTagsCmd := msCmd.NewSubCommand("tags", "t").
		Description("tags operations of a micro-service", "").
		Group("")

	// cTags.NewFlag(cmdr.OptFlagTypeString).
	// 	Titles("n", "name").
	// 	Description("name of the service", "").
	// 	Group("").
	// 	DefaultValue("", "NAME")
	//
	// cTags.NewFlag(cmdr.OptFlagTypeString).
	// 	Titles("i", "id").
	// 	Description("unique id of the service", "").
	// 	Group("").
	// 	DefaultValue("", "ID")
	//
	// cTags.NewFlag(cmdr.OptFlagTypeString).
	// 	Titles("a", "addr").
	// 	Description("", "").
	// 	Group("").
	// 	DefaultValue("consul.ops.local", "ADDR")

	attachConsulConnectFlags(msTagsCmd)

	// ms tags ls

	msTagsCmd.NewSubCommand("list", "ls", "l", "lst", "dir").
		Description("list tags").
		Group("2333.List").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	// ms tags add

	tagsAdd := msTagsCmd.NewSubCommand("add", "a", "new", "create").
		Description("add tags").
		Deprecated("0.2.1").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	tagsAdd.NewFlagV([]string{}, "list", "ls", "l", "lst", "dir").
		Description("a comma list to be added").
		Group("").
		Placeholder("LIST")

	c1 := tagsAdd.NewSubCommand("check", "c", "chk").
		Description("[sub] check").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	c2 := c1.NewSubCommand("check-point", "pt", "chk-pt").
		Description("[sub][sub] checkpoint").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	c2.NewFlagV([]string{}, "add", "a", "add-list").
		Description("a comma list to be added.").
		Placeholder("LIST").
		Group("List")
	c2.NewFlagV([]string{}, "remove", "r", "rm-list", "rm", "del", "delete").
		Description("a comma list to be removed.", ``).
		Placeholder("LIST").
		Group("List")

	c3 := c1.NewSubCommand("check-in", "in", "chk-in").
		Description("[sub][sub] check-in").
		Group("")

	c3.NewFlag(cmdr.OptFlagTypeString).
		Titles("n", "name").
		Description("a string to be added.").
		DefaultValue("", "")

	c3.NewSubCommand("demo-1", "d1").
		Description("[sub][sub] check-in sub").
		Group("")

	c3.NewSubCommand("demo-2", "d2").
		Description("[sub][sub] check-in sub").
		Group("")

	c3.NewSubCommand("demo-3", "d3").
		Description("[sub][sub] check-in sub").
		Group("")

	c1.NewSubCommand("check-out", "out", "chk-out").
		Description("[sub][sub] check-out").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	// ms tags rm

	tagsRm := msTagsCmd.NewSubCommand("rm", "r", "remove", "delete", "del", "erase").
		Description("remove tags").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	tagsRm.NewFlagV([]string{}, "list", "ls", "l", "lst", "dir").
		Description("a comma list to be added").
		Group("").
		Placeholder("LIST")

	// ms tags modify

	msTagsModifyCmd := msTagsCmd.NewSubCommand("modify", "m", "mod", "modi", "update", "change").
		Description("modify tags of a service.").
		Action(msTagsModify)

	attachModifyFlags(msTagsModifyCmd)

	msTagsModifyCmd.NewFlagV([]string{}, "add", "a", "add-list").
		Description("a comma list to be added.").
		Placeholder("LIST").
		Group("List")
	msTagsModifyCmd.NewFlagV([]string{}, "remove", "r", "rm-list", "rm", "del", "delete").
		Description("a comma list to be removed.").
		Placeholder("LIST").
		Group("List")

	// ms tags toggle

	tagsTog := msTagsCmd.NewSubCommand("toggle", "t", "tog", "switch").
		Description("toggle tags").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	attachModifyFlags(tagsTog)

	tagsTog.NewFlagV([]string{}, "set", "s").
		Description("a comma list to be set").
		Group("").
		Placeholder("LIST")

	tagsTog.NewFlagV([]string{}, "unset", "un").
		Description("a comma list to be unset").
		Group("").
		Placeholder("LIST")

	tagsTog.NewFlagV("", "address", "a", "addr").
		Description("the address of the service (by id or name)").
		Placeholder("HOST:PORT")

	return
}

const (
	appName   = "fluent"
	copyright = "fluent is an effective devops tool"
	desc      = "fluent is an effective devops tool. It make an demo application for `cmdr`."
	longDesc  = "fluent is an effective devops tool. It make an demo application for `cmdr`."
	examples  = `
$ {{.AppName}} gen shell [--bash|--zsh|--auto]
  generate bash/shell completion scripts
$ {{.AppName}} gen man
  generate linux man page 1
$ {{.AppName}} --help
  show help screen.
`
	overview = ``

	zero = 0
)
