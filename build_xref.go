/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"github.com/hedzr/cmdr/conf"
	"strings"
	"time"
)

func (w *ExecWorker) buildRootCrossRefs(root *RootCommand) {
	// initializes the internal variables/members
	w.ensureCmdMembers(&root.Command)

	conf.AppName = root.AppName
	conf.Version = root.Version
	if len(conf.Buildstamp) == 0 {
		conf.Buildstamp = time.Now().Format(time.RFC1123)
	}

	w.attachVersionCommands(root)
	w.attachHelpCommands(root)
	w.attachVerboseCommands(root)
	w.attachGeneratorsCommands(root)
	w.attachCmdrCommands(root)

	w.buildCrossRefs(&root.Command)
}

func (w *ExecWorker) attachVersionCommands(root *RootCommand) {
	if EnableVersionCommands {
		if _, ok := root.allCmds[SysMgmtGroup]["version"]; !ok {
			cx := &Command{
				BaseOpt: BaseOpt{
					Full:        "version",
					Aliases:     []string{"ver"},
					Description: "Show the version of this app.",
					Action: func(cmd *Command, args []string) (err error) {
						w.showVersion()
						return ErrShouldBeStopException
					},
					Hidden: true,
					owner:  &root.Command,
				},
			}
			root.SubCommands = append(root.SubCommands, cx)
			root.allCmds[SysMgmtGroup]["version"] = cx
			root.allCmds[SysMgmtGroup]["ver"] = cx
		}
		if _, ok := root.allFlags[SysMgmtGroup]["version"]; !ok {
			root.allFlags[SysMgmtGroup]["version"] = &Flag{
				BaseOpt: BaseOpt{
					Short:       "V",
					Full:        "version",
					Aliases:     []string{"ver"},
					Description: "Show the version of this app.",
					// Hidden:      true,
					Action: func(cmd *Command, args []string) (err error) {
						w.showVersion()
						return ErrShouldBeStopException
					},
				},
				DefaultValue: false,
			}
			root.plainLongFlags["version"] = root.allFlags[SysMgmtGroup]["version"]
			root.plainLongFlags["ver"] = root.allFlags[SysMgmtGroup]["version"]
			root.plainShortFlags["V"] = root.allFlags[SysMgmtGroup]["version"]
		}
		if _, ok := root.allFlags[SysMgmtGroup]["version-sim"]; !ok {
			root.allFlags[SysMgmtGroup]["version-sim"] = &Flag{
				BaseOpt: BaseOpt{
					Full:        "version-sim",
					Aliases:     []string{"version-simulate"},
					Description: "Simulate a faked version number for this app.",
					Hidden:      true,
					Action: func(cmd *Command, args []string) (err error) {
						conf.Version = GetStringR("version-sim")
						Set("version", conf.Version) // set into option 'app.version' too.
						return
					},
				},
				DefaultValue: "",
			}
			root.plainLongFlags["version-sim"] = root.allFlags[SysMgmtGroup]["version-sim"]
			root.plainLongFlags["version-simulate"] = root.allFlags[SysMgmtGroup]["version-sim"]
		}
		if _, ok := root.allFlags[SysMgmtGroup]["build-info"]; !ok {
			root.allFlags[SysMgmtGroup]["build-info"] = &Flag{
				BaseOpt: BaseOpt{
					Full:        "#",
					Aliases:     []string{},
					Description: "Show the building information of this app.",
					Hidden:      true,
					Action: func(cmd *Command, args []string) (err error) {
						w.showBuildInfo()
						return ErrShouldBeStopException
					},
				},
				DefaultValue: false,
			}
			root.plainShortFlags["#"] = root.allFlags[SysMgmtGroup]["build-info"]
			root.plainLongFlags["build-info"] = root.allFlags[SysMgmtGroup]["build-info"]
		}
	}
}

func (w *ExecWorker) attachHelpCommands(root *RootCommand) {
	if EnableHelpCommands {
		if _, ok := root.allFlags[SysMgmtGroup]["help"]; !ok {
			root.allFlags[SysMgmtGroup]["help"] = &Flag{
				BaseOpt: BaseOpt{
					Short:       "h",
					Full:        "help",
					Aliases:     []string{"?", "helpme", "info", "usage"},
					Description: "Show this help screen",
					Hidden:      true,
					Action: func(cmd *Command, args []string) (err error) {
						// logrus.Debugf("-- helpCommand hit. printHelp and stop.")
						// printHelp(cmd)
						// return ErrShouldBeStopException
						return nil
					},
					owner: &root.Command,
				},
				DefaultValue: false,
			}
			root.plainLongFlags["help"] = root.allFlags[SysMgmtGroup]["help"]
			root.plainLongFlags["helpme"] = root.allFlags[SysMgmtGroup]["help"]
			root.plainLongFlags["info"] = root.allFlags[SysMgmtGroup]["help"]
			root.plainLongFlags["usage"] = root.allFlags[SysMgmtGroup]["help"]
			root.plainShortFlags["h"] = root.allFlags[SysMgmtGroup]["help"]
			root.plainShortFlags["?"] = root.allFlags[SysMgmtGroup]["help"]

			root.allFlags[SysMgmtGroup]["help-zsh"] = &Flag{
				BaseOpt: BaseOpt{
					Full:        "help-zsh",
					Description: "show help with zsh format, or others",
					Hidden:      true,
					owner:       &root.Command,
				},
				DefaultValue:            0,
				DefaultValuePlaceholder: "LEVEL",
			}
			root.allFlags[SysMgmtGroup]["help-bash"] = &Flag{
				BaseOpt: BaseOpt{
					Full:        "help-bash",
					Description: "show help with bash format, or others",
					Hidden:      true,
					owner:       &root.Command,
				},
				DefaultValue: false,
			}
			root.plainLongFlags["help-zsh"] = root.allFlags[SysMgmtGroup]["help-zsh"]
			root.plainLongFlags["help-bash"] = root.allFlags[SysMgmtGroup]["help-bash"]

			root.allFlags[SysMgmtGroup]["tree"] = &Flag{
				BaseOpt: BaseOpt{
					Full:        "tree",
					Description: "show a tree for all commands",
					Hidden:      true,
					owner:       &root.Command,
					Action:      dumpTreeForAllCommands,
				},
				DefaultValue: false,
			}
			root.plainLongFlags["tree"] = root.allFlags[SysMgmtGroup]["tree"]
		}
		if _, ok := root.allFlags[SysMgmtGroup]["config"]; !ok {
			root.allFlags[SysMgmtGroup]["config"] = &Flag{
				BaseOpt: BaseOpt{
					Full:        "config",
					Aliases:     []string{},
					Description: "load config files from where you specified",
					Action: func(cmd *Command, args []string) (err error) {
						// logrus.Debugf("-- --config hit. printHelp and stop.")
						// return ErrShouldBeStopException
						return nil
					},
					owner: &root.Command,
					// TODO how to display examples section for a flag?
					Examples: `
$ {{.AppName}} --configci/etc/demo-yy ~~debug
	try loading config from 'ci/etc/demo-yy', noted that assumes a child folder 'conf.d' should be exists
$ {{.AppName}} --config=ci/etc/demo-yy/any.yml ~~debug
	try loading config from 'ci/etc/demo-yy/any.yml', noted that assumes a child folder 'conf.d' should be exists
`,
				},
				DefaultValue:            "",
				DefaultValuePlaceholder: "[Locations of config files]",
			}
			root.plainLongFlags["config"] = root.allFlags[SysMgmtGroup]["config"]
		}
	}
}

func (w *ExecWorker) attachVerboseCommands(root *RootCommand) {
	if EnableVerboseCommands {
		if _, ok := root.allFlags[SysMgmtGroup]["verbose"]; !ok {
			root.allFlags[SysMgmtGroup]["verbose"] = &Flag{
				BaseOpt: BaseOpt{
					Short:       "v",
					Full:        "verbose",
					Aliases:     []string{"vv", "vvv"},
					Description: "Show this help screen",
					// Hidden:      true,
					owner: &root.Command,
				},
				DefaultValue: false,
			}
			root.plainLongFlags["verbose"] = root.allFlags[SysMgmtGroup]["verbose"]
			root.plainLongFlags["vvv"] = root.allFlags[SysMgmtGroup]["verbose"]
			root.plainLongFlags["vv"] = root.allFlags[SysMgmtGroup]["verbose"]
			root.plainShortFlags["v"] = root.allFlags[SysMgmtGroup]["verbose"]
		}
		if _, ok := root.allFlags[SysMgmtGroup]["quiet"]; !ok {
			root.allFlags[SysMgmtGroup]["quiet"] = &Flag{
				BaseOpt: BaseOpt{
					Short:       "q",
					Full:        "quiet",
					Aliases:     []string{},
					Description: "No more screen output.",
					// Hidden:      true,
					owner: &root.Command,
				},
				DefaultValue: false,
			}
			root.plainLongFlags["quiet"] = root.allFlags[SysMgmtGroup]["quiet"]
			root.plainShortFlags["q"] = root.allFlags[SysMgmtGroup]["quiet"]
		}
		if _, ok := root.allFlags[SysMgmtGroup]["debug"]; !ok {
			root.allFlags[SysMgmtGroup]["debug"] = &Flag{
				BaseOpt: BaseOpt{
					Short:       "D",
					Full:        "debug",
					Aliases:     []string{},
					Description: "Get into debug mode.",
					Hidden:      true,
					owner:       &root.Command,
				},
				DefaultValue: false,
			}
			root.plainLongFlags["debug"] = root.allFlags[SysMgmtGroup]["debug"]
			root.plainShortFlags["D"] = root.allFlags[SysMgmtGroup]["debug"]
		}
	}
}

func (w *ExecWorker) attachCmdrCommands(root *RootCommand) {
	if EnableCmdrCommands {
		if _, ok := root.allFlags[SysMgmtGroup]["strict-mode"]; !ok {
			root.allFlags[SysMgmtGroup]["strict-mode"] = &Flag{
				BaseOpt: BaseOpt{
					Full:        "strict-mode",
					Description: "strict mode for `cmdr`.",
					Hidden:      true,
					owner:       &root.Command,
				},
				DefaultValue: false,
			}
			root.plainLongFlags["strict-mode"] = root.allFlags[SysMgmtGroup]["strict-mode"]
		}
		if _, ok := root.allFlags[SysMgmtGroup]["no-env-overrides"]; !ok {
			root.allFlags[SysMgmtGroup]["no-env-overrides"] = &Flag{
				BaseOpt: BaseOpt{
					Full:        "no-env-overrides",
					Description: "No env var overrrides for `cmdr`.",
					Hidden:      true,
					owner:       &root.Command,
				},
				DefaultValue: false,
			}
			root.plainLongFlags["no-env-overrides"] = root.allFlags[SysMgmtGroup]["no-env-overrides"]
		}
		if _, ok := root.allFlags[SysMgmtGroup]["no-color"]; !ok {
			root.allFlags[SysMgmtGroup]["no-color"] = &Flag{
				BaseOpt: BaseOpt{
					Full:        "no-color",
					Description: "No color output for `cmdr`.",
					Hidden:      true,
					owner:       &root.Command,
				},
				DefaultValue: false,
			}
			root.plainLongFlags["no-color"] = root.allFlags[SysMgmtGroup]["no-color"]
		}
	}
}

func (w *ExecWorker) attachGeneratorsCommands(root *RootCommand) {
	if EnableGenerateCommands {
		found := false
		for _, sc := range root.SubCommands {
			if sc.Full == generatorCommands.Full {
				found = true
				return
			}
		}
		if !found {
			root.SubCommands = append(root.SubCommands, generatorCommands)
		}
	}
}

func (w *ExecWorker) forFlagNames(flg *Flag, cmd *Command, singleFlagNames, stringFlagNames map[string]bool) {
	if len(flg.Short) != 0 {
		if _, ok := singleFlagNames[flg.Short]; ok {
			ferr("flag char '%v' was been used. (command: %v)", flg.Short, w.backtraceCmdNames(cmd))
		} else {
			singleFlagNames[flg.Short] = true
		}
	}
	if len(flg.Full) != 0 {
		if _, ok := stringFlagNames[flg.Full]; ok {
			ferr("flag '%v' was been used. (command: %v)", flg.Full, w.backtraceCmdNames(cmd))
		} else {
			stringFlagNames[flg.Full] = true
		}
	}
	if len(flg.Short) == 0 && len(flg.Full) == 0 && len(flg.Name) != 0 {
		if _, ok := stringFlagNames[flg.Name]; ok {
			ferr("flag '%v' was been used. (command: %v)", flg.Name, w.backtraceCmdNames(cmd))
		} else {
			stringFlagNames[flg.Name] = true
		}
	}
}

func (w *ExecWorker) buildCrossRefsForFlag(flg *Flag, cmd *Command, singleFlagNames, stringFlagNames map[string]bool) {
	w.forFlagNames(flg, cmd, singleFlagNames, stringFlagNames)

	for _, sz := range flg.Aliases {
		if _, ok := stringFlagNames[sz]; ok {
			ferr("flag alias name '%v' was been used. (command: %v)", sz, w.backtraceCmdNames(cmd))
		} else {
			stringFlagNames[sz] = true
		}
	}
	if len(flg.Group) == 0 {
		flg.Group = UnsortedGroup
	}
	if _, ok := cmd.allFlags[flg.Group]; !ok {
		cmd.allFlags[flg.Group] = make(map[string]*Flag)
	}
	for _, sz := range flg.GetShortTitleNamesArray() {
		cmd.plainShortFlags[sz] = flg
	}
	for _, sz := range flg.GetLongTitleNamesArray() {
		cmd.plainLongFlags[sz] = flg
	}
	if flg.HeadLike {
		cmd.headLikeFlag = flg
	}
	cmd.allFlags[flg.Group][flg.GetTitleName()] = flg
}

func (w *ExecWorker) forCommandNames(cx, cmd *Command, singleCmdNames, stringCmdNames map[string]bool) {
	if len(cx.Short) != 0 {
		if _, ok := singleCmdNames[cx.Short]; ok {
			ferr("command char '%v' was been used. (command: %v)", cx.Short, w.backtraceCmdNames(cmd))
		} else {
			singleCmdNames[cx.Short] = true
		}
	}
	if len(cx.Full) != 0 {
		if _, ok := stringCmdNames[cx.Full]; ok {
			ferr("command '%v' was been used. (command: %v)", cx.Full, w.backtraceCmdNames(cmd))
		} else {
			stringCmdNames[cx.Full] = true
		}
	}
	if len(cx.Short) == 0 && len(cx.Full) == 0 && len(cx.Name) != 0 {
		if _, ok := stringCmdNames[cx.Name]; ok {
			ferr("command '%v' was been used. (command: %v)", cx.Name, w.backtraceCmdNames(cmd))
		} else {
			stringCmdNames[cx.Name] = true
		}
		cmd.plainCmds[cx.Name] = cx
	}
}

func (w *ExecWorker) buildCrossRefsForCommand(cx, cmd *Command, singleCmdNames, stringCmdNames map[string]bool) {
	w.forCommandNames(cx, cmd, singleCmdNames, stringCmdNames)

	for _, sz := range cx.Aliases {
		if len(sz) != 0 {
			if _, ok := stringCmdNames[sz]; ok {
				ferr("command alias name '%v' was been used. (command: %v)", sz, w.backtraceCmdNames(cmd))
			} else {
				stringCmdNames[sz] = true
			}
		}
	}

	if len(cx.Group) == 0 {
		cx.Group = UnsortedGroup
	}
	if _, ok := cmd.allCmds[cx.Group]; !ok {
		cmd.allCmds[cx.Group] = make(map[string]*Command)
	}
	for _, sz := range cx.GetTitleNamesArray() {
		cmd.plainCmds[sz] = cx
	}
	cmd.allCmds[cx.Group][cx.GetTitleName()] = cx
}

func (w *ExecWorker) buildCrossRefs(cmd *Command) {
	w.ensureCmdMembers(cmd)

	singleFlagNames := make(map[string]bool)
	stringFlagNames := make(map[string]bool)
	singleCmdNames := make(map[string]bool)
	stringCmdNames := make(map[string]bool)

	for _, flg := range cmd.Flags {
		flg.owner = cmd

		w.buildCrossRefsForFlag(flg, cmd, singleFlagNames, stringFlagNames)

		// opt.Children[flg.Full] = &OptOne{Value: flg.DefaultValue,}
		w.rxxtOptions.Set(w.backtraceFlagNames(flg), flg.DefaultValue)
	}

	for _, cx := range cmd.SubCommands {
		cx.owner = cmd

		w.buildCrossRefsForCommand(cx, cmd, singleCmdNames, stringCmdNames)
		// opt.Children[cx.Full] = newOpt()

		w.rxxtOptions.Set(w.backtraceCmdNames(cx), nil)
		// buildCrossRefs(cx, opt.Children[cx.Full])
		w.buildCrossRefs(cx)
	}

}

func (w *ExecWorker) backtraceFlagNames(flg *Flag) (str string) {
	var a []string
	a = append(a, flg.Full)
	for p := flg.owner; p != nil && p.owner != nil; {
		a = append(a, p.Full)
		p = p.owner
	}

	// reverse it
	i := 0
	j := len(a) - 1
	for i < j {
		a[i], a[j] = a[j], a[i]
		i++
		j--
	}

	str = strings.Join(a, ".")
	return
}

func (w *ExecWorker) backtraceCmdNames(cmd *Command) (str string) {
	var a []string
	a = append(a, cmd.GetTitleName())
	for p := cmd.owner; p != nil && p.owner != nil; {
		a = append(a, p.GetTitleName())
		p = p.owner
	}

	// reverse it
	i := 0
	j := len(a) - 1
	for i < j {
		a[i], a[j] = a[j], a[i]
		i++
		j--
	}

	str = strings.Join(a, ".")
	return
}

func (w *ExecWorker) ensureCmdMembers(cmd *Command) *Command {
	if cmd.allFlags == nil {
		cmd.allFlags = make(map[string]map[string]*Flag)
		cmd.allFlags[UnsortedGroup] = make(map[string]*Flag)
		cmd.allFlags[SysMgmtGroup] = make(map[string]*Flag)
	}

	if cmd.allCmds == nil {
		cmd.allCmds = make(map[string]map[string]*Command)
		cmd.allCmds[UnsortedGroup] = make(map[string]*Command)
		cmd.allCmds[SysMgmtGroup] = make(map[string]*Command)
	}

	if cmd.plainCmds == nil {
		cmd.plainCmds = make(map[string]*Command)
	}

	if cmd.plainLongFlags == nil {
		cmd.plainLongFlags = make(map[string]*Flag)
	}

	if cmd.plainShortFlags == nil {
		cmd.plainShortFlags = make(map[string]*Flag)
	}

	if cmd.root == nil {
		cmd.root = w.rootCommand
	}
	return cmd
}
