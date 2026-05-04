package sim

import "math/rand"

// Network stores the global state of the simulation
type Network struct {
	Nodes     map[string]*Node
	TickCount uint64
}

// Tick advances the simulation state by one step
func (this *Network) Tick() {
	for _, node := range this.Nodes {
		switch node.Type {

		// 1. Standard processiong (Servers and Databases)
		case TypeServer:
		case TypeDatabase:
			packetsProcessedThisTick := 0
			for len(node.Queue) > 0 && packetsProcessedThisTick < node.ProcessPower {
				packet := node.Queue[0]
				node.Queue = node.Queue[1:]
				packet.Status = StatusProcessed
				node.ProcessedCount++
				packetsProcessedThisTick++
			}

		// 2. Cache processing (Probability case)
		case TypeCache:
			packetsProcessedThisTick := 0
			if len(node.Queue) > 0 && packetsProcessedThisTick < node.ProcessPower {
				packet := node.Queue[0]
				node.Queue = node.Queue[1:]

				if rand.Float32() <= 0.8 {
					// CACHE HIT (80% Hit chance): Packet is Processed instantly
					packet.Status = StatusProcessed
					packetsProcessedThisTick++
				} else {
					// CACHE MISS (20% Chance): Packet is forwarded to DB or dropped
					if len(node.Outbound) > 0 {
						target := node.Outbound[0]
						if len(target.Queue) >= target.MaxRam {
							packet.Status = StatusDropped
							target.DroppedCount++
						} else {
							target.Queue = append(target.Queue, packet)
						}
					} else {
						packet.Status = StatusDropped
						node.DroppedCount++
					}
				}
			}
			packetsProcessedThisTick++

		// 3. INSTANT ROUTING (Loadbalancer)
		case TypeLoadBalancer:
			for len(node.Queue) > 0 {
				packet := node.Queue[0]
				node.Queue = node.Queue[1:]

				if len(node.Outbound) == 0 {
					packet.Status = StatusDropped
					node.DroppedCount++
					continue
				}

				target := node.Outbound[node.roundRobinIdx]
				node.roundRobinIdx = (node.roundRobinIdx + 1) % len(node.Outbound)

				if len(target.Queue) >= target.MaxRam {
					packet.Status = StatusDropped
					target.DroppedCount++
				} else {
					target.Queue = append(target.Queue, packet)
				}
			}

		// 4. THROTTLED ROUTING (Message Queue)
		case TypeMessageQueue:
			routedThisTick := 0
			for len(node.Queue) > 0 && routedThisTick < node.ProcessPower {
				packet := node.Queue[0]
				node.Queue = node.Queue[1:]

				if len(node.Outbound) == 0 {
					packet.Status = StatusDropped
					node.DroppedCount++
					routedThisTick++
					continue
				}

				target := node.Outbound[node.roundRobinIdx]
				node.roundRobinIdx = (node.roundRobinIdx + 1) % len(node.Outbound)

				if len(target.Queue) >= target.MaxRam {
					packet.Status = StatusDropped
					target.DroppedCount++
				} else {
					target.Queue = append(target.Queue, packet)
				}

				routedThisTick++
			}
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
