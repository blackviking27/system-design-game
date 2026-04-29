package engine

import (
	"fmt"
	"time"

	"github.com/blackviking27/system-design-game/internal/sim"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameState string

const (
	StatePlaying  GameState = "Playing"
	StateGameOver GameState = "GameOver"
	StateVictory  GameState = "Victory"
)

// level defines the parameters and win conditions for a scenario
type Level struct {
	Name              string
	TargetUptimeTicks int // How many ticks the user needs to survive
	MaxDroppedPackets int // How many dropped packets cause a game over
	BaseTrafficRate   int // Packets generated per tick
	StartingBudget    int // Starting budget of the game
}

// Gameplay scene
type GameplayScene struct {
	Network   *sim.Network
	tickTimer int

	// Game level
	Level Level
	State GameState

	// level budget
	CurrentBudget int

	// Game controls
	draggingNode *sim.Node
	dragOffsetX  float64
	dragOffsetY  float64
	linkingNode  *sim.Node
	mouseX       int
	mouseY       int
}

func NewGameplayScene() *GameplayScene {
	// Initializing the sim network
	network := &sim.Network{Nodes: make(map[string]*sim.Node)}

	// Create a load balancer
	lb := sim.NewNode("LB-Main", sim.TypeLoadBalancer, 100, 100, 0)
	lb.X, lb.Y = 400, 150

	// Create 2 weak servers
	serverA := sim.NewNode("Server-A", sim.TypeServer, 10, 1, 0)
	serverA.X, serverA.Y = 200, 450

	serverB := sim.NewNode("Server-B", sim.TypeServer, 10, 1, 0)
	serverB.X, serverB.Y = 600, 450

	// Creating connections
	lb.LinkTo(serverA)
	lb.LinkTo(serverB)

	network.Nodes[lb.ID] = lb
	network.Nodes[serverA.ID] = serverA
	network.Nodes[serverB.ID] = serverB

	// Defining the level
	levelOne := Level{
		Name:              "Level 1",
		TargetUptimeTicks: 200,
		MaxDroppedPackets: 50,
		BaseTrafficRate:   3,
		StartingBudget:    400,
	}

	// Creating the engine game wrapper
	return &GameplayScene{
		Network:       network,
		Level:         levelOne,
		State:         StatePlaying,
		CurrentBudget: levelOne.StartingBudget,
	}
}

func (this *GameplayScene) Update() (Scene, error) {
	// Handling user input
	this.HandleInput()

	// only run the game if we are not in a terminal game state
	if this.State != StatePlaying {
		// Game restart logic
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			return &MainMenuScene{}, nil
		}
		return this, nil
	}

	this.tickTimer += 1

	if this.tickTimer >= framesPerTick {
		this.Network.Tick()

		// Simulating Traffic
		// Continously feed traffic to the load balancer for simulation
		trafficRate := this.Level.BaseTrafficRate

		// Manually increasing traffic
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			trafficRate = 20
		}

		for _, node := range this.Network.Nodes {
			if node.Type == sim.TypeLoadBalancer {
				for i := 0; i < trafficRate; i++ {
					node.Queue = append(node.Queue, &sim.Packet{ID: fmt.Sprintf("pkt-%v", time.Now().Unix()/int64(time.Microsecond))})
				}
			}
		}

		// Check win or lose condition
		this.checkWinOrLoseCondition()
		this.tickTimer = 0
	}
	return this, nil
}

func (this *GameplayScene) checkWinOrLoseCondition() {
	// Counting total dropped packets in the current state of game
	totalDroppedPacket := 0
	for _, node := range this.Network.Nodes {
		totalDroppedPacket += node.DroppedCount
	}

	// Loss condition: Too many packets dropped
	if totalDroppedPacket >= this.Level.MaxDroppedPackets {
		this.State = StateGameOver
	}

	// Win condition: Survived for the duration
	if int(this.Network.TickCount) >= this.Level.TargetUptimeTicks {
		this.State = StateVictory
	}

}

func (this *GameplayScene) Draw(screen *ebiten.Image) {
	DrawNetwork(screen, this)
}
