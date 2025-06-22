package hs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	errorsv3 "gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/store"

	"github.com/hedzr/cmdr/v2/conf"
)

func rootCmdForTesting() (app *appS, root *cli.RootCommand, err error) { //nolint:funlen,revive //for test
	app = &appS{
		Runner: newTestRunner(),
	}

	serverCommands := serverCommandsGet()

	kvCommands := kvCommandsGet()

	msCommands := msCommandsGet()

	root = &cli.RootCommand{
		AppName: "consul-tags",
		Version: "0.0.1",
		// app:     app,
		// Header:  `dsjlfsdjflsdfjlsdjflksjdfdsfsd`,
		// Version:    consul_tags.Version,
		// VersionInt: consul_tags.VersionInt,
		Copyright: "consul-tags is an effective devops tool",
		Author:    "Hedzr Yeh <hedzr@duck.com>",
		Cmd:       &cli.CmdS{},
	}
	root.SetApp(app)
	root.SetName(root.AppName)
	cc := root.Cmd.(*cli.CmdS)
	cc.SetName("consul-tags")

	ff := &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "t",
			Long:  "retry",
		},
	}
	ff.SetDesc("ss")
	ff.SetExamples(`random examples`)
	ff.SetDeprecated("1.2.3")
	ff.SetDefaultValue(1)
	ff.SetPlaceHolder("RETRY")
	err = cc.AddFlag(ff)
	if err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "s",
			Long:  "s",
		},
	}
	ff.SetOnMatchedHandler(func(flg *cli.Flag, position int, hitState *cli.MatchState) (err error) {
		if flg.GetDescZsh() != "ss" {
			err = errors.New("err `t`.GetDescZsh()")
		}
		if ttl, _ := flg.GetTitleNamesBy(","); ttl == "" {
			err = errors.New("err ss.GetTitleNamesBy()")
		}
		if len(flg.GetTitleZshFlagNamesArray()) != 2 {
			err = errors.New("err ss.GetTitleZshFlagNamesArray()")
		}
		// flg = cmd.Flags[1]
		// if flg.GetDescZsh() == "" {
		// 	err = errors.New("err sss.GetDescZsh()")
		// }
		// if flg.GetTitleZshNamesBy(",", false, false) == "" {
		// 	err = errors.New("err ss.GetTitleZshNamesBy()")
		// }
		// if len(flg.GetTitleZshFlagNamesArray()) != 2 {
		// 	err = errors.New("err ss.GetTitleZshFlagNamesArray()")
		// }
		// flg = cmd.Flags[2]
		// if flg.GetDescZsh() == "" {
		// 	err = errors.New("err ssss.GetDescZsh()")
		// }
		// if flg.GetTitleZshFlagName() == "" {
		// 	err = errors.New("err ss.GetTitleZshFlagName()")
		// }
		// if flg.GetTitleZshFlagShortName() == "" {
		// 	err = errors.New("err ss.GetTitleZshFlagShortName()")
		// }
		// if len(flg.GetTitleZshFlagNamesArray()) != 2 {
		// 	err = errors.New("err ss.GetTitleZshFlagNamesArray()")
		// }
		return
	})
	ff.SetDesc("s")
	ff.SetDefaultValue(uint(1))
	err = cc.AddFlag(ff)
	if err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "ff",
			Long:  "float",
		},
	}
	ff.SetDesc("")
	ff.SetDefaultValue(float64(1))
	err = cc.AddFlag(ff)
	if err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "pp",
			Long:  "spasswd",
		},
	}
	ff.SetDesc("")
	ff.SetDefaultValue("")
	ff.SetExternalEditor(cli.ExternalToolPasswordInput)
	ff.SetOnMatchedHandler(func(flg *cli.Flag, position int, hitState *cli.MatchState) (err error) {
		_, _ = fmt.Println("**** -pp action running")

		// f.owner.Runner.showVersions()
		// PrintBuildInfo()
		// cmd.PrintBuildInfo()
		// cmd.GetTitleZshNames()

		// SetCustomShowVersion(nil)
		// SetCustomShowBuildInfo(nil)
		_, _ = fmt.Println("**** -pp action end")
		return
	})
	err = cc.AddFlag(ff)
	if err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "qq",
			Long:  "qqpasswd",
		},
	}
	ff.SetDesc("")
	ff.SetDefaultValue("567")
	ff.SetExternalEditor(cli.ExternalToolEditor)
	err = cc.AddFlag(ff)
	if err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "dd",
			Long:  "ddduration",
		},
	}
	ff.SetDesc("")
	ff.SetDefaultValue(time.Second)
	ff.SetExternalEditor(cli.ExternalToolEditor)
	ff.SetOnChangingHandler(func(f *cli.Flag, oldVal, newVal any) (err error) {
		return
	})
	ff.SetOnChangedHandler(func(f *cli.Flag, oldVal, newVal any) {})
	ff.SetOnParseValueHandler(func(f *cli.Flag, position int, hitCaption string, hitValue string, moreArgs []string) (newVal any, remainPartInHitValue string, err error) {
		return
	})
	ff.SetOnSetHandler(func(f *cli.Flag, oldVal, newVal any) {})
	err = cc.AddFlag(ff)
	if err != nil {
		return
	}

	err = cc.AddSubCommand(serverCommands)
	if err != nil {
		return
	}
	err = cc.AddSubCommand(msCommands)
	if err != nil {
		return
	}
	err = cc.AddSubCommand(kvCommands)
	if err != nil {
		return
	}

	c1 := &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Short: "ls",
			Long:  "list",
		},
	}
	c1.SetDesc("list to Consul's KV store, from a a JSON/YAML backup file")

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "t",
			Long:  "retry",
		},
	}
	ff.SetDesc("ss")
	ff.SetExamples(`random examples`)
	ff.SetDeprecated("1.2.3")
	ff.SetDefaultValue(1)
	ff.SetPlaceHolder("RETRY")
	err = c1.AddFlag(ff)
	if err != nil {
		return
	}

	err = cc.AddSubCommand(c1)
	if err != nil {
		return
	}

	ctx := context.TODO()
	app.root = root
	root.Cmd.(*cli.CmdS).EnsureTree(ctx, app, root)
	root.Cmd.(*cli.CmdS).EnsureXref(ctx)
	return
}

//

//

//

func consulConnectFlagsGet() (flags []*cli.Flag) { //nolint:funlen,revive //for test
	c1 := &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Short: "ls",
			Long:  "list",
		},
	}

	flags = c1.Flags()

	var err error

	ff := &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "a",
			Long:  "addr",
		},
	}
	ff.SetDesc("Consul ip/host and port: HOST[:PORT] (No leading 'http(s)://')")
	ff.SetDefaultValue("consul.app.local")
	ff.SetPlaceHolder("HOST[:PORT]")
	err = c1.AddFlag(ff)
	if err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "p",
			Long:  "port",
		},
	}
	ff.SetDesc("Consul port")
	ff.SetDefaultValue(8500)
	ff.SetPlaceHolder("PORT")
	err = c1.AddFlag(ff)
	if err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "ui",
			Long:  "uint",
		},
	}
	ff.SetDesc("uint flag")
	ff.SetDefaultValue(uint(357))
	ff.SetPlaceHolder("NUM")
	err = c1.AddFlag(ff)
	if err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "dur",
			Long:  "duration",
		},
	}
	ff.SetDesc("duration flag")
	ff.SetDefaultValue(time.Second * 5)
	ff.SetPlaceHolder("DURATION")
	err = c1.AddFlag(ff)
	if err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "flt",
			Long:  "float",
		},
	}
	ff.SetDesc("float flag")
	ff.SetDefaultValue(float32(357))
	ff.SetPlaceHolder("NUM")
	err = c1.AddFlag(ff)
	if err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "K",
			Long:  "insecure",
		},
	}
	ff.SetDesc("Skip TLS host verification")
	ff.SetDefaultValue(true)
	ff.SetPlaceHolder("")
	err = c1.AddFlag(ff)
	if err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "",
			Long:  "prefix",
		},
	}
	ff.SetDesc("Root key prefix")
	ff.SetDefaultValue("/")
	ff.SetPlaceHolder("ROOT")
	err = c1.AddFlag(ff)
	if err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "",
			Long:  "cacert",
		},
	}
	ff.SetDesc("Client CA cert")
	ff.SetDefaultValue("")
	ff.SetPlaceHolder("FILE")
	err = c1.AddFlag(ff)
	if err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "",
			Long:  "cert",
		},
	}
	ff.SetDesc("Client cert")
	ff.SetDefaultValue("")
	ff.SetPlaceHolder("FILE")
	err = c1.AddFlag(ff)
	if err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "",
			Long:  "scheme",
		},
	}
	ff.SetDesc("Consul connection scheme (HTTP or HTTPS)")
	ff.SetDefaultValue("")
	ff.SetPlaceHolder("SCHEME")
	err = c1.AddFlag(ff)
	if err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "u",
			Long:  "username",
		},
	}
	ff.SetDesc("HTTP Basic auth user")
	ff.SetDefaultValue("")
	ff.SetPlaceHolder("USERNAME")
	err = c1.AddFlag(ff)
	if err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short:   "pw",
			Long:    "password",
			Aliases: []string{"passwd", "pwd"},
		},
	}
	ff.SetDesc("HTTP Basic auth password")
	ff.SetDefaultValue("")
	ff.SetPlaceHolder("PASSWORD")
	err = c1.AddFlag(ff)
	if err != nil {
		return
	}

	return
}

func serverCommandsGet() (cmd *cli.CmdS) { //nolint:funlen,revive //for test
	cmd = &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			// name:        "server",
			Short:   "s",
			Long:    "server",
			Aliases: []string{"serve", "svr"},
		},
	}
	cmd.SetDesc("server ops: for linux service/daemon.")
	cmd.SetDeprecated("1.0")
	cmd.SetExamples(`random examples`)

	var err error

	ff := &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "h",
			Long:  "head",
		},
	}
	ff.SetDesc("head -1 like")
	ff.SetDefaultValue(0)
	ff.SetHeadLike(true)
	err = cmd.AddFlag(ff)
	if err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "l",
			Long:  "tail",
		},
	}
	ff.SetDesc("tail -1 like [not support]")
	ff.SetDefaultValue(0)
	ff.SetPlaceHolder("")
	err = cmd.AddFlag(ff)
	if err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "e",
			Long:  "enum",
		},
	}
	ff.SetDesc("enum tests")
	ff.SetDefaultValue("apple")
	ff.SetValidArgs("apple", "banana", "orange")
	ff.SetPlaceHolder("")
	err = cmd.AddFlag(ff)
	if err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "tt",
			Long:  "retry",
		},
	}
	ff.SetDesc("ss")
	ff.SetDefaultValue(1)
	ff.SetPlaceHolder("RETRY")
	err = cmd.AddFlag(ff)
	if err != nil {
		return
	}

	cc := &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Short:   "s",
			Long:    "start",
			Aliases: []string{"run", "startup"},
		},
	}
	cc.SetDesc("startup this system service/daemon.")
	err = cmd.AddSubCommand(cc)
	if err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Long:    "foreground",
			Short:   "f",
			Aliases: []string{"fg"},
		},
	}
	ff.SetDesc("run at foreground")
	ff.SetDefaultValue(false)
	ff.SetPlaceHolder("FOREGROUND")
	err = cc.AddFlag(ff)
	if err != nil {
		return
	}

	cc = &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Short:   "nf", // parent no Full
			Aliases: []string{"run1", "startup1"},
		},
	}
	cc.SetDesc("dup test: startup this system service/daemon.")
	err = cmd.AddSubCommand(cc)
	if err != nil {
		return
	}

	c1 := &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Short:   "nf", // parent no Full
			Aliases: []string{"run1", "startup1"},
		},
	}
	c1.SetDesc("dup test: startup this system service/daemon.")
	err = cc.AddSubCommand(c1)
	if err != nil {
		return
	}

	c2 := &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Short:   "nf", // parent no Full
			Aliases: []string{"run1", "startup1"},
		},
	}
	c2.SetDesc("dup test: startup this system service/daemon.")
	err = c1.AddSubCommand(c2)
	if err != nil {
		return
	}

	cc = &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Short:   "t",
			Long:    "stop",
			Aliases: []string{"stp", "halt", "pause"},
		},
	}
	cc.SetDesc("stop this system service/daemon.")
	err = cmd.AddSubCommand(cc)
	if err != nil {
		return
	}

	cc = &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Short:   "r",
			Long:    "restart",
			Aliases: []string{"reload"},
		},
	}
	cc.SetDesc("restart this system service/daemon.")
	err = cmd.AddSubCommand(cc)
	if err != nil {
		return
	}

	cc = &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Long:    "status",
			Aliases: []string{"st"},
		},
	}
	cc.SetDesc("display its running status as a system service/daemon.")
	err = cmd.AddSubCommand(cc)
	if err != nil {
		return
	}

	cc = &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Short:   "i",
			Long:    "install",
			Aliases: []string{"setup"},
		},
	}
	cc.SetDesc("install as a system service/daemon.")
	err = cmd.AddSubCommand(cc)
	if err != nil {
		return
	}

	cc = &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Short:   "u",
			Long:    "uninstall",
			Aliases: []string{"remove"},
		},
	}
	cc.SetDesc("remove from a system service/daemon.")
	err = cmd.AddSubCommand(cc)
	if err != nil {
		return
	}

	return
}

func kvCommandsGet() (cmd *cli.CmdS) { //nolint:funlen,revive //for test
	cmd = &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Long:    "kv",
			Aliases: []string{"kvstore"},
		},
	}
	cmd.SetName("kvstore")
	cmd.SetDesc("consul kv store operations...")

	// flags: consulConnectFlagsGet(), // *Clone(&consulConnectFlags, &[]*Flag{}).(*[]*Flag),
	for _, f := range consulConnectFlagsGet() {
		if err := cmd.AddFlag(f); err != nil {
			return
		}
	}

	cc := &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Short:   "b",
			Long:    "backup",
			Aliases: []string{"bk", "bf", "bkp"},
		},
	}
	cc.SetDesc("Dump Consul's KV database to a JSON/YAML file")
	cc.SetGroup("bbb")
	cc.SetAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
		return
	})
	cc.SetPreActions(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
		return
	})
	cc.SetPostActions(func(ctx context.Context, cmd cli.Cmd, args []string, errInvoked error) (err error) {
		return
	})
	err := cmd.AddSubCommand(cc)
	if err != nil {
		return
	}

	ff := &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "o",
			Long:  "output",
		},
	}
	ff.SetDesc("Write output to a file (*.json / *.yml)")
	ff.SetDeprecated("2.0")
	ff.SetDefaultValue("consul-backup.json")
	ff.SetPlaceHolder("FILE")
	err = cc.AddFlag(ff)
	if err != nil {
		return
	}

	cc = &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Short: "r",
			Long:  "restore",
		},
	}
	cc.SetDesc("restore to Consul's KV store, from a a JSON/YAML backup file")
	cc.SetGroup("bbb")
	cc.SetAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
		return
	})
	err = cmd.AddSubCommand(cc)
	if err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "in",
			Long:  "input",
		},
	}
	ff.SetDesc("Read input from a file (*.json / *.yml)")
	ff.SetDeprecated("2.0")
	ff.SetDefaultValue("consul-backup.json")
	ff.SetPlaceHolder("FILE")
	err = cc.AddFlag(ff)
	if err != nil {
		return
	}

	cc = &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Short: "hh",
			Long:  "hidden-cmd",
		},
	}
	cc.SetDesc("restore to Consul's KV store, from a a JSON/YAML backup file")
	cc.SetGroup("bbb")
	cc.SetHidden(true, false)
	err = cmd.AddSubCommand(cc)
	if err != nil {
		return
	}

	cc = &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Short: "ls",
			Long:  "list",
		},
	}
	cc.SetDesc("list to Consul's KV store, from a a JSON/YAML backup file")
	cc.SetGroup("bbb")
	err = cmd.AddSubCommand(cc)
	if err != nil {
		return
	}

	return
}

func tagsCommandsGet() (cmd *cli.CmdS) { //nolint:funlen,revive //for test
	cmd = &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			// Short:       "t",
			Long:    "tags",
			Aliases: []string{},
		},
	}
	cmd.SetDesc("tags operations...")
	for _, f := range consulConnectFlagsGet() {
		if err := cmd.AddFlag(f); err != nil {
			return
		}
	}

	cc := &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Short:   "ls",
			Long:    "list",
			Aliases: []string{"l", "lst", "dir"},
		},
	}
	cc.SetDesc("list tags ...")
	cc.SetGroup("2323.List")
	if err := cmd.AddSubCommand(cc); err != nil {
		return
	}

	cc = &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Short:   "a",
			Long:    "add",
			Aliases: []string{"create", "new"},
		},
	}
	cc.SetDesc("add tags ...")
	cc.SetGroup("311Z.Add")
	if err := cmd.AddSubCommand(cc); err != nil {
		return
	}

	ff := &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short:   "ls",
			Long:    "list",
			Aliases: []string{"l", "lst", "dir"},
		},
	}
	ff.SetDesc("a comma list to be added")
	ff.SetPlaceHolder("LIST")
	ff.SetDefaultValue([]string{})
	if err := cc.AddFlag(ff); err != nil {
		return
	}

	cc = &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Short:   "r",
			Long:    "rm",
			Aliases: []string{"remove", "erase", "delete", "del"},
		},
	}
	cc.SetDesc("remove tags ...")
	cc.SetGroup("311Z.Add")
	if err := cmd.AddSubCommand(cc); err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short:   "ls",
			Long:    "list",
			Aliases: []string{"l", "lst", "dir"},
		},
	}
	ff.SetDesc("a comma list to be added")
	ff.SetPlaceHolder("LIST")
	ff.SetDefaultValue([]string{})
	if err := cc.AddFlag(ff); err != nil {
		return
	}

	cc = &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Short:   "m",
			Long:    "modify",
			Aliases: []string{"mod", "update", "change"},
		},
	}
	cc.SetDesc("modify tags ...")
	cc.SetGroup("313Z.Modify")
	if err := cmd.AddSubCommand(cc); err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "a",
			Long:  "add",
		},
	}
	ff.SetDesc("a comma list to be added")
	ff.SetPlaceHolder("LIST")
	ff.SetDefaultValue([]string{})
	if err := cc.AddFlag(ff); err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short:   "r",
			Long:    "rm",
			Aliases: []string{"remove", "erase", "del"},
		},
	}
	ff.SetDesc("a comma list to be removed.")
	ff.SetPlaceHolder("LIST")
	ff.SetDefaultValue([]string{})
	if err := cc.AddFlag(ff); err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "u",
			Long:  "ued",
		},
	}
	ff.SetDesc("a comma list to be removed.")
	ff.SetPlaceHolder("LIST")
	ff.SetDefaultValue("7,99")
	if err := cc.AddFlag(ff); err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "w",
			Long:  "wed",
		},
	}
	ff.SetDesc("a comma list to be removed.")
	ff.SetPlaceHolder("LIST")
	ff.SetDefaultValue([]string{"2", "3"})
	if err := cc.AddFlag(ff); err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "z",
			Long:  "zed",
		},
	}
	ff.SetDesc("a comma list to be removed.")
	ff.SetPlaceHolder("LIST")
	ff.SetDefaultValue([]uint{2, 3})
	if err := cc.AddFlag(ff); err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "x",
			Long:  "xed",
		},
	}
	ff.SetDesc("a comma list to be removed.")
	ff.SetPlaceHolder("LIST")
	ff.SetDefaultValue([]int{4, 5})
	if err := cc.AddFlag(ff); err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "v",
			Long:  "ved",
		},
	}
	ff.SetDesc("a comma list to be removed.")
	ff.SetPlaceHolder("LIST")
	ff.SetDefaultValue(2 * time.Second)
	if err := cc.AddFlag(ff); err != nil {
		return
	}

	cc = &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Short:   "t",
			Long:    "toggle",
			Aliases: []string{"tog", "switch"},
		},
	}
	cc.SetDesc("toggle tags for ms.")
	cc.SetGroup("313Z.Modify")
	if err := cmd.AddSubCommand(cc); err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short:   "x",
			Long:    "address",
			Aliases: []string{"addr"},
		},
	}
	ff.SetDesc("the address of the service (by id or name)")
	ff.SetPlaceHolder("HOST:PORT")
	ff.SetDefaultValue("")
	if err := cc.AddFlag(ff); err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "s",
			Long:  "set",
		},
	}
	ff.SetDesc("set to `tag` which service specified by --address")
	ff.SetPlaceHolder("LIST")
	ff.SetDefaultValue([]string{})
	if err := cc.AddFlag(ff); err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short:   "u",
			Long:    "unset",
			Aliases: []string{"reset"},
		},
	}
	ff.SetDesc("and reset the others service nodes to `tag`")
	ff.SetPlaceHolder("LIST")
	ff.SetDefaultValue([]string{})
	if err := cc.AddFlag(ff); err != nil {
		return
	}

	return
}

func msCommandsGet() (cmd *cli.CmdS) { //nolint:funlen,revive //for test
	cmd = &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Long:    "microservices",
			Aliases: []string{"microservice", "micro-service"},
		},
	}
	cmd.SetName("ms")
	cmd.SetDesc("micro-service operations...")

	ff := &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "n",
			Long:  "name",
		},
	}
	ff.SetDesc("name of the service")
	ff.SetDescription("name of the service", `fdhsjsfhdsjk`)
	ff.SetExamples(`fdhsjsfhdsjk`)
	ff.SetDefaultValue("")
	ff.SetPlaceHolder("NAME")
	if err := cmd.AddFlag(ff); err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "i",
			Long:  "id",
		},
	}
	ff.SetDesc("unique id of the service")
	ff.SetDefaultValue("")
	ff.SetPlaceHolder("ID")
	if err := cmd.AddFlag(ff); err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "a",
			Long:  "all",
		},
	}
	ff.SetDesc("all services")
	ff.SetDefaultValue(false)
	if err := cmd.AddFlag(ff); err != nil {
		return
	}

	ff = &cli.Flag{
		BaseOpt: cli.BaseOpt{
			Short: "cc",
			Long:  "",
		},
	}
	ff.SetDesc("unique id of the service")
	ff.SetDefaultValue("")
	if err := cmd.AddFlag(ff); err != nil {
		return
	}

	if err := cmd.AddSubCommand(tagsCommandsGet()); err != nil {
		return
	}

	cc := &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Short:   "l",
			Long:    "list",
			Aliases: []string{"ls", "lst"},
		},
	}
	cc.SetDesc("list services")
	cc.SetGroup("")
	if err := cmd.AddSubCommand(cc); err != nil {
		return
	}

	cc = &cli.CmdS{
		BaseOpt: cli.BaseOpt{},
	}
	cc.SetDesc("an empty subcommand for testing - list services.")
	cc.SetGroup("56.vvvvvv")
	if err := cmd.AddSubCommand(cc); err != nil {
		return
	}

	cc = &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Short: "dr",
		},
	}
	cc.SetDesc("list services.")
	cc.SetGroup("56.vvvvvv")
	if err := cmd.AddSubCommand(cc); err != nil {
		return
	}

	cc = &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Short: "dz",
		},
	}
	cc.SetDesc("list services.")
	cc.SetGroup("56.vvvvvv")
	if err := cmd.AddSubCommand(cc); err != nil {
		return
	}

	c1 := &cli.CmdS{
		BaseOpt: cli.BaseOpt{
			Short: "dz",
		},
	}
	c1.SetDesc("list services.")
	c1.SetGroup("56.vvvvvv")
	if err := cc.AddSubCommand(c1); err != nil {
		return
	}

	return
}

//

//

//

func newTestRunner() cli.Runner {
	return &workerS{store: store.New()}
}

// workerS for testing only
type workerS struct {
	store   store.Store
	retCode int
}

func (w *workerS) SetSuggestRetCode(ret int) {
	w.retCode = ret
}

func (w *workerS) InitGlobally(ctx context.Context) {}
func (w *workerS) Ready() bool                      { return true }
func (w *workerS) DumpErrors(wr io.Writer)          {}             //nolint:revive
func (w *workerS) Error() errorsv3.Error            { return nil } //nolint:revive
func (w *workerS) Recycle(errs ...error)            {}             //
func (w *workerS) Store(prefix ...string) store.Store {
	if len(prefix) > 0 {
		return w.store.WithPrefix(prefix...)
	}
	return w.store
}
func (w *workerS) Run(ctx context.Context, opts ...cli.Opt) (err error) { return }            //nolint:revive
func (w *workerS) Actions() (ret map[string]bool)                       { return }            //nolint:revive
func (w *workerS) Name() string                                         { return "for-test" } //
func (*workerS) Version() string                                        { return "v0.0.0" }
func (*workerS) Root() *cli.RootCommand                                 { return nil }
func (*workerS) Args() []string                                         { return nil }       //
func (w *workerS) SuggestRetCode() int                                  { return w.retCode } //
func (w *workerS) ParsedState() cli.ParsedState                         { return nil }
func (w *workerS) LoadedSources() (results []cli.LoadedSources)         { return }

func (w *workerS) SetCancelFunc(cancelFunc func()) {}
func (w *workerS) CancelFunc() func()              { return nil }

func (w *workerS) DoBuiltinAction(ctx context.Context, action cli.ActionEnum, args ...any) (handled bool, err error) {
	return
}

//

//

//

// appS for testing only
type appS struct {
	cli.Runner
	root  *cli.RootCommand
	args  []string
	inCmd bool
	inFlg bool
}

func (s *appS) GetRunner() cli.Runner { return s.Runner }

func (s *appS) NewCommandBuilder(longTitle string, titles ...string) cli.CommandBuilder {
	return s.Cmd(longTitle, titles...)
}

func (s *appS) NewFlagBuilder(longTitle string, titles ...string) cli.FlagBuilder {
	return s.Flg(longTitle, titles...)
}

func (s *appS) Cmd(longTitle string, titles ...string) cli.CommandBuilder { //nolint:revive
	s.inCmd = true
	// return newCommandBuilder(s, longTitle, titles...)
	return nil
}

func (s *appS) With(cb func(app cli.App)) { //nolint:revive
	cb(s)
}

func (s *appS) WithOpts(opts ...cli.Opt) cli.App {
	return s
}

func (s *appS) Flg(longTitle string, titles ...string) cli.FlagBuilder { //nolint:revive
	s.inFlg = true
	// return newFlagBuilder(s, longTitle, titles...)
	return nil
}

func (s *appS) AddCmd(f func(b cli.CommandBuilder)) cli.App { //nolint:revive
	// b := newCommandBuilder(s, "")
	// defer b.Build()
	// cb(b)
	return s
}

func (s *appS) AddFlg(cb func(b cli.FlagBuilder)) cli.App { //nolint:revive
	// b := newFlagBuilder(s, "")
	// defer b.Build()
	// cb(b)
	return s
}

func (s *appS) NewCmdFrom(from *cli.CmdS, cb func(b cli.CommandBuilder)) cli.App { //nolint:revive
	// b := newCommandBuilderFrom(from, s, "")
	// defer b.Build()
	// cb(b)
	return s
}

func (s *appS) NewFlgFrom(from *cli.CmdS, defaultValue any, cb func(b cli.FlagBuilder)) cli.App { //nolint:revive
	// b := newFlagBuilderFrom(from, s, defaultValue, "")
	// defer b.Build()
	// cb(b)
	return s
}

func (s *appS) ToggleableFlags(toggleGroupName string, items ...cli.BatchToggleFlag) {}

func (s *appS) RootBuilder(cb func(b cli.CommandBuilder)) cli.App { return s }

func (s *appS) addCommand(child *cli.CmdS) { //nolint:unused
	s.inCmd = false
	if cx, ok := s.root.Cmd.(*cli.CmdS); ok {
		cx.AddSubCommand(child)
	}
}

func (s *appS) addFlag(child *cli.Flag) { //nolint:unused
	s.inFlg = false
	if cx, ok := s.root.Cmd.(*cli.CmdS); ok {
		cx.AddFlag(child)
	}
}

func (s *appS) Info(name, version string, desc ...string) cli.App {
	s.ensureNewApp()
	if s.root.AppName == "" {
		s.root.AppName = name
	}
	if s.root.Version == "" {
		s.root.Version = version
	}
	if cx, ok := s.root.Cmd.(*cli.CmdS); ok {
		cx.SetDescription("", desc...)
	}
	return s
}

func (s *appS) Examples(examples ...string) cli.App {
	s.ensureNewApp()
	if cx, ok := s.root.Cmd.(*cli.CmdS); ok {
		cx.SetExamples(examples...)
	}
	return s
}

func (s *appS) Copyright(copyright string) cli.App {
	s.ensureNewApp()
	s.root.Copyright = copyright
	return s
}

func (s *appS) Author(author string) cli.App {
	s.ensureNewApp()
	s.root.Author = author
	return s
}

func (s *appS) Description(desc string) cli.App {
	s.ensureNewApp()
	s.root.SetDesc(desc)
	return s
}

func (s *appS) Header(headerLine string) cli.App {
	s.ensureNewApp()
	s.root.HeaderLine = headerLine
	return s
}

func (s *appS) Footer(footerLine string) cli.App {
	s.ensureNewApp()
	s.root.FooterLine = footerLine
	return s
}

func (s *appS) OnAction(handler cli.OnInvokeHandler) cli.App {
	return s
}

func (s *appS) SetRootCommand(root *cli.RootCommand) cli.App {
	s.root = root
	return s
}

func (s *appS) WithRootCommand(cb func(root *cli.RootCommand)) cli.App {
	cb(s.root)
	return s
}

func (s *appS) RootCommand() *cli.RootCommand { return s.root }

func (s *appS) Name() string                    { return s.root.AppName }
func (s *appS) Version() string                 { return s.root.Version }
func (s *appS) Worker() cli.Runner              { return s.Runner }
func (s *appS) Root() *cli.RootCommand          { return s.root }
func (s *appS) Args() []string                  { return s.args }
func (w *appS) SetCancelFunc(cancelFunc func()) {}

func (s *appS) ensureNewApp() cli.App { //nolint:unparam
	if s.root == nil {
		s.root = &cli.RootCommand{
			AppName: conf.AppName,
			Version: conf.Version,
			// Copyright:  "",
			// Author:     "",
			// HeaderLine: "",
			// FooterLine: "",
			// CmdS:    nil,
		}
		s.root.SetApp(s)
	}
	if s.root.Cmd == nil {
		s.root.Cmd = new(cli.CmdS)
		s.root.Cmd.SetName(s.root.AppName)
	}
	return s
}

func (s *appS) Build() {
	type setRoot interface {
		SetRoot(root *cli.RootCommand, args []string)
	}
	if sr, ok := s.Runner.(setRoot); ok {
		ctx := context.Background()
		if cx, ok := s.root.Cmd.(*cli.CmdS); ok {
			cx.EnsureTree(ctx, s, s.root)
		}
		sr.SetRoot(s.root, s.args)
	}
}

func (s *appS) Run(ctx context.Context, opts ...cli.Opt) (err error) {
	if s.inCmd {
		return errors.New("a NewCommandBuilder()/Cmd() call needs ending with Build()")
	}
	if s.inFlg {
		return errors.New("a NewFlagBuilder()/Flg() call needs ending with Build()")
	}

	if s.root == nil || s.root.Cmd == nil {
		return errors.New("the RootCommand hasn't been built")
	}

	s.Build() // set rootCommand into worker

	s.Runner.InitGlobally(ctx) // let worker starts initializations

	if !s.Runner.Ready() {
		return errors.New("the RootCommand hasn't been built, or Init() failed. Has builder.App.Build() called? ")
	}

	err = s.Runner.Run(ctx, opts...)

	// if err != nil {
	// 	s.Runner.DumpErrors(os.Stderr)
	// }

	return
}
