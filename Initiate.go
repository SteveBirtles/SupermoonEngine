package main

import (
	"time"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
	"os"
	"encoding/gob"
	"math"
)

func floor(x float64) float64 {
	return math.Floor(x)
}

const gridCentre = 128
const outsideGrid = gridCentre + 1
const maxUndo = 10000
const clipboardSize = 64

var (
	windowTitlePrefix = "Go Pixel & Lua Test"
	frames            = 0
	undoFrame         = 1
	second            = time.Tick(time.Second)
	win               *pixelgl.Window
	tilePic           pixel.Picture
	tileSprite        [17]*pixel.Sprite
	playerSprite      [12]*pixel.Sprite
	grid            [2*gridCentre][2*gridCentre][16][2]uint16
	clipboard       [10][clipboardSize][clipboardSize][16][2]uint16
	clipboardWidth  [10]int
	clipboardHeight [10]int
	currentClipboard = 1
	previewClipboard = -1
	clobber 		 = true
	undo            [maxUndo][6]int //0 frame,  1 x,  2 y,  3 z,  4 base,  5 front
	undoCounter      = 0
	scale            = 0.5
	aspect           = 0.5
	hScale           = 64.0
	vScale           = hScale * aspect
	lastTileX        = outsideGrid
	lastTileY        = 0
	selectedTile1   uint16 = 4
	selectedTile2   uint16 = 0
	cameraX          = 0.0 //128.0*gridCentre
	cameraY          = 0.0 //128.0*gridCentre
	mouseX           = 0.0
	mouseY           = 0.0
	tileX            = 0
	tileY            = 0
	tileZ            = 0
	showGrid         = 1
	xPressed         = false
	zPressed         = false
	selectionStartX  = 0
	selectionStartY  = 0
	selectionEndX    = 0
	selectionEndY    = 0
	selectionLive = false
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const levelFile = "resources/level.dat"
const clipboardFile = "resources/clipboards.dat"


func save() {
	f1, err := os.Create(levelFile)
	check(err)
	encoder1 := gob.NewEncoder(f1)
	check(encoder1.Encode(grid))

	f2, err := os.Create(clipboardFile)
	check(err)
	encoder2 := gob.NewEncoder(f2)
	check(encoder2.Encode(clipboard))
	check(encoder2.Encode(clipboardWidth))
	check(encoder2.Encode(clipboardHeight))
}

func load() {
	f1, err := os.Open(levelFile)
	if err == nil {
		decoder1 := gob.NewDecoder(f1)
		check(decoder1.Decode(&grid))
	}

	f2, err := os.Open(clipboardFile)
	if err == nil {
		decoder2 := gob.NewDecoder(f2)
		check(decoder2.Decode(&clipboard))
		check(decoder2.Decode(&clipboardWidth))
		check(decoder2.Decode(&clipboardHeight))
	}
}

func backup() {
	f, err := os.Create(levelFile + ".bak")
	if err == nil {
		encoder := gob.NewEncoder(f)
		check(encoder.Encode(grid))
	}
}

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




}