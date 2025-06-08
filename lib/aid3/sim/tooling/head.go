package tooling

type Head interface {
	Pos() *Point
	MoveTo(p *Point)
	//MoveBy(a *reality.Affine)
}
