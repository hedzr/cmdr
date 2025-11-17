package common

import (
	"context"
	"fmt"

	"github.com/hedzr/cmdr/v2/cli"
)

func AddNegatableFlag(parent cli.CommandBuilder)           { AddNegatableFlagImpl(parent) }
func AddNegatableFlagWithoutCmd(parent cli.CommandBuilder) { AddNegatableFlagImpl(parent, true) }

func AddNegatableFlagImpl(parent cli.CommandBuilder, dontHandlingParentCmd ...bool) { //nolint:revive
	parent.Flg("warning", "w").
		Description("negatable flag: <code>--no-warning</code> is available", "").
		Group("Negatable").
		LeadingPlusSign(true).
		Negatable(true).
		Default(false).
		Build()

	parent.Flg("warnings", "W").
		Description("gcc-style negatable flag: <code>-Wunused-variable</code> and -Wno-unused-variable", "").
		Group("Negatable").
		Negatable(true, "unused-variable", "unused-parameter", "unused-function", "unused-but-set-variable", "unused-private-field", "unused-label").
		Default(false).
		Build()

	var noHpc bool
	for _, p := range dontHandlingParentCmd {
		if p {
			noHpc = true
		}
	}
	if !noHpc {
		// give root command an action to handle it
		action := func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
			println("Dumping Store ----------\n", cmd.Set().Dump())

			cmd.App().DoBuiltinAction(ctx, cli.ActionDefault)

			cs := cmd.Store()

			v := cs.MustBool("warning")
			w := cs.MustBool("no-warning")
			println("warning flag: ", v)
			println("no-warning flag: ", w)
			wf := cmd.FlagBy("warning")
			if wf != nil && wf.LeadingPlusSign() {
				println("warning flag with leading plus sign: +w FOUND.")
			}

			sw1 := cs.MustBool("warnings.unused-variable")
			sw2 := cs.MustBool("warnings.unused-parameter")
			sw3 := cs.MustBool("warnings.unused-function")
			sw4 := cs.MustBool("warnings.unused-but-set-variable")
			sw5 := cs.MustBool("warnings.unused-private-field")
			sw6 := cs.MustBool("warnings.unused-label")
			fmt.Printf(`

--warnings, -W:
    > TG: %q
    > TG.selected: %q
	unused-variable:		%v, (no-): %v
	unused-parameter:		%v, (no-): %v
	unused-function:		%v, (no-): %v
	unused-but-set-variable:	%v, (no-): %v
	unused-private-field:		%v, (no-): %v
	unused-label:			%v, (no-): %v
`,
				cs.MustString("warnings"),
				cs.MustStringSlice("warnings.selected"),
				sw1, cs.MustBool("warnings.no-unused-variable"),
				sw2, cs.MustBool("warnings.no-unused-parameter"),
				sw3, cs.MustBool("warnings.no-unused-function"),
				sw4, cs.MustBool("warnings.no-unused-but-set-variable"),
				sw5, cs.MustBool("warnings.no-unused-private-field"),
				sw6, cs.MustBool("warnings.no-unused-label"),
			)
			return
		}

		parent.OnAction(action)
	}

	parent.Examples(`Try to use negatable flag,
	  
	  $ $APP --no-warning
	    <code>cmd.Store().MustBool("no-warning")</code> will be '<code>true</code>'.
	  $ $APP --warning
	    <code>cmd.Store().MustBool("warning")</code> will be '<code>true</code>'.
	  $ $APP -Wunused-variable -Wno-unused-parameter
	    1. <code>cmd.Store().MustBool("warnings.unused-variable")</code> will be '<code>true</code>',
	    2. <code>cmd.Store().MustBool("warnings.no-unused-parameter")</code> will be '<code>true</code>'.
	       <code>cmd.Store().MustBool("warnings.unused-parameter")</code> will be '<code>false</code>'.
	  $ $APP -Wno-unused-function -Wunused-but-set-variable
	    1. <code>cmd.Store().MustBool("warnings.unused-function")</code> will be '<code>false</code>',
	       <code>cmd.Store().MustBool("warnings.no-unused-function")</code> will be '<code>true</code>',
	    2. <code>cmd.Store().MustBool("warnings.unused-but-set-variable")</code> will be '<code>true</code>'.
	  $ $APP +warning
	    <code>cmd.Store().MustBool("warning")</code> will be '<code>true</code>'.
	`)
}
