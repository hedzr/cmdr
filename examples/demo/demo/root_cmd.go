/*
 * Copyright © 2019 Hedzr Yeh.
 */

package demo

import (
	"fmt"

	"github.com/hedzr/cmdr"
)

const (
	appName   = "demo"
	copyright = "demo.austr is an effective devops tool"
	desc      = "demo.austr is an effective devops tool. It make an demo application for `cmdr`."
	longDesc  = "demo.austr is an effective devops tool. It make an demo application for `cmdr`."
	examples  = `
$ {{.AppName}} gen shell [--bash|--zsh|--auto]
  generate bash/shell completion scripts
$ {{.AppName}} gen man
  generate linux man page 1
$ {{.AppName}} --help
  show help screen.
`
	// overview = ``
)

var (
	rootCmd = &cmdr.RootCommand{
		Command: cmdr.Command{
			BaseOpt: cmdr.BaseOpt{
				Name:            appName,
				Description:     desc,
				LongDescription: longDesc,
				Examples:        examples,
			},
			Flags: []*cmdr.Flag{},
			SubCommands: []*cmdr.Command{
				// generatorCommands,
				// serverCommands,
				msCommands,
				testCommands,
				{
					BaseOpt: cmdr.BaseOpt{
						Short:       "xy",
						Full:        "xy-print",
						Description: `test terminal control sequences`,
						Action: func(cmd *cmdr.Command, args []string) (err error) {
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
						},
					},
				},
				{
					BaseOpt: cmdr.BaseOpt{
						Short:       "mx",
						Full:        "mx-test",
						Description: `test new features`,
						Action: func(cmd *cmdr.Command, args []string) (err error) {
							fmt.Printf("*** Got pp: %s\n", cmdr.GetString("app.mx-test.password"))
							fmt.Printf("*** Got msg: %s\n", cmdr.GetString("app.mx-test.message"))
							return
						},
					},
					Flags: []*cmdr.Flag{
						{
							BaseOpt: cmdr.BaseOpt{
								Short:       "pp",
								Full:        "password",
								Description: "the password requesting.",
							},
							DefaultValue: "",
							ExternalTool: cmdr.ExternalToolPasswordInput,
						},
						{
							BaseOpt: cmdr.BaseOpt{
								Short:       "m",
								Full:        "message",
								Description: "the message requesting.",
							},
							DefaultValue: "",
							ExternalTool: cmdr.ExternalToolEditor,
						},
					},
				},
			},
		},

		AppName:    appName,
		Version:    cmdr.Version,
		VersionInt: cmdr.VersionInt,
		Copyright:  copyright,
		Author:     "Hedzr Yeh <hedzr@duck.com>",
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
					Short:       "m",
					Full:        "money",
					Description: "a placeholder flag",
				},
				DefaultValue: false,
			},
		},
		SubCommands: []*cmdr.Command{
			{
				BaseOpt: cmdr.BaseOpt{
					// Short:       "t",
					Full:        "tags",
					Aliases:     []string{},
					Description: "tags op.",
				},
				Flags: []*cmdr.Flag{
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "n",
							Full:        "name",
							Description: "name of the service",
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
				},
				SubCommands: []*cmdr.Command{
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "ls",
							Full:        "list",
							Aliases:     []string{"l", "lst", "dir"},
							Description: "list tags.",
							Group:       "2333.List",
						},
					},
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "a",
							Full:        "add",
							Aliases:     []string{"create", "new"},
							Description: "add tags.",
							Deprecated:  "0.2.1",
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
							Short:       "t",
							Full:        "toggle",
							Aliases:     []string{"tog", "switch"},
							Description: "toggle tags for ms.",
						},
						Flags: []*cmdr.Flag{
							{
								BaseOpt: cmdr.BaseOpt{
									Short: "s",
									Full:  "set",
								},
								DefaultValue:            []string{},
								DefaultValuePlaceholder: "LIST",
							},
							{
								BaseOpt: cmdr.BaseOpt{
									Short: "u",
									Full:  "unset",
								},
								DefaultValue:            []string{},
								DefaultValuePlaceholder: "LIST",
							},
						},
					},
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "l",
					Full:        "list",
					Aliases:     []string{"ls", "lst", "dir"},
					Description: "list services.",
				},
			},
		},
	}

	testCommands = &cmdr.Command{
		BaseOpt: cmdr.BaseOpt{
			Short:       "t",
			Full:        "test",
			Description: "test operations...",
		}, Flags: []*cmdr.Flag{
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "m",
					Full:        "message",
					Description: "a placeholder flag",
				},
				DefaultValue: "",
				ExternalTool: cmdr.ExternalToolEditor,
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "g",
					Full:        "good",
					Description: "good is a placeholder flag",
					Deprecated:  "0.0.1",
				},
				DefaultValue: "",
				ExternalTool: cmdr.ExternalToolEditor,
			},
		},
	}
)
