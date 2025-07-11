
PLATFORM       ?= linux
ARCH           ?= amd64
BUILD_DIR      ?= bin
LOGS_DIR       ?= ./logs


GO             := $(shell which go)
GOOS           := $(shell go env GOOS)
GOARCH         := $(shell go env GOARCH)
GOPROXY        := $(shell go env GOPROXY)
GOVERSION      := $(shell go version)
DEFAULT_TARGET := $(GOOS)-$(GOARCH)
W_PKG          := github.com/hedzr/cmdr/v2/conf
CMDR_SETTING   := \
	-X '$(W_PKG).Buildstamp=$(TIMESTAMP)' \
	-X '$(W_PKG).Githash=$(GIT_REVISION)' \
	-X '$(W_PKG).GitSummary=$(GIT_SUMMARY)' \
	-X '$(W_PKG).GitDesc=$(GIT_DESC)' \
	-X '$(W_PKG).BuilderComments=$(BUILDER_COMMENT)' \
	-X '$(W_PKG).GoVersion=$(GOVERSION)' \
	-X '$(W_PKG).Version=$(GIT_VERSION)' \
	-X '$(W_PKG).AppName=$(NAME)'
GOBUILD := CGO_ENABLED=0 \
	$(GO) build \
	-tags "cmdr hzstudio sec antonal" \
	-trimpath \
	-ldflags="-s -w $(CMDR_SETTING)" \
	-o $(BUILD_DIR)

