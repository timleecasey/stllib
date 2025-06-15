package sim

import (
	"github.com/timleecasey/stllib/lib/aid3/gcode"
	"github.com/timleecasey/stllib/lib/aid3/sim/tooling"
	"math"
)

func cmdLinear(s *Sim, cn *gcode.CmdNode) {
	curPt := s.ToolHead.Pos()
	curFeedRate := s.Tool.FeedRate() // mm/s?
	coords := cn.Cmd.Coords()
	toPt := CmdToXYZ(coords, curPt)
	slice := s.TimeSlice // s
	distPerSlice := curFeedRate * slice
	//
	// The x,y,z diff over the time slice
	//
	diffPt := &tooling.Point{
		X: toPt.X - curPt.X,
		Y: toPt.Y - curPt.Y,
		Z: toPt.Z - curPt.Z,
	}
	dist := math.Sqrt(math.Pow(toPt.X, 2) + math.Pow(toPt.Y, 2) + math.Pow(toPt.Z, 2))
	numIntersMoving := (dist / s.Tool.FeedRate()) / s.TimeSlice // (mm / (mm/s) -> s) / s -> count

	if dist != 0 {
		diffPt.X = diffPt.X / numIntersMoving
		diffPt.Y = diffPt.Y / numIntersMoving
		diffPt.Z = diffPt.Z / numIntersMoving

		// The affine to apply per iteration
		affine := tooling.Translate(
			diffPt.X,
			diffPt.Y,
			diffPt.Z)

		if distPerSlice > 0 {
			//
			// This may not be fine enough for a simulation when
			// there is deflection and heat, among others.
			// This is directly calculating an incremental jump in the
			// simulation time slice, vs calculating a step to a tolerance
			// and then morking points when past a time slice.
			//
			// Also, there is accumulated timeslice errors between linear moves.
			// A partial time slice is not taken into account at the start of the
			// next move.  (Diff of 0 is always used, vs diff of previous)
			//
			runLinearAffine(s, affine, toPt, diffPt)
		}
	}

}
