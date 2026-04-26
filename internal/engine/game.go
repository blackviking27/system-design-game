package engine

import (
	"fmt"
	"time"

	"github.com/blackviking27/system-design-game/internal/sim"
	"github.com/hajimehoshi/ebiten/v2"
)

// Defining the game structure
type Game struct {
	Network   *sim.Network
	tickTimer int

	// Game level
	Level Level
	State GameState

	// Game controls
	draggingNode *sim.Node
	dragOffsetX  float64
	dragOffsetY  float64
	linkingNode  *sim.Node
	mouseX       int
	mouseY       int
}

// Runs the simulation 6 ticks per second
const framesPerTick = 10

func (this *Game) Update() error {
	// Handling user input
	this.HandleInput()

	// only run the game if we are not in a terminal game state
	if this.State != StatePlaying {
		// Game restart logic
		if ebiten.IsKeyPressed(ebiten.KeySpace) {

		}
		return nil
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
	return nil
}

func (this *Game) checkWinOrLoseCondition() {
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

func (this *Game) Draw(screen *ebiten.Image) {
	DrawNetwork(screen, this)
}

func (this *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}
