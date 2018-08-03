// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mgl32

import (
	"math/rand"
	"testing"
	"time"
)

/* Only floats are tested because the double versions are simply a find->replace on floats */

func Test2DVecAdd(t *testing.T) {
	v1 := Vec2{1.0, 2.5}
	v2 := Vec2{0.0, 1.0}

	v3 := v1.Add(v2)

	if !FloatEqual(v3[0], 1.0) || !FloatEqual(v3[1], 3.5) {
		t.Errorf("Add not adding properly")
	}

	v4 := v2.Add(v1)

	if !FloatEqual(v3[0], v4[0]) || !FloatEqual(v3[1], v4[1]) {
		t.Errorf("Addition is somehow not commutative")
	}

}

func Test3DVecAdd(t *testing.T) {
	v1 := Vec3{1.0, 2.5, 1.1}
	v2 := Vec3{0.0, 1.0, 9.9}

	v3 := v1.Add(v2)

	if !FloatEqual(v3[0], 1.0) || !FloatEqual(v3[1], 3.5) || !FloatEqual(v3[2], 11.0) {
		t.Errorf("Add not adding properly")
	}

	v4 := v2.Add(v1)

	if !FloatEqual(v3[0], v4[0]) || !FloatEqual(v3[1], v4[1]) || !FloatEqual(v3[2], v4[2]) {
		t.Errorf("Addition is somehow not commutative")
	}

}

func Test4DVecAdd(t *testing.T) {
	v1 := Vec4{1.0, 2.5, 1.1, 2.0}
	v2 := Vec4{0.0, 1.0, 9.9, 100.0}

	v3 := v1.Add(v2)

	if !FloatEqual(v3[0], 1.0) || !FloatEqual(v3[1], 3.5) || !FloatEqual(v3[2], 11.0) || !FloatEqual(v3[3], 102.0) {
		t.Errorf("Add not adding properly")
	}

	v4 := v2.Add(v1)

	if !FloatEqual(v3[0], v4[0]) || !FloatEqual(v3[1], v4[1]) || !FloatEqual(v3[2], v4[2]) || !FloatEqual(v3[3], v4[3]) {
		t.Errorf("Addition is somehow not commutative")
	}

}

func Test2DVecSub(t *testing.T) {
	v1 := Vec2{1.0, 2.5}
	v2 := Vec2{0.0, 1.0}

	v3 := v1.Sub(v2)

	if !FloatEqual(v3[0], 1.0) || !FloatEqual(v3[1], 1.5) {
		t.Errorf("Sub not subtracting properly [%f, %f]", v3[0], v3[1])
	}

}

func Test3DVecSub(t *testing.T) {
	v1 := Vec3{1.0, 2.5, 1.1}
	v2 := Vec3{0.0, 1.0, 9.9}

	v3 := v1.Sub(v2)

	// 1.1-9.9 does stupid things to floats, so we need threshold
	if !FloatEqual(v3[0], 1.0) || !FloatEqual(v3[1], 1.5) || !FloatEqualThreshold(v3[2], -8.8, 1e-5) {
		t.Errorf("Sub not subtracting properly [%f, %f, %f]", v3[0], v3[1], v3[2])
	}

}

func Test4DVecSub(t *testing.T) {
	v1 := Vec4{1.0, 2.5, 1.1, 2.0}
	v2 := Vec4{0.0, 1.0, 9.9, 100.0}

	v3 := v1.Sub(v2)

	// 1.1-9.9 does stupid things to floats, so we need a more tolerant threshold
	if !FloatEqual(v3[0], 1.0) || !FloatEqual(v3[1], 1.5) || !FloatEqualThreshold(v3[2], -8.8, 1e-5) || !FloatEqual(v3[3], -98.0) {
		t.Errorf("Sub not subtracting properly [%f, %f, %f, %f]", v3[0], v3[1], v3[2], v3[3])
	}

}

func TestVecMul(t *testing.T) {
	v := Vec2{1.0, 0.0}
	v = v.Mul(15.0)

	if !FloatEqual(v[0], 15.0) || !FloatEqual(v[1], 0.0) {
		t.Errorf("Vec mul does something weird [%f, %f]", v[0], v[1])
	}

	v2 := Vec3{1.0, 0.0, 100.1}
	v2 = v2.Mul(15.0)

	if !FloatEqual(v2[0], 15.0) || !FloatEqual(v2[1], 0.0) || !FloatEqual(v2[2], 1501.5) {
		t.Errorf("Vec mul does something weird [%f, %f, %f]", v2[0], v2[1], v2[2])
	}

	v3 := Vec4{1.0, 0.0, 100.1, -1.0}
	v3 = v3.Mul(15.0)

	if !FloatEqual(v3[0], 15.0) || !FloatEqual(v3[1], 0.0) || !FloatEqual(v3[2], 1501.5) || !FloatEqual(v3[3], -15.0) {
		t.Errorf("Vec mul does something weird [%f, %f, %f, %f]", v3[0], v3[1], v3[2], v3[3])
	}
}

func TestVecOuterProd(t *testing.T) {
	v1 := Vec3{1, 2, 3}
	v2 := Vec2{10, 11}

	correct := Mat3x2FromRows(
		Vec2{10, 11},
		Vec2{20, 22},
		Vec2{30, 33},
	)

	result := v1.OuterProd2(v2)

	if !correct.ApproxEqualThreshold(correct, 1e-4) {
		t.Errorf("Vector outer product isn't working. Got: %v. Expected: %v", result, correct)
	}
}

func TestVecCrossProduct(t *testing.T) {
	v1 := Vec3{1, 2, 3}
	v2 := Vec3{10, 11, 12}
	expected := Vec3{-9, 18, -9}
	result := v1.Cross(v2)

	if !expected.ApproxEqual(result) {
		t.Errorf("Vec3 cross product %v x %v Got: %v. Expected: %v.",
			v1, v2, result, expected)
	}
}

func TestVecDotProduct(t *testing.T) {
	mustEqual := func(result float32, expected float32, name string) {
		if !FloatEqual(result, expected) {
			t.Errorf("%v.Dot(%v) failed. Got: %v. Expected %v.",
				name, name, result, expected)
		}
	}

	mustEqual(Vec2{1, 2}.Dot(Vec2{3, 4}), 11, "Vec2")
	mustEqual(Vec3{-1, -5, -7}.Dot(Vec3{10, 20, 30}), -320, "Vec3")
	mustEqual(Vec4{1, 3, 5, 7}.Dot(Vec4{10, 20, 30, 40}), 500, "Vec3")
}

func TestVecLen(t *testing.T) {
	mustEqual := func(result float32, expected float32, name string) {
		if !FloatEqual(result, expected) {
			t.Errorf("%v failed. Got: %v. Expected %v.",
				name, result, expected)
		}
	}

	mustEqual(Vec2{3, 4}.Len(), 5, "Vec2.Len()")
	mustEqual(Vec3{2, -5, 4}.Len(), 6.708203932499, "Vec3.Len()")
	mustEqual(Vec4{2, 3, 5, 7}.Len(), 9.3273790530888, "Vec4.Len()")
}

func Test2DVecNormalize(t *testing.T) {
	v := Vec2{3, 4}
	norm := v.Normalize()
	expected := Vec2{0.6, 0.8}
	if !norm.ApproxEqual(expected) {
		t.Errorf("%v.Normalize() failed. Got: %v. Expected: %v",
			v, norm, expected)
	}
}

func TestVecElemAccessors(t *testing.T) {
	mustEqual := func(desc string, expected float32, results ...float32) {
		for _, r := range results {
			if !FloatEqual(expected, r) {
				t.Errorf("%v failed. Got: %v. Expected %v...",
					desc, results, expected)
			}
		}
	}

	const x, y, z, w = 1, 2, 3, 4
	v2 := Vec2{x, y}
	v3 := Vec3{x, y, z}
	v4 := Vec4{x, y, z, w}

	mustEqual("Vec.X()", x, v2.X(), v3.X(), v4.X())
	mustEqual("Vec.Y()", y, v2.Y(), v3.Y(), v4.Y())
	mustEqual("Vec.Z()", z, v3.Z(), v4.Z())
	mustEqual("Vec.W()", w, v4.W())

	x2, y2 := v2.Elem()
	x3, y3, z3 := v3.Elem()
	x4, y4, z4, w4 := v4.Elem()
	mustEqual("Vec.Elem() -> x", x, x2, x3, x4)
	mustEqual("Vec.Elem() -> y", y, y2, y3, y4)
	mustEqual("Vec.Elem() -> z", z, z3, z4)
	mustEqual("Vec.Elem() -> w", w, w4)
}

func TestVecEqual(t *testing.T) {
	assert := func(res bool, desc string) {
		if !res {
			t.Errorf("%v failed.", desc)
		}
	}

	v2, errV2 := Vec2{1, 2}, Vec2{1, 0}
	v3, errV3 := Vec3{1, 2, 3}, Vec3{1, 2, 0}
	v4, errV4 := Vec4{1, 2, 3, 4}, Vec4{1, 2, 0, 4}

	assert(v2.ApproxEqual(v2), "Vec2.ApproxEqual")
	assert(!v2.ApproxEqual(errV2), "Vec2.ApproxEqual")

	assert(v2.ApproxEqualThreshold(v2, 0.1), "Vec2.ApproxEqualThreshold")
	assert(!v2.ApproxEqualThreshold(errV2, 0.1), "Vec2.ApproxEqualThreshold")

	assert(v2.ApproxFuncEqual(v2, FloatEqual), "Vec2.ApproxFuncEq")
	assert(!v2.ApproxFuncEqual(errV2, FloatEqual), "Vec2.ApproxFuncEq")

	assert(v3.ApproxEqual(v3), "Vec3.ApproxEqual")
	assert(!v3.ApproxEqual(errV3), "Vec3.ApproxEqual")

	assert(v3.ApproxEqualThreshold(v3, 0.1), "Vec3.ApproxEqualThreshold")
	assert(!v3.ApproxEqualThreshold(errV3, 0.1), "Vec3.ApproxEqualThreshold")

	assert(v3.ApproxFuncEqual(v3, FloatEqual), "Vec3.ApproxFuncEq")
	assert(!v3.ApproxFuncEqual(errV3, FloatEqual), "Vec3.ApproxFuncEq")

	assert(v4.ApproxEqual(v4), "Vec4.ApproxEqual")
	assert(!v4.ApproxEqual(errV4), "Vec4.ApproxEqual")

	assert(v4.ApproxEqualThreshold(v4, 0.1), "Vec4.ApproxEqualThreshold")
	assert(!v4.ApproxEqualThreshold(errV4, 0.1), "Vec4.ApproxEqualThreshold")

	assert(v4.ApproxFuncEqual(v4, FloatEqual), "Vec4.ApproxFuncEq")
	assert(!v4.ApproxFuncEqual(errV4, FloatEqual), "Vec4.ApproxFuncEq")
}

func BenchmarkVec4Add(b *testing.B) {
	b.StopTimer()
	r := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		v1 := Vec4{r.Float32(), r.Float32(), r.Float32(), r.Float32()}
		v2 := Vec4{r.Float32(), r.Float32(), r.Float32(), r.Float32()}
		b.StartTimer()

		v1.Add(v2)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var v1 = Vec4{rand.Float32(), rand.Float32(), rand.Float32(), rand.Float32()}
var v2 = Vec4{rand.Float32(), rand.Float32(), rand.Float32(), rand.Float32()}

func BenchmarkVec4Sub(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v1.Sub(v2)
	}
}

func BenchmarkVec4Scale(b *testing.B) {
	b.StopTimer()
	r := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		v1 := Vec4{r.Float32(), r.Float32(), r.Float32(), r.Float32()}
		c := r.Float32()
		b.StartTimer()

		v1.Mul(c)
	}
}

func BenchmarkVec4Dot(b *testing.B) {
	b.StopTimer()
	r := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		v1 := Vec4{r.Float32(), r.Float32(), r.Float32(), r.Float32()}
		v2 := Vec4{r.Float32(), r.Float32(), r.Float32(), r.Float32()}
		b.StartTimer()

		v1.Dot(v2)
	}
}

func BenchmarkVec4Len(b *testing.B) {
	b.StopTimer()
	r := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		v1 := Vec4{r.Float32(), r.Float32(), r.Float32(), r.Float32()}
		b.StartTimer()

		v1.Len()
	}
}

func BenchmarkVec4Normalize(b *testing.B) {
	b.StopTimer()
	r := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		v1 := Vec4{r.Float32(), r.Float32(), r.Float32(), r.Float32()}
		b.StartTimer()

		v1.Normalize()
	}
}

func BenchmarkVecCross(b *testing.B) {
	b.StopTimer()
	r := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		v1 := Vec3{r.Float32(), r.Float32(), r.Float32()}
		v2 := Vec3{r.Float32(), r.Float32(), r.Float32()}
		b.StartTimer()

		v1.Cross(v2)
	}
}
