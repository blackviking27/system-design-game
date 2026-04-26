package ui

import (
	"fmt"
	"image/color"

	"github.com/blackviking27/system-design-game/internal/sim"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	colorOverlay = color.RGBA{0, 0, 0, 200}
	colorText    = color.RGBA{255, 255, 255, 255}
)

func DrawHUD(screen *ebiten.Image, network *sim.Network, levelName string, targetTick, maxDropped int, isGameOver, isVictory bool) {
	// Calculate the total drop packets
	totalDropped := 0
	for _, node := range network.Nodes {
		totalDropped += node.DroppedCount
	}

	// Top left stats
	stats := fmt.Sprintf("Level:%s\nUPTIME: %d / %d Ticks\nDROPPED: %d / %d Max\n\n[HOLD SPACE FOR TRAFFIC SPIKE]",
		levelName,
		network.TickCount,
		targetTick,
		totalDropped,
		maxDropped,
	)

	ebitenutil.DebugPrintAt(screen, stats, 10, 10)

	// Game overlay for loss or victory
	if isGameOver {
		screen.Fill(colorOverlay)
		ebitenutil.DebugPrintAt(screen, "CRITICAL SYSTEM FAILURE\n\nToo many packets dropped.\nGAME OVER", 330, 280)
	} else if isVictory {
		screen.Fill(colorOverlay)
		ebitenutil.DebugPrintAt(screen, "SYSTEM STABLE\n\nYou survived the traffic surge.\nVICTORY!", 330, 280)

	}

}
