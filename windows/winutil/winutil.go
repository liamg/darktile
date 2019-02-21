package winutil

import (
	"golang.org/x/sys/windows"
	"unicode/utf16"
	"unsafe"
)

var (
	kernel                = windows.MustLoadDLL("kernel32.dll")
	getModuleFileNameProc = kernel.MustFindProc("GetModuleFileNameW")
)

func GetExecutablePath() (string, error) {
	var n uint32
	b := make([]uint16, windows.MAX_PATH)
	size := uint32(len(b))
	bPtr := uintptr(unsafe.Pointer(&b[0]))
	r0, _, e1 := getModuleFileNameProc.Call(0, bPtr, uintptr(size))
	n = uint32(r0)
	if n == 0 {
		return "", e1
	}
	return string(utf16.Decode(b[0:n])), nil
}