package glfw

//#define GLFW_EXPOSE_NATIVE_WIN32
//#define GLFW_EXPOSE_NATIVE_WGL
//#include "glfw/include/GLFW/glfw3.h"
//#include "glfw/include/GLFW/glfw3native.h"
import "C"

// GetWin32Adapter returns the adapter device name of the monitor.
//
// This function may be called from any thread. Access is not synchronized.
func (m *Monitor) GetWin32Adapter() string {
	ret := C.glfwGetWin32Adapter(m.data)
	return C.GoString(ret)
}

// GetWin32Monitor returns the display device name of the monitor.
//
// This function may be called from any thread. Access is not synchronized.
func (m *Monitor) GetWin32Monitor() string {
	ret := C.glfwGetWin32Monitor(m.data)
	return C.GoString(ret)
}

// GetWin32Window returns the HWND of the window.
//
// This function may be called from any thread. Access is not synchronized.
func (w *Window) GetWin32Window() C.HWND {
	ret := C.glfwGetWin32Window(w.data)
	return ret
}

// GetWGLContext returns the HGLRC of the window.
//
// This function may be called from any thread. Access is not synchronized.
func (w *Window) GetWGLContext() C.HGLRC {
	ret := C.glfwGetWGLContext(w.data)
	return ret
}
