package tooling

// Simple3d
// This simulates a 3-axis head tool
// Each increment in the simulation
// is a position change, represented
// in the tool head as a list of visited
// points.
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
	pos    *Point
	path   []*Point
	curVel *Velocity
}

func BuildCnc() Cnc {
	ret := &Simple3d{}

	head := &SimpleHead{
		pos:    &Point{0, 0, 0},
		path:   make([]*Point, 0),
		curVel: Still(),
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
func (s3d *Simple3d) Plane() int {
	return s3d.plane
}

func (s3d *Simple3d) Reset() {
	s3d.zero = &Point{}
	s3d.plane = PLANE_XY
	s3d.spindleSpeed = 0
	s3d.feedMode = FEED_PER_MINUTE
	s3d.feed = s3d.FastFeedRate()
	s3d.units = UNIT_MM
	s3d.head.Reset(s3d.zero)
}

func (s3d *Simple3d) Units(units int) {
	s3d.units = units
}

// Head
func (h *SimpleHead) Pos() *Point {
	return h.pos
}

func (h *SimpleHead) MoveTo(p *Point) {
	h.MarkVelocity(h.pos, p)
	h.pos = p
	h.path = append(h.path, p)
}

func (h *SimpleHead) PointCount() int {
	return len(h.path)
}

func (h *SimpleHead) Path(f func(p *Point)) {
	for i := range h.path {
		f(h.path[i])
	}
}

func (h *SimpleHead) CurVelocity() *Velocity {
	return h.curVel
}

func (h *SimpleHead) Reset(zero *Point) {
	h.MoveTo(zero)
	h.curVel = Still()
}

func (h *SimpleHead) MarkVelocity(fr *Point, to *Point) {
	h.curVel.X = to.X - fr.X
	h.curVel.Y = to.Y - fr.Y
	h.curVel.Z = to.Z - fr.Z
}
