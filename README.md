# Aminal - A Modern Terminal Emulator

[![CircleCI](https://circleci.com/gh/liamg/aminal/tree/master.svg?style=svg)](https://circleci.com/gh/liamg/aminal/tree/master)
[![GoReportCard](https://goreportcard.com/badge/github.com/liamg/aminal)](https://goreportcard.com/report/github.com/liamg/aminal)
[![Github Release](https://img.shields.io/github/release/liamg/aminal.svg)](https://github.com/liamg/aminal/releases)

Aminal is a modern terminal emulator for Mac/Linux implemented in Golang and utilising OpenGL. 

![Demo GIF](demo.gif)

The project is experimental at the moment, so you probably won't want to rely on Aminal as your main terminal for a while.

Ensure you have your latest graphics card drivers installed before use.

## Features

- Unicode support
- OpenGL rendering
- Customisation options
- True colour support
- Support for common ANSI escape sequences a la xterm
- Scrollback buffer
- Clipboard access
- Clickable URLs
- Multi platform support (Windows coming soon...)
- Sixel support
- Hints/overlays
- Built-in patched fonts for powerline
- Retina display support

## Quick Start

### Installation

#### Prebuilt Binaries

Prebuilt binaries are available for Linux and OSX on the [releases](https://github.com/liamg/aminal/releases) page. 

Download the binary and `sudo cp aminal-* /usr/local/bin/aminal`.

#### Install with Go

```
go get -u github.com/liamg/aminal
```

### Build 

#### Dependencies

- On macOS, you need Xcode or Command Line Tools for Xcode (`xcode-select --install`) for required headers and libraries.
- On Ubuntu/Debian-like Linux distributions, you need `libgl1-mesa-dev xorg-dev`.
- On CentOS/Fedora-like Linux distributions, you need `libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel`.

#### Building Locally

There are various make targets available, the most obvious being:

```
make test
make build
make install
```

As long as you have your `GOBIN` environment variable set up properly (and in `PATH`), you should be able to run `aminal`.

## Keyboard/Mouse Shortcuts

| Operation            | Key(s)               |
| -------------------- | -------------------- |
| Select text          | click + drag         |
| Select word          | double click         |
| Select line          | triple click         |
| Copy                 | `ctrl + shift + c` (Mac: `super + c`) |
| Paste                | `ctrl + shift + v` (Mac: `super + v`) |
| Search online for selected text | `ctrl + shift + g` (Mac: `super + g`) |
| Toggle debug display | `ctrl + shift + d` (Mac: `super + d`) |
| Toggle slomo         | `ctrl + shift + ;` (Mac: `super + ;`) |
| Report bug in aminal | `ctrl + shift + r` (Mac: `super + r`) |

## Configuration

Aminal looks for a config file in `~/.aminal.toml`, and will write one there the first time it runs, if it doesn't already exist.

You can ignore the config and use defaults by specifying `--ignore-config` as a CLI flag.

### Config File

```toml
debug = false               # Enable debug logging to stdout. Defaults to false.
slomo = false               # Enable slow motion output mode, useful for debugging shells/terminal GUI apps etc. Defaults to false.
shell = "/bin/bash"         # The shell to run for the terminal session. Defaults to the users shell.
search_url = "https://www.google.com/search?q=$QUERY" # The search engine to use for the "search selected text" action. Defaults to google. Set this to your own search url using $QUERY as the keywords to replcae when searching.

[colours]
  cursor        = "#e8dfd6" 
  foreground    = "#e8dfd6" 
  background    = "#021b21" 
  black         = "#032c36" 
  red           = "#c2454e" 
  green         = "#7cbf9e"
  yellow        = "#8a7a63"
  blue          = "#065f73"
  magenta       = "#ff5879"
  cyan          = "#44b5b1"
  light_grey    = "#f2f1b9"
  dark_grey     = "#3e4360"
  light_red     = "#ef5847"
  light_green   = "#a2db91"
  light_yellow  = "#beb090"
  light_blue    = "#61778d"
  light_magenta = "#ff99a1"
  light_cyan    = "#9ed9d8"
  white         = "#f6f6c9"
  selection     = "#333366" # Mouse selection background colour

[keys]
  copy      = "ctrl + shift + c"    # Copy highlighted text to system clipboard
  paste     = "ctrl + shift + v"    # Paste text from system clipboard
  debug     = "ctrl + shift + d"    # Toggle debug panel overlay
  google    = "ctrl + shift + g"    # Google selected text
  report    = "ctrl + shift + r"    # Send bug report
  slomo     = "ctrl + shift + ;"    # Toggle slow motion output mode (useful for debugging)
```

### CLI Flags

| Flag              | Description                                                                                                                   |
| ----------------- | ----------------------------------------------------------------------------------------------------------------------------- |
| `--debug`         | Enable debug mode, with debug logging and debug info terminal overlay.
| `--slomo`         | Enable slomo mode, delay the handling of each incoming byte (or escape sequence) from the pty by 100ms. Useful for debugging.
| `--shell [shell]` | Use the specified shell program instead of the user's usual one. 
| `--version`       | Show the version of aminal and exit.
