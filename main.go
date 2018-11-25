package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/kr/pty"
	"github.com/liamg/aminal/gui"
	"github.com/liamg/aminal/terminal"
	"github.com/riywo/loginshell"
)

func main() {

	conf := getConfig()
	logger, err := getLogger(conf)
	if err != nil {
		fmt.Printf("Failed to create logger: %s\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Infof("Allocating pty...")
	pty, tty, err := pty.Open()
	if err != nil {
		logger.Fatalf("Failed to allocate pty: %s", err)
	}

	shellStr, err := loginshell.Shell()
	if err != nil {
		logger.Fatalf("Failed to ascertain your shell: %s", err)
	}

	if conf.Shell != "" {
		shellStr = conf.Shell
	}

	os.Setenv("TERM", "xterm-256color") // contraversial! easier than installing terminfo everywhere, but obviously going to be slightly different to xterm functionality, so we'll see...
	os.Setenv("COLORTERM", "truecolor")

	shell := exec.Command(shellStr)
	shell.Stdout = tty
	shell.Stdin = tty
	shell.Stderr = tty
	shell.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}
	if err := shell.Start(); err != nil {
		pty.Close()
		logger.Fatalf("Failed to start your shell: %s", err)
	}

	logger.Infof("Creating terminal...")
	terminal := terminal.New(pty, logger, conf)

	g, err := gui.New(conf, terminal, logger)
	if err != nil {
		logger.Fatalf("Cannot start: %s", err)
	}
	if err := g.Render(); err != nil {
		logger.Fatalf("Render error: %s", err)
	}

}
