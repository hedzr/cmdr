
APP_NAME = demo
BIN := ../../bin

build:
	go build --tags=hzstudio -o $(BIN)/$(APP_NAME) ./
	ls -la $(BIN)/$(APP_NAME)

run: build
	@# go run ./
	[ -x $(BIN)/$(APP_NAME) ] && $(BIN)/$(APP_NAME)

