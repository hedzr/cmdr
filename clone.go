/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bytes"
	"encoding/gob"
	"github.com/hedzr/log"
	"gopkg.in/hedzr/errors.v2"
	"reflect"
	"strings"
)

type (
	// Copier interface
	// Copier is based on from github.com/jinzhu/copier
	Copier interface {
		SetIgnoreNames(ignoreNames ...string) Copier
		SetEachFieldAlways(b bool) Copier
		Copy(toValue interface{}, fromValue interface{}, ignoreNames ...string) (err error)
	}

	// copierImpl impl
	copierImpl struct {
		KeepIfFromIsNil  bool // 源字段值为nil指针时，目标字段的值保持不变
		KeepIfFromIsZero bool // 源字段值为未初始化的零值时，目标字段的值保持不变 // 此条尚未实现
		ZeroIfEqualsFrom bool // 源和目标字段值相同时，目标字段被清除为未初始化的零值
		IgnoreNames      []string
		EachFieldAlways  bool
		IgnoreIfNotEqual bool
	}
)

var (
	// GormDefaultCopier used for gorm
	GormDefaultCopier = &copierImpl{KeepIfFromIsNil: true, ZeroIfEqualsFrom: true, KeepIfFromIsZero: true, EachFieldAlways: true}
	// StandardCopier is a normal copier
	StandardCopier = &copierImpl{}
)

// CloneViaGob do deep-clone with gob supports
func CloneViaGob(to, from interface{}) (err error) {
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	dec := gob.NewDecoder(buff)

	err = enc.Encode(from)
	if err != nil {
		return
	}
	err = dec.Decode(to)
	//if err != nil {
	//	return
	//}
	return
}

// Clone deep copy source to target
func Clone(fromValue, toValue interface{}) interface{} {
	_ = StandardCopier.Copy(toValue, fromValue)
	return toValue
}

// SetIgnoreNames give a group of ignored fieldNames
func (s *copierImpl) SetIgnoreNames(ignoreNames ...string) Copier {
	s.IgnoreNames = ignoreNames
	return s
}

func (s *copierImpl) SetEachFieldAlways(b bool) Copier {
	s.EachFieldAlways = b
	return s
}

// Copy copy things
func (s *copierImpl) Copy(toValue interface{}, fromValue interface{}, ignoreNames ...string) (err error) {
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
		if from.Type().AssignableTo(to.Type()) && !s.EachFieldAlways {
			to.Set(from)
			return
		}
		if to.Kind() == reflect.Struct {
			amount = 1
		}
	}

	err = s.copyAll(amount, isSlice, from, to, fromType, toType, append(ignoreNames, s.IgnoreNames...))
	return
}

func (s *copierImpl) copyAll(amount int, isSlice bool, from, to reflect.Value, fromType, toType reflect.Type, ignoreNames []string) error {
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
		if err := s.copyFieldToField(dest, source, fromType, ignoreNames); err != nil {
			return err
		}

		// Copy from method to field
		if err := s.copyMethodToField(dest, source, toType); err != nil {
			return err
		}

		if isSlice {
			if dest.Addr().Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest.Addr()))
			} else if dest.Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest))
			}
		}
	}
	return nil
}

func safetyTarget(dest reflect.Value, fromType reflect.Type, ignoreNames []string) {
	if fromType.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < fromType.NumField(); i++ {
		fld := fromType.Field(i)
		if fld.Anonymous && fld.Type.Kind() == reflect.Ptr {
			v := dest.Field(i)
			if v.IsValid() && v.IsNil() && v.CanSet() {
				v.Set(reflect.New(fld.Type.Elem()))
				safetyTarget(v, fld.Type, ignoreNames)
			}
		}
	}
}

func (s *copierImpl) copyFieldToField(dest, source reflect.Value, fromType reflect.Type, ignoreNames []string) error {
	var names []string
	var name string
	defer func() {
		if e := recover(); e != nil {
			log.Errorf("failed on copying field %q : %v", name, e)
			log.Errorf("    past fields: %v", names)
			panic(e)
		}
	}()

	safetyTarget(dest, fromType, ignoreNames)

	// Copy from field to field or method
	for _, field := range s.deepFields(fromType) {
		name = field.Name
		names = append(names, name)
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
	return nil
}

func (s *copierImpl) copyMethodToField(dest, source reflect.Value, toType reflect.Type) error {
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
	return nil
}

func (s *copierImpl) deepFields(reflectType reflect.Type) []reflect.StructField {
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

func (s *copierImpl) indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func (s *copierImpl) indirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

func contains(names []string, name string) bool {
	for _, n := range names {
		if strings.EqualFold(n, name) {
			return true
		}
	}
	return false
}

func partialContains(names []string, partialNeedle string) (index int, matched string, contains bool) {
	for ix, n := range names {
		if strings.Contains(n, partialNeedle) {
			return ix, n, true
		}
	}
	return -1, "", false
}

func equal(to, from reflect.Value) bool {
	switch to.Kind() {
	case reflect.Bool:
		return from.Bool() == to.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return from.Int() == to.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return from.Uint() == to.Uint()
	// case reflect.Uintptr:
	// 	return from.Pointer() == to.Pointer()
	case reflect.Float32, reflect.Float64:
		return from.Float() == to.Float()
	case reflect.Complex64, reflect.Complex128:
		return from.Complex() == to.Complex()
	case reflect.Array:
		return equalArray(to, from)

	// case reflect.Chan:
	// 	// cmdr.Logger.Warnf("unrecognized type: %v", to.Kind())
	// case reflect.Func:
	// 	// cmdr.Logger.Warnf("unrecognized type: %v", to.Kind())
	// case reflect.Interface:
	// 	// cmdr.Logger.Warnf("unrecognized type: %v", to.Kind())
	// case reflect.Map:
	// 	// cmdr.Logger.Warnf("unrecognized type: %v", to.Kind())
	// case reflect.Ptr:
	// 	return from.Pointer() == to.Pointer()
	case reflect.Slice:
		return equalSlice(to, from)

	case reflect.String:
		return strings.EqualFold(from.String(), to.String())

		// case reflect.Struct:
		// 	// cmdr.Logger.Warnf("unrecognized type: %v", to.Kind())
		// case reflect.UnsafePointer:
		// 	// cmdr.Logger.Warnf("unrecognized type: %v", to.Kind())
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

func equalArray(to, from reflect.Value) bool {
	if from.Len() != to.Len() {
		return false
	}
	if from.Len() == 0 {
		return true
	}
	//for i := 0; i < from.Len(); i++ {
	//	if !equal(from.Slice(i, i+1), to.Slice(i, i+1)) {
	//		return false
	//	}
	//}
	//return true

	x := make(map[interface{}]bool)
	for i := 0; i < from.Len(); i++ {
		x[from.Index(i).Interface()] = true
	}
	for i := 0; i < from.Len(); i++ {
		delete(x, to.Index(i).Interface())
	}
	return len(x) == 0
}

func equalSlice(to, from reflect.Value) bool {
	if from.Len() != to.Len() {
		return false
	}
	if from.Len() == 0 {
		return true
	}

	x := make(map[interface{}]bool)
	for i := 0; i < from.Len(); i++ {
		x[from.Index(i).Interface()] = true
	}
	for i := 0; i < from.Len(); i++ {
		v := to.Index(i).Interface()
		delete(x, v)
	}
	return len(x) == 0
}

func setDefault(to reflect.Value) {
	switch to.Kind() {
	case reflect.Bool:
		to.SetBool(false)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		to.SetInt(0)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		to.SetUint(0)
	// case reflect.Uintptr:
	// 	to.SetPointer(nil)
	case reflect.Float32, reflect.Float64:
		to.SetFloat(0)
	case reflect.Complex64, reflect.Complex128:
		to.SetComplex(0)
	case reflect.Array:
		for i := 0; i < to.Len(); i++ {
			setDefault(to.Index(i))
		}
	// case reflect.Chan:
	// 	// cmdr.Logger.Warnf("unrecognized type: %v", to.Kind())
	// case reflect.Func:
	// 	// cmdr.Logger.Warnf("unrecognized type: %v", to.Kind())
	// case reflect.Interface:
	// 	// cmdr.Logger.Warnf("unrecognized type: %v", to.Kind())
	// case reflect.Map:
	// 	// cmdr.Logger.Warnf("unrecognized type: %v", to.Kind())
	// case reflect.Ptr:
	// 	to.SetPointer(nil)
	case reflect.Slice:
		to.SetLen(0)
	case reflect.String:
		to.SetString("")
		// case reflect.Struct:
		// 	// cmdr.Logger.Warnf("unrecognized type: %v", to.Kind())
		// case reflect.UnsafePointer:
		// 	// cmdr.Logger.Warnf("unrecognized type: %v", to.Kind())
	}
}

func (s *copierImpl) setCvt(to, from reflect.Value) {
	if !(s.KeepIfFromIsNil && isNil(from)) {
		if !(s.KeepIfFromIsZero && IsZero(from)) {
			if equal(to, from) && s.ZeroIfEqualsFrom {
				setDefault(to)
			} else if s.IgnoreIfNotEqual == false {
				to.Set(from.Convert(to.Type()))
			}
		}
	}
}

func (s *copierImpl) set(to, from reflect.Value) bool {
	if from.IsValid() {
		if to.Kind() == reflect.Ptr {
			// if s.setPtr(to, from) {
			// 	return true
			// }

			// set `to` to nil if from is nil
			if from.Kind() == reflect.Ptr && from.IsNil() {
				if !s.KeepIfFromIsNil && !s.KeepIfFromIsZero && !s.IgnoreIfNotEqual {
					to.Set(reflect.Zero(to.Type()))
				}
				return true
			} else if to.IsNil() {
				to.Set(reflect.New(to.Type().Elem()))
			}
			to = to.Elem()
		}

		if from.Type().ConvertibleTo(to.Type()) {
			s.setCvt(to, from)
			// } else if scanner, ok := to.Addr().Interface().(sql.Scanner); ok {
			// 	err := scanner.Scan(from.Interface())
			// 	if err != nil {
			// 		return false
			// 	}
		} else if from.Kind() == reflect.Ptr {
			if !s.IgnoreIfNotEqual {
				return s.set(to, from.Elem())
			}
		} else {
			return false
		}
	}
	return true
}
