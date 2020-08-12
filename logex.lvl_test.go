// Copyright © 2020 Hedzr Yeh.

package cmdr

import (
	"testing"
)

func TestLevels(t *testing.T) {
	for _, l := range AllLevels {
		t.Logf("level: %v", l)
	}

	GetLoggerLevel()
	if Level(uint32(1000)).String() != "unknown" {
		t.Fail()
	}
	_, e := Level(uint32(1000)).MarshalText()
	t.Logf("- level %q: %v", Level(uint32(1000)), e)

	var l = DebugLevel
	e = (&l).UnmarshalText([]byte("XX"))
	t.Logf("- level XX: %v", e)
	e = (&l).UnmarshalText([]byte("TRACE"))
	t.Logf("- level TRACE: %v", e)

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
		_ = internalGetWorker().getWithLogexInitializer(DebugLevel)(&rootCmdX.Command, []string{})
	}

	Set("logger.target", "journal")
	Set("logger.format", "json")
	_ = internalGetWorker().getWithLogexInitializer(DebugLevel)(&rootCmdX.Command, []string{})
}
