package main

import (
	"golang.org/x/image/colornames"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math"
)

func mainLoop() {

	initiate()

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
		cameraX = 0.0
		cameraY = 0.0
	)

	for !win.Closed() {

		mouseX := float64(win.MousePosition().X - screenWidth/2)
		mouseY := float64(screenHeight/2 - win.MousePosition().Y)

		if win.MouseScroll().Y != 0 {
			scale /= 1 - win.MouseScroll().Y/10
			if scale < 0.2 { scale = 0.2 }
			if scale > 2.0 { scale = 2.0 }
			hScale = 128 * scale
			vScale = 128 * aspect * scale
		}

		cursorX := float64(math.Floor(mouseX / hScale)) * hScale
		cursorY := float64(math.Floor(mouseY / vScale)) * vScale

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

		if win.Pressed(pixelgl.KeyW) { cameraY -= 25 }
		if win.Pressed(pixelgl.KeyS) { cameraY += 25 }
		if win.Pressed(pixelgl.KeyD) { cameraX -= 25 }
		if win.Pressed(pixelgl.KeyA) { cameraX += 25 }

		tileX := int((cursorX - cameraX*scale) / hScale)
		tileY := int((cursorY + cameraY*scale) / vScale)
		onGrid := tileX > -gridCentre && tileY > -gridCentre && tileX < gridCentre && tileY < gridCentre

		leftDown := win.Pressed(pixelgl.MouseButtonLeft)
		rightDown := win.Pressed(pixelgl.MouseButtonRight)
		middleDown := win.Pressed(pixelgl.MouseButtonMiddle)

		if onGrid {
			if middleDown {

				selectedTile = grid[tileX+gridCentre][tileY+gridCentre][0]

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

		win.Clear(colornames.Black)
		win.SetComposeMethod(pixel.ComposeOver)

		iRange := float64(math.Floor(screenWidth/(2*hScale))) + 1
		jRange := float64(math.Floor(screenHeight/(2*vScale))) + 1

		for i := -iRange; i < iRange; i++ {
			for j:= -jRange; j < jRange ; j++ {

				if int(i) > -gridCentre && int(j) > -gridCentre && int(i) < gridCentre && int(j) < gridCentre {

					tileNo := grid[int(i)+gridCentre][int(j)+gridCentre][0]


					matrix := pixel.IM.
						Moved(pixel.V(cameraX, cameraY)).
						ScaledXY(pixel.ZV, pixel.V(scale, scale*aspect)).
						Moved(pixel.V(screenWidth/2+float64(i*hScale)+hScale/2, screenHeight/2+(-vScale/2-float64(j*vScale))))

					if tileNo > 0 {

						tileSprite[tileNo-1].Draw(win, matrix)

					} else {

						tileSprite[16].Draw(win, matrix)

					}


				}
			}
		}

		matrix := pixel.IM.
			Moved(pixel.V(cameraX, cameraY)).
			ScaledXY(pixel.ZV, pixel.V(scale,scale*aspect)).
			Moved(pixel.V(screenWidth/2+float64(float64(tileX)*hScale)+hScale/2, screenHeight/2+(-vScale/2-float64(float64(tileY)*vScale))))
		tileSprite[selectedTile-1].Draw(win, matrix)
		tileSprite[16].Draw(win, matrix)

		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d | X: %d | Y: %d | Aspect: %d%% | Camera: %d, %d", windowTitlePrefix, frames, tileX, tileY, int(100*(1-aspect)), int(cameraX), int(cameraY)))
			frames = 0
			default:
		}

	}
}
