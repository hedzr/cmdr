// Copyright © 2020 Hedzr Yeh.

package cmdr

import (
	"gopkg.in/hedzr/errors.v2"
	"testing"
)

func TestErrorForCmdr_As(t *testing.T) {
	a := &ErrorForCmdr{
		Ignorable: false,
		causer:    nil,
		msg:       "test error 1",
		livedArgs: []interface{}{2},
	}
	t.Logf("a is: %v, %T", a, a)

	e := newError(false, a)
	if _, ok := e.(*errors.WithStackInfo); !ok {
		t.Fatal(e)
	}

	e1 := newErr("test error %d", 2)
	//if _, ok := e1.(*errors.WithStackInfo); !ok {
	//	t.Fatal(e1)
	//}

	var et *errors.WithStackInfo
	if errors.Is(e, et) {
		if !errors.As(e, &et) {
			t.Fatal("cannot errors.As(e, -> *errors.WithStackInfo)")
		}
	} else {
		t.Fatal("Is detection failed")
	}

	errors.Is(e, nil)

	t.Logf("e has Causer: %v / %v | unwrapped: %v | ",
		e.(*errors.WithStackInfo).Cause(),
		errors.Cause(e),
		errors.Unwrap(e),
	)
	t.Logf("e1 is: %v, %T", e1, e1)
}
