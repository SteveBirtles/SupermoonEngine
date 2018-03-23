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

	initiateAPI()

}


func initiateAPI() {

	linkToLua(L, APILoadMap, "LoadMap")
	linkToLua(L, APIGetTile, "GetTile")
	linkToLua(L, APISetTile, "SetTile")
	linkToLua(L, APISetFocus, "SetFocus")
	linkToLua(L, APISetZoom, "SetZoom")

	linkToLua(L, APIGetId, "GetId")
	linkToLua(L, APIGetEntities, "GetEntities")
	linkToLua(L, APISetEntityPosition, "SetEntityPosition")
	linkToLua(L, APIGetEntityPosition, "GetEntityPosition")
	linkToLua(L, APISetEntityVelocity, "SetEntityVelocity")
	linkToLua(L, APIPathFind, "PathFind")
	linkToLua(L, APICreateEntity, "CreateEntity")
	linkToLua(L, APISetEntitySprite, "SetEntitySprite")
	linkToLua(L, APIGetEntitySprite, "GetEntitySprite")
	linkToLua(L, APISetEntityAnimation, "SetEntityAnimation")
	linkToLua(L, APISetEntityScript, "SetEntityScript")
	linkToLua(L, APIGetEntityScript, "GetEntityScript")
	linkToLua(L, APISetEntityProperty, "SetEntityProperty")
	linkToLua(L, APIGetEntityProperty, "GetEntityProperty")
	linkToLua(L, APIListEntityProperties, "ListEntityProperties")
	linkToLua(L, APIDeleteEntity, "DeleteEntity")
	linkToLua(L, APIEntityProximity, "EntityProximity")
	linkToLua(L, APISetFocusEntity, "SetFocusEntity")
	linkToLua(L, APIGetEntityClass, "GetEntityClass")
	linkToLua(L, APISetEntityClass, "SetEntityClass")
	linkToLua(L, APIOverrideClassScript, "OverrideClassScript")
	linkToLua(L, APISetEntityActive, "SetEntityActive")
	linkToLua(L, APISetClassActive, "SetClassActive")
	linkToLua(L, APISetAllActive, "SetAllActive")

	linkToLua(L, APIKeyPressed, "KeyPressed")
	linkToLua(L, APIKeyJustPressed, "KeyJustPressed")
	linkToLua(L, APIDisplayText, "DisplayText")
	linkToLua(L, APIDisplayOptions, "DisplayOptions")
	linkToLua(L, APIPlaySound, "PlaySound")
	linkToLua(L, APIPlayMusic, "PlayMusic")
	linkToLua(L, APIPauseMusic, "PauseMusic")
	linkToLua(L, APIEndGame, "EndGame")

	linkToLua(L, APIStartTimer, "StartTimer")
	linkToLua(L, APIGetTimer, "GetTimer")
	linkToLua(L, APICancelTimer, "CancelTimer")
	linkToLua(L, APISetPersistent, "SetPersistent")
	linkToLua(L, APIGetPersistent, "GetPersistent")

	linkToLua(L, luaPrint, "print")

}
