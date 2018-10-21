package glfw

//#include "glfw/include/GLFW/glfw3.h"
import "C"

// GetTime returns the value of the GLFW timer. Unless the timer has been set
// using SetTime, the timer measures time elapsed since GLFW was initialized.
//
// The resolution of the timer is system dependent, but is usually on the order
// of a few micro- or nanoseconds. It uses the highest-resolution monotonic time
// source on each supported platform.
//
// This function may be called from any thread. Reading and writing of the
// internal timer offset is not atomic, so it needs to be externally
// synchronized with calls to SetTime.
func GetTime() float64 {
	return float64(C.glfwGetTime())
}

// SetTime sets the value of the GLFW timer. It then continues to count up from
// that value. The value must be a positive finite number less than or equal to
// 18446744073.0, which is approximately 584.5 years.
//
// This function may be called from any thread. Reading and writing of the
// internal timer offset is not atomic, so it needs to be externally
// synchronized with calls to GetTime.
func SetTime(time float64) {
	C.glfwSetTime(C.double(time))
}

// GetTimerValue returns the current value of the raw timer, measured in
// 1 / frequency seconds.
//
// This function may be called from any thread.
func GetTimerValue() uint64 {
	return uint64(C.glfwGetTimerValue())
}

// GetTimerFrequency returns frequency of the timer, in Hz, or zero if an
// error occurred.
//
// This function may be called from any thread.
func GetTimerFrequency() uint64 {
	return uint64(C.glfwGetTimerFrequency())
}
