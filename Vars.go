package main

import (
	"time"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
)

const gridCentre = 128
const outsideGrid = gridCentre + 1
const maxUndo = 10000
const clipboardSize = 64
const clipboardFile = "cache/clipboard.dat"

const tileSheetWidth = 2048
const tileSheetHeight = 1280
const totalTiles = (tileSheetWidth /128)*(tileSheetHeight /128)

const spriteSheetWidth = 2048
const spriteSheetHeight = 3072
const totalSprites = (spriteSheetWidth /128)*(spriteSheetHeight /128)


var (
	windowTitlePrefix  = "Supermoon Engine"
	frameCounter            = 0
	gameFrame                  = 0
	undoFrame                  = 1
	frameLength       float64
	frameStart        time.Time
	second                        = time.Tick(time.Second)
	luaTick           			  = time.Tick(time.Second/10)
	win               *pixelgl.Window
	textRenderer      *text.Text
	textLine          int
	tilePic           pixel.Picture
	tileTexture       [totalTiles]*pixel.Sprite
	tileOverlay       *pixel.Batch
	tileOverlayWidth  uint16
	tileOverlayHeight uint16
	spritePic         pixel.Picture
	spriteTexture     [totalSprites]*pixel.Sprite
	imGrid            *imdraw.IMDraw
	imUI              *imdraw.IMDraw
	tileBatch         *pixel.Batch
	spriteBatch       *pixel.Batch
	grid              [2*gridCentre][2*gridCentre][16][2]uint16
	gridBackup        [2*gridCentre][2*gridCentre][16][2]uint16
	clipboard         [10][clipboardSize][clipboardSize][16][2]uint16
	clipboardWidth    [10]int
	clipboardHeight   [10]int
	clipboardShift                             = 0
	currentClipboard               = 1
	previewClipboard               = -1
	clobber                                             	   = false
	undo              [maxUndo][6]int //0 frame,  1 X,  2 Y,  3 Z,  4 base,  5 front
	undoCounter                          = 0
	scale                               = 0.5
	aspect                             = 0.5
	viewDirection               = 0
	compass                         = [4]string{"North", "East", "South", "West"}
	reverseCompass                         = map[byte]int{'N':0 , 'E':1, 'S':2, 'W':3}
	hScale                           = 64.0
	vScale                                             = hScale * aspect
	lastTileX                                      = outsideGrid
	lastTileY                                      = 0
	selectedTile1 uint16 = 4
	tileRow1      uint16 = 0
	selectedTile2 uint16 = 0
	tileRow2      uint16 = 0
	cameraX                                                = 0.0 //128.0*gridCentre
	cameraY                                                = 0.0 //128.0*gridCentre
	iStart        float64
	jStart        float64
	iEnd          float64
	jEnd          float64
	iOffset       float64
	jOffset       float64
	cameraAdjX    float64
	cameraAdjY    float64
	mouseX                                                     = 0.0
	mouseY                                                     = 0.0
	tileX                                                          = 0
	tileY                                                          = 0
	tileZ                                                          = 0
	hideTile                                           = false
	showShadows                            = true
	showGrid                                           = 1
	xRay                                                               = false
	zRay                                                     = false
	flipX                                        = false
	flipY                                        = false
	selectionStartX          = 0
	selectionStartY         = 0
	selectionEndX             = 0
	selectionEndY             = 0
	selectionLive             = false
	leftAltPressed         = false
	rightAltPressed        = false
	quit                   = 0
	help                   = false
	editing				   = true
	entityDebugData 	   = false

	gameKeyDownStart            map[pixelgl.Button]bool
	gameKeyDownEnd              map[pixelgl.Button]bool
	gameKeyWasPressed           map[pixelgl.Button]bool
	gameKeyJustPressed          map[pixelgl.Button]bool
	gameKeyTimeSinceLastPressed map[pixelgl.Button]int

	gameKeys = map[pixelgl.Button]string{
		pixelgl.Key0:              "0",
		pixelgl.Key1:              "1",
		pixelgl.Key2:              "2",
		pixelgl.Key3:              "3",
		pixelgl.Key4:              "4",
		pixelgl.Key5:              "5",
		pixelgl.Key6:              "6",
		pixelgl.Key7:              "7",
		pixelgl.Key8:              "8",
		pixelgl.Key9:              "9",
		pixelgl.KeyA:              "A",
		pixelgl.KeyB:              "B",
		pixelgl.KeyC:              "C",
		pixelgl.KeyD:              "D",
		pixelgl.KeyE:              "E",
		pixelgl.KeyF:              "F",
		pixelgl.KeyG:              "G",
		pixelgl.KeyH:              "H",
		pixelgl.KeyI:              "I",
		pixelgl.KeyJ:              "J",
		pixelgl.KeyK:              "K",
		pixelgl.KeyL:              "L",
		pixelgl.KeyM:              "M",
		pixelgl.KeyN:              "N",
		pixelgl.KeyO:              "O",
		pixelgl.KeyP:              "P",
		pixelgl.KeyQ:              "Q",
		pixelgl.KeyR:              "R",
		pixelgl.KeyS:              "S",
		pixelgl.KeyT:              "T",
		pixelgl.KeyU:              "U",
		pixelgl.KeyV:              "V",
		pixelgl.KeyW:              "W",
		pixelgl.KeyX:              "X",
		pixelgl.KeyY:              "Y",
		pixelgl.KeyZ:              "Z",
		pixelgl.KeyRight:          "RIGHT",
		pixelgl.KeyLeft:           "LEFT",
		pixelgl.KeyDown:           "DOWN",
		pixelgl.KeyUp:             "UP",
		pixelgl.KeySpace:          "SPACE",
		pixelgl.KeyEnter:          "ENTER",
		pixelgl.KeyEscape:         "ESCAPE",
		pixelgl.KeyBackspace:      "BACKSPACE",
		pixelgl.KeyDelete:         "DELETE",
		pixelgl.KeyInsert:         "INSERT",
		pixelgl.KeyPageUp:         "PAPGEUP",
		pixelgl.KeyPageDown:       "PAGEDOWN",
		pixelgl.KeyHome:           "HOME",
		pixelgl.KeyEnd:            "END",
		pixelgl.KeyLeftControl:    "LEFTCONTROL",
		pixelgl.KeyLeftShift:      "LEFTSHIFT",
		pixelgl.KeyLeftAlt:    	   "LEFTALT",
		pixelgl.KeyRightControl:   "RIGHTCONTROL",
		pixelgl.KeyRightShift:     "RIGHTSHIFT",
		pixelgl.KeyRightAlt:       "RIGHTALT",
	}
)