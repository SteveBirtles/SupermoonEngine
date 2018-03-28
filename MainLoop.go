package main

func mainLoop() {

	initiateEngine()
	initiateAPI()
	load()
	backup()

	for !win.Closed() && quit==0 {

		startFrame()
		processInputs()
		if !editing {
			updateEntities()
		}
		preRenderEntities()
		render()
		endFrame()
		
	}

	if quit >= 0 { save() }

}
