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

type Condition struct {
	Metric     string  `json:"metric"`
	Operator   string  `json:"operator"`
	Value      float64 `json:"value"`
	PacketType string  `json:"packet_type,omitempty"`
	AfterTick  int     `json:"after_tick,omitempty"`
}

// GameState represents the current state of the game loop
type Level struct {
	ID                  string                  `json:"id"`
	Name                string                  `json:"name"`
	StartingBudget      int                     `json:"starting_budget"`
	TargetUptimeTicks   int                     `json:"target_uptime_ticks"`
	MaxDroppedPackets   int                     `json:"max_dropped_packets"`
	UnlockedComponents  []string                `json:"unlocked_components"`
	TrafficPattern      []TrafficEvent          `json:"traffic_pattern"`
	NodeTemplates       map[string]NodeTemplate `json:"node_templates"`
	TerminalPacketTypes []string                `json:"terminal_packet_types"`
	WinConditions       []Condition             `json:"win_conditions"`
	LossConditions      []Condition             `json:"loss_conditions"`
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
