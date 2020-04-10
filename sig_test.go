// Copyright © 2019 Hedzr Yeh.

//+build linux

package cmdr_test

import (
	"github.com/hedzr/cmdr"
	"os"
	"testing"
	"time"
)

// func RaiseInterrupt(t *testing.T, timeout int) {
// 	go func() {
// 		time.Sleep(time.Duration(timeout) * time.Second)
// 		p, err := os.FindProcess(os.Getpid())
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		err = p.Signal(os.Interrupt)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	}()
// }

func TestTrapSignals(t *testing.T) {

	if cmdr.SavedOsArgs == nil {
		cmdr.SavedOsArgs = os.Args
	}
	defer func() {
		os.Args = cmdr.SavedOsArgs
	}()

	// RaiseInterrupt(t, 6)
	go func() {
		time.Sleep(6 * time.Second)
		cmdr.SignalQuitSignal()
	}()
	cmdr.TrapSignals(func(s os.Signal) {
		//
	})

	go func() {
		time.Sleep(6 * time.Second)
		cmdr.SignalTermSignal()
	}()
	cmdr.TrapSignals(func(s os.Signal) {
		//
	})

	_ = cmdr.RemoveDirRecursive("docs")

	// testTypes(t)
}
