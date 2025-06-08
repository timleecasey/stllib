package sim

import (
	"fmt"
	"github.com/timleecasey/stllib/lib/aid3/gcode"
	"github.com/timleecasey/stllib/lib/aid3/sim/reality"
	"github.com/timleecasey/stllib/lib/aid3/sim/tooling"
	"math"
)

// Sim
// TimeSlice is the time unit for running the sim
// Head is the current position
// Velocity is the Velocity of tool head
// The unit is meters, so Velocity is m/s and Head position is measured in meter with respect to the tool zero point
type Sim struct {
	TimeSlice float64
	Steps     uint
	Tool      tooling.Cnc
	ToolHead  tooling.Head
	Velocity  *reality.Velocity
}

func (s *Sim) Start() {
	s.TimeSlice = 0.001
	s.Steps = 10

	tool := tooling.BuildCnc()
	head := tool.Head()

	s.ToolHead = head
	s.Tool = tool

	zero := tool.ZeroPoint()
	head.MoveTo(zero)

	s.Velocity = reality.Still()
}

func (s *Sim) Run(tree *gcode.ParseTree) {

	head := s.ToolHead

	for i := range s.Steps {
		fmt.Printf("%v @ %v\n", i, head.Pos())
		diff := reality.Still()
		diff.X = s.Velocity.X * s.TimeSlice
		diff.Y = s.Velocity.Y * s.TimeSlice
		diff.Z = s.Velocity.Z * s.TimeSlice
		affine := reality.Translate(diff.X, diff.Y, diff.Z)
		newPos := affine.MultiplyPoint(head.Pos())
		head.MoveTo(newPos)
	}

	affine := reality.Identity()
	timeStep := s.TimeSlice
	tree.TraverseCmds(func(cn *gcode.CmdNode) error {
		if cn.Cmd.Coords().F != 0 {
			s.Tool.AssignFeedRate(cn.Cmd.Coords().F)
		}
		switch cn.Cmd.CmdType() {
		case gcode.CMD_LINEAR:
			curPt := s.ToolHead.Pos()
			curFeedRate := s.Tool.FeedRate() // mm/s?
			coords := cn.Cmd.Coords()
			toPt := CmdToXYZ(coords, curPt)
			slice := timeStep // s
			distPerSlice := curFeedRate * slice
			//
			// The x,y,z diff over the time slice
			//
			diffPt := &tooling.Point{
				X: toPt.X - curPt.X,
				Y: toPt.Y - curPt.Y,
				Z: toPt.Z - curPt.Z,
			}

			affine = reality.Translate(distPerSlice, distPerSlice, distPerSlice)

			runLinearAffine(s, affine, toPt, diffPt, distPerSlice)
		}
		return nil
	})
}

func notClipped(to *tooling.Point, fr *tooling.Point, dist float64) bool {
	return math.Abs(to.X-fr.X) > math.Abs(dist) &&
		math.Abs(to.Y-fr.Y) > math.Abs(dist) &&
		math.Abs(to.Z-fr.Z) > math.Abs(dist)
}

func runLinearAffine(s *Sim, affine *reality.Affine, toPt *tooling.Point, diffPt *tooling.Point, distPerSlice float64) {
	for notClipped(s.ToolHead.Pos(), toPt, distPerSlice) {
		h := s.ToolHead
		affine.MoveHeadBy(h)
	}
	s.ToolHead.MoveTo(toPt)
}

func CmdToXYZ(c *gcode.Coords, curPt *tooling.Point) *tooling.Point {
	ret := &tooling.Point{
		X: c.X,
		Y: c.Y,
		Z: c.Z,
	}

	if ret.X == 0 {
		ret.X = curPt.X
	}
	if ret.Y == 0 {
		ret.Y = curPt.Y
	}
	if ret.Z == 0 {
		ret.Z = curPt.Z
	}

	return ret
}
