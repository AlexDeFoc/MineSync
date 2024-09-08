package gui

import (
	"mineSync/globals"
	"strconv"
	"strings"

	color "image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func init() {
    rl.SetTraceLogLevel(rl.LogError)
}

func hov (b box) bool{
    tc := vec{
        X: b.pos.X - b.size.X / 2,
        Y: b.pos.Y - b.size.Y / 2,
    }

    bc := vec{
        X: b.pos.X + b.size.X / 2,
        Y: b.pos.Y + b.size.Y / 2,
    }

    m := rl.GetMousePosition()

    if m.X >= tc.X && m.X <= bc.X && m.Y >= tc.Y && m.Y <= bc.Y {
        return true
    }

    return false
}

func clk (b box) bool{
    if hov(b) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
        return true
    }

    return false
}

func scrl(b *box, t *text) {
    if rl.GetMouseWheelMoveV().Y != 0 {
        (*b).pos.Y -= rl.GetMouseWheelMoveV().Y * 20
        (*t).pos.Y -= rl.GetMouseWheelMoveV().Y * 20
        globals.Debug("Scrolling worlds list by this amount", rl.GetMouseWheelMoveV().Y)
    }
}

func typ(b *box, t *text) {
    // Center text constantly while typing
    t.pos.X = win.X / 2

    // Adjust field box size based on how many chars it contains
    b.size.X = 60 + rl.MeasureTextEx(t.font, "W", t.size, 0).X * float32(len(t.text))

    // Get char pressed
    key := rl.GetCharPressed()

    // Check if typing char from table (0-9, A-Z, a-z, '-', '.')
    if (key >= 48 && key <= 57) ||
       (key >= 65 && key <= 90) ||
       (key >= 97 && key <= 122)||
        key == 45 ||
        key == 46{
            t.text += string(key)
            globals.NetworkURL += string(key)
        }

    // Delete char
    if len(t.text) > 0 && (rl.IsKeyPressedRepeat(rl.KeyBackspace) || rl.IsKeyPressed(rl.KeyBackspace)) {
        t.text = t.text[:len(t.text)-1]
        globals.NetworkURL = globals.NetworkURL[:len(globals.NetworkURL)-1]
    }
}

/// Draw fun
func (b *box) Draw(fill ...int) {
    if fill == nil {
        r := rl.Rectangle{
            X: b.pos.X - b.size.X / 2,
            Y: b.pos.Y - b.size.Y / 2,
            Width: b.size.X,
            Height: b.size.Y,
        }

        thickness := float32(2)

        rl.DrawRectangleLinesEx(r, thickness, b.color)
    } else {
        s := rlVec(b.size)
        p := rl.Vector2{
            X: b.pos.X - b.size.X / 2,
            Y: b.pos.Y - b.size.Y / 2,
        }

        rl.DrawRectangleV(p, s, colorPallet.hoverButton)
    }
}

func (t *text) Draw(inverse ...int) {
    s := rl.MeasureTextEx(t.font, t.text, t.size, 1)
    p := rl.Vector2{
        X: t.pos.X - s.X / 2,
        Y: t.pos.Y - s.Y / 2,
    }

    c := t.color

    if inverse != nil {
        c = colorPallet.hoverText
    }

    rl.DrawTextEx(t.font, t.text, p, t.size, 1, c)
}

// UTILITY
func rlVec(v vec) rl.Vector2 {
    return rl.Vector2 {
        X: v.X,
        Y: v.Y,
    }
}

func hex(value string) color.RGBA{
	// Remove the hash at the start if present
	if strings.HasPrefix(value, "#") {
		value = value[1:]
	}

	// Check length
	if len(value) != 6 && len(value) != 8 {
		return color.RGBA{} // Return an empty color if the value code is invalid
	}

	// Parse RGB values
	var r, g, b, a uint32

	r = convertHex(value[0:2])
	g = convertHex(value[2:4])
	b = convertHex(value[4:6])

	// Parse Alpha if present
	if len(value) == 8 {
		a = convertHex(value[6:8])
	} else {
		a = 255 // Default alpha value if not specified
	}

	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
}

// Helper function to parse a 2-character value string to a uint32
func convertHex(value string) uint32 {
	val, _ := strconv.ParseUint(value, 16, 8)
	return uint32(val)
}
