// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mgl32

import (
	"math"
	"testing"
)

func TestHomogRotate3D(t *testing.T) {
	tests := []struct {
		Description string
		Angle       float32
		Axis        Vec3
		Expected    Mat4
	}{
		{
			"forward",
			0, Vec3{0, 0, 0},
			Ident4(),
		},
		{
			"heading 90 degree",
			DegToRad(90), Vec3{0, 1, 0},
			Mat4{
				0, 0, -1, 0,
				0, 1, 0, 0,
				1, 0, 0, 0,
				0, 0, 0, 1,
			},
		},
		{
			"heading 180 degree",
			DegToRad(180), Vec3{0, 1, 0},
			Mat4{
				-1, 0, 0, 0,
				0, 1, 0, 0,
				0, 0, -1, 0,
				0, 0, 0, 1,
			},
		},
		{
			"attitude 90 degree",
			DegToRad(90), Vec3{0, 0, 1},
			Mat4{
				0, 1, 0, 0,
				-1, 0, 0, 0,
				0, 0, 1, 0,
				0, 0, 0, 1,
			},
		},
		{
			"bank 90 degree",
			DegToRad(90), Vec3{1, 0, 0},
			Mat4{
				1, 0, 0, 0,
				0, 0, 1, 0,
				0, -1, 0, 0,
				0, 0, 0, 1,
			},
		},
		{
			"heading and attitude 90 degree",
			DegToRad(90), Vec3{0, 1, 1}.Normalize(),
			Mat4{
				0, 0.707107, -0.707107, 0,
				-0.707107, 0.5, 0.5, 0,
				0.707107, 0.5, 0.5, 0,
				0, 0, 0, 1,
			},
		},
		{
			"bank, heading and attitude 180 degree",
			DegToRad(180), Vec3{1, 1, 1}.Normalize(),
			Mat4{
				-1 / 3.0, 2 / 3.0, 2 / 3.0, 0,
				2 / 3.0, -1 / 3.0, 2 / 3.0, 0,
				2 / 3.0, 2 / 3.0, -1 / 3.0, 0,
				0, 0, 0, 1,
			},
		},
	}

	threshold := float32(math.Pow(10, -2))
	for _, c := range tests {
		if r := HomogRotate3D(c.Angle, c.Axis); !r.ApproxEqualThreshold(c.Expected, threshold) {
			t.Errorf("%v failed: HomogRotate3D(%v, %v) != %v (got %v)", c.Description, c.Angle, c.Axis, c.Expected, r)
		}
	}
}

func TestExtract3DScale(t *testing.T) {
	tests := []struct {
		M       Mat4
		X, Y, Z float32
	}{
		{
			Ident4(),
			1, 1, 1,
		}, {
			Scale3D(1, 2, 3),
			1, 2, 3,
		}, {
			Translate3D(10, 12, -5).Mul4(HomogRotate3D(math.Pi/2, Vec3{1, 0, 0})).Mul4(Scale3D(2, 3, 4)),
			2, 3, 4,
		},
	}

	eq := FloatEqualFunc(1e-6)
	for _, c := range tests {
		if x, y, z := Extract3DScale(c.M); !eq(x, c.X) || !eq(y, c.Y) || !eq(z, c.Z) {
			t.Errorf("ExtractScale(%v) != %v, %v, %v (got %v, %v, %v)", c.M, c.X, c.Y, c.Z, x, y, z)
		}
	}
}

func TestExtractMaxScale(t *testing.T) {
	tests := []struct {
		M Mat4
		V float32
	}{
		{
			Ident4(),
			1,
		}, {
			Scale3D(1, 2, 3),
			3,
		}, {
			Translate3D(10, 12, -5).Mul4(HomogRotate3D(math.Pi/2, Vec3{1, 0, 0})).Mul4(Scale3D(2, 3, 4)),
			4,
		},
	}

	eq := FloatEqualFunc(1e-6)
	for _, c := range tests {
		if r := ExtractMaxScale(c.M); !eq(r, c.V) {
			t.Errorf("ExtractMaxScale(%v) != %v (got %v)", c.M, c.V, r)
		}
	}
}

func TestTransformCoordinate(t *testing.T) {
	tests := [...]struct {
		v Vec3
		m Mat4

		out Vec3
	}{
		{Vec3{1, 1, 1}, Ident4(), Vec3{1, 1, 1}},
		{Vec3{1, 1, 1}, Translate3D(0, 1, 1).Mul4(Scale3D(2, 2, 2)), Vec3{2, 3, 3}},
		{Vec3{1, 1, -1}, Perspective(DegToRad(90), 4/3, 0, 100), Vec3{1, 1, 1}},
		{Vec3{0, 0, -100}, Perspective(DegToRad(45), 4/3, 0, 100), Vec3{0, 0, 1}},
		{Vec3{2, 2, -2}, Perspective(DegToRad(45), 4/3, 0, 100), Vec3{2.4142, 2.4142, 1}}, // sqrt(1^2+1^2)+1
	}

	for _, test := range tests {
		if v := TransformCoordinate(test.v, test.m); !test.out.ApproxEqualThreshold(v, 1e-4) {
			t.Errorf("TransformCoordinate on vector %v and matrix %v fails to give result %v (got %v)", test.v, test.m, test.out, v)
		}
	}
}

func TestTransformNormal(t *testing.T) {
	tests := [...]struct {
		v Vec3
		m Mat4

		out Vec3
	}{
		{Vec3{1, 1, 1}, Ident4(), Vec3{1, 1, 1}},
		{Vec3{1, 1, 1}, Translate3D(0, 1, 1).Mul4(Scale3D(2, 2, 2)), Vec3{2, 2, 2}},
	}

	for _, test := range tests {
		if v := TransformNormal(test.v, test.m); !test.out.ApproxEqualThreshold(v, 1e-4) {
			t.Errorf("TransformNormal on vector %v and matrix %v fails to give result %v (got %v)", test.v, test.m, test.out, v)
		}
	}
}
