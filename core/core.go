package core

import (
	"crypto/sha512"
	"mineSync/globals"
	"os"
	"path/filepath"
	"strings"
    "time"
    "encoding/json"
    "errors"
)

func Write(path string, data []byte) error {
	// Set for errors
	globals.SetFileName("core.go")
	globals.SetFuncName("Write")

    // Create path with worlds folder and selected world
    world_folder := globals.WorldsFolder
    clean(&world_folder)
    absolute_path := filepath.Join(world_folder, globals.SelectedWorld, path)

    // Create folder if needed
    err := os.MkdirAll(filepath.Dir(absolute_path), 0777)
    if err != nil {
        globals.Error("Error creating folder", err)
        return err
    }

    // Create file
    file, err := os.Create(absolute_path)
    if err != nil {
        globals.Error("Error creating file", err)
        return err
    }

    // Write data
    _, err = file.Write(data)
    if err != nil {
        globals.Error("Error writing to file", err)
        return err
    }

    // End
    file.Close()
    return nil
}

func LoadChanged(dataMap *map[string][]byte, changed []string) error {
	// Set for errors
	globals.SetFileName("core.go")
	globals.SetFuncName("LoadChanged")

    // for each CHANGED LOAD DATA
    for _, path := range changed {
        // Create path containing worlds folder and selected world
        world_folder := globals.WorldsFolder
        clean(&world_folder)
        absolute_path := filepath.Join(world_folder, globals.SelectedWorld, path)

        // Read file data
        file, err := os.ReadFile(absolute_path)
        if err != nil {
            globals.Error("Error reading file", err)
            return err
        }

        // Load into map
        (*dataMap)[path] = file
    }

    // End
    return nil
}

func CompareMetadata(changed, removed *[]string, client, server map[string][]byte) error {
	// Set for errors
	globals.SetFileName("core.go")
	globals.SetFuncName("CompareMetadata")

    // Changed/New
    for path, serverHash := range server {
        if clientHash, exists := client[path]; !exists || string(clientHash) != string(serverHash) {
            (*changed) = append((*changed), path)
        }
    }

    // Removed
    for path := range client {
        if _, exists := server[path]; !exists {
            (*removed) = append((*removed), path)
        }
    }

    // End
    return nil
}

func GetMetadata() (map[string][]byte, error) {
	// Set for errors
	globals.SetFileName("core.go")
	globals.SetFuncName("GetMetadata")

	// Create empty metadata
    metadata := make(map[string][]byte, 0)

    // Create path
    world_folder := globals.WorldsFolder
    clean(&world_folder)

    world_path := filepath.Join(world_folder, globals.SelectedWorld)

    // Walk & fill metadata
    err := filepath.WalkDir(world_path, func(path string, dir os.DirEntry, err error) error {
        // Check error from WalkDir
        if err != nil {
            globals.Error("Error from WalkDir", err)
            return err
        }

        // Add file to metadata
        if !dir.IsDir() {
            // Open & read file
            file, err := os.ReadFile(path)
            if err != nil {
                globals.Error("Error reading file", err)
                return err
            }

            // Create hash writer
            w := sha512.New()

            // Write to writer file contents
            _, err = w.Write(file)
            if err != nil {
                globals.Error("Error making hash", err)
                return err
            }

            // Remove parent folder
            cur_path := path
            clean(&cur_path)
            relative_path := strings.ReplaceAll(cur_path, world_path + `\`, ``)

            // Add file to metadata
            metadata[relative_path] = w.Sum(nil)
        }

        return nil
    })

    if err != nil {
        globals.Error("Error filling metadata", err)
        return metadata, err
    }

    // End
    return metadata, nil
}

func ConfigFile () error {
    // Set for errors
    globals.SetFileName("core.go")
    globals.SetFuncName("ConfigFile")

    // Default config structure
    WorldsFolder := struct {
        Path string `json:"Worlds Folder"`
        LogLevel string `json:"Log Level"`
    }{}

    // If file DOES NOT exists
    if !checkFileExists("config.json") {
        // Create file
        file, err := os.Create("config.json")
        if err != nil {
            globals.Error("Can't create config file", err)
            return err
        }

        // Apply default log level
        WorldsFolder.LogLevel = "Info"

        // Create default format
        contents, err := json.MarshalIndent(WorldsFolder, "", "  ")
        if err != nil {
            globals.Error("Can't make default config data", err)
            return err
        }

        // Write default format to file
        _, err = file.Write(contents)
        if err != nil {
            globals.Error("Can't write default formatted config to file", err)
            return err
        }
        file.Close()

        globals.Info("Created default config file")
        globals.Info("Please add in it your minecraft saved worlds folder path")
        time.Sleep(15 * time.Second)
        os.Exit(0)
    } else {
        // If exists clean worlds path, set worlds folder in globals

        // Open file & read from it
        previewFile, err := os.ReadFile("config.json")
        if err != nil {
            globals.Error("Can't read config contents", err)
            return err
        }

        // Clean previewed contents
        preview := string(previewFile[:])
        clean(&preview)
		if strings.Contains(preview, `\\`) || strings.Contains(preview, `\`) {
            globals.Info("Cleaned config path")
		}

        // New config file
        file, err := os.Create("config.json")
        if err != nil {
            globals.Error("Can't create a new config file", err)
            return err
        }

        // Write cleaned path to file
        _, err = file.Write([]byte(preview))
        if err != nil {
            globals.Error("Can't write cleaned path to config file", err)
        }
        file.Close()

        // Open file & read from it
        file, err = os.Open("config.json")
        if err != nil {
            globals.Error("Can't open config file", err)
            return err
        }

        // Decode file contents
        err = json.NewDecoder(file).Decode(&WorldsFolder)
        if err != nil {
            globals.Error("Error decoding config contents", err)
            return err
        }
        file.Close()

        // Check if path is absent
        if WorldsFolder.Path == "" {
            err = errors.New("No path found")
            globals.Error("Please provide a saved worlds folder path in the config file", err)
            return err
        }

        globals.WorldsFolder = WorldsFolder.Path
        globals.LogLevel = WorldsFolder.LogLevel

        globals.Info("Loaded config file")
    }

    return nil
}

func LoadWorldsList () error {
    // Set for errors
    globals.SetFileName("core.go")
    globals.SetFuncName("LoadWorldsList")

    // Config file structure
    WorldsFolder := struct {
        Path string `json:"Worlds Folder"`
        LogLevel string `json:"Log Level"`
    }{}

    // Open file & read from it
    file, err := os.Open("config.json")
    if err != nil {
        globals.Error("Can't open config file", err)
        return err
    }

    // Decode file contents
    err = json.NewDecoder(file).Decode(&WorldsFolder)
    if err != nil {
        globals.Error("Error decoding config contents", err)
        return err
    }
    file.Close()

    // Get worlds folder path
    path := WorldsFolder.Path

    // Get directories in the worlds folder path
    folders, err := os.ReadDir(path)
    if err != nil {
        globals.Error("Error getting folders from worlds folder path", err)
        return err
    }

    // Load worlds folders into globals list
    for i := range folders {
        globals.WorldsList = append(globals.WorldsList, folders[i].Name())
    }

    // End
    globals.Info("Loaded a list full of worlds")

    return nil
}

func Remove(files []string) error {
	// Set for errors
	globals.SetFileName("core.go")
	globals.SetFuncName("Remove")

	// Root folder path
	rootPath := filepath.Join(globals.WorldsFolder, globals.SelectedWorld)

	// Loop through file list
	for _, file := range files {
		// Set path
		path := filepath.Join(rootPath, file)

		// Remove file
		err := os.Remove(path)
		if err != nil {
			globals.Error("Failed removing file", err)
			return err
		}

		// Print which file got deleted
		globals.Info("Removed:", file)
	}

	// Remove empty directories
	err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && path != rootPath {
			entries, err := os.ReadDir(path)
			if err != nil {
				return err
			}

			if len(entries) == 0 {
				err = os.Remove(path)
				if err != nil {
					return err
				}
				globals.Info("Removed empty directory:", path)
			}
		}

		return nil
	})

	if err != nil {
		globals.Error("Failed removing empty directories", err)
		return err
	}

	return nil
}

func clean(path *string) {
    if strings.Contains(*path, `\\`) {
        *path = strings.ReplaceAll(*path, `\\`, `\`)
    }
}

func checkFileExists(path string) bool {
	_, error := os.Stat(path)
	return !errors.Is(error, os.ErrNotExist)
}
