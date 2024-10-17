package exec

import (
	"os"
	"runtime"
	"strings"
)

// InvokeShellScripts invokes a shell script fragments with internal
// invoker (typically it's hedzr/log.exec.Run).
//
// InvokeShellScripts finds and prepares the proper shell (bash, or
// powershell, etc.), and call it by invoker.
//
// The invoker can be customized by yours, use WithScriptInvoker.
func InvokeShellScripts(scripts string, opts ...ISSOpt) (err error) {
	return invokeShellScripts(scripts, opts...)
}

type issCtx struct {
	knownShell  string
	invoker     func(command string, args ...string) (err error)
	isFile      bool
	delayedOpts []func(c *issCtx) // not yet
	expander    func(source string) string
	// cmd        *CmdS
	// args       []string
}

type ISSOpt func(c *issCtx)

// WithScriptShell provides a predefined shell executable (short if it's
// in $PATH, or full-path by you risks).
//
// The knownShell looks like shebang.
// Such as: '/bin/bash', or '/usr/bin/env bash', ...
func WithScriptShell(knownShell string) ISSOpt {
	return func(c *issCtx) {
		c.knownShell = knownShell
	}
}

// // WithCmdrEnviron provides the current command hit with its args.
// // Its generally come from cmdr.CmdS.Action of a cmdr CmdS.
// func WithCmdrEnviron(cmd *CmdS, args []string) ISSOpt {
//	return func(c *issCtx) {
//		c.cmd, c.args = cmd, args
//	}
// }

// WithScriptInvoker provides a custom runner to run the shell and scripts.
//
// The default is exec.Run in hedzr/log package.
//
// For example:
//
//	err = InvokeShellScripts("ls -l /",
//		WithScriptShell("/bin/bash"),
//		WithScriptIsFile(false),
//		WithScriptInvoker(func(command string, args ...string) (err error) {
//			err = exec.New().
//				WithCommandArgs(command, args...).
//				WithOnOK(func(retCode int, stdoutText string) {
//					t.Logf("%v", LeftPad(stdoutText, 4))
//				}).
//				RunAndCheckError()
//			return
//		}),
//	)
//	if err != nil {
//		t.Errorf("%v", err)
//	}
func WithScriptInvoker(invoker func(command string, args ...string) (err error)) ISSOpt {
	return func(c *issCtx) {
		c.invoker = invoker
	}
}

// WithScriptExpander providers a string expander for the given script.
//
// You may specify a special one (such as os.ExpandEnv) rather than
// internal default (a dummy functor to return the source directly).
func WithScriptExpander(expander func(source string) string) ISSOpt {
	return func(c *issCtx) {
		c.expander = expander
	}
}

// WithScriptIsFile provides a bool flag for flagging the given
// scripts is a shell scripts fragments or a file.
func WithScriptIsFile(isFile bool) ISSOpt {
	return func(c *issCtx) {
		c.isFile = isFile
	}
}

func invokeShellScripts(scripts string, opts ...ISSOpt) (err error) {

	var c issCtx

	for _, opt := range opts {
		opt(&c)
	}

	var a []string

	if c.knownShell == "" {
		if runtime.GOOS == "windows" {
			c.knownShell = "powershell.exe"
		} else {
			c.knownShell = "/bin/bash"
		}
	} else if strings.Contains(c.knownShell, "/env ") {
		cc := strings.Split(c.knownShell, " ")
		c.knownShell, a = cc[0], append(a, cc[1:]...)
	}

	var scriptFragments string
	// if c.cmd != nil {
	//	scriptFragments = internalGetWorker().expandTmplWithExecutiveEnv(scripts, c.cmd, c.args)
	// } else {
	//	scriptFragments = scripts
	// }
	if c.expander == nil {
		c.expander = os.ExpandEnv
	}
	scriptFragments = c.expander(scripts)

	for _, opt := range c.delayedOpts {
		opt(&c)
	}

	if strings.Contains(c.knownShell, "powershell") {
		a = append(a, "-NoProfile", "-NonInteractive")
		if c.isFile {
			a = append(a, "-File", scriptFragments)
		} else {
			a = append(a, "-CmdS", scriptFragments)
		}
	} else {
		if c.isFile {
			a = append(a, scriptFragments)
		} else {
			a = append(a, "-c", scriptFragments)
		}
	}

	if c.invoker == nil {
		c.invoker = Run
	}

	err = c.invoker(c.knownShell, a...)

	return
}
