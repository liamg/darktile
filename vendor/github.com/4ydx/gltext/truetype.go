// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gltext

import (
	"errors"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"io"
	"io/ioutil"
	"sort"
)

// RuneRanges specify the rune ranges for ordered disjoint subsets of the ttf
// EG 32 - 127, 5000 - 6000 will created a more compact bitmap that holds the
// specified ranges of runes.
type RuneRange struct {
	Low, High rune
}

type RuneRanges []RuneRange

func (rr RuneRanges) Len() int           { return len(rr) }
func (rr RuneRanges) Swap(i, j int)      { rr[i], rr[j] = rr[j], rr[i] }
func (rr RuneRanges) Less(i, j int) bool { return rr[i].Low < rr[j].Low }

func (rr RuneRanges) Validate() bool {
	sort.Sort(rr)
	previousMax := rune(0)
	for _, r := range rr {
		if r.Low <= previousMax {
			return false
		}
		if r.Low > r.High {
			return false
		}
		previousMax = r.High
	}
	return true
}

// GetGlyphIndex returns the location of the glyph data within
// the compressed rune ranges covered by the font
// EG if runes 0-25, 100-110 are supported by the font then
// the actual location of 100 will be in position 26 in the png image
func (rr RuneRanges) GetGlyphIndex(char rune) rune {
	var index, offset rune
	index = -1
	for _, runes := range rr {
		if char >= runes.Low && char <= runes.High {
			index = char - runes.Low + offset
		}
		offset += runes.High - runes.Low + 1
	}
	return index
}

// http://www.freetype.org/freetype2/docs/tutorial/step2.html

// LoadTruetype loads a truetype font from the given stream and
// applies the given font scale in points.
//
// The low and high values determine the lower and upper rune limits
// we should load for this font. For standard ASCII this would be: 32, 127.
func NewTruetypeFontConfig(r io.Reader, scale fixed.Int26_6, runeRanges RuneRanges, runesPerRow fixed.Int26_6) (*FontConfig, error) {
	if !runeRanges.Validate() {
		return nil, errors.New("Invalid rune ranges supplied.")
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Read the truetype font.
	ttf, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}

	// Create our FontConfig type.
	fc := &FontConfig{}
	length := rune(0)
	for _, r := range runeRanges {
		length += r.High - r.Low + 1
	}
	fc.RuneRanges = runeRanges
	fc.Glyphs = make(Charset, int(length))

	// Create an image, large enough to store all requested glyphs.
	// The resulting image is set to power of 2 dimensions so it might be wise to adjust the runesPerRow
	// parameter to ensure that unnecessary space isn't created based on the character set being used
	gc := fixed.Int26_6(len(fc.Glyphs))
	runesPerCol := (gc / runesPerRow) + 1

	gb := ttf.Bounds(scale)
	gw := (gb.Max.X - gb.Min.X)
	gh := (gb.Max.Y - gb.Min.Y)

	iw := Pow2(uint32(gw * runesPerRow))
	ih := Pow2(uint32(gh * runesPerCol))

	fg, bg := image.White, image.Transparent
	rect := image.Rect(0, 0, int(iw), int(ih))
	fc.Image = image.NewNRGBA(rect)
	draw.Draw(fc.Image, fc.Image.Bounds(), bg, image.ZP, draw.Src)

	// Use a freetype context to do the drawing.
	c := freetype.NewContext()
	c.SetDPI(72) // Do not change this.  It is required in order to have a properly aligned bounding box!!!
	c.SetFont(ttf)
	c.SetFontSize(float64(scale))
	c.SetClip(fc.Image.Bounds())
	c.SetDst(fc.Image)
	c.SetSrc(fg)

	// Iterate over all relevant glyphs in the truetype font and draw them all to the image buffer
	// Add Glyph objects to track various glyph values
	var gi fixed.Int26_6
	var gx, gy fixed.Int26_6

	for _, runeRange := range fc.RuneRanges {
		for ch := runeRange.Low; ch <= runeRange.High; ch++ {
			index := ttf.Index(ch)
			metric := ttf.HMetric(scale, index)

			if gi%runesPerRow == 0 {
				gx = 0
				if gi > 0 {
					gy += gh
				}
			} else {
				gx += gw
			}
			fc.Glyphs[gi].Advance = int(metric.AdvanceWidth)
			fc.Glyphs[gi].X = int(gx)
			fc.Glyphs[gi].Y = int(gy)
			fc.Glyphs[gi].Width = int(gw)
			fc.Glyphs[gi].Height = int(gh)

			pt := freetype.Pt(int(gx), int(gy)+int(c.PointToFixed(float64(scale))>>6))
			c.DrawString(string(ch), pt)
			gi++
		}
	}
	return fc, nil
}

func LoadTruetypeFontConfig(rootPath, name string) (*FontConfig, error) {
	fc := &FontConfig{}
	fc.Name = name

	err := fc.Load(rootPath)
	if err != nil {
		return nil, err
	}
	return fc, nil
}
