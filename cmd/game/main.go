package main

import (
	"github.com/blackviking27/system-design-game/internal/engine"
	"github.com/blackviking27/system-design-game/internal/sim"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	// Initializing the sim network
	network := &sim.Network{Nodes: make(map[string]*sim.Node)}

	// Create a load balancer
	lb := sim.NewNode("LB-Main", sim.TypeLoadBalancer, 100, 100)
	lb.X, lb.Y = 400, 150

	// Create 2 weak servers
	serverA := sim.NewNode("Server-A", sim.TypeServer, 10, 1)
	serverA.X, serverA.Y = 200, 450

	serverB := sim.NewNode("Server-B", sim.TypeServer, 10, 1)
	serverB.X, serverB.Y = 600, 450

	// Creating connections
	lb.LinkTo(serverA)
	lb.LinkTo(serverB)

	network.Nodes[lb.ID] = lb
	network.Nodes[serverA.ID] = serverA
	network.Nodes[serverB.ID] = serverB

	// Defining the level
	levelOne := engine.Level{
		Name:              "Level 1",
		TargetUptimeTicks: 500,
		MaxDroppedPackets: 50,
		BaseTrafficRate:   3,
	}

	// Creating the engine game wrapper
	game := &engine.Game{
		Network: network,
		Level:   levelOne,
		State:   engine.StatePlaying,
	}

	// Configure window and run
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("System design sandbox")

	// Running the game
	err := ebiten.RunGame(game)
	if err != nil {
		panic(err)
	}
}
