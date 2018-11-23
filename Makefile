SHELL := /bin/bash
BINARY := aminal
VERSION ?= vlatest
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
	sudo install -m 0755 aminal "/usr/local/bin/${BINARY}"

.PHONY: install-tools
install-tools:
	which dep || curl -L https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	which packr || go get -u github.com/gobuffalo/packr/packr
	which github-release || go get -u github.com/aktau/github-release

.PHONY: update-fonts
update-fonts: install-tools
	curl -L https://github.com/ryanoasis/nerd-fonts/raw/master/patched-fonts/Hack/Regular/complete/Hack%20Regular%20Nerd%20Font%20Complete.ttf -o "${FONTPATH}/Hack Regular Nerd Font Complete.ttf"
	curl -L https://github.com/ryanoasis/nerd-fonts/raw/master/patched-fonts/Hack/Bold/complete/Hack%20Bold%20Nerd%20Font%20Complete.ttf -o "${FONTPATH}/Hack Bold Nerd Font Complete.ttf"
	packr -v

.PHONY: release
release: test install-tools
	echo -n "Enter a version: "
	read -s VERSION
	if [[ "${VERSION}" == "" ]]; then
		exit 1
	fi
	mkdir -p release/bin/darwin/amd64/
	mkdir -p release/bin/linux/amd64/
	mkdir -p release/bin/linux/i386/
	GOOS=darwin GOARCH=amd64 go build -o release/bin/darwin/amd64/${BINARY}
	GOOS=linux GOARCH=amd64 go build -o release/bin/linux/amd64/${BINARY}
	GOOS=linux GOARCH=386 go build -o release/bin/linux/386/${BINARY}
	git tag "${VERSION}"
	git push origin "${VERSION}"
	github-release release \
		--user liamg \
		--repo aminal \
		--tag "${VERSION}" \
		--name "Aminal ${VERSION}"
	github-release upload \
		--user liamg \
		--repo aminal \
		--tag "${VERSION}" \
		--name "${BINARY}-osx-amd64" \
		--file release/bin/darwin/amd64/${BINARY}
	github-release upload \
		--user liamg \
		--repo aminal \
		--tag "${VERSION}" \
		--name "${BINARY}-linux-amd64" \
		--file release/bin/linux/amd64/${BINARY}
	github-release upload \
		--user liamg \
		--repo aminal \
		--tag "${VERSION}" \
		--name "${BINARY}-linux-386" \
		--file release/bin/linux/386/${BINARY}
	