package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/liamg/aminal/config"
	"github.com/liamg/aminal/gui"
	"github.com/liamg/aminal/platform"
	"github.com/liamg/aminal/terminal"
	"github.com/riywo/loginshell"
)

type callback func(terminal *terminal.Terminal, g *gui.GUI)

func init() {
	runtime.LockOSThread()
}

func main() {
	initialize(nil, nil)
}

func initialize(unitTestfunc callback, configOverride *config.Config) {
	conf := maybeGetConfig(configOverride)

	logger, err := getLogger(conf)
	if err != nil {
		fmt.Printf("Failed to create logger: %s\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	if conf.CPUProfile != "" {
		logger.Infof("Starting CPU profiling...")
		stop := startCPUProf(conf.CPUProfile)
		defer stop()
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

func startCPUProf(filename string) func() {
	profileFile, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(profileFile)
	return pprof.StopCPUProfile
}
