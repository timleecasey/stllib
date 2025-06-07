package voxel

import "testing"

func TestBasicOverlap(t *testing.T) {
	boxcenter := &[3]float64{0, 0, 0}
	boxhalfsize := &[3]float64{1, 1, 1}
	triverts := &[3][3]float64{
		{0.5, 0.5, 0.5},
		{-0.5, -0.5, 0.5},
		{0.5, -0.5, -0.5},
	}

	result := TriBoxOverlap(boxcenter, boxhalfsize, triverts)
	if result {
		println("Triangle and box overlap!")
	} else {
		println("No overlap.")
	}
}

// Helper function to run a single test case
func runTriBoxTest(t *testing.T, name string, boxCenter, boxHalfSize [3]float64, triverts [3][3]float64, expected bool) {
	result := TriBoxOverlap(&boxCenter, &boxHalfSize, &triverts)
	if result != expected {
		t.Errorf("%s → Expected: %v, Got: %v", name, expected, result)
	}
}

// Test Case: Triangle fully inside the box
func TestTriangleInsideBox(t *testing.T) {
	runTriBoxTest(t, "Triangle fully inside the box",
		[3]float64{0, 0, 0}, [3]float64{1, 1, 1},
		[3][3]float64{
			{0.1, 0.1, 0.1},
			{-0.1, -0.1, 0.1},
			{0.1, -0.1, -0.1},
		}, true)
}

// Test Case: Triangle completely outside the box
func TestTriangleOutsideBox(t *testing.T) {
	runTriBoxTest(t, "Triangle completely outside the box",
		[3]float64{0, 0, 0}, [3]float64{1, 1, 1},
		[3][3]float64{
			{2.0, 2.0, 2.0},
			{2.5, 2.5, 2.0},
			{2.0, 2.5, 2.5},
		}, false)
}

// Test Case: Triangle intersecting the box
func TestTriangleIntersectingBox(t *testing.T) {
	runTriBoxTest(t, "Triangle intersecting the box",
		[3]float64{0, 0, 0}, [3]float64{1, 1, 1},
		[3][3]float64{
			{1.5, 0.5, 0.5},
			{0.5, -0.5, -0.5},
			{-0.5, 0.5, -0.5},
		}, true)
}

// Test Case: Triangle lying exactly on a face of the box
func TestTriangleOnFaceOfBox(t *testing.T) {
	runTriBoxTest(t, "Triangle lying exactly on a face of the box",
		[3]float64{0, 0, 0}, [3]float64{1, 1, 1},
		[3][3]float64{
			{-1.0, 0.5, 0.5},
			{-1.0, -0.5, -0.5},
			{-1.0, 0.5, -0.5},
		}, true)
}

// Test Case: Triangle with one vertex inside and two outside
func TestTriangleOneVertexInside(t *testing.T) {
	runTriBoxTest(t, "Triangle with one vertex inside and two outside",
		[3]float64{0, 0, 0}, [3]float64{1, 1, 1},
		[3][3]float64{
			{0.5, 0.5, 0.5}, // Inside
			{2.0, 2.0, 2.0}, // Outside
			{2.5, 2.5, 2.5}, // Outside
		}, true)
}

// Test Case: Degenerate triangle (all points are the same)
func TestDegenerateTriangleInside(t *testing.T) {
	runTriBoxTest(t, "Degenerate triangle (all points the same)",
		[3]float64{0, 0, 0}, [3]float64{1, 1, 1},
		[3][3]float64{
			{0.5, 0.5, 0.5},
			{0.5, 0.5, 0.5},
			{0.5, 0.5, 0.5},
		}, true) // Expect true since it is inside
}

// Test Case: Degenerate triangle (all points are the same)
func TestDegenerateTriangleOutside(t *testing.T) {
	runTriBoxTest(t, "Degenerate triangle (all points the same)",
		[3]float64{0, 0, 0}, [3]float64{1, 1, 1},
		[3][3]float64{
			{1.1, 1.1, 1.1},
			{1.1, 1.1, 1.1},
			{1.1, 1.1, 1.1},
		}, false) // Expect false since it is outside
}

// Test Case: Triangle covering the entire box
func TestTriangleCoveringBox(t *testing.T) {
	runTriBoxTest(t, "Triangle covering the entire box",
		[3]float64{0, 0, 0}, [3]float64{1, 1, 1},
		[3][3]float64{
			{-2.0, -2.0, 0},
			{2.0, -2.0, 0},
			{0, 2.0, 0},
		}, true)
}

// Test Case: Triangle exactly at the edge of the box
func TestTriangleAtBoxEdge(t *testing.T) {
	runTriBoxTest(t, "Triangle exactly at the edge of the box",
		[3]float64{0, 0, 0}, [3]float64{1, 1, 1},
		[3][3]float64{
			{1.0, 0.5, 0.5},
			{1.0, -0.5, -0.5},
			{1.0, 0.5, -0.5},
		}, true)
}

// Test Case: Large triangle that should definitely intersect
func TestLargeTriangleIntersectingBox(t *testing.T) {
	runTriBoxTest(t, "Large triangle intersecting the box",
		[3]float64{0, 0, 0}, [3]float64{1, 1, 1},
		[3][3]float64{
			{-5.0, -5.0, -5.0},
			{5.0, 0, 0},
			{0, 5.0, 0},
		}, true)
}
