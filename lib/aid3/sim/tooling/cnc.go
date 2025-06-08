package tooling

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

type Cnc interface {
	Axis() []int
	ZeroPoint() *Point
	Head() Head
	FeedRate() float64
	AssignFeedRate(f float64)
}
