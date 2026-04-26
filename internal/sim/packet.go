package sim

type PacketStatus string

const (
	StatusPending   PacketStatus = "Pending"
	StatusProcessed PacketStatus = "Processed"
	StatusDropped   PacketStatus = "Dropped"
)

type Packet struct {
	ID     string
	Status PacketStatus
}
