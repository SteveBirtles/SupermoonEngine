package main

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"io/ioutil"
)

func initiateEngine() {

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

	luaRenderer = text.New(pixel.ZV, textAtlas)
	luaRenderer.LineHeight = textAtlas.LineHeight() * 1.5
	luaRenderer.Color = colornames.Orangered
	luaRenderer.Orig = pixel.V(screenWidth*0.75, screenHeight-22)

	tileOverlayWidth  = uint16((screenWidth) / 75)
	tileOverlayHeight = totalTiles /tileOverlayWidth

	textureImage, initError := loadImageFile("textures/tiles.png")
	check(initError)

	tilePic = pixel.PictureDataFromImage(textureImage)
	n := 0
	for j := 0; j < tileSheetHeight/128; j++ {
		for i := 0; i < tileSheetWidth/128; i++ {
			if n >= totalTiles { panic("Tiles error!") }
			u := float64(i*128)
			v := float64(tileSheetHeight -(j+1)*128)
			tileTexture[n] = pixel.NewSprite(tilePic, pixel.Rect{Min: pixel.V(u, v), Max: pixel.V(u + 128, v + 128)})
			n++
		}
	}

	spriteImage, initError := loadImageFile("textures/sprites.png")
	check(initError)

	spritePic = pixel.PictureDataFromImage(spriteImage)
	n = 0
	for j := 0; j < spriteSheetHeight/128; j++ {
		for i := 0; i < spriteSheetWidth/128; i++ {
			if n >= totalSprites { panic("Sprites error!") }
			u := float64(i*128)
			v := float64(spriteSheetHeight -(j+1)*128)
			spriteTexture[n] = pixel.NewSprite(spritePic, pixel.Rect{Min: pixel.V(u, v), Max: pixel.V(u + 128, v + 128)})
			n++
		}
	}

	tileOverlay = pixel.NewBatch(&pixel.TrianglesData{}, tilePic)
	imGrid = imdraw.New(nil)
	tileBatch = pixel.NewBatch(&pixel.TrianglesData{}, tilePic)
	spriteBatch = pixel.NewBatch(&pixel.TrianglesData{}, spritePic)
	imUI = imdraw.New(nil)

	scripts, err := ioutil.ReadDir("scripts/")
	check(err)

	b := 0
	entityClassBlockCount = 1
	for _, f := range scripts {
		s := f.Name()
		if s[len(s)-4:] != ".lua" { continue }
		lastEntityClass++
		entityClass = append(entityClass, s[:len(s)-4])
		b++
		if b > 12 {
			entityClassBlockCount++
			b = 1
		}
	}

}
