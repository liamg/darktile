package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/liamg/aminal/config"
)

var Version string

func getConfig() *config.Config {

	showVersion := false
	flag.BoolVar(&showVersion, "version", showVersion, "Output version information")

	ignore := false
	flag.BoolVar(&ignore, "ignore-config", ignore, "Ignore user config files and use defauls")
	if ignore {
		return &config.DefaultConfig
	}

	conf := loadConfigFile()

	flag.StringVar(&conf.Shell, "shell", conf.Shell, "Specify the shell to use")
	flag.BoolVar(&conf.DebugMode, "debug", conf.DebugMode, "Enable debug logging")
	flag.BoolVar(&conf.Slomo, "slomo", conf.Slomo, "Render in slow motion (useful for debugging)")

	flag.Parse()

	if showVersion {
		if Version == "" {
			Version = "development"
		}
		fmt.Printf("Aminal %s\n", Version)
		os.Exit(0)
	}

	return conf
}

func loadConfigFile() *config.Config {

	home := os.Getenv("HOME")
	if home == "" {
		return &config.DefaultConfig
	}

	places := []string{
		fmt.Sprintf("%s/.aminal.toml", home),
	}

	for _, place := range places {
		if b, err := ioutil.ReadFile(place); err == nil {
			if c, err := config.Parse(b); err == nil {
				return c
			}

			fmt.Printf("Invalid config at %s: %s\n", place, err)
		}
	}

	if b, err := config.DefaultConfig.Encode(); err != nil {
		fmt.Printf("Failed to encode config file: %s\n", err)
	} else {
		if err := ioutil.WriteFile(fmt.Sprintf("%s/.aminal.toml", home), b, 0644); err != nil {
			fmt.Printf("Failed to encode config file: %s\n", err)
		}
	}
	return &config.DefaultConfig
}
