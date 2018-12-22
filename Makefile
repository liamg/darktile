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
install: build
	go install -ldflags "-X github.com/liamg/aminal/version.Version=`git describe --tags`"

.PHONY: install-tools
install-tools:
	which dep || curl -L https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	which packr || go get -u github.com/gobuffalo/packr/packr

.PHONY:	build-linux
build-linux:
	mkdir -p bin/linux
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o bin/linux/${BINARY}-linux-amd64 -ldflags "-X github.com/liamg/aminal/version.Version=${TRAVIS_TAG}"

.PHONY: windows-cross-compile
windows-cross-compile:
	mkdir -p bin/windows
	GOOS=windows GOARCH=386 CGO_ENABLED=1 CXX=i686-w64-mingw32-g++ CC=i686-w64-mingw32-gcc go build -o bin/windows/${BINARY}-windows-386.exe -ldflags "-X github.com/liamg/aminal/version.Version=${TRAVIS_TAG}"

.PHONY:	build-windows
build-windows:
	go build -o ${BINARY}-windows-amd64.exe -ldflags "-X github.com/liamg/aminal/version.Version=${TRAVIS_TAG}"


.PHONY:	build-darwin-native
build-darwin-native:
	mkdir -p bin/darwin
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o bin/darwin/${BINARY}-darwin-amd64 -ldflags "-X github.com/liamg/aminal/version.Version=${TRAVIS_TAG}"

.PHONY:	build-darwin
build-darwin:
	mkdir -p bin/darwin
	xgo -x -v -ldflags "-X github.com/liamg/aminal/version.Version=${TRAVIS_TAG}" --targets=darwin/amd64 -out bin/darwin/${BINARY} .

.PHONY:	package-debian
package-debian: build-linux
	./scripts/package-debian.sh "${TRAVIS_TAG}" bin/linux/${BINARY}-linux-amd64
