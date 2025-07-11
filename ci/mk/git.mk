
GIT_VERSION    := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
GIT_REVISION   := $(shell git rev-parse --short HEAD)
GIT_SUMMARY    := $(shell git describe --tags --dirty --always)
GIT_DESC       := $(shell git log --oneline -1)
GIT_HASH       := $(shell git rev-parse HEAD)

DIRTY_VERSION  := $(shell git describe --dirty --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
#COMMIT         := $(shell git rev-parse HEAD)
