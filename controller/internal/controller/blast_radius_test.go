package controller

import (
	"testing"

	v1 "obzev0/controller/api/v1"

	"google.golang.org/grpc"
)

func makeConn(podName, namespace string, labels map[string]string) *PodConnection {
	return &PodConnection{
		Conn:      &grpc.ClientConn{},
		PodName:   podName,
		Namespace: namespace,
		Labels:    labels,
	}
}

func buildRegistry(pcs ...*PodConnection) map[string]*PodConnection {
	m := make(map[string]*PodConnection, len(pcs))
	for i, pc := range pcs {
		m[pc.PodName+":"+string(rune('0'+i))] = pc
	}
	return m
}

func TestFilterConnections_NoFilter(t *testing.T) {
	reg := buildRegistry(
		makeConn("pod-a", "default", nil),
		makeConn("pod-b", "staging", nil),
	)
	result := filterConnections(reg, v1.BlastRadius{})
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestFilterConnections_ByNamespace(t *testing.T) {
	reg := buildRegistry(
		makeConn("pod-a", "default", nil),
		makeConn("pod-b", "staging", nil),
		makeConn("pod-c", "staging", nil),
	)
	result := filterConnections(reg, v1.BlastRadius{Namespaces: []string{"staging"}})
	if len(result) != 2 {
		t.Fatalf("expected 2 staging pods, got %d", len(result))
	}
	for _, pc := range result {
		if pc.Namespace != "staging" {
			t.Errorf("unexpected namespace %q", pc.Namespace)
		}
	}
}

func TestFilterConnections_ByLabel(t *testing.T) {
	reg := buildRegistry(
		makeConn("pod-a", "default", map[string]string{"tier": "frontend"}),
		makeConn("pod-b", "default", map[string]string{"tier": "backend"}),
		makeConn("pod-c", "default", map[string]string{"tier": "backend"}),
	)
	result := filterConnections(reg, v1.BlastRadius{LabelSelector: map[string]string{"tier": "backend"}})
	if len(result) != 2 {
		t.Fatalf("expected 2 backend pods, got %d", len(result))
	}
}

func TestFilterConnections_Percentage(t *testing.T) {
	reg := buildRegistry(
		makeConn("pod-a", "default", nil),
		makeConn("pod-b", "default", nil),
		makeConn("pod-c", "default", nil),
		makeConn("pod-d", "default", nil),
	)
	result := filterConnections(reg, v1.BlastRadius{Percentage: 50})
	if len(result) != 2 {
		t.Fatalf("expected 2 pods (50%% of 4), got %d", len(result))
	}
}

func TestFilterConnections_PercentageRoundsUpToAtLeastOne(t *testing.T) {
	reg := buildRegistry(makeConn("pod-a", "default", nil))
	result := filterConnections(reg, v1.BlastRadius{Percentage: 10})
	if len(result) != 1 {
		t.Fatalf("expected at least 1, got %d", len(result))
	}
}

func TestFilterConnections_NamespaceAndLabel(t *testing.T) {
	reg := buildRegistry(
		makeConn("pod-a", "staging", map[string]string{"tier": "backend"}),
		makeConn("pod-b", "staging", map[string]string{"tier": "frontend"}),
		makeConn("pod-c", "default", map[string]string{"tier": "backend"}),
	)
	result := filterConnections(reg, v1.BlastRadius{
		Namespaces:    []string{"staging"},
		LabelSelector: map[string]string{"tier": "backend"},
	})
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].PodName != "pod-a" {
		t.Errorf("expected pod-a, got %s", result[0].PodName)
	}
}
