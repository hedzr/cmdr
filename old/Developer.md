

### For Developer

To build and test `cmdr`:

```bash
$ make help   # see all available sub-targets
$ make info   # display building environment
$ make build  # build binary files for examples
$ make gocov  # test

# customizing
$ GOPROXY_CUSTOM=https://goproxy.io make info
$ GOPROXY_CUSTOM=https://goproxy.io make build
$ GOPROXY_CUSTOM=https://goproxy.io make gocov
```

### Build your cli app with `cmdr`

```bash
APP_NAME=your-app-name
APP_VERSION=your-app-version

W_PKG=github.com/hedzr/cmdr/conf

TIMESTAMP=$(date -u '+%Y-%m-%d_%I:%M:%S%p')
GITHASH=$(git rev-parse HEAD)
GOVERSION=$(go version)

LDFLAGS="-s -w -X '$W_PKG.Buildstamp=$TIMESTAMP' -X '$W_PKG.Githash=$GITHASH' -X '$W_PKG.GoVersion=$GOVERSION' -X '$W_PKG.Version=$APP_VERSION' -X '$W_PKG.AppName=$APP_NAME"

go build -ldflags "$LDFLAGS" -o bin/app-name ./cli
```

