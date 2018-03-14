package main

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

func initiate() {

	var initError error

	cfg := pixelgl.WindowConfig{
		Bounds: pixel.R(0, 0, screenWidth, screenHeight),
		VSync:  true,
	}

	win, initError = pixelgl.NewWindow(cfg)
	check(initError)

	textFace, initError := loadTTF("resources/font.ttf", 14)
	check(initError)

	textAtlas := text.NewAtlas(textFace, text.ASCII)

	textRenderer = text.New(pixel.ZV, textAtlas)
	textRenderer.LineHeight = textAtlas.LineHeight() * 1.5
	textRenderer.Color = colornames.Limegreen
	textRenderer.Orig = pixel.V(10, screenHeight-22)

	spriteImage, initError := loadImageFile("textures/blocks.png")
	check(initError)

	playerImage, initError := loadImageFile("textures/player.png")
	check(initError)

	tilePic = pixel.PictureDataFromImage(spriteImage)
	for i := 0; i <= 16; i++ {
		tileSprite[i] = pixel.NewSprite(tilePic, pixel.Rect{Min: pixel.V(0,float64((16-i)*128)), Max: pixel.V(128,128+float64((16-i)*128))})
	}

	playerPic := pixel.PictureDataFromImage(playerImage)

	for j := 0; j < 2; j++ {
		for i := 0; i < 6; i++ {
			playerSprite[i+j*6] = pixel.NewSprite(playerPic, pixel.Rect{Min: pixel.V(float64(i*104), float64((1-j)*151)), Max: pixel.V(float64((i+1)*104), 151+float64((1-j)*151))})
		}
	}

	tileOverlay = pixel.NewBatch(&pixel.TrianglesData{}, tilePic)
	imd1 = imdraw.New(nil)
	batch = pixel.NewBatch(&pixel.TrianglesData{}, tilePic)
	imd2 = imdraw.New(nil)

}