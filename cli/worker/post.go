package worker

import (
	"context"

	"github.com/hedzr/cmdr/v2/pkg/logz"
)

func (w *workerS) postProcess(ctx context.Context, pc *parseCtx) (err error) {
	logz.VerboseContext(ctx, "post-processing...")
	_, _ = pc, ctx
	return
}
