package main

func mainLoop() {

	initiateEngine()
	initiateAPI()
	load()
	backup()
	createEntities()

	for !win.Closed() && quit==0 {

		startFrame()
		processEditorInputs()
		updateEntities()
		renderOutputs()
		endFrame()
		
	}

	if quit >= 0 { save() }

}
