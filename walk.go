/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import "fmt"

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
				if err == ErrShouldBeStopException {
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
		if action == nil && assumeDefaultAction {
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

func defaultActionImpl(cmd *Command, args []string) (err error) {
	fmt.Printf(`
    Command: %v
Description: %q
       Args: %v
      Flags:

`,
		cmd.GetDottedNamePath(), cmd.Description, args)
	for _, f := range GetHitFlags() {
		kp := f.GetDottedNamePath()
		v := GetR(kp)
		fmt.Printf(`  %v: %v
`, kp, v)
	}

	println()

	return
}
