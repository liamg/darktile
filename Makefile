SHELL := /bin/bash
BINARY := aminal
FONTPATH := ./gui/packed-fonts

.PHONY: build
build: 
	./build.sh `git describe --tags`

.PHONY: test
test:
	go test -v ./...
	go vet -v

.PHONY: install
install: build install-tools
	packr -v
	go install -ldflags "-X github.com/liamg/aminal/version.Version=`git describe --tags`"

.PHONY: install-tools
install-tools:
	which dep || curl -L https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	which packr || go get -u github.com/gobuffalo/packr/packr

.PHONY:	build-linux
build-linux:
	mkdir -p bin/linux
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o bin/linux/${BINARY}-linux-amd64 -ldflags "-X github.com/liamg/aminal/version.Version=${CIRCLE_TAG}"

.PHONY:	build-darwin
build-darwin:
	mkdir -p bin/darwin
	xgo -x -v -ldflags "-X github.com/liamg/aminal/version.Version=${CIRCLE_TAG}" --targets=darwin/amd64 -out bin/darwin/${BINARY} .

.PHONY:	package-debian
package-debian: build-linux
	./scripts/package-debian.sh "${CIRCLE_TAG}" bin/linux/${BINARY}-linux-amd64
