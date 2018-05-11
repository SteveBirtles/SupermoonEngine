package main

import (
	"io/ioutil"
	"github.com/faiface/pixel/pixelgl"
	"time"
	"math"
	"strings"
	"strconv"
	"sort"
	"fmt"
)

var (
	entityUID               uint32
	entityDynamicID         uint32
	entityClassActiveRadius map[string]int
	entityClass             []string
	lastEntityClass         int
	entityClassBlock        int
	entityClassBlockCount   int
	entityGrid              [2*gridCentre][2*gridCentre][]Entity
	entities                [2][]Entity
	focusEntity             uint32
	modalEntity             uint32
	currentEntity           uint32
	entityCount				int
	activeEntityCount		int
	entityExecutionTime		[]execTime
	sortedEntityExecutionTimes	[]execTime
	totalExecutionTime			int64
	cumulativeExecution			int64
	executionLastSecond			int64
)

type execTime struct {
	id   uint32
	time int64
	class string
}

type Entity struct {

	//exported:

	Id    uint32
	X     float64 // position
	Y     float64
	Z     float64
	Class string // corresponds to a lua file

	//not exported:

	active   bool
	new      bool
	resetNew bool
	deleteMe bool

	lastX float64
	lastY float64
	lastZ float64

	targetX float64
	targetY float64
	targetZ float64

	progress   float64
	onTile     bool
	onTileX    float64
	onTileY    float64
	onTileZ    float64
	justOnTile bool

	nextDirection byte
	nextVelocity float64 // velocity
	nextDistance int     // number of squares to continue at that velocity

	direction byte
	velocity float64 // velocity
	distance int     // number of squares to continue at that velocity

	script string            // their lua script
	flags  map[string]string // entity flags map
	timers map[string]time.Time

	sprite          [4]int
	animationSpeed  [4]float64
	firstSprite     [4]int
	lastSprite      [4]int
	staticAnimation [4]bool
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

	updateFocus()

}

func updateFocus() {

	if !editing && focusEntity > 0 {

		x := 0.0
		y := 0.0

		found := false

		for _, e := range entities[1] {
			if e.Id == focusEntity {
				x = e.X
				y = e.Y
				found = true
			}
		}

		if found {
			cameraX = -x*128 - 64
			cameraY = y*128 + 64
		}

	}

}

func updateEntities() {

	for i := range entities[1] {
		entities[1][i].active = true

		if radius, featured := entityClassActiveRadius[entities[1][i].Class]; featured {

			if radius == 0 {
				entities[1][i].active = false
				break
			}

			x := -entities[1][i].X*128 - 64
			y := entities[1][i].Y*128 + 64

			if math.Pow(x - cameraX, 2) + math.Pow(y - cameraY, 2) >= math.Pow(float64(radius)*128, 2) {
				entities[1][i].active = false
			}

		}
	}

	select {
	case <-luaTick:

		for k := range gameKeys {
			gameKeyDownStart[k] = gameKeyDownEnd[k]
			gameKeyDownEnd[k] = win.Pressed(k)
			gameKeyWasPressed[k] = gameKeyJustPressed[k]
			gameKeyJustPressed[k] = false
		}

		//fmt.Print("--------------------------------------------------------------------------\n")

		entityCount = 0
		activeEntityCount = 0
		totalExecutionTime = 0
		entityExecutionTime	= make([]execTime, 0)

		entityIds := make([]uint32, 0)
		sortedEntityExecutionTimes = make([]execTime, 0)

		for i, e := range entities[1] {

			for _, id := range entityIds {
				if id == e.Id {
					fmt.Printf("Entity id clash: %d", id)
				}
			}
			entityIds = append(entityIds, e.Id)

			entityCount++

			//fmt.Printf("%d (%d, %d, %d) %s\n", e.Id, int(e.targetX), int(e.targetY), int(e.targetZ), e.Class)

			if modalEntity == 0 && e.active || modalEntity > 0 && e.Id == modalEntity && !e.deleteMe {

				activeEntityCount++

				script := ""

				inBlock := 0
				blockKey := ""
				blockValue := ""

				for _, line := range strings.Split(e.script, "\n") {

					line = strings.TrimSpace(line)

					if line == "" || strings.HasPrefix(line, "--") { continue }

					if line == "#if_new" {
						inBlock = 1
						continue
					} else if strings.HasPrefix(line, "#if_step") {
						inBlock = 2
						continue
					} else if strings.HasPrefix(line, "#if_focus") {
						inBlock = 3
						continue
					} else if strings.HasPrefix(line, "#if_flag") {
						inBlock = 4
						pair := strings.Split(line, " ")
						if len(pair) != 3 {
							blockKey = ""
							blockValue = ""
						} else {
							blockKey = pair[1]
							blockValue = pair[2]
						}
						continue
					} else if strings.HasPrefix(line,  "#if_timer") {
						inBlock = 5
						pair := strings.Split(line, " ")
						if len(pair) != 3 {
							blockKey = ""
							blockValue = ""
						} else {
							blockKey = pair[1]
							blockValue = pair[2]
						}
						continue
					} else if strings.HasPrefix(line, "#always") {
						inBlock = 0
						continue
					}

					includeLine := true

					switch inBlock {
					case 1:
						if !e.new { includeLine = false }
					case 2:
						if !e.justOnTile { includeLine = false }
					case 3:
						if e.Id != focusEntity { includeLine = false }
					case 4:
						if value, err1 := e.flags[blockKey]; err1 {
							if value != blockValue {
								includeLine = false
							}
						} else {
							includeLine = false
						}
					case 5:
						if value, err1 := e.timers[blockKey]; err1 {
							t0 := time.Now().Sub(value).Seconds()
							t1, err2 := strconv.ParseFloat(blockValue, 64)
							if err2 == nil && t0 < t1 {
								includeLine = false
							}
						} else {
							includeLine = false
						}
					}


					if includeLine {
						script += line + "\n"
					}

				}

				startTime := time.Now().UnixNano()

				if script != "" {
					currentEntity = e.Id
					executeLua(L, "do\n"+script+"\nend\n")
				}

				t := time.Now().UnixNano() - startTime
				entityExecutionTime = append(entityExecutionTime, execTime{e.Id, t, e.Class})
				totalExecutionTime += t

				entities[1][i].new = entities[1][i].resetNew
				entities[1][i].resetNew = false
				entities[1][i].onTile = false

			}
		}

		cumulativeExecution += totalExecutionTime

		for _, e := range entityExecutionTime {
			sortedEntityExecutionTimes = append(sortedEntityExecutionTimes, execTime{e.id, e.time, e.class})
		}

		sort.Slice(sortedEntityExecutionTimes, func(i, j int) bool {
			return sortedEntityExecutionTimes[i].time > sortedEntityExecutionTimes[j].time
		})

		for i := 0; i < len(entities[1]); {
			if entities[1][i].deleteMe {
				entities[1] = append(entities[1][:i], entities[1][i+1:]...)
			} else {
				i++
			}
		}

	default:
	}

	for i := range entities[1] {

		if modalEntity == 0 && entities[1][i].active || modalEntity > 0 && entities[1][i].Id == modalEntity {

			if entities[1][i].progress+entities[1][i].velocity/60 > 1 {

				entities[1][i].justOnTile = true
				entities[1][i].progress = 0
				entities[1][i].X = entities[1][i].targetX
				entities[1][i].Y = entities[1][i].targetY
				entities[1][i].Z = entities[1][i].targetZ
				entities[1][i].lastX = entities[1][i].targetX
				entities[1][i].lastY = entities[1][i].targetY
				entities[1][i].lastZ = entities[1][i].targetZ


			} else if entities[1][i].targetX != entities[1][i].lastX || entities[1][i].targetY != entities[1][i].lastY || entities[1][i].targetZ != entities[1][i].lastZ {

				entities[1][i].progress += entities[1][i].velocity / 60
				entities[1][i].X = entities[1][i].lastX + (entities[1][i].targetX-entities[1][i].lastX)*entities[1][i].progress
				entities[1][i].Y = entities[1][i].lastY + (entities[1][i].targetY-entities[1][i].lastY)*entities[1][i].progress
				entities[1][i].Z = entities[1][i].lastZ + (entities[1][i].targetZ-entities[1][i].lastZ)*entities[1][i].progress

			}

			if entities[1][i].progress == 0 {

				if entities[1][i].onTile {
					entities[1][i].justOnTile = false
				} else {
					entities[1][i].onTile = true
					entities[1][i].onTileX = entities[1][i].X
					entities[1][i].onTileY = entities[1][i].Y
					entities[1][i].onTileZ = entities[1][i].Z
				}

				if entities[1][i].nextDirection != 0 && entities[1][i].distance > 0 {

					dx := 0
					dy := 0
					entities[1][i].distance--

					entities[1][i].direction = entities[1][i].nextDirection
					entities[1][i].velocity = entities[1][i].nextVelocity

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

					gX := int(entities[1][i].X) + dx + gridCentre
					gY := int(entities[1][i].Y) + dy + gridCentre
					gZ := int(entities[1][i].Z)

					if !(gX < 0 || gY < 0 || gX >= 2*gridCentre || gY >= 2*gridCentre) && grid[gX][gY][gZ][1] == 0 {
						entities[1][i].targetX = entities[1][i].X + float64(dx)
						entities[1][i].targetY = entities[1][i].Y + float64(dy)
						entities[1][i].targetZ = entities[1][i].Z
					}
				}
			}

		}

	}

}

func resetEntities() {

	entities[1] = entities[1][:0]
	entities[1] = make([]Entity, len(entities[0]))
	copy(entities[1], entities[0])

	for i := range entities[1] {
		script, err := ioutil.ReadFile("scripts/" + entities[1][i].Class + ".lua")
		check(err)
		entities[1][i].active = true
		entities[1][i].new = true
		entities[1][i].script = string(script)
		entities[1][i].flags = make(map[string]string)
		entities[1][i].timers = make(map[string]time.Time)
		entities[1][i].direction = 'S'
		entities[1][i].distance = 0
		entities[1][i].velocity = 0
	}

	gameKeyDownStart = make(map[pixelgl.Button]bool)
	gameKeyDownEnd = make(map[pixelgl.Button]bool)
	gameKeyWasPressed = make(map[pixelgl.Button]bool)
	gameKeyJustPressed = make(map[pixelgl.Button]bool)
	gameKeyTimeSinceLastPressed = make(map[pixelgl.Button]int)
	entityClassActiveRadius = make(map[string]int)

	for k := range gameKeys {
		gameKeyDownStart[k] = false
		gameKeyDownEnd[k] = false
		gameKeyWasPressed[k] = false
		gameKeyJustPressed[k] = false
		gameKeyTimeSinceLastPressed[k] = 0
	}

	entityDynamicID = entityUID

}
