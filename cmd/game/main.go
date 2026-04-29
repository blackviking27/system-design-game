package main

import (
	"github.com/blackviking27/system-design-game/internal/engine"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	game := &engine.Game{
		CurrentScene: &engine.MainMenuScene{},
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("System design Sandbox")

	// Running the game
	err := ebiten.RunGame(game)
	if err != nil {
		panic(err)
	}
}
