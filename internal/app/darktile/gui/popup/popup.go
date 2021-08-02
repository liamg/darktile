package popup

import (
	"image/color"
	"time"
)

type Message struct {
	Text       string
	Expiry     time.Time
	Foreground color.Color
	Background color.Color
}
