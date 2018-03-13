package main

import (
	"math"
	"image"
	"os"
	"fmt"
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

func endFrame() {

	frames++
	select {
	case <-second:
		undoFrame++
		win.SetTitle(fmt.Sprintf("%s | FPS: %d | X: %d | Y: %d | Aspect: %d%% | Camera: %d, %d",
			windowTitlePrefix, frames, tileX, tileY, int(100*(1-aspect)), int(cameraX), int(cameraY)))
		frames = 0
	default:
	}

}
