package tdm

import "github.com/timleecasey/stllib/lib/stl"

// a 3d point
type Point struct {
	X float64
	Y float64
	Z float64
}

type Side int

const (
	In Side = iota
	Out
	Boundary
)

// Represents the minimal and maximal point of 3d volume.
type Dim struct {
	From Point
	To   Point
}

type Volume interface {
	Bounds() *Dim
	Sidedness(p *Point, epsilon float64) Side
	Intersect(stl *stl.Model)

	Lookup
}
