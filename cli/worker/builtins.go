package worker

import (
	"context"
	"fmt"
	"runtime"

	"github.com/hedzr/is/states"
	logz "github.com/hedzr/logg/slog"

	"github.com/hedzr/cmdr/v2/cli"
)

func (w *workerS) addBuiltinCommands(root *cli.RootCommand) (err error) { //nolint:unparam //unified form
	app, cmd := root.App(), root.Command
	w.builtinCmdrs(app, cmd)
	w.builtinSBOM(app, cmd)
	w.builtinGenerators(app, cmd)
	w.builtinVerboses(app, cmd)
	w.builtinVersions(app, cmd)
	w.builtinHelps(app, cmd)
	return
}

func (w *workerS) builtinVersions(app cli.App, p *cli.Command) {
	app.NewCmdFrom(p, func(b cli.CommandBuilder) {
		b.Titles("version", "ver", "versions").
			Description("Show app versions information").
			Group(cli.SysMgmtGroup).
			Hidden(true, false).
			OnAction(func(ctx context.Context, cmd cli.BaseOptI, args []string) (err error) { //nolint:revive
				w.actionsMatched |= actionShowVersion
				// w.showVersion(cmd, args)
				// return cli.ErrShouldStop
				return
			})
	})
	app.NewFlgFrom(p, false, func(b cli.FlagBuilder) {
		b.Titles("version", "V", "ver", "versions").
			Description("Show app versions information").
			Group(cli.SysMgmtGroup).
			Hidden(true, false).
			OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) { //nolint:revive
				// app.Store().Set("app.cmd.show.version", true)
				w.actionsMatched |= actionShowVersion
				return
			}).
			CompCircuitBreak(true).
			CompJustOnce(true)
	})
	app.NewFlgFrom(p, "", func(b cli.FlagBuilder) {
		b.Titles("version-sim", "VS", "ver-sim", "version-simulate").
			Description("Simulate a faked version for this app").
			Group(cli.SysMgmtGroup).
			Hidden(true, true).
			OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) { //nolint:revive
				// app.Store().Set("app.version.simulate", value)
				w.versionSimulate = fmt.Sprint(hitState.Value)
				return
			})
	})
	app.NewFlgFrom(p, false, func(b cli.FlagBuilder) {
		b.Titles("built-info", "#", "bi").
			Description("Show the building information of this app").
			Group(cli.SysMgmtGroup).
			Hidden(true, true).
			OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) { //nolint:revive
				// app.Store().Set("app.version.simulate", value)
				w.actionsMatched |= actionShowBuiltInfo
				return
			}).
			CompCircuitBreak(true).
			CompJustOnce(true)
	})
}

func (w *workerS) builtinHelps(app cli.App, p *cli.Command) { //nolint:revive
	app.NewCmdFrom(p, func(b cli.CommandBuilder) {
		b.Titles("help", "h", "info", "usage", "__completion", "__complete").
			Description("Show help system for commands").
			Group(cli.SysMgmtGroup).
			Hidden(true, false).
			OnMatched(func(c *cli.Command, position int, hitState *cli.MatchState) (err error) { //nolint:revive
				w.inCompleting = true
				logz.SetLevel(logz.ErrorLevel) // disable trace, debug, info, warn messages
				return
			}).
			OnAction(func(ctx context.Context, cmd cli.BaseOptI, args []string) (err error) {
				err = w.helpSystemAction(ctx, cmd, args)
				// w.actionsMatched |= actionShowHelpScreen
				return // return cli.ErrShouldStop
			})
	})

	app.NewFlgFrom(p, false, func(b cli.FlagBuilder) {
		b.Titles("help", "h", "info", "usage").
			Description("Show this help screen (-?)").
			Group(cli.SysMgmtGroup).
			Hidden(false, false).
			OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) { //nolint:revive
				// app.Store().Set("app.version.simulate", value)
				w.actionsMatched |= actionShowHelpScreen
				return
			}).
			CompCircuitBreak(true).
			CompJustOnce(true).
			EnvVars("HELP").
			ExtraShorts("?")
	})

	m := map[string]bool{
		"linux": true, "darwin": true,
	}
	if _, ok := m[runtime.GOOS]; ok {
		app.NewFlgFrom(p, false, func(b cli.FlagBuilder) {
			b.Titles("manual", "man").
				Description("Show help screen in manpage format (INSTALL NEEDED!)").
				Group(cli.SysMgmtGroup).
				Hidden(true, true).
				OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) { //nolint:revive
					// app.Store().Set("app.version.simulate", value)
					w.actionsMatched |= actionShowHelpScreenAsMan
					return
				}).
				EnvVars("MAN")
		})
	}

	app.NewFlgFrom(p, false, func(b cli.FlagBuilder) {
		b.Titles("tree").
			Description("Show help screen in manpage format (INSTALL NEEDED!)").
			Group(cli.SysMgmtGroup).
			Hidden(true, true).
			OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) { //nolint:revive
				// app.Store().Set("app.version.simulate", value)
				w.actionsMatched |= actionShowTree
				return
			}).
			EnvVars("TREE").
			DoubleTildeOnly(true)
	})

	// find config file loader at first
	found := false
	for _, l := range w.Loaders {
		if _, found = l.(interface {
			LoadFile(filename string, app cli.App) (err error)
		}); found {
			break
		}
	}
	if found {
		app.NewFlgFrom(p, "", func(b cli.FlagBuilder) {
			b.Titles("config").
				Description("Load your config file").
				Group(cli.SysMgmtGroup).
				Hidden(true, false).
				PlaceHolder("FILE").
				Examples(`
$ {{.AppName}} --configci/etc/demo-yy ~~debug
	try loading config from 'ci/etc/demo-yy', noted that assumes a child folder 'conf.d' should be exists
$ {{.AppName}} --config=ci/etc/demo-yy/any.yml ~~debug
	try loading config from 'ci/etc/demo-yy/any.yml', noted that assumes a child folder 'conf.d' should be exists
`).
				OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) { //nolint:revive
					// app.Store().Set("app.version.simulate", value)
					var ok bool
					w.configFile, ok = hitState.Value.(string)
					if !ok {
						err = fmt.Errorf("value is not a string. [value=%v]", hitState.Value)
					}
					return
				}).
				EnvVars("CONFIG", "CONF_FILE")
		})
	}
}

func (w *workerS) builtinVerboses(app cli.App, p *cli.Command) { //nolint:revive
	app.NewFlgFrom(p, false, func(b cli.FlagBuilder) {
		b.Titles("verbose", "v").
			Description("Show more progress/debug info with verbose mode").
			Group(cli.SysMgmtGroup).
			Hidden(true, true).
			EnvVars("VERBOSE").
			OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) { //nolint:revive
				// app.Store().Set("app.mode.verbose", value)
				if v, ok := hitState.Value.(bool); ok {
					states.Env().SetVerboseMode(v)
					states.Env().SetVerboseCount(f.GetTriggeredTimes())
				}
				return
			})
	})
	app.NewFlgFrom(p, false, func(b cli.FlagBuilder) {
		b.Titles("quiet", "q").
			Description("No more screen output").
			Group(cli.SysMgmtGroup).
			Hidden(true, true).
			EnvVars("QUIET", "SILENT").
			OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) { //nolint:revive
				// app.Store().Set("app.mode.verbose", value)
				if v, ok := hitState.Value.(bool); ok {
					states.Env().SetQuietMode(v)
					states.Env().SetQuietCount(f.GetTriggeredTimes())
				}
				return
			})
	})

	app.NewFlgFrom(p, false, func(b cli.FlagBuilder) {
		b.Titles("debug", "D").
			Description("Get into debug mode").
			Group(cli.SysMgmtGroup).
			Hidden(true, true).
			EnvVars("DEBUG").
			OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) { //nolint:revive
				// app.Store().Set("app.mode.verbose", value)
				if v, ok := hitState.Value.(bool); ok {
					states.Env().SetDebugMode(v)
					states.Env().SetDebugLevel(hitState.HitTimes)
					if hitState.DblTilde {
						w.actionsMatched |= actionShowDebug // ~~debug to show debug states screen
					}
				}
				return
			})
	})

	app.NewFlgFrom(p, "", func(b cli.FlagBuilder) {
		b.Titles("debug-output", "DO").
			Description("Store the ~~debug outputs into file").
			Group(cli.SysMgmtGroup).
			Hidden(true, true).
			EnvVars("DEBUG_OUTPUT").
			OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) { //nolint:revive
				// app.Store().Set("app.mode.verbose", value)
				if v, ok := hitState.Value.(string); ok {
					w.debugOutputFile = v
				}
				return
			})
	})

	mutualExclusives := []string{"raw", "value-type", "more", "env"}

	app.NewFlgFrom(p, false, func(b cli.FlagBuilder) {
		b.Titles("env").
			Description("Dump environment info in '~~debug' mode").
			Group(cli.SysMgmtGroup).
			Hidden(true, true).
			// EnvVars("ENV").
			OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) { //nolint:revive
				w.actionsMatched |= actionShowDebugEnv
				return
			}).
			CompPrerequisites("debug").
			CompMutualExclusives(mutualExclusives...)
	})
	app.NewFlgFrom(p, false, func(b cli.FlagBuilder) {
		b.Titles("more").
			Description("Dump more info in '~~debug' mode").
			Group(cli.SysMgmtGroup).
			Hidden(true, true).
			// EnvVars("MORE").
			OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) { //nolint:revive
				w.actionsMatched |= actionShowDebugMore
				return
			}).
			CompPrerequisites("debug").
			CompMutualExclusives(mutualExclusives...)
	})
	app.NewFlgFrom(p, false, func(b cli.FlagBuilder) {
		b.Titles("raw").
			Description("Dump the option value in raw mode (with golang data structure, without envvar expanding)").
			Group(cli.SysMgmtGroup).
			Hidden(true, true).
			EnvVars("RAW").
			OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) { //nolint:revive
				w.actionsMatched |= actionShowDebugRaw
				return
			}).
			CompPrerequisites("debug").
			CompMutualExclusives(mutualExclusives...)
	})
	app.NewFlgFrom(p, false, func(b cli.FlagBuilder) {
		b.Titles("value-type").
			Description("Dump the option value type in '~~debug' mode").
			Group(cli.SysMgmtGroup).
			Hidden(true, true).
			// EnvVars("RAW").
			OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) { //nolint:revive
				w.actionsMatched |= actionShowDebugValueType
				return
			}).
			CompPrerequisites("debug").
			CompMutualExclusives(mutualExclusives...)
	})
}

func (w *workerS) builtinCmdrs(app cli.App, p *cli.Command) { //nolint:revive
	app.NewFlgFrom(p, false, func(b cli.FlagBuilder) {
		b.Titles("strict-mode", "").
			Description("<mark>Strict mode</mark> for '<code>cmdr</code>'").
			Group(cli.SysMgmtGroup).
			Hidden(true, true).
			EnvVars("STRICT").
			OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) { //nolint:revive
				if v, ok := hitState.Value.(bool); ok {
					w.strictMode = v
					w.strictModeLevel = hitState.HitTimes
				}
				return
			})
	})
	app.NewFlgFrom(p, false, func(b cli.FlagBuilder) {
		b.Titles("no-env-overrides", "").
			Description("No env var overrides for '<code>cmdr</code>'").
			Group(cli.SysMgmtGroup).
			Hidden(true, true).
			Deprecated("v0.1.1").
			OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) { //nolint:revive
				if v, ok := hitState.Value.(bool); ok {
					w.noLoadEnv = v
				}
				return
			})
	})
	app.NewFlgFrom(p, false, func(b cli.FlagBuilder) {
		b.Titles("no-color", "nc").
			Description("<i>No color</i> output for '<code>cmdr</code>'").
			Group(cli.SysMgmtGroup).
			Hidden(true, true).
			EnvVars("NO_COLOR", "NOCOLOR").
			OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) { //nolint:revive
				// app.Store().Set("app.mode.verbose", value)
				if v, ok := hitState.Value.(bool); ok {
					states.Env().SetNoColorMode(v)
					states.Env().SetNoColorCount(hitState.HitTimes)
				}
				return
			})
	})
}

func (w *workerS) builtinGenerators(app cli.App, p *cli.Command) { //nolint:revive
	app.NewCmdFrom(p, func(bb cli.CommandBuilder) {
		bb.Titles("generate", "g", "gen", "generator").
			Description("Generators for this app", `
[cmdr] includes multiple generators like:

- linux man page generator
- shell completion script generator
- more...

			`).
			Examples(`
$ {{.AppName}} gen sh --bash
			generate bash completion script
$ {{.AppName}} gen shell --auto
			generate shell completion script with detecting on current shell environment.
$ {{.AppName}} gen sh
			generate shell completion script with detecting on current shell environment.
$ {{.AppName}} gen man
			generate linux manual (man page)
			`).
			Group(cli.SysMgmtGroup).
			Hidden(true, true).
			OnMatched(func(c *cli.Command, position int, hitState *cli.MatchState) (err error) { //nolint:revive
				return
			}).
			OnAction((&genS{}).onAction)

		bb.Cmd("manual", "m", "man").
			Description("Generate Linux Manual Documentations").
			Group(cli.SysMgmtGroup).
			Hidden(true, true).
			OnMatched(func(c *cli.Command, position int, hitState *cli.MatchState) (err error) { //nolint:revive
				return
			}).
			OnAction((&genManS{}).onAction).
			With(func(b cli.CommandBuilder) {
				b.Flg("dir", "d").
					Default("").
					Description("The output directory").
					Group("Output").
					Hidden(true, true).
					PlaceHolder("DIR").
					Build()

				b.Flg("type", "t").
					Default(1).
					Description("Linux man type").
					Hidden(true, true).
					HeadLike(true, 1, 9).
					Build()
			})

		bb.Cmd("doc", "d", "docx", "tex", "pdf", "markdown").
			Description("Generate documentations").
			Group(cli.SysMgmtGroup).
			Hidden(true, true).
			OnMatched(func(c *cli.Command, position int, hitState *cli.MatchState) (err error) { //nolint:revive
				return
			}).
			OnAction((&genDocS{}).onAction).
			With(func(b cli.CommandBuilder) {
				b.Flg("dir", "d").
					Default("").
					Description("The output directory").
					Group("Output").
					Hidden(true, true).
					PlaceHolder("DIR").
					Build()
			})

		bb.Cmd("shell", "s", "sh", "bash", "zsh", "fish", "elvish", "fig", "powershell", "ps").
			Description("Generate the shell completion script or install it").
			Group(cli.SysMgmtGroup).
			Hidden(true, true).
			OnMatched(func(c *cli.Command, position int, hitState *cli.MatchState) (err error) { //nolint:revive
				return
			}).
			OnAction((&genShS{}).onAction).
			With(func(b cli.CommandBuilder) {
				b.Flg("dir", "d").
					Default("").
					Description("The output directory").
					Group("Output").
					Hidden(true, true).
					PlaceHolder("DIR").
					Build()

				b.Flg("output", "o").
					Default("").
					Description("The output filename").
					Group("Output").
					Hidden(true, true).
					PlaceHolder("FILE").
					Build()

				b.Flg("auto", "a").
					Default(true).
					Description("Generate auto completion script to fit for your current env").
					ToggleGroup("Shell").
					Hidden(true, true).
					Build()

				b.Flg("zsh", "z").
					Default(false).
					Description("Generate auto completion script for Zsh").
					ToggleGroup("Shell").
					Hidden(true, true).
					Build()

				b.Flg("bash", "b").
					Default(false).
					Description("Generate auto completion script for Bash").
					ToggleGroup("Shell").
					Hidden(true, true).
					Build()

				b.Flg("fish", "f").
					Default(false).
					Description("Generate auto completion script for Fish").
					ToggleGroup("Shell").
					Hidden(true, true).
					Build()

				b.Flg("powershell", "p").
					Default(false).
					Description("Generate auto completion script for PowerShell").
					ToggleGroup("Shell").
					Hidden(true, true).
					Build()

				b.Flg("elvish", "e").
					Default(false).
					Description("Generate auto completion script for Elvish [TODO]").
					ToggleGroup("Shell").
					Hidden(true, true).
					Build()

				b.Flg("fig", "f").
					Default(false).
					Description("Generate auto completion script for fig-shell [TODO]").
					ToggleGroup("Shell").
					Hidden(true, true).
					Build()
			})
	})
}

func (w *workerS) builtinSBOM(app cli.App, p *cli.Command) { //nolint:revive
	app.NewCmdFrom(p, func(bb cli.CommandBuilder) {
		bb.Titles("sbom", "", "").
			Description("Show SBOM Info", ``).
			Group(cli.SysMgmtGroup).
			Hidden(true, true).
			OnMatched(func(c *cli.Command, position int, hitState *cli.MatchState) (err error) { //nolint:revive
				return
			}).
			OnAction((&sbomS{}).onAction)
	})
}

func (w *workerS) uniAddCmd(cmd *cli.Command, callbacks ...func(cc *cli.Command)) { //nolint:revive,unused
	cmd.AddSubCommand(new(cli.Command), callbacks...)
}

//

//

//
