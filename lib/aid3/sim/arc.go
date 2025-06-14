package sim

import (
	"github.com/timleecasey/stllib/lib/aid3/gcode"
	"github.com/timleecasey/stllib/lib/aid3/sim/tooling"
	"log"
	"math"
)

// Format 1:
// (G17) G02/03 X__ Y__ I__ J__ F__
// (G18) G02/03 X__ Z__ I__ K__ F__
// (G19) G02/03 Y__ Z__ J__ K__ F__
//
// I,J,K specify the current plane analogous to X,Y,Z
// Only one of the three planes is used at any one time
// X/Y/Z I/J -> xy plane and z == k
// X/Y/Z I/K -> xz plane and y == j
// X/Y/Z J/K -> yz plane and x == i
//

// Format 2:
// (G17)G02/03 X__ Y__ R__ F__
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

const (
	ClockWise        = 1.0
	CounterClockWise = -1.0
)

func cmdCwArch(s *Sim, cn *gcode.CmdNode) {

	coords := cn.Cmd.Coords()
	plane := s.Tool.Plane()

	gcodeForm := figureOutArcCmdFormat(coords)
	planeConst := establishPlaneConstant(plane)

	pointsAlongCurve(s, coords, gcodeForm, planeConst, ClockWise)
}

func cmdCcwArch(s *Sim, cn *gcode.CmdNode) {
	coords := cn.Cmd.Coords()
	plane := s.Tool.Plane()

	gcodeForm := figureOutArcCmdFormat(coords)
	planeConst := establishPlaneConstant(plane)

	pointsAlongCurve(s, coords, gcodeForm, planeConst, CounterClockWise)
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

// Given the start and coords, find the start and end angles and the angle increment to match the given
// tolerance.
func pointsAlongCurve(s *Sim, c *gcode.Coords, cmdForm int, planeConst int, dir float64) {
	outOfTolerance := true
	var midAngle float64
	var radius float64
	var arc *tooling.Point
	fr := s.ToolHead.Pos()
	tolerance := s.Tolerance

	curFeedRate := s.Tool.FeedRate() // mm/s?
	slice := s.TimeSlice             // s
	distPerSlice := curFeedRate * slice

	to := &tooling.Point{
		X: c.X,
		Y: c.Y,
		Z: c.Z,
	}

	// abs or offset from cur pos?
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
		if c.R < 0 {
			//big = true
		}
		break
	case formatCoords:
		rX := fr.X - center.X
		rY := fr.Y - center.Y
		rZ := fr.Z - center.Z
		radius = math.Sqrt(rX*rX + rY*rY + rZ*rZ)
		break
	}

	cnt := 1000

	diffPt := &tooling.Point{}
	diffPt.X = to.X
	diffPt.Y = to.Y
	diffPt.Z = to.Z

	if debugArc {
		log.Printf("ARC fr %v to %v center %v", fr, to, center)
	}

	startAng := math.Atan2(fr.Y-center.Y, fr.X-center.X)
	endAng := math.Atan2(to.Y-center.Y, to.X-center.X)
	conv := 180 / math.Pi

	if debugArc {
		log.Printf("ARC START %3.2f END %3.2f", startAng*conv, endAng*conv)
	}

	for outOfTolerance {
		cnt--
		if cnt <= 0 {
			break
		}

		midPt := tooling.MidPoint(fr, diffPt)

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
			midAngle = math.Atan2(midPt.X-center.X, midPt.Z-center.Z)
			arc = tooling.PointAt(center, radius, midAngle)
			break
		case planeYConst:
			midAngle = math.Atan2(midPt.Z-center.Z, midPt.X-center.X)
			arc = tooling.PointAt(center, radius, midAngle)
			break

		case planeZConst:
			midAngle = math.Atan2(midPt.Y-center.Y, midPt.X-center.X)
			arc = tooling.PointAt(center, radius, midAngle)
			break
		}

		arcSegmentDiff := arc.Dist(diffPt)

		if debugArc {
			log.Printf("ARC TOL @%v DIFF: %3.4f ANGLE: %3.2f\n", cnt, arcSegmentDiff, midAngle*conv)
			log.Printf("ARC     DIFF: %v <-- ARC: %v\n", diffPt, arc)
		}

		diffPt.X = arc.X
		diffPt.Y = arc.Y
		diffPt.Z = arc.Z

		if arcSegmentDiff > tolerance {
			outOfTolerance = true
		} else {
			outOfTolerance = false
		}
	}

	midAngle = math.Abs(midAngle - startAng)
	approx := math.Abs(startAng-endAng) / math.Abs(midAngle)
	if debugArc {
		log.Printf("ARC angle %3.6f ~# %v", midAngle, approx)
	}

	cnt = 0
	dist := 0.0
	prev := fr
	var p *tooling.Point
	posted := 0

	//
	// This is a bit excessive.
	// We know the arc length, we know the angle increment
	// We could calculate the number of arc segments per
	// distPerSlice and just post those points directly.
	//
	for curAng := startAng; ; curAng += dir * midAngle {
		if dir > 0 && curAng > endAng {
			break
		}
		if dir < 0 && curAng < endAng {
			break
		}

		cnt++

		switch planeConst {
		case planeXConst:
			y := center.Y + radius*math.Cos(curAng)
			z := center.Z + radius*math.Sin(curAng)
			p = &tooling.Point{X: fr.X, Y: y, Z: z}
			break

		case planeYConst:
			x := center.X + radius*math.Cos(curAng)
			z := center.Z + radius*math.Sin(curAng)
			p = &tooling.Point{X: x, Y: fr.Y, Z: z}
			break

		case planeZConst:
			x := center.X + radius*math.Cos(curAng)
			y := center.Y + radius*math.Sin(curAng)
			p = &tooling.Point{X: x, Y: y, Z: fr.Z}
			break

		}

		dist += prev.Dist(p)
		if dist > distPerSlice {
			s.ToolHead.MoveTo(p)
			dist = 0.0
			posted++
		}
		prev = p
	}

	if dist > 0 {
		s.ToolHead.MoveTo(p)
		posted++
	}

	if debugPts {
		log.Printf("ARC PTs simmed %3v posted %v DIST %v", cnt, posted, distPerSlice)
	}

}
