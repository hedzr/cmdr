//go:build !delve && !test
// +build !delve,!test

// go/:build delve || test
// +/build delve test

package cmdr

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/tool"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func cmdrMoreCommandsForTest(root cmdr.OptCmd) {

	// test/debug build, many multilevel subcommands here

	cmdrXyPrint(root)
	cmdrKbPrint(root)
	cmdrSoundex(root)
	cmdrTtySize(root)

	tgCommand(root)
	mxCommand(root)

	cmdr.NewBool(false).
		Titles("enable-ueh", "ueh").
		Description("Enables the unhandled exception handler?").
		AttachTo(root)
	cmdrPanic(root)

	cmdrMultiLevelTest(root)
	cmdrManyCommandsTest(root)

	// pprof.AttachToCmdr(root.RootCmdOpt())
}

func tgCommand(root cmdr.OptCmd) {

	// toggle-group-test - without a default choice

	fx := cmdr.NewSubCmd().
		Titles("tg-test", "tg", "toggle-group-test").
		Description("tg test new features", "tg test new features,\nverbose long descriptions here.").
		Group("Test").
		Action(func(cmd *cmdr.Command, args []string) (err error) {

			fmt.Printf("*** Got fruit (toggle group): %v\n", cmdr.GetString("app.tg-test.fruit"))

			fmt.Printf("> STDIN MODE: %v \n", cmdr.GetBoolR("mx-test.stdin"))
			fmt.Println()

			// logrus.Debug("debug")
			// logrus.Info("debug")
			// logrus.Warning("debug")
			// logrus.WithField(logex.SKIP, 1).Warningf("dsdsdsds")

			return
		}).
		AttachTo(root)
	cmdr.NewBool().Titles("apple", "").
		Description("the test text.", "").
		ToggleGroup("fruit").
		AttachTo(fx)
	cmdr.NewBool().Titles("banana", "").
		Description("the test text.", "").
		ToggleGroup("fruit").
		AttachTo(fx)
	cmdr.NewBool(true).Titles("orange", "").
		Description("the test text.", "").
		ToggleGroup("fruit").
		AttachTo(fx)

	// tg2 - with a default choice

	fx2 := cmdr.NewSubCmd().
		Titles("tg-test2", "tg2", "toggle-group-test2").
		Description("tg2 test new features", "tg2 test new features,\nverbose long descriptions here.").
		Group("Test").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			fmt.Printf("*** Got fruit (toggle group): %v\n", cmdr.GetString("app.tg-test2.fruit"))

			fmt.Printf("> STDIN MODE: %v \n", cmdr.GetBoolR("mx-test.stdin"))
			fmt.Println()
			return
		}).
		AttachTo(root)
	cmdr.NewBool().Titles("apple", "").
		Description("the test text.", "").
		ToggleGroup("fruit").
		AttachTo(fx2)
	cmdr.NewBool().Titles("banana", "").
		Description("the test text.", "").
		ToggleGroup("fruit").
		AttachTo(fx2)
	cmdr.NewBool().Titles("orange", "").
		Description("the test text.", "").
		ToggleGroup("fruit").
		AttachTo(fx2)

}

func mxCommand(root cmdr.OptCmd) {

	// mx-test

	mx := cmdr.NewSubCmd().
		Titles("mx-test", "mx").
		Description("mx test new features", "mx test new features,\nverbose long descriptions here.").
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
			fmt.Printf("*** Got fruit (valid args): %v\n", cmdr.GetString("app.mx-test.fruit"))
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

			// logrus.Debug("debug")
			// logrus.Info("debug")
			// logrus.Warning("debug")
			// logrus.WithField(logex.SKIP, 1).Warningf("dsdsdsds")

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
		}).
		AttachTo(root)
	cmdr.NewString().Titles("test", "t").
		Description("the test text.", "").
		EnvKeys("COOLT", "TEST").
		Group("").
		AttachTo(mx)

	cmdr.NewString().Titles("password", "pp").
		Description("the password requesting.", "").
		Group("").
		Placeholder("PASSWORD").
		ExternalTool(cmdr.ExternalToolPasswordInput).
		AttachTo(mx)

	cmdr.NewString().Titles("message", "m", "msg").
		Description("the message requesting.", "").
		Group("").
		Placeholder("MESG").
		ExternalTool(cmdr.ExternalToolEditor).
		AttachTo(mx)

	cmdr.NewString().Titles("fruit", "fr").
		Description("the message.", "").
		Group("").
		Placeholder("FRUIT").
		ValidArgs("apple", "banana", "orange").
		AttachTo(mx)

	cmdr.NewInt(1).Titles("head", "hd").
		Description("the head lines.", "").
		Group("").
		Placeholder("LINES").
		HeadLike(true, 1, 3000).
		AttachTo(mx)

	cmdr.NewBool().Titles("stdin", "c").
		Description("read file content from stdin.", "").
		Group("").
		AttachTo(mx)

}

func cmdrXyPrint(root cmdr.OptCmd) {

	// xy-print

	cmdr.NewSubCmd().Titles("xy-print", "xy").
		Description("test terminal control sequences", "xy-print test terminal control sequences,\nverbose long descriptions here.").
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
		}).
		AttachTo(root)

}

func cmdrKbPrint(root cmdr.OptCmd) {

	// kb-print

	kb := cmdr.NewSubCmd().Titles("kb-print", "kb").
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
		}).
		AttachTo(root)

	cmdr.NewString("1k").Titles("size", "s").
		Description("max message size. Valid formats: 2k, 2kb, 2kB, 2KB. Suffixes: k, m, g, t, p, e.", "").
		Group("").
		AttachTo(kb)

}

func cmdrPanic(root cmdr.OptCmd) {
	// panic test

	pa := cmdr.NewSubCmd().
		Titles("panic-test", "pa").
		Description("test panic inside cmdr actions", "").
		Group("Test").
		AttachTo(root)

	val := 9
	zeroVal := zero

	cmdr.NewSubCmd().
		Titles("division-by-zero", "dz").
		Description("").
		Group("Test").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			fmt.Println(val / zeroVal)
			return
		}).
		AttachTo(pa)

	cmdr.NewSubCmd().
		Titles("panic", "pa").
		Description("").
		Group("Test").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			panic(9)
			return
		}).
		AttachTo(pa)

}

func cmdrSoundex(root cmdr.OptCmd) {

	cmdr.NewSubCmd().Titles("soundex", "snd", "sndx", "sound").
		Description("soundex test").
		Group("Test").
		TailPlaceholder("[text1, text2, ...]").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			for ix, s := range args {
				fmt.Printf("%5d. %s => %s\n", ix, s, tool.Soundex(s))
			}
			return
		}).
		AttachTo(root)

}

func cmdrTtySize(root cmdr.OptCmd) {

	cmdr.NewSubCmd().Titles("cols", "rows", "tty-size").
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
		}).
		AttachTo(root)

}

func cmdrManyCommandsTest(root cmdr.OptCmd) {
	for i := 1; i <= 23; i++ {
		t := fmt.Sprintf("subcmd-%v", i)
		var ttls []string
		for o, l := root, cmdr.OptCmd(nil); o != nil && o != l; {
			ttls = append(ttls, o.ToCommand().GetTitleName())
			l, o = o, o.OwnerCommand()
		}
		ttl := strings.Join(tool.ReverseStringSlice(ttls), ".")

		c := cmdr.NewSubCmd().
			Titles(t, fmt.Sprintf("sc%v", i)).
			Description(fmt.Sprintf("subcommands %v.sc%v test", ttl, i)).
			Group("Test").
			AttachTo(root)
		cmdrAddFlags(c)
	}
}

func cmdrMultiLevelTest(root cmdr.OptCmd) {

	cmd := cmdr.NewSubCmd().
		Titles("mls", "mls").
		Description("multi-level subcommands test").
		Group("Test").
		// Sets(func(cmd cmdr.OptCmd) {
		//	cmdrAddFlags(cmd)
		// }).
		AttachTo(root)
	// cmd := root.NewSubCommand("mls", "mls").
	//	Description("multi-level subcommands test").
	//	Group("Test")
	cmdrAddFlags(cmd)
	cmdrMultiLevel(cmd, 1)

}

func cmdrMultiLevel(parent cmdr.OptCmd, depth int) {
	if depth > 3 {
		return
	}

	for i := 1; i < 4; i++ {
		t := fmt.Sprintf("subcmd-%v", i)
		var ttls []string
		for o, l := parent, cmdr.OptCmd(nil); o != nil && o != l; {
			ttls = append(ttls, o.ToCommand().GetTitleName())
			l, o = o, o.OwnerCommand()
		}
		ttl := strings.Join(tool.ReverseStringSlice(ttls), ".")

		cc := cmdr.NewSubCmd().
			Titles(t, fmt.Sprintf("sc%v", i)).
			// cc := parent.NewSubCommand(t, fmt.Sprintf("sc%v", i)).
			Description(fmt.Sprintf("subcommands %v.sc%v test", ttl, i)).
			Group("Test").
			AttachTo(parent)
		cmdrAddFlags(cc)
		cmdrMultiLevel(cc, depth+1)
	}
}

func cmdrAddFlags(c cmdr.OptCmd) {

	cmdr.NewBool().Titles("apple", "").
		Description("the test text.", "").
		ToggleGroup("fruit").
		AttachTo(c)
	cmdr.NewBool().Titles("banana", "").
		Description("the test text.", "").
		ToggleGroup("fruit").
		AttachTo(c)
	cmdr.NewBool().Titles("orange", "").
		Description("the test text.", "").
		ToggleGroup("fruit").
		AttachTo(c)

	cmdr.NewString().Titles("message", "m", "msg").
		Description("the message requesting.", "").
		Group("").
		Placeholder("MESG").
		ExternalTool(cmdr.ExternalToolEditor).
		AttachTo(c)

	cmdr.NewString().Titles("fruit", "fr").
		Description("the message.", "").
		Group("").
		Placeholder("FRUIT").
		ValidArgs("apple", "banana", "orange").
		AttachTo(c)

	cmdr.NewBool(false).
		Titles("bool", "b").
		Description("A bool flag", "").
		Group("").
		AttachTo(c)

	cmdr.NewInt(1).
		Titles("int", "i").
		Description("A int flag", "").
		Group("1000.Integer").
		AttachTo(c)
	cmdr.NewInt64(2).
		Titles("int64", "i64").
		Description("A int64 flag", "").
		Group("1000.Integer").
		AttachTo(c)
	cmdr.NewUint(3).
		Titles("uint", "u").
		Description("A uint flag", "").
		Group("1000.Integer").
		AttachTo(c)
	cmdr.NewUint64(4).
		Titles("uint64", "u64").
		Description("A uint64 flag", "").
		Group("1000.Integer").
		AttachTo(c)

	cmdr.NewFloat32(2.71828).
		Titles("float32", "f", "float").
		Description("A float32 flag with 'e' value", "").
		Group("2000.Float").
		AttachTo(c)
	cmdr.NewFloat64(3.14159265358979323846264338327950288419716939937510582097494459230781640628620899).
		Titles("float64", "f64").
		Description("A float64 flag with a `PI` value", "").
		Group("2000.Float").
		AttachTo(c)
	cmdr.NewComplex64(3.14+9i).
		Titles("complex64", "c64").
		Description("A complex64 flag", "").
		Group("2010.Complex").
		AttachTo(c)
	cmdr.NewComplex64(3.14+9i).
		Titles("complex128", "c128").
		Description("A complex128 flag", "").
		Group("2010.Complex").
		AttachTo(c)

	// a set of booleans

	cmdr.NewBool(false).
		Titles("single", "s").
		Description("A bool flag: single", "").
		Group("Boolean").
		EnvKeys("").
		AttachTo(c)

	cmdr.NewBool(false).
		Titles("double", "d").
		Description("A bool flag: double", "").
		Group("Boolean").
		EnvKeys("").
		AttachTo(c)

	cmdr.NewBool(false).
		Titles("norway", "n", "nw").
		Description("A bool flag: norway", "").
		Group("Boolean").
		EnvKeys("").
		AttachTo(c)

	cmdr.NewBool(false).
		Titles("mongo", "mongo").
		Description("A bool flag: mongo", "").
		Group("Boolean").
		EnvKeys("").
		AttachTo(c)

}
