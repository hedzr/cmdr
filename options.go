/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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
	}
}

// NewOptionsWith returns an `Options` structure pointer
func NewOptionsWith(entries map[string]interface{}) *Options {
	return &Options{
		entries:   entries,
		hierarchy: make(map[string]interface{}),
		rw:        new(sync.RWMutex),
	}
}

// Get returns the generic value of an `Option` key. Such as:
// ```golang
// cmdr.Get("app.logger.level") => 'DEBUG',...
// ```
//
func Get(key string) interface{} {
	return rxxtOptions.Get(key)
}

// GetBool returns the bool value of an `Option` key.
func GetBool(key string) bool {
	return rxxtOptions.GetBool(key)
}

// GetBoolP returns the bool value of an `Option` key.
func GetBoolP(prefix, key string) bool {
	return rxxtOptions.GetBool(fmt.Sprintf("%s.%s", prefix, key))
}

// GetInt returns the int value of an `Option` key.
func GetInt(key string) int {
	return int(rxxtOptions.GetInt(key))
}

// GetIntP returns the int value of an `Option` key.
func GetIntP(prefix, key string) int {
	return int(rxxtOptions.GetInt(fmt.Sprintf("%s.%s", prefix, key)))
}

// GetInt64 returns the int64 value of an `Option` key.
func GetInt64(key string) int64 {
	return rxxtOptions.GetInt(key)
}

// GetInt64P returns the int64 value of an `Option` key.
func GetInt64P(prefix, key string) int64 {
	return rxxtOptions.GetInt(fmt.Sprintf("%s.%s", prefix, key))
}

// GetUint returns the uint value of an `Option` key.
func GetUint(key string) uint {
	return uint(rxxtOptions.GetUint(key))
}

// GetUintP returns the uint value of an `Option` key.
func GetUintP(prefix, key string) uint {
	return uint(rxxtOptions.GetUint(fmt.Sprintf("%s.%s", prefix, key)))
}

// GetUint64 returns the uint64 value of an `Option` key.
func GetUint64(key string) uint64 {
	return rxxtOptions.GetUint(key)
}

// GetUint64P returns the uint64 value of an `Option` key.
func GetUint64P(prefix, key string) uint64 {
	return rxxtOptions.GetUint(fmt.Sprintf("%s.%s", prefix, key))
}

// GetString returns the string value of an `Option` key.
func GetString(key string) string {
	return rxxtOptions.GetString(key)
}

// GetStringP returns the string value of an `Option` key.
func GetStringP(prefix, key string) string {
	return rxxtOptions.GetString(fmt.Sprintf("%s.%s", prefix, key))
}

// GetStringSlice returns the string slice value of an `Option` key.
func GetStringSlice(key string) []string {
	return rxxtOptions.GetStringSlice(key)
}

// GetStringSliceP returns the string slice value of an `Option` key.
func GetStringSliceP(prefix, key string) []string {
	return rxxtOptions.GetStringSlice(fmt.Sprintf("%s.%s", prefix, key))
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

// GetBool returns the bool value of an `Option` key.
func (s *Options) GetBool(key string) (ret bool) {
	switch strings.ToLower(s.GetString(key)) {
	case "1", "y", "t", "yes", "true", "ok", "on":
		ret = true
	}
	return
}

// GetInt returns the int64 value of an `Option` key.
func (s *Options) GetInt(key string) (ir int64) {
	if ir64, err := strconv.ParseInt(s.GetString(key), 10, 64); err == nil {
		ir = ir64
	}
	return
}

// GetUint returns the uint64 value of an `Option` key.
func (s *Options) GetUint(key string) (ir uint64) {
	if ir64, err := strconv.ParseUint(s.GetString(key), 10, 64); err == nil {
		ir = ir64
	}
	return
}

// GetStringSlice returns the string slice value of an `Option` key.
func (s *Options) GetStringSlice(key string) (ir []string) {
	// envkey := s.envKey(key)
	// if s, ok := os.LookupEnv(envkey); ok {
	// 	ir = strings.Split(s, ",")
	// }

	defer s.rw.RUnlock()
	s.rw.RLock()

	if v, ok := s.entries[key]; ok {
		switch reflect.ValueOf(v).Kind() {
		case reflect.String:
			ir = strings.Split(v.(string), ",")
		case reflect.Slice:
			if r, ok := v.([]string); ok {
				ir = r
			} else if ri, ok := v.([]int); ok {
				for _, rii := range ri {
					ir = append(ir, strconv.Itoa(rii))
				}
			}
		default:
			ir = strings.Split(fmt.Sprintf("%v", v), ",")
		}
	}
	return
}

func stringSliceToIntSlice(in []string) (out []int) {
	for _, ii := range in {
		if i, err := strconv.Atoi(ii); err == nil {
			out = append(out, i)
		}
	}
	return
}

// GetIntSlice returns the int slice value of an `Option` key.
func GetIntSlice(key string) []int {
	return rxxtOptions.GetIntSlice(key)
}

// GetIntSliceP returns the int slice value of an `Option` key.
func GetIntSliceP(prefix, key string) []int {
	return rxxtOptions.GetIntSlice(fmt.Sprintf("%s.%s", prefix, key))
}

// GetIntSlice returns the string slice value of an `Option` key.
func (s *Options) GetIntSlice(key string) (ir []int) {
	// envkey := s.envKey(key)
	// if s, ok := os.LookupEnv(envkey); ok {
	// 	ir = stringSliceToIntSlice(strings.Split(s, ","))
	// }

	defer s.rw.RUnlock()
	s.rw.RLock()

	if v, ok := s.entries[key]; ok {
		switch reflect.ValueOf(v).Kind() {
		case reflect.String:
			ir = stringSliceToIntSlice(strings.Split(v.(string), ","))
		case reflect.Slice:
			if r, ok := v.([]string); ok {
				ir = stringSliceToIntSlice(r)
			} else if ri, ok := v.([]int); ok {
				ir = ri
			}
		default:
			ir = stringSliceToIntSlice(strings.Split(fmt.Sprintf("%v", v), ","))
		}
	}
	return
}

// GetDuration returns the int slice value of an `Option` key.
func GetDuration(key string) time.Duration {
	return rxxtOptions.GetDuration(key)
}

// GetDurationP returns the int slice value of an `Option` key.
func GetDurationP(prefix, key string) time.Duration {
	return rxxtOptions.GetDuration(fmt.Sprintf("%s.%s", prefix, key))
}

// GetDuration returns the time duration value of an `Option` key.
func (s *Options) GetDuration(key string) (ir time.Duration) {
	str := s.GetString(key)
	ir, _ = time.ParseDuration(str)
	return
}

// GetString returns the string value of an `Option` key.
func (s *Options) GetString(key string) (ret string) {
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
		default:
			ret = fmt.Sprintf("%v", v)
		}
	}
	return
}

func (s *Options) buildAutomaticEnv(rootCmd *RootCommand) (err error) {
	// p := strings.Join(EnvPrefix,"_")
	p := strings.Join(RxxtPrefix, ".")
	for key := range s.entries {
		ek := s.envKey(key)
		if v, ok := os.LookupEnv(ek); ok {
			if strings.HasPrefix(key, p) {
				s.Set(key[len(p)+1:], v)
			} else {
				s.Set(key, v)
			}
		}
	}
	return
}

func (s *Options) envKey(key string) (envkey string) {
	key = strings.ReplaceAll(key, ".", "_")
	key = strings.ReplaceAll(key, "-", "_")
	envkey = strings.Join(append(EnvPrefix, strings.ToUpper(key)), "_")
	return
}

func wrapRxxtPrefix(key string) string {
	p := strings.Join(RxxtPrefix, ".")
	if len(p) == 0 {
		return key
	}
	if len(key) == 0 {
		return p
	}
	return p + "." + key
}

// Set set the value of an `Option` key (with prefix auto-wrap).
// ```golang
// cmdr.Set("logger.level", "DEBUG")
// cmdr.Set("ms.tags.port", 8500)
// ...
// ```
//
func Set(key string, val interface{}) {
	rxxtOptions.Set(key, val)
}

// SetNx but without prefix auto-wrapped.
// `rxxtPrefix` is a string slice to define the prefix string array, default is ["app"].
// So, cmdr.Set("debug", true) will put an real entry with (`app.debug`, true).
func SetNx(key string, val interface{}) {
	rxxtOptions.SetNx(key, val)
}

// Set set the value of an `Option` key.
func (s *Options) Set(key string, val interface{}) {
	defer s.rw.Unlock()
	s.rw.Lock()

	k := wrapRxxtPrefix(key)
	s.entries[k] = val
	a := strings.Split(k, ".")
	mergeMap(s.hierarchy, a[0], et(a, 1, val))
}

// SetNx but without prefix auto-wrapped.
// `rxxtPrefix` is a string slice to define the prefix string array, default is ["app"].
// So, cmdr.Set("debug", true) will put an real entry with (`app.debug`, true).
func (s *Options) SetNx(key string, val interface{}) {
	defer s.rw.Unlock()
	s.rw.Lock()

	s.entries[key] = val
	a := strings.Split(key, ".")
	mergeMap(s.hierarchy, a[0], et(a, 1, val))
}

func mergeMap(m map[string]interface{}, key string, val interface{}) map[string]interface{} {
	if z, ok := m[key]; ok {
		if zm, ok := z.(map[string]interface{}); ok {
			if vm, ok := val.(map[string]interface{}); ok {
				for k, v := range vm {
					zm = mergeMap(zm, k, v)
				}
				m[key] = zm
			} else {
				m[key] = val
			}
		} else {
			m[key] = val
		}
	} else {
		m[key] = val
	}
	return m
}

func et(keys []string, ix int, val interface{}) interface{} {
	if ix <= len(keys)-1 {
		p := make(map[string]interface{})
		p[keys[ix]] = et(keys, ix+1, val)
		return p
	}
	return val
}

// ResetOptions to reset the exists `Options`, so that you could follow a `LoadConfigFile()` with it.
func ResetOptions() {
	rxxtOptions.Reset()
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
			s.SetNx(mx(kdot, k), v)
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
			s.SetNx(mxIx(kdot, k), v)
		}
	}
	return
}

// DumpAsString for debugging.
func DumpAsString() (str string) {
	return rxxtOptions.DumpAsString()
}

// SaveAsYaml to Save all config entries as a yaml file
func SaveAsYaml(filename string) (err error) {
	obj := rxxtOptions.GetHierarchyList()

	b, err := yaml.Marshal(obj)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(filename, b, 0644)
	return
}

// SaveAsJSON to Save all config entries as a json file
func SaveAsJSON(filename string) (err error) {
	obj := rxxtOptions.GetHierarchyList()

	b, err := json.Marshal(obj)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(filename, b, 0644)
	return
}

// SaveAsToml to Save all config entries as a toml file
func SaveAsToml(filename string) (err error) {
	obj := rxxtOptions.GetHierarchyList()
	err = SaveObjAsToml(obj, filename)
	return
}

// SaveObjAsToml to Save an object as a toml file
func SaveObjAsToml(obj interface{}, filename string) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return
	}

	e := toml.NewEncoder(bufio.NewWriter(f))
	if err = e.Encode(obj); err != nil {
		return
	}

	// err = ioutil.WriteFile(filename, b, 0644)
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
		str += string(b)
	}
	return
}

// GetHierarchyList returns the hierarchy data for dumping
func (s *Options) GetHierarchyList() map[string]interface{} {
	defer s.rw.RUnlock()
	s.rw.RLock()
	return s.hierarchy
}
