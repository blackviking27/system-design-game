package ui

import (
	"fmt"
	"image/color"

	"github.com/blackviking27/system-design-game/internal/sim"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func DrawTray(screen *ebiten.Image, budget int) {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	trayHeight := 100
	trayY := h - trayHeight

	// Draw background tray
	vector.FillRect(screen, 0, float32(trayY), float32(w), float32(trayHeight), color.RGBA{40, 40, 40, 255}, true)

	// Draw budget
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("BUDGET: $%d", budget), 20, trayY+40)

	// Draw catalog items
	startX := float32(150)
	for i, template := range sim.Catalog {
		x := startX + float32(i*180)
		y := float32(trayY + 20)

		// Color based on type
		c := color.RGBA{100, 255, 150, 255} // Default Green
		switch template.Type {
		case sim.TypeLoadBalancer:
			c = color.RGBA{100, 150, 255, 255}
		case sim.TypeMessageQueue:
			c = color.RGBA{0, 255, 255, 255}
		case sim.TypeDatabase:
			c = color.RGBA{200, 100, 255, 255}
		case sim.TypeCache:
			c = color.RGBA{255, 200, 100, 255}
		}

		// Draw icon
		vector.FillRect(screen, x, y, 40, 40, c, true)

		// Draw label and Cost
		label := fmt.Sprintf("%s\n%d", template.Name, template.Cost)
		ebitenutil.DebugPrintAt(screen, label, int(x+50), int(y+50))
	}

}
