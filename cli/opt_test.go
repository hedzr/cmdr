package cli

import (
	"testing"
)

func TestScreen(t *testing.T) {
	root := rootCmdForTesting()
	(&helpPrinter{}).Print(nil, root.Command)
	t.Log()
}

func TestDottedPathToCommandOrFlag(t *testing.T) {
	root := rootCmdForTesting()

	t.Logf("cmd: %v", root.commands[0])
	cmd, ff := root.dottedPathToCommandOrFlag("list")
	t.Logf("L1: cmd: %v; ff: %v", cmd, ff)
	cmd, ff = root.dottedPathToCommandOrFlag("list.retry")
	t.Logf("L2: cmd: %v; ff: %v", cmd, ff)

	// DottedPathToCommandOrFlag, backtraceCmdNames
	cmd, ff = DottedPathToCommandOrFlag("list.retry", root.Command)
	t.Logf("cmd: %v (%v); ff: %v (%v)", cmd, cmd.GetDottedPathFull(), ff, ff.GetDottedPathFull())

	cmd, ff = root.dottedPathToCommandOrFlag("microservices.tags.toggle.address")
	t.Logf("L3: cmd: %v; ff: %v", cmd, ff)
	t.Logf("L3: cmd: %v (%v); ff: %v (%v)", cmd, cmd.GetDottedPathFull(), ff, ff.GetDottedPathFull())

	cmd, ff = root.dottedPathToCommandOrFlag("microservices.tags.toggle.address")
	t.Logf("L: cmd: %v; ff: %v", cmd, ff)
	t.Logf("L: cmd: %v (%v); ff: %v (%v)", cmd, cmd.GetDottedPathFull(), ff, ff.GetDottedPathFull())
}
