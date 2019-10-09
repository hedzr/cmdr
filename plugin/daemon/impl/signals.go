/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"errors"
	"log"
	"os"
	"os/signal"
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

	signals := []os.Signal{syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT, syscall.SIGKILL, syscall.SIGUSR1, syscall.SIGUSR2}
	if onSetTermHandler != nil {
		signals = onSetTermHandler()
	}
	SetSigHandler(termHandler, signals...)

	signals = []os.Signal{syscall.Signal(0x7)}
	if onSetSigEmtHandler != nil {
		signals = onSetSigEmtHandler()
	}
	SetSigHandler(sigEmtHandler, signals...)

	signals = []os.Signal{syscall.SIGHUP}
	if onSetReloadHandler != nil {
		signals = onSetReloadHandler()
	}
	SetSigHandler(reloadHandler, signals...)
}

// ServeSignals calls handlers for system signals.
// before invoking ServeSignals(), you should run SetupSignals() at first.
func ServeSignals() (err error) {
	if len(handlers) == 0 {
		setupSignals()
	}
	if len(handlers) == 0 {
		return // no handlers, skip the os signal listening
	}

	signals := make([]os.Signal, 0, len(handlers))
	for sig := range handlers {
		signals = append(signals, sig)
	}

	defer func() {
		removePID()
	}()
	
	ch := make(chan os.Signal, 8)
	signal.Notify(ch, signals...)

	for sig := range ch {
		log.Printf(".. signal caught: %v", sig)
		err = handlers[sig](sig)
		if err != nil {
			break
		}
	}

	signal.Stop(ch)

	if err == ErrStop {
		err = nil
	}

	return
}

func HandleSignalCaughtEvent() bool {
	select {
	case <-stop:
		log.Print("stop ch received. send done ch.")
		done <- struct{}{}
		return true
	default:
		return false
	}
}

func termHandler(sig os.Signal) error {
	log.Println("terminating...")
	stop <- struct{}{}
	log.Println("  - send stop ch")
	if sig == syscall.SIGQUIT {
		log.Println("  - waiting for done ch...")
		<-done
		log.Println("  - done ch received.")
	}
	return ErrStop
}

func sigEmtHandler(sig os.Signal) error {
	log.Println("terminating (SIGEMT)...")
	stop <- struct{}{}
	if sig == syscall.SIGQUIT {
		<-done
	}
	return ErrStop
}

func reloadHandler(sig os.Signal) error {
	log.Println("configuration reloaded")
	// if daemonImpl != nil {
	// 	daemonImpl.OnReload()
	// }
	return nil
}

var (
	// ErrStop should be returned signal handler function
	// for termination of handling signals.
	ErrStop = errors.New("stop serve signals")

	handlers = make(map[os.Signal]SignalHandlerFunc)

	child *os.Process

	onSetTermHandler   func() []os.Signal
	onSetSigEmtHandler func() []os.Signal
	onSetReloadHandler func() []os.Signal

	stop = make(chan struct{})
	done = make(chan struct{})
)

const (
	ErrnoForkAndDaemonFailed = -1
	envvarInDaemonized       = "__DAEMON"
)

// SignalHandlerFunc is the interface for signal handler functions.
type SignalHandlerFunc func(sig os.Signal) (err error)
