package config

import "runtime"

var DefaultConfig = Config{
	DebugMode: false,
	ColourScheme: ColourScheme{
		Cursor:       strToColourNoErr("#e8dfd6"),
		Foreground:   strToColourNoErr("#e8dfd6"),
		Background:   strToColourNoErr("#021b21"),
		Black:        strToColourNoErr("#032c36"),
		Red:          strToColourNoErr("#c2454e"),
		Green:        strToColourNoErr("#7cbf9e"),
		Yellow:       strToColourNoErr("#8a7a63"),
		Blue:         strToColourNoErr("#065f73"),
		Magenta:      strToColourNoErr("#ff5879"),
		Cyan:         strToColourNoErr("#44b5b1"),
		LightGrey:    strToColourNoErr("#f2f1b9"),
		DarkGrey:     strToColourNoErr("#3e4360"),
		LightRed:     strToColourNoErr("#ef5847"),
		LightGreen:   strToColourNoErr("#a2db91"),
		LightYellow:  strToColourNoErr("#beb090"),
		LightBlue:    strToColourNoErr("#61778d"),
		LightMagenta: strToColourNoErr("#ff99a1"),
		LightCyan:    strToColourNoErr("#9ed9d8"),
		White:        strToColourNoErr("#f6f6c9"),
		Selection:    strToColourNoErr("#333366"),
	},
	KeyMapping: KeyMappingConfig(map[string]string{}),
	SearchURL:  "https://www.google.com/search?q=$QUERY",
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
