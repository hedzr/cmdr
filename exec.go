/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"fmt"
	"github.com/hedzr/cmdr/conf"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"
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
		matched   bool
		// helpFlag = rootCmd.allFlags[UnsortedGroup]["help"]
	)

	if rootCommand == nil {
		setRootCommand(rootCmd)
	}

	defer func() {
		_ = rootCmd.ow.Flush()
		_ = rootCmd.oerr.Flush()
	}()

	err = preprocess(rootCmd, args)

	if err == nil {
		for pkg.i = 1; pkg.i < len(args); pkg.i++ {
			pkg.Reset()
			pkg.a = args[pkg.i]

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
			matched, stop, err = xxTestCmd(pkg, &goCommand, rootCmd, args)
			if e, ok := err.(*ErrorForCmdr); ok {
				ferr("%v", e)
				if !e.Ignorable {
					return
				}
			}
			if stop {
				if pkg.lastCommandHeld || (matched && pkg.flg == nil) {
					err = afterInternalExec(pkg, rootCmd, goCommand, args)
				}
				return
			}
		}

		err = afterInternalExec(pkg, rootCmd, goCommand, args)
	}
	return
}

func xxTestCmd(pkg *ptpkg, goCommand **Command, rootCmd *RootCommand, args []string) (matched, stop bool, err error) {
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
		if pkg.flg != nil && pkg.found {
			matched = true
			return
		}

		// fn + val
		// fn: short,
		// fn: long
		// fn: short||val: such as '-t3'
		// fn: long=val, long='val', long="val", long val, long 'val', long "val"
		// fn: longval, long'val', long"val"

		pkg.savedGoCommand = *goCommand
		cc := *goCommand
		if matched, stop, err = flagsMatching(pkg, cc, goCommand, args); stop || err != nil {
			return
		}

	} else {
		// testing the next command, but the last one has already been the end of command series.
		if pkg.lastCommandHeld {
			pkg.i--
			stop = true
			return
		}

		// or, keep going on...
		if matched, stop, err = cmdMatching(pkg, goCommand, args); stop || err != nil {
			return
		}
	}
	return
}

func preprocess(rootCmd *RootCommand, args []string) (err error) {
	for _, x := range beforeXrefBuilding {
		x(rootCmd, args)
	}

	err = buildXref(rootCmd)

	if err == nil {
		err = rxxtOptions.buildAutomaticEnv(rootCmd)
	}

	if err == nil {
		for _, x := range afterXrefBuilt {
			x(rootCmd, args)
		}
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
					if err = rootCmd.PreAction(goCommand, args); err == ErrShouldBeStopException {
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

		// and now, loading the external configuration files
		err = loadFromPredefinedLocation(rootCmd)
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
			conf.CfgFile = fn
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
	lastCommandHeld   bool
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

func (pkg *ptpkg) Reset() {
	pkg.assigned = false
	pkg.found = false
	pkg.short = false

	pkg.savedFn = ""
	pkg.savedVal = ""
	pkg.fn = ""
	pkg.val = ""
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
		if _, ok := pkg.flg.DefaultValue.(time.Duration); ok {
			if err = processTypeDuration(pkg, args); err != nil {
				ferr("wrong time.Duration: flag=%v, value=%v", pkg.fn, pkg.val)
				return
			}
			// ferr("wrong time.Duration: flag=%v, value=%v", pkg.fn, pkg.val)
			return
		}
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

	var v = pkg.flg.DefaultValue
	var keyPath = backtraceFlagNames(pkg.flg)
	xxSet(keyPath, v, pkg)
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
			content = []byte("demo for testing")
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
	return processTypeIntCore(pkg, args)
}

func processTypeDuration(pkg *ptpkg, args []string) (err error) {
	if err = preprocessPkg(pkg, args); err == nil {
		var v time.Duration
		v, err = time.ParseDuration(pkg.val)
		if err == nil {
			var keyPath = backtraceFlagNames(pkg.flg)
			xxSet(keyPath, v, pkg)
		}
	}
	return
}

func xxSet(keyPath string, v interface{}, pkg *ptpkg) {
	if pkg.a[0] == '~' {
		rxxtOptions.SetNx(keyPath, v)
	} else {
		rxxtOptions.Set(keyPath, v)
	}
	if pkg.flg != nil && pkg.flg.onSet != nil {
		pkg.flg.onSet(keyPath, v)
	}
	pkg.found = true
}

func processTypeIntCore(pkg *ptpkg, args []string) (err error) {
	v, err := strconv.ParseInt(pkg.val, 10, 64)
	if err != nil {
		if _, ok := pkg.flg.DefaultValue.(time.Duration); ok {
			err = processTypeDuration(pkg, args)
			return
		}
		ferr("wrong number: flag=%v, number=%v", pkg.fn, pkg.val)
		err = fmt.Errorf("wrong number: flag=%v, number=%v, inner error is: %v", pkg.fn, pkg.val, err)
	}

	var keyPath = backtraceFlagNames(pkg.flg)
	xxSet(keyPath, v, pkg)
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

	var keyPath = backtraceFlagNames(pkg.flg)
	xxSet(keyPath, v, pkg)
	return
}

func processTypeString(pkg *ptpkg, args []string) (err error) {
	if err = preprocessPkg(pkg, args); err != nil {
		return
	}

	if len(pkg.flg.ValidArgs) > 0 {
		// validate for enum
		for _, w := range pkg.flg.ValidArgs {
			if pkg.val == w {
				goto SAVE_IT
			}
		}
		pkg.found = true
		err = NewError(ShouldIgnoreWrongEnumValue, errWrongEnumValue, pkg.val, pkg.fn, pkg.flg.owner.GetName())
		return
	}

SAVE_IT:
	var v = pkg.val
	var keyPath = backtraceFlagNames(pkg.flg)
	xxSet(keyPath, v, pkg)
	return
}

func processTypeStringSlice(pkg *ptpkg, args []string) (err error) {
	if err = preprocessPkg(pkg, args); err != nil {
		return
	}

	var v = strings.Split(pkg.val, ",")
	var keyPath = backtraceFlagNames(pkg.flg)
	xxSet(keyPath, v, pkg)
	return
}

func processTypeIntSlice(pkg *ptpkg, args []string) (err error) {
	if err = preprocessPkg(pkg, args); err != nil {
		return
	}

	v := make([]int64, 0)
	for _, x := range strings.Split(pkg.val, ",") {
		if xi, err := strconv.ParseInt(x, 10, 64); err == nil {
			v = append(v, xi)
		}
	}

	var keyPath = backtraceFlagNames(pkg.flg)
	xxSet(keyPath, v, pkg)
	return
}

func processTypeUintSlice(pkg *ptpkg, args []string) (err error) {
	if err = preprocessPkg(pkg, args); err != nil {
		return
	}

	v := make([]uint64, 0)
	for _, x := range strings.Split(pkg.val, ",") {
		if xi, err := strconv.ParseUint(x, 10, 64); err == nil {
			v = append(v, xi)
		}
	}

	var keyPath = backtraceFlagNames(pkg.flg)
	xxSet(keyPath, v, pkg)
	return
}
