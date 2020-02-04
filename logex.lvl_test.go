// Copyright © 2020 Hedzr Yeh.

package cmdr

import (
	"testing"
)

func TestLevels(t *testing.T) {
	for _, l := range AllLevels {
		t.Logf("level: %v", l)
	}

	for _, x := range []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL", "PANIC", "OFF", "XX"} {
		l, err := ParseLevel(x)
		if err == nil {
			t.Logf("level: %s => %v", x, l)
		}
	}
}

func TestLog(t *testing.T) {
	var rootCmdX = &RootCommand{
		Command: Command{
			BaseOpt: BaseOpt{
				Name: "consul-tags",
			},
		},
	}

	for _, x := range []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL", "PANIC", "OFF", "XX"} {
		Set("logger.level", x)
		_ = internalGetWorker().getWithLogexInitializor(DebugLevel)(&rootCmdX.Command, []string{})
	}

	Set("logger.target", "journal")
	Set("logger.format", "json")
	_ = internalGetWorker().getWithLogexInitializor(DebugLevel)(&rootCmdX.Command, []string{})
}
