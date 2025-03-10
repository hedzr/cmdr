package examples

import (
	"github.com/hedzr/cmdr/v2/cli"
)

func AddExternalEditorFlag(c cli.CommandBuilder) { //nolint:revive
	c.Flg("message", "m", "msg").
		Default("").
		Description("the message requesting.", "").
		Group("External Editor").
		PlaceHolder("MESG").
		ExternalEditor(cli.ExternalToolEditor).
		Build()
}
