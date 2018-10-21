#include "_cgo_export.h"

void glfwSetMonitorCallbackCB() {
	glfwSetMonitorCallback((GLFWmonitorfun) goMonitorCB);
}

GLFWmonitor *GetMonitorAtIndex(GLFWmonitor **monitors, int index) {
	return monitors[index];
}

GLFWvidmode GetVidmodeAtIndex(GLFWvidmode *vidmodes, int index) {
	return vidmodes[index];
}

unsigned int GetGammaAtIndex(unsigned short *color, int i) {
	return color[i];
}

void SetGammaAtIndex(unsigned short *color, int i, unsigned short value) {
	color[i] = value;
}