package engine

import "github.com/hajimehoshi/ebiten/v2"

// Represents a single state/screen in game
type Scene interface {
	Update() (Scene, error)
	Draw(screen *ebiten.Image)
}
