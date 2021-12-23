GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
RESULT_PROTO_FILE=$(shell find ./result -name *.proto)

.PHONY: init
# init env
init:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go

.PHONY: result-proto
# generate result-proto proto
result-proto:
	protoc --proto_path=. \
 	       --go_out=paths=source_relative:. \
	       --proto_path=./third_party \
		   --python_out=. \
	       $(RESULT_PROTO_FILE)


.PHONY: all
# generate all
all:
	make result-proto;

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
