// Copyright © 2019 Hedzr Yeh.

package cmdr

import (
	"encoding"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/hedzr/evendeep"

	"github.com/hedzr/cmdr/tool"
	"github.com/hedzr/log"
	"github.com/hedzr/log/detects"
	"github.com/hedzr/log/dir"
)

// newOptions returns an `Options` structure pointer
func newOptions() *Options {
	return &Options{
		entries:   make(map[string]interface{}),
		hierarchy: make(map[string]interface{}),
		rw:        new(sync.RWMutex),

		onConfigReloadedFunctions: make(map[ConfigReloaded]bool),
		rwlCfgReload:              new(sync.RWMutex),
	}
}

// newOptionsWith returns an `Options` structure pointer
func newOptionsWith(entries map[string]interface{}) *Options {
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
	_, ok = s.hasKey(key)
	return
}

func (s *Options) hasKey(key string) (v interface{}, ok bool) {
	defer s.rw.RUnlock()
	s.rw.RLock()

	v, ok = s.entries[key]
	return
}

// DeleteKey deletes a key from cmdr options store
func DeleteKey(key string) {
	currentOptions().Delete(key)
}

// Delete deletes a key from cmdr options store
func (s *Options) Delete(key string) {
	defer s.rw.RUnlock()
	s.rw.RLock()

	val := s.entries[key]
	a := strings.Split(key, ".")
	s.deleteWithKey(s.hierarchy, a[0], "", et(a, 1, val))
	// return
}

func (s *Options) deleteWithKey(m map[string]interface{}, key, path string, val interface{}) (ret map[string]interface{}) {
	if len(path) > 0 {
		path = fmt.Sprintf("%v.%v", path, key)
	} else {
		path = key
	}

	if z, ok := m[key]; ok {
		if zm, ok := z.(map[string]interface{}); ok {
			switch vm := val.(type) {
			case map[string]interface{}:
				for k, v := range vm {
					zm = s.deleteWithKey(zm, k, path, v)
				}
				// delete(m, key)
				// delete(s.entries, path)
				return
			case map[interface{}]interface{}:
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
func (s *Options) Get(key string) (ret interface{}) {
	defer s.rw.RUnlock()
	s.rw.RLock()
	ret = s.entries[key]
	return
}

// GetMap an `Option` by key string, it returns a hierarchy map or nil
func (s *Options) GetMap(key string) map[string]interface{} {
	defer s.rw.RUnlock()
	s.rw.RLock()

	return s.getMapNoLock(key)
}

func (s *Options) getMapNoLock(key string) (m map[string]interface{}) {
	if v, ok := s.entries[key]; ok {
		if vv, ok := v.(map[string]interface{}); ok {
			return vv
		}
		tv, ts := reflect.TypeOf(v), "<nil>"
		if tv != nil {
			ts = tv.String()
		}
		fwrn("need attention: getMapNoLock(%q) not found. Or got: '%v'", key, ts)
	}

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
	ret = toBool(s.GetString(key, ""), defaultVal...)
	return
}

// ToBool translate a value to boolean
func ToBool(val interface{}, defaultVal ...bool) (ret bool) {
	if v, ok := val.(bool); ok {
		return v
	}
	if v, ok := val.(int); ok {
		return v != 0
	}
	if v, ok := val.(string); ok {
		return toBool(v, defaultVal...)
	}
	for _, vv := range defaultVal {
		ret = vv
	}
	return
}

func toBool(val string, defaultVal ...bool) (ret bool) {
	// ret = ToBool(val, defaultVal...)
	switch strings.ToLower(val) {
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
	if ir64, err := strconv.ParseInt(s.GetString(key, ""), 0, 64); err == nil {
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
	if ir64, err := strconv.ParseInt(s.GetString(key, ""), 0, 64); err == nil {
		ir = ir64
	} else {
		for _, vv := range defaultVal {
			ir = vv
		}
	}
	return
}

// GetKibibytesEx returns the uint64 value of an `Option` key based kibibyte format.
//
// kibibyte format is for human readable. In this format, number presentations
// are: 2k, 8m, 3g, 5t, 6p, 7e. optional 'b' can be appended, such as: 2kb, 5tb, 7EB.
// All of them is case-insensitive.
//
// kibibyte is based 1024. That means:
// 1 KiB = 1k = 1024 bytes
//
// See also: https://en.wikipedia.org/wiki/Kibibyte
// Its related word is kilobyte, refer to: https://en.wikipedia.org/wiki/Kilobyte
//
// The pure number part can be golang presentation, such as 0x99, 0001b, 0700.
func (s *Options) GetKibibytesEx(key string, defaultVal ...uint64) (ir64 uint64) {
	sz := s.GetString(key, "")
	if sz == "" {
		for _, v := range defaultVal {
			ir64 = v
		}
		return
	}
	return s.FromKibiBytes(sz)
}

// FromKibiBytes convert string to the uint64 value based kibibyte format.
//
// kibibyte format is for human readable. In this format, number presentations
// are: 2k, 8m, 3g, 5t, 6p, 7e. optional 'b' can be appended, such as: 2kb, 5tb, 7EB.
// All of them is case-insensitive.
//
// kibibyte is based 1024. That means:
// 1 KiB = 1k = 1024 bytes
//
// See also: https://en.wikipedia.org/wiki/Kibibyte
// Its related word is kilobyte, refer to: https://en.wikipedia.org/wiki/Kilobyte
//
// The pure number part can be golang presentation, such as 0x99, 0001b, 0700.
func (s *Options) FromKibiBytes(sz string) (ir64 uint64) {
	// var suffixes = []string {"B","KB","MB","GB","TB","PB","EB","ZB","YB"}
	const suffix = "kmgtpezyKMGTPEZY"
	sz = strings.TrimSpace(sz)
	sz = strings.TrimRight(sz, "B")
	sz = strings.TrimRight(sz, "b")
	szr := strings.TrimSpace(strings.TrimRightFunc(sz, func(r rune) bool {
		return strings.ContainsRune(suffix, r)
	}))

	var if64 float64
	var err error
	if strings.ContainsRune(szr, '.') {
		if if64, err = strconv.ParseFloat(szr, 64); err == nil {
			r := []rune(sz)[len(sz)-1]
			ir64 = uint64(if64 * float64(s.fromKibiBytes(r)))
		}
	} else {
		if ir64, err = strconv.ParseUint(szr, 0, 64); err == nil {
			r := []rune(sz)[len(sz)-1]
			ir64 *= s.fromKibiBytes(r)
		}
	}
	return
}

func (s *Options) fromKibiBytes(r rune) (times uint64) {
	switch r {
	case 'k', 'K':
		return 1024
	case 'm', 'M':
		return 1024 * 1024
	case 'g', 'G':
		return 1024 * 1024 * 1024
	case 't', 'T':
		return 1024 * 1024 * 1024 * 1024
	case 'p', 'P':
		return 1024 * 1024 * 1024 * 1024 * 1024
	case 'e', 'E':
		return 1024 * 1024 * 1024 * 1024 * 1024 * 1024
	// case 'z', 'Z':
	// 	ir64 *= 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024
	// case 'y', 'Y':
	// 	ir64 *= 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024
	default:
		return 1
	}
}

// GetKilobytesEx returns the uint64 value of an `Option` key with kilobyte format.
//
// kilobyte format is for human readable. In this format, number presentations
// are: 2k, 8m, 3g, 5t, 6p, 7e. optional 'b' can be appended, such as: 2kb, 5tb, 7EB.
// All of them is case-insensitive.
//
// kilobyte is based 1000. That means:
// 1 KB = 1k = 1000 bytes
//
// See also: https://en.wikipedia.org/wiki/Kilobyte
// Its related word is kibibyte, refer to: https://en.wikipedia.org/wiki/Kibibyte
//
// The pure number part can be golang presentation, such as 0x99, 0001b, 0700.
func (s *Options) GetKilobytesEx(key string, defaultVal ...uint64) (ir64 uint64) {
	sz := s.GetString(key, "")
	if sz == "" {
		for _, v := range defaultVal {
			ir64 = v
		}
		return
	}
	return s.FromKilobytes(sz)
}

// FromKilobytes convert string to the uint64 value based kilobyte format.
//
// kilobyte format is for human readable. In this format, number presentations
// are: 2k, 8m, 3g, 5t, 6p, 7e. optional 'b' can be appended, such as: 2kb, 5tb, 7EB.
// All of them is case-insensitive.
//
// kilobyte is based 1000. That means:
// 1 KB = 1k = 1000 bytes
//
// See also: https://en.wikipedia.org/wiki/Kilobyte
// Its related word is kibibyte, refer to: https://en.wikipedia.org/wiki/Kibibyte
//
// The pure number part can be golang presentation, such as 0x99, 0001b, 0700.
func (s *Options) FromKilobytes(sz string) (ir64 uint64) {
	// var suffixes = []string {"B","KB","MB","GB","TB","PB","EB","ZB","YB"}
	const suffix = "kmgtpezyKMGTPEZY"
	sz = strings.TrimSpace(sz)
	sz = strings.TrimRight(sz, "B")
	sz = strings.TrimRight(sz, "b")
	szr := strings.TrimSpace(strings.TrimRightFunc(sz, func(r rune) bool {
		return strings.ContainsRune(suffix, r)
	}))

	var if64 float64
	var err error
	if strings.ContainsRune(szr, '.') {
		if if64, err = strconv.ParseFloat(szr, 64); err == nil {
			r := []rune(sz)[len(sz)-1]
			ir64 = uint64(if64 * float64(s.fromKilobytes(r)))
		}
	} else {
		if ir64, err = strconv.ParseUint(szr, 0, 64); err == nil {
			r := []rune(sz)[len(sz)-1]
			ir64 *= s.fromKilobytes(r)
		}
	}
	return
}

func (s *Options) fromKilobytes(r rune) (times uint64) {
	switch r {
	case 'k', 'K':
		return 1000
	case 'm', 'M':
		return 1000 * 1000
	case 'g', 'G':
		return 1000 * 1000 * 1000
	case 't', 'T':
		return 1000 * 1000 * 1000 * 1000
	case 'p', 'P':
		return 1000 * 1000 * 1000 * 1000 * 1000
	case 'e', 'E':
		return 1000 * 1000 * 1000 * 1000 * 1000 * 1000
	// case 'z', 'Z':
	// 	ir64 *= 1000 * 1000 * 1000 * 1000 * 1000 * 1000 * 1000
	// case 'y', 'Y':
	// 	ir64 *= 1000 * 1000 * 1000 * 1000 * 1000 * 1000 * 1000 * 1000
	default:
		return 1
	}
}

// GetUintEx returns the uint64 value of an `Option` key.
func (s *Options) GetUintEx(key string, defaultVal ...uint) (ir uint) {
	if ir64, err := strconv.ParseUint(s.GetString(key, ""), 0, 64); err == nil {
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
	if ir64, err := strconv.ParseUint(s.GetString(key, ""), 0, 64); err == nil {
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
	if ir64, err := strconv.ParseFloat(s.GetString(key, ""), 32); err == nil {
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
	if ir64, err := strconv.ParseFloat(s.GetString(key, ""), 64); err == nil {
		ir = ir64
	} else {
		for _, vv := range defaultVal {
			ir = vv
		}
	}
	return
}

// GetComplex64 returns the complex64 value of an `Option` key.
func (s *Options) GetComplex64(key string, defaultVal ...complex64) (ir complex64) {
	if ir128, err := tool.ParseComplexX(s.GetString(key, "")); err == nil {
		ir = complex64(ir128)
	} else {
		for _, vv := range defaultVal {
			ir = vv
		}
	}
	return
}

// GetComplex128 returns the complex128 value of an `Option` key.
func (s *Options) GetComplex128(key string, defaultVal ...complex128) (ir complex128) {
	if ir128, err := tool.ParseComplexX(s.GetString(key, "")); err == nil {
		ir = ir128
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
			ir = strings.Split(os.ExpandEnv(v.(string)), ",")
		case reflect.Slice:
			switch rr := v.(type) {
			case []string:
				// ir = r
				for _, xx := range rr {
					ir = append(ir, os.ExpandEnv(xx))
				}
			case []int:
				for _, rii := range rr {
					ir = append(ir, os.ExpandEnv(strconv.Itoa(rii)))
				}
			case []byte:
				ir = strings.Split(os.ExpandEnv(string(rr)), ",")
			default:
				for i := 0; i < vvv.Len(); i++ {
					ir = append(ir, os.ExpandEnv(fmt.Sprintf("%v", vvv.Index(i).Interface())))
				}
			}
		default:
			ir = strings.Split(os.ExpandEnv(fmt.Sprintf("%v", v)), ",")
		}
	} else if len(defaultVal) > 0 {
		for _, xx := range defaultVal {
			ir = append(ir, os.ExpandEnv(xx))
		}
	}
	// ret = os.ExpandEnv(ret)
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
			switch rr := v.(type) {
			case []string:
				ir = stringSliceToIntSlice(rr)
			case []int:
				ir = rr
			case []int8:
				ir = int8SliceToIntSlice(rr)
			case []int16:
				ir = int16SliceToIntSlice(rr)
			case []int32:
				ir = int32SliceToIntSlice(rr)
			case []int64:
				ir = int64SliceToIntSlice(rr)
			case []uint:
				ir = uintSliceToIntSlice(rr)
			// case []uint8:
			// 	ir = uint8SliceToIntSlice(rr)
			case []uint16:
				ir = uint16SliceToIntSlice(rr)
			case []uint32:
				ir = uint32SliceToIntSlice(rr)
			case []uint64:
				ir = uint64SliceToIntSlice(rr)
			case []byte:
				xx := strings.Split(string(rr), ",")
				ir = stringSliceToIntSlice(xx)
			default:
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
			switch rr := v.(type) {
			case []string:
				ir = stringSliceToInt64Slice(rr)
			case []int:
				ir = intSliceToInt64Slice(rr)
			case []int8:
				ir = int8SliceToInt64Slice(rr)
			case []int16:
				ir = int16SliceToInt64Slice(rr)
			case []int32:
				ir = int32SliceToInt64Slice(rr)
			case []int64:
				ir = rr
			case []uint:
				ir = uintSliceToInt64Slice(rr)
			// case []uint8:
			// 	ir = uint8SliceToInt64Slice(rr)
			case []uint16:
				ir = uint16SliceToInt64Slice(rr)
			case []uint32:
				ir = uint32SliceToInt64Slice(rr)
			case []uint64:
				ir = uint64SliceToInt64Slice(rr)
			case []byte:
				xx := strings.Split(string(rr), ",")
				ir = stringSliceToInt64Slice(xx)
			default:
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
			switch rr := v.(type) {
			case []string:
				ir = stringSliceToUint64Slice(rr)
			case []int:
				ir = intSliceToUint64Slice(rr)
			case []int8:
				ir = int8SliceToUint64Slice(rr)
			case []int16:
				ir = int16SliceToUint64Slice(rr)
			case []int32:
				ir = int32SliceToUint64Slice(rr)
			case []int64:
				ir = int64SliceToUint64Slice(rr)
			case []uint:
				ir = uintSliceToUint64Slice(rr)
			case []uint16:
				ir = uint16SliceToUint64Slice(rr)
			case []uint32:
				ir = uint32SliceToUint64Slice(rr)
			case []uint64:
				ir = rr
			case []byte:
				xx := strings.Split(string(rr), ",")
				ir = stringSliceToUint64Slice(xx)
			default:
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

// Wrap wraps a key with [RxxtPrefix], for [GetXxx(key)] and [GetXxxP(prefix,key)]
func (s *Options) Wrap(key string) string {
	return wrapWithRxxtPrefix(key)
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
			if ret, ok = v.(string); ok && ret == "" {
				for _, v := range defaultVal {
					ret = v
				}
			}
		default:
			if v != nil {
				ret = fmt.Sprint(v)
			} else {
				for _, vv := range defaultVal {
					ret = vv
				}
			}
		}
	} else {
		for _, vv := range defaultVal {
			ret = vv
		}
	}
	return
}

func (s *Options) oneenv(rootCmd *RootCommand, prefix, key string) {
	ek := s.envKey(key)
	if v, ok := os.LookupEnv(ek); ok {
		if strings.HasPrefix(key, prefix) {
			s.Set(key[len(prefix)+1:], v)
		} else {
			s.Set(key, v)
		}
	}
	// Logger.Printf("buildAutomaticEnv: %v", key)
	if flg := s.lookupFlag(key, rootCmd); flg != nil {
		// flog("    [cmdr] lookupFlag for %q: %v", key, flg.GetTitleName())
		//
		// if key == "app.mx-test.test" {
		// 	Logger.Debugf("                 : flag=%+v", flg)
		// }
		for _, ek = range flg.EnvVars {
			if v, ok := os.LookupEnv(ek); ok { //nolint:gocritic //can't reduce
				// flog("    [cmdr][buildAutomaticEnv] envvar %q found (flg=%v): %v", ek, flg.GetTitleName(), v)
				k := key
				if !strings.HasPrefix(k, prefix) {
					k = wrapWithRxxtPrefix(key)
				}
				// Logger.Printf("setnx: %v <-- %v", key, v)
				s.setNxNoLock(k, v)
				if flg.ToggleGroup != "" {
					s.tryResetOthersInTG(flg, k)
				}
				// Logger.Printf("setnx: %v", s.GetString(key))

				if flg.onSet != nil {
					flg.onSet(key, v)
				}
			}
		}
	}
}

func (s *Options) buildAutomaticEnv(rootCmd *RootCommand) (err error) {
	// Logger.SetLevel(logrus.DebugLevel)

	s.rwCB.RLock()
	defer s.rwCB.RUnlock()

	s.rw.RLock()
	defer s.rw.RUnlock()

	// prefix := strings.Join(EnvPrefix,"_")
	prefix := internalGetWorker().getPrefix() // strings.Join(RxxtPrefix, ".")
	for key := range s.entries {
		s.oneenv(rootCmd, prefix, key)
	}

	// // fmt.Printf("EXE = %v, PWD = %v, CURRDIR = %v\n", GetExecutableDir(), os.Getenv("PWD"), GetCurrentDir())
	// // _ = os.Setenv("THIS", GetExecutableDir())
	// for k, v := range uniqueWorker.envVarToValueMap {
	// 	_ = os.Setenv(k, v())
	// }
	internalGetWorker().setupFromEnvvarMap()

	for _, h := range internalGetWorker().afterAutomaticEnv {
		h(rootCmd, s)
	}
	return
}

func (s *Options) tryResetOthersInTG(flg *Flag, fullKey string) {
	tgs := make(map[string][]*Flag)
	for _, c := range flg.owner.Flags {
		if c.ToggleGroup != "" {
			tgs[c.ToggleGroup] = append(tgs[c.ToggleGroup], c)
		}
	}

	for _, c := range tgs[flg.ToggleGroup] {
		if c == flg {
			continue
		}

		k1 := strings.Split(fullKey, ".")
		k1[len(k1)-1] = c.Full
		k2 := strings.Join(k1, ".")
		s.setNxNoLock(k2, false)

		if f := s.lookupFlag(k2, flg.owner.GetRoot()); f != nil {
			f.DefaultValueType = "bool"
			f.DefaultValue = false
		}
	}

	k1 := flg.owner.GetDottedNamePath() + "." + flg.ToggleGroup
	k2 := wrapWithRxxtPrefix(k1)
	s.setNxNoLock(k2, flg.Full)
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

func (s *Options) envKey(key string) (envKey string) {
	key = replaceAll(key, ".", "_")
	key = replaceAll(key, "-", "_")
	envKey = strings.Join(append(internalGetWorker().envPrefixes, strings.ToUpper(key)), "_")
	return
}

// Set sets the value of an `Option` key. The key MUST not have an `app` prefix.
//
// For example:
// ```golang
// cmdr.Set("debug", true)
// cmdr.GetBool("app.debug") => true
// ```
//
// Set runs in slice appending mode: for a slice, Set(key, newSlice) will append
// newSlice into original slice.
func (s *Options) Set(key string, val interface{}) {
	k := wrapWithRxxtPrefix(key)
	s.setNx(k, val)
}

// SetOverwrite sets the value of an `Option` key. The key MUST not have an `app` prefix.
// It replaces the old value on a slice, instead of the default append mode.
//
// SetOverwrite(key, newSlice) will replace the original value with newSlice.
func (s *Options) SetOverwrite(key string, val interface{}) {
	k := wrapWithRxxtPrefix(key)
	s.setNxOverwrite(k, val)
}

// SetNx likes Set but without prefix auto-wrapped.
// `rxxtPrefix` is a string slice to define the prefix string array, default is ["app"].
// So, cmdr.SetNx("debug", true) will put a real entry with (`debug`, true).
func (s *Options) SetNx(key string, val interface{}) {
	s.setNx(key, val)
}

// SetNxOverwrite likes SetOverwrite but without prefix auto-wrapped.
// It replaces the old value on a slice, instead of the default append mode.
func (s *Options) SetNxOverwrite(key string, val interface{}) {
	s.setNxOverwrite(key, val)
}

// SetRaw but without prefix auto-wrapped.
// So, cmdr.SetRaw("debug", true) will put a real entry with (`debug`, true).
func (s *Options) SetRaw(key string, val interface{}) {
	s.setNx(key, val)
}

// SetRawOverwrite likes SetOverwrite but without prefix auto-wrapped.
// It replaces the old value on a slice, instead of the default append mode.
func (s *Options) SetRawOverwrite(key string, val interface{}) {
	s.setNxOverwrite(key, val)
}

func (s *Options) setNx(key string, val interface{}) (oldVal interface{}, modi bool) {
	defer s.rw.Unlock()
	s.rw.Lock()
	return s.setNxNoLock(key, val)
}

func (s *Options) setNxOverwrite(key string, val interface{}) (oldVal interface{}, modi bool) {
	defer s.rw.Unlock()
	s.rw.Lock()
	old := s.setAppendMode(modi)
	defer s.resetAppendMode(old)
	return s.setNxNoLock(key, val)
}

func (s *Options) setAppendMode(bn bool) bool {
	b := s.appendMode
	s.appendMode = bn
	return b
}

func (s *Options) resetAppendMode(b bool) { s.appendMode = b }
func (s *Options) setToAppendMode()       { s.appendMode = true }
func (s *Options) setToReplaceMode()      { s.appendMode = false }

func (s *Options) setNxNoLock(key string, val interface{}) (oldVal interface{}, modi bool) {
	if val == nil && s.getMapNoLock(key) != nil {
		// don't set a branch node to nil if it have children.
		return
	}

	oldVal = s.entries[key]
	leaf := isLeaf(oldVal, val)
	if leaf {
		isComparable := (oldVal == nil || oldVal != nil &&
			reflect.TypeOf(oldVal).Comparable()) &&
			(val == nil || (val != nil && reflect.TypeOf(val).Comparable()))
		if isComparable && oldVal != val {
			newVal := s.closelyConvert(val, oldVal)
			s.entries[key] = newVal
			a := strings.Split(key, ".")
			s.mergeMap(s.hierarchy, a[0], "", et(a, 1, newVal))
			s.internalRaiseOnSetCB(key, newVal, oldVal)
			modi = true
			return
		}
		if isEmptySlice(val) && isSlice(oldVal) {
			newVal := s.closelyConvert(val, oldVal)
			s.entries[key] = newVal
			s.internalRaiseOnSetCB(key, newVal, oldVal)
			modi = true
			return
		}
	}

	modi = s.setNxNoLock2(key, oldVal, val, leaf)
	return
}

func (s *Options) setNxNoLock2(key string, oldVal, val interface{}, leaf bool) (modi bool) {
	if leaf {
		if isSlice(oldVal) && isSlice(val) {
			newVal := s.mergeSlice(oldVal, val)
			val = newVal
		} else if isMap(oldVal) && isMap(val) {
			newVal := mergeTwoMapNoRecursive(oldVal, val)
			val = newVal
		}
	}

	newVal := s.closelyConvert(val, oldVal)
	s.entries[key] = newVal
	a := strings.Split(key, ".")
	s.mergeMap(s.hierarchy, a[0], "", et(a, 1, newVal))
	s.internalRaiseOnSetCB(key, newVal, oldVal)
	modi = true
	return
}

func (s *Options) closelyConvert(val, oldVal interface{}) (newVal interface{}) {
	newVal = val
	if oldVal == nil {
		return
	}

	// convert val to newVal by the type of oldVal.
	// if can't, replace newVal with val.
	if ts, ok := val.(string); ok {
		switch tu := oldVal.(type) {
		case *time.Time:
			if tm, err := smartParseTime(ts); err != nil {
				log.Errorf("closelyConvert: can't parse time format: %v", err)
			} else {
				newVal = &tm
			}
		case time.Time:
			if tm, err := smartParseTime(ts); err != nil {
				log.Errorf("closelyConvert: can't parse time format: %v", err)
			} else {
				newVal = tm
			}
		case encoding.TextUnmarshaler: // deal with TextVar here
			if err := tu.UnmarshalText([]byte(ts)); err != nil {
				log.Errorf("closelyConvert: can't unmarshal text from %q: %v", ts, err)
			} else {
				newVal = tu
			}
		}
	}
	return
}

func smartParseTime(str string) (tm time.Time, err error) {
	knownFormats := []string{
		"2006-01-02 15:04:05.999999999 -0700",
		"2006-01-02 15:04:05.999999999Z07:00",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02 15:04:05.999",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"2006/01/02",
		"01/02/2006",
		"01-02",

		"2006-1-2 15:4:5.999999999 -0700",
		"2006-1-2 15:4:5.999999999Z07:00",
		"2006-1-2 15:4:5.999999999",
		"2006-1-2 15:4:5.999",
		"2006-1-2 15:4:5",
		"2006-1-2",
		"2006/1/2",
		"1/2/2006",
		"1-2",

		"15:04:05",
		"15:04",

		time.RFC3339,
		time.RFC3339Nano,
		time.RFC1123Z,
		time.RFC1123,
		time.RFC850,
		time.RFC822Z,
		time.RFC822,
		time.RubyDate,
		time.UnixDate,
		time.ANSIC,
	}

	for _, layout := range knownFormats {
		if tm, err = time.Parse(layout, str); err == nil {
			break
		}
	}
	return
}

func isLeaf(oldVal, val interface{}) (leaf bool) {
	if _, ok := oldVal.(map[string]interface{}); !ok {
		if _, ok = val.(map[string]interface{}); !ok {
			leaf = true
		}
	}
	return
}

func mergeTwoMapNoRecursive(v1, v2 interface{}) (v3 interface{}) {
	if !isMap(v1) || !isMap(v2) {
		return
	}
	return v2 // todo, merge two map as a new one, we need a clean redesigned map merger
}

func isMap(v interface{}) bool {
	x := reflect.ValueOf(v)
	return x.Kind() == reflect.Map
}

func (s *Options) _cvtToTgt(typ reflect.Type, t, x4 reflect.Value) reflect.Value {
	tn := t.Convert(typ)
	var found bool
	for zi := 0; zi < x4.Len(); zi++ {
		if found = reflect.DeepEqual(tn.Interface(), x4.Index(zi).Interface()); found {
			return x4
		}
	}
	return reflect.Append(x4, tn)
}

func (s *Options) _cvtToStr(typ reflect.Type, t, x4 reflect.Value) reflect.Value {
	str := fmt.Sprintf("%v", t.Interface())
	var found bool
	for zi := 0; zi < x4.Len(); zi++ {
		if found = reflect.DeepEqual(str, x4.Index(zi).Interface()); found {
			return x4
		}
	}
	return reflect.Append(x4, reflect.ValueOf(str))
}

func (s *Options) mergeSlice(v1, v2 interface{}) (v3 interface{}) {
	if !isSlice(v1) || !isSlice(v2) {
		return
	}

	if !s.appendMode {
		return v2 // just replace the old value
	}

	x1, x2 := reflect.ValueOf(v1), reflect.ValueOf(v2)
	if et1, et2 := x1.Type().Elem(), x2.Type().Elem(); et1 == et2 || et1.ConvertibleTo(et2) {
		// x1 = reflect.AppendSlice(x1, x2)
		// v3 = x1.Interface()
		// } else if et1.ConvertibleTo(et2) {
		typ := et1
		if et1.Kind() < et2.Kind() {
			typ = et2
		}

		x3 := reflect.New(reflect.SliceOf(typ)) // reflect.MakeSlice(typ, x1.Len(), x1.Len())
		x4 := x3.Elem()
		for i := 0; i < x1.Len(); i++ {
			x4 = reflect.Append(x4, x1.Index(i).Convert(typ))
		}
		for i := 0; i < x2.Len(); i++ {
			t := x2.Index(i)
			if canConvert(t, typ) {
				x4 = s._cvtToTgt(typ, t, x4)
			} else if typ.Kind() == reflect.String {
				x4 = s._cvtToStr(typ, t, x4)
			}

			// Just for go1.17+
			//
			// if x2.Index(i).CanConvert(typ) {
			//	x4 = reflect.Append(x4, x2.Index(i).Convert(typ))
			// } else {
			//	if typ.Kind() == reflect.String {
			//		str := fmt.Sprintf("%v", x2.Interface())
			//		x4 = reflect.Append(x4, reflect.ValueOf(str))
			//	} else {
			//		ferr("cannot convert '%v' to type: %v", x2.Interface(), typ.Kind())
			//	}
			// }
		}
		v3 = x4.Interface()
	} else {
		ferr("cannot convert type between two slice elements, cannot merge slice %v -> %v", v2, v1)
		v3 = v1
	}
	return
}

// canConvert reports whether the value v can be converted to type t.
// If v.CanConvert(t) returns true then v.Convert(t) will not panic.
func canConvert(v reflect.Value, t reflect.Type) bool {
	vt := v.Type()
	if !vt.ConvertibleTo(t) { //nolint:gosimple //keep it
		return false
	}

	// go1.17+ only:
	//
	// // Currently the only conversion that is OK in terms of type
	// // but that can panic depending on the value is converting
	// // from slice to pointer-to-array.
	// if vt.Kind() == reflect.Slice && t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Array {
	//	n := t.Elem().Len()
	//	h := (*unsafeheader.Slice)(v.ptr)
	//	if n > h.Len {
	//		return false
	//	}
	// }
	return true
}

func isSlice(v interface{}) bool {
	x := reflect.ValueOf(v)
	return x.Kind() == reflect.Slice
}

func isEmptySlice(v interface{}) bool {
	x := reflect.ValueOf(v)
	if x.Kind() == reflect.Slice {
		return x.Len() == 0
	}
	return false
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
			switch vm := val.(type) {
			case map[string]interface{}:
				for k, v := range vm {
					zm = s.mergeMap(zm, k, path, v)
				}
				// hierarchy[key] = zm
				// s.entries[path] = zm
				val = zm
			case map[interface{}]interface{}:
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
		if _, ok = val.(map[string]interface{}); !ok {
			leaf = true
		}
	}
	if leaf {
		iscomparable := oldval != nil && reflect.TypeOf(oldval).Comparable() &&
			val != nil && reflect.TypeOf(val).Comparable()
		if iscomparable {
			if oldval != val {
				// defer s.rw.Unlock()
				// s.rw.Lock()
				s.entries[path] = val
				m[key] = val

				s.internalRaiseOnMergingSetCB(path, val, oldval)
				// Logger.Debugf("%%-> s.entries[%q] = m[%q] = %v", path, key, val)
				return
			}
		}
	}
	s.entries[path] = val
	m[key] = val
}

func (s *Options) setCB(onMergingSet, onSet func(keyPath string, value, oldVal interface{})) {
	s.rwCB.Lock()
	defer s.rwCB.Unlock()
	s.onMergingSet = onMergingSet
	s.onSet = onSet
}

func (s *Options) internalRaiseOnMergingSetCB(path string, val, oldval interface{}) {
	if s.onMergingSet != nil {
		s.rwCB.RLock()
		defer s.rwCB.RUnlock()
		s.onMergingSet(path, val, oldval)
	}
}

func (s *Options) internalRaiseOnSetCB(path string, val, oldval interface{}) {
	if s.onSet != nil {
		s.rwCB.RLock()
		defer s.rwCB.RUnlock()
		s.onSet(path, val, oldval)
	}
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

// Flush writes all changes back to the alternative config file
func (s *Options) Flush() {
	if len(s.usedAlterConfigFile) > 0 {
		var err error
		var fi os.FileInfo
		if fi, err = os.Stat(os.ExpandEnv(s.usedAlterConfigFile)); err == nil && dir.IsModeWriteOwner(fi.Mode()) {
			// // str := AsYaml() // s.DumpAsString(false)
			// var b []byte
			// var err error
			// obj := s.GetHierarchyList()
			// defer handleSerializeError(&err)
			// b, err = yaml.Marshal(obj)

			var updated bool
			var m map[string]interface{}
			// var mc map[string]interface{}
			m, err = s.loadConfigFileAsMap(s.usedAlterConfigFile)
			if err == nil {
				updated, _, err = s.updateMap("", m)
				if updated && err == nil {
					var b []byte
					defer handleSerializeError(&err)
					b, err = yaml.Marshal(m)

					err = dir.WriteFile(s.usedAlterConfigFile, b, 0o600)
					if err == nil {
						flog("config file %q updated.", s.usedAlterConfigFile)
					}
				}
			}
		}

		if err != nil {
			log.Errorf("err: %v", err)
		}
	}
}

func (s *Options) updateMap(keyDot string, m map[string]interface{}) (updated bool, mc map[string]interface{}, err error) {
	for k, v := range m {
		key := mxIx(keyDot, k)
		if v == nil {
			continue
			// mc[k], m[k] = v, v
			// nothing to do
		}
		switch vm := v.(type) {
		case map[interface{}]interface{}:
			if err = s.updateIxMap(key, vm); err != nil {
				return
			}
		case map[string]interface{}:
			if updated, mc, err = s.updateMap(key, vm); err != nil {
				return
			}
		default:
			if sv, ok := s.entries[key]; ok {
				tsv, tv := reflect.TypeOf(sv), reflect.TypeOf(v)
				ne := !tsv.Comparable() || !tv.Comparable()
				if tsv.Comparable() && tv.Comparable() && sv != v {
					ne = true
				}
				if ne {
					updated = true
					if mc == nil {
						mc = make(map[string]interface{})
					}
					mc[k], m[k] = sv, sv
				}
			}
		}
	}
	return
}

func (s *Options) updateIxMap(key string, vm map[interface{}]interface{}) (err error) {
	// TODO for k, v in map[interface{}]interface{}
	return
}

// Reset the exists `Options`, so that you could follow a `LoadConfigFile()` with it.
func (s *Options) Reset() {
	defer s.rw.Unlock()
	s.rw.Lock()

	s.entries = nil
	s.hierarchy = nil
	time.Sleep(100 * time.Millisecond)
	s.entries = make(map[string]interface{})
	s.hierarchy = make(map[string]interface{})
}

func mx(pre, k string) string {
	if pre == "" {
		return k
	}
	return pre + "." + k
}

func mxIx(pre string, k interface{}) string {
	if pre == "" {
		return fmt.Sprintf("%v", k)
	}
	return fmt.Sprintf("%v.%v", pre, k)
}

func (s *Options) loopMapMap(keyDot string, m map[string]map[string]interface{}) (err error) {
	for k, v := range m {
		if err = s.loopMap(mx(keyDot, k), v); err != nil {
			return
		}
	}
	return
}

func (s *Options) loopMap(keyDot string, m map[string]interface{}) (err error) {
	defer s.mapOrphans()
	for k, v := range m {
		switch vm := v.(type) {
		case map[interface{}]interface{}:
			key := mx(keyDot, k)
			if err = s.loopIxMap(key, vm); err != nil {
				return
			}
		case map[string]interface{}:
			key := mx(keyDot, k)
			if err = s.loopMap(key, vm); err != nil {
				return
			}
		default:
			// s.SetNx(mx(kDot, k), v)
			key := mxIx(keyDot, k)
			if oldVal, modified := s.setNx(key, v); modified {
				s.rw.Lock()
				v = s.entries[key]
				s.rw.Unlock()
				s.internalRaiseOnMergingSetCB(k, v, oldVal)
			}
		}
	}
	return
}

func (s *Options) mapOrphans() {
	s.rw.Lock()
	defer s.rw.Unlock()

	// flog("mapOrphans")
	if s.batchMerging {
		return
	}

retryChecking:
	var keysSorted []string
	for k := range s.entries {
		keysSorted = append(keysSorted, k)
	}
	tool.SortAsDottedSliceReverse(keysSorted)

	for _, k := range keysSorted {
		// flog("mapOrphans: %v => %v", k, v)
		keys := strings.Split(k, ".")
		for i := 1; i < len(keys); i++ {
			ks := strings.Join(keys[:i], ".")
			if vz, ok := s.entries[ks]; !ok || vz == nil {
				flog("mapOrphans: %v: %q not exists!!", k, ks)

				vv := make(map[string]interface{})
				for kk, v := range s.entries {
					kks := strings.Split(kk, ".")
					if strings.Contains(kk, ks) && len(kks) == i+1 {
						vv[kks[len(kks)-1]] = v
					}
				}
				s.entries[ks] = vv
				goto retryChecking
			}
		}
	}

	if detects.InDebugging() {
		keysSorted = nil
		for k := range s.entries {
			keysSorted = append(keysSorted, k)
		}
		tool.SortAsDottedSliceReverse(keysSorted)
		flog("mapOrphans: END")
	}
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
			if oldval, modi := s.setNx(key, v); modi {
				s.internalRaiseOnMergingSetCB(key, v, oldval)
			}
		}
	}
	return
}

// DumpAsString for debugging.
func (s *Options) DumpAsString(showType bool) (str string) {
	k3 := make([]string, 0)
	for k := range s.entries {
		k3 = append(k3, k)
	}
	sort.Strings(k3)

	for _, k := range k3 {
		if showType {
			str += fmt.Sprintf("%-48v => %v (%T)\n", k, s.entries[k], s.entries[k])
		} else {
			str += fmt.Sprintf("%-48v => %v\n", k, s.entries[k])
		}
	}

	str += s.dumpAsStringYamlPart(showType)
	return
}

func (s *Options) dumpAsStringYamlPart(showType bool) (str string) {
	str += "---------------------------------\n"
	var err error
	var sb strings.Builder
	defer handleSerializeError(&err)
	e := yaml.NewEncoder(&sb)
	e.SetIndent(2)
	if err = e.Encode(s.hierarchy); err == nil {
		err = e.Close()
		// var b []byte
		// b, err = yaml.Marshal(s.hierarchy)
		if err == nil {
			if s.GetBoolEx("raw") {
				str += sb.String() // string(b)
			} else {
				ss := sb.String() // string(b)
				ss = os.ExpandEnv(ss)
				str += ss
			}
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

// SaveCheckpoint make a snapshot of the current Option Store
//
// You may ResetOptions after SaveCheckpoint:
//
//	func x(aMap map[string]interface{}){
//	    defer cmdr.RestoreCheckpoint()
//	    cmdr.SaveCheckpoint()
//	    cmdr.ResetOptions()
//	    cmdr.MergeWith(map[string]interface{}{
//	      "app": map[string]interface{}{
//	        conf.AppName: map[string]interface{}{
//	          "a-map": aMap,
//	        }
//	      }
//	    }
//	    cmdr.SaveAsYaml("a-setting.yml")
//	}
func (s *Options) SaveCheckpoint() (err error) {
	defer s.rw.RUnlock()
	s.rw.RLock()
	no := newOptions()
	if err = evendeep.DefaultCopyController.CopyTo(s, no); err == nil {
		w := internalGetWorker()
		w.savedOptions = append(w.savedOptions, no)
	}
	// if err = StandardCopier.Copy(no, s); err == nil {
	// 	w := internalGetWorker()
	// 	w.savedOptions = append(w.savedOptions, no)
	// }
	return
}

// RestoreCheckpoint restore 1 or n checkpoint(s) from snapshots history.
// see also SaveCheckpoint
func (s *Options) RestoreCheckpoint(n ...int) (err error) {
	nn := 1
	for _, n1 := range n {
		nn = n1
	}

	if 1 <= nn && nn <= s.CheckpointSize() {
		w := internalGetWorker()
		var no *Options
		for i := 0; i < nn; i++ {
			no = w.savedOptions[len(w.savedOptions)-1]
		}
		w.savedOptions = w.savedOptions[0 : len(w.savedOptions)-nn]

		defer s.rw.RUnlock()
		s.rw.RLock()
		w.rxxtOptions = no
	}
	return
}

// ClearCheckpoints removes all checkpoints from snapshot history
// see also SaveCheckpoint
func (s *Options) ClearCheckpoints() {
	w := internalGetWorker()
	w.savedOptions = nil
}

// CheckpointSize returns how many snapshots made
// see also SaveCheckpoint
func (s *Options) CheckpointSize() int { return len(internalGetWorker().savedOptions) }

// SaveCheckpoint make a snapshot of the current Option Store
//
// You may ResetOptions after SaveCheckpoint:
//
//	func x(aMap map[string]interface{}){
//	    defer cmdr.RestoreCheckpoint()
//	    cmdr.SaveCheckpoint()
//	    cmdr.ResetOptions()
//	    cmdr.MergeWith(map[string]interface{}{
//	      "app": map[string]interface{}{
//	        conf.AppName: map[string]interface{}{
//	          "a-map": aMap,
//	        }
//	      }
//	    }
//	    cmdr.SaveAsYaml("a-setting.yml")
//	}
func SaveCheckpoint() (err error) { return currentOptions().SaveCheckpoint() }

// RestoreCheckpoint restore 1 or n checkpoint(s) from snapshots history.
// see also SaveCheckpoint
func RestoreCheckpoint(n ...int) (err error) {
	return currentOptions().RestoreCheckpoint(n...)
}

// ClearCheckpoints removes all checkpoints from snapshot history
// see also SaveCheckpoint
func ClearCheckpoints() { currentOptions().ClearCheckpoints() }

// CheckpointSize returns how many snapshots made
// see also SaveCheckpoint
func CheckpointSize() int { return currentOptions().CheckpointSize() }
