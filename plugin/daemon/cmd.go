/*
 * Copyright © 2019 Hedzr Yeh.
 */

package daemon

import "github.com/hedzr/cmdr"

var (
	// DaemonServerCommands defines a group of sub-commands for daemon operations.
	DaemonServerCommands = &cmdr.Command{
		BaseOpt: cmdr.BaseOpt{
			// Name:        "server",
			Short:       "s",
			Full:        "server",
			Aliases:     []string{"serve", "svr", "daemon"},
			Description: "server ops: for linux daemon.",
			Group:       "Daemonization",
		},
		SubCommands: []*cmdr.Command{
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "s",
					Full:        "start",
					Aliases:     []string{"run", "startup"},
					Description: "startup this system service/daemon.",
					Action:      daemonStart,
					LongDescription: `**start** command make program running as a daemon background.
**run** command make program running in current tty foreground.
`,
					Examples: `
$ {{.AppName}} start
					make program running as a daemon background.
$ {{.AppName}} start --foreground
					make program running in current tty foreground.
$ {{.AppName}} run
					make program running in current tty foreground.
$ {{.AppName}} stop
					stop daemonized program.
$ {{.AppName}} reload
					send signal to trigger program reload its configurations.
$ {{.AppName}} status
					display the daemonized program running status.
$ {{.AppName}} install [--systemd]
					install program as a systemd service.
$ {{.AppName}} uninstall
					remove the installed systemd service.
`,
				},
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
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "t",
					Full:        "stop",
					Aliases:     []string{"stp", "halt", "pause"},
					Description: "stop this system service/daemon.",
					Action:      daemonStop,
				},
				Flags: []*cmdr.Flag{
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "1",
							Full:        "hup",
							Description: "send SIGHUP - to reload service",
						},
						DefaultValue: false,
					},
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "3",
							Full:        "quit",
							Description: "send SIGQUIT - to quit service gracefully",
						},
						DefaultValue: false,
					},
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "9",
							Full:        "kill",
							Description: "send SIGKILL - to quit service unconditionally",
						},
						DefaultValue: false,
					},
					{
						BaseOpt: cmdr.BaseOpt{
							Short:       "15",
							Full:        "term",
							Description: "send SIGTERM - to quit service gracefully",
						},
						DefaultValue: false,
					},
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
					Short:       "ss",
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
					Group:       "Config",
					Action:      daemonInstall,
				},
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
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "u",
					Full:        "uninstall",
					Aliases:     []string{"remove"},
					Description: "remove from a system service/daemon.",
					Group:       "Config",
					Action:      daemonUninstall,
				},
			},
		},
	}
)
