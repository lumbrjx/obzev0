package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	v1 "obzev0/controller/api/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ExperimentReport is the structured summary written after every experiment ends.
type ExperimentReport struct {
	ExperimentName  string            `json:"experimentName"`
	Namespace       string            `json:"namespace"`
	StartTime       time.Time         `json:"startTime"`
	EndTime         time.Time         `json:"endTime"`
	DurationSeconds float64           `json:"durationSeconds"`
	AffectedPods    []string          `json:"affectedPods"`
	ChaosTypes      []string          `json:"chaosTypes,omitempty"`
	StopReason      string            `json:"stopReason"`
	StepCount       int               `json:"stepCount,omitempty"`
	MetricsBefore   map[string]string `json:"metricsBefore,omitempty"`
	MetricsAfter    map[string]string `json:"metricsAfter,omitempty"`
}

// publishReport logs the report as structured JSON and writes it as a Kubernetes ConfigMap.
func publishReport(cr *v1.Obzev0Resource, report ExperimentReport) {
	data, err := json.Marshal(report)
	if err != nil {
		klog.Errorf("Failed to marshal experiment report: %v", err)
		return
	}

	klog.Infof("ExperimentReport: %s", string(data))

	if ctrlClient == nil {
		return
	}

	ts := time.Now().Unix()
	cmName := fmt.Sprintf("obzev0-report-%s-%d", report.ExperimentName, ts)
	if len(cmName) > 253 {
		cmName = cmName[:253]
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cmName,
			Namespace: report.Namespace,
			Labels: map[string]string{
				"obzev0.io/report":     "true",
				"obzev0.io/experiment": report.ExperimentName,
			},
		},
		Data: map[string]string{
			"report.json": string(data),
		},
	}

	if err := ctrlClient.Create(context.Background(), cm, &client.CreateOptions{}); err != nil {
		klog.Errorf("Failed to create report ConfigMap %s: %v", cmName, err)
		return
	}

	klog.Infof("Report ConfigMap created: %s/%s", report.Namespace, cmName)
}
