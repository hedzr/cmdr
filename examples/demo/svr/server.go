/*
 * Copyright © 2019 Hedzr Yeh.
 */

package svr

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/plugin/daemon"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func NewDaemon() daemon.Daemon {
	return &DaemonImpl{}
}

type DaemonImpl struct {
}

func (*DaemonImpl) OnRun(cmd *cmdr.Command, args []string, stopCh, doneCh chan struct{}) (err error) {
	logrus.Debugf("demo daemon OnRun, pid = %v, ppid = %v", os.Getpid(), os.Getppid())
	go worker(stopCh, doneCh)
	return
}

func worker(stopCh, doneCh chan struct{}) {
LOOP:
	for {
		time.Sleep(time.Second) // this is work to be done by worker.
		select {
		case <-stopCh:
			break LOOP
		default:
		}
	}
	doneCh <- struct{}{}
}

func (*DaemonImpl) OnStop(cmd *cmdr.Command, args []string) (err error) {
	logrus.Debugf("demo daemon OnStop")
	return
}

func (*DaemonImpl) OnReload() {
	logrus.Debugf("demo daemon OnReload")
}

func (*DaemonImpl) OnStatus(cxt *daemon.Context, cmd *cmdr.Command, p *os.Process) (err error) {
	fmt.Printf("%v v%v\n", cmd.GetRoot().AppName, cmd.GetRoot().Version)
	fmt.Printf("PID=%v\nLOG=%v\n", cxt.PidFileName, cxt.LogFileName)
	return
}

func (*DaemonImpl) OnInstall(cxt *daemon.Context, cmd *cmdr.Command, args []string) (err error) {
	logrus.Debugf("demo daemon OnInstall")
	return
	// panic("implement me")
}

func (*DaemonImpl) OnUninstall(cxt *daemon.Context, cmd *cmdr.Command, args []string) (err error) {
	logrus.Debugf("demo daemon OnUninstall")
	return
	// panic("implement me")
}
