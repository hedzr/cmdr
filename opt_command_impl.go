/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"reflect"
	"time"
)

type optCommandImpl struct {
	working *Command
	parent  OptCmd
}

func (s *optCommandImpl) ToCommand() *Command {
	return s.working
}

func (s *optCommandImpl) AddOptFlag(flag OptFlag) {
	s.working.Flags = append(s.working.Flags, flag.ToFlag())
}

func (s *optCommandImpl) AddFlag(flag *Flag) {
	s.working.Flags = append(s.working.Flags, flag)
}

func (s *optCommandImpl) AddOptCmd(opt OptCmd) {
	cmd := opt.ToCommand()

	// optCtx.current = cmd

	s.working.SubCommands = append(s.working.SubCommands, cmd)

	// opt = &subCmdOpt{optCommandImpl: optCommandImpl{working: cmd, parent: s}}
}

func (s *optCommandImpl) AddCommand(cmd *Command) {
	s.working.SubCommands = append(s.working.SubCommands, cmd)
}

func (s *optCommandImpl) AttachTo(opt OptCmd) {
	opt.AddOptCmd(s)
}

func (s *optCommandImpl) AttachToCommand(cmd *Command) {
	cmd.SubCommands = append(cmd.SubCommands, s.working)
}

func (s *optCommandImpl) AttachToRoot(root *RootCommand) {
	root.SubCommands = append(root.SubCommands, s.working)
}

func (s *optCommandImpl) Titles(short, long string, aliases ...string) (opt OptCmd) {
	s.working.Short = short
	s.working.Full = long
	s.working.Aliases = append(s.working.Aliases, aliases...)
	opt = s
	return
}

func (s *optCommandImpl) Short(short string) (opt OptCmd) {
	s.working.Short = short
	opt = s
	return
}

func (s *optCommandImpl) Long(long string) (opt OptCmd) {
	s.working.Full = long
	opt = s
	return
}

func (s *optCommandImpl) Aliases(aliases ...string) (opt OptCmd) {
	s.working.Aliases = append(s.working.Aliases, aliases...)
	opt = s
	return
}

func (s *optCommandImpl) Description(oneLine, long string) (opt OptCmd) {
	s.working.Description = oneLine
	s.working.LongDescription = long
	opt = s
	return
}

func (s *optCommandImpl) Examples(examples string) (opt OptCmd) {
	s.working.Examples = examples
	opt = s
	return
}

func (s *optCommandImpl) Group(group string) (opt OptCmd) {
	s.working.Group = group
	opt = s
	return
}

func (s *optCommandImpl) Hidden(hidden bool) (opt OptCmd) {
	s.working.Hidden = hidden
	opt = s
	return
}

func (s *optCommandImpl) Deprecated(deprecation string) (opt OptCmd) {
	s.working.Deprecated = deprecation
	opt = s
	return
}

func (s *optCommandImpl) Action(action func(cmd *Command, args []string) (err error)) (opt OptCmd) {
	s.working.Action = action
	opt = s
	return
}

func (s *optCommandImpl) PreAction(pre func(cmd *Command, args []string) (err error)) (opt OptCmd) {
	// s.workingFlag.ExternalTool = envKeyName
	s.working.PreAction = pre
	opt = s
	return
}

func (s *optCommandImpl) PostAction(post func(cmd *Command, args []string)) (opt OptCmd) {
	// s.workingFlag.ExternalTool = envKeyName
	s.working.PostAction = post
	opt = s
	return
}

func (s *optCommandImpl) TailPlaceholder(placeholder string) (opt OptCmd) {
	// s.workingFlag.ExternalTool = envKeyName
	s.working.TailPlaceHolder = placeholder
	opt = s
	return
}

func (s *optCommandImpl) Bool() (opt OptFlag) {
	flg := &Flag{}
	s.working.Flags = append(s.working.Flags, flg)
	return &boolOpt{optFlagImpl: optFlagImpl{working: flg, parent: s}}
}

func (s *optCommandImpl) String() (opt OptFlag) {
	flg := &Flag{}
	s.working.Flags = append(s.working.Flags, flg)
	return &stringOpt{optFlagImpl: optFlagImpl{working: flg, parent: s}}
}

func (s *optCommandImpl) StringSlice() (opt OptFlag) {
	flg := &Flag{}
	s.working.Flags = append(s.working.Flags, flg)
	return &stringSliceOpt{optFlagImpl: optFlagImpl{working: flg, parent: s}}
}

func (s *optCommandImpl) IntSlice() (opt OptFlag) {
	flg := &Flag{}
	s.working.Flags = append(s.working.Flags, flg)
	return &intSliceOpt{optFlagImpl: optFlagImpl{working: flg, parent: s}}
}

func (s *optCommandImpl) Int() (opt OptFlag) {
	flg := &Flag{}
	s.working.Flags = append(s.working.Flags, flg)
	return &intOpt{optFlagImpl: optFlagImpl{working: flg, parent: s}}
}

func (s *optCommandImpl) Uint() (opt OptFlag) {
	flg := &Flag{}
	s.working.Flags = append(s.working.Flags, flg)
	return &uintOpt{optFlagImpl: optFlagImpl{working: flg, parent: s}}
}

func (s *optCommandImpl) Int64() (opt OptFlag) {
	flg := &Flag{}
	s.working.Flags = append(s.working.Flags, flg)
	return &int64Opt{optFlagImpl: optFlagImpl{working: flg, parent: s}}
}

func (s *optCommandImpl) Uint64() (opt OptFlag) {
	flg := &Flag{}
	s.working.Flags = append(s.working.Flags, flg)
	return &uint64Opt{optFlagImpl: optFlagImpl{working: flg, parent: s}}
}

func (s *optCommandImpl) Float32() (opt OptFlag) {
	flg := &Flag{}
	s.working.Flags = append(s.working.Flags, flg)
	return &float32Opt{optFlagImpl: optFlagImpl{working: flg, parent: s}}
}

func (s *optCommandImpl) Float64() (opt OptFlag) {
	flg := &Flag{}
	s.working.Flags = append(s.working.Flags, flg)
	return &float64Opt{optFlagImpl: optFlagImpl{working: flg, parent: s}}
}

func (s *optCommandImpl) Duration() (opt OptFlag) {
	flg := &Flag{}
	s.working.Flags = append(s.working.Flags, flg)
	return &durationOpt{optFlagImpl: optFlagImpl{working: flg, parent: s}}
}

func (s *optCommandImpl) NewFlag(typ OptFlagType) (opt OptFlag) {
	var flg OptFlag

	switch typ {
	case OptFlagTypeInt:
		flg = s.Int()
	case OptFlagTypeUint:
		flg = s.Uint()
	case OptFlagTypeInt64:
		flg = s.Int64()
	case OptFlagTypeUint64:
		flg = s.Uint64()
	case OptFlagTypeString:
		flg = s.String()
	case OptFlagTypeStringSlice:
		flg = s.StringSlice()
	case OptFlagTypeIntSlice:
		flg = s.IntSlice()
	case OptFlagTypeFloat32:
		flg = s.Float32()
	case OptFlagTypeFloat64:
		flg = s.Float64()
	case OptFlagTypeDuration:
		flg = s.Duration()
	default:
		flg = s.Bool()
	}

	flg.SetOwner(s)

	opt = flg
	return
}

func (s *optCommandImpl) newFlagVC(vv reflect.Type, defaultValue interface{}) (flg OptFlag) {
	switch vv.Kind() {
	case reflect.Int, reflect.Int16, reflect.Int32:
		if _, ok := defaultValue.(time.Duration); ok {
			flg = s.Duration()
		} else {
			flg = s.Int()
		}
	case reflect.Uint, reflect.Uint16, reflect.Uint32:
		flg = s.Uint()
	case reflect.Int64:
		flg = s.Int64()
	case reflect.Uint64:
		flg = s.Uint64()
	case reflect.String:
		flg = s.String()
	case reflect.Slice:
		var elt = vv.Elem()
		switch elt.Kind() {
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			flg = s.IntSlice()
		case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			// flg = s.UintSlice()
			flg = s.IntSlice()
		case reflect.String:
			flg = s.StringSlice()
		}
	case reflect.Float32:
		flg = s.Float32()
	case reflect.Float64:
		flg = s.Float64()
	default:
		flg = s.Bool()
	}
	return
}

func (s *optCommandImpl) NewFlagV(defaultValue interface{}) (opt OptFlag) {
	var flg OptFlag
	var vv = reflect.TypeOf(defaultValue)
	flg = s.newFlagVC(vv, defaultValue)
	if flg != nil {
		flg.DefaultValue(defaultValue, "")
		flg.SetOwner(s)
	}
	opt = flg
	return
}

func (s *optCommandImpl) NewSubCommand() (opt OptCmd) {
	cmd := &Command{root: uniqueWorker.rootCommand}

	optCtx.current = cmd

	s.working.SubCommands = append(s.working.SubCommands, cmd)

	opt = &subCmdOpt{optCommandImpl: optCommandImpl{working: cmd, parent: s}}
	return
}

func (s *optCommandImpl) OwnerCommand() (opt OptCmd) {
	opt = s.parent
	return
}

func (s *optCommandImpl) SetOwner(opt OptCmd) {
	s.parent = opt
	return
}

func (s *optCommandImpl) RootCommand() (root *RootCommand) {
	root = optCtx.root
	return
}
