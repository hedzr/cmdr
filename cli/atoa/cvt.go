package atoa

import (
	"reflect"
	"sync"
	"time"

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
