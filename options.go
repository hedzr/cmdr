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
	"time"
)

//
//
//

// HasKey detects whether a key exists in cmdr options store or not
func HasKey(key string) (ok bool) {
	return uniqueWorker.rxxtOptions.Has(key)
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

// GetBool returns the bool value of an `Option` key. Such as:
// ```golang
// cmdr.GetBool("app.logger.enable", false) => true,...
// ```
//
func GetBool(key string, defaultVal ...bool) bool {
	return uniqueWorker.rxxtOptions.GetBoolEx(key, defaultVal...)
}

// GetBoolP returns the bool value of an `Option` key. Such as:
// ```golang
// cmdr.GetBoolP("app.logger", "enable", false) => true,...
// ```
func GetBoolP(prefix, key string, defaultVal ...bool) bool {
	return uniqueWorker.rxxtOptions.GetBoolEx(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetBoolR returns the bool value of an `Option` key with [WrapWithRxxtPrefix]. Such as:
// ```golang
// cmdr.GetBoolR("logger.enable", false) => true,...
// ```
//
func GetBoolR(key string, defaultVal ...bool) bool {
	return uniqueWorker.rxxtOptions.GetBoolEx(wrapWithRxxtPrefix(key), defaultVal...)
}

// GetBoolRP returns the bool value of an `Option` key with [WrapWithRxxtPrefix]. Such as:
// ```golang
// cmdr.GetBoolRP("logger", "enable", false) => true,...
// ```
func GetBoolRP(prefix, key string, defaultVal ...bool) bool {
	return uniqueWorker.rxxtOptions.GetBoolEx(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// GetInt returns the int value of an `Option` key.
func GetInt(key string, defaultVal ...int) int {
	return uniqueWorker.rxxtOptions.GetIntEx(key, defaultVal...)
}

// GetIntP returns the int value of an `Option` key.
func GetIntP(prefix, key string, defaultVal ...int) int {
	return uniqueWorker.rxxtOptions.GetIntEx(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetIntR returns the int value of an `Option` key with [WrapWithRxxtPrefix].
func GetIntR(key string, defaultVal ...int) int {
	return uniqueWorker.rxxtOptions.GetIntEx(wrapWithRxxtPrefix(key), defaultVal...)
}

// GetIntRP returns the int value of an `Option` key with [WrapWithRxxtPrefix].
func GetIntRP(prefix, key string, defaultVal ...int) int {
	return uniqueWorker.rxxtOptions.GetIntEx(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// GetInt64 returns the int64 value of an `Option` key.
func GetInt64(key string, defaultVal ...int64) int64 {
	return uniqueWorker.rxxtOptions.GetInt64Ex(key, defaultVal...)
}

// GetInt64P returns the int64 value of an `Option` key.
func GetInt64P(prefix, key string, defaultVal ...int64) int64 {
	return uniqueWorker.rxxtOptions.GetInt64Ex(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetInt64R returns the int64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetInt64R(key string, defaultVal ...int64) int64 {
	return uniqueWorker.rxxtOptions.GetInt64Ex(wrapWithRxxtPrefix(key), defaultVal...)
}

// GetInt64RP returns the int64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetInt64RP(prefix, key string, defaultVal ...int64) int64 {
	return uniqueWorker.rxxtOptions.GetInt64Ex(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// GetUint returns the uint value of an `Option` key.
func GetUint(key string, defaultVal ...uint) uint {
	return uniqueWorker.rxxtOptions.GetUintEx(key, defaultVal...)
}

// GetUintP returns the uint value of an `Option` key.
func GetUintP(prefix, key string, defaultVal ...uint) uint {
	return uniqueWorker.rxxtOptions.GetUintEx(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetUintR returns the uint value of an `Option` key with [WrapWithRxxtPrefix].
func GetUintR(key string, defaultVal ...uint) uint {
	return uniqueWorker.rxxtOptions.GetUintEx(wrapWithRxxtPrefix(key), defaultVal...)
}

// GetUintRP returns the uint value of an `Option` key with [WrapWithRxxtPrefix].
func GetUintRP(prefix, key string, defaultVal ...uint) uint {
	return uniqueWorker.rxxtOptions.GetUintEx(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// GetUint64 returns the uint64 value of an `Option` key.
func GetUint64(key string, defaultVal ...uint64) uint64 {
	return uniqueWorker.rxxtOptions.GetUint64Ex(key, defaultVal...)
}

// GetUint64P returns the uint64 value of an `Option` key.
func GetUint64P(prefix, key string, defaultVal ...uint64) uint64 {
	return uniqueWorker.rxxtOptions.GetUint64Ex(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetUint64R returns the uint64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetUint64R(key string, defaultVal ...uint64) uint64 {
	return uniqueWorker.rxxtOptions.GetUint64Ex(wrapWithRxxtPrefix(key), defaultVal...)
}

// GetUint64RP returns the uint64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetUint64RP(prefix, key string, defaultVal ...uint64) uint64 {
	return uniqueWorker.rxxtOptions.GetUint64Ex(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// GetFloat32 returns the float32 value of an `Option` key.
func GetFloat32(key string, defaultVal ...float32) float32 {
	return float32(uniqueWorker.rxxtOptions.GetFloat32Ex(key, defaultVal...))
}

// GetFloat32P returns the float32 value of an `Option` key.
func GetFloat32P(prefix, key string, defaultVal ...float32) float32 {
	return float32(uniqueWorker.rxxtOptions.GetFloat32Ex(fmt.Sprintf("%s.%s", prefix, key), defaultVal...))
}

// GetFloat32R returns the float32 value of an `Option` key with [WrapWithRxxtPrefix].
func GetFloat32R(key string, defaultVal ...float32) float32 {
	return float32(uniqueWorker.rxxtOptions.GetFloat32Ex(wrapWithRxxtPrefix(key), defaultVal...))
}

// GetFloat32RP returns the float32 value of an `Option` key with [WrapWithRxxtPrefix].
func GetFloat32RP(prefix, key string, defaultVal ...float32) float32 {
	return float32(uniqueWorker.rxxtOptions.GetFloat32Ex(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...))
}

// GetFloat64 returns the float64 value of an `Option` key.
func GetFloat64(key string, defaultVal ...float64) float64 {
	return uniqueWorker.rxxtOptions.GetFloat64Ex(key, defaultVal...)
}

// GetFloat64P returns the float64 value of an `Option` key.
func GetFloat64P(prefix, key string, defaultVal ...float64) float64 {
	return uniqueWorker.rxxtOptions.GetFloat64Ex(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetFloat64R returns the float64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetFloat64R(key string, defaultVal ...float64) float64 {
	return uniqueWorker.rxxtOptions.GetFloat64Ex(wrapWithRxxtPrefix(key), defaultVal...)
}

// GetFloat64RP returns the float64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetFloat64RP(prefix, key string, defaultVal ...float64) float64 {
	return uniqueWorker.rxxtOptions.GetFloat64Ex(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// GetString returns the string value of an `Option` key.
func GetString(key string, defaultVal ...string) string {
	return uniqueWorker.rxxtOptions.GetString(key, defaultVal...)
}

// GetStringP returns the string value of an `Option` key.
func GetStringP(prefix, key string, defaultVal ...string) string {
	return uniqueWorker.rxxtOptions.GetString(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetStringR returns the string value of an `Option` key with [WrapWithRxxtPrefix].
func GetStringR(key string, defaultVal ...string) string {
	return uniqueWorker.rxxtOptions.GetString(wrapWithRxxtPrefix(key), defaultVal...)
}

// GetStringRP returns the string value of an `Option` key with [WrapWithRxxtPrefix].
func GetStringRP(prefix, key string, defaultVal ...string) string {
	return uniqueWorker.rxxtOptions.GetString(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
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

// // GetDuration returns the int slice value of an `Option` key.
// func GetDuration(key string) time.Duration {
// 	return uniqueWorker.rxxtOptions.GetDuration(key, 0)
// }
//
// // GetDurationP returns the int slice value of an `Option` key.
// func GetDurationP(prefix, key string) time.Duration {
// 	return uniqueWorker.rxxtOptions.GetDuration(fmt.Sprintf("%s.%s", prefix, key), 0)
// }
//
// // GetDurationR returns the int slice value of an `Option` key.
// func GetDurationR(key string) time.Duration {
// 	return uniqueWorker.rxxtOptions.GetDuration(wrapWithRxxtPrefix(key), 0)
// }
//
// // GetDurationRP returns the int slice value of an `Option` key.
// func GetDurationRP(prefix, key string) time.Duration {
// 	return uniqueWorker.rxxtOptions.GetDuration(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), 0)
// }

// GetDuration returns the int slice value of an `Option` key.
func GetDuration(key string, defaultVal ...time.Duration) time.Duration {
	return uniqueWorker.rxxtOptions.GetDuration(key, defaultVal...)
}

// GetDurationP returns the int slice value of an `Option` key.
func GetDurationP(prefix, key string, defaultVal ...time.Duration) time.Duration {
	return uniqueWorker.rxxtOptions.GetDuration(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetDurationR returns the int slice value of an `Option` key.
func GetDurationR(key string, defaultVal ...time.Duration) time.Duration {
	return uniqueWorker.rxxtOptions.GetDuration(wrapWithRxxtPrefix(key), defaultVal...)
}

// GetDurationRP returns the int slice value of an `Option` key.
func GetDurationRP(prefix, key string, defaultVal ...time.Duration) time.Duration {
	return uniqueWorker.rxxtOptions.GetDuration(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
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

// MergeWith will merge a map recursive.
// You could merge a yaml/json/toml options into cmdr Hierarchy Options.
func MergeWith(m map[string]interface{}) (err error) {
	err = uniqueWorker.rxxtOptions.MergeWith(m)
	return
}

// ResetOptions to reset the exists `Options`, so that you could follow a `LoadConfigFile()` with it.
func ResetOptions() {
	uniqueWorker.rxxtOptions.Reset()
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

// GetHierarchyList returns the hierarchy data
func GetHierarchyList() map[string]interface{} {
	return uniqueWorker.rxxtOptions.GetHierarchyList()
}
