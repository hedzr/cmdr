/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"fmt"
	"os"
	"path"
	"plugin"
	"regexp"
	"runtime"
	"strings"

	cmdrbase "github.com/hedzr/cmdr-base"
	"github.com/hedzr/cmdr/conf"
	"github.com/hedzr/log/closers"
	"github.com/hedzr/log/dir"
	"github.com/hedzr/log/exec"
	"gopkg.in/hedzr/errors.v3"
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

func (w *ExecWorker) buildXref(rootCmd *RootCommand, args []string) (err error) {
	flog("--> preprocess / buildXref")

	if rootCmd != nil {
		// build xref for root command and its all sub-commands and flags
		// and build the default values
		w.buildRootCrossRefs(rootCmd)
		w.buildAddonsCrossRefs(rootCmd)
		w.buildExtensionsCrossRefs(rootCmd)

		w.setupFromEnvvarMap()
	}

	flog("--> before-config-file-loading")
	for _, x := range w.beforeConfigFileLoading {
		if x != nil {
			x(rootCmd, args)
		}
	}

	if rootCmd != nil && !w.doNotLoadingConfigFiles {
		// flog("--> buildXref: loadFromPredefinedLocations()")

		// pre-detects for `--config xxx`, `--config=xxx`, `--configxxx`
		// if err = w.parsePredefinedLocation(); err != nil {
		//	return
		// }
		_ = w.parsePredefinedLocation()

		w.rxxtOptions.setToAppendMode()

		// and now, loading the external configuration files
		_ = w.loadFromPredefinedLocations(rootCmd)
		_ = w.loadFromSecondaryLocations(rootCmd)
		_ = w.loadFromAlterLocations(rootCmd)

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

		w.buildAliasesCrossRefs(rootCmd)
	}

	flog("--> after-config-file-loading")
	for _, x := range w.afterConfigFileLoading {
		if x != nil {
			x(rootCmd, args)
		}
	}
	return
}

type aliasesCommands struct {
	Group    string
	Commands []*Command
}

// buildAliasesCrossRefs binds aliases into cmdr command
// system.
//
// The aliases are defined in config file at path 'app.aliases'.
// For more details, lookup the real sample in 'examples/fluent':
//
//	ci/etc/fluent/conf.d/91.cmd-aliases.yml
//
// A load-failure aliases state will be ignored by the mainly
// cmdr parsing process (but print an error message).
//
//goland:noinspection GoUnusedParameter
func (w *ExecWorker) buildAliasesCrossRefs(root *RootCommand) {
	var (
		aliases = new(aliasesCommands)
		err     error
	)
	if err = GetSectionFrom("aliases", &aliases); err == nil {
		err = w._addCommandsForAliasesGroup(root, aliases)
	}
	if err != nil {
		Logger.Errorf("buildAliasesCrossRefs error: %v", err)
	}
}

//goland:noinspection GoUnusedParameter
func (w *ExecWorker) _addCommandsForAliasesGroup(root *RootCommand, aliases *aliasesCommands) (err error) {
	flog("aliases:\n%v\n", aliases)
	if root == nil || aliases == nil {
		err = errors.New("bad param")
		return
	}

	if aliases.Group == "" {
		aliases.Group = AliasesGroup
	}

	for _, cmd := range aliases.Commands {
		w.ensureCmdMembers(cmd, root)
		err = w._toolAddCmd(&root.Command, aliases.Group, cmd)
	}
	w._buildCrossRefs(&root.Command, root)
	return
}

func (w *ExecWorker) _toolChkFlags(cc *Command) (err error) {
	// for _, f := range cc.Flags {
	//	if f.owner == nil {
	//		f.owner = cc
	//	}
	//	if f.DefaultValueType != "" && f.DefaultValue == nil {
	//		//var kind = reflect.String
	//		//for x := reflect.Bool; x <= reflect.Complex128; x++ {
	//		//	if x.String() == f.DefaultValueType {
	//		//		kind = x
	//		//		break
	//		//	}
	//		//}
	//
	//		//typ := reflect.Kind.String()
	//		//f.DefaultValue = reflect.New(typ)
	//
	//		Logger.Errorf("fatal error")
	//	}
	// }
	return
}

func (w *ExecWorker) _toolAddCmd(parent *Command, groupName string, cc *Command) (err error) {
	if cc.Group == "" {
		cc.Group = groupName
	} else {
		groupName = cc.Group
	}

	if _, ok := parent.allCmds[groupName]; !ok {
		parent.allCmds[groupName] = make(map[string]*Command)
	}

	err = w._toolChkFlags(cc)

	cmdName := cc.GetTitleName()
	if _, ok := parent.allCmds[groupName][cmdName]; !ok {
		if cc.Action == nil {
			w.bindInvokeToAction(cc)
		}
		if cc.Short != "" && cc.Short != cmdName {
			if _, ok := parent.plainCmds[cc.Short]; !ok {
				parent.plainCmds[cc.Short] = cc
			}
		}
		for _, n := range cc.Aliases {
			if _, ok := parent.plainCmds[n]; !ok {
				parent.plainCmds[n] = cc
			}
		}
		parent.SubCommands = uniAddCmd(parent.SubCommands, cc)
		parent.allCmds[groupName][cmdName] = cc
		parent.plainCmds[cmdName] = cc

		for _, c := range cc.SubCommands {
			if c.owner == nil {
				c.owner = cc
			}
			if len(c.SubCommands) > 0 {
				err = w._toolAddCmd(cc, groupName, c)
			} else if c.Action == nil {
				w.bindInvokeToAction(c)
			}
		}
	}
	return
}

func (w *ExecWorker) bindInvokeToAction(c *Command) {
	if c.Invoke != "" {
		cmdPathParts := strings.Split(c.Invoke, " ")
		if len(cmdPathParts) > 1 {
			c.Invoke, c.presetCmdLines = cmdPathParts[0], cmdPathParts[1:]
		}
		c.Action = w.getInvokeAction(c)
	}
	if c.InvokeProc != "" {
		c.Action = w.getInvokeProcAction(c)
	}
	if c.InvokeShell != "" {
		c.Action = w.getInvokeShellAction(c)
	}
}

func (w *ExecWorker) locateCommand(cmdPath string, from *Command) (cmd *Command, matched bool) {
	if from == nil {
		from = &w.rootCommand.Command
	}
	cmdPathParts := strings.Split(cmdPath, " ")
	if len(cmdPathParts) == 0 {
		return
	}
	parts := strings.Split(cmdPathParts[0], "/")
	for i, pp := range parts {
		if pp == "." {
			continue
		}
		if pp == ".." {
			from = from.GetOwner()
			continue
		}
		if pp == "" {
			if i == 0 {
				from = &w.rootCommand.Command
			}
			continue
		}
		if cmd, matched = from.plainCmds[pp]; matched {
			from = cmd
		}
	}
	return
}

func (w *ExecWorker) getInvokeAction(from *Command) Handler {
	return func(cmd *Command, args []string) (err error) {
		invoke := w.expandTmplWithExecutiveEnv(cmd.Invoke, cmd, args)
		if cx, matched := w.locateCommand(invoke, cmd); matched {
			if cx.Action != nil {
				err = cx.Action(cmd, args)
			}
		}
		return
	}
}

func (w *ExecWorker) getInvokeProcAction(from *Command) Handler {
	return func(cmd *Command, args []string) (err error) {
		invokeProc := w.expandTmplWithExecutiveEnv(cmd.InvokeProc, cmd, args)
		cmdParts := strings.Split(invokeProc, " ")
		c, args := cmdParts[0], cmdParts[1:]
		err = exec.Run(c, args...)
		return
	}
}

func (w *ExecWorker) getInvokeShellAction(from *Command) Handler {
	return func(cmd *Command, args []string) (err error) {
		// cmdParts := strings.Split(from.InvokeShell, " ")
		// c, args := cmdParts[0], cmdParts[1:]
		// err = exec.Run(c, args...)

		// NOTE: cmd == from

		// var a []string
		// shell := cmd.Shell
		// if shell == "" {
		//	if runtime.GOOS == "windows" {
		//		shell = "powershell.exe"
		//	} else {
		//		shell = "/bin/bash"
		//	}
		// } else if strings.Contains(shell, "/env ") {
		//	c := strings.Split(shell, " ")
		//	shell, a = c[0], append(a, c[1:]...)
		// }
		//
		// scriptFragments := w.expandTmplWithExecutiveEnv(cmd.InvokeShell, cmd, args)
		//
		// if strings.Contains(shell, "powershell") {
		//	a = append(a, "-NoProfile", "-NonInteractive", scriptFragments)
		// } else {
		//	a = append(a, "-c", scriptFragments)
		// }
		//
		// err = exec.Run(shell, a...)

		return exec.InvokeShellScripts(cmd.InvokeShell,
			exec.WithScriptExpander(func(source string) string {
				return w.expandTmplWithExecutiveEnv(cmd.InvokeShell, cmd, args)
			}),
			exec.WithScriptShell(cmd.Shell),
		)
	}
}

func (w *ExecWorker) expandTmplWithExecutiveEnv(source string, cmd *Command, args []string) (text string) {
	text = tplApply(source, struct { //nolint:govet //just a testcase
		Cmd        *Command
		Args       []string
		ArgsString string
		Store      *Options
	}{
		cmd,
		args,
		strings.Join(args, " "),
		w.rxxtOptions,
	})
	return
}

// buildAddonsCrossRefs for cmdr addons.
//
// A cmdr addon, which is a golang plugin, can be integrated
// into host-app better than an extension.
//
//goland:noinspection GoUnusedParameter
func (w *ExecWorker) buildAddonsCrossRefs(root *RootCommand) {
	// var cwd = dir.GetCurrentDir()
	// flog("    - preprocess / buildXref / buildAddonsCrossRefs...%q, %q", cwd, conf.AppName)
	flog("    - preprocess / buildXref / buildAddonsCrossRefs...")
	for _, d := range w.pluginsLocations {
		dirExpanded := os.ExpandEnv(d)
		// Logger.Debugf("      -> addons.dir: %v", dirExpanded)
		if dir.FileExists(dirExpanded) {
			err := dir.ForFileMax(dirExpanded, 0, 1, func(depth int, cwd string, fi os.FileInfo) (stop bool, err error) {
				if !fi.IsDir() {
					ok := true // = strings.HasPrefix(fi.Name(), prefix)
					// Logger.Debugf("      -> addons.dir: %v, file: %v", dirExpanded, fi.Name())
					if ok && dir.IsModeExecAny(fi.Mode()) {
						err = w._addOnAddIt(fi, root, cwd)
					}
				}
				return
			})
			if err != nil {
				Logger.Warnf("  warn - error in buildExtensionsCrossRefs.ForDir(): %v", err)
			}
		}
	}
}

func (w *ExecWorker) _addOnAddIt(fi os.FileInfo, root *RootCommand, cwd string) (err error) {
	if fi.Mode().IsRegular() {
		// name := fi.Name()[:len(prefix)]
		name := fi.Name()
		exe := path.Join(cwd, fi.Name())
		// if strings.HasPrefix(name, "-") || strings.HasPrefix(name, "_") {
		//	name = name[1:]
		//	Logger.Debugf("      -> addons.dir: %v, file: %v", dirExpanded, fi.Name())
		err = w._addonAsSubCmd(root, name, exe)
		// }
	}
	return
}

func (w *ExecWorker) _addonAsSubCmd(root *RootCommand, cmdName, cmdPath string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("wrapped error").WithData(e)
		}
	}()

	desc := fmt.Sprintf("execute %q", cmdPath)

	var p *plugin.Plugin
	p, err = plugin.Open(cmdPath)
	if err != nil {
		return
	}

	var newAddonSymbol plugin.Symbol
	newAddonSymbol, err = p.Lookup("NewAddon")
	if err != nil {
		return
	}

	var newAddon cmdrbase.PluginEntry
	if newAddonEntryFunc, ok := newAddonSymbol.(func() cmdrbase.PluginEntry); ok {
		newAddon = newAddonEntryFunc()
		w.addons = append(w.addons, newAddon)
	}

	// add command into .
	err = w._addonAddCmd(&root.Command, cmdName, desc, newAddon, newAddon)
	return
}

func (w *ExecWorker) _addonAddCmd(parent *Command, cmdName, desc string, addon cmdrbase.PluginEntry, cmd cmdrbase.PluginCmd) (err error) {
	if cmd.Name() != "" {
		cmdName = cmd.Name()
	}
	if cmd.Description() != "" {
		desc = cmd.Description()
	}

	if _, ok := parent.allCmds[AddonsGroup]; !ok {
		parent.allCmds[AddonsGroup] = make(map[string]*Command)
	}
	if _, ok := parent.allCmds[AddonsGroup][cmdName]; !ok {
		cx := &Command{
			BaseOpt: BaseOpt{
				Full:        cmdName,
				Short:       cmd.ShortName(),
				Aliases:     cmd.Aliases(),
				Description: desc,
				Action: func(cmdMatched *Command, args []string) (err error) {
					Logger.Infof("pre - hello, args: %v", args)
					err = cmd.Action(args)
					return
				},
				Hidden: false,
				Group:  AddonsGroup,
				owner:  parent,
			},
		}
		parent.SubCommands = uniAddCmd(parent.SubCommands, cx)
		parent.allCmds[AddonsGroup][cmdName] = cx
		parent.plainCmds[cmdName] = cx
		if cmd.ShortName() != "" && cmd.ShortName() != cmdName {
			if _, ok := parent.plainCmds[cmd.ShortName()]; !ok {
				parent.plainCmds[cmd.ShortName()] = cx
			}
		}
		for _, n := range cmd.Aliases() {
			if _, ok := parent.plainCmds[n]; !ok {
				parent.plainCmds[n] = cx
			}
		}

		// add flags
		for _, ff := range cmd.Flags() {
			err = w._addonAddFlag(cx, addon, ff)
		}

		// children: sub-commands
		for _, cc := range cmd.SubCommands() {
			err = w._addonAddCmd(cx, "", "", addon, cc)
		}
	}
	return
}

//goland:noinspection GoUnusedParameter
func (w *ExecWorker) _addonAddFlag(parent *Command, addon cmdrbase.PluginEntry, flg cmdrbase.PluginFlag) (err error) {
	name, short := flg.Name(), flg.ShortName()
	cx := &Flag{
		BaseOpt: BaseOpt{
			Full:        name,
			Short:       short,
			Aliases:     flg.Aliases(),
			Description: flg.Description(),
			Action: func(cmd *Command, args []string) (err error) {
				err = flg.Action()
				return
			},
			Hidden: false,
			Group:  AddonsGroup,
			owner:  parent,
		},
		DefaultValue:            flg.DefaultValue(),
		DefaultValuePlaceholder: flg.PlaceHolder(),
	}
	w.ensureCommandMaps(parent, "", AddonsGroup)
	parent.Flags = uniAddFlg(parent.Flags, cx)
	parent.allFlags[AddonsGroup][name] = cx
	parent.plainLongFlags[name] = cx
	for _, as := range cx.Aliases {
		parent.plainLongFlags[as] = cx
	}
	if short != "" {
		parent.plainShortFlags[short] = cx
	}
	return
}

// buildExtensionsCrossRefs for cmdr extensions.
//
// An extension, which is an external script or an executable
// typically, can be integrated into host-app cmdr command
// system.
//
//goland:noinspection ALL
func (w *ExecWorker) buildExtensionsCrossRefs(root *RootCommand) {
	flog("    - preprocess / buildXref / buildExtensionsCrossRefs...")
	// prefix := conf.AppName
	for _, d := range w.extensionsLocations {
		dirExpanded := os.ExpandEnv(d)
		// Logger.Debugf("      -> ext.dir: %v", dirExpanded)
		if dir.FileExists(dirExpanded) {
			err := dir.ForFileMax(dirExpanded, 0, 1, func(depth int, cwd string, fi os.FileInfo) (stop bool, err error) {
				if !fi.IsDir() {
					ok := true // = strings.HasPrefix(fi.Name(), prefix)
					// Logger.Debugf("      -> ext.dir: %v, file: %v", dirExpanded, fi.Name())
					if ok && dir.IsModeExecAny(fi.Mode()) {
						err = w._addIt(fi, root, cwd)
					}
				}
				return
			})
			if err != nil {
				Logger.Warnf("  warn - error in buildExtensionsCrossRefs.ForDir(): %v", err)
			}
		}
	}
}

func (w *ExecWorker) _addIt(fi os.FileInfo, root *RootCommand, cwd string) (err error) {
	if fi.Mode().IsRegular() {
		// name := fi.Name()[:len(prefix)]
		name := fi.Name()
		exe := path.Join(cwd, fi.Name())
		// if strings.HasPrefix(name, "-") || strings.HasPrefix(name, "_") {
		//	name = name[1:]
		//	Logger.Debugf("      -> addons.dir: %v, file: %v", dirExpanded, fi.Name())
		w._addAsSubCmd(&root.Command, name, exe)
		// }
	}
	return
}

func (w *ExecWorker) _addAsSubCmd(parent *Command, cmdName, cmdPath string) {
	desc := fmt.Sprintf("execute %q", cmdPath)
	if _, ok := parent.allCmds[ExtGroup]; !ok {
		parent.allCmds[ExtGroup] = make(map[string]*Command)
	}
	if _, ok := parent.allCmds[ExtGroup][cmdName]; !ok {
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
				owner:  parent,
			},
		}
		parent.SubCommands = uniAddCmd(parent.SubCommands, cx)
		parent.allCmds[ExtGroup][cmdName] = cx
		parent.plainCmds[cmdName] = cx
	}
}

//goland:noinspection GoUnusedParameter
func (w *ExecWorker) buildRootCrossRefs(root *RootCommand) {
	flog("    - preprocess / buildXref / buildRootCrossRefs...")

	// initializes the internal variables/members
	w.ensureCmdMembers(&root.Command, root)

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

	w._buildCrossRefs(&root.Command, root)
}

func (w *ExecWorker) ensureCommandMaps(parent *Command, group, flgGroup string) {
	if group == "" {
		group = UnsortedGroup
	}
	if flgGroup == "" {
		flgGroup = UnsortedGroup
	}

	if parent.allCmds == nil {
		parent.allCmds = make(map[string]map[string]*Command)
	}
	if _, ok := parent.allCmds[group]; !ok {
		parent.allCmds[group] = make(map[string]*Command)
	}

	if parent.plainCmds == nil {
		parent.plainCmds = make(map[string]*Command)
	}

	if parent.allFlags == nil {
		parent.allFlags = make(map[string]map[string]*Flag)
	}
	if _, ok := parent.allFlags[flgGroup]; !ok {
		parent.allFlags[flgGroup] = make(map[string]*Flag)
	}

	if parent.plainLongFlags == nil {
		parent.plainLongFlags = make(map[string]*Flag)
	}
	if parent.plainShortFlags == nil {
		parent.plainShortFlags = make(map[string]*Flag)
	}
}

func (w *ExecWorker) attachVersionCommands(root *RootCommand) {
	if w.enableVersionCommands {
		w._cmdAdd(root, "version", "Show the version of this app.", func(cx *Command) {
			cx.Aliases = []string{"ver", "versions"}
			cx.Action = func(cmd *Command, args []string) (err error) {
				w.showVersion(cmd)
				return ErrShouldBeStopException
			}
		})
		w._boolFlgAdd(root, "version", "Show the version of this app.", SysMgmtGroup, func(ff *Flag) {
			ff.Short = "V"
			ff.Aliases = []string{"ver", "versions"}
			ff.Action = func(cmd *Command, args []string) (err error) {
				w.showVersion(cmd)
				return ErrShouldBeStopException
			}
			ff.circuitBreak = true
			ff.justOnce = true
		})
		w._stringFlgAdd(root, "version-sim", "Simulate a faked version number for this app.", SysMgmtGroup, func(ff *Flag) {
			ff.Aliases = []string{"version-simulate"}
			ff.Action = func(cmd *Command, args []string) (err error) {
				conf.Version = GetStringR("version-sim")
				Set("version", conf.Version) // set into option 'app.version' too.
				return
			}
		})
		w._boolFlgAdd(root, "build-info", "Show the building information of this app.", SysMgmtGroup, func(ff *Flag) {
			ff.Short = "#"
			ff.Action = func(cmd *Command, args []string) (err error) {
				w.showBuildInfo(cmd)
				return ErrShouldBeStopException
			}
			ff.circuitBreak = true
			ff.justOnce = true
		})
	}
}

func (w *ExecWorker) attachHelpCommands(root *RootCommand) {
	if w.enableHelpCommands {
		if _, ok := root.allFlags[SysMgmtGroup]["help"]; !ok {
			w._boolFlgAdd(root, "help", "Show this help screen", SysMgmtGroup, func(ff *Flag) {
				ff.Short = "h"
				ff.Aliases = []string{"?", "helpme", "info", "usage"}
				ff.Action = func(cmd *Command, args []string) (err error) {
					// cmdr.Logger.Debugf("-- helpCommand hit. printHelp and stop.")
					// printHelp(cmd)
					// return ErrShouldBeStopException
					return nil
				}
				ff.EnvVars = []string{"HELP"}
				ff.circuitBreak = true
				ff.justOnce = true
			})
			root.plainShortFlags["?"] = root.allFlags[SysMgmtGroup]["help"]

			// w._intFlgAdd(root, "help-zsh", "Show help with zsh completion format, or others", func(ff *Flag) { ff.DefaultValuePlaceholder = "LEVEL" })
			// w._intFlgAdd(root, "help-bash", "Show help with bash completion format, or others", func(ff *Flag) { ff.DefaultValuePlaceholder = "LEVEL" })

			m := map[string]bool{
				"linux": true, "darwin": true,
			}
			if _, ok := m[runtime.GOOS]; ok {
				w._boolFlgAdd(root, "man", "Show help screen in manpage format (INSTALL NEEDED!)", SysMgmtGroup, func(ff *Flag) {
					ff.Action = func(cmd *Command, args []string) (err error) {
						str := strings.ReplaceAll(backtraceCmdNames(cmd, false), ".", "-")
						if !cmd.IsRoot() {
							str = cmd.root.Name + "-" + str
						}
						if e := exec.Run("man", str); e != nil {
							// Logger.Warnf("%v", errors.Unwrap(e).Error())
							if fn, e2 := genManualForCommand(cmd); e2 == nil {
								defer func() { _ = dir.DeleteFile(fn) }()
								closers.RegisterCloseFns(func() {
									d := path.Dir(fn)
									_ = dir.RemoveDirRecursive(d)
								})
								_ = exec.Run("man", fn)
							}
						}
						return ErrShouldBeStopException
					}
					// ff.Hidden = false
					// ff.dblTildeOnly = true
				})
			}

			w._boolFlgAdd(root, "tree", "Show a tree for all commands", SysMgmtGroup, func(ff *Flag) {
				ff.Action = dumpTreeForAllCommands
				ff.dblTildeOnly = true
				ff.VendorHidden = true
			})

			if enableShellCompletionCommand {
				w._cmdAdd(root, "help", "Completion Help system", func(cx *Command) {
					cx.Short = "h"
					cx.Aliases = []string{"__completion", "__complete"}
					cx.VendorHidden = true
					cx.Action = w.helpSystemAction
					cx.onMatched = func(cmd *Command, args []string) (err error) {
						w.inCompleting = true
						// disable trace, debug, info, warn messages
						w.setLoggerLevel(ErrorLevel)
						return
					}
				})
				root.plainCmds["Ø"] = root.allCmds[SysMgmtGroup]["help"]
			}
		}
		w._stringFlgAdd(root, "config", "load config files from where you specified", SysMgmtGroup, func(ff *Flag) {
			ff.Action = func(cmd *Command, args []string) (err error) {
				// cmdr.Logger.Debugf("-- --config hit. printHelp and stop.")
				// return ErrShouldBeStopException
				return nil
			}
			ff.Examples = `
$ {{.AppName}} --configci/etc/demo-yy ~~debug
	try loading config from 'ci/etc/demo-yy', noted that assumes a child folder 'conf.d' should be exists
$ {{.AppName}} --config=ci/etc/demo-yy/any.yml ~~debug
	try loading config from 'ci/etc/demo-yy/any.yml', noted that assumes a child folder 'conf.d' should be exists
`
			ff.DefaultValuePlaceholder = "[Location,...]"
			ff.Hidden = false
		})
	}
}

func (w *ExecWorker) attachVerboseCommands(root *RootCommand) {
	if w.enableVerboseCommands {
		w._boolFlgAdd(root, "verbose", "Show more progress/debug info", SysMgmtGroup, func(ff *Flag) {
			ff.Short = "v"
			ff.EnvVars = []string{"VERBOSE"}
		})

		w._boolFlgAdd(root, "quiet", "No more screen output", SysMgmtGroup, func(ff *Flag) {
			ff.Short = "q"
			ff.EnvVars = []string{"QUIET", "SILENT"}
		})

		w._boolFlgAdd(root, "debug", "Get into debug mode.", SysMgmtGroup, func(ff *Flag) {
			ff.Short = "D"
			ff.EnvVars = []string{"DEBUG"}
			// ff.Action = func(cmd *Command, args []string) (err error) {
			// 	return
			// }
			// ff.onSet = func(keyPath string, value interface{}) {
			// 	flog("--debug: %v => %v", keyPath, value)
			// }
		})
		w._stringFlgAdd(root, "debug-output", "Store the ~~debug outputs into file.", SysMgmtGroup, func(ff *Flag) {
			ff.DefaultValue = "" // "dbg.log"
			ff.EnvVars = []string{"DEBUG_OUTPUT"}
		})

		mutualExclusives := []string{"raw", "value-type", "more", "env"}
		w._boolFlgAdd(root, "env", "Dump environment info in '~~debug' mode.", SysMgmtGroup, func(ff *Flag) {
			ff.prerequisites = []string{"debug"}
			ff.mutualExclusives = mutualExclusives
		})
		w._boolFlgAdd(root, "more", "Dump more info in '~~debug' mode.", SysMgmtGroup, func(ff *Flag) {
			ff.prerequisites = []string{"debug"}
			ff.mutualExclusives = mutualExclusives
		})
		// w._boolFlgAdd(root, "yaml", "Dump as a tree in '~~debug' mode.", SysMgmtGroup, func(ff *Flag) {
		// 	ff.prerequisites = []string{"debug"}
		// 	ff.mutualExclusives = mutualExclusives
		// })
		w._boolFlgAdd(root, "raw", "Dump the option value in raw mode (with golang data structure, without envvar expanding).", SysMgmtGroup, func(ff *Flag) {
			ff.prerequisites = []string{"debug"}
			ff.mutualExclusives = mutualExclusives
		})
		w._boolFlgAdd(root, "value-type", "Dump the option value type.", SysMgmtGroup, func(ff *Flag) {
			ff.prerequisites = []string{"debug"}
			ff.mutualExclusives = mutualExclusives
		})
	}
}

func (w *ExecWorker) attachCmdrCommands(root *RootCommand) {
	if w.enableCmdrCommands {
		w._boolFlgAdd(root, "strict-mode", "Strict mode for 'cmdr'.", SysMgmtGroup, func(ff *Flag) {
			ff.EnvVars, ff.VendorHidden = []string{"STRICT"}, true
		})
		w._boolFlgAdd(root, "no-env-overrides", "No env var overrides for 'cmdr'.", SysMgmtGroup, func(ff *Flag) {
			ff.VendorHidden = true
		})
		w._boolFlgAdd(root, "no-color", "No color output for 'cmdr'.", SysMgmtGroup, func(ff *Flag) {
			ff.Short = "nc"
			ff.Aliases = []string{"nocolor"}
			ff.EnvVars, ff.VendorHidden = []string{"NOCOLOR", "NO_COLOR"}, true
		})

		sbomAttach(w, root)
	}
}

//nolint:funlen //for test
func (w *ExecWorker) attachGeneratorsCommands(root *RootCommand) {
	if w.enableGenerateCommands {
		found := false
		for _, sc := range root.SubCommands {
			if sc.Full == "generate" { // generatorCommands.Full {
				found = true
				break
			}
		}
		if !found {
			// root.SubCommands = append(root.SubCommands, generatorCommands)
			w._cmdAdd(root, "generate", "Generators for this app.", func(cx1 *Command) {
				cx1.Short = "g"
				cx1.Aliases = []string{"gen"}
				cx1.LongDescription = `
[cmdr] includes multiple generators like:

- linux man page generator
- shell completion script generator
- more...

			`
				// - markdown generator

				cx1.Examples = `
$ {{.AppName}} gen sh --bash
			generate bash completion script
$ {{.AppName}} gen shell --auto
			generate shell completion script with detecting on current shell environment.
$ {{.AppName}} gen sh
			generate shell completion script with detecting on current shell environment.
$ {{.AppName}} gen man
			generate linux manual (man page)
			`
				// $ {{.AppName}} gen doc
				// generate document, default markdown.
				// $ {{.AppName}} gen doc --markdown
				// generate markdown.
				// $ {{.AppName}} gen doc --pdf
				// generate pdf.
				// $ {{.AppName}} gen markdown
				// generate markdown.
				// $ {{.AppName}} gen pdf
				// generate pdf.

				w._cmdAdd1(cx1, "shell", "Generate the bash/zsh auto-completion script or install it.", func(cx *Command) {
					cx.Short = "s"
					cx.Aliases = []string{"sh"}
					cx.Action = genShell
					cx.Hidden = false

					w._stringFlgAdd1(cx, "dir", "The output directory", "Output", func(ff *Flag) {
						ff.Short = "d"
						ff.DefaultValue = "."
						ff.DefaultValuePlaceholder = dirString
						ff.Hidden = false
					})

					w._stringFlgAdd1(cx, "output", "The output filename", "Output", func(ff *Flag) {
						ff.Short = "o"
						ff.DefaultValue = "" // os.ExpandEnv("$AppName")
						ff.DefaultValuePlaceholder = "FILENAME"
						ff.Hidden = false
					})

					w._boolFlgAdd1(cx, "auto", "Generate auto completion script to fit for your current env.", shTypeGroup, func(ff *Flag) {
						ff.Short = "a"
						ff.DefaultValue = true
						ff.ToggleGroup = shTypeGroup
						ff.Hidden = false
					})
					w._boolFlgAdd1(cx, "bash", "Generate auto completion script for Bash", shTypeGroup, func(ff *Flag) {
						ff.Short = "b"
						ff.ToggleGroup = shTypeGroup
						ff.Hidden = false
					})
					w._boolFlgAdd1(cx, "zsh", "Generate auto completion script for Zsh", shTypeGroup, func(ff *Flag) {
						ff.Short = "z"
						ff.ToggleGroup = shTypeGroup
						ff.Hidden = false
					})
					w._boolFlgAdd1(cx, "fish", "Generate auto completion script for Fish-shell", shTypeGroup, func(ff *Flag) {
						ff.Short = "f"
						ff.ToggleGroup = shTypeGroup
						ff.Hidden = false
					})
					w._boolFlgAdd1(cx, "powershell", "Generate auto completion script for Powershell", shTypeGroup, func(ff *Flag) {
						ff.Short = "ps"
						ff.ToggleGroup = shTypeGroup
						ff.Hidden = false
					})
					w._boolFlgAdd1(cx, "elvish", "Generate auto completion script for Elvish-shell [TODO]", shTypeGroup, func(ff *Flag) {
						ff.Short = "el"
						ff.ToggleGroup = shTypeGroup
						ff.VendorHidden = true
					})
					w._boolFlgAdd1(cx, "fig", "Generate auto completion script for fig-shell [TODO]", shTypeGroup, func(ff *Flag) {
						ff.Short = "fig"
						ff.ToggleGroup = shTypeGroup
						ff.VendorHidden = true
					})
					w._boolFlgAdd1(cx, "force-bash", "Just for --auto", shTypeGroup, func(ff *Flag) {
						ff.Hidden = false
						ff.prerequisites = []string{"auto"}
					})
				})
				w._cmdAdd1(cx1, "manual", "Generate linux man page.", func(cx *Command) {
					cx.Short = "m"
					cx.Aliases = []string{"man"}
					cx.Action = genManual
					cx.Hidden = false

					w._stringFlgAdd1(cx, "dir", "The output directory", "Output", func(ff *Flag) {
						ff.Short = "d"
						ff.DefaultValue = "./man1"
						ff.DefaultValuePlaceholder = dirString
						ff.Hidden = false
					})
				})
				w._cmdAdd1(cx1, "doc", "Generate a markdown document, or: pdf/TeX/...", func(cx *Command) {
					cx.Short = "d"
					cx.Aliases = []string{"pdf", "docx", "tex", "markdown"}
					cx.Action = genDoc
					cx.VendorHidden = true
					cx.Deprecated = "1.9.9"

					w._stringFlgAdd1(cx, "dir", "The output directory", "Output", func(ff *Flag) {
						ff.Short = "d"
						ff.DefaultValue = "./docs"
						ff.DefaultValuePlaceholder = dirString
						ff.Hidden = false
					})
					w._boolFlgAdd1(cx, "markdown", "To generate a markdown file", tgDocType, func(ff *Flag) {
						ff.Short = "md"
						ff.Aliases = []string{"mkd", "m"}
						ff.ToggleGroup = tgDocType
						ff.DefaultValue = true
					})
					w._boolFlgAdd1(cx, "pdf", "To generate a PDF file", tgDocType, func(ff *Flag) {
						ff.Short = "p"
						ff.ToggleGroup = tgDocType
					})
					w._boolFlgAdd1(cx, "docx", "To generate a Word (.docx) file", tgDocType, func(ff *Flag) {
						ff.Aliases = []string{"doc"}
						ff.ToggleGroup = tgDocType
					})
					w._boolFlgAdd1(cx, "tex", "To generate a LaTeX file", tgDocType, func(ff *Flag) {
						ff.Short = "t"
						ff.ToggleGroup = tgDocType
					})
				})
			})
		}
	}
}

const (
	shTypeGroup = "ShellType"

	tgDocType = "DocType"

	dirString = "DIR"
)

// enableShellCompletionCommand _
var enableShellCompletionCommand = true

func (w *ExecWorker) _boolFlgAdd(root *RootCommand, full, desc, group string, adding func(ff *Flag)) {
	w._boolFlgAdd1(&root.Command, full, desc, group, adding)
}

func (w *ExecWorker) _boolFlgAdd1(parent *Command, full, desc, group string, adding func(ff *Flag)) {
	if group == "" {
		group = UnsortedGroup
	}
	if _, ok := parent.allFlags[group][full]; !ok {
		ff := &Flag{
			BaseOpt: BaseOpt{
				Full:        full,
				Description: desc,
				Hidden:      true,
				Group:       group,
				owner:       parent,
			},
			DefaultValue: false,
		}
		if adding != nil {
			adding(ff)
		}
		parent.Flags = append(parent.Flags, ff)
		if group != "" {
			w.ensureCommandMaps(parent, "", group)
			parent.allFlags[group][full] = ff
		}
		parent.plainLongFlags[full] = ff
		if ff.Short != "" {
			// NOTE: dup short title would be ignored
			if _, ok := parent.plainShortFlags[ff.Short]; !ok {
				parent.plainShortFlags[ff.Short] = ff
			}
		}
		for _, t := range ff.Aliases {
			if t != "" {
				// NOTE: dup aliases would be ignored
				if _, ok := parent.plainLongFlags[t]; !ok {
					parent.plainLongFlags[t] = ff
				}
			}
		}
	} else {
		Logger.Warnf("duplicated bool flag %q had been adding.", full)
	}
}

func (w *ExecWorker) _intFlgAdd(root *RootCommand, full, desc, group string, adding func(ff *Flag)) {
	w._intFlgAdd1(&root.Command, full, desc, group, adding)
}

func (w *ExecWorker) _intFlgAdd1(parent *Command, full, desc, group string, adding func(ff *Flag)) {
	if group == "" {
		group = UnsortedGroup
	}
	if _, ok := parent.allFlags[group][full]; !ok {
		ff := &Flag{
			BaseOpt: BaseOpt{
				Full:        full,
				Description: desc,
				Hidden:      true,
				Group:       group,
				owner:       parent,
			},
			DefaultValue: 0,
		}
		if adding != nil {
			adding(ff)
		}
		parent.Flags = append(parent.Flags, ff)
		w.ensureCommandMaps(parent, "", group)
		parent.allFlags[group][full] = ff
		parent.plainLongFlags[full] = ff
		if ff.Short != "" {
			// NOTE: dup short title would be ignored
			if _, ok := parent.plainShortFlags[ff.Short]; !ok {
				parent.plainShortFlags[ff.Short] = ff
			}
		}
		for _, t := range ff.Aliases {
			if t != "" {
				// NOTE: dup aliases would be ignored
				if _, ok := parent.plainLongFlags[t]; !ok {
					parent.plainLongFlags[t] = ff
				}
			}
		}
	} else {
		Logger.Warnf("duplicated int flag %q had been adding.", full)
	}
}

func (w *ExecWorker) _stringFlgAdd(root *RootCommand, full, desc, group string, adding func(ff *Flag)) {
	w._stringFlgAdd1(&root.Command, full, desc, group, adding)
}

func (w *ExecWorker) _stringFlgAdd1(parent *Command, full, desc, group string, adding func(ff *Flag)) {
	if group == "" {
		group = UnsortedGroup
	}
	if _, ok := parent.allFlags[group][full]; !ok {
		ff := &Flag{
			BaseOpt: BaseOpt{
				Full:        full,
				Description: desc,
				Hidden:      true,
				Group:       group,
				owner:       parent,
			},
			DefaultValue: "",
		}
		if adding != nil {
			adding(ff)
		}
		parent.Flags = append(parent.Flags, ff)
		w.ensureCommandMaps(parent, "", group)
		parent.allFlags[group][full] = ff
		parent.plainLongFlags[full] = ff
		if ff.Short != "" {
			// NOTE: dup short title would be ignored
			if _, ok := parent.plainShortFlags[ff.Short]; !ok {
				parent.plainShortFlags[ff.Short] = ff
			}
		}
		for _, t := range ff.Aliases {
			if t != "" {
				// NOTE: dup aliases would be ignored
				if _, ok := parent.plainLongFlags[t]; !ok {
					parent.plainLongFlags[t] = ff
				}
			}
		}
	} else {
		Logger.Warnf("duplicated string flag %q had been adding.", full)
	}
}

func (w *ExecWorker) _cmdAdd(root *RootCommand, full, desc string, adding func(cx *Command)) {
	w._cmdAdd1(&root.Command, full, desc, adding)
}

func (w *ExecWorker) _cmdAdd1(parent *Command, full, desc string, adding func(cx *Command)) {
	if _, ok := parent.allCmds[SysMgmtGroup][full]; !ok {
		cx := &Command{
			BaseOpt: BaseOpt{
				Full:        full,
				Description: desc,
				Hidden:      true,
				Group:       SysMgmtGroup,
				owner:       parent,
			},
		}
		w.ensureCommandMaps(parent, cx.Group, "")
		if adding != nil {
			adding(cx)
		}
		w.ensureCommandMaps(parent, cx.Group, "")
		parent.SubCommands = uniAddCmd(parent.SubCommands, cx)
		parent.allCmds[cx.Group][full] = cx
		parent.plainCmds[full] = cx
		if cx.Short != "" {
			if _, ok := parent.allCmds[cx.Group][cx.Short]; !ok {
				// parent.allCmds[cx.Group][cx.Short] = cx
				parent.plainCmds[cx.Short] = cx
			}
		}
		for _, t := range cx.Aliases {
			if t != "" {
				if _, ok := parent.allCmds[cx.Group][t]; !ok {
					// parent.allCmds[cx.Group][t] = cx
					parent.plainCmds[t] = cx
				}
			}
		}
	} else {
		Logger.Warnf("duplicated command %q had been adding.", full)
	}
}

func (w *ExecWorker) _buildCrossRefs(cmd *Command, root *RootCommand) {
	w.ensureCmdMembers(cmd, root)

	singleFlagNames := make(map[string]bool)
	stringFlagNames := make(map[string]bool)
	singleCmdNames := make(map[string]bool)
	stringCmdNames := make(map[string]bool)
	tgs := make(map[string]bool)

	bre := regexp.MustCompile("`(.+)`")

	for _, flg := range cmd.Flags {
		flg.owner = cmd

		if flg.ToggleGroup != "" {
			if flg.Group == "" {
				flg.Group = flg.ToggleGroup
			}
			tgs[flg.ToggleGroup] = true
		}

		if b := bre.Find([]byte(flg.Description)); flg.DefaultValuePlaceholder == "" && len(b) > 2 {
			ph := strings.ToUpper(strings.Trim(string(b), "`"))
			flg.DefaultValuePlaceholder = ph
		}

		w._buildCrossRefsForFlag(flg, cmd, singleFlagNames, stringFlagNames)

		// opt.Children[flg.Full] = &OptOne{Value: flg.DefaultValue,}
		w.rxxtOptions.Set(backtraceFlagNames(flg), flg.DefaultValue)
	}

	for _, cx := range cmd.SubCommands {
		cx.owner = cmd

		w._buildCrossRefsForCommand(cx, cmd, singleCmdNames, stringCmdNames)
		// opt.Children[cx.Full] = newOpt()

		w.rxxtOptions.Set(backtraceCmdNames(cx, false), nil)
		// buildCrossRefs(cx, opt.Children[cx.Full])
		w._buildCrossRefs(cx, root)
	}

	for tg := range tgs {
		w.buildToggleGroup(tg, cmd)
	}
}

func (w *ExecWorker) _buildCrossRefsForFlag(flg *Flag, cmd *Command, singleFlagNames, stringFlagNames map[string]bool) {
	w.forFlagNames(flg, cmd, singleFlagNames, stringFlagNames)

	for _, sz := range flg.Aliases {
		if _, ok := stringFlagNames[sz]; ok {
			ferr("\nNOTE: flag alias name '%v' has been used. (command: %v)", sz, backtraceCmdNames(cmd, false))
		} else {
			stringFlagNames[sz] = true
		}
	}
	if flg.Group == "" {
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
	if flg.Short != "" {
		if _, ok := singleFlagNames[flg.Short]; ok {
			ferr("\nNOTE: flag char '%v' has been used. (command: %v)", flg.Short, backtraceCmdNames(cmd, false))
		} else {
			singleFlagNames[flg.Short] = true
		}
	}
	if flg.Full != "" {
		if _, ok := stringFlagNames[flg.Full]; ok {
			ferr("\nNOTE: flag '%v' has been used. (command: %v)", flg.Full, backtraceCmdNames(cmd, false))
		} else {
			stringFlagNames[flg.Full] = true
		}
	}
	if flg.Short == "" && flg.Full == "" && flg.Name != "" {
		if _, ok := stringFlagNames[flg.Name]; ok {
			ferr("\nNOTE: flag '%v' has been used. (command: %v)", flg.Name, backtraceCmdNames(cmd, false))
		} else {
			stringFlagNames[flg.Name] = true
		}
	}
}

func (w *ExecWorker) _buildCrossRefsForCommand(cx, cmd *Command, singleCmdNames, stringCmdNames map[string]bool) {
	w.forCommandNames(cx, cmd, singleCmdNames, stringCmdNames)

	for _, sz := range cx.Aliases {
		if sz != "" {
			if _, ok := stringCmdNames[sz]; ok {
				ferr("\nNOTE: command alias name '%v' has been used. (command: %v)", sz, backtraceCmdNames(cmd, false))
			} else {
				stringCmdNames[sz] = true
			}
		}
	}

	if cx.Group == "" {
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
	if cx.Short != "" {
		if _, ok := singleCmdNames[cx.Short]; ok {
			ferr("\nNOTE: command char '%v' has been used. (command: %v)", cx.Short, backtraceCmdNames(cmd, false))
		} else {
			singleCmdNames[cx.Short] = true
		}
	}
	if cx.Full != "" {
		if _, ok := stringCmdNames[cx.Full]; ok {
			ferr("\nNOTE: command '%v' has been used. (command: %v)", cx.Full, backtraceCmdNames(cmd, false))
		} else {
			stringCmdNames[cx.Full] = true
		}
	}
	if cx.Short == "" && cx.Full == "" && cx.Name != "" {
		if _, ok := stringCmdNames[cx.Name]; ok {
			ferr("\nNOTE: command '%v' has been used. (command: %v)", cx.Name, backtraceCmdNames(cmd, false))
		} else {
			stringCmdNames[cx.Name] = true
		}
		cmd.plainCmds[cx.Name] = cx
	}
}

func (w *ExecWorker) buildToggleGroup(tg string, cmd *Command) {
	for _, f := range cmd.Flags {
		if tg == f.ToggleGroup && f.DefaultValue == true {
			w.rxxtOptions.Set(backtraceFlagNames(f), true)
			w.rxxtOptions.Set(backtraceCmdNames(cmd, false)+"."+tg, f.Full)
			break
		}
	}
}

// DottedPathToCommand searches the matched Command with the specified dotted-path.
// The searching will start from root if anyCmd is nil.
func DottedPathToCommand(dottedPath string, anyCmd *Command) (cc *Command) {
	return dottedPathToCommand(dottedPath, anyCmd)
}

// dottedPathToCommand searches the matched Command with the specified dotted-path.
// The searching will start from root if anyCmd is nil.
func dottedPathToCommand(dottedPath string, anyCmd *Command) (cc *Command) {
	var c *Command
	if anyCmd == nil {
		c = &internalGetWorker().rootCommand.Command
	} else {
		c = anyCmd
	}

	if err := walkFromCommand(c, 0, 0,
		func(cmd *Command, index, level int) (err error) {
			if cmd.GetDottedNamePath() == dottedPath {
				cc, err = cmd, ErrShouldBeStopException
			}
			return
		}); err == nil && cc != nil {
		return
	}

	return
}

// DottedPathToCommandOrFlag searches the matched Command or Flag with the specified dotted-path.
// The searching will start from root if anyCmd is nil.
func DottedPathToCommandOrFlag(dottedPath string, anyCmd *Command) (cc *Command, ff *Flag) {
	return dottedPathToCommandOrFlag(dottedPath, anyCmd)
}

// dottedPathToCommandOrFlag searches the matched Command or Flag with the specified dotted-path.
// The searching will start from root if anyCmd is nil.
func dottedPathToCommandOrFlag(dottedPath string, anyCmd *Command) (cc *Command, ff *Flag) {
	if anyCmd == nil {
		anyCmd = &internalGetWorker().rootCommand.Command
	}

	if !strings.HasPrefix(dottedPath, anyCmd.root.AppName) {
		dottedPath = anyCmd.root.AppName + "." + dottedPath
	}

	if err := walkFromCommand(anyCmd, 0, 0,
		func(cmd *Command, index, level int) (err error) {
			kp := cmd.GetDottedNamePath()
			if !strings.HasPrefix(kp, anyCmd.root.AppName) {
				kp = anyCmd.root.AppName + "." + kp
			}

			if kp == dottedPath {
				cc, err = cmd, ErrShouldBeStopException
				return
			}

			if strings.HasPrefix(dottedPath, kp) {
				parts := strings.TrimPrefix(dottedPath, kp+".")
				if !strings.Contains(parts, ".") {
					// try matching flags in this command
					for _, f := range cmd.Flags {
						if parts == f.Full {
							ff, err = f, ErrShouldBeStopException
							return
						}
					}
				}
			}
			return
		}); err == nil && cc != nil {
		return
	}

	return
}

// func (w *ExecWorker) backtraceFlagNames(flg *Flag) (str string) { return backtraceFlagNames(flg) }

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

// backtraceCmdNames returns the sequences of a sub-command from
// top-level.
//
// - if verboseLast = false, got 'microservices.tags.list' for sub-cmd microservice/tags/list.
//
// - if verboseLast = true,  got 'microservices.tags.[ls|list|l|lst|dir]'.
//
// - at root command, it returns 'appName' or ” when verboseLast is true.
func backtraceCmdNames(cmd *Command, verboseLast bool) (str string) {
	var a []string
	if verboseLast {
		va := cmd.GetTitleNamesArray()
		if len(va) > 0 {
			vas := strings.Join(va, "|")
			a = append(a, "["+vas+"]")
		}
	} else {
		a = append(a, cmd.GetTitleName())
	}
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

func (w *ExecWorker) ensureCmdMembers(cmd *Command, root *RootCommand) *Command {
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
		cmd.root = root // w.rootCommand
	}
	return cmd
}
