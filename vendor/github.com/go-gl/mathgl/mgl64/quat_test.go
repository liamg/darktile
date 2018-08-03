// This file is generated from mgl32/quat_test.go; DO NOT EDIT

// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mgl64

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestQuatMulIdentity(t *testing.T) {
	t.Parallel()

	i1 := Quat{1.0, Vec3{0, 0, 0}}
	i2 := QuatIdent()
	i3 := QuatIdent()

	mul := i2.Mul(i3)

	if !FloatEqual(mul.W, 1.0) {
		t.Errorf("Multiplication of identities does not yield identity")
	}

	for i := range mul.V {
		if mul.V[i] != i1.V[i] {
			t.Errorf("Multiplication of identities does not yield identity")
		}
	}
}

func TestQuatRotateOnAxis(t *testing.T) {
	t.Parallel()

	var angleDegrees float64 = 30.0
	axis := Vec3{1, 0, 0}

	i1 := QuatRotate(DegToRad(angleDegrees), axis)

	rotatedAxis := i1.Rotate(axis)

	for i := range rotatedAxis {
		if !FloatEqualThreshold(rotatedAxis[i], axis[i], 1e-4) {
			t.Errorf("Rotation of axis does not yield identity")
		}
	}
}

func TestQuatRotateOffAxis(t *testing.T) {
	t.Parallel()

	var angleRads float64 = DegToRad(30.0)
	axis := Vec3{1, 0, 0}

	i1 := QuatRotate(angleRads, axis)

	vector := Vec3{0, 1, 0}
	rotatedVector := i1.Rotate(vector)

	s, c := math.Sincos(float64(angleRads))
	answer := Vec3{0, float64(c), float64(s)}

	for i := range rotatedVector {
		if !FloatEqualThreshold(rotatedVector[i], answer[i], 1e-4) {
			t.Errorf("Rotation of vector does not yield answer")
		}
	}
}

func TestQuatIdentityToMatrix(t *testing.T) {
	t.Parallel()

	quat := QuatIdent()
	matrix := quat.Mat4()
	answer := Ident4()

	if !matrix.ApproxEqual(answer) {
		t.Errorf("Identity quaternion does not yield identity matrix")
	}
}

func TestQuatRotationToMatrix(t *testing.T) {
	t.Parallel()

	var angle float64 = DegToRad(45.0)

	axis := Vec3{1, 2, 3}.Normalize()
	quat := QuatRotate(angle, axis)
	matrix := quat.Mat4()
	answer := HomogRotate3D(angle, axis)

	if !matrix.ApproxEqualThreshold(answer, 1e-4) {
		t.Errorf("Rotation quaternion does not yield correct rotation matrix; got: %v expected: %v", matrix, answer)
	}
}

// Taken from the Matlab AnglesToQuat documentation example
func TestAnglesToQuatZYX(t *testing.T) {
	t.Parallel()

	q := AnglesToQuat(.7854, 0.1, 0, ZYX)

	t.Log("Calculated quaternion: ", q, "\n")

	if !FloatEqualThreshold(q.W, .9227, 1e-3) {
		t.Errorf("Quaternion W incorrect. Got: %f Expected: %f", q.W, .9227)
	}

	if !q.V.ApproxEqualThreshold(Vec3{-0.0191, 0.0462, 0.3822}, 1e-3) {
		t.Errorf("Quaternion V incorrect. Got: %v, Expected: %v", q.V, Vec3{-0.0191, 0.0462, 0.3822})
	}
}

func TestQuatMatRotateY(t *testing.T) {
	t.Parallel()

	q := QuatRotate(float64(math.Pi), Vec3{0, 1, 0})
	q = q.Normalize()
	v := Vec3{1, 0, 0}

	result := q.Rotate(v)

	expected := Rotate3DY(float64(math.Pi)).Mul3x1(v)
	t.Logf("Computed from rotation matrix: %v", expected)
	if !result.ApproxEqualThreshold(expected, 1e-4) {
		t.Errorf("Quaternion rotating vector doesn't match 3D matrix method. Got: %v, Expected: %v", result, expected)
	}

	expected = q.Mul(Quat{0, v}).Mul(q.Conjugate()).V
	t.Logf("Computed from conjugate method: %v", expected)
	if !result.ApproxEqualThreshold(expected, 1e-4) {
		t.Errorf("Quaternion rotating vector doesn't match slower conjugate method. Got: %v, Expected: %v", result, expected)
	}

	expected = Vec3{-1, 0, 0}
	if !result.ApproxEqualThreshold(expected, 4e-4) { // The result we get for z is like 8e-8, but a 1e-4 threshold juuuuuust causes it to freak out when compared to 0.0
		t.Errorf("Quaternion rotating vector doesn't match hand-computed result. Got: %v, Expected: %v", result, expected)
	}
}

func BenchmarkQuatRotateOptimized(b *testing.B) {
	b.StopTimer()
	rand := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		q := QuatRotate(rand.Float64(), Vec3{rand.Float64(), rand.Float64(), rand.Float64()})
		v := Vec3{rand.Float64(), rand.Float64(), rand.Float64()}
		q = q.Normalize()
		b.StartTimer()

		v = q.Rotate(v)
	}
}

func BenchmarkQuatRotateConjugate(b *testing.B) {
	b.StopTimer()
	rand := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		q := QuatRotate(rand.Float64(), Vec3{rand.Float64(), rand.Float64(), rand.Float64()})
		v := Vec3{rand.Float64(), rand.Float64(), rand.Float64()}
		q = q.Normalize()
		b.StartTimer()

		v = q.Mul(Quat{0, v}).Mul(q.Conjugate()).V
	}
}

func BenchmarkQuatArrayAccess(b *testing.B) {
	b.StopTimer()
	rand := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		q := QuatRotate(rand.Float64(), Vec3{rand.Float64(), rand.Float64(), rand.Float64()})
		b.StartTimer()

		_ = q.V[0]
	}
}

func BenchmarkQuatFuncElementAccess(b *testing.B) {
	b.StopTimer()
	rand := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		q := QuatRotate(rand.Float64(), Vec3{rand.Float64(), rand.Float64(), rand.Float64()})
		b.StartTimer()

		_ = q.X()
	}
}

func TestMat4ToQuat(t *testing.T) {
	// http://www.euclideanspace.com/maths/geometry/rotations/conversions/matrixToQuaternion/examples/index.htm

	tests := []struct {
		Description string
		Rotation    Mat4
		Expected    Quat
	}{
		{
			"forward",
			Ident4(),
			QuatIdent(),
		},
		{
			"heading 90 degree",
			Mat4{
				0, 0, -1, 0,
				0, 1, 0, 0,
				1, 0, 0, 0,
				0, 0, 0, 1,
			},
			Quat{0.7071, Vec3{0, 0.7071, 0}},
		},
		{
			"heading 180 degree",
			Mat4{
				-1, 0, 0, 0,
				0, 1, 0, 0,
				0, 0, -1, 0,
				0, 0, 0, 1,
			},
			Quat{0, Vec3{0, 1, 0}},
		},
		{
			"attitude 90 degree",
			Mat4{
				0, 1, 0, 0,
				-1, 0, 0, 0,
				0, 0, 1, 0,
				0, 0, 0, 1,
			},
			Quat{0.7071, Vec3{0, 0, 0.7071}},
		},
		{
			"bank 90 degree",
			Mat4{
				1, 0, 0, 0,
				0, 0, 1, 0,
				0, -1, 0, 0,
				0, 0, 0, 1,
			},
			Quat{0.7071, Vec3{0.7071, 0, 0}},
		},
	}

	threshold := float64(math.Pow(10, -2))
	for _, c := range tests {
		if r := Mat4ToQuat(c.Rotation); !r.ApproxEqualThreshold(c.Expected, threshold) {
			t.Errorf("%v failed: Mat4ToQuat(%v) != %v (got %v)", c.Description, c.Rotation, c.Expected, r)
		}
	}
}

func TestQuatRotate(t *testing.T) {
	tests := []struct {
		Description string
		Angle       float64
		Axis        Vec3
		Expected    Quat
	}{
		{
			"forward",
			0, Vec3{0, 0, 0},
			QuatIdent(),
		},
		{
			"heading 90 degree",
			DegToRad(90), Vec3{0, 1, 0},
			Quat{0.7071, Vec3{0, 0.7071, 0}},
		},
		{
			"heading 180 degree",
			DegToRad(180), Vec3{0, 1, 0},
			Quat{0, Vec3{0, 1, 0}},
		},
		{
			"attitude 90 degree",
			DegToRad(90), Vec3{0, 0, 1},
			Quat{0.7071, Vec3{0, 0, 0.7071}},
		},
		{
			"bank 90 degree",
			DegToRad(90), Vec3{1, 0, 0},
			Quat{0.7071, Vec3{0.7071, 0, 0}},
		},
	}

	threshold := float64(math.Pow(10, -2))
	for _, c := range tests {
		if r := QuatRotate(c.Angle, c.Axis); !r.OrientationEqualThreshold(c.Expected, threshold) {
			t.Errorf("%v failed: QuatRotate(%v, %v) != %v (got %v)", c.Description, c.Angle, c.Axis, c.Expected, r)
		}
	}
}

func TestQuatLookAtV(t *testing.T) {
	// http://www.euclideanspace.com/maths/algebra/realNormedAlgebra/quaternions/transforms/examples/index.htm

	tests := []struct {
		Description     string
		Eye, Center, Up Vec3
		Expected        Quat
	}{
		{
			"forward",
			Vec3{0, 0, 0},
			Vec3{0, 0, -1},
			Vec3{0, 1, 0},
			QuatIdent(),
		},
		{
			"heading 90 degree",
			Vec3{0, 0, 0},
			Vec3{1, 0, 0},
			Vec3{0, 1, 0},
			Quat{0.7071, Vec3{0, 0.7071, 0}},
		},
		{
			"heading 180 degree",
			Vec3{0, 0, 0},
			Vec3{0, 0, 1},
			Vec3{0, 1, 0},
			Quat{0, Vec3{0, 1, 0}},
		},
		{
			"attitude 90 degree",
			Vec3{0, 0, 0},
			Vec3{0, 0, -1},
			Vec3{1, 0, 0},
			Quat{0.7071, Vec3{0, 0, 0.7071}},
		},
		{
			"bank 90 degree",
			Vec3{0, 0, 0},
			Vec3{0, -1, 0},
			Vec3{0, 0, -1},
			Quat{0.7071, Vec3{0.7071, 0, 0}},
		},
	}

	threshold := float64(math.Pow(10, -2))
	for _, c := range tests {
		if r := QuatLookAtV(c.Eye, c.Center, c.Up); !r.OrientationEqualThreshold(c.Expected, threshold) {
			t.Errorf("%v failed: QuatLookAtV(%v, %v, %v) != %v (got %v)", c.Description, c.Eye, c.Center, c.Up, c.Expected, r)
		}
	}
}

func TestCompareLookAt(t *testing.T) {
	type OrigExp [2]Vec3

	tests := []struct {
		Description     string
		Eye, Center, Up Vec3
		Pos             []OrigExp
	}{
		{
			"forward, identity rotation",
			// looking from viewer into screen z-, up y+
			Vec3{0, 0, 0}, Vec3{0, 0, -1}, Vec3{0, 1, 0},
			[]OrigExp{
				{Vec3{1, 2, 3}, Vec3{1, 2, 3}},
			},
		},
		{
			"heading -90 degree, look right",
			// look x+
			// rotate around y -90 deg
			Vec3{0, 0, 0}, Vec3{1, 0, 0}, Vec3{0, 1, 0},
			[]OrigExp{
				{Vec3{1, 2, 3}, Vec3{3, 2, -1}},

				{Vec3{1, 1, -1}, Vec3{-1, 1, -1}},
				{Vec3{1, 1, 1}, Vec3{1, 1, -1}},
				{Vec3{1, -1, 1}, Vec3{1, -1, -1}},
				{Vec3{1, -1, -1}, Vec3{-1, -1, -1}},

				{Vec3{-1, 1, -1}, Vec3{-1, 1, 1}},
				{Vec3{-1, 1, 1}, Vec3{1, 1, 1}},
				{Vec3{-1, -1, 1}, Vec3{1, -1, 1}},
				{Vec3{-1, -1, -1}, Vec3{-1, -1, 1}},
			},
		},
		{
			"heading 180 degree",
			Vec3{0, 0, 0}, Vec3{0, 0, 1}, Vec3{0, 1, 0},
			[]OrigExp{
				{Vec3{1, 2, 3}, Vec3{-1, 2, -3}},
			},
		},
		{
			"attitude 90 degree",
			Vec3{0, 0, 0}, Vec3{0, 0, -1}, Vec3{1, 0, 0},
			[]OrigExp{
				{Vec3{1, 2, 3}, Vec3{-2, 1, 3}},
			},
		},
		{
			"bank 90 degree, look down",
			// look y-
			// rotate around x -90 deg
			// up toward z-
			Vec3{0, 0, 0}, Vec3{0, -1, 0}, Vec3{0, 0, -1},
			[]OrigExp{
				{Vec3{1, 2, 3}, Vec3{1, -3, 2}},

				{Vec3{1, 1, -1}, Vec3{1, 1, 1}},
				{Vec3{1, 1, 1}, Vec3{1, -1, 1}},
				{Vec3{1, -1, 1}, Vec3{1, -1, -1}},
				{Vec3{1, -1, -1}, Vec3{1, 1, -1}},

				{Vec3{-1, 1, -1}, Vec3{-1, 1, 1}},
				{Vec3{-1, 1, 1}, Vec3{-1, -1, 1}},
				{Vec3{-1, -1, 1}, Vec3{-1, -1, -1}},
				{Vec3{-1, -1, -1}, Vec3{-1, 1, -1}},
			},
		},
		{
			"half roll",
			// immelmann turn without the half roll
			// looking from screen to viewer z+
			// upside down, y-
			Vec3{0, 0, 0}, Vec3{0, 0, 1}, Vec3{0, -1, 0},
			[]OrigExp{
				{Vec3{1, 1, -1}, Vec3{1, -1, 1}},
				{Vec3{1, 1, 1}, Vec3{1, -1, -1}},
				{Vec3{1, -1, 1}, Vec3{1, 1, -1}},
				{Vec3{1, -1, -1}, Vec3{1, 1, 1}},

				{Vec3{-1, 1, -1}, Vec3{-1, -1, 1}},
				{Vec3{-1, 1, 1}, Vec3{-1, -1, -1}},
				{Vec3{-1, -1, 1}, Vec3{-1, 1, -1}},
				{Vec3{-1, -1, -1}, Vec3{-1, 1, 1}},
			},
		},
		{
			"roll left",
			// look x-
			// rotate around y 90 deg
			// up toward viewer z+
			Vec3{0, 0, 0}, Vec3{-1, 0, 0}, Vec3{0, 0, 1},
			[]OrigExp{
				{Vec3{1, 1, -1}, Vec3{1, -1, 1}},
				{Vec3{1, 1, 1}, Vec3{1, 1, 1}},
				{Vec3{1, -1, 1}, Vec3{-1, 1, 1}},
				{Vec3{1, -1, -1}, Vec3{-1, -1, 1}},

				{Vec3{-1, 1, -1}, Vec3{1, -1, -1}},
				{Vec3{-1, 1, 1}, Vec3{1, 1, -1}},
				{Vec3{-1, -1, 1}, Vec3{-1, 1, -1}},
				{Vec3{-1, -1, -1}, Vec3{-1, -1, -1}},
			},
		},
	}

	threshold := float64(math.Pow(10, -2))
	for _, c := range tests {
		m := LookAtV(c.Eye, c.Center, c.Up)
		q := QuatLookAtV(c.Eye, c.Center, c.Up)

		for i, p := range c.Pos {
			t.Log(c.Description, i)
			o, e := p[0], p[1]
			rm := m.Mul4x1(o.Vec4(0)).Vec3()
			rq := q.Rotate(o)

			if !rq.ApproxEqualThreshold(rm, threshold) {
				t.Errorf("%v failed: QuatLookAtV() != LookAtV()", c.Description)
			}

			if !e.ApproxEqualThreshold(rm, threshold) {
				t.Errorf("%v failed: (%v).Mul4x1(%v) != %v (got %v)", c.Description, m, o, e, rm)
			}

			if !e.ApproxEqualThreshold(rq, threshold) {
				t.Errorf("%v failed: (%v).Rotate(%v) != %v (got %v)", c.Description, q, o, e, rq)
			}
		}
	}
}

func TestQuatMatConversion(t *testing.T) {
	tests := []struct {
		Angle float64
		Axis  Vec3
	}{}

	for a := 0.0; a <= math.Pi*2; a += math.Pi / 4.0 {
		af := float64(a)
		tests = append(tests, []struct {
			Angle float64
			Axis  Vec3
		}{
			{af, Vec3{1, 0, 0}},
			{af, Vec3{0, 1, 0}},
			{af, Vec3{0, 0, 1}},
		}...)
	}

	for _, c := range tests {
		m1 := HomogRotate3D(c.Angle, c.Axis)
		q1 := Mat4ToQuat(m1)
		q2 := QuatRotate(c.Angle, c.Axis)

		if !FloatEqualThreshold(Abs(q1.Dot(q2)), 1, 1e-4) {
			t.Errorf("Quaternions for %v %v do not match:\n%v\n%v", RadToDeg(c.Angle), c.Axis, q1, q2)
		}
	}
}

func TestQuatGetter(t *testing.T) {
	tests := []Quat{
		{0, Vec3{0, 0, 0}},
		{1, Vec3{2, 3, 4}},
		{-4, Vec3{-3, -2, -1}},
	}

	for _, q := range tests {
		if r := q.X(); !FloatEqualThreshold(r, q.V[0], 1e-4) {
			t.Errorf("Quat(%v).X() != %v (got %v)", q, q.V[0], r)
		}

		if r := q.Y(); !FloatEqualThreshold(r, q.V[1], 1e-4) {
			t.Errorf("Quat(%v).Y() != %v (got %v)", q, q.V[1], r)
		}

		if r := q.Z(); !FloatEqualThreshold(r, q.V[2], 1e-4) {
			t.Errorf("Quat(%v).Z() != %v (got %v)", q, q.V[2], r)
		}
	}
}

func TestQuatEqual(t *testing.T) {
	tests := []struct {
		A, B     Quat
		Expected bool
	}{
		{Quat{1, Vec3{0, 0, 0}}, Quat{1, Vec3{0, 0, 0}}, true},
		{Quat{1, Vec3{2, 3, 4}}, Quat{1, Vec3{2, 3, 4}}, true},
		{Quat{0.0000000000001, Vec3{0, 0, 0}}, Quat{0, Vec3{0, 0, 0}}, true},
		{Quat{MaxValue, Vec3{1, 0, 0}}, Quat{MaxValue, Vec3{1, 0, 0}}, true},
		{Quat{0, Vec3{0, 1, 0}}, Quat{1, Vec3{0, 0, 0}}, false},
		{Quat{1, Vec3{2, 3, 0}}, Quat{-4, Vec3{5, 6, 0}}, false},
	}

	for _, c := range tests {
		if r := c.A.ApproxEqualThreshold(c.B, 1e-4); r != c.Expected {
			t.Errorf("Quat(%v).ApproxEqualThreshold(Quat(%v), 1e-4) != %v (got %v)", c.A, c.B, c.Expected, r)
		}
	}
}

func TestQuatOrientationEqual(t *testing.T) {
	tests := []struct {
		A, B     Quat
		Expected bool
	}{
		{Quat{1, Vec3{0, 0, 0}}, Quat{1, Vec3{0, 0, 0}}, true},
		{Quat{0, Vec3{0, 1, 0}}, Quat{0, Vec3{0, -1, 0}}, true},
		{Quat{0, Vec3{0, 1, 0}}, Quat{1, Vec3{0, 0, 0}}, false},
		{Quat{1, Vec3{2, 3, 0}}, Quat{-4, Vec3{5, 6, 0}}, false},
	}

	for _, c := range tests {
		if r := c.A.OrientationEqualThreshold(c.B, 1e-4); r != c.Expected {
			t.Errorf("Quat(%v).OrientationEqualThreshold(Quat(%v), 1e-4) != %v (got %v)", c.A, c.B, c.Expected, r)
		}
	}
}

func TestQuatAdd(t *testing.T) {
	tests := []struct {
		A, B     Quat
		Expected Quat
	}{
		{Quat{0, Vec3{0, 0, 0}}, Quat{0, Vec3{0, 0, 0}}, Quat{0, Vec3{0, 0, 0}}},
		{Quat{1, Vec3{0, 0, 0}}, Quat{1, Vec3{0, 0, 0}}, Quat{2, Vec3{0, 0, 0}}},
		{Quat{1, Vec3{2, 3, 4}}, Quat{5, Vec3{6, 7, 8}}, Quat{6, Vec3{8, 10, 12}}},
	}

	for _, c := range tests {
		if r := c.A.Add(c.B); !r.ApproxEqualThreshold(c.Expected, 1e-4) {
			t.Errorf("Quat(%v).Add(Quat(%v)) != %v (got %v)", c.A, c.B, c.Expected, r)
		}
	}
}

func TestQuatSub(t *testing.T) {
	tests := []struct {
		A, B     Quat
		Expected Quat
	}{
		{Quat{0, Vec3{0, 0, 0}}, Quat{0, Vec3{0, 0, 0}}, Quat{0, Vec3{0, 0, 0}}},
		{Quat{1, Vec3{0, 0, 0}}, Quat{1, Vec3{0, 0, 0}}, Quat{0, Vec3{0, 0, 0}}},
		{Quat{1, Vec3{2, 3, 4}}, Quat{5, Vec3{6, 7, 8}}, Quat{-4, Vec3{-4, -4, -4}}},
	}

	for _, c := range tests {
		if r := c.A.Sub(c.B); !r.ApproxEqualThreshold(c.Expected, 1e-4) {
			t.Errorf("Quat(%v).Sub(Quat(%v)) != %v (got %v)", c.A, c.B, c.Expected, r)
		}
	}
}

func TestQuatScale(t *testing.T) {
	tests := []struct {
		Rotation Quat
		Scalar   float64
		Expected Quat
	}{
		{Quat{0, Vec3{0, 0, 0}}, 1, Quat{0, Vec3{0, 0, 0}}},
		{Quat{1, Vec3{0, 0, 0}}, 2, Quat{2, Vec3{0, 0, 0}}},
		{Quat{1, Vec3{2, 3, 4}}, 3, Quat{3, Vec3{6, 9, 12}}},
	}

	for _, c := range tests {
		if r := c.Rotation.Scale(c.Scalar); !r.ApproxEqualThreshold(c.Expected, 1e-4) {
			t.Errorf("Quat(%v).Scale(%v) != %v (got %v)", c.Rotation, c.Scalar, c.Expected, r)
		}
	}
}

func TestQuatLen(t *testing.T) {
	tests := []struct {
		Rotation Quat
		Expected float64
	}{
		{Quat{0, Vec3{1, 0, 0}}, 1},
		{Quat{0, Vec3{0.0000000000001, 0, 0}}, 0},
		{Quat{0, Vec3{MaxValue, 1, 0}}, InfPos},
		{Quat{4, Vec3{1, 2, 3}}, float64(math.Sqrt(1*1 + 2*2 + 3*3 + 4*4))},
		{Quat{0, Vec3{3.1, 4.2, 1.3}}, float64(math.Sqrt(3.1*3.1 + 4.2*4.2 + 1.3*1.3))},
	}

	for _, c := range tests {
		if r := c.Rotation.Len(); !FloatEqualThreshold(c.Expected, r, 1e-4) {
			t.Errorf("Quat(%v).Len() != %v (got %v)", c.Rotation, c.Expected, r)
		}

		if !FloatEqualThreshold(c.Rotation.Len(), c.Rotation.Norm(), 1e-4) {
			t.Error("Quat().Len() != Quat().Norm()")
		}
	}
}

func TestQuatNormalize(t *testing.T) {
	tests := []struct {
		Rotation Quat
		Expected Quat
	}{
		{Quat{0, Vec3{0, 0, 0}}, Quat{1, Vec3{0, 0, 0}}},
		{Quat{0, Vec3{1, 0, 0}}, Quat{0, Vec3{1, 0, 0}}},
		{Quat{0, Vec3{0.0000000000001, 0, 0}}, Quat{0, Vec3{1, 0, 0}}},
		{Quat{0, Vec3{MaxValue, 1, 0}}, Quat{0, Vec3{1, 0, 0}}},
		{Quat{4, Vec3{1, 2, 3}}, Quat{4.0 / 5.477, Vec3{1.0 / 5.477, 2.0 / 5.477, 3.0 / 5.477}}},
		{Quat{0, Vec3{3.1, 4.2, 1.3}}, Quat{0, Vec3{3.1 / 5.3795, 4.2 / 5.3795, 1.3 / 5.3795}}},
	}

	for _, c := range tests {
		if r := c.Rotation.Normalize(); !r.ApproxEqualThreshold(c.Expected, 1e-4) {
			t.Errorf("Quat(%v).Normalize() != %v (got %v)", c.Rotation, c.Expected, r)
		}
	}
}

func TestQuatInverse(t *testing.T) {
	tests := []struct {
		Rotation Quat
		Expected Quat
	}{
		{Quat{0, Vec3{1, 0, 0}}, Quat{0, Vec3{-1, 0, 0}}},
		{Quat{3, Vec3{-1, 4, 3}}, Quat{3.0 / 35.0, Vec3{1.0 / 35.0, -4.0 / 35.0, -3.0 / 35.0}}},
		{Quat{1, Vec3{0, 0, 2}}, Quat{1.0 / 5.0, Vec3{0, 0, -2.0 / 5.0}}},
	}

	for _, c := range tests {
		if r := c.Rotation.Inverse(); !r.ApproxEqualThreshold(c.Expected, 1e-4) {
			t.Errorf("Quat(%v).Inverse() != %v (got %v)", c.Rotation, c.Expected, r)
		}
	}
}

func TestQuatSlerp(t *testing.T) {
	tests := []struct {
		A, B     Quat
		Scalar   float64
		Expected Quat
	}{
		{Quat{0, Vec3{0, 0, 0}}, Quat{0, Vec3{0, 0, 0}}, 0, Quat{1, Vec3{0, 0, 0}}},
		{Quat{0, Vec3{1, 0, 0}}, Quat{0, Vec3{1, 0, 0}}, 0.5, Quat{0, Vec3{1, 0, 0}}},
		{Quat{1, Vec3{0, 0, 0}}, Quat{0, Vec3{1, 0, 0}}, 0.5, Quat{0.7071067811865475, Vec3{0.7071067811865475, 0, 0}}},
		{Quat{0.5, Vec3{-0.5, -0.5, 0.5}}, Quat{0.996, Vec3{-0.080, -0.080, 0}}, 1, Quat{0.996, Vec3{-0.080, -0.080, 0}}},
		{Quat{0.5, Vec3{-0.5, -0.5, 0.5}}, Quat{0.996, Vec3{-0.080, -0.080, 0}}, 0, Quat{0.5, Vec3{-0.5, -0.5, 0.5}}},
		{Quat{0.5, Vec3{-0.5, -0.5, 0.5}}, Quat{0.996, Vec3{-0.080, -0.080, 0}}, 0.2, Quat{0.6553097459373098, Vec3{-0.44231939784548874, -0.44231939784548874, 0.4237176207195655}}},
		{Quat{0.996, Vec3{-0.080, -0.080, 0}}, Quat{0.5, Vec3{-0.5, -0.5, 0.5}}, 0.8, Quat{0.6553097459373098, Vec3{-0.44231939784548874, -0.44231939784548874, 0.4237176207195655}}},
		{Quat{1, Vec3{0, 0, 0}}, Quat{-0.9999999, Vec3{0, 0, 0}}, 0, Quat{1, Vec3{0, 0, 0}}},
	}

	for _, c := range tests {
		if r := QuatSlerp(c.A, c.B, c.Scalar); !r.ApproxEqualThreshold(c.Expected, 1e-2) {
			t.Errorf("QuatSlerp(%v, %v, %v) != %v (got %v)", c.A, c.B, c.Scalar, c.Expected, r)
		}
	}
}

func TestQuatDot(t *testing.T) {
	tests := []struct {
		A, B     Quat
		Expected float64
	}{
		{Quat{0, Vec3{0, 0, 0}}, Quat{0, Vec3{0, 0, 0}}, 0},
		{Quat{0, Vec3{1, 2, 3}}, Quat{0, Vec3{4, 5, 6}}, 32},
		{Quat{4, Vec3{1, 2, 3}}, Quat{8, Vec3{5, 6, 7}}, 70},
	}

	for _, c := range tests {
		if r := c.A.Dot(c.B); !FloatEqualThreshold(r, c.Expected, 1e-4) {
			t.Errorf("Quat(%v).Dot(Quat(%v)) != %v (got %v)", c.A, c.B, c.Expected, r)
		}
	}
}
