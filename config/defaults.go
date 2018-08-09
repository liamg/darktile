package config

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
		Blue:         strToColourNoErr("#2e3340"),
		Magenta:      strToColourNoErr("#ff5879"),
		Cyan:         strToColourNoErr("#44b5b1"),
		LightGrey:    strToColourNoErr("#f2f1b9"),
		DarkGrey:     strToColourNoErr("#065f73"),
		LightRed:     strToColourNoErr("#ef5847"),
		LightGreen:   strToColourNoErr("#a2db91"),
		LightYellow:  strToColourNoErr("#beb090"),
		LightBlue:    strToColourNoErr("#61778d"),
		LightMagenta: strToColourNoErr("#ff99a1"),
		LightCyan:    strToColourNoErr("#9ed9d8"),
		White:        strToColourNoErr("#f6f6c9"),
	},
}
