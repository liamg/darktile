# GLFW 3.3 for Go [![Build Status](https://travis-ci.org/go-gl/glfw.svg?branch=master)](https://travis-ci.org/go-gl/glfw) [![GoDoc](https://godoc.org/github.com/go-gl/glfw/v3.3/glfw?status.svg)](https://godoc.org/github.com/go-gl/glfw/v3.3/glfw)

## Installation

* GLFW C library source is included and built automatically as part of the Go package. But you need to make sure you have dependencies of GLFW:
	* On macOS, you need Xcode or Command Line Tools for Xcode (`xcode-select --install`) for required headers and libraries.
	* On Ubuntu/Debian-like Linux distributions, you need `libgl1-mesa-dev` and `xorg-dev` packages.
	* On CentOS/Fedora-like Linux distributions, you need `libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel` packages.
	* See [here](http://www.glfw.org/docs/latest/compile.html#compile_deps) for full details.
* Go 1.4+ is required on Windows (otherwise you must use MinGW v4.8.1 exactly, see [Go issue 8811](https://github.com/golang/go/issues/8811)).

```
go get -u github.com/go-gl/glfw/v3.3/glfw
```

## Usage

```Go
package main

import (
	"runtime"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	if !glfw.Init() {
		panic(glfw.GetError())
	}
	defer glfw.Terminate()

	window := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if window == nil {
		panic(glfw.GetError())
	}

	window.MakeContextCurrent()

	for !window.ShouldClose() {
		// Do OpenGL stuff.
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
```

## Changelog

* Internal error callback is now removed since GLFW now has a method called `GetError`. You can either set a custom error callback via `SetErrorCallback` or you can check for the last error via `GetError`. Due to this, some backward incompatible API changes needed to be made. See below for details.
* Joystick functions now uses receivers instead of passing the joystick ID as argument.
* Vulkan methods are intentionally not implemented. `Window.Handle` can be used to create a Vulkan surface via the [this](https://github.com/vulkan-go/vulkan) package.

### GLFW 3.3 Specific Changes
* Renamed `Window.GLFWWindow` to `Window.Handle`
* Added function `SetErrorCallback`.
* Added function `GetError`.
* Added function `Window.RequestAttention`.
* Added function `Window.SetAttrib`.
* Added function `Window.GetContentScale`.
* Added function `Window.GetOpacity`.
* Added function `Window.SetOpacity`.
* Added function `Window.SetMaximizeCallback`.
* Added function `Window.SetContentScaleCallback`.
* Added function `Monitor.GetContentScale`.
* Added function `Monitor.SetUserPointer`.
* Added function `Monitor.GetUserPointer`.
* Added function `GetKeyScancode`.
* Added function `InitHint`.
* Added function `Joystick.GetHats`.
* Added function `Joystick.IsGamepad`.
* Added function `Joystick.GetGUID`.
* Added function `Joystick.GetGamepadName`.
* Added function `Joystick.GetGamepadState`.
* Added function `Joystick.SetUserPointer`.
* Added function `Joystick.GetUserPointer`.
* Added function `UpdateGamepadMappings`.
* Added function `SetX11SelectionString`.
* Added function `GetX11SelectionString`.
* Added function `WindowHintString`.
* Added gamepad button IDs.
* Added gamepad axis IDs.
* Added joystick hat state IDs.
* Added hint `Hovered`.
* Added hint `CenterCursor`.
* Added hint `JoystickHatButtons`.
* Added hint `CocoaChdirResources`.
* Added hint `CocoaMenubar`.
* Added hint `TransparentFramebuffer`.
* Added hint value `OSMesaContextAPI`.
* `MonitorEvent` renamed to `PeripheralEvent` for reuse with joystick events.
* `Init` Returns `bool` instead of error.
* `Joystick.GetButtons` Returns `[]Action` instead of `[]byte`.
* `SetMonitorCallback` Returns `MonitorCallback`.
* `Focus` No longer returns an error.
* `Iconify` No longer returns an error.
* `Maximize` No longer returns an error.
* `Restore` No longer returns an error.
* `GetClipboardString` No longer returns an error.

### GLFW 3.2 Specfic Changes
* Added function `Window.SetSizeLimits`.
* Added function `Window.SetAspectRatio`.
* Added function `Window.SetMonitor`.
* Added function `Window.Maximize`.
* Added function `Window.SetIcon`.
* Added function `Window.Focus`.
* Added function `GetKeyName`.
* Added function `VulkanSupported`.
* Added function `GetTimerValue`.
* Added function `GetTimerFrequency`.
* Added function `WaitEventsTimeout`.
* Added function `SetJoystickCallback`.
* Added window hint `Maximized`.
* Added hint `NoAPI`.
* Added hint `NativeContextAPI`.
* Added hint `EGLContextAPI`.

### GLFW 3.1 Specfic Changes
* Added type `Cursor`.
* Added function `Window.SetDropCallback`.
* Added function `Window.SetCharModsCallback`.
* Added function `PostEmptyEvent`.
* Added function `CreateCursor`.
* Added function `CreateStandardCursor`.
* Added function `Cursor.Destroy`.
* Added function `Window.SetCursor`.
* Added function `Window.GetFrameSize`.
* Added window hint `Floating`.
* Added window hint `AutoIconify`.
* Added window hint `ContextReleaseBehavior`.
* Added window hint `DoubleBuffer`.
* Added hint value `AnyReleaseBehavior`.
* Added hint value `ReleaseBehaviorFlush`.
* Added hint value `ReleaseBehaviorNone`.
* Added hint value `DontCare`.

### API changes
* `Window.Iconify` Returns an error.
* `Window.Restore` Returns an error.
* `Init` Returns an error instead of `bool`.
* `GetJoystickAxes` No longer returns an error.
* `GetJoystickButtons` No longer returns an error.
* `GetJoystickName` No longer returns an error.
* `GetMonitors` No longer returns an error.
* `GetPrimaryMonitor` No longer returns an error.
* `Monitor.GetGammaRamp` No longer returns an error.
* `Monitor.GetVideoMode` No longer returns an error.
* `Monitor.GetVideoModes` No longer returns an error.
* `GetCurrentContext` No longer returns an error.
* `Window.SetCharCallback` Accepts `rune` instead of `uint`.
* Added type `Error`.
* Removed `SetErrorCallback`.
* Removed error code `NotInitialized`.
* Removed error code `NoCurrentContext`.
* Removed error code `InvalidEnum`.
* Removed error code `InvalidValue`.
* Removed error code `OutOfMemory`.
* Removed error code `PlatformError`.
* Removed `KeyBracket`.
* Renamed `Window.SetCharacterCallback` to `Window.SetCharCallback`.
* Renamed `Window.GetCursorPosition` to `GetCursorPos`.
* Renamed `Window.SetCursorPosition` to `SetCursorPos`.
* Renamed `CursorPositionCallback` to `CursorPosCallback`.
* Renamed `Window.SetCursorPositionCallback` to `SetCursorPosCallback`.
* Renamed `VideoMode` to `VidMode`.
* Renamed `Monitor.GetPosition` to `Monitor.GetPos`.
* Renamed `Window.GetPosition` to `Window.GetPos`.
* Renamed `Window.SetPosition` to `Window.SetPos`.
* Renamed `Window.GetAttribute` to `Window.GetAttrib`.
* Renamed `Window.SetPositionCallback` to `Window.SetPosCallback`.
* Renamed `PositionCallback` to `PosCallback`.
* Ranamed `Cursor` to `CursorMode`.
* Renamed `StickyKeys` to `StickyKeysMode`.
* Renamed `StickyMouseButtons` to `StickyMouseButtonsMode`.
* Renamed `ApiUnavailable` to `APIUnavailable`.
* Renamed `ClientApi` to `ClientAPI`.
* Renamed `OpenglForwardCompatible` to `OpenGLForwardCompatible`.
* Renamed `OpenglDebugContext` to `OpenGLDebugContext`.
* Renamed `OpenglProfile` to `OpenGLProfile`.
* Renamed `SrgbCapable` to `SRGBCapable`.
* Renamed `OpenglApi` to `OpenGLAPI`.
* Renamed `OpenglEsApi` to `OpenGLESAPI`.
* Renamed `OpenglAnyProfile` to `OpenGLAnyProfile`.
* Renamed `OpenglCoreProfile` to `OpenGLCoreProfile`.
* Renamed `OpenglCompatProfile` to `OpenGLCompatProfile`.
* Renamed `KeyKp...` to `KeyKP...`.
