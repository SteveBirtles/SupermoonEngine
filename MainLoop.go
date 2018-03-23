package main

func mainLoop() {

	initiate()
	load()
	backup()

	for !win.Closed() && quit==0 {

		startFrame()

		processEditorInputs()

		updateEntities()

		renderEditorOutputs()

		endFrame()
	}

	if quit >= 0 { save() }

}
