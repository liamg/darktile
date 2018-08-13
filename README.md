Raft is a modern terminal emulator utilising OpenGL.

The project is experimental at the moment, so you probably won't want to rely on Raft as your main terminal for a while.

Ensure you have your latest graphics card drivers installed before use.

## Aims

- Full unicode support
- OpenGL rendering
- Full customisation options
- True colour support
- Support for commmon ANSI escape sequences a la xterm
- Scrollback buffer
- Clipboard access
- Clickable URLs
- Resize logic that wraps/unwraps lines _correctly_
- Bullshit graphical effects
- Multi platform support

## What isn't supported?

- Suspend/Continue (^S, ^Q). This is archaic bullshit that annoys more people than it helps. Basically:

<span style="display:block;text-align:center">
![Overheating](https://imgs.xkcd.com/comics/workflow.png)
</span>


## Build Dependencies

- Go 1.10.3+
- On macOS, you need Xcode or Command Line Tools for Xcode (`xcode-select --install`) for required headers and libraries.
- On Ubuntu/Debian-like Linux distributions, you need `libgl1-mesa-dev xorg-dev`.
- On CentOS/Fedora-like Linux distributions, you need `libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel`.


## Platform Support

| Platform | Supported  |
|----------|------------|
| Linux    | ✔
| MacOSX   | ⏳
| Windows  | ⏳


## Planned Features

| Feature                     | Done | Notes |
|-----------------------------|------|-------|
| Pty allocation              | ✔    | Needs work for OSX + Windows
| OpenGL rendering            | ✔    |
| Resizing/content reordering | ⏳    | 
| ANSI escape codes           | ⏳    | Most of these are handled now
| UTF-8 input                 | ✔    | 
| UTF-8 output                | ✔    | Works as long as the font in use supports the relevant characters.
| Copy/paste                  |      | Paste working, no mouse interaction for copy
| Customisable colour schemes | ✔    | Complete, but the config file has no entry for this yet 
| Config file                 | ⏳    | Minimal options atm
| Scrolling                   | ⏳    | Infinite buffer implemented, need GUI scrollbar & render updates
| Mouse interaction           |      | 
| Sweet render effects        |      | 

## Keyboard Shortcuts

| Operation          | Key(s)              |
|--------------------|---------------------|
| Paste              | ctrl + shift + v
| Toggle slomo       | ctrl + shift + ;
| Interrupt (SIGINT) | ctrl + c

## Configuration

Raft looks for a config file in `~/.raft.toml`, and will write one there the first time it runs, if it doesn't already exist.

You can ignore the config and use defauls by specifying `--ignore-config` as a CLI flag.

### Config Options/CLI Flags

| CLI Flag           | Config Section      | Config Name            | Type    | Default      | Description |
|--------------------|---------------------|------------------------|---------|--------------|-------------|
| --debug            | _root_              | debug                  | boolean | false        | Enable debug mode, with debug logging and debug info terminal overlay.
| --slomo            | _root_              | slomo                  | boolean | false        | Enable slomo mode, delay the handling of each incoming byte (or escape sequence) from the pty by 100ms. Useful for debugging.
| --always-repaint   | rendering           | always_repaint         | boolean | false        | Redraw the terminal GUI constantly, even when no changes have occurred.


