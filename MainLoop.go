package main

import (
	"golang.org/x/image/colornames"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math"
	"github.com/faiface/pixel/imdraw"
	"os"
	"encoding/gob"
	"image/color"
)

func floor(x float64) float64 {
	return math.Floor(x)
}


const gridCentre = 512
const outsideGrid = gridCentre + 1

var (
	grid [2*gridCentre][2*gridCentre][2]int
	scale = 0.5
	aspect = 1.0
	hScale = 64.0
	vScale = hScale * aspect
	lastTileX = outsideGrid
	lastTileY = 0
	selectedTile = 4
	cameraX = 0.0 //128.0*gridCentre
	cameraY = 0.0 //128.0*gridCentre
	ignoreMouse = false
	lastMouseX = 0.0
	lastMouseY = 0.0
	mouseX = 0.0
	mouseY = 0.0
	tileX = 0
	tileY = 0
)

func centreCursor() {
	cameraX = -float64(tileX)*128 - 64
	cameraY = float64(tileY)*128 + 64
	lastMouseX = float64(win.MousePosition().X - screenWidth/2)
	lastMouseY = float64(screenHeight/2 - win.MousePosition().Y)
	mouseX = 0.0
	mouseY = 0.0
	ignoreMouse = true
	win.SetCursorVisible(false)
	tileX = int(floor((mouseX - cameraX*scale) / hScale))
	tileY = int(floor((mouseY + cameraY*scale*aspect) / vScale))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const levelFile = "resources/level.dat"

func save() {
	f, err := os.Create(levelFile)
	check(err)
	encoder := gob.NewEncoder(f)
	check(encoder.Encode(grid))
}

func load() {
	f, err := os.Open(levelFile)
	if err == nil {
		decoder := gob.NewDecoder(f)
		check(decoder.Decode(&grid))
	}
}

func backup() {
	f, err := os.Create(levelFile + ".bak")
	if err == nil {
		encoder := gob.NewEncoder(f)
		check(encoder.Encode(grid))
	}
}

func mainLoop() {

	initiate()
	load()
	backup()

	centreCursor()

	for !win.Closed() {

		if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyQ) {
			break
		}

		if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyN) {
			grid = [2*gridCentre][2*gridCentre][2]int{}
		}

		if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyS) {
			save()
		}

		if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyL) {
			load()
		}

		altPressed := win.Pressed(pixelgl.KeyLeftAlt)

		if altPressed {

			selectedTile = int(16 * win.MousePosition().X / screenWidth) + 1
			if selectedTile < 1 { selectedTile = 1 }
			if selectedTile > 16 { selectedTile = 16 }

		} else {

			if ignoreMouse {
				if lastMouseX != float64(win.MousePosition().X-screenWidth/2) ||
					lastMouseY != float64(screenHeight/2-win.MousePosition().Y) {
					ignoreMouse = false
					win.SetCursorVisible(true)
				}
			}

			if !ignoreMouse {
				mouseX = float64(win.MousePosition().X - screenWidth/2)
				mouseY = float64(screenHeight/2 - win.MousePosition().Y)
			}

			tileX = int(floor((mouseX - cameraX*scale) / hScale))
			tileY = int(floor((mouseY + cameraY*scale*aspect) / vScale))

			if win.MouseScroll().Y != 0 {
				scale /= 1 - win.MouseScroll().Y/10
				if scale < 0.1 {
					scale = 0.1
				}
				if scale > 2.0 {
					scale = 2.0
				}
				hScale = 128 * scale
				vScale = 128 * aspect * scale
				if win.MouseScroll().Y > 0 {
					centreCursor()
				}
			}


			if win.JustPressed(pixelgl.KeyMinus) {
				selectedTile++
				if selectedTile > 16 { selectedTile = 1 }
			} else if win.JustPressed(pixelgl.KeyEqual) {
				selectedTile--
				if selectedTile < 1 { selectedTile = 16 }
			}

			if win.JustPressed(pixelgl.KeyPageDown) {
				aspect += 0.1
				if aspect > 1.0 { aspect = 1.0 }
				vScale = hScale * aspect
			} else if win.JustPressed(pixelgl.KeyPageUp) {
				aspect -= 0.1
				if aspect < 0.5 { aspect = 0.5 }
				vScale = hScale * aspect
			}

			if win.Pressed(pixelgl.KeyW) { cameraY -= 10/scale }
			if win.Pressed(pixelgl.KeyS) { cameraY += 10/scale }
			if win.Pressed(pixelgl.KeyD) { cameraX -= 10/scale }
			if win.Pressed(pixelgl.KeyA) { cameraX += 10/scale }

		}


		onGrid := tileX >= -gridCentre && tileY >= -gridCentre && tileX < gridCentre && tileY < gridCentre

		leftDown := win.Pressed(pixelgl.MouseButtonLeft) || win.Pressed(pixelgl.KeySpace)
		rightDown := win.Pressed(pixelgl.MouseButtonRight) || win.Pressed(pixelgl.KeyDelete)
		middleDown := win.JustPressed(pixelgl.MouseButtonMiddle) || win.Pressed(pixelgl.KeyTab)

		if onGrid {
			if middleDown {

				t := grid[tileX+gridCentre][tileY+gridCentre][0]
				if t > 0 { selectedTile = t }
				centreCursor()

			} else if leftDown || rightDown {

				newValue := selectedTile
				if rightDown {
					newValue = 0
				}

				if lastTileX != outsideGrid {

					if math.Abs(float64(tileX-lastTileX)) > 1 || math.Abs(float64(tileY-lastTileY)) > 1 {
						d := 1.0 / float64(math.Abs(float64(lastTileX-tileX))+math.Abs(float64(lastTileY-tileY)))
						if d > 0 && d < 100 {
							dx := float64(lastTileX - tileX)
							dy := float64(lastTileY - tileY)
							for s := 0.0; s < 1.0; s += d {
								grid[tileX+int(s*dx)+gridCentre][tileY+int(s*dy)+gridCentre][0] = newValue
							}
						}
					}

				}

				grid[tileX+gridCentre][tileY+gridCentre][0] = newValue

				lastTileX = tileX
				lastTileY = tileY

			} else {
				lastTileX = outsideGrid
			}
		}


		iRange := floor(screenWidth/(2*hScale)) + 2
		jRange := floor(screenHeight/(2*vScale)) + 2

		iOffset := -floor(scale * cameraX / hScale)
		jOffset := floor(scale * aspect * cameraY / vScale)

		imd := imdraw.New(nil)
		batch := pixel.NewBatch(&pixel.TrianglesData{}, tilePic)

		for i := -iRange + iOffset; i <= iRange + iOffset; i++ {
			for j:= -jRange + jOffset; j <= jRange + jOffset ; j++ {

				matrix := pixel.IM.
					Moved(pixel.V(cameraX, cameraY)).
					ScaledXY(pixel.ZV, pixel.V(scale, scale*aspect)).
					Moved(pixel.V(screenWidth/2+float64(i*hScale)+hScale/2, screenHeight/2+(-vScale/2-float64(j*vScale))))

					tileDrawn := false

				if int(i) >= -gridCentre && int(j) >= -gridCentre && int(i) < gridCentre && int(j) < gridCentre {
					tileNo := grid[int(i)+gridCentre][int(j)+gridCentre][0]
					if tileNo > 0 {
						tileSprite[tileNo-1].Draw(batch, matrix)
						tileDrawn = true
					}
				}

				if !tileDrawn {

					imd.SetMatrix(matrix)

					gridIntensity := math.Sqrt(scale / 2)

					if int(i) == 0 || int(i) == 1 {
						imd.Color = pixel.RGB(gridIntensity*2, gridIntensity*2, 0)
					} else if int(i)+gridCentre < 0 || int(i)+gridCentre > 2*gridCentre-1 ||
						int(j)+gridCentre <= 0 || int(j)+gridCentre > 2*gridCentre-1 {
						imd.Color = pixel.RGB(gridIntensity, 0, 0)
					} else {
						imd.Color = pixel.RGB(0, gridIntensity, 0)
					}

					imd.Push(pixel.V(-64, -64))
					imd.Push(pixel.V(-64, 64))
					imd.Line(1.0 / scale)

					if int(j) == 0 || int(j) == -1 {
						imd.Color = pixel.RGB(gridIntensity*2, gridIntensity*2, 0)
					} else if int(i)+gridCentre < 0 || int(i)+gridCentre >= 2*gridCentre-1 ||
						int(j)+gridCentre < 0 || int(j)+gridCentre > 2*gridCentre-1 {
						imd.Color = pixel.RGB(gridIntensity, 0, 0)
					} else {
						imd.Color = pixel.RGB(0, gridIntensity, 0)
					}

					imd.Push(pixel.V(-64, -64))
					imd.Push(pixel.V(64, -64))
					imd.Line(2.0 / scale)

				}

			}
		}

		win.Clear(colornames.Black)

		win.SetComposeMethod(pixel.ComposeOver)
		batch.Draw(win)

		if !altPressed {
			matrix := pixel.IM.
				Moved(pixel.V(cameraX, cameraY)).
				ScaledXY(pixel.ZV, pixel.V(scale, scale*aspect)).
				Moved(pixel.V(screenWidth/2+float64(float64(tileX)*hScale)+hScale/2, screenHeight/2+(-vScale/2-float64(float64(tileY)*vScale))))
			tileSprite[selectedTile-1].Draw(win, matrix)
			tileSprite[16].Draw(win, matrix)
		}

		win.SetComposeMethod(pixel.ComposeOver)
		imd.Draw(win)

		if altPressed {

			tileOverlay := pixel.NewBatch(&pixel.TrianglesData{}, tilePic)

			for i := 0; i < 16; i++ {

				if i == selectedTile - 1 {
					tileOverlay.SetColorMask(color.RGBA{R: 255, G: 255, B: 255, A: 255})
				} else {
					tileOverlay.SetColorMask(color.RGBA{R: 128, G: 128, B: 128, A: 128})
				}

				matrix := pixel.IM.Moved(pixel.V(float64(i)*150+150, 100)).ScaledXY(pixel.ZV, pixel.V(0.5, 0.5))
				tileSprite[i].Draw(tileOverlay, matrix)
			}

			win.SetComposeMethod(pixel.ComposeOver)
			tileOverlay.Draw(win)
		}

		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d | X: %d | Y: %d | Aspect: %d%% | Camera: %d, %d", windowTitlePrefix, frames, tileX, tileY, int(100*(1-aspect)), int(cameraX), int(cameraY)))
			frames = 0
			default:
		}

	}

	save()

}
