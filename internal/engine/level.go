package engine

// GameState represents the current state of the game loop

type GameState string

const (
	StatePlaying  GameState = "Playing"
	StateGameOver GameState = "Game Over"
	StateVictory  GameState = "Victory"
)

// level defines the parameters and win conditions for a scenario
type Level struct {
	Name              string
	TargetUptimeTicks int // How many ticks the user needs to survive
	MaxDroppedPackets int // How many dropped packets cause a game over
	BaseTrafficRate   int // Packets generated per tick
}
