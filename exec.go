/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

//

//

func Exec(rootCmd *RootCommand) (err error) {
	rootCommand = rootCmd
	rootCommand.ow = bufio.NewWriterSize(os.Stdout, 16384)
	rootCommand.oerr = bufio.NewWriterSize(os.Stderr, 16384)
	defer func() {
		_ = rootCommand.ow.Flush()
		_ = rootCommand.oerr.Flush()
	}()

	buildRootCrossRefs(rootCommand)
	for _, s := range []string{"./ci/etc/%s/%s.yml", "/etc/%s/%s.yml", "/usr/local/etc/%s/%s.yml", os.Getenv("HOME") + "/.%s/%s.yml"} {
		if FileExists(s) {
			err = RxxtOptions.LoadConfigFile(fmt.Sprintf(s, rootCmd.AppName, rootCmd.AppName))
			if err != nil {
				return
			}
		}
	}

	goCommand := &rootCommand.Command
	// helpFlag := rootCommand.allFlags[UNSORTED_GROUP]["help"]

	// logrus.Debug("----------------------- args:")
	// for _, a := range os.Args {
	// 	logrus.Infof(" - %v", a)
	// }

	var (
		pkg           *ptpkg = new(ptpkg)
		ok            bool
		needHelp      bool
		needFlagsHelp bool
		// disableParser bool
		unknownCmds []string
	)

	for pkg.i = 1; pkg.i < len(os.Args); pkg.i++ {
		pkg.a = os.Args[pkg.i]
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
		// --name=consul, --name consul: opt with a string, int, string slice argument
		// -nconsul, -n consul: opt with an argument.
		// -nconsul is not good format, but it could get somewhat works.
		// -t3: opt with an argument.
		if pkg.a[0] == '-' || pkg.a[0] == '/' || pkg.a[0] == '~' {
			if len(pkg.a) == 1 {
				needHelp = true
				needFlagsHelp = true
				continue
			}

			// flag
			if len(pkg.a) > 1 && (pkg.a[1] == '-' || pkg.a[1] == '~') {
				if len(pkg.a) == 2 {
					// disableParser = true // '--': ignore the following args
					break
				}

				// long flag
				pkg.fn = pkg.a[2:]

				if strings.Contains(pkg.fn, "=") {
					aa := strings.Split(pkg.fn, "=")
					pkg.fn = aa[0]
					pkg.val = aa[1]
					pkg.assigned = true
				}
			} else {
				pkg.fn = pkg.a[1:2]
				pkg.savedFn = pkg.a[2:]
				pkg.short = true

				if strings.HasPrefix(pkg.savedFn, "=") {
					pkg.val = pkg.savedFn[1:]
					pkg.savedFn = ""
					pkg.assigned = true
				}
			}

			var suffix = pkg.fn[len(pkg.fn)-1]
			if suffix == '+' || suffix == '-' {
				pkg.fn = pkg.fn[0 : len(pkg.fn)-1]
			} else {
				suffix = 0
			}

			// fn + val
			// fn: short,
			// fn: long
			// fn: short||val: such as '-t3'
			// fn: long=val

			pkg.savedGoCommand = goCommand
			cc := goCommand
		GO_UP:
			pkg.found = false
			if pkg.flg, ok = cc.plainFlags[pkg.fn]; ok {
				if _, ok := pkg.flg.DefaultValue.(bool); ok {
					if suffix == '+' {
						pkg.flg.DefaultValue = true
					} else if suffix == '-' {
						pkg.flg.DefaultValue = false
					} else {
						pkg.flg.DefaultValue = true
					}

					if pkg.a[0] == '~' {
						RxxtOptions.SetNx(backtraceFlagNames(pkg.flg), pkg.flg.DefaultValue)
					} else {
						RxxtOptions.Set(backtraceFlagNames(pkg.flg), pkg.flg.DefaultValue)
					}
					pkg.found = true

					// } else if flg.DefaultValue == nil {
					// 	// nothing to do

				} else {
					vv := reflect.ValueOf(pkg.flg.DefaultValue)
					kind := vv.Kind()
					switch kind {
					case reflect.String:
						processTypeString(pkg)

					case reflect.Slice:
						typ := reflect.TypeOf(pkg.flg.DefaultValue).Elem()
						if typ.Kind() == reflect.String {
							processTypeStringSlice(pkg)
						} else if isTypeSInt(typ.Kind()) {
							processTypeIntSlice(pkg)
						} else if isTypeSInt(typ.Kind()) {
							processTypeUintSlice(pkg)
						}

					default:
						if isTypeSInt(kind) {
							processTypeInt(pkg)
						} else if isTypeUInt(kind) {
							processTypeUint(pkg)
						} else {
							ferr("Unacceptable default value kind=%v", kind)
						}
					}
				}

				if pkg.found {
					// if !GetBool("app.quiet") {
					// 	logrus.Debugf("-- flag '%v' hit, go ahead...", pkg.flg.GetTitleName())
					// }
					if pkg.flg.Action != nil {
						if err = pkg.flg.Action(goCommand, getArgs(pkg)); err == ShouldBeStopException {
							return nil
						}
					}

					if !pkg.assigned {
						if len(pkg.savedFn) > 0 && len(pkg.savedVal) == 0 {
							pkg.fn = pkg.savedFn[0:1]
							pkg.savedFn = pkg.savedFn[1:]
							goto GO_UP
						}
					}
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
						goCommand = pkg.savedGoCommand
						goto GO_UP
					}
				}
				ferr("Unknown flag: %v", pkg.a)
			}

		} else {
			// command, files
			if cmd, ok := goCommand.plainCmds[pkg.a]; ok {
				goCommand = cmd
				// logrus.Debugf("-- command '%v' hit, go ahead...", cmd.GetTitleName())
				if cmd.PreAction != nil {
					if err = cmd.PreAction(goCommand, getArgs(pkg)); err == ShouldBeStopException {
						return nil
					}
				}
			} else {
				// ferr("Unknown command: %v (disableParser=%v)", pkg.a, disableParser)
				unknownCmds = append(unknownCmds, pkg.a)
			}
		}
	}

	if !needHelp {
		needHelp = GetBool("app.help")
	}

	if !needHelp && len(unknownCmds) == 0 {
		if goCommand.Action != nil {
			args := getArgs(pkg)

			if goCommand != &rootCommand.Command {
				if rootCommand.PostAction != nil {
					defer rootCommand.PostAction(goCommand, args)
				}
				if rootCommand.PreAction != nil {
					if err = rootCommand.PreAction(goCommand, getArgs(pkg)); err == ShouldBeStopException {
						return nil
					}
				}
			}

			if goCommand.PostAction != nil {
				defer goCommand.PostAction(goCommand, args)
			}

			if err = goCommand.Action(goCommand, args); err == ShouldBeStopException {
				return nil
			}

			return
		}
	}

	if GetInt("app.help-zsh") > 0 || GetBool("app.help-bash") {
		if len(goCommand.SubCommands) == 0 && !needFlagsHelp {
			// needFlagsHelp = true
		}
	}

	printHelp(goCommand, needFlagsHelp)

	return
}

func getArgs(pkg *ptpkg) []string {
	var args []string
	if pkg.i+1 < len(os.Args) {
		args = os.Args[pkg.i+1:]
	}
	return args
}

func isTypeInt(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
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

func isTypeUInt(kind reflect.Kind) bool {
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
}

func preprocessPkg(pkg *ptpkg) {
	if !pkg.assigned {
		if len(pkg.savedVal) > 0 {
			pkg.val = pkg.savedVal
			pkg.savedVal = ""
		} else if len(pkg.savedFn) > 0 {
			pkg.val = pkg.savedFn
			pkg.savedFn = ""
		} else {
			pkg.i++
			pkg.val = os.Args[pkg.i]
		}
		pkg.assigned = true
	}
}

func processTypeInt(pkg *ptpkg) {
	preprocessPkg(pkg)

	v, err := strconv.ParseInt(pkg.val, 10, 64)
	if err != nil {
		ferr("wrong number: flag=%v, number=%v", pkg.fn, pkg.val)
	}

	if pkg.a[0] == '~' {
		RxxtOptions.SetNx(backtraceFlagNames(pkg.flg), v)
	} else {
		RxxtOptions.Set(backtraceFlagNames(pkg.flg), v)
	}
	pkg.found = true
	return
}

func processTypeUint(pkg *ptpkg) {
	preprocessPkg(pkg)

	v, err := strconv.ParseUint(pkg.val, 10, 64)
	if err != nil {
		ferr("wrong number: flag=%v, number=%v", pkg.fn, pkg.val)
	}

	if pkg.a[0] == '~' {
		RxxtOptions.SetNx(backtraceFlagNames(pkg.flg), v)
	} else {
		RxxtOptions.Set(backtraceFlagNames(pkg.flg), v)
	}
	pkg.found = true
	return
}

func processTypeString(pkg *ptpkg) {
	preprocessPkg(pkg)
	if pkg.a[0] == '~' {
		RxxtOptions.SetNx(backtraceFlagNames(pkg.flg), pkg.val)
	} else {
		RxxtOptions.Set(backtraceFlagNames(pkg.flg), pkg.val)
	}
	pkg.found = true
}

func processTypeStringSlice(pkg *ptpkg) {
	preprocessPkg(pkg)
	if pkg.a[0] == '~' {
		RxxtOptions.SetNx(backtraceFlagNames(pkg.flg), strings.Split(pkg.val, ","))
	} else {
		RxxtOptions.Set(backtraceFlagNames(pkg.flg), strings.Split(pkg.val, ","))
	}
	pkg.found = true
}

func processTypeIntSlice(pkg *ptpkg) {
	preprocessPkg(pkg)

	valary := make([]int64, 0)
	for _, x := range strings.Split(pkg.val, ",") {
		if xi, err := strconv.ParseInt(x, 10, 64); err == nil {
			valary = append(valary, xi)
		}
	}

	if pkg.a[0] == '~' {
		RxxtOptions.SetNx(backtraceFlagNames(pkg.flg), valary)
	} else {
		RxxtOptions.Set(backtraceFlagNames(pkg.flg), valary)
	}
	pkg.found = true
}

func processTypeUintSlice(pkg *ptpkg) {
	preprocessPkg(pkg)

	valary := make([]uint64, 0)
	for _, x := range strings.Split(pkg.val, ",") {
		if xi, err := strconv.ParseUint(x, 10, 64); err == nil {
			valary = append(valary, xi)
		}
	}

	if pkg.a[0] == '~' {
		RxxtOptions.SetNx(backtraceFlagNames(pkg.flg), valary)
	} else {
		RxxtOptions.Set(backtraceFlagNames(pkg.flg), valary)
	}
	pkg.found = true
}
