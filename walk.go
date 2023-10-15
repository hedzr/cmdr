/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"encoding"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

// WalkAllCommands loops for all commands, starting from root.
func WalkAllCommands(walk func(cmd *Command, index, level int) (err error)) (err error) {
	err = walkFromCommand(nil, 0, 0, walk)
	return
}

func walkFromCommand(cmd *Command, index, level int, walk func(cmd *Command, index, level int) (err error)) (err error) {
	if cmd == nil {
		cmd = &internalGetWorker().rootCommand.Command
	}

	// run callback for this command at first
	err = walk(cmd, index, level)

	if err == nil {
		for ix, cc := range cmd.SubCommands {
			if err = walkFromCommand(cc, ix, level+1, walk); err != nil {
				if err == ErrShouldBeStopException { //nolint:errorlint //like it
					err = nil // not an error
				}
				return
			}
		}
	}
	return
}

// InvokeCommand invokes a sub-command internally.
func InvokeCommand(dottedCommandPath string, extraArgs ...string) (err error) {
	cc := dottedPathToCommand(dottedCommandPath, nil)
	if cc != nil {
		w := internalGetWorker()
		action := cc.Action
		if isAllowDefaultAction(action == nil) {
			action = defaultAction
		}
		err = w.doInvokeCommand(w.rootCommand, action, cc, extraArgs)
	}
	return
}

var (
	assumeDefaultAction bool
	defaultAction       = defaultActionImpl
)

func isAllowDefaultAction(nilAction bool) bool {
	if val, ok := os.LookupEnv("FORCE_DEFAULT_ACTION"); ok {
		if toBool(val, false) {
			return true
		}
	}
	return nilAction && assumeDefaultAction
	// return toBool(os.Getenv("FORCE_DEFAULT_ACTION"))
}

func typfmtv(v *reflect.Value) string {
	if v == nil || !v.IsValid() {
		return "<invalid>"
	}
	t := v.Type()
	return fmt.Sprintf("%v (%v)", t, t.Kind())
}

func defaultActionImpl(cmd *Command, args []string) (err error) {
	var sb strings.Builder
	_, _ = fmt.Fprintf(&sb, `
 Command-Line: %v
      Command: %v
  Description: %v
         Args: %v
        Flags:
`,
		strings.Join(os.Args, " "),
		cmd.GetDottedNamePath(),
		getCPT().Translate(cmd.Description, 0),
		args)

	for _, f := range GetHitFlags() {
		kp := f.GetDottedNamePathEx()
		v := GetR(kp)
		switch tv := v.(type) {
		case *time.Time:
			_, _ = fmt.Fprintf(&sb, `
	- %q: %#v (%T) (hit = %q, %v)`, kp, tv.Format("2006-01-02 15:04:05.999999999 -0700"), v, f.GetHitStr(), f.times)
		case time.Time:
			_, _ = fmt.Fprintf(&sb, `
	- %q: %#v (%T) (hit = %q, %v)`, kp, tv.Format("2006-01-02 15:04:05.999999999 -0700"), v, f.GetHitStr(), f.times)
		case encoding.TextMarshaler:
			var txt []byte
			txt, err = tv.MarshalText()
			_, _ = fmt.Fprintf(&sb, `
	- %q: %#v (%T) (hit = %q, %v)`, kp, string(txt), v, f.GetHitStr(), f.times)
		default:
			rv := reflect.ValueOf(v)
			_, _ = fmt.Fprintf(&sb, `
	- %q: %v (hit = %q, %v)`, kp, typfmtv(&rv), f.GetHitStr(), f.times)
		}
	}

	sb.WriteRune('\n')

	_, _ = fmt.Fprintf(os.Stdout, sb.String())

	w := internalGetWorker()
	initTabStop(defaultTabStop)
	w.currentHelpPainter.FpPrintHelpTailLine(cmd)

	return
}
