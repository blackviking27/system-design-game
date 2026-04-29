package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Defining the game structure
type Game struct {
	CurrentScene Scene
}

// Runs the simulation 6 ticks per second
const framesPerTick = 10

func (this *Game) Update() error {
	// Update the current scene. It will return the next scene to play
	nextScene, err := this.CurrentScene.Update()
	if err != nil {
		return err
	}

	if nextScene != nil {
		this.CurrentScene = nextScene
	}
	return nil
}

func (this *Game) Draw(screen *ebiten.Image) {
	if this.CurrentScene != nil {
		this.CurrentScene.Draw(screen)
	}
}

func (this *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}
