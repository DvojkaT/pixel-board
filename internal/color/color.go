package color

import "math/rand"

type Color string

const (
	RED    Color = "#FF0000"
	BLUE   Color = "#2E27F5"
	YELLOW Color = "#F2F527"
	GREEN  Color = "#00FF04"
	PINK   Color = "#EA00FF"
	CYAN   Color = "#00FFFF"
	BLACK  Color = "#000000"
)

var Colors = []Color{RED, BLUE, YELLOW, GREEN, PINK, CYAN, BLACK}

func Random() Color {
	return Colors[rand.Intn(len(Colors))]
}
