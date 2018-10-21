package glfw

//#include <stdlib.h>
//#include "glfw/include/GLFW/glfw3.h"
//void glfwSetWindowPosCallbackCB(GLFWwindow *window);
//void glfwSetWindowSizeCallbackCB(GLFWwindow *window);
//void glfwSetFramebufferSizeCallbackCB(GLFWwindow *window);
//void glfwSetWindowCloseCallbackCB(GLFWwindow *window);
//void glfwSetWindowRefreshCallbackCB(GLFWwindow *window);
//void glfwSetWindowFocusCallbackCB(GLFWwindow *window);
//void glfwSetWindowIconifyCallbackCB(GLFWwindow *window);
//void glfwSetWindowMaximizeCallbackCB(GLFWwindow *window);
//void glfwSetWindowContentScaleCallbackCB(GLFWwindow *window);
import "C"

import (
	"image"
	"image/draw"
	"sync"
	"unsafe"
)

// Internal window list stuff
type windowList struct {
	l sync.Mutex
	m map[*C.GLFWwindow]*Window
}

var windows = windowList{m: map[*C.GLFWwindow]*Window{}}

func (w *windowList) put(wnd *Window) {
	w.l.Lock()
	defer w.l.Unlock()
	w.m[wnd.data] = wnd
}

func (w *windowList) remove(wnd *C.GLFWwindow) {
	w.l.Lock()
	defer w.l.Unlock()
	delete(w.m, wnd)
}

func (w *windowList) get(wnd *C.GLFWwindow) *Window {
	w.l.Lock()
	defer w.l.Unlock()
	return w.m[wnd]
}

// Hint can correspond to hints that can be set before initializing GLFW,
// creating a window or to the attributes of the window that can be get after
// its creation.
type Hint int

// Init related hints.
const (
	JoystickHatButtons  Hint = C.GLFW_JOYSTICK_HAT_BUTTONS
	CocoaChdirResources Hint = C.GLFW_COCOA_CHDIR_RESOURCES
	CocoaMenubar        Hint = C.GLFW_COCOA_MENUBAR
)

// Window related hints/attribs.
const (
	// Specifies whether the window will be given input focus when created.
	// This hint is ignored for full screen and initially hidden windows.
	Focused Hint = C.GLFW_FOCUSED

	// Specifies whether the window will be minimized.
	Iconified Hint = C.GLFW_ICONIFIED

	// Specifies whether the window is maximized.
	Maximized Hint = C.GLFW_MAXIMIZED

	// Specifies whether the cursor is currently directly over the client area
	// of the window, with no other windows between
	Hovered Hint = C.GLFW_HOVERED

	// Specifies whether the window will be initially visible.
	Visible Hint = C.GLFW_VISIBLE

	// Specifies whether the window will be resizable by the user.
	Resizable Hint = C.GLFW_RESIZABLE

	// Specifies whether the window will have window decorations such as a
	// border, a close widget, etc.
	Decorated Hint = C.GLFW_DECORATED

	// Specifies whether the window will be always-on-top.
	Floating Hint = C.GLFW_FLOATING

	// Specifies whether fullscreen windows automatically iconify (and restore
	// the previous video mode) on focus loss.
	AutoIconify Hint = C.GLFW_AUTO_ICONIFY

	// Specifies whether the cursor should be centered over newly created full
	// screen windows. This hint is ignored for windowed mode windows.
	CenterCursor Hint = C.GLFW_CENTER_CURSOR
)

// Context related hints.
const (
	// Specifies which client API to create the context for. Hard constraint.
	ClientAPI Hint = C.GLFW_CLIENT_API

	// Specifies the client API version that the created context must be
	// compatible with.
	ContextVersionMajor Hint = C.GLFW_CONTEXT_VERSION_MAJOR

	// Specifies the client API version that the created context must be
	// compatible with.
	ContextVersionMinor Hint = C.GLFW_CONTEXT_VERSION_MINOR

	// Specifies the robustness strategy to be used by the context.
	ContextRobustness Hint = C.GLFW_CONTEXT_ROBUSTNESS

	// Specifies the release behavior to be used by the context.
	ContextReleaseBehavior Hint = C.GLFW_CONTEXT_RELEASE_BEHAVIOR

	// Specifies whether the OpenGL context should be forward-compatible.
	// Hard constraint.
	OpenGLForwardCompatible Hint = C.GLFW_OPENGL_FORWARD_COMPAT

	// Specifies whether to create a debug OpenGL context, which may have
	// additional error and performance issue reporting functionality. If
	// OpenGL ES is requested, this hint is ignored.
	OpenGLDebugContext Hint = C.GLFW_OPENGL_DEBUG_CONTEXT

	// Specifies which OpenGL profile to create the context for.
	// Hard constraint.
	OpenGLProfile Hint = C.GLFW_OPENGL_PROFILE

	// Specifies which context creation API to use to create the context.
	ContextCreationAPI Hint = C.GLFW_CONTEXT_CREATION_API
)

// Framebuffer related hints.
const (
	ContextRevision Hint = C.GLFW_CONTEXT_REVISION

	// Specifies the desired bit depth of the default framebuffer.
	RedBits Hint = C.GLFW_RED_BITS

	// Specifies the desired bit depth of the default framebuffer.
	GreenBits Hint = C.GLFW_GREEN_BITS

	// Specifies the desired bit depth of the default framebuffer.
	BlueBits Hint = C.GLFW_BLUE_BITS

	// Specifies the desired bit depth of the default framebuffer.
	AlphaBits Hint = C.GLFW_ALPHA_BITS

	// Specifies the desired bit depth of the default framebuffer.
	DepthBits Hint = C.GLFW_DEPTH_BITS

	// Specifies the desired bit depth of the default framebuffer.
	StencilBits Hint = C.GLFW_STENCIL_BITS

	// Specifies the desired bit depth of the accumulation buffer.
	AccumRedBits Hint = C.GLFW_ACCUM_RED_BITS

	// Specifies the desired bit depth of the accumulation buffer.
	AccumGreenBits Hint = C.GLFW_ACCUM_GREEN_BITS

	// Specifies the desired bit depth of the accumulation buffer.
	AccumBlueBits Hint = C.GLFW_ACCUM_BLUE_BITS

	// Specifies the desired bit depth of the accumulation buffer.
	AccumAlphaBits Hint = C.GLFW_ACCUM_ALPHA_BITS

	// Specifies the desired number of auxiliary buffers.
	AuxBuffers Hint = C.GLFW_AUX_BUFFERS

	// Specifies whether to use stereoscopic rendering. Hard constraint.
	Stereo Hint = C.GLFW_STEREO

	// Specifies the desired number of samples to use for multisampling. Zero
	// disables multisampling.
	Samples Hint = C.GLFW_SAMPLES

	// Specifies whether the framebuffer should be sRGB capable.
	SRGBCapable Hint = C.GLFW_SRGB_CAPABLE

	// Specifies the desired refresh rate for full screen windows. If set to
	// zero, the highest available refresh rate will be used. This hint is
	// ignored for windowed mode windows.
	RefreshRate Hint = C.GLFW_REFRESH_RATE

	// Specifies whether the framebuffer should be double buffered. You nearly
	// always want to use double buffering. This is a hard constraint.
	DoubleBuffer Hint = C.GLFW_DOUBLEBUFFER

	// Specifies whether the framebuffer should be transparent.
	TransparentFramebuffer Hint = C.GLFW_TRANSPARENT_FRAMEBUFFER
)

// Values for the ClientAPI hint.
const (
	OpenGLAPI   int = C.GLFW_OPENGL_API
	OpenGLESAPI int = C.GLFW_OPENGL_ES_API
	NoAPI       int = C.GLFW_NO_API
)

// Values for ContextCreationAPI hint.
const (
	NativeContextAPI int = C.GLFW_NATIVE_CONTEXT_API
	EGLContextAPI    int = C.GLFW_EGL_CONTEXT_API
	OSMesaContextAPI int = C.GLFW_OSMESA_CONTEXT_API
)

// Values for the ContextRobustness hint.
const (
	NoRobustness        int = C.GLFW_NO_ROBUSTNESS
	NoResetNotification int = C.GLFW_NO_RESET_NOTIFICATION
	LoseContextOnReset  int = C.GLFW_LOSE_CONTEXT_ON_RESET
)

// Values for ContextReleaseBehavior hint.
const (
	AnyReleaseBehavior   int = C.GLFW_ANY_RELEASE_BEHAVIOR
	ReleaseBehaviorFlush int = C.GLFW_RELEASE_BEHAVIOR_FLUSH
	ReleaseBehaviorNone  int = C.GLFW_RELEASE_BEHAVIOR_NONE
)

// Values for the OpenGLProfile hint.
const (
	OpenGLAnyProfile    int = C.GLFW_OPENGL_ANY_PROFILE
	OpenGLCoreProfile   int = C.GLFW_OPENGL_CORE_PROFILE
	OpenGLCompatProfile int = C.GLFW_OPENGL_COMPAT_PROFILE
)

// Other values.
const (
	True     int = C.GL_TRUE
	False    int = C.GL_FALSE
	DontCare int = C.GLFW_DONT_CARE
)

// PosCallback is the function signature for window position callback functions.
type PosCallback func(w *Window, xpos int, ypos int)

// SizeCallback is the function signature for window size callback functions.
type SizeCallback func(w *Window, width int, height int)

// CloseCallback is the function signature for window close callback functions.
type CloseCallback func(w *Window)

// RefreshCallback is the function signature for window refresh callback
// functions.
type RefreshCallback func(w *Window)

// FocusCallback is the function signature for window focus callback functions.
type FocusCallback func(w *Window, focused bool)

// IconifyCallback is the function signature for window iconification callback
// functions.
type IconifyCallback func(w *Window, iconified bool)

// MaximizeCallback is the function signature for window maximize callback
// functions.
type MaximizeCallback func(w *Window, iconified bool)

// FramebufferSizeCallback is the function signature for framebuffer size
// callback functions.
type FramebufferSizeCallback func(w *Window, width int, height int)

// ContentScaleCallback is the function signature for window content scale
// callback functions.
type ContentScaleCallback func(w *Window, x float32, y float32)

// Window represents a window.
type Window struct {
	data *C.GLFWwindow

	// Window
	fPosHolder             PosCallback
	fSizeHolder            SizeCallback
	fFramebufferSizeHolder FramebufferSizeCallback
	fCloseHolder           CloseCallback
	fMaximizeHolder        MaximizeCallback
	fRefreshHolder         RefreshCallback
	fFocusHolder           FocusCallback
	fIconifyHolder         IconifyCallback
	fContentScaleHolder    ContentScaleCallback

	// Input
	fMouseButtonHolder MouseButtonCallback
	fCursorPosHolder   CursorPosCallback
	fCursorEnterHolder CursorEnterCallback
	fScrollHolder      ScrollCallback
	fKeyHolder         KeyCallback
	fCharHolder        CharCallback
	fCharModsHolder    CharModsCallback
	fDropHolder        DropCallback
}

//export goWindowPosCB
func goWindowPosCB(window unsafe.Pointer, xpos, ypos C.int) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fPosHolder(w, int(xpos), int(ypos))
}

//export goWindowSizeCB
func goWindowSizeCB(window unsafe.Pointer, width, height C.int) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fSizeHolder(w, int(width), int(height))
}

//export goFramebufferSizeCB
func goFramebufferSizeCB(window unsafe.Pointer, width, height C.int) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fFramebufferSizeHolder(w, int(width), int(height))
}

//export goWindowCloseCB
func goWindowCloseCB(window unsafe.Pointer) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fCloseHolder(w)
}

//export goWindowMaximizeCB
func goWindowMaximizeCB(window unsafe.Pointer, iconified C.int) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fMaximizeHolder(w, glfwbool(iconified))
}

//export goWindowRefreshCB
func goWindowRefreshCB(window unsafe.Pointer) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fRefreshHolder(w)
}

//export goWindowFocusCB
func goWindowFocusCB(window unsafe.Pointer, focused C.int) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fFocusHolder(w, glfwbool(focused))
}

//export goWindowIconifyCB
func goWindowIconifyCB(window unsafe.Pointer, iconified C.int) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fIconifyHolder(w, glfwbool(iconified))
}

//export goWindowContentScaleCB
func goWindowContentScaleCB(window unsafe.Pointer, x C.float, y C.float) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fContentScaleHolder(w, float32(x), float32(y))
}

// DefaultWindowHints resets all window hints to their default values.
//
// This function must only be called from the main thread.
func DefaultWindowHints() {
	C.glfwDefaultWindowHints()
}

// WindowHint sets hints for the next call to CreateWindow. The hints,
// once set, retain their values until changed by a call to WindowHint or
// DefaultWindowHints, or until the library is terminated with Terminate.
//
// Only integer value hints can be set with this function. String value hints
// are set with WindowHintString.
//
// This function does not check whether the specified hint values are valid. If
// you set hints to invalid values this will instead be reported by the next
// call to CreateWindow.
//
// Some hints are platform specific. These may be set on any platform but they
// will only affect their specific platform. Other platforms will ignore them.
// Setting these hints requires no platform specific headers or functions.
//
// This function must only be called from the main thread.
func WindowHint(hint Hint, value int) {
	C.glfwWindowHint(C.int(hint), C.int(value))
}

// WindowHintString sets hints for the next call to CreateWindow. The hints,
// once set, retain their values until changed by a call to this function or
// DefaultWindowHints, or until the library is terminated.
//
// Only string type hints can be set with this function. Integer value hints are
// set with WindowHint.
//
// This function does not check whether the specified hint values are valid. If
// you set hints to invalid values this will instead be reported by the next
// call to CreateWindow.
//
// Some hints are platform specific. These may be set on any platform but they
// will only affect their specific platform. Other platforms will ignore them.
// Setting these hints requires no platform specific headers or functions.
//
// This function must only be called from the main thread.
func WindowHintString(hint Hint, value string) {
	str := C.CString(value)
	defer C.free(unsafe.Pointer(str))
	C.glfwWindowHintString(C.int(hint), str)
}

// CreateWindow creates a window and its associated OpenGL or OpenGL ES context.
// Most of the options controlling how the window and its context should be
// created are specified with window hints.
//
// Successful creation does not change which context is current. Before you can
// use the newly created context, you need to make it current. For information
// about the share parameter, see Context object sharing.
//
// The created window, framebuffer and context may differ from what you
// requested, as not all parameters and hints are hard constraints. This
// includes the size of the window, especially for full screen windows. To query
// the actual attributes of the created window, framebuffer and context, see
// GetAttrib, GetSize and GetFramebufferSize.
//
// To create a full screen window, you need to specify the monitor the window
// will cover. If no monitor is specified, the window will be windowed mode.
// Unless you have a way for the user to choose a specific monitor, it is
// recommended that you pick the primary monitor. For more information on how to
// query connected monitors, see Retrieving monitors.
//
// For full screen windows, the specified size becomes the resolution of the
// window's desired video mode. As long as a full screen window is not
// iconified, the supported video mode most closely matching the desired video
// mode is set for the specified monitor. For more information about full screen
// windows, including the creation of so called windowed full screen or
// borderless full screen windows, see "Windowed full screen" windows.
//
// Once you have created the window, you can switch it between windowed and full
// screen mode with SetMonitor. This will not affect its OpenGL or OpenGL ES
// context.
//
// By default, newly created windows use the placement recommended by the window
// system. To create the window at a specific position, make it initially
// invisible using the Visible window hint, set its position and then show it.
//
// As long as at least one full screen window is not iconified, the screensaver
// is prohibited from starting.
//
// Window systems put limits on window sizes. Very large or very small window
// dimensions may be overridden by the window system on creation. Check the
// actual size after creation.
//
// The swap interval is not set during window creation and the initial value may
// vary depending on driver settings and defaults.
//
// This function must only be called from the main thread.
func CreateWindow(
	width, height int,
	title string,
	monitor *Monitor,
	share *Window,
) *Window {
	var (
		m *C.GLFWmonitor
		s *C.GLFWwindow
	)

	t := C.CString(title)
	defer C.free(unsafe.Pointer(t))

	if monitor != nil {
		m = monitor.data
	}

	if share != nil {
		s = share.data
	}

	w := C.glfwCreateWindow(C.int(width), C.int(height), t, m, s)
	if w == nil {
		return nil
	}

	wnd := &Window{data: w}
	windows.put(wnd)
	return wnd
}

// Destroy destroys the specified window and its context. On calling this
// function, no further callbacks will be called for that window.
//
// If the context of the specified window is current on the main thread, it is
// detached before being destroyed.
//
// This function must only be called from the main thread.
func (w *Window) Destroy() {
	windows.remove(w.data)
	C.glfwDestroyWindow(w.data)
}

// ShouldClose reports the value of the close flag of the specified window.
//
// This function may be called from any thread. Access is not synchronized.
func (w *Window) ShouldClose() bool {
	return glfwbool(C.glfwWindowShouldClose(w.data))
}

// SetShouldClose sets the value of the close flag of the window. This can be
// used to override the user's attempt to close the window, or to signal that it
// should be closed.
//
// This function may be called from any thread. Access is not synchronized.
func (w *Window) SetShouldClose(value bool) {
	if !value {
		C.glfwSetWindowShouldClose(w.data, C.GL_FALSE)
	} else {
		C.glfwSetWindowShouldClose(w.data, C.GL_TRUE)
	}
}

// SetTitle sets the window title, encoded as UTF-8, of the window.
//
// This function must only be called from the main thread.
func (w *Window) SetTitle(title string) {
	t := C.CString(title)
	defer C.free(unsafe.Pointer(t))
	C.glfwSetWindowTitle(w.data, t)
}

// SetIcon sets the icon of the specified window. If passed an array of
// candidate images, those of or closest to the sizes desired by the system are
// selected. If no images are specified, the window reverts to its default icon.
//
// The pixels are 32-bit, little-endian, non-premultiplied RGBA, i.e. eight bits
// per channel with the red channel first. They are arranged canonically as
// packed sequential rows, starting from the top-left corner.
//
// The desired image sizes varies depending on platform and system settings. The
// selected images will be rescaled as needed. Good sizes include 16x16, 32x32
// and 48x48.
//
// This function must only be called from the main thread.
func (w *Window) SetIcon(images []image.Image) {
	count := len(images)
	cimages := make([]C.GLFWimage, count)
	freePixels := make([]func(), count)

	for i, img := range images {
		var pixels []uint8
		b := img.Bounds()

		switch img := img.(type) {
		case *image.NRGBA:
			pixels = img.Pix
		default:
			m := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
			draw.Draw(m, m.Bounds(), img, b.Min, draw.Src)
			pixels = m.Pix
		}

		pix, free := bytes(pixels)
		freePixels[i] = free

		cimages[i].width = C.int(b.Dx())
		cimages[i].height = C.int(b.Dy())
		cimages[i].pixels = (*C.uchar)(pix)
	}

	var p *C.GLFWimage
	if count > 0 {
		p = &cimages[0]
	}
	C.glfwSetWindowIcon(w.data, C.int(count), p)

	for _, v := range freePixels {
		v()
	}
}

// GetPos returns the position, in screen coordinates, of the upper-left
// corner of the client area of the window.
//
// This function must only be called from the main thread.
func (w *Window) GetPos() (x, y int) {
	var xpos, ypos C.int
	C.glfwGetWindowPos(w.data, &xpos, &ypos)
	return int(xpos), int(ypos)
}

// SetPos sets the position, in screen coordinates, of the upper-left corner of
// the client area of the specified windowed mode window. If the window is a
// full screen window, this function does nothing.
//
// Do not use this function to move an already visible window unless you have
// very good reasons for doing so, as it will confuse and annoy the user.
//
// The window manager may put limits on what positions are allowed. GLFW cannot
// and should not override these limits.
//
// This function must only be called from the main thread.
func (w *Window) SetPos(xpos, ypos int) {
	C.glfwSetWindowPos(w.data, C.int(xpos), C.int(ypos))
}

// GetSize retrieves the size, in screen coordinates, of the client area of the
// specified window. If you wish to retrieve the size of the framebuffer of the
// window in pixels, see GetFramebufferSize.
//
// This function must only be called from the main thread.
func (w *Window) GetSize() (width, height int) {
	var wi, h C.int
	C.glfwGetWindowSize(w.data, &wi, &h)
	return int(wi), int(h)
}

// SetSizeLimits sets the size limits of the client area of the specified
// window. If the window is full screen, the size limits only take effect once
// it is made windowed. If the window is not resizable, this function does
// nothing.
//
// The size limits are applied immediately to a windowed mode window and may
// cause it to be resized.
//
// The maximum dimensions must be greater than or equal to the minimum
// dimensions and all must be greater than or equal to zero.
//
// This function may only be called from the main thread.
func (w *Window) SetSizeLimits(minw, minh, maxw, maxh int) {
	C.glfwSetWindowSizeLimits(
		w.data,
		C.int(minw),
		C.int(minh),
		C.int(maxw),
		C.int(maxh),
	)
}

// SetAspectRatio sets the required aspect ratio of the client area of the
// specified window. If the window is full screen, the aspect ratio only takes
// effect once it is made windowed. If the window is not resizable, this
// function does nothing.
//
// The aspect ratio is specified as a numerator and a denominator and both
// values must be greater than zero. For example, the common 16:9 aspect ratio
// is specified as 16 and 9, respectively.
//
// If the numerator and denominator is set to DontCare then the aspect ratio
// limit is disabled.
//
// The aspect ratio is applied immediately to a windowed mode window and may
// cause it to be resized.
//
// This function may only be called from the main thread.
func (w *Window) SetAspectRatio(numer, denom int) {
	C.glfwSetWindowAspectRatio(w.data, C.int(numer), C.int(denom))
}

// SetSize sets the size, in screen coordinates, of the client area of the
// specified window.
//
// For full screen windows, this function updates the resolution of its desired
// video mode and switches to the video mode closest to it, without affecting
// the window's context. As the context is unaffected, the bit depths of the
// framebuffer remain unchanged.
//
// If you wish to update the refresh rate of the desired video mode in addition
// to its resolution, see SetMonitor.
//
// The window manager may put limits on what sizes are allowed. GLFW cannot and
// should not override these limits.
//
// This function may only be called from the main thread.
func (w *Window) SetSize(width, height int) {
	C.glfwSetWindowSize(w.data, C.int(width), C.int(height))
}

// GetFramebufferSize retrieves the size, in pixels, of the framebuffer of the
// specified window. If you wish to retrieve the size of the window in screen
// coordinates, see GetSize.
//
// This function may only be called from the main thread.
func (w *Window) GetFramebufferSize() (width, height int) {
	var wi, h C.int
	C.glfwGetFramebufferSize(w.data, &wi, &h)
	return int(wi), int(h)
}

// GetFrameSize retrieves the size, in screen coordinates, of each edge of the
// frame of the specified window. This size includes the title bar, if the
// window has one. The size of the frame may vary depending on the
// window-related hints used to create it.
//
// Because this function retrieves the size of each window frame edge and not
// the offset along a particular coordinate axis, the retrieved values will
// always be zero or positive.
//
// This function may only be called from the main thread.
func (w *Window) GetFrameSize() (left, top, right, bottom int) {
	var l, t, r, b C.int
	C.glfwGetWindowFrameSize(w.data, &l, &t, &r, &b)
	return int(l), int(t), int(r), int(b)
}

// GetContentScale function retrieves the content scale for the specified
// window. The content scale is the ratio between the current DPI and the
// platform's default DPI. If you scale all pixel dimensions by this scale then
// your content should appear at an appropriate size. This is especially
// important for text and any UI elements.
//
// This function may only be called from the main thread.
func (w *Window) GetContentScale() (float32, float32) {
	var x, y C.float
	C.glfwGetWindowContentScale(w.data, &x, &y)
	return float32(x), float32(y)
}

// GetOpacity function returns the opacity of the window, including any
// decorations.
//
// The opacity (or alpha) value is a positive finite number between zero and
// one, where zero is fully transparent and one is fully opaque. If the system
// does not support whole window transparency, this function always returns one.
//
// The initial opacity value for newly created windows is one.
//
// This function may only be called from the main thread.
func (w *Window) GetOpacity() float32 {
	return float32(C.glfwGetWindowOpacity(w.data))
}

// SetOpacity function sets the opacity of the window, including any
// decorations. The opacity (or alpha) value is a positive finite number between
// zero and one, where zero is fully transparent and one is fully opaque.
//
// The initial opacity value for newly created windows is one.
//
// A window created with framebuffer transparency may not use whole window
// transparency. The results of doing this are undefined.
//
// This function may only be called from the main thread.
func (w *Window) SetOpacity(opacity float32) {
	C.glfwSetWindowOpacity(w.data, C.float(opacity))
}

// Iconify iconifies (minimizes) the specified window if it was previously
// restored. If the window is already iconified, this function does nothing.
//
// If the specified window is a full screen window, the original monitor
// resolution is restored until the window is restored.
//
// This function may only be called from the main thread.
func (w *Window) Iconify() {
	C.glfwIconifyWindow(w.data)
}

// Restore restores the specified window if it was previously iconified
// (minimized) or maximized. If the window is already restored, this function
// does nothing.
//
// If the specified window is a full screen window, the resolution chosen for
// the window is restored on the selected monitor.
//
// This function may only be called from the main thread.
func (w *Window) Restore() {
	C.glfwRestoreWindow(w.data)
}

// Maximize maximizes the specified window if it was previously not maximized.
// If the window is already maximized, this function does nothing.
//
// If the specified window is a full screen window, this function does nothing.
//
// This function may only be called from the main thread.
func (w *Window) Maximize() {
	C.glfwMaximizeWindow(w.data)
}

// Show makes the specified window visible if it was previously hidden. If the
// window is already visible or is in full screen mode, this function does
// nothing.
//
// This function may only be called from the main thread.
func (w *Window) Show() {
	C.glfwShowWindow(w.data)
}

// Hide hides the specified window if it was previously visible. If the window
// is already hidden or is in full screen mode, this function does nothing.
//
// This function may only be called from the main thread.
func (w *Window) Hide() {
	C.glfwHideWindow(w.data)
}

// Focus brings the specified window to front and sets input focus. The window
// should already be visible and not iconified.
//
// By default, both windowed and full screen mode windows are focused when
// initially created. Set the Focused to disable this behavior.
//
// Do not use this function to steal focus from other applications unless you
// are certain that is what the user wants. Focus stealing can be extremely
// disruptive.
//
// For a less disruptive way of getting the user's attention, see attention
// requests.
//
// This function may only be called from the main thread.
func (w *Window) Focus() {
	C.glfwFocusWindow(w.data)
}

// RequestAttention function requests user attention to the specified window.
// On platforms where this is not supported, attention is requested to the
// application as a whole.
//
// Once the user has given attention, usually by focusing the window or
// application, the system will end the request automatically.
//
// This function may only be called from the main thread.
func (w *Window) RequestAttention() {
	C.glfwRequestWindowAttention(w.data)
}

// GetMonitor returns the handle of the monitor that the window is in
// fullscreen on.
//
// Returns nil if the window is in windowed mode.
//
// This function may only be called from the main thread.
func (w *Window) GetMonitor() *Monitor {
	m := C.glfwGetWindowMonitor(w.data)
	if m == nil {
		return nil
	}
	return &Monitor{m}
}

// SetMonitor sets the monitor that the window uses for full screen mode or,
// if the monitor is NULL, makes it windowed mode.
//
// When setting a monitor, this function updates the width, height and refresh
// rate of the desired video mode and switches to the video mode closest to it.
// The window position is ignored when setting a monitor.
//
// When the monitor is NULL, the position, width and height are used to place
// the window client area. The refresh rate is ignored when no monitor is
// specified. If you only wish to update the resolution of a full screen window
// or the size of a windowed mode window, see SetSize.
//
// When a window transitions from full screen to windowed mode, this function
// restores any previous window settings such as whether it is decorated,
// floating, resizable, has size or aspect ratio limits, etc..
//
// This function may only be called from the main thread.
func (w *Window) SetMonitor(
	monitor *Monitor,
	xpos, ypos, width, height, refreshRate int,
) {
	var m *C.GLFWmonitor
	if monitor == nil {
		m = nil
	} else {
		m = monitor.data
	}
	C.glfwSetWindowMonitor(
		w.data,
		m,
		C.int(xpos),
		C.int(ypos),
		C.int(width),
		C.int(height),
		C.int(refreshRate),
	)
}

// GetAttrib returns the value of an attribute of the specified window or its
// OpenGL or OpenGL ES context.
//
// This function may only be called from the main thread.
func (w *Window) GetAttrib(attrib Hint) int {
	return int(C.glfwGetWindowAttrib(w.data, C.int(attrib)))
}

// SetAttrib function sets the value of an attribute of the specified window.
//
// The supported attributes are Decorated, Resizeable, Floating and AutoIconify.
//
// Some of these attributes are ignored for full screen windows. The new value
// will take effect if the window is later made windowed.
//
// Some of these attributes are ignored for windowed mode windows. The new value
// will take effect if the window is later made full screen.
//
// This function may only be called from the main thread.
func (w *Window) SetAttrib(attrib Hint, value int) {
	C.glfwSetWindowAttrib(w.data, C.int(attrib), C.int(value))
}

// SetUserPointer sets the user-defined pointer of the window. The current value
// is retained until the window is destroyed. The initial value is nil.
//
// This function may be called from any thread. Access is not synchronized.
func (w *Window) SetUserPointer(pointer unsafe.Pointer) {
	C.glfwSetWindowUserPointer(w.data, pointer)
}

// GetUserPointer returns the current value of the user-defined pointer of the
// window. The initial value is nil.
//
// This function may be called from any thread. Access is not synchronized.
func (w *Window) GetUserPointer() unsafe.Pointer {
	return C.glfwGetWindowUserPointer(w.data)
}

// SetPosCallback sets the position callback of the window, which is called
// when the window is moved. The callback is provided with the screen position
// of the upper-left corner of the client area of the window.
//
// This function must only be called from the main thread.
func (w *Window) SetPosCallback(cbfun PosCallback) PosCallback {
	previous := w.fPosHolder
	w.fPosHolder = cbfun
	if cbfun == nil {
		C.glfwSetWindowPosCallback(w.data, nil)
	} else {
		C.glfwSetWindowPosCallbackCB(w.data)
	}
	return previous
}

// SetSizeCallback sets the size callback of the window, which is called when
// the window is resized. The callback is provided with the size, in screen
// coordinates, of the client area of the window.
//
// This function must only be called from the main thread.
func (w *Window) SetSizeCallback(cbfun SizeCallback) SizeCallback {
	previous := w.fSizeHolder
	w.fSizeHolder = cbfun
	if cbfun == nil {
		C.glfwSetWindowSizeCallback(w.data, nil)
	} else {
		C.glfwSetWindowSizeCallbackCB(w.data)
	}
	return previous
}

// SetFramebufferSizeCallback sets the framebuffer resize callback of the
// specified window, which is called when the framebuffer of the specified
// window is resized.
//
// This function must only be called from the main thread.
func (w *Window) SetFramebufferSizeCallback(
	cbfun FramebufferSizeCallback,
) FramebufferSizeCallback {
	previous := w.fFramebufferSizeHolder
	w.fFramebufferSizeHolder = cbfun
	if cbfun == nil {
		C.glfwSetFramebufferSizeCallback(w.data, nil)
	} else {
		C.glfwSetFramebufferSizeCallbackCB(w.data)
	}
	return previous
}

// SetCloseCallback sets the close callback of the window, which is called when
// the user attempts to close the window, for example by clicking the close
// widget in the title bar.
//
// The close flag is set before this callback is called, but you can modify it
// at any time with SetShouldClose.
//
// The close callback is not triggered by Destroy.
//
// This function must only be called from the main thread.
func (w *Window) SetCloseCallback(cbfun CloseCallback) CloseCallback {
	previous := w.fCloseHolder
	w.fCloseHolder = cbfun
	if cbfun == nil {
		C.glfwSetWindowCloseCallback(w.data, nil)
	} else {
		C.glfwSetWindowCloseCallbackCB(w.data)
	}
	return previous
}

// SetMaximizeCallback sets the maximization callback of the specified window,
// which is called when the window is maximized or restored.
//
// This function must only be called from the main thread.
func (w *Window) SetMaximizeCallback(cbfun MaximizeCallback) MaximizeCallback {
	previous := w.fMaximizeHolder
	w.fMaximizeHolder = cbfun
	if cbfun == nil {
		C.glfwSetWindowMaximizeCallback(w.data, nil)
	} else {
		C.glfwSetWindowMaximizeCallbackCB(w.data)
	}
	return previous
}

// SetContentScaleCallback function sets the window content scale callback of
// the specified window, which is called when the content scale of the specified
// window changes.
//
// This function must only be called from the main thread.
func (w *Window) SetContentScaleCallback(
	cbfun ContentScaleCallback,
) ContentScaleCallback {
	previous := w.fContentScaleHolder
	w.fContentScaleHolder = cbfun
	if cbfun == nil {
		C.glfwSetWindowContentScaleCallback(w.data, nil)
	} else {
		C.glfwSetWindowContentScaleCallbackCB(w.data)
	}
	return previous
}

// SetRefreshCallback sets the refresh callback of the window, which
// is called when the client area of the window needs to be redrawn, for example
// if the window has been exposed after having been covered by another window.
//
// On compositing window systems such as Aero, Compiz or Aqua, where the window
// contents are saved off-screen, this callback may be called only very
// infrequently or never at all.
//
// This function must only be called from the main thread.
func (w *Window) SetRefreshCallback(cbfun RefreshCallback) RefreshCallback {
	previous := w.fRefreshHolder
	w.fRefreshHolder = cbfun
	if cbfun == nil {
		C.glfwSetWindowRefreshCallback(w.data, nil)
	} else {
		C.glfwSetWindowRefreshCallbackCB(w.data)
	}
	return previous
}

// SetFocusCallback sets the focus callback of the window, which is called when
// the window gains or loses focus.
//
// After the focus callback is called for a window that lost focus, synthetic
// key and mouse button release events will be generated for all such that had
// been pressed. For more information, see SetKeyCallback and
// SetMouseButtonCallback.
//
// This function must only be called from the main thread.
func (w *Window) SetFocusCallback(cbfun FocusCallback) FocusCallback {
	previous := w.fFocusHolder
	w.fFocusHolder = cbfun
	if cbfun == nil {
		C.glfwSetWindowFocusCallback(w.data, nil)
	} else {
		C.glfwSetWindowFocusCallbackCB(w.data)
	}
	return previous
}

// SetIconifyCallback sets the iconification callback of the window, which is
// called when the window is iconified or restored.
//
// This function must only be called from the main thread.
func (w *Window) SetIconifyCallback(cbfun IconifyCallback) IconifyCallback {
	previous := w.fIconifyHolder
	w.fIconifyHolder = cbfun
	if cbfun == nil {
		C.glfwSetWindowIconifyCallback(w.data, nil)
	} else {
		C.glfwSetWindowIconifyCallbackCB(w.data)
	}
	return previous
}

// SetClipboardString sets the system clipboard to the specified UTF-8 encoded
// string.
//
// This function must only be called from the main thread.
func (w *Window) SetClipboardString(str string) {
	cp := C.CString(str)
	defer C.free(unsafe.Pointer(cp))
	C.glfwSetClipboardString(w.data, cp)
}

// GetClipboardString returns the contents of the system clipboard, if it
// contains or is convertible to a UTF-8 encoded string.
//
// This function must only be called from the main thread.
func (w *Window) GetClipboardString() string {
	cs := C.glfwGetClipboardString(w.data)
	if cs == nil {
		return ""
	}
	return C.GoString(cs)
}

// PollEvents processes only those events that have already been received and
// then returns immediately. Processing events will cause the window and input
// callbacks associated with those events to be called.
//
// This function is not required for joystick input to work.
//
// This function may not be called from a callback.
//
// This function must only be called from the main thread.
func PollEvents() {
	C.glfwPollEvents()
}

// WaitEvents puts the calling thread to sleep until at least one event has been
// received. Once one or more events have been recevied, it behaves as if
// PollEvents was called, i.e. the events are processed and the function then
// returns immediately. Processing events will cause the window and input
// callbacks associated with those events to be called.
//
// Since not all events are associated with callbacks, this function may return
// without a callback having been called even if you are monitoring all
// callbacks.
//
// This function may not be called from a callback.
//
// This function must only be called from the main thread.
func WaitEvents() {
	C.glfwWaitEvents()
}

// WaitEventsTimeout puts the calling thread to sleep until at least one event
// is available in the event queue, or until the specified timeout is reached.
// If one or more events are available, it behaves exactly like PollEvents, i.e.
// the events in the queue are processed and the function then returns
// immediately. Processing events will cause the window and input callbacks
// associated with those events to be called.
//
// The timeout value must be a positive finite number.
//
// Since not all events are associated with callbacks, this function may return
// without a callback having been called even if you are monitoring all
// callbacks.
//
// On some platforms, a window move, resize or menu operation will cause event
// processing to block. This is due to how event processing is designed on those
// platforms. You can use the window refresh callback to redraw the contents of
// your window when necessary during such operations.
//
// On some platforms, certain callbacks may be called outside of a call to one
// of the event processing functions.
//
// If no windows exist, this function returns immediately. For synchronization
// of threads in applications that do not create windows, use native Go
// primitives.
//
// Event processing is not required for joystick input to work.
//
// This function must only be called from the main thread.
func WaitEventsTimeout(timeout float64) {
	C.glfwWaitEventsTimeout(C.double(timeout))
}

// PostEmptyEvent posts an empty event from the current thread to the main
// thread event queue, causing WaitEvents to return.
//
// If no windows exist, this function returns immediately. For synchronization
// of threads in applications that do not create windows, use native Go
// primitives.
//
// This function may be called from any thread.
func PostEmptyEvent() {
	C.glfwPostEmptyEvent()
}

// SwapBuffers swaps the front and back buffers of the specified window when
// rendering with OpenGL or OpenGL ES. If the swap interval is greater than
// zero, the GPU driver waits the specified number of screen updates before
// swapping the buffers.
//
// The specified window must have an OpenGL or OpenGL ES context. Specifying a
// window without a context will generate a NoWindowContext error.
//
// This function does not apply to Vulkan.
//
// This function may be called from any thread.
func (w *Window) SwapBuffers() {
	C.glfwSwapBuffers(w.data)
}

// Handle returns a *C.GLFWwindow reference (i.e. the GLFW window itself).
// This can be used for passing the GLFW window handle to external libraries
// like vulkan-go.
func (w *Window) Handle() uintptr {
	return uintptr(unsafe.Pointer(w.data))
}
