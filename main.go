package main

import (
	"fmt"
	"os"
	"time"

	"gitlab.com/liamg/raft/config"
	"gitlab.com/liamg/raft/gui"
	"gitlab.com/liamg/raft/pty"
	"gitlab.com/liamg/raft/terminal"
	"go.uber.org/zap"
)

func main() {

	// parse this
	conf := config.Config{DebugMode: true}

	var logger *zap.Logger
	var err error
	if conf.DebugMode {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		fmt.Printf("Failed to create logger: %s", err)
		os.Exit(1)
	}
	sugaredLogger := logger.Sugar()
	defer sugaredLogger.Sync()

	sugaredLogger.Infof("Allocationg pty...")
	pty, err := pty.NewPtyWithShell()
	if err != nil {
		panic(err)
	}

	sugaredLogger.Infof("Creating terminal...")
	terminal := terminal.New(pty, sugaredLogger)

	go func() {
		time.Sleep(time.Second * 5)
		terminal.Write([]byte("tput cols && tput lines\n"))
		terminal.Write([]byte("ls -la\n"))
	}()

	g := gui.New(conf, terminal, sugaredLogger)
	if err := g.Render(); err != nil {
		sugaredLogger.Fatalf("Render error: %s", err)
	}

	//go io.Copy(pty, os.Stdin)
	//io.Copy(os.Stdout, pty)

	//	return pty, err
}
