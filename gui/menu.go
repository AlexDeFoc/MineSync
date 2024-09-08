package gui

import (
	"mineSync/client"
	"mineSync/globals"
	"mineSync/server"
)

func mainMenu() error {
    // Set values for logger
    globals.SetFileName("gui.go")
    globals.SetFuncName("mainMenu")

	// Inverse colors on hover
	if hov(connectButton) {
		connectButton.Draw(1)
		connectText.Draw(1)
	} else {
		connectButton.Draw()
		connectText.Draw()
	}

	if hov(worldsButton) {
		worldsButton.Draw(1)
		worldsText.Draw(1)
	} else {
		worldsButton.Draw()
		worldsText.Draw()
	}

    // Change scenes & move & change text of worlds button when clicking worlds button
    if clk(worldsButton) {
        scene = "Worlds"
        worldsButton.pos.Y += o_worlds.Y
        worldsText.pos.Y += o_worlds.Y
        worldsText.text = "Back"

        globals.Info("Worlds button clicked")
    }

    // Render addr field title
    addrFieldTitle.Draw()

    // Render the addr field
    addrFieldBox.Draw(1)
    addrFieldText.Draw(1)

    // DEBUG: Get url text
    if clk(addrFieldBox) {
        globals.Debug("URL:", addrFieldText.text)
    }

    // Type and draw new text
    typ(&addrFieldBox, &addrFieldText)

    // Change Host -> Connect if len of URL is more then zero
    if len(addrFieldText.text) > 0 {
        if connectText.text != "Connect" {
            connectText.text = "Connect"
        }
    } else {
        if connectText.text != "Host" {
            connectText.text = "Host"
        }
    }

    // Make an error chan for the server & client calls
    errChan := make(chan error, 1)

    // If host btn pressed call server | If connect btn pressed call client
    if clk(connectButton) && len(globals.SelectedWorld) != 0 {
        if connectText.text == "Host" {
            globals.Info("Host button clicked")
            go func () {
                if err = server.Run(); err != nil {
                    globals.Error("Error starting server")
                    errChan <- err
                } else {
                    errChan <- nil
                }
            }()
        } else if connectText.text == "Connect" {
            globals.Info("Connect button clicked")
            go func () {
                if err = client.Run(); err != nil {
                    globals.Error("Error connecting to server")
                    errChan <- err
                } else {
                    errChan <- nil
                }
            }()
        }
    }

    // Check if the err chan is empty | using select to not block the main thread if no err is present
    select {
    case err = <- errChan:
        if err != nil {
            return err
        }
    default:
        // No error, continue running
    }

    return nil
}
