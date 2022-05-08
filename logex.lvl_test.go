// Copyright © 2020 Hedzr Yeh.

package cmdr_test

import (
	"testing"

	"github.com/hedzr/cmdr"
)

func TestLevels(t *testing.T) {
	for _, l := range cmdr.AllLevels {
		t.Logf("level: %v", l)
	}

	cmdr.GetLoggerLevel()
	cmdr.SetLogger(cmdr.Logger)
	if cmdr.Level(uint32(1000)).String() != "unknown" {
		t.Fail()
	}
	_, e := cmdr.Level(uint32(1000)).MarshalText()
	t.Logf("- level %q: %v", cmdr.Level(uint32(1000)), e)

	l := cmdr.DebugLevel
	e = (&l).UnmarshalText([]byte("XX"))
	t.Logf("- level XX: %v", e)
	e = (&l).UnmarshalText([]byte("TRACE"))
	t.Logf("- level TRACE: %v", e)

	for _, x := range []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL", "PANIC", "OFF", "XX"} {
		l, err := cmdr.ParseLevel(x)
		if err == nil {
			t.Logf("level: %s => %v", x, l)
		}
	}
}

func TestLog(t *testing.T) {
	// var rootCmdX = &cmdr.RootCommand{
	//	Command: cmdr.Command{
	//		BaseOpt: cmdr.BaseOpt{
	//			Name: "consul-tags",
	//		},
	//	},
	// }
	//
	// for _, x := range []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL", "PANIC", "OFF", "XX"} {
	//	cmdr.Set("logger.level", x)
	//	_ = cmdr.Worker().GetWithLogexInitializer(cmdr.DebugLevel)(&rootCmdX.Command, []string{})
	// }
	//
	// cmdr.Set("logger.target", "journal")
	// cmdr.Set("logger.format", "json")
	// _ = cmdr.Worker().GetWithLogexInitializer(cmdr.DebugLevel)(&rootCmdX.Command, []string{})
}
