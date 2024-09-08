package gui

import (
	color "image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type vec struct{
    X float32
    Y float32
}

type box struct{
    size vec
    pos vec
    color color.RGBA
}

type text struct{
    pos vec
    color color.RGBA

    font rl.Font
    text string
    size float32
}

type Color struct{
    window color.RGBA
    text color.RGBA
    button color.RGBA
    hoverButton color.RGBA
    hoverText color.RGBA
    field color.RGBA
}

type worldItem struct{
    name text
    box box
}
