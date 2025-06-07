package threed

import "github.com/timleecasey/stllib/lib/aid3/sim/tooling"

const (
	In Side = iota
	Out
	Boundary
)

//// a 3d point
//type Point struct {
//	X float64
//	Y float64
//	Z float64
//}

type Side int

// Represents the minimal and maximal point of 3d volume.
type Dim struct {
	From tooling.Point
	To   tooling.Point
}
