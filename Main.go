package main

import (
	_ "image/png"
	"github.com/yuin/gopher-lua"
	"github.com/faiface/pixel/pixelgl"
	"os"
	"flag"
	"log"
	"runtime/pprof"
	"runtime"
)

var (
	L *lua.LState
	levelFile = "maps/user.dat"
	defaultLevelFile = "maps/default.dat"
	screenWidth = 1280.0
	screenHeight = 720.0
)

var mapfile = flag.String("map", "", "loads a give map file")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
var vres = flag.String("vres", "", "choose between 720 or 1080 vertical resolution")

func main() {

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if *mapfile != "" {
		levelFile = "maps/" + *mapfile
	}

	if *vres != "" {
		if *vres == "1080" {
			screenWidth = 1920.0
			screenHeight = 1080.0
		}
	}

	L = lua.NewState()
	defer L.Close()
	
	pixelgl.Run(mainLoop)

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}

}