package engine

import (
	"fmt"
	"image/color"

	"github.com/blackviking27/system-design-game/internal/sim"
	"github.com/blackviking27/system-design-game/internal/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
		nodeColor := colorServerOK

		// Color based on server and RAM limits
		if len(node.Queue) >= node.MaxRam {
			nodeColor = colorServerFailing
		} else {
			switch node.Type {
			case sim.TypeLoadBalancer:
				nodeColor = colorLB
			case sim.TypeMessageQueue:
				nodeColor = colorMessageQueue
			case sim.TypeDatabase:
				nodeColor = colorDB
			case sim.TypeCache:
				nodeColor = colorCache
			}
		}

		// Center the rectangle on the x,y coordinates
		w, h := float32(80), float32(50)
		startX, startY := float32(node.X)-w/2, float32(node.Y)-h/2

		// Draw the node block
		vector.FillRect(screen, startX, startY, w, h, nodeColor, true)

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
