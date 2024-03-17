package tool

import (
	"errors"
	"reflect"
	"strconv"

	"golang.org/x/exp/constraints"
)

func Must[T constraints.Float | constraints.Integer](s string) T {
	t, err := N[T](s)
	if err != nil {
		panic(err)
	}
	return t
}

func N[T constraints.Float | constraints.Integer](s string) (ret T, err error) {
	rt := reflect.TypeOf(ret)
	switch k := rt.Kind(); k {
	case reflect.Float32, reflect.Float64:
		var t float64
		t, err = strconv.ParseFloat(s, rt.Bits())
		if err != nil {
			return
		}
		ret = T(t)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var t int64
		t, err = strconv.ParseInt(s, 0, rt.Bits())
		if err != nil {
			return
		}
		ret = T(t)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		var t uint64
		t, err = strconv.ParseUint(s, 0, rt.Bits())
		if err != nil {
			return
		}
		ret = T(t)
	default:
		err = errors.New("cannot parse")
	}
	return
}
