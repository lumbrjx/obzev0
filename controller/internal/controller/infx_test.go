package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "obzev0/controller/api/v1"
)

func TestQueryPrometheus_ParsesScalarResult(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"resultType": "vector",
				"result": []map[string]interface{}{
					{"value": []interface{}{1234567890.0, "0.042"}},
				},
			},
		})
	}))
	defer srv.Close()

	val, err := queryPrometheus(srv.URL, "up")
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%.3f", val) != "0.042" {
		t.Fatalf("expected 0.042, got %f", val)
	}
}

func TestQueryPrometheus_EmptyResult(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"resultType": "vector",
				"result":     []map[string]interface{}{},
			},
		})
	}))
	defer srv.Close()

	val, err := queryPrometheus(srv.URL, "nonexistent_metric")
	if err != nil {
		t.Fatal(err)
	}
	if val != 0 {
		t.Fatalf("expected 0 for empty result, got %f", val)
	}
}

func TestQueryPrometheus_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("internal error"))
	}))
	defer srv.Close()

	_, err := queryPrometheus(srv.URL, "up")
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestQueryPrometheus_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer srv.Close()

	_, err := queryPrometheus(srv.URL, "up")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestHasRollbackPolicy(t *testing.T) {
	tests := []struct {
		name     string
		policy   v1.RollbackPolicy
		expected bool
	}{
		{"empty", v1.RollbackPolicy{}, false},
		{"maxDuration only", v1.RollbackPolicy{MaxDurationSeconds: 60}, true},
		{"prometheus only", v1.RollbackPolicy{PrometheusURL: "http://prom:9090"}, true},
		{"both", v1.RollbackPolicy{MaxDurationSeconds: 60, PrometheusURL: "http://prom:9090"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasRollbackPolicy(tt.policy)
			if got != tt.expected {
				t.Errorf("hasRollbackPolicy(%+v) = %v, want %v", tt.policy, got, tt.expected)
			}
		})
	}
}

func TestEmitEvent_NilRecorderDoesNotPanic(t *testing.T) {
	old := eventRecorder
	eventRecorder = nil
	defer func() { eventRecorder = old }()

	emitEvent(nil, "Normal", "Test", "should not panic: %s", "ok")
}

func TestPodNameList(t *testing.T) {
	conns := []*PodConnection{
		{PodName: "alpha"},
		{PodName: "bravo"},
		{PodName: "charlie"},
	}
	names := podNameList(conns)
	if len(names) != 3 {
		t.Fatalf("expected 3, got %d", len(names))
	}
	if names[0] != "alpha" || names[1] != "bravo" || names[2] != "charlie" {
		t.Errorf("unexpected names: %v", names)
	}
}
