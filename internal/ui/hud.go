package ui

import (
	"fmt"
	"image/color"

	"github.com/blackviking27/system-design-game/internal/sim"
	"github.com/hajimehoshi/ebiten/v2"
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

	DrawText(screen, stats, 10, 20, 14, colorText)

	// Game overlay for loss or victory
	if isGameOver || isVictory {
		screen.Fill(colorOverlay)
		w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
		msg := "CRITICAL SYSTEM FAILURE\n\nToo many packets dropped.\nGAME OVER"
		if isVictory {
			msg = "SYSTEM STABLE\n\nYou survived the traffic surge.\nVICTORY!"
		}
		// Approximate centering
		DrawText(screen, msg, float64(w/2-100), float64(h/2-20), 20, colorText)
	}
}
