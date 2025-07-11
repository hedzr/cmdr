# Use `make V=1 all` to print commands.
$(V).SILENT:


# MAIN_BUILD_PKG  ?= . 
# MAIN_APPS       ?= cli
# SUB_APPS        ?= cmdr-cli loop
# MAIN_ENTRY_FILE ?= .                 # Or: "main.go"

# NAME           := blueprint
# PACKAGE_NAME   := github.com/hedzr/cmdr/v2
# ENTRY_PKG      := ./examples/blueprint

# PLATFORM       ?= linux
# ARCH           ?= amd64
# BUILD_DIR      ?= bin
# LOGS_DIR       ?= ./logs

-include ./ci/mk/env.mk
-include ./ci/mk/cc.mk
-include ./ci/mk/git.mk

-include .env
-include .env.local

-include ./ci/mk/vars.mk

.PHONY: all $(BUILD_DIR)/$(NAME) release release-all test build primary-target main
all: build
normal: clean $(BUILD_DIR)/$(NAME)

clean:
	rm -rf $(BUILD_DIR)
	rm -f *.zip

## test: run go test
test: cov
## cov: run go coverage
cov: | $(LOGS_DIR)
	@$(HEADLINE) "Running go coverage..."
	$(GO) test ./... -v -race -cover -coverprofile=$(LOGS_DIR)/coverage-cl.txt -covermode=atomic -test.short -vet=off 2>&1 | tee $(LOGS_DIR)/cover-cl.log && echo "RET-CODE OF TESTING: $?"

$(LOGS_DIR):
	@mkdir -pv $@

.PHONY: directories
directories: | $(BUILD_DIR) $(LOGS_DIR)

## build: build executable for current OS and CPU (arch)
build: $(BUILD_DIR)/$(NAME)
	#$(LL) $(BUILD_DIR)/$(NAME)
	cp $(BUILD_DIR)/$(NAME) ~/go/bin/ && echo "- INSTALL TO ~/go/bin OK."
	@echo "- BUILD $(BUILD_DIR)/$(NAME) OK."

## build: build executable for current OS and CPU (arch)
build-default: $(DEFAULT_TARGET)
	@echo BUILD OK

# primary-target: build executable for current GOOS & GOARCH
primary-target: main
# main: build executable for the TARGET GOOS & GOARCH (see also PLATFORM & ARCH vars)
main:
	@-$(MAKE) $(BUILD_DIR)/$(NAME) GOOS=$(PLATFORM) GOARCH=$(ARCH) 

# bin/cmdr is the default executable for running under your current GOOS & GOARCH.
# But you can override them and cross-build whatever targets you want, just like
# what `make cmdr-cli` does.
$(BUILD_DIR)/$(NAME): | $(BUILD_DIR)
	@# mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(ENTRY_PKG)
	$(LL) $(BUILD_DIR)/$(NAME)

# install: install executable with 'make INSTALL=1 install'
install:
	@$(INSTALL_CMD) mkdir -pv $(INSTALL_PREFIX)/etc/$(NAME)
	@$(INSTALL_CMD) mkdir -pv $(INSTALL_PREFIX)/share/$(NAME)
	@$(INSTALL_CMD) cp $(BUILD_DIR)/$(NAME) $(INSTALL_PREFIX)/bin/$(NAME)
	# @$(INSTALL_CMD) cp example/*.json $(INSTALL_PREFIX)/etc/$(NAME)
	# @$(INSTALL_CMD) cp example/$(NAME).service $(INSTALL_PREFIX)/lib/systemd/system/
	# @$(INSTALL_CMD) cp example/$(NAME)@.service $(INSTALL_PREFIX)/lib/systemd/system/
	# @$(INSTALL_CMD) systemctl daemon-reload
	# @$(INSTALL_CMD) ln -fs $(INSTALL_PREFIX)/share/$(NAME)/geoip.dat /usr/bin/
	# @$(INSTALL_CMD) ln -fs $(INSTALL_PREFIX)/share/$(NAME)/geoip-only-cn-private.dat /usr/bin/
	# @$(INSTALL_CMD) ln -fs $(INSTALL_PREFIX)/share/$(NAME)/geosite.dat /usr/bin/
	@$(INSTALL_HELP)

# uninstall: uninstall executable with 'make UNINSTALL=1 uninstall'
uninstall:
	# @$(UNINSTALL_CMD) rm $(INSTALL_PREFIX)/lib/systemd/system/$(NAME).service
	# @$(UNINSTALL_CMD) rm $(INSTALL_PREFIX)/lib/systemd/system/$(NAME)@.service
	# @$(UNINSTALL_CMD) systemctl daemon-reload
	@$(UNINSTALL_CMD) rm $(INSTALL_PREFIX)/bin/$(NAME)
	@$(UNINSTALL_CMD) rm -rd $(INSTALL_PREFIX)/etc/$(NAME)
	@$(UNINSTALL_CMD) rm -rd $(INSTALL_PREFIX)/share/$(NAME)
	@$(UNINSTALL_HELP)

%.zip: %
	@zip -du $(NAME)-$@ -j $(BUILD_DIR)/$</*
	# @zip -du $(NAME)-$@ example/*
	# @-zip -du $(NAME)-$@ *.dat
	@echo "<<< ---- $(NAME)-$@"

release: \
	darwin-amd64.zip darwin-arm64.zip \
	linux-amd64.zip linux-arm64.zip \
	windows-amd64.zip windows-arm64.zip
	$(LL) $(NAME)-*.*

release-all: darwin-arm64.zip linux-386.zip linux-amd64.zip \
	linux-arm.zip linux-armv5.zip linux-armv6.zip linux-armv7.zip linux-armv8.zip \
	linux-mips-softfloat.zip linux-mips-hardfloat.zip linux-mipsle-softfloat.zip linux-mipsle-hardfloat.zip \
	linux-mips64.zip linux-mips64le.zip \
	freebsd-386.zip freebsd-amd64.zip \
	windows-386.zip windows-amd64.zip windows-arm.zip windows-armv6.zip \
	windows-armv7.zip windows-arm64.zip

-include ./ci/mk/go-targets.mk
-include ./ci/mk/help.mk

help-extras:
	@echo
	@echo "              GO = $(GO)"
	@echo "            GOOS = $(GOOS)"
	@echo "          GOARCH = $(GOARCH)"
	@echo "         GOPROXY = $(GOPROXY)"
	@echo
	@echo "  DEFAULT_TARGET = $(DEFAULT_TARGET)"
	@echo "       TIMESTAMP = $(TIMESTAMP)"
	@echo "   VERSION (git) = $(GIT_VERSION)"
	@echo " VERSION (dirty) = $(DIRTY_VERSION)"
	@echo "     COMMIT/HASH = $(GIT_HASH)"
	@echo "    GIT_REVISION = $(GIT_REVISION)"
	@echo "        GIT_DESC = $(GIT_DESC)"
	@echo
	@echo "            NAME = $(NAME)"
	@echo
	@echo " To build $(NAME) under current os and arch, use: 'make V=1 build', or 'make build'."
	@echo " To build $(NAME) for others, use: 'make build linux-amd64', or other os and architect names."
