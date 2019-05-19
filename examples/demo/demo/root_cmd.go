/*
 * Copyright © 2019 Hedzr Yeh.
 */

package demo

import (
	"github.com/hedzr/cmdr"
)

var (
	rootCmd = &cmdr.RootCommand{
		Command: cmdr.Command{
			BaseOpt: cmdr.BaseOpt{
				Name:  "demo",
				Flags: []*cmdr.Flag{},
			},
			SubCommands: []*cmdr.Command{
				// generatorCommands,
				serverCommands,
				msCommands,
			},
		},

		AppName:    "demo",
		Version:    cmdr.Version,
		VersionInt: cmdr.VersionInt,
		Copyright:  "austr is an effective devops tool",
		Author:     "Hedzr Yeh <hedzrz@gmail.com>",
	}

	serverCommands = &cmdr.Command{
		BaseOpt: cmdr.BaseOpt{
			// Name:        "server",
			Short:       "s",
			Full:        "server",
			Aliases:     []string{"serve", "svr"},
			Description: "server ops: for linux service/daemon.",
		},
		SubCommands: []*cmdr.Command{
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "s",
					Full:        "start",
					Aliases:     []string{"run", "startup"},
					Description: "startup this system service/daemon.",
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "t",
					Full:        "stop",
					Aliases:     []string{"stp", "halt", "pause"},
					Description: "stop this system service/daemon.",
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "r",
					Full:        "restart",
					Aliases:     []string{"reload"},
					Description: "restart this system service/daemon.",
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
