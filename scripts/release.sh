#!/bin/bash
set -e

BINARY=aminal

if [[ "${CIRCLE_TAG}" == "" ]]; then
    exit 0 # no tag, nothing to release
fi
mkdir -p release/bin/darwin/amd64/
mkdir -p release/bin/linux/amd64/
mkdir -p release/bin/linux/i386/

# build for osx using xgo - this cannot be used for linux builds due to missing deps in the xgo containers
# xgo --targets=darwin/amd64 --dest=release/bin/darwin/amd64 -out ${BINARY} .

GOOS=linux GOARCH=386 CGO_ENABLED=1 go build -o release/bin/linux/i386/aminal
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o release/bin/linux/amd64/aminal

exit 0


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
