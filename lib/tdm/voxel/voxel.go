package voxel

import (
	"github.com/timleecasey/stllib/lib/aid3/sim/tooling"
	"github.com/timleecasey/stllib/lib/stl"
	"github.com/timleecasey/stllib/lib/threed"
)

//type Voxel interface {
//	VolumeSize() uint
//}

type cube struct {
	side threed.Side
	dim  *threed.Dim
}

func (c *cube) Bounds() *threed.Dim {
	return nil
}
func (c *cube) Sidedness(p *tooling.Point) threed.Side {
	return threed.Out
}
func (c *cube) Intersect(stl *stl.Model) {
}

func makeCube(dim *threed.Dim) *cube {
	ret := cube{
		side: threed.Out,
		dim:  dim,
	}
	return &ret
}

type Voxel struct {
	size   uint
	cubes  *[][][]*cube // pointer to an array of pointers to cube
	bounds *threed.Dim
	model  *stl.Model
}

func MakeVoxel(size uint, stl *stl.Model) *Voxel {
	cubes := make([][][]*cube, size)
	bounds := stl.Bounds()

	diffX := (bounds.From.X - bounds.To.X) / float64(size)
	diffY := (bounds.From.Y - bounds.To.Y) / float64(size)
	diffZ := (bounds.From.Z - bounds.To.Z) / float64(size)

	init := tooling.Point{
		X: bounds.From.X,
		Y: bounds.From.Y,
		Z: bounds.From.Z,
	}

	ret := Voxel{
		size:   size,
		cubes:  &cubes,
		bounds: bounds,
		model:  stl,
	}

	var x, y, z uint
	for x = 0; x < size; x++ {
		for y = 0; y < size; y++ {
			for z = 0; z < size; z++ {
				cubeDim := threed.Dim{
					From: tooling.Point{
						X: init.X + float64(x)*diffX,
						Y: init.Y + float64(y)*diffY,
						Z: init.Z + float64(z)*diffZ,
					},
					To: tooling.Point{
						X: init.X + float64(x)*diffX + diffX,
						Y: init.Y + float64(y)*diffY + diffY,
						Z: init.Z + float64(z)*diffZ + diffZ,
					},
				}
				cubes[x][y][z] = makeCube(&cubeDim)
			}
		}
	}
	return &ret
}

func (v *Voxel) Bounds() *threed.Dim {
	return v.bounds
}

func (v *Voxel) Sidedness(p *tooling.Point, e float64) threed.Side {
	return threed.Out
}

func (v *Voxel) Intersect(stl *stl.Model) {

}
