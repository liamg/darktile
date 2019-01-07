package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"github.com/liamg/aminal/config"
	"github.com/liamg/aminal/version"
)

func getConfig() *config.Config {

	showVersion := false
	flag.BoolVar(&showVersion, "version", showVersion, "Output version information")

	ignore := false
	flag.BoolVar(&ignore, "ignore-config", ignore, "Ignore user config files and use defaults")
	if ignore {
		return &config.DefaultConfig
	}

	conf := loadConfigFile()

	flag.StringVar(&conf.Shell, "shell", conf.Shell, "Specify the shell to use")
	flag.BoolVar(&conf.DebugMode, "debug", conf.DebugMode, "Enable debug logging")
	flag.BoolVar(&conf.Slomo, "slomo", conf.Slomo, "Render in slow motion (useful for debugging)")

	flag.Parse()

	if showVersion {
		v := version.Version
		if v == "" {
			v = "development"
		}
		fmt.Println(v)
		os.Exit(0)
	}

	return conf
}

func loadConfigFile() *config.Config {

	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Failed to get current user information: %s\n", err)
		return &config.DefaultConfig
	}

	home := usr.HomeDir
	if home == "" {
		return &config.DefaultConfig
	}

	places := []string{}

	xdgHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgHome != "" {
		places = append(places, filepath.Join(xdgHome, "aminal/config.toml"))
	}

	places = append(places, filepath.Join(home, ".config/aminal/config.toml"))
	places = append(places, filepath.Join(home, ".aminal.toml"))

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
		err = os.MkdirAll(filepath.Dir(places[0]), 0744)
		if err != nil {
			fmt.Printf("Failed to create config file directory: %s\n", err)
		} else {
			if err = ioutil.WriteFile(places[0], b, 0644); err != nil {
				fmt.Printf("Failed to encode config file: %s\n", err)
			}
		}
	}

	return &config.DefaultConfig
}
