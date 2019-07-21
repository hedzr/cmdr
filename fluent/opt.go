/*
 * Copyright © 2019 Hedzr Yeh.
 */

package fluent

type (
	// Value is the interface to the dynamic value stored in a flag.
	// (The default value is represented as a string.)
	//
	// If a Value has an IsBoolFlag() bool method returning true,
	// the command-line parser makes -name equivalent to -name=true
	// rather than using the next command-line argument.
	//
	// Set is called once, in command line order, for each flag present.
	// The flag package may call the String method with a zero-valued receiver,
	// such as a nil pointer.
	Value interface {
		String() string
		Set(string) error
	}

	// Getter is an interface that allows the contents of a Value to be retrieved.
	// It wraps the Value interface, rather than being part of it, because it
	// appeared after Go 1 and its compatibility rules. All Value types provided
	// by this package satisfy the Getter interface.
	Getter interface {
		Value
		Get() interface{}
	}

	// ErrorHandling defines how FlagSet.Parse behaves if the parse fails.
	ErrorHandling int
)
