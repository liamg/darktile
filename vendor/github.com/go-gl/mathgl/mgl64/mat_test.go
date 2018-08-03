// This file is generated from mgl32/mat_test.go; DO NOT EDIT

// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mgl64

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestMulIdent(t *testing.T) {
	t.Parallel()

	i1 := [...]float64{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
	i2 := Ident4()
	i3 := Ident4()

	mul := i2.Mul4(i3)

	for i := range mul {
		if mul[i] != i1[i] {
			t.Errorf("Multiplication of identities does not yield identity")
		}
	}
}

// M>N
func TestMatRowsTall(t *testing.T) {
	t.Parallel()

	row0 := Vec3{1, 2, 3}
	row1 := Vec3{4, 5, 6}
	rows := [2]Vec3{row0, row1}

	m1 := Mat2x3FromRows(rows[0], rows[1])

	t.Logf("2x3 matrix as built from rows: %v", m1)
	for r := 0; r < 2; r++ {
		for c := 0; c < 3; c++ {
			if !FloatEqualThreshold(m1.At(r, c), rows[r][c], 1e-5) {
				t.Errorf("Matrix element at (%d,%d) wrong when built from rows. Got: %f, Expected: %f", r, c, m1.At(r, c), rows[r][c])
			}
		}
	}

	row0, row1 = m1.Rows()
	rows2 := [2]Vec3{row0, row1}

	t.Logf("2x3 matrix returned rows: %v", rows2)
	for r := 0; r < 2; r++ {
		for c := 0; c < 3; c++ {
			if !FloatEqualThreshold(rows2[r][c], rows[r][c], 1e-5) {
				t.Errorf("Matrix element at (%d,%d) wrong when rows are gotten. Got: %f, Expected: %f", r, c, rows2[r][c], rows[r][c])
			}
		}
	}
}

// M<N
func TestMatRowsWide(t *testing.T) {
	t.Parallel()

	v0 := Vec3{1, 2, 3}
	v1 := Vec3{4, 5, 6}
	v2 := Vec3{7, 8, 9}
	v3 := Vec3{10, 11, 12}

	rows := [4]Vec3{v0, v1, v2, v3}

	m1 := Mat4x3FromRows(v0, v1, v2, v3)

	t.Logf("4x3 matrix as built from rows: %v", m1)
	for r := 0; r < 4; r++ {
		for c := 0; c < 3; c++ {
			if !FloatEqualThreshold(m1.At(r, c), rows[r][c], 1e-5) {
				t.Errorf("Matrix element at (%d,%d) wrong when built from rows. Got: %f, Expected: %f", r, c, m1.At(r, c), rows[r][c])
			}
		}
	}

	v0, v1, v2, v3 = m1.Rows()
	r2 := [4]Vec3{v0, v1, v2, v3}

	t.Logf("4x3 matrix returned rows: %v", r2)
	for r := 0; r < 4; r++ {
		for c := 0; c < 3; c++ {
			if !FloatEqualThreshold(r2[r][c], rows[r][c], 1e-5) {
				t.Errorf("Matrix element at (%d,%d) wrong when rows are gotten. Got: %f, Expected: %f", r, c, r2[r][c], rows[r][c])
			}
		}
	}
}

// Square matrix
func TestMatRowsSquare(t *testing.T) {
	t.Parallel()

	v0 := Vec4{1, 2, 3, 4}
	v1 := Vec4{5, 6, 7, 8}
	v2 := Vec4{9, 10, 11, 12}
	v3 := Vec4{13, 14, 15, 16}
	rows := [4]Vec4{v0, v1, v2, v3}
	m1 := Mat4FromRows(v0, v1, v2, v3)

	t.Logf("4x4 matrix as built from rows: %v", m1)
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			if !FloatEqualThreshold(m1.At(r, c), rows[r][c], 1e-5) {
				t.Errorf("Matrix element at (%d,%d) wrong when built from rows. Got: %f, Expected: %f", r, c, m1.At(r, c), rows[r][c])
			}
		}
	}

	v0, v1, v2, v3 = m1.Rows()
	r2 := [4]Vec4{v0, v1, v2, v3}

	t.Logf("4x4 matrix returned rows: %v", r2)
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			if !FloatEqualThreshold(r2[r][c], rows[r][c], 1e-5) {
				t.Errorf("Matrix element at (%d,%d) wrong when rows are gotten. Got: %f, Expected: %f", r, c, r2[r][c], rows[r][c])
			}
		}
	}
}

// M<N
func TestMatColsTall(t *testing.T) {

	v0 := Vec3{1, 2, 3}
	v1 := Vec3{4, 5, 6}
	cols := [2]Vec3{v0, v1}
	m1 := Mat3x2FromCols(v0, v1)

	t.Logf("3x2 matrix as built from cols: %v", m1)
	for r := 0; r < 3; r++ {
		for c := 0; c < 2; c++ {
			if !FloatEqualThreshold(m1.At(r, c), cols[c][r], 1e-5) {
				t.Errorf("Matrix element at (%d,%d) wrong when built from rows. Got: %f, Expected: %f", r, c, m1.At(r, c), cols[c][r])
			}
		}
	}

	v0, v1 = m1.Cols()
	r2 := [2]Vec3{v0, v1}

	t.Logf("3x2 matrix returned cols: %v", r2)
	for r := 0; r < 3; r++ {
		for c := 0; c < 2; c++ {
			if !FloatEqualThreshold(r2[c][r], cols[c][r], 1e-5) {
				t.Errorf("Matrix element at (%d,%d) wrong when rows are gotten. Got: %f, Expected: %f", r, c, r2[c][r], cols[c][r])
			}
		}
	}
}

// M>N
func TestMatColsWide(t *testing.T) {
	t.Parallel()

	v0 := Vec3{1, 2, 3}
	v1 := Vec3{4, 5, 6}
	v2 := Vec3{7, 8, 9}
	v3 := Vec3{10, 11, 12}
	cols := [4]Vec3{v0, v1, v2, v3}
	m1 := Mat3x4FromCols(v0, v1, v2, v3)

	t.Logf("3x4 matrix as built from cols: %v", m1)
	for r := 0; r < 3; r++ {
		for c := 0; c < 4; c++ {
			if !FloatEqualThreshold(m1.At(r, c), cols[c][r], 1e-5) {
				t.Errorf("Matrix element at (%d,%d) wrong when built from cols. Got: %f, Expected: %f", r, c, m1.At(r, c), cols[c][r])
			}
		}
	}

	v0, v1, v2, v3 = m1.Cols()
	r2 := [4]Vec3{v0, v1, v2, v3}

	t.Logf("3x4 matrix returned cols: %v", r2)
	for r := 0; r < 3; r++ {
		for c := 0; c < 4; c++ {
			if !FloatEqualThreshold(r2[c][r], cols[c][r], 1e-5) {
				t.Errorf("Matrix element at (%d,%d) wrong when rows are gotten. Got: %f, Expected: %f", r, c, r2[c][r], cols[c][r])
			}
		}
	}
}

// Square matrix
func TestMatColsSquare(t *testing.T) {
	t.Parallel()

	v0 := Vec4{1, 2, 3, 4}
	v1 := Vec4{5, 6, 7, 8}
	v2 := Vec4{9, 10, 11, 12}
	v3 := Vec4{13, 14, 15, 16}
	cols := [4]Vec4{v0, v1, v2, v3}
	m1 := Mat4FromCols(v0, v1, v2, v3)

	t.Logf("4x4 matrix as built from cols: %v", m1)
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			if !FloatEqualThreshold(m1.At(r, c), cols[c][r], 1e-5) {
				t.Errorf("Matrix element at (%d,%d) wrong when built from rows. Got: %f, Expected: %f", r, c, m1.At(r, c), cols[c][r])
			}
		}
	}

	v0, v1, v2, v3 = m1.Cols()
	r2 := [4]Vec4{v0, v1, v2, v3}

	t.Logf("4x4 matrix returned cols: %v", r2)
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			if !FloatEqualThreshold(r2[c][r], cols[c][r], 1e-5) {
				t.Errorf("Matrix element at (%d,%d) wrong when rows are gotten. Got: %f, Expected: %f", r, c, r2[c][r], cols[c][r])
			}
		}
	}
}

func TestTransposeTall(t *testing.T) {
	t.Parallel()

	m := Mat3x2FromCols(
		Vec3{1, 2, 3},
		Vec3{4, 5, 6},
	)

	transpose := m.Transpose()

	correct := Mat2x3FromRows(
		Vec3{1, 2, 3},
		Vec3{4, 5, 6},
	)

	if !correct.ApproxEqualThreshold(transpose, 1e-4) {
		t.Errorf("Transpose not correct. Got: %v, expected: %v", transpose, correct)
	}
}

func TestTransposeWide(t *testing.T) {
	t.Parallel()

	m := Mat2x3FromCols(
		Vec2{1, 2},
		Vec2{3, 4},
		Vec2{5, 6},
	)

	transpose := m.Transpose()

	correct := Mat3x2FromRows(
		Vec2{1, 2},
		Vec2{3, 4},
		Vec2{5, 6},
	)

	if !correct.ApproxEqualThreshold(transpose, 1e-4) {
		t.Errorf("Transpose not correct. Got: %v, expected: %v", transpose, correct)
	}
}

func TestTransposeSquare(t *testing.T) {
	t.Parallel()

	m := Mat3FromCols(
		Vec3{1, 2, 3},
		Vec3{4, 5, 6},
		Vec3{7, 8, 9},
	)

	transpose := m.Transpose()

	correct := Mat3FromRows(
		Vec3{1, 2, 3},
		Vec3{4, 5, 6},
		Vec3{7, 8, 9},
	)

	if !correct.ApproxEqualThreshold(transpose, 1e-4) {
		t.Errorf("Transpose not correct. Got: %v, expected: %v", transpose, correct)
	}
}

func TestAtSet(t *testing.T) {
	t.Parallel()

	m := Mat3{1, 2, 3, 4, 5, 6, 7, 8, 9}

	v := m.At(0, 2)

	if !FloatEqualThreshold(v, 7, 1e-4) {
		t.Errorf("Incorrect value gotten by At: %v, expected %v", v, 3)
	}

	m.Set(0, 2, 9001)

	v = m.At(0, 2)

	if !FloatEqualThreshold(v, 9001, 1e-4) {
		t.Errorf("Value set by Set not gotten by At: %v, expected %v", v, 9001)
	}

	correctMat := Mat3{1, 2, 3, 4, 5, 6, 9001, 8, 9}

	if !correctMat.ApproxEqualThreshold(m, 1e-4) {
		t.Errorf("After set, not equal to matrix that should be identical. Got: %v, expected: %v", m, correctMat)
	}
}

func TestDiagTrace(t *testing.T) {
	t.Parallel()

	m := Diag4(Vec4{1, 2, 3, 4})

	tr := m.Trace()

	if !FloatEqualThreshold(tr, 10, 1e-4) {
		t.Errorf("Trace of matrix seeded with diagonal vector {1,2,3,4} not equal to 10. Got %v", tr)
	}
}

func TestMatAbs(t *testing.T) {
	t.Parallel()

	m := Mat3{1, -3, 4, 5, -6, 8, -9, 10, 0}
	result := Mat3{1, 3, 4, 5, 6, 8, 9, 10, 0}

	m = m.Abs()

	if !result.ApproxEqualThreshold(m, 1e-6) {
		t.Errorf("Matrix absolute value does not work properly. Got: %v, Expected: %v", m, result)
	}
}

func TestString(t *testing.T) {
	m := Ident4()

	str := fmt.Sprintf(` %[2]f %[1]f %[1]f %[1]f
 %[1]f %[2]f %[1]f %[1]f
 %[1]f %[1]f %[2]f %[1]f
 %[1]f %[1]f %[1]f %[2]f
`, 0.0, 1.0)

	if str != m.String() {
		t.Errorf("Mat string conversion not working got %q expected %q", m.String(), str)
	}
}

func BenchmarkMatAdd(b *testing.B) {
	b.StopTimer()
	rand := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		m1 := Mat4{}
		m2 := Mat4{}

		for j := 0; j < len(m1); j++ {
			m1[j], m2[j] = rand.Float64(), rand.Float64()
		}
		b.StartTimer()

		m1 = m1.Add(m2)
	}
}

func BenchmarkMatScale(b *testing.B) {
	b.StopTimer()
	rand := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		m1 := Mat4{}

		for j := 0; j < len(m1); j++ {
			m1[j] = rand.Float64()
		}
		c := rand.Float64()
		b.StartTimer()

		m1 = m1.Mul(c)
	}
}

func BenchmarkMatMul(b *testing.B) {
	b.StopTimer()
	rand := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		m1 := Mat4{}
		m2 := Mat4{}

		for j := 0; j < len(m1); j++ {
			m1[j], m2[j] = rand.Float64(), rand.Float64()
		}
		b.StartTimer()

		m1 = m1.Mul4(m2)
	}
}

func BenchmarkMatTranspose(b *testing.B) {
	b.StopTimer()
	rand := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		m1 := Mat4{}

		for j := 0; j < len(m1); j++ {
			m1[j] = rand.Float64()
		}
		b.StartTimer()

		_ = m1.Transpose()
	}
}

func BenchmarkMatDet(b *testing.B) {
	b.StopTimer()
	rand := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		m1 := Mat4{}

		for j := 0; j < len(m1); j++ {
			m1[j] = rand.Float64()
		}
		b.StartTimer()

		_ = m1.Det()
	}
}

func BenchmarkMatInv(b *testing.B) {
	b.StopTimer()
	rand := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		m1 := Mat4{}

		for j := 0; j < len(m1); j++ {
			m1[j] = rand.Float64()
		}
		b.StartTimer()

		m1 = m1.Inv()
	}
}
