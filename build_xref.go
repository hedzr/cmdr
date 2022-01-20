/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"fmt"
	cmdrbase "github.com/hedzr/cmdr-base"
	"github.com/hedzr/cmdr/conf"
	"github.com/hedzr/log/dir"
	"github.com/hedzr/log/exec"
	"os"
	"path"
	"plugin"
	"regexp"
	"runtime"
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

func (w *ExecWorker) buildXref(rootCmd *RootCommand, args []string) (err error) {
	flog("--> preprocess / buildXref")

	// build xref for root command and its all sub-commands and flags
	// and build the default values
	w.buildRootCrossRefs(rootCmd)
	w.buildAddonsCrossRefs(rootCmd)
	w.buildExtensionsCrossRefs(rootCmd)

	w.setupFromEnvvarMap()

	flog("--> before-config-file-loading")
	for _, x := range w.beforeConfigFileLoading {
		if x != nil {
			x(rootCmd, args)
		}
	}

	if !w.doNotLoadingConfigFiles {
		// flog("--> buildXref: loadFromPredefinedLocations()")

		// pre-detects for `--config xxx`, `--config=xxx`, `--configxxx`
		//if err = w.parsePredefinedLocation(); err != nil {
		//	return
		//}
		_ = w.parsePredefinedLocation()

		// and now, loading the external configuration files
		err = w.loadFromPredefinedLocations(rootCmd)

		err = w.loadFromAlterLocations(rootCmd)

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

//goland:noinspection GoUnusedParameter
func (w *ExecWorker) buildAliasesCrossRefs(root *RootCommand) {
	var (
		aliases *aliasesCommands = new(aliasesCommands)
		err     error
	)
	err = GetSectionFrom("aliases", &aliases)
	if err == nil {
		err = w._addCommandsForAliasesGroup(root, aliases)
	}
	if err != nil {
		Logger.Errorf("buildAliasesCrossRefs error: %v", err)
	}
}

//goland:noinspection GoUnusedParameter
func (w *ExecWorker) _addCommandsForAliasesGroup(root *RootCommand, aliases *aliasesCommands) (err error) {
	flog("aliases:\n%v\n", aliases)
	if aliases.Group == "" {
		aliases.Group = AliasesGroup
	}

	for _, cmd := range aliases.Commands {
		w.ensureCmdMembers(cmd)
		err = w._toolAddCmd(&root.Command, aliases.Group, cmd)
	}
	w._buildCrossRefs(&root.Command)
	return
}

func (w *ExecWorker) _toolAddCmd(parent *Command, groupName string, cc *Command) (err error) {
	if _, ok := parent.allCmds[groupName]; !ok {
		parent.allCmds[groupName] = make(map[string]*Command)
	}
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
				err = w._toolAddCmd(c, groupName, c)
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
		if cx, matched := w.locateCommand(from.Invoke, cmd); matched {
			if cx.Action != nil {
				err = cx.Action(cmd, args)
			}
		}
		return
	}
}

func (w *ExecWorker) getInvokeProcAction(from *Command) Handler {
	return func(cmd *Command, args []string) (err error) {
		cmdParts := strings.Split(from.InvokeProc, " ")
		c, args := cmdParts[0], cmdParts[1:]
		err = exec.Run(c, args...)
		return
	}
}

func (w *ExecWorker) getInvokeShellAction(from *Command) Handler {
	return func(cmd *Command, args []string) (err error) {
		cmdParts := strings.Split(from.InvokeShell, " ")
		c, args := cmdParts[0], cmdParts[1:]
		err = exec.Run(c, args...)
		return
	}
}

// buildAddonsCrossRefs for cmdr addons.
//
// A cmdr addon, which is a golang plugin, can be integrated into host-app better than an extension.
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
			err := dir.ForDirMax(dirExpanded, 0, 1, func(depth int, cwd string, fi os.FileInfo) (stop bool, err error) {
				if fi.IsDir() {
					return
				}
				var ok bool // = strings.HasPrefix(fi.Name(), prefix)
				ok = true
				// Logger.Debugf("      -> addons.dir: %v, file: %v", dirExpanded, fi.Name())
				if ok && fi.Mode().IsRegular() && dir.IsModeExecAny(fi.Mode()) {
					//name := fi.Name()[:len(prefix)]
					name := fi.Name()
					exe := path.Join(cwd, fi.Name())
					//if strings.HasPrefix(name, "-") || strings.HasPrefix(name, "_") {
					//	name = name[1:]
					//	Logger.Debugf("      -> addons.dir: %v, file: %v", dirExpanded, fi.Name())
					err = w._addonAsSubCmd(root, name, exe)
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

func (w *ExecWorker) _addonAsSubCmd(root *RootCommand, cmdName, cmdPath string) (err error) {
	var desc string
	desc = fmt.Sprintf("execute %q", cmdPath)

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

	newAddonEntryFunc := newAddonSymbol.(func() cmdrbase.PluginEntry)
	newAddon := newAddonEntryFunc()
	w.addons = append(w.addons, newAddon)

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
			err = w._addonAddFlg(cx, addon, ff)
		}

		// children: sub-commands
		for _, cc := range cmd.SubCommands() {
			err = w._addonAddCmd(cx, "", "", addon, cc)
		}
	}
	return
}

//goland:noinspection GoUnusedParameter
func (w *ExecWorker) _addonAddFlg(parent *Command, addon cmdrbase.PluginEntry, flg cmdrbase.PluginFlag) (err error) {
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

//goland:noinspection ALL
func (w *ExecWorker) buildExtensionsCrossRefs(root *RootCommand) {
	flog("    - preprocess / buildXref / buildExtensionsCrossRefs...")
	// prefix := conf.AppName
	for _, d := range w.extensionsLocations {
		dirExpanded := os.ExpandEnv(d)
		// Logger.Debugf("      -> ext.dir: %v", dirExpanded)
		if dir.FileExists(dirExpanded) {
			err := dir.ForDirMax(dirExpanded, 0, 1, func(depth int, cwd string, fi os.FileInfo) (stop bool, err error) {
				if fi.IsDir() {
					return
				}
				var ok bool // = strings.HasPrefix(fi.Name(), prefix)
				ok = true
				// Logger.Debugf("      -> ext.dir: %v, file: %v", dirExpanded, fi.Name())
				if ok && fi.Mode().IsRegular() && dir.IsModeExecAny(fi.Mode()) {
					//name := fi.Name()[:len(prefix)]
					name := fi.Name()
					exe := path.Join(cwd, fi.Name())
					//if strings.HasPrefix(name, "-") || strings.HasPrefix(name, "_") {
					//	name = name[1:]
					//	Logger.Debugf("      -> ext.dir: %v, file: %v", dirExpanded, fi.Name())
					w._addAsSubCmd(&root.Command, name, exe)
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

func (w *ExecWorker) _addAsSubCmd(parent *Command, cmdName, cmdPath string) {
	var desc string
	desc = fmt.Sprintf("execute %q", cmdPath)
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

func (w *ExecWorker) _boolFlgAdd(root *RootCommand, full, desc string, adding func(ff *Flag)) {
	w._boolFlgAdd1(&root.Command, full, desc, adding)
}

func (w *ExecWorker) _boolFlgAdd1(parent *Command, full, desc string, adding func(ff *Flag)) {
	if _, ok := parent.allFlags[SysMgmtGroup][full]; !ok {
		ff := &Flag{
			BaseOpt: BaseOpt{
				Full:        full,
				Description: desc,
				Hidden:      true,
				Group:       SysMgmtGroup,
				owner:       parent,
			},
			DefaultValue: false,
		}
		if adding != nil {
			adding(ff)
		}
		parent.Flags = append(parent.Flags, ff)
		w.ensureCommandMaps(parent, "", SysMgmtGroup)
		parent.allFlags[SysMgmtGroup][full] = ff
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

func (w *ExecWorker) _intFlgAdd(root *RootCommand, full, desc string, adding func(ff *Flag)) {
	w._intFlgAdd1(&root.Command, full, desc, adding)
}

func (w *ExecWorker) _intFlgAdd1(parent *Command, full, desc string, adding func(ff *Flag)) {
	if _, ok := parent.allFlags[SysMgmtGroup][full]; !ok {
		ff := &Flag{
			BaseOpt: BaseOpt{
				Full:        full,
				Description: desc,
				Hidden:      true,
				Group:       SysMgmtGroup,
				owner:       parent,
			},
			DefaultValue: 0,
		}
		if adding != nil {
			adding(ff)
		}
		parent.Flags = append(parent.Flags, ff)
		w.ensureCommandMaps(parent, "", SysMgmtGroup)
		parent.allFlags[SysMgmtGroup][full] = ff
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

func (w *ExecWorker) _stringFlgAdd(root *RootCommand, full, desc string, adding func(ff *Flag)) {
	w._stringFlgAdd1(&root.Command, full, desc, adding)
}

func (w *ExecWorker) _stringFlgAdd1(parent *Command, full, desc string, adding func(ff *Flag)) {
	if _, ok := parent.allFlags[SysMgmtGroup][full]; !ok {
		ff := &Flag{
			BaseOpt: BaseOpt{
				Full:        full,
				Description: desc,
				Hidden:      true,
				Group:       SysMgmtGroup,
				owner:       parent,
			},
			DefaultValue: "",
		}
		if adding != nil {
			adding(ff)
		}
		parent.Flags = append(parent.Flags, ff)
		w.ensureCommandMaps(parent, "", SysMgmtGroup)
		parent.allFlags[SysMgmtGroup][full] = ff
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

func (w *ExecWorker) attachVersionCommands(root *RootCommand) {
	if w.enableVersionCommands {
		w._cmdAdd(root, "version", "Show the version of this app.", func(cx *Command) {
			cx.Aliases = []string{"ver", "versions"}
			cx.Action = func(cmd *Command, args []string) (err error) {
				w.showVersion()
				return ErrShouldBeStopException
			}
		})
		w._boolFlgAdd(root, "version", "Show the version of this app.", func(ff *Flag) {
			ff.Short = "V"
			ff.Aliases = []string{"ver", "versions"}
			ff.Action = func(cmd *Command, args []string) (err error) {
				w.showVersion()
				return ErrShouldBeStopException
			}
			ff.circuitBreak = true
			ff.justOnce = true
		})
		w._stringFlgAdd(root, "version-sim", "Simulate a faked version number for this app.", func(ff *Flag) {
			ff.Aliases = []string{"version-simulate"}
			ff.Action = func(cmd *Command, args []string) (err error) {
				conf.Version = GetStringR("version-sim")
				Set("version", conf.Version) // set into option 'app.version' too.
				return
			}
		})
		w._boolFlgAdd(root, "build-info", "Show the building information of this app.", func(ff *Flag) {
			ff.Short = "#"
			ff.Action = func(cmd *Command, args []string) (err error) {
				w.showBuildInfo()
				return ErrShouldBeStopException
			}
			ff.circuitBreak = true
			ff.justOnce = true
		})
		//if _, ok := root.allCmds[SysMgmtGroup]["version"]; !ok {
		//	cx := &Command{
		//		BaseOpt: BaseOpt{
		//			Full:        "version",
		//			Aliases:     []string{"ver", "versions"},
		//			Description: "Show the version of this app.",
		//			Action: func(cmd *Command, args []string) (err error) {
		//				w.showVersion()
		//				return ErrShouldBeStopException
		//			},
		//			Hidden: true,
		//			Group:  SysMgmtGroup,
		//			owner:  &root.Command,
		//		},
		//	}
		//	root.SubCommands = uniAddCmd(root.SubCommands, cx)
		//	root.allCmds[SysMgmtGroup]["version"] = cx
		//	root.allCmds[SysMgmtGroup]["versions"] = cx
		//	root.allCmds[SysMgmtGroup]["ver"] = cx
		//	root.plainCmds["version"] = cx
		//	root.plainCmds["versions"] = cx
		//	root.plainCmds["ver"] = cx
		//}
		//if _, ok := root.allFlags[SysMgmtGroup]["version"]; !ok {
		//	ff := &Flag{
		//		BaseOpt: BaseOpt{
		//			Short:       "V",
		//			Full:        "version",
		//			Aliases:     []string{"ver", "versions"},
		//			Description: "Show the version of this app.",
		//			Hidden:      true,
		//			Group:       SysMgmtGroup,
		//			owner:       &root.Command,
		//			Action: func(cmd *Command, args []string) (err error) {
		//				w.showVersion()
		//				return ErrShouldBeStopException
		//			},
		//		},
		//		DefaultValue: false,
		//	}
		//	root.Flags = append(root.Flags, ff)
		//	root.allFlags[SysMgmtGroup]["version"] = ff
		//	root.plainLongFlags["version"] = ff
		//	root.plainLongFlags["versions"] = ff
		//	root.plainLongFlags["ver"] = ff
		//	root.plainShortFlags["V"] = ff
		//}
		//if _, ok := root.allFlags[SysMgmtGroup]["version-sim"]; !ok {
		//	ff := &Flag{
		//		BaseOpt: BaseOpt{
		//			Full:        "version-sim",
		//			Aliases:     []string{"version-simulate"},
		//			Description: "Simulate a faked version number for this app.",
		//			Hidden:      true,
		//			Group:       SysMgmtGroup,
		//			owner:       &root.Command,
		//			Action: func(cmd *Command, args []string) (err error) {
		//				conf.Version = GetStringR("version-sim")
		//				Set("version", conf.Version) // set into option 'app.version' too.
		//				return
		//			},
		//		},
		//		DefaultValue: "",
		//	}
		//	root.Flags = append(root.Flags, ff)
		//	root.allFlags[SysMgmtGroup]["version-sim"] = ff
		//	root.plainLongFlags["version-sim"] = ff
		//	root.plainLongFlags["version-simulate"] = ff
		//}
		//if _, ok := root.allFlags[SysMgmtGroup]["build-info"]; !ok {
		//	root.allFlags[SysMgmtGroup]["build-info"] = &Flag{
		//		BaseOpt: BaseOpt{
		//			Full:        "build-info",
		//			Short:       "#",
		//			Aliases:     []string{},
		//			Description: "Show the building information of this app.",
		//			Hidden:      true,
		//			Group:       SysMgmtGroup,
		//			owner:       &root.Command,
		//			Action: func(cmd *Command, args []string) (err error) {
		//				w.showBuildInfo()
		//				return ErrShouldBeStopException
		//			},
		//		},
		//		DefaultValue: false,
		//	}
		//	root.plainShortFlags["#"] = root.allFlags[SysMgmtGroup]["build-info"]
		//	root.plainLongFlags["build-info"] = root.allFlags[SysMgmtGroup]["build-info"]
		//}
	}
}

func (w *ExecWorker) attachHelpCommands(root *RootCommand) {
	if w.enableHelpCommands {
		if _, ok := root.allFlags[SysMgmtGroup]["help"]; !ok {
			w._boolFlgAdd(root, "help", "Show this help screen", func(ff *Flag) {
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
			//ff := &Flag{
			//	BaseOpt: BaseOpt{
			//		Short:       "h",
			//		Full:        "help",
			//		Aliases:     []string{"?", "helpme", "info", "usage"},
			//		Description: "Show this help screen",
			//		Hidden:      true,
			//		Group:       SysMgmtGroup,
			//		owner:       &root.Command,
			//		Action: func(cmd *Command, args []string) (err error) {
			//			// cmdr.Logger.Debugf("-- helpCommand hit. printHelp and stop.")
			//			// printHelp(cmd)
			//			// return ErrShouldBeStopException
			//			return nil
			//		},
			//	},
			//	DefaultValue: false,
			//	EnvVars:      []string{"HELP"},
			//}
			//root.Flags = append(root.Flags, ff)
			//root.allFlags[SysMgmtGroup]["help"] = ff
			//root.plainLongFlags["help"] = ff
			//root.plainLongFlags["helpme"] = ff
			//root.plainLongFlags["info"] = ff
			//root.plainLongFlags["usage"] = ff
			//root.plainShortFlags["h"] = ff
			//root.plainShortFlags["?"] = ff

			//w._intFlgAdd(root, "help-zsh", "show help with zsh completion format, or others", func(ff *Flag) { ff.DefaultValuePlaceholder = "LEVEL" })
			//w._intFlgAdd(root, "help-bash", "show help with bash completion format, or others", func(ff *Flag) { ff.DefaultValuePlaceholder = "LEVEL" })

			if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
				w._boolFlgAdd(root, "man", "show help screen in manpage format", func(ff *Flag) {
					ff.Action = func(cmd *Command, args []string) (err error) {
						str := strings.ReplaceAll(w.backtraceCmdNames(cmd, false), ".", "-")
						if !cmd.IsRoot() {
							str = cmd.root.Name + "-" + str
						}
						if e := exec.Run("man", str); e != nil {
							return e
						}
						return ErrShouldBeStopException
					}
					ff.Hidden = false
					//ff.dblTildeOnly = true
				})
			}

			w._boolFlgAdd(root, "tree", "show a tree for all commands", func(ff *Flag) {
				ff.Action = dumpTreeForAllCommands
				ff.dblTildeOnly = true
			})

			if enableShellCompletionCommand {
				w._cmdAdd(root, "help", "help system", func(cx *Command) {
					cx.Short = "h"
					cx.Action = w.helpSystemPrint
				})
				root.plainCmds["Ø"] = root.allCmds[SysMgmtGroup]["help"]
			}
		}
		w._stringFlgAdd(root, "config", "load config files from where you specified", func(ff *Flag) {
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

// enableShellCompletionCommand NOT YET, to-do
var enableShellCompletionCommand bool

func (w *ExecWorker) attachVerboseCommands(root *RootCommand) {
	if w.enableVerboseCommands {
		w._boolFlgAdd(root, "verbose", "Show more progress/debug info", func(ff *Flag) {
			ff.Short = "v"
			ff.EnvVars = []string{"VERBOSE"}
		})
		w._boolFlgAdd(root, "quiet", "No more screen output", func(ff *Flag) {
			ff.Short = "q"
			ff.EnvVars = []string{"QUIET", "SILENT"}
		})
		w._boolFlgAdd(root, "debug", "Get into debug mode.", func(ff *Flag) {
			ff.Short = "D"
			ff.EnvVars = []string{"DEBUG"}
		})
		w._stringFlgAdd(root, "debug-output", "store the ~~debug outputs into file.", func(ff *Flag) {
			ff.DefaultValue = "dbg.log"
			ff.EnvVars = []string{"DEBUG_OUTPUT"}
		})
		w._boolFlgAdd(root, "env", "Dump environment info in `~~debug` mode.", func(ff *Flag) {
			ff.prerequisites = []string{"debug"}
			ff.mutualExclusives = []string{"raw", "value-type", "more", "env"}
		})
		w._boolFlgAdd(root, "raw", "Dump the option value in raw mode (with golang data structure, without envvar expanding).", func(ff *Flag) {
			ff.prerequisites = []string{"debug"}
		})
		w._boolFlgAdd(root, "value-type", "Dump the option value type.", func(ff *Flag) {
			ff.prerequisites = []string{"debug"}
		})
		w._boolFlgAdd(root, "more", "Dump more info in `~~debug` mode.", func(ff *Flag) {
			ff.prerequisites = []string{"debug"}
		})
		//if _, ok := root.allFlags[SysMgmtGroup]["verbose"]; !ok {
		//	ff := &Flag{
		//		BaseOpt: BaseOpt{
		//			Short: "v",
		//			Full:  "verbose",
		//			// Aliases:     []string{"vv", "vvv"},
		//			Description: "Show this help screen",
		//			// Hidden:      true,
		//			Group: SysMgmtGroup,
		//			owner: &root.Command,
		//			// Action: func(cmd *Command, args []string) (err error) {
		//			// 	if f := FindFlag("verbose", cmd); f != nil {
		//			// 		f.times++
		//			// 		// fmt.Println("verbose++: ", f.times)
		//			// 	}
		//			// 	return
		//			// },
		//
		//			// Action: func(cmd *Command, args []string) (err error) {
		//			// 	if f := FindFlag("verbose", cmd); f != nil {
		//			// 		fmt.Println("verbose++: ", f.times)
		//			// 	}
		//			// 	return
		//			// },
		//		},
		//		DefaultValue: false,
		//		EnvVars:      []string{"VERBOSE"},
		//	}
		//	root.Flags = append(root.Flags, ff)
		//	root.allFlags[SysMgmtGroup]["verbose"] = ff
		//	root.plainLongFlags["verbose"] = root.allFlags[SysMgmtGroup]["verbose"]
		//	// root.plainLongFlags["vvv"] = root.allFlags[SysMgmtGroup]["verbose"]
		//	// root.plainLongFlags["vv"] = root.allFlags[SysMgmtGroup]["verbose"]
		//	root.plainShortFlags["v"] = root.allFlags[SysMgmtGroup]["verbose"]
		//}
		//if _, ok := root.allFlags[SysMgmtGroup]["quiet"]; !ok {
		//	ff := &Flag{
		//		BaseOpt: BaseOpt{
		//			Short:       "q",
		//			Full:        "quiet",
		//			Aliases:     []string{},
		//			Description: "No more screen output.",
		//			// Hidden:      true,
		//			Group: SysMgmtGroup,
		//			owner: &root.Command,
		//		},
		//		DefaultValue: false,
		//		EnvVars:      []string{"QUITE"},
		//	}
		//	root.Flags = append(root.Flags, ff)
		//	root.allFlags[SysMgmtGroup]["quiet"] = ff
		//	root.plainLongFlags["quiet"] = root.allFlags[SysMgmtGroup]["quiet"]
		//	root.plainShortFlags["q"] = root.allFlags[SysMgmtGroup]["quiet"]
		//}
		//if _, ok := root.allFlags[SysMgmtGroup]["debug"]; !ok {
		//	ff := &Flag{
		//		BaseOpt: BaseOpt{
		//			Short:       "D",
		//			Full:        "debug",
		//			Aliases:     []string{},
		//			Description: "Get into debug mode.",
		//			Hidden:      true,
		//			Group:       SysMgmtGroup,
		//			owner:       &root.Command,
		//		},
		//		DefaultValue: false,
		//		EnvVars:      []string{"DEBUG"},
		//	}
		//	root.Flags = append(root.Flags, ff)
		//	root.allFlags[SysMgmtGroup]["debug"] = ff
		//	root.plainLongFlags["debug"] = root.allFlags[SysMgmtGroup]["debug"]
		//	root.plainShortFlags["D"] = root.allFlags[SysMgmtGroup]["debug"]
		//}
		//if _, ok := root.allFlags[SysMgmtGroup]["debug-output"]; !ok {
		//	ff := &Flag{
		//		BaseOpt: BaseOpt{
		//			Short:       "",
		//			Full:        "debug-output",
		//			Aliases:     []string{},
		//			Description: "store the ~~debug outputs into file.",
		//			Hidden:      true,
		//			Group:       SysMgmtGroup,
		//			owner:       &root.Command,
		//		},
		//		DefaultValue: "dbg.log",
		//		EnvVars:      []string{"DEBUG_OUTPUT"},
		//	}
		//	root.Flags = append(root.Flags, ff)
		//	root.allFlags[SysMgmtGroup]["debug-output"] = ff
		//	root.plainLongFlags["debug-output"] = root.allFlags[SysMgmtGroup]["debug-output"]
		//}
		//if _, ok := root.allFlags[SysMgmtGroup]["env"]; !ok {
		//	ff := &Flag{
		//		BaseOpt: BaseOpt{
		//			Short:       "",
		//			Full:        "env",
		//			Aliases:     []string{},
		//			Description: "Dump environment info in `~~debug` mode.",
		//			Hidden:      true,
		//			Group:       SysMgmtGroup,
		//			owner:       &root.Command,
		//		},
		//		DefaultValue: false,
		//	}
		//	root.Flags = append(root.Flags, ff)
		//	root.allFlags[SysMgmtGroup]["env"] = ff
		//	root.plainLongFlags["env"] = root.allFlags[SysMgmtGroup]["env"]
		//}
		//if _, ok := root.allFlags[SysMgmtGroup]["raw"]; !ok {
		//	ff := &Flag{
		//		BaseOpt: BaseOpt{
		//			Short:       "",
		//			Full:        "raw",
		//			Aliases:     []string{},
		//			Description: "Dump the option value in raw mode (with golang data structure, without envvar expanding).",
		//			Hidden:      true,
		//			Group:       SysMgmtGroup,
		//			owner:       &root.Command,
		//		},
		//		DefaultValue: false,
		//	}
		//	root.Flags = append(root.Flags, ff)
		//	root.allFlags[SysMgmtGroup]["raw"] = ff
		//	root.plainLongFlags["raw"] = root.allFlags[SysMgmtGroup]["raw"]
		//}
		//if _, ok := root.allFlags[SysMgmtGroup]["value-type"]; !ok {
		//	ff := &Flag{
		//		BaseOpt: BaseOpt{
		//			Short:       "",
		//			Full:        "value-type",
		//			Aliases:     []string{},
		//			Description: "Dump the option value type.",
		//			Hidden:      true,
		//			Group:       SysMgmtGroup,
		//			owner:       &root.Command,
		//		},
		//		DefaultValue: false,
		//	}
		//	root.Flags = append(root.Flags, ff)
		//	root.allFlags[SysMgmtGroup]["value-type"] = ff
		//	root.plainLongFlags["value-type"] = root.allFlags[SysMgmtGroup]["value-type"]
		//}
		//if _, ok := root.allFlags[SysMgmtGroup]["more"]; !ok {
		//	ff := &Flag{
		//		BaseOpt: BaseOpt{
		//			Short:       "",
		//			Full:        "more",
		//			Aliases:     []string{},
		//			Description: "Dump more info in `~~debug` mode.",
		//			Hidden:      true,
		//			Group:       SysMgmtGroup,
		//			owner:       &root.Command,
		//		},
		//		DefaultValue: false,
		//	}
		//	root.Flags = append(root.Flags, ff)
		//	root.allFlags[SysMgmtGroup]["more"] = ff
		//	root.plainLongFlags["more"] = root.allFlags[SysMgmtGroup]["more"]
		//}
	}
}

func (w *ExecWorker) attachCmdrCommands(root *RootCommand) {
	if w.enableCmdrCommands {
		w._boolFlgAdd(root, "strict-more", "strict mode for `cmdr`.", func(ff *Flag) {
			ff.EnvVars = []string{"STRICT"}
		})
		w._boolFlgAdd(root, "no-env-overrides", "No env var overrides for `cmdr`.", nil)
		w._boolFlgAdd(root, "no-color", "No color output for `cmdr`.", func(ff *Flag) {
			ff.Short = "nc"
			ff.Aliases = []string{"nc"}
			ff.EnvVars = []string{"NOCOLOR", "NO_COLOR"}
		})
		//if _, ok := root.allFlags[SysMgmtGroup]["strict-mode"]; !ok {
		//	ff := &Flag{
		//		BaseOpt: BaseOpt{
		//			Full:        "strict-mode",
		//			Description: "strict mode for `cmdr`.",
		//			Hidden:      true,
		//			Group:       SysMgmtGroup,
		//			owner:       &root.Command,
		//		},
		//		DefaultValue: false,
		//		EnvVars:      []string{"STRICT"},
		//	}
		//	root.Flags = append(root.Flags, ff)
		//	root.allFlags[SysMgmtGroup]["strict-mode"] = ff
		//	root.plainLongFlags["strict-mode"] = ff
		//}
		//if _, ok := root.allFlags[SysMgmtGroup]["no-env-overrides"]; !ok {
		//	ff := &Flag{
		//		BaseOpt: BaseOpt{
		//			Full:        "no-env-overrides",
		//			Description: "No env var overrides for `cmdr`.",
		//			Hidden:      true,
		//			Group:       SysMgmtGroup,
		//			owner:       &root.Command,
		//		},
		//		DefaultValue: false,
		//	}
		//	root.Flags = append(root.Flags, ff)
		//	root.allFlags[SysMgmtGroup]["no-env-overrides"] = ff
		//	root.plainLongFlags["no-env-overrides"] = ff
		//}
		//if _, ok := root.allFlags[SysMgmtGroup]["no-color"]; !ok {
		//	ff := &Flag{
		//		BaseOpt: BaseOpt{
		//			Full:        "no-color",
		//			Description: "No color output for `cmdr`.",
		//			Hidden:      true,
		//			Group:       SysMgmtGroup,
		//			owner:       &root.Command,
		//		},
		//		DefaultValue: false,
		//		EnvVars:      []string{"NOCOLOR", "NO_COLOR"},
		//	}
		//	root.Flags = append(root.Flags, ff)
		//	root.allFlags[SysMgmtGroup]["no-color"] = ff
		//	root.plainLongFlags["no-color"] = ff
		//}
	}
}

func (w *ExecWorker) attachGeneratorsCommands(root *RootCommand) {
	if w.enableGenerateCommands {
		found := false
		for _, sc := range root.SubCommands {
			if sc.Full == "generate" { // generatorCommands.Full {
				found = true
				return
			}
		}
		if !found {
			// root.SubCommands = append(root.SubCommands, generatorCommands)
			w._cmdAdd(root, "generate", "generators for this app.", func(cx1 *Command) {
				cx1.Short = "g"
				cx1.Aliases = []string{"gen"}
				cx1.LongDescription = `
[cmdr] includes multiple generators like:

- linux man page generator
- shell completion script generator
- markdown generator
- more...

			`
				cx1.Examples = `
$ {{.AppName}} gen sh --bash
			generate bash completion script
$ {{.AppName}} gen shell --auto
			generate shell completion script with detecting on current shell environment.
$ {{.AppName}} gen sh
			generate shell completion script with detecting on current shell environment.
$ {{.AppName}} gen man
			generate linux manual (man page)
$ {{.AppName}} gen doc
			generate document, default markdown.
$ {{.AppName}} gen doc --markdown
			generate markdown.
$ {{.AppName}} gen doc --pdf
			generate pdf.
$ {{.AppName}} gen markdown
			generate markdown.
$ {{.AppName}} gen pdf
			generate pdf.
			`
				w._cmdAdd1(cx1, "shell", "generate the bash/zsh auto-completion script or install it.", func(cx *Command) {
					cx.Short = "s"
					cx.Aliases = []string{"sh"}
					cx.Action = genShell
					cx.Hidden = false

					w._stringFlgAdd1(cx, "dir", "the output directory", func(ff *Flag) {
						ff.Short = "d"
						ff.Group = "output"
						ff.DefaultValue = "."
						ff.DefaultValuePlaceholder = "DIR"
						ff.Hidden = false
					})
					w._boolFlgAdd1(cx, "auto", "generate auto completion script to fit for your current env.", func(ff *Flag) {
						ff.Short = "a"
						ff.ToggleGroup = "ShellType"
						ff.DefaultValue = true
						ff.Hidden = false
					})
					w._boolFlgAdd1(cx, "force-bash", "just for --auto", func(ff *Flag) {
						ff.Group = ""
						ff.DefaultValue = true
						ff.Hidden = false
						ff.prerequisites = []string{"auto"}
					})
					w._boolFlgAdd1(cx, "bash", "generate auto completion script for Bash", func(ff *Flag) {
						ff.Short = "b"
						ff.ToggleGroup = "ShellType"
						ff.Hidden = false
					})
					w._boolFlgAdd1(cx, "zsh", "generate auto completion script for Zsh", func(ff *Flag) {
						ff.Short = "z"
						ff.ToggleGroup = "ShellType"
						ff.Hidden = false
					})
					//w._boolFlgAdd1(cx, "ash", "generate auto completion script for Zsh", func(ff *Flag) {
					//	ff.ToggleGroup = "ShellType"
					//ff.Hidden = false
					//})
					//w._boolFlgAdd1(cx, "rbash", "generate auto completion script for Zsh", func(ff *Flag) {
					//	ff.ToggleGroup = "ShellType"
					//ff.Hidden = false
					//})
					//w._boolFlgAdd1(cx, "csh", "generate auto completion script for Zsh", func(ff *Flag) {
					//	ff.ToggleGroup = "ShellType"
					//ff.Hidden = false
					//})
					//w._boolFlgAdd1(cx, "dash", "generate auto completion script for Zsh", func(ff *Flag) {
					//	ff.ToggleGroup = "ShellType"
					//ff.Hidden = false
					//})
					//w._boolFlgAdd1(cx, "ksh", "generate auto completion script for Zsh", func(ff *Flag) {
					//	ff.ToggleGroup = "ShellType"
					//ff.Hidden = false
					//})
					//w._boolFlgAdd1(cx, "tcsh", "generate auto completion script for Zsh", func(ff *Flag) {
					//	ff.ToggleGroup = "ShellType"
					//ff.Hidden = false
					//})
					//w._boolFlgAdd1(cx, "elvish", "generate auto completion script for Fish", func(ff *Flag) {
					//	ff.ToggleGroup = "ShellType"
					//ff.Hidden = false
					//})
					//w._boolFlgAdd1(cx, "fig", "generate auto completion script for Fish", func(ff *Flag) {
					//	ff.ToggleGroup = "ShellType"
					//ff.Hidden = false
					//})
					w._boolFlgAdd1(cx, "fish", "generate auto completion script for Fish", func(ff *Flag) {
						ff.Short = "f"
						ff.ToggleGroup = "ShellType"
						ff.Hidden = false
					})
					w._boolFlgAdd1(cx, "powershell", "generate auto completion script for Powershell", func(ff *Flag) {
						ff.Short = "p"
						ff.ToggleGroup = "ShellType"
						ff.Hidden = false
					})
				})
				w._cmdAdd1(cx1, "manual", "generate linux man page.", func(cx *Command) {
					cx.Short = "m"
					cx.Aliases = []string{"man"}
					cx.Action = genManual
					cx.Hidden = false

					w._stringFlgAdd1(cx, "dir", "the output directory", func(ff *Flag) {
						ff.Short = "d"
						ff.Group = "output"
						ff.DefaultValue = "./man1"
						ff.DefaultValuePlaceholder = "DIR"
						ff.Hidden = false
					})
				})
				w._cmdAdd1(cx1, "doc", "generate a markdown document, or: pdf/TeX/...", func(cx *Command) {
					cx.Short = "d"
					cx.Aliases = []string{"pdf", "docx", "tex", "markdown"}
					cx.Action = genDoc
					cx.Hidden = false

					w._stringFlgAdd1(cx, "dir", "the output directory", func(ff *Flag) {
						ff.Short = "d"
						ff.Group = "output"
						ff.DefaultValue = "./docs"
						ff.DefaultValuePlaceholder = "DIR"
						ff.Hidden = false
					})
					w._boolFlgAdd1(cx, "markdown", "to generate a markdown file", func(ff *Flag) {
						ff.Short = "md"
						ff.Aliases = []string{"mkd", "m"}
						ff.ToggleGroup = "DocType"
						ff.DefaultValue = true
						ff.Hidden = false
					})
					w._boolFlgAdd1(cx, "pdf", "to generate a PDF file", func(ff *Flag) {
						ff.Short = "p"
						ff.ToggleGroup = "DocType"
						ff.Hidden = false
					})
					w._boolFlgAdd1(cx, "docx", "to generate a Word (.docx) file", func(ff *Flag) {
						ff.Aliases = []string{"doc"}
						ff.ToggleGroup = "DocType"
						ff.Hidden = false
					})
					w._boolFlgAdd1(cx, "tex", "to generate a LaTeX file", func(ff *Flag) {
						ff.Short = "t"
						ff.ToggleGroup = "DocType"
						ff.Hidden = false
					})
				})
			})

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

	bre := regexp.MustCompile("`(.+)`")

	for _, flg := range cmd.Flags {
		flg.owner = cmd

		if len(flg.ToggleGroup) > 0 {
			if len(flg.Group) == 0 {
				flg.Group = flg.ToggleGroup
			}
			tgs[flg.ToggleGroup] = true
		}

		if b := bre.Find([]byte(flg.Description)); len(flg.DefaultValuePlaceholder) == 0 && len(b) > 2 {
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

		w.rxxtOptions.Set(w.backtraceCmdNames(cx, false), nil)
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
			ferr("\nNOTE: flag alias name '%v' has been used. (command: %v)", sz, w.backtraceCmdNames(cmd, false))
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
			ferr("\nNOTE: flag char '%v' has been used. (command: %v)", flg.Short, w.backtraceCmdNames(cmd, false))
		} else {
			singleFlagNames[flg.Short] = true
		}
	}
	if len(flg.Full) != 0 {
		if _, ok := stringFlagNames[flg.Full]; ok {
			ferr("\nNOTE: flag '%v' has been used. (command: %v)", flg.Full, w.backtraceCmdNames(cmd, false))
		} else {
			stringFlagNames[flg.Full] = true
		}
	}
	if len(flg.Short) == 0 && len(flg.Full) == 0 && len(flg.Name) != 0 {
		if _, ok := stringFlagNames[flg.Name]; ok {
			ferr("\nNOTE: flag '%v' has been used. (command: %v)", flg.Name, w.backtraceCmdNames(cmd, false))
		} else {
			stringFlagNames[flg.Name] = true
		}
	}
}

// helpSystemPrint NOT YET, to-do
func (w *ExecWorker) helpSystemPrint(cmd *Command, args []string) (err error) {
	var matchedPrecededList, matchedList map[string]*Command
	var matchedCmd *Command
	matchedCmd, matchedPrecededList, matchedList, err = w.lookupForHelpSystem(cmd, args)
	if err == nil || err == ErrShouldBeStopException {
		var x strings.Builder

		for _, c := range matchedPrecededList {
			_visitCommandTitles(c, func(c *Command, title string) (stopNow bool) {
				x.WriteString(fmt.Sprintf("%v\t%v\n", title, c.Description))
				return
			})
		}

		if matchedCmd != nil && len(matchedPrecededList) == 1 {
			// print, ret
			fp("%v", x.String())
			return
		}

		for _, c := range matchedList {
			_visitCommandTitles(c, func(c *Command, title string) (stopNow bool) {
				x.WriteString(fmt.Sprintf("%v\t%v\n", title, c.Description))
				return
			})
		}

		directive := 4
		x.WriteString(fmt.Sprintf(":%d\n", directive))

		fp("%v", x.String())
	}
	return
}

// lookupForHelpSystem NOT YET, to-do
func (w *ExecWorker) lookupForHelpSystem(cmdHelp *Command, args []string) (matchedCmd *Command, matchedPrecededList, matchedList map[string]*Command, err error) {
	err = ErrShouldBeStopException
	matchedPrecededList = make(map[string]*Command)
	matchedList = make(map[string]*Command)
	cmd := &cmdHelp.root.Command

	for i := 0; i < len(args); i++ {
		ttl := args[i]

		if len(ttl) > 0 && strings.Contains(w.switchCharset, ttl[0:1]) {
			// flags
			continue
		}

		// sub-commands
		for _, c := range cmd.SubCommands {
			var matchedPrecededList1, matchedList1 map[string]*Command
			if exact, ok := _matchCommandTitle(c, ttl, matchedPrecededList1, matchedList1); exact && ok {
				cmd = c
				if i == len(args)-1 {
					matchedCmd = c
					for _, c1 := range c.SubCommands {
						matchedList[c1.Full] = c1
					}
				}
				break
			} else {
				return
			}
		}
	}
	return
}

func _matchCommandTitle(c *Command, titleChecking string, matchedPrecededList, matchedList map[string]*Command) (exact, ok bool) {
	_visitCommandTitles(c, func(c *Command, title string) (stopNow bool) {
		if exact, ok = _matchTitle(c, title, titleChecking, matchedPrecededList, matchedList); ok {
			stopNow = true
		}
		return
	})
	return
}

var noPartialMatching = true

func _matchTitle(c *Command, source, title string, matchedPrecededList, matchedList map[string]*Command) (exact, ok bool) {
	if title == source {
		matchedPrecededList[c.Full] = c
		ok = true
	} else if strings.HasPrefix(source, title) {
		matchedPrecededList[c.Full] = c
		ok = true
	} else if noPartialMatching == false && strings.Contains(source, title) {
		matchedList[c.Full] = c
		ok = true
	}
	return
}

func _visitCommandTitles(c *Command, fn func(c *Command, title string) (stopNow bool)) (ok bool) {
	if c.Full != "" && fn(c, c.Full) {
		return
	}
	if c.Short != "" && fn(c, c.Short) {
		return
	}
	for _, t := range c.Aliases {
		if t != "" && fn(c, t) {
			return
		}
	}
	return true
}

func _visitFlagTitles(c *Flag, fn func(f *Flag, title string) (stopNow bool)) (ok bool) {
	if c.Full != "" && fn(c, c.Full) {
		return
	}
	if c.Short != "" && fn(c, c.Short) {
		return
	}
	for _, t := range c.Aliases {
		if t != "" && fn(c, t) {
			return
		}
	}
	return true
}

func (w *ExecWorker) _buildCrossRefsForCommand(cx, cmd *Command, singleCmdNames, stringCmdNames map[string]bool) {
	w.forCommandNames(cx, cmd, singleCmdNames, stringCmdNames)

	for _, sz := range cx.Aliases {
		if len(sz) != 0 {
			if _, ok := stringCmdNames[sz]; ok {
				ferr("\nNOTE: command alias name '%v' has been used. (command: %v)", sz, w.backtraceCmdNames(cmd, false))
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
			ferr("\nNOTE: command char '%v' has been used. (command: %v)", cx.Short, w.backtraceCmdNames(cmd, false))
		} else {
			singleCmdNames[cx.Short] = true
		}
	}
	if len(cx.Full) != 0 {
		if _, ok := stringCmdNames[cx.Full]; ok {
			ferr("\nNOTE: command '%v' has been used. (command: %v)", cx.Full, w.backtraceCmdNames(cmd, false))
		} else {
			stringCmdNames[cx.Full] = true
		}
	}
	if len(cx.Short) == 0 && len(cx.Full) == 0 && len(cx.Name) != 0 {
		if _, ok := stringCmdNames[cx.Name]; ok {
			ferr("\nNOTE: command '%v' has been used. (command: %v)", cx.Name, w.backtraceCmdNames(cmd, false))
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
			w.rxxtOptions.Set(w.backtraceCmdNames(cmd, false)+"."+tg, f.Full)
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

// backtraceCmdNames returns the sequences of a sub-command from
// top-level.
//
// - if verboseLast = false, got 'microservices.tags.list' for sub-cmd microservice/tags/list.
//
// - if verboseLast = true,  got 'microservices.tags.[ls|list|l|lst|dir]'.
//
// - at root command, it returns 'appName' or '' when verboseLast is true.
func (w *ExecWorker) backtraceCmdNames(cmd *Command, verboseLast bool) (str string) {
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
