package sim

import (
	"github.com/timleecasey/stllib/lib/aid3/gcode"
	"strconv"
)

func cmdSpindleSpeed(s *Sim, cn *gcode.CmdNode) error {
	src := cn.Cmd.Src()
	speedStr := src[1:]
	if speed, err := strconv.ParseInt(speedStr, 10, 64); err != nil {
		s.Tool.SpindleSpeed(speed)
	} else {
		return err
	}
	return nil
}
