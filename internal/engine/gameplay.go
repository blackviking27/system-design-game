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

func (this *GameplayScene) legacyCheckWinOrLoseCondition() {
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

func (this *GameplayScene) evaluateCondition(condition Condition) bool {
	if condition.AfterTick > 0 && int(this.Network.TickCount) < condition.AfterTick {
		return false
	}

	actual, ok := this.metricValue(condition)
	if !ok {
		return false
	}
	return compareMetric(actual, condition.Operator, condition.Value)
}

func (this *GameplayScene) checkForLoss(conditions []Condition) bool {
	for _, condition := range conditions {
		if this.evaluateCondition(condition) {
			return true
		}
	}
	return false
}

func (this *GameplayScene) metricValue(condition Condition) (float64, bool) {
	metrics := this.Network.Metrics

	switch condition.Metric {
	case "uptime_ticks":
		return float64(this.Network.TickCount), true
	case "generated_packets":
		if condition.PacketType != "" {
			return float64(metrics.GeneratedByType[condition.PacketType]), true
		}
		return float64(metrics.GeneratedPackets), true
	case "completed_packets":
		if condition.PacketType != "" {
			return float64(metrics.CompletedByType[condition.PacketType]), true
		}
		return float64(metrics.CompletedPackets), true
	case "dropped_packets":
		if condition.PacketType != "" {
			return float64(metrics.DroppedByType[condition.PacketType]), true
		}
		return float64(metrics.DroppedPackets), true
	case "processed_packets":
		return float64(metrics.ProcessedPackets), true
	case "success_rate":
		generated := metrics.GeneratedPackets
		completed := metrics.CompletedPackets

		if condition.PacketType != "" {
			generated = metrics.GeneratedByType[condition.PacketType]
			completed = metrics.CompletedByType[condition.PacketType]
		}

		if generated == 0 {
			return 0, true
		}
		return float64(completed) / float64(generated), true
	case "drop_rate":
		generated := metrics.GeneratedPackets
		dropped := metrics.DroppedPackets

		if condition.PacketType != "" {
			generated = metrics.GeneratedByType[condition.PacketType]
			dropped = metrics.DroppedByType[condition.PacketType]
		}

		if generated == 0 {
			return 0, true
		}
		return float64(dropped) / float64(generated), true
	case "avg_latency_ticks":
		if condition.PacketType != "" {
			completed := metrics.CompletedByType[condition.PacketType]
			if completed == 0 {
				return 0, true
			}
			return float64(metrics.TotalLatencyTicksByType[condition.PacketType]) / float64(completed), true
		}

		if metrics.CompletedPackets == 0 {
			return 0, true
		}
		return float64(metrics.TotalLatencyTicks) / float64(metrics.CompletedPackets), true
	case "queue_depth":
		total := 0
		for _, node := range this.Network.Nodes {
			total += len(node.Queue)
		}
		return float64(total), true
	default:
		return 0, false
	}
}

func compareMetric(actual float64, operator string, expected float64) bool {
	switch operator {
	case ">":
		return actual > expected
	case ">=":
		return actual >= expected
	case "<":
		return actual < expected
	case "<=":
		return actual <= expected
	case "==":
		return actual == expected
	case "!=":
		return actual != expected
	default:
		return false
	}
}

func (this *GameplayScene) checkForWin(conditions []Condition) bool {
	for _, condition := range conditions {
		if !this.evaluateCondition(condition) {
			return false
		}
	}
	return true
}

func (this *GameplayScene) checkWinOrLoseCondition() {
	if this.checkForLoss(this.Level.LossConditions) {
		this.State = types.StateGameOver
		return
	}

	if len(this.Level.WinConditions) > 0 && this.checkForWin(this.Level.WinConditions) {
		this.State = types.StateVictory
		return
	}

	if len(this.Level.WinConditions) == 0 {
		this.legacyCheckWinOrLoseCondition()
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

	// Resetting the metrics
	this.Network.Metrics = sim.NewMetrics()
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

	// Add metrics to the network w.r.t packet generated
	packetType := pickPacketType(this.currentPacketMix)
	this.Network.RecordGenerated(packetType)

	return &sim.Packet{
		ID:          id,
		TraceId:     id,
		Type:        packetType,
		InitialType: packetType,
		CreatedAt:   this.Network.TickCount,
		Payload:     make(map[string]any),
		Metadata:    make(map[string]int),
	}
}

func NewGameplayScene(levelPath string) *GameplayScene {
	// Load level json data
	lvl, err := LoadLevel(levelPath)
	if err != nil {
		panic(err)
	}

	// Initializing the sim network
	network := &sim.Network{
		Nodes:               make(map[string]*sim.Node),
		Metrics:             sim.NewMetrics(),
		TerminalPacketTypes: make(map[string]bool),
	}

	for _, terminalPackets := range lvl.TerminalPacketTypes {
		network.TerminalPacketTypes[terminalPackets] = true
	}

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
	// Create a load balancer on the left, clear of the HUD and bottom tray.
	lb := scene.CreateNodeFromTemplate(lbCatalogNode, 180, 320)
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
