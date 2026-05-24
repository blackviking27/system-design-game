package sim

import (
	"errors"
	"testing"
)

type passthroughProcessor struct{}

func (passthroughProcessor) Process(node *Node, packet *Packet) ([]*Packet, error) {
	return []*Packet{packet}, nil
}

type errorProcessor struct{}

func (errorProcessor) Process(node *Node, packet *Packet) ([]*Packet, error) {
	return nil, errors.New("processor failed")
}

func TestRecordGeneratedInitializesMetricsMaps(t *testing.T) {
	network := &Network{}

	network.RecordGenerated("READ")

	if network.Metrics.GeneratedPackets != 1 {
		t.Fatalf("expected 1 generated packet, got %d", network.Metrics.GeneratedPackets)
	}
	if network.Metrics.GeneratedByType["READ"] != 1 {
		t.Fatalf("expected generated READ count to be 1, got %d", network.Metrics.GeneratedByType["READ"])
	}
}

func TestTickRecordsTerminalPacketCompletion(t *testing.T) {
	node := NewNode("server", TypeServer, 10, 1, 0)
	node.Processor = passthroughProcessor{}
	node.Router = &RoundRobinRouter{}
	node.Queue = append(node.Queue, &Packet{
		ID:          "pkt-1",
		Type:        "DONE",
		InitialType: "READ",
		CreatedAt:   3,
	})

	network := &Network{
		Nodes:               map[string]*Node{node.ID: node},
		TickCount:           10,
		TerminalPacketTypes: map[string]bool{"DONE": true},
	}

	network.Tick()

	if network.Metrics.CompletedPackets != 1 {
		t.Fatalf("expected 1 completed packet, got %d", network.Metrics.CompletedPackets)
	}
	if network.Metrics.CompletedByType["READ"] != 1 {
		t.Fatalf("expected completed READ count to be 1, got %d", network.Metrics.CompletedByType["READ"])
	}
	if network.Metrics.TotalLatencyTicks != 7 {
		t.Fatalf("expected total latency 7, got %d", network.Metrics.TotalLatencyTicks)
	}
	if network.Metrics.TotalLatencyTicksByType["READ"] != 7 {
		t.Fatalf("expected READ latency 7, got %d", network.Metrics.TotalLatencyTicksByType["READ"])
	}
}

func TestTickRecordsProcessorErrorDrops(t *testing.T) {
	node := NewNode("server", TypeServer, 10, 1, 0)
	node.Processor = errorProcessor{}
	node.Router = &RoundRobinRouter{}
	node.Queue = append(node.Queue, &Packet{
		ID:          "pkt-1",
		Type:        "READ",
		InitialType: "READ",
	})

	network := &Network{Nodes: map[string]*Node{node.ID: node}}

	network.Tick()

	if node.DroppedCount != 1 {
		t.Fatalf("expected node drop count 1, got %d", node.DroppedCount)
	}
	if network.Metrics.DroppedPackets != 1 {
		t.Fatalf("expected 1 dropped packet, got %d", network.Metrics.DroppedPackets)
	}
	if network.Metrics.DroppedByType["READ"] != 1 {
		t.Fatalf("expected dropped READ count to be 1, got %d", network.Metrics.DroppedByType["READ"])
	}
}

func TestResetStateClearsRuntimeState(t *testing.T) {
	node := NewNode("server", TypeServer, 10, 1, 0)
	node.Queue = append(node.Queue, &Packet{ID: "pkt-1"})
	node.WaitingRoom["trace-1"] = &Packet{ID: "pkt-2"}
	node.State["key"] = "value"
	node.ProcessedCount = 1
	node.DroppedCount = 1

	node.ResetState()

	if len(node.Queue) != 0 {
		t.Fatalf("expected queue to be empty, got %d", len(node.Queue))
	}
	if len(node.WaitingRoom) != 0 {
		t.Fatalf("expected waiting room to be empty, got %d", len(node.WaitingRoom))
	}
	if len(node.State) != 0 {
		t.Fatalf("expected state to be empty, got %d", len(node.State))
	}
	if node.ProcessedCount != 0 || node.DroppedCount != 0 {
		t.Fatalf("expected counts to reset, got processed=%d dropped=%d", node.ProcessedCount, node.DroppedCount)
	}
}
