package gui

import "mineSync/globals"

func worldsMenu() {
	// Inverse colors on hover
	if hov(worldsButton) {
		worldsButton.Draw(1)
		worldsText.Draw(1)
	} else {
		worldsButton.Draw()
		worldsText.Draw()
	}

	// Change scenes when clicking worlds button
	if clk(worldsButton) {
		scene = "Main"
		worldsButton.pos.Y -= o_worlds.Y
		worldsText.pos.Y -= o_worlds.Y
		worldsText.text = "Worlds"

        globals.Info("Back button clicked")
	}

	// Render worlds list
	for i := range worldsList {
		// Inverse on click & hover
		if clk(worldsList[i].box) {
			worldsList[i].box.Draw()
			worldsList[i].name.Draw()

			// Change "SelectedWorld" on click IF not already the same
            if globals.SelectedWorld != worldsList[i].name.text{
                globals.SelectedWorld = worldsList[i].name.text
                globals.Info("Selected world with name:", worldsList[i].name.text)
            }

		} else if hov(worldsList[i].box) {
			worldsList[i].box.Draw(1)
			worldsList[i].name.Draw(1)
		} else {
			worldsList[i].box.Draw()
			worldsList[i].name.Draw()
		}

		// Move on scroll each world
		scrl(&worldsList[i].box, &worldsList[i].name)
	}
}
