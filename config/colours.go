package config

import (
	"encoding/hex"
	"fmt"
	"math"
	"strings"
)

type Colour [3]float32

func (c *Colour) UnmarshalText(data []byte) error {

	hexStr := string(data)

	if strings.HasPrefix(hexStr, "#") {
		hexStr = hexStr[1:]
	}

	if len(hexStr) != 6 {
		return fmt.Errorf("Invalid colour format. Should be like #ffffff")
	}

	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return err
	}

	c[0] = float32(bytes[0]) / 255
	c[1] = float32(bytes[1]) / 255
	c[2] = float32(bytes[2]) / 255

	return nil
}

func (c Colour) MarshalText() (text []byte, err error) {
	return []byte(fmt.Sprintf(
		"#%02x%02x%02x",
		uint8(math.Floor(float64(255*c[0]))),
		uint8(math.Floor(float64(255*c[1]))),
		uint8(math.Floor(float64(255*c[2]))),
	)), nil
}

type ColourScheme struct {
	Cursor         Colour `toml:"cursor"`
	DefaultFg      Colour `toml:"default_fg"`
	BlackFg        Colour `toml:"black_fg"`
	RedFg          Colour `toml:"red_fg"`
	GreenFg        Colour `toml:"green_fg"`
	YellowFg       Colour `toml:"yellow_fg"`
	BlueFg         Colour `toml:"blue_fg"`
	MagentaFg      Colour `toml:"magenta_fg"`
	CyanFg         Colour `toml:"cyan_fg"`
	LightGreyFg    Colour `toml:"light_grey_fg"`
	DarkGreyFg     Colour `toml:"dark_grey_fg"`
	LightRedFg     Colour `toml:"light_red_fg"`
	LightGreenFg   Colour `toml:"light_green_fg"`
	LightYellowFg  Colour `toml:"light_yellow_fg"`
	LightBlueFg    Colour `toml:"light_blue_fg"`
	LightMagentaFg Colour `toml:"light_magenta_fg"`
	LightCyanFg    Colour `toml:"light_cyan_fg"`
	WhiteFg        Colour `toml:"white_fg"`
	DefaultBg      Colour `toml:"default_bg"`
	BlackBg        Colour `toml:"black_bg"`
	RedBg          Colour `toml:"red_bg"`
	GreenBg        Colour `toml:"green_bg"`
	YellowBg       Colour `toml:"yellow_bg"`
	BlueBg         Colour `toml:"blue_bg"`
	MagentaBg      Colour `toml:"magenta_bg"`
	CyanBg         Colour `toml:"cyan_bg"`
	LightGreyBg    Colour `toml:"light_grey_bg"`
	DarkGreyBg     Colour `toml:"dark_grey_bg"`
	LightRedBg     Colour `toml:"light_red_bg"`
	LightGreenBg   Colour `toml:"light_green_bg"`
	LightYellowBg  Colour `toml:"light_yellow_bg"`
	LightBlueBg    Colour `toml:"light_blue_bg"`
	LightMagentaBg Colour `toml:"light_magenta_bg"`
	LightCyanBg    Colour `toml:"light_cyan_bg"`
	WhiteBg        Colour `toml:"white_bg"`
}
