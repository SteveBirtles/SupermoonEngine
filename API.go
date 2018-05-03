package main

import (
	"github.com/yuin/gopher-lua"
	"fmt"
	"os"
	"github.com/faiface/pixel/pixelgl"
	"strings"
	"math"
	"time"
	"io/ioutil"
)

func initiateAPI() {

	linkToLua(L, APILoadMap, "LoadMap")
	linkToLua(L, APIGetTile, "GetTile")
	linkToLua(L, APISetTile, "SetTile")
	linkToLua(L, APIGetId, "GetId")
	linkToLua(L, APINearby, "Nearby")
	linkToLua(L, APIProximity, "Proximity")
	linkToLua(L, APISetFocus, "SetFocus")
	linkToLua(L, APIGetFocus, "GetFocus")
	linkToLua(L, APISetModal, "SetModal")
	linkToLua(L, APISetClassActive, "SetClassActive")
	linkToLua(L, APISetView, "SetView")
	linkToLua(L, APIGetPosition, "GetPosition")
	linkToLua(L, APIGetDistance, "GetDistance")
	linkToLua(L, APISetPosition, "SetPosition")
	linkToLua(L, APISetVelocity, "SetVelocity")
	linkToLua(L, APIKeyPressed, "KeyPressed")
	linkToLua(L, APISetFlag, "SetFlag")
	linkToLua(L, APIGetFlag, "GetFlag")
	linkToLua(L, APIListFlags, "ListFlags")
	linkToLua(L, APIStartTimer, "StartTimer")
	linkToLua(L, APIGetTimer, "GetTimer")
	linkToLua(L, APICreate, "Create")
	linkToLua(L, APIDelete, "Delete")
	linkToLua(L, APIGetClass, "GetClass")
	linkToLua(L, APISetClass, "SetClass")
	linkToLua(L, APIGetScript, "GetScript")
	linkToLua(L, APIOverride, "Override")
	linkToLua(L, APIReset, "Reset")
	linkToLua(L, APIEndGame, "EndGame")

	linkToLua(L, luaPrint, "print")

}

func APICreate(L *lua.LState) int {

	x := float64(L.ToInt(1))
	y := float64(L.ToInt(2))
	z := float64(L.ToInt(3))
	class := L.ToString(4)

	script, err := ioutil.ReadFile("scripts/" + class + ".lua")
	check(err)

	entityDynamicID++

	e := Entity{id: entityDynamicID,
		active: true,
		sprite: 0,
		velocity: 0,
		direction: 0,
		distance: 0,
		class: class,
		x: x,
		y: y,
		z: z,
		lastX: x,
		lastY: y,
		lastZ: z,
		targetX: x,
		targetY: y,
		targetZ: z,
		progress: 0,
		script: "do\n" +string(script) + "\nend\n",
		flags: make(map[string]float64),
		timers: make(map[string]time.Time),
	}

	entities[1] = append(entities[1], e)

	return 0

}


func APIDelete(L *lua.LState) int {

	id := uint32(L.ToNumber(1))
	if id == 0 { fmt.Println("Lua error: id not specified") }

	index := -1

	for i, e := range entities[1] {
		if e.id == id {
			index = i
			break
		}
	}

	if index > -1 {
		entities[1] = append(entities[1][:index], entities[1][index+1:]...)
	}

	return 0
}

func APILoadMap(L *lua.LState) int {

	originalLevelFile := levelFile
	candidateLevelFile := "maps/" + L.ToString(1)
	loadEntities := L.ToBool(2)

	if _, err := os.Stat(candidateLevelFile); err == nil {
		levelFile = candidateLevelFile
		load()
		levelFile = originalLevelFile

		if loadEntities {
			resetEntities()
		}

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

func APINearby(L *lua.LState) int {

	ids := L.NewTable()

	id := L.ToInt(1)
	radius := L.ToInt(2)

	if id == 0 { fmt.Println("Lua error: id not specified") }
	if radius == 0 { fmt.Println("Lua error: radius 0 or not specified") }

	for _, e1 := range entities[1] {
		if e1.id == uint32(id) {
			for _, e2 := range entities[1] {
				if e2.id != uint32(id) {

					if math.Pow(float64(e2.x-e1.x), 2) + math.Pow(float64(e2.y-e1.y), 2) <= math.Pow(float64(radius), 2) {
						ids.Append(lua.LNumber(int(e2.id)))
					}
				}
			}
			break
		}
	}

	L.Push(ids)
	return 1
}

func APIProximity(L *lua.LState) int {

	id1 := L.ToInt(1)
	id2 := L.ToInt(1)

	if id1 == 0 { fmt.Println("Lua error: id1 not specified") }
	if id2 == 0 { fmt.Println("Lua error: id2 not specified") }

	if id1 != id2 {
		for _, e1 := range entities[1] {
			if e1.id == uint32(id1) {
				for _, e2 := range entities[1] {
					if e2.id == uint32(id2) {

						d := math.Sqrt(math.Pow(float64(e2.x-e1.x), 2) + math.Pow(float64(e2.y-e1.y), 2))

						L.Push(lua.LNumber(d))
						return 1
					}
				}
			}
		}
	}

	L.Push(lua.LNumber(0))
	return 1

}

func APISetFocus(L *lua.LState) int {

	id := uint32(L.ToNumber(1))
	follow := bool(L.ToBool(2))

	if id == 0 { fmt.Println("Lua error: id not specified") }

	x := 0.0
	y := 0.0

	found := false

	for _, e := range entities[1] {
		if e.id == id {
			x = e.x
			y = e.y
			found = true
		}
	}

	if found {
		cameraX = -x*128 - 64
		cameraY = y*128 + 64

		if follow {
			focusEntity = id
		}

	}

	return 0
}

func APISetModal(L *lua.LState) int {

	id := uint32(L.ToNumber(1))

	for _, e := range entities[1] {
		if e.id == id {
			modalEntity = id
			return 0
		}
	}

	modalEntity = 0
	return 0
}


func APISetClassActive(L *lua.LState) int {

	class := L.ToString(1)
	radius := L.ToInt(2)

	if class == "" {
		fmt.Println("Lua error: class not specified")
	} else {
		if radius >= gridCentre*2 {
			delete(entityClassActiveRadius, class)
		} else {
			entityClassActiveRadius[class] = radius
		}
	}

	return 0
}

func APIGetFocus(L *lua.LState) int {
	L.Push(lua.LNumber(focusEntity))
	return 1
}

func APIGetClass(L *lua.LState) int {

	id := L.ToInt(1)

	if id == 0 { fmt.Println("Lua error: id not specified") }

	for _, e := range entities[1] {

		if e.id == uint32(id) {
			L.Push(lua.LString(e.class))
			return 1
		}

	}

	L.Push(lua.LString(""))

	return 1
}

func APISetClass(L *lua.LState) int {

	id := L.ToInt(1)
	class := L.ToString(2)

	if id == 0 { fmt.Println("Lua error: id not specified") }

	for i, e := range entities[1] {
		if e.id == uint32(id) {
			entities[1][i].class = class
			return 0
		}
	}
	return 0

}


func APIReset(L *lua.LState) int {

	id := L.ToInt(1)

	if id == 0 { fmt.Println("Lua error: id not specified") }

	for i, e := range entities[1] {

		if e.id == uint32(id) {
			script, err := ioutil.ReadFile("scripts/" + entities[1][i].class + ".lua")
			check(err)
			entities[1][i].script = "do\n" +string(script) + "\nend\n"
			entities[1][i].flags = make(map[string]float64)
			entities[1][i].timers = make(map[string]time.Time)
			return 0
		}

	}

	return 0
}

func APIGetScript(L *lua.LState) int {

	id := L.ToInt(1)

	if id == 0 { fmt.Println("Lua error: id not specified") }

	for _, e := range entities[1] {

		if e.id == uint32(id) {
			L.Push(lua.LString(e.script))
			return 1
		}

	}

	L.Push(lua.LString(""))

	return 1
}

func APIOverride(L *lua.LState) int {

	id := L.ToInt(1)
	script := L.ToString(2)

	if id == 0 { fmt.Println("Lua error: id not specified") }

	for i, e := range entities[1] {
		if e.id == uint32(id) {
			entities[1][i].script = script
			return 0
		}
	}
	return 0

}

func APISetView(L *lua.LState) int {

	direction := strings.ToUpper(L.ToString(1))

	switch direction {
	case "N":
		viewDirection = 0
	case "E":
		viewDirection = 1
	case "S":
		viewDirection = 2
	case "W":
		viewDirection = 3
	}

	requestedScale := float64(L.ToNumber(2))

	if requestedScale > 0 {

		scale = requestedScale

		if scale < 0.1 {
			scale = 0.1
		}
		if scale > 2.0 {
			scale = 2.0
		}

		hScale = 128 * scale
		vScale = 128 * aspect * scale

	}

	return 0
}

func APIGetPosition(L *lua.LState) int {

	id := L.ToInt(1)

	if id == 0 { fmt.Println("Lua error: id not specified") }

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

func APIGetDistance(L *lua.LState) int {

	id := L.ToInt(1)

	if id == 0 { fmt.Println("Lua error: id not specified") }

	for _, e := range entities[1] {

		if e.id == uint32(id) {
			L.Push(lua.LNumber(e.distance))
			return 1
		}

	}

	L.Push(lua.LNumber(0))
	return 1

}


func APISetPosition(L *lua.LState) int {

	id := L.ToInt(1)
	x := L.ToInt(2)
	y := L.ToInt(3)
	z := L.ToInt(4)

	if id == 0 { fmt.Println("Lua error: id not specified") }

	for i, e := range entities[1] {

		if e.id == uint32(id) {

			entities[1][i].x = float64(x)
			entities[1][i].y = float64(y)
			entities[1][i].z = float64(z)

			entities[1][i].targetX = entities[1][i].x
			entities[1][i].targetY = entities[1][i].y
			entities[1][i].targetZ = entities[1][i].z
			entities[1][i].lastX = entities[1][i].targetX
			entities[1][i].lastY = entities[1][i].targetY
			entities[1][i].lastZ = entities[1][i].targetZ
			entities[1][i].distance = 0
			entities[1][i].progress = 0

			updateFocus()

			return 0

		}

	}

	return 0

}


func APIKeyPressed(L *lua.LState) int {

	keyString := strings.ToUpper(L.ToString(1))
	keyJust := L.ToBool(2)

	if keyString == "" { fmt.Println("Lua error: key not specified") }

	var key pixelgl.Button = -1

	for k, v := range gameKeys {
		if v == keyString {
			key = k
		}
	}

	if key != -1 {
		if keyJust {
			if gameKeyWasPressed[key] {
				L.Push(lua.LTrue)
			} else {
				L.Push(lua.LFalse)
			}
		} else {
			if gameKeyWasPressed[key] || (gameKeyDownStart[key] && gameKeyDownEnd[key]) {
				L.Push(lua.LTrue)
			} else {
				L.Push(lua.LFalse)
			}
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

	if id == 0 { fmt.Println("Lua error: id not specified") }
	if flag == "" { fmt.Println("Lua error: flag not specified") }

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

	if id == 0 { fmt.Println("Lua error: id not specified") }
	if flag == "" { fmt.Println("Lua error: flag not specified") }

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

func APIListFlags(L *lua.LState) int {

	flags := L.NewTable()

	id := L.ToInt(1)

	if id == 0 { fmt.Println("Lua error: id not specified") }

	for _, e := range entities[1] {
		if e.id == uint32(id) {
			for f := range e.flags {
				flags.Append(lua.LString(f))
			}
			break
		}
	}

	L.Push(flags)
	return 1
}


func APIStartTimer(L *lua.LState) int {
	id := uint32(L.ToInt(1))
	timer := string(L.ToString(2))

	if id == 0 { fmt.Println("Lua error: id not specified") }
	if timer == "" { fmt.Println("Lua error: timer not specified") }

	for _, e := range entities[1] {
		if e.id == id {
			e.timers[timer] = time.Now()
			break
		}
	}

	return 0
}

func APIGetTimer(L *lua.LState) int {
	id := uint32(L.ToInt(1))
	timer := string(L.ToString(2))

	if id == 0 { fmt.Println("Lua error: id not specified") }
	if timer == "" { fmt.Println("Lua error: timer not specified") }

	for _, e := range entities[1] {
		if e.id == id {
			value, ok := e.timers[timer]
			if ok {
				elapsed := time.Now().Sub(value)
				L.Push(lua.LNumber(elapsed.Seconds()))
			} else {
				L.Push(lua.LNumber(0))
			}
			break
		}
	}

	return 1
}

func APISetVelocity(L *lua.LState) int {
	id := uint32(L.ToInt(1))
	dirString := strings.ToUpper(L.ToString(2)) + "-"
	dir := []byte(dirString[0:1])
	vel := float64(L.ToNumber(3))
	dist := int(L.ToInt(4))

	if id == 0 { fmt.Println("Lua error: id not specified") }
	if dir[0] == '-' { fmt.Println("Lua error: direction not specified") }
	if vel == 0 { fmt.Println("Lua error: velocity not specified") }
	if dist == 0 { fmt.Println("Lua error: distance not specified") }

	if dir[0] == 'N' || dir[0] == 'E' || dir[0] == 'S' || dir[0] == 'W' {
		for i, e := range entities[1] {
			if e.id == id {
				entities[1][i].nextDirection = dir[0]
				entities[1][i].nextVelocity = vel
				entities[1][i].distance = dist
			}
		}
	}

	return 0
}


func APIEndGame(_ *lua.LState) int {

	copyGrid(&gridBackup, &grid)
	editing = true
	resetViewState()

	return 0
}


/* TEMPLATE
func APIxxx(L *lua.LState) int {
	x := L.ToString(1)
	fmt.Println(x)
	return 0
}
*/

