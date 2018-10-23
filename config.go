package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/liamg/aminal/config"
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
