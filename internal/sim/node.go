// Core entity of infrastrcuture
// Handles own RAM limit and CPU limit

package sim

type NodeType string

const (
	TypeLoadBalancer NodeType = "LoadBalancer"
	TypeServer       NodeType = "Server"
	TypeMessageQueue NodeType = "MessageQueue"
	TypeDatabase     NodeType = "DataBase"
	TypeCache        NodeType = "Cache"
)

type Node struct {
	ID           string
	Type         NodeType
	MaxRam       int       // Max number of packets the node can store in memory
	ProcessPower int       // Number of packets this node can process in one tick
	Queue        []*Packet // the current memory buffer
	Outbound     []*Node   // Connections to downstream nodes

	// internal state of routing
	roundRobinIdx int

	// Metrics for game hud
	ProcessedCount int
	DroppedCount   int

	// Canvas coordinates
	X, Y float64

	// Cost of the system
	Cost int
}

// Create a new node
func NewNode(id string, t NodeType, maxRam, processPower, cost int) *Node {
	return &Node{
		ID:           id,
		Type:         t,
		MaxRam:       maxRam,
		ProcessPower: processPower,
		Queue:        make([]*Packet, 0),
		Outbound:     make([]*Node, 0),
		Cost:         cost,
	}
}

// Function to add downstream node
func (this *Node) LinkTo(dest *Node) {
	this.Outbound = append(this.Outbound, dest)
}

// Node component catalogue type
type NodeTemplate struct {
	Type         NodeType
	Name         string
	Cost         int
	MaxRam       int
	ProcessPower int
}

var Catalog = []NodeTemplate{
	{Type: TypeServer, Name: "Lite server\n(2 pkts/tick)", Cost: 100, MaxRam: 10, ProcessPower: 2},
	{Type: TypeServer, Name: "Heavy server\n(5 pkts/tick)", Cost: 100, MaxRam: 50, ProcessPower: 5},
	{Type: TypeLoadBalancer, Name: "Load Balancer", Cost: 500, MaxRam: 5000, ProcessPower: 100},

	{Type: TypeMessageQueue, Name: "Message Queue\n(Buffer)", Cost: 250, MaxRam: 1000, ProcessPower: 5},
	{Type: TypeDatabase, Name: "Data Base\n(SQL)", Cost: 400, MaxRam: 5000, ProcessPower: 1},
	{Type: TypeCache, Name: "Cache\n(Redis)", Cost: 200, MaxRam: 20, ProcessPower: 15},
}
