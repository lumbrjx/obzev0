package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"sync"
	"time"

	v1 "obzev0/controller/api/v1"
	latencyproto "obzev0/common/proto/latency"
	pcaproto "obzev0/common/proto/packetManipulation"
	tcaproto "obzev0/common/proto/tcAnalyser"

	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

// schedulers tracks active cron jobs keyed by "namespace/name".
var schedulers sync.Map

// ── K8s Event helpers ─────────────────────────────────────────────────────────

func emitEvent(cr *v1.Obzev0Resource, eventType, reason, msgFmt string, args ...interface{}) {
	if eventRecorder == nil {
		return
	}
	msg := fmt.Sprintf(msgFmt, args...)
	eventRecorder.Eventf(cr, eventType, reason, "%s", msg)
}

// ── Blast Radius ─────────────────────────────────────────────────────────────

// filterConnections returns the subset of connections that match the BlastRadius spec.
func filterConnections(all map[string]*PodConnection, br v1.BlastRadius) []*PodConnection {
	var matched []*PodConnection

	for _, pc := range all {
		if !matchesNamespaces(pc, br.Namespaces) {
			continue
		}
		if !matchesLabels(pc, br.LabelSelector) {
			continue
		}
		matched = append(matched, pc)
	}

	if br.Percentage > 0 && br.Percentage < 100 && len(matched) > 0 {
		// Deterministic subset: sort by pod name so the same pods are always chosen.
		sort.Slice(matched, func(i, j int) bool {
			return matched[i].PodName < matched[j].PodName
		})
		n := int(float64(len(matched)) * float64(br.Percentage) / 100.0)
		if n < 1 {
			n = 1
		}
		matched = matched[:n]
	}

	return matched
}

func matchesNamespaces(pc *PodConnection, namespaces []string) bool {
	if len(namespaces) == 0 {
		return true
	}
	for _, ns := range namespaces {
		if pc.Namespace == ns {
			return true
		}
	}
	return false
}

func matchesLabels(pc *PodConnection, selector map[string]string) bool {
	if len(selector) == 0 {
		return true
	}
	for k, v := range selector {
		if pc.Labels[k] != v {
			return false
		}
	}
	return true
}

func podNameList(conns []*PodConnection) []string {
	names := make([]string, len(conns))
	for i, pc := range conns {
		names[i] = pc.PodName
	}
	return names
}

// ── Experiment dispatch ───────────────────────────────────────────────────────

func processCustomResource(cr *v1.Obzev0Resource, targets []*PodConnection) {
	name := cr.GetName()
	namespace := cr.GetNamespace()
	klog.Infof("Processing CR: %s/%s against %d pod(s)", namespace, name, len(targets))
	emitEvent(cr, corev1.EventTypeNormal, "ExperimentStarted",
		"Chaos experiment started on %d pod(s)", len(targets))

	if cr.Spec.Schedule != "" {
		scheduleExperiment(cr, targets)
		return
	}

	if len(cr.Spec.Steps) > 0 {
		executeWorkflow(cr, targets, cr.Spec.RollbackPolicy)
		return
	}

	// Flat (legacy) single-shot mode.
	cfg := GrpcServiceConfig{
		LatencyConfig: cr.Spec.LatencyServiceConfig,
		TcAConfig:     cr.Spec.TcAnalyserServiceConfig,
		PctmConfig:    cr.Spec.PacketManipulationServiceConfig,
	}
	startTime := time.Now()
	for _, pc := range targets {
		if err := callGrpcServices(pc.Conn, cfg); err != nil {
			log.Printf("Error calling gRPC services on %s: %v", pc.PodName, err)
		}
	}

	if hasRollbackPolicy(cr.Spec.RollbackPolicy) {
		go monitorAndRollback(cr, targets, cr.Spec.RollbackPolicy, cfg, startTime)
	}

	klog.Infof("Experiment dispatched for CR: %s/%s", namespace, name)
}

// ── Workflow execution ────────────────────────────────────────────────────────

func executeWorkflow(cr *v1.Obzev0Resource, conns []*PodConnection, policy v1.RollbackPolicy) {
	steps := cr.Spec.Steps
	startTime := time.Now()

	for i, step := range steps {
		klog.Infof("Workflow step %d/%d: %s", i+1, len(steps), step.Name)
		emitEvent(cr, corev1.EventTypeNormal, "WorkflowStep",
			"Step %d/%d: %s", i+1, len(steps), step.Name)

		cfg := GrpcServiceConfig{
			LatencyConfig: step.LatencyServiceConfig,
			TcAConfig:     step.TcAnalyserServiceConfig,
			PctmConfig:    step.PacketManipulationServiceConfig,
		}

		for _, pc := range conns {
			if err := callGrpcServices(pc.Conn, cfg); err != nil {
				log.Printf("Step %q error on %s: %v", step.Name, pc.PodName, err)
			}
		}

		if step.DurationSeconds > 0 {
			d := time.Duration(step.DurationSeconds) * time.Second
			klog.Infof("Step %q running for %s", step.Name, d)
			time.Sleep(d)
			stopExperiments(cr, conns, cfg, "workflowStep")
			klog.Infof("Step %q stopped", step.Name)
		}
	}

	publishReport(cr, ExperimentReport{
		ExperimentName:  cr.Name,
		Namespace:       cr.Namespace,
		StartTime:       startTime,
		EndTime:         time.Now(),
		DurationSeconds: time.Since(startTime).Seconds(),
		AffectedPods:    podNameList(conns),
		StopReason:      "workflowComplete",
		StepCount:       len(steps),
	})

	emitEvent(cr, corev1.EventTypeNormal, "WorkflowComplete",
		"Workflow completed %d step(s) in %.1fs", len(steps), time.Since(startTime).Seconds())
}

// ── Scheduling ────────────────────────────────────────────────────────────────

func scheduleExperiment(cr *v1.Obzev0Resource, conns []*PodConnection) {
	key := cr.Namespace + "/" + cr.Name
	stopScheduler(key)

	c := cron.New()
	_, err := c.AddFunc(cr.Spec.Schedule, func() {
		klog.Infof("Scheduled experiment firing: %s", key)
		if len(cr.Spec.Steps) > 0 {
			executeWorkflow(cr, conns, cr.Spec.RollbackPolicy)
		} else {
			cfg := GrpcServiceConfig{
				LatencyConfig: cr.Spec.LatencyServiceConfig,
				TcAConfig:     cr.Spec.TcAnalyserServiceConfig,
				PctmConfig:    cr.Spec.PacketManipulationServiceConfig,
			}
			startTime := time.Now()
			for _, pc := range conns {
				if err := callGrpcServices(pc.Conn, cfg); err != nil {
					log.Printf("Scheduled gRPC error on %s: %v", pc.PodName, err)
				}
			}
			if hasRollbackPolicy(cr.Spec.RollbackPolicy) {
				go monitorAndRollback(cr, conns, cr.Spec.RollbackPolicy, cfg, startTime)
			}
		}
	})
	if err != nil {
		log.Printf("Invalid cron schedule %q for CR %s: %v", cr.Spec.Schedule, key, err)
		return
	}

	c.Start()
	schedulers.Store(key, c)
	klog.Infof("Scheduled experiment %q with schedule: %s", key, cr.Spec.Schedule)
}

func stopScheduler(key string) {
	if v, ok := schedulers.LoadAndDelete(key); ok {
		v.(*cron.Cron).Stop()
		klog.Infof("Stopped scheduler for %s", key)
	}
}

// ── Auto-rollback ─────────────────────────────────────────────────────────────

func hasRollbackPolicy(p v1.RollbackPolicy) bool {
	return p.MaxDurationSeconds > 0 || p.PrometheusURL != ""
}

func monitorAndRollback(cr *v1.Obzev0Resource, conns []*PodConnection, policy v1.RollbackPolicy, cfg GrpcServiceConfig, startTime time.Time) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stopReason := "unknown"

	if policy.MaxDurationSeconds > 0 {
		go func() {
			time.Sleep(time.Duration(policy.MaxDurationSeconds) * time.Second)
			klog.Info("Max experiment duration reached; triggering rollback")
			emitEvent(cr, corev1.EventTypeWarning, "RollbackTriggered",
				"Max duration of %ds reached", policy.MaxDurationSeconds)
			stopReason = "maxDuration"
			cancel()
		}()
	}

	if policy.PrometheusURL != "" && policy.MetricQuery != "" {
		interval := time.Duration(policy.PollIntervalSeconds) * time.Second
		if policy.PollIntervalSeconds <= 0 {
			interval = 10 * time.Second
		}
		go func() {
			pollPrometheus(ctx, policy, func(val float64) {
				klog.Infof("Metric %.4f exceeded threshold %.4f; triggering rollback", val, policy.Threshold)
				emitEvent(cr, corev1.EventTypeWarning, "RollbackTriggered",
					"Metric value %.4f exceeded threshold %.4f", val, policy.Threshold)
				stopReason = "metricThreshold"
				cancel()
			}, interval)
		}()
	}

	<-ctx.Done()
	stopExperiments(cr, conns, cfg, stopReason)
}

func pollPrometheus(ctx context.Context, policy v1.RollbackPolicy, onThreshold func(float64), interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			val, err := queryPrometheus(policy.PrometheusURL, policy.MetricQuery)
			if err != nil {
				log.Printf("Prometheus query error: %v", err)
				continue
			}
			if val > policy.Threshold {
				onThreshold(val)
				return
			}
		}
	}
}

type promQueryResult struct {
	Data struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Value [2]interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

func queryPrometheus(baseURL, query string) (float64, error) {
	endpoint := fmt.Sprintf("%s/api/v1/query", baseURL)
	params := url.Values{"query": {query}}
	resp, err := http.Get(endpoint + "?" + params.Encode())
	if err != nil {
		return 0, fmt.Errorf("prometheus request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("reading prometheus response: %w", err)
	}

	var result promQueryResult
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("parsing prometheus response: %w", err)
	}

	if len(result.Data.Result) == 0 {
		return 0, nil
	}

	raw, ok := result.Data.Result[0].Value[1].(string)
	if !ok {
		return 0, fmt.Errorf("unexpected metric value type")
	}

	var val float64
	if _, err := fmt.Sscanf(raw, "%f", &val); err != nil {
		return 0, fmt.Errorf("parsing metric value %q: %w", raw, err)
	}
	return val, nil
}

// stopExperiments calls the Stop RPC on every targeted pod for each enabled service.
func stopExperiments(cr *v1.Obzev0Resource, conns []*PodConnection, cfg GrpcServiceConfig, reason string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	startTime := time.Now()
	podNames := podNameList(conns)

	for _, pc := range conns {
		if cfg.LatencyConfig.Enabled {
			client := latencyproto.NewLatencyServiceClient(pc.Conn)
			if _, err := client.StopTcpServer(ctx, &latencyproto.StopRequest{Reason: reason}); err != nil {
				log.Printf("StopTcpServer error on %s: %v", pc.PodName, err)
			}
		}
		if cfg.TcAConfig.Enabled {
			client := tcaproto.NewTcAnalyserServiceClient(pc.Conn)
			if _, err := client.StopUserSpace(ctx, &tcaproto.StopRequest{Reason: reason}); err != nil {
				log.Printf("StopUserSpace error on %s: %v", pc.PodName, err)
			}
		}
		if cfg.PctmConfig.Enabled {
			client := pcaproto.NewPacketManipulationServiceClient(pc.Conn)
			if _, err := client.StopManipulationProxy(ctx, &pcaproto.StopRequest{Reason: reason}); err != nil {
				log.Printf("StopManipulationProxy error on %s: %v", pc.PodName, err)
			}
		}
	}

	if cr != nil {
		emitEvent(cr, corev1.EventTypeNormal, "ExperimentStopped",
			"Experiment stopped: reason=%s pods=%v", reason, podNames)

		publishReport(cr, ExperimentReport{
			ExperimentName:  cr.Name,
			Namespace:       cr.Namespace,
			StartTime:       startTime,
			EndTime:         time.Now(),
			DurationSeconds: time.Since(startTime).Seconds(),
			AffectedPods:    podNames,
			StopReason:      reason,
		})
	}
}

// ── Legacy helpers kept for compatibility ─────────────────────────────────────

func handleAdd(obj interface{}, conn *grpc.ClientConn) {
	CheckConnection(conn)
	obz, ok := obj.(*v1.Obzev0Resource)
	if !ok {
		klog.Errorf("Error converting object to Obzev0Resource: %v", obj)
		return
	}
	processCustomResource(obz, []*PodConnection{{Conn: conn}})
}

func handleUpdate(newObj interface{}, conn *grpc.ClientConn) {
	CheckConnection(conn)
	obz, ok := newObj.(*v1.Obzev0Resource)
	if !ok {
		klog.Errorf("Error converting object to Obzev0Resource: %v", newObj)
		return
	}
	processCustomResource(obz, []*PodConnection{{Conn: conn}})
}

func handleDelete(obj interface{}) {
	obz, ok := obj.(*v1.Obzev0Resource)
	if !ok {
		klog.Errorf("Error converting object to Obzev0Resource: %v", obj)
		return
	}
	klog.Infof("Custom Resource deleted: %s/%s", obz.GetNamespace(), obz.GetName())
}
