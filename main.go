package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"

	"github.com/kr/pty"
	"github.com/riywo/loginshell"
	"gitlab.com/liamg/raft/config"
	"gitlab.com/liamg/raft/gui"
	"gitlab.com/liamg/raft/terminal"
	"go.uber.org/zap"
)

func getConfig() *config.Config {
	ignore := false
	flag.BoolVar(&ignore, "ignore-config", ignore, "Ignore user config files and use defauls")
	if ignore {
		return &config.DefaultConfig
	}

	conf := loadConfigFile()

	flag.BoolVar(&conf.DebugMode, "debug", conf.DebugMode, "Enable debug logging")
	flag.BoolVar(&conf.Slomo, "slomo", conf.Slomo, "Render in slow motion (useful for debugging)")
	flag.BoolVar(&conf.Rendering.AlwaysRepaint, "always-repaint", conf.Rendering.AlwaysRepaint, "Always repaint the window, even when no changes have occurred")

	flag.Parse()
	return conf
}

func loadConfigFile() *config.Config {

	home := os.Getenv("HOME")
	if home == "" {
		return &config.DefaultConfig
	}

	places := []string{
		//fmt.Sprintf("%s/.config/raft.yml", home),
		fmt.Sprintf("%s/.raft.toml", home),
	}

	for _, place := range places {
		if b, err := ioutil.ReadFile(place); err == nil {
			if c, err := config.Parse(b); err == nil {
				return c
			} else {
				fmt.Printf("Invalid config at %s: %s\n", place, err)
			}
		}
	}

	if b, err := config.DefaultConfig.Encode(); err != nil {
		fmt.Printf("Failed to encode config file: %s\n", err)
	} else {
		if err := ioutil.WriteFile(fmt.Sprintf("%s/.raft.toml", home), b, 0644); err != nil {
			fmt.Printf("Failed to encode config file: %s\n", err)
		}
	}
	return &config.DefaultConfig
}

func getLogger(conf *config.Config) (*zap.SugaredLogger, error) {

	var logger *zap.Logger
	var err error
	if conf.DebugMode {
		logger, err = zap.NewDevelopment()
	} else {
		loggerConfig := zap.NewProductionConfig()
		loggerConfig.Encoding = "console"
		logger, err = loggerConfig.Build()
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to create logger: %s", err)
	}
	return logger.Sugar(), nil
}

func main() {

	// parse this
	conf := getConfig()

	os.Setenv("TERM", "xterm-256color")

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

	g := gui.New(conf, terminal, logger)
	if err := g.Render(); err != nil {
		logger.Fatalf("Render error: %s", err)
	}

}
