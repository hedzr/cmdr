package builder

import (
	"testing"

	"github.com/hedzr/cmdr/v2/cli"
)

func TestNew(t *testing.T) {
	v := New(nil)
	t.Logf("v = %v", v)
}

func TestAppS_NewFlagBuilder(t *testing.T) {

}

func TestCcb_NewFlagBuilder(t *testing.T) {

}

func TestFfb_NewFlagBuilder(t *testing.T) {
	testNewFlagBuilder(t)
}

func TestCcb_NewCommandBuilder(t *testing.T) {
	testNewCommandBuilder(t)
}

func testNewCommandBuilder(t *testing.T) {
	b := buildable(nil)
	bb := newCommandBuilderShort(b, "help", "h", "info")

	bb.Titles("verbose", "v", "verbose-mode", "non-quiet-mode")
	bb.ExtraShorts("V")
	bb.Description("verbose mode", "verbose mode (long desc)")
	bb.Examples(`$APP --verbose|$APP -v`)
	bb.Group("zzz9.Misc")
	bb.Deprecated("v2.0")
	bb.Hidden(true, true)
	bb.TailPlaceHolders("")
	bb.RedirectTo("help-system.commands")
	bb.OnMatched(nil)
	bb.OnAction(nil)
	bb.OnPreAction(nil)
	bb.OnPostAction(nil)
	bb.PresetCmdLines("")
	bb.InvokeProc("")
	bb.InvokeShell("")
	bb.UseShell("/bin/bash")

	cb := bb.NewCommandBuilder("command", "c", "cc", "cmd")
	cb.UseShell("/bin/zsh")

	fb := bb.NewFlagBuilder("flag", "f", "ff", "flg")
	fb.OnMatched(nil)

	bb.Build()

	bb.AddCmd(func(b cli.CommandBuilder) {
		b.OnMatched(nil)
	})
	bb.AddFlg(func(b cli.FlagBuilder) {
		b.OnMatched(nil)
	})

	cb = bb.Cmd("dash", "d")
	cb.UseShell("/bin/dash")

	fb = bb.Flg("cool", "c")
	fb.OnMatched(nil)
}

func testNewFlagBuilder(t *testing.T) {
	b := buildable(nil)
	bb := newFlagBuilderShort(b, "verbose", "v", "verbose-mode")

	app := buildable(nil)
	bb.SetApp(app)

	bb.Titles("verbose", "v", "verbose-mode", "non-quiet-mode")
	bb.Default(true)
	bb.DefaultValue(false)
	bb.ExtraShorts("V")
	bb.Description("verbose mode", "verbose mode (long desc)")
	bb.Examples(`$APP --verbose|$APP -v`)
	bb.Group("zzz9.Misc")
	bb.Deprecated("v2.0")
	bb.Hidden(true, true)
	bb.ToggleGroup("zzz9.Misc")
	bb.PlaceHolder("")
	bb.EnvVars("VERBOSE")
	bb.AppendEnvVars("V", "VERBOSE_MODE")
	bb.ExternalEditor("EDITOR")
	bb.ExternalEditor("")
	bb.ValidArgs("")
	bb.AppendValidArgs("", "")
	bb.Range(0, 0)
	bb.HeadLike(false, 0, 100)
	bb.Required(false)
	bb.CompJustOnce(false)
	bb.CompActionStr("")
	bb.CompMutualExclusives("quiet")
	bb.CompPrerequisites("")
	bb.CompCircuitBreak(false)
	bb.DoubleTildeOnly(false)
	bb.OnParseValue(nil)
	bb.OnMatched(nil)
	bb.OnChanging(nil)
	bb.OnChanged(nil)
	bb.OnSet(nil)

	bb.Build()
}
