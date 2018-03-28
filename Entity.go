package main

import (
	"time"
	"math/rand"
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
	vN       int     // number of squares to continue at that velocity

	class      string            // corresponds to a lua file
	script     string            // their lua script
	properties map[string]string // entity properties map
	timers     map[string]uint16

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

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for i := range entities[1] {

		noTarget := entities[1][i].lastX == entities[1][i].targetX &&
			entities[1][i].lastY == entities[1][i].targetY &&
			entities[1][i].lastZ == entities[1][i].targetZ

		if noTarget || entities[1][i].progress+entities[1][i].velocity/60 > 1 {

			d := -1

			if entities[1][i].targetY < entities[1][i].lastY {
				d = 0
			}
			if entities[1][i].targetX < entities[1][i].lastX {
				d = 1
			}
			if entities[1][i].targetY > entities[1][i].lastY {
				d = 2
			}
			if entities[1][i].targetX > entities[1][i].lastX {
				d = 3
			}

			entities[1][i].progress = 0
			entities[1][i].x = entities[1][i].targetX
			entities[1][i].y = entities[1][i].targetY
			entities[1][i].z = entities[1][i].targetZ
			entities[1][i].lastX = entities[1][i].x
			entities[1][i].lastY = entities[1][i].y
			entities[1][i].lastZ = entities[1][i].z

			for failCount := 0; failCount < 10; failCount++ {

				dx := 0.0
				dy := 0.0

				if failCount > 0 || d == -1 {
					d = r.Intn(4)
				}

				switch d {
				case 0:
					dy = -1
				case 1:
					dx = -1
				case 2:
					dy = 1
				case 3:
					dx = 1
				}

				gX := int(entities[1][i].x + dx + gridCentre)
				gY := int(entities[1][i].y + dy + gridCentre)
				gZ := int(entities[1][i].z)

				if gX < 0 || gY < 0 || gX >= 2*gridCentre || gY >= 2*gridCentre {
					continue
				}
				if grid[gX][gY][gZ][1] != 0 {
					continue
				}

				entities[1][i].targetX = entities[1][i].x + dx
				entities[1][i].targetY = entities[1][i].y + dy
				entities[1][i].targetZ = entities[1][i].z

				break

			}

		} else {

			entities[1][i].progress += entities[1][i].velocity / 60
			entities[1][i].x = entities[1][i].lastX + (entities[1][i].targetX-entities[1][i].lastX)*entities[1][i].progress
			entities[1][i].y = entities[1][i].lastY + (entities[1][i].targetY-entities[1][i].lastY)*entities[1][i].progress
			entities[1][i].z = entities[1][i].lastZ + (entities[1][i].targetZ-entities[1][i].lastZ)*entities[1][i].progress

		}

	}

}

func resetEntities() {

	entities[1] = entities[1][:0]
	entities[1] = make([]Entity, len(entities[0]))
	copy(entities[1], entities[0])

}
