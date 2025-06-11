package tooling

import "fmt"

const (
	M_CncMilling = iota
	M_CncLathe
	M_CncGrinder
	M_CncDrill
	M_CncRouter
	M_CncLaserCutter
	M_CncPlasmaCutter
	M_CncWaterJetCutting
	M_CncElectricDischargeMachines // (EDMs), but not the music
)

type Point struct {
	X float64
	Y float64
	Z float64
}

func (p *Point) String() string {
	return fmt.Sprintf("X:%v, Y:%v, Z:%v", p.X, p.Y, p.Z)
}

const (
	FEED_INVERSE_TIME = iota
	FEED_PER_MINUTE
	FEED_PER_REVOLUTION
)

const (
	PLANE_NONE = iota
	PLANE_XY
	PLANE_XZ
	PLANE_YZ
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
	ToolChangeTo(tool int)
	SelectPlane(plane int)
	Reset()
}
