// Copyright © 2020 Hedzr Yeh.

// +build go113

package cmdr

import "reflect"

// IsZero reports whether v is the zero value for its type.
// It panics if the argument is invalid.
func IsZero(from reflect.Value) bool {
	return from.IsZero()
}
