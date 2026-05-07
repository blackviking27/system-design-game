package engine

import (
	"image/color"

	"github.com/blackviking27/system-design-game/internal/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type MainMenuScene struct {
	screenWidth  int
	screenHeight int
}

func (this *MainMenuScene) Update() (Scene, error) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		sw, sh := this.screenWidth, this.screenHeight
		if sw == 0 || sh == 0 {
			sw, sh = 800, 600
		}
		bx, by, bw, bh := sw/2-100, sh/2-30, 200, 60

		currentLevel := "levels/01.json"
		if x >= bx && x <= bx+bw && y >= by && y <= by+bh {
			return NewGameplayScene(currentLevel), nil
		}

	}
	return this, nil
}

func (this *MainMenuScene) Draw(screen *ebiten.Image) {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	this.screenWidth = w
	this.screenHeight = h

	// Draw title
	ui.DrawText(screen, "=== SYSTEM DESIGN SANDBOX ===", float64(w/2-150), float64(h/2-150), 20, color.White)

	// Draw Button Background
	bx, by, bw, bh := w/2-100, h/2-30, 200, 60
	vector.FillRect(screen, float32(bx), float32(by), float32(bw), float32(bh), color.RGBA{100, 150, 255, 255}, true)

	// Draw Button Text
	ui.DrawText(screen, "START LEVEL 1", float64(w/2-55), float64(h/2-8), 16, color.White)
}
