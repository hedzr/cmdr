/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"reflect"
	"strconv"
	"strings"
	"syscall"
)

//

//

// Exec is main entry of `cmdr`.
func Exec(rootCmd *RootCommand) (err error) {
	err = InternalExecFor(rootCmd, os.Args)
	return
}

// ExecWith is main entry of `cmdr`.
func ExecWith(rootCmd *RootCommand, beforeXrefBuildingX, afterXrefBuiltX HookXrefFunc) (err error) {
	if beforeXrefBuildingX != nil {
		beforeXrefBuilding = append(beforeXrefBuilding, beforeXrefBuildingX)
	}
	if afterXrefBuiltX != nil {
		afterXrefBuilt = append(afterXrefBuilt, afterXrefBuiltX)
	}
	err = InternalExecFor(rootCmd, os.Args)
	return
}

// InternalExecFor is an internal helper, esp for debugging
func InternalExecFor(rootCmd *RootCommand, args []string) (err error) {
	var (
		pkg       = new(ptpkg)
		goCommand = &rootCmd.Command
		stop      bool
		// helpFlag = rootCmd.allFlags[UnsortedGroup]["help"]
	)

	if rootCommand == nil {
		setRootCommand(rootCmd)
	}

	defer func() {
		_ = rootCmd.ow.Flush()
		_ = rootCmd.oerr.Flush()
	}()

	if err = preprocess(rootCmd, args); err != nil {
		return
	}

	for pkg.i = 1; pkg.i < len(args); pkg.i++ {
		pkg.a = args[pkg.i]
		pkg.assigned = false
		pkg.short = false
		pkg.savedFn = ""
		pkg.savedVal = ""

		// --debug: long opt
		// -D:      short opt
		// -nv:     double chars short opt
		// ~~debug: long opt without opt-entry prefix.
		// ~D:      short opt without opt-entry prefix.
		// -abc:    the combined short opts
		// -nvabc, -abnvc: a,b,c,nv the four short opts, if no -n & -v defined.
		// --name=consul, --name consul, --name=consul: opt with a string, int, string slice argument
		// -nconsul, -n consul, -n=consul: opt with an argument.
		//  - -nconsul is not good format, but it could get somewhat works.
		//  - -n'consul', -n"consul" could works too.
		// -t3: opt with an argument.
		stop, err = xxTestCmd(pkg, &goCommand, rootCmd, args)
		if stop {
			break
		}
	}

	err = afterInternalExec(pkg, rootCmd, goCommand, args)

	return
}

func xxTestCmd(pkg *ptpkg, goCommand **Command, rootCmd *RootCommand, args []string) (stop bool, err error) {
	if pkg.a[0] == '-' || pkg.a[0] == '/' || pkg.a[0] == '~' {
		if len(pkg.a) == 1 {
			pkg.needHelp = true
			pkg.needFlagsHelp = true
			return
		}

		// flag
		if stop, err = flagsPrepare(pkg, goCommand, args); stop || err != nil {
			return
		}
		// if len(pkg.a) > 1 && (pkg.a[1] == '-' || pkg.a[1] == '~') {
		// 	if len(pkg.a) == 2 {
		// 		// disableParser = true // '--': ignore the following args
		// 		stop = true
		// 		return
		// 	}
		//
		// 	// long flag
		// 	pkg.fn = pkg.a[2:]
		// 	findValueAttached(pkg, &pkg.fn)
		// } else {
		// 	pkg.suffix = pkg.a[len(pkg.a)-1]
		// 	if pkg.suffix == '+' || pkg.suffix == '-' {
		// 		pkg.a = pkg.a[0 : len(pkg.a)-1]
		// 	} else {
		// 		pkg.suffix = 0
		// 	}
		//
		// 	if i := matchShortFlag(pkg, *goCommand, pkg.a); i >= 0 {
		// 		pkg.fn = pkg.a[1:i]
		// 		pkg.savedFn = pkg.a[i:]
		// 	} else {
		// 		pkg.fn = pkg.a[1:2]
		// 		pkg.savedFn = pkg.a[2:]
		// 	}
		// 	pkg.short = true
		// 	findValueAttached(pkg, &pkg.savedFn)
		// }

		// fn + val
		// fn: short,
		// fn: long
		// fn: short||val: such as '-t3'
		// fn: long=val, long='val', long="val", long val, long 'val', long "val"
		// fn: longval, long'val', long"val"

		pkg.savedGoCommand = *goCommand
		cc := *goCommand
		if stop, err = flagsMatching(pkg, cc, goCommand, args); stop || err != nil {
			return
		}

	} else {
		if stop, err = cmdMatching(pkg, goCommand, args); stop || err != nil {
			return
		}
		// // command, files
		// if cmd, ok := (*goCommand).plainCmds[pkg.a]; ok {
		// 	cmd.strHit = pkg.a
		// 	*goCommand = cmd
		// 	// logrus.Debugf("-- command '%v' hit, go ahead...", cmd.GetTitleName())
		// 	if cmd.PreAction != nil {
		// 		if err = cmd.PreAction(*goCommand, getArgs(pkg, args)); err == ErrShouldBeStopException {
		// 			return false, nil
		// 		}
		// 	}
		// } else {
		// 	if (*goCommand).Action != nil && len((*goCommand).SubCommands) == 0 {
		// 		// the args remained are files, not sub-commands.
		// 		pkg.i--
		// 		stop = true
		// 		return
		// 	}
		//
		// 	ferr("Unknown command: %v", pkg.a)
		// 	pkg.unknownCmds = append(pkg.unknownCmds, pkg.a)
		// }
	}
	return
}

func cmdMatching(pkg *ptpkg, goCommand **Command, args []string) (stop bool, err error) {
	// command, files
	if cmd, ok := (*goCommand).plainCmds[pkg.a]; ok {
		cmd.strHit = pkg.a
		*goCommand = cmd
		// logrus.Debugf("-- command '%v' hit, go ahead...", cmd.GetTitleName())
		stop, err = cmdMatched(pkg, *goCommand, args)
	} else {
		if (*goCommand).Action != nil && len((*goCommand).SubCommands) == 0 {
			// the args remained are files, not sub-commands.
			pkg.i--
			stop = true
			return
		}

		pkg.unknownCmds = append(pkg.unknownCmds, pkg.a)
		unknownCommand(pkg, *goCommand, args)
	}
	return
}

func cmdMatched(pkg *ptpkg, goCommand *Command, args []string) (stop bool, err error) {
	if goCommand.PreAction != nil {
		if err = goCommand.PreAction(goCommand, getArgs(pkg, args)); err == ErrShouldBeStopException {
			return false, nil
		}
	}
	return
}

func flagsPrepare(pkg *ptpkg, goCommand **Command, args []string) (stop bool, err error) {
	if len(pkg.a) > 1 && (pkg.a[1] == '-' || pkg.a[1] == '~') {
		if len(pkg.a) == 2 {
			// disableParser = true // '--': ignore the following args
			stop = true
			return
		}

		// long flag
		pkg.fn = pkg.a[2:]
		findValueAttached(pkg, &pkg.fn)
	} else {
		pkg.suffix = pkg.a[len(pkg.a)-1]
		if pkg.suffix == '+' || pkg.suffix == '-' {
			pkg.a = pkg.a[0 : len(pkg.a)-1]
		} else {
			pkg.suffix = 0
		}

		if i := matchShortFlag(pkg, *goCommand, pkg.a); i >= 0 {
			pkg.fn = pkg.a[1:i]
			pkg.savedFn = pkg.a[i:]
		} else {
			pkg.fn = pkg.a[1:2]
			pkg.savedFn = pkg.a[2:]
		}
		pkg.short = true
		findValueAttached(pkg, &pkg.savedFn)
	}
	return
}

func flagsMatching(pkg *ptpkg, cc *Command, goCommand **Command, args []string) (stop bool, err error) {
	var ok, upLevel bool
GO_UP:
	pkg.found = false
	if pkg.short {
		pkg.flg, ok = cc.plainShortFlags[pkg.fn]
	} else {
		ok = matchForLongFlags(pkg.fn, cc, pkg)
	}

	if ok {
		if upLevel, stop, err = flagsMatched(pkg, *goCommand, args); stop || err != nil {
			return
		}
		if upLevel {
			goto GO_UP
		}
	} else {
		if cc.owner != nil {
			// match the flag within parent's flags set.
			cc = cc.owner
			goto GO_UP
		}
		if !pkg.assigned && pkg.short {
			// try matching 2-chars short opt
			if len(pkg.savedFn) > 0 {
				fnf := pkg.fn + pkg.savedFn
				pkg.fn = fnf[0:2]
				pkg.savedFn = fnf[2:]
				*goCommand = pkg.savedGoCommand
				goto GO_UP
			}
		}

		pkg.unknownFlags = append(pkg.unknownFlags, pkg.a)
		unknownFlag(pkg, *goCommand, args)
	}
	return
}

func flagsMatched(pkg *ptpkg, goCommand *Command, args []string) (upLevel, stop bool, err error) {
	if err = tryExtractingValue(pkg, args); err != nil {
		stop = true
		return
	}

	if pkg.found {
		// if !GetBoolP(getPrefix(), "quiet") {
		// 	logrus.Debugf("-- flag '%v' hit, go ahead...", pkg.flg.GetTitleName())
		// }
		if pkg.flg.Action != nil {
			if err = pkg.flg.Action(goCommand, getArgs(pkg, args)); err == ErrShouldBeStopException {
				stop = true
				err = nil
				return
			}
		}
		if isBool(pkg.flg.DefaultValue) || isNil1(pkg.flg.DefaultValue) {
			toggleGroup(pkg)
		}

		if !pkg.assigned {
			if len(pkg.savedFn) > 0 && len(pkg.savedVal) == 0 {
				pkg.fn = pkg.savedFn[0:1]
				pkg.savedFn = pkg.savedFn[1:]
				// goto GO_UP
				upLevel = true
			}
		}
	}
	return
}

func unknownCommand(pkg *ptpkg, cmd *Command, args []string) {
	ferr("\n\x1b[%dmUnknown command:\x1b[0m %v", BgBoldOrBright, pkg.a)
	unknownCommandDetector(pkg, cmd, args)
}

func unknownFlag(pkg *ptpkg, cmd *Command, args []string) {
	ferr("\n\x1b[%dmUnknown flag:\x1b[0m %v", BgBoldOrBright, pkg.a)
	unknownFlagDetector(pkg, cmd, args)
}

func unknownCommandDetector(pkg *ptpkg, cmd *Command, args []string) {
	sndSrc := soundex(pkg.a)
	ever := false
	for k := range cmd.plainCmds {
		snd := soundex(k)
		if sndSrc == snd {
			ferr("  - do you mean: %v", k)
			ever = true
			// } else {
			// 	ferr("  . %v -> %v: --%v -> %v", pkg.a, sndSrc, k, snd)
		}
	}
	if !ever && cmd.HasParent() {
		unknownCommandDetector(pkg, cmd.GetOwner(), args)
	}
}

func unknownFlagDetector(pkg *ptpkg, cmd *Command, args []string) {
	sndSrc := soundex(pkg.a)
	if !pkg.short {
		ever := false
		for k := range cmd.plainLongFlags {
			snd := soundex(k)
			if sndSrc == snd {
				ferr("  - do you mean: --%v", k)
				ever = true
				// } else {
				// 	ferr("  . %v -> %v: --%v -> %v", pkg.a, sndSrc, k, snd)
			}
		}
		if !ever && cmd.HasParent() {
			unknownFlagDetector(pkg, cmd.GetOwner(), args)
		}
	}
}

func matchForLongFlags(fn string, cc *Command, pkg *ptpkg) (ok bool) {
	var ln = len(fn)
	for ; ln > 1; ln-- {
		fn = pkg.fn[0:ln]
		pkg.flg, ok = cc.plainLongFlags[fn]
		if ok {
			if ln < len(pkg.fn) {
				pkg.val = pkg.fn[ln:]
				pkg.fn = fn
				pkg.assigned = true
			}
			break
		}
	}
	return
}

func preprocess(rootCmd *RootCommand, args []string) (err error) {
	for _, x := range beforeXrefBuilding {
		x(rootCmd, args)
	}

	if err = buildXref(rootCmd); err != nil {
		return
	}

	if err = rxxtOptions.buildAutomaticEnv(rootCmd); err != nil {
		return
	}

	for _, x := range afterXrefBuilt {
		x(rootCmd, args)
	}

	return
}

func afterInternalExec(pkg *ptpkg, rootCmd *RootCommand, goCommand *Command, args []string) (err error) {
	if !pkg.needHelp {
		pkg.needHelp = GetBoolP(getPrefix(), "help")
	}

	if !pkg.needHelp && len(pkg.unknownCmds) == 0 && len(pkg.unknownFlags) == 0 {
		if goCommand.Action != nil {
			args := getArgs(pkg, args)

			if goCommand != &rootCmd.Command {
				if rootCmd.PostAction != nil {
					defer rootCmd.PostAction(goCommand, args)
				}
				if rootCmd.PreAction != nil {
					if err = rootCmd.PreAction(goCommand, getArgs(pkg, args)); err == ErrShouldBeStopException {
						return nil
					}
				}
			}

			if goCommand.PostAction != nil {
				defer goCommand.PostAction(goCommand, args)
			}

			if err = goCommand.Action(goCommand, args); err == ErrShouldBeStopException {
				return nil
			}

			return
		}
	}

	// if GetIntP(getPrefix(), "help-zsh") > 0 || GetBoolP(getPrefix(), "help-bash") {
	// 	if len(goCommand.SubCommands) == 0 && !pkg.needFlagsHelp {
	// 		// pkg.needFlagsHelp = true
	// 	}
	// }

	printHelp(goCommand, pkg.needFlagsHelp)
	return
}

func buildXref(rootCmd *RootCommand) (err error) {
	// build xref for root command and its all sub-commands and flags
	// and build the default values
	buildRootCrossRefs(rootCmd)

	if !doNotLoadingConfigFiles {
		// pre-detects for `--config xxx`, `--config=xxx`, `--configxxx`
		if err = parsePredefinedLocation(); err != nil {
			return
		}
		// if ix, str, yes := partialContains(os.Args, "--config"); yes {
		// 	var location string
		// 	if i := strings.Index(str, "="); i > 0 {
		// 		location = str[i+1:]
		// 	} else if len(str) > 8 {
		// 		location = str[8:]
		// 	} else if ix+1 < len(os.Args) {
		// 		location = os.Args[ix+1]
		// 	}
		//
		// 	location = trimQuotes(location)
		//
		// 	if len(location) > 0 && FileExists(location) {
		// 		if yes, err = IsDirectory(location); yes {
		// 			if FileExists(location + "/conf.d") {
		// 				SetPredefinedLocations([]string{location + "/%s.yml"})
		// 			} else {
		// 				SetPredefinedLocations([]string{location + "/%s/%s.yml"})
		// 			}
		// 		} else if yes, err = IsRegularFile(location); yes {
		// 			SetPredefinedLocations([]string{location})
		// 		}
		// 	}
		// }

		// and now, loading the external configuration files
		err = loadFromPredefinedLocation(rootCmd)
		// for _, s := range getExpandedPredefinedLocations() {
		// 	fn := s
		// 	switch strings.Count(fn, "%s") {
		// 	case 2:
		// 		fn = fmt.Sprintf(s, rootCmd.AppName, rootCmd.AppName)
		// 	case 1:
		// 		fn = fmt.Sprintf(s, rootCmd.AppName)
		// 	}
		//
		// 	if FileExists(fn) {
		// 		err = rxxtOptions.LoadConfigFile(fn)
		// 		if err != nil {
		// 			return
		// 		}
		// 		break
		// 	}
		// }
	}
	return
}

func parsePredefinedLocation() (err error) {
	// pre-detects for `--config xxx`, `--config=xxx`, `--configxxx`
	if ix, str, yes := partialContains(os.Args, "--config"); yes {
		var location string
		if i := strings.Index(str, "="); i > 0 {
			location = str[i+1:]
		} else if len(str) > 8 {
			location = str[8:]
		} else if ix+1 < len(os.Args) {
			location = os.Args[ix+1]
		}

		location = trimQuotes(location)

		if len(location) > 0 && FileExists(location) {
			if yes, err = IsDirectory(location); yes {
				if FileExists(location + "/conf.d") {
					SetPredefinedLocations([]string{location + "/%s.yml"})
				} else {
					SetPredefinedLocations([]string{location + "/%s/%s.yml"})
				}
			} else if yes, err = IsRegularFile(location); yes {
				SetPredefinedLocations([]string{location})
			}
		}
	}
	return
}

func loadFromPredefinedLocation(rootCmd *RootCommand) (err error) {
	// and now, loading the external configuration files
	for _, s := range getExpandedPredefinedLocations() {
		fn := s
		switch strings.Count(fn, "%s") {
		case 2:
			fn = fmt.Sprintf(s, rootCmd.AppName, rootCmd.AppName)
		case 1:
			fn = fmt.Sprintf(s, rootCmd.AppName)
		}

		if FileExists(fn) {
			err = rxxtOptions.LoadConfigFile(fn)
			if err != nil {
				return
			}
			break
		}
	}
	return
}

// AddOnBeforeXrefBuilding add hook func
func AddOnBeforeXrefBuilding(cb HookXrefFunc) {
	beforeXrefBuilding = append(beforeXrefBuilding, cb)
}

// AddOnAfterXrefBuilt add hook func
func AddOnAfterXrefBuilt(cb HookXrefFunc) {
	afterXrefBuilt = append(afterXrefBuilt, cb)
}

func getPrefix() string {
	return strings.Join(RxxtPrefix, ".")
}

func setRootCommand(rootCmd *RootCommand) {
	rootCommand = rootCmd

	rootCommand.ow = defaultStdout
	rootCommand.oerr = defaultStderr
}

func getArgs(pkg *ptpkg, args []string) []string {
	var a []string
	if pkg.i+1 < len(args) {
		a = args[pkg.i+1:]
	}
	return a
}

// func isTypeInt(kind reflect.Kind) bool {
// 	switch kind {
// 	case reflect.Int:
// 	case reflect.Int8:
// 	case reflect.Int16:
// 	case reflect.Int32:
// 	case reflect.Int64:
// 	case reflect.Uint:
// 	case reflect.Uint8:
// 	case reflect.Uint16:
// 	case reflect.Uint32:
// 	case reflect.Uint64:
// 	default:
// 		return false
// 	}
// 	return true
// }

func isTypeUint(kind reflect.Kind) bool {
	switch kind {
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
	default:
		return false
	}
	return true
}

func isTypeSInt(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	default:
		return false
	}
	return true
}

func isBool(v interface{}) bool {
	_, ok := v.(bool)
	return ok
}

func isNil1(v interface{}) bool {
	return v == nil
}

type ptpkg struct {
	assigned          bool
	found             bool
	short             bool
	fn, val           string
	savedFn, savedVal string
	i                 int
	a                 string
	flg               *Flag
	savedGoCommand    *Command
	needHelp          bool
	needFlagsHelp     bool
	suffix            uint8
	unknownCmds       []string
	unknownFlags      []string
}

func toggleGroup(pkg *ptpkg) {
	tg := pkg.flg.ToggleGroup
	for _, f := range pkg.flg.owner.Flags {
		if f.ToggleGroup == tg && (isBool(f.DefaultValue) || isNil1(f.DefaultValue)) {
			rxxtOptions.SetNx(backtraceFlagNames(pkg.flg), false)
		}
	}
}

func findValueAttached(pkg *ptpkg, fn *string) {
	if strings.Contains(*fn, "=") {
		aa := strings.Split(*fn, "=")
		*fn = aa[0]
		pkg.val = trimQuotes(aa[1])
		pkg.assigned = true
	} else {
		splitQuotedValueIfNecessary(pkg, fn)
	}
}

func splitQuotedValueIfNecessary(pkg *ptpkg, fn *string) {
	if pos := strings.Index(*fn, "'"); pos >= 0 {
		pkg.val = trimQuotes((*fn)[pos:])
		*fn = (*fn)[0:pos]
		pkg.assigned = true
	} else if pos := strings.Index(*fn, "\""); pos >= 0 {
		pkg.val = trimQuotes((*fn)[pos:])
		*fn = (*fn)[0:pos]
		pkg.assigned = true
		// } else {
		// --xVALUE need to be parsed.
	}
}

func matchShortFlag(pkg *ptpkg, goCommand *Command, a string) (i int) {
	for i = len(a); i > 1; i-- {
		fn := a[1:i]
		if _, ok := goCommand.plainShortFlags[fn]; ok {
			return
		}
	}
	return -1
}

func tryExtractingValue(pkg *ptpkg, args []string) (err error) {
	if _, ok := pkg.flg.DefaultValue.(bool); ok {
		return tryExtractingBoolValue(pkg)
	}

	vv := reflect.ValueOf(pkg.flg.DefaultValue)
	kind := vv.Kind()
	switch kind {
	case reflect.String:
		if err = processTypeString(pkg, args); err != nil {
			return
		}

	case reflect.Slice:
		err = tryExtractingSliceValue(pkg, args)

	default:
		err = tryExtractingOthers(pkg, args, kind)
	}

	return
}

func tryExtractingOthers(pkg *ptpkg, args []string, kind reflect.Kind) (err error) {
	if isTypeSInt(kind) {
		if err = processTypeInt(pkg, args); err != nil {
			return
		}
	} else if isTypeUint(kind) {
		if err = processTypeUint(pkg, args); err != nil {
			return
		}
	} else {
		ferr("Unacceptable default value kind=%v", kind)
	}
	return
}

func tryExtractingSliceValue(pkg *ptpkg, args []string) (err error) {
	typ := reflect.TypeOf(pkg.flg.DefaultValue).Elem()
	if typ.Kind() == reflect.String {
		if err = processTypeStringSlice(pkg, args); err != nil {
			return
		}
	} else if isTypeSInt(typ.Kind()) {
		if err = processTypeIntSlice(pkg, args); err != nil {
			return
		}
	} else if isTypeUint(typ.Kind()) {
		if err = processTypeUintSlice(pkg, args); err != nil {
			return
		}
	}
	return
}

func tryExtractingBoolValue(pkg *ptpkg) (err error) {
	// bool flag, -D+, -D-

	if pkg.suffix == '+' {
		pkg.flg.DefaultValue = true
	} else if pkg.suffix == '-' {
		pkg.flg.DefaultValue = false
	} else {
		pkg.flg.DefaultValue = true
	}

	if pkg.a[0] == '~' {
		rxxtOptions.SetNx(backtraceFlagNames(pkg.flg), pkg.flg.DefaultValue)
	} else {
		rxxtOptions.Set(backtraceFlagNames(pkg.flg), pkg.flg.DefaultValue)
	}
	pkg.found = true
	return
}

func preprocessPkg(pkg *ptpkg, args []string) (err error) {
	if !pkg.assigned {
		if len(pkg.savedVal) > 0 {
			pkg.val = pkg.savedVal
			pkg.savedVal = ""
		} else if len(pkg.savedFn) > 0 {
			pkg.val = pkg.savedFn
			pkg.savedFn = ""
		} else {
			if pkg.i < len(args)-1 && args[pkg.i+1][0] != '-' && args[pkg.i+1][0] != '~' {
				pkg.i++
				pkg.val = args[pkg.i]
			} else {
				if len(pkg.flg.ExternalTool) > 0 {
					err = processExternalTool(pkg)
				} else if GetStrictMode() {
					err = fmt.Errorf("unexpect end of command line [i=%v,args=(%v)], need more args for %v", pkg.i, args, pkg)
					return
				}
			}
		}
		pkg.assigned = true
	}
	return
}

func processExternalTool(pkg *ptpkg) (err error) {
	switch pkg.flg.ExternalTool {
	case ExternalToolPasswordInput:
		fmt.Print("Password: ")
		var bytePassword []byte
		if InTesting() {
			bytePassword = []byte("demo")
		} else {
			if bytePassword, err = terminal.ReadPassword(int(syscall.Stdin)); err != nil {
				fmt.Println() // it's necessary to add a new line after user's input
				return
			}
			fmt.Println() // it's necessary to add a new line after user's input
		}
		pkg.val = string(bytePassword)

	default:
		editor := os.Getenv(pkg.flg.ExternalTool)
		if len(editor) == 0 {
			editor = DefaultEditor
		}
		var content []byte
		if InTesting() {
			content = []byte("demo")
		} else {
			if content, err = LaunchEditor(editor); err != nil {
				return
			}
		}
		pkg.val = string(content)
	}
	return
}

func processTypeInt(pkg *ptpkg, args []string) (err error) {
	if err = preprocessPkg(pkg, args); err != nil {
		return
	}

	v, err := strconv.ParseInt(pkg.val, 10, 64)
	if err != nil {
		ferr("wrong number: flag=%v, number=%v", pkg.fn, pkg.val)
	}

	if pkg.a[0] == '~' {
		rxxtOptions.SetNx(backtraceFlagNames(pkg.flg), v)
	} else {
		rxxtOptions.Set(backtraceFlagNames(pkg.flg), v)
	}
	pkg.found = true
	return
}

func processTypeUint(pkg *ptpkg, args []string) (err error) {
	if err = preprocessPkg(pkg, args); err != nil {
		return
	}

	v, err := strconv.ParseUint(pkg.val, 10, 64)
	if err != nil {
		ferr("wrong number: flag=%v, number=%v", pkg.fn, pkg.val)
	}

	if pkg.a[0] == '~' {
		rxxtOptions.SetNx(backtraceFlagNames(pkg.flg), v)
	} else {
		rxxtOptions.Set(backtraceFlagNames(pkg.flg), v)
	}
	pkg.found = true
	return
}

func processTypeString(pkg *ptpkg, args []string) (err error) {
	if err = preprocessPkg(pkg, args); err != nil {
		return
	}

	if pkg.a[0] == '~' {
		rxxtOptions.SetNx(backtraceFlagNames(pkg.flg), pkg.val)
	} else {
		rxxtOptions.Set(backtraceFlagNames(pkg.flg), pkg.val)
	}
	pkg.found = true
	return
}

func processTypeStringSlice(pkg *ptpkg, args []string) (err error) {
	if err = preprocessPkg(pkg, args); err != nil {
		return
	}

	if pkg.a[0] == '~' {
		rxxtOptions.SetNx(backtraceFlagNames(pkg.flg), strings.Split(pkg.val, ","))
	} else {
		rxxtOptions.Set(backtraceFlagNames(pkg.flg), strings.Split(pkg.val, ","))
	}
	pkg.found = true
	return
}

func processTypeIntSlice(pkg *ptpkg, args []string) (err error) {
	if err = preprocessPkg(pkg, args); err != nil {
		return
	}

	valary := make([]int64, 0)
	for _, x := range strings.Split(pkg.val, ",") {
		if xi, err := strconv.ParseInt(x, 10, 64); err == nil {
			valary = append(valary, xi)
		}
	}

	if pkg.a[0] == '~' {
		rxxtOptions.SetNx(backtraceFlagNames(pkg.flg), valary)
	} else {
		rxxtOptions.Set(backtraceFlagNames(pkg.flg), valary)
	}
	pkg.found = true
	return
}

func processTypeUintSlice(pkg *ptpkg, args []string) (err error) {
	if err = preprocessPkg(pkg, args); err != nil {
		return
	}

	valary := make([]uint64, 0)
	for _, x := range strings.Split(pkg.val, ",") {
		if xi, err := strconv.ParseUint(x, 10, 64); err == nil {
			valary = append(valary, xi)
		}
	}

	if pkg.a[0] == '~' {
		rxxtOptions.SetNx(backtraceFlagNames(pkg.flg), valary)
	} else {
		rxxtOptions.Set(backtraceFlagNames(pkg.flg), valary)
	}
	pkg.found = true
	return
}
