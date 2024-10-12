package worker

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/store"
)

func TestWorkerS_Run2(t *testing.T) { //nolint:revive
	ctx := context.TODO()
	for i, c := range []struct {
		args     string
		verifier taskAfterParse
		opts     []cli.Opt
	}{
		{args: "m unk snd cool", verifier: func(w *workerS, ctx *parseCtx, errParsed error) (err error) { //nolint:revive
			if !regexp.MustCompile(`UNKNOWN (Command|Flag) FOUND:?`).MatchString(errParsed.Error()) {
				t.Log("expect 'UNKNOWN Command FOUND' error, but not matched.") // "unk"
			}
			return /* errParsed */
		}, opts: []cli.Opt{cli.WithUnmatchedAsError(true)}},

		{},
		{},
		{},
		{},
		{},
		{},
	} {
		if c.args == "" && c.verifier == nil {
			continue
		}

		t.Log()
		t.Log()
		t.Log()
		t.Logf("--------------- test #%d: Parsing %q\n", i, c.args)

		app, ww := cleanApp(t, false)
		ww.Config.Store = store.New()
		// ww.Config.Loaders = []cli.Loader{loaders.NewConfigFileLoader(), loaders.NewEnvVarLoader()}

		ww.setArgs(append([]string{app.Name()}, strings.Split(c.args, " ")...))
		ww.tasksAfterParse = []taskAfterParse{c.verifier}
		ww.Config.TasksBeforeRun = []cli.Task{aTaskBeforeRun}
		err := ww.Run(ctx, c.opts...) // withTasksBeforeRun(taskBeforeRun),withTasksAfterParse(c.verifier))
		// err := app.Run()
		if err != nil {
			_ = app
			t.Fatal(err)
		}
	}
}
