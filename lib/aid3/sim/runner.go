package sim

import (
	"fmt"
	"github.com/timleecasey/stllib/lib/aid3/gcode"
	"github.com/timleecasey/stllib/lib/aid3/sim/reality"
	"github.com/timleecasey/stllib/lib/aid3/sim/tooling"
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
}
