package gui

import (
	"mineSync/globals"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	colorPallet Color

	// Buttons
	connectButton box
	worldsButton box

    // Texts
    addrFieldTitle text
    connectText text
    worldsText text

    // Field
    addrFieldBox box
    addrFieldText text

    // Worlds
    worldsList []worldItem
)

// OFF-SETS
var (
    o_worlds vec
    o_field vec
)

func Init() {
    ///// Main Menu

	// Color Pallet
	colorPallet = Color{
		window:      hex("0a0a0a"),
		button:      hex("ffffff"),
		text:        hex("ffffff"),
		hoverButton: hex("ffffff"),
		hoverText:   hex("000000"),
        field:       hex("ffffff"),
	}

	// Buttons
	connectButton = box{
		size: vec{
			X: 200,
			Y: 50,
		},
		pos: vec{
			X: win.X / 2,
			Y: win.Y / 2,
		},
		color: colorPallet.button,
	}

	worldsButton = box{
		size: vec{
			X: 200,
			Y: 50,
		},
		pos: vec{
			X: win.X / 2,
			Y: win.Y / 2,
		},
		color: colorPallet.button,
	}

	// Texts
	connectText = text{
		pos: vec{
			X: win.X / 2,
			Y: win.Y / 2,
		},
		color: colorPallet.text,

		font: rl.GetFontDefault(),
        text: "Host",
        size: 25,
	}

	worldsText = text{
		pos: vec{
			X: win.X / 2,
			Y: win.Y / 2,
		},
		color: colorPallet.text,

		font: rl.GetFontDefault(),
        text: "Worlds",
        size: 25,
	}

	addrFieldTitle = text{
		pos: vec{
			X: win.X / 2,
			Y: win.Y / 2,
		},
		color: colorPallet.text,

		font: rl.GetFontDefault(),
        text: "URL",
        size: 25,
	}

    // Addr Field
    addrFieldBox = box{
        size: vec{
            X: 60,
            Y: 40,
        },
        pos: vec{
            X: win.X / 2,
            Y: win.Y / 2,
        },
        color: colorPallet.field,
    }

    addrFieldText = text{
        pos: vec{
            X: win.X / 2,
            Y: win.Y / 2,
        },
        color: colorPallet.text,

        font: rl.GetFontDefault(),
        size: 25,
    }

    ///// Worlds Menu
    worldsList = make([]worldItem, len(globals.WorldsList))
    for i, entry := range globals.WorldsList {
        // Make default entry items
        worldsList[i].name = text{
            pos: vec{
                X: win.X / 2, 
                Y: 100,
            },
            color: colorPallet.text,

            font: rl.GetFontDefault(),
            size: 20,
        }

        worldsList[i].box = box{
            size: vec{
                X: 1.5 * rl.MeasureTextEx(worldsList[i].name.font, entry, worldsList[i].name.size, 0).X + 20,
                Y: rl.MeasureTextEx(worldsList[i].name.font, entry, worldsList[i].name.size, 0).Y + 20,
            },
            pos: vec{
                X: win.X / 2,
                Y: 100,
            },
            color: colorPallet.button,
        }

        worldsList[i].name.text = entry

        // Adjust each based on the other, the height
        worldsList[i].name.pos.Y += 2 * worldsList[i].box.size.Y * float32(i)
        worldsList[i].box.pos.Y += 2 * worldsList[i].box.size.Y * float32(i)

    }

    // OFF-SETS
    o_worlds = vec{
        X: 0,
        Y: 2 * worldsButton.size.Y,
    }

    o_field = vec{
        X: 0,
        Y: 2.25 * connectButton.size.Y,
    }

    // // set offsets
    worldsButton.pos.Y += o_worlds.Y
    worldsText.pos.Y += o_worlds.Y

    addrFieldBox.pos.Y -= o_field.Y
    addrFieldText.pos.Y -= o_field.Y

    addrFieldTitle.pos.Y = addrFieldBox.pos.Y - addrFieldBox.size.Y * (3/2)
}
