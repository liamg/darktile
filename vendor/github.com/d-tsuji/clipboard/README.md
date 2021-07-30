# clipboard

[![Actions Status](https://github.com/d-tsuji/clipboard/workflows/test/badge.svg)](https://github.com/d-tsuji/clipboard/actions)
[![Doc](https://img.shields.io/badge/doc-reference-blue.svg)](https://pkg.go.dev/github.com/d-tsuji/clipboard)
[![Go Report Card](https://goreportcard.com/badge/github.com/d-tsuji/clipboard)](https://goreportcard.com/report/github.com/d-tsuji/clipboard)

This is a multi-platform clipboard library in Go.

## Abstract

- This is clipboard library in Go, which runs on multiple platforms.
- External clipboard package is not required.

## Supported Platforms

- Windows
- macOS
- Linux, Unix (X11)

## Installation

```
go get github.com/d-tsuji/clipboard
```

## API

```go
package clipboard

// Get returns the current text data of the clipboard.
func Get() (string, error)

// Set sets the current text data of the clipboard.
func Set(text string) error
```

