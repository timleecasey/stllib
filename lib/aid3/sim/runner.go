package sim

import (
	"github.com/timleecasey/stllib/lib/aid3/gcode"
	"github.com/timleecasey/stllib/lib/aid3/sim/reality"
	"github.com/timleecasey/stllib/lib/aid3/sim/tooling"
	"log"
	"math"
)

var debugMove = true

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
	//Velocity  *reality.Velocity
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
	s.Tool.AssignFeedRate(s.Tool.FastFeedRate())

	//s.Velocity = reality.Still()
}

func (s *Sim) Run(tree *gcode.ParseTree) {

	log.Printf("Start %v\n", s.Tool.Head().Pos())

	cnt := 0

	tree.TraverseCmds(func(cn *gcode.CmdNode) error {
		if cn.Cmd.Coords().F != 0 {
			s.Tool.AssignFeedRate(cn.Cmd.Coords().F)
		}
		switch cn.Cmd.CmdType() {
		case gcode.CMD_FAST:
			s.Tool.AssignFeedRate(s.Tool.FastFeedRate())
			cmdLinear(s, cn)
			cnt++
			break
		case gcode.CMD_LINEAR:
			cmdLinear(s, cn)
			cnt++
			break

		case gcode.CMD_CW_ARC:
			cmdCwArch(s, cn)
			cnt++
			break
		case gcode.CMD_CCW_ARC:
			cmdCcwArch(s, cn)
			cnt++
			break

		case gcode.CMD_SPINDLE_SPEED:
			cmdSpindleSpeed(s, cn)
			cnt++
			break
		case gcode.CMD_SPINDLE_OFF:
			s.Tool.SpindleSpeed(0)
			cnt++
			break

		case gcode.CMD_FEED_PER_MIN_MODE:
			s.Tool.FeedMode(tooling.FEED_PER_MINUTE)
			cnt++
		case gcode.CMD_INVERSE_TIME_FEED:
			s.Tool.FeedMode(tooling.FEED_INVERSE_TIME)
			cnt++
		case gcode.CMD_FEED_PER_REVOLUTION:
			s.Tool.FeedMode(tooling.FEED_PER_REVOLUTION)
			cnt++

		}
		if debugMove {
			log.Printf("After %v %v F: %v\n", cn.Cmd.Src(), s.Tool.Head().Pos(), s.Tool.FeedRate())
		}
		return nil
	})
	log.Printf("Ran %v commands\n", cnt)

}

func notClipped(to *tooling.Point, fr *tooling.Point, diffPt *tooling.Point) bool {
	// So long as we have some distance to move greater than the diffPoint, in at least one direction, then move
	return math.Abs(to.X-fr.X) > math.Abs(diffPt.X) ||
		math.Abs(to.Y-fr.Y) > math.Abs(diffPt.Y) ||
		math.Abs(to.Z-fr.Z) > math.Abs(diffPt.Z)
}

func runLinearAffine(s *Sim, affine *reality.Affine, toPt *tooling.Point, diffPt *tooling.Point) {
	for notClipped(s.ToolHead.Pos(), toPt, diffPt) {
		//log.Printf("LINEAR %v\n", s.ToolHead.Pos())
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
