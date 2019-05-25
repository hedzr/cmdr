/*
 * Copyright © 2019 Hedzr Yeh.
 */

package daemon

import "github.com/hedzr/cmdr"

var (
	DaemonServerCommands = &cmdr.Command{
		BaseOpt: cmdr.BaseOpt{
			// Name:        "server",
			Short:       "s",
			Full:        "server",
			Aliases:     []string{"serve", "svr", "daemon"},
			Description: "server ops: for linux daemon.",
		},
		SubCommands: []*cmdr.Command{
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "s",
					Full:        "start",
					Aliases:     []string{"run", "startup"},
					Description: "startup this system service/daemon.",
					Action:      daemonStart,
					Flags: []*cmdr.Flag{
						{
							BaseOpt: cmdr.BaseOpt{
								Short:       "f",
								Full:        "foreground",
								Aliases:     []string{"fg"},
								Description: "run on foreground, NOT daemonized.",
							},
							DefaultValue: false,
						},
					},
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "t",
					Full:        "stop",
					Aliases:     []string{"stp", "halt", "pause"},
					Description: "stop this system service/daemon.",
					Action:      daemonStop,
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "re",
					Full:        "restart",
					Aliases:     []string{"reload"},
					Description: "restart this system service/daemon.",
					Action:      daemonRestart,
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Full:        "status",
					Aliases:     []string{"st"},
					Description: "display its running status as a system service/daemon.",
					Action:      daemonStatus,
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "i",
					Full:        "install",
					Aliases:     []string{"setup"},
					Description: "install as a system service/daemon.",
					Action:      daemonInstall,
					Flags: []*cmdr.Flag{
						{
							BaseOpt: cmdr.BaseOpt{
								Short:       "s",
								Full:        "systemd",
								Aliases:     []string{"sys"},
								Description: "install as a systemd service.",
							},
							DefaultValue: true,
						},
					},
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "u",
					Full:        "uninstall",
					Aliases:     []string{"remove"},
					Description: "remove from a system service/daemon.",
					Action:      daemonUninstall,
				},
			},
		},
	}
)
