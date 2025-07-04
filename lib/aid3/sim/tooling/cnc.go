package tooling

import (
	"fmt"
	"math"
)

const (
	CncMilling = iota
	CncLathe
	CncGrinder
	CncDrill
	CncRouter
	CncLaserCutter
	CncPlasmaCutter
	CncWaterJetCutting
	CncElectricDischargeMachines // (EDMs), but not the music
)

type Point struct {
	X float64
	Y float64
	Z float64
}

func (p *Point) String() string {
	return fmt.Sprintf("X:%4.3f, Y:%4.3f, Z:%4.3f", p.X, p.Y, p.Z)
}

func (p *Point) Dist(to *Point) float64 {
	diffX := p.X - to.X
	diffY := p.Y - to.Y
	diffZ := p.Z - to.Z
	return math.Sqrt((diffX * diffX) + (diffY * diffY) + (diffZ * diffZ))
}

func PointAt(center *Point, radius float64, angle float64) *Point {
	ret := &Point{
		X: center.X + radius*math.Cos(angle),
		Y: center.Y + radius*math.Sin(angle),
		Z: center.Z,
	}

	return ret
}

func MidPoint(fr *Point, to *Point) *Point {
	return &Point{
		X: (fr.X + to.X) / 2,
		Y: (fr.Y + to.Y) / 2,
		Z: (fr.Z + to.Z) / 2,
	}
}

const (
	FEED_NONE = iota
	FEED_INVERSE_TIME
	FEED_PER_MINUTE
	FEED_PER_REVOLUTION
)

const (
	PLANE_NONE = iota
	PLANE_XY
	PLANE_XZ
	PLANE_YZ
)

const (
	UNIT_NONE = iota
	UNIT_INCH
	UNIT_MM
)

type Cnc interface {
	Axis() []int
	ZeroPoint() *Point
	Head() Head
	FeedRate() float64
	AssignFeedRate(f float64)
	FastFeedRate() float64
	FeedMode(mode int)
	SpindleSpeed(speed int64)
	ToolChangeTo(tool int64)
	SelectPlane(plane int)
	Plane() int
	Reset()
	Units(units int)
	WorkVolume() Volume
	Material() Material
}
