package tooling

type Simple3d struct {
	head Head
	zero *Point
}
type SimpleHead struct {
	pos *Point
}

func BuildCnc() Cnc {
	ret := &Simple3d{
		zero: &Point{},
	}

	head := &SimpleHead{
		pos: &Point{0, 0, 0},
	}

	ret.head = head

	return ret
}

//
// Cnc
//

func (s3d Simple3d) Axis() []int {
	ret := make([]int, 1)
	ret[0] = 3
	return ret
}

func (s3d Simple3d) ZeroPoint() *Point {
	return s3d.zero
}

func (s3d Simple3d) Head() Head {
	return s3d.head
}

// Head
func (h SimpleHead) Pos() *Point {
	return h.pos
}

func (h SimpleHead) MoveTo(p *Point) {
	h.pos = p
}
