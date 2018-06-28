// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gltext

// A Glyph describes metrics for a single font glyph.
// These indicate which area of a given image contains the
// glyph data and how the glyph should be spaced in a rendered string.
type Point struct {
	X float32
	Y float32
}

type Glyph struct {
	X      int `json:"x"`      // The x location of the glyph on a sprite sheet.
	Y      int `json:"y"`      // The y location of the glyph on a sprite sheet.
	Width  int `json:"width"`  // The width of the glyph on a sprite sheet.
	Height int `json:"height"` // The height of the glyph on a sprite sheet.

	// Advance determines the distance to the next glyph.
	// This is used to properly align non-monospaced fonts.
	Advance int `json:"advance"`
}

func (g *Glyph) GetTexturePositions(font FontLike) (tP1, tP2 Point) {
	// Quad width/height

	// Originally the ttf width value was being used.  This, however, differs from the Advance value.
	// This has been changed to advance so that the resulting quads that are generated for text to not
	// overlap one another.
	vw := float32(g.Advance)

	vh := float32(g.Height)

	// Unfortunately with the current font, if I don't add a small offset to the Y axis location
	// the bottom edge of the character above might appear.
	//
	// EG:
	// Wrapping 16 characters per line:
	// runesPerRow := fixed.Int26_6(16)
	// runeRanges := make(gltext.RuneRanges, 0)
	// runeRange := gltext.RuneRange{Low: 1, High: 128}
	// runeRanges = append(runeRanges, runeRange)
	//
	// The resulting image file will place "g" above "w".  The very bottom edge of "g" will show up
	// when using the "w" character in a line of text. So the dirty hack is to remove just a bit of
	// the original top as per below.  This is not ideal.  Either I am not understanding something
	// about the glyph layout or this will have to be tweaked based on the font being used.
	// See the file example_image.png.

	// texture point 1
	tP1 = Point{X: float32(g.X) / font.GetTextureWidth(), Y: float32(g.Y) / font.GetTextureHeight()}

	// texture point 2
	tP2 = Point{X: (float32(g.X) + vw) / font.GetTextureWidth(), Y: (float32(g.Y) + vh) / font.GetTextureHeight()}

	return
}

// A Charset represents a set of glyph descriptors for a font.
// Each glyph descriptor holds glyph metrics which are used to
// properly align the given glyph in the resulting rendered string.
type Charset []Glyph

// Scale scales all glyphs by the given factor and repositions them
// appropriately. A scale of 1 retains the original size. A scale of 2
// doubles the size of each glyph, etc.
//
// This is useful when the accompanying sprite sheet is scaled by the
// same factor. In this case, we want the glyph data to match up with the
// new image.
func (c Charset) Scale(factor int) {
	if factor <= 1 {
		// A factor of zero results in zero-sized glyphs and
		// is therefore not valid. A factor of 1 does not change
		// the glyphs, so we can ignore it.
		return
	}

	// Multiply each glyph field by the given factor
	// to scale them up to the new size.
	for i := range c {
		c[i].X *= factor
		c[i].Y *= factor
		c[i].Width *= factor
		c[i].Height *= factor
		c[i].Advance *= factor
	}
}
