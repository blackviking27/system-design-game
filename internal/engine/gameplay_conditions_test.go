package engine

import (
	"testing"

	"github.com/blackviking27/system-design-game/internal/sim"
	"github.com/blackviking27/system-design-game/internal/types"
)

func testSceneWithMetrics(metrics sim.Metrics) *GameplayScene {
	return &GameplayScene{
		Network: &sim.Network{
			Nodes:     make(map[string]*sim.Node),
			TickCount: 100,
			Metrics:   metrics,
		},
		Level: &Level{},
		State: types.StateSimulating,
	}
}

func TestEvaluateConditionRejectsUnknownMetric(t *testing.T) {
	scene := testSceneWithMetrics(sim.NewMetrics())

	if scene.evaluateCondition(Condition{Metric: "typo_metric", Operator: "<=", Value: 10}) {
		t.Fatal("expected unknown metric condition to evaluate false")
	}
}

func TestEvaluateConditionRejectsUnknownOperator(t *testing.T) {
	scene := testSceneWithMetrics(sim.NewMetrics())

	if scene.evaluateCondition(Condition{Metric: "uptime_ticks", Operator: "around", Value: 100}) {
		t.Fatal("expected unknown operator condition to evaluate false")
	}
}

func TestEvaluateConditionHonorsAfterTick(t *testing.T) {
	scene := testSceneWithMetrics(sim.NewMetrics())

	condition := Condition{Metric: "uptime_ticks", Operator: ">=", Value: 10, AfterTick: 200}
	if scene.evaluateCondition(condition) {
		t.Fatal("expected after_tick to prevent early condition evaluation")
	}
}

func TestMetricValueSupportsPacketTypeLatency(t *testing.T) {
	metrics := sim.NewMetrics()
	metrics.CompletedPackets = 3
	metrics.CompletedByType["READ"] = 2
	metrics.TotalLatencyTicks = 18
	metrics.TotalLatencyTicksByType["READ"] = 10
	scene := testSceneWithMetrics(metrics)

	actual, ok := scene.metricValue(Condition{Metric: "avg_latency_ticks", PacketType: "READ"})
	if !ok {
		t.Fatal("expected avg_latency_ticks metric to be valid")
	}
	if actual != 5 {
		t.Fatalf("expected READ avg latency 5, got %v", actual)
	}
}

func TestCheckWinOrLoseConditionUsesConfiguredConditions(t *testing.T) {
	metrics := sim.NewMetrics()
	metrics.GeneratedPackets = 10
	metrics.CompletedPackets = 9
	scene := testSceneWithMetrics(metrics)
	scene.Level = &Level{
		WinConditions: []Condition{
			{Metric: "success_rate", Operator: ">=", Value: 0.9},
		},
	}

	scene.checkWinOrLoseCondition()

	if scene.State != types.StateVictory {
		t.Fatalf("expected victory, got %v", scene.State)
	}
}

func TestCheckWinOrLoseConditionAnyLossWins(t *testing.T) {
	metrics := sim.NewMetrics()
	metrics.GeneratedPackets = 10
	metrics.DroppedPackets = 6
	scene := testSceneWithMetrics(metrics)
	scene.Level = &Level{
		LossConditions: []Condition{
			{Metric: "drop_rate", Operator: ">", Value: 0.5},
		},
		WinConditions: []Condition{
			{Metric: "uptime_ticks", Operator: ">=", Value: 1},
		},
	}

	scene.checkWinOrLoseCondition()

	if scene.State != types.StateGameOver {
		t.Fatalf("expected game over, got %v", scene.State)
	}
}

func TestCheckWinOrLoseConditionFallsBackToLegacy(t *testing.T) {
	scene := testSceneWithMetrics(sim.NewMetrics())
	scene.Network.TickCount = 10
	scene.Level = &Level{TargetUptimeTicks: 10, MaxDroppedPackets: 99}

	scene.checkWinOrLoseCondition()

	if scene.State != types.StateVictory {
		t.Fatalf("expected legacy victory, got %v", scene.State)
	}
}
