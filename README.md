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

## Configuration

Raft looks for a config file in the following places: `~/.raft.yml`, `~/.raft/config.yml`, `~/.config/raft/config.yml` (earlier in the list prioritised).

Example config:
```
debug: False
```

The following options are available:

| Name          | Type    | Default |Description            |
|---------------|---------|---------|-----------------------|
| debug         | bool    | False   | Enables debug logging 


## Flags

| Name            | Type    | Default |Description            |
|-----------------|---------|---------|-----------------------|
| --debug         | bool    | False   | Enables debug logging |
| --ignore-config | bool    | False   | Ignores user config files and uses defaults

## Keyboard Shortcuts

| Operation | Key(s)              |
|-----------|---------------------|
| Paste     | ctrl + alt + v

