/*
 * Copyright © 2019 Hedzr Yeh.
 */

package main

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/examples/demo/svr"
	"github.com/hedzr/cmdr/plugin/daemon"
	"github.com/sirupsen/logrus"
	"log"
)

var optAddTraceOption cmdr.ExecOption

func init() {
	// attaches `--trace` to root command
	optAddTraceOption = cmdr.WithXrefBuildingHooks(func(root *cmdr.RootCommand, args []string) {
		cmdr.NewBool().
			Titles("tr", "trace").
			Description("enable trace mode for tcp/mqtt send/recv data dump", "").
			DefaultValue(false, "").
			AttachToRoot(root)
	}, nil)
}

func main() {
	// logrus.SetLevel(logrus.DebugLevel)
	// logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})

	if err := cmdr.Exec(buildRootCmd(),

		// To disable internal commands and flags, uncomment the following codes
		// cmdr.WithBuiltinCommands(false, false, false, false, false),
		daemon.WithDaemon(svr.NewDaemon(), nil, nil, nil),

		// integrate with logex library
		cmdr.WithLogex(logrus.DebugLevel),
		cmdr.WithLogexPrefix("logger"),

		cmdr.WithHelpTabStop(50),

		optAddTraceOption,
	); err != nil {
		logrus.Fatal(err)
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

	root := cmdr.Root(appName, "1.0.1").
		Header("fluent - test for cmdr - no version - hedzr").
		Description(desc, longDesc).
		Examples(examples)
	rootCmd = root.RootCommand()

	// soundex

	root.NewSubCommand().
		Titles("snd", "soundex", "sndx", "sound").
		Description("", "soundex test").
		Group("Test").
		TailPlaceholder("[text1, text2, ...]").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			for ix, s := range args {
				fmt.Printf("%5d. %s => %s\n", ix, s, cmdr.Soundex(s))
			}
			return
		})

	// xy-print

	root.NewSubCommand().
		Titles("xy", "xy-print").
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

	mx := root.NewSubCommand().
		Titles("mx", "mx-test").
		Description("test new features", "test new features,\nverbose long descriptions here.").
		Group("Test").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			// cmdr.Set("test.1", 8)
			cmdr.Set("test.deep.branch.1", "test")
			z := cmdr.GetString("app.test.deep.branch.1")
			fmt.Printf("*** Got app.test.deep.branch.1: %s\n", z)
			if z != "test" {
				log.Fatalf("err, expect 'test', but got '%v'", z)
			}

			cmdr.DeleteKey("app.test.deep.branch.1")
			if cmdr.HasKey("app.test.deep.branch.1") {
				log.Fatalf("FAILED, expect key not found, but found: %v", cmdr.Get("app.test.deep.branch.1"))
			}
			fmt.Printf("*** Got app.test.deep.branch.1 (after deleted): %s\n", cmdr.GetString("app.test.deep.branch.1"))

			fmt.Printf("*** Got pp: %s\n", cmdr.GetString("app.mx-test.password"))
			fmt.Printf("*** Got msg: %s\n", cmdr.GetString("app.mx-test.message"))
			fmt.Printf("*** Got fruit (toggle group): %v\n", cmdr.GetString("app.mx-test.fruit"))
			fmt.Printf("*** Got head (head-like): %v\n", cmdr.GetInt("app.mx-test.head"))
			fmt.Println()
			fmt.Printf("InTesting: args[0]=%v ", cmdr.SavedOsArgs[0])
			return
		})
	mx.NewFlagV("").
		Titles("pp", "password").
		Description("the password requesting.", "").
		Group("").
		Placeholder("PASSWORD").
		ExternalTool(cmdr.ExternalToolPasswordInput)
	mx.NewFlagV("").
		Titles("m", "message", "msg").
		Description("the message requesting.", "").
		Group("").
		Placeholder("MESG").
		ExternalTool(cmdr.ExternalToolEditor)
	mx.NewFlagV("").
		Titles("fr", "fruit").
		Description("the message.", "").
		Group("").
		Placeholder("FRUIT").
		ValidArgs("apple", "banana", "orange")
	mx.NewFlagV(1).
		Titles("hd", "head").
		Description("the head lines.", "").
		Group("").
		Placeholder("LINES").
		HeadLike(true, 1, 3000)

	// kv

	kvCmd := root.NewSubCommand().
		Titles("kv", "kvstore").
		Description("consul kv store operations...", ``)

	attachConsulConnectFlags(kvCmd)

	kvBackupCmd := kvCmd.NewSubCommand().
		Titles("b", "backup", "bk", "bf", "bkp").
		Description("Dump Consul's KV database to a JSON/YAML file", ``).
		Action(kvBackup)
	kvBackupCmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("o", "output").
		Description("Write output to a file (*.json / *.yml)", ``).
		DefaultValue("consul-backup.json", "FILE")

	kvRestoreCmd := kvCmd.NewSubCommand().
		Titles("r", "restore").
		Description("restore to Consul's KV store, from a a JSON/YAML backup file", ``).
		Action(kvRestore)
	kvRestoreCmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("i", "input").
		Description("Read the input file (*.json / *.yml)", ``).
		DefaultValue("consul-backup.json", "FILE")

	// ms

	msCmd := root.NewSubCommand().
		Titles("ms", "micro-service", "microservice").
		Description("micro-service operations...", "").
		Group("")

	msCmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("mm", "money").
		Description("A placeholder flag.", "").
		Group("").
		DefaultValue(false, "")

	msCmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("n", "name").
		Description("name of the service", ``).
		DefaultValue("", "NAME")
	msCmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("i", "id", "ID").
		Description("unique id of the service", ``).
		DefaultValue("", "ID")
	msCmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("a", "all").
		Description("all services", ``).
		DefaultValue(false, "")

	msCmd.NewFlag(cmdr.OptFlagTypeUint).
		Titles("t", "retry").
		Description("", "").
		Group("").
		DefaultValue(3, "RETRY")

	// ms ls

	msCmd.NewSubCommand().
		Titles("ls", "list", "l", "lst", "dir").
		Description("list tags", "").
		Group("2333.List").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	// ms tags

	msTagsCmd := msCmd.NewSubCommand().
		Titles("t", "tags").
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

	msTagsCmd.NewSubCommand().
		Titles("ls", "list", "l", "lst", "dir").
		Description("list tags", "").
		Group("2333.List").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	// ms tags add

	tagsAdd := msTagsCmd.NewSubCommand().
		Titles("a", "add", "new", "create").
		Description("add tags", "").
		Deprecated("0.2.1").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	tagsAdd.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("ls", "list", "l", "lst", "dir").
		Description("a comma list to be added", "").
		Group("").
		DefaultValue([]string{}, "LIST")

	c1 := tagsAdd.NewSubCommand().
		Titles("c", "check", "chk").
		Description("[sub] check", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	c2 := c1.NewSubCommand().
		Titles("pt", "check-point", "chk-pt").
		Description("[sub][sub] checkpoint", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	c2.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("a", "add", "add-list").
		Description("a comma list to be added.", ``).
		DefaultValue([]string{}, "LIST").
		Group("List")
	c2.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("r", "remove", "rm-list", "rm", "del", "delete").
		Description("a comma list to be removed.", ``).
		DefaultValue([]string{}, "LIST").
		Group("List")

	c3 := c1.NewSubCommand().
		Titles("in", "check-in", "chk-in").
		Description("[sub][sub] check-in", "").
		Group("")

	c3.NewFlag(cmdr.OptFlagTypeString).
		Titles("n", "name").
		Description("a string to be added.", ``).
		DefaultValue("", "")

	c3.NewSubCommand().
		Titles("d1", "demo-1").
		Description("[sub][sub] check-in sub", "").
		Group("")

	c3.NewSubCommand().
		Titles("d2", "demo-2").
		Description("[sub][sub] check-in sub", "").
		Group("")

	c3.NewSubCommand().
		Titles("d3", "demo-3").
		Description("[sub][sub] check-in sub", "").
		Group("")

	c1.NewSubCommand().
		Titles("out", "check-out", "chk-out").
		Description("[sub][sub] check-out", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	// ms tags rm

	tagsRm := msTagsCmd.NewSubCommand().
		Titles("r", "rm", "remove", "delete", "del", "erase").
		Description("remove tags", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	tagsRm.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("ls", "list", "l", "lst", "dir").
		Description("a comma list to be added", "").
		Group("").
		DefaultValue([]string{}, "LIST")

	// ms tags modify

	msTagsModifyCmd := msTagsCmd.NewSubCommand().
		Titles("m", "modify", "mod", "modi", "update", "change").
		Description("modify tags of a service.", ``).
		Action(msTagsModify)

	attachModifyFlags(msTagsModifyCmd)

	msTagsModifyCmd.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("a", "add", "add-list").
		Description("a comma list to be added.", ``).
		DefaultValue([]string{}, "LIST").
		Group("List")
	msTagsModifyCmd.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("r", "remove", "rm-list", "rm", "del", "delete").
		Description("a comma list to be removed.", ``).
		DefaultValue([]string{}, "LIST").
		Group("List")

	// ms tags toggle

	tagsTog := msTagsCmd.NewSubCommand().
		Titles("t", "toggle", "tog", "switch").
		Description("toggle tags", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	attachModifyFlags(tagsTog)

	tagsTog.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("s", "set").
		Description("a comma list to be set", "").
		Group("").
		DefaultValue([]string{}, "LIST")

	tagsTog.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("u", "unset", "un").
		Description("a comma list to be unset", "").
		Group("").
		DefaultValue([]string{}, "LIST")

	tagsTog.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("a", "address", "addr").
		Description("the address of the service (by id or name)", ``).
		DefaultValue("", "HOST:PORT")

	return
}

func attachModifyFlags(cmd cmdr.OptCmd) {
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("d", "delim").
		Description("delimitor char in `non-plain` mode.", ``).
		DefaultValue("=", "")

	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("c", "clear").
		Description("clear all tags.", ``).
		DefaultValue(false, "").
		Group("Operate")

	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("g", "string", "string-mode").
		Description("In 'String Mode', default will be disabled: default, a tag string will be split by comma(,), and treated as a string list.", ``).
		DefaultValue(false, "").
		Group("Mode")

	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("m", "meta", "meta-mode").
		Description("In 'Meta Mode', service 'NodeMeta' field will be updated instead of 'Tags'. (--plain assumed false).", ``).
		DefaultValue(false, "").
		Group("Mode")

	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("2", "both", "both-mode").
		Description("In 'Both Mode', both of 'NodeMeta' and 'Tags' field will be updated.", ``).
		DefaultValue(false, "").
		Group("Mode")

	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("p", "plain", "plain-mode").
		Description("In 'Plain Mode', a tag be NOT treated as `key=value` or `key:value`, and modify with the `key`.", ``).
		DefaultValue(false, "").
		Group("Mode")

	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("t", "tag", "tag-mode").
		Description("In 'Tag Mode', a tag be treated as `key=value` or `key:value`, and modify with the `key`.", ``).
		DefaultValue(true, "").
		Group("Mode")

}

func attachConsulConnectFlags(cmd cmdr.OptCmd) {
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("a", "addr").
		Description("Consul ip/host and port: HOST[:PORT] (No leading 'http(s)://')", ``).
		DefaultValue("localhost", "HOST[:PORT]").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeInt).
		Titles("p", "port").
		Description("Consul port", ``).
		DefaultValue(8500, "PORT").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("K", "insecure").
		Description("Skip TLS host verification", ``).
		DefaultValue(true, "").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("px", "prefix").
		Description("Root key prefix", ``).
		DefaultValue("/", "ROOT").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("", "cacert").
		Description("Consul Client CA cert)", ``).
		DefaultValue("", "FILE").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("", "cert").
		Description("Consul Client cert", ``).
		DefaultValue("", "FILE").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("", "scheme").
		Description("Consul connection protocol", ``).
		DefaultValue("http", "SCHEME").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("u", "username", "user", "usr", "uid").
		Description("HTTP Basic auth user", ``).
		DefaultValue("", "USERNAME").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("pw", "password", "passwd", "pass", "pwd").
		Description("HTTP Basic auth password", ``).
		DefaultValue("", "PASSWORD").
		Group("Consul").
		ExternalTool(cmdr.ExternalToolPasswordInput)

}

func kvBackup(cmd *cmdr.Command, args []string) (err error) {
	// err = consul.Backup()
	return
}

func kvRestore(cmd *cmdr.Command, args []string) (err error) {
	// err = consul.Restore()
	return
}

func msList(cmd *cmdr.Command, args []string) (err error) {
	// err = consul.ServiceList()
	return
}

func msTagsList(cmd *cmdr.Command, args []string) (err error) {
	// err = consul.TagsList()
	return
}

func msTagsAdd(cmd *cmdr.Command, args []string) (err error) {
	// err = consul.Tags()
	return
}

func msTagsRemove(cmd *cmdr.Command, args []string) (err error) {
	// err = consul.Tags()
	return
}

func msTagsModify(cmd *cmdr.Command, args []string) (err error) {
	// err = consul.Tags()
	return
}

func msTagsToggle(cmd *cmdr.Command, args []string) (err error) {
	// err = consul.TagsToggle()
	return
}

const (
	appName   = "fluent"
	copyright = "fluent.austr is an effective devops tool"
	desc      = "fluent.austr is an effective devops tool. It make an demo application for `cmdr`."
	longDesc  = "fluent.austr is an effective devops tool. It make an demo application for `cmdr`."
	examples  = `
$ {{.AppName}} gen shell [--bash|--zsh|--auto]
  generate bash/shell completion scripts
$ {{.AppName}} gen man
  generate linux man page 1
$ {{.AppName}} --help
  show help screen.
`
	overview = ``
)
