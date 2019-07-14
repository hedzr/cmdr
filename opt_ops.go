/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

type (
	// RootCmdOpt for fluent api
	RootCmdOpt struct {
		optCommandImpl
	}

	// CmdOpt for fluent api
	CmdOpt struct {
		optCommandImpl
	}

	// SubCmdOpt for fluent api
	SubCmdOpt struct {
		optCommandImpl
	}

	// StringOpt for fluent api
	StringOpt struct {
		optFlagImpl
	}

	// StringSliceOpt for fluent api
	StringSliceOpt struct {
		optFlagImpl
	}

	// IntSliceOpt for fluent api
	IntSliceOpt struct {
		optFlagImpl
	}

	// BoolOpt for fluent api
	BoolOpt struct {
		optFlagImpl
	}

	// IntOpt for fluent api
	IntOpt struct {
		optFlagImpl
	}

	// UintOpt for fluent api
	UintOpt struct {
		optFlagImpl
	}

	// Int64Opt for fluent api
	Int64Opt struct {
		optFlagImpl
	}

	// Uint64Opt for fluent api
	Uint64Opt struct {
		optFlagImpl
	}

	// Float32Opt for fluent api
	Float32Opt struct {
		optFlagImpl
	}

	// Float64Opt for fluent api
	Float64Opt struct {
		optFlagImpl
	}

	// DurationOpt for fluent api
	DurationOpt struct {
		optFlagImpl
	}
)

// Header for fluent api
func (s *RootCmdOpt) Header(header string) *RootCmdOpt {
	optCtx.root.Header = header
	return s
}

// Copyright for fluent api
func (s *RootCmdOpt) Copyright(copyright, author string) *RootCmdOpt {
	optCtx.root.Copyright = copyright
	optCtx.root.Author = author
	return s
}

// func (s *RootCmdOpt) Command(cmdOpt *CmdOpt) *RootCmdOpt {
// 	optCtx.root.Command = *cmdOpt.workingFlag
// 	return s
// }

// func (s *RootCmdOpt) SubCmd() (opt OptCmd) {
// 	cmd := &Command{}
// 	optCtx.root.SubCommands = append(optCtx.root.SubCommands, cmd)
// 	optCtx.current = cmd
// 	return &SubCmdOpt{optCommandImpl: optCommandImpl{workingFlag: cmd},}
// }

// NewCmdFrom for fluent api
func NewCmdFrom(cmd *Command) (opt OptCmd) {
	optCtx.current = cmd
	return &CmdOpt{optCommandImpl: optCommandImpl{working: optCtx.current}}
}

// NewCmd for fluent api
func NewCmd() (opt OptCmd) {
	optCtx.current = &Command{}
	return &CmdOpt{optCommandImpl: optCommandImpl{working: optCtx.current}}
}

// NewSubCmd for fluent api
func NewSubCmd() (opt OptCmd) {
	cmd := &Command{}
	optCtx.current.SubCommands = append(optCtx.current.SubCommands, cmd)
	return &SubCmdOpt{optCommandImpl: optCommandImpl{working: cmd}}
}

// NewBool for fluent api
func NewBool() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &BoolOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// NewString for fluent api
func NewString() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &StringOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// NewStringSlice for fluent api
func NewStringSlice() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &StringSliceOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// NewIntSlice for fluent api
func NewIntSlice() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &IntSliceOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// NewInt for fluent api
func NewInt() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &IntOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// NewUint for fluent api
func NewUint() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &UintOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// NewInt64 for fluent api
func NewInt64() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &Int64Opt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// NewUint64 for fluent api
func NewUint64() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &Uint64Opt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// NewFloat32 for fluent api
func NewFloat32() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &Float32Opt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// NewFloat64 for fluent api
func NewFloat64() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &Float64Opt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// NewDuration for fluent api
func NewDuration() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &DurationOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}
