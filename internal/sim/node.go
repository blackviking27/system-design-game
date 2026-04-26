// Core entity of infrastrcuture
// Handles own RAM limit and CPU limit

package sim

type NodeType string

const (
	TypeLoadBalancer NodeType = "LoadBalancer"
	TypeServer       NodeType = "Server"
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
}

// Create a new node
func NewNode(id string, t NodeType, maxRam, processPower int) *Node {
	return &Node{
		ID:           id,
		Type:         t,
		MaxRam:       maxRam,
		ProcessPower: processPower,
		Queue:        make([]*Packet, 0),
		Outbound:     make([]*Node, 0),
	}
}

// Function to add downstream node
func (this *Node) LinkTo(dest *Node) {
	this.Outbound = append(this.Outbound, dest)
}
