package main

import (
	_ "image/png"
	"github.com/yuin/gopher-lua"
	"github.com/faiface/pixel/pixelgl"
	"os"
	"strings"
)

var (
	L *lua.LState
	levelFile = "maps/default.dat"
	screenWidth = 1280.0
	screenHeight = 720.0
)

func main() {

	args := os.Args[1:]

	if len(args) >= 1 {
		levelFile = "maps/" + args[0]
	}

	if len(args) >= 2 {
		switch strings.ToLower(args[1]) {
		case "hd":
			screenWidth = 1920.0
			screenHeight = 1080.0
		}
	}

	L = lua.NewState()
	defer L.Close()

	pixelgl.Run(mainLoop)

}