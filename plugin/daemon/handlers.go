/*
 * Copyright © 2019 Hedzr Yeh.
 */

package daemon

// import (
// 	"github.com/sevlyar/go-daemon"
// 	"log"
// 	"os"
// 	"syscall"
// )
//
// var (
// 	stop = make(chan struct{})
// 	done = make(chan struct{})
// )
//
// func termHandler(sig os.Signal) error {
// 	log.Println("terminating...")
// 	stop <- struct{}{}
// 	if sig == syscall.SIGQUIT {
// 		<-done
// 	}
// 	return daemon.ErrStop
// }
//
// func sigEmtHandler(sig os.Signal) error {
// 	log.Println("terminating (SIGEMT)...")
// 	stop <- struct{}{}
// 	if sig == syscall.SIGQUIT {
// 		<-done
// 	}
// 	return daemon.ErrStop
// }
//
// func reloadHandler(sig os.Signal) error {
// 	log.Println("configuration reloaded")
// 	if daemonImpl != nil {
// 		daemonImpl.OnReload()
// 	}
// 	return nil
// }
