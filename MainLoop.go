package main

func mainLoop() {

	initiate()
	load()
	backup()

	for !win.Closed() && !quit {
		processEditorInputs()
		renderEditorOuputs()
		endFrame()
	}

	save()

}
