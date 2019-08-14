/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

type (
	// RootCmdOpt for fluent api
	RootCmdOpt struct {
		optCommandImpl
	}

	// cmdOpt for fluent api
	cmdOpt struct {
		optCommandImpl
	}

	// subCmdOpt for fluent api
	subCmdOpt struct {
		optCommandImpl
	}

	// stringOpt for fluent api
	stringOpt struct {
		optFlagImpl
	}

	// stringSliceOpt for fluent api
	stringSliceOpt struct {
		optFlagImpl
	}

	// intSliceOpt for fluent api
	intSliceOpt struct {
		optFlagImpl
	}

	// boolOpt for fluent api
	boolOpt struct {
		optFlagImpl
	}

	// intOpt for fluent api
	intOpt struct {
		optFlagImpl
	}

	// uintOpt for fluent api
	uintOpt struct {
		optFlagImpl
	}

	// int64Opt for fluent api
	int64Opt struct {
		optFlagImpl
	}

	// uint64Opt for fluent api
	uint64Opt struct {
		optFlagImpl
	}

	// float32Opt for fluent api
	float32Opt struct {
		optFlagImpl
	}

	// float64Opt for fluent api
	float64Opt struct {
		optFlagImpl
	}

	// durationOpt for fluent api
	durationOpt struct {
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

// func (s *RootCmdOpt) Command(cmdOpt *cmdOpt) *RootCmdOpt {
// 	optCtx.root.Command = *cmdOpt.workingFlag
// 	return s
// }

// func (s *RootCmdOpt) SubCmd() (opt OptCmd) {
// 	cmd := &Command{}
// 	optCtx.root.SubCommands = append(optCtx.root.SubCommands, cmd)
// 	optCtx.current = cmd
// 	return &subCmdOpt{optCommandImpl: optCommandImpl{workingFlag: cmd},}
// }

// NewCmdFrom for fluent api
func NewCmdFrom(cmd *Command) (opt OptCmd) {
	optCtx.current = cmd
	return &cmdOpt{optCommandImpl: optCommandImpl{working: optCtx.current}}
}

// NewCmd for fluent api
func NewCmd() (opt OptCmd) {
	optCtx.current = &Command{}
	return &cmdOpt{optCommandImpl: optCommandImpl{working: optCtx.current}}
}

// NewSubCmd for fluent api
func NewSubCmd() (opt OptCmd) {
	cmd := &Command{}
	optCtx.current.SubCommands = append(optCtx.current.SubCommands, cmd)
	return &subCmdOpt{optCommandImpl: optCommandImpl{working: cmd}}
}

// NewBool for fluent api
func NewBool() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &boolOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// NewString for fluent api
func NewString() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &stringOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// NewStringSlice for fluent api
func NewStringSlice() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &stringSliceOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// NewIntSlice for fluent api
func NewIntSlice() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &intSliceOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// NewInt for fluent api
func NewInt() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &intOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// NewUint for fluent api
func NewUint() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &uintOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// NewInt64 for fluent api
func NewInt64() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &int64Opt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// NewUint64 for fluent api
func NewUint64() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &uint64Opt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// NewFloat32 for fluent api
func NewFloat32() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &float32Opt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// NewFloat64 for fluent api
func NewFloat64() (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &float64Opt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}

// NewDuration for fluent api
func NewDuration() (opt OptFlag) {
	opt = NewDurationFrom(&Flag{})
	return
}

// NewDurationFrom for fluent api
func NewDurationFrom(flg *Flag) (opt OptFlag) {
	optCtx.workingFlag = flg
	optCtx.current.Flags = append(optCtx.current.Flags, optCtx.workingFlag)
	return &durationOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}
