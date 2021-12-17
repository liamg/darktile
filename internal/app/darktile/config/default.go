package config

import (
	"encoding/hex"
	"fmt"
	"image/color"

	"github.com/liamg/darktile/internal/app/darktile/termutil"
)

var defaultConfig = Config{
	Opacity: 1.0,
	Font: Font{
		Family:    "", // internally packed font will be loaded by default
		Size:      18.0,
		DPI:       72.0,
		Ligatures: true,
	},
}

var defaultTheme = Theme{
	Black:               "#1d1f21",
	Red:                 "#cc6666",
	Green:               "#b5bd68",
	Yellow:              "#f0c674",
	Blue:                "#81a2be",
	Magenta:             "#b294bb",
	Cyan:                "#8abeb7",
	White:               "#c5c8c6",
	BrightBlack:         "#666666",
	BrightRed:           "#d54e53",
	BrightGreen:         "#b9ca4a",
	BrightYellow:        "#e7c547",
	BrightBlue:          "#7aa6da",
	BrightMagenta:       "#c397d8",
	BrightCyan:          "#70c0b1",
	BrightWhite:         "#eaeaea",
	Background:          "#000000",
	Foreground:          "#c5c8c6",
	SelectionBackground: "#33aa33",
	SelectionForeground: "#ffffff",
	CursorForeground:    "#1d1f21",
	CursorBackground:    "#c5c8c6",
}

func DefaultConfig() *Config {
	copiedConf := defaultConfig
	return &copiedConf
}

func DefaultTheme(conf *Config) (*termutil.Theme, error) {
	return loadThemeFromConf(conf, &defaultTheme)
}

func LoadTheme(conf *Config) (*termutil.Theme, error) {

	themeConf, err := loadTheme("")
	if err != nil {
		return nil, err
	}

	return loadThemeFromConf(conf, themeConf)
}

func LoadThemeFromPath(conf *Config, path string) (*termutil.Theme, error) {

	themeConf, err := loadTheme(path)
	if err != nil {
		return nil, err
	}

	return loadThemeFromConf(conf, themeConf)
}

func loadThemeFromConf(conf *Config, themeConf *Theme) (*termutil.Theme, error) {

	factory := termutil.NewThemeFactory()

	colours := map[termutil.Colour]string{
		termutil.ColourBlack:               themeConf.Black,
		termutil.ColourRed:                 themeConf.Red,
		termutil.ColourGreen:               themeConf.Green,
		termutil.ColourYellow:              themeConf.Yellow,
		termutil.ColourBlue:                themeConf.Blue,
		termutil.ColourMagenta:             themeConf.Magenta,
		termutil.ColourCyan:                themeConf.Cyan,
		termutil.ColourWhite:               themeConf.White,
		termutil.ColourBrightBlack:         themeConf.BrightBlack,
		termutil.ColourBrightRed:           themeConf.BrightRed,
		termutil.ColourBrightGreen:         themeConf.BrightGreen,
		termutil.ColourBrightYellow:        themeConf.BrightYellow,
		termutil.ColourBrightBlue:          themeConf.BrightBlue,
		termutil.ColourBrightMagenta:       themeConf.BrightMagenta,
		termutil.ColourBrightCyan:          themeConf.BrightCyan,
		termutil.ColourBrightWhite:         themeConf.BrightWhite,
		termutil.ColourBackground:          themeConf.Background,
		termutil.ColourForeground:          themeConf.Foreground,
		termutil.ColourSelectionBackground: themeConf.SelectionBackground,
		termutil.ColourSelectionForeground: themeConf.SelectionForeground,
		termutil.ColourCursorForeground:    themeConf.CursorForeground,
		termutil.ColourCursorBackground:    themeConf.CursorBackground,
	}

	for key, colHex := range colours {
		col, err := colourFromHex(colHex, conf.Opacity)
		if err != nil {
			return nil, fmt.Errorf("invalid hex value '%s' in theme", colHex)
		}
		factory.WithColour(
			key,
			col,
		)
	}

	return factory.Build(), nil

}

func colourFromHex(hexadecimal string, opacity float64) (color.Color, error) {
	if len(hexadecimal) == 0 {
		return nil, fmt.Errorf("colour value cannot be empty")
	}
	if hexadecimal[0] != '#' || len(hexadecimal) != 7 {
		return nil, fmt.Errorf("colour values should start with '#' and contain an RGB value encoded in hex, for example #ffffff")
	}

	decoded, err := hex.DecodeString(hexadecimal[1:])
	if err != nil {
		return nil, err
	}

	return color.RGBA{
		R: decoded[0],
		G: decoded[1],
		B: decoded[2],
		A: uint8(opacity * 0xff),
	}, nil
}
