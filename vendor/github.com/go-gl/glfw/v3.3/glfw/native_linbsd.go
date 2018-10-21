// +build linux,!wayland freebsd,!wayland

package glfw

//#define GLFW_EXPOSE_NATIVE_X11
//#define GLFW_EXPOSE_NATIVE_GLX
//#include <stdlib.h>
//#include "glfw/include/GLFW/glfw3.h"
//#include "glfw/include/GLFW/glfw3native.h"
import "C"
import "unsafe"

// GetX11Display returns X11 display handle
//
// This function may be called from any thread. Access is not synchronized.
func GetX11Display() *C.Display {
	return C.glfwGetX11Display()
}

// GetX11Adapter returns the RRCrtc of the monitor.
//
// This function may be called from any thread. Access is not synchronized.
func (m *Monitor) GetX11Adapter() C.RRCrtc {
	return C.glfwGetX11Adapter(m.data)
}

// GetX11Monitor returns the RROutput of the monitor.
//
// This function may be called from any thread. Access is not synchronized.
func (m *Monitor) GetX11Monitor() C.RROutput {
	return C.glfwGetX11Monitor(m.data)
}

// GetX11Window returns the Window of the window.
//
// This function may be called from any thread. Access is not synchronized.
func (w *Window) GetX11Window() C.Window {
	return C.glfwGetX11Window(w.data)
}

// GetGLXContext returns the GLXContext of the window.
//
// This function may be called from any thread. Access is not synchronized.
func (w *Window) GetGLXContext() C.GLXContext {
	return C.glfwGetGLXContext(w.data)
}

// GetGLXWindow returns the GLXWindow of the window.
//
// This function may be called from any thread. Access is not synchronized.
func (w *Window) GetGLXWindow() C.GLXWindow {
	return C.glfwGetGLXWindow(w.data)
}

// SetX11SelectionString sets the X11 selection string.
//
// This function must only be called from the main thread.
func SetX11SelectionString(str string) {
	s := C.CString(str)
	defer C.free(unsafe.Pointer(s))
	C.glfwSetX11SelectionString(s)
}

// GetX11SelectionString gets the X11 selection string.
//
// This function must only be called from the main thread.
func GetX11SelectionString() string {
	s := C.glfwGetX11SelectionString()
	return C.GoString(s)
}
