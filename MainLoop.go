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


const gridCentre = 128
const outsideGrid = gridCentre + 1

var (
	grid [2*gridCentre][2*gridCentre][16][2]int
	scale = 0.5
	aspect = 0.5
	hScale = 64.0
	vScale = hScale * aspect
	lastTileX = outsideGrid
	lastTileY = 0
	selectedTile1 = 4
	selectedTile2 = 0
	cameraX = 0.0 //128.0*gridCentre
	cameraY = 0.0 //128.0*gridCentre
	mouseX = 0.0
	mouseY = 0.0
	tileX = 0
	tileY = 0
	tileZ = 0
	showGrid = true
	xPressed = false
	zPressed = false
)

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

	cameraX = -float64(tileX)*128 - 64
	cameraY = float64(tileY)*128 + 64
	tileX = int(floor((mouseX - cameraX*scale) / hScale))
	tileY = int(floor((mouseY + cameraY*scale*aspect) / vScale))

	tileOverlay := pixel.NewBatch(&pixel.TrianglesData{}, tilePic)
	imd := imdraw.New(nil)
	batch := pixel.NewBatch(&pixel.TrianglesData{}, tilePic)

	for !win.Closed() {

		if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyQ) {
			break
		}

		if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyG) {
			showGrid = !showGrid
			xPressed = false
			zPressed = false
		}

		if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyN) {
			grid = [2 * gridCentre][2 * gridCentre][16][2]int{}
		}

		if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyS) {
			save()
		}

		if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyL) {
			load()
		}

		leftAltPressed := win.Pressed(pixelgl.KeyLeftAlt)
		rightAltPressed := win.Pressed(pixelgl.KeyRightAlt)
		backspacePressed := win.Pressed(pixelgl.KeyBackspace)

		if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyX) {
			zPressed = !zPressed
			showGrid = true
			xPressed = false
		}

		if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyZ) {
			xPressed = !xPressed
			showGrid = true
			zPressed = false
		}


		if backspacePressed {

			selectedTile2 = 0

		} else if leftAltPressed {

			selectedTile1 = int(16*win.MousePosition().X/screenWidth) + 1
			if selectedTile1 < 1 {
				selectedTile1 = 1
			}
			if selectedTile1 > 16 {
				selectedTile1 = 16
			}

		} else if rightAltPressed {

			selectedTile2 = int(16*win.MousePosition().X/screenWidth) + 1
			if selectedTile2 < 1 {
				selectedTile2 = 1
			}
			if selectedTile2 > 16 {
				selectedTile2 = 16
			}

		} else {

			mouseX = float64(win.MousePosition().X - screenWidth/2)
			mouseY = float64(screenHeight/2 - win.MousePosition().Y)

			tileX = int(floor((mouseX - cameraX*scale) / hScale))
			tileY = int(floor((mouseY + cameraY*scale*aspect) / vScale))

			if win.MouseScroll().Y != 0 {
				lastScale := scale
				scale /= 1 - win.MouseScroll().Y/10
				if scale < 0.1 {
					scale = 0.1
				}
				if scale > 2.0 {
					scale = 2.0
				}
				hScale = 128 * scale
				vScale = 128 * aspect * scale

				deltaX := cameraX - (-float64(tileX)*128 - 64)
				deltaY := cameraY - (+float64(tileY)*128 + 64)

				cameraX -= deltaX * (1 - lastScale/scale)
				cameraY -= deltaY * (1 - lastScale/scale)

			}

			if win.Pressed(pixelgl.KeyRightBracket) {
				aspect += 0.01
				if aspect > 1.0 {
					aspect = 1.0
				}
				vScale = hScale * aspect
			} else if win.Pressed(pixelgl.KeyLeftBracket) {
				aspect -= 0.01
				if aspect < 0.5 {
					aspect = 0.5
				}
				vScale = hScale * aspect
			}

			if win.JustPressed(pixelgl.KeyPageUp) {
				tileZ += 1
				if tileZ > 15 { tileZ = 15 }
			} else if win.JustPressed(pixelgl.KeyPageDown) {
				tileZ -= 1
				if tileZ < 0 { tileZ = 0 }
			} else if win.JustPressed(pixelgl.KeyHome) {
				tileZ = 15
			} else if win.JustPressed(pixelgl.KeyEnd) {
				tileZ = 0
			}

			if win.Pressed(pixelgl.KeyW) {
				cameraY -= 10 / scale
			}
			if win.Pressed(pixelgl.KeyS) {
				cameraY += 10 / scale
			}
			if win.Pressed(pixelgl.KeyD) {
				cameraX -= 10 / scale
			}
			if win.Pressed(pixelgl.KeyA) {
				cameraX += 10 / scale
			}

		}

		onGrid := tileX >= -gridCentre && tileY >= -gridCentre && tileX < gridCentre && tileY < gridCentre

		leftDown := win.Pressed(pixelgl.MouseButtonLeft) || win.Pressed(pixelgl.KeySpace)
		rightDown := win.Pressed(pixelgl.MouseButtonRight) || win.Pressed(pixelgl.KeyDelete)
		middleDown := win.JustPressed(pixelgl.MouseButtonMiddle)

		if onGrid {
			if middleDown {

				t := grid[tileX+gridCentre][tileY+gridCentre][tileZ][0]
				if t > 0 {
					selectedTile1 = t
				}
				t = grid[tileX+gridCentre][tileY+gridCentre][tileZ][1]
				if t > 0 {
					selectedTile2 = t
				}

			} else if leftDown || rightDown {

				newValue1 := selectedTile1
				newValue2 := selectedTile2
				if rightDown {
					newValue1 = 0
					newValue2 = 0
				}

				if lastTileX != outsideGrid {

					if math.Abs(float64(tileX-lastTileX)) > 1 || math.Abs(float64(tileY-lastTileY)) > 1 {
						d := 1.0 / float64(math.Abs(float64(lastTileX-tileX))+math.Abs(float64(lastTileY-tileY)))
						if d > 0 && d < 100 {
							dx := float64(lastTileX - tileX)
							dy := float64(lastTileY - tileY)
							for s := 0.0; s < 1.0; s += d {
								grid[tileX+int(s*dx)+gridCentre][tileY+int(s*dy)+gridCentre][tileZ][0] = newValue1
								grid[tileX+int(s*dx)+gridCentre][tileY+int(s*dy)+gridCentre][tileZ][1] = newValue2
							}
						}
					}

				}

				grid[tileX+gridCentre][tileY+gridCentre][tileZ][0] = newValue1
				grid[tileX+gridCentre][tileY+gridCentre][tileZ][1] = newValue2

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

		imd.Clear()
		batch.Clear()

		for i := -iRange + iOffset; i <= iRange+iOffset; i++ {
			for j := -jRange + jOffset; j <= jRange+jOffset; j++ {
				for k := 0.0; k < 16; k++ {

					if (zPressed && int(k) != tileZ) || (xPressed && int(j) != tileY) {
						continue
					}


					if int(i) >= -gridCentre && int(j) >= -gridCentre && int(i) < gridCentre && int(j) < gridCentre {

						baseTile := grid[int(i)+gridCentre][int(j)+gridCentre][int(k)][0]

						if baseTile > 0 || (selectedTile1 > 0 && int(i) == tileX && int(j) == tileY && int(k) == tileZ) {

							s := 4*(1-aspect)

							cam := pixel.V(cameraX, cameraY)
							pos := pixel.V(screenWidth/2+float64(i*hScale)+hScale/2, screenHeight/2+(-vScale/2-float64((j-k*s)*vScale)))

							frontTile := grid[int(i)+gridCentre][int(j)+gridCentre][int(k)][1]

							if frontTile > 0 || (selectedTile2 > 0 && int(i) == tileX && int(j) == tileY && int(k) == tileZ) {

								matrix := pixel.IM.Moved(cam).ScaledXY(pixel.ZV, pixel.V(scale, scale*aspect)).Moved(pos).
									Moved(pixel.V(0, vScale*(1-aspect)*4))

								if int(i) == tileX && int(j) == tileY && int(k) == tileZ {
									tileSprite[selectedTile1-1].Draw(batch, matrix)
								} else {
									tileSprite[baseTile-1].Draw(batch, matrix)
								}

								if aspect < 1 {

									matrix := pixel.IM.
										ScaledXY(pixel.ZV, pixel.V(1, s)).
										Moved(cam).
										ScaledXY(pixel.ZV, pixel.V(scale, scale*aspect)).
										Moved(pos).
										Moved(pixel.V(0, vScale/2 - 2*(aspect-0.5)*vScale))

									if selectedTile2 > 0 && int(i) == tileX && int(j) == tileY && int(k) == tileZ {
										tileSprite[selectedTile2-1].Draw(batch, matrix)
									} else {
										tileSprite[frontTile-1].Draw(batch, matrix)
									}
								}


							} else {

								matrix := pixel.IM.Moved(cam).ScaledXY(pixel.ZV, pixel.V(scale, scale*aspect)).Moved(pos)

								if int(i) == tileX && int(j) == tileY && int(k) == tileZ {
									tileSprite[selectedTile1-1].Draw(batch, matrix)
								} else {
									tileSprite[baseTile-1].Draw(batch, matrix)
								}

							}

						}

					}

					if k == 0 && showGrid && !(xPressed || zPressed) {

						cam := pixel.V(cameraX, cameraY)
						pos := pixel.V(screenWidth/2+float64(i*hScale)+hScale/2, screenHeight/2+(-vScale/2-float64(j*vScale)))

						matrix := pixel.IM.Moved(cam).ScaledXY(pixel.ZV, pixel.V(scale, scale*aspect)).Moved(pos)
						imd.SetMatrix(matrix)

						gridIntensity := math.Sqrt(scale / 2)

						if (int(i) == tileX || int(i) == tileX+1) && int(j) == tileY {
							imd.Color = pixel.RGB(255, 255, 255)
						} else if int(i) == 0 || int(i) == 1 {
							imd.Color = pixel.RGB(gridIntensity, gridIntensity, 0)
						} else if int(i)+gridCentre < 0 || int(i)+gridCentre > 2*gridCentre-1 ||
							int(j)+gridCentre <= 0 || int(j)+gridCentre > 2*gridCentre-1 {
							imd.Color = pixel.RGB(gridIntensity, 0, 0)
						} else {
							imd.Color = pixel.RGB(0, gridIntensity, 0)
						}

						imd.Push(pixel.V(-64, -64))
						imd.Push(pixel.V(-64, 64))
						imd.Line(1.0 / scale)

						if int(i) == tileX && (int(j) == tileY || int(j) == tileY-1) {
							imd.Color = pixel.RGB(255, 255, 255)
						} else if int(j) == 0 || int(j) == -1 {
							imd.Color = pixel.RGB(gridIntensity, gridIntensity, 0)
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
		}

		win.Clear(colornames.Black)

		win.SetComposeMethod(pixel.ComposeOver)
		imd.Draw(win)

		win.SetComposeMethod(pixel.ComposeOver)
		batch.Draw(win)

		if leftAltPressed || rightAltPressed {

			tileOverlay.Clear()

			for i := 0; i < 16; i++ {

				if leftAltPressed && i == selectedTile1-1 {
					tileOverlay.SetColorMask(color.RGBA{R: 255, G: 255, B: 255, A: 255})
				} else if rightAltPressed && i == selectedTile2-1 {
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
