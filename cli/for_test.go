package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	errorsv3 "gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/store"

	"github.com/hedzr/cmdr/v2/conf"
)

func rootCmdForTesting() (root *RootCommand) { //nolint:funlen,revive //for test
	app := &appS{
		Runner: newTestRunner(),
	}

	serverCommands := serverCommandsGet()

	kvCommands := kvCommandsGet()

	msCommands := msCommandsGet()

	root = &RootCommand{
		AppName: "consul-tags",
		Version: "0.0.1",
		app:     app,
		// Header:  `dsjlfsdjflsdfjlsdjflksjdfdsfsd`,
		// Version:    consul_tags.Version,
		// VersionInt: consul_tags.VersionInt,
		Copyright: "consul-tags is an effective devops tool",
		Author:    "Hedzr Yeh <hedzr@duck.com>",

		Cmd: &CmdS{
			BaseOpt: BaseOpt{
				name: "consul-tags",
			},
			flags: []*Flag{
				// global options here.
				{
					BaseOpt: BaseOpt{
						Short:       "t",
						Long:        "retry",
						description: "ss",
						examples:    `random examples`,
						deprecated:  "1.2.3",
					},
					defaultValue: 1,
					placeHolder:  "RETRY",
				},
				{
					BaseOpt: BaseOpt{
						Short:       "s",
						Long:        "s",
						description: "s",
					},
					defaultValue: uint(1),
					onMatched: func(flg *Flag, position int, hitState *MatchState) (err error) {
						if flg.GetDescZsh() != "ss" {
							err = errors.New("err `t`.GetDescZsh()")
						}
						if flg.GetTitleNamesBy(",") == "" {
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
					},
				},
				{
					BaseOpt: BaseOpt{
						Short:       "ff",
						Long:        "float",
						description: "",
					},
					defaultValue: float64(0),
				},
				// "consul-tags -cc 3.14159-2.56i": func(t *testing.T) error {
				// 	if GetComplex128("app.complex") != 3.14159-2.56i {
				// 		return errors.New("something wrong complex. |expected %v|got %v|", 3.14159-2.56i, GetComplex128("app.complex"))
				// 	}
				// 	fmt.Println("consul-tags kv b ~ -------- no errors")
				// 	return nil
				// },
				// {
				// 	BaseOpt: BaseOpt{
				// 		Short:       "cc",
				// 		Long:        "complex",
				// 		description: "",
				// 	},
				// 	defaultValue: complex128(0),
				// },
				{
					BaseOpt: BaseOpt{
						Short:       "pp",
						Long:        "spasswd",
						description: "",
					},
					defaultValue:   "",
					externalEditor: ExternalToolPasswordInput,
					onMatched: func(f *Flag, position int, hitState *MatchState) (err error) {
						_, _ = fmt.Println("**** -pp action running")

						// f.owner.Runner.showVersions()
						// PrintBuildInfo()
						// cmd.PrintBuildInfo()
						// cmd.GetTitleZshNames()

						// SetCustomShowVersion(nil)
						// SetCustomShowBuildInfo(nil)
						_, _ = fmt.Println("**** -pp action end")
						return
					},
				},
				{
					BaseOpt: BaseOpt{
						Short:       "qq",
						Long:        "qqpasswd",
						description: "",
					},
					defaultValue:   "567",
					externalEditor: ExternalToolEditor,
				},
				{
					BaseOpt: BaseOpt{
						Short:       "dd",
						Long:        "ddduration",
						description: "",
					},
					defaultValue: time.Second,
				},
			},
			preActions: []OnPreInvokeHandler{func(ctx context.Context, cmd Cmd, args []string) (err error) {
				return
			}},
			postActions: []OnPostInvokeHandler{func(ctx context.Context, cmd Cmd, args []string, errInvoked error) (err error) {
				return
			}},
			commands: []*CmdS{
				// dnsCommands,
				// playCommand,
				// generatorCommands,

				serverCommands,
				msCommands,
				kvCommands,

				{
					BaseOpt: BaseOpt{
						Short:       "ls",
						Long:        "list",
						description: "list to Consul's KV store, from a a JSON/YAML backup file",
					},
					flags: []*Flag{
						// global options here.
						{
							BaseOpt: BaseOpt{
								Short:       "t",
								Long:        "retry",
								description: "ss",
								examples:    `random examples`,
								deprecated:  "1.2.3",
							},
							defaultValue: 1,
							placeHolder:  "RETRY",
						},
					},
				},
			},
		},
	}

	ctx := context.TODO()
	app.root = root
	root.Cmd.(*CmdS).EnsureTree(ctx, app, root)
	root.Cmd.(*CmdS).EnsureXref(ctx)
	return
}

//

//

//

func consulConnectFlagsGet() []*Flag { //nolint:funlen,revive //for test
	consulConnectFlags := []*Flag{
		{
			BaseOpt: BaseOpt{
				Short:       "a",
				Long:        "addr",
				description: "Consul ip/host and port: HOST[:PORT] (No leading 'http(s)://')",
			},
			defaultValue: "consul.ops.local",
			placeHolder:  "HOST[:PORT]",
		},
		{
			BaseOpt: BaseOpt{
				Short:       "p",
				Long:        "port",
				description: "Consul port",
			},
			defaultValue: 8500,
			placeHolder:  "PORT",
		},
		{
			BaseOpt: BaseOpt{
				Short:       "ui",
				Long:        "uint",
				description: "uint flag",
			},
			defaultValue: uint(357),
			placeHolder:  "NUM",
		},
		{
			BaseOpt: BaseOpt{
				Short:       "dur",
				Long:        "duration",
				description: "duration flag",
			},
			defaultValue: time.Second,
			placeHolder:  "DURATION",
		},
		{
			BaseOpt: BaseOpt{
				Short:       "flt",
				Long:        "float",
				description: "float flag",
			},
			defaultValue: float32(357),
			placeHolder:  "NUM",
		},
		{
			BaseOpt: BaseOpt{
				Short:       "K",
				Long:        "insecure",
				description: "Skip TLS host verification",
			},
			defaultValue: true,
			placeHolder:  "PORT",
		},
		{
			BaseOpt: BaseOpt{
				Short:       "",
				Long:        "prefix",
				description: "Root key prefix",
			},
			defaultValue: "/",
			placeHolder:  "ROOT",
		},
		{
			BaseOpt: BaseOpt{
				Short:       "",
				Long:        "cacert",
				description: "Client CA cert",
			},
			defaultValue: "",
			placeHolder:  "FILE",
		},
		{
			BaseOpt: BaseOpt{
				Short:       "",
				Long:        "cert",
				description: "Client cert",
			},
			defaultValue: "",
			placeHolder:  "FILE",
		},
		{
			BaseOpt: BaseOpt{
				Short:       "",
				Long:        "scheme",
				description: "Consul connection scheme (HTTP or HTTPS)",
			},
			defaultValue: "",
			placeHolder:  "SCHEME",
		},
		{
			BaseOpt: BaseOpt{
				Short:       "u",
				Long:        "username",
				description: "HTTP Basic auth user",
			},
			defaultValue: "",
			placeHolder:  "USERNAME",
		},
		{
			BaseOpt: BaseOpt{
				Short:       "pw",
				Long:        "password",
				Aliases:     []string{"passwd", "pwd"},
				description: "HTTP Basic auth password",
			},
			defaultValue: "",
			placeHolder:  "PASSWORD",
		},
	}
	return consulConnectFlags
}

func serverCommandsGet() *CmdS { //nolint:funlen,revive //for test
	serverCommands := &CmdS{
		BaseOpt: BaseOpt{
			// name:        "server",
			Short:       "s",
			Long:        "server",
			Aliases:     []string{"serve", "svr"},
			description: "server ops: for linux service/daemon.",
			deprecated:  "1.0",
			examples:    `random examples`,
		},
		flags: []*Flag{
			{
				BaseOpt: BaseOpt{
					Short:       "h",
					Long:        "head",
					description: "head -1 like",
				},
				defaultValue: 0,
				headLike:     true,
			},
			{
				BaseOpt: BaseOpt{
					Short:       "l",
					Long:        "tail",
					description: "tail -1 like [not support]",
				},
				defaultValue: 0,
				headLike:     false,
			},
			{
				BaseOpt: BaseOpt{
					Short:       "e",
					Long:        "enum",
					description: "enum tests",
				},
				defaultValue: "apple",
				validArgs:    []string{"apple", "banana", "orange"},
			},
			{
				BaseOpt: BaseOpt{
					Short:       "tt",
					Long:        "retry",
					description: "ss",
				},
				defaultValue: 1,
				placeHolder:  "RETRY",
			},
		},
		commands: []*CmdS{
			{
				BaseOpt: BaseOpt{
					Short:       "s",
					Long:        "start",
					Aliases:     []string{"run", "startup"},
					description: "startup this system service/daemon.",
					// Action:impl.ServerStart,
				},
				// preActions: []OnPreInvokeHandler{impl.ServerStartPre},
				// PostAction: impl.ServerStartPost,
				flags: []*Flag{
					// {
					// 	BaseOpt: BaseOpt{
					// 		Short:       "t",
					// 		Long:        "retry",
					// 		description: "ss",
					// 	},
					// 	defaultValue:            1,
					// 	placeHolder: "RETRY",
					// },
					// {
					// 	BaseOpt: BaseOpt{
					// 		Short:       "t",
					// 		Long:        "retry",
					// 		description: "ss: dup test",
					// 	},
					// 	defaultValue:            1,
					// 	placeHolder: "RETRY",
					// },
					{
						BaseOpt: BaseOpt{
							Long:        "foreground",
							Short:       "f",
							description: "run foreground",
						},
						defaultValue: false,
					},
				},
			},
			// {
			// 	BaseOpt: BaseOpt{
			// 		Short:       "s",
			// 		Long:        "start",
			// 		Aliases:     []string{"run", "startup"},
			// 		description: "dup test: startup this system service/daemon.",
			// 		// Action:impl.ServerStart,
			// 	},
			// 	// preActions: []OnPreInvokeHandler{impl.ServerStartPre},
			// 	// PostAction: impl.ServerStartPost,
			// },
			{
				BaseOpt: BaseOpt{
					Short:       "nf", // parent no Full
					Aliases:     []string{"run1", "startup1"},
					description: "dup test: startup this system service/daemon.",
				},
				preActions: []OnPreInvokeHandler{func(ctx context.Context, cmd Cmd, args []string) (err error) {
					_, _ = fmt.Println(cmd.OwnerCmd().Name(), cmd.Name(), cmd.GroupTitle())
					_, _ = fmt.Println(cmd.Root().Name(), cmd.Root().OwnerCmd())
					_, _ = fmt.Println(cmd.App().Name())
					return
				}},
				commands: []*CmdS{
					{
						BaseOpt: BaseOpt{
							Short:       "nf", // parent no Full
							Aliases:     []string{"run", "startup"},
							description: "dup test: startup this system service/daemon.",
						},
						preActions: []OnPreInvokeHandler{func(ctx context.Context, cmd Cmd, args []string) (err error) {
							_, _ = fmt.Println(cmd.OwnerCmd().Name(), cmd.Name(), cmd.GroupTitle())
							return
						}},
					},
				},
			},
			{
				BaseOpt: BaseOpt{
					Short:       "t",
					Long:        "stop",
					Aliases:     []string{"stp", "halt", "pause"},
					description: "stop this system service/daemon.",
					// Action:impl.ServerStop,
				},
			},
			{
				BaseOpt: BaseOpt{
					Short:       "r",
					Long:        "restart",
					Aliases:     []string{"reload"},
					description: "restart this system service/daemon.",
					// Action:impl.ServerRestart,
				},
			},
			{
				BaseOpt: BaseOpt{
					Long:        "status",
					Aliases:     []string{"st"},
					description: "display its running status as a system service/daemon.",
				},
			},
			{
				BaseOpt: BaseOpt{
					Short:       "i",
					Long:        "install",
					Aliases:     []string{"setup"},
					description: "install as a system service/daemon.",
				},
			},
			{
				BaseOpt: BaseOpt{
					Short:       "u",
					Long:        "uninstall",
					Aliases:     []string{"remove"},
					description: "remove from a system service/daemon.",
				},
			},
		},
	}
	return serverCommands
}

func kvCommandsGet() *CmdS { //nolint:funlen,revive //for test
	kvCommands := &CmdS{
		BaseOpt: BaseOpt{
			name:        "kvstore",
			Long:        "kv",
			Aliases:     []string{"kvstore"},
			description: "consul kv store operations...",
		},
		flags: consulConnectFlagsGet(), // *Clone(&consulConnectFlags, &[]*Flag{}).(*[]*Flag),
		commands: []*CmdS{
			{
				BaseOpt: BaseOpt{
					Short:       "b",
					Long:        "backup",
					Aliases:     []string{"bk", "bf", "bkp"},
					description: "Dump Consul's KV database to a JSON/YAML file",
					group:       "bbb",
				},
				onInvoke: func(ctx context.Context, cmd Cmd, args []string) (err error) {
					// for gocov

					// cmd.PrintHelp(false)
					// cmd.PrintVersion()

					// if cmd.GetRoot() != copyRootCmd {
					// 	return errors.New(fmt.Sprintf("failed: root is wrong (cmd.GetRoot() != copyRootCmd):
					//      copyRootCmd = %p, cmd.GetRoot() = %p", copyRootCmd, cmd.GetRoot()))
					// }
					// if copyRootCmd.IsRoot() != true {
					// 	return errors.New("failed: root test is wrong")
					// }

					// if cmd.GetHitStr() != "b" {
					// 	return errors.New("failed: GetHitStr() is wrong")
					// }
					// if cmd.GetName() != "backup" {
					// 	return errors.New("failed: GetName() is wrong")
					// }
					// if cmd.GetExpandableNames() != "{backup,b}" {
					// 	return errors.New("failed: GetExpandableNames() is wrong")
					// }
					// if cmd.GetQuotedGroupName() != "[bbb]" {
					// 	return errors.New("failed: GetQuotedGroupName() is wrong")
					// }
					//
					// if cmd.GetParentName() != "kv" {
					// 	return errors.New("failed: GetParentName() is wrong")
					// }
					// if cmd.GetOwner().GetSubCommandNamesBy(",") != "b,backup,bk,bf,bkp,r,restore,ls,list" {
					// 	return errors.New("failed: GetSubCommandNamesBy() is wrong: '%s'",
					// 		cmd.GetOwner().GetSubCommandNamesBy(","))
					// }
					//
					// cmd.PrintHelp(true)
					return
				},
				flags: []*Flag{
					{
						BaseOpt: BaseOpt{
							Short:       "o",
							Long:        "output",
							description: "Write output to a file (*.json / *.yml)",
							deprecated:  "2.0",
						},
						defaultValue: "consul-backup.json",
						placeHolder:  "FILE",
					},
				},
				preActions: []OnPreInvokeHandler{func(ctx context.Context, cmd Cmd, args []string) (err error) {
					return
				}},
				postActions: []OnPostInvokeHandler{func(ctx context.Context, cmd Cmd, args []string, errInvoked error) (err error) {
					return
				}},
			},
			{
				BaseOpt: BaseOpt{
					Short:       "r",
					Long:        "restore",
					description: "restore to Consul's KV store, from a a JSON/YAML backup file",
					// Action:      kvRestore,
				},
				flags: []*Flag{
					{
						BaseOpt: BaseOpt{
							Short:       "i",
							Long:        "input",
							description: "read the input file (*.json / *.yml)",
						},
						defaultValue: "consul-backup.json",
						placeHolder:  "FILE",
					},
				},
			},
			{
				BaseOpt: BaseOpt{
					Short:       "hh",
					Long:        "hidden-cmd",
					description: "restore to Consul's KV store, from a a JSON/YAML backup file",
					hidden:      true,
				},
			},
			{
				BaseOpt: BaseOpt{
					Short:       "ls",
					Long:        "list",
					description: "list to Consul's KV store, from a a JSON/YAML backup file",
				},
			},
		},
	}
	return kvCommands
}

func tagsCommandsGet() *CmdS { //nolint:funlen,revive //for test
	tagsCommands := &CmdS{
		BaseOpt: BaseOpt{
			// Short:       "t",
			Long:        "tags",
			Aliases:     []string{},
			description: "tags op.",
		},
		flags: consulConnectFlagsGet(), // *Clone(&consulConnectFlags, &[]*Flag{}).(*[]*Flag),
		commands: []*CmdS{
			{
				BaseOpt: BaseOpt{
					Short:       "ls",
					Long:        "list",
					Aliases:     []string{"l", "lst", "dir"},
					description: "list tags.",
					// Action:      msTagsList,
					group: "2323.List",
				},
			},
			{
				BaseOpt: BaseOpt{
					Short:       "a",
					Long:        "add",
					Aliases:     []string{"create", "new"},
					description: "add tags.",
					// Action:      msTagsAdd,
					group: "311Z.Add",
				},
				flags: []*Flag{
					{
						BaseOpt: BaseOpt{
							Short:       "ls",
							Long:        "list",
							Aliases:     []string{"l", "lst", "dir"},
							description: "a comma list to be added",
						},
						defaultValue: []string{},
						placeHolder:  "LIST",
					},
				},
			},
			{
				BaseOpt: BaseOpt{
					Short:       "r",
					Long:        "rm",
					Aliases:     []string{"remove", "erase", "delete", "del"},
					description: "remove tags.",
					// Action:      msTagsRemove,
				},
				flags: []*Flag{
					{
						BaseOpt: BaseOpt{
							Short:       "ls",
							Long:        "list",
							Aliases:     []string{"l", "lst", "dir"},
							description: "a comma list to be added.",
						},
						defaultValue: []string{},
						placeHolder:  "LIST",
					},
				},
			},
			{
				BaseOpt: BaseOpt{
					Short:       "m",
					Long:        "modify",
					Aliases:     []string{"mod", "update", "change"},
					description: "modify tags.",
					// Action:      msTagsModify,
				},
				flags: []*Flag{
					{
						BaseOpt: BaseOpt{
							Short:       "a",
							Long:        "add",
							description: "a comma list to be added.",
						},
						defaultValue: []string{},
						placeHolder:  "LIST",
					},
					{
						BaseOpt: BaseOpt{
							Short:       "r",
							Long:        "rm",
							Aliases:     []string{"remove", "erase", "del"},
							description: "a comma list to be removed.",
						},
						defaultValue: []string{},
						placeHolder:  "LIST",
					},
					{
						BaseOpt: BaseOpt{
							Short:       "u",
							Long:        "ued",
							description: "a comma list to be removed.",
						},
						defaultValue: "7,99",
						placeHolder:  "LIST",
					},
					{
						BaseOpt: BaseOpt{
							Short:       "w",
							Long:        "wed",
							description: "a comma list to be removed.",
						},
						defaultValue: []string{"2", "3"},
						placeHolder:  "LIST",
					},
					{
						BaseOpt: BaseOpt{
							Short:       "z",
							Long:        "zed",
							description: "a comma list to be removed.",
						},
						defaultValue: []uint{2, 3},
						placeHolder:  "LIST",
					},
					{
						BaseOpt: BaseOpt{
							Short:       "x",
							Long:        "xed",
							description: "a comma list to be removed.",
						},
						defaultValue: []int{4, 5},
						placeHolder:  "LIST",
					},
					{
						BaseOpt: BaseOpt{
							Short:       "v",
							Long:        "ved",
							description: "a comma list to be removed.",
						},
						defaultValue: 2 * time.Second,
					},
				},
			},
			{
				BaseOpt: BaseOpt{
					Short:       "t",
					Long:        "toggle",
					Aliases:     []string{"tog", "switch"},
					description: "toggle tags for ms.",
					// Action:      msTagsToggle,
				},
				flags: []*Flag{
					{
						BaseOpt: BaseOpt{
							Short:       "x",
							Long:        "address",
							description: "the address of the service (by id or name)",
						},
						defaultValue: "",
						placeHolder:  "HOST:PORT",
					},
					{
						BaseOpt: BaseOpt{
							Short:       "s",
							Long:        "set",
							description: "set to `tag` which service specified by --address",
						},
						defaultValue: []string{},
						placeHolder:  "LIST",
					},
					{
						BaseOpt: BaseOpt{
							Short:       "u",
							Long:        "unset",
							Aliases:     []string{"reset"},
							description: "and reset the others service nodes to `tag`",
						},
						defaultValue: []string{},
						placeHolder:  "LIST",
					},
				},
			},
		},
	}
	return tagsCommands
}

func msCommandsGet() *CmdS { //nolint:funlen,revive //for test
	msCommands := &CmdS{
		BaseOpt: BaseOpt{
			name:        "microservices",
			Long:        "ms",
			Aliases:     []string{"microservice", "micro-service"},
			description: "micro-service operations...",
		},
		flags: []*Flag{
			{
				BaseOpt: BaseOpt{
					Short:       "n",
					Long:        "name",
					description: "name of the service",
					longDesc:    `fdhsjsfhdsjk`,
					examples:    `fsdhjkfsdhk`,
				},
				defaultValue: "",
				placeHolder:  "NAME",
			},
			{
				BaseOpt: BaseOpt{
					Short:       "i",
					Long:        "id",
					description: "unique id of the service",
				},
				defaultValue: "",
				placeHolder:  "ID",
			},
			{
				BaseOpt: BaseOpt{
					Short:       "a",
					Long:        "all",
					description: "all services",
				},
				defaultValue: false,
			},
			{
				BaseOpt: BaseOpt{
					Short:       "cc",
					Long:        "",
					description: "unique id of the service",
				},
				defaultValue: "",
				placeHolder:  "ID",
			},
		},
		commands: []*CmdS{
			tagsCommandsGet(),
			{
				BaseOpt: BaseOpt{
					Short:       "l",
					Long:        "list",
					Aliases:     []string{"ls", "lst"},
					description: "list services.",
					// Action:      msList,
					group: " ",
				},
				preActions: []OnPreInvokeHandler{func(ctx context.Context, cmd Cmd, args []string) (err error) {
					_, _ = fmt.Println(cmd.OwnerCmd().Name(), cmd.Name(), cmd.GroupTitle())
					return
				}},
			},
			{
				BaseOpt: BaseOpt{
					// Short: "",
					// Long:        "",
					// Aliases:     []string{"ls", "lst", "dir"},
					description: "an empty subcommand for testing - list services.",
					group:       "56.vvvvvv",
				},
				preActions: []OnPreInvokeHandler{func(ctx context.Context, cmd Cmd, args []string) (err error) {
					_, _ = fmt.Println(cmd.OwnerCmd().Name(), cmd.Name(), cmd.GroupTitle())
					return
				}},
			},
			{
				BaseOpt: BaseOpt{
					Short: "dr",
					// Long:        "list",
					// Aliases:     []string{"ls", "lst", "dir"},
					description: "list services.",
					group:       "56.vvvvvv",
				},
				preActions: []OnPreInvokeHandler{func(ctx context.Context, cmd Cmd, args []string) (err error) {
					_, _ = fmt.Println(cmd.OwnerCmd().Name(), cmd.Name(), cmd.GroupTitle())
					return
				}},
			},
			{
				BaseOpt: BaseOpt{
					name: "dz",
					// Long:        "list",
					// Aliases:     []string{"ls", "lst", "dir"},
					description: "list services.",
					group:       "56.vvvvvv",
				},
				preActions: []OnPreInvokeHandler{func(ctx context.Context, cmd Cmd, args []string) (err error) {
					_, _ = fmt.Println(cmd, "'s owner is", cmd.OwnerCmd())
					_, _ = fmt.Println(cmd.OwnerCmd().Name(), cmd.Name(), cmd.GroupTitle())
					return
				}},
				commands: []*CmdS{
					{
						BaseOpt: BaseOpt{
							name: "dz",
							// Long:        "list",
							// Aliases:     []string{"ls", "lst", "dir"},
							description: "list services.",
							group:       "56.vvvvvv",
						},
						preActions: []OnPreInvokeHandler{func(ctx context.Context, cmd Cmd, args []string) (err error) {
							_, _ = fmt.Println(cmd, "'s owner is", cmd.OwnerCmd())
							_, _ = fmt.Println(cmd.OwnerCmd().Name(), cmd.Name(), cmd.GroupTitle())
							return
						}},
					},
				},
			},
		},
	}
	return msCommands
}

//

//

//

func newTestRunner() Runner {
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

func (w *workerS) InitGlobally(ctx context.Context)                 {}
func (w *workerS) Ready() bool                                      { return true }
func (w *workerS) DumpErrors(wr io.Writer)                          {}                    //nolint:revive
func (w *workerS) Error() errorsv3.Error                            { return nil }        //nolint:revive
func (w *workerS) Recycle(errs ...error)                            {}                    //
func (w *workerS) Store() store.Store                               { return w.store }    //
func (w *workerS) Run(ctx context.Context, opts ...Opt) (err error) { return }            //nolint:revive
func (w *workerS) Actions() (ret map[string]bool)                   { return }            //nolint:revive
func (w *workerS) Name() string                                     { return "for-test" } //
func (*workerS) Version() string                                    { return "v0.0.0" }
func (*workerS) Root() *RootCommand                                 { return nil }
func (*workerS) Args() []string                                     { return nil }       //
func (w *workerS) SuggestRetCode() int                              { return w.retCode } //
func (w *workerS) ParsedState() ParsedState                         { return nil }
func (w *workerS) LoadedSources() (results []LoadedSources)         { return }

func (w *workerS) DoBuiltinAction(ctx context.Context, action ActionEnum, args ...any) (handled bool, err error) {
	return
}

//

//

//

// appS for testing only
type appS struct {
	Runner
	root  *RootCommand
	args  []string
	inCmd bool
	inFlg bool
}

func (s *appS) NewCommandBuilder(longTitle string, titles ...string) CommandBuilder {
	return s.Cmd(longTitle, titles...)
}

func (s *appS) NewFlagBuilder(longTitle string, titles ...string) FlagBuilder {
	return s.Flg(longTitle, titles...)
}

func (s *appS) Cmd(longTitle string, titles ...string) CommandBuilder { //nolint:revive
	s.inCmd = true
	// return newCommandBuilder(s, longTitle, titles...)
	return nil
}

func (s *appS) With(cb func(app App)) { //nolint:revive
	cb(s)
}

func (s *appS) WithOpts(opts ...Opt) App {
	return s
}

func (s *appS) Flg(longTitle string, titles ...string) FlagBuilder { //nolint:revive
	s.inFlg = true
	// return newFlagBuilder(s, longTitle, titles...)
	return nil
}

func (s *appS) AddCmd(f func(b CommandBuilder)) App { //nolint:revive
	// b := newCommandBuilder(s, "")
	// defer b.Build()
	// cb(b)
	return s
}

func (s *appS) AddFlg(cb func(b FlagBuilder)) App { //nolint:revive
	// b := newFlagBuilder(s, "")
	// defer b.Build()
	// cb(b)
	return s
}

func (s *appS) NewCmdFrom(from *CmdS, cb func(b CommandBuilder)) App { //nolint:revive
	// b := newCommandBuilderFrom(from, s, "")
	// defer b.Build()
	// cb(b)
	return s
}

func (s *appS) NewFlgFrom(from *CmdS, defaultValue any, cb func(b FlagBuilder)) App { //nolint:revive
	// b := newFlagBuilderFrom(from, s, defaultValue, "")
	// defer b.Build()
	// cb(b)
	return s
}

func (s *appS) RootBuilder(cb func(b CommandBuilder)) App { return s }

func (s *appS) addCommand(child *CmdS) { //nolint:unused
	s.inCmd = false
	if cx, ok := s.root.Cmd.(*CmdS); ok {
		cx.AddSubCommand(child)
	}
}

func (s *appS) addFlag(child *Flag) { //nolint:unused
	s.inFlg = false
	if cx, ok := s.root.Cmd.(*CmdS); ok {
		cx.AddFlag(child)
	}
}

func (s *appS) Info(name, version string, desc ...string) App {
	s.ensureNewApp()
	if s.root.AppName == "" {
		s.root.AppName = name
	}
	if s.root.Version == "" {
		s.root.Version = version
	}
	if cx, ok := s.root.Cmd.(*CmdS); ok {
		cx.SetDescription("", desc...)
	}
	return s
}

func (s *appS) Examples(examples ...string) App {
	s.ensureNewApp()
	if cx, ok := s.root.Cmd.(*CmdS); ok {
		cx.SetExamples(examples...)
	}
	return s
}

func (s *appS) Copyright(copyright string) App {
	s.ensureNewApp()
	s.root.Copyright = copyright
	return s
}

func (s *appS) Author(author string) App {
	s.ensureNewApp()
	s.root.Author = author
	return s
}

func (s *appS) Description(desc string) App {
	s.ensureNewApp()
	s.root.SetDesc(desc)
	return s
}

func (s *appS) Header(headerLine string) App {
	s.ensureNewApp()
	s.root.HeaderLine = headerLine
	return s
}

func (s *appS) Footer(footerLine string) App {
	s.ensureNewApp()
	s.root.FooterLine = footerLine
	return s
}

func (s *appS) OnAction(handler OnInvokeHandler) App {
	return s
}

func (s *appS) SetRootCommand(root *RootCommand) App {
	s.root = root
	return s
}

func (s *appS) WithRootCommand(cb func(root *RootCommand)) App {
	cb(s.root)
	return s
}

func (s *appS) RootCommand() *RootCommand { return s.root }

func (s *appS) Name() string       { return s.root.AppName }
func (s *appS) Version() string    { return s.root.Version }
func (s *appS) Worker() Runner     { return s.Runner }
func (s *appS) Root() *RootCommand { return s.root }
func (s *appS) Args() []string     { return s.args }

func (s *appS) ensureNewApp() App { //nolint:unparam
	if s.root == nil {
		s.root = &RootCommand{
			AppName: conf.AppName,
			Version: conf.Version,
			app:     s,
			// Copyright:  "",
			// Author:     "",
			// HeaderLine: "",
			// FooterLine: "",
			// CmdS:    nil,
		}
	}
	if s.root.Cmd == nil {
		s.root.Cmd = new(CmdS)
		s.root.Cmd.SetName(s.root.AppName)
	}
	return s
}

func (s *appS) Build() {
	type setRoot interface {
		SetRoot(root *RootCommand, args []string)
	}
	if sr, ok := s.Runner.(setRoot); ok {
		ctx := context.Background()
		if cx, ok := s.root.Cmd.(*CmdS); ok {
			cx.EnsureTree(ctx, s, s.root)
		}
		sr.SetRoot(s.root, s.args)
	}
}

func (s *appS) Run(ctx context.Context, opts ...Opt) (err error) {
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
