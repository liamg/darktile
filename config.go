package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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

	home := os.Getenv("HOME")
	if home == "" {
		return &config.DefaultConfig
	}

	places := []string{}

	xdgHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgHome != "" {
		places = append(places, fmt.Sprintf("%s/aminal/config.toml", xdgHome))
	}

	places = append(places, fmt.Sprintf("%s/.config/aminal/config.toml", home))
	places = append(places, fmt.Sprintf("%s/.aminal.toml", home))

	for _, place := range places {
		if b, err := ioutil.ReadFile(place); err == nil {
			if c, err := config.Parse(b); err == nil {
				return c
			}

			fmt.Printf("Invalid config at %s: %s\n", place, err)
		}
	}

	parts := strings.Split(places[0], string(os.PathSeparator))
	path := strings.Join(parts[0:len(parts)-1], string(os.PathSeparator))

	err := os.MkdirAll(path, 0744)
	if err != nil {
		panic(err)
	}

	if b, err := config.DefaultConfig.Encode(); err != nil {
		fmt.Printf("Failed to encode config file: %s\n", err)
	} else {
		if err := ioutil.WriteFile(fmt.Sprintf("%s/config.toml", path), b, 0644); err != nil {
			fmt.Printf("Failed to encode config file: %s\n", err)
		}
	}
	return &config.DefaultConfig
}
