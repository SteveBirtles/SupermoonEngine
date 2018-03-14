package main

import (
	"github.com/faiface/pixel"
	"image/color"
	"math"
	"fmt"
)

func renderEditorOutputs() {

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

	imd1.Clear()
	imd2.Clear()
	batch.Clear()

	for i := -iRange + iOffset; i <= iRange+iOffset; i++ {
		for j := -jRange + jOffset; j <= jRange+jOffset; j++ {
			for k := 0.0; k < 16; k++ {

				var alpha uint8 = 255

				if zRay && int(k) > tileZ || xRay && int(j) > tileY {
					continue
				} else {
					if zRay && int(k) < tileZ || xRay && int(j) < tileY {
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

				if k == 0 && (showGrid > 0 || selectionLive) {

					cam := pixel.V(cameraX, cameraY)
					pos := pixel.V(screenWidth/2+float64(i*hScale)+hScale/2, screenHeight/2+(-vScale/2-float64(j*vScale)))

					matrix := pixel.IM.Moved(cam).ScaledXY(pixel.ZV, pixel.V(scale, scale*aspect)).Moved(pos)
					imd1.SetMatrix(matrix)

					gridIntensity := math.Sqrt(scale / 2) * float64(showGrid) * 0.5

					if selectionLive &&
						int(i) >= startX && int(j) >= startY &&
						int(i) <= endX && int(j) <= endY {
						imd2.SetMatrix(matrix)

						imd2.Color = pixel.RGBA{R: 0.0, G: 0.333, B: 0.333}
						imd2.Push(pixel.V(-64, -64))
						imd2.Push(pixel.V(-64, 64))
						imd2.Push(pixel.V(64, 64))
						imd2.Push(pixel.V(64, -64))
						imd2.Polygon(0)

					}

					if showGrid > 0 {

						if (int(i) == tileX || int(i) == tileX+1) && int(j) == tileY {
							imd1.Color = pixel.RGB(255, 255, 255)
						} else if int(i) == 0 || int(i) == 1 {
							imd1.Color = pixel.RGB(gridIntensity, gridIntensity, 0)
						} else if int(i)+gridCentre < 0 || int(i)+gridCentre > 2*gridCentre-1 ||
							int(j)+gridCentre <= 0 || int(j)+gridCentre > 2*gridCentre-1 {
							imd1.Color = pixel.RGB(gridIntensity, 0, 0)
						} else {
							imd1.Color = pixel.RGB(0, gridIntensity, 0)
						}

						imd1.Push(pixel.V(-64, -64))
						imd1.Push(pixel.V(-64, 64))
						imd1.Line(1.0 / scale)

						if int(i) == tileX && (int(j) == tileY || int(j) == tileY-1) {
							imd1.Color = pixel.RGB(255, 255, 255)
						} else if int(j) == 0 || int(j) == -1 {
							imd1.Color = pixel.RGB(gridIntensity, gridIntensity, 0)
						} else if int(i)+gridCentre < 0 || int(i)+gridCentre >= 2*gridCentre-1 ||
							int(j)+gridCentre < 0 || int(j)+gridCentre > 2*gridCentre-1 {
							imd1.Color = pixel.RGB(gridIntensity, 0, 0)
						} else {
							imd1.Color = pixel.RGB(0, gridIntensity, 0)
						}

						imd1.Push(pixel.V(-64, -64))
						imd1.Push(pixel.V(64, -64))
						imd1.Line(2.0 / scale)

					}

				}

			}
		}
	}


	if showGrid == 1 {
		win.SetComposeMethod(pixel.ComposeOver)
		imd1.Draw(win)
	}

	win.SetComposeMethod(pixel.ComposeOver)
	batch.Draw(win)

	if showGrid == 2 {
		win.SetComposeMethod(pixel.ComposeOver)
		imd1.Draw(win)
	}

	imd2.Draw(win)

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

	if !help {

		print ("Level: " + levelFile)
		print ("")
		print(fmt.Sprintf("Cursor: %d, %d, %d", tileX, tileY, tileZ))
		print(fmt.Sprintf("Aspect: %d%%", int(100*(1-aspect))))
		print(fmt.Sprintf("Camera: %d, %d", int(cameraX), int(cameraY)))
		print(fmt.Sprintf("Scale: %.2f", scale))
		print(fmt.Sprintf("Tiles: %d/%d", selectedTile1, selectedTile2))
		switch showGrid {
		case 0:
			print("Grid: Off")
		case 1:
			print("Grid: Behind")
		case 2:
			print("Grid: Front")
		}
		if clobber {
			print(fmt.Sprintf("Clipboard: %d", currentClipboard))
		} else {
			print(fmt.Sprintf("Clipboard: %d*", currentClipboard))
		}

		if xRay {
			print("Vertical slice on")
		}
		if zRay {
			print("Horizontal slice on")
		}
		print("")
		print("H for help...")
	} else {
		print("Left click : Draw tile")
		print("Right click : Clear tile")
		print("Middle click : Pick tile")
		print("PgUp / PgDn / Home / End : Vertical")
		print("Left Alt : Choose base tile")
		print("Right Alt : Choose front tile")
		print("Backspace : Clear front tile")
		print("W / S / A / D : Move camera")
		print("Mouse wheel : Zoom camera")
		print("Ctrl+Q : Save and quit")
		print("Ctrl+Alt+Q : Quit without saving")
		print("Ctrl+Alt+N : Clear map")
		print("Ctrl+S : Save")
		print("Ctrl+L : Load")
		print("Ctrl+G : Change grid mode")
		print("Ctrl+J : Vertical slice")
		print("Ctrl+K : Horizontal slice")
		print("Shift : Selection")
		print("Esc : Reset slice/selection")
		print("Number : Choose clipboard")
		print("Ctrl+E : Clear clipboard")
		print("Ctrl+C : Copy")
		print("Ctrl+X : Cut")
		print("Ctrl+V : Paste")
		print("Number+Click : Quick paste")
		print("Ctrl+B : Toggle pasting blanks")
		print("Ctrl+Del : Clear selection")
		print("Ctrl+Ins : Fill selection")
		print("Ctrl+F : Fill selection gaps")
		print("Ctrl+Z : Undo")
		print("Ctrl+Y : Redo")
		print("[ / ] : Change view angle")

	}


}

