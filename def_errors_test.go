// Copyright © 2020 Hedzr Yeh.

package cmdr

import (
	"gopkg.in/hedzr/errors.v3"
	"testing"
)

func TestErrors(t *testing.T) {
	a := &ErrorForCmdr{
		Ignorable: false,
		causer:    nil,
		msg:       "test error 1",
		liveArgs:  []interface{}{2},
	}
	if _, ignorable, ok := UnwrapError(a); ok && !ignorable {
		// ok
	} else {
		t.Fatal(a)
	}

	for _, err := range []error{
		ErrShouldBeStopException, ErrBadArg, errWrongEnumValue,
	} {
		if _, _, ok := UnwrapError(err); ok {
			// ok
		} else {
			t.Fatal(err)
		}
	}

	//for _, err := range []error{
	//	ErrShouldBeStopException,
	//	AttachErrorsTo(AttachLiveArgsTo(SetAsAnIgnorableError(NewErrorForCmdr("*ErrorForCmdr here"), true), []int{2, 5}, "bad"), io.EOF),
	//} {
	//	if IsIgnorableError(err) {
	//		// ok
	//		liveArgs, e, ok := UnwrapLiveArgsFromCmdrError(err)
	//		t.Logf("liveArgs: %v, e: %v, ok: %v", liveArgs, e, ok)
	//		//t.Logf("Unwrap() returned: %v", e.Unwrap())
	//
	//		errs := UnwrapInnerErrorsFromCmdrError(err)
	//		t.Logf("inner errors are: %v", errs)
	//	} else {
	//		t.Fatal(err)
	//	}
	//}
	//
	//// 1.
	//ec := errors.NewContainer("container of errors")
	//for _, e := range []error{io.EOF, io.ErrClosedPipe} {
	//	ec.Attach(e)
	//}
	//errs := UnwrapInnerErrorsFromCmdrError(&ew{(*errors.WithCauses)(ec)})
	//t.Logf("inner errors are: %v, error is: %v", errs, ec.Error())
	//
	//// 2.
	//ewc := errors.WithCause(io.EOF, "*withCauses object here")
	//if e1 := AttachErrorsTo(ewc, io.EOF, io.ErrShortBuffer); e1 != nil {
	//	errs = UnwrapInnerErrorsFromCmdrError(e1)
	//	t.Logf("inner errors are: %v, error is: %v", errs, e1)
	//}

	for _, err := range []error{
		ErrBadArg, errWrongEnumValue,
	} {
		if !IsIgnorableError(err) {
			// ok
		} else {
			t.Fatal(err)
		}
	}
}

type ew struct {
	msg string
	// *errors.WithCauses
}

func (e *ew) Error() string {
	return e.msg // e.WithCauses.Error().Error()
}

func TestErrorForCmdr(t *testing.T) {
	a := &ErrorForCmdr{
		Ignorable: false,
		causer:    nil,
		msg:       "test error 1",
		liveArgs:  []interface{}{2},
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

	if _, ok := e.(*errors.WithStackInfo); ok {
		var et *errors.WithStackInfo
		if !errors.As(e, &et) {
			t.Fatal("cannot errors.As(e, -> *errors.WithStackInfo)")
		}
	} else {
		t.Fatalf("Is detection failed: e = %+v", e)
	}

	errors.Is(e, nil)

	t.Logf("e has Causer: %v / %v | unwrapped: %v | ",
		e.(*errors.WithStackInfo).Cause(),
		e.Error(), // errors.Cause(e),
		errors.Unwrap(e),
	)
	t.Logf("e1 is: %v, %T", e1, e1)
}
