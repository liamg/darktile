package pty

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/riywo/loginshell"
)

func NewPtyWithShell() (*os.File, error) {

	shellStr, err := loginshell.Shell()
	if err != nil {
		return nil, fmt.Errorf("Failed to ascertain your shell: %s", err)
	}

	pty, tty, err := open()
	if err != nil {
		return nil, err
	}
	defer tty.Close()
	shell := exec.Command(shellStr)
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
