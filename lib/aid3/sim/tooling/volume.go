package tooling

import "math"

// https://github.com/hschendel/stl/blob/master/solid.go
// https://github.com/github/gh-skyline/blob/main/internal/stl/geometry/shapes.go
// https://www.reddit.com/r/golang/comments/az81nt/module_for_reading_and_writing_
type Volume interface {
	BoundingBox() *Point
	Subtract(mesh *Mesh)
	AddTriangle(p1 *Point, p2 *Point, p3 *Point)
}

func MakeVolume(fr *Point, to *Point) Volume {
	ret := &Mesh{
		shape: nil,
		bbMin: &Point{0, 0, 0},
		bbMax: &Point{0, 0, 0},
	}
	ret.bbMin.X = math.Min(fr.X, to.X)
	ret.bbMin.Y = math.Min(fr.Y, to.Y)
	ret.bbMin.Z = math.Min(fr.Z, to.Z)

	ret.bbMax.X = math.Max(fr.X, to.X)
	ret.bbMax.Y = math.Max(fr.Y, to.Y)
	ret.bbMax.Z = math.Max(fr.Z, to.Z)

	ret.offset = &Point{}

	return ret
}

// triangle
// clockwise definition of a triangle
type triangle struct {
	pts    []*Point
	normal *Point
	next   *triangle
}

type Mesh struct {
	// the shape which is currently representing the mesh
	shape *triangle
	// bbMin contains the minimal points for each axis
	bbMin *Point
	// bbMax contains the maximal points for each axis
	bbMax *Point
	// offset says to shift the shape by the point
	offset *Point
}

func (m *Mesh) BoundingBox() *Point {
	return &Point{
		X: m.bbMax.X - m.bbMin.X,
		Y: m.bbMax.Y - m.bbMin.Y,
		Z: m.bbMax.Z - m.bbMin.Z,
	}
}

// Subtract
// Given a volume, subtract the mesh from the volume.
// Intersect, then invert the normal/CC rule, making the
// outside of the subtracted shape the inside, then merge
// in the triangles.
func (m *Mesh) Subtract(vol *Mesh) {

}

// AddTriangle
// Extend the shape by this triangle
func (m *Mesh) AddTriangle(p1 *Point, p2 *Point, p3 *Point) {

}

func (m *Mesh) TranslateTo(p *Point) *Mesh {
	return nil
}
