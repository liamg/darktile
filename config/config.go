package config

import (
	"gitlab.com/liamg/raft/terminal"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	DebugMode    bool `yaml:"debug"`
	ColourScheme terminal.ColourScheme
}

var DefaultConfig = Config{
	DebugMode:    false,
	ColourScheme: terminal.DefaultColourScheme,
}

func Parse(data []byte) (*Config, error) {
	c := DefaultConfig
	err := yaml.Unmarshal(data, &c)
	return &c, err
}
