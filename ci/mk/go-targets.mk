
$(BUILD_DIR):
	mkdir -pv $@

darwin-amd64:
	mkdir -pv $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=darwin $(GOBUILD)/$@ $(ENTRY_PKG)

darwin-arm64: | $(BUILD_DIR)
	@$(HEADLINE) "Running go build for $@ ..."
	mkdir -pv $(BUILD_DIR)/$@
	GOARCH=arm64 GOOS=darwin $(GOBUILD)/$@ $(ENTRY_PKG)

linux-386:
	mkdir -pv $(BUILD_DIR)/$@
	GOARCH=386 GOOS=linux $(GOBUILD)/$@ $(ENTRY_PKG)

linux-amd64:
	mkdir -pv $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=linux $(GOBUILD)/$@ $(ENTRY_PKG)

linux-arm64:
	mkdir -pv $(BUILD_DIR)/$@
	GOARCH=arm64 GOOS=linux $(GOBUILD)/$@ $(ENTRY_PKG)

linux-arm:
	mkdir -pv $(BUILD_DIR)/$@
	GOARCH=arm GOOS=linux $(GOBUILD)/$@ $(ENTRY_PKG)

linux-armv5:
	mkdir -pv $(BUILD_DIR)/$@
	GOARCH=arm GOOS=linux GOARM=5 $(GOBUILD)/$@ $(ENTRY_PKG)

linux-armv6:
	mkdir -pv $(BUILD_DIR)/$@
	GOARCH=arm GOOS=linux GOARM=6 $(GOBUILD)/$@ $(ENTRY_PKG)

linux-armv7:
	mkdir -pv $(BUILD_DIR)/$@
	GOARCH=arm GOOS=linux GOARM=7 $(GOBUILD)/$@ $(ENTRY_PKG)

linux-armv8:
	mkdir -pv $(BUILD_DIR)/$@
	GOARCH=arm64 GOOS=linux $(GOBUILD)/$@ $(ENTRY_PKG)

linux-mips-softfloat:
	mkdir -pv $(BUILD_DIR)/$@
	GOARCH=mips GOMIPS=softfloat GOOS=linux $(GOBUILD)/$@ $(ENTRY_PKG)

linux-mips-hardfloat:
	mkdir -pv $(BUILD_DIR)/$@
	GOARCH=mips GOMIPS=hardfloat GOOS=linux $(GOBUILD)/$@ $(ENTRY_PKG)

linux-mipsle-softfloat:
	mkdir -pv $(BUILD_DIR)/$@
	GOARCH=mipsle GOMIPS=softfloat GOOS=linux $(GOBUILD)/$@ $(ENTRY_PKG)

linux-mipsle-hardfloat:
	mkdir -pv $(BUILD_DIR)/$@
	GOARCH=mipsle GOMIPS=hardfloat GOOS=linux $(GOBUILD)/$@ $(ENTRY_PKG)

linux-mips64:
	mkdir -pv $(BUILD_DIR)/$@
	GOARCH=mips64 GOOS=linux $(GOBUILD)/$@ $(ENTRY_PKG)

linux-mips64le:
	mkdir -pv $(BUILD_DIR)/$@
	GOARCH=mips64le GOOS=linux $(GOBUILD)/$@ $(ENTRY_PKG)

freebsd-386:
	mkdir -pv $(BUILD_DIR)/$@
	GOARCH=386 GOOS=freebsd $(GOBUILD)/$@ $(ENTRY_PKG)

freebsd-amd64:
	mkdir -pv $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=freebsd $(GOBUILD)/$@ $(ENTRY_PKG)

windows-386:
	mkdir -pv $(BUILD_DIR)/$@
	GOARCH=386 GOOS=windows $(GOBUILD)/$@ $(ENTRY_PKG)
	#@mv $(BUILD_DIR)/$@/$(NAME) $(BUILD_DIR)/$@/$(NAME).exe

windows-amd64:
	mkdir -pv $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=windows $(GOBUILD)/$@ $(ENTRY_PKG)
	# @mv $(BUILD_DIR)/$@/$(NAME) $(BUILD_DIR)/$@/$(NAME).exe

windows-arm64:
	mkdir -pv $(BUILD_DIR)/$@
	GOARCH=arm64 GOOS=windows $(GOBUILD)/$@ $(ENTRY_PKG)
	# @mv $(BUILD_DIR)/$@/$(NAME) $(BUILD_DIR)/$@/$(NAME).exe

windows-arm:
	mkdir -pv $(BUILD_DIR)/$@
	GOARCH=arm GOOS=windows $(GOBUILD)/$@ $(ENTRY_PKG)
	#@mv $(BUILD_DIR)/$@/$(NAME) $(BUILD_DIR)/$@/$(NAME).exe

windows-armv6:
	mkdir -pv $(BUILD_DIR)/$@
	GOARCH=arm GOOS=windows GOARM=6 $(GOBUILD)/$@ $(ENTRY_PKG)
	#@mv $(BUILD_DIR)/$@/$(NAME) $(BUILD_DIR)/$@/$(NAME).exe

windows-armv7:
	mkdir -pv $(BUILD_DIR)/$@
	GOARCH=arm GOOS=windows GOARM=7 $(GOBUILD)/$@ $(ENTRY_PKG)
	#@mv $(BUILD_DIR)/$@/$(NAME) $(BUILD_DIR)/$@/$(NAME).exe
