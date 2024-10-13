package main

import (
	"context"
	"os"
	"path"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/pkg/dir"
	"github.com/hedzr/cmdr/v2/pkg/exec"
	"github.com/hedzr/is/term/color"
)

func onEvalJumpSubCommands(ctx context.Context, c cli.BaseOptI) (it cli.EvalIterator, err error) {
	files := make(map[string]*liteCmdS)
	pos := 0
	var keys []string

	baseDir := path.Join("ci", "pkg", "usr.local.lib.tiny-app", "ext")
	err = dir.ForFile(baseDir, func(depth int, dirName string, fi os.DirEntry) (stop bool, err error) {
		if fi.Name()[0] == '.' {
			return
		}
		key := path.Join(dirName, fi.Name())
		files[key] = &liteCmdS{dirName: dirName, fi: fi, depth: depth, owner: c}
		keys = append(keys, key)
		return
	})

	it = func(context.Context) (bo cli.BaseOptI, hasNext bool, err error) {
		if pos < len(keys) {
			key := keys[pos]
			bo = files[key]
			pos++
			hasNext = pos < len(keys)
		}
		return
	}
	return
}

type liteCmdS struct {
	dirName string
	fi      os.DirEntry
	depth   int
	owner   cli.BaseOptI
	group   string

	hitTitle string
	hitTimes int
}

func (s *liteCmdS) name() string { return s.fi.Name() }

func (s *liteCmdS) OwnerOrParent() cli.BacktraceableMin { return s.owner.(cli.Backtraceable) }
func (s *liteCmdS) OwnerIsNil() bool                    { return s.owner == nil }
func (s *liteCmdS) GetDottedPath() string               { return cli.DottedPath(s) }
func (s *liteCmdS) GetTitleName() string                { return s.name() }
func (s *liteCmdS) GetTitleNamesArray() []string        { return []string{s.name()} }
func (s *liteCmdS) GetTitleNames() string               { return s.name() }
func (s *liteCmdS) App() cli.App                        { return nil }

func (s *liteCmdS) String() string { return path.Join(s.dirName, s.name()) }

func (s *liteCmdS) FindSubCommand(ctx context.Context, longName string, wide bool) (res cli.BaseOptI) {
	return
}
func (s *liteCmdS) Walk(ctx context.Context, cb cli.WalkCB) {
	return
}
func (s *liteCmdS) WalkGrouped(ctx context.Context, cb cli.WalkGroupedCB) {
	return
}
func (s *liteCmdS) WalkBackwardsCtx(ctx context.Context, cb cli.WalkBackwardsCB, pc *cli.WalkBackwardsCtx) {
	return
}

func (s *liteCmdS) OwnerCmd() cli.BaseOptI   { return s.owner }
func (s *liteCmdS) Root() *cli.RootCommand   { return s.owner.Root() }
func (s *liteCmdS) Name() string             { return s.name() }
func (s *liteCmdS) ShortName() string        { return s.name() }
func (s *liteCmdS) ShortNames() []string     { return []string{s.name()} }
func (s *liteCmdS) AliasNames() []string     { return nil }
func (s *liteCmdS) Desc() string             { return s.String() }
func (s *liteCmdS) DescLong() string         { return "" }
func (s *liteCmdS) Examples() string         { return "" }
func (s *liteCmdS) TailPlaceHolder() string  { return "" }
func (s *liteCmdS) GetCommandTitles() string { return s.name() }

func (s *liteCmdS) GroupTitle() string     { return cmdr.RemoveOrderedPrefix(s.SafeGroup()) }
func (s *liteCmdS) GroupHelpTitle() string { return s.GroupTitle() }
func (s *liteCmdS) SafeGroup() string {
	if s.group == "" {
		return cli.UnsortedGroup
	}
	return s.group
}
func (s *liteCmdS) AllGroupKeys(chooseFlag, sort bool) []string             { return nil }
func (s *liteCmdS) Hidden() bool                                            { return false }
func (s *liteCmdS) VendorHidden() bool                                      { return false }
func (s *liteCmdS) Deprecated() string                                      { return "" }
func (s *liteCmdS) CountOfCommands() int                                    { return 0 }
func (s *liteCmdS) CommandsInGroup(groupTitle string) (list []cli.BaseOptI) { return nil }
func (s *liteCmdS) FlagsInGroup(groupTitle string) (list []*cli.Flag)       { return nil }
func (s *liteCmdS) Flags() []*cli.Flag                                      { return nil }
func (s *liteCmdS) SubCommands() []*cli.Command                             { return nil }
func (s *liteCmdS) Invoke(ctx context.Context, args []string) (err error) {
	err = exec.Run("sh", "-c", s.fi.Name())
	return
}
func (s *liteCmdS) OnEvalSubcommands() cli.OnEvaluateSubCommands {
	return nil
}
func (s *liteCmdS) OnEvalSubcommandsOnce() cli.OnEvaluateSubCommands {
	return nil
}
func (s *liteCmdS) OnEvalSubcommandsOnceInvoked() bool {
	return false
}
func (s *liteCmdS) OnEvalSubcommandsOnceCache() []cli.BaseOptI {
	return nil
}
func (s *liteCmdS) OnEvalSubcommandsOnceSetCache(list []cli.BaseOptI) {
}

func (s *liteCmdS) OnEvalFlags() cli.OnEvaluateFlags {
	return nil
}
func (s *liteCmdS) OnEvalFlagsOnce() cli.OnEvaluateFlags {
	return nil
}
func (s *liteCmdS) OnEvalFlagsOnceInvoked() bool {
	return false
}
func (s *liteCmdS) OnEvalFlagsOnceCache() []*cli.Flag {
	return nil
}
func (s *liteCmdS) OnEvalFlagsOnceSetCache(list []*cli.Flag) {
}

func (s *liteCmdS) DeprecatedHelpString(trans func(ss string, clr color.Color) string, clr, clrDefault color.Color) (hs, plain string) {
	return
}

func (s *liteCmdS) SetHitTitle(title string) {
	s.hitTitle = title
	s.hitTimes++
}
func (s *liteCmdS) HitTitle() string { return s.hitTitle }
func (s *liteCmdS) HitTimes() int    { return s.hitTimes }
func (s *liteCmdS) ForeachFlags(func(f *cli.Flag) (stop bool)) (stop bool) {
	return
}
