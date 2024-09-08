package main

import (
	"mineSync/core"
	"mineSync/globals"
	"mineSync/gui"
	"sync"
)

var err error

func run() error {
    // Set for errors
    globals.SetFileName("main.go")
    globals.SetFuncName("run")

    // Create a wait group
    wg := sync.WaitGroup{}

    // Make a err chan
    errChan := make(chan error, 1)

    // Create/Load config file
    wg.Add(1)
    go func () {
        if err = core.ConfigFile(); err != nil {
            // Set for errors
            globals.SetFileName("main.go")
            globals.SetFuncName("run")

            // Set error
            globals.Error("Failed to run core function")
            errChan <- err
        } else {
            errChan <- nil
        }
        wg.Done()
    }()

    // Check if any errors present in the err chan
    wg.Wait()
    select {
    case err = <- errChan:
        if err != nil {
            return err
        }
    default:
        // No error present, continue runnin
    }

    // Load worlds list
    if err = core.LoadWorldsList(); err != nil {
        globals.Error("Failed to run core function")
    }

    // Run the GUI
    if err = gui.Run(); err != nil {
        globals.Error("Failed to run gui")
        return err
    }


    return nil
}

func main() {

    // Run the run func && check for any errors
    if err = run(); err != nil {
        // Set for errors
        globals.SetFileName("main.go")
        globals.SetFuncName("main")

        globals.Error("Failed to run the app")
        globals.Print()
    }
}
