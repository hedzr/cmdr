package exec

import (
	"os"
	"runtime"
	"testing"
)

func TestScriptsRun(t *testing.T) {
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {

		var err error

		err = InvokeShellScripts("ls -l /",
			WithScriptShell("/bin/bash"),
			WithScriptIsFile(false),
			WithScriptInvoker(Run),
			WithScriptExpander(os.ExpandEnv),
		)
		if err != nil {
			t.Errorf("%v", err)
		}

		err = InvokeShellScripts("ls -l /",
			WithScriptShell("/bin/bash"),
			WithScriptIsFile(false),
			WithScriptInvoker(func(command string, args ...string) (err error) {
				err = New().
					WithCommandArgs(command, args...).
					WithOnOK(func(retCode int, stdoutText string) {
						t.Logf("%v", LeftPad(stdoutText, 4))
					}).
					RunAndCheckError()
				return
			}),
		)
		if err != nil {
			t.Errorf("%v", err)
		}

	}
}
