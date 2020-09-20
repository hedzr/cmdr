/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"fmt"
	"github.com/hedzr/cmdr/conf"
	"github.com/hedzr/log/exec"
	"os"
	"path"
	"regexp"
	"strings"
)

// AddOnBeforeXrefBuilding add hook func
// daemon plugin needed
func (w *ExecWorker) AddOnBeforeXrefBuilding(cb HookFunc) {
	if cb != nil {
		w.beforeXrefBuilding = append(w.beforeXrefBuilding, cb)
	}
}

// AddOnAfterXrefBuilt add hook func
// daemon plugin needed
func (w *ExecWorker) AddOnAfterXrefBuilt(cb HookFunc) {
	if cb != nil {
		w.afterXrefBuilt = append(w.afterXrefBuilt, cb)
	}
}

func (w *ExecWorker) setupFromEnvvarMap() {
	for k, v := range w.envVarToValueMap {
		_ = os.Setenv(k, v())
	}
}

func (w *ExecWorker) buildXref(rootCmd *RootCommand) (err error) {
	flog("--> preprocess / buildXref")

	// build xref for root command and its all sub-commands and flags
	// and build the default values
	w.buildRootCrossRefs(rootCmd)
	w.buildAddonsCrossRefs(rootCmd)
	w.buildExtensionsCrossRefs(rootCmd)

	w.setupFromEnvvarMap()

	if !w.doNotLoadingConfigFiles {
		// flog("--> buildXref: loadFromPredefinedLocation()")

		// pre-detects for `--config xxx`, `--config=xxx`, `--configxxx`
		//if err = w.parsePredefinedLocation(); err != nil {
		//	return
		//}
		_ = w.parsePredefinedLocation()

		// and now, loading the external configuration files
		err = w.loadFromPredefinedLocation(rootCmd)

		// if len(w.envPrefixes) > 0 {
		// 	EnvPrefix = w.envPrefixes
		// }
		// w.envPrefixes = EnvPrefix
		var envPrefix []string
		eps := GetString("env-prefix", "")
		if eps != "" && strings.Trim(eps, "[]") == eps {
			envPrefix = strings.Split(eps, ".")
		} else {
			envPrefix = GetStringSlice("env-prefix")
		}
		if len(envPrefix) > 0 {
			w.envPrefixes = envPrefix
			flog("--> preprocess / buildXref: env-prefix %v loaded", envPrefix)
		}
	}
	return
}

//goland:noinspection GoUnusedParameter
func (w *ExecWorker) buildAddonsCrossRefs(root *RootCommand) {
	flog("    - preprocess / buildXref / buildAddonsCrossRefs...")
}

//goland:noinspection ALL
func (w *ExecWorker) buildExtensionsCrossRefs(root *RootCommand) {
	flog("    - preprocess / buildXref / buildExtensionsCrossRefs...")
	// prefix := conf.AppName
	for _, dir := range w.extensionsLocations {
		dirExpanded := os.ExpandEnv(dir)
		// Logger.Debugf("      -> ext.dir: %v", dirExpanded)
		if exec.FileExists(dirExpanded) {
			err := exec.ForDirMax(dirExpanded, 0, 1, func(depth int, cwd string, fi os.FileInfo) (stop bool, err error) {
				if fi.IsDir() {
					return
				}
				var ok bool // = strings.HasPrefix(fi.Name(), prefix)
				ok = true
				// Logger.Debugf("      -> ext.dir: %v, file: %v", dirExpanded, fi.Name())
				if ok && fi.Mode().IsRegular() && exec.IsExecAny(fi.Mode()) {
					//name := fi.Name()[:len(prefix)]
					name := fi.Name()
					exe := path.Join(cwd, fi.Name())
					//if strings.HasPrefix(name, "-") || strings.HasPrefix(name, "_") {
					//	name = name[1:]
					//	Logger.Debugf("      -> ext.dir: %v, file: %v", dirExpanded, fi.Name())
					w._addAsSubCmd(root, name, exe)
					//}
				}
				return
			})
			if err != nil {
				Logger.Warnf("  warn - error in buildExtensionsCrossRefs.ForDir(): %v", err)
			}
		}
	}
}

func (w *ExecWorker) _addAsSubCmd(root *RootCommand, cmdName, cmdPath string) {
	var desc string
	desc = fmt.Sprintf("execute %q", cmdPath)
	if _, ok := root.allCmds[ExtGroup]; !ok {
		root.allCmds[ExtGroup] = make(map[string]*Command)
	}
	if _, ok := root.allCmds[ExtGroup][cmdName]; !ok {
		cx := &Command{
			BaseOpt: BaseOpt{
				Full:        cmdName,
				Short:       cmdName,
				Description: desc,
				Action: func(cmd *Command, args []string) (err error) {
					var out string
					_, out, err = exec.RunWithOutput(cmdPath)
					fmt.Print(out)
					return
				},
				Hidden: false,
				Group:  ExtGroup,
				owner:  &root.Command,
			},
		}
		root.SubCommands = uniAddCmd(root.SubCommands, cx)
		root.allCmds[ExtGroup][cmdName] = cx
		root.plainCmds[cmdName] = cx
	}

}

//goland:noinspection GoUnusedParameter
func (w *ExecWorker) buildRootCrossRefs(root *RootCommand) {
	flog("    - preprocess / buildXref / buildRootCrossRefs...")

	// initializes the internal variables/members
	w.ensureCmdMembers(&root.Command)

	// conf.AppName = root.AppName
	// conf.Version = root.Version
	// if len(conf.Buildstamp) == 0 {
	// 	conf.Buildstamp = time.Now().Format(time.RFC1123)
	// }

	w.attachVersionCommands(root)
	w.attachHelpCommands(root)
	w.attachVerboseCommands(root)
	w.attachGeneratorsCommands(root)
	w.attachCmdrCommands(root)

	w._buildCrossRefs(&root.Command)
}

func (w *ExecWorker) attachVersionCommands(root *RootCommand) {
	if w.enableVersionCommands {
		if _, ok := root.allCmds[SysMgmtGroup]["version"]; !ok {
			cx := &Command{
				BaseOpt: BaseOpt{
					Full:        "version",
					Aliases:     []string{"ver", "versions"},
					Description: "Show the version of this app.",
					Action: func(cmd *Command, args []string) (err error) {
						w.showVersion()
						return ErrShouldBeStopException
					},
					Hidden: true,
					Group:  SysMgmtGroup,
					owner:  &root.Command,
				},
			}
			root.SubCommands = uniAddCmd(root.SubCommands, cx)
			root.allCmds[SysMgmtGroup]["version"] = cx
			root.allCmds[SysMgmtGroup]["versions"] = cx
			root.allCmds[SysMgmtGroup]["ver"] = cx
			root.plainCmds["version"] = cx
			root.plainCmds["versions"] = cx
			root.plainCmds["ver"] = cx
		}
		if _, ok := root.allFlags[SysMgmtGroup]["version"]; !ok {
			ff := &Flag{
				BaseOpt: BaseOpt{
					Short:       "V",
					Full:        "version",
					Aliases:     []string{"ver", "versions"},
					Description: "Show the version of this app.",
					Hidden:      true,
					Group:       SysMgmtGroup,
					owner:       &root.Command,
					Action: func(cmd *Command, args []string) (err error) {
						w.showVersion()
						return ErrShouldBeStopException
					},
				},
				DefaultValue: false,
			}
			root.Flags = append(root.Flags, ff)
			root.allFlags[SysMgmtGroup]["version"] = ff
			root.plainLongFlags["version"] = ff
			root.plainLongFlags["versions"] = ff
			root.plainLongFlags["ver"] = ff
			root.plainShortFlags["V"] = ff
		}
		if _, ok := root.allFlags[SysMgmtGroup]["version-sim"]; !ok {
			ff := &Flag{
				BaseOpt: BaseOpt{
					Full:        "version-sim",
					Aliases:     []string{"version-simulate"},
					Description: "Simulate a faked version number for this app.",
					Hidden:      true,
					Group:       SysMgmtGroup,
					owner:       &root.Command,
					Action: func(cmd *Command, args []string) (err error) {
						conf.Version = GetStringR("version-sim")
						Set("version", conf.Version) // set into option 'app.version' too.
						return
					},
				},
				DefaultValue: "",
			}
			root.Flags = append(root.Flags, ff)
			root.allFlags[SysMgmtGroup]["version-sim"] = ff
			root.plainLongFlags["version-sim"] = ff
			root.plainLongFlags["version-simulate"] = ff
		}
		if _, ok := root.allFlags[SysMgmtGroup]["build-info"]; !ok {
			root.allFlags[SysMgmtGroup]["build-info"] = &Flag{
				BaseOpt: BaseOpt{
					Full:        "#",
					Aliases:     []string{},
					Description: "Show the building information of this app.",
					Hidden:      true,
					Group:       SysMgmtGroup,
					owner:       &root.Command,
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
	if w.enableHelpCommands {
		if _, ok := root.allFlags[SysMgmtGroup]["help"]; !ok {
			ff := &Flag{
				BaseOpt: BaseOpt{
					Short:       "h",
					Full:        "help",
					Aliases:     []string{"?", "helpme", "info", "usage"},
					Description: "Show this help screen",
					Hidden:      true,
					Group:       SysMgmtGroup,
					owner:       &root.Command,
					Action: func(cmd *Command, args []string) (err error) {
						// cmdr.Logger.Debugf("-- helpCommand hit. printHelp and stop.")
						// printHelp(cmd)
						// return ErrShouldBeStopException
						return nil
					},
				},
				DefaultValue: false,
				EnvVars:      []string{"HELP"},
			}
			root.Flags = append(root.Flags, ff)
			root.allFlags[SysMgmtGroup]["help"] = ff
			root.plainLongFlags["help"] = ff
			root.plainLongFlags["helpme"] = ff
			root.plainLongFlags["info"] = ff
			root.plainLongFlags["usage"] = ff
			root.plainShortFlags["h"] = ff
			root.plainShortFlags["?"] = ff

			ff = &Flag{
				BaseOpt: BaseOpt{
					Full:        "help-zsh",
					Description: "show help with zsh format, or others",
					Hidden:      true,
					Group:       SysMgmtGroup,
					owner:       &root.Command,
				},
				DefaultValue:            0,
				DefaultValuePlaceholder: "LEVEL",
			}
			root.Flags = append(root.Flags, ff)
			root.allFlags[SysMgmtGroup]["help-zsh"] = ff
			root.plainLongFlags["help-zsh"] = ff
			ff = &Flag{
				BaseOpt: BaseOpt{
					Full:        "help-bash",
					Description: "show help with bash format, or others",
					Hidden:      true,
					Group:       SysMgmtGroup,
					owner:       &root.Command,
				},
				DefaultValue: false,
			}
			root.Flags = append(root.Flags, ff)
			root.allFlags[SysMgmtGroup]["help-bash"] = ff
			root.plainLongFlags["help-bash"] = ff

			ff = &Flag{
				BaseOpt: BaseOpt{
					Full:        "tree",
					Description: "show a tree for all commands",
					Hidden:      true,
					Group:       SysMgmtGroup,
					owner:       &root.Command,
					Action:      dumpTreeForAllCommands,
				},
				DefaultValue: false,
			}
			root.Flags = append(root.Flags, ff)
			root.allFlags[SysMgmtGroup]["tree"] = ff
			root.plainLongFlags["tree"] = ff
		}
		if _, ok := root.allFlags[SysMgmtGroup]["config"]; !ok {
			ff := &Flag{
				BaseOpt: BaseOpt{
					Full:        "config",
					Aliases:     []string{},
					Description: "load config files from where you specified",
					Action: func(cmd *Command, args []string) (err error) {
						// cmdr.Logger.Debugf("-- --config hit. printHelp and stop.")
						// return ErrShouldBeStopException
						return nil
					},
					Group: SysMgmtGroup,
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
			root.Flags = append(root.Flags, ff)
			root.allFlags[SysMgmtGroup]["config"] = ff
			root.plainLongFlags["config"] = ff
		}
	}
}

func (w *ExecWorker) attachVerboseCommands(root *RootCommand) {
	if w.enableVerboseCommands {
		if _, ok := root.allFlags[SysMgmtGroup]["verbose"]; !ok {
			ff := &Flag{
				BaseOpt: BaseOpt{
					Short: "v",
					Full:  "verbose",
					// Aliases:     []string{"vv", "vvv"},
					Description: "Show this help screen",
					// Hidden:      true,
					Group: SysMgmtGroup,
					owner: &root.Command,
					// Action: func(cmd *Command, args []string) (err error) {
					// 	if f := FindFlag("verbose", cmd); f != nil {
					// 		f.times++
					// 		// fmt.Println("verbose++: ", f.times)
					// 	}
					// 	return
					// },

					// Action: func(cmd *Command, args []string) (err error) {
					// 	if f := FindFlag("verbose", cmd); f != nil {
					// 		fmt.Println("verbose++: ", f.times)
					// 	}
					// 	return
					// },
				},
				DefaultValue: false,
				EnvVars:      []string{"VERBOSE"},
			}
			root.Flags = append(root.Flags, ff)
			root.allFlags[SysMgmtGroup]["verbose"] = ff
			root.plainLongFlags["verbose"] = root.allFlags[SysMgmtGroup]["verbose"]
			// root.plainLongFlags["vvv"] = root.allFlags[SysMgmtGroup]["verbose"]
			// root.plainLongFlags["vv"] = root.allFlags[SysMgmtGroup]["verbose"]
			root.plainShortFlags["v"] = root.allFlags[SysMgmtGroup]["verbose"]
		}
		if _, ok := root.allFlags[SysMgmtGroup]["quiet"]; !ok {
			ff := &Flag{
				BaseOpt: BaseOpt{
					Short:       "q",
					Full:        "quiet",
					Aliases:     []string{},
					Description: "No more screen output.",
					// Hidden:      true,
					Group: SysMgmtGroup,
					owner: &root.Command,
				},
				DefaultValue: false,
				EnvVars:      []string{"QUITE"},
			}
			root.Flags = append(root.Flags, ff)
			root.allFlags[SysMgmtGroup]["quiet"] = ff
			root.plainLongFlags["quiet"] = root.allFlags[SysMgmtGroup]["quiet"]
			root.plainShortFlags["q"] = root.allFlags[SysMgmtGroup]["quiet"]
		}
		if _, ok := root.allFlags[SysMgmtGroup]["debug"]; !ok {
			ff := &Flag{
				BaseOpt: BaseOpt{
					Short:       "D",
					Full:        "debug",
					Aliases:     []string{},
					Description: "Get into debug mode.",
					Hidden:      true,
					Group:       SysMgmtGroup,
					owner:       &root.Command,
				},
				DefaultValue: false,
				EnvVars:      []string{"DEBUG"},
			}
			root.Flags = append(root.Flags, ff)
			root.allFlags[SysMgmtGroup]["debug"] = ff
			root.plainLongFlags["debug"] = root.allFlags[SysMgmtGroup]["debug"]
			root.plainShortFlags["D"] = root.allFlags[SysMgmtGroup]["debug"]
		}
		if _, ok := root.allFlags[SysMgmtGroup]["env"]; !ok {
			ff := &Flag{
				BaseOpt: BaseOpt{
					Short:       "",
					Full:        "env",
					Aliases:     []string{},
					Description: "Dump environment info in `~~debug` mode.",
					Hidden:      true,
					Group:       SysMgmtGroup,
					owner:       &root.Command,
				},
				DefaultValue: false,
			}
			root.Flags = append(root.Flags, ff)
			root.allFlags[SysMgmtGroup]["env"] = ff
			root.plainLongFlags["env"] = root.allFlags[SysMgmtGroup]["env"]
		}
		if _, ok := root.allFlags[SysMgmtGroup]["raw"]; !ok {
			ff := &Flag{
				BaseOpt: BaseOpt{
					Short:       "",
					Full:        "raw",
					Aliases:     []string{},
					Description: "Dump the option value in raw mode (with golang data structure, without envvar expanding).",
					Hidden:      true,
					Group:       SysMgmtGroup,
					owner:       &root.Command,
				},
				DefaultValue: false,
			}
			root.Flags = append(root.Flags, ff)
			root.allFlags[SysMgmtGroup]["raw"] = ff
			root.plainLongFlags["raw"] = root.allFlags[SysMgmtGroup]["raw"]
		}
		if _, ok := root.allFlags[SysMgmtGroup]["value-type"]; !ok {
			ff := &Flag{
				BaseOpt: BaseOpt{
					Short:       "",
					Full:        "value-type",
					Aliases:     []string{},
					Description: "Dump the option value type.",
					Hidden:      true,
					Group:       SysMgmtGroup,
					owner:       &root.Command,
				},
				DefaultValue: false,
			}
			root.Flags = append(root.Flags, ff)
			root.allFlags[SysMgmtGroup]["value-type"] = ff
			root.plainLongFlags["value-type"] = root.allFlags[SysMgmtGroup]["value-type"]
		}
		if _, ok := root.allFlags[SysMgmtGroup]["more"]; !ok {
			ff := &Flag{
				BaseOpt: BaseOpt{
					Short:       "",
					Full:        "more",
					Aliases:     []string{},
					Description: "Dump more info in `~~debug` mode.",
					Hidden:      true,
					Group:       SysMgmtGroup,
					owner:       &root.Command,
				},
				DefaultValue: false,
			}
			root.Flags = append(root.Flags, ff)
			root.allFlags[SysMgmtGroup]["more"] = ff
			root.plainLongFlags["more"] = root.allFlags[SysMgmtGroup]["more"]
		}
	}
}

func (w *ExecWorker) attachCmdrCommands(root *RootCommand) {
	if w.enableCmdrCommands {
		if _, ok := root.allFlags[SysMgmtGroup]["strict-mode"]; !ok {
			ff := &Flag{
				BaseOpt: BaseOpt{
					Full:        "strict-mode",
					Description: "strict mode for `cmdr`.",
					Hidden:      true,
					Group:       SysMgmtGroup,
					owner:       &root.Command,
				},
				DefaultValue: false,
				EnvVars:      []string{"STRICT"},
			}
			root.Flags = append(root.Flags, ff)
			root.allFlags[SysMgmtGroup]["strict-mode"] = ff
			root.plainLongFlags["strict-mode"] = ff
		}
		if _, ok := root.allFlags[SysMgmtGroup]["no-env-overrides"]; !ok {
			ff := &Flag{
				BaseOpt: BaseOpt{
					Full:        "no-env-overrides",
					Description: "No env var overrides for `cmdr`.",
					Hidden:      true,
					Group:       SysMgmtGroup,
					owner:       &root.Command,
				},
				DefaultValue: false,
			}
			root.Flags = append(root.Flags, ff)
			root.allFlags[SysMgmtGroup]["no-env-overrides"] = ff
			root.plainLongFlags["no-env-overrides"] = ff
		}
		if _, ok := root.allFlags[SysMgmtGroup]["no-color"]; !ok {
			ff := &Flag{
				BaseOpt: BaseOpt{
					Full:        "no-color",
					Description: "No color output for `cmdr`.",
					Hidden:      true,
					Group:       SysMgmtGroup,
					owner:       &root.Command,
				},
				DefaultValue: false,
				EnvVars:      []string{"NOCOLOR", "NO_COLOR"},
			}
			root.Flags = append(root.Flags, ff)
			root.allFlags[SysMgmtGroup]["no-color"] = ff
			root.plainLongFlags["no-color"] = ff
		}
	}
}

func (w *ExecWorker) attachGeneratorsCommands(root *RootCommand) {
	if w.enableGenerateCommands {
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

func (w *ExecWorker) _buildCrossRefs(cmd *Command) {
	w.ensureCmdMembers(cmd)

	singleFlagNames := make(map[string]bool)
	stringFlagNames := make(map[string]bool)
	singleCmdNames := make(map[string]bool)
	stringCmdNames := make(map[string]bool)
	tgs := make(map[string]bool)

	for _, flg := range cmd.Flags {
		flg.owner = cmd

		if len(flg.ToggleGroup) > 0 {
			if len(flg.Group) == 0 {
				flg.Group = flg.ToggleGroup
			}
			tgs[flg.ToggleGroup] = true
		}

		if b := regexp.MustCompile("`(.+)`").Find([]byte(flg.Description)); len(flg.DefaultValuePlaceholder) == 0 && len(b) > 2 {
			ph := strings.ToUpper(strings.Trim(string(b), "`"))
			flg.DefaultValuePlaceholder = ph
		}

		w._buildCrossRefsForFlag(flg, cmd, singleFlagNames, stringFlagNames)

		// opt.Children[flg.Full] = &OptOne{Value: flg.DefaultValue,}
		w.rxxtOptions.Set(w.backtraceFlagNames(flg), flg.DefaultValue)
	}

	for _, cx := range cmd.SubCommands {
		cx.owner = cmd

		w._buildCrossRefsForCommand(cx, cmd, singleCmdNames, stringCmdNames)
		// opt.Children[cx.Full] = newOpt()

		w.rxxtOptions.Set(w.backtraceCmdNames(cx), nil)
		// buildCrossRefs(cx, opt.Children[cx.Full])
		w._buildCrossRefs(cx)
	}

	for tg := range tgs {
		w.buildToggleGroup(tg, cmd)
	}
}

func (w *ExecWorker) _buildCrossRefsForFlag(flg *Flag, cmd *Command, singleFlagNames, stringFlagNames map[string]bool) {
	w.forFlagNames(flg, cmd, singleFlagNames, stringFlagNames)

	for _, sz := range flg.Aliases {
		if _, ok := stringFlagNames[sz]; ok {
			ferr("\nNOTE: flag alias name '%v' has been used. (command: %v)", sz, w.backtraceCmdNames(cmd))
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

func (w *ExecWorker) forFlagNames(flg *Flag, cmd *Command, singleFlagNames, stringFlagNames map[string]bool) {
	if len(flg.Short) != 0 {
		if _, ok := singleFlagNames[flg.Short]; ok {
			ferr("\nNOTE: flag char '%v' has been used. (command: %v)", flg.Short, w.backtraceCmdNames(cmd))
		} else {
			singleFlagNames[flg.Short] = true
		}
	}
	if len(flg.Full) != 0 {
		if _, ok := stringFlagNames[flg.Full]; ok {
			ferr("\nNOTE: flag '%v' has been used. (command: %v)", flg.Full, w.backtraceCmdNames(cmd))
		} else {
			stringFlagNames[flg.Full] = true
		}
	}
	if len(flg.Short) == 0 && len(flg.Full) == 0 && len(flg.Name) != 0 {
		if _, ok := stringFlagNames[flg.Name]; ok {
			ferr("\nNOTE: flag '%v' has been used. (command: %v)", flg.Name, w.backtraceCmdNames(cmd))
		} else {
			stringFlagNames[flg.Name] = true
		}
	}
}

func (w *ExecWorker) _buildCrossRefsForCommand(cx, cmd *Command, singleCmdNames, stringCmdNames map[string]bool) {
	w.forCommandNames(cx, cmd, singleCmdNames, stringCmdNames)

	for _, sz := range cx.Aliases {
		if len(sz) != 0 {
			if _, ok := stringCmdNames[sz]; ok {
				ferr("\nNOTE: command alias name '%v' has been used. (command: %v)", sz, w.backtraceCmdNames(cmd))
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

func (w *ExecWorker) forCommandNames(cx, cmd *Command, singleCmdNames, stringCmdNames map[string]bool) {
	if len(cx.Short) != 0 {
		if _, ok := singleCmdNames[cx.Short]; ok {
			ferr("\nNOTE: command char '%v' has been used. (command: %v)", cx.Short, w.backtraceCmdNames(cmd))
		} else {
			singleCmdNames[cx.Short] = true
		}
	}
	if len(cx.Full) != 0 {
		if _, ok := stringCmdNames[cx.Full]; ok {
			ferr("\nNOTE: command '%v' has been used. (command: %v)", cx.Full, w.backtraceCmdNames(cmd))
		} else {
			stringCmdNames[cx.Full] = true
		}
	}
	if len(cx.Short) == 0 && len(cx.Full) == 0 && len(cx.Name) != 0 {
		if _, ok := stringCmdNames[cx.Name]; ok {
			ferr("\nNOTE: command '%v' has been used. (command: %v)", cx.Name, w.backtraceCmdNames(cmd))
		} else {
			stringCmdNames[cx.Name] = true
		}
		cmd.plainCmds[cx.Name] = cx
	}
}

func (w *ExecWorker) buildToggleGroup(tg string, cmd *Command) {
	for _, f := range cmd.Flags {
		if tg == f.ToggleGroup && f.DefaultValue == true {
			w.rxxtOptions.Set(w.backtraceFlagNames(f), true)
			w.rxxtOptions.Set(w.backtraceCmdNames(cmd)+"."+tg, f.Full)
			break
		}
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
