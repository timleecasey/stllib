package tooling

type Material interface {
	Hardness() float64
	Stiffness() float64
	Volume() Volume
}

type Wood struct {
	hardness    float64
	stiffness   float64
	startVolume Volume
}

func MakeWood(dim float64) Material {
	return &Wood{
		hardness:    0.0,
		stiffness:   0.0,
		startVolume: MakeVolume(&Point{-dim, -dim, -dim}, &Point{dim, dim, dim}),
	}
}

func (w *Wood) Hardness() float64 {
	return w.hardness
}

func (w *Wood) Stiffness() float64 {
	return w.stiffness
}

func (w *Wood) Volume() Volume {
	return w.startVolume
}
