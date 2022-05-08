/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr_test

import (
	"reflect"
	"testing"

	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/tool"
)

func TestToBool(t *testing.T) {
	cmdr.ToBool(false)
	cmdr.ToBool(1)
	cmdr.ToBool(-1)
	cmdr.ToBool("sss")
	cmdr.ToBool(3.14, false)
}

func TestOptions_GetInt64Ex(t *testing.T) {
	o := cmdr.NewOptionsForTest()
	o.Set("test", "123")
	o.GetInt64Ex("app.test", 3)
	o.GetKibibytesEx("app.test")
	o.GetKilobytesEx("app.test")
	o.GetUintEx("app.test")
	o.GetUint64Ex("app.test")
	o.GetFloat32Ex("app.test")
	o.GetFloat64Ex("app.test")

	o.GetComplex64("app.test")
	o.GetComplex128("app.test")
	o.GetComplex64("app.test-none", 1)
	o.GetComplex128("app.test-none", 1)

	o.Set("ss", []string{"3", "5"})
	o.GetStringSlice("app.ss")
	o.Set("ss", 123)
	o.GetStringSlice("app.ss", "2", "3")
	o.Set("ss", struct {
		v int
	}{9})
	o.GetStringSlice("app.ss", "2", "3")

	o.GetStringSlice("app.ss.none", "2", "3")

	o.Set("ss", []struct {
		v int
	}{{3}, {5}})
	o.GetInt64Slice("app.ss")
}

func TestTypes(t *testing.T) {
	// func TestTyper(t *testing.T) {

	v := uint(3)
	tAssert(t, cmdr.IsTypeUint(reflect.ValueOf(v).Kind()) == true)
	tAssert(t, cmdr.IsTypeSInt(reflect.ValueOf(v).Kind()) == false)
	vc := int(-3)
	tAssert(t, cmdr.IsTypeUint(reflect.ValueOf(vc).Kind()) == false)
	tAssert(t, cmdr.IsTypeSInt(reflect.ValueOf(vc).Kind()) == true)

	f1 := float32(2)
	tAssert(t, cmdr.IsTypeFloat(reflect.ValueOf(f1).Kind()) == true)
	tAssert(t, cmdr.IsTypeFloat(reflect.ValueOf(vc).Kind()) == false)

	f2 := float64(2)
	tAssert(t, cmdr.IsTypeFloat(reflect.ValueOf(f2).Kind()) == true)
	tAssert(t, cmdr.IsTypeFloat(reflect.ValueOf(vc).Kind()) == false)

	c1 := complex(2.0, 3.0)
	tAssert(t, cmdr.IsTypeComplex(reflect.ValueOf(c1).Kind()) == true)
	tAssert(t, cmdr.IsTypeComplex(reflect.ValueOf(vc).Kind()) == false)
	tAssert(t, cmdr.IsTypeComplex(reflect.ValueOf(f2).Kind()) == false)
}

func tAssert(t *testing.T, cond bool) {
	if cond == false {
		t.Error("unwanted assertion")
	} else {
		t.Log("tested.")
	}
}

func TestFindsX(t *testing.T) {
	t.Log("finds",
		cmdr.InTesting(),
		tool.RandomStringPure(5),
		tool.Min(3, 5),
		tool.Min(13, 5),
		tool.StripPrefix("sss", "a"),
		tool.IsDigitHeavy("ds"),
		tool.IsDigitHeavy("3521"))
}

func TestSoundeX(t *testing.T) {
	for _, str := range []string{
		"flush",
		"bug",
		"this",
		"is",
		"a",
		"distinguashing",
		"mam",
		"nurde",
		"worker",
	} {
		t.Logf("soundex of '%v' = %v", str, tool.Soundex(str))
	}
}
