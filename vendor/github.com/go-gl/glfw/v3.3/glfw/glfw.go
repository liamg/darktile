package glfw

//#include "glfw/include/GLFW/glfw3.h"
import "C"

// Version constants.
const (
	// This is incremented when the API is changed in non-compatible ways.
	VersionMajor = C.GLFW_VERSION_MAJOR

	// This is incremented when features are added to the API but it remains
	// backward-compatible.
	VersionMinor = C.GLFW_VERSION_MINOR

	// This is incremented when a bug fix release is made that does not contain
	// any API changes.
	VersionRevision = C.GLFW_VERSION_REVISION
)

// Init initializes the GLFW library. Before most GLFW functions can be used,
// GLFW must be initialized, and before an application terminates GLFW should be
// terminated in order to free any resources allocated during or after
// initialization.
//
// If this function fails, it calls Terminate before returning. If it succeeds,
// you should call Terminate before the application exits.
//
// Additional calls to this function after successful initialization but before
// termination will return true immediately.
//
// This function must only be called from the main thread.
func Init() bool {
	return glfwbool(C.glfwInit())
}

// Terminate function destroys all remaining windows and cursors, restores any
// modified gamma ramps and frees any other allocated resources. Once this
// function is called, you must again call Init successfully before you will be
// able to use most GLFW functions.
//
// If GLFW has been successfully initialized, this function should be called
// before the application exits. If initialization fails, there is no need to
// call this function, as it is called by Init before it returns failure.
//
// This function must only be called from the main thread.
func Terminate() {
	C.glfwTerminate()
}

// InitHint function sets hints for the next initialization of GLFW.
//
// The values you set hints to are never reset by GLFW, but they only take
// effect during initialization. Once GLFW has been initialized, any values you
// set will be ignored until the library is terminated and initialized again.
//
// Some hints are platform specific. These may be set on any platform but they
// will only affect their specific platform. Other platforms will ignore them.
// Setting these hints requires no platform specific headers or functions.
//
// This function must only be called from the main thread.
func InitHint(hint Hint, value int) {
	C.glfwInitHint(C.int(hint), C.int(value))
}

// GetVersion retrieves the major, minor and revision numbers of the GLFW
// library. It is intended for when you are using GLFW as a shared library and
// want to ensure that you are using the minimum required version.
//
// This function may be called from any thread.
func GetVersion() (major, minor, revision int) {
	var (
		maj C.int
		min C.int
		rev C.int
	)

	C.glfwGetVersion(&maj, &min, &rev)
	return int(maj), int(min), int(rev)
}

// GetVersionString returns the compile-time generated version string of the
// GLFW library binary. It describes the version, platform, compiler and any
// platform-specific compile-time options. It should not be confused with the
// OpenGL or OpenGL ES version string, queried with glGetString.
//
// Do not use the version string to parse the GLFW library version.
// The GetVersion function provides the version of the running library binary in
// numerical format.
//
// This function may be called from any thread.
func GetVersionString() string {
	return C.GoString(C.glfwGetVersionString())
}
