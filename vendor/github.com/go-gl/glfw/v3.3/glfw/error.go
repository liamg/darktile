package glfw

//#include "glfw/include/GLFW/glfw3.h"
//void glfwSetErrorCallbackCB();
import "C"

import (
	"fmt"
)

// ErrorCallback is the function signature for error callback functions.
type ErrorCallback func(code ErrorCode, desc string)

var fErrorHolder ErrorCallback

//export goErrorCB
func goErrorCB(code C.int, desc *C.char) {
	fErrorHolder(ErrorCode(code), C.GoString(desc))
}

// ErrorCode corresponds to an error code.
type ErrorCode int

const (
	// NoError occurs when there are no errors.
	NoError ErrorCode = C.GLFW_NO_ERROR

	// NotInitialized occurs if a GLFW function was called that must not be
	// called unless the library is initialized.
	//
	// Application programmer error. Initialize GLFW before calling any function
	// that requires initialization.
	NotInitialized ErrorCode = C.GLFW_NOT_INITIALIZED

	// NoCurrentContext occurs if a GLFW function was called that needs and
	// operates on the current OpenGL or OpenGL ES context but no context is
	// current on the calling thread. One such function is SwapInterval.
	//
	// Application programmer error. Ensure a context is current before calling
	// functions that require a current context.
	NoCurrentContext ErrorCode = C.GLFW_NO_CURRENT_CONTEXT

	// InvalidEnum occurs when one of the arguments to the function is an
	// invalid enum value, for example requesting RedBits with GetWindowAttrib.
	//
	// Application programmer error. Fix the offending call.
	InvalidEnum ErrorCode = C.GLFW_INVALID_ENUM

	// InvalidValue occurs when one of the arguments to the function is an
	// invalid value, for example requesting a non-existent OpenGL or OpenGL ES
	// version like 2.7.
	//
	// Requesting a valid but unavailable OpenGL or OpenGL ES version will
	// instead result in a VersionUnavailable error.
	//
	// Application programmer error. Fix the offending call.
	InvalidValue ErrorCode = C.GLFW_INVALID_VALUE

	// OutOfMemory occurs when a memory allocation fails.
	//
	// A bug in GLFW or the underlying operating system. Report the bug.
	OutOfMemory ErrorCode = C.GLFW_OUT_OF_MEMORY

	// APIUnavailable occurs when GLFW can't find support for the requested API
	// on the system.
	//
	// The installed graphics driver does not support the requested API, or does
	// not support it via the chosen context creation backend. Below are a few
	// examples.
	//
	// Some pre-installed Windows graphics drivers do not support OpenGL. AMD
	// only supports OpenGL ES via EGL, while Nvidia and Intel only support it
	// via a WGL or GLX extension. macOS does not provide OpenGL ES at all.
	// The Mesa EGL, OpenGL and OpenGL ES libraries do not interface with the
	// Nvidia binary driver. Older graphics drivers do not support Vulkan.
	APIUnavailable ErrorCode = C.GLFW_API_UNAVAILABLE

	// VersionUnavailable occurs when the requested OpenGL or OpenGL ES version
	// (including any requested context or framebuffer hints) is not available
	// on this machine.
	//
	// The machine does not support your requirements. If your application is
	// sufficiently flexible, downgrade your requirements and try again.
	// Otherwise, inform the user that their machine does not match your
	// requirements.
	//
	// Future invalid OpenGL and OpenGL ES versions, for example OpenGL 4.8 if
	// 5.0 comes out before the 4.x series gets that far, also fail with this
	// error and not InvalidValue, because GLFW cannot know what future versions
	// will exist.
	VersionUnavailable ErrorCode = C.GLFW_VERSION_UNAVAILABLE

	// PlatformError occurs when a platform-specific error occurs that does not
	// match any of the more specific categories.
	//
	// A bug or configuration error in GLFW, the underlying operating system or
	// its drivers, or a lack of required resources. Report the issue.
	PlatformError ErrorCode = C.GLFW_PLATFORM_ERROR

	// FormatUnavailable if emitted during window creation, means that the
	// requested pixel format is not supported.
	//
	// If emitted when querying the clipboard, the contents of the clipboard
	// could not be converted to the requested format.
	//
	// If emitted during window creation, one or more hard constraints did not
	// match any of the available pixel formats. If your application is
	// sufficiently flexible, downgrade your requirements and try again.
	// Otherwise, inform the user that their machine does not match your
	// requirements.
	//
	// If emitted when querying the clipboard, ignore the error or report it to
	// the user, as appropriate.
	FormatUnavailable ErrorCode = C.GLFW_FORMAT_UNAVAILABLE

	// NoWindowContext occurs when a window that does not have an OpenGL or
	// OpenGL ES context was passed to a function that requires it to have one.
	//
	// Application programmer error. Fix the offending call.
	NoWindowContext ErrorCode = C.GLFW_NO_WINDOW_CONTEXT
)

func (e ErrorCode) String() string {
	switch e {
	case NoError:
		return "NoError"
	case NotInitialized:
		return "NotInitialized"
	case NoCurrentContext:
		return "NoCurrentContext"
	case InvalidEnum:
		return "InvalidEnum"
	case InvalidValue:
		return "InvalidValue"
	case OutOfMemory:
		return "OutOfMemory"
	case APIUnavailable:
		return "APIUnavailable"
	case VersionUnavailable:
		return "VersionUnavailable"
	case PlatformError:
		return "PlatformError"
	case FormatUnavailable:
		return "FormatUnavailable"
	case NoWindowContext:
		return "NoWindowContext"
	default:
		return fmt.Sprintf("ErrorCode(%d)", e)
	}
}

// Error represents a GLFW error.
type Error struct {
	Code        ErrorCode
	Description string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s\n", e.Code, e.Description)
}

// GetError returns and clears the error code of the last error that occurred on
// the calling thread, and optionally a UTF-8 encoded human-readable description
// of it. If no error has occurred since the last call, it returns NoError and
// the description is set to empty string.
//
// This function may be called from any thread.
func GetError() error {
	var (
		desc *C.char
		e    = ErrorCode(C.glfwGetError(&desc))
	)
	if e == NoError {
		return nil
	}
	return &Error{
		Code:        e,
		Description: C.GoString(desc),
	}
}

// SetErrorCallback function sets the error callback, which is called with an
// error code and a human-readable description each time a GLFW error occurs.
//
// The error code is set before the callback is called. Calling GetError from
// the error callback will return the same value as the error code argument.
//
// The error callback is called on the thread where the error occurred. If you
// are using GLFW from multiple threads, your error callback needs to be written
// accordingly.
//
// Because the description string may have been generated specifically for that
// error, it is not guaranteed to be valid after the callback has returned. If
// you wish to use it after the callback returns, you need to make a copy.
//
// Once set, the error callback remains set even after the library has been
// terminated.
//
// This function must only be called from the main thread.
func SetErrorCallback(cbfun ErrorCallback) ErrorCallback {
	previous := fErrorHolder
	fErrorHolder = cbfun
	if cbfun == nil {
		C.glfwSetErrorCallback(nil)
	} else {
		C.glfwSetErrorCallbackCB()
	}
	return previous
}
