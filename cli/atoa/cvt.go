package atoa

import (
	"math"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/hedzr/cmdr/v2/internal/tool"
	"github.com/hedzr/cmdr/v2/pkg/times"
)

type Converter func(str string, targetType reflect.Type) (ret any, err error)

var defaultConverters map[reflect.Type]Converter

var oncedefcvts sync.Once

func defcvts() map[reflect.Type]Converter {
	oncedefcvts.Do(func() {
		defaultConverters = map[reflect.Type]Converter{
			reflect.TypeOf((*time.Time)(nil)).Elem(): toTimeTime,
			reflect.TypeOf((*time.Time)(nil)):        toTimeTimePtr,
		}
	})
	return defaultConverters
}

// RoundedNumber is an Opt for FromString.
//
// For example:
//
//	v = FromString(ctx, "8.97", int32(9), WithFeatures(RoundedNumber))
//	assert.Equal(t, v, int32(9))
//	v = FromString(ctx, "8.97", int32(9))
//	assert.Equal(t, v, int32(8))
func RoundedNumber(str string, targetType reflect.Type) (ret any, err error) {
	parseFloat := func(str string, kind reflect.Kind) (v any, err error) {
		// return strconv.ParseFloat(str, bits)
		var vi float64
		// vi, err = strconv.ParseFloat(str, rt.Bits())
		vi, err = tool.N[float64](str)
		if err == nil {
			f := math.Floor(vi + 0.5)
			v = fltByKind(f, kind)
		}
		return
	}
	parseInt := func(str string, kind reflect.Kind) (v any, err error) {
		if strings.Contains(str, ".") {
			f, err := parseFloat(str, reflect.Float64)
			if err != nil {
				return nil, err
			}
			vi := (int64)(math.Floor(f.(float64) + 0.5))
			return intByKind(vi, kind), nil
		}
		var vi int64
		// vi, err = strconv.ParseInt(str, rt.Bits())
		vi, err = tool.N[int64](str)
		if err == nil {
			v = intByKind(vi, kind)
		}
		return
	}
	parseUint := func(str string, kind reflect.Kind) (v any, err error) {
		if strings.Contains(str, ".") {
			f, err := parseFloat(str, reflect.Float64)
			if err != nil {
				return nil, err
			}
			vi := (uint64)(math.Floor(f.(float64) + 0.5))
			return uintByKind(vi, kind), nil
		}
		var vi uint64
		// vi, err = strconv.ParseUint(str, rt.Bits())
		vi, err = tool.N[uint64](str)
		if err == nil {
			v = uintByKind(vi, kind)
		}
		return
	}

	switch k := targetType.Kind(); k {
	// case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	// case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	// case reflect.Float32, reflect.Float64:

	case reflect.Float32, reflect.Float64:
		ret, err = parseFloat(str, k)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ret, err = parseInt(str, k)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		ret, err = parseUint(str, k)
	}
	return
}

func toTimeTime(str string, targetType reflect.Type) (ret any, err error) {
	var tm time.Time
	tm, err = times.SmartParseTime(str)
	ret, _ = tm, targetType
	return
}

func toTimeTimePtr(str string, targetType reflect.Type) (ret any, err error) {
	var tm time.Time
	tm, err = times.SmartParseTime(str)
	ret, _ = &tm, targetType
	return
}

func toTimeDuration(str string, _ reflect.Type) (ret any, err error) {
	var tm time.Duration
	tm, err = times.ParseDuration(str)
	ret = tm
	return
}

func (s *toS) getcvts() map[reflect.Type]Converter {
	if s.cvts != nil {
		return s.cvts
	}
	return defcvts()
}
