/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"bytes"
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

//
//
//

// GetBool returns the bool value of an `Option` key. Such as:
// ```golang
// cmdr.GetBool("app.logger.enable") => true,...
// ```
//
func GetBool(key string) bool {
	return uniqueWorker.rxxtOptions.GetBoolEx(key, false)
}

// GetBoolEx returns the bool value of an `Option` key. Such as:
// ```golang
// cmdr.GetBoolEx("app.logger.enable", false) => true,...
// ```
//
func GetBoolEx(key string, defaultVal bool) bool {
	return uniqueWorker.rxxtOptions.GetBoolEx(key, defaultVal)
}

// GetBoolP returns the bool value of an `Option` key. Such as:
// ```golang
// cmdr.GetBoolP("app.logger", "enable") => true,...
// ```
func GetBoolP(prefix, key string) bool {
	return uniqueWorker.rxxtOptions.GetBoolEx(fmt.Sprintf("%s.%s", prefix, key), false)
}

// GetBoolExP returns the bool value of an `Option` key. Such as:
// ```golang
// cmdr.GetBoolExP("app.logger", "enable", false) => true,...
// ```
func GetBoolExP(prefix, key string, defaultVal bool) bool {
	return uniqueWorker.rxxtOptions.GetBoolEx(fmt.Sprintf("%s.%s", prefix, key), defaultVal)
}

// GetBoolR returns the bool value of an `Option` key with [WrapWithRxxtPrefix]. Such as:
// ```golang
// cmdr.GetBoolR("logger.enable") => true,...
// ```
//
func GetBoolR(key string) bool {
	return uniqueWorker.rxxtOptions.GetBoolEx(wrapWithRxxtPrefix(key), false)
}

// GetBoolExR returns the bool value of an `Option` key with [WrapWithRxxtPrefix]. Such as:
// ```golang
// cmdr.GetBoolExR("logger.enable", false) => true,...
// ```
//
func GetBoolExR(key string, defaultVal bool) bool {
	return uniqueWorker.rxxtOptions.GetBoolEx(wrapWithRxxtPrefix(key), defaultVal)
}

// GetBoolRP returns the bool value of an `Option` key with [WrapWithRxxtPrefix]. Such as:
// ```golang
// cmdr.GetBoolRP("logger", "enable") => true,...
// ```
func GetBoolRP(prefix, key string) bool {
	return uniqueWorker.rxxtOptions.GetBoolEx(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), false)
}

// GetBoolExRP returns the bool value of an `Option` key with [WrapWithRxxtPrefix]. Such as:
// ```golang
// cmdr.GetBoolExRP("logger", "enable", false) => true,...
// ```
func GetBoolExRP(prefix, key string, defaultVal bool) bool {
	return uniqueWorker.rxxtOptions.GetBoolEx(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal)
}

// GetInt returns the int value of an `Option` key.
func GetInt(key string) int {
	return int(uniqueWorker.rxxtOptions.GetIntEx(key, 0))
}

// GetIntEx returns the int value of an `Option` key.
func GetIntEx(key string, defaultVal int) int {
	return int(uniqueWorker.rxxtOptions.GetIntEx(key, int64(defaultVal)))
}

// GetIntP returns the int value of an `Option` key.
func GetIntP(prefix, key string) int {
	return int(uniqueWorker.rxxtOptions.GetIntEx(fmt.Sprintf("%s.%s", prefix, key), 0))
}

// GetIntExP returns the int value of an `Option` key.
func GetIntExP(prefix, key string, defaultVal int) int {
	return int(uniqueWorker.rxxtOptions.GetIntEx(fmt.Sprintf("%s.%s", prefix, key), int64(defaultVal)))
}

// GetIntR returns the int value of an `Option` key with [WrapWithRxxtPrefix].
func GetIntR(key string) int {
	return int(uniqueWorker.rxxtOptions.GetIntEx(wrapWithRxxtPrefix(key), 0))
}

// GetIntExR returns the int value of an `Option` key with [WrapWithRxxtPrefix].
func GetIntExR(key string, defaultVal int) int {
	return int(uniqueWorker.rxxtOptions.GetIntEx(wrapWithRxxtPrefix(key), int64(defaultVal)))
}

// GetIntRP returns the int value of an `Option` key with [WrapWithRxxtPrefix].
func GetIntRP(prefix, key string) int {
	return int(uniqueWorker.rxxtOptions.GetIntEx(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), 0))
}

// GetIntExRP returns the int value of an `Option` key with [WrapWithRxxtPrefix].
func GetIntExRP(prefix, key string, defaultVal int) int {
	return int(uniqueWorker.rxxtOptions.GetIntEx(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), int64(defaultVal)))
}

// GetInt64 returns the int64 value of an `Option` key.
func GetInt64(key string) int64 {
	return uniqueWorker.rxxtOptions.GetIntEx(key, 0)
}

// GetInt64Ex returns the int64 value of an `Option` key.
func GetInt64Ex(key string, defaultVal int64) int64 {
	return uniqueWorker.rxxtOptions.GetIntEx(key, defaultVal)
}

// GetInt64P returns the int64 value of an `Option` key.
func GetInt64P(prefix, key string) int64 {
	return uniqueWorker.rxxtOptions.GetIntEx(fmt.Sprintf("%s.%s", prefix, key), 0)
}

// GetInt64ExP returns the int64 value of an `Option` key.
func GetInt64ExP(prefix, key string, defaultVal int64) int64 {
	return uniqueWorker.rxxtOptions.GetIntEx(fmt.Sprintf("%s.%s", prefix, key), defaultVal)
}

// GetInt64R returns the int64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetInt64R(key string) int64 {
	return uniqueWorker.rxxtOptions.GetIntEx(wrapWithRxxtPrefix(key), 0)
}

// GetInt64ExR returns the int64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetInt64ExR(key string, defaultVal int64) int64 {
	return uniqueWorker.rxxtOptions.GetIntEx(wrapWithRxxtPrefix(key), defaultVal)
}

// GetInt64RP returns the int64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetInt64RP(prefix, key string) int64 {
	return uniqueWorker.rxxtOptions.GetIntEx(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), 0)
}

// GetInt64ExRP returns the int64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetInt64ExRP(prefix, key string, defaultVal int64) int64 {
	return uniqueWorker.rxxtOptions.GetIntEx(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal)
}

// GetUint returns the uint value of an `Option` key.
func GetUint(key string) uint {
	return uint(uniqueWorker.rxxtOptions.GetUintEx(key, 0))
}

// GetUintEx returns the uint value of an `Option` key.
func GetUintEx(key string, defaultVal uint) uint {
	return uint(uniqueWorker.rxxtOptions.GetUintEx(key, uint64(defaultVal)))
}

// GetUintP returns the uint value of an `Option` key.
func GetUintP(prefix, key string) uint {
	return uint(uniqueWorker.rxxtOptions.GetUintEx(fmt.Sprintf("%s.%s", prefix, key), 0))
}

// GetUintExP returns the uint value of an `Option` key.
func GetUintExP(prefix, key string, defaultVal uint) uint {
	return uint(uniqueWorker.rxxtOptions.GetUintEx(fmt.Sprintf("%s.%s", prefix, key), uint64(defaultVal)))
}

// GetUintR returns the uint value of an `Option` key with [WrapWithRxxtPrefix].
func GetUintR(key string) uint {
	return uint(uniqueWorker.rxxtOptions.GetUintEx(wrapWithRxxtPrefix(key), 0))
}

// GetUintExR returns the uint value of an `Option` key with [WrapWithRxxtPrefix].
func GetUintExR(key string, defaultVal uint) uint {
	return uint(uniqueWorker.rxxtOptions.GetUintEx(wrapWithRxxtPrefix(key), uint64(defaultVal)))
}

// GetUintRP returns the uint value of an `Option` key with [WrapWithRxxtPrefix].
func GetUintRP(prefix, key string) uint {
	return uint(uniqueWorker.rxxtOptions.GetUintEx(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), 0))
}

// GetUintExRP returns the uint value of an `Option` key with [WrapWithRxxtPrefix].
func GetUintExRP(prefix, key string, defaultVal uint) uint {
	return uint(uniqueWorker.rxxtOptions.GetUintEx(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), uint64(defaultVal)))
}

// GetUint64 returns the uint64 value of an `Option` key.
func GetUint64(key string) uint64 {
	return uniqueWorker.rxxtOptions.GetUintEx(key, 0)
}

// GetUint64Ex returns the uint64 value of an `Option` key.
func GetUint64Ex(key string, defaultVal uint64) uint64 {
	return uniqueWorker.rxxtOptions.GetUintEx(key, defaultVal)
}

// GetUint64P returns the uint64 value of an `Option` key.
func GetUint64P(prefix, key string) uint64 {
	return uniqueWorker.rxxtOptions.GetUintEx(fmt.Sprintf("%s.%s", prefix, key), 0)
}

// GetUint64ExP returns the uint64 value of an `Option` key.
func GetUint64ExP(prefix, key string, defaultVal uint64) uint64 {
	return uniqueWorker.rxxtOptions.GetUintEx(fmt.Sprintf("%s.%s", prefix, key), defaultVal)
}

// GetUint64R returns the uint64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetUint64R(key string) uint64 {
	return uniqueWorker.rxxtOptions.GetUintEx(wrapWithRxxtPrefix(key), 0)
}

// GetUint64ExR returns the uint64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetUint64ExR(key string, defaultVal uint64) uint64 {
	return uniqueWorker.rxxtOptions.GetUintEx(wrapWithRxxtPrefix(key), defaultVal)
}

// GetUint64RP returns the uint64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetUint64RP(prefix, key string) uint64 {
	return uniqueWorker.rxxtOptions.GetUintEx(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), 0)
}

// GetUint64ExRP returns the uint64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetUint64ExRP(prefix, key string, defaultVal uint64) uint64 {
	return uniqueWorker.rxxtOptions.GetUintEx(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal)
}

// GetFloat32 returns the float32 value of an `Option` key.
func GetFloat32(key string) float32 {
	return float32(uniqueWorker.rxxtOptions.GetFloat32Ex(key, 0))
}

// GetFloat32Ex returns the float32 value of an `Option` key.
func GetFloat32Ex(key string, defaultVal float32) float32 {
	return float32(uniqueWorker.rxxtOptions.GetFloat32Ex(key, defaultVal))
}

// GetFloat32P returns the float32 value of an `Option` key.
func GetFloat32P(prefix, key string) float32 {
	return float32(uniqueWorker.rxxtOptions.GetFloat32Ex(fmt.Sprintf("%s.%s", prefix, key), 0))
}

// GetFloat32ExP returns the float32 value of an `Option` key.
func GetFloat32ExP(prefix, key string, defaultVal float32) float32 {
	return float32(uniqueWorker.rxxtOptions.GetFloat32Ex(fmt.Sprintf("%s.%s", prefix, key), defaultVal))
}

// GetFloat32R returns the float32 value of an `Option` key with [WrapWithRxxtPrefix].
func GetFloat32R(key string) float32 {
	return float32(uniqueWorker.rxxtOptions.GetFloat32Ex(wrapWithRxxtPrefix(key), 0))
}

// GetFloat32ExR returns the float32 value of an `Option` key with [WrapWithRxxtPrefix].
func GetFloat32ExR(key string, defaultVal float32) float32 {
	return float32(uniqueWorker.rxxtOptions.GetFloat32Ex(wrapWithRxxtPrefix(key), defaultVal))
}

// GetFloat32RP returns the float32 value of an `Option` key with [WrapWithRxxtPrefix].
func GetFloat32RP(prefix, key string) float32 {
	return float32(uniqueWorker.rxxtOptions.GetFloat32Ex(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), 0))
}

// GetFloat32ExRP returns the float32 value of an `Option` key with [WrapWithRxxtPrefix].
func GetFloat32ExRP(prefix, key string, defaultVal float32) float32 {
	return float32(uniqueWorker.rxxtOptions.GetFloat32Ex(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal))
}

// GetFloat64 returns the float64 value of an `Option` key.
func GetFloat64(key string) float64 {
	return uniqueWorker.rxxtOptions.GetFloat64Ex(key, 0)
}

// GetFloat64Ex returns the float64 value of an `Option` key.
func GetFloat64Ex(key string, defaultVal float64) float64 {
	return uniqueWorker.rxxtOptions.GetFloat64Ex(key, defaultVal)
}

// GetFloat64P returns the float64 value of an `Option` key.
func GetFloat64P(prefix, key string) float64 {
	return uniqueWorker.rxxtOptions.GetFloat64Ex(fmt.Sprintf("%s.%s", prefix, key), 0)
}

// GetFloat64ExP returns the float64 value of an `Option` key.
func GetFloat64ExP(prefix, key string, defaultVal float64) float64 {
	return uniqueWorker.rxxtOptions.GetFloat64Ex(fmt.Sprintf("%s.%s", prefix, key), defaultVal)
}

// GetFloat64R returns the float64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetFloat64R(key string) float64 {
	return uniqueWorker.rxxtOptions.GetFloat64Ex(wrapWithRxxtPrefix(key), 0)
}

// GetFloat64ExR returns the float64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetFloat64ExR(key string, defaultVal float64) float64 {
	return uniqueWorker.rxxtOptions.GetFloat64Ex(wrapWithRxxtPrefix(key), defaultVal)
}

// GetFloat64RP returns the float64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetFloat64RP(prefix, key string) float64 {
	return uniqueWorker.rxxtOptions.GetFloat64Ex(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), 0)
}

// GetFloat64ExRP returns the float64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetFloat64ExRP(prefix, key string, defaultVal float64) float64 {
	return uniqueWorker.rxxtOptions.GetFloat64Ex(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal)
}

// GetString returns the string value of an `Option` key.
func GetString(key string) string {
	return uniqueWorker.rxxtOptions.GetString(key, "")
}

// GetStringEx returns the string value of an `Option` key.
func GetStringEx(key, defaultVal string) string {
	return uniqueWorker.rxxtOptions.GetString(key, defaultVal)
}

// GetStringP returns the string value of an `Option` key.
func GetStringP(prefix, key string) string {
	return uniqueWorker.rxxtOptions.GetString(fmt.Sprintf("%s.%s", prefix, key), "")
}

// GetStringExP returns the string value of an `Option` key.
func GetStringExP(prefix, key, defaultVal string) string {
	return uniqueWorker.rxxtOptions.GetString(fmt.Sprintf("%s.%s", prefix, key), defaultVal)
}

// GetStringR returns the string value of an `Option` key with [WrapWithRxxtPrefix].
func GetStringR(key string) string {
	return uniqueWorker.rxxtOptions.GetString(wrapWithRxxtPrefix(key), "")
}

// GetStringExR returns the string value of an `Option` key with [WrapWithRxxtPrefix].
func GetStringExR(key, defaultVal string) string {
	return uniqueWorker.rxxtOptions.GetString(wrapWithRxxtPrefix(key), defaultVal)
}

// GetStringRP returns the string value of an `Option` key with [WrapWithRxxtPrefix].
func GetStringRP(prefix, key string) string {
	return uniqueWorker.rxxtOptions.GetString(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), "")
}

// GetStringExRP returns the string value of an `Option` key with [WrapWithRxxtPrefix].
func GetStringExRP(prefix, key, defaultVal string) string {
	return uniqueWorker.rxxtOptions.GetString(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal)
}

// GetStringSlice returns the string slice value of an `Option` key.
func GetStringSlice(key string, defaultVal ...string) []string {
	return uniqueWorker.rxxtOptions.GetStringSlice(key, defaultVal...)
}

// GetStringSliceP returns the string slice value of an `Option` key.
func GetStringSliceP(prefix, key string, defaultVal ...string) []string {
	return uniqueWorker.rxxtOptions.GetStringSlice(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetStringSliceR returns the string slice value of an `Option` key with [WrapWithRxxtPrefix].
func GetStringSliceR(key string, defaultVal ...string) []string {
	return uniqueWorker.rxxtOptions.GetStringSlice(wrapWithRxxtPrefix(key), defaultVal...)
}

// GetStringSliceRP returns the string slice value of an `Option` key with [WrapWithRxxtPrefix].
func GetStringSliceRP(prefix, key string, defaultVal ...string) []string {
	return uniqueWorker.rxxtOptions.GetStringSlice(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// Get returns the generic value of an `Option` key with [WrapWithRxxtPrefix]. Such as:
// ```golang
// cmdr.Get("app.logger.level") => 'DEBUG',...
// ```
//
func Get(key string) interface{} {
	return uniqueWorker.rxxtOptions.Get(key)
}

// GetR returns the generic value of an `Option` key with [WrapWithRxxtPrefix]. Such as:
// ```golang
// cmdr.GetR("logger.level") => 'DEBUG',...
// ```
//
func GetR(key string) interface{} {
	return uniqueWorker.rxxtOptions.Get(wrapWithRxxtPrefix(key))
}

// GetMap an `Option` by key string, it returns a hierarchy map or nil
func GetMap(key string) map[string]interface{} {
	return uniqueWorker.rxxtOptions.GetMap(key)
}

// GetMapR an `Option` by key string with [WrapWithRxxtPrefix], it returns a hierarchy map or nil
func GetMapR(key string) map[string]interface{} {
	return uniqueWorker.rxxtOptions.GetMap(wrapWithRxxtPrefix(key))
}

// GetSectionFrom returns error while cannot yaml Marshal and Unmarshal
// `cmdr.GetSectionFrom(sectionKeyPath, &holder)` could load all sub-tree nodes from sectionKeyPath and transform them into holder structure, such as:
// ```go
//  type ServerConfig struct {
//    Port int
//    HttpMode int
//    EnableTls bool
//  }
//  var serverConfig = new(ServerConfig)
//  cmdr.GetSectionFrom("server", &serverConfig)
//  assert serverConfig.Port == 7100
// ```
func GetSectionFrom(sectionKeyPath string, holder interface{}) (err error) {
	fObj := GetMapR(sectionKeyPath)
	if fObj != nil {
		var b []byte
		b, err = yaml.Marshal(fObj)
		if err == nil {
			err = yaml.Unmarshal(b, holder)
			// if err == nil {
			// 	logrus.Debugf("configuration section got: %v", configHolder)
			// }
		}
	}
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
func (s *Options) GetBoolEx(key string, defaultVal bool) (ret bool) {
	switch strings.ToLower(s.GetString(key, "")) {
	case "1", "y", "t", "yes", "true", "ok", "on":
		ret = true
	case "":
		ret = defaultVal
	}
	return
}

// GetIntEx returns the int64 value of an `Option` key.
func (s *Options) GetIntEx(key string, defaultVal int64) (ir int64) {
	ir = defaultVal
	if ir64, err := strconv.ParseInt(s.GetString(key, ""), 10, 64); err == nil {
		ir = ir64
	}
	return
}

// GetUintEx returns the uint64 value of an `Option` key.
func (s *Options) GetUintEx(key string, defaultVal uint64) (ir uint64) {
	ir = defaultVal
	if ir64, err := strconv.ParseUint(s.GetString(key, ""), 10, 64); err == nil {
		ir = ir64
	}
	return
}

// GetFloat32Ex returns the float32 value of an `Option` key.
func (s *Options) GetFloat32Ex(key string, defaultVal float32) (ir float32) {
	ir = defaultVal
	if ir64, err := strconv.ParseFloat(s.GetString(key, ""), 10); err == nil {
		ir = float32(ir64)
	}
	return
}

// GetFloat64Ex returns the float64 value of an `Option` key.
func (s *Options) GetFloat64Ex(key string, defaultVal float64) (ir float64) {
	ir = defaultVal
	if ir64, err := strconv.ParseFloat(s.GetString(key, ""), 10); err == nil {
		ir = ir64
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

// GetIntSlice returns the int slice value of an `Option` key.
func GetIntSlice(key string, defaultVal ...int) []int {
	return uniqueWorker.rxxtOptions.GetIntSlice(key, defaultVal...)
}

// GetIntSliceP returns the int slice value of an `Option` key.
func GetIntSliceP(prefix, key string, defaultVal ...int) []int {
	return uniqueWorker.rxxtOptions.GetIntSlice(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetIntSliceR returns the int slice value of an `Option` key with [WrapWithRxxtPrefix].
func GetIntSliceR(key string, defaultVal ...int) []int {
	return uniqueWorker.rxxtOptions.GetIntSlice(wrapWithRxxtPrefix(key), defaultVal...)
}

// GetIntSliceRP returns the int slice value of an `Option` key with [WrapWithRxxtPrefix].
func GetIntSliceRP(prefix, key string, defaultVal ...int) []int {
	return uniqueWorker.rxxtOptions.GetIntSlice(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
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

// // GetUintSlice returns the int slice value of an `Option` key.
// func GetUintSlice(key string, defaultVal ...uint) []uint {
// 	return uniqueWorker.rxxtOptions.GetUintSlice(key, defaultVal...)
// }
//
// // GetUintSliceP returns the int slice value of an `Option` key.
// func GetUintSliceP(prefix, key string, defaultVal ...uint) []uint {
// 	return uniqueWorker.rxxtOptions.GetUintSlice(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
// }
//
// // GetUintSliceR returns the int slice value of an `Option` key with [WrapWithRxxtPrefix].
// func GetUintSliceR(key string, defaultVal ...uint) []uint {
// 	return uniqueWorker.rxxtOptions.GetUintSlice(wrapWithRxxtPrefix(key), defaultVal...)
// }
//
// // GetUintSliceRP returns the int slice value of an `Option` key with [WrapWithRxxtPrefix].
// func GetUintSliceRP(prefix, key string, defaultVal ...uint) []uint {
// 	return uniqueWorker.rxxtOptions.GetUintSlice(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
// }

// GetInt64Slice returns the int slice value of an `Option` key.
func GetInt64Slice(key string, defaultVal ...int64) []int64 {
	return uniqueWorker.rxxtOptions.GetInt64Slice(key, defaultVal...)
}

// GetInt64SliceP returns the int slice value of an `Option` key.
func GetInt64SliceP(prefix, key string, defaultVal ...int64) []int64 {
	return uniqueWorker.rxxtOptions.GetInt64Slice(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetInt64SliceR returns the int slice value of an `Option` key with [WrapWithRxxtPrefix].
func GetInt64SliceR(key string, defaultVal ...int64) []int64 {
	return uniqueWorker.rxxtOptions.GetInt64Slice(wrapWithRxxtPrefix(key), defaultVal...)
}

// GetInt64SliceRP returns the int slice value of an `Option` key with [WrapWithRxxtPrefix].
func GetInt64SliceRP(prefix, key string, defaultVal ...int64) []int64 {
	return uniqueWorker.rxxtOptions.GetInt64Slice(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
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

// GetUint64Slice returns the int slice value of an `Option` key.
func GetUint64Slice(key string, defaultVal ...uint64) []uint64 {
	return uniqueWorker.rxxtOptions.GetUint64Slice(key, defaultVal...)
}

// GetUint64SliceP returns the int slice value of an `Option` key.
func GetUint64SliceP(prefix, key string, defaultVal ...uint64) []uint64 {
	return uniqueWorker.rxxtOptions.GetUint64Slice(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetUint64SliceR returns the int slice value of an `Option` key with [WrapWithRxxtPrefix].
func GetUint64SliceR(key string, defaultVal ...uint64) []uint64 {
	return uniqueWorker.rxxtOptions.GetUint64Slice(wrapWithRxxtPrefix(key), defaultVal...)
}

// GetUint64SliceRP returns the int slice value of an `Option` key with [WrapWithRxxtPrefix].
func GetUint64SliceRP(prefix, key string, defaultVal ...uint64) []uint64 {
	return uniqueWorker.rxxtOptions.GetUint64Slice(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
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

// GetDuration returns the int slice value of an `Option` key.
func GetDuration(key string) time.Duration {
	return uniqueWorker.rxxtOptions.GetDuration(key, 0)
}

// GetDurationP returns the int slice value of an `Option` key.
func GetDurationP(prefix, key string) time.Duration {
	return uniqueWorker.rxxtOptions.GetDuration(fmt.Sprintf("%s.%s", prefix, key), 0)
}

// GetDurationR returns the int slice value of an `Option` key.
func GetDurationR(key string) time.Duration {
	return uniqueWorker.rxxtOptions.GetDuration(wrapWithRxxtPrefix(key), 0)
}

// GetDurationRP returns the int slice value of an `Option` key.
func GetDurationRP(prefix, key string) time.Duration {
	return uniqueWorker.rxxtOptions.GetDuration(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), 0)
}

// GetDurationEx returns the int slice value of an `Option` key.
func GetDurationEx(key string, defaultVal time.Duration) time.Duration {
	return uniqueWorker.rxxtOptions.GetDuration(key, defaultVal)
}

// GetDurationExP returns the int slice value of an `Option` key.
func GetDurationExP(prefix, key string, defaultVal time.Duration) time.Duration {
	return uniqueWorker.rxxtOptions.GetDuration(fmt.Sprintf("%s.%s", prefix, key), defaultVal)
}

// GetDurationExR returns the int slice value of an `Option` key.
func GetDurationExR(key string, defaultVal time.Duration) time.Duration {
	return uniqueWorker.rxxtOptions.GetDuration(wrapWithRxxtPrefix(key), defaultVal)
}

// GetDurationExRP returns the int slice value of an `Option` key.
func GetDurationExRP(prefix, key string, defaultVal time.Duration) time.Duration {
	return uniqueWorker.rxxtOptions.GetDuration(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal)
}

// GetDuration returns the time duration value of an `Option` key.
func (s *Options) GetDuration(key string, defaultVal time.Duration) (ir time.Duration) {
	str := s.GetString(key, "BAD")
	if str == "BAD" {
		ir = defaultVal
	} else {
		var err error
		if ir, err = time.ParseDuration(str); err != nil {
			ir = defaultVal
		}
	}
	return
}

// GetString returns the string value of an `Option` key.
func (s *Options) GetString(key, defaultVal string) (ret string) {
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
			if v != nil {
				ret = fmt.Sprintf("%v", v)
			}
		}
	} else {
		ret = defaultVal
	}
	return
}

func (s *Options) buildAutomaticEnv(rootCmd *RootCommand) (err error) {
	// prefix := strings.Join(EnvPrefix,"_")
	prefix := uniqueWorker.getPrefix() // strings.Join(RxxtPrefix, ".")
	for key := range s.entries {
		ek := s.envKey(key)
		if v, ok := os.LookupEnv(ek); ok {
			if strings.HasPrefix(key, prefix) {
				s.Set(key[len(prefix)+1:], v)
			} else {
				s.Set(key, v)
			}
		}
	}

	// fmt.Printf("EXE = %v, PWD = %v, CURRDIR = %v\n", GetExecutableDir(), os.Getenv("PWD"), GetCurrentDir())
	_ = os.Setenv("THIS", GetExecutableDir())

	for _, h := range uniqueWorker.afterAutomaticEnv {
		h(rootCmd, s)
	}
	return
}

func (s *Options) envKey(key string) (envkey string) {
	key = strings.ReplaceAll(key, ".", "_")
	key = strings.ReplaceAll(key, "-", "_")
	envkey = strings.Join(append(uniqueWorker.envPrefixes, strings.ToUpper(key)), "_")
	return
}

// WrapWithRxxtPrefix wrap an key with [RxxtPrefix], for [GetXxx(key)] and [GetXxxP(prefix,key)]
func WrapWithRxxtPrefix(key string) string {
	return wrapWithRxxtPrefix(key)
}

func wrapWithRxxtPrefix(key string) string {
	if len(uniqueWorker.rxxtPrefixes) == 0 {
		return key
	}
	p := uniqueWorker.getPrefix() // strings.Join(RxxtPrefix, ".")
	if len(key) == 0 {
		return p
	}
	return p + "." + key
}

// Set set the value of an `Option` key (with prefix auto-wrap). The key MUST not have an `app` prefix. eg:
//
//   cmdr.Set("logger.level", "DEBUG")
//   cmdr.Set("ms.tags.port", 8500)
//   ...
//   cmdr.Set("debug", true)
//   cmdr.GetBool("app.debug") => true
//
//
func Set(key string, val interface{}) {
	uniqueWorker.rxxtOptions.Set(key, val)
}

// SetNx but without prefix auto-wrapped.
// `rxxtPrefix` is a string slice to define the prefix string array, default is ["app"].
// So, cmdr.Set("debug", true) will put an real entry with (`app.debug`, true).
func SetNx(key string, val interface{}) {
	uniqueWorker.rxxtOptions.SetNx(key, val)
}

// Set set the value of an `Option` key. The key MUST not have an `app` prefix. eg:
// ```golang
// cmdr.Set("debug", true)
// cmdr.GetBool("app.debug") => true
// ```
func (s *Options) Set(key string, val interface{}) {
	k := wrapWithRxxtPrefix(key)
	s.SetNx(k, val)
}

// SetNx but without prefix auto-wrapped.
// `rxxtPrefix` is a string slice to define the prefix string array, default is ["app"].
// So, cmdr.Set("debug", true) will put an real entry with (`app.debug`, true).
func (s *Options) SetNx(key string, val interface{}) {
	defer s.rw.Unlock()
	s.rw.Lock()

	if val == nil {
		if s.getMapNoLock(key) != nil {
			// don't set a branch node to nil if it have children.
			return
		}
	}

	s.entries[key] = val
	a := strings.Split(key, ".")
	s.mergeMap(s.hierarchy, a[0], "", et(a, 1, val))
}

// MergeWith will merge a map recursive.
// You could merge a yaml/json/toml options into cmdr Hierarchy Options.
func MergeWith(m map[string]interface{}) (err error) {
	err = uniqueWorker.rxxtOptions.MergeWith(m)
	return
}

// MergeWith will merge a map recursive.
func (s *Options) MergeWith(m map[string]interface{}) (err error) {
	for k, v := range m {
		s.mergeMap(s.hierarchy, k, "", v)
	}
	return
}

func (s *Options) mergeMap(m map[string]interface{}, key, path string, val interface{}) map[string]interface{} {
	if len(path) > 0 {
		path = fmt.Sprintf("%v.%v", path, key)
	} else {
		path = key
	}

	if z, ok := m[key]; ok {
		if zm, ok := z.(map[string]interface{}); ok {
			if vm, ok := val.(map[string]interface{}); ok {
				for k, v := range vm {
					zm = s.mergeMap(zm, k, path, v)
				}
				m[key] = zm
				s.entries[path] = zm
			} else if vm, ok := val.(map[interface{}]interface{}); ok {
				for k, v := range vm {
					kk, ok := k.(string)
					if !ok {
						kk = fmt.Sprintf("%v", k)
					}
					zm = s.mergeMap(zm, kk, path, v)
				}
				m[key] = zm
				s.entries[path] = zm
			} else {
				m[key] = val
				s.entries[path] = val
			}
		} else {
			m[key] = val
			s.entries[path] = val
		}
	} else {
		m[key] = val
		s.entries[path] = val
	}
	return m
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

// ResetOptions to reset the exists `Options`, so that you could follow a `LoadConfigFile()` with it.
func ResetOptions() {
	uniqueWorker.rxxtOptions.Reset()
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
	return uniqueWorker.rxxtOptions.DumpAsString()
}

// AsYaml returns a yaml string bytes about all options
func AsYaml() (b []byte) {
	obj := uniqueWorker.rxxtOptions.GetHierarchyList()
	b, _ = yaml.Marshal(obj)
	return
}

// SaveAsYaml to Save all config entries as a yaml file
func SaveAsYaml(filename string) (err error) {
	b := AsYaml()
	err = ioutil.WriteFile(filename, b, 0644)
	return
}

// AsJSON returns a json string bytes about all options
func AsJSON() (b []byte) {
	obj := uniqueWorker.rxxtOptions.GetHierarchyList()
	b, _ = json.Marshal(obj)
	return
}

// SaveAsJSON to Save all config entries as a json file
func SaveAsJSON(filename string) (err error) {
	b := AsJSON()
	err = ioutil.WriteFile(filename, b, 0644)
	return
}

// AsToml returns a toml string bytes about all options
func AsToml() (b []byte) {
	obj := uniqueWorker.rxxtOptions.GetHierarchyList()
	buf := bytes.NewBuffer([]byte{})
	e := toml.NewEncoder(buf)
	if err := e.Encode(obj); err == nil {
		b = buf.Bytes()
	}
	return
}

// SaveAsToml to Save all config entries as a toml file
func SaveAsToml(filename string) (err error) {
	obj := uniqueWorker.rxxtOptions.GetHierarchyList()
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

// GetHierarchyList returns the hierarchy data
func GetHierarchyList() map[string]interface{} {
	return uniqueWorker.rxxtOptions.GetHierarchyList()
}

// GetHierarchyList returns the hierarchy data for dumping
func (s *Options) GetHierarchyList() map[string]interface{} {
	defer s.rw.RUnlock()
	s.rw.RLock()
	return s.hierarchy
}
