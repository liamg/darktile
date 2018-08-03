// This file is generated from mgl32/conv_test.go; DO NOT EDIT

// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mgl64

import (
	"math"
	"testing"
)

func TestCartesianToSphere(t *testing.T) {
	t.Parallel()

	v := Vec3{5, 12, 9}

	r, theta, phi := CartesianToSpherical(v)

	if !FloatEqualThreshold(r, 15.8114, 1e-4) {
		t.Errorf("Got incorrect value for radius. Got: %f, expected: %f", r, 15.8114)
	}

	if !FloatEqualThreshold(theta, 0.965250852, 1e-4) {
		t.Errorf("Got incorrect value for theta. Got: %f, expected: %f", theta, 0.965250852)
	}

	if !FloatEqualThreshold(phi, 1.1760046, 1e-4) {
		t.Errorf("Got incorrect value for phi. Got: %f, expected: %f", phi, 1.1760046)
	}
}

func TestSphereToCartesian(t *testing.T) {
	t.Parallel()

	v := Vec3{5, 12, 9}

	result := SphericalToCartesian(15.8114, 0.965250852, 1.1760046)

	if !v.ApproxEqualThreshold(result, 1e-4) {
		t.Errorf("Got incorrect vector. Got: %v, Expected: %v", result, v)
	}
}

func TestCartesianToCylinder(t *testing.T) {
	t.Parallel()

	v := Vec3{5, 12, 9}

	rho, phi, z := CartesianToCylindical(v)

	if !FloatEqualThreshold(rho, 13, 1e-4) {
		t.Errorf("Got incorrect value for radius. Got: %f, expected: %f", rho, 13.)
	}

	if !FloatEqualThreshold(phi, 1.17601, 1e-4) {
		t.Errorf("Got incorrect value for theta. Got: %f, expected: %f", phi, 1.17601)
	}

	if !FloatEqualThreshold(z, 9, 1e-4) {
		t.Errorf("Got incorrect value for phi. Got: %f, expected: %f", z, 9.)
	}
}

func TestCylinderToCartesian(t *testing.T) {
	t.Parallel()

	v := Vec3{5, 12, 9}

	result := CylindricalToCartesian(13, 1.17601, 9)

	if !v.ApproxEqualThreshold(result, 1e-4) {
		t.Errorf("Got incorrect vector. Got: %v, expected: %v", result, v)
	}
}

func TestCylinderToSphere(t *testing.T) {
	t.Parallel()

	r, theta, phi := CylindircalToSpherical(13, 1.17601, 9)

	if !FloatEqualThreshold(r, 15.8114, 1e-4) {
		t.Errorf("Got incorrect value for radius. Got: %f, expected: %f", r, 15.8114)
	}

	if !FloatEqualThreshold(theta, 0.965250852, 1e-4) {
		t.Errorf("Got incorrect value for theta. Got: %f, expected: %f", theta, 0.965250852)
	}

	if !FloatEqualThreshold(phi, 1.1760046, 1e-4) {
		t.Errorf("Got incorrect value for phi. Got: %f, expected: %f", phi, 1.1760046)
	}
}

func TestSphereToCylinder(t *testing.T) {
	t.Parallel()

	rho, phi, z := SphericalToCylindrical(15.8114, 0.965250852, 1.1760046)

	if !FloatEqualThreshold(rho, 13, 1e-4) {
		t.Errorf("Got incorrect value for radius. Got: %f, expected: %f", rho, 13.)
	}

	if !FloatEqualThreshold(phi, 1.17601, 1e-4) {
		t.Errorf("Got incorrect value for theta. Got: %f, expected: %f", phi, 1.17601)
	}

	if !FloatEqualThreshold(z, 9, 1e-4) {
		t.Errorf("Got incorrect value for phi. Got: %f, expected: %f", z, 9.)
	}
}

func TestDegToRad(t *testing.T) {
	tests := []struct {
		Deg, Rad float64
	}{
		{0, 0},
		{90, math.Pi / 2},
		{180, math.Pi},
		{270, math.Pi + math.Pi/2},
		{360, math.Pi * 2},
		{-90, -math.Pi / 2},
		{-360, -math.Pi * 2},
	}

	for _, c := range tests {
		if r := DegToRad(c.Deg); r != c.Rad {
			t.Errorf("DegToRad(%v) != %v (got %v)", c.Deg, c.Rad, r)
		}
	}
}

func TestRadToDeg(t *testing.T) {
	tests := []struct {
		Deg, Rad float64
	}{
		{0, 0},
		{90, math.Pi / 2},
		{180, math.Pi},
		{270, math.Pi + math.Pi/2},
		{360, math.Pi * 2},
		{-90, -math.Pi / 2},
		{-360, -math.Pi * 2},
	}

	for _, c := range tests {
		if r := RadToDeg(c.Rad); r != c.Deg {
			t.Errorf("RadToDeg(%v) != %v (got %v)", c.Rad, c.Deg, r)
		}
	}
}
