package ui

import (
	"fmt"

	"github.com/blackviking27/system-design-game/internal/sim"
	"github.com/blackviking27/system-design-game/internal/types"
	"github.com/hajimehoshi/ebiten/v2"
)

func DrawHUD(screen *ebiten.Image, network *sim.Network, levelName string, targetTick, maxDropped int, gameState types.GameState) {
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

	DrawText(screen, stats, 10, 20, 14, types.ColorText)

	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()

	if gameState == types.StateDesigning {
		DrawText(screen, "DESIGN MODE - [Press S to Start Sim]", 10, 140, 16, types.ColorTextGreen)
	} else if gameState == types.StateSimulating {
		DrawText(screen, "SIMULATING...", 10, 140, 16, types.ColorTextYellow)
	}

	if gameState == types.StateGameOver {
		DrawText(screen, "GAME OVER...\n\n[Press R to Retry Design]", 10, 140, 16, types.ColorTextRed)
	}
	// Game overlay for loss or victory
	if gameState == types.StateVictory {
		screen.Fill(types.ColorOverlay)
		msg := "SYSTEM STABLE\n\nYou survived the traffic surge.\nVICTORY!"
		// Approximate centering
		DrawText(screen, msg, float64(w/2-100), float64(h/2-20), 20, types.ColorText)
	}
}
