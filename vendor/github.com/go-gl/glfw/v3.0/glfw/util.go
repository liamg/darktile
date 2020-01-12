package glfw

//#define GLFW_INCLUDE_NONE
//#include <GLFW/glfw3.h>
import "C"

func glfwbool(b C.int) bool {
	if b == C.int(True) {
		return true
	}
	return false
}
