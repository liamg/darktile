package glfw

//#include "glfw/include/GLFW/glfw3.h"
import "C"

// VulkanSupported reports whether the Vulkan loader has been found. This check
// is performed by Init.
//
// The availability of a Vulkan loader does not by itself guarantee that window
// surface creation or even device creation is possible.
//
// You can use Window.Handle and vulkan-go library to create Vulkan surface
//
// This function may be called from any thread.
func VulkanSupported() bool {
	return glfwbool(C.glfwVulkanSupported())
}
