/*
 * Copyright © 2019 Hedzr Yeh.
 */

package daemon

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/sevlyar/go-daemon"
	"log"
	"os"
	"strings"
	"syscall"
)

// Enable daemon plugin:
// - add daemon commands and sub-commands: start/run, stop, restart/reload, status, install/uninstall
// - pidfile
// - go-daemon supports
// -
func Enable(daemonImplX Daemon, modifier func(daemonServerCommand *cmdr.Command) *cmdr.Command,
	preAction func(cmd *cmdr.Command, args []string) (err error),
	postAction func(cmd *cmdr.Command, args []string)) {
	daemonImpl = daemonImplX

	cmdr.AddOnBeforeXrefBuilding(func(root *cmdr.RootCommand, args []string) {

		if modifier != nil {
			root.SubCommands = append(root.SubCommands, modifier(DaemonServerCommand))
		} else {
			root.SubCommands = append(root.SubCommands, DaemonServerCommand)
		}

		prefix = strings.Join(append(cmdr.RxxtPrefix, "server"), ".")

		if root.PreAction != nil {
			savedPreAction := root.PreAction
			root.PreAction = func(cmd *cmdr.Command, args []string) (err error) {
				pidfile.Create(cmd)
				logger.Setup(cmd)
				if err = savedPreAction(cmd, args); err != nil {
					return
				}
				if preAction != nil {
					err = preAction(cmd, args)
				}
				return
			}
		} else {
			root.PreAction = func(cmd *cmdr.Command, args []string) (err error) {
				pidfile.Create(cmd)
				logger.Setup(cmd)
				if preAction != nil {
					err = preAction(cmd, args)
				}
				return
			}
		}
		if root.PostAction != nil {
			savedPostAction := root.PostAction
			root.PostAction = func(cmd *cmdr.Command, args []string) {
				if postAction != nil {
					postAction(cmd, args)
				}
				pidfile.Destroy()
				savedPostAction(cmd, args)
				return
			}
		} else {
			root.PostAction = func(cmd *cmdr.Command, args []string) {
				if postAction != nil {
					postAction(cmd, args)
				}
				pidfile.Destroy()
				return
			}
		}

	})
}

func daemonStart(cmd *cmdr.Command, args []string) (err error) {
	if cmdr.GetBoolP(prefix, "foreground") {
		err = run(cmd, args)
	} else if cmd.GetHitStr() == "run" {
		err = run(cmd, args)
	} else {
		err = runAsDaemon(cmd, args)
	}
	return
}

func runAsDaemon(cmd *cmdr.Command, args []string) (err error) {
	cxt := getContext(cmd, args)
	child, err = cxt.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if child != nil {
		fmt.Println("Daemon started. parent stopped.")
		return
	}

	defer cxt.Release()
	err = run(cmd, args)
	return
}

var child *os.Process

func run(cmd *cmdr.Command, args []string) (err error) {
	daemon.SetSigHandler(termHandler, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT)
	daemon.SetSigHandler(reloadHandler, syscall.SIGHUP)

	if daemonImpl != nil {
		if err = daemonImpl.OnRun(cmd, args, stop, done); err != nil {
			return
		}
	}

	log.Printf("daemon ServeSignals, pid = %v", os.Getpid())
	if err = daemon.ServeSignals(); err != nil {
		log.Println("Error:", err)
	}

	if daemonImpl != nil {
		err = daemonImpl.OnStop(cmd, args)
	}

	if err != nil {
		log.Fatal("daemon terminated.", err)
	}
	log.Println("daemon terminated.")
	return
}

func daemonStop(cmd *cmdr.Command, args []string) (err error) {
	getContext(cmd, args)

	p, err := daemonCtx.Search()
	if err != nil {
		fmt.Printf("%v is stopped.\n", cmd.GetRoot().AppName)
		return
	}

	if err = p.Signal(syscall.SIGTERM); err != nil {
		return
	}
	return
}

func daemonRestart(cmd *cmdr.Command, args []string) (err error) {
	getContext(cmd, args)

	p, err := daemonCtx.Search()
	if err != nil {
		fmt.Printf("%v is stopped.\n", cmd.GetRoot().AppName)
	} else {
		if err = p.Signal(syscall.SIGHUP); err != nil {
			return
		}
	}
	return
}

func daemonStatus(cmd *cmdr.Command, args []string) (err error) {
	getContext(cmd, args)

	p, err := daemonCtx.Search()
	if err != nil {
		fmt.Printf("%v is stopped.\n", cmd.GetRoot().AppName)
	} else {
		fmt.Printf("%v is running as %v.\n", cmd.GetRoot().AppName, p.Pid)
	}

	if daemonImpl != nil {
		err = daemonImpl.OnStatus(&Context{Context: *daemonCtx}, cmd, p)
	}
	return
}

func daemonInstall(cmd *cmdr.Command, args []string) (err error) {
	getContext(cmd, args)

	err = runInstaller(cmd, args)
	if err != nil {
		return
	}
	if daemonImpl != nil {
		err = daemonImpl.OnInstall(&Context{Context: *daemonCtx}, cmd, args)
	}
	return
}

func daemonUninstall(cmd *cmdr.Command, args []string) (err error) {
	getContext(cmd, args)

	err = runUninstaller(cmd, args)
	if err != nil {
		return
	}
	if daemonImpl != nil {
		err = daemonImpl.OnUninstall(&Context{Context: *daemonCtx}, cmd, args)
	}
	return
}

var prefix string
