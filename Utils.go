package main

import (
	"math"
	"image"
	"os"
	"fmt"
	"io/ioutil"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/colornames"
	"github.com/faiface/pixel"
	"time"
	"github.com/faiface/beep/wav"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
)

var (
	mixer = new(beep.Mixer)
	sampleRate = beep.SampleRate(48000)
	music beep.StreamSeekCloser
	musicPlaying string
)

func floor(x float64) float64 {
	return math.Floor(x)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadImageFile(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func startFrame() {

	win.Clear(colornames.Black)
	frameStart = time.Now()
	textRenderer.Clear()
	textLine = 0

}

func print(s string) {

	textRenderer.Dot = pixel.V(10, screenHeight-22 - float64(textLine)*22)
	textRenderer.WriteString(s + "\n")
	textLine++

}

func endFrame() {

	textRenderer.Draw(win, pixel.IM)

	if entityDebugData {
		renderLuaText()
	}

	win.Update()

	if music != nil {
		if music.Position() >= music.Len() {
			ASyncPlayMusic(musicPlaying)
		}
	}

	frameCounter++
	gameFrame++
	select {
	case <-second:
		undoFrame++
		executionLastSecond = cumulativeExecution
		cumulativeExecution = 0
		win.SetTitle(fmt.Sprintf("%s | FPS: %d", windowTitlePrefix, frameCounter))
		frameCounter = 0
	default:
	}

	frameLength = time.Now().Sub(frameStart).Seconds()

}

func loadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}


func calculateRenderBounds() {

	iStart = -floor(screenWidth/(2*hScale)) - 2
	jStart = -floor(screenHeight/(2*vScale)) - 2
	iEnd = floor(screenWidth/(2*hScale)) + 2
	jEnd = floor(screenHeight/(2*vScale)) + 20

	switch viewDirection {
	case 0:
		cameraAdjX = cameraX
		cameraAdjY = cameraY
	case 1:
		cameraAdjX = -cameraY
		cameraAdjY = cameraX
	case 2:
		cameraAdjX = -cameraX
		cameraAdjY = -cameraY
	case 3:
		cameraAdjX = cameraY
		cameraAdjY = -cameraX
	}

	iOffset = -floor(scale * cameraAdjX / hScale)
	jOffset = floor(scale * aspect * cameraAdjY / vScale)

	startX = selectionStartX
	startY = selectionStartY
	endX = selectionEndX
	endY = selectionEndY

	if startX > endX {
		temp := startX
		startX = endX
		endX = temp
	}

	if endX - startX > clipboardSize-1 {
		endX = startX + clipboardSize-1
	}

	if startY > endY {
		temp := startY
		startY = endY
		endY = temp
	}

	if endY - startY > clipboardSize-1 {
		endY = startY + clipboardSize-1
	}
}

func calculateViewVectors(i0 float64, j0 float64) (float64, float64) {

	switch viewDirection {
	case 0:
		return i0, j0
	case 1:
		return -j0, i0
	case 2:
		return -i0, -j0
	case 3:
		return j0, -i0
	}
	return 0, 0

}

func copyGrid(source *[2*gridCentre][2*gridCentre][16][2]uint16, destination *[2*gridCentre][2*gridCentre][16][2]uint16) {

	for x := 0; x < 2*gridCentre; x++ {
		for y := 0; y < 2*gridCentre; y++ {
			for z := 0; z < 16; z++ {
				(*destination)[x][y][z][0] = (*source)[x][y][z][0]
				(*destination)[x][y][z][1] = (*source)[x][y][z][1]
			}
		}
	}

}


func ASyncPlayMusic(musicFilename string) {

	musicPlaying = ""

	if music != nil {
		music.Close()
	}

	if musicFilename == "" {
		return
	}

	musicPlaying = musicFilename
	musicFilename = "resources/" + musicFilename

	if file, err := os.Open(musicFilename); err == nil {

		if m, musicFormat, musicErr := mp3.Decode(file); musicErr == nil {
			music = m

			//halfVolume := effects.Volume{Base: 10, Volume: -1000}

			if sampleRate != musicFormat.SampleRate {
				mixer.Play(beep.Resample(3, musicFormat.SampleRate, sampleRate, m)) //halfVolume.Stream(m)))

			} else {
				mixer.Play(m) //halfVolume.Stream(m))
			}

		} else {
			luaConsolePrint(fmt.Sprintf("Music error with %s: %s", musicFilename, musicErr))
		}


	} else {

		luaConsolePrint (fmt.Sprintf("Music file %s not found.", musicFilename))

	}

}

func AsyncPlaySound(soundFilename string) {

	soundFilename = "resources/" + soundFilename

	if file, err := os.Open(soundFilename); err == nil {
		if  s, soundFormat, soundError := wav.Decode(file); soundError == nil {
			if sampleRate != soundFormat.SampleRate {
				mixer.Play(beep.Resample(3, soundFormat.SampleRate, sampleRate, s))
			} else {
				mixer.Play(s)
			}
		} else {
			luaConsolePrint (fmt.Sprintf("Sound error with %s: %s", soundFilename, soundError))
		}
	} else {

		luaConsolePrint (fmt.Sprintf("Sound file %s not found.", soundFilename))

	}

}
