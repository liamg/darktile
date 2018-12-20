//+build !windows

package platform

import (
	"errors"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/kr/pty"
)

type unixPty struct {
	pty                       *os.File
	tty                       *os.File
	platformDependentSettings PlatformDependentSettings
}

type winsize struct {
	Height uint16
	Width  uint16
	x      uint16 //ignored, but necessary for ioctl calls
	y      uint16 //ignored, but necessary for ioctl calls
}

func (p *unixPty) Read(b []byte) (int, error) {
	if p == nil || p.pty == nil {
		return 0, errors.New("Attempted to read from a deallocated pty")
	}
	return p.pty.Read(b)
}

func (p *unixPty) Write(b []byte) (int, error) {
	if p == nil || p.pty == nil {
		return 0, errors.New("Attempted to write to a deallocated pty")
	}
	return p.pty.Write(b)
}

func (p *unixPty) Close() error {
	if p == nil || p.pty == nil {
		return nil
	}
	ret := p.pty.Close()
	p.pty = nil
	p.tty = nil
	return ret
}

func (p *unixPty) Resize(x, y int) error {
	size := winsize{
		Height: uint16(y),
		Width:  uint16(x),
		x:      0,
		y:      0,
	}
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(p.pty.Fd()),
		uintptr(syscall.TIOCSWINSZ), uintptr(unsafe.Pointer(&size)))

	if errno != 0 {
		return errors.New(errno.Error())
	}

	return nil
}

func (p *unixPty) CreateGuestProcess(imagePath string) (Process, error) {
	if p == nil || p.tty == nil {
		return nil, errors.New("Attempted to create a process on a deallocated pty")
	}
	shell := newCmdProc(exec.Command(imagePath))
	shell.cmd.Stdout = p.tty
	shell.cmd.Stdin = p.tty
	shell.cmd.Stderr = p.tty
	shell.cmd.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}
	if err := shell.cmd.Start(); err != nil {
		return nil, err
	}

	return shell, nil
}

func (pty *unixPty) GetPlatformDependentSettings() PlatformDependentSettings {
	return pty.platformDependentSettings
}

func NewPty(x, y int) (Pty, error) {
	innerPty, innerTty, err := pty.Open()
	if err != nil {
		return nil, err
	}
	return &unixPty{
		pty: innerPty,
		tty: innerTty,
		platformDependentSettings: PlatformDependentSettings{
			OSCTerminators: map[rune]struct{}{0x07: {}, 0x5c: {}},
		},
	}, nil
}
