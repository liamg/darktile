// This file is generated from mgl32/project_test.go; DO NOT EDIT

// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mgl64

import (
	"math"
	"testing"
)

func TestProject(t *testing.T) {
	t.Parallel()

	obj := Vec3{1002, 960, 0}
	modelview := Mat4{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 203, 1, 0, 1}
	projection := Mat4{0.0013020833721384406, 0, 0, 0, -0, -0.0020833334419876337, -0, -0, -0, -0, -1, -0, -1, 1, 0, 1}
	initialX, initialY, width, height := 0, 0, 1536, 960
	win := Project(obj, modelview, projection, initialX, initialY, width, height)
	answer := Vec3{1205.0000359117985, -1.0000501200556755, 0.5} // From glu.Project()

	if !win.ApproxEqualThreshold(answer, 1e-4) {
		t.Errorf("Project does something weird, differs from expected by of %v", win.Sub(answer).Len())
	}

	objr, err := UnProject(win, modelview, projection, initialX, initialY, width, height)
	if err != nil {
		t.Errorf("UnProject returned error: %v", err)
	}
	if !objr.ApproxEqualThreshold(obj, 1e-4) {
		t.Errorf("UnProject(%v) != %v (got %v)", win, obj, objr)
	}
}

// Test from
// http://stackoverflow.com/questions/38471708/opengl-glm-project-method-giving-unexpected-results
func TestProjectNonOneW(t *testing.T) {
	t.Parallel()

	obj := Vec3{5, 0, 0}

	projection := Perspective(
		DegToRad(45), // Field of view (45 degrees).
		800.0/600.0,  // Aspect ratio.
		0.1,          // Near Z at 0.1.
		10)           // Far Z at 10.
	camera := LookAtV(
		Vec3{0, 0.1, 10}, // Camera out on Z and slightly above.
		Vec3{0, 0, 0},    // Looking at the origin.
		Vec3{0, 1, 0})    // Up is positive Y.
	model := Ident4()               // Simple model matrix, to avoid confusion.
	modelView := camera.Mul4(model) // The model-view matrix (== camera, here).

	win := Project(obj, modelView, projection, 0, 0, 800, 600)

	t.Logf("Test:   (%v, %v, %v)", win[0], win[1], win[2])

	answer := Vec3{762.114, 300, 1} // verified with glm

	if !win.ApproxEqualThreshold(answer, 1e-4) {
		t.Errorf("Project does not properly do perspective division, differs from expected by %v", win.Sub(answer).Len())
	}
}

func TestUnprojectSingular(t *testing.T) {
	if _, err := UnProject(Vec3{}, Mat4{}, Mat4{}, 0, 0, 2048, 1152); err == nil {
		t.Errorf("Did not get error from UnProject on singular matrix")
	} else {
		t.Logf("Successfully got error on UnProject: %v", err)
	}
}

func TestLookAtV(t *testing.T) {
	// http://www.euclideanspace.com/maths/algebra/matrix/transforms/examples/index.htm

	tests := []struct {
		Description     string
		Eye, Center, Up Vec3
		Expected        Mat4
	}{
		{
			"forward",
			Vec3{0, 0, 0},
			Vec3{0, 0, -1},
			Vec3{0, 1, 0},
			Ident4(),
		},
		{
			"heading 90 degree",
			Vec3{0, 0, 0},
			Vec3{1, 0, 0},
			Vec3{0, 1, 0},
			Mat4{
				0, 0, -1, 0,
				0, 1, 0, 0,
				1, 0, 0, 0,
				0, 0, 0, 1,
			},
		},
		{
			"heading 180 degree",
			Vec3{0, 0, 0},
			Vec3{0, 0, 1},
			Vec3{0, 1, 0},
			Mat4{
				-1, 0, 0, 0,
				0, 1, 0, 0,
				0, 0, -1, 0,
				0, 0, 0, 1,
			},
		},
		{
			"attitude 90 degree",
			Vec3{0, 0, 0},
			Vec3{0, 0, -1},
			Vec3{1, 0, 0},
			Mat4{
				0, 1, 0, 0,
				-1, 0, 0, 0,
				0, 0, 1, 0,
				0, 0, 0, 1,
			},
		},
		{
			"bank 90 degree",
			Vec3{0, 0, 0},
			Vec3{0, -1, 0},
			Vec3{0, 0, -1},
			Mat4{
				1, 0, 0, 0,
				0, 0, 1, 0,
				0, -1, 0, 0,
				0, 0, 0, 1,
			},
		},
	}

	threshold := float64(math.Pow(10, -2))
	for _, c := range tests {
		if r := LookAtV(c.Eye, c.Center, c.Up); !r.ApproxEqualThreshold(c.Expected, threshold) {
			t.Errorf("%v failed: LookAtV(%v, %v, %v) != %v (got %v)", c.Description, c.Eye, c.Center, c.Up, c.Expected, r)
		}

		if r := LookAt(c.Eye[0], c.Eye[1], c.Eye[2], c.Center[0], c.Center[1], c.Center[2], c.Up[0], c.Up[1], c.Up[2]); !r.ApproxEqualThreshold(c.Expected, threshold) {
			t.Errorf("%v failed: LookAt(%v, %v, %v) != %v (got %v)", c.Description, c.Eye, c.Center, c.Up, c.Expected, r)
		}
	}
}

func TestOrtho(t *testing.T) {
	tests := []struct {
		Left, Right,
		Bottom, Top,
		Near, Far float64
		Expected Mat4
	}{
		{
			-1.0, 1.0, -1.0, 1.0, 1.0, -1.0,
			Ident4(),
		}, {
			-10.0, 10.0, -10.0, 10.0, 0.0, 100.0,
			Mat4{0.1, 0.0, 0.0, 0.0, 0.0, 0.1, 0.0, 0.0, 0.0, 0.0, -0.02, 0.0, 0.0, 0.0, -1.0, 1.0},
		}, {
			0.0, 10.0, 0.0, 10.0, 0.0, 100.0,
			Mat4{0.2, 0.0, 0.0, 0.0, 0.0, 0.2, 0.0, 0.0, 0.0, 0.0, -0.02, 0.0, -1.0, -1.0, -1.0, 1.0},
		},
	}

	for _, c := range tests {
		if r := Ortho(c.Left, c.Right, c.Bottom, c.Top, c.Near, c.Far); !r.ApproxEqualThreshold(c.Expected, 1e-4) {
			t.Errorf("Ortho(%v, %v, %v, %v, %v, %v) != %v (got %v)", c.Left, c.Right, c.Bottom, c.Top, c.Near, c.Far, c.Expected, r)
		}
	}
}

func TestOrtho2D(t *testing.T) {
	tests := []struct {
		Left, Right,
		Bottom, Top float64
		Expected Mat4
	}{
		{
			-1.0, 1.0, -1.0, 1.0,
			Mat4{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, -1, 0, 0, 0, 0, 1},
		}, {
			-10.0, 10.0, -10.0, 10.0,
			Mat4{0.1, 0.0, 0.0, 0.0, 0.0, 0.1, 0.0, 0.0, 0.0, 0.0, -1.0, 0.0, 0.0, 0.0, 0.0, 1.0},
		}, {
			0.0, 10.0, 0.0, 10.0,
			Mat4{0.2, 0.0, 0.0, 0.0, 0.0, 0.2, 0.0, 0.0, 0.0, 0.0, -1.0, 0.0, -1.0, -1.0, 0.0, 1.0},
		},
	}

	for _, c := range tests {
		if r := Ortho2D(c.Left, c.Right, c.Bottom, c.Top); !r.ApproxEqualThreshold(c.Expected, 1e-4) {
			t.Errorf("Ortho2D(%v, %v, %v, %v) != %v (got %v)", c.Left, c.Right, c.Bottom, c.Top, c.Expected, r)
		}
	}
}

func TestPerspective(t *testing.T) {
	tests := []struct {
		Fovy, Aspect,
		Near, Far float64
		Expected Mat4
	}{
		{
			DegToRad(450.0), 1.0, -1.0, 1.0,
			Mat4{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, -1, 0, 0, 1, 0},
		}, {
			DegToRad(45.0), 4.0 / 3.0, 0.1, 100.0,
			Mat4{1.810660, 0.0, 0.0, 0.0, 0.0, 2.4142134, 0.0, 0.0, 0.0, 0.0, -1.002002, -1.0, 0.0, 0.0, -0.2002002, 0.0},
		}, {
			DegToRad(90.0), 16.0 / 9.0, -1.0, 1.0,
			Mat4{0.562500, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, -0.0, -1.0, 0.0, 0.0, 1.0, 0.0},
		},
	}

	for _, c := range tests {
		if r := Perspective(c.Fovy, c.Aspect, c.Near, c.Far); !r.ApproxEqualThreshold(c.Expected, 1e-4) {
			t.Errorf("Perspective(%v, %v, %v, %v) != %v (got %v)", c.Fovy, c.Aspect, c.Near, c.Far, c.Expected, r)
		}
	}
}

func TestFrustum(t *testing.T) {
	tests := []struct {
		Left, Right,
		Bottom, Top,
		Near, Far float64
		Expected Mat4
	}{
		{
			-1.0, 1.0, -1.0, 1.0, 1.0, 2.0,
			Mat4{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, -3, -1, 0, 0, -4, 0},
		},
		// TODO: more tests
	}

	for _, c := range tests {
		if r := Frustum(c.Left, c.Right, c.Bottom, c.Top, c.Near, c.Far); !r.ApproxEqualThreshold(c.Expected, 1e-4) {
			t.Errorf("Frustum(%v, %v, %v, %v, %v, %v) != %v (got %v)", c.Left, c.Right, c.Bottom, c.Top, c.Near, c.Far, c.Expected, r)
		}
	}
}
