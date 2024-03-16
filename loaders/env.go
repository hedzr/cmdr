package loaders

import (
	"context"

	"github.com/hedzr/cmdr/v2/cli"
	logz "github.com/hedzr/logg/slog"
	"github.com/hedzr/store"
	"github.com/hedzr/store/providers/env"
)

func NewEnvVarLoader() *envvarloader {
	return &envvarloader{}
}

type envvarloader struct{}

func (w *envvarloader) Load(app cli.App) (err error) {
	conf := app.Store()
	name := app.Name()
	_, err = conf.Load(context.Background(),
		store.WithProvider(env.New(
			env.WithStorePrefix("app.cmd"),
			env.WithPrefix(name+"_", name+"_"),
			env.WithLowerCase(true),
			env.WithUnderlineToDot(true),
		)),
		// store.WithStorePrefix("app.cmd"),
	)
	if err == nil {
		logz.Verbose("envvars loaded")
	}
	return
}
