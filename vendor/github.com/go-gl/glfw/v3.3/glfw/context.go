package glfw

//#include <stdlib.h>
//#include "glfw/include/GLFW/glfw3.h"
import "C"

import (
	"unsafe"
)

// MakeContextCurrent makes the OpenGL or OpenGL ES context of the specified
// window current on the calling thread. A context must only be made current on
// a single thread at a time and each thread can have only a single current
// context at a time.
//
// When moving a context between threads, you must make it non-current on the
// old thread before making it current on the new one.
//
// By default, making a context non-current implicitly forces a pipeline flush.
// On machines that support GL_KHR_context_flush_control, you can control
// whether a context performs this flush by setting the ContextReleaseBehavior
// hint.
//
// The specified window must have an OpenGL or OpenGL ES context. Specifying a
// window without a context will generate a NoWindowContext error.
//
// This function may be called from any thread.
func (w *Window) MakeContextCurrent() {
	C.glfwMakeContextCurrent(w.data)
}

// DetachCurrentContext detaches the current context.
//
// This function may be called from any thread.
func DetachCurrentContext() {
	C.glfwMakeContextCurrent(nil)
}

// GetCurrentContext returns the window whose OpenGL or OpenGL ES context is
// current on the calling thread.
//
// This function may be called from any thread.
func GetCurrentContext() *Window {
	w := C.glfwGetCurrentContext()
	if w == nil {
		return nil
	}
	return windows.get(w)
}

// SwapInterval sets the swap interval for the current OpenGL or OpenGL ES
// context, i.e. the number of screen updates to wait from the time SwapBuffers
// was called before swapping the buffers and returning. This is sometimes
// called vertical synchronization, vertical retrace synchronization or just
// vsync.
//
// A context that supports either of the WGL_EXT_swap_control_tear and
// GLX_EXT_swap_control_tear extensions also accepts negative swap intervals,
// which allows the driver to swap immediately even if a frame arrives a little
// bit late. You can check for these extensions with ExtensionSupported.
//
// A context must be current on the calling thread. Calling this function
// without a current context will cause a NoCurrentContext error.
//
// This function does not apply to Vulkan. If you are rendering with Vulkan,
// see the present mode of your swapchain instead.
//
// This function may be called from any thread.
func SwapInterval(interval int) {
	C.glfwSwapInterval(C.int(interval))
}

// ExtensionSupported returns whether the specified API extension is supported
// by the current OpenGL or OpenGL ES context. It searches both for client API
// extension and context creation API extensions.
//
// A context must be current on the calling thread. Calling this function
// without a current context will cause a NoCurrentContext error.
//
// As this functions retrieves and searches one or more extension strings each
// call, it is recommended that you cache its results if it is going to be used
// frequently. The extension strings will not change during the lifetime of a
// context, so there is no danger in doing this.
//
// This function does not apply to Vulkan.
//
// This function may be called from any thread.
func ExtensionSupported(extension string) bool {
	e := C.CString(extension)
	defer C.free(unsafe.Pointer(e))
	return glfwbool(C.glfwExtensionSupported(e))
}
