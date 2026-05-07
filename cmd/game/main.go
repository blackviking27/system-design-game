package main

import (
	"log"

	"github.com/blackviking27/system-design-game/internal/engine"
	"github.com/blackviking27/system-design-game/internal/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Loading game assets
	if err := ui.LoadAssets(); err != nil {
		log.Fatalf("Failed to load assets: %v", err)
	}

	game := &engine.Game{
		CurrentScene: &engine.MainMenuScene{},
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("System design Sandbox")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Running the game
	err := ebiten.RunGame(game)
	if err != nil {
		panic(err)
	}
}
