package main

import (
	"time"
	"math/rand"
	"fmt"
	"math"
)

type Entity struct {

  id uint32
  active bool

  x float64   // position
  y float64
  z float64

  dx float64  // velocity
  dy float64
  dz float64
  dn int      // number of squares to continue at that velocity

  class string  // corresponds to a lua file
  script string // their lua script
  properties map[string]string   // entity properties map
  timers map[string]uint16

  sprite int
  animated bool
  startSprite int
  endSprite int
  animationSpeed float64

}

func updateEntities() {

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for i := range entities {

		intX := math.Floor(entities[i].x)
		inTransitX := entities[i].x != intX
		intNextX := math.Floor(entities[i].x + entities[i].dx / 60)

		fmt.Printf("x:%f ix:%f dx:%f nx:%f\n", entities[i].x, intX, entities[i].dx, intNextX)

		if intX != intNextX {
			entities[i].x = intNextX
			inTransitX = false
		} else {
			entities[i].x += entities[i].dx / 60
		}

		intY := math.Floor(entities[i].y)
		inTransitY := entities[i].y != intY
		intNextY := math.Floor(entities[i].y + entities[i].dy / 60)

		if intY != intNextY {
			entities[i].y = intNextY
			inTransitY = false
		} else {
			entities[i].y += entities[i].dy / 60
		}

		if !(inTransitX || inTransitY) {

			dx := 0.0
			dy := 0.0

			d := r.Intn(4)

			fmt.Printf("%d\n", d)

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

			entities[i].dx = dx
			entities[i].dy = dy

		}


	}

	var i, j int
	iMin := 9999
	jMin := 9999
	iMax := -9999
	jMax := -9999

	for j0 := jStart + jOffset; j0 <= jEnd+ jOffset; j0++ {
		for i0 := iStart + iOffset; i0 <= iEnd+ iOffset; i0++ {
			switch viewDirection {
			case 0:
				i = int(i0)
				j = int(j0)
			case 1:
				i = int(-j0)
				j = int(i0)
			case 2:
				i = int(-i0)
				j = int(-j0)
			case 3:
				i = int(j0)
				j = int(-i0)
			}

			if i > iMax { iMax = i }
			if j > jMax { jMax = j }
			if i < iMin { iMin = i }
			if j < jMin { jMin = j }

			if i >= -gridCentre && j >= -gridCentre && i < gridCentre && j < gridCentre {
				if len(entityGrid[i+gridCentre][j+gridCentre]) > 0 {
					entityGrid[i+gridCentre][j+gridCentre] = entityGrid[i+gridCentre][j+gridCentre][:0]
				}
			}
		}
	}

	var x, y int

	for _, e := range entities {

		inTransit := e.y != float64(int(e.y))

		switch viewDirection {
		case 0:
			x = int(e.x)
			y = int(e.y)
			if inTransit { y++ }
		case 1:
			x = -int(e.y)
			y = int(e.x)
			if inTransit { x-- }
		case 2:
			x = -int(e.x)
			y = -int(e.y)
			if inTransit { y-- }
		case 3:
			x = int(e.y)
			y = -int(e.x)
			if inTransit { x++ }
		}

		if x < iMin || x > iMax || j < jMin	|| j > jMax { continue }

		if x >= -gridCentre && y >= -gridCentre && x < gridCentre && y < gridCentre {
			entityGrid[x+gridCentre][y+gridCentre] = append(entityGrid[x+gridCentre][y+gridCentre], e)
		}

	}

}