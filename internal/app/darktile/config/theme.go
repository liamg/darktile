package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type Theme struct {
	Black               string
	Red                 string
	Green               string
	Yellow              string
	Blue                string
	Magenta             string
	Cyan                string
	White               string
	BrightBlack         string
	BrightRed           string
	BrightGreen         string
	BrightYellow        string
	BrightBlue          string
	BrightMagenta       string
	BrightCyan          string
	BrightWhite         string
	Background          string
	Foreground          string
	SelectionBackground string
	SelectionForeground string
	CursorForeground    string
	CursorBackground    string
}

func getThemePath() (string, error) {
	return getPath("theme.yaml")
}

func loadTheme(themePath string) (*Theme, error) {

	if themePath == "" {
		var err error
		themePath, err = getThemePath()
		if err != nil {
			return nil, fmt.Errorf("failed to locate theme path: %w", err)
		}
	}

	if _, err := os.Stat(themePath); os.IsNotExist(err) {
		return nil, &ErrorFileNotFound{Path: themePath}
	}

	themeData, err := ioutil.ReadFile(themePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read theme file at '%s': %w", themePath, err)
	}

	theme := defaultTheme
	if err := yaml.Unmarshal(themeData, &theme); err != nil {
		return nil, fmt.Errorf("invalid theme file at '%s': %w", themePath, err)
	}

	return &theme, nil
}

func (t *Theme) Save() (string, error) {
	themePath, err := getThemePath()
	if err != nil {
		return "", fmt.Errorf("failed to locate theme path: %w", err)
	}

	if err := os.MkdirAll(path.Dir(themePath), 0700); err != nil {
		return "", err
	}

	data, err := yaml.Marshal(t)
	if err != nil {
		return "", err
	}

	return themePath, ioutil.WriteFile(themePath, data, 0600)
}
