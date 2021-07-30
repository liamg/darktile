package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Opacity float64
	Font    Font
}

type Font struct {
	Family string
	Size   float64
	DPI    float64
}

type ErrorFileNotFound struct {
	Path string
}

func (e *ErrorFileNotFound) Error() string {
	return fmt.Sprintf("file was not found at '%s'", e.Path)
}

func getConfigPath() (string, error) {
	return getPath("config.yaml")
}

func getPath(filename string) (string, error) {
	baseDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("config directory missing: %w", err)
	}

	return path.Join(baseDir, "darktile", filename), nil
}

func LoadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to locate config path: %w", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, &ErrorFileNotFound{Path: configPath}
	}

	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file at '%s': %w", configPath, err)
	}

	config := defaultConfig
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("invalid config file at '%s': %w", configPath, err)
	}

	return &config, nil
}

func (c *Config) Save() (string, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return "", fmt.Errorf("failed to locate config path: %w", err)
	}

	if err := os.MkdirAll(path.Dir(configPath), 0700); err != nil {
		return "", err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return "", err
	}

	return configPath, ioutil.WriteFile(configPath, data, 0600)
}
