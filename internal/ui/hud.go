package ui

import (
	"fmt"
	"strings"

	"github.com/blackviking27/system-design-game/internal/sim"
	"github.com/blackviking27/system-design-game/internal/types"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	hudX          = 10.0
	statsY        = 20.0
	statsFontSize = 14.0
	stateFontSize = 16.0
)

func DrawHUD(screen *ebiten.Image, network *sim.Network, levelName string, targetTick, maxDropped int, gameState types.GameState) {
	// Game metrics
	generated := network.Metrics.GeneratedPackets
	completed := network.Metrics.CompletedPackets
	dropped := network.Metrics.DroppedPackets

	successRate := 0.0
	dropRate := 0.0
	if generated > 0 {
		successRate = float64(completed) / float64(generated)
		dropRate = float64(dropped) / float64(generated)
	}

	avgLatency := 0.0
	if completed > 0 {
		avgLatency = float64(network.Metrics.TotalLatencyTicks) / float64(completed)
	}

	// Top left stats
	stats := fmt.Sprintf(
		"Level:%s\nUPTIME: %d / %d Ticks\nGENERATED: %d\nCOMPLETED: %d\nSUCCESS: %.1f%%\nDROPPED: %d / %d Max\nDROP RATE: %.1f%%\nAVG LATENCY: %.1f ticks\n\n[HOLD SPACE FOR TRAFFIC SPIKE]",
		levelName,
		network.TickCount,
		targetTick,
		generated,
		completed,
		successRate*100,
		dropped,
		maxDropped,
		dropRate*100,
		avgLatency,
	)

	DrawText(screen, stats, hudX, statsY, statsFontSize, types.ColorText)

	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	statsLines := strings.Count(stats, "\n") + 1
	stateY := statsY + float64(statsLines)*statsFontSize*1.2 + 8

	if gameState == types.StateDesigning {
		DrawText(screen, "DESIGN MODE - [Press S to Start Sim]", hudX, stateY, stateFontSize, types.ColorTextGreen)
	} else if gameState == types.StateSimulating {
		DrawText(screen, "SIMULATING...", hudX, stateY, stateFontSize, types.ColorTextYellow)
	}

	if gameState == types.StateGameOver {
		DrawText(screen, "GAME OVER...\n\n[Press R to Retry Design]", hudX, stateY, stateFontSize, types.ColorTextRed)
	}
	// Game overlay for loss or victory
	if gameState == types.StateVictory {
		screen.Fill(types.ColorOverlay)
		msg := "SYSTEM STABLE\n\nYou survived the traffic surge.\nVICTORY!"
		// Approximate centering
		DrawText(screen, msg, float64(w/2-100), float64(h/2-20), 20, types.ColorText)
	}
}
