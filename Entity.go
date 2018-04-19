package main

import (
	"io/ioutil"
	"github.com/faiface/pixel/pixelgl"
)

var entityUID uint32 = 0

type Entity struct {
	id     uint32
	active bool

	x float64 // position
	y float64
	z float64

	lastX float64
	lastY float64
	lastZ float64

	targetX float64
	targetY float64
	targetZ float64

	progress float64

	velocity float64 // velocity
	direction byte
	distance int     // number of squares to continue at that velocity

	class  string            // corresponds to a lua file
	script string            // their lua script
	flags  map[string]float64 // entity flags map
	timers map[string]uint16

	sprite         int
	animated       bool
	startSprite    int
	endSprite      int
	animationSpeed float64
}

func preRenderEntities() {

	for j := -gridCentre; j < gridCentre; j++ {
		for i := -gridCentre; i < gridCentre; i++ {
			if len(entityGrid[i+gridCentre][j+gridCentre]) > 0 {
				entityGrid[i+gridCentre][j+gridCentre] = entityGrid[i+gridCentre][j+gridCentre][:0]
			}
		}
	}

	live := 1
	if editing { live = 0 }

	for _, e := range entities[live] {

		gx := int(e.lastX)
		gy := int(e.lastY)

		switch viewDirection {
		case 0:
			if e.lastY < e.targetY {
				gy++
			}
		case 1:
			if e.lastX > e.targetX {
				gx--
			}
		case 2:
			if e.lastY > e.targetY {
				gy--
			}
		case 3:
			if e.lastX < e.targetX {
				gx++
			}
		}

		if gx >= -gridCentre && gy >= -gridCentre && gx < gridCentre && gy < gridCentre {
			entityGrid[gx+gridCentre][gy+gridCentre] = append(entityGrid[gx+gridCentre][gy+gridCentre], e)
		}

	}


}

func updateEntities() {

	select {
	case <-luaTick:

		for _, e := range entities[1] {
			currentEntity = e.id
			executeLua(L, e.script)
		}

		for k := range gameKeys {
			if v, ok := gameKeyDown[k]; ok {
				if v {
					gameKeyDownLast[k] = true
				} else {
					gameKeyDownLast[k] = false
				}
			}
		}

	default:
	}

	for i := range entities[1] {

		if entities[1][i].progress+entities[1][i].velocity/60 > 1 {

			entities[1][i].progress = 0
			entities[1][i].x = entities[1][i].targetX
			entities[1][i].y = entities[1][i].targetY
			entities[1][i].z = entities[1][i].targetZ
			entities[1][i].lastX = entities[1][i].targetX
			entities[1][i].lastY = entities[1][i].targetY
			entities[1][i].lastZ = entities[1][i].targetZ

		} else if entities[1][i].targetX != entities[1][i].lastX || entities[1][i].targetY != entities[1][i].lastY || entities[1][i].targetZ != entities[1][i].lastZ {

			entities[1][i].progress += entities[1][i].velocity / 60
			entities[1][i].x = entities[1][i].lastX + (entities[1][i].targetX-entities[1][i].lastX)*entities[1][i].progress
			entities[1][i].y = entities[1][i].lastY + (entities[1][i].targetY-entities[1][i].lastY)*entities[1][i].progress
			entities[1][i].z = entities[1][i].lastZ + (entities[1][i].targetZ-entities[1][i].lastZ)*entities[1][i].progress

		}

		if entities[1][i].progress == 0 && entities[1][i].direction != 0 && entities[1][i].distance > 0 {

			dx := 0
			dy := 0
			entities[1][i].distance--

			switch entities[1][i].direction {
			case 'N':
				dy = -1
			case 'W':
				dx = -1
			case 'S':
				dy = 1
			case 'E':
				dx = 1
			}

			gX := int(entities[1][i].x) + dx + gridCentre
			gY := int(entities[1][i].y) + dy + gridCentre
			gZ := int(entities[1][i].z)

			if !(gX < 0 || gY < 0 || gX >= 2*gridCentre || gY >= 2*gridCentre) && grid[gX][gY][gZ][1] == 0 {
				entities[1][i].targetX = entities[1][i].x + float64(dx)
				entities[1][i].targetY = entities[1][i].y + float64(dy)
				entities[1][i].targetZ = entities[1][i].z
			}

		}

	}

}

func resetEntities() {

	entities[1] = entities[1][:0]
	entities[1] = make([]Entity, len(entities[0]))
	copy(entities[1], entities[0])

	for i := range entities[1] {

		script, err := ioutil.ReadFile("scripts/" + entities[1][i].class + ".lua")
		check(err)
		entities[1][i].script = "do\n" +string(script) + "\nend\n"
		entities[1][i].flags = make(map[string]float64)

	}

	gameKeyDown = make(map[pixelgl.Button]bool)
	gameKeyDownLast = make(map[pixelgl.Button]bool)

	for k := range gameKeys {
		gameKeyDownLast[k] = false
	}

}
