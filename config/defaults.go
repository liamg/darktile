package config

import "runtime"

var DefaultConfig = Config{
	DebugMode: false,
	ColourScheme: ColourScheme{
		Cursor:       strToColourNoErr("#e8dfd6"),
		Foreground:   strToColourNoErr("#e8dfd6"),
		Background:   strToColourNoErr("#021b21"),
		Black:        strToColourNoErr("#000000"),
		Red:          strToColourNoErr("#800000"),
		Green:        strToColourNoErr("#008000"),
		Yellow:       strToColourNoErr("#808000"),
		Blue:         strToColourNoErr("#000080"),
		Magenta:      strToColourNoErr("#800080"),
		Cyan:         strToColourNoErr("#008080"),
		LightGrey:    strToColourNoErr("#f2f2f2"),
		DarkGrey:     strToColourNoErr("#808080"),
		LightRed:     strToColourNoErr("#ff0000"),
		LightGreen:   strToColourNoErr("#00ff00"),
		LightYellow:  strToColourNoErr("#ffff00"),
		LightBlue:    strToColourNoErr("#0000ff"),
		LightMagenta: strToColourNoErr("#ff00ff"),
		LightCyan:    strToColourNoErr("#00ffff"),
		White:        strToColourNoErr("#ffffff"),
		Selection:    strToColourNoErr("#333366"),
	},
	KeyMapping:            KeyMappingConfig(map[string]string{}),
	SearchURL:             "https://www.google.com/search?q=$QUERY",
	MaxLines:              1000,
	CopyAndPasteWithMouse: true,
	ShowVerticalScrollbar: true,
}

func init() {
	DefaultConfig.KeyMapping[string(ActionCopy)] = addMod("c")
	DefaultConfig.KeyMapping[string(ActionPaste)] = addMod("v")
	DefaultConfig.KeyMapping[string(ActionSearch)] = addMod("g")
	DefaultConfig.KeyMapping[string(ActionToggleDebug)] = addMod("d")
	DefaultConfig.KeyMapping[string(ActionToggleSlomo)] = addMod(";")
	DefaultConfig.KeyMapping[string(ActionReportBug)] = addMod("r")
}

func addMod(keys string) string {
	standardMod := "ctrl + shift + "

	if runtime.GOOS == "darwin" {
		standardMod = "super + "
	}

	return standardMod + keys
}
