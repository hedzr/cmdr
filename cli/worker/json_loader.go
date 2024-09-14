package worker

import (
	"context"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/store"
	"github.com/hedzr/store/codecs/json"
	"github.com/hedzr/store/providers/file"
)

type jsonLoaderS struct {
	Watch     bool
	WriteBack bool
	filename  string
	hit       bool
	wbh       writeBackHandler
}

type writeBackHandler interface {
	Save(ctx context.Context) error
}

func (j *jsonLoaderS) Load(app cli.App) (err error) {
	var wr writeBackHandler
	wr, err = app.Store().Load(context.Background(),
		// store.WithStorePrefix("app.yaml"),
		// store.WithPosition("app"),
		store.WithCodec(j.codec()),
		store.WithProvider(file.New(j.filename,
			file.WithWatchEnabled(j.Watch),
			file.WithWriteBackEnabled(j.WriteBack))),
	)
	if err == nil && j.WriteBack && wr != nil {
		j.wbh = wr
		j.hit = true
	}

	// TODO implement me
	// panic("implement me")

	return
}

func (j *jsonLoaderS) Save(ctx context.Context) (err error) {
	if j.hit && j.WriteBack && j.wbh != nil {
		err = j.wbh.Save(ctx)
	}
	return
}

func (j *jsonLoaderS) codec() store.Codec {
	return json.New()
}
