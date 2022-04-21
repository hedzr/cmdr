/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"fmt"
	"github.com/hedzr/cmdr/tool"
	"gopkg.in/hedzr/errors.v3"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type ptpkg struct {
	assigned          bool
	found             bool
	short             bool
	lastCommandHeld   bool
	fn, val           string
	savedFn, savedVal string
	i                 int
	iLastCommand      int
	a                 string
	flg               *Flag
	savedGoCommand    *Command
	needHelp          bool
	needFlagsHelp     bool
	suffix            uint8
	unknownCmds       []string
	unknownFlags      []string
	remainArgs        []string
	aliasCommand      *Command
}

func (pkg *ptpkg) store() *Options                    { return currentOptions() }
func (pkg *ptpkg) storeFrom(wkr *ExecWorker) *Options { return wkr.rxxtOptions }

func (pkg *ptpkg) ResetAnd(n string) (length int) {
	pkg.Reset()
	pkg.a = n
	return len(n)
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

func (pkg *ptpkg) tryToggleGroup() {
	tg := pkg.flg.ToggleGroup
	if len(tg) > 0 {
		wkr := internalGetWorker()
		for _, f := range pkg.flg.owner.Flags {
			if f.ToggleGroup == tg && (isBool(f.DefaultValue) || isNil1(f.DefaultValue)) {
				var store = pkg.storeFrom(wkr)
				if f != pkg.flg {
					store.Set(backtraceFlagNames(f), false)
					f.DefaultValue = false
				} else {
					store.Set(backtraceFlagNames(f), true)
					store.Set(backtraceCmdNames(f.owner, false)+"."+f.ToggleGroup, f.Full)
					f.DefaultValue = true
				}
			}
		}
	}
}

func (pkg *ptpkg) findValueAttached(fn *string) {
	if strings.Contains(*fn, "=") {
		aa := strings.Split(*fn, "=")
		*fn = aa[0]
		pkg.val = tool.StripQuotes(aa[1])
		pkg.assigned = true
	} else {
		pkg.splitQuotedValueIfNecessary(fn)
	}
}

func (pkg *ptpkg) splitQuotedValueIfNecessary(fn *string) {
	if pos := strings.Index(*fn, "'"); pos >= 0 {
		pkg.val = tool.StripQuotes((*fn)[pos:])
		*fn = (*fn)[0:pos]
		pkg.assigned = true
	} else if pos := strings.Index(*fn, "\""); pos >= 0 {
		pkg.val = tool.StripQuotes((*fn)[pos:])
		*fn = (*fn)[0:pos]
		pkg.assigned = true
		// } else {
		// --xVALUE need to be parsed.
	}
}

func (pkg *ptpkg) doExtractingHeadLikeFlag(goCommand **Command, args []string) (ok bool, err error) {
	if (*goCommand).headLikeFlag != nil && tool.IsDigitHeavy(pkg.a[1:]) {
		// println("head-like")
		pkg.short = true
		pkg.flg = (*goCommand).headLikeFlag
		pkg.val = pkg.a[1:]
		pkg.fn = pkg.flg.Short
		pkg.found = true
		err = pkg.processTypeIntCore(args)
		ok = true
	}
	return
}

func (pkg *ptpkg) doParseSuffix() {
	pkg.suffix = pkg.a[len(pkg.a)-1]
	if pkg.suffix == '+' || pkg.suffix == '-' {
		pkg.a = pkg.a[0 : len(pkg.a)-1]
	} else {
		pkg.suffix = 0
	}
}

func (pkg *ptpkg) matchLongFlag(goCommand **Command) {
	if *goCommand == &((*goCommand).root.Command) && pkg.aliasCommand != nil {
		if i := pkg.matchLongFlagsRecursively(pkg.aliasCommand, pkg.a, 2); i >= 0 {
			pkg.fn = pkg.a[2:]
			pkg.findValueAttached(&pkg.fn)
			return
		}
	}

	if i := pkg.matchLongFlagsRecursively(*goCommand, pkg.a, 2); i >= 0 {
		pkg.fn = pkg.a[2:]
		pkg.findValueAttached(&pkg.fn)
		return
	}

	pkg.fn = pkg.a[2:]
	pkg.findValueAttached(&pkg.fn)
}

func (pkg *ptpkg) matchLongFlagsRecursively(cc *Command, a string, startPos int) (i int) {
	var ln int
	var fn string
	var ok bool
	if pkg.assigned {
		if pkg.flg != nil {
			return
		}
	}

retry:
	i, ln = -1, len(a)
	for ; ln > startPos; ln-- {
		fn = a[startPos:ln]
		pkg.flg, ok = cc.plainLongFlags[fn]
		if ok {
			if ln < len(a) {
				pkg.val = tool.StripQuotes(a[ln:])
				pkg.assigned = true
			}
			pkg.fn = fn
			i = ln
			return
		}
	}
	if cc.HasParent() {
		cc = cc.owner
		goto retry
	}
	return
}

func (pkg *ptpkg) matchShortFlag(goCommand **Command) {
	if i := pkg.matchLongestShortFlag(*goCommand, pkg.a, 1); i >= 0 {
		pkg.fn = pkg.a[1:i]
		pkg.savedFn = pkg.a[i:]
		pkg.short = true
		pkg.findValueAttached(&pkg.savedFn)
		return
	}

	if *goCommand == &((*goCommand).root.Command) && pkg.aliasCommand != nil {
		if i := pkg.matchLongestShortFlag(pkg.aliasCommand, pkg.a, 1); i >= 0 {
			pkg.fn = pkg.a[1:i]
			pkg.savedFn = pkg.a[i:]
			pkg.short = true
			pkg.findValueAttached(&pkg.savedFn)
			return
		}
	}

	// no matched short flags found
	if len(pkg.a) >= 2 && pkg.a[0] == '/' && pkg.a[1] != '/' {
		if i := pkg.matchLongFlagsRecursively(*goCommand, pkg.a, 2); i >= 0 {
			return
		}
	}

	pkg.fn = pkg.a[1:2]     // from one char
	pkg.savedFn = pkg.a[2:] // save others

	pkg.short = true
	pkg.findValueAttached(&pkg.savedFn)
}

func (pkg *ptpkg) matchLongestShortFlag(cc *Command, a string, startPos int) (i int) {
	type MS struct {
		index int
		fn    string
	}
	var matched []MS
	var longest = -1
	for i = len(a); i > startPos; i-- {
		fn := a[startPos:i]
		if _, ok := cc.plainShortFlags[fn]; ok {
			matched = append(matched, MS{index: i, fn: fn})
			if longest < i {
				longest = i
			}
		}
	}

	if longest > 0 {
		for _, ms := range matched {
			if ms.index == longest {
				return longest
			}
		}
	}
	return -1
}

func (pkg *ptpkg) tryExtractingValue(args []string) (err error) {
	if _, ok := pkg.flg.DefaultValue.(bool); ok {
		return pkg.tryExtractingBoolValue()
	}

	vv := reflect.ValueOf(pkg.flg.DefaultValue)
	kind := vv.Kind()
	switch kind {
	case reflect.String:
		err = pkg.processTypeString(args)

	case reflect.Slice:
		err = pkg.tryExtractingSliceValue(args)

	default:
		err = pkg.tryExtractingOthers(args, kind)
	}

	// fmt.Println("tryExtractingValue end")
	return
}

func (pkg *ptpkg) tryExtractingOthers(args []string, kind reflect.Kind) (err error) {
	if isTypeSInt(kind) {
		if _, ok := pkg.flg.DefaultValue.(time.Duration); ok {
			if err = pkg.processTypeDuration(args); err != nil {
				ferr("wrong time.Duration: flag=%v, value=%v", pkg.fn, pkg.val)
				return
			}
			// ferr("wrong time.Duration: flag=%v, value=%v", pkg.fn, pkg.val)
			return
		}
		err = pkg.processTypeInt(args)
	} else if isTypeUint(kind) {
		err = pkg.processTypeUint(args)
	} else if isTypeFloat(kind) {
		err = pkg.processTypeFloat(args)
	} else if isTypeComplex(kind) {
		err = pkg.processTypeComplex(args)
	} else {
		fwrn("Unacceptable default value kind=%v", kind)
		// try parsing as a string
		err = pkg.processTypeString(args)
	}
	return
}

func (pkg *ptpkg) tryExtractingSliceValue(args []string) (err error) {
	typ := reflect.TypeOf(pkg.flg.DefaultValue).Elem()
	if typ.Kind() == reflect.String {
		err = pkg.processTypeStringSlice(args)
	} else if isTypeSInt(typ.Kind()) {
		err = pkg.processTypeIntSlice(args)
	} else if isTypeUint(typ.Kind()) {
		err = pkg.processTypeUintSlice(args)
	}
	return
}

func (pkg *ptpkg) tryExtractingBoolValue() (err error) {
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
	pkg.xxSet(keyPath, v, false)
	return
}

func (pkg *ptpkg) preprocessPkg(args []string) (err error) {
	if !pkg.assigned {
		if len(pkg.savedVal) > 0 {
			pkg.val = pkg.savedVal
			pkg.savedVal = ""
		} else if len(pkg.savedFn) > 0 {
			pkg.val = pkg.savedFn
			pkg.savedFn = ""
		} else {
			yes := false
			if pkg.i < len(args)-1 {
				if len(args[pkg.i+1]) == 0 {
					yes = true
				} else if args[pkg.i+1][0] != '-' && (args[pkg.i+1][0] != '~' || args[pkg.i+1][1] != '~') {
					yes = true
				}
			}
			if yes {
				pkg.i++
				pkg.val = args[pkg.i]
			} else {
				if len(pkg.flg.ExternalTool) > 0 {
					err = pkg.processExternalTool()
				} else if GetStrictMode() {
					err = errors.New("unexpected end of command line [i=%v,args=(%v)], need more args for %v", pkg.i, args, pkg)
					return
				}
			}
		}
		pkg.assigned = true
	}
	return
}

func (pkg *ptpkg) processExternalTool() (err error) {
	switch pkg.flg.ExternalTool {
	case ExternalToolPasswordInput:
		var password string
		if InTesting() {
			fmt.Printf("go-demo")
			password = "demo"
		} else {
			// fmt.Printf("InTesting = false,,,\n")
			fmt.Print("Password: ")
			password, err = tool.ReadPassword()
		}
		pkg.val = password

	default:
		editor := os.Getenv(pkg.flg.ExternalTool)
		if len(editor) == 0 {
			editor = DefaultEditor
		}
		var content []byte
		if InTesting() {
			content = []byte("demo for testing")
		} else {
			content, err = tool.LaunchEditor(editor)
		}
		pkg.val = string(content)
	}
	return
}

// xxSet _
// replaceOrAppend: // true: replace old, false: append to old value
func (pkg *ptpkg) xxSet(keyPath string, v interface{}, replaceOrAppend bool) {
	if !(len(pkg.a) > 0 && pkg.a[0] == '~') {
		keyPath = wrapWithRxxtPrefix(keyPath)
	}

	store := pkg.store()
	if replaceOrAppend {
		store.SetNxOverwrite(keyPath, v)
	} else {
		store.SetNx(keyPath, v)
	}

	if pkg.flg != nil && pkg.flg.onSet != nil {
		pkg.flg.onSet(keyPath, v)
	}

	pkg.found = true
}

func (pkg *ptpkg) processTypeInt(args []string) (err error) {
	if err = pkg.preprocessPkg(args); err == nil {
		err = pkg.processTypeIntCore(args)
	}
	return
}

func (pkg *ptpkg) processTypeDuration(args []string) (err error) {
	if err = pkg.preprocessPkg(args); err == nil {
		var v time.Duration
		v, err = time.ParseDuration(pkg.val)
		if err == nil {
			// flog("    .  . [duration] %q => %v", pkg.val, v)
			var keyPath = backtraceFlagNames(pkg.flg)
			pkg.xxSet(keyPath, v, false)
		}
	}
	return
}

func (pkg *ptpkg) processTypeIntCore(args []string) (err error) {
	var v int64
	v, err = strconv.ParseInt(pkg.val, 0, 0)
	if err != nil {
		if _, ok := pkg.flg.DefaultValue.(time.Duration); ok {
			err = pkg.processTypeDuration(args)
			return
		}
		ferr("wrong number (int): flag=%v, number=%v, err: %v", pkg.fn, pkg.val, err)
		err = errors.New("wrong number (int): flag=%v, number=%v, inner error is: %v", pkg.fn, pkg.val, err)
	}

	var keyPath = backtraceFlagNames(pkg.flg)
	pkg.xxSet(keyPath, v, false)
	return
}

func (pkg *ptpkg) processTypeUint(args []string) (err error) {
	if err = pkg.preprocessPkg(args); err == nil {
		var v uint64
		v, err = strconv.ParseUint(pkg.val, 0, 0)
		if err != nil {
			ferr("wrong number (uint): flag=%v, number=%v, err: %v", pkg.fn, pkg.val, err)
			err = errors.New("wrong number (uint): flag=%v, number=%v, inner error is: %v", pkg.fn, pkg.val, err)
			return
		}

		var keyPath = backtraceFlagNames(pkg.flg)
		pkg.xxSet(keyPath, v, false)
	}
	return
}

func (pkg *ptpkg) processTypeFloat(args []string) (err error) {
	if err = pkg.preprocessPkg(args); err == nil {
		var v float64
		v, err = strconv.ParseFloat(pkg.val, 64)
		if err != nil {
			ferr("wrong number (float): flag=%v, number=%v, err: %v", pkg.fn, pkg.val, err)
			err = errors.New("wrong number (float): flag=%v, number=%v, inner error is: %v", pkg.fn, pkg.val, err)
			return
		}

		var keyPath = backtraceFlagNames(pkg.flg)
		pkg.xxSet(keyPath, v, false)
	}
	return
}

func (pkg *ptpkg) processTypeComplex(args []string) (err error) {
	if err = pkg.preprocessPkg(args); err == nil {
		var v complex128
		v, err = tool.ParseComplexX(pkg.val)
		if err != nil {
			ferr("wrong number (complex): flag=%v, number=%v, err: %v", pkg.fn, pkg.val, err)
			err = errors.New("wrong number (complex): flag=%v, number=%v, inner error is: %v", pkg.fn, pkg.val, err)
			return
		}

		var keyPath = backtraceFlagNames(pkg.flg)
		pkg.xxSet(keyPath, v, false)
	}
	return
}

func (pkg *ptpkg) processTypeString(args []string) (err error) {
	if err = pkg.preprocessPkg(args); err == nil {
		var wkr = internalGetWorker()

		if len(pkg.flg.ValidArgs) > 0 {
			// validate for enum
			for _, w := range pkg.flg.ValidArgs {
				if pkg.val == w {
					goto saveIt
				}
			}
			pkg.found = true
			err = newError(wkr.shouldIgnoreWrongEnumValue,
				errWrongEnumValue, // .Format(pkg.val, pkg.fn, pkg.flg.owner.GetName()),
				pkg.val, pkg.flg.GetTitleZshFlagName(), pkg.flg.owner.GetName(),
			)
			return
		}

	saveIt:
		var v = pkg.val
		var keyPath = backtraceFlagNames(pkg.flg)
		pkg.xxSet(keyPath, v, false)
		pkg.found = true

	}
	return
}

func (pkg *ptpkg) processTypeStringSlice(args []string) (err error) {
	if err = pkg.preprocessPkg(args); err == nil {
		var v = strings.Split(pkg.val, ",")

		var keyPath = backtraceFlagNames(pkg.flg)
		// var replaceOrAppend bool // true: replace old, false: append to old value
		// var existedVal = pkg.store().GetStringSlice(wrapWithRxxtPrefix(keyPath))
		// if reflect.DeepEqual(existedVal, pkg.flg.DefaultValue) || pkg.flg.times == 1 { // if first matching
		//	existedVal, replaceOrAppend = nil, true
		// }
		// pkg.xxSet(keyPath, append(existedVal, v...), replaceOrAppend)
		pkg.xxSet(keyPath, v, false)
	}
	return
}

func (pkg *ptpkg) processTypeIntSlice(args []string) (err error) {
	if err = pkg.preprocessPkg(args); err == nil {
		v := make([]int64, 0)
		for _, x := range strings.Split(pkg.val, ",") {
			if xi, err := strconv.ParseInt(x, 0, 64); err == nil {
				v = append(v, xi)
			}
		}

		var keyPath = backtraceFlagNames(pkg.flg)
		// var replaceOrAppend bool
		// pkg.xxSet(keyPath, v)
		// var existedVal = pkg.store().GetInt64Slice(wrapWithRxxtPrefix(keyPath))
		// if reflect.DeepEqual(existedVal, pkg.flg.DefaultValue) || pkg.flg.times == 1 { // if first matching
		//	existedVal, replaceOrAppend = nil, true
		// }
		// pkg.xxSet(keyPath, append(existedVal, v...), replaceOrAppend)
		pkg.xxSet(keyPath, v, false)
	}
	return
}

func (pkg *ptpkg) processTypeUintSlice(args []string) (err error) {
	if err = pkg.preprocessPkg(args); err == nil {
		v := make([]uint64, 0)
		for _, x := range strings.Split(pkg.val, ",") {
			if xi, err := strconv.ParseUint(x, 0, 64); err == nil {
				v = append(v, xi)
			}
		}

		var keyPath = backtraceFlagNames(pkg.flg)
		// var replaceOrAppend bool
		// pkg.xxSet(keyPath, v)
		// var existedVal = pkg.store().GetUint64Slice(wrapWithRxxtPrefix(keyPath))
		// if reflect.DeepEqual(existedVal, pkg.flg.DefaultValue) || pkg.flg.times == 1 { // if first matching
		//	existedVal, replaceOrAppend = nil, true
		// }
		// pkg.xxSet(keyPath, append(existedVal, v...), replaceOrAppend)
		pkg.xxSet(keyPath, v, false)
	}
	return
}
