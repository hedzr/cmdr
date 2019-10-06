/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"fmt"
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

func (pkg *ptpkg) toggleGroup() {
	tg := pkg.flg.ToggleGroup
	if len(tg) > 0 {
		for _, f := range pkg.flg.owner.Flags {
			if f.ToggleGroup == tg && (isBool(f.DefaultValue) || isNil1(f.DefaultValue)) {
				uniqueWorker.rxxtOptions.Set(uniqueWorker.backtraceFlagNames(pkg.flg), false)
			}
		}
	}
}

func (pkg *ptpkg) findValueAttached(fn *string) {
	if strings.Contains(*fn, "=") {
		aa := strings.Split(*fn, "=")
		*fn = aa[0]
		pkg.val = trimQuotes(aa[1])
		pkg.assigned = true
	} else {
		pkg.splitQuotedValueIfNecessary(fn)
	}
}

func (pkg *ptpkg) splitQuotedValueIfNecessary(fn *string) {
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

func (pkg *ptpkg) matchShortFlag(goCommand *Command, a string) (i int) {
	for i = len(a); i > 1; i-- {
		fn := a[1:i]
		if _, ok := goCommand.plainShortFlags[fn]; ok {
			return
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
	} else {
		ferr("Unacceptable default value kind=%v", kind)
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
	var keyPath = uniqueWorker.backtraceFlagNames(pkg.flg)
	pkg.xxSet(keyPath, v)
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
			if pkg.i < len(args)-1 && args[pkg.i+1][0] != '-' && (args[pkg.i+1][0] != '~' || args[pkg.i+1][1] != '~') {
				pkg.i++
				pkg.val = args[pkg.i]
			} else {
				if len(pkg.flg.ExternalTool) > 0 {
					err = pkg.processExternalTool()
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

func (pkg *ptpkg) processExternalTool() (err error) {
	switch pkg.flg.ExternalTool {
	case ExternalToolPasswordInput:
		fmt.Print("Password: ")
		var password string
		if InTesting() {
			password = "demo"
		} else {
			if password, err = readPassword(); err != nil {
				return
			}
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
			if content, err = LaunchEditor(editor); err != nil {
				return
			}
		}
		pkg.val = string(content)
	}
	return
}

func (pkg *ptpkg) processTypeInt(args []string) (err error) {
	if err = pkg.preprocessPkg(args); err != nil {
		return
	}
	return pkg.processTypeIntCore(args)
}

func (pkg *ptpkg) processTypeDuration(args []string) (err error) {
	if err = pkg.preprocessPkg(args); err == nil {
		var v time.Duration
		v, err = time.ParseDuration(pkg.val)
		if err == nil {
			var keyPath = uniqueWorker.backtraceFlagNames(pkg.flg)
			pkg.xxSet(keyPath, v)
		}
	}
	return
}

func (pkg *ptpkg) xxSet(keyPath string, v interface{}) {
	if pkg.a[0] == '~' {
		uniqueWorker.rxxtOptions.SetNx(keyPath, v)
	} else {
		uniqueWorker.rxxtOptions.Set(keyPath, v)
	}
	if pkg.flg != nil && pkg.flg.onSet != nil {
		pkg.flg.onSet(keyPath, v)
	}
	pkg.found = true
}

func (pkg *ptpkg) processTypeIntCore(args []string) (err error) {
	v, err := strconv.ParseInt(pkg.val, 10, 64)
	if err != nil {
		if _, ok := pkg.flg.DefaultValue.(time.Duration); ok {
			err = pkg.processTypeDuration(args)
			return
		}
		ferr("wrong number: flag=%v, number=%v", pkg.fn, pkg.val)
		err = fmt.Errorf("wrong number: flag=%v, number=%v, inner error is: %v", pkg.fn, pkg.val, err)
	}

	var keyPath = uniqueWorker.backtraceFlagNames(pkg.flg)
	pkg.xxSet(keyPath, v)
	return
}

func (pkg *ptpkg) processTypeUint(args []string) (err error) {
	if err = pkg.preprocessPkg(args); err != nil {
		return
	}

	v, err := strconv.ParseUint(pkg.val, 10, 64)
	if err != nil {
		ferr("wrong number: flag=%v, number=%v", pkg.fn, pkg.val)
	}

	var keyPath = uniqueWorker.backtraceFlagNames(pkg.flg)
	pkg.xxSet(keyPath, v)
	return
}

func (pkg *ptpkg) processTypeString(args []string) (err error) {
	if err = pkg.preprocessPkg(args); err != nil {
		return
	}

	if len(pkg.flg.ValidArgs) > 0 {
		// validate for enum
		for _, w := range pkg.flg.ValidArgs {
			if pkg.val == w {
				goto saveIt
			}
		}
		pkg.found = true
		err = NewError(uniqueWorker.shouldIgnoreWrongEnumValue, errWrongEnumValue, pkg.val, pkg.fn, pkg.flg.owner.GetName())
		return
	}

saveIt:
	var v = pkg.val
	var keyPath = uniqueWorker.backtraceFlagNames(pkg.flg)
	pkg.xxSet(keyPath, v)
	return
}

func (pkg *ptpkg) processTypeStringSlice(args []string) (err error) {
	if err = pkg.preprocessPkg(args); err != nil {
		return
	}

	var v = strings.Split(pkg.val, ",")

	var keyPath = uniqueWorker.backtraceFlagNames(pkg.flg)
	var existedVal = uniqueWorker.rxxtOptions.GetStringSlice(wrapWithRxxtPrefix(keyPath))
	if reflect.DeepEqual(existedVal, pkg.flg.DefaultValue) {
		existedVal = nil
	}
	pkg.xxSet(keyPath, append(v, existedVal...))
	return
}

func (pkg *ptpkg) processTypeIntSlice(args []string) (err error) {
	if err = pkg.preprocessPkg(args); err != nil {
		return
	}

	v := make([]int64, 0)
	for _, x := range strings.Split(pkg.val, ",") {
		if xi, err := strconv.ParseInt(x, 10, 64); err == nil {
			v = append(v, xi)
		}
	}

	var keyPath = uniqueWorker.backtraceFlagNames(pkg.flg)
	// pkg.xxSet(keyPath, v)
	var existedVal = uniqueWorker.rxxtOptions.GetInt64Slice(wrapWithRxxtPrefix(keyPath))
	if reflect.DeepEqual(existedVal, pkg.flg.DefaultValue) {
		existedVal = nil
	}
	pkg.xxSet(keyPath, append(v, existedVal...))
	return
}

func (pkg *ptpkg) processTypeUintSlice(args []string) (err error) {
	if err = pkg.preprocessPkg(args); err != nil {
		return
	}

	v := make([]uint64, 0)
	for _, x := range strings.Split(pkg.val, ",") {
		if xi, err := strconv.ParseUint(x, 10, 64); err == nil {
			v = append(v, xi)
		}
	}

	var keyPath = uniqueWorker.backtraceFlagNames(pkg.flg)
	// pkg.xxSet(keyPath, v)
	var existedVal = uniqueWorker.rxxtOptions.GetUint64Slice(wrapWithRxxtPrefix(keyPath))
	if reflect.DeepEqual(existedVal, pkg.flg.DefaultValue) {
		existedVal = nil
	}
	pkg.xxSet(keyPath, append(v, existedVal...))
	return
}
