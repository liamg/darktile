.PHONY: all build test lint clean deps devel-deps

BIN := clipboard
BUILD_LDFLAGS := "-s -w"
GOBIN ?= $(shell go env GOPATH)/bin
export GO111MODULE=on

all: clean build

deps:
	go mod tidy

devel-deps: deps
	GO111MODULE=off go get -u \
	  golang.org/x/lint/golint

build: clean
	go build -ldflags=$(BUILD_LDFLAGS) -o $(BIN)

test: deps
	go test -v -race -count=1 ./...

test-cover: deps
	go test -v -race -count=1 ./... -cover -coverprofile=c.out
	go tool cover -html=c.out -o coverage.html

lint: devel-deps
	go vet ./...
	$(GOBIN)/golint -set_exit_status ./...

clean:
	rm -rf $(BIN)
	go clean
