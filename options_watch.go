/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/hedzr/errors"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

// GetOptions returns the global options instance (rxxtOptions),
// ie. cmdr Options Store
func GetOptions() *Options {
	return internalGetWorker().rxxtOptions
}

// GetUsedConfigFile returns the main config filename (generally
// it's `<appname>.yml`)
func GetUsedConfigFile() string {
	return internalGetWorker().rxxtOptions.usedConfigFile
}

// GetUsedConfigSubDir returns the sub-directory `conf.d` of config files.
// Note that it be always normalized now.
// Sometimes it might be empty string ("") if `conf.d` have not been found.
func GetUsedConfigSubDir() string {
	return internalGetWorker().rxxtOptions.usedConfigSubDir
}

// GetUsingConfigFiles returns all loaded config files, includes
// the main config file and children in sub-directory `conf.d`.
func GetUsingConfigFiles() []string {
	return internalGetWorker().rxxtOptions.configFiles
}

// var rwlCfgReload = new(sync.RWMutex)

// AddOnConfigLoadedListener adds an functor on config loaded
// and merged
func AddOnConfigLoadedListener(c ConfigReloaded) {
	defer internalGetWorker().rxxtOptions.rwlCfgReload.Unlock()
	internalGetWorker().rxxtOptions.rwlCfgReload.Lock()

	// rwlCfgReload.RLock()
	if _, ok := internalGetWorker().rxxtOptions.onConfigReloadedFunctions[c]; ok {
		// rwlCfgReload.RUnlock()
		return
	}

	// rwlCfgReload.RUnlock()
	// rwlCfgReload.Lock()

	// defer rwlCfgReload.Unlock()

	internalGetWorker().rxxtOptions.onConfigReloadedFunctions[c] = true
}

// RemoveOnConfigLoadedListener remove an functor on config
// loaded and merged
func RemoveOnConfigLoadedListener(c ConfigReloaded) {
	w := internalGetWorker()
	defer w.rxxtOptions.rwlCfgReload.Unlock()
	w.rxxtOptions.rwlCfgReload.Lock()
	delete(w.rxxtOptions.onConfigReloadedFunctions, c)
}

// SetOnConfigLoadedListener enable/disable an functor on config
// loaded and merged
func SetOnConfigLoadedListener(c ConfigReloaded, enabled bool) {
	w := internalGetWorker()
	defer w.rxxtOptions.rwlCfgReload.Unlock()
	w.rxxtOptions.rwlCfgReload.Lock()
	w.rxxtOptions.onConfigReloadedFunctions[c] = enabled
}

// LoadConfigFile loads a yaml config file and merge the settings
// into `rxxtOptions`
// and load files in the `conf.d` child directory too.
func LoadConfigFile(file string) (err error) {
	return internalGetWorker().rxxtOptions.LoadConfigFile(file)
}

// LoadConfigFile loads a yaml config file and merge the settings
// into `rxxtOptions`
// and load files in the `conf.d` child directory too.
func (s *Options) LoadConfigFile(file string) (err error) {
	if !FileExists(file) {
		// log.Warnf("%v NOT EXISTS. PWD=%v", file, GetCurrentDir())
		return // not error, just ignore loading
	}

	if err = s.loadConfigFile(file); err != nil {
		return
	}

	s.usedConfigFile = file

	dir := path.Dir(s.usedConfigFile)
	_ = os.Setenv("CFG_DIR", dir)

	enableWatching := internalGetWorker().watchMainConfigFileToo
	dirWatch := dir
	filesWatching := []string{}
	if internalGetWorker().watchMainConfigFileToo {
		filesWatching = append(filesWatching, s.usedConfigFile)
	}

	s.usedConfigSubDir = path.Join(dir, "conf.d")
	if !FileExists(s.usedConfigSubDir) {
		s.usedConfigSubDir = ""
		if len(filesWatching) == 0 {
			return
		}
	}

	s.usedConfigSubDir, err = filepath.Abs(s.usedConfigSubDir)
	if err == nil {
		err = filepath.Walk(s.usedConfigSubDir, s.visit)
		if err == nil {
			if !internalGetWorker().watchMainConfigFileToo {
				dirWatch = s.usedConfigSubDir
			}
			filesWatching = append(filesWatching, s.configFiles...)
			enableWatching = true
		}
		// log.Fatalf("ERROR: filepath.Walk() returned %v\n", err)
	}

	if enableWatching {
		s.watchConfigDir(dirWatch, filesWatching)
	}
	return
}

// Load a yaml config file and merge the settings into `Options`
func (s *Options) loadConfigFile(file string) (err error) {
	var (
		b  []byte
		m  map[string]interface{}
		mm map[string]map[string]interface{}
	)

	b, _ = ioutil.ReadFile(file)

	m = make(map[string]interface{})
	switch path.Ext(file) {
	case ".toml", ".ini", ".conf", "toml":
		mm = make(map[string]map[string]interface{})
		err = toml.Unmarshal(b, &mm)
		if err == nil {
			err = s.loopMapMap("", mm)
		}
		if err != nil {
			return
		}
		return

	case ".json", "json":
		err = json.Unmarshal(b, &m)
	default:
		err = yaml.Unmarshal(b, &m)
	}

	if err == nil {
		err = s.loopMap("", m)
	}
	if err != nil {
		return
	}

	return
}

func (s *Options) mergeConfigFile(fr io.Reader, ext string) (err error) {
	var (
		m   map[string]interface{}
		buf *bytes.Buffer
	)

	buf = new(bytes.Buffer)
	_, err = buf.ReadFrom(fr)

	m = make(map[string]interface{})
	switch ext {
	case ".toml", ".ini", ".conf", "toml":
		err = toml.Unmarshal(buf.Bytes(), &m)
	case ".json", "json":
		err = json.Unmarshal(buf.Bytes(), &m)
	default:
		err = yaml.Unmarshal(buf.Bytes(), &m)
	}

	if err == nil {
		err = s.loopMap("", m)
	}
	if err != nil {
		return
	}

	return
}

func (s *Options) visit(path string, f os.FileInfo, e error) (err error) {
	// fmt.Printf("Visited: %s, e: %v\n", path, e)
	if f != nil && !f.IsDir() && e == nil {
		// log.Infof("    path: %v, ext: %v", path, filepath.Ext(path))
		ext := filepath.Ext(path)
		switch ext {
		case ".yml", ".yaml", ".json", ".toml", ".ini", ".conf": // , "yml", "yaml":
			var file *os.File
			file, err = os.Open(path)
			// if err != nil {
			// log.Warnf("ERROR: os.Open() returned %v", err)
			// } else {
			if err == nil {
				defer file.Close()
				if err = s.mergeConfigFile(bufio.NewReader(file), ext); err != nil {
					err = errors.New("error in merging config file '%s': %v", path, err)
					return
				}
				s.configFiles = append(s.configFiles, path)
			} else {
				err = errors.New("error in merging config file '%s': %v", path, err)
			}
		}
	} else {
		err = e
	}
	return
}

func (s *Options) reloadConfig() {
	// log.Debugf("\n\nConfig file changed: %s\n", e.String())

	defer s.rwlCfgReload.RUnlock()
	s.rwlCfgReload.RLock()

	for x, ok := range s.onConfigReloadedFunctions {
		if ok {
			x.OnConfigReloaded()
		}
	}
}

func (s *Options) watchConfigDir(configDir string, filesWatching []string) {
	if internalGetWorker().doNotWatchingConfigFiles || GetBoolR("no-watch-conf-dir") {
		return
	}

	initWG := &sync.WaitGroup{}
	initWG.Add(1)
	// initExitingChannelForFsWatcher()
	go fsWatcherRoutine(s, configDir, filesWatching, initWG)
	initWG.Wait() // make sure that the go routine above fully ended before returning
	s.SetNx("watching", true)
}

func testCfgSuffix(name string) bool {
	for _, suf := range []string{".yaml", ".yml", ".json", ".toml", ".ini", ".conf"} {
		if strings.HasSuffix(name, suf) {
			return true
		}
	}
	return false
}

func testArrayContains(s string, container []string) (contained bool) {
	for _, ss := range container {
		if ss == s {
			contained = true
			break
		}
	}
	return
}
