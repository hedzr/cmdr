package builder

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hedzr/cmdr/v2/cli"
)

func TestAppS_AddCmd(t *testing.T) {
	a := &appS{}
	a.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("ask", "a")
	})
	assertEqual(t, a.root.Cmd.(*cli.CmdS).Long, "ask")
	assertEqual(t, a.root.Cmd.(*cli.CmdS).Short, "a")

	a.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("bunny", "b")
	})
	child := a.root.SubCommands()[0]
	assertEqual(t, child.Long, "bunny")
	assertEqual(t, child.Short, "b")
}

func TestAppS_Run(t *testing.T) {
	ctx := context.TODO()

	a := &appS{inCmd: 1}
	err := a.Run(ctx)
	assertEqual(t, err != nil, true)

	a = &appS{inFlg: 2}
	err = a.Run(ctx)
	assertEqual(t, err != nil, true)

	a = &appS{}
	err = a.Run(ctx)
	assertEqual(t, err, cli.ErrEmptyRootCommand)
}

func assertEqual(t testing.TB, expect, actual any, msg ...any) { //nolint:govet,unparam //it's a printf/println dual interface
	if reflect.DeepEqual(expect, actual) {
		return
	}

	var mesg string
	if len(msg) > 0 {
		if format, ok := msg[0].(string); ok {
			mesg = fmt.Sprintf(format, msg[1:]...)
		} else {
			mesg = fmt.Sprint(msg...)
		}
	}

	t.Fatalf("assertEqual failed: %v\n    expect: %v\n    actual: %v\n", mesg, expect, actual)
}

// func assertTrue(t testing.TB, cond bool, msg ...any) { //nolint:revive
// 	if cond {
// 		return
// 	}
//
// 	var mesg string
// 	if len(msg) > 0 {
// 		if format, ok := msg[0].(string); ok {
// 			mesg = fmt.Sprintf(format, msg[1:]...)
// 		} else {
// 			mesg = fmt.Sprint(msg...)
// 		}
// 	}
//
// 	t.Fatalf("assertTrue failed: %s", mesg)
// }
//
// func assertFalse(t testing.TB, cond bool, msg ...any) { //nolint:revive
// 	if !cond {
// 		return
// 	}
//
// 	var mesg string
// 	if len(msg) > 0 {
// 		if format, ok := msg[0].(string); ok {
// 			mesg = fmt.Sprintf(format, msg[1:]...)
// 		} else {
// 			mesg = fmt.Sprint(msg...)
// 		}
// 	}
//
// 	t.Fatalf("assertFalse failed: %s", mesg)
// }
