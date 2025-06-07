package tooling

type Point struct {
	X float64
	Y float64
	Z float64
}

type Cnc interface {
	Axis() []int
	ZeroPoint() *Point
	Head() Head
}
