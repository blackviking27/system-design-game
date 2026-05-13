package engine

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/blackviking27/system-design-game/internal/sim"
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

func drawConnectionLine(screen *ebiten.Image, x1, y1, x2, y2 float32, isSimulating bool, color color.RGBA) {
	//1. Drawing the baseline
	vector.StrokeLine(screen, x1, y1, x2, y2, 2, color, true)

	// 2. Vector match for direction and midpoint
	dx := x2 - x1
	dy := y2 - y1
	dist := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	if dist < 1 {
		return
	}

	ux := dx / dist
	uy := dy / dist

	midX := x1 + dx/2
	midY := y1 + dy/2

	// 3. Draw direction arrows
	arrowSize := float32(8)
	px := -uy
	py := ux

	// Arrow head points (V shape)
	ax1 := midX - ux*arrowSize + px*arrowSize*0.6
	ay1 := midY - uy*arrowSize + py*arrowSize*0.6
	ax2 := midX - ux*arrowSize - px*arrowSize*0.6
	ay2 := midY - uy*arrowSize - py*arrowSize*0.6

	vector.StrokeLine(screen, midX, midY, ax1, ay1, 2, color, true)
	vector.StrokeLine(screen, midX, midY, ax2, ay2, 2, color, true)

	// 4. Draw Flow Animation (Pulses)
	if isSimulating {
		// Create 3 pulses moving along the line
		timeScale := float32(time.Now().UnixMilli()%1000) / 1000.0 // 0.0 to 1.0
		for i := 0; i < 3; i++ {
			t := timeScale + float32(i)/3.0
			if t > 1.0 {
				t -= 1.0
			}
			px := x1 + dx*t
			py := y1 + dy*t
			vector.FillCircle(screen, px, py, 3, types.ColorTextYellow, true)
		}
	}

}

func drawNode(screen *ebiten.Image, node *sim.Node, isDragging bool) {
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

		if isDragging {
			// Making it see through when dragging the node
			op.ColorScale.Scale(1, 1, 1, 0.5)
		}

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

func DrawNetwork(screen *ebiten.Image, game *GameplayScene) {
	// Fill the background
	screen.Fill(types.ColorBackground)

	isSimulationRunning := game.State == types.StateSimulating

	// Draw links (lines) first so they render underneath nodes
	for _, node := range game.Network.Nodes {
		for _, out := range node.Outbound {
			drawConnectionLine(screen, float32(node.X), float32(node.Y), float32(out.X), float32(out.Y), isSimulationRunning, types.ColorLine)
		}
	}

	// Draw in-progress link
	if game.linkingNode != nil {
		drawConnectionLine(screen, float32(game.linkingNode.X), float32(game.linkingNode.Y), float32(game.mouseX), float32(game.mouseY), false, types.ColorLine)
	}

	// Draw nodes(servers)
	for _, node := range game.Network.Nodes {
		drawNode(screen, node, false)
	}

	if game.draggingNode != nil {
		drawNode(screen, game.draggingNode, true)
	}

	// Draw the catalog tray
	ui.DrawTray(screen, game.CurrentBudget)

	// Draw the HUD
	ui.DrawHUD(screen, game.Network, game.Level.Name, game.Level.TargetUptimeTicks, game.Level.MaxDroppedPackets, game.State)
}
