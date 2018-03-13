package main

import (
	"github.com/faiface/pixel"
	"image/color"
	"math"
	"golang.org/x/image/colornames"
)

func renderEditorOutputs() {

	win.Clear(colornames.Black)

	iRange := floor(screenWidth/(2*hScale)) + 2
	jRange := floor(screenHeight/(2*vScale)) + 2

	iOffset := -floor(scale * cameraX / hScale)
	jOffset := floor(scale * aspect * cameraY / vScale)

	startX := selectionStartX
	startY := selectionStartY
	endX := selectionEndX
	endY := selectionEndY

	if startX > endX {
		temp := startX
		startX = endX
		endX = temp
	}

	if endX - startX > clipboardSize {
		endX = startX + clipboardSize
	}

	if startY > endY {
		temp := startY
		startY = endY
		endY = temp
	}

	if endY - startY > clipboardSize {
		endY = startY + clipboardSize
	}

	imd.Clear()
	batch.Clear()

	for i := -iRange + iOffset; i <= iRange+iOffset; i++ {
		for j := -jRange + jOffset; j <= jRange+jOffset; j++ {
			for k := 0.0; k < 16; k++ {

				var alpha uint8 = 255

				if zPressed && int(k) > tileZ || xPressed && int(j) > tileY {
					continue
				} else {
					if zPressed && int(k) < tileZ || xPressed && int(j) < tileY {
						alpha = 128
					}
				}

				batch.SetColorMask(color.RGBA{alpha, alpha, alpha, 255})

				deltaX := 0
				deltaY := 0
				preview := false

				if int(i) >= -gridCentre && int(j) >= -gridCentre && int(i) < gridCentre && int(j) < gridCentre {

					baseTile := grid[int(i)+gridCentre][int(j)+gridCentre][int(k)][0]

					if previewClipboard != -1 {
						deltaX = int(i) - tileX
						deltaY = int(j) - tileY
						if deltaX >= 0 && deltaY >= 0 && deltaX <= clipboardWidth[previewClipboard] && deltaY <= clipboardHeight[previewClipboard] {
							if clobber || clipboard[previewClipboard][deltaX][deltaY][int(k)][0] != 0 || clipboard[previewClipboard][deltaX][deltaY][int(k)][1] != 0 {
								preview = true
							}
						}
					}

					if preview {
						baseTile = clipboard[previewClipboard][deltaX][deltaY][int(k)][0]
					}

					if baseTile > 0 || (selectedTile1 > 0 && int(i) == tileX && int(j) == tileY && int(k) == tileZ) {

						s := 4*(1-aspect)
						cam := pixel.V(cameraX, cameraY)
						pos := pixel.V(screenWidth/2+float64(i*hScale)+hScale/2, screenHeight/2+(-vScale/2-float64((j-k*s)*vScale)))

						frontTile := grid[int(i)+gridCentre][int(j)+gridCentre][int(k)][1]

						if preview {
							frontTile = clipboard[previewClipboard][deltaX][deltaY][int(k)][1]
						}

						if frontTile > 0 || (selectedTile2 > 0 && int(i) == tileX && int(j) == tileY && int(k) == tileZ) {

							matrix := pixel.IM.Moved(cam).ScaledXY(pixel.ZV, pixel.V(scale, scale*aspect)).Moved(pos).
								Moved(pixel.V(0, vScale*(1-aspect)*4))

							if int(i) == tileX && int(j) == tileY && int(k) == tileZ {
								batch.SetColorMask(color.RGBA{alpha, alpha, alpha, 255})
								tileSprite[selectedTile1-1].Draw(batch, matrix)
							} else {
								tileSprite[baseTile-1].Draw(batch, matrix)
							}

							if aspect < 1 {

								matrix := pixel.IM.
									ScaledXY(pixel.ZV, pixel.V(1, s)).
									Moved(cam).
									ScaledXY(pixel.ZV, pixel.V(scale, scale*aspect)).
									Moved(pos).
									Moved(pixel.V(0, vScale/2 - 2*(aspect-0.5)*vScale))

								if selectedTile2 > 0 && int(i) == tileX && int(j) == tileY && int(k) == tileZ {
									batch.SetColorMask(color.RGBA{alpha, alpha, alpha, 255})
									tileSprite[selectedTile2-1].Draw(batch, matrix)
								} else {
									tileSprite[frontTile-1].Draw(batch, matrix)
								}
							}


						} else {

							matrix := pixel.IM.Moved(cam).ScaledXY(pixel.ZV, pixel.V(scale, scale*aspect)).Moved(pos)

							if int(i) == tileX && int(j) == tileY && int(k) == tileZ {
								batch.SetColorMask(color.RGBA{alpha, alpha, alpha, 255})
								tileSprite[selectedTile1-1].Draw(batch, matrix)
							} else {
								tileSprite[baseTile-1].Draw(batch, matrix)
							}

						}

					}

				}

				if k == 0 && showGrid > 0 {

					cam := pixel.V(cameraX, cameraY)
					pos := pixel.V(screenWidth/2+float64(i*hScale)+hScale/2, screenHeight/2+(-vScale/2-float64(j*vScale)))

					matrix := pixel.IM.Moved(cam).ScaledXY(pixel.ZV, pixel.V(scale, scale*aspect)).Moved(pos)
					imd.SetMatrix(matrix)

					gridIntensity := math.Sqrt(scale / 2) * float64(showGrid) * 0.5

					if selectionLive &&
						int(i) >= startX && int(j) >= startY &&
						int(i) <= endX && int(j) <= endY {

						imd.Color = pixel.RGBA{R: 0.0, G: 0.333, B: 0.333}
						imd.Push(pixel.V(-64, -64))
						imd.Push(pixel.V(-64, 64))
						imd.Push(pixel.V(64, 64))
						imd.Push(pixel.V(64, -64))
						imd.Polygon(0)

					} else {

						if (int(i) == tileX || int(i) == tileX+1) && int(j) == tileY {
							imd.Color = pixel.RGB(255, 255, 255)
						} else if int(i) == 0 || int(i) == 1 {
							imd.Color = pixel.RGB(gridIntensity, gridIntensity, 0)
						} else if int(i)+gridCentre < 0 || int(i)+gridCentre > 2*gridCentre-1 ||
							int(j)+gridCentre <= 0 || int(j)+gridCentre > 2*gridCentre-1 {
							imd.Color = pixel.RGB(gridIntensity, 0, 0)
						} else {
							imd.Color = pixel.RGB(0, gridIntensity, 0)
						}

						imd.Push(pixel.V(-64, -64))
						imd.Push(pixel.V(-64, 64))
						imd.Line(1.0 / scale)

						if int(i) == tileX && (int(j) == tileY || int(j) == tileY-1) {
							imd.Color = pixel.RGB(255, 255, 255)
						} else if int(j) == 0 || int(j) == -1 {
							imd.Color = pixel.RGB(gridIntensity, gridIntensity, 0)
						} else if int(i)+gridCentre < 0 || int(i)+gridCentre >= 2*gridCentre-1 ||
							int(j)+gridCentre < 0 || int(j)+gridCentre > 2*gridCentre-1 {
							imd.Color = pixel.RGB(gridIntensity, 0, 0)
						} else {
							imd.Color = pixel.RGB(0, gridIntensity, 0)
						}

						imd.Push(pixel.V(-64, -64))
						imd.Push(pixel.V(64, -64))
						imd.Line(2.0 / scale)

					}

				}

			}
		}
	}


	if showGrid == 1 {
		win.SetComposeMethod(pixel.ComposeOver)
		imd.Draw(win)
	}

	win.SetComposeMethod(pixel.ComposeOver)
	batch.Draw(win)

	if showGrid == 2 {
		win.SetComposeMethod(pixel.ComposeOver)
		imd.Draw(win)
	}

	if leftAltPressed || rightAltPressed {

		tileOverlay.Clear()

		for i := uint16(0); i < 16; i++ {

			if leftAltPressed && i == selectedTile1-1 {
				tileOverlay.SetColorMask(color.RGBA{R: 255, G: 255, B: 255, A: 255})
			} else if rightAltPressed && i == selectedTile2-1 {
				tileOverlay.SetColorMask(color.RGBA{R: 255, G: 255, B: 255, A: 255})
			} else {
				tileOverlay.SetColorMask(color.RGBA{R: 128, G: 128, B: 128, A: 255})
			}

			matrix := pixel.IM.Moved(pixel.V(float64(i)*150+150, 100)).ScaledXY(pixel.ZV, pixel.V(0.5, 0.5))
			tileSprite[i].Draw(tileOverlay, matrix)
		}

		win.SetComposeMethod(pixel.ComposeOver)
		tileOverlay.Draw(win)
	}

	win.Update()

}
