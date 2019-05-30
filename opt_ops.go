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

// Cmd for fluent api
func Cmd() (opt OptCmd) {
	optCtx.current = &Command{}
	return &CmdOpt{optCommandImpl: optCommandImpl{working: optCtx.current}}
}

// SubCmd for fluent api
func SubCmd() (opt OptCmd) {
	cmd := &Command{}
	optCtx.current.SubCommands = append(optCtx.current.SubCommands, cmd)
	return &SubCmdOpt{optCommandImpl: optCommandImpl{working: cmd}}
}

// Bool for fluent api
func Bool() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &BoolOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// String for fluent api
func String() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &StringOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// StringSlice for fluent api
func StringSlice() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &StringSliceOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// IntSlice for fluent api
func IntSlice() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &IntSliceOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// Int for fluent api
func Int() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &IntOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// Uint for fluent api
func Uint() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &UintOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// Int64 for fluent api
func Int64() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &Int64Opt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// Uint64 for fluent api
func Uint64() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &Uint64Opt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}
