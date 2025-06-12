package tooling

type Head interface {
	Pos() *Point
	MoveTo(p *Point)
	Path(f func(p *Point))
	CurVelocity() *Velocity
	Reset(zero *Point)
}
