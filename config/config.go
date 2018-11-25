package config

import (
	"bytes"

	"github.com/BurntSushi/toml"
)

type Config struct {
	DebugMode    bool             `toml:"debug"`
	Slomo        bool             `toml:"slomo"`
	ColourScheme ColourScheme     `toml:"colours"`
	Shell        string           `toml:"shell"`
	KeyMapping   KeyMappingConfig `toml:"keys"`
}

type KeyMappingConfig map[string]string

func Parse(data []byte) (*Config, error) {
	c := DefaultConfig
	err := toml.Unmarshal(data, &c)
	if c.KeyMapping == nil {
		c.KeyMapping = KeyMappingConfig(map[string]string{})
	}
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
