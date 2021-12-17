#!/usr/bin/env bash

version=$(git describe --exact-match --tags 2>/dev/null || git describe 2>/dev/null || echo "prerelease")
go build \
    -mod=vendor\
    -ldflags="-X github.com/liamg/darktile/internal/app/darktile/version.Version=${version}" \
    ./cmd/darktile
