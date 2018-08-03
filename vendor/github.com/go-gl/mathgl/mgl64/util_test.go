// This file is generated from mgl32/util_test.go; DO NOT EDIT

// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mgl64

import (
	"math/rand"
	"testing"
	"time"
)

func TestEqual(t *testing.T) {
	t.Parallel()

	var a float64 = 1.5
	var b float64 = 1.0 + .5

	if !FloatEqual(a, a) {
		t.Errorf("Float Equal fails on comparing a number with itself")
	}

	if !FloatEqual(a, b) {
		t.Errorf("Float Equal fails to compare two equivalent numbers with minimal drift")
	} else if !FloatEqual(b, a) {
		t.Errorf("Float Equal is not symmetric for some reason")
	}

	if !FloatEqual(0.0, 0.0) {
		t.Errorf("Float Equal fails to compare zero values correctly")
	}

	if FloatEqual(1.5, 1.51) {
		t.Errorf("Float Equal gives false positive on large difference")
	}

	if FloatEqual(1.5, 1.5000001) {
		t.Errorf("Float Equal gives false positive on small difference")
	}

	if FloatEqual(1.5, 0.0) {
		t.Errorf("Float Equal gives false positive comparing with zero")
	}
}

func TestEqualThreshold(t *testing.T) {
	t.Parallel()

	// |1.0 - 1.01| < .1
	if !FloatEqualThreshold(1.0, 1.01, 1e-1) {
		t.Errorf("Thresholded equal returns negative on threshold")
	}

	// Comes out to |1.0 - 1.01| < .0001
	if FloatEqualThreshold(1.0, 1.01, 1e-3) {
		t.Errorf("Thresholded equal returns false positive on tolerant threshold")
	}
}

func TestEqualThresholdTable(t *testing.T) {
	// http://floating-point-gui.de/errors/NearlyEqualsTest.java

	tests := []struct {
		A, B, Ep float64
		Expected bool
	}{
		{1.0, 1.01, 1e-1, true},
		{1.0, 1.01, 1e-3, false},

		// Regular large numbers
		{1000000.0, 1000001.0, 0.00001, true},
		{1000001.0, 1000000.0, 0.00001, true},
		{10000.0, 10001.0, 0.00001, false},
		{10001.0, 10000.0, 0.00001, false},

		// Negative large numbers
		{-1000000.0, -1000001.0, 0.00001, true},
		{-1000001.0, -1000000.0, 0.00001, true},
		{-10000.0, -10001.0, 0.00001, false},
		{-10001.0, -10000.0, 0.00001, false},

		// Numbers around 1
		{1.0000001, 1.0000002, 0.00001, true},
		{1.0000002, 1.0000001, 0.00001, true},
		{1.0002, 1.0001, 0.00001, false},
		{1.0001, 1.0002, 0.00001, false},

		// Numbers around -1
		{-1.000001, -1.000002, 0.00001, true},
		{-1.000002, -1.000001, 0.00001, true},
		{-1.0001, -1.0002, 0.00001, false},
		{-1.0002, -1.0001, 0.00001, false},

		// Numbers between 1 and 0
		{0.000000001000001, 0.000000001000002, 0.00001, true},
		{0.000000001000002, 0.000000001000001, 0.00001, true},
		{0.000000000001002, 0.000000000001001, 0.00001, false},
		{0.000000000001001, 0.000000000001002, 0.00001, false},

		// Numbers between -1 and 0
		{-0.000000001000001, -0.000000001000002, 0.00001, true},
		{-0.000000001000002, -0.000000001000001, 0.00001, true},
		{-0.000000000001002, -0.000000000001001, 0.00001, false},
		{-0.000000000001001, -0.000000000001002, 0.00001, false},

		// Comparisons involving zero
		{0.0, 0.0, 0.00001, true},
		{0.0, -0.0, 0.00001, true},
		{-0.0, -0.0, 0.00001, true},
		{0.00000001, 0.0, 0.00001, false},
		{0.0, 0.00000001, 0.00001, false},
		{-0.00000001, 0.0, 0.00001, false},
		{0.0, -0.00000001, 0.00001, false},

		// Comparisons involving infinities
		{InfPos, InfPos, 0.00001, true},
		{InfNeg, InfNeg, 0.00001, true},
		{InfNeg, InfPos, 0.00001, false},
		{InfPos, MaxValue, 0.00001, false},
		{InfNeg, -MaxValue, 0.00001, false},

		// Comparisons involving NaN values
		{NaN, NaN, 0.00001, false},
		{0.0, NaN, 0.00001, false},
		{NaN, 0.0, 0.00001, false},
		{-0.0, NaN, 0.00001, false},
		{NaN, -0.0, 0.00001, false},
		{NaN, InfPos, 0.00001, false},
		{InfPos, NaN, 0.00001, false},
		{NaN, InfNeg, 0.00001, false},
		{InfNeg, NaN, 0.00001, false},
		{NaN, MaxValue, 0.00001, false},
		{MaxValue, NaN, 0.00001, false},
		{NaN, -MaxValue, 0.00001, false},
		{-MaxValue, NaN, 0.00001, false},
		{NaN, MinValue, 0.00001, false},
		{MinValue, NaN, 0.00001, false},
		{NaN, -MinValue, 0.00001, false},
		{-MinValue, NaN, 0.00001, false},

		// Comparisons of numbers on opposite sides of 0
		{1.000000001, -1.0, 0.00001, false},
		{-1.0, 1.000000001, 0.00001, false},
		{-1.000000001, 1.0, 0.00001, false},
		{1.0, -1.000000001, 0.00001, false},
		{10 * MinValue, 10 * -MinValue, 0.00001, true},
		{10000 * MinValue, 10000 * -MinValue, 0.00001, true},

		// Comparisons of numbers very close to zero
		{MinValue, -MinValue, 0.00001, true},
		{-MinValue, MinValue, 0.00001, true},
		{MinValue, 0, 0.00001, true},
		{0, MinValue, 0.00001, true},
		{-MinValue, 0, 0.00001, true},
		{0, -MinValue, 0.00001, true},
		{0.000000001, -MinValue, 0.00001, false},
		{0.000000001, MinValue, 0.00001, false},
		{MinValue, 0.000000001, 0.00001, false},
		{-MinValue, 0.000000001, 0.00001, false},
	}

	for _, c := range tests {
		if r := FloatEqualThreshold(c.A, c.B, c.Ep); r != c.Expected {
			t.Errorf("FloatEqualThreshold(%v, %v, %v) != %v (got %v)", c.A, c.B, c.Ep, c.Expected, r)
		}
	}
}

func TestEqual32(t *testing.T) {
	t.Parallel()

	a := float64(1.5)
	b := float64(1.0 + .5)

	if !FloatEqual(a, a) {
		t.Errorf("Float Equal fails on comparing a number with itself")
	}

	if !FloatEqual(a, b) {
		t.Errorf("Float Equal fails to compare two equivalent numbers with minimal drift")
	} else if !FloatEqual(b, a) {
		t.Errorf("Float Equal is not symmetric for some reason")
	}

	if !FloatEqual(0.0, 0.0) {
		t.Errorf("Float Equal fails to compare zero values correctly")
	}

	if FloatEqual(1.5, 1.51) {
		t.Errorf("Float Equal gives false positive on large difference")
	}

	if FloatEqual(1.5, 0.0) {
		t.Errorf("Float Equal gives false positive comparing with zero")
	}
}

func TestClampf(t *testing.T) {
	t.Parallel()

	if !FloatEqual(Clamp(-1.0, 0.0, 1.0), 0.0) {
		t.Errorf("Clamp returns incorrect value for below threshold")
	}

	if !FloatEqual(Clamp(0.0, 0.0, 1.0), 0.0) {
		t.Errorf("Clamp does something weird when value is at threshold")
	}

	if !FloatEqual(Clamp(.14, 0.0, 1.0), .14) {
		t.Errorf("Clamp fails to return correct value when value is within threshold")
	}

	if !FloatEqual(Clamp(1.1, 0.0, 1.0), 1.0) {
		t.Errorf("Clamp fails to return max threshold when appropriate")
	}
}

func TestIsClamped(t *testing.T) {
	t.Parallel()

	if IsClamped(-1.0, 0.0, 1.0) {
		t.Errorf("Test below min is considered clamped")
	}

	if !IsClamped(.15, 0.0, 1.0) {
		t.Errorf("Test in threshold returns false")
	}

	if IsClamped(1.5, 0.0, 1.0) {
		t.Errorf("Test above max threshold returns false positive")
	}
}

/* These benchmarks probably aren't very interesting, there's not really many ways to optimize the functions they're benchmarking */

func BenchmarkEqual(b *testing.B) {
	b.StopTimer()
	r := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		f1 := r.Float64()
		f2 := r.Float64()
		b.StartTimer()

		FloatEqual(f1, f2)
	}
}

// Here just to get a baseline of how much worse the safer equal is
func BenchmarkBuiltinEqual(b *testing.B) {
	b.StopTimer()
	r := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		f1 := r.Float64()
		f2 := r.Float64()
		b.StartTimer()

		_ = f1 == f2
	}
}

func BenchmarkClampf(b *testing.B) {
	b.StopTimer()
	r := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		a := r.Float64()
		t1 := r.Float64()
		t2 := r.Float64()
		b.StartTimer()

		Clamp(a, t1, t2)
	}
}

func TestRound(t *testing.T) {
	tests := []struct {
		Value     float64
		Precision int
		Expected  float64
	}{
		{0.5, 0, 1},
		{0.123, 2, 0.12},
		{9.99999999, 6, 10},
		{-9.99999999, 6, -10},
		{-0.000099, 4, -0.0001},
	}

	for _, c := range tests {
		if r := Round(c.Value, c.Precision); r != c.Expected {
			t.Errorf("Round(%v, %v) != %v (got %v)", c.Value, c.Precision, c.Expected, r)
		}
	}
}

func BenchmarkRound(b *testing.B) {
	b.StopTimer()
	r := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		v := r.Float64()
		p := r.Intn(10)
		b.StartTimer()

		Round(v, p)
	}
}
