package sim

import (
	"github.com/timleecasey/stllib/lib/aid3/gcode"
	"github.com/timleecasey/stllib/lib/aid3/sim/tooling"
	"log"
	"math"
)

// Format 1: (G17) G02/03 X__ Y__ I__ J__ F__
// (G18) G02/03 X__ Z__ I__ K__ F__
// (G19) G02/03 Y__ Z__ J__ K__ F__
//
// I,J,K specify the current plane analogous to X,Y,Z
// Only one of the three planes is used at any one time
// X/Y/Z I/J -> xy plane and z == k
// X/Y/Z I/K -> xz plane and y == j
// X/Y/Z J/K -> yz plane and x == i
//

// Format 2: (G17)G02/03 X__ Y__ R__ F__
// (G18)G02/03 X__ Z__ R__ F__
// (G19)G02/03 Y__ Z__ R__ F__
//
// Same game.
// X/Y/Z X/Y/R xy plane z is constant
// X/Y/Z X/Z/R xz plane y is constant
// X/Y/Z Y/Z/R yz plane x is constant
//

const (
	planeUnknownConst = 0
	planeXConst       = 1
	planeYConst       = 2
	planeZConst       = 3

	formatCoords = 1
	formatRadius = 2
)

func cmdCwArch(s *Sim, cn *gcode.CmdNode) {

	coords := cn.Cmd.Coords()
	plane := s.Tool.Plane()

	gcodeForm := figureOutArcCmdFormat(coords)
	planeConst := establishPlaneConstant(plane)

	angleStep := angleForTolerance(s, coords, gcodeForm, planeConst)
	log.Printf("ANG STEP %v\n", angleStep)

}

func cmdCcwArch(s *Sim, cn *gcode.CmdNode) {

}

func establishPlaneConstant(plane int) int {
	switch plane {
	case tooling.PLANE_XY:
		return planeZConst
	case tooling.PLANE_XZ:
		return planeYConst
	case tooling.PLANE_YZ:
		return planeXConst
	}
	return planeUnknownConst
}

// figureOutArchMode
// return 1 for form 1 and return 2 for form 2.
// format 2 uses radians, as a sweep amount
// format 1 uses coords, and a sweep is calculated.
func figureOutArcCmdFormat(c *gcode.Coords) int {
	if c.R != 0 {
		return formatRadius
	} else {
		return formatCoords
	}
}

func cwArcIjkInXy(s *Sim, start [2]float64, cn *gcode.CmdNode, endAngle float64, clockwise bool, tolerance float64) {
	coords := cn.Cmd.Coords()
	st := s.ToolHead.Pos()
	center := tooling.Point{}
	center.X = st.X + coords.I
	center.Y = st.Y + coords.J
	radius := math.Sqrt(coords.I*coords.I + coords.J*coords.J)

	//short := true
	//if coords.R < 0 {
	//	short := false
	//}

	// Convert the start angle to radians
	startAngle := math.Atan2(coords.J-center.Y, coords.I-center.X) // Calculate the starting angle in radians

	// Adjust the direction based on G02/G03 (clockwise or counterclockwise)
	// G02 is clockwise, G03 is counterclockwise
	direction := 1
	if !clockwise {
		direction = -1
	}

	angleIncrement := math.Pi / 180.0 // Start with 1 degree increment
	for angleIncrement*radius > tolerance {
		angleIncrement /= 2 // Reduce increment for more precision
	}

	// Generate points along the arc
	currentAngle := startAngle
	for currentAngle <= endAngle {
		// Calculate the new x and y position along the arc
		x := center.X + radius*math.Cos(currentAngle)
		y := center.Y + radius*math.Sin(currentAngle)
		s.ToolHead.MoveTo(&tooling.Point{X: x, Y: y, Z: st.Z})
		currentAngle += float64(direction) * angleIncrement
	}

}

func angleForTolerance(s *Sim, c *gcode.Coords, cmdForm int, planeConst int) float64 {
	//func angleForTolerance(center *tooling.Point, fr *tooling.Point, to *tooling.Point, radius float64, tolerance float64) float64 {
	outOfTolerance := true
	var midAngle float64
	var radius float64
	var arc *tooling.Point
	fr := s.ToolHead.Pos()
	tolerance := s.Tolerance

	to := &tooling.Point{
		X: c.X,
		Y: c.Y,
		Z: c.Z,
	}

	center := &tooling.Point{}
	//center.X = fr.X + c.I
	//center.Y = fr.Y + c.J
	//center.Z = fr.Z + c.K
	center.X = c.I
	center.Y = c.J
	center.Z = c.K

	switch planeConst {
	case planeXConst:
		to.X = fr.X
		center.X = fr.X
		break
	case planeYConst:
		to.Y = fr.Y
		center.Y = fr.Y
		break
	case planeZConst:
		to.Z = fr.Z
		center.Z = fr.Z
		break
	}

	switch cmdForm {
	case formatRadius:
		radius = math.Abs(c.R)
		break
	case formatCoords:
		rX := math.Abs(to.X - fr.X)
		rY := math.Abs(to.Y - fr.Y)
		rZ := math.Abs(to.Z - fr.Z)
		radius = math.Sqrt(rX*rX + rY*rY + rZ*rZ)
		break
	}

	cnt := 1000

	segPt := &tooling.Point{}
	segPt.X = to.X
	segPt.Y = to.Y
	segPt.Z = to.Z

	for outOfTolerance {
		cnt--
		if cnt <= 0 {
			break
		}

		//
		// Given a middle angle and the point along the arc
		// Find the distance between the point on the arc and
		// the line segment estimating the curve.  When it goes
		// below the given tolerance, use that incremental angle.
		// 'max' is used to stop bugs from doing too much damange
		// during running.  That is dividing things by 1/2 1000 times
		// is probably always within tolerance.
		switch planeConst {
		case planeXConst:
			midAngle = math.Atan2(segPt.X-center.X, segPt.Z-center.Z)
			arc = &tooling.Point{
				Y: center.Y + radius*math.Cos(midAngle),
				Z: center.Z + radius*math.Sin(midAngle),
				X: center.X,
			}
			break
		case planeYConst:
			midAngle = math.Atan2(segPt.Z-center.Z, segPt.X-center.X)
			arc = &tooling.Point{
				X: center.X + radius*math.Cos(midAngle),
				Z: center.Z + radius*math.Sin(midAngle),
				Y: center.Y,
			}
			break

		case planeZConst:
			midAngle = math.Atan2(segPt.Y-center.Y, segPt.X-center.X)
			arc = &tooling.Point{
				X: center.X + radius*math.Cos(midAngle),
				Y: center.Y + radius*math.Sin(midAngle),
				Z: center.Z,
			}
			break
		}

		arcSegmentDiff := arc.Dist(segPt)

		if debugArc {
			log.Printf("TOL @%v SEG: %v ANGLE: %v\n", cnt, segPt, midAngle)
		}

		segPt.X = arc.X
		segPt.Y = arc.Y
		segPt.Z = arc.Z

		//// Cut the mid point in half again.
		//mid.X = (fr.X + segPt.X) / 2
		//mid.Y = (fr.Y + segPt.Y) / 2
		//mid.Z = (fr.Z + segPt.Z) / 2

		if arcSegmentDiff > tolerance {
			outOfTolerance = true
		} else {
			outOfTolerance = false
		}
	}

	return midAngle
}
