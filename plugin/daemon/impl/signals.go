/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"errors"
	"log"
	"os"
	"os/signal"
)

// ServeSignals calls handlers for system signals.
// before invoking ServeSignals(), you should run SetupSignals() at first.
func ServeSignals() (err error) {
	if len(handlers) == 0 {
		setupSignals()
	}
	if len(handlers) == 0 {
		return // no handlers, skip the os signal listening
	}

	// syscall.Getenv()
	signals := makeHandlers()

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
