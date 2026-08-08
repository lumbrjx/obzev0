package controller

import (
	"encoding/json"
	"testing"
	"time"
)

func TestExperimentReport_JSONMarshal(t *testing.T) {
	start := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	end := start.Add(2 * time.Minute)

	report := ExperimentReport{
		ExperimentName:  "test-latency",
		Namespace:       "staging",
		StartTime:       start,
		EndTime:         end,
		DurationSeconds: 120.0,
		AffectedPods:    []string{"pod-a", "pod-b"},
		ChaosTypes:      []string{"latency", "packetDrop"},
		StopReason:      "maxDuration",
		StepCount:       3,
	}

	data, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatal(err)
	}

	if parsed["experimentName"] != "test-latency" {
		t.Errorf("expected experimentName=test-latency, got %v", parsed["experimentName"])
	}
	if parsed["namespace"] != "staging" {
		t.Errorf("expected namespace=staging, got %v", parsed["namespace"])
	}
	if parsed["stopReason"] != "maxDuration" {
		t.Errorf("expected stopReason=maxDuration, got %v", parsed["stopReason"])
	}
	if parsed["durationSeconds"].(float64) != 120.0 {
		t.Errorf("expected durationSeconds=120, got %v", parsed["durationSeconds"])
	}

	pods := parsed["affectedPods"].([]interface{})
	if len(pods) != 2 {
		t.Errorf("expected 2 pods, got %d", len(pods))
	}
	if pods[0].(string) != "pod-a" || pods[1].(string) != "pod-b" {
		t.Errorf("unexpected pod names: %v", pods)
	}

	steps := parsed["stepCount"].(float64)
	if steps != 3 {
		t.Errorf("expected stepCount=3, got %v", steps)
	}
}

func TestExperimentReport_EmptyOptionalFields(t *testing.T) {
	report := ExperimentReport{
		ExperimentName:  "simple",
		Namespace:       "default",
		StartTime:       time.Now(),
		EndTime:         time.Now(),
		DurationSeconds: 0,
		AffectedPods:    []string{},
		StopReason:      "manual",
	}

	data, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatal(err)
	}

	if _, exists := parsed["chaosTypes"]; exists {
		t.Error("chaosTypes should be omitted when empty")
	}
	if _, exists := parsed["metricsBefore"]; exists {
		t.Error("metricsBefore should be omitted when nil")
	}
	if _, exists := parsed["metricsAfter"]; exists {
		t.Error("metricsAfter should be omitted when nil")
	}
}

func TestExperimentReport_RoundTrip(t *testing.T) {
	original := ExperimentReport{
		ExperimentName:  "roundtrip",
		Namespace:       "production",
		StartTime:       time.Date(2025, 6, 15, 10, 0, 0, 0, time.UTC),
		EndTime:         time.Date(2025, 6, 15, 10, 5, 0, 0, time.UTC),
		DurationSeconds: 300,
		AffectedPods:    []string{"pod-x", "pod-y", "pod-z"},
		ChaosTypes:      []string{"httpFault"},
		StopReason:      "metricThreshold",
		StepCount:       1,
		MetricsBefore:   map[string]string{"error_rate": "0.01"},
		MetricsAfter:    map[string]string{"error_rate": "0.15"},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}

	var decoded ExperimentReport
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}

	if decoded.ExperimentName != original.ExperimentName {
		t.Errorf("name mismatch: %s vs %s", decoded.ExperimentName, original.ExperimentName)
	}
	if decoded.StopReason != original.StopReason {
		t.Errorf("stopReason mismatch: %s vs %s", decoded.StopReason, original.StopReason)
	}
	if len(decoded.AffectedPods) != len(original.AffectedPods) {
		t.Errorf("affectedPods length mismatch: %d vs %d", len(decoded.AffectedPods), len(original.AffectedPods))
	}
	if decoded.MetricsBefore["error_rate"] != "0.01" {
		t.Errorf("metricsBefore mismatch: %v", decoded.MetricsBefore)
	}
	if decoded.MetricsAfter["error_rate"] != "0.15" {
		t.Errorf("metricsAfter mismatch: %v", decoded.MetricsAfter)
	}
}
