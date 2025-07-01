package builder

import (
	"context"
	"testing"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/pkg/logz"
	"github.com/hedzr/is"
	"github.com/hedzr/logg/slog"
)

func TestStructBuilder_FromStruct(t *testing.T) {
	// New will initialize appS{} struct and make a new
	// rootCommand object into it.
	var w cli.Runner // an empty dummy runner for testing
	a := New(w).Info("demo-app", "0.3.1").Author("hedzr")
	app := a.(*appS)

	if is.DebuggerAttached() {
		// logz.SetLevel(logz.DebugLevel)
		logz.SetLevel(slog.TraceLevel)
	}

	// FromStruct assumes creating a command system from RootCommand.Cmd
	// since a bracketed longTitle "(...)" passed.
	b := app.FromStruct(R{
		F2: "/tmp/value",
	})
	b.Build()

	assertEqual(t, int32(0), app.inCmd)
	assertEqual(t, int32(0), app.inFlg)

	root := app.root.Cmd.(*cli.CmdS)
	assertEqual(t, "/tmp/value", root.Flags()[1].DefaultValue())
	assertEqual(t, "a-cmd", root.SubCommands()[0].Long)
	assertEqual(t, "b", root.SubCommands()[1].Long)
	assertEqual(t, "c", root.SubCommands()[2].Long)

	acmd := root.SubCommands()[0]
	assertEqual(t, "d", acmd.SubCommands()[0].Long)
	assertEqual(t, "a-cmd", acmd.Long)
	assertEqual(t, "a-short", acmd.Short+"-short")
	assertEqual(t, []string{"a", "a1", "a2"}, acmd.Shorts())
	assertEqual(t, []string{"a1-cmd", "a2-cmd"}, acmd.Aliases)

	dcmd := acmd.SubCommands()[0]
	assertEqual(t, "e", dcmd.SubCommands()[0].Long)
	assertEqual(t, "from-now-on", dcmd.SubCommands()[1].Long)
	ecmd := dcmd.SubCommands()[0]
	// F3 bool `title:"f3" shorts:"ff" alias:"f3ff" desc:"A flag for demo" required:"true"`
	assertEqual(t, "Flg{'a-cmd.d.e.f3'}", ecmd.Flags()[0].String())
	assertEqual(t, "ff", ecmd.Flags()[0].Short)
	assertEqual(t, true, ecmd.Flags()[0].Required())
	assertEqual(t, []string{"f3ff"}, ecmd.Flags()[0].Aliases)
	assertEqual(t, "f4", ecmd.Flags()[1].Long)
	assertEqual(t, "f4", ecmd.Flags()[1].Short)
	fcmd := dcmd.SubCommands()[1]
	assertEqual(t, "f5", fcmd.Flags()[0].Long)
	assertEqual(t, "f6", fcmd.Flags()[1].Long)

	pnt, cmd := b.Parent(), b.Building()
	t.Logf("Parent command: %+v", pnt)   // parent should be nil
	t.Logf("Building command: %+v", cmd) // print the RootCommand.Cmd
}

type A struct {
	D
	F1 int
	F2 string
}
type B struct {
	F2 int
	F3 string
}
type C struct {
	F3 bool
	F4 string
}
type D struct {
	E
	FromNowOn F
	F3        bool
	F4        string
}
type E struct {
	F3 bool `title:"f3" shorts:"ff" alias:"f3ff" desc:"A flag for demo" required:"true"`
	F4 string
}
type F struct {
	F5    uint
	F6    byte
	Files []string `cmdr:"positional"`
}

type R struct {
	b   bool // unexported values ignored
	Int int  `cmdr:"-"` // ignored
	A   `title:"a-cmd" shorts:"a,a1,a2" alias:"a1-cmd,a2-cmd" desc:"A command for demo" required:"true"`
	B   `env:"B"`
	C
	F1 int
	F2 string
}

// a --f1 1 --f2 str
// --a.f1 1 --a.f2 str

func (A) With(cb cli.CommandBuilder) {
	// customize for A command, for instance: fb.ExtraShorts("ff")
	logz.Info(".   - A.With() invoked.", "cmdbuilder", cb)
}
func (A) F1With(fb cli.FlagBuilder) {
	// customize for A.F1 flag, for instance: fb.ExtraShorts("ff")
	logz.Info(".   - A.F1With() invoked.", "flgbuilder", fb)
}

// Action method will be called if end-user type subcmd for it (like `app a d e --f3`).
func (E) Action(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
	logz.Info(".   - E.Action() invoked.", "cmd", cmd, "args", args)
	_, err = cmd.App().DoBuiltinAction(ctx, cli.ActionDefault, stringArrayToAnyArray(args)...)
	return
}

// Action method will be called if end-user type subcmd for it (like `app a d f --f5=7`).
func (s F) Action(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
	(&s).Inc()
	logz.Info(".   - F.Action() invoked.", "cmd", cmd, "args", args, "F5", s.F5)
	_, err = cmd.App().DoBuiltinAction(ctx, cli.ActionDefault, stringArrayToAnyArray(args)...)
	return
}

func (s *F) Inc() {
	s.F5++
}

func stringArrayToAnyArray(args []string) (ret []any) {
	for _, it := range args {
		ret = append(ret, it)
	}
	return
}
