

.PHONY: run-demo run-large


run-demo:
	$(MAKE) -C ./examples/demo

run-large:
	$(MAKE) -C ./examples/large

cov:
	go test ./... -v -race -cover -coverprofile=./logs/coverage-cl.txt -covermode=atomic -test.short -vet=off 2>&1 | tee ./logs/cover-cl.log 




