package tdm

import (
	"github.com/timleecasey/stllib/lib/aid3/sim/tooling"
	"github.com/timleecasey/stllib/lib/stl"
	"github.com/timleecasey/stllib/lib/threed"
)

type Volume interface {
	Bounds() *threed.Dim
	Sidedness(p tooling.Point, epsilon float64) threed.Side
	Intersect(stl *stl.Model)
}
