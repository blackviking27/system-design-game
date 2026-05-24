package test

import (
	"strconv"
	"testing"

	"github.com/blackviking27/system-design-game/internal/sim"
)

func TestBottleNeck(t *testing.T) {
	net := &sim.Network{Nodes: make(map[string]*sim.Node), Metrics: sim.NewMetrics()}

	// Create a LB with enough processing power to overwhelm one weak server.
	lb := sim.NewNode("lb-1", sim.TypeLoadBalancer, 1000, 20, 0)
	lb.Processor = &sim.WorkflowProcessor{Workflows: map[string][]sim.WorkflowStep{}}
	lb.Router = &sim.RoundRobinRouter{}

	// Create 2 weak servers. Only server A is linked so this test is deterministic
	// even though Network.Nodes is a map.
	serverA := sim.NewNode("srv-A", sim.TypeServer, 5, 2, 0)
	serverB := sim.NewNode("srv-B", sim.TypeServer, 5, 2, 0)

	// Wiring server A to the load balancer.
	lb.LinkTo(serverA)

	// Adding the server to the network
	net.Nodes[lb.ID] = lb
	net.Nodes[serverA.ID] = serverA
	net.Nodes[serverB.ID] = serverB

	// Simulating massive traffic spike: 10 packets hit the load balancer.
	for i := range 10 {
		lb.Queue = append(lb.Queue, &sim.Packet{ID: strconv.Itoa(i), InitialType: "READ", Status: sim.StatusPending})
	}

	net.Tick()

	// LB should route 5 packets into server A and drop the remaining 5 because
	// server A's queue is full.
	if lb.DroppedCount != 5 {
		t.Errorf("Expected LB to drop 5 packets due to downstream RAM limits, got %d", lb.DroppedCount)
	}
	if len(serverA.Queue) != 5 {
		t.Errorf("Expected Server A RAM to be full at 5, got %d", len(serverA.Queue))
	}
	if net.Metrics.DroppedPackets != 5 {
		t.Errorf("Expected network to record 5 dropped packets, got %d", net.Metrics.DroppedPackets)
	}
}
