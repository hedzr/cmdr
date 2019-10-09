// +build darwin dragonfly freebsd linux netbsd openbsd windows aix arm_linux solaris
// +build !nacl

/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

func fsWatcherRoutine(s *Options, configDir string, initWG *sync.WaitGroup) {
	watcher, err := fsnotify.NewWatcher()
	if err == nil {
		defer watcher.Close()

		eventsWG := &sync.WaitGroup{}
		eventsWG.Add(1)
		go fsWatchRunner(s, configDir, watcher, eventsWG)
		_ = watcher.Add(configDir)
		initWG.Done()   // done initializing the watch in this go routine, so the parent routine can move on...
		eventsWG.Wait() // now, wait for event loop to end in this go-routine...
	}
}

func fsWatchRunner(s *Options, configDir string, watcher *fsnotify.Watcher, eventsWG *sync.WaitGroup) {
	defer func() {
		eventsWG.Done()
	}()
	for {
		select {
		case event, ok := <-watcher.Events:
			// ok == false: 'Events' channel is closed
			if ok {
				// log.Debugf("ooo | fsnotify.watcher |%v", event.String())
				// currentConfigFile, _ := filepath.EvalSymlinks(filename)
				// we only care about the config file with the following cases:
				// 1 - if the config file was modified or created
				// 2 - if the real path to the config file changed (eg: k8s ConfigMap replacement)
				const writeOrCreateMask = fsnotify.Write | fsnotify.Create
				if strings.HasPrefix(filepath.Clean(event.Name), configDir) &&
					event.Op&writeOrCreateMask != 0 &&
					(testCfgSuffix(event.Name)) {
					file, err := os.Open(event.Name)
					if err != nil {
						log.Printf("ERROR: os.Open() returned %v\n", err)
					} else {
						err = s.mergeConfigFile(bufio.NewReader(file), path.Ext(event.Name))
						if err != nil {
							log.Printf("ERROR: os.Open() returned %v\n", err)
						}
						s.reloadConfig()
						file.Close()
					}
				}
			}

		case err, ok := <-watcher.Errors:
			if ok { // 'Errors' channel is not closed
				// log.Printf("watcher error: %v\n", err)
				log.Printf("Watcher error: %v\n", err)
			}
			return
		}
	}
}
