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
	which github-release || go get -u github.com/aktau/github-release
	which xgo || go get github.com/karalabe/xgo

.PHONY: update-fonts
update-fonts: install-tools
	curl -L https://github.com/ryanoasis/nerd-fonts/raw/master/patched-fonts/Hack/Regular/complete/Hack%20Regular%20Nerd%20Font%20Complete.ttf -o "${FONTPATH}/Hack Regular Nerd Font Complete.ttf"
	curl -L https://github.com/ryanoasis/nerd-fonts/raw/master/patched-fonts/Hack/Bold/complete/Hack%20Bold%20Nerd%20Font%20Complete.ttf -o "${FONTPATH}/Hack Bold Nerd Font Complete.ttf"
	packr -v

.PHONY: release
release: test install-tools
	./scripts/release.sh
	