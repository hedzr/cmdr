/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"net"
	"os"
)

// IsRunningInDemonizedMode returns true if you are running under demonized mode.
// false means that you're running in normal console/tty mode.
func IsRunningInDemonizedMode() bool {
	// return cmdr.GetBoolR(envvarInDaemonized)
	return isDemonized()
}

// SetTermSignals allows an functor to provide a list of Signals
func SetTermSignals(sig func() []os.Signal) {
	onSetTermHandler = sig
}

// SetSigEmtSignals allows an functor to provide a list of Signals
func SetSigEmtSignals(sig func() []os.Signal) {
	onSetSigEmtHandler = sig
}

// SetReloadSignals allows an functor to provide a list of Signals
func SetReloadSignals(sig func() []os.Signal) {
	onSetReloadHandler = sig
}

// SetHotReloadSignals allows an functor to provide a list of Signals
func SetHotReloadSignals(sig func() []os.Signal) {
	onSetHotReloadHandler = sig
}

// SetOnGetListener returns tcp/http listener for daemon hot-restarting
func SetOnGetListener(fn func() net.Listener) {
	onGetListener = fn
}

// SetSigHandler sets handler for the given signals.
// SIGTERM has the default handler, he returns ErrStop.
func SetSigHandler(handler SignalHandlerFunc, signals ...os.Signal) {
	for _, sig := range signals {
		handlers[sig] = handler
	}
}
