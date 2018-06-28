// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mgl32

import (
	"math"
)

// Converts 3-dimensional cartesian coordinates (x,y,z) to spherical
// coordinates with radius r, inclination theta, and azimuth phi.
//
// All angles are in radians.
func CartesianToSpherical(coord Vec3) (r, theta, phi float32) {
	r = coord.Len()
	theta = float32(math.Acos(float64(coord[2] / r)))
	phi = float32(math.Atan2(float64(coord[1]), float64(coord[0])))

	return
}

// Converts 3-dimensional cartesian coordinates (x,y,z) to cylindrical
// coordinates with radial distance r, azimuth phi, and height z.
//
// All angles are in radians.
func CartesianToCylindical(coord Vec3) (rho, phi, z float32) {
	rho = float32(math.Hypot(float64(coord[0]), float64(coord[1])))

	phi = float32(math.Atan2(float64(coord[1]), float64(coord[0])))

	z = coord[2]

	return
}

// Converts spherical coordinates with radius r, inclination theta,
// and azimuth phi to cartesian coordinates (x,y,z).
//
// Angles are in radians.
func SphericalToCartesian(r, theta, phi float32) Vec3 {
	st, ct := math.Sincos(float64(theta))
	sp, cp := math.Sincos(float64(phi))

	return Vec3{r * float32(st*cp), r * float32(st*sp), r * float32(ct)}
}

// Converts spherical coordinates with radius r, inclination theta,
// and azimuth phi to cylindrical coordinates with radial distance r,
// azimuth phi, and height z.
//
// Angles are in radians
func SphericalToCylindrical(r, theta, phi float32) (rho, phi2, z float32) {
	s, c := math.Sincos(float64(theta))

	rho = r * float32(s)
	z = r * float32(c)
	phi2 = phi

	return
}

// Converts cylindrical coordinates with radial distance r,
// azimuth phi, and height z to spherical coordinates with radius r,
// inclination theta, and azimuth phi.
//
// Angles are in radians
func CylindircalToSpherical(rho, phi, z float32) (r, theta, phi2 float32) {
	r = float32(math.Hypot(float64(rho), float64(z)))
	phi2 = phi
	theta = float32(math.Atan2(float64(rho), float64(z)))

	return
}

// Converts cylindrical coordinates with radial distance r,
// azimuth phi, and height z to cartesian coordinates (x,y,z)
//
// Angles are in radians.
func CylindricalToCartesian(rho, phi, z float32) Vec3 {
	s, c := math.Sincos(float64(phi))

	return Vec3{rho * float32(c), rho * float32(s), z}
}

// Converts degrees to radians
func DegToRad(angle float32) float32 {
	return angle * float32(math.Pi) / 180
}

// Converts radians to degrees
func RadToDeg(angle float32) float32 {
	return angle * 180 / float32(math.Pi)
}
