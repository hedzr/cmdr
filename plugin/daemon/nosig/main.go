/*
 * Copyright © 2019 Hedzr Yeh.
 */

package main

import (
	"fmt"
	"github.com/hedzr/cmdr/flag"
	"github.com/hedzr/cmdr/plugin/daemon/impl"
	"log"
	"os"
	"time"
)

func main() {
	fmt.Printf("-- args: %v\n", os.Args)

	flag.Parse() // enable cmdr working table

	daemonContext := impl.DefaultContext()

	// https://blog.csdn.net/fyxichen/article/details/50541449
	// daemon(false, true)
	// BAD: daemonNew()
	if err := impl.Demonize(daemonContext); err != nil {
		log.Printf("Unable to create child process: %+v", err)
		os.Exit(impl.ErrnoForkAndDaemonFailed)
	}

	log.Println("Keep going...")

	// var doneCh, exitCh = make(chan struct{}), make(chan struct{})
	go func() {
		log.Println("routine running...")
		defer func() {
			log.Println("routine stopped.")
		}()

		counter := 0
		for {
			if impl.HandleSignalCaughtEvent() {
				break
			}

			fmt.Println("hello ", counter)
			counter++
			time.Sleep(2 * time.Second)

			// select {
			// case <-exitCh:
			// 	doneCh <- struct{}{}
			// 	return
			// default:
			// 	fmt.Println("hello")
			// 	time.Sleep(1 * time.Second)
			// }
		}
	}()

	log.Println("For Signals...")
	time.Sleep(5 * time.Second)

	if err := impl.ServeSignals(daemonContext); err != nil {
		log.Printf("error at ServeSignals: %+v", err)
	}
	// exitCh <- struct{}{}
	// <-doneCh
	log.Println("DONE.")
}
