package cli

import (
	"testing"
)

func TestScreen(t *testing.T) {
	root := rootCmdForTesting()
	(&helpPrinter{}).Print(nil, root.Cmd)
	t.Log()
}

func TestDottedPathToCommandOrFlag(t *testing.T) {
	root := rootCmdForTesting()

	r := root.Cmd.(*CmdS)
	t.Logf("cmd: %v", r.commands[0])
	cmd, ff := r.dottedPathToCommandOrFlag("list")
	t.Logf("L1: cmd: %v; ff: %v", cmd, ff)
	cmd, ff = r.dottedPathToCommandOrFlag("list.retry")
	t.Logf("L2: cmd: %v; ff: %v", cmd, ff)

	// DottedPathToCommandOrFlag, backtraceCmdNames
	cmd, ff = DottedPathToCommandOrFlag1("list.retry", root.Cmd)
	t.Logf("cmd: %v (%v); ff: %v (%v)", cmd, cmd.GetDottedPath(), ff, ff.GetDottedPath())

	cmd, ff = r.dottedPathToCommandOrFlag("microservices.tags.toggle.address")
	t.Logf("L3: cmd: %v; ff: %v", cmd, ff)
	t.Logf("L3: cmd: %v (%v); ff: %v (%v)", cmd, cmd.GetDottedPath(), ff, ff.GetDottedPath())

	cmd, ff = r.dottedPathToCommandOrFlag("microservices.tags.toggle.address")
	t.Logf("L: cmd: %v; ff: %v", cmd, ff)
	t.Logf("L: cmd: %v (%v); ff: %v (%v)", cmd, cmd.GetDottedPath(), ff, ff.GetDottedPath())
}
