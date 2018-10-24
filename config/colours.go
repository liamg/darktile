package config

import (
	"encoding/hex"
	"fmt"
	"math"
	"strings"
)

type Colour [3]float32

func strToColourNoErr(hexStr string) Colour {
	c, _ := strToColour(hexStr)
	return c
}

func strToColour(hexStr string) (Colour, error) {

	c := [3]float32{0, 0, 0}

	if strings.HasPrefix(hexStr, "#") {
		hexStr = hexStr[1:]
	}

	if len(hexStr) != 6 {
		return c, fmt.Errorf("Invalid colour format. Should be like #ffffff")
	}

	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return c, err
	}

	c[0] = float32(bytes[0]) / 255
	c[1] = float32(bytes[1]) / 255
	c[2] = float32(bytes[2]) / 255

	return c, nil
}

func (c *Colour) UnmarshalText(data []byte) error {
	var err error
	*c, err = strToColour(string(data))
	return err
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
	Cursor       Colour `toml:"cursor"`
	Foreground   Colour `toml:"foreground"`
	Background   Colour `toml:"background"`
	Black        Colour `toml:"black"`
	Red          Colour `toml:"red"`
	Green        Colour `toml:"green"`
	Yellow       Colour `toml:"yellow"`
	Blue         Colour `toml:"blue"`
	Magenta      Colour `toml:"magenta"`
	Cyan         Colour `toml:"cyan"`
	LightGrey    Colour `toml:"light_grey"`
	DarkGrey     Colour `toml:"dark_grey"`
	LightRed     Colour `toml:"light_red"`
	LightGreen   Colour `toml:"light_green"`
	LightYellow  Colour `toml:"light_yellow"`
	LightBlue    Colour `toml:"light_blue"`
	LightMagenta Colour `toml:"light_magenta"`
	LightCyan    Colour `toml:"light_cyan"`
	White        Colour `toml:"white"`
	Selection    Colour `toml:"selection"`
}
