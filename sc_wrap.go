// +build darwin dragonfly freebsd linux netbsd openbsd windows aix arm_linux plan9 solaris
// +build !nacl

/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)


// TrapSignal is a helper for simplify your infinite loop codes.
//
// Usage
//
//  func aLoop() {
// 	  waiter := trapSignals(func(s os.Signal) {
// 	    logrus.Debugf("receive signal '%v' in onTrapped()", s)
// 	  })
// 	  go waiter()
//  }
//
//
// 
func TrapSignals(onTrapped func(s os.Signal)) (waiter func()) {
	done := make(chan struct{}, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)

	go func() {
		s := <-sigs
		logrus.Debugf("receive signal '%v'", s)

		onTrapped(s)
		
		// for _, s := range servers {
		// 	s.Stop()
		// }
		// logrus.Infof("done")
		done <- struct{}{}
	}()
	
	waiter = func() {
		for {
			select {
			case <-done:
				// os.Exit(1)
				// logrus.Infof("done got.")
				return
			}
		}
	}
	
	return
}
