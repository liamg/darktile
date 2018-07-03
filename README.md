Raft is a terminal emulator utilising OpenGL v4.1.

The project is purely a learning exercise right now.

Ensure you have your latest graphics card drivers installed before use.

## Aims

- Full unicode support
- OpenGL rendering
- Full customisation options
- 256 colour support
- Clipboard access
- Clickable URLs


## Build Dependencies

- Go 1.10.3+
- On macOS, you need Xcode or Command Line Tools for Xcode (`xcode-select --install`) for required headers and libraries.
- On Ubuntu/Debian-like Linux distributions, you need `libgl1-mesa-dev xorg-dev`.
- On CentOS/Fedora-like Linux distributions, you need `libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel`.


## Platform Support

| Platform | Supported  |
|----------|------------|
| Linux    | ✔
| MacOSX   | 
| Windows  | 


## Planned Features

| Feature                     | Done | Notes |
|-----------------------------|------|-------|
| Pty allocation              | ✔    | Needs work for other platforms
| OpenGL rendering            | ✔    |
| Resizing/content reordering | ✔    | 
| ANSI escape codes           |      | Most of these are handled now
| UTF-8 input                 | ✔    | 
| UTF-8 output                | ✔    | Works as long as the font in use supports the relevant characters.
| Copy/paste                  |      | Paste working, no mouse interaction for copy
| Customisable colour schemes | ✔    | Complete, but the config file has no entry for this yet 
| Config file                 | ✔    | Minimal options atm
| Scrolling                   |      | Infinite buffer implemented, need GUI scrollbar & render updates
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

