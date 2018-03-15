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

	tileOverlayWidth  = uint16((screenWidth) / 75)
	tileOverlayHeight = superTiles/tileOverlayWidth

	textureImage, initError := loadImageFile("textures/super.png")
	check(initError)

	tilePic = pixel.PictureDataFromImage(textureImage)
	n := 0
	for j := 0; j < superHeight/128; j++ {
		for i := 0; i < superWidth/128; i++ {
			if n >= superTiles { panic("Super texture error!") }
			u := float64(i*128)
			v := float64(superHeight-(j+1)*128)
			tileTexture[n] = pixel.NewSprite(tilePic, pixel.Rect{Min: pixel.V(u, v), Max: pixel.V(u + 128, v + 128)})
			n++
		}
	}

	tileOverlay = pixel.NewBatch(&pixel.TrianglesData{}, tilePic)
	imd1 = imdraw.New(nil)
	batch = pixel.NewBatch(&pixel.TrianglesData{}, tilePic)
	imd2 = imdraw.New(nil)

}