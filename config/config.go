package config

import (
	"bytes"

	"github.com/BurntSushi/toml"
)

type Config struct {
	DebugMode    bool            `toml:"debug"`
	Rendering    RenderingConfig `toml:"rendering"`
	Slomo        bool            `toml:"slomo"`
	ColourScheme ColourScheme    `toml:"colours"`
	Shell        string          `toml:"shell"`
}

type RenderingConfig struct {
	AlwaysRepaint bool `toml:"always_repaint"`
}

func Parse(data []byte) (*Config, error) {
	c := DefaultConfig
	err := toml.Unmarshal(data, &c)
	return &c, err
}

func (c *Config) Encode() ([]byte, error) {
	var buf bytes.Buffer
	e := toml.NewEncoder(&buf)
	err := e.Encode(c)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
