// Copyright © 2019 Hedzr Yeh.

package cmdr

import (
	"bytes"
	"fmt"
	"github.com/hedzr/errors"
)

var (
	// ErrShouldBeStopException tips `Exec()` cancelled the following actions after `PreAction()`
	ErrShouldBeStopException = newErrorWithMsg("stop me right now")

	// ErrBadArg is a generic error for user
	ErrBadArg = newErrorWithMsg("bad argument")

	errWrongEnumValue = newErrTmpl("unexpected enumerable value '%s' for option '%s', under command '%s'")
)

// ErrorForCmdr structure
type ErrorForCmdr struct {
	// Inner     error
	// Msg       string
	Ignorable bool
	errors.ExtErr
}

// newError formats a ErrorForCmdr object
func newError(ignorable bool, sourceTemplate *ErrorForCmdr, args ...interface{}) *ErrorForCmdr {
	// log.Printf("--- newError: sourceTemplate args: %v", args)
	e := sourceTemplate.Format(args...)
	e.Ignorable = ignorable
	return e
	// if len(args) > 0 {
	// 	return &ErrorForCmdr{Inner: nil, Ignorable: ignorable, Msg: fmt.Sprintf(inner.Error(), args...)}
	// }
	// return &ErrorForCmdr{Inner: inner, Ignorable: ignorable}
}

// newErrorWithMsg formats a ErrorForCmdr object
func newErrorWithMsg(msg string, inners ...error) *ErrorForCmdr {
	return newErr(msg).Attach(inners...)
	// return &ErrorForCmdr{Inner: inner, Ignorable: false, Msg: msg}
}

// func (s *ErrorForCmdr) Error() string {
// 	if s.Inner != nil {
// 		return fmt.Sprintf("Error: %v. Inner: %v", s.Msg, s.Inner.Error())
// 	}
// 	return s.Msg
// }

func newErr(msg string, args ...interface{}) *ErrorForCmdr {
	return &ErrorForCmdr{ExtErr: *errors.New(msg, args...)}
}

func newErrTmpl(tmpl string) *ErrorForCmdr {
	return &ErrorForCmdr{ExtErr: *errors.NewTemplate(tmpl)}
}

func (e *ErrorForCmdr) Error() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%v|%s", e.Ignorable, e.ExtErr.Error()))
	return buf.String()
}

// Template setup a string format template.
// Coder could compile the error object with formatting args later.
//
// Note that `ExtErr.Template()` had been overrided here
func (e *ErrorForCmdr) Template(tmpl string) *ErrorForCmdr {
	_ = e.ExtErr.Template(tmpl)
	return e
}

// Format compiles the final msg with string template and args
//
// Note that `ExtErr.Template()` had been overridden here
func (e *ErrorForCmdr) Format(args ...interface{}) *ErrorForCmdr {
	_ = e.ExtErr.Format(args...)
	return e
}

// Msg encodes a formattable msg with args into ErrorForCmdr
//
// Note that `ExtErr.Template()` had been overridden here
func (e *ErrorForCmdr) Msg(msg string, args ...interface{}) *ErrorForCmdr {
	_ = e.ExtErr.Msg(msg, args...)
	return e
}

// Attach attaches the nested errors into ErrorForCmdr
//
// Note that `ExtErr.Template()` had been overridden here
func (e *ErrorForCmdr) Attach(errors ...error) *ErrorForCmdr {
	_ = e.ExtErr.Attach(errors...)
	return e
}

// Nest attaches the nested errors into ErrorForCmdr
//
// Note that `ExtErr.Template()` had been overridden here
func (e *ErrorForCmdr) Nest(errors ...error) *ErrorForCmdr {
	_ = e.ExtErr.Nest(errors...)
	return e
}
