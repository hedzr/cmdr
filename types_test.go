/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"reflect"
	"testing"
)

func TestTypes(t *testing.T) {
	// func TestTyper(t *testing.T) {

	v := uint(3)
	tAssert(t, isTypeUint(reflect.ValueOf(v).Kind()) == true)
	vc := int(-3)
	tAssert(t, isTypeUint(reflect.ValueOf(vc).Kind()) == false)
}

func tAssert(t *testing.T, cond bool) {
	if cond == false {
		t.Error("unwanted assertion")
	} else {
		t.Log("tested.")
	}
}

func TestFindsX(t *testing.T) {
	t.Log("finds", InTesting(), randomStringPure(5), min(3, 5), min(13, 5), stripPrefix("sss", "a"), IsDigitHeavy("ds"), IsDigitHeavy("3521"))
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
		t.Logf("soundex of '%v' = %v", str, Soundex(str))
	}
}
