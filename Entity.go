package main

import (
	"time"
	"math/rand"
)

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

func updateEntities() {

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for i := range entities {

		noTarget := entities[i].lastX == entities[i].targetX &&
			entities[i].lastY == entities[i].targetY &&
			entities[i].lastZ == entities[i].targetZ

		if noTarget || entities[i].progress+entities[i].velocity/60 > 1 {

			d := -1

			if entities[i].targetY < entities[i].lastY {
				d = 0
			}
			if entities[i].targetX < entities[i].lastX {
				d = 1
			}
			if entities[i].targetY > entities[i].lastY {
				d = 2
			}
			if entities[i].targetX > entities[i].lastX {
				d = 3
			}

			entities[i].progress = 0
			entities[i].x = entities[i].targetX
			entities[i].y = entities[i].targetY
			entities[i].z = entities[i].targetZ
			entities[i].lastX = entities[i].x
			entities[i].lastY = entities[i].y
			entities[i].lastZ = entities[i].z

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

				gX := int(entities[i].x + dx + gridCentre)
				gY := int(entities[i].y + dy + gridCentre)
				gZ := int(entities[i].z)

				if gX < 0 || gY < 0 || gX >= 2*gridCentre || gY >= 2*gridCentre {
					continue
				}
				if grid[gX][gY][gZ][1] != 0 {
					continue
				}

				entities[i].targetX = entities[i].x + dx
				entities[i].targetY = entities[i].y + dy
				entities[i].targetZ = entities[i].z

				break

			}

		} else {

			entities[i].progress += entities[i].velocity / 60
			entities[i].x = entities[i].lastX + (entities[i].targetX-entities[i].lastX)*entities[i].progress
			entities[i].y = entities[i].lastY + (entities[i].targetY-entities[i].lastY)*entities[i].progress
			entities[i].z = entities[i].lastZ + (entities[i].targetZ-entities[i].lastZ)*entities[i].progress

		}

	}

	for j := -gridCentre; j < gridCentre; j++ {
		for i := -gridCentre; i < gridCentre; i++ {
			if len(entityGrid[i+gridCentre][j+gridCentre]) > 0 {
				entityGrid[i+gridCentre][j+gridCentre] = entityGrid[i+gridCentre][j+gridCentre][:0]
			}
		}
	}

	for _, e := range entities {

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

func createEntities() {

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for i := uint32(0); i < 10; i ++ {

		e := Entity{id: i,
			active: true,
			sprite: 0,
			z: 0,
			velocity: 5}

		for failCount := 0; failCount < 10; failCount++ {

			e.x = float64(int(r.Float64()*10 - 5))
			e.y = float64(int(r.Float64()*10 - 5))

			gX := int(e.x + gridCentre)
			gY := int(e.y + gridCentre)
			gZ := int(e.z)

			if gX < 0 || gY < 0 || gX >= 2*gridCentre || gY >= 2*gridCentre {
				continue
			}
			if grid[gX][gY][gZ][1] != 0 {
				continue
			}

		}

		dx := 0.0
		dy := 0.0
		d := r.Intn(4)

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

		e.lastX = e.x
		e.lastY = e.y
		e.lastZ = e.z
		e.targetX = e.x + dx
		e.targetY = e.y + dy
		e.targetZ = e.z
		e.progress = 0
		entities = append(entities, e)

	}

}
