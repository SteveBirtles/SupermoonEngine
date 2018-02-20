package main

import (
	"golang.org/x/image/colornames"
	"fmt"
	"github.com/faiface/pixel"
	"time"
	"engo/math"
)

var (
	tick = time.Tick(time.Millisecond * 100)
	playerFrame = 0
	playerX = 0.0
)

func mainLoop() {

	initiate()

	for !win.Closed() {

		//executeLuaFile(L, "sonic.lua")

		//angleLua := executeLuaFunction(L, "getAngle", []lua.LValue{})
		//angle, _ := strconv.ParseFloat(angleLua.String(), 64)
		//matrix := pixel.IM.Rotated(pixel.ZV, angle).Scaled(pixel.ZV, 0.2).Moved(pixel.Vec{X: x, Y: y})

		win.Clear(colornames.Black)
		win.SetComposeMethod(pixel.ComposeOver)


		for i := 0; i < 16; i++ {
			for j:= 0; j < 24; j++ {
				matrix := pixel.IM.ScaledXY(pixel.ZV, pixel.V(0.5,0.25)).Moved(pixel.V(float64(i*64)+32, 768-16-float64(j*32)))
				tileSprite[0].Draw(win, matrix)
			}
		}

		for i := 4; i < 13; i++ {

			for j:= 12; j > 2*int(math.Abs(float32(i) - 8)) + 2; j-=2 {

				matrix := pixel.IM.ScaledXY(pixel.ZV, pixel.V(0.5, 0.25)).Moved(pixel.V(float64((i+2)*64)+32, 768-16-float64(j*32)))
				tileSprite[8].Draw(win, matrix)

				matrix = pixel.IM.ScaledXY(pixel.ZV, pixel.V(0.5, 0.5)).Moved(pixel.V(float64((i+2)*64)+32, 768-32-float64((j+1)*32)))
				tileSprite[7].Draw(win, matrix)

			}

		}

		matrix := pixel.IM.Scaled(pixel.ZV, 0.9).Moved( pixel.V(playerX,338))
		playerSprite[playerFrame].Draw(win, matrix)

		for i := 4; i < 13; i++ {

			for j:= 14; j > 2*int(math.Abs(float32(i) - 8)) + 4; j-=2 {

				matrix := pixel.IM.ScaledXY(pixel.ZV, pixel.V(0.5, 0.25)).Moved(pixel.V(float64((i-3)*64)+32, 768-16-float64(j*32)))
				tileSprite[5].Draw(win, matrix)

				matrix = pixel.IM.ScaledXY(pixel.ZV, pixel.V(0.5, 0.5)).Moved(pixel.V(float64((i-3)*64)+32, 768-32-float64((j+1)*32)))
				tileSprite[6].Draw(win, matrix)

			}

		}

		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", windowTitlePrefix, frames))
			frames = 0
			default:
		}

		select {
		case <-tick:
			playerFrame = (playerFrame + 1) % 6
			playerX += 16
			if playerX > 1050 {
				playerX = -50
			}
		default:
		}

	}
}
