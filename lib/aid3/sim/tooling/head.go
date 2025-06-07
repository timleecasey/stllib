package tooling

type Head interface {
	Pos() *Point
	MoveTo(p *Point)
}
