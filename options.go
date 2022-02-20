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
	"gopkg.in/hedzr/errors.v3"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

//
//
//

// HasKey detects whether a key exists in cmdr options store or not
func HasKey(key string) (ok bool) {
	return currentOptions().Has(key)
}

// Get returns the generic value of an `Option` key with [WrapWithRxxtPrefix]. Such as:
// ```golang
// cmdr.Get("app.logger.level") => 'DEBUG',...
// ```
//
func Get(key string) interface{} {
	return currentOptions().Get(key)
}

// GetR returns the generic value of an `Option` key with [WrapWithRxxtPrefix]. Such as:
// ```golang
// cmdr.GetR("logger.level") => 'DEBUG',...
// ```
//
func GetR(key string) interface{} {
	w := internalGetWorker()
	return w.rxxtOptions.Get(w.wrapWithRxxtPrefix(key))
}

// GetBool returns the bool value of an `Option` key. Such as:
// ```golang
// cmdr.GetBool("app.logger.enable", false) => true,...
// ```
//
func GetBool(key string, defaultVal ...bool) bool {
	return currentOptions().GetBoolEx(key, defaultVal...)
}

// GetBoolP returns the bool value of an `Option` key. Such as:
// ```golang
// cmdr.GetBoolP("app.logger", "enable", false) => true,...
// ```
func GetBoolP(prefix, key string, defaultVal ...bool) bool {
	return currentOptions().GetBoolEx(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetBoolR returns the bool value of an `Option` key with [WrapWithRxxtPrefix]. Such as:
// ```golang
// cmdr.GetBoolR("logger.enable", false) => true,...
// ```
//
func GetBoolR(key string, defaultVal ...bool) bool {
	w := internalGetWorker()
	return w.rxxtOptions.GetBoolEx(w.wrapWithRxxtPrefix(key), defaultVal...)
}

// GetBoolRP returns the bool value of an `Option` key with [WrapWithRxxtPrefix]. Such as:
// ```golang
// cmdr.GetBoolRP("logger", "enable", false) => true,...
// ```
func GetBoolRP(prefix, key string, defaultVal ...bool) bool {
	w := internalGetWorker()
	return w.rxxtOptions.GetBoolEx(w.wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// GetInt returns the int value of an `Option` key.
func GetInt(key string, defaultVal ...int) int {
	return currentOptions().GetIntEx(key, defaultVal...)
}

// GetIntP returns the int value of an `Option` key.
func GetIntP(prefix, key string, defaultVal ...int) int {
	return currentOptions().GetIntEx(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetIntR returns the int value of an `Option` key with [WrapWithRxxtPrefix].
func GetIntR(key string, defaultVal ...int) int {
	w := internalGetWorker()
	return w.rxxtOptions.GetIntEx(w.wrapWithRxxtPrefix(key), defaultVal...)
}

// GetIntRP returns the int value of an `Option` key with [WrapWithRxxtPrefix].
func GetIntRP(prefix, key string, defaultVal ...int) int {
	w := internalGetWorker()
	return w.rxxtOptions.GetIntEx(w.wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// GetInt64 returns the int64 value of an `Option` key.
func GetInt64(key string, defaultVal ...int64) int64 {
	return currentOptions().GetInt64Ex(key, defaultVal...)
}

// GetInt64P returns the int64 value of an `Option` key.
func GetInt64P(prefix, key string, defaultVal ...int64) int64 {
	return currentOptions().GetInt64Ex(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetInt64R returns the int64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetInt64R(key string, defaultVal ...int64) int64 {
	w := internalGetWorker()
	return w.rxxtOptions.GetInt64Ex(w.wrapWithRxxtPrefix(key), defaultVal...)
}

// GetInt64RP returns the int64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetInt64RP(prefix, key string, defaultVal ...int64) int64 {
	w := internalGetWorker()
	return w.rxxtOptions.GetInt64Ex(w.wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// GetUint returns the uint value of an `Option` key.
func GetUint(key string, defaultVal ...uint) uint {
	return currentOptions().GetUintEx(key, defaultVal...)
}

// GetUintP returns the uint value of an `Option` key.
func GetUintP(prefix, key string, defaultVal ...uint) uint {
	return currentOptions().GetUintEx(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetUintR returns the uint value of an `Option` key with [WrapWithRxxtPrefix].
func GetUintR(key string, defaultVal ...uint) uint {
	w := internalGetWorker()
	return w.rxxtOptions.GetUintEx(w.wrapWithRxxtPrefix(key), defaultVal...)
}

// GetUintRP returns the uint value of an `Option` key with [WrapWithRxxtPrefix].
func GetUintRP(prefix, key string, defaultVal ...uint) uint {
	w := internalGetWorker()
	return w.rxxtOptions.GetUintEx(w.wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// GetUint64 returns the uint64 value of an `Option` key.
func GetUint64(key string, defaultVal ...uint64) uint64 {
	return currentOptions().GetUint64Ex(key, defaultVal...)
}

// GetUint64P returns the uint64 value of an `Option` key.
func GetUint64P(prefix, key string, defaultVal ...uint64) uint64 {
	return currentOptions().GetUint64Ex(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetUint64R returns the uint64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetUint64R(key string, defaultVal ...uint64) uint64 {
	w := internalGetWorker()
	return w.rxxtOptions.GetUint64Ex(w.wrapWithRxxtPrefix(key), defaultVal...)
}

// GetUint64RP returns the uint64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetUint64RP(prefix, key string, defaultVal ...uint64) uint64 {
	w := internalGetWorker()
	return w.rxxtOptions.GetUint64Ex(w.wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// GetKibibytes returns the uint64 value of an `Option` key.
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
func GetKibibytes(key string, defaultVal ...uint64) uint64 {
	return currentOptions().GetKibibytesEx(key, defaultVal...)
}

// GetKibibytesP returns the uint64 value of an `Option` key.
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
func GetKibibytesP(prefix, key string, defaultVal ...uint64) uint64 {
	return currentOptions().GetKibibytesEx(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetKibibytesR returns the uint64 value of an `Option` key with [WrapWithRxxtPrefix].
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
func GetKibibytesR(key string, defaultVal ...uint64) uint64 {
	w := internalGetWorker()
	return w.rxxtOptions.GetKibibytesEx(w.wrapWithRxxtPrefix(key), defaultVal...)
}

// GetKibibytesRP returns the uint64 value of an `Option` key with [WrapWithRxxtPrefix].
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
func GetKibibytesRP(prefix, key string, defaultVal ...uint64) uint64 {
	w := internalGetWorker()
	return w.rxxtOptions.GetKibibytesEx(w.wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// GetKilobytes returns the uint64 value of an `Option` key.
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
func GetKilobytes(key string, defaultVal ...uint64) uint64 {
	return currentOptions().GetKilobytesEx(key, defaultVal...)
}

// GetKilobytesP returns the uint64 value of an `Option` key.
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
func GetKilobytesP(prefix, key string, defaultVal ...uint64) uint64 {
	return currentOptions().GetKilobytesEx(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetKilobytesR returns the uint64 value of an `Option` key with [WrapWithRxxtPrefix].
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
func GetKilobytesR(key string, defaultVal ...uint64) uint64 {
	w := internalGetWorker()
	return w.rxxtOptions.GetKilobytesEx(w.wrapWithRxxtPrefix(key), defaultVal...)
}

// GetKilobytesRP returns the uint64 value of an `Option` key with [WrapWithRxxtPrefix].
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
func GetKilobytesRP(prefix, key string, defaultVal ...uint64) uint64 {
	w := internalGetWorker()
	return w.rxxtOptions.GetKilobytesEx(w.wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// GetFloat32 returns the float32 value of an `Option` key.
func GetFloat32(key string, defaultVal ...float32) float32 {
	return currentOptions().GetFloat32Ex(key, defaultVal...)
}

// GetFloat32P returns the float32 value of an `Option` key.
func GetFloat32P(prefix, key string, defaultVal ...float32) float32 {
	return currentOptions().GetFloat32Ex(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetFloat32R returns the float32 value of an `Option` key with [WrapWithRxxtPrefix].
func GetFloat32R(key string, defaultVal ...float32) float32 {
	w := internalGetWorker()
	return w.rxxtOptions.GetFloat32Ex(w.wrapWithRxxtPrefix(key), defaultVal...)
}

// GetFloat32RP returns the float32 value of an `Option` key with [WrapWithRxxtPrefix].
func GetFloat32RP(prefix, key string, defaultVal ...float32) float32 {
	w := internalGetWorker()
	return w.rxxtOptions.GetFloat32Ex(w.wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// GetFloat64 returns the float64 value of an `Option` key.
func GetFloat64(key string, defaultVal ...float64) float64 {
	return currentOptions().GetFloat64Ex(key, defaultVal...)
}

// GetFloat64P returns the float64 value of an `Option` key.
func GetFloat64P(prefix, key string, defaultVal ...float64) float64 {
	return currentOptions().GetFloat64Ex(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetFloat64R returns the float64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetFloat64R(key string, defaultVal ...float64) float64 {
	w := internalGetWorker()
	return w.rxxtOptions.GetFloat64Ex(w.wrapWithRxxtPrefix(key), defaultVal...)
}

// GetFloat64RP returns the float64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetFloat64RP(prefix, key string, defaultVal ...float64) float64 {
	w := internalGetWorker()
	return w.rxxtOptions.GetFloat64Ex(w.wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// GetComplex64 returns the complex64 value of an `Option` key.
func GetComplex64(key string, defaultVal ...complex64) complex64 {
	return currentOptions().GetComplex64(key, defaultVal...)
}

// GetComplex64P returns the complex64 value of an `Option` key.
func GetComplex64P(prefix, key string, defaultVal ...complex64) complex64 {
	return currentOptions().GetComplex64(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetComplex64R returns the complex64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetComplex64R(key string, defaultVal ...complex64) complex64 {
	w := internalGetWorker()
	return w.rxxtOptions.GetComplex64(w.wrapWithRxxtPrefix(key), defaultVal...)
}

// GetComplex64RP returns the complex64 value of an `Option` key with [WrapWithRxxtPrefix].
func GetComplex64RP(prefix, key string, defaultVal ...complex64) complex64 {
	w := internalGetWorker()
	return w.rxxtOptions.GetComplex64(w.wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// GetComplex128 returns the complex128 value of an `Option` key.
func GetComplex128(key string, defaultVal ...complex128) complex128 {
	return currentOptions().GetComplex128(key, defaultVal...)
}

// GetComplex128P returns the complex128 value of an `Option` key.
func GetComplex128P(prefix, key string, defaultVal ...complex128) complex128 {
	return currentOptions().GetComplex128(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetComplex128R returns the complex128 value of an `Option` key with [WrapWithRxxtPrefix].
func GetComplex128R(key string, defaultVal ...complex128) complex128 {
	w := internalGetWorker()
	return w.rxxtOptions.GetComplex128(w.wrapWithRxxtPrefix(key), defaultVal...)
}

// GetComplex128RP returns the complex128 value of an `Option` key with [WrapWithRxxtPrefix].
func GetComplex128RP(prefix, key string, defaultVal ...complex128) complex128 {
	w := internalGetWorker()
	return w.rxxtOptions.GetComplex128(w.wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// GetString returns the string value of an `Option` key.
func GetString(key string, defaultVal ...string) string {
	return currentOptions().GetString(key, defaultVal...)
}

// GetStringP returns the string value of an `Option` key.
func GetStringP(prefix, key string, defaultVal ...string) string {
	return currentOptions().GetString(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetStringR returns the string value of an `Option` key with [WrapWithRxxtPrefix].
func GetStringR(key string, defaultVal ...string) string {
	if noResetWorker {
		return uniqueWorker.rxxtOptions.GetString(wrapWithRxxtPrefix(key), defaultVal...)
	}
	w := internalGetWorker()
	return w.rxxtOptions.GetString(w.wrapWithRxxtPrefix(key), defaultVal...)
}

// GetStringRP returns the string value of an `Option` key with [WrapWithRxxtPrefix].
func GetStringRP(prefix, key string, defaultVal ...string) string {
	w := internalGetWorker()
	return w.rxxtOptions.GetString(w.wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// GetStringNoExpand returns the string value of an `Option` key.
func GetStringNoExpand(key string, defaultVal ...string) string {
	return currentOptions().GetStringNoExpand(key, defaultVal...)
}

// GetStringNoExpandP returns the string value of an `Option` key.
func GetStringNoExpandP(prefix, key string, defaultVal ...string) string {
	return currentOptions().GetStringNoExpand(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetStringNoExpandR returns the string value of an `Option` key with [WrapWithRxxtPrefix].
func GetStringNoExpandR(key string, defaultVal ...string) string {
	w := internalGetWorker()
	return w.rxxtOptions.GetStringNoExpand(w.wrapWithRxxtPrefix(key), defaultVal...)
}

// GetStringNoExpandRP returns the string value of an `Option` key with [WrapWithRxxtPrefix].
func GetStringNoExpandRP(prefix, key string, defaultVal ...string) string {
	w := internalGetWorker()
	return w.rxxtOptions.GetStringNoExpand(w.wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// GetStringSlice returns the string slice value of an `Option` key.
func GetStringSlice(key string, defaultVal ...string) []string {
	return currentOptions().GetStringSlice(key, defaultVal...)
}

// GetStringSliceP returns the string slice value of an `Option` key.
func GetStringSliceP(prefix, key string, defaultVal ...string) []string {
	return currentOptions().GetStringSlice(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetStringSliceR returns the string slice value of an `Option` key with [WrapWithRxxtPrefix].
func GetStringSliceR(key string, defaultVal ...string) []string {
	w := internalGetWorker()
	return w.rxxtOptions.GetStringSlice(w.wrapWithRxxtPrefix(key), defaultVal...)
}

// GetStringSliceRP returns the string slice value of an `Option` key with [WrapWithRxxtPrefix].
func GetStringSliceRP(prefix, key string, defaultVal ...string) []string {
	w := internalGetWorker()
	return w.rxxtOptions.GetStringSlice(w.wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// GetMap an `Option` by key string, it returns a hierarchy map or nil
func GetMap(key string) map[string]interface{} {
	return currentOptions().GetMap(key)
}

// GetMapR an `Option` by key string with [WrapWithRxxtPrefix], it returns a hierarchy map or nil
func GetMapR(key string) map[string]interface{} {
	w := internalGetWorker()
	return w.rxxtOptions.GetMap(w.wrapWithRxxtPrefix(key))
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
	w := internalGetWorker()
	key := w.wrapWithRxxtPrefix(sectionKeyPath)
	fObj := w.rxxtOptions.GetMap(key) // GetMapR(sectionKeyPath)
	if fObj != nil {
		defer handleSerializeError(&err)
		var b []byte
		b, err = yaml.Marshal(fObj)
		if err == nil {
			err = yaml.Unmarshal(b, holder)
			// if err == nil {
			// 	logrus.Debugf("configuration section got: %v", configHolder)
			// }
		}
	} else {
		fObj := w.rxxtOptions.Get(key) // GetR(sectionKeyPath)
		if fObj != nil {
			defer handleSerializeError(&err)
			var b []byte
			b, err = yaml.Marshal(fObj)
			if err == nil {
				err = yaml.Unmarshal(b, holder)
				// if err == nil {
				// 	logrus.Debugf("configuration section got: %v", configHolder)
				// }
			}
		}
	}
	return
}

// GetIntSlice returns the int slice value of an `Option` key.
func GetIntSlice(key string, defaultVal ...int) []int {
	return currentOptions().GetIntSlice(key, defaultVal...)
}

// GetIntSliceP returns the int slice value of an `Option` key.
func GetIntSliceP(prefix, key string, defaultVal ...int) []int {
	return currentOptions().GetIntSlice(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetIntSliceR returns the int slice value of an `Option` key with [WrapWithRxxtPrefix].
func GetIntSliceR(key string, defaultVal ...int) []int {
	w := internalGetWorker()
	return w.rxxtOptions.GetIntSlice(w.wrapWithRxxtPrefix(key), defaultVal...)
}

// GetIntSliceRP returns the int slice value of an `Option` key with [WrapWithRxxtPrefix].
func GetIntSliceRP(prefix, key string, defaultVal ...int) []int {
	w := internalGetWorker()
	return w.rxxtOptions.GetIntSlice(w.wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// // GetUintSlice returns the int slice value of an `Option` key.
// func GetUintSlice(key string, defaultVal ...uint) []uint {
// 	return currentOptions().GetUintSlice(key, defaultVal...)
// }
//
// // GetUintSliceP returns the int slice value of an `Option` key.
// func GetUintSliceP(prefix, key string, defaultVal ...uint) []uint {
// 	return currentOptions().GetUintSlice(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
// }
//
// // GetUintSliceR returns the int slice value of an `Option` key with [WrapWithRxxtPrefix].
// func GetUintSliceR(key string, defaultVal ...uint) []uint {
// 	return currentOptions().GetUintSlice(wrapWithRxxtPrefix(key), defaultVal...)
// }
//
// // GetUintSliceRP returns the int slice value of an `Option` key with [WrapWithRxxtPrefix].
// func GetUintSliceRP(prefix, key string, defaultVal ...uint) []uint {
// 	return currentOptions().GetUintSlice(wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
// }

// GetInt64Slice returns the int slice value of an `Option` key.
func GetInt64Slice(key string, defaultVal ...int64) []int64 {
	return currentOptions().GetInt64Slice(key, defaultVal...)
}

// GetInt64SliceP returns the int slice value of an `Option` key.
func GetInt64SliceP(prefix, key string, defaultVal ...int64) []int64 {
	return currentOptions().GetInt64Slice(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetInt64SliceR returns the int slice value of an `Option` key with [WrapWithRxxtPrefix].
func GetInt64SliceR(key string, defaultVal ...int64) []int64 {
	w := internalGetWorker()
	return w.rxxtOptions.GetInt64Slice(w.wrapWithRxxtPrefix(key), defaultVal...)
}

// GetInt64SliceRP returns the int slice value of an `Option` key with [WrapWithRxxtPrefix].
func GetInt64SliceRP(prefix, key string, defaultVal ...int64) []int64 {
	w := internalGetWorker()
	return w.rxxtOptions.GetInt64Slice(w.wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// GetUint64Slice returns the int slice value of an `Option` key.
func GetUint64Slice(key string, defaultVal ...uint64) []uint64 {
	return currentOptions().GetUint64Slice(key, defaultVal...)
}

// GetUint64SliceP returns the int slice value of an `Option` key.
func GetUint64SliceP(prefix, key string, defaultVal ...uint64) []uint64 {
	return currentOptions().GetUint64Slice(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetUint64SliceR returns the int slice value of an `Option` key with [WrapWithRxxtPrefix].
func GetUint64SliceR(key string, defaultVal ...uint64) []uint64 {
	w := internalGetWorker()
	return w.rxxtOptions.GetUint64Slice(w.wrapWithRxxtPrefix(key), defaultVal...)
}

// GetUint64SliceRP returns the int slice value of an `Option` key with [WrapWithRxxtPrefix].
func GetUint64SliceRP(prefix, key string, defaultVal ...uint64) []uint64 {
	w := internalGetWorker()
	return w.rxxtOptions.GetUint64Slice(w.wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
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
	return currentOptions().GetDuration(key, defaultVal...)
}

// GetDurationP returns the int slice value of an `Option` key.
func GetDurationP(prefix, key string, defaultVal ...time.Duration) time.Duration {
	return currentOptions().GetDuration(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetDurationR returns the int slice value of an `Option` key.
func GetDurationR(key string, defaultVal ...time.Duration) time.Duration {
	w := internalGetWorker()
	return w.rxxtOptions.GetDuration(w.wrapWithRxxtPrefix(key), defaultVal...)
}

// GetDurationRP returns the int slice value of an `Option` key.
func GetDurationRP(prefix, key string, defaultVal ...time.Duration) time.Duration {
	w := internalGetWorker()
	return w.rxxtOptions.GetDuration(w.wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// WrapWithRxxtPrefix wrap an key with [RxxtPrefix], for [GetXxx(key)] and [GetXxxP(prefix,key)]
func WrapWithRxxtPrefix(key string) string {
	return wrapWithRxxtPrefix(key)
}

func wrapWithRxxtPrefix(key string) string {
	w := internalGetWorker()
	return w.wrapWithRxxtPrefix(key)
}

func (w *ExecWorker) wrapWithRxxtPrefix(key string) string {
	if len(w.rxxtPrefixes) == 0 {
		return key
	}
	p := w.getPrefix() // strings.Join(RxxtPrefix, ".")
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
	currentOptions().Set(key, val)
}

// SetOverwrite sets the value of an `Option` key. The key MUST not have an `app` prefix.
// It replaces the old value on a slice, instead of the default append mode.
//
// SetOverwrite(key, newSlice) will replace the original value with newSlice.
func SetOverwrite(key string, val interface{}) {
	currentOptions().SetOverwrite(key, val)
}

// SetNx but without prefix auto-wrapped.
// `rxxtPrefix` is a string slice to define the prefix string array, default is ["app"].
// So, cmdr.Set("debug", true) will put a real entry with (`debug`, true).
func SetNx(key string, val interface{}) {
	currentOptions().SetNx(key, val)
}

// SetNxOverwrite likes SetOverwrite but without prefix auto-wrapped.
// It replaces the old value on a slice, instead of the default append mode.
func SetNxOverwrite(key string, val interface{}) {
	currentOptions().SetNxOverwrite(key, val)
}

// SetRaw but without prefix auto-wrapped.
// So, cmdr.SetRaw("debug", true) will put a real entry with (`debug`, true).
//
//    cmdr.Set("debug", true)
//    cmdr.GetBool("debug") => true
//
func SetRaw(key string, val interface{}) {
	currentOptions().SetRaw(key, val)
}

// MergeWith will merge a map into Option Store recursively.
// You could merge a yaml/json/toml options into cmdr Hierarchy Options.
func MergeWith(m map[string]interface{}) (err error) {
	err = currentOptions().MergeWith(m)
	return
}

// ResetOptions to reset the current `Options Store`, so that you could follow a `LoadConfigFile()` with it.
func ResetOptions() {
	currentOptions().Reset()
}

// DumpAsString for debugging.
func DumpAsString() (str string) {
	return currentOptions().DumpAsString(false)
}

// AsYaml returns a yaml string bytes about all options
func AsYaml() (b []byte) {
	b, _ = AsYamlExt()
	return
}

// AsYamlExt returns a yaml string bytes about all options
func AsYamlExt() (b []byte, err error) {
	obj := currentOptions().GetHierarchyList()
	defer handleSerializeError(&err)

	var sb strings.Builder
	e := yaml.NewEncoder(&sb)
	e.SetIndent(2)
	err = e.Encode(obj)
	if err != nil {
		return nil, err
	}
	err = e.Close()
	if err != nil {
		return nil, err
	}

	b = []byte(sb.String())
	// b, err = yaml.Marshal(obj)
	return
}

// SaveAsYaml to Save all config entries as a yaml file
func SaveAsYaml(filename string) (err error) {
	var b []byte
	b, err = AsYamlExt()
	if err == nil {
		err = ioutil.WriteFile(filename, b, 0644)
	}
	return
}

// AsJSON returns a json string bytes about all options
func AsJSON() (b []byte) {
	b, _ = AsJSONExt(false)
	return
}

// AsJSONExt returns a json string bytes about all options
func AsJSONExt(prettyFormat bool) (b []byte, err error) {
	obj := currentOptions().GetHierarchyList()
	defer handleSerializeError(&err)
	if prettyFormat {
		b, err = json.MarshalIndent(obj, "", "  ")
	} else {
		b, err = json.Marshal(obj)
	}
	return
}

// SaveAsJSON to Save all config entries as a json file
func SaveAsJSON(filename string) (err error) {
	b := AsJSON()
	err = ioutil.WriteFile(filename, b, 0644)
	return
}

// SaveAsJSONExt to Save all config entries as a json file
func SaveAsJSONExt(filename string, prettyFormat bool) (err error) {
	var b []byte
	b, err = AsJSONExt(prettyFormat)
	if err == nil {
		err = ioutil.WriteFile(filename, b, 0644)
	}
	return
}

// AsToml returns a toml string bytes about all options
func AsToml() (b []byte) {
	b, _ = AsTomlExt()
	return
}

// AsTomlExt returns a toml string bytes about all options
func AsTomlExt() (b []byte, err error) {
	obj := currentOptions().GetHierarchyList()
	buf := bytes.NewBuffer([]byte{})
	e := toml.NewEncoder(buf)
	defer handleSerializeError(&err)
	if err = e.Encode(obj); err == nil {
		b = buf.Bytes()
	}
	return
}

// SaveAsToml to Save all config entries as a toml file
func SaveAsToml(filename string) (err error) {
	obj := currentOptions().GetHierarchyList()
	err = SaveObjAsToml(obj, filename)
	return
}

// SaveObjAsToml to Save an object as a toml file
func SaveObjAsToml(obj interface{}, filename string) (err error) {
	var f *os.File
	f, err = os.Create(filename)
	if err == nil {

		defer handleSerializeError(&err)
		e := toml.NewEncoder(bufio.NewWriter(f))
		err = e.Encode(obj)

		// err = ioutil.WriteFile(filename, b, 0644)
	}
	return
}

// GetHierarchyList returns the hierarchy data
func GetHierarchyList() map[string]interface{} {
	return currentOptions().GetHierarchyList()
}

func handleSerializeError(err *error) {
	if v := recover(); v != nil {
		if e, ok := v.(error); ok {
			*err = e
		} else if s, ok := v.(string); ok {
			*err = errors.New(s)
			// if s == "cannot marshal type: complex128" {
			// 	err = errors.New(s)
			// }
		} else {
			panic(v)
		}
	}
}
