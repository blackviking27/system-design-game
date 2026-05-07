package engine

import (
	"fmt"
	"image/color"

	"github.com/blackviking27/system-design-game/internal/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	imgWidth  int = 64
	imgHeight int = 64
)

var (
	colorLB            = color.RGBA{R: 100, G: 150, B: 255, A: 255}
	colorServerOK      = color.RGBA{R: 100, G: 255, B: 250, A: 255}
	colorServerFailing = color.RGBA{R: 255, G: 100, B: 100, A: 255}
	colorLine          = color.RGBA{R: 150, G: 150, B: 150, A: 255}
	colorLineDrawing   = color.RGBA{R: 255, G: 200, B: 0, A: 255}

	colorMessageQueue = color.RGBA{0, 255, 255, 255}
	colorDB           = color.RGBA{200, 100, 255, 255}
	colorCache        = color.RGBA{255, 200, 100, 255}
)

func DrawNetwork(screen *ebiten.Image, game *GameplayScene) {
	// Draw links (lines) first so they render underneath nodes
	for _, node := range game.Network.Nodes {
		for _, out := range node.Outbound {
			vector.StrokeLine(screen, float32(node.X), float32(node.Y), float32(out.X), float32(out.Y), 2, colorLine, true)
		}
	}

	// Draw in-progress link
	if game.linkingNode != nil {
		vector.StrokeLine(screen, float32(game.linkingNode.X), float32(game.linkingNode.Y), float32(game.mouseX), float32(game.mouseY), 2, colorLineDrawing, true)
	}

	// Draw nodes(servers)
	for _, node := range game.Network.Nodes {
		isFailing := len(node.Queue) >= node.MaxRam

		var img *ebiten.Image
		if isFailing {
			img = ui.NodeImagesRed[node.Type]
		} else {
			img = ui.NodeImages[node.Type]
		}

		// Center the img on the x,y coordinates
		w, h := float64(imgWidth), float64(imgHeight)
		startX, startY := float64(node.X)-w/2, float64(node.Y)-h/2

		if img != nil {
			op := &ebiten.DrawImageOptions{}

			// Optional: Scale image to fit 80x50 exactly if it's too big/small
			bounds := img.Bounds()
			scaleX := w / float64(bounds.Dx())
			scaleY := h / float64(bounds.Dy())
			op.GeoM.Scale(scaleX, scaleY)

			op.GeoM.Translate(startX, startY)
			screen.DrawImage(img, op)
		} else {
			// Fallback
			vector.FillRect(screen, float32(startX), float32(startY), float32(w), float32(h), color.RGBA{100, 255, 250, 255}, true)
		}

		// Stats for node
		stats := fmt.Sprintf("%s\nRAM: %d/%d\nDrop: %d", node.ID, len(node.Queue), node.MaxRam, node.DroppedCount)

		// Adjust text above rectangle
		ebitenutil.DebugPrintAt(screen, stats, int(startX), int(startY)-50)
	}

	// Draw the catalog tray
	ui.DrawTray(screen, game.CurrentBudget)

	// Draw the HUD
	isGameOver := game.State == StateGameOver
	isVictory := game.State == StateVictory
	ui.DrawHUD(screen, game.Network, game.Level.Name, game.Level.TargetUptimeTicks, game.Level.MaxDroppedPackets, isGameOver, isVictory)
}
