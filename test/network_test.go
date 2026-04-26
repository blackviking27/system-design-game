package test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/blackviking27/system-design-game/internal/sim"
)

func TestBottleNeck(t *testing.T) {
	net := &sim.Network{Nodes: make(map[string]*sim.Node)}

	// Create a LB with infinite capacity
	lb := sim.NewNode("lb-1", sim.TypeLoadBalancer, 1000, 1000)

	// Create 2 weak servers (can only hold 5 packets in RAM and process 2 at a time)
	serverA := sim.NewNode("srv-A", sim.TypeServer, 5, 2)
	serverB := sim.NewNode("srv-B", sim.TypeServer, 5, 2)

	// Wiring the servers to the load balancer
	lb.LinkTo(serverA)
	lb.LinkTo(serverB)

	// Adding the server to the network
	net.Nodes[lb.ID] = lb
	net.Nodes[serverA.ID] = serverA
	net.Nodes[serverB.ID] = serverB

	// Simulating massive traffic spike: 20 packets hit the load balancer
	for i := range 20 {
		lb.Queue = append(lb.Queue, &sim.Packet{ID: strconv.Itoa(i), Status: sim.StatusPending})
	}

	// TICK 1
	net.Tick()

	// LB should route 10 to A, and 10 to B.
	// Since A and B only have 5 RAM capacity each, they should both drop 5 packets.

	fmt.Printf("Server A: %v", serverA)
	fmt.Printf("Server B: %v", serverB)

	if serverA.DroppedCount != 5 {
		t.Errorf("Expected Server A to drop 5 packets due to RAM limits, got %d", serverA.DroppedCount)
	}
	if len(serverA.Queue) != 5 {
		t.Errorf("Expected Server A RAM to be full at 5, got %d", len(serverA.Queue))
	}

	// TICK 2

	net.Tick()

	fmt.Printf("Server A: %v", serverA)
	fmt.Printf("Server B: %v", serverB)

	// Servers should now process packets from their RAM, freeing up space.
	// They can process 2 per tick.
	if serverA.ProcessedCount != 2 {
		t.Errorf("Expected Server A to have processed 2 packets, got %d", serverA.ProcessedCount)
	}
	if len(serverA.Queue) != 3 {
		t.Errorf("Expected Server A RAM queue to reduce to 3, got %d", len(serverA.Queue))
	}

}
