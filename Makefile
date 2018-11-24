SHELL := /bin/bash
BINARY := aminal
FONTPATH := ./gui/packed-fonts

.PHONY: build
build: test install-tools
	mkdir -p release
	packr -v
	go build

.PHONY: test
test:
	go test -v ./...
	go vet -v

.PHONY: install
install: build
	install -m 0755 aminal "${GOBIN}/${BINARY}"

.PHONY: install-tools
install-tools:
	which dep || curl -L https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	which packr || go get -u github.com/gobuffalo/packr/packr

.PHONY: update-fonts
update-fonts: install-tools
	curl -L https://github.com/ryanoasis/nerd-fonts/raw/master/patched-fonts/Hack/Regular/complete/Hack%20Regular%20Nerd%20Font%20Complete.ttf -o "${FONTPATH}/Hack Regular Nerd Font Complete.ttf"
	curl -L https://github.com/ryanoasis/nerd-fonts/raw/master/patched-fonts/Hack/Bold/complete/Hack%20Bold%20Nerd%20Font%20Complete.ttf -o "${FONTPATH}/Hack Bold Nerd Font Complete.ttf"
	packr -v

.PHONY:	build-linux
build-linux:
	mkdir -p release/bin/linux
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o bin/${BINARY}-linux-amd64 -ldflags "-X main.Version=`${CIRCLE_TAG}`"

.PHONY:	build-darwin
build-darwin:
	mkdir -p release/bin/darwin
	xgo -ldflags "-X main.Version=`${CIRCLE_TAG}`" --targets=darwin/amd64 --dest=bin/darwin -out ${BINARY}-darwin-amd64 .
