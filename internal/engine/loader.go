package engine

import (
	"encoding/json"
	"os"
)

type TrafficEvent struct {
	StartTick int `json:"start_tick"`
	Rate      int `json:"rate"`
}

type Level struct {
	ID                 string         `json:"id"`
	Name               string         `json:"name"`
	StartingBudget     int            `json:"starting_budget"`
	TargetUptimeTicks  int            `json:"target_uptime_ticks"`
	MaxDroppedPackets  int            `json:"max_dropped_packets"`
	UnlockedComponents []string       `json:"unlocked_components"`
	TrafficPattern     []TrafficEvent `json:"traffic_pattern"`
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
