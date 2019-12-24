// Copyright © 2019 Hedzr Yeh.

package cmdr

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// NewOptions returns an `Options` structure pointer
func NewOptions() *Options {
	return &Options{
		entries:   make(map[string]interface{}),
		hierarchy: make(map[string]interface{}),
		rw:        new(sync.RWMutex),

		onConfigReloadedFunctions: make(map[ConfigReloaded]bool),
		rwlCfgReload:              new(sync.RWMutex),
	}
}

// NewOptionsWith returns an `Options` structure pointer
func NewOptionsWith(entries map[string]interface{}) *Options {
	return &Options{
		entries:   entries,
		hierarchy: make(map[string]interface{}),
		rw:        new(sync.RWMutex),

		onConfigReloadedFunctions: make(map[ConfigReloaded]bool),
		rwlCfgReload:              new(sync.RWMutex),
	}
}

// Has detects whether a key exists in cmdr options store or not
func (s *Options) Has(key string) (ok bool) {
	defer s.rw.RUnlock()
	s.rw.RLock()
	_, ok = s.entries[key]
	return
}

// DeleteKey deletes a key from cmdr options store
func DeleteKey(key string) {
	internalGetWorker().rxxtOptions.Delete(key)
}

// Delete deletes a key from cmdr options store
func (s *Options) Delete(key string) {
	defer s.rw.RUnlock()
	s.rw.RLock()

	val := s.entries[key]
	a := strings.Split(key, ".")
	s.deleteWithKey(s.hierarchy, a[0], "", et(a, 1, val))
	return
}

func (s *Options) deleteWithKey(m map[string]interface{}, key, path string, val interface{}) (ret map[string]interface{}) {
	if len(path) > 0 {
		path = fmt.Sprintf("%v.%v", path, key)
	} else {
		path = key
	}

	if z, ok := m[key]; ok {
		if zm, ok := z.(map[string]interface{}); ok {
			if vm, ok := val.(map[string]interface{}); ok {
				for k, v := range vm {
					zm = s.deleteWithKey(zm, k, path, v)
				}
				// delete(m, key)
				// delete(s.entries, path)
				return
			} else if vm, ok := val.(map[interface{}]interface{}); ok {
				for k, v := range vm {
					kk, ok := k.(string)
					if !ok {
						kk = fmt.Sprintf("%v", k)
					}
					zm = s.deleteWithKey(zm, kk, path, v)
				}
				// delete(m, key)
				// delete(s.entries, path)
				return
			}
		}
	}

	delete(m, key)
	delete(s.entries, path)
	return
}

// Get an `Option` by key string, eg:
// ```golang
// cmdr.Get("app.logger.level") => 'DEBUG',...
// ```
//
func (s *Options) Get(key string) interface{} {
	defer s.rw.RUnlock()
	s.rw.RLock()
	return s.entries[key]
}

// GetMap an `Option` by key string, it returns a hierarchy map or nil
func (s *Options) GetMap(key string) map[string]interface{} {
	defer s.rw.RUnlock()
	s.rw.RLock()

	return s.getMapNoLock(key)
}

func (s *Options) getMapNoLock(key string) (m map[string]interface{}) {
	a := strings.Split(key, ".")
	if len(a) > 0 {
		m = s.getMap(s.hierarchy, a[0], a[1:]...)
	}
	return
}

func (s *Options) getMap(vp map[string]interface{}, key string, remains ...string) map[string]interface{} {
	if len(remains) > 0 {
		if v, ok := vp[key]; ok {
			if vm, ok := v.(map[string]interface{}); ok {
				return s.getMap(vm, remains[0], remains[1:]...)
			}
		}
		return nil
	}

	if v, ok := vp[key]; ok {
		if vm, ok := v.(map[string]interface{}); ok {
			return vm
		}
		return vp
	}
	return nil
}

// GetBoolEx returns the bool value of an `Option` key.
func (s *Options) GetBoolEx(key string, defaultVal ...bool) (ret bool) {
	switch strings.ToLower(s.GetString(key, "")) {
	case "1", "y", "t", "yes", "true", "ok", "on":
		ret = true
	case "":
		for _, vv := range defaultVal {
			ret = vv
		}
	}
	return
}

// GetIntEx returns the int64 value of an `Option` key.
func (s *Options) GetIntEx(key string, defaultVal ...int) (ir int) {
	if ir64, err := strconv.ParseInt(s.GetString(key, ""), 10, 64); err == nil {
		ir = int(ir64)
	} else {
		for _, vv := range defaultVal {
			ir = vv
		}
	}
	return
}

// GetInt64Ex returns the int64 value of an `Option` key.
func (s *Options) GetInt64Ex(key string, defaultVal ...int64) (ir int64) {
	if ir64, err := strconv.ParseInt(s.GetString(key, ""), 10, 64); err == nil {
		ir = ir64
	} else {
		for _, vv := range defaultVal {
			ir = vv
		}
	}
	return
}

// GetUintEx returns the uint64 value of an `Option` key.
func (s *Options) GetUintEx(key string, defaultVal ...uint) (ir uint) {
	if ir64, err := strconv.ParseUint(s.GetString(key, ""), 10, 64); err == nil {
		ir = uint(ir64)
	} else {
		for _, vv := range defaultVal {
			ir = vv
		}
	}
	return
}

// GetUint64Ex returns the uint64 value of an `Option` key.
func (s *Options) GetUint64Ex(key string, defaultVal ...uint64) (ir uint64) {
	if ir64, err := strconv.ParseUint(s.GetString(key, ""), 10, 64); err == nil {
		ir = ir64
	} else {
		for _, vv := range defaultVal {
			ir = vv
		}
	}
	return
}

// GetFloat32Ex returns the float32 value of an `Option` key.
func (s *Options) GetFloat32Ex(key string, defaultVal ...float32) (ir float32) {
	if ir64, err := strconv.ParseFloat(s.GetString(key, ""), 10); err == nil {
		ir = float32(ir64)
	} else {
		for _, vv := range defaultVal {
			ir = vv
		}
	}
	return
}

// GetFloat64Ex returns the float64 value of an `Option` key.
func (s *Options) GetFloat64Ex(key string, defaultVal ...float64) (ir float64) {
	if ir64, err := strconv.ParseFloat(s.GetString(key, ""), 10); err == nil {
		ir = ir64
	} else {
		for _, vv := range defaultVal {
			ir = vv
		}
	}
	return
}

// GetStringSlice returns the string slice value of an `Option` key.
func (s *Options) GetStringSlice(key string, defaultVal ...string) (ir []string) {
	// envkey := s.envKey(key)
	// if s, ok := os.LookupEnv(envkey); ok {
	// 	ir = strings.Split(s, ",")
	// }

	defer s.rw.RUnlock()
	s.rw.RLock()

	if v, ok := s.entries[key]; ok {
		vvv := reflect.ValueOf(v)
		switch vvv.Kind() {
		case reflect.String:
			ir = strings.Split(v.(string), ",")
		case reflect.Slice:
			if r, ok := v.([]string); ok {
				ir = r
			} else if ri, ok := v.([]int); ok {
				for _, rii := range ri {
					ir = append(ir, strconv.Itoa(rii))
				}
			} else if ri, ok := v.([]byte); ok {
				ir = strings.Split(string(ri), ",")
			} else {
				for i := 0; i < vvv.Len(); i++ {
					ir = append(ir, fmt.Sprintf("%v", vvv.Index(i).Interface()))
				}
			}
		default:
			ir = strings.Split(fmt.Sprintf("%v", v), ",")
		}
	} else {
		ir = defaultVal
	}
	return
}

// GetIntSlice returns the string slice value of an `Option` key.
func (s *Options) GetIntSlice(key string, defaultVal ...int) (ir []int) {
	// envkey := s.envKey(key)
	// if s, ok := os.LookupEnv(envkey); ok {
	// 	ir = stringSliceToIntSlice(strings.Split(s, ","))
	// }

	defer s.rw.RUnlock()
	s.rw.RLock()

	if v, ok := s.entries[key]; ok {
		vvv := reflect.ValueOf(v)
		switch vvv.Kind() {
		case reflect.String:
			ir = stringSliceToIntSlice(strings.Split(v.(string), ","))
		case reflect.Slice:
			if r, ok := v.([]string); ok {
				ir = stringSliceToIntSlice(r)
			} else if ri, ok := v.([]int); ok {
				ir = ri
			} else if ri, ok := v.([]int64); ok {
				ir = int64SliceToIntSlice(ri)
			} else if ri, ok := v.([]uint64); ok {
				ir = uint64SliceToIntSlice(ri)
			} else if ri, ok := v.([]byte); ok {
				xx := strings.Split(string(ri), ",")
				ir = stringSliceToIntSlice(xx)
			} else {
				var xx []string
				for i := 0; i < vvv.Len(); i++ {
					xx = append(xx, fmt.Sprintf("%v", vvv.Index(i).Interface()))
				}
				ir = stringSliceToIntSlice(xx)
			}
		default:
			ir = stringSliceToIntSlice(strings.Split(fmt.Sprintf("%v", v), ","))
		}
	} else {
		ir = defaultVal
	}
	return
}

// GetInt64Slice returns the string slice value of an `Option` key.
func (s *Options) GetInt64Slice(key string, defaultVal ...int64) (ir []int64) {
	// envkey := s.envKey(key)
	// if s, ok := os.LookupEnv(envkey); ok {
	// 	ir = stringSliceToIntSlice(strings.Split(s, ","))
	// }

	defer s.rw.RUnlock()
	s.rw.RLock()

	if v, ok := s.entries[key]; ok {
		vvv := reflect.ValueOf(v)
		switch vvv.Kind() {
		case reflect.String:
			ir = stringSliceToInt64Slice(strings.Split(v.(string), ","))
		case reflect.Slice:
			if r, ok := v.([]string); ok {
				ir = stringSliceToInt64Slice(r)
			} else if ri, ok := v.([]int); ok {
				ir = intSliceToInt64Slice(ri)
			} else if ri, ok := v.([]int64); ok {
				ir = ri
			} else if ri, ok := v.([]uint64); ok {
				ir = uint64SliceToInt64Slice(ri)
			} else if ri, ok := v.([]byte); ok {
				xx := strings.Split(string(ri), ",")
				ir = stringSliceToInt64Slice(xx)
			} else {
				var xx []string
				for i := 0; i < vvv.Len(); i++ {
					xx = append(xx, fmt.Sprintf("%v", vvv.Index(i).Interface()))
				}
				ir = stringSliceToInt64Slice(xx)
			}
		default:
			ir = stringSliceToInt64Slice(strings.Split(fmt.Sprintf("%v", v), ","))
		}
	} else {
		ir = defaultVal
	}
	return
}

// GetUint64Slice returns the string slice value of an `Option` key.
func (s *Options) GetUint64Slice(key string, defaultVal ...uint64) (ir []uint64) {
	// envkey := s.envKey(key)
	// if s, ok := os.LookupEnv(envkey); ok {
	// 	ir = stringSliceToIntSlice(strings.Split(s, ","))
	// }

	defer s.rw.RUnlock()
	s.rw.RLock()

	if v, ok := s.entries[key]; ok {
		vvv := reflect.ValueOf(v)
		switch vvv.Kind() {
		case reflect.String:
			ir = stringSliceToUint64Slice(strings.Split(v.(string), ","))
		case reflect.Slice:
			if r, ok := v.([]string); ok {
				ir = stringSliceToUint64Slice(r)
			} else if ri, ok := v.([]int); ok {
				ir = intSliceToUint64Slice(ri)
			} else if ri, ok := v.([]int64); ok {
				ir = int64SliceToUint64Slice(ri)
			} else if ri, ok := v.([]uint64); ok {
				ir = ri
			} else if ri, ok := v.([]byte); ok {
				xx := strings.Split(string(ri), ",")
				ir = stringSliceToUint64Slice(xx)
			} else {
				var xx []string
				for i := 0; i < vvv.Len(); i++ {
					xx = append(xx, fmt.Sprintf("%v", vvv.Index(i).Interface()))
				}
				ir = stringSliceToUint64Slice(xx)
			}
		default:
			ir = stringSliceToUint64Slice(strings.Split(fmt.Sprintf("%v", v), ","))
		}
	}
	return
}

// GetDuration returns the time duration value of an `Option` key.
func (s *Options) GetDuration(key string, defaultVal ...time.Duration) (ir time.Duration) {
	str := s.GetString(key, "BAD")
	if str == "BAD" {
		for _, vv := range defaultVal {
			ir = vv
		}
	} else {
		var err error
		if ir, err = time.ParseDuration(str); err != nil {
			for _, vv := range defaultVal {
				ir = vv
			}
		}
	}
	return
}

// GetString returns the string value of an `Option` key.
func (s *Options) GetString(key string, defaultVal ...string) (ret string) {
	ret = s.GetStringNoExpand(key, defaultVal...)
	ret = os.ExpandEnv(ret)
	return
}

// GetStringNoExpand returns the string value of an `Option` key.
func (s *Options) GetStringNoExpand(key string, defaultVal ...string) (ret string) {
	// envkey := s.envKey(key)
	// if s, ok := os.LookupEnv(envkey); ok {
	// 	ret = s
	// }

	defer s.rw.RUnlock()
	s.rw.RLock()

	if v, ok := s.entries[key]; ok {
		switch reflect.ValueOf(v).Kind() {
		case reflect.String:
			ret = v.(string)
			if len(ret) == 0 {
				for _, v := range defaultVal {
					ret = v
				}
			}
		default:
			if v != nil {
				ret = fmt.Sprint(v)
			}
		}
	} else {
		for _, vv := range defaultVal {
			ret = vv
		}
	}
	return
}

func (s *Options) buildAutomaticEnv(rootCmd *RootCommand) (err error) {
	// logrus.SetLevel(logrus.DebugLevel)

	// prefix := strings.Join(EnvPrefix,"_")
	prefix := internalGetWorker().getPrefix() // strings.Join(RxxtPrefix, ".")
	for key := range s.entries {
		ek := s.envKey(key)
		if v, ok := os.LookupEnv(ek); ok {
			if strings.HasPrefix(key, prefix) {
				s.Set(key[len(prefix)+1:], v)
			} else {
				s.Set(key, v)
			}
		}
		// logrus.Printf("buildAutomaticEnv: %v", key)
		if flg := s.lookupFlag(key, rootCmd); flg != nil {
			// // logrus.Debugf("buildAutomaticEnv: %v matched", key)
			// if key == "app.mx-test.test" {
			// 	logrus.Debugf("                 : flag=%+v", flg)
			// }
			for _, ek := range flg.EnvVars {
				if v, ok := os.LookupEnv(ek); ok {
					// logrus.Debugf("buildAutomaticEnv: envvar %v found", ek)
					// logrus.Debugf("                 : flag=%+v", flg)
					if strings.HasPrefix(key, prefix) {
						// logrus.Printf("setnx: %v <-- %v", key, v)
						s.SetNx(key, v)
						// logrus.Printf("setnx: %v", s.GetString(key))
					} else {
						// logrus.Printf("set: %v <-- %v", key, v)
						s.Set(key, v)
					}
				}
			}
		}
	}

	// // fmt.Printf("EXE = %v, PWD = %v, CURRDIR = %v\n", GetExecutableDir(), os.Getenv("PWD"), GetCurrentDir())
	// // _ = os.Setenv("THIS", GetExecutableDir())
	// for k, v := range uniqueWorker.envvarToValueMap {
	// 	_ = os.Setenv(k, v())
	// }
	internalGetWorker().setupFromEnvvarMap()

	for _, h := range internalGetWorker().afterAutomaticEnv {
		h(rootCmd, s)
	}
	return
}

func (s *Options) lookupFlag(keyPath string, rootCmd *RootCommand) (flg *Flag) {
	flg = s.loopForLookupFlag(strings.Split(keyPath, ".")[len(internalGetWorker().envPrefixes):], &rootCmd.Command)
	return
}

func (s *Options) loopForLookupFlag(keys []string, cmd *Command) (flg *Flag) {
	switch len(keys) {
	case 0:
		return
	case 1:
		for _, f := range cmd.Flags {
			if f.Full == keys[0] {
				flg = f
				return
			}
		}
	default:
		tmpkeys := keys[1:]
		for _, sc := range cmd.SubCommands {
			if flg = s.loopForLookupFlag(tmpkeys, sc); flg != nil {
				return
			}
		}
	}
	return
}

func (s *Options) envKey(key string) (envkey string) {
	key = strings.ReplaceAll(key, ".", "_")
	key = strings.ReplaceAll(key, "-", "_")
	envkey = strings.Join(append(internalGetWorker().envPrefixes, strings.ToUpper(key)), "_")
	return
}

// Set set the value of an `Option` key. The key MUST not have an `app` prefix. eg:
// ```golang
// cmdr.Set("debug", true)
// cmdr.GetBool("app.debug") => true
// ```
func (s *Options) Set(key string, val interface{}) {
	k := wrapWithRxxtPrefix(key)
	s.setNx(k, val)
}

// SetNx but without prefix auto-wrapped.
// `rxxtPrefix` is a string slice to define the prefix string array, default is ["app"].
// So, cmdr.Set("debug", true) will put an real entry with (`app.debug`, true).
func (s *Options) SetNx(key string, val interface{}) {
	s.setNx(key, val)
}

func (s *Options) setNx(key string, val interface{}) (oldval interface{}, modi bool) {
	defer s.rw.Unlock()
	s.rw.Lock()

	if val == nil {
		if s.getMapNoLock(key) != nil {
			// don't set a branch node to nil if it have children.
			return
		}
	}

	oldval = s.entries[key]
	var leaf bool
	if _, ok := oldval.(map[string]interface{}); !ok {
		if _, ok := val.(map[string]interface{}); !ok {
			leaf = true
		}
	}
	if leaf {
		comparable := (oldval == nil || oldval != nil && reflect.TypeOf(oldval).Comparable()) && (val == nil || (val != nil && reflect.TypeOf(val).Comparable()))
		if comparable {
			if oldval != val {
				s.entries[key] = val
				a := strings.Split(key, ".")
				s.mergeMap(s.hierarchy, a[0], "", et(a, 1, val))
				if s.onSet != nil {
					s.onSet(key, val, oldval)
				}
				modi = true
				return
			}
		}
	}
	s.entries[key] = val
	return
}

// MergeWith will merge a map recursive.
func (s *Options) MergeWith(m map[string]interface{}) (err error) {
	defer s.rw.Unlock()
	s.rw.Lock()
	for k, v := range m {
		s.mergeMap(s.hierarchy, k, "", v)
	}
	return
}

func (s *Options) mergeMap(hierarchy map[string]interface{}, key, path string, val interface{}) map[string]interface{} {
	if len(path) > 0 {
		path = fmt.Sprintf("%v.%v", path, key)
	} else {
		path = key
	}

	if z, ok := hierarchy[key]; ok {
		if zm, ok := z.(map[string]interface{}); ok {
			if vm, ok := val.(map[string]interface{}); ok {
				for k, v := range vm {
					zm = s.mergeMap(zm, k, path, v)
				}
				// hierarchy[key] = zm
				// s.entries[path] = zm
				val = zm
			} else if vm, ok := val.(map[interface{}]interface{}); ok {
				for k, v := range vm {
					kk, ok := k.(string)
					if !ok {
						kk = fmt.Sprintf("%v", k)
					}
					zm = s.mergeMap(zm, kk, path, v)
				}
				// hierarchy[key] = zm
				// s.entries[path] = zm
				val = zm
				// } else {
				// 	hierarchy[key] = val
				// 	s.entries[path] = val
			}
			// } else {
			// 	hierarchy[key] = val
			// 	s.entries[path] = val
		}
		// } else {
		// 	hierarchy[key] = val
		// 	s.entries[path] = val
	}

	s.mmset(hierarchy, key, path, val)
	return hierarchy
}

func (s *Options) mmset(m map[string]interface{}, key, path string, val interface{}) {
	oldval := s.entries[path]

	var leaf bool
	if _, ok := oldval.(map[string]interface{}); !ok {
		if _, ok := val.(map[string]interface{}); !ok {
			leaf = true
		}
	}
	if leaf {
		comparable := oldval != nil && reflect.TypeOf(oldval).Comparable() && val != nil && reflect.TypeOf(val).Comparable()
		if comparable {
			if oldval != val {
				// defer s.rw.Unlock()
				// s.rw.Lock()
				s.entries[path] = val
				m[key] = val
				if s.onMergingSet != nil {
					s.onMergingSet(path, val, oldval)
				}
				// logrus.Debugf("%%-> s.entries[%q] = m[%q] = %v", path, key, val)
				return
			}
		}
	}
	s.entries[path] = val
	m[key] = val
}

// et will eat the left part string from `keys[ix:]`
func et(keys []string, ix int, val interface{}) interface{} {
	if ix <= len(keys)-1 {
		p := make(map[string]interface{})
		p[keys[ix]] = et(keys, ix+1, val)
		return p
	}
	return val
}

// Reset the exists `Options`, so that you could follow a `LoadConfigFile()` with it.
func (s *Options) Reset() {
	defer s.rw.Unlock()
	s.rw.Lock()

	s.entries = nil
	time.Sleep(100 * time.Millisecond)
	s.entries = make(map[string]interface{})
}

func mx(pre, k string) string {
	if len(pre) == 0 {
		return k
	}
	return pre + "." + k
}

func mxIx(pre string, k interface{}) string {
	if len(pre) == 0 {
		return fmt.Sprintf("%v", k)
	}
	return fmt.Sprintf("%v.%v", pre, k)
}

func (s *Options) loopMapMap(kdot string, m map[string]map[string]interface{}) (err error) {
	for k, v := range m {
		if err = s.loopMap(mx(kdot, k), v); err != nil {
			return
		}
	}
	return
}

func (s *Options) loopMap(kdot string, m map[string]interface{}) (err error) {
	for k, v := range m {
		if vm, ok := v.(map[interface{}]interface{}); ok {
			if err = s.loopIxMap(mx(kdot, k), vm); err != nil {
				return
			}
		} else if vm, ok := v.(map[string]interface{}); ok {
			if err = s.loopMap(mx(kdot, k), vm); err != nil {
				return
			}
		} else {
			// s.SetNx(mx(kdot, k), v)
			key := mxIx(kdot, k)
			oldval, modi := s.setNx(key, v)
			if s.onMergingSet != nil && modi {
				s.onMergingSet(key, v, oldval)
			}
		}
	}
	return
}

func (s *Options) loopIxMap(kdot string, m map[interface{}]interface{}) (err error) {
	for k, v := range m {
		if vm, ok := v.(map[interface{}]interface{}); ok {
			if err = s.loopIxMap(mxIx(kdot, k), vm); err != nil {
				return
			}
			// } else if vm, ok := v.(map[string]interface{}); ok {
			// 	if err = s.loopMap(mxIx(kdot, k), vm); err != nil {
			// 		return
			// 	}
		} else {
			// s.SetNx(mx(kdot, k), v)
			key := mxIx(kdot, k)
			oldval, modi := s.setNx(key, v)
			if s.onMergingSet != nil && modi {
				s.onMergingSet(key, v, oldval)
			}
		}
	}
	return
}

// DumpAsString for debugging.
func (s *Options) DumpAsString() (str string) {
	k3 := make([]string, 0)
	for k := range s.entries {
		k3 = append(k3, k)
	}
	sort.Strings(k3)

	for _, k := range k3 {
		str = str + fmt.Sprintf("%-48v => %v\n", k, s.entries[k])
	}
	str += "---------------------------------\n"

	b, err := yaml.Marshal(s.hierarchy)
	if err == nil {
		if s.GetBoolEx("raw") {
			str += string(b)
		} else {
			ss := string(b)
			ss = os.ExpandEnv(ss)
			str += ss
		}
	}
	return
}

// GetHierarchyList returns the hierarchy data for dumping
func (s *Options) GetHierarchyList() map[string]interface{} {
	defer s.rw.RUnlock()
	s.rw.RLock()
	return s.hierarchy
}
