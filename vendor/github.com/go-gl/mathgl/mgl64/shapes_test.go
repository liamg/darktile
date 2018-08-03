// This file is generated from mgl32/shapes_test.go; DO NOT EDIT

// Copyright 2018 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.package mgl64_test

package mgl64

import (
	"testing"
)

func TestScreenToGLCoords(t *testing.T) {
	// use a small screen size in order to minimize errors due to fp rounding.
	const (
		sw = 100
		sh = 100
	)
	x, y := ScreenToGLCoords(0, sh-1, sw, sh)
	if x != -1 {
		t.Errorf("x = %f, expected -1.0", x)
	}
	if y != -1 {
		t.Errorf("y = %f, expected -1.0", y)
	}

	x, y = ScreenToGLCoords(sw-1, 0, sw, sh)
	if x != 1 {
		t.Errorf("x = %f, expected 1.0", x)
	}
	if y != 1 {
		t.Errorf("y = %f, expected 1.0", y)
	}
}

func TestGLToScreenCoords(t *testing.T) {
	const (
		sw = 100
		sh = 100
	)
	x, y := GLToScreenCoords(-1, -1, sw, sh)
	if x != 0 {
		t.Errorf("x = %d, expected 0", x)
	}
	if y != sh-1 {
		t.Errorf("y = %d, expected %d", y, sh-1)
	}

	x, y = GLToScreenCoords(1, 1, sw, sh)
	if x != sw-1 {
		t.Errorf("x = %d, expected %d", x, sw-1)
	}
	if y != 0 {
		t.Errorf("y = %d, expected 0", y)
	}
}
