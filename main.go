package main

import (
	"fmt"
	"github.com/liamg/aminal/gui"
	"github.com/liamg/aminal/platform"
	"github.com/liamg/aminal/terminal"
	"github.com/riywo/loginshell"
	"os"
	"runtime"
)

type callback func(terminal *terminal.Terminal, g *gui.GUI)

func init() {
	runtime.LockOSThread()
}

func main() {
	initialize(nil)
}

func initialize(unitTestfunc callback) {
	conf := getConfig()
	logger, err := getLogger(conf)
	if err != nil {
		fmt.Printf("Failed to create logger: %s\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	if unitTestfunc != nil {
		// Force the scrollbar not showing when running unit tests
		conf.ShowVerticalScrollbar = false
	}

	logger.Infof("Allocating pty...")

	pty, err := platform.NewPty(80, 25)
	if err != nil {
		logger.Fatalf("Failed to allocate pty: %s", err)
	}

	shellStr := conf.Shell
	if shellStr == "" {
		loginShell, err := loginshell.Shell()
		if err != nil {
			logger.Fatalf("Failed to ascertain your shell: %s", err)
		}
		shellStr = loginShell
	}

	os.Setenv("TERM", "xterm-256color") // controversial! easier than installing terminfo everywhere, but obviously going to be slightly different to xterm functionality, so we'll see...
	os.Setenv("COLORTERM", "truecolor")

	guestProcess, err := pty.CreateGuestProcess(shellStr)
	if err != nil {
		pty.Close()
		logger.Fatalf("Failed to start your shell: %s", err)
	}
	defer guestProcess.Close()

	logger.Infof("Creating terminal...")
	terminal := terminal.New(pty, logger, conf)

	g, err := gui.New(conf, terminal, logger)
	if err != nil {
		logger.Fatalf("Cannot start: %s", err)
	}
	defer g.Free()

	if unitTestfunc != nil {
		go unitTestfunc(terminal, g)
	} else {
		go func() {
			if err := guestProcess.Wait(); err != nil {
				logger.Fatalf("Failed to wait for guest process: %s", err)
			}
			g.Close()
		}()
	}

	if err := g.Render(); err != nil {
		logger.Fatalf("Render error: %s", err)
	}
}
