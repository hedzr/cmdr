// Package atoa convert a string to any type following a given meme.
//
// For example:
//
//	v, _ = atoa.Parse("8.97", int32(9))
//	assert.Equal(t, v, int32(8))
//	v = atoa.MustParse("8.97", int32(9), atoa.WithFeatures(atoa.RoundNumbers))
//	assert.Equal(t, v, int32(9))
//	v = atoa.MustParse("apple=1, banana=2, orange=3", map[string]int{})
//	assert.Equal(t, v, map[string]int{"apple": 1, "banana": 2, "orange": 3})
//	v = atoa.MustParse("8,9,7", []int{})
//	assert.Equal(t, v, []int{8, 9, 7})
package atoa

import (
	"context"
	"encoding"
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hedzr/evendeep/ref"
	logz "github.com/hedzr/logg/slog"

	"github.com/hedzr/cmdr/v2/internal/tool"
	"github.com/hedzr/cmdr/v2/pkg/exec"
)

// Parse convert a string to any type following a given meme.
//
// For example:
//
//	v, _ := Parse("8.97", int32(9))
//	assert.Equal(t, v, int32(8))
//	v = MustParse("8.97", int32(9), WithFeatures(RoundNumbers))
//	assert.Equal(t, v, int32(9))
//	v = MustParse("apple=1, banana=2, orange=3", map[string]int{})
//	assert.Equal(t, v, map[string]int{"apple": 1, "banana": 2, "orange": 3})
//	v = MustParse("8,9,7", []int{})
//	assert.Equal(t, v, []int{8, 9, 7})
func Parse(str string, meme any, opts ...Opt) (v any, err error) {
	var s toS
	for _, opt := range opts {
		opt(&s)
	}

	return s.Parse(str, meme)
}

// MustParse convert a string to any type following a given meme.
//
// For example:
//
//	v = MustParse("8.97", int32(9), WithFeatures(RoundNumbers))
//	assert.Equal(t, v, int32(9))
func MustParse(str string, meme any, opts ...Opt) (v any) {
	if vv, err := Parse(str, meme, opts...); err == nil {
		v = vv
	}
	return
}

func New(opts ...Opt) *toS {
	var s toS
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}

type Opt func(s *toS)

type toS struct {
	cvts map[reflect.Type]Converter
}

func (s *toS) Parse(str string, meme any) (v any, err error) { //nolint:revive
	rt := reflect.TypeOf(meme)
	return s.parseImpl(str, rt, meme)
}

func (s *toS) parseImpl(str string, rt reflect.Type, meme any) (v any, err error) { //nolint:revive
	if meme == nil {
		err = errors.New("meme must be a valid value rather than nil")
		return
	}

	if str == "" {
		return
	}

	if cvt, ok := s.getcvts()[rt]; ok {
		v, err = cvt(str, rt)
		return
	}

	// some types cannot be recognized via reflect types
	switch meme.(type) {
	case time.Duration:
		v, err = toTimeDuration(str, rt)
		return
	case encoding.TextUnmarshaler:
		if rt.Kind() == reflect.Pointer {
			rt = rt.Elem() //nolint:revive
			rv := reflect.New(rt)
			ctx := context.Background()
			logz.DebugContext(ctx, "[cmdr] toS.parseImpl - rv", "rv", ref.Valfmt(&rv))
			err = rv.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(str))
			if err == nil {
				v = rv.Interface()
			}
			return
		}

		rv := reflect.New(rt)
		err = rv.Elem().Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(str))
		if err == nil {
			v = rv.Elem().Interface()
		}
		return
	}

	switch kind := rt.Kind(); kind {
	case reflect.Interface, reflect.Pointer:
		t := rt.Elem()
		return s.parseImpl(str, t, meme)

	case reflect.Complex128, reflect.Complex64:
		var vi complex128
		vi, err = strconv.ParseComplex(str, rt.Bits())
		if err == nil {
			v = comByKind(vi, kind)
		}
	case reflect.Float32, reflect.Float64:
		var vi float64
		// vi, err = strconv.ParseFloat(str, rt.Bits())
		vi, err = tool.N[float64](str)
		if err == nil {
			v = fltByKind(vi, kind)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var vi int64
		// vi, err = strconv.ParseInt(str, 0, rt.Bits())
		vi, err = tool.N[int64](str)
		if err == nil {
			v = intByKind(vi, kind)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		var vi uint64
		// vi, err = strconv.ParseUint(str, 0, rt.Bits())
		vi, err = tool.N[uint64](str)
		if err == nil {
			v = uintByKind(vi, kind)
		}

	case reflect.Chan, reflect.UnsafePointer:
		// not support
		v, err = str, errNotSupport

	case reflect.Func:
		// calc it
		v, err = str, errNotSupport

	case reflect.Array:
		return s.parseArray(kind, rt, str, meme)
	case reflect.Slice:
		return s.parseSlice(kind, rt, str, meme)
	case reflect.Map:
		return s.parseMap(kind, rt, str, meme)
	case reflect.Struct:
		return s.parseStruct(kind, rt, str, meme)

	case reflect.String:
		v = exec.StripQuotes(str)
	case reflect.Bool:
		v = tool.ToBool(str)

	default:
		v, err = str, errParse
	}
	return
}

func (s *toS) parseStruct(preferKind reflect.Kind, typStruct reflect.Type, str string, meme any) (v any, err error) {
	// check types: time.Time, bytes.Buffer, ...
	return s.parseArray(preferKind, typStruct, str, meme)
}

func (s *toS) parseMap(preferKind reflect.Kind, typMap reflect.Type, str string, meme any) (v any, err error) {
	return s.parseArray(preferKind, typMap, str, meme)
}

func (s *toS) parseSlice(preferKind reflect.Kind, typSlice reflect.Type, str string, meme any) (v any, err error) {
	return s.parseArray(preferKind, typSlice, str, meme)
}

func (s *toS) parseArray(preferKind reflect.Kind, typArray reflect.Type, str string, meme any) (v any, err error) {
	runes := []rune(str)
	var pos, position int
	position, v, err = s.stepComplexObject(preferKind, typArray, runes, pos, meme)
	if position == pos || position < len(runes) {
		ctx := context.Background()
		logz.WarnContext(ctx, "[cmdr] the given string is empty or has too much data?", "str", str, "posAfterParsed", position)
	}
	return
}

func (s *toS) parseMapOld(rt reflect.Type, str string, meme any) (v any, err error) { //nolint:revive,unused
	var ate bool
	str, ate = eat(str, '{') //nolint:revive
	if ate {
		str, ate = eatTail(str, '}') //nolint:revive,staticcheck,ineffassign
	}
	ssa := strings.Split(str, ",")

	rv := reflect.MakeMapWithSize(rt, len(ssa))
	kt, vt := rt.Key(), rt.Elem()

	re := regexp.MustCompile(`(.*)[=:](.*)`)
	for _, txt := range ssa {
		a := re.FindAllStringSubmatch(txt, -1)
		if len(a) > 0 {
			b := a[0]
			if len(b) > 1 {
				k1, v1 := strings.TrimSpace(b[1]), strings.TrimSpace(b[2])
				kv, err1 := s.parseImpl(k1, kt, meme)
				if err1 != nil {
					err = err1
					return
				}
				vv, err2 := s.parseImpl(v1, vt, meme)
				if err2 != nil {
					err = err2
					return
				}
				rv.SetMapIndex(reflect.ValueOf(kv), reflect.ValueOf(vv))
			}
		}
	}
	v = rv.Interface()
	return
}

func (s *toS) parseSliceOld(rt reflect.Type, str string, meme any) (v any, err error) { //nolint:revive,unused
	var ate bool
	str, ate = eat(str, '[') //nolint:revive
	if ate {
		str, ate = eatTail(str, ']') //nolint:revive,staticcheck,ineffassign
	}
	ssa := strings.Split(str, ",")

	rv := reflect.MakeSlice(rt, 0, len(ssa))
	elt := rt.Elem()

	for _, txt := range ssa {
		var vv any
		vv, err = s.parseImpl(strings.TrimSpace(txt), elt, meme)
		if err == nil {
			rv = reflect.Append(rv, reflect.ValueOf(vv))
		}
	}
	v = rv.Interface()
	return
}

var (
	errParse      = errors.New("cannot parse")
	errNotSupport = errors.New("not support")
)
