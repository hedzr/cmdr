package dyncmd

import (
	"context"
	"os"
	"path"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/pkg/dir"
	"github.com/hedzr/is/exec"
	"github.com/hedzr/is/term/color"
	"github.com/hedzr/store"
)

// onEvalJumpSubCommands querys shell scripts in EXT directory
// (typically it is `/usr/local/lib/<app-name>/ext/`) and build
// as subcommands dynamically.
//
// In this demo app, we looks up `./ci/pkg/usr.local.lib.large-app/ext`
// with hard-code.
//
// EXT directory: see the [cmdr.UsrLibDir()] for its location.
func onEvalJumpSubCommands(ctx context.Context, c cli.Cmd) (it cli.EvalIterator, err error) {
	files := make(map[string]*liteCmdS)
	pos := 0
	var keys []string

	baseDir := cmdr.UsrLibDir()
	if dir.FileExists(baseDir) {
		baseDir = path.Join(baseDir, "ext")
	} else {
		baseDir = path.Join("ci", "pkg", "usr.local.lib", c.App().Name(), "ext")
	}
	if !dir.FileExists(baseDir) {
		return
	}

	err = dir.ForFile(baseDir, func(depth int, dirName string, fi os.DirEntry) (stop bool, err error) {
		if fi.Name()[0] == '.' {
			return
		}
		key := path.Join(dirName, fi.Name())
		files[key] = &liteCmdS{dirName: dirName, fi: fi, depth: depth, owner: c}
		keys = append(keys, key)
		return
	})

	it = func(context.Context) (bo cli.Cmd, hasNext bool, err error) {
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
	owner   cli.Cmd
	group   string

	hitTitle string
	hitTimes int
}

var _ cli.Cmd = (*liteCmdS)(nil)

// var _ cli.CmdPriv = (*liteCmdS)(nil)

func (s *liteCmdS) name() string { return s.fi.Name() }

func (s *liteCmdS) String() string { return path.Join(s.dirName, s.name()) }

func (s *liteCmdS) GetDottedPath() string        { return cli.DottedPath(s) }
func (s *liteCmdS) GetTitleName() string         { return s.name() }
func (s *liteCmdS) GetTitleNamesArray() []string { return []string{s.name()} }
func (s *liteCmdS) GetTitleNames() string        { return s.name() }

func (s *liteCmdS) App() cli.App       { return nil }
func (s *liteCmdS) Set() store.Store   { return s.Root().App().Store() }
func (s *liteCmdS) Store() store.Store { return cmdr.Store() }

func (s *liteCmdS) OwnerIsValid() bool {
	if s.OwnerIsNotNil() {
		if cx, ok := s.owner.(*liteCmdS); ok {
			return cx != s
		}
	}
	return false
}
func (s *liteCmdS) OwnerIsNil() bool                    { return s.owner == nil }
func (s *liteCmdS) OwnerIsNotNil() bool                 { return s.owner != nil }
func (s *liteCmdS) OwnerCmd() cli.Cmd                   { return s.owner }
func (s *liteCmdS) SetOwnerCmd(c cli.Cmd)               { s.owner = c }
func (s *liteCmdS) Root() *cli.RootCommand              { return s.owner.Root() }
func (s *liteCmdS) SetRoot(*cli.RootCommand)            {}
func (s *liteCmdS) OwnerOrParent() cli.BacktraceableMin { return s.owner.(cli.Backtraceable) }

func (s *liteCmdS) Name() string             { return s.name() }
func (s *liteCmdS) SetName(string)           {}
func (s *liteCmdS) ShortTitle() string       { return s.name() }
func (s *liteCmdS) LongTitle() string        { return s.name() }
func (s *liteCmdS) ShortNames() []string     { return []string{s.name()} }
func (s *liteCmdS) AliasNames() []string     { return nil }
func (s *liteCmdS) Desc() string             { return s.String() }
func (s *liteCmdS) DescLong() string         { return "" }
func (s *liteCmdS) SetDesc(desc string)      {}
func (s *liteCmdS) Examples() string         { return "" }
func (s *liteCmdS) TailPlaceHolder() string  { return "" }
func (s *liteCmdS) GetCommandTitles() string { return s.name() }

func (s *liteCmdS) GroupTitle() string { return cmdr.RemoveOrderedPrefix(s.SafeGroup()) }
func (s *liteCmdS) GroupHelpTitle() string {
	tmp := s.SafeGroup()
	if tmp == cli.UnsortedGroup {
		return ""
	}
	return cmdr.RemoveOrderedPrefix(tmp)
}
func (s *liteCmdS) SafeGroup() string {
	if s.group == "" {
		return cli.UnsortedGroup
	}
	return s.group
}
func (s *liteCmdS) AllGroupKeys(chooseFlag, sort bool) []string { return nil }
func (s *liteCmdS) Hidden() bool                                { return false }
func (s *liteCmdS) VendorHidden() bool                          { return false }
func (s *liteCmdS) Deprecated() string                          { return "" }
func (s *liteCmdS) DeprecatedHelpString(trans func(ss string, clr color.Color) string, clr, clrDefault color.Color) (hs, plain string) {
	return
}

func (s *liteCmdS) CountOfCommands() int                               { return 0 }
func (s *liteCmdS) CommandsInGroup(groupTitle string) (list []cli.Cmd) { return nil }
func (s *liteCmdS) FlagsInGroup(groupTitle string) (list []*cli.Flag)  { return nil }
func (s *liteCmdS) SubCommands() []*cli.CmdS                           { return nil }
func (s *liteCmdS) Flags() []*cli.Flag                                 { return nil }

func (s *liteCmdS) HeadLikeFlag() *cli.Flag   { return nil }
func (s *liteCmdS) SetHeadLikeFlag(*cli.Flag) {}

func (s *liteCmdS) SetHitTitle(title string) {
	s.hitTitle = title
	s.hitTimes++
}
func (s *liteCmdS) HitTitle() string { return s.hitTitle }
func (s *liteCmdS) HitTimes() int    { return s.hitTimes }

func (s *liteCmdS) RedirectTo() (dottedPath string) { return }
func (s *liteCmdS) SetRedirectTo(dottedPath string) {}

func (s *liteCmdS) PresetCmdLines() []string         { return nil }
func (s *liteCmdS) InvokeProc() string               { return "" }
func (s *liteCmdS) InvokeShell() string              { return "" }
func (s *liteCmdS) Shell() string                    { return "" }
func (c *liteCmdS) SetPresetCmdLines(args ...string) {}
func (c *liteCmdS) SetInvokeProc(str string)         {}
func (c *liteCmdS) SetInvokeShell(str string)        {}
func (c *liteCmdS) SetShell(str string)              {}

func (s *liteCmdS) CanInvoke() bool {
	return s.fi.Type().IsRegular()
}

func (s *liteCmdS) Invoke(ctx context.Context, args []string) (err error) {
	fullPath := path.Join(s.dirName, s.name())
	err = exec.Run("sh", "-c", fullPath)
	return
}

func (c *liteCmdS) OnEvaluateSubCommandsFromConfig() string {
	return ""
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
func (s *liteCmdS) OnEvalSubcommandsOnceCache() []cli.Cmd {
	return nil
}
func (s *liteCmdS) OnEvalSubcommandsOnceSetCache(list []cli.Cmd) {
}

func (c *liteCmdS) IsDynamicCommandsLoading() bool { return false }
func (c *liteCmdS) IsDynamicFlagsLoading() bool    { return false }

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

func (s *liteCmdS) findSubCommandIn(ctx context.Context, cc cli.Cmd, children []cli.Cmd, longName string, wide bool) (res cli.Cmd) {
	return
}
func (s *liteCmdS) findFlagIn(ctx context.Context, cc cli.Cmd, children []cli.Cmd, longName string, wide bool) (res *cli.Flag) {
	return
}
func (s *liteCmdS) findFlagBackwardsIn(ctx context.Context, cc cli.Cmd, children []cli.Cmd, longName string) (res *cli.Flag) {
	return
}
func (s *liteCmdS) partialMatchFlag(context.Context, string, bool, bool, map[string]*cli.Flag) (matched, remains string, ff *cli.Flag, err error) {
	return
}

func (s *liteCmdS) Match(ctx context.Context, title string) (short bool, cc cli.Cmd) {
	return
}
func (s *liteCmdS) TryOnMatched(position int, hitState *cli.MatchState) (handled bool, err error) {
	return
}
func (s *liteCmdS) MatchFlag(ctx context.Context, vp *cli.FlagValuePkg) (ff *cli.Flag, err error) { //nolint:revive
	return
}

func (s *liteCmdS) FindSubCommand(ctx context.Context, longName string, wide bool) (res cli.Cmd) {
	return
}
func (s *liteCmdS) FindFlagBackwards(ctx context.Context, longName string) (res *cli.Flag) {
	return
}
func (c *liteCmdS) SubCmdBy(longName string) (res cli.Cmd) { return }
func (c *liteCmdS) FlagBy(longName string) (res *cli.Flag) { return }
func (s *liteCmdS) ForeachFlags(context.Context, func(f *cli.Flag) (stop bool)) (stop bool) {
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
func (s *liteCmdS) WalkEverything(ctx context.Context, cb cli.WalkEverythingCB) {
}
func (s *liteCmdS) WalkFast(ctx context.Context, cb cli.WalkFastCB) (stop bool) { return }

func (s *liteCmdS) DottedPathToCommandOrFlag(dottedPath string) (cc cli.Backtraceable, ff *cli.Flag) {
	return
}
