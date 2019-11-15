/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

type optFlagImpl struct {
	working *Flag
	parent  OptCmd
}

func (s *optFlagImpl) ToFlag() *Flag {
	return s.working
}

func (s *optFlagImpl) AttachTo(opt OptCmd) {
	opt.AddOptFlag(s)
}

func (s *optFlagImpl) AttachToCommand(cmd *Command) {
	cmd.Flags = append(cmd.Flags, s.ToFlag())

}

func (s *optFlagImpl) AttachToRoot(root *RootCommand) {
	root.Command.Flags = append(root.Command.Flags, s.ToFlag())
}

func (s *optFlagImpl) Titles(short, long string, aliases ...string) (opt OptFlag) {
	s.working.Short = short
	s.working.Full = long
	s.working.Aliases = append(s.working.Aliases, aliases...)
	opt = s
	return
}

func (s *optFlagImpl) Short(short string) (opt OptFlag) {
	s.working.Short = short
	opt = s
	return
}

func (s *optFlagImpl) Long(long string) (opt OptFlag) {
	s.working.Full = long
	opt = s
	return
}

func (s *optFlagImpl) Aliases(aliases ...string) (opt OptFlag) {
	s.working.Aliases = append(s.working.Aliases, aliases...)
	opt = s
	return
}

func (s *optFlagImpl) Description(oneLine, long string) (opt OptFlag) {
	s.working.Description = oneLine
	s.working.LongDescription = long
	opt = s
	return
}

func (s *optFlagImpl) Examples(examples string) (opt OptFlag) {
	s.working.Examples = examples
	opt = s
	return
}

func (s *optFlagImpl) Group(group string) (opt OptFlag) {
	s.working.Group = group
	opt = s
	return
}

func (s *optFlagImpl) Hidden(hidden bool) (opt OptFlag) {
	s.working.Hidden = hidden
	opt = s
	return
}

func (s *optFlagImpl) Deprecated(deprecation string) (opt OptFlag) {
	s.working.Deprecated = deprecation
	opt = s
	return
}

func (s *optFlagImpl) Action(action func(cmd *Command, args []string) (err error)) (opt OptFlag) {
	s.working.Action = action
	opt = s
	return
}

func (s *optFlagImpl) ToggleGroup(group string) (opt OptFlag) {
	s.working.ToggleGroup = group
	opt = s
	return
}

func (s *optFlagImpl) DefaultValue(val interface{}, placeholder string) (opt OptFlag) {
	s.working.DefaultValue = val
	s.working.DefaultValuePlaceholder = placeholder
	opt = s
	return
}

func (s *optFlagImpl) ExternalTool(envKeyName string) (opt OptFlag) {
	s.working.ExternalTool = envKeyName
	opt = s
	return
}

func (s *optFlagImpl) ValidArgs(list ...string) (opt OptFlag) {
	s.working.ValidArgs = list
	opt = s
	return
}

func (s *optFlagImpl) HeadLike(enable bool, min, max int64) (opt OptFlag) {
	s.working.HeadLike = enable
	s.working.Min, s.working.Max = min, max
	opt = s
	return
}

func (s *optFlagImpl) OnSet(f func(keyPath string, value interface{})) (opt OptFlag) {
	s.working.onSet = f
	opt = s
	return
}

func (s *optFlagImpl) SetOwner(opt OptCmd) {
	s.parent = opt
	return
}

func (s *optFlagImpl) OwnerCommand() (opt OptCmd) {
	opt = s.parent
	return
}

func (s *optFlagImpl) RootCommand() (root *RootCommand) {
	root = optCtx.root
	return
}
