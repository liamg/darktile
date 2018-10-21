package glfw

//#include <stdlib.h>
//#include "glfw/include/GLFW/glfw3.h"
//void glfwSetJoystickCallbackCB();
//void glfwSetKeyCallbackCB(GLFWwindow *window);
//void glfwSetCharCallbackCB(GLFWwindow *window);
//void glfwSetCharModsCallbackCB(GLFWwindow *window);
//void glfwSetMouseButtonCallbackCB(GLFWwindow *window);
//void glfwSetCursorPosCallbackCB(GLFWwindow *window);
//void glfwSetCursorEnterCallbackCB(GLFWwindow *window);
//void glfwSetScrollCallbackCB(GLFWwindow *window);
//void glfwSetDropCallbackCB(GLFWwindow *window);
//float GetAxisAtIndex(float *axis, int i);
//unsigned char GetButtonsAtIndex(unsigned char *buttons, int i);
//float GetGamepadAxisAtIndex(GLFWgamepadstate *gp, int i);
//unsigned char GetGamepadButtonAtIndex(GLFWgamepadstate *gp, int i);
import "C"

import (
	"image"
	"image/draw"
	"unsafe"
)

// Joystick corresponds to a joystick.
type Joystick int

// Joystick IDs.
const (
	Joystick1    Joystick = C.GLFW_JOYSTICK_1
	Joystick2    Joystick = C.GLFW_JOYSTICK_2
	Joystick3    Joystick = C.GLFW_JOYSTICK_3
	Joystick4    Joystick = C.GLFW_JOYSTICK_4
	Joystick5    Joystick = C.GLFW_JOYSTICK_5
	Joystick6    Joystick = C.GLFW_JOYSTICK_6
	Joystick7    Joystick = C.GLFW_JOYSTICK_7
	Joystick8    Joystick = C.GLFW_JOYSTICK_8
	Joystick9    Joystick = C.GLFW_JOYSTICK_9
	Joystick10   Joystick = C.GLFW_JOYSTICK_10
	Joystick11   Joystick = C.GLFW_JOYSTICK_11
	Joystick12   Joystick = C.GLFW_JOYSTICK_12
	Joystick13   Joystick = C.GLFW_JOYSTICK_13
	Joystick14   Joystick = C.GLFW_JOYSTICK_14
	Joystick15   Joystick = C.GLFW_JOYSTICK_15
	Joystick16   Joystick = C.GLFW_JOYSTICK_16
	JoystickLast Joystick = C.GLFW_JOYSTICK_LAST
)

// JoystickHatState corresponds to joystick hat states.
type JoystickHatState int

// Joystick Hat State IDs.
const (
	HatCentered  JoystickHatState = C.GLFW_HAT_CENTERED
	HatUp        JoystickHatState = C.GLFW_HAT_UP
	HatRight     JoystickHatState = C.GLFW_HAT_RIGHT
	HatDown      JoystickHatState = C.GLFW_HAT_DOWN
	HatLeft      JoystickHatState = C.GLFW_HAT_LEFT
	HatRightUp   JoystickHatState = C.GLFW_HAT_RIGHT_UP
	HatRightDown JoystickHatState = C.GLFW_HAT_RIGHT_DOWN
	HatLeftUp    JoystickHatState = C.GLFW_HAT_LEFT_UP
	HatLeftDown  JoystickHatState = C.GLFW_HAT_LEFT_DOWN
)

// GamepadAxis corresponds to a gamepad axis.
type GamepadAxis int

// Gamepad axis IDs.
const (
	AxisLeftX        GamepadAxis = C.GLFW_GAMEPAD_AXIS_LEFT_X
	AxisLeftY        GamepadAxis = C.GLFW_GAMEPAD_AXIS_LEFT_Y
	AxisRightX       GamepadAxis = C.GLFW_GAMEPAD_AXIS_RIGHT_X
	AxisRightY       GamepadAxis = C.GLFW_GAMEPAD_AXIS_RIGHT_Y
	AxisLeftTrigger  GamepadAxis = C.GLFW_GAMEPAD_AXIS_LEFT_TRIGGER
	AxisRightTrigger GamepadAxis = C.GLFW_GAMEPAD_AXIS_RIGHT_TRIGGER
	AxisLast         GamepadAxis = C.GLFW_GAMEPAD_AXIS_LAST
)

// GamepadButton corresponds to a gamepad button.
type GamepadButton int

// Gamepad button IDs.
const (
	ButtonA           GamepadButton = C.GLFW_GAMEPAD_BUTTON_A
	ButtonB           GamepadButton = C.GLFW_GAMEPAD_BUTTON_B
	ButtonX           GamepadButton = C.GLFW_GAMEPAD_BUTTON_X
	ButtonY           GamepadButton = C.GLFW_GAMEPAD_BUTTON_Y
	ButtonLeftBumper  GamepadButton = C.GLFW_GAMEPAD_BUTTON_LEFT_BUMPER
	ButtonRightBumper GamepadButton = C.GLFW_GAMEPAD_BUTTON_RIGHT_BUMPER
	ButtonBack        GamepadButton = C.GLFW_GAMEPAD_BUTTON_BACK
	ButtonStart       GamepadButton = C.GLFW_GAMEPAD_BUTTON_START
	ButtonGuide       GamepadButton = C.GLFW_GAMEPAD_BUTTON_GUIDE
	ButtonLeftThumb   GamepadButton = C.GLFW_GAMEPAD_BUTTON_LEFT_THUMB
	ButtonRightThumb  GamepadButton = C.GLFW_GAMEPAD_BUTTON_RIGHT_THUMB
	ButtonDpadUp      GamepadButton = C.GLFW_GAMEPAD_BUTTON_DPAD_UP
	ButtonDpadRight   GamepadButton = C.GLFW_GAMEPAD_BUTTON_DPAD_RIGHT
	ButtonDpadDown    GamepadButton = C.GLFW_GAMEPAD_BUTTON_DPAD_DOWN
	ButtonDpadLeft    GamepadButton = C.GLFW_GAMEPAD_BUTTON_DPAD_LEFT
	ButtonLast        GamepadButton = C.GLFW_GAMEPAD_BUTTON_LAST
	ButtonCross       GamepadButton = C.GLFW_GAMEPAD_BUTTON_CROSS
	ButtonCircle      GamepadButton = C.GLFW_GAMEPAD_BUTTON_CIRCLE
	ButtonSquare      GamepadButton = C.GLFW_GAMEPAD_BUTTON_SQUARE
	ButtonTriangle    GamepadButton = C.GLFW_GAMEPAD_BUTTON_TRIANGLE
)

// Key corresponds to a keyboard key.
type Key int

// These key codes are inspired by the USB HID Usage Tables v1.12 (p. 53-60),
// but re-arranged to map to 7-bit ASCII for printable keys (function keys are
// put in the 256+ range).
const (
	KeyUnknown      Key = C.GLFW_KEY_UNKNOWN
	KeySpace        Key = C.GLFW_KEY_SPACE
	KeyApostrophe   Key = C.GLFW_KEY_APOSTROPHE
	KeyComma        Key = C.GLFW_KEY_COMMA
	KeyMinus        Key = C.GLFW_KEY_MINUS
	KeyPeriod       Key = C.GLFW_KEY_PERIOD
	KeySlash        Key = C.GLFW_KEY_SLASH
	Key0            Key = C.GLFW_KEY_0
	Key1            Key = C.GLFW_KEY_1
	Key2            Key = C.GLFW_KEY_2
	Key3            Key = C.GLFW_KEY_3
	Key4            Key = C.GLFW_KEY_4
	Key5            Key = C.GLFW_KEY_5
	Key6            Key = C.GLFW_KEY_6
	Key7            Key = C.GLFW_KEY_7
	Key8            Key = C.GLFW_KEY_8
	Key9            Key = C.GLFW_KEY_9
	KeySemicolon    Key = C.GLFW_KEY_SEMICOLON
	KeyEqual        Key = C.GLFW_KEY_EQUAL
	KeyA            Key = C.GLFW_KEY_A
	KeyB            Key = C.GLFW_KEY_B
	KeyC            Key = C.GLFW_KEY_C
	KeyD            Key = C.GLFW_KEY_D
	KeyE            Key = C.GLFW_KEY_E
	KeyF            Key = C.GLFW_KEY_F
	KeyG            Key = C.GLFW_KEY_G
	KeyH            Key = C.GLFW_KEY_H
	KeyI            Key = C.GLFW_KEY_I
	KeyJ            Key = C.GLFW_KEY_J
	KeyK            Key = C.GLFW_KEY_K
	KeyL            Key = C.GLFW_KEY_L
	KeyM            Key = C.GLFW_KEY_M
	KeyN            Key = C.GLFW_KEY_N
	KeyO            Key = C.GLFW_KEY_O
	KeyP            Key = C.GLFW_KEY_P
	KeyQ            Key = C.GLFW_KEY_Q
	KeyR            Key = C.GLFW_KEY_R
	KeyS            Key = C.GLFW_KEY_S
	KeyT            Key = C.GLFW_KEY_T
	KeyU            Key = C.GLFW_KEY_U
	KeyV            Key = C.GLFW_KEY_V
	KeyW            Key = C.GLFW_KEY_W
	KeyX            Key = C.GLFW_KEY_X
	KeyY            Key = C.GLFW_KEY_Y
	KeyZ            Key = C.GLFW_KEY_Z
	KeyLeftBracket  Key = C.GLFW_KEY_LEFT_BRACKET
	KeyBackslash    Key = C.GLFW_KEY_BACKSLASH
	KeyRightBracket Key = C.GLFW_KEY_RIGHT_BRACKET
	KeyGraveAccent  Key = C.GLFW_KEY_GRAVE_ACCENT
	KeyWorld1       Key = C.GLFW_KEY_WORLD_1
	KeyWorld2       Key = C.GLFW_KEY_WORLD_2
	KeyEscape       Key = C.GLFW_KEY_ESCAPE
	KeyEnter        Key = C.GLFW_KEY_ENTER
	KeyTab          Key = C.GLFW_KEY_TAB
	KeyBackspace    Key = C.GLFW_KEY_BACKSPACE
	KeyInsert       Key = C.GLFW_KEY_INSERT
	KeyDelete       Key = C.GLFW_KEY_DELETE
	KeyRight        Key = C.GLFW_KEY_RIGHT
	KeyLeft         Key = C.GLFW_KEY_LEFT
	KeyDown         Key = C.GLFW_KEY_DOWN
	KeyUp           Key = C.GLFW_KEY_UP
	KeyPageUp       Key = C.GLFW_KEY_PAGE_UP
	KeyPageDown     Key = C.GLFW_KEY_PAGE_DOWN
	KeyHome         Key = C.GLFW_KEY_HOME
	KeyEnd          Key = C.GLFW_KEY_END
	KeyCapsLock     Key = C.GLFW_KEY_CAPS_LOCK
	KeyScrollLock   Key = C.GLFW_KEY_SCROLL_LOCK
	KeyNumLock      Key = C.GLFW_KEY_NUM_LOCK
	KeyPrintScreen  Key = C.GLFW_KEY_PRINT_SCREEN
	KeyPause        Key = C.GLFW_KEY_PAUSE
	KeyF1           Key = C.GLFW_KEY_F1
	KeyF2           Key = C.GLFW_KEY_F2
	KeyF3           Key = C.GLFW_KEY_F3
	KeyF4           Key = C.GLFW_KEY_F4
	KeyF5           Key = C.GLFW_KEY_F5
	KeyF6           Key = C.GLFW_KEY_F6
	KeyF7           Key = C.GLFW_KEY_F7
	KeyF8           Key = C.GLFW_KEY_F8
	KeyF9           Key = C.GLFW_KEY_F9
	KeyF10          Key = C.GLFW_KEY_F10
	KeyF11          Key = C.GLFW_KEY_F11
	KeyF12          Key = C.GLFW_KEY_F12
	KeyF13          Key = C.GLFW_KEY_F13
	KeyF14          Key = C.GLFW_KEY_F14
	KeyF15          Key = C.GLFW_KEY_F15
	KeyF16          Key = C.GLFW_KEY_F16
	KeyF17          Key = C.GLFW_KEY_F17
	KeyF18          Key = C.GLFW_KEY_F18
	KeyF19          Key = C.GLFW_KEY_F19
	KeyF20          Key = C.GLFW_KEY_F20
	KeyF21          Key = C.GLFW_KEY_F21
	KeyF22          Key = C.GLFW_KEY_F22
	KeyF23          Key = C.GLFW_KEY_F23
	KeyF24          Key = C.GLFW_KEY_F24
	KeyF25          Key = C.GLFW_KEY_F25
	KeyKP0          Key = C.GLFW_KEY_KP_0
	KeyKP1          Key = C.GLFW_KEY_KP_1
	KeyKP2          Key = C.GLFW_KEY_KP_2
	KeyKP3          Key = C.GLFW_KEY_KP_3
	KeyKP4          Key = C.GLFW_KEY_KP_4
	KeyKP5          Key = C.GLFW_KEY_KP_5
	KeyKP6          Key = C.GLFW_KEY_KP_6
	KeyKP7          Key = C.GLFW_KEY_KP_7
	KeyKP8          Key = C.GLFW_KEY_KP_8
	KeyKP9          Key = C.GLFW_KEY_KP_9
	KeyKPDecimal    Key = C.GLFW_KEY_KP_DECIMAL
	KeyKPDivide     Key = C.GLFW_KEY_KP_DIVIDE
	KeyKPMultiply   Key = C.GLFW_KEY_KP_MULTIPLY
	KeyKPSubtract   Key = C.GLFW_KEY_KP_SUBTRACT
	KeyKPAdd        Key = C.GLFW_KEY_KP_ADD
	KeyKPEnter      Key = C.GLFW_KEY_KP_ENTER
	KeyKPEqual      Key = C.GLFW_KEY_KP_EQUAL
	KeyLeftShift    Key = C.GLFW_KEY_LEFT_SHIFT
	KeyLeftControl  Key = C.GLFW_KEY_LEFT_CONTROL
	KeyLeftAlt      Key = C.GLFW_KEY_LEFT_ALT
	KeyLeftSuper    Key = C.GLFW_KEY_LEFT_SUPER
	KeyRightShift   Key = C.GLFW_KEY_RIGHT_SHIFT
	KeyRightControl Key = C.GLFW_KEY_RIGHT_CONTROL
	KeyRightAlt     Key = C.GLFW_KEY_RIGHT_ALT
	KeyRightSuper   Key = C.GLFW_KEY_RIGHT_SUPER
	KeyMenu         Key = C.GLFW_KEY_MENU
	KeyLast         Key = C.GLFW_KEY_LAST
)

// ModifierKey corresponds to a modifier key.
type ModifierKey int

// Modifier keys.
const (
	ModShift   ModifierKey = C.GLFW_MOD_SHIFT
	ModControl ModifierKey = C.GLFW_MOD_CONTROL
	ModAlt     ModifierKey = C.GLFW_MOD_ALT
	ModSuper   ModifierKey = C.GLFW_MOD_SUPER
)

// MouseButton corresponds to a mouse button.
type MouseButton int

// Mouse buttons.
const (
	MouseButton1      MouseButton = C.GLFW_MOUSE_BUTTON_1
	MouseButton2      MouseButton = C.GLFW_MOUSE_BUTTON_2
	MouseButton3      MouseButton = C.GLFW_MOUSE_BUTTON_3
	MouseButton4      MouseButton = C.GLFW_MOUSE_BUTTON_4
	MouseButton5      MouseButton = C.GLFW_MOUSE_BUTTON_5
	MouseButton6      MouseButton = C.GLFW_MOUSE_BUTTON_6
	MouseButton7      MouseButton = C.GLFW_MOUSE_BUTTON_7
	MouseButton8      MouseButton = C.GLFW_MOUSE_BUTTON_8
	MouseButtonLast   MouseButton = C.GLFW_MOUSE_BUTTON_LAST
	MouseButtonLeft   MouseButton = C.GLFW_MOUSE_BUTTON_LEFT
	MouseButtonRight  MouseButton = C.GLFW_MOUSE_BUTTON_RIGHT
	MouseButtonMiddle MouseButton = C.GLFW_MOUSE_BUTTON_MIDDLE
)

// StandardCursor corresponds to a standard cursor icon.
type StandardCursor int

// Standard cursors
const (
	ArrowCursor     StandardCursor = C.GLFW_ARROW_CURSOR
	IBeamCursor     StandardCursor = C.GLFW_IBEAM_CURSOR
	CrosshairCursor StandardCursor = C.GLFW_CROSSHAIR_CURSOR
	HandCursor      StandardCursor = C.GLFW_HAND_CURSOR
	HResizeCursor   StandardCursor = C.GLFW_HRESIZE_CURSOR
	VResizeCursor   StandardCursor = C.GLFW_VRESIZE_CURSOR
)

// Action corresponds to a key or button action.
type Action int

// Action types.
const (
	Release Action = C.GLFW_RELEASE // The key or button was released.
	Press   Action = C.GLFW_PRESS   // The key or button was pressed.
	Repeat  Action = C.GLFW_REPEAT  // The key was held down until it repeated.
)

// InputMode corresponds to an input mode.
type InputMode int

// Input modes.
const (
	CursorMode             InputMode = C.GLFW_CURSOR
	StickyKeysMode         InputMode = C.GLFW_STICKY_KEYS
	StickyMouseButtonsMode InputMode = C.GLFW_STICKY_MOUSE_BUTTONS
)

// Cursor mode values.
const (
	CursorNormal   int = C.GLFW_CURSOR_NORMAL
	CursorHidden   int = C.GLFW_CURSOR_HIDDEN
	CursorDisabled int = C.GLFW_CURSOR_DISABLED
)

// Cursor represents a cursor.
type Cursor struct {
	data *C.GLFWcursor
}

// GamepadState describes the input state of a gamepad.
type GamepadState struct {
	Buttons [15]Action
	Axes    [6]float32
}

// MouseButtonCallback is the function signature for mouse button callback
// functions.
type MouseButtonCallback func(
	w *Window,
	button MouseButton,
	action Action,
	mod ModifierKey,
)

// CursorPosCallback is the function signature for cursor position callback
// functions.
type CursorPosCallback func(w *Window, xpos float64, ypos float64)

// CursorEnterCallback is the function signature for cursor enter/leave callback
// functions.
type CursorEnterCallback func(w *Window, entered bool)

// ScrollCallback is the function signature for scroll callback functions.
type ScrollCallback func(w *Window, xoff float64, yoff float64)

// KeyCallback is the function signature for keyboard key callback functions.
type KeyCallback func(
	w *Window,
	key Key,
	scancode int,
	action Action,
	mods ModifierKey,
)

// CharCallback is the function signature for Unicode character callback
// functions.
type CharCallback func(w *Window, char rune)

// CharModsCallback is the function signature for Unicode character with
// modifiers callback functions. It is called for each input character,
// regardless of what modifier keys are held down.
type CharModsCallback func(w *Window, char rune, mods ModifierKey)

// DropCallback is the function signature for file drop callbacks.
type DropCallback func(w *Window, names []string)

// JoystickCallback is the function signature for joystick configuration
// callback functions.
type JoystickCallback func(joy Joystick, event PeripheralEvent)

var fJoystickHolder JoystickCallback

//export goJoystickCB
func goJoystickCB(joy, event C.int) {
	fJoystickHolder(Joystick(joy), PeripheralEvent(event))
}

//export goMouseButtonCB
func goMouseButtonCB(window unsafe.Pointer, button, action, mods C.int) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fMouseButtonHolder(
		w, MouseButton(button),
		Action(action),
		ModifierKey(mods),
	)
}

//export goCursorPosCB
func goCursorPosCB(window unsafe.Pointer, xpos, ypos C.double) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fCursorPosHolder(w, float64(xpos), float64(ypos))
}

//export goCursorEnterCB
func goCursorEnterCB(window unsafe.Pointer, entered C.int) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fCursorEnterHolder(w, glfwbool(entered))
}

//export goScrollCB
func goScrollCB(window unsafe.Pointer, xoff, yoff C.double) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fScrollHolder(w, float64(xoff), float64(yoff))
}

//export goKeyCB
func goKeyCB(window unsafe.Pointer, key, scancode, action, mods C.int) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fKeyHolder(w, Key(key), int(scancode), Action(action), ModifierKey(mods))
}

//export goCharCB
func goCharCB(window unsafe.Pointer, character C.uint) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fCharHolder(w, rune(character))
}

//export goCharModsCB
func goCharModsCB(window unsafe.Pointer, character C.uint, mods C.int) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fCharModsHolder(w, rune(character), ModifierKey(mods))
}

//export goDropCB
func goDropCB(window unsafe.Pointer, count C.int, names **C.char) {
	// TODO: Make this better. This part is unfinished, hacky, probably not
	// correct, and not idiomatic.
	w := windows.get((*C.GLFWwindow)(window))
	namesSlice := make([]string, int(count))
	for i := 0; i < int(count); i++ {
		var x *C.char
		p := (**C.char)(unsafe.Pointer(
			uintptr(unsafe.Pointer(names)) + uintptr(i)*unsafe.Sizeof(x)),
		)
		namesSlice[i] = C.GoString(*p)
	}
	w.fDropHolder(w, namesSlice)
}

// SetKeyCallback sets the key callback of the specified window, which is called
// when a key is pressed, repeated or released.
//
// The key functions deal with physical keys, with layout independent key tokens
// named after their values in the standard US keyboard layout. If you want to
// input text, use the character callback instead.
//
// When a window loses input focus, it will generate synthetic key release
// events for all pressed keys. You can tell these events from user-generated
// events by the fact that the synthetic ones are generated after the focus loss
// event has been processed, i.e. after the window focus callback has been
// called.
//
// The scancode of a key is specific to that platform or sometimes even to that
// machine. Scancodes are intended to allow users to bind keys that don't have a
// GLFW key token. Such keys have key set to KeyUnknown, their state is not
// saved and so it cannot be queried with GetKey.
//
// Sometimes GLFW needs to generate synthetic key events, in which case the
// scancode may be zero.
//
// This function must only be called from the main thread.
func (w *Window) SetKeyCallback(cbfun KeyCallback) KeyCallback {
	previous := w.fKeyHolder
	w.fKeyHolder = cbfun
	if cbfun == nil {
		C.glfwSetKeyCallback(w.data, nil)
	} else {
		C.glfwSetKeyCallbackCB(w.data)
	}
	return previous
}

// SetCharCallback sets the character callback of the specified window, which is
// called when a Unicode character is input.
//
// The character callback is intended for Unicode text input. As it deals with
// characters, it is keyboard layout dependent, whereas the key callback is not.
// Characters do not map 1:1 to physical keys, as a key may produce zero, one or
// more characters. If you want to know whether a specific physical key was
// pressed or released, see the key callback instead.
//
// The character callback behaves as system text input normally does and will
// not be called if modifier keys are held down that would prevent normal text
// input on that platform, for example a Super (Command) key on macOS or Alt key
// on Windows. There is a character with modifiers callback that receives
// these events.
//
// This function must only be called from the main thread.
func (w *Window) SetCharCallback(cbfun CharCallback) CharCallback {
	previous := w.fCharHolder
	w.fCharHolder = cbfun
	if cbfun == nil {
		C.glfwSetCharCallback(w.data, nil)
	} else {
		C.glfwSetCharCallbackCB(w.data)
	}
	return previous
}

// SetCharModsCallback sets the character with modifiers callback of the
// specified window, which is called when a Unicode character is input
// regardless of what modifier keys are used.
//
// The character with modifiers callback is intended for implementing custom
// Unicode character input. For regular Unicode text input, see the character
// callback. Like the character callback, the character with modifiers callback
// deals with characters and is keyboard layout dependent. Characters do not map
// 1:1 to physical keys, as a key may produce zero, one or more characters.
// If you want to know whether a specific physical key was pressed or released,
// see the key callback instead.
//
// This function must only be called from the main thread.
func (w *Window) SetCharModsCallback(cbfun CharModsCallback) CharModsCallback {
	previous := w.fCharModsHolder
	w.fCharModsHolder = cbfun
	if cbfun == nil {
		C.glfwSetCharModsCallback(w.data, nil)
	} else {
		C.glfwSetCharModsCallbackCB(w.data)
	}
	return previous
}

// SetMouseButtonCallback sets the mouse button callback of the specified
// window, which is called when a mouse button is pressed or released.
//
// When a window loses input focus, it will generate synthetic mouse button
// release events for all pressed mouse buttons. You can tell these events from
// user-generated events by the fact that the synthetic ones are generated after
// the focus loss event has been processed, i.e. after the window focus callback
// has been called.
//
// This function must only be called from the main thread.
func (w *Window) SetMouseButtonCallback(
	cbfun MouseButtonCallback,
) MouseButtonCallback {
	previous := w.fMouseButtonHolder
	w.fMouseButtonHolder = cbfun
	if cbfun == nil {
		C.glfwSetMouseButtonCallback(w.data, nil)
	} else {
		C.glfwSetMouseButtonCallbackCB(w.data)
	}
	return previous
}

// SetCursorPosCallback sets the cursor position callback which is called
// when the cursor is moved. The callback is provided with the position relative
// to the upper-left corner of the client area of the window.
//
// This function must only be called from the main thread.
func (w *Window) SetCursorPosCallback(
	cbfun CursorPosCallback,
) CursorPosCallback {
	previous := w.fCursorPosHolder
	w.fCursorPosHolder = cbfun
	if cbfun == nil {
		C.glfwSetCursorPosCallback(w.data, nil)
	} else {
		C.glfwSetCursorPosCallbackCB(w.data)
	}
	return previous
}

// SetCursorEnterCallback the cursor boundary crossing callback of the specified
// window, which is called when the cursor enters or leaves the client area of
// the window.
//
// This function must only be called from the main thread.
func (w *Window) SetCursorEnterCallback(
	cbfun CursorEnterCallback,
) CursorEnterCallback {
	previous := w.fCursorEnterHolder
	w.fCursorEnterHolder = cbfun
	if cbfun == nil {
		C.glfwSetCursorEnterCallback(w.data, nil)
	} else {
		C.glfwSetCursorEnterCallbackCB(w.data)
	}
	return previous
}

// SetScrollCallback sets the scroll callback of the specified window, which is
// called when a scrolling device is used, such as a mouse wheel or scrolling
// area of a touchpad.
//
// The scroll callback receives all scrolling input, like that from a mouse
// wheel or a touchpad scrolling area.
//
// This function must only be called from the main thread.
func (w *Window) SetScrollCallback(cbfun ScrollCallback) ScrollCallback {
	previous := w.fScrollHolder
	w.fScrollHolder = cbfun
	if cbfun == nil {
		C.glfwSetScrollCallback(w.data, nil)
	} else {
		C.glfwSetScrollCallbackCB(w.data)
	}
	return previous
}

// SetDropCallback sets the file drop callback of the specified window, which is
// called when one or more dragged files are dropped on the window.
//
// Because the path array and its strings may have been generated specifically
// for that event, they are not guaranteed to be valid after the callback has
// returned. If you wish to use them after the callback returns, you need to
// make a deep copy.
//
// This function must only be called from the main thread.
func (w *Window) SetDropCallback(cbfun DropCallback) DropCallback {
	previous := w.fDropHolder
	w.fDropHolder = cbfun
	if cbfun == nil {
		C.glfwSetDropCallback(w.data, nil)
	} else {
		C.glfwSetDropCallbackCB(w.data)
	}
	return previous
}

// SetJoystickCallback sets the joystick configuration callback, or removes the
// currently set callback. This is called when a joystick is connected to or
// disconnected from the system.
//
// For joystick connection and disconnection events to be delivered on all
// platforms, you need to call one of the event processing functions. Joystick
// disconnection may also be detected and the callback called by joystick
// functions. The function will then return whatever it returns if the joystick
// is not present.
//
// This function must only be called from the main thread.
func SetJoystickCallback(cbfun JoystickCallback) JoystickCallback {
	previous := fJoystickHolder
	fJoystickHolder = cbfun
	if cbfun == nil {
		C.glfwSetJoystickCallback(nil)
	} else {
		C.glfwSetJoystickCallbackCB()
	}
	return previous
}

// GetInputMode returns the value of an input option for the specified window.
//
// This function must only be called from the main thread.
func (w *Window) GetInputMode(mode InputMode) int {
	return int(C.glfwGetInputMode(w.data, C.int(mode)))
}

// SetInputMode sets an input mode option for the specified window.
//
// This function must only be called from the main thread.
func (w *Window) SetInputMode(mode InputMode, value int) {
	C.glfwSetInputMode(w.data, C.int(mode), C.int(value))
}

// GetKeyName returns the name of the specified printable key, encoded as UTF-8.
// This is typically the character that key would produce without any modifier
// keys, intended for displaying key bindings to the user. For dead keys, it is
// typically the diacritic it would add to a character.
//
// Do not use this function for text input. You will break text input for many
// languages even if it happens to work for yours.
//
// If the key is KeyUnknown, the scancode is used to identify the key, otherwise
// the scancode is ignored. If you specify a non-printable key, or KeyUnknown
// and a scancode that maps to a non-printable key, this function returns
// empty string but does not emit an error.
//
// This behavior allows you to always pass in the arguments in the key callback
// without modification.
//
// This function must only be called from the main thread.
func GetKeyName(key Key, scancode int) string {
	ret := C.glfwGetKeyName(C.int(key), C.int(scancode))
	return C.GoString(ret)
}

// GetKeyScancode function returns the platform-specific scancode of the
// specified key.
//
// If the key is KeyUnknown or does not exist on the keyboard this method will
// return -1.
//
// This function may be called from any thread.
func GetKeyScancode(key Key) int {
	return int(C.glfwGetKeyScancode(C.int(key)))
}

// GetKey returns the last state reported for the specified key to the specified
// window. The returned state is one of Press or Release. The higher-level
// action Repeat is only reported to the key callback.
//
// If the StickyKeys input mode is enabled, this function returns Press the
// first time you call it for a key that was pressed, even if that key has
// already been released.
//
// The key functions deal with physical keys, with key tokens named after their
// use on the standard US keyboard layout. If you want to input text, use the
// Unicode character callback instead.
//
// The modifier key bit masks are not key tokens and cannot be used with this
// function.
//
// Do not use this function to implement text input.
//
// This function must only be called from the main thread.
func (w *Window) GetKey(key Key) Action {
	return Action(C.glfwGetKey(w.data, C.int(key)))
}

// GetMouseButton returns the last state reported for the specified mouse button
// to the specified window. The returned state is one of Press or Release.
//
// If the StickyMouseButtons input mode is enabled, this function returns Press
// the first time you call it for a mouse button that was pressed, even if that
// mouse button has already been released.
//
// This function must only be called from the main thread.
func (w *Window) GetMouseButton(button MouseButton) Action {
	return Action(C.glfwGetMouseButton(w.data, C.int(button)))
}

// GetCursorPos returns the position of the cursor, in screen coordinates,
// relative to the upper-left corner of the client area of the specified window.
//
// If the cursor is disabled (with CursorDisabled) then the cursor position is
// unbounded and limited only by the minimum and maximum values of a double.
//
// The coordinate can be converted to their integer equivalents with the floor
// function. Casting directly to an integer type works for positive coordinates,
// but fails for negative ones.
//
// This function must only be called from the main thread.
func (w *Window) GetCursorPos() (x, y float64) {
	var xpos, ypos C.double
	C.glfwGetCursorPos(w.data, &xpos, &ypos)
	return float64(xpos), float64(ypos)
}

// SetCursorPos sets the position, in screen coordinates, of the cursor relative
// to the upper-left corner of the client area of the specified window.
// The window must have input focus. If the window does not have input focus
// when this function is called, it fails silently.
//
// Do not use this function to implement things like camera controls. GLFW
// already provides the CursorDisabled cursor mode that hides the cursor,
// transparently re-centers it and provides unconstrained cursor motion. See
// SetInputMode for more information.
//
// If the cursor mode is CursorDisabled then the cursor position is
// unconstrained and limited only by the minimum and maximum values of a double.
//
// This function must only be called from the main thread.
func (w *Window) SetCursorPos(xpos, ypos float64) {
	C.glfwSetCursorPos(w.data, C.double(xpos), C.double(ypos))
}

// CreateCursor creates a new custom cursor image that can be set for a window
// with SetCursor. The cursor can be destroyed with DestroyCursor. Any remaining
// cursors are destroyed by Terminate.
//
// The pixels are 32-bit, little-endian, non-premultiplied RGBA, i.e. eight bits
// per channel with the red channel first. They are arranged canonically as
// packed sequential rows, starting from the top-left corner.
//
// The cursor hotspot is specified in pixels, relative to the upper-left corner
// of the cursor image. Like all other coordinate systems in GLFW, the X-axis
// points to the right and the Y-axis points down.
//
// This function must only be called from the main thread.
func CreateCursor(img image.Image, xhot, yhot int) *Cursor {
	var imgC C.GLFWimage
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

	imgC.width = C.int(b.Dx())
	imgC.height = C.int(b.Dy())
	imgC.pixels = (*C.uchar)(pix)

	c := C.glfwCreateCursor(&imgC, C.int(xhot), C.int(yhot))

	free()

	return &Cursor{c}
}

// CreateStandardCursor returns a cursor with a standard shape, that can be set
// for a window with SetCursor.
//
// This function must only be called from the main thread.
func CreateStandardCursor(shape StandardCursor) *Cursor {
	c := C.glfwCreateStandardCursor(C.int(shape))
	return &Cursor{c}
}

// Destroy destroys a cursor previously created with CreateCursor. Any remaining
// cursors will be destroyed by Terminate.
//
// If the specified cursor is current for any window, that window will be
// reverted to the default cursor. This does not affect the cursor mode.
//
// This function must only be called from the main thread.
func (c *Cursor) Destroy() {
	C.glfwDestroyCursor(c.data)
}

// SetCursor sets the cursor image to be used when the cursor is over the client
// area of the specified window. The set cursor will only be visible when the
// cursor mode of the window is CursorNormal.
//
// On some platforms, the set cursor may not be visible unless the window also
// has input focus.
//
// This function must only be called from the main thread.
func (w *Window) SetCursor(c *Cursor) {
	if c == nil {
		C.glfwSetCursor(w.data, nil)
	} else {
		C.glfwSetCursor(w.data, c.data)
	}
}

// Present returns whether the specified joystick is present.
//
// There is no need to call this function before other methods of Joystick type
// as they all check for presence before performing any other work.
//
// This function must only be called from the main thread.
func (joy Joystick) Present() bool {
	return glfwbool(C.glfwJoystickPresent(C.int(joy)))
}

// GetAxes returns the values of all axes of the specified joystick. Each
// element in the array is a value between -1.0 and 1.0.
//
// If the specified joystick is not present this function will return nil but
// will not generate an error. This can be used instead of first calling
// Present.
//
// This function must only be called from the main thread.
func (joy Joystick) GetAxes() []float32 {
	var length int

	axis := C.glfwGetJoystickAxes(C.int(joy), (*C.int)(unsafe.Pointer(&length)))
	if axis == nil {
		return nil
	}

	a := make([]float32, length)
	for i := 0; i < length; i++ {
		a[i] = float32(C.GetAxisAtIndex(axis, C.int(i)))
	}
	return a
}

// GetButtons returns the state of all buttons of the specified joystick. Each
// element in the array is either Press or Release.
//
// For backward compatibility with earlier versions that did not have GetHats,
// the button array also includes all hats, each represented as four buttons.
// The hats are in the same order as returned by GetHats and are in the order
// up, right, down and left. To disable these extra buttons, set the
// JoystickHatButtons init hint before initialization.
//
// If the specified joystick is not present this function will return nil but
// will not generate an error. This can be used instead of first calling
// Present.
//
// This function must only be called from the main thread.
func (joy Joystick) GetButtons() []Action {
	var length int

	buttons := C.glfwGetJoystickButtons(
		C.int(joy),
		(*C.int)(unsafe.Pointer(&length)),
	)
	if buttons == nil {
		return nil
	}

	b := make([]Action, length)
	for i := 0; i < length; i++ {
		b[i] = Action(C.GetButtonsAtIndex(buttons, C.int(i)))
	}
	return b
}

// GetHats returns the state of all hats of the specified joystick.
//
// If the specified joystick is not present this function will return nil but
// will not generate an error. This can be used instead of first calling
// Present.
//
// This function must only be called from the main thread.
func (joy Joystick) GetHats() []JoystickHatState {
	var length int

	hats := C.glfwGetJoystickHats(C.int(joy), (*C.int)(unsafe.Pointer(&length)))
	if hats == nil {
		return nil
	}

	b := make([]JoystickHatState, length)
	for i := 0; i < length; i++ {
		b[i] = JoystickHatState(C.GetButtonsAtIndex(hats, C.int(i)))
	}
	return b
}

// GetName returns the name, encoded as UTF-8, of the specified joystick.
//
// If the specified joystick is not present this function will return nil but
// will not generate an error. This can be used instead of first calling
// Present.
//
// This function must only be called from the main thread.
func (joy Joystick) GetName() string {
	jn := C.glfwGetJoystickName(C.int(joy))
	return C.GoString(jn)
}

// GetGUID returns the SDL compatible GUID, as a UTF-8 encoded
// hexadecimal string, of the specified joystick.
//
// The GUID is what connects a joystick to a gamepad mapping. A connected
// joystick will always have a GUID even if there is no gamepad mapping
// assigned to it.
//
// If the specified joystick is not present this function will return empty
// string but will not generate an error. This can be used instead of first
// calling JoystickPresent.
//
// The GUID uses the format introduced in SDL 2.0.5. This GUID tries to uniquely
// identify the make and model of a joystick but does not identify a specific
// unit, e.g. all wired Xbox 360 controllers will have the same GUID on that
// platform. The GUID for a unit may vary between platforms depending on what
// hardware information the platform specific APIs provide.
//
// This function must only be called from the main thread.
func (joy Joystick) GetGUID() string {
	guid := C.glfwGetJoystickGUID(C.int(joy))
	return C.GoString(guid)
}

// SetUserPointer sets the user-defined pointer of the joystick. The current value
// is retained until the joystick is disconnected. The initial value is nil.
//
// This function may be called from the joystick callback, even for a joystick
// that is being disconnected.
//
// This function may be called from any thread. Access is not synchronized.
func (joy Joystick) SetUserPointer(pointer unsafe.Pointer) {
	C.glfwSetJoystickUserPointer(C.int(joy), pointer)
}

// GetUserPointer returns the current value of the user-defined pointer of the
// joystick. The initial value is nil.
//
// This function may be called from the joystick callback, even for a joystick
// that is being disconnected.
//
// This function may be called from any thread. Access is not synchronized.
func (joy Joystick) GetUserPointer() unsafe.Pointer {
	return C.glfwGetJoystickUserPointer(C.int(joy))
}

// IsGamepad returns whether the specified joystick is both present and
// has a gamepad mapping.
//
// If the specified joystick is present but does not have a gamepad mapping this
// function will return false but will not generate an error. Call Present to
// check if a joystick is present regardless of whether it has a mapping.
//
// This function must only be called from the main thread.
func (joy Joystick) IsGamepad() bool {
	return glfwbool(C.glfwJoystickIsGamepad(C.int(joy)))
}

// UpdateGamepadMappings parses the specified ASCII encoded string and updates
// the internal list with any gamepad mappings it finds. This string may contain
// either a single gamepad mapping or many mappings separated by newlines. The
// parser supports the full format of the gamecontrollerdb.txt source file
// including empty lines and comments.
//
// See Gamepad mappings for a description of the format.
//
// If there is already a gamepad mapping for a given GUID in the internal list,
// it will be replaced by the one passed to this function. If the library is
// terminated and re-initialized the internal list will revert to the built-in
// default.
//
// This function must only be called from the main thread.
func UpdateGamepadMappings(mapping string) bool {
	m := C.CString(mapping)
	defer C.free(unsafe.Pointer(m))
	return glfwbool(C.glfwUpdateGamepadMappings(m))
}

// GetGamepadName returns the human-readable name of the gamepad from the
// gamepad mapping assigned to the specified joystick.
//
// If the specified joystick is not present or does not have a gamepad mapping
// this function will return empty string but will not generate an error. Call
// Present to check whether it is present regardless of whether it has a
// mapping.
//
// This function must only be called from the main thread.
func (joy Joystick) GetGamepadName() string {
	gn := C.glfwGetGamepadName(C.int(joy))
	return C.GoString(gn)
}

// GetGamepadState retrives the state of the specified joystick remapped to an
// Xbox-like gamepad.
//
// If the specified joystick is not present or does not have a gamepad mapping
// this function will return nil but will not generate an error. Call
// Present to check whether it is present regardless of whether it has a
// mapping.
//
// The Guide button may not be available for input as it is often hooked by the
// system or the Steam client.
//
// Not all devices have all the buttons or axes provided by GamepadState.
// Unavailable buttons and axes will always report Release and 0.0 respectively.
//
// This function must only be called from the main thread.
func (joy Joystick) GetGamepadState() *GamepadState {
	var (
		gs  GamepadState
		cgs C.GLFWgamepadstate
	)

	ret := C.glfwGetGamepadState(C.int(joy), &cgs)
	if ret == C.GLFW_FALSE {
		return nil
	}

	for i := 0; i < 15; i++ {
		gs.Buttons[i] = Action(C.GetGamepadButtonAtIndex(&cgs, C.int(i)))
	}

	for i := 0; i < 6; i++ {
		gs.Axes[i] = float32(C.GetGamepadAxisAtIndex(&cgs, C.int(i)))
	}

	return &gs
}
