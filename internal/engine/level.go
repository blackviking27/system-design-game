package engine

import (
	"encoding/json"
	"os"

	"github.com/blackviking27/system-design-game/internal/sim"
)

type PacketMixEntry struct {
	Type   string `json:"type"`
	Weight int    `json:"weight"`
}

type TrafficEvent struct {
	StartTick int              `json:"start_tick"`
	Rate      int              `json:"rate"`
	PacketMix []PacketMixEntry `json:"packet_mix"`
}

type NodeTemplate struct {
	Type      sim.NodeType                  `json:"type"`
	Workflows map[string][]sim.WorkflowStep `json:"workflows"`
}

// GameState represents the current state of the game loop
type Level struct {
	ID                 string                  `json:"id"`
	Name               string                  `json:"name"`
	StartingBudget     int                     `json:"starting_budget"`
	TargetUptimeTicks  int                     `json:"target_uptime_ticks"`
	MaxDroppedPackets  int                     `json:"max_dropped_packets"`
	UnlockedComponents []string                `json:"unlocked_components"`
	TrafficPattern     []TrafficEvent          `json:"traffic_pattern"`
	NodeTemplates      map[string]NodeTemplate `json:"node_templates"`
}

func LoadLevel(filepath string) (*Level, error) {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var data Level
	err = json.Unmarshal(bytes, &data)
	return &data, err
}
