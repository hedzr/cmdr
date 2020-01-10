/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import "time"

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

	// complex64Opt for fluent api
	complex64Opt struct {
		optFlagImpl
	}

	// complex128Opt for fluent api
	complex128Opt struct {
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
	optCtx.current.SubCommands = uniAddCmd(optCtx.current.SubCommands, cmd)
	return &subCmdOpt{optCommandImpl: optCommandImpl{working: cmd}}
}

// NewBool for fluent api
func NewBool(defaultValue bool) (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &boolOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
	opt.DefaultValue(defaultValue, "")
	return
}

func uniAddCmd(cmds []*Command, cmd *Command) []*Command {
	for _, f := range cmds {
		if f == cmd {
			return cmds
		}
	}
	return append(cmds, cmd)
}

func uniAddFlg(flags []*Flag, flg *Flag) []*Flag {
	for _, f := range flags {
		if f == flg {
			return flags
		}
	}
	return append(flags, flg)
}

func uniAddStr(a []string, s string) []string {
	for _, f := range a {
		if f == s {
			return a
		}
	}
	return append(a, s)
}

func uniAddStrs(a []string, ss ...string) []string {
	for _, s := range ss {
		found := false
		for _, f := range a {
			if f == s {
				found = true
				break
			}
		}
		if !found {
			a = append(a, s)
		}
	}
	return a
}

// NewString for fluent api
func NewString(defaultValue string) (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &stringOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
	opt.DefaultValue(defaultValue, "")
	return
}

// NewStringSlice for fluent api
func NewStringSlice(defaultValue []string) (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &stringSliceOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
	opt.DefaultValue(defaultValue, "")
	return
}

// NewIntSlice for fluent api
func NewIntSlice(defaultValue []int) (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &intSliceOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
	opt.DefaultValue(defaultValue, "")
	return
}

// NewInt for fluent api
func NewInt(defaultValue int) (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &intOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
	opt.DefaultValue(defaultValue, "")
	return
}

// NewUint for fluent api
func NewUint(defaultValue uint) (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &uintOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
	opt.DefaultValue(defaultValue, "")
	return
}

// NewInt64 for fluent api
func NewInt64(defaultValue int64) (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &int64Opt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
	opt.DefaultValue(defaultValue, "")
	return
}

// NewUint64 for fluent api
func NewUint64(defaultValue uint64) (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &uint64Opt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
	opt.DefaultValue(defaultValue, "")
	return
}

// NewFloat32 for fluent api
func NewFloat32(defaultValue float32) (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &float32Opt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
	opt.DefaultValue(defaultValue, "")
	return
}

// NewFloat64 for fluent api
func NewFloat64(defaultValue float64) (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &float64Opt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
	opt.DefaultValue(defaultValue, "")
	return
}

// NewComplex64 for fluent api
func NewComplex64(defaultValue complex64) (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &complex64Opt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
	opt.DefaultValue(defaultValue, "")
	return
}

// NewComplex128 for fluent api
func NewComplex128(defaultValue complex128) (opt OptFlag) {
	optCtx.workingFlag = &Flag{}
	optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &complex128Opt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
	opt.DefaultValue(defaultValue, "")
	return
}

// NewDuration for fluent api
func NewDuration(defaultValue time.Duration) (opt OptFlag) {
	opt = NewDurationFrom(&Flag{})
	opt.DefaultValue(defaultValue, "")
	return
}

// NewDurationFrom for fluent api
func NewDurationFrom(flg *Flag) (opt OptFlag) {
	optCtx.workingFlag = flg
	optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	return &durationOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}
