package sim

import (
	"fmt"
	"github.com/timleecasey/stllib/lib/aid3/gcode"
	"github.com/timleecasey/stllib/lib/aid3/sim/reality"
	"github.com/timleecasey/stllib/lib/aid3/sim/tooling"
	"log"
	"math"
	"os"
)

var debugMove = true

// Sim
// TimeSlice is the time unit increment for running the sim
// Head is the current position
// The unit is meters, so Velocity is m/s and Head position is measured in meter with respect to the tool zero point
type Sim struct {
	TimeSlice float64
	Tool      tooling.Cnc
	ToolHead  tooling.Head
}

func (s *Sim) Start() {
	s.TimeSlice = 0.001

	tool := tooling.BuildCnc()
	head := tool.Head()
	tool.Reset()

	s.ToolHead = head
	s.Tool = tool
}

func (s *Sim) Run(tree *gcode.ParseTree) {

	log.Printf("Start %v\n", s.Tool.Head().Pos())

	cnt := 0

	tree.TraverseCmds(func(cn *gcode.CmdNode) error {
		var err error
		err = nil

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

		case gcode.CMD_TOOL_CHANGE:
			err = cmdToolChange(s, cn)
			cnt++
			break

		case gcode.CMD_SPINDLE_SPEED:
			err = cmdSpindleSpeed(s, cn)
			cnt++
			break
		case gcode.CMD_SPINDLE_OFF:
			s.Tool.SpindleSpeed(0)
			cnt++
			break

		case gcode.CMD_FEED_PER_MIN_MODE:
			s.Tool.FeedMode(tooling.FEED_PER_MINUTE)
			cnt++
			break
		case gcode.CMD_INVERSE_TIME_FEED: // feed is time instead of rate
			s.Tool.FeedMode(tooling.FEED_INVERSE_TIME)
			cnt++
			break
		case gcode.CMD_FEED_PER_REVOLUTION:
			s.Tool.FeedMode(tooling.FEED_PER_REVOLUTION)
			cnt++
			break

		case gcode.CMD_PLANE_XY:
			s.Tool.SelectPlane(tooling.PLANE_XY)
			cnt++
			break
		case gcode.CMD_PLANE_XZ:
			s.Tool.SelectPlane(tooling.PLANE_XZ)
			cnt++
			break
		case gcode.CMD_PLANE_YZ:
			s.Tool.SelectPlane(tooling.PLANE_YZ)
			cnt++
			break
		}
		if debugMove {
			log.Printf("After %v %v F: %v\n", cn.Cmd.Src(), s.Tool.Head().Pos(), s.Tool.FeedRate())
		}
		return err
	})

	writePathPoints(s.ToolHead)
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
		h := s.ToolHead
		moveHeadBy(h, affine)
	}
	s.ToolHead.MoveTo(toPt)
}

func moveHeadBy(h tooling.Head, affine *reality.Affine) {
	cur := h.Pos()
	newPt := affine.MultiplyPoint(cur)
	h.MoveTo(newPt)
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

func writePathPoints(h tooling.Head) {

	if f, err := os.Create("path.gcode"); err == nil {
		defer f.Close()
		h.Path(func(p *tooling.Point) {
			ptStr := fmt.Sprintf("G1 %v %v %v\n", p.X, p.Y, p.Z)
			_, err = f.WriteString(ptStr)
		})
	} else {
		log.Printf("Could not write path.gcode : %v", err)
	}
}
