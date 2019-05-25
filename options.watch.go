/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"bytes"
	"github.com/fsnotify/fsnotify"
	"log"
	"path"

	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// GetUsedConfigFile returns the main config filename (generally it's `<appname>.yml`)
func GetUsedConfigFile() string {
	return usedConfigFile
}

// GetUsedConfigSubDir returns the sub-directory `conf.d` of config files
func GetUsedConfigSubDir() string {
	return usedConfigSubDir
}

// AddOnConfigLoadedListener add an functor on config loaded and merged
func AddOnConfigLoadedListener(c ConfigReloaded) {
	if _, ok := onConfigReloadedFunctions[c]; ok {
		//
	} else {
		onConfigReloadedFunctions[c] = true
	}
}

// RemoveOnConfigLoadedListener remove an functor on config loaded and merged
func RemoveOnConfigLoadedListener(c ConfigReloaded) {
	delete(onConfigReloadedFunctions, c)
}

// SetOnConfigLoadedListener enable/disable an functor on config loaded and merged
func SetOnConfigLoadedListener(c ConfigReloaded, enabled bool) {
	onConfigReloadedFunctions[c] = enabled
}

// LoadConfigFile Load a yaml config file and merge the settings into `rxxtOptions`
// and load files in the `conf.d` child directory too.
func LoadConfigFile(file string) (err error) {
	return rxxtOptions.LoadConfigFile(file)
}

// LoadConfigFile Load a yaml config file and merge the settings into `rxxtOptions`
// and load files in the `conf.d` child directory too.
func (s *Options) LoadConfigFile(file string) (err error) {
	if !FileExists(file) {
		// log.Warnf("%v NOT EXISTS. PWD=%v", file, GetCurrentDir())
		return // not error, just ignore loading
	}

	if err = s.loadConfigFile(file); err != nil {
		return
	}

	usedConfigFile = file

	usedConfigSubDir = path.Join(path.Dir(usedConfigFile), "conf.d")
	if !FileExists(usedConfigSubDir) {
		usedConfigSubDir = ""
		return
	}

	if err = filepath.Walk(usedConfigSubDir, s.visit); err != nil {
		log.Fatalf("ERROR: filepath.Walk() returned %v\n", err)
	}

	s.watchConfigDir(usedConfigSubDir)

	return
}

// Load a yaml config file and merge the settings into `Options`
func (s *Options) loadConfigFile(file string) (err error) {
	var (
		b []byte
		m map[string]interface{}
	)

	b, err = ioutil.ReadFile(file)
	if err != nil {
		return
	}

	m = make(map[string]interface{})
	err = yaml.Unmarshal(b, m)
	if err != nil {
		return
	}

	err = s.loopMap("", m)
	if err != nil {
		return
	}

	return
}

func (s *Options) mergeConfigFile(fr io.Reader) (err error) {
	var (
		m   map[string]interface{}
		buf *bytes.Buffer
	)

	buf = new(bytes.Buffer)
	_, err = buf.ReadFrom(fr)

	m = make(map[string]interface{})
	err = yaml.Unmarshal(buf.Bytes(), m)
	if err != nil {
		return
	}

	err = s.loopMap("", m)
	if err != nil {
		return
	}

	return
}

func (s *Options) visit(path string, f os.FileInfo, e error) (err error) {
	// fmt.Printf("Visited: %s, e: %v\n", path, e)
	if f != nil && !f.IsDir() && e == nil {
		// log.Infof("    path: %v, ext: %v", path, filepath.Ext(path))
		switch filepath.Ext(path) {
		case ".yml", ".yaml": // , "yml", "yaml":
			var file *os.File
			file, err = os.Open(path)
			if err != nil {
				// log.Warnf("ERROR: os.Open() returned %v", err)
			} else {
				defer file.Close()
				err = s.mergeConfigFile(bufio.NewReader(file))
				configFiles = append(configFiles, path)
				// env := viper.Get("app.registrar.env")
				// key := fmt.Sprintf("app.registrar.consul.%s.addr", env)
				// log.Infof("%s = %s", key, viper.Get(key))
			}
		}
	} else {
		err = e
	}
	return err
}

func (s *Options) reloadConfig(e fsnotify.Event) {
	// log.Debugf("\n\nConfig file changed: %s\n", e.String())

	for x, ok := range onConfigReloadedFunctions {
		if ok {
			x.OnConfigReloaded()
		}
	}
}

func (s *Options) watchConfigDir(configDir string) {
	initWG := sync.WaitGroup{}
	initWG.Add(1)
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		eventsWG := sync.WaitGroup{}
		eventsWG.Add(1)
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok { // 'Events' channel is closed
						eventsWG.Done()
						return
					}
					// log.Debugf("ooo | fsnotify.watcher |%v", event.String())
					// currentConfigFile, _ := filepath.EvalSymlinks(filename)
					// we only care about the config file with the following cases:
					// 1 - if the config file was modified or created
					// 2 - if the real path to the config file changed (eg: k8s ConfigMap replacement)
					const writeOrCreateMask = fsnotify.Write | fsnotify.Create
					if strings.HasPrefix(filepath.Clean(event.Name), configDir) &&
						event.Op&writeOrCreateMask != 0 &&
						(strings.HasSuffix(event.Name, ".yml") || strings.HasSuffix(event.Name, ".yaml")) {
						file, err := os.Open(event.Name)
						if err != nil {
							log.Printf("ERROR: os.Open() returned %v\n", err)
						} else {
							err = s.mergeConfigFile(bufio.NewReader(file))
							if err != nil {
								log.Printf("ERROR: os.Open() returned %v\n", err)
							}
							s.reloadConfig(event)
							file.Close()
						}
					}

				case err, ok := <-watcher.Errors:
					if ok { // 'Errors' channel is not closed
						// log.Printf("watcher error: %v\n", err)
						log.Printf("Watcher error: %v\n", err)
					}
					eventsWG.Done()
					return
				}
			}
		}()
		_ = watcher.Add(configDir)
		initWG.Done()   // done initializing the watch in this go routine, so the parent routine can move on...
		eventsWG.Wait() // now, wait for event loop to end in this go-routine...
	}()
	initWG.Wait() // make sure that the go routine above fully ended before returning
}
