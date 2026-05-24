package sim

type Metrics struct {
	GeneratedPackets        int
	CompletedPackets        int
	DroppedPackets          int
	ProcessedPackets        int
	GeneratedByType         map[string]int
	CompletedByType         map[string]int
	DroppedByType           map[string]int
	TotalLatencyTicks       int
	TotalLatencyTicksByType map[string]int
}

// Network stores the global state of the simulation
type Network struct {
	Nodes               map[string]*Node
	TickCount           uint64
	Metrics             Metrics
	TerminalPacketTypes map[string]bool
}

// Tick advances the simulation state by one step
func (this *Network) Tick() {
	this.EnsureMetrics()

	for _, node := range this.Nodes {
		// Skip nodes that aren't really configured
		if node.Processor == nil || node.Router == nil {
			continue
		}

		// 1. PROCESS: Handle packets based on process power
		packetProcessedThisTick := 0
		for len(node.Queue) > 0 && packetProcessedThisTick < node.ProcessPower {
			packet := node.Queue[0]
			node.Queue = node.Queue[1:]

			// Delegate logic to processor (workflow)\
			outboundPackets, err := node.Processor.Process(node, packet)
			if err != nil {
				packet.Status = StatusDropped
				node.DroppedCount++
				this.RecordDropped(packet)
				continue
			}

			// 2. ROUTE: Send the response to the next destination
			for _, outPkt := range outboundPackets {
				if this.IsTerminalPacket(outPkt) {
					this.RecordCompleted(outPkt)
					continue
				}

				target := node.Router.Route(node, outPkt)
				if target != nil && len(target.Queue) < target.MaxRam {
					target.Queue = append(target.Queue, outPkt)
				} else {
					outPkt.Status = StatusDropped
					node.DroppedCount++
					this.RecordDropped(outPkt)
				}
			}
			node.ProcessedCount++
			packetProcessedThisTick++
			this.Metrics.ProcessedPackets++
		}
	}
	this.TickCount += 1
}

func (this *Network) RemoveNode(nodeID string) {

	_, exists := this.Nodes[nodeID]
	if !exists {
		return
	}

	// remove the node from the network
	delete(this.Nodes, nodeID)

	// remove the connection to node from other
	for _, node := range this.Nodes {
		var newOutBoundConnections []*Node
		for _, targetNode := range node.Outbound {
			if targetNode.ID != nodeID {
				newOutBoundConnections = append(newOutBoundConnections, targetNode)
			}
		}
		node.Outbound = newOutBoundConnections
	}
}

func (this *Network) IsTerminalPacket(packet *Packet) bool {
	if packet == nil {
		return false
	}
	return this.TerminalPacketTypes[packet.Type]
}

func (this *Network) EnsureMetrics() {
	if this.Metrics.GeneratedByType == nil {
		this.Metrics.GeneratedByType = make(map[string]int)
	}
	if this.Metrics.CompletedByType == nil {
		this.Metrics.CompletedByType = make(map[string]int)
	}
	if this.Metrics.DroppedByType == nil {
		this.Metrics.DroppedByType = make(map[string]int)
	}
	if this.Metrics.TotalLatencyTicksByType == nil {
		this.Metrics.TotalLatencyTicksByType = make(map[string]int)
	}
}

func (this *Network) RecordGenerated(packetType string) {
	this.EnsureMetrics()
	this.Metrics.GeneratedPackets++
	this.Metrics.GeneratedByType[packetType]++
}

func (this *Network) RecordCompleted(packet *Packet) {
	if packet == nil {
		return
	}

	this.EnsureMetrics()
	packet.Status = StatusProcessed
	latency := int(this.TickCount - packet.CreatedAt)

	this.Metrics.CompletedPackets++
	this.Metrics.CompletedByType[packet.InitialType]++
	this.Metrics.TotalLatencyTicks += latency
	this.Metrics.TotalLatencyTicksByType[packet.InitialType] += latency
}

func (this *Network) RecordDropped(packet *Packet) {
	if packet == nil {
		return
	}

	this.EnsureMetrics()
	packet.Status = StatusDropped
	this.Metrics.DroppedPackets++
	this.Metrics.DroppedByType[packet.InitialType]++
}

func NewMetrics() Metrics {
	return Metrics{
		GeneratedByType:         make(map[string]int),
		CompletedByType:         make(map[string]int),
		DroppedByType:           make(map[string]int),
		TotalLatencyTicksByType: make(map[string]int),
	}
}
