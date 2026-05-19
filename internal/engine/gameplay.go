package engine

import (
	"fmt"
	"math/rand"
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
	currentPacketMix    []PacketMixEntry

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
			trafficLevel := this.Level.TrafficPattern[this.nextTrafficEventIdx]
			this.currentTrafficRate = trafficLevel.Rate
			this.currentPacketMix = trafficLevel.PacketMix
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
					node.Queue = append(node.Queue, this.GeneratePacket())
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
	this.currentPacketMix = nil
	this.State = types.StateDesigning

	for _, node := range this.Network.Nodes {
		node.ResetState()
	}

}

func (this *GameplayScene) CreateNodeFromTemplate(catalogNodeTemplate sim.CatalogNodeTemplate, x, y float64) *sim.Node {
	// 1. Get node workflow from loaded level data
	workflowForNode, exists := this.Level.NodeTemplates[string(catalogNodeTemplate.Type)]
	if !exists {
		return nil
	}

	node := sim.NewNode(
		generateId(),
		catalogNodeTemplate.Type,
		catalogNodeTemplate.MaxRam,
		catalogNodeTemplate.ProcessPower,
		catalogNodeTemplate.Cost,
	)

	node.X, node.Y = x, y

	// 3. Assigning workflow
	node.Processor = &sim.WorkflowProcessor{
		Workflows: workflowForNode.Workflows,
	}

	// 4. Assign router
	node.Router = &sim.RoundRobinRouter{}

	return node
}

func (this *GameplayScene) GeneratePacket() *sim.Packet {
	id := generateId()

	return &sim.Packet{
		ID:       id,
		TraceId:  id,
		Type:     pickPacketType(this.currentPacketMix),
		Payload:  make(map[string]any),
		Metadata: make(map[string]int),
	}
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

	lbCatalogNode, exists := sim.GetCatalogTemplateForType(sim.TypeLoadBalancer)
	if !exists {
		panic("No Loadbalancer catalog node found")
	}
	// Create a load balancer
	lb := scene.CreateNodeFromTemplate(lbCatalogNode, 400, 150)
	network.Nodes[lb.ID] = lb

	return scene
}

func generateId() string {
	return fmt.Sprintf("node-%v", time.Now().UnixNano())
}

func pickPacketType(packetMix []PacketMixEntry) string {
	totalWeight := 0
	for _, entry := range packetMix {
		if entry.Weight > 0 {
			totalWeight += entry.Weight
		}
	}

	if totalWeight == 0 {
		return ""
	}

	pick := rand.Intn(totalWeight)
	cumulative := 0
	for _, entry := range packetMix {
		if entry.Weight <= 0 {
			continue
		}

		cumulative += entry.Weight
		if pick < cumulative {
			return entry.Type
		}
	}
	return ""
}
