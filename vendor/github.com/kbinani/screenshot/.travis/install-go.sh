#!/bin/bash

[ -n "$GO_VERSION" ] || exit 1
[ -n "$TRAVIS_REPO_SLUG" ] || exit 6

(
	cd "$HOME"
	OS=$(uname -s | tr A-Z a-z)
	ARCH=$(uname -m)
	if [ "$ARCH" = "x86_64" ]; then
		ARCH=amd64
	elif [ "$ARCH" = "i386" ]; then
		ARCH=386
	else
		exit 7
	fi
	mkdir -p "cache/$TRAVIS_REPO_SLUG"
	if [ ! -f "cache/${TRAVIS_REPO_SLUG}/go${GO_VERSION}.${OS}-${ARCH}.tar.gz" ]; then
		wget "https://dl.google.com/go/go${GO_VERSION}.${OS}-${ARCH}.tar.gz" -O "cache/${TRAVIS_REPO_SLUG}/go${GO_VERSION}.${OS}-${ARCH}.tar.gz" || exit 2
	fi

	if [ ! -f "go${GO_VERSION}/bin/go" ]; then
		rm -rf "go${GO_VERSION}"
		rm -rf "gobootstrap${GO_VERSION}"

		rm -rf go
		tar zxf "cache/${TRAVIS_REPO_SLUG}/go${GO_VERSION}.${OS}-${ARCH}.tar.gz" || exit 3

		mv go "go${GO_VERSION}" || exit 4
		cp -R "go${GO_VERSION}" "gobootstrap${GO_VERSION}" || exit 5
	fi
) || exit $?
