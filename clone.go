/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"errors"
	"reflect"
	"strings"
)

type (
	// Copier interface
	// Copier is based on from github.com/jinzhu/copier
	Copier interface {
		Copy(toValue interface{}, fromValue interface{}, ignoreNames ...string) (err error)
	}

	// CopierImpl impl
	CopierImpl struct {
		KeepIfFromIsNil  bool // 源字段值为nil指针时，目标字段的值保持不变
		ZeroIfEqualsFrom bool // 源和目标字段值相同时，目标字段被清除为未初始化的零值
		KeepIfFromIsZero bool // 源字段值为未初始化的零值时，目标字段的值保持不变 // 此条尚未实现
	}
)

var (
	// GormDefaultCopier used for gorm
	GormDefaultCopier = &CopierImpl{true, true, true}
	// StandardCopier is a normal copier
	StandardCopier = &CopierImpl{false, false, false}
)

// Clone deep copy source to target
func Clone(fromValue, toValue interface{}) interface{} {
	_ = StandardCopier.Copy(toValue, fromValue)
	return toValue
}

// Copy copy things
func (s *CopierImpl) Copy(toValue interface{}, fromValue interface{}, ignoreNames ...string) (err error) {
	var (
		isSlice bool
		amount  = 1
		from    = s.indirect(reflect.ValueOf(fromValue))
		to      = s.indirect(reflect.ValueOf(toValue))
	)

	if !to.CanAddr() {
		return errors.New("copy to value is unaddressable")
	}

	// Return is from value is invalid
	if !from.IsValid() {
		return
	}

	fromType := s.indirectType(from.Type())
	toType := s.indirectType(to.Type())

	if fromType.Kind() != reflect.Struct || toType.Kind() != reflect.Struct {
		return
	}

	if to.Kind() == reflect.Slice {
		isSlice = true
		if from.Kind() == reflect.Slice {
			amount = from.Len()
		}
	} else {
		// Just set it if possible to assign
		if from.Type().AssignableTo(to.Type()) {
			to.Set(from)
			return
		}
	}

	for i := 0; i < amount; i++ {
		var dest, source reflect.Value

		if isSlice {
			// source
			if from.Kind() == reflect.Slice {
				source = s.indirect(from.Index(i))
			} else {
				source = s.indirect(from)
			}

			// dest
			dest = s.indirect(reflect.New(toType).Elem())
		} else {
			source = s.indirect(from)
			dest = s.indirect(to)
		}

		// Copy from field to field or method
		for _, field := range s.deepFields(fromType) {
			name := field.Name
			if contains(ignoreNames, name) {
				continue
			}

			if fromField := source.FieldByName(name); fromField.IsValid() {
				// has field
				if toField := dest.FieldByName(name); toField.IsValid() {
					if toField.CanSet() {
						if !s.set(toField, fromField) {
							if err := s.Copy(toField.Addr().Interface(), fromField.Interface()); err != nil {
								return err
							}
						}
					}
				} else {
					// try to set to method
					var toMethod reflect.Value
					if dest.CanAddr() {
						toMethod = dest.Addr().MethodByName(name)
					} else {
						toMethod = dest.MethodByName(name)
					}

					if toMethod.IsValid() && toMethod.Type().NumIn() == 1 && fromField.Type().AssignableTo(toMethod.Type().In(0)) {
						toMethod.Call([]reflect.Value{fromField})
					}
				}
			}
		}

		// Copy from method to field
		for _, field := range s.deepFields(toType) {
			name := field.Name

			var fromMethod reflect.Value
			if source.CanAddr() {
				fromMethod = source.Addr().MethodByName(name)
			} else {
				fromMethod = source.MethodByName(name)
			}

			if fromMethod.IsValid() && fromMethod.Type().NumIn() == 0 && fromMethod.Type().NumOut() == 1 {
				if toField := dest.FieldByName(name); toField.IsValid() && toField.CanSet() {
					values := fromMethod.Call([]reflect.Value{})
					if len(values) >= 1 {
						s.set(toField, values[0])
					}
				}
			}
		}

		if isSlice {
			if dest.Addr().Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest.Addr()))
			} else if dest.Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest))
			}
		}
	}
	return
}

func (s *CopierImpl) deepFields(reflectType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	if reflectType = s.indirectType(reflectType); reflectType.Kind() == reflect.Struct {
		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			if v.Anonymous {
				fields = append(fields, s.deepFields(v.Type)...)
			} else {
				fields = append(fields, v)
			}
		}
	}

	return fields
}

func (s *CopierImpl) indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func (s *CopierImpl) indirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

func contains(names []string, name string) bool {
	for _, n := range names {
		if strings.EqualFold(name, n) {
			return true
		}
	}
	return false
}

func equal(to, from reflect.Value) bool {
	switch to.Kind() {
	case reflect.Bool:
		return from.Bool() == to.Bool()
	case reflect.Int:
		return from.Int() == to.Int()
	case reflect.Int8:
		return from.Int() == to.Int()
	case reflect.Int16:
		return from.Int() == to.Int()
	case reflect.Int32:
		return from.Int() == to.Int()
	case reflect.Int64:
		return from.Int() == to.Int()
	case reflect.Uint:
		return from.Uint() == to.Uint()
	case reflect.Uint8:
		return from.Uint() == to.Uint()
	case reflect.Uint16:
		return from.Uint() == to.Uint()
	case reflect.Uint32:
		return from.Uint() == to.Uint()
	case reflect.Uint64:
		return from.Uint() == to.Uint()
	case reflect.Uintptr:
		return from.Pointer() == to.Pointer()
	case reflect.Float32:
		return from.Float() == to.Float()
	case reflect.Float64:
		return from.Float() == to.Float()
	case reflect.Complex64:
		return from.Complex() == to.Complex()
	case reflect.Complex128:
		return from.Complex() == to.Complex()
	case reflect.Array:
		if from.Len() != to.Len() {
			return false
		}
		if from.Len() == 0 {
			return true
		}
		for i := 0; i < from.Len(); i++ {
			if !equal(from.Slice(i, i+1), to.Slice(i, i+1)) {
				return false
			}
		}
		return true

	case reflect.Chan:
		// logrus.Warnf("unrecognized type: %v", to.Kind())
	case reflect.Func:
		// logrus.Warnf("unrecognized type: %v", to.Kind())
	case reflect.Interface:
		// logrus.Warnf("unrecognized type: %v", to.Kind())
	case reflect.Map:
		// logrus.Warnf("unrecognized type: %v", to.Kind())
	case reflect.Ptr:
		return from.Pointer() == to.Pointer()
	case reflect.Slice:
		if from.Len() != to.Len() {
			return false
		}
		if from.Len() == 0 {
			return true
		}
		for i := 0; i < from.Len(); i++ {
			if !equal(from.Slice(i, i+1), to.Slice(i, i+1)) {
				return false
			}
		}
		return true

	case reflect.String:
		return strings.EqualFold(from.String(), to.String())

	case reflect.Struct:
		// logrus.Warnf("unrecognized type: %v", to.Kind())
	case reflect.UnsafePointer:
		// logrus.Warnf("unrecognized type: %v", to.Kind())
	}

	// if to.IsNil() && from.IsNil() {
	// 	return true
	// }
	// if to.IsNil() || from.IsNil() {
	// 	return false
	// }
	// if to.Interface() == from.Interface() {
	// 	return true
	// }

	// deep test
	return false
}

func setDefault(to reflect.Value) {
	switch to.Kind() {
	case reflect.Bool:
		to.SetBool(false)
	case reflect.Int:
		to.SetInt(0)
	case reflect.Int8:
		to.SetInt(0)
	case reflect.Int16:
		to.SetInt(0)
	case reflect.Int32:
		to.SetInt(0)
	case reflect.Int64:
		to.SetInt(0)
	case reflect.Uint:
		to.SetUint(0)
	case reflect.Uint8:
		to.SetUint(0)
	case reflect.Uint16:
		to.SetUint(0)
	case reflect.Uint32:
		to.SetUint(0)
	case reflect.Uint64:
		to.SetUint(0)
	case reflect.Uintptr:
		to.SetPointer(nil)
	case reflect.Float32:
		to.SetFloat(0)
	case reflect.Float64:
		to.SetFloat(0)
	case reflect.Complex64:
		to.SetComplex(0)
	case reflect.Complex128:
		to.SetComplex(0)
	case reflect.Array:
		to.SetLen(0)
	case reflect.Chan:
		// logrus.Warnf("unrecognized type: %v", to.Kind())
	case reflect.Func:
		// logrus.Warnf("unrecognized type: %v", to.Kind())
	case reflect.Interface:
		// logrus.Warnf("unrecognized type: %v", to.Kind())
	case reflect.Map:
		// logrus.Warnf("unrecognized type: %v", to.Kind())
	case reflect.Ptr:
		to.SetPointer(nil)
	case reflect.Slice:
		to.SetLen(0)
	case reflect.String:
		to.SetString("")
	case reflect.Struct:
		// logrus.Warnf("unrecognized type: %v", to.Kind())
	case reflect.UnsafePointer:
		// logrus.Warnf("unrecognized type: %v", to.Kind())
	}
}

func isNil(to reflect.Value) bool {
	switch to.Kind() {
	case reflect.Uintptr:
		return to.IsNil()
	case reflect.Array:
	case reflect.Chan:
		return to.IsNil()
	case reflect.Func:
		return to.IsNil()
	case reflect.Interface:
		return to.IsNil()
	case reflect.Map:
		return to.IsNil()
	case reflect.Ptr:
		return to.IsNil()
	case reflect.Slice:
	case reflect.String:
	case reflect.Struct:
		return to.IsNil()
	case reflect.UnsafePointer:
		return to.IsNil()
	}
	return false
}

func (s *CopierImpl) set(to, from reflect.Value) bool {
	if from.IsValid() {
		if to.Kind() == reflect.Ptr {
			// set `to` to nil if from is nil
			if from.Kind() == reflect.Ptr && from.IsNil() {
				to.Set(reflect.Zero(to.Type()))
				return true
			} else if to.IsNil() {
				to.Set(reflect.New(to.Type().Elem()))
			}
			to = to.Elem()
		}

		if from.Type().ConvertibleTo(to.Type()) {
			if !(s.KeepIfFromIsNil && isNil(from)) {
				if equal(to, from) && s.ZeroIfEqualsFrom {
					setDefault(to)
				} else {
					to.Set(from.Convert(to.Type()))
				}
			}
			// } else if scanner, ok := to.Addr().Interface().(sql.Scanner); ok {
			// 	err := scanner.Scan(from.Interface())
			// 	if err != nil {
			// 		return false
			// 	}
		} else if from.Kind() == reflect.Ptr {
			return s.set(to, from.Elem())
		} else {
			return false
		}
	}
	return true
}
