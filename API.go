package main

import (
	"github.com/yuin/gopher-lua"
	"fmt"
	"os"
	"github.com/faiface/pixel/pixelgl"
)

func initiateAPI() {

	linkToLua(L, APILoadMap, "LoadMap")
	linkToLua(L, APIGetTile, "GetTile")
	linkToLua(L, APISetTile, "SetTile")
	linkToLua(L, APIGetId, "GetId")
	linkToLua(L, APISetFocus, "SetFocus")
	linkToLua(L, APISetZoom, "SetZoom")
	linkToLua(L, APIEntityPosition, "GetEntityPosition")
	linkToLua(L, APIKeyPressed, "KeyPressed")
	linkToLua(L, APISetFlag, "SetFlag")
	linkToLua(L, APIGetFlag, "GetFlag")

	linkToLua(L, luaPrint, "print")

}

func APILoadMap(L *lua.LState) int {

	originalLevelFile := levelFile
	candidateLevelFile := "maps/" + L.ToString(1)

	if _, err := os.Stat(candidateLevelFile); err == nil {
		levelFile = candidateLevelFile
		load()
		levelFile = originalLevelFile
	} else {
		luaConsolePrint (fmt.Sprintf("Map file %s not found.", candidateLevelFile))
	}

	return 0
}

func APIGetTile(L *lua.LState) int {

	x := L.ToInt(1) + gridCentre
	y := L.ToInt(2) + gridCentre
	z := L.ToInt(3)

	if x < 0 || y < 0 || z < 0 || x >= 2*gridCentre || y >= 2*gridCentre || z > 15 {
		L.Push(lua.LNumber(0))
		L.Push(lua.LNumber(0))
	} else {
		L.Push(lua.LNumber(grid[x][y][z][0]))
		L.Push(lua.LNumber(grid[x][y][z][1]))
	}

	return 2

}

func APISetTile(L *lua.LState) int {

	x := L.ToInt(1) + gridCentre
	y := L.ToInt(2) + gridCentre
	z := L.ToInt(3)
	a := L.ToInt(4)
	b := L.ToInt(5)

	if x < 0 || y < 0 || z < 0 || x >= 2*gridCentre || y >= 2*gridCentre || z > 15 {
		return 0
	}

	if a < 0 || b < 0 {
		return 0
	}

	grid[x][y][z][0] = uint16(a)
	grid[x][y][z][1] = uint16(b)

	return 0

}

func APIGetId(L *lua.LState) int {
	L.Push(lua.LNumber(currentEntity))
	return 1
}

func APISetFocus(L *lua.LState) int {

	x := float64(L.ToNumber(1))
	y := float64(L.ToNumber(2))

	cameraX = -x*128 - 64
	cameraY = y*128 + 64

	return 0
}

func APISetZoom(L *lua.LState) int {

	scale = float64(L.ToNumber(1))

	if scale < 0.1 {
		scale = 0.1
	}
	if scale > 2.0 {
		scale = 2.0
	}

	hScale = 128 * scale
	vScale = 128 * aspect * scale

	return 0
}

func APIEntityPosition(L *lua.LState) int {

	id := L.ToInt(1)

	for _, e := range entities[1] {

		if e.id == uint32(id) {
			L.Push(lua.LNumber(e.x))
			L.Push(lua.LNumber(e.y))
			L.Push(lua.LNumber(e.z))
			return 3
		}

	}

	L.Push(lua.LNumber(0))
	L.Push(lua.LNumber(0))
	L.Push(lua.LNumber(0))
	return 3

}

func APIKeyPressed(L *lua.LState) int {

	keyString := L.ToString(1)
	keyJust := L.ToBool(2)

	var key pixelgl.Button = -1

	for k, v := range gameKeys {
		if v == keyString {
			key = k
		}
	}

	if key != -1 {
		isPressed, ok := gameKeyDown[key]
		if ok && isPressed == keyJust {
			L.Push(lua.LTrue)
		} else {
			L.Push(lua.LFalse)
		}
	} else {
		L.Push(lua.LFalse)
	}

	return 1

}

func APISetFlag(L *lua.LState) int {
	id := uint32(L.ToInt(1))
	flag := string(L.ToString(2))
	value := float64(L.ToNumber(3))

	for _, e := range entities[1] {
		if e.id == id {
			e.flags[flag] = value
			break
		}
	}

	return 0
}

func APIGetFlag(L *lua.LState) int {
	id := uint32(L.ToInt(1))
	flag := string(L.ToString(2))

	for _, e := range entities[1] {
		if e.id == id {
			value, ok := e.flags[flag]
			if ok {
				L.Push(lua.LNumber(value))
			} else {
				L.Push(lua.LNumber(0))
			}
			break
		}
	}

	return 1
}


/* TEMPLATE
func APIxxx(L *lua.LState) int {
	x := L.ToString(1)
	fmt.Println(x)
	return 0
}
*/

