package worker

import (
	"context"

	"github.com/hedzr/cmdr/v2/cli"
)

type genS struct{}

func (w *genS) onAction(ctx context.Context, cmd *cli.Command, args []string) (err error) { //nolint:revive,unused
	return
}

type genShS struct{}

func (w *genShS) onAction(ctx context.Context, cmd *cli.Command, args []string) (err error) { //nolint:revive,unused
	return
}

type genDocS struct{}

func (w *genDocS) onAction(ctx context.Context, cmd *cli.Command, args []string) (err error) { //nolint:revive,unused
	return
}

type genManS struct{}

func (w *genManS) onAction(ctx context.Context, cmd *cli.Command, args []string) (err error) { //nolint:revive,unused
	return
}
