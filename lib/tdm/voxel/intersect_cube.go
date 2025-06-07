package voxel

/*
Copyright 2020 Tomas Akenine-Möller

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
documentation files (the "Software"), to deal in the Software without restriction, including without limitation
the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and
to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial
portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS
OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT
OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

// From: https://fileadmin.cs.lth.se/cs/Personal/Tomas_Akenine-Moller/code/tribox3.txt

import (
	"math"
)

// Constants for indexing
const (
	X = 0
	Y = 1
	Z = 2
)

// cross computes the cross product of v1 and v2, storing it in dest.
func cross(dest, v1, v2 *[3]float64) {
	dest[0] = v1[1]*v2[2] - v1[2]*v2[1]
	dest[1] = v1[2]*v2[0] - v1[0]*v2[2]
	dest[2] = v1[0]*v2[1] - v1[1]*v2[0]
}

// dot computes the dot product of v1 and v2.
func dot(v1, v2 [3]float64) float64 {
	return v1[0]*v2[0] + v1[1]*v2[1] + v1[2]*v2[2]
}

// sub computes the vector subtraction v1 - v2 and stores it in dest.
func sub(dest, v1, v2 *[3]float64) {
	dest[0] = v1[0] - v2[0]
	dest[1] = v1[1] - v2[1]
	dest[2] = v1[2] - v2[2]
}

// findMinMax finds the min and max values among x0, x1, and x2.
func findMinMax(x0, x1, x2 float64) (float64, float64) {
	min, max := x0, x0
	if x1 < min {
		min = x1
	}
	if x1 > max {
		max = x1
	}
	if x2 < min {
		min = x2
	}
	if x2 > max {
		max = x2
	}
	return min, max
}

// planeBoxOverlap tests if the box overlaps the plane of the triangle.
func planeBoxOverlap(normal, vert [3]float64, boxHalfSize *[3]float64) bool {
	var vmin, vmax [3]float64
	for q := X; q <= Z; q++ {
		v := vert[q]
		if normal[q] > 0.0 {
			vmin[q] = -boxHalfSize[q] - v
			vmax[q] = boxHalfSize[q] - v
		} else {
			vmin[q] = boxHalfSize[q] - v
			vmax[q] = -boxHalfSize[q] - v
		}
	}
	if dot(normal, vmin) > 0.0 {
		return false
	}
	if dot(normal, vmax) >= 0.0 {
		return true
	}
	return false
}

// Axis test functions
func axisTestX01(a, b, fa, fb float64, v0, v2 [3]float64, boxHalfSize *[3]float64) bool {
	p0 := a*v0[Y] - b*v0[Z]
	p2 := a*v2[Y] - b*v2[Z]
	min, max := p0, p2
	if p0 > p2 {
		min, max = p2, p0
	}
	rad := fa*boxHalfSize[Y] + fb*boxHalfSize[Z]
	return !(min > rad || max < -rad)
}

func axisTestX2(a, b, fa, fb float64, v0, v1 [3]float64, boxHalfSize *[3]float64) bool {
	p0 := a*v0[Y] - b*v0[Z]
	p1 := a*v1[Y] - b*v1[Z]
	min, max := p0, p1
	if p0 > p1 {
		min, max = p1, p0
	}
	rad := fa*boxHalfSize[Y] + fb*boxHalfSize[Z]
	return !(min > rad || max < -rad)
}

func axisTestY02(a, b, fa, fb float64, v0, v2 [3]float64, boxHalfSize *[3]float64) bool {
	p0 := -a*v0[X] + b*v0[Z]
	p2 := -a*v2[X] + b*v2[Z]
	min, max := p0, p2
	if p0 > p2 {
		min, max = p2, p0
	}
	rad := fa*boxHalfSize[X] + fb*boxHalfSize[Z]
	return !(min > rad || max < -rad)
}

func axisTestY1(a, b, fa, fb float64, v0, v1 [3]float64, boxHalfSize *[3]float64) bool {
	p0 := -a*v0[X] + b*v0[Z]
	p1 := -a*v1[X] + b*v1[Z]
	min, max := p0, p1
	if p0 > p1 {
		min, max = p1, p0
	}
	rad := fa*boxHalfSize[X] + fb*boxHalfSize[Z]
	return !(min > rad || max < -rad)
}

func axisTestZ12(a, b, fa, fb float64, v1, v2 [3]float64, boxHalfSize *[3]float64) bool {
	p1 := a*v1[X] - b*v1[Y]
	p2 := a*v2[X] - b*v2[Y]
	min, max := p2, p1
	if p2 > p1 {
		min, max = p1, p2
	}
	rad := fa*boxHalfSize[X] + fb*boxHalfSize[Y]
	return !(min > rad || max < -rad)
}

func axisTestZ0(a, b, fa, fb float64, v0, v1 [3]float64, boxHalfSize *[3]float64) bool {
	p0 := a*v0[X] - b*v0[Y]
	p1 := a*v1[X] - b*v1[Y]
	min, max := p0, p1
	if p0 > p1 {
		min, max = p1, p0
	}
	rad := fa*boxHalfSize[X] + fb*boxHalfSize[Y]
	return !(min > rad || max < -rad)
}

func TriBoxOverlap(boxCenter, boxHalfSize *[3]float64, triverts *[3][3]float64) bool {
	var v0, v1, v2 [3]float64
	var e0, e1, e2 [3]float64
	var normal [3]float64
	var min, max, fex, fey, fez float64

	sub(&v0, &triverts[0], boxCenter)
	sub(&v1, &triverts[1], boxCenter)
	sub(&v2, &triverts[2], boxCenter)

	//// **Degeneracy Check: If all three vertices are identical, it's a single point**
	//if v0 == v1 && v1 == v2 {
	//	// Check if the point is inside the box
	//	if math.Abs(v0[X]) <= boxHalfSize[X] &&
	//		math.Abs(v0[Y]) <= boxHalfSize[Y] &&
	//		math.Abs(v0[Z]) <= boxHalfSize[Z] {
	//		return true // The point is inside the box
	//	}
	//	return false // The point is outside the box
	//}

	sub(&e0, &v1, &v0)
	sub(&e1, &v2, &v1)
	sub(&e2, &v0, &v2)

	fex, fey, fez = math.Abs(e0[X]), math.Abs(e0[Y]), math.Abs(e0[Z])
	if !axisTestX01(e0[Z], e0[Y], fez, fey, v0, v2, boxHalfSize) ||
		!axisTestY02(e0[Z], e0[X], fez, fex, v0, v2, boxHalfSize) ||
		!axisTestZ12(e0[Y], e0[X], fey, fex, v1, v2, boxHalfSize) {
		return false
	}

	fex, fey, fez = math.Abs(e1[X]), math.Abs(e1[Y]), math.Abs(e1[Z])
	if !axisTestX2(e1[Z], e1[Y], fez, fey, v0, v1, boxHalfSize) ||
		!axisTestY1(e1[Z], e1[X], fez, fex, v0, v1, boxHalfSize) ||
		!axisTestZ0(e1[Y], e1[X], fey, fex, v0, v1, boxHalfSize) {
		return false
	}

	min, max = findMinMax(v0[X], v1[X], v2[X])
	if min > boxHalfSize[X] || max < -boxHalfSize[X] {
		return false
	}

	cross(&normal, &e0, &e1)
	return planeBoxOverlap(normal, v0, boxHalfSize)
}
