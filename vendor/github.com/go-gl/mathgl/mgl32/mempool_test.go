// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mgl32

import (
	"testing"
)

func TestBinLog(t *testing.T) {
	tests := []struct {
		in int

		// out
		val   int
		exact bool
	}{
		{-256, -1, false},
		{0, -1, false},
		{1, 0, true},
		{2, 1, true},
		{3, 1, false},
		{32, 5, true},
		{37, 5, false},
	}

	for _, test := range tests {
		outV, outE := binLog(test.in)
		if outV != test.val || outE != test.exact {
			t.Errorf("binLog gives incorrect result for input %v. Got: (%v,%v); Expected: (%v,%v)", test.in, outV, outE, test.val, test.exact)
		}
	}
}

func TestGetPool(t *testing.T) {
	slicePools = nil
	pool := getPool(3)

	if len(slicePools) != 4 || pool == nil {
		t.Errorf("Couldn't get pool. Size of slice %v (should be 4)", len(slicePools))
	}

	slice, ok := pool.Get().([]float32)
	if slice == nil || !ok || cap(slice) != 1<<3 {
		t.Errorf("Slice from pool either not allocated, not ok, or of wrong cap. Got slice: %v, ok: %v, cap: %v", slice, ok, cap(slice))
	}
}

func TestGrabFromPool(t *testing.T) {
	slicePools = nil
	slice := grabFromPool(17)

	if slice == nil || len(slice) != 17 || cap(slice) != 32 {
		t.Errorf("Got bad, ill sized, or badly capped slice from grabFromPool. Slice: %v, len: %v, cap: %v", slice, len(slice), cap(slice))
	}
}

func BenchmarkBinLogReasonable(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = binLog(100)
	}
}

func BenchmarkBinLogBig(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = binLog(1<<30 + 1)
	}
}

func BenchmarkBinLogSmall(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = binLog(10)
	}
}
