// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package v41

import (
	"github.com/4ydx/gltext"
	"github.com/go-gl/mathgl/mgl32"
	"testing"
)

func TestHasRune(t *testing.T) {
	f := &Font{}
	f.Config = &gltext.FontConfig{}
	f.Config.Glyphs = make(gltext.Charset, 100)
	f.Config.RuneRanges = make(gltext.RuneRanges, 0)

	r := gltext.RuneRange{Low: 30, High: 40}
	f.Config.RuneRanges = append(f.Config.RuneRanges, r)
	r = gltext.RuneRange{Low: 100, High: 400}
	f.Config.RuneRanges = append(f.Config.RuneRanges, r)

	if !f.Config.RuneRanges.Validate() {
		t.Error("Not validating properly.")
	}
	text := &Text{}
	text.Font = f
	if !text.HasRune(40) {
		t.Error("Missing rune 40.")
	}
	if text.HasRune(41) {
		t.Error("Should not have 41.")
	}
}

// TestClickedCharacter tests a hypothetical string of length 3 with variable width chars
func TestClickedCharacter(t *testing.T) {
	text := &Text{}
	text.Font = &Font{}
	text.Font.WindowWidth = 100
	text.X1.X = -20
	text.String = "ABC"

	// click was just around the middle of the screen
	xPos := float64(51)

	// -20 to -10 is A
	// -10 to +10 is B
	// +10 to +20 is C
	text.CharSpacing = make([]float32, 0)
	text.CharSpacing = append(text.CharSpacing, 10)
	text.CharSpacing = append(text.CharSpacing, 20)
	text.CharSpacing = append(text.CharSpacing, 10)

	index, side := text.ClickedCharacter(xPos, 0)
	if index != 1 {
		t.Error("Expecting index 1")
	}
	if side != CSRight {
		t.Error("Expecting right side click")
	}
}

func TestBoundingBox(t *testing.T) {
	text := &Text{}
	text.X1 = gltext.Point{-10, -10}
	text.X2 = gltext.Point{+10, +10}
	text.Font = &Font{}
	v := mgl32.Vec2{10, 5}
	text.SetPosition(v)
	x1, x2 := text.GetBoundingBox()
	if x1.X != 0 || x1.Y != -5 {
		t.Error(x1)
	}
	if x2.X != 20 || x2.Y != 15 {
		t.Error(x2)
	}
}
