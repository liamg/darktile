package pty

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// todo: port these for darwin/windows:
func open() (*os.File, *os.File, error) {

	pty, err := getpt()
	if err != nil {
		panic(err)
	}

	ptsName, err := ptsname(pty)
	if err != nil {
		panic(err)
	}

	err = grantpt(pty)
	if err != nil {
		return nil, nil, err
	}

	err = unlockpt(pty)
	if err != nil {
		return nil, nil, err
	}

	tty, err := os.OpenFile(ptsName, os.O_RDWR, 0)
	if err != nil {
		return nil, nil, err
	}
	return pty, tty, nil

}

func getpt() (file *os.File, err error) {
	return os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
}

func ptsname(file *os.File) (name string, err error) {
	n, err := ioctl(file, syscall.TIOCPTYGNAME, 0)
	return fmt.Sprintf("/dev/pts/%d", n), err
}

func ioctl(file *os.File, command uint, arg int) (int, error) {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(file.Fd()),
		uintptr(command), uintptr(unsafe.Pointer(&arg)))
	if err != 0 {
		return 0, fmt.Errorf("Error no %d", err)
	}
	return arg, nil
}

func grantpt(f *os.File) error {
	_, err := ioctl(f, syscall.TIOCPTYGRANT, 0)
	return err
}

func unlockpt(f *os.File) error {
	_, err := ioctl(f, syscall.TIOCPTYUNLK, 0)
	return err
}
