package terminal

type ColourScheme struct {
	DefaultFg      [3]float32
	BlackFg        [3]float32
	RedFg          [3]float32
	GreenFg        [3]float32
	YellowFg       [3]float32
	BlueFg         [3]float32
	MagentaFg      [3]float32
	CyanFg         [3]float32
	LightGreyFg    [3]float32
	DarkGreyFg     [3]float32
	LightRedFg     [3]float32
	LightGreenFg   [3]float32
	LightYellowFg  [3]float32
	LightBlueFg    [3]float32
	LightMagentaFg [3]float32
	LightCyanFg    [3]float32
	WhiteFg        [3]float32
	DefaultBg      [3]float32
	BlackBg        [3]float32
	RedBg          [3]float32
	GreenBg        [3]float32
	YellowBg       [3]float32
	BlueBg         [3]float32
	MagentaBg      [3]float32
	CyanBg         [3]float32
	LightGreyBg    [3]float32
	DarkGreyBg     [3]float32
	LightRedBg     [3]float32
	LightGreenBg   [3]float32
	LightYellowBg  [3]float32
	LightBlueBg    [3]float32
	LightMagentaBg [3]float32
	LightCyanBg    [3]float32
	WhiteBg        [3]float32
}

var DefaultColourScheme = ColourScheme{
	//fg
	DefaultFg:      [3]float32{1, 1, 1},
	BlackFg:        [3]float32{0, 0, 0},
	RedFg:          [3]float32{1, 0, 0},
	GreenFg:        [3]float32{0, 1, 0},
	YellowFg:       [3]float32{1, 1, 0},
	BlueFg:         [3]float32{0, 0, 1},
	MagentaFg:      [3]float32{1, 0, 1},
	CyanFg:         [3]float32{0, 1, 1},
	LightGreyFg:    [3]float32{0.7, 0.7, 0.7},
	DarkGreyFg:     [3]float32{0.3, 0.3, 0.3},
	LightRedFg:     [3]float32{1, 0.5, 0.5},
	LightGreenFg:   [3]float32{0.5, 1, 0.5},
	LightYellowFg:  [3]float32{1, 1, 0.5},
	LightBlueFg:    [3]float32{0.5, 0.5, 1},
	LightMagentaFg: [3]float32{1, 0.5, 1},
	LightCyanFg:    [3]float32{0.5, 1, 1},
	WhiteFg:        [3]float32{1, 1, 1},
	// bg
	DefaultBg:      [3]float32{0.1, 0.1, 0.1},
	BlackBg:        [3]float32{0, 0, 0},
	RedBg:          [3]float32{1, 0, 0},
	GreenBg:        [3]float32{0, 1, 0},
	YellowBg:       [3]float32{1, 1, 0},
	BlueBg:         [3]float32{0, 0, 1},
	MagentaBg:      [3]float32{1, 0, 1},
	CyanBg:         [3]float32{0, 1, 1},
	LightGreyBg:    [3]float32{0.7, 0.7, 0.7},
	DarkGreyBg:     [3]float32{0.3, 0.3, 0.3},
	LightRedBg:     [3]float32{1, 0.5, 0.5},
	LightGreenBg:   [3]float32{0.5, 1, 0.5},
	LightYellowBg:  [3]float32{1, 1, 0.5},
	LightBlueBg:    [3]float32{0.5, 0.5, 1},
	LightMagentaBg: [3]float32{1, 0.5, 1},
	LightCyanBg:    [3]float32{0.5, 1, 1},
	WhiteBg:        [3]float32{1, 1, 1},
}
