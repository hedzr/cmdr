/*
 * Copyright © 2019 Hedzr Yeh.
 */

package demo

import (
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
	overview = ``
)

var (
	rootCmd = &cmdr.RootCommand{
		Command: cmdr.Command{
			BaseOpt: cmdr.BaseOpt{
				Name:            appName,
				Description:     desc,
				LongDescription: longDesc,
				Examples:        examples,
				Flags:           []*cmdr.Flag{},
			},
			SubCommands: []*cmdr.Command{
				// generatorCommands,
				// serverCommands,
				msCommands,
			},
		},

		AppName:    appName,
		Version:    cmdr.Version,
		VersionInt: cmdr.VersionInt,
		Copyright:  copyright,
		Author:     "Hedzr Yeh <hedzrz@gmail.com>",
	}

	msCommands = &cmdr.Command{
		BaseOpt: cmdr.BaseOpt{
			Name:        "microservices",
			Full:        "ms",
			Aliases:     []string{"microservice", "micro-service"},
			Description: "micro-service operations...",
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
		},
		SubCommands: []*cmdr.Command{
			{
				BaseOpt: cmdr.BaseOpt{
					// Short:       "t",
					Full:        "tags",
					Aliases:     []string{},
					Description: "tags op.",
					Flags: []*cmdr.Flag{
						{
							BaseOpt: cmdr.BaseOpt{
								Short:                   "n",
								Full:                    "name",
								Description:             "name of the service",
								DefaultValuePlaceholder: "NAME",
							},
							DefaultValue: "",
						},
						{
							BaseOpt: cmdr.BaseOpt{
								Short:                   "i",
								Full:                    "id",
								Description:             "unique id of the service",
								DefaultValuePlaceholder: "ID",
							},
							DefaultValue: "",
						},
					},
				},
				SubCommands: []*cmdr.Command{
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "ls",
							Full:        "list",
							Aliases:     []string{"l", "lst", "dir"},
							Description: "list tags.",
						},
					},
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "a",
							Full:        "add",
							Aliases:     []string{"create", "new"},
							Description: "add tags.",
							Flags: []*cmdr.Flag{
								{
									BaseOpt: cmdr.BaseOpt{
										Short:                   "ls",
										Full:                    "list",
										Aliases:                 []string{"l", "lst", "dir"},
										Description:             "a comma list to be added",
										DefaultValuePlaceholder: "LIST",
									},
									DefaultValue: []string{},
								},
							},
						},
					},
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "r",
							Full:        "rm",
							Aliases:     []string{"remove", "erase", "delete", "del"},
							Description: "remove tags.",
							Flags: []*cmdr.Flag{
								{
									BaseOpt: cmdr.BaseOpt{
										Short:                   "ls",
										Full:                    "list",
										Aliases:                 []string{"l", "lst", "dir"},
										Description:             "a comma list to be added.",
										DefaultValuePlaceholder: "LIST",
									},
									DefaultValue: []string{},
								},
							},
						},
					},
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "t",
							Full:        "toggle",
							Aliases:     []string{"tog", "switch"},
							Description: "toggle tags for ms.",
							Flags: []*cmdr.Flag{
								{
									BaseOpt: cmdr.BaseOpt{
										Short:                   "s",
										Full:                    "set",
										DefaultValuePlaceholder: "LIST",
									},
									DefaultValue: []string{},
								},
								{
									BaseOpt: cmdr.BaseOpt{
										Short:                   "u",
										Full:                    "unset",
										DefaultValuePlaceholder: "LIST",
									},
									DefaultValue: []string{},
								},
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
)
