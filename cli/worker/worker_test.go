package worker

import (
	"context"
	"regexp"
	"testing"

	"github.com/hedzr/cmdr/v2/cli"
)

func TestWorkerS_Run2(t *testing.T) { //nolint:revive
	cases := cmdrRunTests{[]cmdrRunTest{
		{args: "m unk snd cool", verifier: func(ctx context.Context, w *workerS, pc *parseCtx, errParsed error) (err error) { //nolint:revive
			if !regexp.MustCompile(`UNKNOWN (Cmd|Flag) FOUND:?`).MatchString(errParsed.Error()) {
				t.Log("expect 'UNKNOWN Cmd FOUND' error, but not matched.") // "unk"
			}
			return /* errParsed */
		}, opts: []cli.Opt{cli.WithUnmatchedAsError(true)}},

		{},
		{},
		{},
	}}

	ctx := context.TODO()
	testWorkerS_Parse(ctx, t, cases)
}
