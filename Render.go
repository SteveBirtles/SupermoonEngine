package main

import (
	"github.com/faiface/pixel"
	"image/color"
	"math"
	"fmt"
	"github.com/faiface/pixel/pixelgl"
)

var (
	startX int
	startY int
	endX int
	endY int
)

func renderEntities(i float64, j float64, i0 float64, j0 float64) {

	if int(i) >= -gridCentre && int(j) >= -gridCentre && int(i) < gridCentre && int(j) < gridCentre {

		if len(entityGrid[int(i)+gridCentre][int(j)+gridCentre]) > 0 {

			s := 4 * (1 - aspect)
			cam := pixel.V(cameraAdjX, cameraAdjY)

			var eX, eY float64

			for _, e := range entityGrid[int(i)+gridCentre][int(j)+gridCentre] {

				switch viewDirection {
				case 0:
					eX = e.X - i
					eY = e.Y - j
				case 1:
					eX = e.Y - j
					eY = -e.X + i
				case 2:
					eX = -e.X + i
					eY = -e.Y + j
				case 3:
					eX = -e.Y + j
					eY = e.X - i
				}

				k := e.Z
				spriteNo := -1

				d := reverseCompass[e.direction]

				moving := e.distance > 0 || e.nextDistance > 0 || e.progress != 0 || e.justOnTile

				if !editing {
					if (e.staticAnimation[d] || moving) && e.animationSpeed[d] > 0 {
						spriteNo = e.firstSprite[d] + int(math.Mod(float64(gameFrame) * e.animationSpeed[d] / 60.0, float64(e.lastSprite[d] - e.firstSprite[d] + 1)))
					} else {
						spriteNo = e.sprite[d]
					}
				} else {
					spriteNo = totalSprites-1
				}

				if spriteNo > -1 {

					pos := pixel.V(screenWidth/2+float64(i0*hScale)+hScale/2, screenHeight/2+(-vScale/2-float64((j0-k*s)*vScale)))

					spriteMatrix := pixel.IM.
						ScaledXY(pixel.ZV, pixel.V(1, s)).
						Moved(cam).
						ScaledXY(pixel.ZV, pixel.V(scale, scale*aspect)).
						Moved(pos).
						Moved(pixel.V(0, 2*(vScale/2-2*(aspect-0.5)*vScale))).
						Moved(pixel.V(eX*hScale, -eY*vScale))

					if editing {
						spriteBatch.SetColorMask(color.RGBA{128, 128, 128, 128})
					} else {
						spriteBatch.SetColorMask(color.RGBA{255, 255, 255, 255})
					}

					spriteTexture[spriteNo].Draw(spriteBatch, spriteMatrix)

				}

			}

		}
	}
}

func calculateShadows(i float64, j float64, k float64, frontTile uint16) bool {

	for s := 1.0; k+s < 16; s++ {

		if int(j-s) >= -gridCentre && int(j-s) < gridCentre {

			if frontTile == 0 && (grid[int(i)+gridCentre][int(j-s)+gridCentre][int(k+s-1)][1] > 0 ||
				grid[int(i)+gridCentre][int(j-s)+gridCentre][int(k+s)][0] > 0) ||
				frontTile > 0 && (grid[int(i)+gridCentre][int(j-s)+gridCentre][int(k+s)][0] > 0 ||
					grid[int(i)+gridCentre][int(j-s+1)+gridCentre][int(k+s)][0] > 0) {
				return  true
			}

		}
	}

	return false

}

func renderWorld(i float64, j float64, k float64, i0 float64, j0 float64) {

	var (
		alpha uint8 = 255
		beta uint8 = 128
		gamma uint8 = 192
	)

	if viewDirection == 2 {
		beta = 224
	}

	if zRay && int(k) > tileZ || xRay && int(j) > tileY {
		return
	} else {
		if zRay && int(k) < tileZ || xRay && int(j) < tileY {
			alpha = 128
			if viewDirection == 2 {
				beta = 112
			} else {
				beta = 64
			}
			gamma = 96
		}
	}

	deltaX := 0
	deltaY := 0
	kC := k - float64(clipboardShift)
	preview := false

	if int(i) >= -gridCentre && int(j) >= -gridCentre && int(i) < gridCentre && int(j) < gridCentre {

		baseTile := grid[int(i)+gridCentre][int(j)+gridCentre][int(k)][0]
		frontTile := grid[int(i)+gridCentre][int(j)+gridCentre][int(k)][1]

		inShadow := false

		if baseTile > 0 && showShadows {
			inShadow = calculateShadows(i, j, k, frontTile)
		}

		if previewClipboard != -1 && !(kC < 0 || kC > 15) {
			deltaX = int(i) - tileX
			deltaY = int(j) - tileY

			if flipX {
				deltaX = clipboardWidth[currentClipboard] - 1 - (int(i) - tileX)
			}
			if flipY {
				deltaY = clipboardHeight[currentClipboard] - 1 - (int(j) - tileY)
			}

			if deltaX >= 0 && deltaY >= 0 && deltaX < clipboardWidth[previewClipboard] && deltaY < clipboardHeight[previewClipboard] {
				if clobber || clipboard[previewClipboard][deltaX][deltaY][int(kC)][0] != 0 || clipboard[previewClipboard][deltaX][deltaY][int(kC)][1] != 0 {
					preview = true
				}
			}
		}

		if preview {
			baseTile = clipboard[previewClipboard][deltaX][deltaY][int(kC)][0]
		} else if inShadow {
			alpha = gamma
		}

		s := 4 * (1 - aspect)
		cam := pixel.V(cameraAdjX, cameraAdjY)

		pos := pixel.V( screenWidth/2+float64(i0*hScale)+hScale/2,
						screenHeight/2+(-vScale/2-float64((j0-k*s)*vScale)))

		if baseTile > 0 && baseTile <= totalTiles || (editing && selectedTile1 > 0 && int(i) == tileX && int(j) == tileY && int(k) == tileZ && !hideTile) {

			frontMatrix := pixel.IM.
				ScaledXY(pixel.ZV, pixel.V(1, s)).
				Moved(cam).
				ScaledXY(pixel.ZV, pixel.V(scale, scale*aspect)).
				Moved(pos).
				Moved(pixel.V(0, vScale/2-2*(aspect-0.5)*vScale))

			if preview {
				frontTile = clipboard[previewClipboard][deltaX][deltaY][int(kC)][1]
			}

			if frontTile > 0 && frontTile <= totalTiles || (editing && selectedTile2 > 0 && int(i) == tileX && int(j) == tileY && int(k) == tileZ && !hideTile) {

				baseMatrix := pixel.IM.
					Moved(cam).
					ScaledXY(pixel.ZV, pixel.V(scale, scale*aspect)).
					Moved(pos).
					Moved(pixel.V(0, vScale*(1-aspect)*4))

				if baseTile > 0 {
					tileBatch.SetColorMask(color.RGBA{alpha, alpha, alpha, 255})
					tileTexture[baseTile-1].Draw(tileBatch, baseMatrix)
				}

				if editing && int(i) == tileX && int(j) == tileY && int(k) == tileZ && !hideTile {
					tileBatch.SetColorMask(color.RGBA{alpha, alpha, alpha, 128})
					if selectedTile2 == 0 {
						tileTexture[selectedTile1-1].Draw(tileBatch, baseMatrix.Moved(pixel.V(0, -vScale*(1-aspect)*4)))
					} else {
						tileTexture[selectedTile1-1].Draw(tileBatch, baseMatrix)
					}

				}

				if aspect < 1 {

					if frontTile > 0 {
						tileBatch.SetColorMask(color.RGBA{beta, beta, beta, 255})
						tileTexture[frontTile-1].Draw(tileBatch, frontMatrix)
					}

					if editing && selectedTile2 > 0 && int(i) == tileX && int(j) == tileY && int(k) == tileZ && !hideTile {
						tileBatch.SetColorMask(color.RGBA{beta, beta, beta, 128})
						tileTexture[selectedTile2-1].Draw(tileBatch, frontMatrix)

					}
				}

			} else {

				matrix := pixel.IM.Moved(cam).ScaledXY(pixel.ZV, pixel.V(scale, scale*aspect)).Moved(pos)

				if baseTile > 0 {
					tileBatch.SetColorMask(color.RGBA{alpha, alpha, alpha, 255})
					tileTexture[baseTile-1].Draw(tileBatch, matrix)
				}

				if editing && int(i) == tileX && int(j) == tileY && int(k) == tileZ && !hideTile {
					tileBatch.SetColorMask(color.RGBA{alpha, alpha, alpha, 128})
					tileTexture[selectedTile1-1].Draw(tileBatch, matrix)
				}

			}

		}

	}

}

func renderGrid(i float64, j float64, k float64, i0 float64, j0 float64) {

	if k == 0 && (showGrid > 0 || selectionLive) {

		cam := pixel.V(cameraAdjX, cameraAdjY)
		pos := pixel.V(screenWidth/2+float64(i0*hScale)+hScale/2, screenHeight/2+(-vScale/2-float64(j0*vScale)))

		matrix := pixel.IM.Moved(cam).ScaledXY(pixel.ZV, pixel.V(scale, scale*aspect)).Moved(pos)
		imGrid.SetMatrix(matrix)

		gridIntensity := math.Sqrt(scale/2) * float64(showGrid) * 0.5

		if selectionLive &&
			int(i) >= startX && int(j) >= startY &&
			int(i) <= endX && int(j) <= endY {
			imUI.SetMatrix(matrix)

			imUI.Color = pixel.RGBA{R: 0.0, G: 0.333, B: 0.333}
			imUI.Push(pixel.V(-64, -64))
			imUI.Push(pixel.V(-64, 64))
			imUI.Push(pixel.V(64, 64))
			imUI.Push(pixel.V(64, -64))
			imUI.Polygon(0)
		}

		if int(i) == tileX && int(j) == tileY {

			if tileZ > 0 {
				imUI.SetMatrix(matrix)
				imUI.Color = pixel.RGBA{R: 0.2, G: 0.2, B: 0}
				imUI.Push(pixel.V(-64, -64))
				imUI.Push(pixel.V(-64, 64+(1-aspect)*(float64(tileZ)-0.5)*512.0))
				imUI.Push(pixel.V(64, 64+(1-aspect)*(float64(tileZ)-0.5)*512.0))
				imUI.Push(pixel.V(64, -64))
				imUI.Polygon(0)
			}
			if selectedTile2 > 0 {
				imUI.SetMatrix(matrix)
				imUI.Color = pixel.RGBA{R: 0.2, G: 0.2, B: 0.2}
				imUI.Push(pixel.V(-64, -64+(1-aspect)*float64(tileZ)*512.0))
				imUI.Push(pixel.V(-64, 64+(1-aspect)*float64(tileZ+1)*512.0))
				imUI.Push(pixel.V(64, 64+(1-aspect)*float64(tileZ+1)*512.0))
				imUI.Push(pixel.V(64, -64+(1-aspect)*float64(tileZ)*512.0))
				imUI.Polygon(0)
			} else {
				imUI.SetMatrix(matrix)
				imUI.Color = pixel.RGBA{R: 0.2, G: 0.2, B: 0.2}
				imUI.Push(pixel.V(-64, -64+(1-aspect)*float64(tileZ)*512.0))
				imUI.Push(pixel.V(-64, 64+(1-aspect)*float64(tileZ)*512.0))
				imUI.Push(pixel.V(64, 64+(1-aspect)*float64(tileZ)*512.0))
				imUI.Push(pixel.V(64, -64+(1-aspect)*float64(tileZ)*512.0))
				imUI.Polygon(0)

				if int(i) >= -gridCentre && int(j) >= -gridCentre && int(i) < gridCentre && int(j) < gridCentre {
					currentFront := grid[int(i)+gridCentre][int(j)+gridCentre][int(tileZ)][1]
					if currentFront != 0 {
						imUI.SetMatrix(matrix)
						imUI.Color = pixel.RGBA{R: 0.2, G: 0.0, B: 0.0}
						imUI.Push(pixel.V(-64, 64+(1-aspect)*float64(tileZ)*512.0))
						imUI.Push(pixel.V(-64, 64+(1-aspect)*float64(tileZ+1)*512.0))
						imUI.Push(pixel.V(64, 64+(1-aspect)*float64(tileZ+1)*512.0))
						imUI.Push(pixel.V(64, 64+(1-aspect)*float64(tileZ)*512.0))
						imUI.Polygon(0)
					}
				}
			}

		}

		if showGrid > 0 {

			if !(int(i) >= -gridCentre && int(j) >= -gridCentre && int(i) < gridCentre && int(j) < gridCentre) {

				imUI.SetMatrix(matrix)

				imUI.Color = pixel.RGBA{R: gridIntensity, G: 0.0, B: 0.0}
				imUI.Push(pixel.V(-64, -64))
				imUI.Push(pixel.V(-64, 64))
				imUI.Push(pixel.V(64, 64))
				imUI.Push(pixel.V(64, -64))
				imUI.Polygon(0)

			}

			imGrid.Color = pixel.RGB(0, gridIntensity, 0)
			imGrid.Push(pixel.V(-64, -64))
			imGrid.Push(pixel.V(-64, 64))
			imGrid.Line(2.5 / scale)

			imGrid.Color = pixel.RGB(0, gridIntensity, 0)
			imGrid.Push(pixel.V(-64, -64))
			imGrid.Push(pixel.V(64, -64))
			imGrid.Line(2.5 / scale)

		}

	}

}

func renderTileOverlay() {

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

}

func renderHelp() {

	if !help {

		print ("Level: " + levelFile)
		print ("")
		print(fmt.Sprintf("Cursor: X %d, Y %d, Z %d", tileX, tileY, tileZ))
		print(fmt.Sprintf("View direction: %s", compass[viewDirection]))
		print(fmt.Sprintf("Aspect: %d%%", int(100*(1-aspect))))
		print(fmt.Sprintf("Camera: %d, %d", int(cameraAdjX), int(cameraAdjY)))
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
		print(fmt.Sprintf("Clipboard: %d (%d, %d)", currentClipboard, clipboardWidth[currentClipboard], clipboardHeight[currentClipboard]))
		print(fmt.Sprintf("Clipboard shift: %d", clipboardShift))
		if flipX {
			print("Clipboard flipped in X direction")
		}
		if flipY {
			print("Clipboard flipped in Y direction")
		}
		if clobber {
			print("Blank tile pasting on")
		}
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

		var entitiesAtCursor []int

		for i := len(entities[0]) - 1; i >= 0; i-- {

			if int(entities[0][i].X) == tileX &&
				int(entities[0][i].Y) == tileY &&
				int(entities[0][i].Z) == tileZ {
				entitiesAtCursor = append(entitiesAtCursor, i)
			}

		}

		if len(entitiesAtCursor) == 0 {
			print("No entity at cursor")
		} else {
			print("Entity at cursor:")
			for _, i := range entitiesAtCursor {
				print(fmt.Sprintf("-> %s", entities[0][i].Class))
			}
		}


		print("")
		print(fmt.Sprintf("Scripts available (page %d of %d)", entityClassBlock + 1, entityClassBlockCount))

		for i := 1; i <= 12; i++ {
			j := i + 12 * entityClassBlock - 1
			if j < len(entityClass) {
				print(fmt.Sprintf("F%d : %s", i, entityClass[j]))
			}
		}

		print("")
		print("PgUp/Dn : Prev/Next page")
		print("Backspace : Delete entity")
		print("H : Additional help")
		print("")
		print("Tab : Toggle play mode")

	} else {

		print("Left click : Draw tile")
		print("Right click : Clear tile")
		print("Middle click : Pick tile")
		print("Mouse wheel : Cursor up/down")
		print("Home/End : Cursor high/low")
		print("Left/Right Alt : Choose base/front tile")
		print("Alt+Mouse Wheel : Cycle tile row")
		print("Right Ctrl : Clear front tile")
		print("Ctrl+O/P : Toggle shadow/cursor visibility")
		print("W/S/A/D : Move camera")
		print("Ctrl+Mouse wheel : Zoom camera")
		print("Ctrl+H : Centre camera")
		print("Ctrl+G : Change grid mode")
		print("Ctrl+J/K : Vertical/Horizontal slice")
		print("Shift : Make selection")
		print("Number : Choose & preview clipboard")
		print("-/= : Clipboard vertical shift")
		print("Ctrl+E : Clear clipboard")
		print("Ctrl+C/X/V : Copy/Cut/Paste")
		print("Number+Click : Quick paste")
		print("Ctrl+B : Toggle pasting blanks")
		print("Ctrl+F/Del : Fill/Del selection to floor")
		print("Ctrl+Ins : Fill selection gaps at Z level")
		print("Ctrl+R/T : Toggle mirror paste X/Y")
		print("Ctrl+Z/Y : Undo/Redo")
		print("[/]/Arrows : Change view angle/direction")
		print("Esc : Reset selection, view & toggles")
		print("Ctrl+S/Ctrl+L : Save/Load")
		print("Ctrl+Alt+N : Clear map")
		print("Ctrl+Q : Save & quit")
		print("Ctrl+Alt+Q : Quit without saving")

	}

}

func renderEntityStats() {

	print(fmt.Sprintf("Entities : %d (%d active)", entityCount, activeEntityCount))

	for i, e := range sortedEntityExecutionTimes {
		if i >= 5 { break }
		print(fmt.Sprintf("%d : \t%4.2fms - \t%s", e.id, float64(e.time)/1000000, e.class))
	}

	print (fmt.Sprintf("Total : %4.2fms (%4.2f%%)", float64(totalExecutionTime)/1000000,  100*float64(executionLastSecond)/1000000000))

}

func render() {

	preRenderEntities()

	calculateRenderBounds()

	imUI.Clear()

	for j0 := jStart + jOffset; j0 <= jEnd+ jOffset; j0++ {

		imGrid.Clear()
		tileBatch.Clear()
		spriteBatch.Clear()

		for i0 := iStart + iOffset; i0 <= iEnd+ iOffset; i0++ {

			i, j := calculateViewVectors(i0, j0)

			for k := 0.0; k < 16; k++ {
				renderWorld(i, j, k, i0, j0)
				if editing {
					renderGrid(i, j, k, i0, j0)
				}
			}
			renderEntities(i, j, i0, j0)

		}

		if editing && showGrid == 1 {
			win.SetComposeMethod(pixel.ComposeOver)
			imGrid.Draw(win)
		}

		win.SetComposeMethod(pixel.ComposeOver)

		tileBatch.Draw(win)
		spriteBatch.Draw(win)

	}

	if showGrid == 2 {
		win.SetComposeMethod(pixel.ComposeOver)
		imGrid.Draw(win)
	}

	imUI.Draw(win)
	if editing {
		renderTileOverlay()
		renderHelp()
	} else if win.Pressed(pixelgl.KeyGraveAccent) {
		renderEntityStats()
	}

}
