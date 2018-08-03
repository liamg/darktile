// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mgl32

import (
	"runtime"
	"testing"
)

func TestMxNTransposeWide(t *testing.T) {
	m := Mat2x3FromCols(
		Vec2{1, 2},
		Vec2{3, 4},
		Vec2{5, 6},
	)

	mn := NewMatrixFromData(m[:], 2, 3)

	transpose := m.Transpose()

	transposeMN := mn.Transpose(nil)

	correct := NewMatrixFromData(transpose[:], 3, 2)

	if !correct.ApproxEqualThreshold(transposeMN, 1e-4) {
		t.Errorf("Transpose gives incorrect result; got: %v, expected: %v", transposeMN, correct)
	}
}

func TestMxNTransposeTall(t *testing.T) {
	m := Mat3x2FromCols(
		Vec3{1, 2, 3},
		Vec3{4, 5, 6},
	)

	mn := NewMatrixFromData(m[:], 3, 2)

	transpose := m.Transpose()

	transposeMN := mn.Transpose(nil)

	correct := NewMatrixFromData(transpose[:], 2, 3)

	if !correct.ApproxEqualThreshold(transposeMN, 1e-4) {
		t.Errorf("Transpose gives incorrect result; got: %v, expected: %v", transposeMN, correct)
	}
}

func TestMxNTransposeSquare(t *testing.T) {
	m := Mat3FromCols(
		Vec3{1, 2, 3},
		Vec3{4, 5, 6},
		Vec3{7, 8, 9},
	)

	mn := NewMatrixFromData(m[:], 3, 3)

	transpose := m.Transpose()

	transposeMN := mn.Transpose(nil)

	correct := NewMatrixFromData(transpose[:], 3, 3)

	if !correct.ApproxEqualThreshold(transposeMN, 1e-4) {
		t.Errorf("Transpose gives incorrect result; got: %v, expected: %v", transposeMN, correct)
	}
}

func TestMxNAtSet(t *testing.T) {
	m := Mat3{1, 2, 3, 4, 5, 6, 7, 8, 9}

	mn := NewMatrixFromData(m[:], 3, 3)

	v := mn.At(0, 2)

	if !FloatEqualThreshold(v, 7, 1e-4) {
		t.Errorf("Incorrect value gotten by At: %v, expected %v", v, 7)
	}

	mn.Set(0, 2, 9001)

	v = mn.At(0, 2)

	if !FloatEqualThreshold(v, 9001, 1e-4) {
		t.Errorf("Incorrect value set by Set: %v, expected %v", v, 9001)
	}

	correct := Mat3{1, 2, 3, 4, 5, 6, 9001, 8, 9}
	correctMN := NewMatrixFromData(correct[:], 3, 3)

	if !correctMN.ApproxEqualThreshold(mn, 1e-4) {
		t.Errorf("Set matrix does not equal correct matrix. Got: %v, expected: %v", mn, correctMN)
	}
}

func TestMxNMulMxN(t *testing.T) {
	m := Ident4()
	r := HomogRotate3DX(DegToRad(45))
	tr := Translate3D(1, 0, 0)
	s := Scale3D(2, 2, 2)

	correct := tr.Mul4(r.Mul4(s.Mul4(m))) // tr*r*s
	correctMN := NewMatrixFromData(correct[:], 4, 4)

	mn := NewMatrixFromData(m[:], 4, 4)
	rmn := NewMatrixFromData(r[:], 4, 4)
	trmn := NewMatrixFromData(tr[:], 4, 4)
	smn := NewMatrixFromData(s[:], 4, 4)

	result := trmn.MulMxN(nil, rmn.MulMxN(nil, smn.MulMxN(nil, mn)))

	if !result.ApproxEqualThreshold(correctMN, 1e-4) {
		t.Errorf("Multiplication of MxN matrix and 4x4 matrix not the same. Got: %v expected: %v", result, correctMN)
	}
}

func TestMxNMulMxNErrorHandling(t *testing.T) {
	mn := NewMatrix(4, 12)
	mn2 := NewMatrix(9, 3)

	result := mn2.MulMxN(nil, mn)

	if result != nil {
		t.Errorf("Nil not returned for bad matrix multiplication, got %v instead", result)
	}
}

func TestMxNMul(t *testing.T) {
	m := Mat3{2, 4, 6, 1, 9, 12, 7, 4, 3}
	mn := NewMatrixFromData(m[:], 3, 3)

	correct := m.Mul(15)
	correctMN := NewMatrixFromData(correct[:], 3, 3)

	result := mn.Mul(nil, 15)

	if !correctMN.ApproxEqualThreshold(result, 1e-4) {
		t.Errorf("Scaling a matrix produces weird results got: %v, expected: %v", result, correct)
	}
}

func TestMxNMulNx1(t *testing.T) {
	m := Ident4()
	r := HomogRotate3DX(DegToRad(45))
	tr := Translate3D(1, 0, 0)
	s := Scale3D(2, 2, 2)

	model := tr.Mul4(r.Mul4(s.Mul4(m)))

	v := Vec4{5, 5, 5, 1}
	correct := model.Mul4x1(v)
	correctn := NewVecNFromData(correct[:])

	modelMN := NewMatrixFromData(model[:], 4, 4)
	vn := NewVecNFromData(v[:])

	result := modelMN.MulNx1(nil, vn)

	if !correctn.ApproxEqualThreshold(result, 1e-4) {
		t.Errorf("Multiplying N-dim vector and MxN matrix produces bad result. Got: %v, expected: %v", result, correct)
	}
}

func TestMxNMulNx1Rectangular(t *testing.T) {

	m := NewMatrixFromData([]float32{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
		1, 2, 3,
		4, 5, 6,
	}, 3, 5)

	v := NewVecNFromData([]float32{
		1, 2, 3, 4, 5,
	})

	expected := Vec3{
		25, 35, 45,
	}

	result := m.MulNx1(nil, v).Vec3()

	if expected != result {
		t.Errorf("Multiplying MxN matrix and N-dim vector produces bad result: %v, expected: %v", result, expected)
	}
}

func TestMxNTrace(t *testing.T) {
	m := DiagN(nil, NewVecNFromData([]float32{1, 2, 3, 4, 5}))

	if !FloatEqualThreshold(m.Trace(), 15, 1e-4) {
		t.Errorf("MatMxN's trace of a diagonal with 1,2,3,4,5 is not 15. Got: %v", m.Trace())
	}
}

func complexOperations() {
	m := NewMatrix(15, 20)
	t := m.Transpose(nil)

	t.MulMxN(m, m).MulMxN(m, t)

	t = t.Transpose(t)
}

func BenchmarkMxNWithPooling(b *testing.B) {
	shouldPool = true
	slicePools = nil

	for n := 0; n < b.N; n++ {
		complexOperations()
	}

	b.StopTimer()
	runtime.GC()
}

func BenchmarkMxNWithoutPooling(b *testing.B) {
	shouldPool = false

	for n := 0; n < b.N; n++ {
		complexOperations()
	}

	b.StopTimer()
	runtime.GC()
}
