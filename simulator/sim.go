package main

import (
	"github.com/timleecasey/stllib/lib/aid3/gcode"
	"github.com/timleecasey/stllib/lib/aid3/sim"
	"log"
	"os"
)

func main() {
	gcodeFileNm := os.Args[1]
	tree, err := gcode.Parse(gcodeFileNm)
	if err != nil {
		log.Printf("Could not parse %v: %v", gcodeFileNm, err)
	}
	s := &sim.Sim{}
	s.Start()
	s.Run(tree)
}
