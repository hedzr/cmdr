/*
 * Copyright © 2019 Hedzr Yeh.
 */

package daemon

import (
	"bytes"
	"github.com/hedzr/cmdr/plugin/daemon/impl"
	"log"
	"os"
	"os/exec"
	"text/template"
)

func tplApply(tmpl string, data interface{}) string {
	var w = new(bytes.Buffer)
	var tpl = template.Must(template.New("y").Parse(tmpl))
	if err := tpl.Execute(w, data); err != nil {
		log.Fatalf("tpl execute error: %v", err)
	}
	return w.String()
}

func isRoot() bool {
	return os.Getuid() == 0
}

func shellRunAuto(name string, arg ...string) error {
	output, err := shellRun(name, arg...)
	if err != nil {
		log.Fatalf("shellRunAuto err: %v\n\noutput:\n%v", err, output.String())
	}
	return err
}

func shellRun(name string, arg ...string) (output bytes.Buffer, err error) {
	cmd := exec.Command(name, arg...)
	// cmd.Stdin = strings.NewReader("some input")
	cmd.Stdout = &output
	err = cmd.Run()
	return
}

// IsRunningInDemonizedMode returns true if you are running under demonized mode.
// false means that you're running in normal console/tty mode.
func IsRunningInDemonizedMode() bool {
	// return cmdr.GetBoolR(DaemonizedKey)
	return impl.IsRunningInDemonizedMode()
}

// SetTermSignals allows an functor to provide a list of Signals
func SetTermSignals(sig func() []os.Signal) {
	// onSetTermHandler = sig
	impl.SetTermSignals(sig)
}

// SetSigEmtSignals allows an functor to provide a list of Signals
func SetSigEmtSignals(sig func() []os.Signal) {
	// onSetSigEmtHandler = sig
	impl.SetSigEmtSignals(sig)
}

// SetReloadSignals allows an functor to provide a list of Signals
func SetReloadSignals(sig func() []os.Signal) {
	// onSetReloadHandler = sig
	impl.SetReloadSignals(sig)
}
