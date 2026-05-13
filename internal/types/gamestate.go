package types

type GameState string

const (
	StateDesigning  GameState = "Designing"
	StateSimulating GameState = "Simulating"
	StateGameOver   GameState = "GameOver"
	StateVictory    GameState = "Victory"
)
