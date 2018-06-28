package pty

import (
	"os"
	"os/exec"
	"syscall"
)

func NewPtyWithShell() (*os.File, error) {
	pty, tty, err := open()
	if err != nil {
		return nil, err
	}
	defer tty.Close()
	shell := exec.Command("/bin/bash")
	shell.Stdout = tty
	shell.Stdin = tty
	shell.Stderr = tty
	shell.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}
	if err := shell.Start(); err != nil {
		pty.Close()
		return nil, err
	}
	return pty, nil
}
