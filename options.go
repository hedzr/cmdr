/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

func NewOptions() *Options {
	return &Options{
		make(map[string]interface{}),
		make(map[string]interface{}),
	}
}

func NewOptionsWith(entries map[string]interface{}) *Options {
	return &Options{
		entries:   entries,
		hierarchy: make(map[string]interface{}),
	}
}

func Get(key string) interface{} {
	return RxxtOptions.Get(key)
}

// Get an `Option` by key string, eg:
// ```golang
// cmdr.Get("app.logger.level") => 'DEBUG',...
// ```
//
func (s *Options) Get(key string) interface{} {
	return s.entries[key]
}

func GetBool(key string) bool {
	return RxxtOptions.GetBool(key)
}

func GetInt(key string) int64 {
	return RxxtOptions.GetInt(key)
}

func GetString(key string) string {
	return RxxtOptions.GetString(key)
}

func GetStringSlice(key string) []string {
	return RxxtOptions.GetStringSlice(key)
}

func (s *Options) GetBool(key string) (ret bool) {
	switch strings.ToLower(s.GetString(key)) {
	case "1", "y", "t", "yes", "true", "ok", "on":
		ret = true
	}
	return
}

func (s *Options) GetInt(key string) (ir int64) {
	if ir64, err := strconv.ParseInt(s.GetString(key), 10, 64); err == nil {
		ir = ir64
	}
	return
}

func (s *Options) GetStringSlice(key string) (ir []string) {
	ss := s.GetString(key)
	ir = strings.Split(ss, ",")
	return
}

func (s *Options) GetString(key string) (ret string) {
	envkey := strings.Join(append(EnvPrefix, strings.ToUpper(strings.ReplaceAll(key, ".", "_"))), "_")
	if s, ok := os.LookupEnv(envkey); ok {
		ret = s
	}

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

func (s *Options) Set(key string, val interface{}) {
	k := wrapRxxtPrefix(key)
	s.entries[k] = val
	a := strings.Split(k, ".")
	mergeMap(s.hierarchy, a[0], et(a, 1, val))
}

// Set() but without prefix auto-wrapped.
// `RxxtPrefix` is a string slice to define the prefix string array, default is ["app"].
// So, cmdr.Set("debug", true) will put an real entry with (`app.debug`, true).
func (s *Options) SetNx(key string, val interface{}) {
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

// Reset the exists `Options`, so that you could follow a `LoadConfigFile()` with it.
func (s *Options) Reset() {
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

func (s *Options) loopMap(kdot string, m map[string]interface{}) (err error) {
	for k, v := range m {
		if vm, ok := v.(map[interface{}]interface{}); ok {
			err = s.loopIxMap(mx(kdot, k), vm)
			if err != nil {
				return
			}
		} else if vm, ok := v.(map[string]interface{}); ok {
			err = s.loopMap(mx(kdot, k), vm)
			if err != nil {
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
			err = s.loopIxMap(mxIx(kdot, k), vm)
			if err != nil {
				return
			}
		} else if vm, ok := v.(map[string]interface{}); ok {
			err = s.loopMap(mxIx(kdot, k), vm)
			if err != nil {
				return
			}
		} else {
			s.SetNx(mxIx(kdot, k), v)
		}
	}
	return
}

func DumpAsString() (str string) {
	return RxxtOptions.DumpAsString()
}

func (s *Options) DumpAsString() (str string) {
	k3 := make([]string, 0)
	for k, _ := range s.entries {
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
