/*
 * Copyright © 2019 Hedzr Yeh.
 */

package daemon

import (
	"os"
	"os/signal"
	"syscall"
)

var quitSignal chan os.Signal

// QuitSignals return a channel for quit signal raising up.
func QuitSignals() chan os.Signal {
	// return []os.Signal{
	// 	syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT,
	// }
	if quitSignal == nil {
		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 10 seconds.
		quitSignal = make(chan os.Signal)
		signal.Notify(quitSignal, // os.Interrupt, os.Kill, syscall.SIGHUP,
			syscall.SIGQUIT, syscall.SIGTERM,
			// syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGTSTP
			syscall.SIGABRT, syscall.SIGINT)
	}
	return quitSignal
}
