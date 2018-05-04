package main

import (
	"github.com/faiface/pixel/pixelgl"
	"math"
	"math/rand"
	"time"
)

func processSystemInputs() {

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyQ) {

		if !editing {
			copyGrid(&gridBackup, &grid)
		}

		if win.Pressed(pixelgl.KeyLeftAlt) || win.Pressed(pixelgl.KeyRightAlt) {
			quit = -1
		} else {
			quit = 1
		}
	}

	if win.JustPressed(pixelgl.KeyTab) {

		scale = 0.5
		hScale = 128 * scale
		vScale = 128 * aspect * scale

		if editing {
			copyGrid(&grid, &gridBackup)
		} else {
			copyGrid(&gridBackup, &grid)
		}

		editing = !editing
		resetViewState()

		if !editing {
			resetEntities()
		}

	}

}

func processEditorToggles() {

	if win.Pressed(pixelgl.KeyLeftControl) {

		if win.JustPressed(pixelgl.KeyG) {
			showGrid = (showGrid + 1) % 3
		}

		if win.JustPressed(pixelgl.KeyJ) {
			xRay = !xRay
		}

		if win.JustPressed(pixelgl.KeyK) {
			zRay = !zRay
		}

		if win.JustPressed(pixelgl.KeyP) {
			hideTile = !hideTile
		}

		if win.JustPressed(pixelgl.KeyO) {
			showShadows = !showShadows
		}

		if win.JustPressed(pixelgl.KeyR) {
			flipX = !flipX
		}

		if win.JustPressed(pixelgl.KeyT) {
			flipY = !flipY
		}

		if win.JustPressed(pixelgl.KeyE) {
			clipboardWidth[currentClipboard] = -1
		}

		if win.JustPressed(pixelgl.KeyB) {
			clobber = !clobber
		}

	}

	previewClipboard = -1

	if win.JustPressed(pixelgl.Key0) {
		currentClipboard = 0
	}
	if win.JustPressed(pixelgl.Key1) {
		currentClipboard = 1
	}
	if win.JustPressed(pixelgl.Key2) {
		currentClipboard = 2
	}
	if win.JustPressed(pixelgl.Key3) {
		currentClipboard = 3
	}
	if win.JustPressed(pixelgl.Key4) {
		currentClipboard = 4
	}
	if win.JustPressed(pixelgl.Key5) {
		currentClipboard = 5
	}
	if win.JustPressed(pixelgl.Key6) {
		currentClipboard = 6
	}
	if win.JustPressed(pixelgl.Key7) {
		currentClipboard = 7
	}
	if win.JustPressed(pixelgl.Key8) {
		currentClipboard = 8
	}
	if win.JustPressed(pixelgl.Key9) {
		currentClipboard = 9
	}

	if win.Pressed(pixelgl.Key0) {
		previewClipboard = 0
	}
	if win.Pressed(pixelgl.Key1) {
		previewClipboard = 1
	}
	if win.Pressed(pixelgl.Key2) {
		previewClipboard = 2
	}
	if win.Pressed(pixelgl.Key3) {
		previewClipboard = 3
	}
	if win.Pressed(pixelgl.Key4) {
		previewClipboard = 4
	}
	if win.Pressed(pixelgl.Key5) {
		previewClipboard = 5
	}
	if win.Pressed(pixelgl.Key6) {
		previewClipboard = 6
	}
	if win.Pressed(pixelgl.Key7) {
		previewClipboard = 7
	}
	if win.Pressed(pixelgl.Key8) {
		previewClipboard = 8
	}
	if win.Pressed(pixelgl.Key9) {
		previewClipboard = 9
	}

}

func processTileChoiceInputs() {

	if !win.Pressed(pixelgl.KeyLeftControl) {

		if win.JustPressed(pixelgl.KeyLeftAlt) {
			tileRow1 = uint16(selectedTile1 / tileOverlayWidth)
		}

		if win.JustPressed(pixelgl.KeyRightAlt) {
			tileRow2 = uint16(selectedTile2 / tileOverlayWidth)
		}

	}

	leftAltPressed = win.Pressed(pixelgl.KeyLeftAlt) && !win.Pressed(pixelgl.KeyLeftControl)
	rightAltPressed = win.Pressed(pixelgl.KeyRightAlt) && !win.Pressed(pixelgl.KeyLeftControl)

	if win.Pressed(pixelgl.KeyRightControl) {

		selectedTile2 = 0

	} else if leftAltPressed {

		if win.MouseScroll().Y > 0 {
			if tileRow1 == 0 {
				tileRow1 = tileOverlayHeight - 1
			} else {
				tileRow1 -= 1
			}
		} else if win.MouseScroll().Y < 0 {
			tileRow1 += 1
			if tileRow1 >= tileOverlayHeight {
				tileRow1 = 0
			}
		}

		selectedTile1 = uint16(float64(tileOverlayWidth)*float64(win.MousePosition().X/screenWidth)) +
			tileOverlayWidth*tileRow1 + 1

		if selectedTile1 < 1 {
			selectedTile1 = 1
		}
		if selectedTile1 > totalTiles {
			selectedTile1 = totalTiles
		}

	} else if rightAltPressed {

		if win.MouseScroll().Y > 0 {
			if tileRow2 == 0 {
				tileRow2 = tileOverlayHeight - 1
			} else {
				tileRow2 -= 1
			}
		} else if win.MouseScroll().Y < 0 {
			tileRow2 += 1
			if tileRow2 >= tileOverlayHeight {
				tileRow2 = 0
			}
		}

		selectedTile2 = uint16(float64(tileOverlayWidth)*float64(win.MousePosition().X/screenWidth)) +
			tileOverlayWidth*tileRow2 + 1

		if selectedTile2 < 1 {
			selectedTile2 = 1
		}
		if selectedTile2 > totalTiles {
			selectedTile2 = totalTiles
		}

	}

}

func processLevelFileInputs() {

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyR) {
		if win.Pressed(pixelgl.KeyLeftAlt) || win.Pressed(pixelgl.KeyRightAlt) {

			s := rand.NewSource(time.Now().UnixNano())
			r := rand.New(s)
			grid = [2 * gridCentre][2 * gridCentre][16][2]uint16{}
			for i := 0; i < 2*gridCentre; i++ {
				for j := 0; j < 2*gridCentre; j++ {
					h := r.Intn(15)
					h -= 8
					if h < 1 {
						h = 1
					}
					for k := 0; k < h; k++ {
						u := r.Intn(totalTiles) + 1
						grid[i][j][k][0] = uint16(u)
						if h > 1 {
							v := r.Intn(totalTiles) + 1
							grid[i][j][k][1] = uint16(v)
						}
					}
				}
			}

		}
	}

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyN) {
		if win.Pressed(pixelgl.KeyLeftAlt) || win.Pressed(pixelgl.KeyRightAlt) {
			backup()
			grid = [2 * gridCentre][2 * gridCentre][16][2]uint16{}
		}
	}

	if win.JustPressed(pixelgl.KeyH) {
		if win.Pressed(pixelgl.KeyLeftControl) {
			cameraX = 0
			cameraY = 0
		} else {
			help = !help
		}
	}

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyS) {
		save()
	}

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyL) {
		load()
	}

}

func processClipboardActions() {


	cpy := win.JustPressed(pixelgl.KeyC)
	cut := win.JustPressed(pixelgl.KeyX)
	clr := win.JustPressed(pixelgl.KeyDelete)
	flood := win.JustPressed(pixelgl.KeyInsert)
	fill := win.JustPressed(pixelgl.KeyF)
	paste := win.JustPressed(pixelgl.KeyV) && win.Pressed(pixelgl.KeyLeftControl) ||
		previewClipboard != -1 && win.JustPressed(pixelgl.MouseButtonLeft)

	if win.Pressed(pixelgl.KeyLeftControl) && (cpy || cut || clr || flood || fill) {

		startX := tileX + gridCentre
		startY := tileY + gridCentre
		endX := tileX + gridCentre
		endY := tileY + gridCentre

		startZ := 0
		endZ := 15

		if flood {
			startZ = tileZ
			endZ = tileZ
		} else if fill {
			startZ = 0
			endZ = tileZ
		}

		if selectionLive {

			startX = selectionStartX + gridCentre
			startY = selectionStartY + gridCentre
			endX = selectionEndX + gridCentre
			endY = selectionEndY + gridCentre

			if startX > endX {
				temp := startX
				startX = endX
				endX = temp
			}

			if endX-startX > clipboardSize {
				endX = startX + clipboardSize
			}

			if startY > endY {
				temp := startY
				startY = endY
				endY = temp
			}

			if endY-startY > clipboardSize {
				endY = startY + clipboardSize
			}

		}

		if cpy || cut {
			clipboardWidth[currentClipboard] = endX - startX + 1
			clipboardHeight[currentClipboard] = endY - startY + 1
		}

		for i := startX; i <= endX; i++ {
			for j := startY; j <= endY; j++ {

				if i < 2*gridCentre && j < 2*gridCentre && i-startX < clipboardSize && j-startY < clipboardSize {

					for k := startZ; k <= endZ; k++ {

						if flood && (grid[i][j][k][0] != 0 || grid[i][j][k][1] != 0) {
							continue
						}

						if cpy || cut {
							clipboard[currentClipboard][i-startX][j-startY][k][0] = grid[i][j][k][0]
							clipboard[currentClipboard][i-startX][j-startY][k][1] = grid[i][j][k][1]
						}

						temp1 := uint16(0)
						temp2 := uint16(0)

						if flood || fill {
							temp1 = selectedTile1
							temp2 = selectedTile2
						}

						if !cpy && (grid[i][j][k][0] != temp1 || grid[i][j][k][1] != temp2) {

							undoCounter = (undoCounter + 1) % maxUndo
							for i := 0; i < maxUndo; i++ {
								if undo[i][0] < 0 {
									undo[i][0] = 0
								}
							}
							undo[undoCounter][0] = undoFrame
							undo[undoCounter][1] = i
							undo[undoCounter][2] = j
							undo[undoCounter][3] = k
							undo[undoCounter][4] = int(grid[i][j][k][0])
							undo[undoCounter][5] = int(grid[i][j][k][1])

							grid[i][j][k][0] = temp1
							grid[i][j][k][1] = temp2
						}

					}
				}

			}
		}

		selectionLive = false

	} else if paste && clipboardWidth[currentClipboard] >= 0 {

		for i := tileX; i < tileX+clipboardWidth[currentClipboard]; i++ {
			for j := tileY; j < tileY+clipboardHeight[currentClipboard]; j++ {

				ii := i
				jj := j

				if flipX {
					ii = tileX + (clipboardWidth[currentClipboard] - 1 - (i - tileX))
				}
				if flipY {
					jj = tileY + (clipboardHeight[currentClipboard] - 1 - (j - tileY))
				}

				if i < gridCentre && j < gridCentre {
					for k0 := 0; k0 < 16; k0++ {

						k := k0 + clipboardShift
						if k < 0 || k > 15 {
							continue
						}

						if grid[i+gridCentre][j+gridCentre][k][0] != clipboard[currentClipboard][ii-tileX][jj-tileY][k0][0] ||
							grid[i+gridCentre][j+gridCentre][k][1] != clipboard[currentClipboard][ii-tileX][jj-tileY][k0][1] {

							if !clobber && clipboard[currentClipboard][ii-tileX][jj-tileY][k0][0] == 0 && clipboard[currentClipboard][ii-tileX][jj-tileY][k0][1] == 0 {
								continue
							}

							undoCounter = (undoCounter + 1) % maxUndo
							for i := 0; i < maxUndo; i++ {
								if undo[i][0] < 0 {
									undo[i][0] = 0
								}
							}
							undo[undoCounter][0] = undoFrame
							undo[undoCounter][1] = i + gridCentre
							undo[undoCounter][2] = j + gridCentre
							undo[undoCounter][3] = k
							undo[undoCounter][4] = int(grid[i+gridCentre][j+gridCentre][k][0])
							undo[undoCounter][5] = int(grid[i+gridCentre][j+gridCentre][k][1])

							grid[i+gridCentre][j+gridCentre][k][0] = clipboard[currentClipboard][ii-tileX][jj-tileY][k0][0]
							grid[i+gridCentre][j+gridCentre][k][1] = clipboard[currentClipboard][ii-tileX][jj-tileY][k0][1]
						}
					}
				}
			}
		}

	}

}

func processUndoInputs() {

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyZ) {
		m := 0
		for i := 0; i < maxUndo; i++ {
			if undo[i][0] > m {
				m = undo[i][0]
			}
		}

		if m != 0 {
			for i := 0; i < maxUndo; i++ {
				if undo[i][0] == m {
					temp1 := grid[undo[i][1]][undo[i][2]][undo[i][3]][0]
					temp2 := grid[undo[i][1]][undo[i][2]][undo[i][3]][1]
					grid[undo[i][1]][undo[i][2]][undo[i][3]][0] = uint16(undo[i][4])
					grid[undo[i][1]][undo[i][2]][undo[i][3]][1] = uint16(undo[i][5])
					undo[i][4] = int(temp1)
					undo[i][5] = int(temp2)
					undo[i][0] = -m
				}
			}
		}
	}

	if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyY) {
		n := -undoFrame - 1
		for i := 0; i < maxUndo; i++ {
			if undo[i][0] < 0 && undo[i][0] > n {
				n = undo[i][0]
			}
		}

		for i := 0; i < maxUndo; i++ {
			if undo[i][0] == n {
				temp1 := grid[undo[i][1]][undo[i][2]][undo[i][3]][0]
				temp2 := grid[undo[i][1]][undo[i][2]][undo[i][3]][1]
				grid[undo[i][1]][undo[i][2]][undo[i][3]][0] = uint16(undo[i][4])
				grid[undo[i][1]][undo[i][2]][undo[i][3]][1] = uint16(undo[i][5])
				undo[i][4] = int(temp1)
				undo[i][5] = int(temp2)
				undo[i][0] = -n
			}
		}

	}

}

func processMouseClicks() {

	onGrid := tileX >= -gridCentre && tileY >= -gridCentre && tileX < gridCentre && tileY < gridCentre

	leftDown := previewClipboard == -1 && win.Pressed(pixelgl.MouseButtonLeft)
	rightDown := previewClipboard == -1 && win.Pressed(pixelgl.MouseButtonRight)
	middleDown := win.JustPressed(pixelgl.MouseButtonMiddle)

	if onGrid {
		if middleDown {

			t := grid[tileX+gridCentre][tileY+gridCentre][tileZ][0]
			if t > 0 {
				selectedTile1 = t
				selectedTile2 = 0
			}
			t = grid[tileX+gridCentre][tileY+gridCentre][tileZ][1]
			if t > 0 {
				selectedTile2 = t
			}

		} else if leftDown || rightDown {

			newValue1 := selectedTile1
			newValue2 := selectedTile2
			if rightDown {
				newValue1 = 0
				newValue2 = 0
			}

			if lastTileX != outsideGrid {

				if math.Abs(float64(tileX-lastTileX)) > 1 || math.Abs(float64(tileY-lastTileY)) > 1 {

					d := 1.0 / float64(math.Abs(float64(lastTileX-tileX))+math.Abs(float64(lastTileY-tileY)))
					if d > 0 && d < 100 {
						dx := float64(lastTileX - tileX)
						dy := float64(lastTileY - tileY)
						for s := 0.0; s < 1.0; s += d {

							if grid[tileX+int(s*dx)+gridCentre][tileY+int(s*dy)+gridCentre][tileZ][0] != newValue1 ||
								grid[tileX+int(s*dx)+gridCentre][tileY+int(s*dy)+gridCentre][tileZ][1] != newValue2 {

								undoCounter = (undoCounter + 1) % maxUndo
								for i := 0; i < maxUndo; i++ {
									if undo[i][0] < 0 {
										undo[i][0] = 0
									}
								}
								undo[undoCounter][0] = undoFrame
								undo[undoCounter][1] = tileX + int(s*dx) + gridCentre
								undo[undoCounter][2] = tileY + int(s*dy) + gridCentre
								undo[undoCounter][3] = tileZ
								undo[undoCounter][4] = int(grid[tileX+int(s*dx)+gridCentre][tileY+int(s*dy)+gridCentre][tileZ][0])
								undo[undoCounter][5] = int(grid[tileX+int(s*dx)+gridCentre][tileY+int(s*dy)+gridCentre][tileZ][1])

								grid[tileX+int(s*dx)+gridCentre][tileY+int(s*dy)+gridCentre][tileZ][0] = newValue1
								grid[tileX+int(s*dx)+gridCentre][tileY+int(s*dy)+gridCentre][tileZ][1] = newValue2
							}
						}
					}
				}
			}

			if grid[tileX+gridCentre][tileY+gridCentre][tileZ][0] != newValue1 ||
				grid[tileX+gridCentre][tileY+gridCentre][tileZ][1] != newValue2 {

				undoCounter = (undoCounter + 1) % maxUndo
				for i := 0; i < maxUndo; i++ {
					if undo[i][0] < 0 {
						undo[i][0] = 0
					}
				}
				undo[undoCounter][0] = undoFrame
				undo[undoCounter][1] = tileX + gridCentre
				undo[undoCounter][2] = tileY + gridCentre
				undo[undoCounter][3] = tileZ
				undo[undoCounter][4] = int(grid[tileX+gridCentre][tileY+gridCentre][tileZ][0])
				undo[undoCounter][5] = int(grid[tileX+gridCentre][tileY+gridCentre][tileZ][1])

				grid[tileX+gridCentre][tileY+gridCentre][tileZ][0] = newValue1
				grid[tileX+gridCentre][tileY+gridCentre][tileZ][1] = newValue2

			}

			lastTileX = tileX
			lastTileY = tileY

		} else {

			lastTileX = outsideGrid

		}
	}


}

func processPositionInputs() {

	if win.Pressed(pixelgl.KeyLeftControl) && win.MouseScroll().Y != 0 {
		lastScale := scale
		scale /= 1 - win.MouseScroll().Y/10

		if scale < 0.1 {
			scale = 0.1
		}
		if scale > 2.0 {
			scale = 2.0
		}

		deltaX := cameraX - (-float64(tileX)*128 - 64)
		deltaY := cameraY - (+float64(tileY)*128 + 64)

		cameraX -= deltaX * (1 - lastScale/scale)
		cameraY -= deltaY * (1 - lastScale/scale)

	}

	hScale = 128 * scale
	vScale = 128 * aspect * scale

	if !(leftAltPressed || rightAltPressed) {

		mouseX = float64(win.MousePosition().X - screenWidth/2)
		mouseY = float64(screenHeight/2 - win.MousePosition().Y)

		switch  viewDirection {
		case 0:
			tileX = int(floor((mouseX - cameraX*scale) / hScale))
			tileY = int(floor((mouseY + cameraY*scale*aspect) / vScale))
		case 1:
			tileX = -int(floor((mouseY + cameraX*scale*aspect) / vScale))
			tileY = int(floor((mouseX + cameraY*scale) / hScale))
		case 2:
			tileX = -int(floor((mouseX + cameraX*scale) / hScale))
			tileY = -int(floor((mouseY - cameraY*scale*aspect) / vScale))
		case 3:
			tileX = int(floor((mouseY - cameraX*scale*aspect) / vScale))
			tileY = -int(floor((mouseX - cameraY*scale) / hScale))
		}

		if tileX > gridCentre-1 {
			tileX = gridCentre - 1
		}
		if tileY > gridCentre-1 {
			tileY = gridCentre - 1
		}
		if tileX < -gridCentre {
			tileX = -gridCentre
		}
		if tileY < -gridCentre {
			tileY = -gridCentre
		}

		if viewDirection == 0 && win.Pressed(pixelgl.KeyW) ||
			viewDirection == 1 && win.Pressed(pixelgl.KeyA) ||
			viewDirection == 2 && win.Pressed(pixelgl.KeyS) ||
			viewDirection == 3 && win.Pressed(pixelgl.KeyD) {
			cameraY -= 10 / scale
		}
		if viewDirection == 0 && win.Pressed(pixelgl.KeyS) ||
			viewDirection == 1 && win.Pressed(pixelgl.KeyD) ||
			viewDirection == 2 && win.Pressed(pixelgl.KeyW) ||
			viewDirection == 3 && win.Pressed(pixelgl.KeyA) {
			cameraY += 10 / scale
		}
		if viewDirection == 0 && win.Pressed(pixelgl.KeyD) ||
			viewDirection == 1 && win.Pressed(pixelgl.KeyW) ||
			viewDirection == 2 && win.Pressed(pixelgl.KeyA) ||
			viewDirection == 3 && win.Pressed(pixelgl.KeyS) {
			cameraX -= 10 / scale
		}
		if viewDirection == 0 && win.Pressed(pixelgl.KeyA) ||
			viewDirection == 1 && win.Pressed(pixelgl.KeyS) ||
			viewDirection == 2 && win.Pressed(pixelgl.KeyD) ||
			viewDirection == 3 && win.Pressed(pixelgl.KeyW) {
			cameraX += 10 / scale
		}
	}

}

func processEditorDirectives() {

	if !win.Pressed(pixelgl.KeyLeftControl) {

		if win.Pressed(pixelgl.KeyRightBracket) {
			aspect += 0.01
			if aspect > 1.0 {
				aspect = 1.0
			}
			vScale = hScale * aspect
		}
		if win.Pressed(pixelgl.KeyLeftBracket) {
			aspect -= 0.01
			if aspect < 0.5 {
				aspect = 0.5
			}
			vScale = hScale * aspect
		}

		if !(win.Pressed(pixelgl.KeyLeftAlt) || win.Pressed(pixelgl.KeyRightAlt)) {
			if win.MouseScroll().Y > 0 {
				tileZ += 1
				if tileZ > 15 {
					tileZ = 15
				}
			} else if win.MouseScroll().Y < 0 {
				tileZ -= 1
				if tileZ < 0 {
					tileZ = 0
				}
			} else if win.JustPressed(pixelgl.KeyEnd) {

				tileZ = 0
				for k := 0; k <= 15; k++ {
					if int(grid[tileX+gridCentre][tileY+gridCentre][k][0]) > 0 {
						tileZ = k
						break
					}
				}

			} else if win.JustPressed(pixelgl.KeyHome) {

				tileZ = 0
				for k := 15; k >= 0; k-- {
					if int(grid[tileX+gridCentre][tileY+gridCentre][k][0]) > 0 {
						tileZ = k
						break
					}
				}

			}
		}

	}

	if win.JustPressed(pixelgl.KeyUp) {
		viewDirection = 0
	}
	if win.JustPressed(pixelgl.KeyRight) {
		viewDirection = 1
	}
	if win.JustPressed(pixelgl.KeyDown) {
		viewDirection = 2
	}
	if win.JustPressed(pixelgl.KeyLeft) {
		viewDirection = 3
	}

	if win.JustPressed(pixelgl.KeyEqual) {
		clipboardShift += 1
		if clipboardShift > 15 {
			clipboardShift = 15
		}
	} else if win.JustPressed(pixelgl.KeyMinus) {
		clipboardShift -= 1
		if clipboardShift < -15 {
			clipboardShift = -15
		}
	}

	if win.JustPressed(pixelgl.KeyLeftShift) {
		selectionStartX = tileX
		selectionStartY = tileY
		selectionLive = true
	}
	if win.Pressed(pixelgl.KeyLeftShift) {
		selectionEndX = tileX
		selectionEndY = tileY
	} else {
		if selectionLive && selectionStartX == selectionEndX && selectionStartY == selectionEndY {
			selectionLive = false
		}
	}

	if win.JustPressed(pixelgl.KeyEscape) {
		resetViewState()
	}

	if win.JustPressed(pixelgl.KeyBackspace) {

		var deleteList []int

		for i := len(entities[0]) - 1; i >= 0; i-- {

			if int(entities[0][i].X) == tileX &&
				int(entities[0][i].Y) == tileY &&
				int(entities[0][i].Z) == tileZ {
					deleteList = append(deleteList, i)
			}

		}

		for _, i := range deleteList {

			entities[0] = append(entities[0][0:i], entities[0][i+1:]...)

		}

	}

	if win.JustPressed(pixelgl.KeyPageUp) {
		entityClassBlock--
		if entityClassBlock < 0 { entityClassBlock = entityClassBlockCount - 1 }
	}
	if win.JustPressed(pixelgl.KeyPageDown) {
		entityClassBlock++
		if entityClassBlock >= entityClassBlockCount { entityClassBlock = 0 }
	}

	fKeyPressed := 0

	if win.JustPressed(pixelgl.KeyF1)  { fKeyPressed = 1  }
	if win.JustPressed(pixelgl.KeyF2)  { fKeyPressed = 2  }
	if win.JustPressed(pixelgl.KeyF3)  { fKeyPressed = 3  }
	if win.JustPressed(pixelgl.KeyF4)  { fKeyPressed = 4  }

	if win.JustPressed(pixelgl.KeyF5)  { fKeyPressed = 5  }
	if win.JustPressed(pixelgl.KeyF6)  { fKeyPressed = 6  }
	if win.JustPressed(pixelgl.KeyF7)  { fKeyPressed = 7  }
	if win.JustPressed(pixelgl.KeyF8)  { fKeyPressed = 8  }

	if win.JustPressed(pixelgl.KeyF9)  { fKeyPressed = 9  }
	if win.JustPressed(pixelgl.KeyF10) { fKeyPressed = 10 }
	if win.JustPressed(pixelgl.KeyF11) { fKeyPressed = 11 }
	if win.JustPressed(pixelgl.KeyF12) { fKeyPressed = 12 }

	s := fKeyPressed + 12 * entityClassBlock - 1

	if fKeyPressed > 0 && s < lastEntityClass {

		eX := float64(tileX)
		eY := float64(tileY)
		eZ := float64(tileZ)

		gX := int(eX + gridCentre)
		gY := int(eY + gridCentre)
		gZ := int(eZ)

		valid := true

		if gX < 0 || gY < 0 || gX >= 2*gridCentre || gY >= 2*gridCentre || grid[gX][gY][gZ][1] != 0 {
			valid = false
		}

		for i := len(entities[0]) - 1; i >= 0; i-- {

			if int(entities[0][i].X) == tileX &&
				int(entities[0][i].Y) == tileY &&
				int(entities[0][i].Z) == tileZ {
					valid = false
			}

		}

		if valid {

			entityUID++

			e := Entity{Id: entityUID,
				active: true,
				firstSprite: -1,
				velocity: 0,
				direction: 'S',
				distance: 0,
				Class: entityClass[s],
				X: eX,
				Y: eY,
				Z: eZ,
				lastX: eX,
				lastY: eY,
				lastZ: eZ,
				targetX: eX,
				targetY: eY,
				targetZ: eZ,
				progress: 0}

			entities[0] = append(entities[0], e)

		}

	}


}

func resetViewState() {

	xRay = false
	zRay = false
	selectionLive = false
	help = false
	viewDirection = 0
	tileZ = 0
	flipX = false
	flipY = false
	showShadows = true
	hideTile = false
	clobber = false
	aspect = 0.5

}

func processGameKeys() {

	for k := range gameKeys {
		pressed := win.Pressed(k)
		if pressed {
			if gameKeyTimeSinceLastPressed[k] > 2 {
				gameKeyJustPressed[k] = true
				gameKeyTimeSinceLastPressed[k] = 0
			}
		} else {
			gameKeyTimeSinceLastPressed[k]++
		}
	}

}

func processInputs() {

	processSystemInputs()

	if editing {
		processLevelFileInputs()
		processEditorToggles()
		processUndoInputs()
		processEditorDirectives()
		processPositionInputs()
		processTileChoiceInputs()
		processClipboardActions()
		processMouseClicks()
	} else {
		processGameKeys()
	}

}
