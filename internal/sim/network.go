package sim

// Network stores the global state of the simulation
type Network struct {
	Nodes     map[string]*Node
	TickCount uint64
}

// Tick advances the simulation state by one step
func (this *Network) Tick() {
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
				node.DroppedCount++
				continue
			}

			// 2. ROUTE: Send the response to the next destination
			for _, outPkt := range outboundPackets {
				target := node.Router.Route(node, outPkt)
				if target != nil && len(target.Queue) < target.MaxRam {
					target.Queue = append(target.Queue, outPkt)
				} else {
					node.DroppedCount++
				}
			}
			node.ProcessedCount++
			packetProcessedThisTick++
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
