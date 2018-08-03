// This file is generated from mgl32/quat.go; DO NOT EDIT

// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mgl64

import (
	"math"
)

// A rotation order is the order in which
// rotations will be transformed for the purposes of AnglesToQuat
type RotationOrder int

const (
	XYX RotationOrder = iota
	XYZ
	XZX
	XZY
	YXY
	YXZ
	YZY
	YZX
	ZYZ
	ZYX
	ZXZ
	ZXY
)

// A Quaternion is an extension of the imaginary numbers; there's all sorts of
// interesting theory behind it. In 3D graphics we mostly use it as a cheap way of
// representing rotation since quaternions are cheaper to multiply by, and easier to
// interpolate than matrices.
//
// A Quaternion has two parts: W, the so-called scalar component,
// and "V", the vector component. The vector component is considered to
// be the part in 3D space, while W (loosely interpreted) is its 4D coordinate.
type Quat struct {
	W float64
	V Vec3
}

// The quaternion identity: W=1; V=(0,0,0).
//
// As with all identities, multiplying any quaternion by this will yield the same
// quaternion you started with.
func QuatIdent() Quat {
	return Quat{1., Vec3{0, 0, 0}}
}

// Creates an angle from an axis and an angle relative to that axis.
//
// This is cheaper than HomogRotate3D.
func QuatRotate(angle float64, axis Vec3) Quat {
	// angle = (float32(math.Pi) * angle) / 180.0

	c, s := float64(math.Cos(float64(angle/2))), float64(math.Sin(float64(angle/2)))

	return Quat{c, axis.Mul(s)}
}

// A convenient alias for q.V[0]
func (q Quat) X() float64 {
	return q.V[0]
}

// A convenient alias for q.V[1]
func (q Quat) Y() float64 {
	return q.V[1]
}

// A convenient alias for q.V[2]
func (q Quat) Z() float64 {
	return q.V[2]
}

// Adds two quaternions. It's no more complicated than
// adding their W and V components.
func (q1 Quat) Add(q2 Quat) Quat {
	return Quat{q1.W + q2.W, q1.V.Add(q2.V)}
}

// Subtracts two quaternions. It's no more complicated than
// subtracting their W and V components.
func (q1 Quat) Sub(q2 Quat) Quat {
	return Quat{q1.W - q2.W, q1.V.Sub(q2.V)}
}

// Multiplies two quaternions. This can be seen as a rotation. Note that
// Multiplication is NOT commutative, meaning q1.Mul(q2) does not necessarily
// equal q2.Mul(q1).
func (q1 Quat) Mul(q2 Quat) Quat {
	return Quat{q1.W*q2.W - q1.V.Dot(q2.V), q1.V.Cross(q2.V).Add(q2.V.Mul(q1.W)).Add(q1.V.Mul(q2.W))}
}

// Scales every element of the quaternion by some constant factor.
func (q1 Quat) Scale(c float64) Quat {
	return Quat{q1.W * c, Vec3{q1.V[0] * c, q1.V[1] * c, q1.V[2] * c}}
}

// Returns the conjugate of a quaternion. Equivalent to
// Quat{q1.W, q1.V.Mul(-1)}
func (q1 Quat) Conjugate() Quat {
	return Quat{q1.W, q1.V.Mul(-1)}
}

// Returns the Length of the quaternion, also known as its Norm. This is the same thing as
// the Len of a Vec4
func (q1 Quat) Len() float64 {
	return float64(math.Sqrt(float64(q1.W*q1.W + q1.V[0]*q1.V[0] + q1.V[1]*q1.V[1] + q1.V[2]*q1.V[2])))
}

// Norm() is an alias for Len() since both are very common terms.
func (q1 Quat) Norm() float64 {
	return q1.Len()
}

// Normalizes the quaternion, returning its versor (unit quaternion).
//
// This is the same as normalizing it as a Vec4.
func (q1 Quat) Normalize() Quat {
	length := q1.Len()

	if FloatEqual(1, length) {
		return q1
	}
	if length == 0 {
		return QuatIdent()
	}
	if length == InfPos {
		length = MaxValue
	}

	return Quat{q1.W * 1 / length, q1.V.Mul(1 / length)}
}

// The inverse of a quaternion. The inverse is equivalent
// to the conjugate divided by the square of the length.
//
// This method computes the square norm by directly adding the sum
// of the squares of all terms instead of actually squaring q1.Len(),
// both for performance and precision.
func (q1 Quat) Inverse() Quat {
	return q1.Conjugate().Scale(1 / q1.Dot(q1))
}

// Rotates a vector by the rotation this quaternion represents.
// This will result in a 3D vector. Strictly speaking, this is
// equivalent to q1.v.q* where the "."" is quaternion multiplication and v is interpreted
// as a quaternion with W 0 and V v. In code:
// q1.Mul(Quat{0,v}).Mul(q1.Conjugate()), and
// then retrieving the imaginary (vector) part.
//
// In practice, we hand-compute this in the general case and simplify
// to save a few operations.
func (q1 Quat) Rotate(v Vec3) Vec3 {
	cross := q1.V.Cross(v)
	// v + 2q_w * (q_v x v) + 2q_v x (q_v x v)
	return v.Add(cross.Mul(2 * q1.W)).Add(q1.V.Mul(2).Cross(cross))
}

// Returns the homogeneous 3D rotation matrix corresponding to the quaternion.
func (q1 Quat) Mat4() Mat4 {
	w, x, y, z := q1.W, q1.V[0], q1.V[1], q1.V[2]
	return Mat4{
		1 - 2*y*y - 2*z*z, 2*x*y + 2*w*z, 2*x*z - 2*w*y, 0,
		2*x*y - 2*w*z, 1 - 2*x*x - 2*z*z, 2*y*z + 2*w*x, 0,
		2*x*z + 2*w*y, 2*y*z - 2*w*x, 1 - 2*x*x - 2*y*y, 0,
		0, 0, 0, 1,
	}
}

// The dot product between two quaternions, equivalent to if this was a Vec4
func (q1 Quat) Dot(q2 Quat) float64 {
	return q1.W*q2.W + q1.V[0]*q2.V[0] + q1.V[1]*q2.V[1] + q1.V[2]*q2.V[2]
}

// Returns whether the quaternions are approximately equal, as if
// FloatEqual was called on each matching element
func (q1 Quat) ApproxEqual(q2 Quat) bool {
	return FloatEqual(q1.W, q2.W) && q1.V.ApproxEqual(q2.V)
}

// Returns whether the quaternions are approximately equal with a given tolerence, as if
// FloatEqualThreshold was called on each matching element with the given epsilon
func (q1 Quat) ApproxEqualThreshold(q2 Quat, epsilon float64) bool {
	return FloatEqualThreshold(q1.W, q2.W, epsilon) && q1.V.ApproxEqualThreshold(q2.V, epsilon)
}

// Returns whether the quaternions are approximately equal using the given comparison function, as if
// the function had been called on each individual element
func (q1 Quat) ApproxEqualFunc(q2 Quat, f func(float64, float64) bool) bool {
	return f(q1.W, q2.W) && q1.V.ApproxFuncEqual(q2.V, f)
}

// Returns whether the quaternions represents the same orientation
//
// Different values can represent the same orientation (q == -q) because quaternions avoid singularities
// and discontinuities involved with rotation in 3 dimensions by adding extra dimensions
func (q1 Quat) OrientationEqual(q2 Quat) bool {
	return q1.OrientationEqualThreshold(q2, Epsilon)
}

// Returns whether the quaternions represents the same orientation with a given tolerence
func (q1 Quat) OrientationEqualThreshold(q2 Quat, epsilon float64) bool {
	return Abs(q1.Normalize().Dot(q2.Normalize())) > 1-epsilon
}

// Slerp is *S*pherical *L*inear Int*erp*olation, a method of interpolating
// between two quaternions. This always takes the straightest path on the sphere between
// the two quaternions, and maintains constant velocity.
//
// However, it's expensive and QuatSlerp(q1,q2) is not the same as QuatSlerp(q2,q1)
func QuatSlerp(q1, q2 Quat, amount float64) Quat {
	q1, q2 = q1.Normalize(), q2.Normalize()
	dot := q1.Dot(q2)

	// If the inputs are too close for comfort, linearly interpolate and normalize the result.
	if dot > 0.9995 {
		return QuatNlerp(q1, q2, amount)
	}

	// This is here for precision errors, I'm perfectly aware that *technically* the dot is bound [-1,1], but since Acos will freak out if it's not (even if it's just a liiiiitle bit over due to normal error) we need to clamp it
	dot = Clamp(dot, -1, 1)

	theta := float64(math.Acos(float64(dot))) * amount
	c, s := float64(math.Cos(float64(theta))), float64(math.Sin(float64(theta)))
	rel := q2.Sub(q1.Scale(dot)).Normalize()

	return q1.Scale(c).Add(rel.Scale(s))
}

// *L*inear Int*erp*olation between two Quaternions, cheap and simple.
//
// Not excessively useful, but uses can be found.
func QuatLerp(q1, q2 Quat, amount float64) Quat {
	return q1.Add(q2.Sub(q1).Scale(amount))
}

// *Normalized* *L*inear Int*erp*olation between two Quaternions. Cheaper than Slerp
// and usually just as good. This is literally Lerp with Normalize() called on it.
//
// Unlike Slerp, constant velocity isn't maintained, but it's much faster and
// Nlerp(q1,q2) and Nlerp(q2,q1) return the same path. You should probably
// use this more often unless you're suffering from choppiness due to the
// non-constant velocity problem.
func QuatNlerp(q1, q2 Quat, amount float64) Quat {
	return QuatLerp(q1, q2, amount).Normalize()
}

// Performs a rotation in the specified order. If the order is not
// a valid RotationOrder, this function will panic
//
// The rotation "order" is more of an axis descriptor. For instance XZX would
// tell the function to interpret angle1 as a rotation about the X axis, angle2 about
// the Z axis, and angle3 about the X axis again.
//
// Based off the code for the Matlab function "angle2quat", though this implementation
// only supports 3 single angles as opposed to multiple angles.
func AnglesToQuat(angle1, angle2, angle3 float64, order RotationOrder) Quat {
	var s [3]float64
	var c [3]float64

	s[0], c[0] = math.Sincos(float64(angle1 / 2))
	s[1], c[1] = math.Sincos(float64(angle2 / 2))
	s[2], c[2] = math.Sincos(float64(angle3 / 2))

	ret := Quat{}
	switch order {
	case ZYX:
		ret.W = float64(c[0]*c[1]*c[2] + s[0]*s[1]*s[2])
		ret.V = Vec3{float64(c[0]*c[1]*s[2] - s[0]*s[1]*c[2]),
			float64(c[0]*s[1]*c[2] + s[0]*c[1]*s[2]),
			float64(s[0]*c[1]*c[2] - c[0]*s[1]*s[2]),
		}
	case ZYZ:
		ret.W = float64(c[0]*c[1]*c[2] - s[0]*c[1]*s[2])
		ret.V = Vec3{float64(c[0]*s[1]*s[2] - s[0]*s[1]*c[2]),
			float64(c[0]*s[1]*c[2] + s[0]*s[1]*s[2]),
			float64(s[0]*c[1]*c[2] + c[0]*c[1]*s[2]),
		}
	case ZXY:
		ret.W = float64(c[0]*c[1]*c[2] - s[0]*s[1]*s[2])
		ret.V = Vec3{float64(c[0]*s[1]*c[2] - s[0]*c[1]*s[2]),
			float64(c[0]*c[1]*s[2] + s[0]*s[1]*c[2]),
			float64(c[0]*s[1]*s[2] + s[0]*c[1]*c[2]),
		}
	case ZXZ:
		ret.W = float64(c[0]*c[1]*c[2] - s[0]*c[1]*s[2])
		ret.V = Vec3{float64(c[0]*s[1]*c[2] + s[0]*s[1]*s[2]),
			float64(s[0]*s[1]*c[2] - c[0]*s[1]*s[2]),
			float64(c[0]*c[1]*s[2] + s[0]*c[1]*c[2]),
		}
	case YXZ:
		ret.W = float64(c[0]*c[1]*c[2] + s[0]*s[1]*s[2])
		ret.V = Vec3{float64(c[0]*s[1]*c[2] + s[0]*c[1]*s[2]),
			float64(s[0]*c[1]*c[2] - c[0]*s[1]*s[2]),
			float64(c[0]*c[1]*s[2] - s[0]*s[1]*c[2]),
		}
	case YXY:
		ret.W = float64(c[0]*c[1]*c[2] - s[0]*c[1]*s[2])
		ret.V = Vec3{float64(c[0]*s[1]*c[2] + s[0]*s[1]*s[2]),
			float64(s[0]*c[1]*c[2] + c[0]*c[1]*s[2]),
			float64(c[0]*s[1]*s[2] - s[0]*s[1]*c[2]),
		}
	case YZX:
		ret.W = float64(c[0]*c[1]*c[2] - s[0]*s[1]*s[2])
		ret.V = Vec3{float64(c[0]*c[1]*s[2] + s[0]*s[1]*c[2]),
			float64(c[0]*s[1]*s[2] + s[0]*c[1]*c[2]),
			float64(c[0]*s[1]*c[2] - s[0]*c[1]*s[2]),
		}
	case YZY:
		ret.W = float64(c[0]*c[1]*c[2] - s[0]*c[1]*s[2])
		ret.V = Vec3{float64(s[0]*s[1]*c[2] - c[0]*s[1]*s[2]),
			float64(c[0]*c[1]*s[2] + s[0]*c[1]*c[2]),
			float64(c[0]*s[1]*c[2] + s[0]*s[1]*s[2]),
		}
	case XYZ:
		ret.W = float64(c[0]*c[1]*c[2] - s[0]*s[1]*s[2])
		ret.V = Vec3{float64(c[0]*s[1]*s[2] + s[0]*c[1]*c[2]),
			float64(c[0]*s[1]*c[2] - s[0]*c[1]*s[2]),
			float64(c[0]*c[1]*s[2] + s[0]*s[1]*c[2]),
		}
	case XYX:
		ret.W = float64(c[0]*c[1]*c[2] - s[0]*c[1]*s[2])
		ret.V = Vec3{float64(c[0]*c[1]*s[2] + s[0]*c[1]*c[2]),
			float64(c[0]*s[1]*c[2] + s[0]*s[1]*s[2]),
			float64(s[0]*s[1]*c[2] - c[0]*s[1]*s[2]),
		}
	case XZY:
		ret.W = float64(c[0]*c[1]*c[2] + s[0]*s[1]*s[2])
		ret.V = Vec3{float64(s[0]*c[1]*c[2] - c[0]*s[1]*s[2]),
			float64(c[0]*c[1]*s[2] - s[0]*s[1]*c[2]),
			float64(c[0]*s[1]*c[2] + s[0]*c[1]*s[2]),
		}
	case XZX:
		ret.W = float64(c[0]*c[1]*c[2] - s[0]*c[1]*s[2])
		ret.V = Vec3{float64(c[0]*c[1]*s[2] + s[0]*c[1]*c[2]),
			float64(c[0]*s[1]*s[2] - s[0]*s[1]*c[2]),
			float64(c[0]*s[1]*c[2] + s[0]*s[1]*s[2]),
		}
	default:
		panic("Unsupported rotation order")
	}
	return ret
}

// Mat4ToQuat converts a pure rotation matrix into a quaternion
func Mat4ToQuat(m Mat4) Quat {
	// http://www.euclideanspace.com/maths/geometry/rotations/conversions/matrixToQuaternion/index.htm

	if tr := m[0] + m[5] + m[10]; tr > 0 {
		s := float64(0.5 / math.Sqrt(float64(tr+1.0)))
		return Quat{
			0.25 / s,
			Vec3{
				(m[6] - m[9]) * s,
				(m[8] - m[2]) * s,
				(m[1] - m[4]) * s,
			},
		}
	}

	if (m[0] > m[5]) && (m[0] > m[10]) {
		s := float64(2.0 * math.Sqrt(float64(1.0+m[0]-m[5]-m[10])))
		return Quat{
			(m[6] - m[9]) / s,
			Vec3{
				0.25 * s,
				(m[4] + m[1]) / s,
				(m[8] + m[2]) / s,
			},
		}
	}

	if m[5] > m[10] {
		s := float64(2.0 * math.Sqrt(float64(1.0+m[5]-m[0]-m[10])))
		return Quat{
			(m[8] - m[2]) / s,
			Vec3{
				(m[4] + m[1]) / s,
				0.25 * s,
				(m[9] + m[6]) / s,
			},
		}

	}

	s := float64(2.0 * math.Sqrt(float64(1.0+m[10]-m[0]-m[5])))
	return Quat{
		(m[1] - m[4]) / s,
		Vec3{
			(m[8] + m[2]) / s,
			(m[9] + m[6]) / s,
			0.25 * s,
		},
	}
}

// QuatLookAtV creates a rotation from an eye vector to a center vector
//
// It assumes the front of the rotated object at Z- and up at Y+
func QuatLookAtV(eye, center, up Vec3) Quat {
	// http://www.opengl-tutorial.org/intermediate-tutorials/tutorial-17-quaternions/#I_need_an_equivalent_of_gluLookAt__How_do_I_orient_an_object_towards_a_point__
	// https://bitbucket.org/sinbad/ogre/src/d2ef494c4a2f5d6e2f0f17d3bfb9fd936d5423bb/OgreMain/src/OgreCamera.cpp?at=default#cl-161

	direction := center.Sub(eye).Normalize()

	// Find the rotation between the front of the object (that we assume towards Z-,
	// but this depends on your model) and the desired direction
	rotDir := QuatBetweenVectors(Vec3{0, 0, -1}, direction)

	// Recompute up so that it's perpendicular to the direction
	// You can skip that part if you really want to force up
	//right := direction.Cross(up)
	//up = right.Cross(direction)

	// Because of the 1rst rotation, the up is probably completely screwed up.
	// Find the rotation between the "up" of the rotated object, and the desired up
	upCur := rotDir.Rotate(Vec3{0, 1, 0})
	rotUp := QuatBetweenVectors(upCur, up)

	rotTarget := rotUp.Mul(rotDir) // remember, in reverse order.
	return rotTarget.Inverse()     // camera rotation should be inversed!
}

// QuatBetweenVectors calculates the rotation between two vectors
func QuatBetweenVectors(start, dest Vec3) Quat {
	// http://www.opengl-tutorial.org/intermediate-tutorials/tutorial-17-quaternions/#I_need_an_equivalent_of_gluLookAt__How_do_I_orient_an_object_towards_a_point__
	// https://github.com/g-truc/glm/blob/0.9.5/glm/gtx/quaternion.inl#L225
	// https://bitbucket.org/sinbad/ogre/src/d2ef494c4a2f5d6e2f0f17d3bfb9fd936d5423bb/OgreMain/include/OgreVector3.h?at=default#cl-654

	start = start.Normalize()
	dest = dest.Normalize()
	epsilon := float64(0.001)

	cosTheta := start.Dot(dest)
	if cosTheta < -1.0+epsilon {
		// special case when vectors in opposite directions:
		// there is no "ideal" rotation axis
		// So guess one; any will do as long as it's perpendicular to start
		axis := Vec3{1, 0, 0}.Cross(start)
		if axis.Dot(axis) < epsilon {
			// bad luck, they were parallel, try again!
			axis = Vec3{0, 1, 0}.Cross(start)
		}

		return QuatRotate(math.Pi, axis.Normalize())
	}

	axis := start.Cross(dest)
	s := float64(math.Sqrt(float64(1.0+cosTheta) * 2.0))

	return Quat{
		s * 0.5,
		axis.Mul(1.0 / s),
	}
}
