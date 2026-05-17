package engine

import (
	"fmt"
	"time"

	"github.com/blackviking27/system-design-game/internal/sim"
	"github.com/blackviking27/system-design-game/internal/types"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Gameplay scene
type GameplayScene struct {
	Network   *sim.Network
	tickTimer int

	// Game level
	Level *Level
	State types.GameState

	// level budget
	CurrentBudget int

	// Traffic rate
	currentTrafficRate  int
	nextTrafficEventIdx int

	// Game controls
	draggingNode *sim.Node
	dragOffsetX  float64
	dragOffsetY  float64
	linkingNode  *sim.Node
	mouseX       int
	mouseY       int

	// Screen dimensions
	screenWidth  int
	screenHeight int
}

func NewGameplayScene(levelPath string) *GameplayScene {
	// Load level json data
	lvl, err := LoadLevel(levelPath)
	if err != nil {
		panic(err)
	}

	// Initializing the sim network
	network := &sim.Network{Nodes: make(map[string]*sim.Node)}

	// Creating the engine game wrapper
	scene := &GameplayScene{
		Network:       network,
		Level:         lvl,
		State:         types.StateDesigning,
		CurrentBudget: lvl.StartingBudget,
	}

	// Create a load balancer
	lb := scene.CreateNodeFromTemplate(string(sim.TypeLoadBalancer), 400, 150)
	network.Nodes[lb.ID] = lb

	return scene
}

func generateId() string {
	return fmt.Sprintf("node-%v", time.Now().UnixNano())
}

func (this *GameplayScene) Update() (Scene, error) {
	// Handling user input
	this.HandleInput()

	if this.State != types.StateSimulating {
		// Phase transition shortcuts
		if this.State == types.StateDesigning && inpututil.IsKeyJustPressed(ebiten.KeyS) {
			this.State = types.StateSimulating
		}

		if (this.State == types.StateGameOver || this.State == types.StateVictory) && ebiten.IsKeyPressed(ebiten.KeyR) {
			this.Reset()
		}

		return this, nil
	}

	this.tickTimer += 1

	if this.tickTimer >= framesPerTick {
		this.Network.Tick()

		// Dynamic traffic rate
		for this.nextTrafficEventIdx < len(this.Level.TrafficPattern) && this.Network.TickCount > uint64(this.Level.TrafficPattern[this.nextTrafficEventIdx].StartTick) {
			this.currentTrafficRate = this.Level.TrafficPattern[this.nextTrafficEventIdx].Rate
			this.nextTrafficEventIdx++
		}

		trafficRate := this.currentTrafficRate
		// Manually increasing traffic
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			trafficRate = 50
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
		this.State = types.StateGameOver
	}

	// Win condition: Survived for the duration
	if int(this.Network.TickCount) >= this.Level.TargetUptimeTicks {
		this.State = types.StateVictory
	}

}

func (this *GameplayScene) Draw(screen *ebiten.Image) {
	this.screenWidth = screen.Bounds().Dx()
	this.screenHeight = screen.Bounds().Dy()
	DrawNetwork(screen, this)
}

func (this *GameplayScene) Reset() {
	this.Network.TickCount = 0
	this.tickTimer = 0
	this.nextTrafficEventIdx = 0
	this.currentTrafficRate = 0
	this.State = types.StateDesigning

	for _, node := range this.Network.Nodes {
		node.ResetState()
	}

}

func (this *GameplayScene) CreateNodeFromTemplate(templateName string, x, y float64) *sim.Node {
	// 1. Get blueprint from loaded level data
	template, exists := this.Level.NodeTemplates[templateName]
	if !exists {
		return nil
	}

	// 2. Create the node object
	node := sim.NewNode(
		generateId(),
		template.Type,
		100,
		5,
		100,
	)

	node.X, node.Y = x, y

	// 3. Assigning workflow
	node.Processor = &sim.WorkflowProcessor{
		Workflows: template.Workflows,
	}

	// 4. Assign router
	node.Router = &sim.RoundRobinRouter{}

	return node
}
