package types

import "image/color"

var (
	ColorLine        = color.RGBA{R: 150, G: 150, B: 150, A: 255}
	ColorLineDrawing = color.RGBA{R: 255, G: 200, B: 0, A: 255}

	// Game play scene bg
	ColorBackground = color.RGBA{17, 17, 27, 255}

	// Game victory scene bg
	ColorOverlay = color.RGBA{49, 50, 68, 255}

	ColorText       = color.RGBA{255, 255, 255, 255}
	ColorTextGreen  = color.RGBA{166, 227, 161, 255}
	ColorTextYellow = color.RGBA{249, 226, 175, 255}
	ColorTextRed    = color.RGBA{243, 139, 168, 255}
)
