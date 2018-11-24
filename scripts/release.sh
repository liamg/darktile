#!/bin/bash
set -e



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
