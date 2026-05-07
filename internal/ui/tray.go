package ui

import (
	"fmt"
	"image/color"

	"github.com/blackviking27/system-design-game/internal/sim"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func DrawTray(screen *ebiten.Image, budget int) {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	trayHeight := 150
	trayY := h - trayHeight

	// Draw background tray
	vector.FillRect(screen, 0, float32(trayY), float32(w), float32(trayHeight), color.RGBA{40, 40, 40, 255}, true)

	// Draw budget
	DrawText(screen, fmt.Sprintf("BUDGET: $%d", budget), 20, float64(trayY+40), 16, color.White)

	// Draw catalog items
	startX := float32(150)
	for i, template := range sim.Catalog {
		x := startX + float32(i*180)
		y := float32(trayY + 20)

		// Draw icon (40x40)
		img := NodeImages[template.Type]
		if img != nil {
			op := &ebiten.DrawImageOptions{}
			// Scale image to 40x40
			bounds := img.Bounds()
			scaleX := 40.0 / float64(bounds.Dx())
			scaleY := 40.0 / float64(bounds.Dy())
			op.GeoM.Scale(scaleX, scaleY)
			op.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(img, op)
		} else {
			// Fallback if image not found
			vector.FillRect(screen, x, y, 40, 40, color.RGBA{100, 255, 150, 255}, true)
		}

		// Draw label and Cost below the icon
		label := fmt.Sprintf("%s\n$%d", template.Name, template.Cost)
		// Approximate center alignment: icon is at x, width 40. Text starts a bit to the left of center.
		DrawText(screen, label, float64(x-10), float64(y+50), 11, color.White)
	}

}
