package main

import (
	"github.com/faiface/pixel"
	"image/color"
	"math"
	"fmt"
)

func renderEditorOutputs() {

	iStart := -floor(screenWidth/(2*hScale)) - 2
	jStart := -floor(screenHeight/(2*vScale)) - 2
	iEnd := floor(screenWidth/(2*hScale)) + 2
	jEnd := floor(screenHeight/(2*vScale)) + 20

	var cX, cY float64

	switch viewDirection {
	case 0:
		cX = cameraX
		cY = cameraY
	case 1:
		cX = -cameraY
		cY = cameraX
	case 2:
		cX = -cameraX
		cY = -cameraY
	case 3:
		cX = cameraY
		cY = -cameraX
	}

	iOffset := -floor(scale * cX / hScale)
	jOffset := floor(scale * aspect * cY / vScale)

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

	if viewDirection == 0 {
		viewTileX = tileX
		viewTileY = tileY
	}

	for i0 := iStart + iOffset; i0 <= iEnd+ iOffset; i0++ {
		for j0 := jStart + jOffset; j0 <= jEnd+ jOffset; j0++ {

			var i, j float64

			switch viewDirection {
			case 0:
				i = i0
				j = j0
			case 1:
				i = -j0
				j = i0
			case 2:
				i = -i0
				j = -j0
			case 3:
				i = j0
				j = -i0
			}

			for k := 0.0; k < 16; k++ {

				kC := k - float64(clipboardShift)

				var alpha uint8 = 255
				var beta uint8 = 128

				if zRay && int(k) > tileZ || xRay && int(j) > viewTileY {
					continue
				} else {
					if zRay && int(k) < tileZ || xRay && int(j) < viewTileY {
						alpha = 128
						beta = 64
					}
				}

				//batch.SetColorMask(color.RGBA{alpha, alpha, alpha, 255})

				deltaX := 0
				deltaY := 0
				preview := false

				if int(i) >= -gridCentre && int(j) >= -gridCentre && int(i) < gridCentre && int(j) < gridCentre {

					baseTile := grid[int(i)+gridCentre][int(j)+gridCentre][int(k)][0]

					if previewClipboard != -1 && !(kC < 0 || kC > 15) {
						deltaX = int(i) - viewTileX
						deltaY = int(j) - viewTileY
						if deltaX >= 0 && deltaY >= 0 && deltaX <= clipboardWidth[previewClipboard] && deltaY <= clipboardHeight[previewClipboard] {
							if clobber || clipboard[previewClipboard][deltaX][deltaY][int(kC)][0] != 0 || clipboard[previewClipboard][deltaX][deltaY][int(kC)][1] != 0 {
								preview = true
							}
						}
					}

					if preview {
						baseTile = clipboard[previewClipboard][deltaX][deltaY][int(kC)][0]
					}

					if baseTile > 0 && baseTile <= superTiles || (selectedTile1 > 0 && int(i) == viewTileX && int(j) == viewTileY && int(k) == tileZ && !hideTile) {

						s := 4*(1-aspect)
						cam := pixel.V(cX, cY)
						pos := pixel.V(screenWidth/2+float64(i0*hScale)+hScale/2, screenHeight/2+(-vScale/2-float64((j0-k*s)*vScale)))

						frontTile := grid[int(i)+gridCentre][int(j)+gridCentre][int(k)][1]

						if preview {
							frontTile = clipboard[previewClipboard][deltaX][deltaY][int(kC)][1]
						}

						if frontTile > 0 && frontTile <= superTiles || (selectedTile2 > 0 && int(i) == viewTileX && int(j) == viewTileY && int(k) == tileZ && !hideTile) {

							matrix := pixel.IM.Moved(cam).ScaledXY(pixel.ZV, pixel.V(scale, scale*aspect)).Moved(pos).
								Moved(pixel.V(0, vScale*(1-aspect)*4))

							if baseTile > 0 {
								batch.SetColorMask(color.RGBA{alpha, alpha, alpha, 255})
								tileTexture[baseTile-1].Draw(batch, matrix)
							}

							if int(i) == viewTileX && int(j) == viewTileY && int(k) == tileZ && !hideTile {
								batch.SetColorMask(color.RGBA{alpha, alpha, alpha, 128})
								tileTexture[selectedTile1-1].Draw(batch, matrix)
							}

							if aspect < 1 {

								matrix := pixel.IM.
									ScaledXY(pixel.ZV, pixel.V(1, s)).
									Moved(cam).
									ScaledXY(pixel.ZV, pixel.V(scale, scale*aspect)).
									Moved(pos).
									Moved(pixel.V(0, vScale/2 - 2*(aspect-0.5)*vScale))

								if frontTile > 0 {
									batch.SetColorMask(color.RGBA{beta, beta, beta, 255})
									tileTexture[frontTile-1].Draw(batch, matrix)
								}

								if selectedTile2 > 0 && int(i) == viewTileX && int(j) == viewTileY && int(k) == tileZ && !hideTile {
									batch.SetColorMask(color.RGBA{beta, beta, beta, 128})
									tileTexture[selectedTile2-1].Draw(batch, matrix)
								}
							}


						} else {

							matrix := pixel.IM.Moved(cam).ScaledXY(pixel.ZV, pixel.V(scale, scale*aspect)).Moved(pos)

							if baseTile > 0 {
								batch.SetColorMask(color.RGBA{alpha, alpha, alpha, 255})
								tileTexture[baseTile-1].Draw(batch, matrix)
							}

							if int(i) == viewTileX && int(j) == viewTileY && int(k) == tileZ && !hideTile {
								batch.SetColorMask(color.RGBA{alpha, alpha, alpha, 128})
								tileTexture[selectedTile1-1].Draw(batch, matrix)
							}

						}

					}

				}

				if k == 0 && (showGrid > 0 || selectionLive) {

					cam := pixel.V(cX, cY)
					pos := pixel.V(screenWidth/2+float64(i0*hScale)+hScale/2, screenHeight/2+(-vScale/2-float64(j0*vScale)))

					matrix := pixel.IM.Moved(cam).ScaledXY(pixel.ZV, pixel.V(scale, scale*aspect)).Moved(pos)
					imd1.SetMatrix(matrix)

					gridIntensity := math.Sqrt(scale / 2) * float64(showGrid) * 0.5

					if selectionLive && viewDirection == 0 &&
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

						if viewDirection != 0 {
							imd1.Color = pixel.RGB(0, gridIntensity, gridIntensity)
						} else if (int(i) == viewTileX || int(i) == viewTileX+1) && int(j) == viewTileY {
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
						imd1.Line(2.5 / scale)

						if viewDirection != 0 {
							imd1.Color = pixel.RGB(0, gridIntensity, gridIntensity)
						} else if int(i) == viewTileX && (int(j) == viewTileY || int(j) == viewTileY-1) {
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
						imd1.Line(2.5 / scale)

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

		n := 0
		
		for j := uint16(0); j < tileOverlayHeight; j++ {

			yes := leftAltPressed && int(j) == int(float64(selectedTile1-1) / float64(tileOverlayWidth)) ||
				rightAltPressed && int(j) == int(float64(selectedTile2-1) / float64(tileOverlayWidth))

			if !yes {
				n += int(tileOverlayWidth)
				continue
			}

			for i := uint16(0); i < tileOverlayWidth; i++ {

				if leftAltPressed && uint16(n) == selectedTile1-1 {
					tileOverlay.SetColorMask(color.RGBA{R: 255, G: 255, B: 255, A: 255})
				} else if rightAltPressed && uint16(n) == selectedTile2-1 {
					tileOverlay.SetColorMask(color.RGBA{R: 255, G: 255, B: 255, A: 255})
				} else {
					tileOverlay.SetColorMask(color.RGBA{R: 128, G: 128, B: 128, A: 255})
				}

				u := (float64(i) + 0.5) * 150 + (screenWidth - float64(tileOverlayWidth) * 75)
				v := 100.0

				matrix := pixel.IM.Moved(pixel.V(u, v)).ScaledXY(pixel.ZV, pixel.V(0.5, 0.5))
				tileTexture[n].Draw(tileOverlay, matrix)
				n++

			}
		}

		win.SetComposeMethod(pixel.ComposeOver)
		tileOverlay.Draw(win)
	}

	if !help {

		print ("Level: " + levelFile)
		print ("")
		if viewDirection == 0 {
			print(fmt.Sprintf("Cursor: %d, %d, %d", viewTileX, viewTileY, tileZ))
		} else {
			print(fmt.Sprintf("View direction: %d", viewDirection))
		}
		print(fmt.Sprintf("Aspect: %d%%", int(100*(1-aspect))))
		print(fmt.Sprintf("Camera: %d, %d", int(cX), int(cY)))
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
		print(fmt.Sprintf("Clipboard shift: %d", clipboardShift))

		if xRay {
			print("Vertical slice on")
		}
		if zRay {
			print("Horizontal slice on")
		}
		if hideTile {
			print("Cursor tile hidden")
		}
		print("")
		print("H for help...")

	} else {

		print("Left click : Draw tile")
		print("Right click : Clear tile")
		print("Middle click : Pick tile")
		print("PgUp / PgDn : Cursor up / down")
		print("Home / End : Cursor top / bottom")
		print("Tab : Toggle cursor tile visibility")
		print("Left / Right Alt : Choose base / front tile")
		print("Alt + Mouse Wheel : Cycle tile row")
		print("Backspace : Clear front tile")
		print("W / S / A / D : Move camera")
		print("Mouse wheel : Zoom camera")
		print("Ctrl+Q : Save and quit")
		print("Ctrl+Alt+Q : Quit without saving")
		print("Ctrl+Alt+N : Clear map")
		print("Ctrl+S / Ctrl+L : Save / Load")
		print("Ctrl+G : Change grid mode")
		print("Ctrl+J : Vertical slice")
		print("Ctrl+K : Horizontal slice")
		print("Shift : Make selection")
		print("Esc : Reset slice/selection")
		print("Number : Choose & preview clipboard")
		print("- / = : Clipboard vertical shift")
		print("Ctrl+E : Clear clipboard")
		print("Ctrl+C / Ctrl+X / Ctrl+V : Copy / Cut / Paste")
		print("Number+Click : Quick paste")
		print("Ctrl+B : Toggle pasting blanks")
		print("Ctrl+Del : Clear selection")
		print("Ctrl+Ins : Fill selection")
		print("Ctrl+F : Fill selection gaps")
		print("Ctrl+Z / Ctrl+Y : Undo / Redo")
		print("[ / ] : Change view angle")
		print("Arrows : Change view direction")

	}


}
