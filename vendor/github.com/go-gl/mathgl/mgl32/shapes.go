// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mgl32

import (
	"fmt"
	"math"
)

// Generates a circle centered at (0,0) with a given radius.
// The radii are assumed to be in GL's coordinate sizing.
//
// Technically this draws an ellipse with two axes that match with the X and Y axes, the reason it has a radiusX and radiusY is because GL's coordinate system
// is proportional to screen width and screen height. So if you have a non-square viewport, a single radius will appear
// to "squash" the circle in one direction (usually the Y direction), so the X and Y radius allow for a circle to be made.
// A good way to get the correct radii is with mathgl.ScreenToGLCoords(radius, radius, screenWidth, screenHeight) which will get you the correct
// proportional GL coords.
//
// The numSlices argument specifies how many triangles you want your circle divided into, setting this
// number to too low a value may cause problem (and too high will cause it to take a lot of memory and time to compute
// without much gain in resolution).
//
// This uses discrete triangles, not a triangle fan
func Circle(radiusX, radiusY float32, numSlices int) []Vec2 {
	twoPi := float32(2.0 * math.Pi)

	circlePoints := make([]Vec2, 0, numSlices*3)
	center := Vec2{0.0, 0.0}
	previous := Vec2{radiusX, 0.0}

	for theta := twoPi / float32(numSlices); !FloatEqual(theta, twoPi); theta = Clamp(theta+twoPi/float32(numSlices), 0.0, twoPi) {
		sin, cos := math.Sincos(float64(theta))
		curr := Vec2{float32(cos) * radiusX, float32(sin) * radiusY}

		circlePoints = append(circlePoints, center, previous, curr)
		previous = curr
	}

	// Now add the final point at theta=2pi
	circlePoints = append(circlePoints, center, previous, Vec2{radiusX, 0.0})
	return circlePoints
}

// Generates a 2-triangle rectangle for use with GL_TRIANGLES. The width and height should use GL's proportions (that is, where a width of 1.0
// is equivalent to half of the width of the render target); however, the y-coordinates grow downwards, not upwards. That is, it
// assumes you want the origin of the rectangle with the top-left corner at (0.0,0.0).
//
// Keep in mind that GL's coordinate system is proportional, so width=height will not result in a square unless your viewport is square.
// If you want to maintain proportionality regardless of screen size, use the results of w,h := ScreenToGLCoordsf(absoluteWidth, absoluteHeight, screenWidth, screenHeight);
// w,h=w+1,h-1 in the call to this function. (The w+1,h-1 step maps the coordinates to start at 0.0 rather than -1.0)
func Rect(width, height float32) []Vec2 {
	return []Vec2{
		{0.0, 0.0},
		{0.0, -height},
		{width, -height},

		{0.0, 0.0},
		{width, -height},
		{width, 0.0},
	}
}

func QuadraticBezierCurve2D(t float32, cPoint1, cPoint2, cPoint3 Vec2) Vec2 {
	if t < 0.0 || t > 1.0 {
		panic("Can't interpolate on bezier curve with t out of range [0.0,1.0]")
	}

	return cPoint1.Mul((1.0 - t) * (1.0 - t)).Add(cPoint2.Mul(2 * (1 - t) * t)).Add(cPoint3.Mul(t * t))
}

func QuadraticBezierCurve3D(t float32, cPoint1, cPoint2, cPoint3 Vec3) Vec3 {
	if t < 0.0 || t > 1.0 {
		panic("Can't interpolate on bezier curve with t out of range [0.0,1.0]")
	}

	return cPoint1.Mul((1.0 - t) * (1.0 - t)).Add(cPoint2.Mul(2 * (1 - t) * t)).Add(cPoint3.Mul(t * t))
}

func CubicBezierCurve2D(t float32, cPoint1, cPoint2, cPoint3, cPoint4 Vec2) Vec2 {
	if t < 0.0 || t > 1.0 {
		panic("Can't interpolate on bezier curve with t out of range [0.0,1.0]")
	}

	return cPoint1.Mul((1 - t) * (1 - t) * (1 - t)).Add(cPoint2.Mul(3 * (1 - t) * (1 - t) * t)).Add(cPoint3.Mul(3 * (1 - t) * t * t)).Add(cPoint4.Mul(t * t * t))
}

func CubicBezierCurve3D(t float32, cPoint1, cPoint2, cPoint3, cPoint4 Vec3) Vec3 {
	if t < 0.0 || t > 1.0 {
		panic("Can't interpolate on bezier curve with t out of range [0.0,1.0]")
	}

	return cPoint1.Mul((1 - t) * (1 - t) * (1 - t)).Add(cPoint2.Mul(3 * (1 - t) * (1 - t) * t)).Add(cPoint3.Mul(3 * (1 - t) * t * t)).Add(cPoint4.Mul(t * t * t))
}

// Returns the point at point t along an n-control point Bezier curve
//
// t must be in the range 0.0 and 1.0 or this function will panic. Consider [0.0,1.0] to be similar to a percentage,
// 0.0 is first control point, and the point at 1.0 is the last control point. Any point in between is how far along the path you are between 0 and 1.
//
// This function is not sensitive to the coordinate system of the control points. It will correctly interpolate regardless of whether they're in screen coords,
// gl coords, or something else entirely
func BezierCurve2D(t float32, cPoints []Vec2) Vec2 {
	if t < 0.0 || t > 1.0 {
		panic("Input to bezier has t not in range [0,1]. If you think this is a precision error, use mathgl.Clamp[f|d] before calling this function")
	}

	n := len(cPoints) - 1
	point := cPoints[0].Mul(float32(math.Pow(float64(1.0-t), float64(n))))

	for i := 1; i <= n; i++ {
		point = point.Add(cPoints[i].Mul(float32(float64(choose(n, i)) * math.Pow(float64(1-t), float64(n-i)) * math.Pow(float64(t), float64(i))))) // P += P_i * nCi * (1-t)^(n-i) * t^i
	}

	return point
}

// Same as the 2D version, except the line is in 3D space
func BezierCurve3D(t float32, cPoints []Vec3) Vec3 {
	if t < 0.0 || t > 1.0 {
		panic("Input to bezier has t not in range [0,1]. If you think this is a precision error, use mathgl.Clamp[f|d] before calling this function")
	}

	n := len(cPoints) - 1
	point := cPoints[0].Mul(float32(math.Pow(float64(1.0-t), float64(n))))

	for i := 1; i <= n; i++ {
		point = point.Add(cPoints[i].Mul(float32(float64(choose(n, i)) * math.Pow(float64(1-t), float64(n-i)) * math.Pow(float64(t), float64(i))))) // P += P_i * nCi * (1-t)^(n-i) * t^i
	}

	return point
}

// Generates a bezier curve with controlPoints cPoints. The numPoints argument
// determines how many "samples" it makes along the line. For instance, a
// call to this with numPoints 2 will have exactly two points: the start and end points
// For any points above 2 it will divide it into numPoints-1 chunks (which means it will generate numPoints-2 vertices other than the beginning and end).
// So for 3 points it will divide it in half, 4 points into thirds, and so on.
//
// This is likely to get rather expensive for anything over perhaps a cubic curve.
func MakeBezierCurve2D(numPoints int, cPoints []Vec2) (line []Vec2) {
	line = make([]Vec2, numPoints)
	if numPoints == 0 {
		return
	} else if numPoints == 1 {
		line[0] = cPoints[0]
		return
	} else if numPoints == 2 {
		line[0] = cPoints[0]
		line[1] = cPoints[len(cPoints)-1]
		return
	}

	line[0] = cPoints[0]
	for i := 1; i < numPoints-1; i++ {
		line[i] = BezierCurve2D(Clamp(float32(i)/float32(numPoints-1), 0.0, 1.0), cPoints)
	}
	line[numPoints-1] = cPoints[len(cPoints)-1]

	return
}

// Same as the 2D version, except with the line in 3D space
func MakeBezierCurve3D(numPoints int, cPoints []Vec3) (line []Vec3) {
	line = make([]Vec3, numPoints)
	if numPoints == 0 {
		return
	} else if numPoints == 1 {
		line[0] = cPoints[0]
		return
	} else if numPoints == 2 {
		line[0] = cPoints[0]
		line[1] = cPoints[len(cPoints)-1]
		return
	}

	line[0] = cPoints[0]
	for i := 1; i < numPoints-1; i++ {
		line[i] = BezierCurve3D(Clamp(float32(i)/float32(numPoints-1), 0.0, 1.0), cPoints)
	}
	line[numPoints-1] = cPoints[len(cPoints)-1]

	return
}

// Creates a 2-dimensional Bezier surface of arbitrary degree in 3D Space
// Like the curve functions, if u or v are not in the range [0.0,1.0] the function will panic, use Clamp[f|d]
// to ensure it is correct.
//
// The control point matrix must not be jagged, or this will end up panicking from an index out of bounds exception
func BezierSurface(u, v float32, cPoints [][]Vec3) Vec3 {
	if u < 0.0 || u > 1.0 || v < 1.0 || v > 1.0 {
		panic("u or v not in range [0.0,1.0] in BezierSurface")
	}

	n := len(cPoints) - 1
	m := len(cPoints[0]) - 1

	point := cPoints[0][0].Mul(float32(math.Pow(float64(1.0-u), float64(n)) * math.Pow(float64(1.0-v), float64(m))))

	for i := 0; i <= n; i++ {
		for j := 0; j <= m; j++ {
			if i == 0 && j == 0 {
				continue
			}

			point = point.Add(cPoints[i][j].Mul(float32(float64(choose(n, i)) * math.Pow(float64(u), float64(i)) * math.Pow(float64(1.0-u), float64(n-i)) * float64(choose(m, j)) * math.Pow(float64(v), float64(j)) * math.Pow(float64(1.0-v), float64(m-j)))))
		}
	}

	return point
}

// Does interpolation over a spline of several bezier curves. Each bezier curve must have a finite range,
// though the spline may be disjoint. The bezier curves are not required to be in any particular order.
//
// If t is out of the range of all given curves, this function will panic
func BezierSplineInterpolate2D(t float32, ranges [][2]float32, cPoints [][]Vec2) Vec2 {
	if len(ranges) != len(cPoints) {
		panic("Each bezier curve needs a range")
	}

	for i, curveRange := range ranges {
		if t >= curveRange[0] && t <= curveRange[1] {
			return BezierCurve2D((t-curveRange[0])/(curveRange[1]-curveRange[0]), cPoints[i])
		}
	}

	panic("t is out of the range of all bezier curves in this spline")
}

// Does interpolation over a spline of several bezier curves. Each bezier curve must have a finite range,
// though the spline may be disjoint. The bezier curves are not required to be in any particular order.
//
// If t is out of the range of all given curves, this function will panic
func BezierSplineInterpolate3D(t float32, ranges [][2]float32, cPoints [][]Vec3) Vec3 {
	if len(ranges) != len(cPoints) {
		panic("Each bezier curve needs a range")
	}

	for i, curveRange := range ranges {
		if t >= curveRange[0] && t <= curveRange[1] {
			return BezierCurve3D((t-curveRange[0])/(curveRange[1]-curveRange[0]), cPoints[i])
		}
	}

	panic("t is out of the range of all bezier curves in this spline")
}

// Reticulates ALL the Splines
//
// For the overly serious: the function is just for fun. It does nothing except prints a Maxis reference. Technically you could "reticulate splines"
// by joining a bunch of splines together, but that ruins the joke.
func ReticulateSplines(ranges [][][2]float32, cPoints [][][]Vec2, withLlamas bool) {
	if !withLlamas {
		fmt.Println("You can't reticulate splines without llamas, silly.")
	} else {
		fmt.Println("Actually, you can't even reticulate splines WITH llamas")
	}
}

// Transform from pixel coordinates to GL coordinates.
//
// This assumes that your pixel coordinate system considers its origin to be in the top left corner (GL's is in the bottom left).
// The coordinates x and y may be out of the range [0,screenWidth-1] and [0,screeneHeight-1].
//
// GL's coordinate system maps [screenWidth-1,0] to [1.0,1.0] and [0,screenHeight-1] to [-1.0,-1.0]. If x and y are out of the range, they'll still
// be mapped correctly, just off the screen. (e.g. if y = 2*(screenHeight-1) you'll get -3.0 for yOut)
//
// This is similar to Unproject, except for 2D cases and much simpler (especially since an inverse may always be found)
func ScreenToGLCoords(x, y int, screenWidth, screenHeight int) (xOut, yOut float32) {
	xOut = 2.0*float32(x)/float32(screenWidth-1) - 1.0
	yOut = -2.0*float32(y)/float32(screenHeight-1) + 1.0

	return
}

// Transform from GL's proportional system to pixel coordinates.
//
// Assumes the pixel coordinate system has its origin in the top left corner. (GL's is in the bottom left)
//
// GL's coordinate system maps [screenWidth-1,0] to [1.0,1.0] and [0,screenHeight-1] to [-1.0,-1.0]. If x and y are out of the range, they'll still
// be mapped correctly, just off the screen. (e.g. if y=-3.0, you'll get 2*(screenHeight-1) for yOut)
//
// This is similar to Project, except for 2D cases and much simpler
func GLToScreenCoords(x, y float32, screenWidth, screenHeight int) (xOut, yOut int) {
	xOut = int((x + 1.0) * float32(screenWidth-1) / 2.0)
	yOut = int((1.0 - y) * float32(screenHeight-1) / 2.0)

	return
}

func choose(n, k int) (result int) {
	if k == 0 {
		return 1
	} else if n == 0 {
		return 0
	}
	result = (n - (k - 1))
	for i := 2; i <= k; i++ {
		result *= (n - (k - i)) / i
	}

	return result
}
