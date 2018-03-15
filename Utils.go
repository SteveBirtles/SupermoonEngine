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
	luaRenderer.Draw(win, pixel.IM)

	win.Update()

	frames++
	select {
	case <-second:
		undoFrame++
		win.SetTitle(fmt.Sprintf("%s | FPS: %d", windowTitlePrefix, frames))
		frames = 0
	default:
	}

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