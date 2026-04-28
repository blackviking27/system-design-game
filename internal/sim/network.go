package sim

// Network stores the global state of the simulation
type Network struct {
	Nodes     map[string]*Node
	TickCount uint64
}

// Tick advances the simulation state by one step
func (this *Network) Tick() {
	// Step 1: Servers process their existing queues to free up RAM
	for _, node := range this.Nodes {
		if node.Type == TypeServer {
			packetsProcessedThisTick := 0
			// Process up the node's CPU limit
			for len(node.Queue) > 0 && packetsProcessedThisTick < node.ProcessPower {
				packet := node.Queue[0]
				node.Queue = node.Queue[1:]

				packet.Status = StatusProcessed
				node.ProcessedCount += 1
				packetsProcessedThisTick += 1
			}
		}
	}

	// Step 2: Load balancers distribute their incoming traffic
	for _, node := range this.Nodes {
		if node.Type == TypeLoadBalancer {
			// A load balancer attempts to route everything in its queue
			for len(node.Queue) > 0 {
				packet := node.Queue[0]
				node.Queue = node.Queue[1:]

				// No connection established wit the load balancer
				if len(node.Outbound) == 0 {
					packet.Status = StatusDropped
					node.DroppedCount += 1
					continue
				}

				// Execute round robin algorithm
				target := node.Outbound[node.roundRobinIdx]
				node.roundRobinIdx = (node.roundRobinIdx + 1) % len(node.Outbound)

				// Enforce bottleneck: Does the target have enough RAM to process the current packet
				if len(target.Queue) >= target.MaxRam {
					packet.Status = StatusDropped
					target.DroppedCount += 1
				} else {
					target.Queue = append(target.Queue, packet)
				}
			}
		}
	}

	this.TickCount += 1
}
