# Aminal - A Modern Terminal Emulator

[![Travis Build Status](https://travis-ci.org/liamg/aminal.svg?branch=master)](https://travis-ci.org/liamg/aminal)
[![GoReportCard](https://goreportcard.com/badge/github.com/liamg/aminal)](https://goreportcard.com/report/github.com/liamg/aminal)
[![Github Release](https://img.shields.io/github/release/liamg/aminal.svg)](https://github.com/liamg/aminal/releases)
[![Slack](https://img.shields.io/badge/slack-%23aminal-%23ffcc00.svg)](http://gophers.slack.com/messages/aminal)
![License](https://img.shields.io/github/license/liamg/aminal.svg)

Aminal is a modern terminal emulator for Mac/Linux/Windows implemented in Golang and utilising OpenGL. 

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
- Multi platform support (Windows, Linux, OSX)
- Sixel support
- Hints/overlays
- Built-in patched fonts for powerline
- Retina display support

## Installation

### MacOS

```
brew tap liamg/aminal
brew install aminal
```

### Windows

You can run Aminal with a binary from the [releases](https://github.com/liamg/aminal/releases) page.

Dev environment setup instructions are available [here](windows.md).

### Prebuilt Binaries

Prebuilt binaries are available for Linux, OSX and Windows on the [releases](https://github.com/liamg/aminal/releases) page. 

Download the binary and `sudo cp aminal-* /usr/local/bin/aminal && chmod +x /usr/local/bin/aminal`.

### Install with Go

```
go get -u github.com/liamg/aminal
```

## Build 

### Dependencies

- On macOS, you need Xcode or Command Line Tools for Xcode (`xcode-select --install`) for required headers and libraries.
- On Ubuntu/Debian-like Linux distributions, you need `libgl1-mesa-dev xorg-dev`.
- On CentOS/Fedora-like Linux distributions, you need `libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel`.

### Building Locally

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

Aminal looks for a config file in the following places, and stops when it finds one:

* `$XDG_CONFIG_HOME/aminal/config.toml`
* `$HOME/.config/aminal/config.toml`
* `$HOME/.aminal.toml`

Note that on Windows Aminal uses `%USERPROFILE%` environment variable instead of `$HOME`.

It will write a config file to whichever of those directories exists (preferring the top of the list) the first time it runs, if one doesn't already exist.

You can ignore the config and use defaults by specifying `--ignore-config` as a CLI flag.

### Config File

```toml
debug = false               # Enable debug logging to stdout. Defaults to false.
slomo = false               # Enable slow motion output mode, useful for debugging shells/terminal GUI apps etc. Defaults to false.
shell = "/bin/bash"         # The shell to run for the terminal session. Defaults to the users shell.
search_url = "https://www.google.com/search?q=$QUERY" # The search engine to use for the "search selected text" action. Defaults to google. Set this to your own search url using $QUERY as the keywords to replace when searching.
max_lines = 1000            # Maximum number of lines in the terminal buffer.
copy_and_paste_with_mouse = true # Text selected with the mouse is copied to the clipboard on end selection, and is pasted on right mouse button click.
show_vertical_scrollbar = true   # Whether to show the vertical scrollbar
dpi-scale = 0.0             # Override DPI scale. Defaults to 0.0 (let Aminal determine the DPI scale itself).

[colours]
  cursor        = "#e8dfd6" 
  foreground    = "#e8dfd6" 
  background    = "#021b21" 
  black         = "#000000" 
  red           = "#800000" 
  green         = "#008000"
  yellow        = "#808000"
  blue          = "#000080"
  magenta       = "#800080"
  cyan          = "#008080"
  light_grey    = "#f2f2f2"
  dark_grey     = "#808080"
  light_red     = "#ff0000"
  light_green   = "#00ff00"
  light_yellow  = "#ffff00"
  light_blue    = "#0000ff"
  light_magenta = "#ff00ff"
  light_cyan    = "#00ffff"
  white         = "#ffffff"
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

# Contributors

[![](https://sourcerer.io/fame/liamg/liamg/aminal/images/0)](https://sourcerer.io/fame/liamg/liamg/aminal/links/0)[![](https://sourcerer.io/fame/liamg/liamg/aminal/images/1)](https://sourcerer.io/fame/liamg/liamg/aminal/links/1)[![](https://sourcerer.io/fame/liamg/liamg/aminal/images/2)](https://sourcerer.io/fame/liamg/liamg/aminal/links/2)[![](https://sourcerer.io/fame/liamg/liamg/aminal/images/3)](https://sourcerer.io/fame/liamg/liamg/aminal/links/3)[![](https://sourcerer.io/fame/liamg/liamg/aminal/images/4)](https://sourcerer.io/fame/liamg/liamg/aminal/links/4)[![](https://sourcerer.io/fame/liamg/liamg/aminal/images/5)](https://sourcerer.io/fame/liamg/liamg/aminal/links/5)[![](https://sourcerer.io/fame/liamg/liamg/aminal/images/6)](https://sourcerer.io/fame/liamg/liamg/aminal/links/6)[![](https://sourcerer.io/fame/liamg/liamg/aminal/images/7)](https://sourcerer.io/fame/liamg/liamg/aminal/links/7)
