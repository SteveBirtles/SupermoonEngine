package main

import (
	"golang.org/x/image/colornames"
	"fmt"
	"github.com/faiface/pixel"
	"engo/math"
	"github.com/faiface/pixel/pixelgl"
)

func mainLoop() {

	initiate()

	var (
		scale = 0.5
		hScale = 64.0
		vScale = 32.0
		lastTileX = -9000
		lastTileY = 0
	)

	gridCentre := 512
	var grid [1024][1024]int

	for !win.Closed() {

		mouseX := float64(win.MousePosition().X - screenWidth/2)
		mouseY := float64(screenHeight/2 - win.MousePosition().Y)

		if win.MouseScroll().Y != 0 {
			scale /= 1 - win.MouseScroll().Y/10
			if scale < 0.2 { scale = 0.2 }
			if scale > 2.0 { scale = 2.0 }
			hScale = 128 * scale
			vScale = 64 * scale
		}

		cursorX := float64(math.Floor(float32(mouseX / hScale))) * hScale
		cursorY := float64(math.Floor(float32(mouseY / vScale))) * vScale

		tileX := int(cursorX / hScale)
		tileY := int(cursorY / vScale)
		onGrid := tileX > -gridCentre && tileY > -gridCentre && tileX < gridCentre && tileY < gridCentre

		leftDown := win.Pressed(pixelgl.MouseButtonLeft)
		rightDown := win.Pressed(pixelgl.MouseButtonRight)
		middleDown := win.Pressed(pixelgl.MouseButtonMiddle)

		if (leftDown || rightDown || middleDown) && onGrid {

			newValue := 3
			if rightDown { newValue = 0 }

			if lastTileX != -9000 {

				if math.Abs(float32(tileX - lastTileX)) > 1 || math.Abs(float32(tileY - lastTileY)) > 1 {

					d := 1.0 / float64(math.Abs(float32(lastTileX-tileX))+math.Abs(float32(lastTileY-tileY)))

					if d > 0 && d < 100 {

						dx := float64(lastTileX - tileX)
						dy := float64(lastTileY - tileY)

						for s := 0.0; s < 1.0; s += d {
							grid[tileX+int(s*dx)+gridCentre][tileY+int(s*dy)+gridCentre] = newValue
						}
					}

				}

			}

			grid[tileX+gridCentre][tileY+gridCentre] = newValue

			lastTileX = tileX
			lastTileY = tileY

		} else {
			lastTileX = -9000
		}

		win.Clear(colornames.Black)
		win.SetComposeMethod(pixel.ComposeOver)

		iRange := float64(math.Floor(float32(screenWidth/(2*hScale)))) + 1
		jRange := float64(math.Floor(float32(screenHeight/(2*vScale)))) + 1

		for i := -iRange; i < iRange; i++ {
			for j:= -jRange; j < jRange; j++ {

				if int(i) > -gridCentre && int(j) > -gridCentre && int(i) < gridCentre && int(j) < gridCentre {

					tileNo := grid[int(i)+gridCentre][int(j)+gridCentre]

					matrix := pixel.IM.ScaledXY(pixel.ZV, pixel.V(scale, scale/2)).Moved(pixel.V(screenWidth/2+float64(i*hScale)+vScale, screenHeight/2+(-vScale/2-float64(j*vScale))))
					tileSprite[tileNo].Draw(win, matrix)

				}
			}
		}

		matrix := pixel.IM.ScaledXY(pixel.ZV, pixel.V(scale,scale/2)).Moved(pixel.V(screenWidth/2 + cursorX + hScale/2, screenHeight/2 - (cursorY + vScale/2)))
		tileSprite[16].Draw(win, matrix)


		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d | X: %d | Y: %d", windowTitlePrefix, frames, tileX, tileY))
			frames = 0
			default:
		}

	}
}
