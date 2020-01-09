// Copyright © 2019 Hedzr Yeh.

package impl

import (
	"log"
	"os"
	"syscall"
)

// // SetupSignals initialize all signal handlers
// func SetupSignals() {
// 	setupSignals()
// }

func setupSignals() {
	// for i := 1; i < 34; i++ {
	// 	daemon.SetSigHandler(termHandler, syscall.Signal(i))
	// }

	// signals := []os.Signal{syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT, syscall.SIGKILL, syscall.SIGUSR1, syscall.SIGUSR2}
	// if onSetTermHandler != nil {
	// 	signals = onSetTermHandler()
	// }
	// SetSigHandler(termHandler, signals...)
	//
	// signals = []os.Signal{syscall.Signal(0x7)}
	// if onSetSigEmtHandler != nil {
	// 	signals = onSetSigEmtHandler()
	// }
	// SetSigHandler(sigEmtHandler, signals...)
	//
	// signals = []os.Signal{syscall.SIGHUP}
	// if onSetReloadHandler != nil {
	// 	signals = onSetReloadHandler()
	// }
	// SetSigHandler(reloadHandler, signals...)
}

func termHandler(sig os.Signal) error {
	log.Println("terminating...")
	stop <- struct{}{}
	// log.Println("  - send stop ch")
	// if sig == syscall.SIGQUIT {
	// 	log.Println("  - waiting for done ch...")
	// 	<-done
	// 	log.Println("  - done ch received.")
	// }
	return ErrStop
}

func sigEmtHandler(sig os.Signal) error {
	log.Println("terminating (SIGEMT)...")
	stop <- struct{}{}
	// if sig == syscall.SIGQUIT {
	// 	<-done
	// }
	return ErrStop
}

func makeHandlers() (signals []os.Signal) {
	signals = make([]os.Signal, 0, len(handlers))
	for sig := range handlers {
		signals = append(signals, sig)
	}
	return
}

func nilSigSend(process *os.Process) error {
	return process.Signal(syscall.Signal(0))
	// return nil
}

func sigSendHUP(process *os.Process) error {
	return process.Signal(syscall.SIGHUP)
}

func sigSendUSR1(process *os.Process) error {
	return process.Signal(syscall.SIGUSR1)
}

func sigSendUSR2(process *os.Process) error {
	return process.Signal(syscall.SIGUSR2)
}

func sigSendTERM(process *os.Process) error {
	return process.Signal(syscall.SIGTERM)
}
