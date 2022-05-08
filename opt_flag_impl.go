/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"regexp"
	"strings"

	"github.com/hedzr/cmdr/tool"
)

type optFlagImpl struct {
	working *Flag
	parent  OptCmd
}

func (s *optFlagImpl) ToFlag() *Flag {
	return s.working
}

func (s *optFlagImpl) AttachTo(parent OptCmd) (opt OptFlag) {
	if s != nil && s.working != nil && parent != nil {
		parent.AddOptFlag(s)
		s.parent = parent
	}
	return s
}

func (s *optFlagImpl) AttachToCommand(cmd *Command) (opt OptFlag) {
	if cmd != nil {
		f := s.ToFlag()
		f.owner = cmd
		cmd.Flags = uniAddFlg(cmd.Flags, f)
		s.parent = NewCmdFrom(cmd)
	}
	return s
}

func (s *optFlagImpl) AttachToRoot(root *RootCommand) (opt OptFlag) {
	if root != nil {
		f := s.ToFlag()
		f.owner = &root.Command
		root.Command.Flags = uniAddFlg(root.Command.Flags, f)
		s.parent = RootFrom(root)
	}
	return s
}

func (s *optFlagImpl) Titles(long, short string, aliases ...string) (opt OptFlag) {
	s.working.Short = short
	s.working.Full = long
	if tool.HasOrderPrefix(long) {
		s.working.Full = tool.StripOrderPrefix(long)
		s.working.Name = long
	}
	s.working.Aliases = uniAddStrs(s.working.Aliases, aliases...)
	return s
}

func (s *optFlagImpl) Short(short string) (opt OptFlag) {
	s.working.Short = short
	return s
}

func (s *optFlagImpl) Long(long string) (opt OptFlag) {
	s.working.Full = long
	if tool.HasOrderPrefix(long) {
		s.working.Full = tool.StripOrderPrefix(long)
		s.working.Name = long
	}
	return s
}

func (s *optFlagImpl) Name(name string) (opt OptFlag) {
	s.working.Name = name
	return s
}

func (s *optFlagImpl) Aliases(aliases ...string) (opt OptFlag) {
	s.working.Aliases = uniAddStrs(s.working.Aliases, aliases...)
	return s
}

func (s *optFlagImpl) Description(oneLineDesc string, longDesc ...string) (opt OptFlag) {
	s.working.Description = oneLineDesc

	for _, long := range longDesc {
		s.working.LongDescription = long

		if s.working.Description == "" {
			s.working.Description = long
		}
	}

	if b := regexp.MustCompile("`(.+)`").Find([]byte(s.working.Description)); len(b) > 2 {
		ph := strings.ToUpper(strings.Trim(string(b), "`"))
		s.Placeholder(ph)
	}

	return s
}

func (s *optFlagImpl) Examples(examples string) (opt OptFlag) {
	s.working.Examples = examples
	return s
}

func (s *optFlagImpl) Group(group string) (opt OptFlag) {
	s.working.Group = group
	return s
}

func (s *optFlagImpl) Hidden(hidden bool) (opt OptFlag) {
	s.working.Hidden = hidden
	return s
}

func (s *optFlagImpl) VendorHidden(hidden bool) (opt OptFlag) {
	s.working.VendorHidden = hidden
	return s
}

func (s *optFlagImpl) Deprecated(deprecation string) (opt OptFlag) {
	s.working.Deprecated = deprecation
	return s
}

func (s *optFlagImpl) Action(action Handler) (opt OptFlag) {
	s.working.Action = action
	return s
}

func (s *optFlagImpl) ToggleGroup(group string) (opt OptFlag) {
	s.working.ToggleGroup = group
	return s
}

func (s *optFlagImpl) DefaultValue(val interface{}, placeholder string) (opt OptFlag) {
	s.working.DefaultValue = val
	s.working.DefaultValuePlaceholder = placeholder
	return s
}

// Placeholder to specify the text string that will be appended
// to the end of a flag expr; it is used into help screen.
// For example, `Placeholder("PASSWORD")` will take the form like:
//   -p, --password=PASSWORD, --pwd, --passwd    to input password
func (s *optFlagImpl) Placeholder(placeholder string) (opt OptFlag) {
	s.working.DefaultValuePlaceholder = placeholder
	return s
}

func (s *optFlagImpl) CompletionActionStr(str string) (opt OptFlag) {
	s.working.actionStr = str
	return s
}

func (s *optFlagImpl) CompletionMutualExclusiveFlags(flags ...string) (opt OptFlag) {
	s.working.mutualExclusives = append(s.working.mutualExclusives, flags...)
	return s
}

func (s *optFlagImpl) CompletionPrerequisitesFlags(flags ...string) (opt OptFlag) {
	s.working.prerequisites = append(s.working.prerequisites, flags...)
	return s
}

func (s *optFlagImpl) CompletionJustOnce(once bool) (opt OptFlag) {
	s.working.justOnce = once
	return s
}

func (s *optFlagImpl) CompletionCircuitBreak(ccb bool) (opt OptFlag) {
	s.working.circuitBreak = ccb
	return s
}

func (s *optFlagImpl) DoubleTildeOnly(dto bool) (opt OptFlag) {
	s.working.dblTildeOnly = dto
	return s
}

func (s *optFlagImpl) ExternalTool(etName string) (opt OptFlag) {
	s.working.ExternalTool = etName
	return s
}

func (s *optFlagImpl) ValidArgs(list ...string) (opt OptFlag) {
	s.working.ValidArgs = list
	return s
}

func (s *optFlagImpl) HeadLike(enable bool, min, max int64) (opt OptFlag) {
	s.working.HeadLike = enable
	s.working.Min, s.working.Max = min, max
	return s
}

func (s *optFlagImpl) EnvKeys(keys ...string) (opt OptFlag) {
	s.working.EnvVars = uniAddStrs(s.working.EnvVars, keys...)
	return s
}

func (s *optFlagImpl) Required(required ...bool) (opt OptFlag) {
	b := true
	for _, bb := range required {
		b = bb
	}
	s.working.Required = b
	return s
}

func (s *optFlagImpl) OnSet(f func(keyPath string, value interface{})) (opt OptFlag) {
	s.working.onSet = f
	return s
}

func (s *optFlagImpl) SetOwner(opt OptCmd) {
	s.parent = opt
	if s.working != nil && opt != nil {
		s.working.owner = opt.ToCommand()
	} else if s.working != nil {
		s.working.owner = nil
	}
}

func (s *optFlagImpl) OwnerCommand() (opt OptCmd) {
	opt = s.parent
	return
}

func (s *optFlagImpl) RootCommand() (root *RootCommand) {
	root = optCtx.root
	return
}
