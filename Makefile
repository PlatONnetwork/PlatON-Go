# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: geth android ios evm all test clean

GOBIN = $(shell pwd)/build/bin
GO ?= latest
GPATH = $(shell go env GOPATH)
GORUN = env GO111MODULE=on GOPATH=$(GPATH) go run

platon:
	build/build_deps.sh
	$(GORUN) build/ci.go install ./cmd/platon
	@echo "Done building."
	@echo "Run \"$(GOBIN)/platon\" to launch platon."


all:
	build/build_deps.sh
	$(GORUN) build/ci.go install
	@mv $(GOBIN)/keytool $(GOBIN)/platonkey



test: all
	$(GORUN) build/ci.go test

lint: ## Run linters.
	$(GORUN) build/ci.go lint

clean:
	./build/clean_deps.sh
	env GO111MODULE=on go clean -cache
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go install golang.org/x/tools/cmd/stringer@latest
	env GOBIN= go install github.com/fjl/gencodec@latest
	env GOBIN= go install github.com/golang/protobuf/protoc-gen-go@latest
	env GOBIN= go install ./cmd/abigen
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'
