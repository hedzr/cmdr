// Copyright © 2019 Hedzr Yeh.

package cmdr

import (
	"bytes"
	"fmt"
	"reflect"

	"gopkg.in/hedzr/errors.v3"

	errorsStd "errors"
)

var (
	// ErrShouldBeStopException tips `Exec()` canceled the following actions after `PreAction()`
	ErrShouldBeStopException = SetAsAnIgnorableError(newErrorWithMsg("stop me right now"), true)

	// ErrNotImpl requests cmdr launch the internal default action (see defaultAction) after returning from user's action.
	//
	// There's 3 ways to make default impl: set action to nil; use envvar `FORCE_DEFAULT_ACTION=1`; return `cmdr.ErrNotImpl` in your action handler:
	//    ```go
	//    listCmd := cmdr.NewSubCmd().Titles("list", "ls").
	//    Description("list projects in .poi.yml").
	//    Action(func(cmd *cmdr.Command, args []string) (err error) {
	//    	err = cmdr.ErrNotImpl // notify caller fallback to defaultAction
	//    	return
	//    }).
	//    AttachTo(root)
	//    ```
	ErrNotImpl = errorsStd.New("not implements")

	// ErrBadArg is a generic error for user
	ErrBadArg = newErrorWithMsg("bad argument")

	errWrongEnumValue = newErrTmpl("unexpected enumerable value '%s' for option '%s', under command '%s'")
)

// NewErrorForCmdr creates a *ErrorForCmdr object with ignorable field supports.
//
// For examples:
//
//	err := cmdr.NewErrorForCmdr("error here")
//	cmdr.SetAsAnIgnorableError(err, true)
//	if cmdr.IsIgnorable(err) {
//	    // do jump out and ignore it
//	}
func NewErrorForCmdr(msg string, args ...interface{}) error {
	return newErr(msg, args...)
}

// UnwrapError try extracting the *ErrorForCmdr object from a given error
// object even if it'd been wrapped deeply.
func UnwrapError(err error) (e *ErrorForCmdr, ignorable, ok bool) {
	ok = errors.As(err, &e)
	if ok && e != nil {
		ignorable = e.Ignorable
	}
	return
}

// UnwrapLiveArgsFromCmdrError try extracting the *ErrorForCmdr object
// from a given error object, and return the liveArgs field for usages.
func UnwrapLiveArgsFromCmdrError(err error) (liveArgs []interface{}, e *ErrorForCmdr, ok bool) {
	ok = errors.As(err, &e)
	if ok && e != nil {
		liveArgs = e.liveArgs
	}
	return
}

// // UnwrapInnerErrorsFromCmdrError try extracting the *ErrorForCmdr object
// // from a given error object, and return the inner errors attached for usages.
// func UnwrapInnerErrorsFromCmdrError(err error) (errs []error) {
//	var e *ErrorForCmdr
//	ok := errors.As(err, &e)
//	if ok && e != nil {
//		errs = []error{e.causer}
//
//	} else {
//		if ewc, ok := err.(interface{ Causes() (errs []error) }); ok {
//			return ewc.Causes()
//		}
//		if ewc, ok := err.(interface{ Cause() (err error) }); ok {
//			return []error{ewc.Cause()}
//		}
//
//		defer func() {
//			recover() // for errors.As v2.1.9 and lower
//		}()
//		var e1 *errors.WithCauses
//		ok = errors.As(err, &e1)
//		if ok && e1 != nil {
//			errs = e1.Causes()
//		}
//	}
//	return
// }

// IsIgnorableError tests if an error is a *ErrorForCmdr and its Ignorable field is true
func IsIgnorableError(err error) bool {
	if err == nil {
		return true
	}
	var e *ErrorForCmdr
	ok := errors.As(err, &e)
	return ok && e.Ignorable
}

// SetAsAnIgnorableError checks err is a *ErrorForCmdr object and set its
// Ignorable field.
func SetAsAnIgnorableError(err error, ignorable bool) error {
	var e *ErrorForCmdr
	ok := errors.As(err, &e)
	if ok {
		e.Ignorable = ignorable
	}
	return err
}

// // AttachErrorsTo wraps innerErrors into err if it's a *ErrorForCmdr as
// // a container. For the general error object, AttachErrorTo forwards it
// // to hedzr/errors to try to attach the causes.
// func AttachErrorsTo(err error, causes ...error) error {
//	var e *ErrorForCmdr
//	ok := errors.As(err, &e)
//	if ok {
//		e.Attach(causes...)
//	} else if errors.CanAttach(err) {
//		//if z, ok := err.(interface{ Attach(errs ...error) }); ok {
//		//	z.Attach(causes...)
//		//}
//
//		if ewc, ok := err.(interface{ Attach(errs ...error) }); ok {
//			ewc.Attach(causes...)
//		} else if eWC, ok := err.(interface{ Attach(errs ...error) bool }); ok {
//			eWC.Attach(causes...)
//		}
//
//	}
//	return err
// }

// AttachLiveArgsTo wraps liveArgs into err if it's a *ErrorForCmdr as
// a container.
func AttachLiveArgsTo(err error, liveArgs ...interface{}) error {
	var e *ErrorForCmdr
	ok := errors.As(err, &e)
	if ok {
		e.liveArgs = append(e.liveArgs, liveArgs...)
	}
	return err
}

// ErrorForCmdr structure
type ErrorForCmdr struct {
	// Inner     error
	// Msg       string

	// Ignorable represents this is a software logical signal rather
	// than a really programming error state.
	//
	// cmdr provides a standard Ignorable error object: ErrShouldBeStopException
	causer   error
	msg      string
	liveArgs []interface{}

	Ignorable bool
}

// newError formats a ErrorForCmdr object
func newError(ignorable bool, srcTemplate error, livedArgs ...interface{}) error {
	// log.Printf("--- newError: sourceTemplate args: %v", livedArgs)
	var e error
	switch v := srcTemplate.(type) { //nolint:errorlint //like it
	case *ErrorForCmdr:
		e = v.FormatNew(ignorable, livedArgs...)
	case *errors.WithStackInfo:
		if ex, ok := v.Cause().(*ErrorForCmdr); ok { //nolint:errorlint //like it
			e = ex.FormatNew(ignorable, livedArgs...)
		}
	}
	return e
	// if len(args) > 0 {
	// 	return &ErrorForCmdr{Inner: nil, Ignorable: ignorable, Msg: fmt.Sprintf(inner.Error(), args...)}
	// }
	// return &ErrorForCmdr{Inner: inner, Ignorable: ignorable}
}

// newErrorWithMsg formats a ErrorForCmdr object
func newErrorWithMsg(msg string, inners ...error) error {
	return newErr(msg).WithErrors(inners...)
	// return &ErrorForCmdr{Inner: inner, Ignorable: false, Msg: msg}
}

// func (s *ErrorForCmdr) Error() string {
// 	if s.Inner != nil {
// 		return fmt.Sprintf("Error: %v. Inner: %v", s.Msg, s.Inner.Error())
// 	}
// 	return s.Msg
// }

// newErr creates a *errors.WithStackInfo object
func newErr(msg string, args ...interface{}) *errors.WithStackInfo {
	// return &ErrorForCmdr{ExtErr: *errors.New(msg, args...)}
	return withIgnorable(false, nil, msg, args...).(*errors.WithStackInfo) //nolint:errorlint //like it
}

// newErrTmpl creates a *errors.WithStackInfo object
func newErrTmpl(tmpl string) *errors.WithStackInfo {
	return withIgnorable(false, nil, tmpl).(*errors.WithStackInfo) //nolint:errorlint //like it
}

// withIgnorable formats a wrapped error object with error code.
func withIgnorable(ignorable bool, err error, message string, args ...interface{}) error {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	err = &ErrorForCmdr{
		Ignorable: ignorable,
		causer:    err,
		msg:       message,
	}
	// n := errors.WithStack(err).(*errors.WithStackInfo)
	// return n.SetCause(err)
	return errors.WithStack(err)
}

func (w *ErrorForCmdr) Error() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprint(w.Ignorable))
	if len(w.msg) > 0 {
		buf.WriteRune('|')
		buf.WriteString(w.msg)
	}
	if w.causer != nil {
		buf.WriteRune('|')
		buf.WriteString(w.causer.Error())
	}
	return buf.String()
}

// FormatNew creates a new error object based on this error template 'w'.
//
// Example:
//
//	errTmpl1001 := BUG1001.NewTemplate("something is wrong %v")
//	err4 := errTmpl1001.FormatNew("ok").Attach(errBug1)
//	fmt.Println(err4)
//	fmt.Printf("%+v\n", err4)
func (w *ErrorForCmdr) FormatNew(ignorable bool, livedArgs ...interface{}) *errors.WithStackInfo {
	x := withIgnorable(ignorable, w.causer, w.msg, livedArgs...).(*errors.WithStackInfo) //nolint:errcheck,errorlint //like it
	x.Cause().(*ErrorForCmdr).liveArgs = livedArgs                                       //nolint:errcheck,errorlint //like it
	return x
}

// Attach appends errs.
// For ErrorForCmdr, only one last error will be kept.
func (w *ErrorForCmdr) Attach(errs ...error) {
	for _, err := range errs {
		w.causer = err
	}
}

// Cause returns the underlying cause of the error recursively,
// if possible.
func (w *ErrorForCmdr) Cause() error {
	return w.causer
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func (w *ErrorForCmdr) Unwrap() error {
	return w.causer
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true.
func (w *ErrorForCmdr) As(target interface{}) bool {
	if target == nil {
		panic("errors: target cannot be nil")
	}
	val := reflect.ValueOf(target)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr || val.IsNil() {
		panic("errors: target must be a non-nil pointer")
	}
	if e := typ.Elem(); e.Kind() != reflect.Interface && !e.Implements(errorType) {
		panic("errors: *target must be interface or implement error")
	}
	targetType := typ.Elem()
	err := w.causer
	for err != nil {
		if reflect.TypeOf(err).AssignableTo(targetType) {
			val.Elem().Set(reflect.ValueOf(err))
			return true
		}
		if x, ok := err.(interface{ As(interface{}) bool }); ok && x.As(target) { //nolint:errorlint //like it
			return true
		}
		err = errors.Unwrap(err)
	}
	return false
}

// Is reports whether any error in err's chain matches target.
func (w *ErrorForCmdr) Is(target error) bool {
	if target == nil {
		return w.causer == target //nolint:errorlint //like it
	}

	isComparable := reflect.TypeOf(target).Comparable()
	for {
		if isComparable && w.causer == target { //nolint:errorlint //like it
			return true
		}
		if x, ok := w.causer.(interface{ Is(error) bool }); ok && x.Is(target) { //nolint:errorlint //like it
			return true
		}
		// TODO: consider supporting target.Is(err). This would allow
		// user-definable predicates, but also may allow for coping with sloppy
		// APIs, thereby making it easier to get away with them.
		if err := errors.Unwrap(w.causer); err == nil {
			return false
		}
	}
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()
