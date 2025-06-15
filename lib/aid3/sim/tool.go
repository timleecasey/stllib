package sim

import (
	"fmt"
)

func cmdToolChange(s *Sim, tool int64) {
	fmt.Printf("CHANGE TOOL %v\n", tool)
	s.Tool.ToolChangeTo(tool)
}
