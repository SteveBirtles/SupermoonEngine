package main

import (
	"golang.org/x/image/colornames"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math"
	"github.com/faiface/pixel/imdraw"
	"image/color"
)


func mainLoop() {

	initiate()
	load()
	backup()

	tileOverlay := pixel.NewBatch(&pixel.TrianglesData{}, tilePic)
	imd := imdraw.New(nil)
	batch := pixel.NewBatch(&pixel.TrianglesData{}, tilePic)

	for !win.Closed() {

		if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyQ) {
			break
		}

		if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyN) {
			backup()
			grid = [2 * gridCentre][2 * gridCentre][16][2]uint16{}
		}

		if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyS) {
			save()
		}

		if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyL) {
			load()
		}

		backspacePressed := win.Pressed(pixelgl.KeyBackspace)

		if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyG) {
			showGrid = (showGrid + 1) % 3
		}

		if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyJ) {
			xPressed = !xPressed
		}

		if win.Pressed(pixelgl.KeyLeftControl) && win.JustPressed(pixelgl.KeyK) {
			zPressed = !zPressed
		}

		if win.JustPressed(pixelgl.KeyLeftShift) {
			selectionStartX = tileX
			selectionStartY = tileY
			selectionLive = true
		}
		if win.Pressed(pixelgl.KeyLeftShift) {
			selectionEndX = tileX
			selectionEndY = tileY
		}

		if win.JustPressed(pixelgl.KeyEscape) {
			xPressed = false
			zPressed = false
			selectionLive = false
		}

		previewClipboard = -1

		if win.JustPressed(pixelgl.Key0) { currentClipboard = 0 }
		if win.JustPressed(pixelgl.Key1) { currentClipboard = 1 }
		if win.JustPressed(pixelgl.Key2) { currentClipboard = 2 }
		if win.JustPressed(pixelgl.Key3) { currentClipboard = 3 }
		if win.JustPressed(pixelgl.Key4) { currentClipboard = 4 }
		if win.JustPressed(pixelgl.Key5) { currentClipboard = 5 }
		if win.JustPressed(pixelgl.Key6) { currentClipboard = 6 }
		if win.JustPressed(pixelgl.Key7) { currentClipboard = 7 }
		if win.JustPressed(pixelgl.Key8) { currentClipboard = 8 }
		if win.JustPressed(pixelgl.Key9) { currentClipboard = 9 }

		if win.Pressed(pixelgl.Key0) { previewClipboard = 0 }
		if win.Pressed(pixelgl.Key1) { previewClipboard = 1 }
		if win.Pressed(pixelgl.Key2) { previewClipboard = 2 }
		if win.Pressed(pixelgl.Key3) { previewClipboard = 3 }
		if win.Pressed(pixelgl.Key4) { previewClipboard = 4 }
		if win.Pressed(pixelgl.Key5) { previewClipboard = 5 }
		if win.Pressed(pixelgl.Key6) { previewClipboard = 6 }
		if win.Pressed(pixelgl.Key7) { previewClipboard = 7 }
		if win.Pressed(pixelgl.Key8) { previewClipboard = 8 }
		if win.Pressed(pixelgl.Key9) { previewClipboard = 9 }

		if win.JustPressed(pixelgl.KeyE) {
			clipboardWidth[currentClipboard] = -1
		}

		if win.JustPressed(pixelgl.KeyB) {
			clobber = !clobber
		}

		leftAltPressed := win.Pressed(pixelgl.KeyLeftAlt)
		rightAltPressed := win.Pressed(pixelgl.KeyRightAlt)

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

		if backspacePressed {

			selectedTile2 = 0

		} else if leftAltPressed {

			selectedTile1 = uint16(16*win.MousePosition().X/screenWidth) + 1
			if selectedTile1 < 1 {
				selectedTile1 = 1
			}
			if selectedTile1 > 16 {
				selectedTile1 = 16
			}

		} else if rightAltPressed {

			selectedTile2 = uint16(16*win.MousePosition().X/screenWidth) + 1
			if selectedTile2 < 1 {
				selectedTile2 = 1
			}
			if selectedTile2 > 16 {
				selectedTile2 = 16
			}

		} else {

			mouseX = float64(win.MousePosition().X - screenWidth/2)
			mouseY = float64(screenHeight/2 - win.MousePosition().Y)

			tileX = int(floor((mouseX - cameraX*scale) / hScale))
			tileY = int(floor((mouseY + cameraY*scale*aspect) / vScale))

			if tileX > gridCentre-1 { tileX = gridCentre-1 }
			if tileY > gridCentre-1 { tileY = gridCentre-1 }
			if tileX < -gridCentre { tileX = -gridCentre }
			if tileY < -gridCentre { tileY = -gridCentre }

			if win.MouseScroll().Y != 0 {
				lastScale := scale
				scale /= 1 - win.MouseScroll().Y/10
				if scale < 0.1 {
					scale = 0.1
				}
				if scale > 2.0 {
					scale = 2.0
				}
				hScale = 128 * scale
				vScale = 128 * aspect * scale

				deltaX := cameraX - (-float64(tileX)*128 - 64)
				deltaY := cameraY - (+float64(tileY)*128 + 64)

				cameraX -= deltaX * (1 - lastScale/scale)
				cameraY -= deltaY * (1 - lastScale/scale)

			}

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

			if win.JustPressed(pixelgl.KeyPageUp) {
				tileZ += 1
				if tileZ > 15 { tileZ = 15 }
			} else if win.JustPressed(pixelgl.KeyPageDown) {
				tileZ -= 1
				if tileZ < 0 { tileZ = 0 }
			} else if win.JustPressed(pixelgl.KeyHome) {
				tileZ = 15
			} else if win.JustPressed(pixelgl.KeyEnd) {
				tileZ = 0
			}

			if win.Pressed(pixelgl.KeyW) {
				cameraY -= 10 / scale
			}
			if win.Pressed(pixelgl.KeyS) {
				cameraY += 10 / scale
			}
			if win.Pressed(pixelgl.KeyD) {
				cameraX -= 10 / scale
			}
			if win.Pressed(pixelgl.KeyA) {
				cameraX += 10 / scale
			}

		}

		cpy := win.JustPressed(pixelgl.KeyC)
		cut := win.JustPressed(pixelgl.KeyX)
		clr := win.JustPressed(pixelgl.KeyDelete)
		fill := win.JustPressed(pixelgl.KeyInsert)
		bfill := win.JustPressed(pixelgl.KeyF)

		if selectionLive && win.Pressed(pixelgl.KeyLeftControl) && (cpy || cut || clr || fill || bfill) {

			startX := selectionStartX + gridCentre
			startY := selectionStartY + gridCentre
			endX := selectionEndX + gridCentre
			endY := selectionEndY + gridCentre

			if startX > endX {
				temp := startX
				startX = endX
				endX = temp
			}
			if startY > endY {
				temp := startY
				startY = endY
				endY = temp
			}

			if cpy || cut {
				clipboardWidth[currentClipboard] = endX - startX
				if clipboardWidth[currentClipboard] > clipboardSize { clipboardWidth[currentClipboard] = clipboardSize }
				clipboardHeight[currentClipboard] = endY - startY
				if clipboardHeight[currentClipboard] > clipboardSize { clipboardHeight[currentClipboard] = clipboardSize }
			}

			startZ := 0
			endZ := 15

			if fill || bfill {
				startZ = tileZ
				endZ = tileZ+1
			}

			for i := startX; i <= endX; i++ {
				for j := startY; j <= endY; j++ {

					if i < 2*gridCentre && j < 2 * gridCentre && i-startX < clipboardSize && j-startY < clipboardSize {

						for k := startZ; k < endZ; k++ {

							if bfill && (grid[i][j][k][0] != 0 || grid[i][j][k][1] != 0) {
								continue
							}

							if cpy || cut {
								clipboard[currentClipboard][i-startX][j-startY][k][0] = grid[i][j][k][0]
								clipboard[currentClipboard][i-startX][j-startY][k][1] = grid[i][j][k][1]
							}

							temp1 := uint16(0)
							temp2 := uint16(0)

							if fill || bfill {
								temp1 = selectedTile1
								temp2 = selectedTile2
							}

							if !cpy && (grid[i][j][k][0] != temp1 || grid[i][j][k][1] != temp2) {

								undoCounter = (undoCounter + 1) % maxUndo
								for i := 0; i < maxUndo; i++ {
									if undo[i][0] < 0 { undo[i][0] = 0 }
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

		}

		paste := win.JustPressed(pixelgl.KeyV) && win.Pressed(pixelgl.KeyLeftControl) ||
			previewClipboard != -1 && win.JustPressed(pixelgl.MouseButtonLeft)

		if paste && clipboardWidth[currentClipboard] >= 0 {

			for i := tileX; i <= tileX + clipboardWidth[currentClipboard]; i++ {
				for j := tileY; j <= tileY + clipboardHeight[currentClipboard]; j++ {
					if i < gridCentre && j < gridCentre {
						for k := 0; k < 16; k++ {

							if grid[i+gridCentre][j+gridCentre][k][0] != clipboard[currentClipboard][i-tileX][j-tileY][k][0] ||
								grid[i+gridCentre][j+gridCentre][k][1] != clipboard[currentClipboard][i-tileX][j-tileY][k][1] {

								if !clobber && clipboard[currentClipboard][i-tileX][j-tileY][k][0] == 0 && clipboard[currentClipboard][i-tileX][j-tileY][k][1] == 0 {
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

								grid[i+gridCentre][j+gridCentre][k][0] = clipboard[currentClipboard][i-tileX][j-tileY][k][0]
								grid[i+gridCentre][j+gridCentre][k][1] = clipboard[currentClipboard][i-tileX][j-tileY][k][1]
							}
						}
					}
				}
			}

		}

		onGrid := tileX >= -gridCentre && tileY >= -gridCentre && tileX < gridCentre && tileY < gridCentre

		leftDown := previewClipboard == -1 && win.Pressed(pixelgl.MouseButtonLeft)
		rightDown := previewClipboard == -1 && win.Pressed(pixelgl.MouseButtonRight)
		middleDown := win.JustPressed(pixelgl.MouseButtonMiddle)

		if onGrid {
			if middleDown {

				t := grid[tileX+gridCentre][tileY+gridCentre][tileZ][0]
				if t > 0 {
					selectedTile1 = t
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
										if undo[i][0] < 0 { undo[i][0] = 0 }
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
						if undo[i][0] < 0 { undo[i][0] = 0 }
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

		win.Clear(colornames.Black)

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

		frames++
		select {
		case <-second:
			undoFrame++
			win.SetTitle(fmt.Sprintf("%s | FPS: %d | X: %d | Y: %d | Aspect: %d%% | Camera: %d, %d", windowTitlePrefix, frames, tileX, tileY, int(100*(1-aspect)), int(cameraX), int(cameraY)))
			frames = 0
		default:
		}

	}

	save()

}
