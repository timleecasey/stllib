package stl

import (
	"github.com/timleecasey/stllib/lib/threed"
	"log"
	"neilpa.me/go-stl"
	"os"
)

// a triangle within an stl file
type Trap struct {
	A, B, C threed.Point
	Normal  threed.Point
}

// visit the traps within the model
type TrapVisitor func(t *Trap)

// The representation of a STL model
type Model struct {
	Objs   *[]*Trap
	bounds *threed.Dim
}

func (m *Model) Bounds() *threed.Dim {
	return m.bounds
}

const (
	X_PT   = 0
	Y_PT   = 1
	Z_PT   = 2
	FIRST  = 0
	SECOND = 1
	THIRD  = 2
)

// Given a point, adjust the dimension to form a bounding box
func boundsOnPoint(d *threed.Dim, p *threed.Point) {
	// minimal point
	if p.X < d.From.X {
		d.From.X = p.X
	}
	if p.Y < d.From.Y {
		d.From.Y = p.Y
	}
	if p.Z < d.From.Z {
		d.From.Z = p.Z
	}

	// maximal point
	if p.X > d.To.X {
		d.To.X = p.X
	}
	if p.Y > d.To.Y {
		d.To.Y = p.Y
	}
	if p.Z > d.To.Z {
		d.To.Z = p.Z
	}

}

func LoadModel(nm string) (*Model, error) {
	ret := Model{}
	objectList := make([]*Trap, 10)
	ret.Objs = &objectList
	err := ret.openStl(nm)
	if err != nil {
		return nil, err
	}

	var bounds threed.Dim
	ret.traverse(func(t *Trap) {
		boundsOnPoint(&bounds, &t.A)
		boundsOnPoint(&bounds, &t.B)
		boundsOnPoint(&bounds, &t.C)
	})

	ret.bounds = &bounds

	return &ret, nil
}

func (m *Model) openStl(nm string) error {
	f, err := os.Open(nm)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer f.Close()

	mesh, err := stl.Decode(f)
	if err != nil {
		log.Fatal(err)
		return err
	}

	for _, face := range mesh.Faces {
		x := threed.Point{
			X: float64(face.Verts[X_PT][FIRST]),
			Y: float64(face.Verts[Y_PT][FIRST]),
			Z: float64(face.Verts[Z_PT][FIRST]),
		}
		y := threed.Point{
			X: float64(face.Verts[X_PT][SECOND]),
			Y: float64(face.Verts[Y_PT][SECOND]),
			Z: float64(face.Verts[Z_PT][SECOND]),
		}
		z := threed.Point{
			X: float64(face.Verts[X_PT][THIRD]),
			Y: float64(face.Verts[Y_PT][THIRD]),
			Z: float64(face.Verts[Z_PT][THIRD]),
		}

		normalPt := threed.Point{
			X: float64(face.Normal[X_PT]),
			Y: float64(face.Normal[Y_PT]),
			Z: float64(face.Normal[Z_PT]),
		}

		t := Trap{
			A:      x,
			B:      y,
			C:      z,
			Normal: normalPt,
		}
		ta := append(*m.Objs, &t)
		m.Objs = &ta
	}

	return nil
}

func (m *Model) traverse(v TrapVisitor) {
	for i := range *m.Objs {
		t := (*m.Objs)[i]
		v(t)
	}
}
