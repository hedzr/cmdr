package cmdr_test

import (
	"context"
	"testing"

	cmdr "github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"

	"gopkg.in/hedzr/errors.v3"
)

func TestDottedPathToCommandOrFlag(t *testing.T) {
	ctx := context.Background()
	app := cmdr.Create("app", "v1", `author`, `desc`,
		cli.WithArgs("test-app", "--debug"),
	).OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
		cc, ff := cmdr.DottedPathToCommandOrFlag("generate.shell.zsh", nil)
		if cc == nil || cc.GetTitleName() != "shell" || ff.Title() != "zsh" {
			t.Fail()
		}
		cc, ff = cmdr.DottedPathToCommandOrFlag("generate.doc", nil)
		if ff != nil || cc.GetTitleName() != "doc" {
			t.Fail()
		}
		return
	}).
		Build()
	err := app.Run(ctx)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestStoreGetSectionFrom(t *testing.T) {
	ctx := context.Background()
	app := cmdr.Create("app", "v1", `author`, `desc`,
		cli.WithArgs("test-app", "--debug"),
	).OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
		b := cmdr.Store().MustBool("debug")
		println("debug flag: ", b)
		if !b {
			t.Fail()
		}
		println(cmdr.Set().Dump())

		type manS struct {
			Dir  string
			Type int
		}
		type genS struct {
			Manual manS
		}
		var v genS
		err = cmdr.Set().GetSectionFrom("cmd.generate", &v)
		if err != nil {
			t.Fail()
		}
		if v.Manual.Type != 1 {
			t.Fail()
		}
		return
	}).
		Build()
	err := app.Run(ctx)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestTo(t *testing.T) {
	ctx := context.Background()
	app := cmdr.Create("app", "v1", `author`, `desc`,
		cli.WithArgs("test-app", "--debug"),
	).OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
		b := cmdr.Store().MustBool("debug")
		println("debug flag: ", b)
		if !b {
			t.Fail()
		}
		println(cmdr.Set().Dump())

		type manS struct {
			Dir  string
			Type int
		}
		type genS struct {
			Manual manS
		}
		var v genS
		err = cmdr.To("cmd.generate", &v)
		if err != nil {
			t.Fail()
		}
		if v.Manual.Type != 1 {
			t.Fail()
		}
		return
	}).
		Build()
	err := app.Run(ctx)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestExecNoRoot(t *testing.T) {
	if err := cmdr.Exec(nil); !errors.Iss(err, cli.ErrEmptyRootCommand) {
		t.Errorf("Error: %v", err)
	}
}

// func TestExecSimple(t *testing.T) {
// 	if err := cmdr.Exec(testdata.BuildCommands(true)); !errors.Iss(err, cli.ErrEmptyRootCommand) {
// 		t.Errorf("Error: %v", err)
// 	}
// }

func TestGetSet(t *testing.T) {
	ctx := context.Background()
	app := cmdr.Create("app", "v1", `author`, `desc`,
		cli.WithArgs("test-app", "--debug"),
	).
		With(func(app cli.App) {
			app.OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
				b := cmdr.Store().MustBool("debug")
				println("debug flag: ", b)
				if !b {
					t.Fail()
				}
				// println(cmdr.Set().Dump())
				return
			})
		}).
		Build()
	err := app.Run(ctx)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}
