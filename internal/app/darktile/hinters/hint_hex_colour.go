package hinters

import (
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"regexp"

	"github.com/liamg/darktile/internal/app/darktile/termutil"
)

func init() {
	register(&HexColourHinter{}, PriorityLow)
}

var hexColourRegex = regexp.MustCompile(`#[0-9A-Fa-f]{6}`)

type HexColourHinter struct {
}

func (h *HexColourHinter) Match(text string, cursorIndex int) (matched bool, offset int, length int) {
	matches := hexColourRegex.FindAllStringIndex(text, -1)
	for _, match := range matches {
		if match[0] <= cursorIndex && match[1] > cursorIndex {
			return true, match[0], match[1] - match[0]
		}
	}
	return
}

func (h *HexColourHinter) Activate(api HintAPI, match string, start termutil.Position, end termutil.Position) error {
	colourBytes, err := hex.DecodeString(match[1:])
	if err != nil {
		return err
	}

	cellSize := api.CellSize()
	size := image.Rectangle{image.Point{}, image.Point{
		X: cellSize.X * 18,
		Y: cellSize.Y,
	}}
	img := image.NewRGBA(size)
	for x := 0; x < size.Dx(); x++ {
		for y := 0; y < size.Dy(); y++ {
			img.SetRGBA(x, y, color.RGBA{
				R: colourBytes[0],
				G: colourBytes[1],
				B: colourBytes[2],
				A: 0xff,
			})
		}
	}

	api.Highlight(start, end, fmt.Sprintf(
		`Hex: %s
RGB: %d, %d, %d`,
		match,
		colourBytes[0],
		colourBytes[1],
		colourBytes[2],
	), img)
	return nil
}

func (h *HexColourHinter) Deactivate(api HintAPI) error {
	api.ClearHighlight()
	return nil
}

func (h *HexColourHinter) Click(api HintAPI) error {
	return nil
}
