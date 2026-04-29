package engine

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type MainMenuScene struct{}

func (this *MainMenuScene) Update() (Scene, error) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		if x >= 300 && y <= 500 && y >= 250 && y <= 310 {
			return NewGameplayScene(), nil
		}

	}
	return this, nil
}

func (this *MainMenuScene) Draw(screen *ebiten.Image) {
	// Draw title
	ebitenutil.DebugPrintAt(screen, "=== SYSTEM DESIGN SANDBOX ===", 310, 150)

	// Draw Button Background
	vector.FillRect(screen, 300, 250, 200, 60, color.RGBA{100, 150, 255, 255}, true)

	// Draw Button Text
	ebitenutil.DebugPrintAt(screen, "START LEVEL 1", 355, 272)
}
