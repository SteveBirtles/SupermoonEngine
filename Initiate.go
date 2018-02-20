package main

import (
	"time"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
)

var (
	windowTitlePrefix  = "Go Pixel & Lua Test"
	frames                        = 0
	second                        = time.Tick(time.Second)
	win               *pixelgl.Window
	tileSprite        [16]*pixel.Sprite
	playerSprite      [12]*pixel.Sprite
)

func initiate() {

	var initError error

	cfg := pixelgl.WindowConfig{
		Bounds: pixel.R(0, 0, screenWidth, screenHeight),
		VSync:  true,
	}

	win, initError = pixelgl.NewWindow(cfg)
	if initError != nil {
		panic(initError)
	}


	spriteImage, initError := loadImageFile("textures/blocks.png")
	if initError != nil {
		panic(initError)
	}

	playerImage, initError := loadImageFile("textures/player.png")
	if initError != nil {
		panic(initError)
	}


	spritePic := pixel.PictureDataFromImage(spriteImage)
	for i := 0; i < 16; i++ {
		tileSprite[i] = pixel.NewSprite(spritePic, pixel.Rect{Min: pixel.V(0,float64((15-i)*128)), Max: pixel.V(128,128+float64((15-i)*128))})
	}

	playerPic := pixel.PictureDataFromImage(playerImage)

	for j := 0; j < 2; j++ {
		for i := 0; i < 6; i++ {
			playerSprite[i+j*6] = pixel.NewSprite(playerPic, pixel.Rect{Min: pixel.V(float64(i*104), float64((1-j)*151)), Max: pixel.V(float64((i+1)*104), 151+float64((1-j)*151))})
		}
	}




}