package main

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
				entityGrid[i+gridCentre][j+gridCentre] =  []Entity{}
			}
		}
	}

	for _, e := range entities {

		x := int(e.x)
		y := int(e.y)

		if x < iMin || x > iMax || j < jMin	|| j > jMax { continue }

		if x >= -gridCentre && y >= -gridCentre && x < gridCentre && y < gridCentre {
			entityGrid[x+gridCentre][y+gridCentre] = append(entityGrid[x+gridCentre][y+gridCentre], e)
		}

	}

}