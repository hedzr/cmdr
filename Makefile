PROJECTNAME1=$(shell basename "$(PWD)")
PROJECTNAME=$(PROJECTNAME1:go-%=%)
include .env
-include .env.local
# ref: https://kodfabrik.com/journal/a-good-makefile-for-go/

  # https://www.gnu.org/savannah-checkouts/gnu/make/manual/html_node/Text-Functions.html
  # https://stackoverflow.com/questions/19571391/remove-prefix-with-make

APPNAME=$(shell grep -E "AppName[ \t]+=[ \t]+" doc.go|grep -Eo "\\\".+\\\"")
VERSION=$(shell grep -E "Version[ \t]+=[ \t]+" doc.go|grep -Eo "[0-9.]+")

# Go related variables.
GOBASE = $(shell pwd)
##GOPATH="$(GOBASE)/vendor:$(GOBASE)"
##GOPATH=$(GOBASE)/vendor:$(GOBASE):$(shell dirname $(GOBASE))
#GOPATH2= $(shell dirname $(GOBASE))
#GOPATH1= $(shell dirname $(GOPATH2))
#GOPATH0= $(shell dirname $(GOPATH1))
#GOPATH = $(shell dirname $(GOPATH0))
GOBIN  = $(GOBASE)/bin
GOFILES= $(wildcard *.go)
BIN    = $(GOPATH)/bin
GOLINT = $(BIN)/golint
GOCYCLO= $(BIN)/gocyclo
GOYOLO = $(BIN)/yolo

GO111MODULE = on
GOPROXY     = $(or $(GOPROXY_CUSTOM),https://athens.azurefd.net)

# Redirect error output to a file, so we can show it in development mode.
STDERR      = $(or $(STDERR_CUSTOM),/tmp/.$(PROJECTNAME)-stderr.txt)

# PID file will keep the process id of the server
PID         = $(or $(PID_CUSTOM),/tmp/.$(PROJECTNAME).pid)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

#
#LDFLAGS=
M = $(shell printf "\033[34;1m▶\033[0m")
ADDR = ":5q5q"
SERVER_START_ARG=server run
SERVER_STOP_ARG=server stop


goarch=amd64
W_PKG=github.com/hedzr/cmdr/conf
TIMESTAMP=$(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
GITHASH=$(shell git rev-parse HEAD)
GOVERSION=$(shell go version)
LDFLAGS := -s -w -X '$(W_PKG).Buildstamp=$(TIMESTAMP)' -X '$(W_PKG).Githash=$(GITHASH)' -X '$(W_PKG).GoVersion=$(GOVERSION)' -X '$(W_PKG).Version=$(VERSION)'
 # -X '$(W_PKG).AppName=$(APPNAME)'

ifeq ($(OS),Windows_NT)
    LS_OPT=
    CCFLAGS += -D WIN32
    ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
        CCFLAGS += -D AMD64
    else
        ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
            CCFLAGS += -D AMD64
        endif
        ifeq ($(PROCESSOR_ARCHITECTURE),x86)
            CCFLAGS += -D IA32
        endif
    endif
else
    LS_OPT=
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        CCFLAGS += -D LINUX
        LS_OPT=--color
    endif
    ifeq ($(UNAME_S),Darwin)
        CCFLAGS += -D OSX
        LS_OPT=-G
    endif
    UNAME_P := $(shell uname -p)
    ifeq ($(UNAME_P),x86_64)
        CCFLAGS += -D AMD64
    endif
    ifneq ($(filter %86,$(UNAME_P)),)
        CCFLAGS += -D IA32
    endif
    ifneq ($(filter arm%,$(UNAME_P)),)
        CCFLAGS += -D ARM
    endif
endif

CN = hedzr/$(N)




.PHONY: build compile exec clean
.PHONY: run build-linux build-ci
.PHONY: go-build go-generate go-mod-download go-get go-install go-clean
.PHONY: godoc format fmt lint cov gocov coverage codecov cyclo bench


## build: Compile the binary. Synonym of `compile`
build: compile

## build-linux: build to linux executable, for LAN deploy manually.
build-linux:
	@echo "  >  Building linux binary..."
	@echo "  >  LDFLAGS = $(LDFLAGS)"
	$(foreach an, fluent demo ffdemo short wget-demo, \
	$(foreach os, linux, \
	  echo "     Building $(GOBIN)/$(an)_$(os)_$(goarch)...$(os)"; \
	    GOARCH="$(goarch)" GOOS="$(os)" \
	    GOPATH="$(GOPATH)" GOBIN="$(GOBIN)" GO111MODULE="$(GO111MODULE)" GOPROXY="$(GOPROXY)" \
	      go build -ldflags "$(LDFLAGS) -X '$(W_PKG).AppName=$(an)'" -o $(GOBIN)/$(an)_$(os)_$(goarch) $(GOBASE)/examples/$(an)/main.go; \
	    chmod +x $(GOBIN)/$(an)_$(os)_$(goarch); \
	    ls -la $(LS_OPT) $(GOBIN)/$(an)_$(os)_$(goarch); \
	) \
	)
	#@ls -la $(LS_OPT) $(GOBIN)/*linux*

## build-nacl: build to linux executable, for LAN deploy manually.
build-nacl:
	@echo "  >  Building linux binary..."
	@echo "  >  LDFLAGS = $(LDFLAGS)"
	$(foreach an, short ffdemo, \
	$(foreach os, nacl, \
	$(foreach goarch, 386 amd64p32 arm, \
	  echo "     >> Building $(GOBIN)/$(an)_$(os)_$(goarch)...$(os)"; \
	    GOARCH="$(goarch)" GOOS="$(os)" \
	    GOPATH="$(GOPATH)" GOBIN="$(GOBIN)" GO111MODULE="$(GO111MODULE)" GOPROXY="$(GOPROXY)" \
	      go build -ldflags "$(LDFLAGS) -X '$(W_PKG).AppName=$(an)'" -o $(GOBIN)/$(an)_$(os)_$(goarch) $(GOBASE)/examples/$(an)/main.go; \
	    chmod +x $(GOBIN)/$(an)_$(os)_$(goarch); \
	    ls -la $(LS_OPT) $(GOBIN)/$(an)_$(os)_$(goarch); \
	) \
	) \
	)
	#@ls -la $(LS_OPT) $(GOBIN)/*linux*

## build-ci: run build-ci task. just for CI tools
build-ci:
	@echo "  >  Building binaries in CI flow..."
	@echo "  >  LDFLAGS = $(LDFLAGS)"
	$(foreach an, fluent ffdemo demo short wget-demo, \
	  echo "  >  APPNAME = $(APPNAME)|$(an)"; \
	  $(foreach os, darwin linux windows, \
	    echo "     Building $(GOBIN)/$(an)_$(os)_$(goarch)...$(os)"; \
	      GOARCH="$(goarch)" GOOS="$(os)" \
	      GOPATH="$(GOPATH)" GOBIN="$(GOBIN)" GO111MODULE="$(GO111MODULE)" GOPROXY="$(GOPROXY)" \
	        go build -ldflags "$(LDFLAGS) -X '$(W_PKG).AppName=$(an)'" -o $(GOBIN)/$(an)_$(os)_$(goarch) $(GOBASE)/examples/$(an)/main.go; \
	        gzip -f $(GOBIN)/$(an)_$(os)_$(goarch); \
	  ) \
	)
	@ls -la $(LS_OPT) $(GOBIN)/*




## compile: Compile the binary.
compile: go-clean go-generate
	@-touch $(STDERR)
	@-rm $(STDERR)
	# @-$(MAKE) info
	@-$(MAKE) -s go-build 2> $(STDERR)
	# @cat $(STDERR) | sed -e '1s/.*/\nError:\n/'  | sed 's/make\[.*/ /' | sed "/^/s/^/     /" 1>&2
	#
	@cat $(STDERR) | sed -e '1s/.*/\nError:\n/' 1>&2

## exec: Run given cmd, wrapped with custom GOPATH. eg; make exec run="go test ./..."
exec:
	@GOPATH=$(GOPATH) GOBIN=$(BIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	$(run)

## clean: Clean build files. Runs `go clean` internally.
clean:
	@(MAKEFILE) go-clean

# go-compile: go-clean go-generate go-build


## run: go run xxx
run:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	go run -ldflags "$(LDFLAGS)" $(GOBASE)/cli/main.go

go-build:
	@echo "  >  Building binary '$(GOBIN)/$(APPNAME)'..."
	# demo short wget-demo 
	$(foreach an, fluent ffdemo demo issue2, \
	  echo "  >  +race. APPNAME = $(APPNAME)|$(an), LDFLAGS = $(LDFLAGS)"; \
	  GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	    go build -v -race -ldflags "$(LDFLAGS) -X '$(W_PKG).AppName=$(an)'" -o $(GOBIN)/$(an) $(GOBASE)/examples/$(an)/main.go; \
	  ls -la $(LS_OPT) $(GOBIN)/$(an); \
	)
	ls -la $(LS_OPT) $(GOBIN)/
	# go build -o $(GOBIN)/$(APPNAME) $(GOFILES)
	# chmod +x $(GOBIN)/*

go-generate:
	@echo "  >  Generating dependency files..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	go generate $(generate)
	# @echo "     done"

go-mod-download:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	go mod download

go-get:
	# Runs `go get` internally. e.g; make install get=github.com/foo/bar
	@echo "  >  Checking if there is any missing dependencies...$(get)"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	go get $(get)

go-install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	go install $(GOFILES)

go-clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	go clean
	# @echo "     Clean done"



$(BIN)/golint: | $(GOBASE)   # # # ❶
	@echo "  >  installing golint ..."
	#@-mkdir -p $(GOPATH)/src/golang.org/x/lint/golint
	#@cd $(GOPATH)/src/golang.org/x/lint/golint
	#@pwd
	#@GOPATH=$(GOPATH) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	#go get -v golang.org/x/lint/golint
	@echo "  >  installing golint ..."
	@GOPATH=$(GOPATH) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	go install golang.org/x/lint/golint
	@cd $(GOBASE)

$(BIN)/gocyclo: | $(GOBASE)  # # # ❶
	@echo "  >  installing gocyclo ..."
	@GOPATH=$(GOPATH) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	go install github.com/fzipp/gocyclo

$(BIN)/yolo: | $(GOBASE)     # # # ❶
	@echo "  >  installing yolo ..."
	@GOPATH=$(GOPATH) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	go install github.com/azer/yolo

$(BIN)/godoc: | $(GOBASE)     # # # ❶
	@echo "  >  installing godoc ..."
	@GOPATH=$(GOPATH) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	go install golang.org/x/tools/cmd/godoc

$(BASE):
	# @mkdir -p $(dir $@)
	# @ln -sf $(CURDIR) $@


## godoc: run godoc server at "localhost;6060"
godoc: | $(GOBASE) $(BIN)/godoc
	@echo "  >  PWD = $(shell pwd)"
	@echo "  >  started godoc server at :6060: http://localhost:6060/pkg/github.com/hedzr/$(PROJECTNAME1) ..."
	@echo "  $  GOPATH=$(GOPATH) godoc -http=:6060 -index -notes '(BUG|TODO|DONE|Deprecated)' -play -timestamps"
	@GOPATH=$(GOPATH) \
	$(BIN)/godoc -http=:6060 -notes '(BUG|TODO|DONE)' -play -timestamps
	# -goroot $(GOPATH) -index
	# https://medium.com/@elliotchance/godoc-tips-tricks-cda6571549b


## fmt: =`format`, run gofmt tool
fmt: format

## format: run gofmt tool
format: | $(GOBASE)
	@echo "  >  gofmt ..."
	@GOPATH=$(GOPATH) GOBIN=$(BIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	gofmt -l -w -s .

## lint: run golint tool
lint: | $(GOBASE) $(GOLINT)
	@echo "  >  golint ..."
	@GOPATH=$(GOPATH) GOBIN=$(BIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	$(GOLINT) ./...

## cov: =`coverage`, run go coverage test
cov: coverage

## gocov: =`coverage`, run go coverage test
gocov: coverage

## coverage: run go coverage test
coverage: | $(GOBASE)
	@echo "  >  gocov ..."
	@GOPATH=$(GOPATH) GOBIN=$(BIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	go test -v -race -coverprofile=coverage.txt -covermode=atomic
	@GOPATH=$(GOPATH) GOBIN=$(BIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	go tool cover -html=coverage.txt -o cover.html
	@open cover.html

## codecov: run go test for codecov; (codecov.io)
codecov: | $(GOBASE)
	@echo "  >  codecov ..."
	@GOPATH=$(GOPATH) GOBIN=$(BIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	go test -v -race -coverprofile=coverage.txt -covermode=atomic
	@bash <(curl -s https://codecov.io/bash) -t $(CODECOV_TOKEN)

## cyclo: run gocyclo tool
cyclo: | $(GOBASE) $(GOCYCLO)
	@echo "  >  gocyclo ..."
	@GOPATH=$(GOPATH) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	$(GOCYCLO) -top 20 .

## bench: benchmark test
bench:
	@echo "  >  benchmark testing ..."
	@GOPATH=$(GOPATH) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	go test -bench="." -run=^$ -benchtime=10s ./...
	# go test -bench "." -run=none -test.benchtime 10s
	# todo: go install golang.org/x/perf/cmd/benchstat




.PHONY: printvars info help all
printvars:
	$(foreach V, $(sort $(filter-out .VARIABLES,$(.VARIABLES))), $(info $(v) = $($(v))) )
	# Simple:
	#   (foreach v, $(filter-out .VARIABLES,$(.VARIABLES)), $(info $(v) = $($(v))) )

print-%:
	@echo $* = $($*)

info:
	@echo "     GOBASE: $(GOBASE)"
	@echo "      GOBIN: $(GOBIN)"
	@echo "     GOROOT: $(GOROOT)"
	@echo "     GOPATH: $(GOPATH)"
	@echo "GO111MODULE: $(GO111MODULE)"
	@echo "    GOPROXY: $(GOPROXY)"
	@echo "PROJECTNAME: $(PROJECTNAME)"
	@echo "    APPNAME: $(APPNAME)"
	@echo "    VERSION: $(VERSION)"
	@echo
	@echo "export GO111MODULE=on"
	@echo "export GOPROXY=$(GOPROXY)"
	@echo "export GOPATH=$(GOPATH)"

all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

