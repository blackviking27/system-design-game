package sim

type PacketStatus string

const (
	StatusPending   PacketStatus = "Pending"
	StatusProcessed PacketStatus = "Processed"
	StatusDropped   PacketStatus = "Dropped"
)

type Packet struct {
	ID           string
	TraceId      string // Links request and responses together
	OriginNodeId string // Allows a response to finds its way home
	Status       PacketStatus
	Type         string         // Trigger for workflow steps
	Payload      map[string]any // Simulated data
	Metadata     map[string]int // Latency, hop count etc
	RouteTarget  string
	InitialType  string
	CreatedAt    uint64
}
