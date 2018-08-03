// This file is generated from mgl32/vecn.go; DO NOT EDIT

// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mgl64

import (
	"math"
)

// A vector of N elements backed by a slice
//
// As with MatMxN, this is not for hardcore linear algebra with large dimensions. Use github.com/gonum/matrix
// or something like BLAS/LAPACK for that. This is for corner cases in 3D math where you require
// something a little bigger that 4D, but still relatively small.
//
// This VecN uses several sync.Pool objects as a memory pool. The rule is that for any sized vector, the backing slice
// has CAPACITY (not length) of 2^p where p is Ceil(log_2(N)) -- or in other words, rounding up the base-2
// log of the size of the vector. E.G. a VecN of size 17 will have a backing slice of Cap 32.
type VecN struct {
	vec []float64
}

// Creates a new vector with a backing slice filled with the contents
// of initial. It is NOT backed by initial, but rather a slice with cap
// 2^p where p is Ceil(log_2(len(initial))), with the data from initial copied into
// it.
func NewVecNFromData(initial []float64) *VecN {
	if initial == nil {
		return &VecN{}
	}
	var internal []float64
	if shouldPool {
		internal = grabFromPool(len(initial))
	} else {
		internal = make([]float64, len(initial))
	}
	copy(internal, initial)
	return &VecN{vec: internal}
}

// Creates a new vector with a backing slice of
// 2^p where p = Ceil(log_2(n))
func NewVecN(n int) *VecN {
	if shouldPool {
		return &VecN{vec: grabFromPool(n)}
	} else {
		return &VecN{vec: make([]float64, n)}
	}
}

// Returns the raw slice backing the VecN
//
// This may be sent back to the memory pool at any time
// and you aren't advised to rely on this value
func (vn VecN) Raw() []float64 {
	return vn.vec
}

// Gets the element at index i from the vector.
// This does not bounds check, and will panic if i is
// out of range.
func (vn VecN) Get(i int) float64 {
	return vn.vec[i]
}

func (vn *VecN) Set(i int, val float64) {
	vn.vec[i] = val
}

// Sends the allocated memory through the callback if it exists
func (vn *VecN) destroy() {
	if vn == nil || vn.vec == nil {
		return
	}

	if shouldPool {
		returnToPool(vn.vec)
	}
	vn.vec = nil
}

// Resizes the underlying slice to the desired amount, reallocating or retrieving from the pool
// if necessary. The values after a Resize cannot be expected to be related to the values before a Resize.
//
// If the caller is a nil pointer, this returns a value as if NewVecN(n) had been called,
// otherwise it simply returns the caller.
func (vn *VecN) Resize(n int) *VecN {
	if vn == nil {
		return NewVecN(n)
	}

	if n <= cap(vn.vec) {
		if vn.vec != nil {
			vn.vec = vn.vec[:n]
		} else {
			vn.vec = []float64{}
		}
		return vn
	}

	if shouldPool && vn.vec != nil {
		returnToPool(vn.vec)
	}
	*vn = (*NewVecN(n))

	return vn
}

// Sets the vector's backing slice to the given
// new one.
func (vn *VecN) SetBackingSlice(newSlice []float64) {
	vn.vec = newSlice
}

// Return the len of the vector's underlying slice.
// This is not titled Len because it conflicts the package's
// convention of calling the Norm the Len.
func (vn *VecN) Size() int {
	return len(vn.vec)
}

// Returns the cap of the vector's underlying slice.
func (vn *VecN) Cap() int {
	return cap(vn.vec)
}

// Sets the vector's size to n and zeroes out the vector.
// If n is bigger than the vector's size, it will realloc.
func (vn *VecN) Zero(n int) {
	vn.Resize(n)
	for i := range vn.vec {
		vn.vec[i] = 0
	}
}

// Adds vn and addend, storing the result in dst.
// If dst does not have sufficient size it will be resized
// Dst may be one of the other arguments. If dst is nil, it will be allocated.
// The value returned is dst, for easier method chaining
//
// If vn and addend are not the same size, this function will add min(vn.Size(), addend.Size())
// elements.
func (vn *VecN) Add(dst *VecN, subtrahend *VecN) *VecN {
	if vn == nil || subtrahend == nil {
		return nil
	}
	size := intMin(len(vn.vec), len(subtrahend.vec))
	dst = dst.Resize(size)

	for i := 0; i < size; i++ {
		dst.vec[i] = vn.vec[i] + subtrahend.vec[i]
	}

	return dst
}

// Subtracts addend from vn, storing the result in dst.
// If dst does not have sufficient size it will be resized
// Dst may be one of the other arguments. If dst is nil, it will be allocated.
// The value returned is dst, for easier method chaining
//
// If vn and addend are not the same size, this function will add min(vn.Size(), addend.Size())
// elements.
func (vn *VecN) Sub(dst *VecN, addend *VecN) *VecN {
	if vn == nil || addend == nil {
		return nil
	}
	size := intMin(len(vn.vec), len(addend.vec))
	dst = dst.Resize(size)

	for i := 0; i < size; i++ {
		dst.vec[i] = vn.vec[i] - addend.vec[i]
	}

	return dst
}

// Takes the binary cross product of vn and other, and stores it in dst.
// If either vn or other are not of size 3 this function will panic
//
// If dst is not of sufficient size, or is nil, a new slice is allocated.
// Dst is permitted to be one of the other arguments
func (vn *VecN) Cross(dst *VecN, other *VecN) *VecN {
	if vn == nil || other == nil {
		return nil
	}
	if len(vn.vec) != 3 || len(other.vec) != 3 {
		panic("Cannot take binary cross product of non-3D elements (7D cross product not implemented)")
	}

	dst = dst.Resize(3)
	dst.vec[0], dst.vec[1], dst.vec[2] = vn.vec[1]*other.vec[2]-vn.vec[2]*other.vec[1], vn.vec[2]*other.vec[0]-vn.vec[0]*other.vec[2], vn.vec[0]*other.vec[1]-vn.vec[1]*other.vec[0]

	return dst
}

func intMin(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// Computes the dot product of two VecNs, if
// the two vectors are not of the same length -- this
// will return NaN.
func (vn *VecN) Dot(other *VecN) float64 {
	if vn == nil || other == nil || len(vn.vec) != len(other.vec) {
		return float64(math.NaN())
	}

	var result float64 = 0.0
	for i, el := range vn.vec {
		result += el * other.vec[i]
	}

	return result
}

// Computes the vector length (also called the Norm) of the
// vector. Equivalent to math.Sqrt(vn.Dot(vn)) with the appropriate
// type conversions.
//
// If vn is nil, this returns NaN
func (vn *VecN) Len() float64 {
	if vn == nil {
		return float64(math.NaN())
	}
	if len(vn.vec) == 0 {
		return 0
	}

	return float64(math.Sqrt(float64(vn.Dot(vn))))
}

// Normalizes the vector and stores the result in dst, which
// will be returned. Dst will be appropraitely resized to the
// size of vn.
//
// The destination can be vn itself and nothing will go wrong.
//
// This is equivalent to vn.Mul(dst, 1/vn.Len())
func (vn *VecN) Normalize(dst *VecN) *VecN {
	if vn == nil {
		return nil
	}

	return vn.Mul(dst, 1/vn.Len())
}

// Multiplied the vector by some scalar value and stores the result in dst, which
// will be returned. Dst will be appropraitely resized to the
// size of vn.
//
// The destination can be vn itself and nothing will go wrong.
func (vn *VecN) Mul(dst *VecN, c float64) *VecN {
	if vn == nil {
		return nil
	}
	dst = dst.Resize(len(vn.vec))

	for i, el := range vn.vec {
		dst.vec[i] = el * c
	}

	return dst
}

// Performs the vector outer product between vn and v2.
// The outer product is like a "reverse" dot product. Where the dot product
// aligns both vectors with the "sized" part facing "inward" (Vec3*Vec3=Mat1x3*Mat3x1=Mat1x1=Scalar).
// The outer product multiplied them with it facing "outward"
// (Vec3*Vec3=Mat3x1*Mat1x3=Mat3x3).
//
// The matrix dst will be Reshaped to the correct size, if vn or v2 are nil,
// this returns nil.
func (vn *VecN) OuterProd(dst *MatMxN, v2 *VecN) *MatMxN {
	if vn == nil || v2 == nil {
		return nil
	}

	dst = dst.Reshape(len(vn.vec), len(v2.vec))

	for c, el1 := range v2.vec {
		for r, el2 := range vn.vec {
			dst.Set(r, c, el1*el2)
		}
	}

	return dst
}

func (vn *VecN) ApproxEqual(vn2 *VecN) bool {
	if vn == nil || vn2 == nil || len(vn.vec) != len(vn2.vec) {
		return false
	}

	for i, el := range vn.vec {
		if !FloatEqual(el, vn2.vec[i]) {
			return false
		}
	}

	return true
}

func (vn *VecN) ApproxEqualThreshold(vn2 *VecN, epsilon float64) bool {
	if vn == nil || vn2 == nil || len(vn.vec) != len(vn2.vec) {
		return false
	}

	for i, el := range vn.vec {
		if !FloatEqualThreshold(el, vn2.vec[i], epsilon) {
			return false
		}
	}

	return true
}

func (vn *VecN) ApproxEqualFunc(vn2 *VecN, comp func(float64, float64) bool) bool {
	if vn == nil || vn2 == nil || len(vn.vec) != len(vn2.vec) {
		return false
	}

	for i, el := range vn.vec {
		if !comp(el, vn2.vec[i]) {
			return false
		}
	}

	return true
}

func (vn *VecN) Vec2() Vec2 {
	raw := vn.Raw()
	return Vec2{raw[0], raw[1]}
}

func (vn *VecN) Vec3() Vec3 {
	raw := vn.Raw()
	return Vec3{raw[0], raw[1], raw[2]}
}

func (vn *VecN) Vec4() Vec4 {
	raw := vn.Raw()
	return Vec4{raw[0], raw[1], raw[2], raw[3]}
}
