// Copyright © 2020 Hedzr Yeh.

// +build !go1.13

package cmdr

import "math"

// IsZero reports whether v is the zero value for its type.
// It panics if the argument is invalid.
func IsZero(v Value) bool {
	switch v.kind() {
	case Bool:
		return !v.Bool()
	case Int, Int8, Int16, Int32, Int64:
		return v.Int() == 0
	case Uint, Uint8, Uint16, Uint32, Uint64, Uintptr:
		return v.Uint() == 0
	case Float32, Float64:
		return math.Float64bits(v.Float()) == 0
	case Complex64, Complex128:
		c := v.Complex()
		return math.Float64bits(real(c)) == 0 && math.Float64bits(imag(c)) == 0
	case Array:
		return isZeroArray(v)
	case Chan, Func, Interface, Map, Ptr, Slice, UnsafePointer:
		return v.IsNil()
	case String:
		return v.Len() == 0
	case Struct:
		return isZeroStruct(v)
	default:
		// This should never happens, but will act as a safeguard for
		// later, as a default value doesn't makes sense here.
		panic(&ValueError{"reflect.Value.IsZero", v.Kind()})
	}
}

func isZeroArray(v Value) bool {
	for i := 0; i < v.Len(); i++ {
		if !v.Index(i).IsZero() {
			return false
		}
	}
	return true
}

func isZeroStruct(v Value) bool {
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsZero() {
			return false
		}
	}
	return true
}
