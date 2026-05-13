package engine

import (
	"fmt"
	"image/color"
	"time"

	"github.com/blackviking27/system-design-game/internal/types"
	"github.com/blackviking27/system-design-game/internal/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	imgWidth  int = 64
	imgHeight int = 64
)

func DrawNetwork(screen *ebiten.Image, game *GameplayScene) {
	// Fill the background
	screen.Fill(types.ColorBackground)

	// Draw links (lines) first so they render underneath nodes
	for _, node := range game.Network.Nodes {
		for _, out := range node.Outbound {
			vector.StrokeLine(screen, float32(node.X), float32(node.Y), float32(out.X), float32(out.Y), 2, types.ColorLine, true)
		}
	}

	// Draw in-progress link
	if game.linkingNode != nil {
		vector.StrokeLine(screen, float32(game.linkingNode.X), float32(game.linkingNode.Y), float32(game.mouseX), float32(game.mouseY), 2, types.ColorLineDrawing, true)
	}

	// Draw nodes(servers)
	for _, node := range game.Network.Nodes {
		isFailing := len(node.Queue) >= node.MaxRam

		// 1. Pick the correct animation array
		var frames []*ebiten.Image
		if isFailing {
			frames = ui.NodeRedFrames[node.Type]
		} else {
			frames = ui.NodeFrames[node.Type]
		}

		// 2. Determin which frame to show in the selected array
		var frameindex int
		if len(frames) > 0 {
			animIndex := (time.Now().UnixMilli() / 500) % int64(len(frames))
			frameindex = int(animIndex)
		}

		// Center the img on the x,y coordinates
		w, h := float64(imgWidth), float64(imgHeight)
		startX, startY := float64(node.X)-w/2, float64(node.Y)-h/2

		// 3. Draw the frame
		if len(frames) > frameindex && frames[frameindex] != nil {
			img := frames[frameindex]
			op := &ebiten.DrawImageOptions{}

			// Scaling to fit imgWidh X imgHeight respectively
			bounds := img.Bounds()
			scaleX := w / float64(bounds.Dx())
			scaleY := h / float64(bounds.Dy())
			op.GeoM.Scale(scaleX, scaleY)

			op.GeoM.Translate(startX, startY)
			screen.DrawImage(img, op)
		} else {
			var fallbackColor color.RGBA
			if isFailing {
				fallbackColor = color.RGBA{255, 100, 100, 255} // Red fallback
			} else {
				fallbackColor = color.RGBA{100, 255, 250, 255} // Normal fallback
			}
			vector.FillRect(screen, float32(startX), float32(startY), float32(w), float32(h), fallbackColor, true)
		}

		// Stats for node
		stats := fmt.Sprintf("%s\nRAM: %d/%d\nDrop: %d", node.ID, len(node.Queue), node.MaxRam, node.DroppedCount)

		// Adjust text above rectangle
		ebitenutil.DebugPrintAt(screen, stats, int(startX), int(startY)-50)
	}

	// Draw the catalog tray
	ui.DrawTray(screen, game.CurrentBudget)

	// Draw the HUD
	ui.DrawHUD(screen, game.Network, game.Level.Name, game.Level.TargetUptimeTicks, game.Level.MaxDroppedPackets, game.State)
}
