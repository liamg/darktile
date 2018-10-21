package glfw

//#define GLFW_EXPOSE_NATIVE_COCOA
//#define GLFW_EXPOSE_NATIVE_NSGL
//#include "glfw/include/GLFW/glfw3.h"
//#include "glfw/include/GLFW/glfw3native.h"
//void *workaround_glfwGetCocoaWindow(GLFWwindow *w);
//void *workaround_glfwGetNSGLContext(GLFWwindow *w);
import "C"

// GetCocoaMonitor returns the CGDirectDisplayID of the monitor.
//
// This function may be called from any thread. Access is not synchronized.
func (m *Monitor) GetCocoaMonitor() uintptr {
	return uintptr(C.glfwGetCocoaMonitor(m.data))
}

// GetCocoaWindow returns the NSWindow of the window.
//
// This function may be called from any thread. Access is not synchronized.
func (w *Window) GetCocoaWindow() uintptr {
	return uintptr(C.workaround_glfwGetCocoaWindow(w.data))
}

// GetNSGLContext returns the NSOpenGLContext of the window.
//
// This function may be called from any thread. Access is not synchronized.
func (w *Window) GetNSGLContext() uintptr {
	return uintptr(C.workaround_glfwGetNSGLContext(w.data))
}
