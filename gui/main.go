package gui

import (
	"mineSync/globals"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var err error

func Run() error {
    globals.SetFileName("gui.go")
    globals.SetFuncName("Run")
    globals.Info("GUI Started")

    rl.InitWindow(int32(win.X), int32(win.Y), "MineSync")
    rl.SetTargetFPS(60)

    Init()

    for !rl.WindowShouldClose() {
        rl.BeginDrawing()
        rl.ClearBackground(colorPallet.window)

        switch scene {
        case "Main":
            if err = mainMenu(); err != nil {
                globals.Error("Error starting main menu")
                return err
            }
        case "Worlds":
            worldsMenu()
        }

        rl.EndDrawing()
    }

    return nil
}
