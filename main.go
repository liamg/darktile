package main

import (
	"fmt"
	"os"

	"gitlab.com/liamg/terminal/config"
	"gitlab.com/liamg/terminal/gui"
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

	/*
		sugaredLogger.Infof("Allocationg pty...")
			pty, err := pty.NewPtyWithShell()
			if err != nil {
				panic(err)
			}
	*/

	g := gui.New(conf, sugaredLogger)
	if err := g.Render(); err != nil {
		sugaredLogger.Fatalf("Render error: %s", err)
	}

	//go io.Copy(pty, os.Stdin)
	//io.Copy(os.Stdout, pty)

	//	return pty, err
}
