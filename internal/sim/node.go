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

	WaitingRoom   map[string]*Packet // Key: TraceId
	MaxConnection int                // To simulate thread exhaustion
	State         map[string]any     // Presistent key / value storage

	// Pluggable behavior
	Processor Processor
	Router    Router
}

// Create a new node
func NewNode(id string, t NodeType, maxRam, processPower, cost int) *Node {
	return &Node{
		ID:            id,
		Type:          t,
		MaxRam:        maxRam,
		ProcessPower:  processPower,
		Queue:         make([]*Packet, 0),
		WaitingRoom:   make(map[string]*Packet, 0),
		MaxConnection: maxRam, // Default Limit
		State:         make(map[string]interface{}, 0),
		Outbound:      make([]*Node, 0),
		Cost:          cost,
	}
}

// Function to add downstream node
func (this *Node) LinkTo(dest *Node) {
	this.Outbound = append(this.Outbound, dest)
}

// Reset node state
func (this *Node) ResetState() {
	this.Queue = make([]*Packet, 0)
	this.ProcessedCount = 0
	this.DroppedCount = 0
	this.roundRobinIdx = 0
}

// Catalogue node type
type CatalogNodeTemplate struct {
	ID           string
	Type         NodeType
	Name         string
	Cost         int
	MaxRam       int
	ProcessPower int
}

var Catalog = []CatalogNodeTemplate{
	{ID: "lite_server", Type: TypeServer, Name: "Lite server\n(4 pkts/tick)", Cost: 100, MaxRam: 30, ProcessPower: 4},
	{ID: "heavy_server", Type: TypeServer, Name: "Heavy server\n(10 pkts/tick)", Cost: 250, MaxRam: 100, ProcessPower: 10},
	{ID: "load_blancer", Type: TypeLoadBalancer, Name: "Load Balancer\n(40 pkts/tick)", Cost: 300, MaxRam: 200, ProcessPower: 40},
	{ID: "message_queue", Type: TypeMessageQueue, Name: "Message Queue\n(Buffer)", Cost: 300, MaxRam: 800, ProcessPower: 12},
	{ID: "data_base", Type: TypeDatabase, Name: "Data Base\n(SQL)", Cost: 500, MaxRam: 300, ProcessPower: 3},
	{ID: "cache", Type: TypeCache, Name: "Cache\n(Redis)", Cost: 250, MaxRam: 150, ProcessPower: 25},
}

func GetCatalogTemplateForType(t NodeType) (CatalogNodeTemplate, bool) {
	for _, template := range Catalog {
		if template.Type == t {
			return template, true
		}
	}
	return CatalogNodeTemplate{}, false
}
