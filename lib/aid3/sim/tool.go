package sim

import (
	"github.com/timleecasey/stllib/lib/aid3/gcode"
	"strconv"
)

func cmdToolChange(s *Sim, cn *gcode.CmdNode) error {
	src := cn.Cmd.Src()
	toolStr := src[1:]
	if tool, err := strconv.ParseInt(toolStr, 10, 32); err != nil {
		s.Tool.SpindleSpeed(tool)
	} else {
		return err
	}
	return nil

}
