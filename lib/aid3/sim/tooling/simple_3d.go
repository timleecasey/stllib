package tooling

type Simple3d struct {
	head         Head
	zero         *Point
	feed         float64
	feedMode     int
	spindleSpeed int64
	curTool      int
	plane        int
	units        int
}

type SimpleHead struct {
	pos  *Point
	path []*Point
}

func BuildCnc() Cnc {
	ret := &Simple3d{}

	head := &SimpleHead{
		pos:  &Point{0, 0, 0},
		path: make([]*Point, 0),
	}

	ret.head = head

	return ret
}

//
// Cnc
//

func (s3d *Simple3d) FeedRate() float64 {
	return s3d.feed
}

func (s3d *Simple3d) AssignFeedRate(f float64) {
	s3d.feed = f
}

func (s3d *Simple3d) Axis() []int {
	ret := make([]int, 1)
	ret[0] = 3
	return ret
}

func (s3d *Simple3d) ZeroPoint() *Point {
	return s3d.zero
}

func (s3d *Simple3d) Head() Head {
	return s3d.head
}

func (s3d *Simple3d) FastFeedRate() float64 {
	return 1000
}

func (s3d *Simple3d) FeedMode(mode int) {
	s3d.feedMode = mode
}

func (s3d *Simple3d) SpindleSpeed(speed int64) {
	s3d.spindleSpeed = speed
}

func (s3d *Simple3d) ToolChangeTo(tool int) {
	s3d.curTool = tool
}

func (s3d *Simple3d) SelectPlane(plane int) {
	s3d.plane = plane
}

func (s3d *Simple3d) Reset() {
	s3d.zero = &Point{}
	s3d.plane = PLANE_XY
	s3d.spindleSpeed = 0
	s3d.feedMode = FEED_PER_MINUTE
	s3d.feed = s3d.FastFeedRate()
	s3d.head.MoveTo(s3d.zero)
	s3d.units = UNIT_MM
}

func (s3d *Simple3d) Units(units int) {
	s3d.units = units
}

// Head
func (h *SimpleHead) Pos() *Point {
	return h.pos
}

func (h *SimpleHead) MoveTo(p *Point) {
	h.pos = p
	h.path = append(h.path, p)
}

func (h *SimpleHead) Path(f func(p *Point)) {
	for i := range h.path {
		f(h.path[i])
	}
}

//func (h SimpleHead) MoveBy(a *reality.Affine) {
//	newPos := a.MultiplyPoint(h.pos)
//	h.pos = newPos
//}
