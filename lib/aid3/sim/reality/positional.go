package reality

import "github.com/timleecasey/stllib/lib/aid3/sim/tooling"

const (
	TRANS_ROW = 3
	TRANS_X   = 0
	TRANS_Y   = 1
	TRANS_Z   = 2

	X = 0
	Y = 1
	Z = 2
)

// Affine
// A 4x4 matrix [row, col] for 3d affines
type Affine struct {
	m [][]float64
}

type Velocity struct {
	X, Y, Z float64
}

func Still() *Velocity {
	return &Velocity{
		X: 0,
		Y: 0,
		Z: 0,
	}
}

func Identity() *Affine {
	m := make([][]float64, 4)

	for i := range 4 {
		m[i] = make([]float64, 4)
		for j := range 4 {
			if i == j {
				m[i][j] = 1
			} else {
				m[i][j] = 0
			}
		}
	}

	return &Affine{
		m: m,
	}
}

func Translate(x float64, y float64, z float64) *Affine {
	id := Identity()
	id.m[TRANS_ROW][TRANS_X] = x
	id.m[TRANS_ROW][TRANS_Y] = y
	id.m[TRANS_ROW][TRANS_Z] = z
	return id
}

func (a *Affine) MultiplyPoint(p *tooling.Point) *tooling.Point {
	var ret tooling.Point
	ret.X = a.m[0][0]*p.X + a.m[0][1]*p.Y + a.m[0][2]*p.Z + a.m[0][3]
	ret.Y = a.m[1][0]*p.X + a.m[1][1]*p.Y + a.m[1][2]*p.Z + a.m[1][3]
	ret.Z = a.m[2][0]*p.X + a.m[2][1]*p.Y + a.m[2][2]*p.Z + a.m[2][3]
	W := a.m[3][0]*p.X + a.m[3][1]*p.Y + a.m[3][2]*p.Z + a.m[3][3]

	ret.X = ret.X / W
	ret.Y = ret.Y / W
	ret.Z = ret.Z / W

	return &ret
}

func (a *Affine) MoveHeadBy(h tooling.Head) {
	cur := h.Pos()
	newPt := a.MultiplyPoint(cur)
	h.MoveTo(newPt)
}
