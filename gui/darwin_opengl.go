//+build darwin

package gui

/*
#cgo darwin CFLAGS: -x objective-c -Wno-deprecated-declarations
#cgo darwin LDFLAGS: -framework Foundation
#include <Cocoa/Cocoa.h>
void cocoa_update_nsgl_context(void* id) {
    NSOpenGLContext *ctx = id;
    [ctx update];
}
*/
import "C"

import (
	"unsafe"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var nsglContextUpdateCounter int

func UpdateNSGLContext(window *glfw.Window) {
	if nsglContextUpdateCounter < 2 {
		ctx := window.GetNSGLContext()
		C.cocoa_update_nsgl_context(unsafe.Pointer(ctx))
		nsglContextUpdateCounter++
	}
}
