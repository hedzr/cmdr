/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"strings"
)

func buildRootCrossRefs(root *RootCommand) {
	ensureCmdMembers(&root.Command)

	if EnableVersionCommands {
		if _, ok := root.allCmds[SysMgmtGroup]["version"]; !ok {
			root.allCmds[SysMgmtGroup]["version"] = &Command{
				BaseOpt: BaseOpt{
					Full:        "version",
					Aliases:     []string{"ver"},
					Description: "Show the version of this app.",
					Action: func(cmd *Command, args []string) (err error) {
						showVersion()
						return ErrShouldBeStopException
					},
				},
			}
		}
		if _, ok := root.allFlags[SysMgmtGroup]["version"]; !ok {
			root.allFlags[SysMgmtGroup]["version"] = &Flag{
				BaseOpt: BaseOpt{
					Short:       "V",
					Full:        "version",
					Aliases:     []string{"vv"},
					Description: "Show the version of this app.",
					// Hidden:      true,
					Action: func(cmd *Command, args []string) (err error) {
						showVersion()
						return ErrShouldBeStopException
					},
				},
				DefaultValue: false,
			}
			root.plainFlags["version"] = root.allFlags[SysMgmtGroup]["version"]
			root.plainFlags["v"] = root.allFlags[SysMgmtGroup]["version"]
		}
		if _, ok := root.allFlags[SysMgmtGroup]["build-info"]; !ok {
			root.allFlags[SysMgmtGroup]["build-info"] = &Flag{
				BaseOpt: BaseOpt{
					Full:        "#",
					Aliases:     []string{},
					Description: "Show the building information of this app.",
					Hidden:      true,
					Action: func(cmd *Command, args []string) (err error) {
						showVersion()
						return ErrShouldBeStopException
					},
				},
				DefaultValue: false,
			}
			root.plainFlags["#"] = root.allFlags[SysMgmtGroup]["build-info"]
			root.plainFlags["build-info"] = root.allFlags[SysMgmtGroup]["build-info"]
		}
	}
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
				},
				DefaultValue: false,
			}
			root.plainFlags["help"] = root.allFlags[SysMgmtGroup]["help"]
			root.plainFlags["helpme"] = root.allFlags[SysMgmtGroup]["help"]
			root.plainFlags["h"] = root.allFlags[SysMgmtGroup]["help"]
			root.plainFlags["?"] = root.allFlags[SysMgmtGroup]["help"]
			root.plainFlags["info"] = root.allFlags[SysMgmtGroup]["help"]
			root.plainFlags["usage"] = root.allFlags[SysMgmtGroup]["help"]

			root.allFlags[SysMgmtGroup]["help-zsh"] = &Flag{
				BaseOpt: BaseOpt{
					Full:                    "help-zsh",
					Description:             "show help with zsh format, or others",
					Hidden:                  true,
					DefaultValuePlaceholder: "LEVEL",
				},
				DefaultValue: 0,
			}
			root.allFlags[SysMgmtGroup]["help-bash"] = &Flag{
				BaseOpt: BaseOpt{
					Full:        "help-bash",
					Description: "show help with bash format, or others",
					Hidden:      true,
				},
				DefaultValue: false,
			}
			root.plainFlags["help-zsh"] = root.allFlags[SysMgmtGroup]["help-zsh"]
			root.plainFlags["help-bash"] = root.allFlags[SysMgmtGroup]["help-bash"]
		}
		if _, ok := root.allFlags[SysMgmtGroup]["config"]; !ok {
			root.allFlags[SysMgmtGroup]["config"] = &Flag{
				BaseOpt: BaseOpt{
					Full:        "config",
					Aliases:     []string{},
					Description: "load config files from where you specified",
					Action: func(cmd *Command, args []string) (err error) {
						// logrus.Debugf("-- helpCommand hit. printHelp and stop.")
						// return ErrShouldBeStopException
						return nil
					},
					DefaultValuePlaceholder: "[Location of config file]",
				},
				DefaultValue: nil,
			}
			root.plainFlags["help"] = root.allFlags[SysMgmtGroup]["help"]
			root.plainFlags["helpme"] = root.allFlags[SysMgmtGroup]["help"]
			root.plainFlags["h"] = root.allFlags[SysMgmtGroup]["help"]
			root.plainFlags["?"] = root.allFlags[SysMgmtGroup]["help"]
			root.plainFlags["info"] = root.allFlags[SysMgmtGroup]["help"]
			root.plainFlags["usage"] = root.allFlags[SysMgmtGroup]["help"]
		}
	}
	if EnableVerboseCommands {
		if _, ok := root.allFlags[SysMgmtGroup]["verbose"]; !ok {
			root.allFlags[SysMgmtGroup]["verbose"] = &Flag{
				BaseOpt: BaseOpt{
					Short:       "v",
					Full:        "verbose",
					Aliases:     []string{},
					Description: "Show this help screen",
					// Hidden:      true,
				},
				DefaultValue: false,
			}
			root.plainFlags["verbose"] = root.allFlags[SysMgmtGroup]["verbose"]
			root.plainFlags["v"] = root.allFlags[SysMgmtGroup]["verbose"]
		}
		if _, ok := root.allFlags[SysMgmtGroup]["quiet"]; !ok {
			root.allFlags[SysMgmtGroup]["quiet"] = &Flag{
				BaseOpt: BaseOpt{
					Short:       "q",
					Full:        "quiet",
					Aliases:     []string{},
					Description: "No more screen output.",
					// Hidden:      true,
				},
				DefaultValue: false,
			}
			root.plainFlags["quiet"] = root.allFlags[SysMgmtGroup]["quiet"]
			root.plainFlags["q"] = root.allFlags[SysMgmtGroup]["quiet"]
		}
		if _, ok := root.allFlags[SysMgmtGroup]["debug"]; !ok {
			root.allFlags[SysMgmtGroup]["debug"] = &Flag{
				BaseOpt: BaseOpt{
					Short:       "D",
					Full:        "debug",
					Aliases:     []string{},
					Description: "Get into debug mode.",
					Hidden:      true,
				},
				DefaultValue: false,
			}
			root.plainFlags["debug"] = root.allFlags[SysMgmtGroup]["debug"]
			root.plainFlags["D"] = root.allFlags[SysMgmtGroup]["debug"]
		}
	}

	if EnableGenerateCommands {
		root.SubCommands = append(root.SubCommands, generatorCommands)
	}

	// rootOptions = newOpt()
	// buildCrossRefs(&root.Command, rootOptions)
	buildCrossRefs(&root.Command)
}

// func newOpt() *OptOne {
// 	return &OptOne{
// 		Children: make(map[string]*OptOne),
// 	}
// }

func newCmd() *Command {
	return ensureCmdMembers(&Command{})
}

func ensureCmdMembers(cmd *Command) *Command {
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

	if cmd.plainFlags == nil {
		cmd.plainFlags = make(map[string]*Flag)
	}

	if cmd.root == nil {
		cmd.root = rootCommand
	}
	return cmd
}

func buildCrossRefs(cmd *Command) {
	ensureCmdMembers(cmd)

	singleFlagNames := make(map[string]bool)
	stringFlagNames := make(map[string]bool)
	singleCmdNames := make(map[string]bool)
	stringCmdNames := make(map[string]bool)

	for _, flg := range cmd.Flags {
		flg.owner = cmd
		if len(flg.Short) != 0 {
			if _, ok := singleFlagNames[flg.Short]; ok {
				ferr("flag char '%c' was been used.", flg.Short)
			} else {
				singleFlagNames[flg.Short] = true
			}
		}
		if len(flg.Full) != 0 {
			if _, ok := stringFlagNames[flg.Full]; ok {
				ferr("flag '%s' was been used.", flg.Full)
			} else {
				stringFlagNames[flg.Full] = true
			}
		}
		for _, sz := range flg.Aliases {
			if _, ok := stringFlagNames[sz]; ok {
				ferr("flag alias name '%s' was been used.", sz)
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
		for _, sz := range flg.GetTitleNamesArray() {
			cmd.plainFlags[sz] = flg
		}
		cmd.allFlags[flg.Group][flg.GetTitleName()] = flg

		// opt.Children[flg.Full] = &OptOne{Value: flg.DefaultValue,}
		rxxtOptions.Set(backtraceFlagNames(flg), flg.DefaultValue)
	}

	for _, cx := range cmd.SubCommands {
		cx.owner = cmd

		if _, ok := singleCmdNames[cx.Short]; ok {
			ferr("command char '%c' was been used.", cx.Short)
		} else {
			singleCmdNames[cx.Short] = true
		}
		if _, ok := stringCmdNames[cx.Full]; ok {
			ferr("command '%s' was been used.", cx.Full)
		} else {
			stringCmdNames[cx.Full] = true
		}
		for _, sz := range cx.Aliases {
			if _, ok := stringCmdNames[sz]; ok {
				ferr("command alias name '%s' was been used.", sz)
			} else {
				stringCmdNames[sz] = true
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

		// opt.Children[cx.Full] = newOpt()

		rxxtOptions.Set(backtraceCmdNames(cx), nil)
		// buildCrossRefs(cx, opt.Children[cx.Full])
		buildCrossRefs(cx)
	}

}

func backtraceFlagNames(flg *Flag) (str string) {
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

func backtraceCmdNames(cmd *Command) (str string) {
	var a []string
	a = append(a, cmd.Full)
	for p := cmd.owner; p != nil && p.owner != nil; {
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
