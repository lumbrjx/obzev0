# obzev0 — New Features

This document covers all features added across two implementation phases. The first phase added **Blast Radius Controls**, **Auto-Rollback / Circuit Breaker**, and **Chaos Workflows + Scheduling**. The second phase added **K8s Event Logging**, **Chaos Report Generation**, **Helm Chart fix**, **Network Chaos** (bandwidth, DNS, TCP RST), **HTTP Fault Injection**, a **Web UI Dashboard**, and the **obzevMini TUI**. Every change is explained at the code level.

---

## Table of Contents

1. [Overview of changed files](#overview-of-changed-files)
2. [Feature 1 — Blast Radius Controls](#feature-1--blast-radius-controls)
3. [Feature 2 — Auto-Rollback / Circuit Breaker](#feature-2--auto-rollback--circuit-breaker)
4. [Feature 3 — Chaos Workflows + Scheduling](#feature-3--chaos-workflows--scheduling)
5. [Proto layer changes](#proto-layer-changes)
6. [Daemon service changes](#daemon-service-changes)
7. [Pre-existing bug fix](#pre-existing-bug-fix)
8. [Unit tests](#unit-tests)
9. [Feature 9 — Helm Chart Fix](#feature-9--helm-chart-fix)
10. [Feature 11 — Kubernetes Event Logging](#feature-11--kubernetes-event-logging)
11. [Feature 10 — Chaos Report Generation](#feature-10--chaos-report-generation)
12. [Feature 5 — Network Chaos (bandwidth, DNS, TCP RST)](#feature-5--network-chaos-bandwidth-dns-tcp-rst)
13. [Feature 6 — HTTP Fault Injection](#feature-6--http-fault-injection)
14. [Feature 7 — Web UI Dashboard](#feature-7--web-ui-dashboard)
15. [Feature 8 — obzevMini TUI](#feature-8--obzevmini-tui)
16. [How to use — YAML examples](#how-to-use--yaml-examples)

---

## Overview of changed files

| File | What changed |
|---|---|
| `controller/api/v1/obzev0resource_types.go` | Added `BlastRadius`, `RollbackPolicy`, `ExperimentStep` types; extended `Obzev0ResourceSpec` |
| `controller/api/v1/zz_generated.deepcopy.go` | Updated deep-copy methods for the new types |
| `controller/internal/controller/obzev0resource_controller.go` | `gRPCConnections` now stores `*PodConnection`; `handleCRDelete` stops cron schedulers |
| `controller/internal/controller/infx.go` | Entirely new: blast-radius filter, workflow engine, cron scheduler, rollback monitor |
| `common/proto/latency/latency.proto` | Added `StopTcpServer` RPC + `StopRequest`/`StopResponse` messages |
| `common/proto/packetManipulation/packetManipulation.proto` | Added `StopManipulationProxy` RPC |
| `common/proto/tcAnalyser/tcAnalyser.proto` | Added `StopUserSpace` RPC |
| `common/proto/latency/*.pb.go` | Regenerated from proto |
| `common/proto/packetManipulation/*.pb.go` | Regenerated from proto |
| `common/proto/tcAnalyser/*.pb.go` | Regenerated from proto |
| `daemon/api/grpc/latency/service.go` | Added `StopTcpServer`, `activeCancel atomic.Value`, context-based cancellation |
| `daemon/api/grpc/latency/latencySvc.go` | `LaunchTcp` now accepts `context.Context`; uses `LatencyInternalConfig` |
| `daemon/api/grpc/packetManipulation/service.go` | Added `StopManipulationProxy`, `activeCancel`, context cancellation |
| `daemon/api/grpc/packetManipulation/proxy.go` | `Proxy` now accepts `context.Context` |
| `daemon/api/grpc/tcAnalyser/service.go` | Added `StopUserSpace`, `activeCancel`, context cancellation |
| `daemon/api/grpc/tcAnalyser/bpfManager.go` | `bpfLoader` now accepts `context.Context`; shutdown listens on ctx + OS signal |
| `common/definitions/yamlConfig.go` | Added `DelaysConfig`, `ServerConfig`, `ClientConfig`, `LatencyInternalConfig` |
| `controller/internal/controller/blast_radius_test.go` | New: 6 unit tests for `filterConnections` |
| `chart/values.yaml` | Replaced `repository: nginx` placeholder with proper daemon/controller image structure |
| `chart/templates/daemonset.yaml` | Image pulled from `{{ .Values.daemon.image.repository }}:{{ .Values.daemon.image.tag }}` |
| `chart/templates/controller.yaml` | Image pulled from `{{ .Values.controller.image.repository }}:{{ .Values.controller.image.tag }}` |
| `controller/cmd/main.go` | Added `--ui-port` flag; creates event recorder; wires `StartAPIServer` |
| `controller/internal/controller/obzev0resource_controller.go` | Added package-level `eventRecorder` and `ctrlClient` singletons; updated `SetupInformers` signature |
| `controller/internal/controller/infx.go` | Added `emitEvent` helper; all experiment functions now emit K8s events; `publishReport` called at end of each experiment |
| `controller/internal/controller/report.go` | New: `ExperimentReport` struct; `publishReport` writes JSON to klog and creates a ConfigMap |
| `controller/internal/api/server.go` | New: REST API + static-file server for the web dashboard |
| `controller/internal/api/embed.go` | New: `//go:embed ui/dist` binding `uiFS embed.FS` |
| `controller/internal/api/ui/dist/` | Placeholder for the compiled React app (replaced by `npm run build`) |
| `controller/api/v1/obzev0resource_types.go` | Added `NetworkChaosConfig`, `HTTPFaultConfig`; extended `Obzev0ResourceSpec` and `ExperimentStep` |
| `controller/api/v1/zz_generated.deepcopy.go` | DeepCopy methods for `NetworkChaosConfig`, `HTTPFaultConfig`; updated `ExperimentStep` and `Obzev0ResourceSpec` |
| `common/proto/networkChaos/networkChaos.proto` | New: `NetworkChaosService` (bandwidth, DNS chaos, TCP RST) + generated Go |
| `common/proto/httpFault/httpFault.proto` | New: `HTTPFaultService` + generated Go |
| `daemon/api/grpc/networkChaos/` | New package: `bandwidthSvc.go`, `dnsChaos.go`, `tcpResetSvc.go`, `service.go` |
| `daemon/api/grpc/httpFault/` | New package: `proxy.go` (fault transport), `service.go` |
| `daemon/api/grpc/agent.go` | Registered `NetworkChaosService` and `HTTPFaultService` |
| `ui/` | New React + TypeScript + Vite app (source only; build outputs to `controller/internal/api/ui/dist`) |
| `cmd/cli/tui.go` | New: bubbletea 4-state TUI (menu → configure → running → stopped) |
| `cmd/cli/cli.go` | Added `tui` command |
| `cmd/cli/go.mod` | Added `charmbracelet/bubbletea`, `bubbles`, `lipgloss` |

---

## Feature 1 — Blast Radius Controls

**Problem before:** When any CR was created, `handleCREvent` iterated over every single gRPC connection and applied the chaos config to all daemon pods on all nodes. There was no way to scope an experiment.

### New CRD types — `controller/api/v1/obzev0resource_types.go`

```go
type BlastRadius struct {
    Namespaces    []string          `json:"namespaces,omitempty"`
    LabelSelector map[string]string `json:"labelSelector,omitempty"`
    Percentage    int32             `json:"percentage,omitempty"`
}
```

- `Namespaces` — only target daemon pods running in these namespaces. Empty = all.
- `LabelSelector` — only target pods whose labels contain **all** these key/value pairs. Empty = all.
- `Percentage` — after namespace/label filtering, take only this fraction of the matching pods (1–100). `0` means 100%. The selection is deterministic: pods are sorted alphabetically by name before slicing, so repeated CR events always hit the same subset.

`BlastRadius` is added to `Obzev0ResourceSpec`:

```go
type Obzev0ResourceSpec struct {
    // ... existing service configs ...
    BlastRadius BlastRadius `json:"blastRadius,omitempty"`
}
```

### PodConnection — `controller/internal/controller/obzev0resource_controller.go`

The connection registry was changed from `map[string]*grpc.ClientConn` to `map[string]*PodConnection`:

```go
type PodConnection struct {
    Conn      *grpc.ClientConn
    PodName   string
    Namespace string
    Labels    map[string]string
    NodeName  string
}
```

When a pod connects, its metadata is now stored alongside the connection:

```go
connections[address] = &PodConnection{
    Conn:      conn,
    PodName:   pod.Name,
    Namespace: pod.Namespace,
    Labels:    pod.Labels,
    NodeName:  pod.Spec.NodeName,
}
```

### Filtering logic — `controller/internal/controller/infx.go`

`filterConnections` is the core of the feature. It takes the full connection map and the `BlastRadius` spec and returns the filtered slice:

```go
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
        sort.Slice(matched, func(i, j int) bool {
            return matched[i].PodName < matched[j].PodName
        })
        n := int(float64(len(matched)) * float64(br.Percentage) / 100.0)
        if n < 1 { n = 1 }
        matched = matched[:n]
    }

    return matched
}
```

`matchesNamespaces` returns `true` if the pod's namespace is in the list (or the list is empty).
`matchesLabels` returns `true` only if the pod has **every** key/value pair in the selector map.

`handleCREvent` now calls `filterConnections` before dispatching and logs the ratio:

```go
func handleCREvent(obj interface{}, connections map[string]*PodConnection) {
    cr := obj.(*v1.Obzev0Resource)
    targets := filterConnections(connections, cr.Spec.BlastRadius)
    setupLog.Info("Dispatching chaos experiment",
        "totalPods", len(connections),
        "targetedPods", len(targets),
    )
    go processCustomResource(cr, targets)
}
```

---

## Feature 2 — Auto-Rollback / Circuit Breaker

**Problem before:** Once a chaos experiment was dispatched via gRPC, there was no way to stop it. Services ran goroutines indefinitely (or until their own internal timeout). There was no circuit-breaker that could react to degraded metrics.

### New CRD type — `controller/api/v1/obzev0resource_types.go`

```go
type RollbackPolicy struct {
    MaxDurationSeconds  int32   `json:"maxDurationSeconds,omitempty"`
    PrometheusURL       string  `json:"prometheusURL,omitempty"`
    MetricQuery         string  `json:"metricQuery,omitempty"`
    Threshold           float64 `json:"threshold,omitempty"`
    PollIntervalSeconds int32   `json:"pollIntervalSeconds,omitempty"`
}
```

- `MaxDurationSeconds` — hard time limit. After this many seconds the experiment is stopped regardless of metric state.
- `PrometheusURL` — base URL of your Prometheus instance, e.g. `http://prometheus:9090`.
- `MetricQuery` — a PromQL instant query whose scalar result is compared to `Threshold`.
- `Threshold` — if the metric value exceeds this, rollback triggers immediately.
- `PollIntervalSeconds` — how often to check the metric (default: 10 s).

### Rollback monitor — `controller/internal/controller/infx.go`

`monitorAndRollback` runs as a goroutine after each experiment is dispatched. It owns a shared `context.CancelFunc`. Two independent goroutines race to cancel it:

```go
func monitorAndRollback(conns []*PodConnection, policy v1.RollbackPolicy, cfg GrpcServiceConfig) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Arm a timer-based rollback
    if policy.MaxDurationSeconds > 0 {
        go func() {
            time.Sleep(time.Duration(policy.MaxDurationSeconds) * time.Second)
            cancel()
        }()
    }

    // Arm a metric-based rollback
    if policy.PrometheusURL != "" && policy.MetricQuery != "" {
        go pollPrometheus(ctx, policy, cancel, interval)
    }

    <-ctx.Done()         // block until either goroutine fires
    stopExperiments(conns, cfg)
}
```

`pollPrometheus` queries the Prometheus HTTP API on a ticker, parses the scalar result, and calls `cancel()` if the threshold is crossed:

```go
func pollPrometheus(ctx context.Context, policy v1.RollbackPolicy, cancel context.CancelFunc, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            val, err := queryPrometheus(policy.PrometheusURL, policy.MetricQuery)
            if val > policy.Threshold {
                cancel()
                return
            }
        }
    }
}
```

`queryPrometheus` hits `/api/v1/query` using standard `net/http` (no extra dependency):

```go
func queryPrometheus(baseURL, query string) (float64, error) {
    endpoint := fmt.Sprintf("%s/api/v1/query", baseURL)
    params := url.Values{"query": {query}}
    resp, err := http.Get(endpoint + "?" + params.Encode())
    // ... parse JSON, extract scalar value ...
}
```

### Stop RPCs — proto + daemon

To actually stop a running experiment, each gRPC service received a new `Stop*` RPC. `stopExperiments` in the controller calls all of them:

```go
func stopExperiments(conns []*PodConnection, cfg GrpcServiceConfig) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    for _, pc := range conns {
        if cfg.LatencyConfig.Enabled {
            client := pb.NewLatencyServiceClient(pc.Conn)
            client.StopTcpServer(ctx, &latency.StopRequest{Reason: "rollback"})
        }
        if cfg.TcAConfig.Enabled {
            client := tca.NewTcAnalyserServiceClient(pc.Conn)
            client.StopUserSpace(ctx, &tca.StopRequest{Reason: "rollback"})
        }
        if cfg.PctmConfig.Enabled {
            client := pca.NewPacketManipulationServiceClient(pc.Conn)
            client.StopManipulationProxy(ctx, &pca.StopRequest{Reason: "rollback"})
        }
    }
}
```

---

## Feature 3 — Chaos Workflows + Scheduling

**Problem before:** The CR spec was a flat bag of service configs — there was no concept of sequencing ("do A, then B") or recurring runs ("run every 5 minutes").

### New CRD types — `controller/api/v1/obzev0resource_types.go`

```go
type ExperimentStep struct {
    Name                            string                   `json:"name"`
    DurationSeconds                 int32                    `json:"durationSeconds,omitempty"`
    LatencyServiceConfig            TcpConfig                `json:"latencySvcConfig,omitempty"`
    TcAnalyserServiceConfig         TcAnalyserConfig         `json:"tcAnalyserSvcConfig,omitempty"`
    PacketManipulationServiceConfig PacketManipulationConfig `json:"packetManipulationSvcConfig,omitempty"`
}
```

Added to `Obzev0ResourceSpec`:

```go
Schedule string           `json:"schedule,omitempty"`
Steps    []ExperimentStep `json:"steps,omitempty"`
```

- When `Steps` is non-empty, the top-level `latencySvcConfig` / `packetManipulationSvcConfig` / `tcAnalyserSvcConfig` fields are ignored. Each step has its own configs.
- When `Schedule` is set, the experiment is not run immediately on CR creation; it is registered with a cron scheduler and fires according to the expression.
- Both can be combined: `schedule + steps` runs the full workflow on a cron.
- When neither is set, behaviour is identical to the original flat single-shot mode (backward compatible).

### Experiment dispatch — `controller/internal/controller/infx.go`

`processCustomResource` now acts as a router:

```go
func processCustomResource(cr *v1.Obzev0Resource, targets []*PodConnection) {
    if cr.Spec.Schedule != "" {
        scheduleExperiment(cr, targets)   // register with cron, return immediately
        return
    }
    if len(cr.Spec.Steps) > 0 {
        executeWorkflow(cr.Spec.Steps, targets, cr.Spec.RollbackPolicy)
        return
    }
    // legacy flat mode
    cfg := GrpcServiceConfig{ ... }
    for _, pc := range targets {
        callGrpcServices(pc.Conn, cfg)
    }
    if hasRollbackPolicy(cr.Spec.RollbackPolicy) {
        go monitorAndRollback(targets, cr.Spec.RollbackPolicy, cfg)
    }
}
```

### Workflow engine — `controller/internal/controller/infx.go`

`executeWorkflow` iterates steps sequentially. For each step it dispatches gRPC, sleeps for `DurationSeconds`, then calls `stopExperiments` before moving to the next step:

```go
func executeWorkflow(steps []v1.ExperimentStep, conns []*PodConnection, policy v1.RollbackPolicy) {
    for i, step := range steps {
        cfg := GrpcServiceConfig{
            LatencyConfig: step.LatencyServiceConfig,
            TcAConfig:     step.TcAnalyserServiceConfig,
            PctmConfig:    step.PacketManipulationServiceConfig,
        }
        for _, pc := range conns {
            callGrpcServices(pc.Conn, cfg)
        }
        if step.DurationSeconds > 0 {
            time.Sleep(time.Duration(step.DurationSeconds) * time.Second)
            stopExperiments(conns, cfg)   // clean stop before next step starts
        }
    }
}
```

### Cron scheduler — `controller/internal/controller/infx.go`

`scheduleExperiment` uses [`github.com/robfig/cron/v3`](https://github.com/robfig/cron) (added as a dependency). A `sync.Map` keyed by `"namespace/name"` stores one `*cron.Cron` per CR so that updating or deleting a CR correctly replaces or removes its scheduler:

```go
var schedulers sync.Map   // key = "namespace/name", value = *cron.Cron

func scheduleExperiment(cr *v1.Obzev0Resource, conns []*PodConnection) {
    key := cr.Namespace + "/" + cr.Name
    stopScheduler(key)    // stop any previous scheduler for this CR

    c := cron.New()
    c.AddFunc(cr.Spec.Schedule, func() {
        // fires on each cron tick
        if len(cr.Spec.Steps) > 0 {
            executeWorkflow(cr.Spec.Steps, conns, cr.Spec.RollbackPolicy)
        } else {
            // flat mode with optional rollback
        }
    })
    c.Start()
    schedulers.Store(key, c)
}

func stopScheduler(key string) {
    if v, ok := schedulers.LoadAndDelete(key); ok {
        v.(*cron.Cron).Stop()
    }
}
```

When a CR is deleted, `handleCRDelete` (a new handler wired to `crInformer.DeleteFunc`) calls `stopScheduler` to prevent orphaned goroutines.

---

## Proto layer changes

Three messages and three RPCs were added, one per service. The pattern is identical across all three:

```protobuf
// Added to each .proto file
message StopRequest  { string reason = 1; }
message StopResponse { string message = 1; }

// Added to each service definition
rpc StopTcpServer(StopRequest)         returns (StopResponse) {};   // latency
rpc StopManipulationProxy(StopRequest) returns (StopResponse) {};   // packetManipulation
rpc StopUserSpace(StopRequest)         returns (StopResponse) {};   // tcAnalyser
```

All three `.pb.go` / `_grpc.pb.go` / `.pb.validate.go` files were regenerated using the existing `make generate-proto PROTO_PATH=<name>` command.

---

## Daemon service changes

### Pattern applied to all three services

Every service previously launched chaos as a goroutine that could not be stopped. The new pattern is:

1. On `Start*` — create a cancellable context and store its `CancelFunc` in a package-level `atomic.Value`.
2. On `Stop*` — load the `CancelFunc` and call it.
3. The underlying function (`LaunchTcp`, `Proxy`, `bpfLoader`) now accepts `context.Context` and respects cancellation.

#### Latency — `daemon/api/grpc/latency/service.go`

```go
var activeCancel atomic.Value

func (s *LatencyService) StartTcpServer(ctx context.Context, req *proto.RequestForTcp) (*proto.ResponseFromTcp, error) {
    expCtx, cancel := context.WithCancel(context.Background())
    activeCancel.Store(cancel)
    go LaunchTcp(conf, expCtx)
    // ...
}

func (s *LatencyService) StopTcpServer(ctx context.Context, req *proto.StopRequest) (*proto.StopResponse, error) {
    if cancel, ok := activeCancel.Load().(context.CancelFunc); ok {
        cancel()
    }
    return &proto.StopResponse{Message: "latency experiment stopped"}, nil
}
```

#### `LaunchTcp` — `daemon/api/grpc/latency/latencySvc.go`

Signature changed from `LaunchTcp(conf)` to `LaunchTcp(conf, ctx context.Context)`. The internal `context.WithTimeout` now derives from the external context, so whichever fires first (external cancel or 10 s timeout) shuts down the TCP listener:

```go
func LaunchTcp(conf definitions.LatencyInternalConfig, ctx context.Context) error {
    // Combine caller context with built-in 10s timeout
    timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    // ... listener loop checks timeoutCtx.Done() ...
}
```

#### Packet manipulation — `daemon/api/grpc/packetManipulation/proxy.go`

Same pattern. `Proxy` signature changed from `Proxy(conf ProxyConfig)` to `Proxy(conf ProxyConfig, extCtx context.Context)`. The timeout context is derived from the external one:

```go
func Proxy(conf ProxyConfig, extCtx context.Context) error {
    ctx, cancel := context.WithTimeout(extCtx, conf.Timeout)
    defer cancel()
    // ...
}
```

#### BPF analyser — `daemon/api/grpc/tcAnalyser/bpfManager.go`

`bpfLoader` previously blocked on a raw OS signal channel. It now uses a `select` that also listens on `ctx.Done()`:

```go
func bpfLoader(interf string, ctx context.Context) error {
    // ... attach eBPF programs ...

    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

    select {
    case <-sigCh:
        fmt.Println("Received OS signal, cleaning up...")
    case <-ctx.Done():
        fmt.Println("Context cancelled, cleaning up...")
    }
    // ... detach filters and qdisc ...
}
```

This means either a `kubectl delete` of the CR or an OS signal will trigger a clean eBPF detach.

---

## Pre-existing bug fix

The daemon referenced `definitions.DelaysConfig`, `definitions.ServerConfig`, `definitions.ClientConfig`, and `definitions.Config` with `Delays`/`Server`/`Client` sub-fields — none of which existed in `common/definitions/yamlConfig.go`. The daemon was not buildable before this fix.

Added to `common/definitions/yamlConfig.go`:

```go
type DelaysConfig struct {
    ReqDelay int32
    ResDelay int32
}

type ServerConfig struct {
    Port string
}

type ClientConfig struct {
    Port string
}

// LatencyInternalConfig is used by LaunchTcp inside the daemon.
// Separate from LatencySvcConfig to avoid coupling the daemon to the CLI YAML schema.
type LatencyInternalConfig struct {
    Delays DelaysConfig
    Server ServerConfig
    Client ClientConfig
}
```

All usages of `definitions.Config` in the latency package were updated to `definitions.LatencyInternalConfig`.

---

## Unit tests

**File:** `controller/internal/controller/blast_radius_test.go`

Six table-driven tests cover all filtering dimensions of `filterConnections`:

| Test | What it checks |
|---|---|
| `TestFilterConnections_NoFilter` | Empty `BlastRadius` returns all connections |
| `TestFilterConnections_ByNamespace` | Filters to only the specified namespace |
| `TestFilterConnections_ByLabel` | Filters to pods matching all label key/value pairs |
| `TestFilterConnections_Percentage` | 50% of 4 pods = 2 pods |
| `TestFilterConnections_PercentageRoundsUpToAtLeastOne` | 10% of 1 pod = 1 pod (never 0) |
| `TestFilterConnections_NamespaceAndLabel` | Namespace AND label filter combined |

Run with:
```
cd controller && go test ./internal/controller/ -run TestFilter -v
```

---

## How to use — YAML examples

### Blast radius only

```yaml
apiVersion: batch.github.com/v1
kind: Obzev0Resource
metadata:
  name: target-staging-backend
spec:
  blastRadius:
    namespaces:
      - staging
    labelSelector:
      tier: backend
    percentage: 50          # hit only half the matching pods
  latencySvcConfig:
    enabled: true
    reqDelay: 200
    resDelay: 100
    server: "8080"
    client: "localhost:3000"
```

### Auto-rollback with max duration

```yaml
spec:
  rollbackPolicy:
    maxDurationSeconds: 120   # experiment stops after 2 minutes
  packetManipulationSvcConfig:
    enabled: true
    server: "9090"
    client: "localhost:3001"
    durationSeconds: 300
    dropRate: "0.3"
```

### Auto-rollback on Prometheus metric

```yaml
spec:
  rollbackPolicy:
    prometheusURL: "http://prometheus.monitoring:9090"
    metricQuery: 'rate(http_requests_total{status=~"5.."}[1m])'
    threshold: 0.05           # stop if error rate exceeds 5%
    pollIntervalSeconds: 15
  latencySvcConfig:
    enabled: true
    reqDelay: 500
    resDelay: 500
    server: "8080"
    client: "localhost:3000"
```

### Multi-step workflow

```yaml
spec:
  blastRadius:
    namespaces: [staging]
  steps:
    - name: "warm-up-latency"
      durationSeconds: 30
      latencySvcConfig:
        enabled: true
        reqDelay: 100
        resDelay: 50
        server: "8080"
        client: "localhost:3000"
    - name: "escalate-packet-loss"
      durationSeconds: 60
      packetManipulationSvcConfig:
        enabled: true
        server: "9090"
        client: "localhost:3001"
        durationSeconds: 60
        dropRate: "0.2"
    - name: "combined-stress"
      durationSeconds: 45
      latencySvcConfig:
        enabled: true
        reqDelay: 300
        resDelay: 300
        server: "8080"
        client: "localhost:3000"
      packetManipulationSvcConfig:
        enabled: true
        server: "9090"
        client: "localhost:3001"
        durationSeconds: 45
        dropRate: "0.1"
```

### Recurring experiment on a schedule

```yaml
spec:
  schedule: "0 */2 * * *"   # every 2 hours
  blastRadius:
    percentage: 25            # rotate through 25% of nodes each run
  rollbackPolicy:
    maxDurationSeconds: 300
  steps:
    - name: "periodic-latency-check"
      durationSeconds: 60
      latencySvcConfig:
        enabled: true
        reqDelay: 150
        resDelay: 75
        server: "8080"
        client: "localhost:3000"
```

When the CR is deleted, the scheduler is automatically stopped and no further runs will occur.

### Bandwidth throttle

```yaml
spec:
  blastRadius:
    namespaces: [production]
    percentage: 30
  networkChaosSvcConfig:
    enabled: true
    bandwidthRateKbps: 1000   # cap egress to 1 Mbps
    interface: eth0
  rollbackPolicy:
    maxDurationSeconds: 180
```

### DNS failure injection

```yaml
spec:
  networkChaosSvcConfig:
    enabled: true
    dnsChaosMode: nxdomain      # all DNS queries return NXDOMAIN
    dnsListenAddr: "127.0.0.1:5353"
    dnsUpstream: "8.8.8.8:53"  # used only for mode=delay
```

For `mode: delay`, replace `nxdomain` with `delay` and add `dnsDelayMs: 500`.

### TCP RST injection

```yaml
spec:
  networkChaosSvcConfig:
    enabled: true
    tcpResetListenAddr: ":9080"
    tcpResetRate: 0.4    # 40% of connections receive RST instead of FIN
```

### HTTP fault injection

```yaml
spec:
  httpFaultSvcConfig:
    enabled: true
    listenAddr: ":8880"
    targetURL: "http://backend-service:9090"
    errorRate: 0.2     # 20% of requests get a 503
    errorCode: 503
    delayMs: 100       # every request is slowed by 100ms
    abortRate: 0.05    # 5% of requests are dropped at the connection level
  rollbackPolicy:
    prometheusURL: "http://prometheus.monitoring:9090"
    metricQuery: 'histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[1m]))'
    threshold: 2.0     # stop if p99 latency exceeds 2 s
    pollIntervalSeconds: 10
```

### Multi-step workflow with new chaos types

```yaml
spec:
  blastRadius:
    namespaces: [staging]
  steps:
    - name: "bandwidth-squeeze"
      durationSeconds: 30
      networkChaosSvcConfig:
        enabled: true
        bandwidthRateKbps: 500
        interface: eth0
    - name: "add-latency"
      durationSeconds: 30
      latencySvcConfig:
        enabled: true
        reqDelay: 300
        resDelay: 300
        server: "8080"
        client: "localhost:3000"
    - name: "http-errors"
      durationSeconds: 60
      httpFaultSvcConfig:
        enabled: true
        listenAddr: ":8880"
        targetURL: "http://backend:9090"
        errorRate: 0.3
        errorCode: 500
```

### Checking reports after an experiment

```bash
# List all reports
kubectl get configmaps -l obzev0.io/report=true

# Get the JSON for a specific experiment
kubectl get configmap obzev0-report-my-experiment-1718000000 -o jsonpath='{.data.report\.json}' | jq .

# Or view in the dashboard
kubectl port-forward deployment/obzev0-controller 8090:8090
# open http://localhost:8090 → Reports tab
```

---

## Feature 9 — Helm Chart Fix

**Problem before:** `chart/values.yaml` had `repository: nginx` as a placeholder image. Both the daemonset and controller templates referenced this single value, making the chart impossible to use out of the box.

### `chart/values.yaml`

Replaced with a properly structured values file:

```yaml
daemon:
  image:
    repository: lumbrjx/obzev0-grpc-daemon
    tag: latest
    pullPolicy: IfNotPresent

controller:
  image:
    repository: lumbrjx/obzev0-k8s-controller
    tag: latest
    pullPolicy: IfNotPresent
  uiPort: 8090
```

### `chart/templates/daemonset.yaml`

The image line was changed to:

```yaml
image: "{{ .Values.daemon.image.repository }}:{{ .Values.daemon.image.tag }}"
imagePullPolicy: {{ .Values.daemon.image.pullPolicy }}
```

### `chart/templates/controller.yaml`

Same treatment for the controller deployment:

```yaml
image: "{{ .Values.controller.image.repository }}:{{ .Values.controller.image.tag }}"
imagePullPolicy: {{ .Values.controller.image.pullPolicy }}
```

Override at deploy time with `--set daemon.image.tag=sha-abc123` without touching the chart files.

---

## Feature 11 — Kubernetes Event Logging

**Problem before:** Experiments ran silently. `kubectl describe obzev0resource <name>` showed no Events section, making it impossible to know when an experiment started, advanced through steps, or triggered a rollback.

### `controller/cmd/main.go`

Creates an event recorder and passes it into `SetupInformers`:

```go
recorder := mgr.GetEventRecorderFor("obzev0-controller")
controller.SetupInformers(mgr, recorder, mgr.GetClient())
```

`GetEventRecorderFor` registers the controller as the reporting component in all events, so they appear under `From` in `kubectl describe`.

### `controller/internal/controller/obzev0resource_controller.go`

`SetupInformers` signature extended to receive both the recorder and the controller-runtime client:

```go
var (
    eventRecorder record.EventRecorder
    ctrlClient    client.Client
)

func SetupInformers(mgr ctrl.Manager, recorder record.EventRecorder, c client.Client) {
    eventRecorder = recorder
    ctrlClient = c
    // ... informer wiring unchanged ...
}
```

Both variables are package-level singletons written once at startup and read by the dispatch functions in `infx.go` and `report.go`.

### `controller/internal/controller/infx.go`

A small helper wraps the recorder to guard against `nil` (useful in tests):

```go
func emitEvent(cr *v1.Obzev0Resource, eventType, reason, msgFmt string, args ...interface{}) {
    if eventRecorder == nil {
        return
    }
    msg := fmt.Sprintf(msgFmt, args...)
    eventRecorder.Eventf(cr, eventType, reason, "%s", msg)
}
```

Events are emitted at these points:

| Call site | Type | Reason | Message |
|---|---|---|---|
| `processCustomResource` start | `Normal` | `ExperimentStarted` | "Chaos experiment started on N pod(s)" |
| `executeWorkflow` each step | `Normal` | `WorkflowStep` | "Step i/n: stepName" |
| `executeWorkflow` finish | `Normal` | `WorkflowComplete` | "Workflow completed N step(s) in X.Xs" |
| `monitorAndRollback` max-duration | `Warning` | `RollbackTriggered` | "Max duration of Ns reached" |
| `monitorAndRollback` metric | `Warning` | `RollbackTriggered` | "Metric value X exceeded threshold Y" |
| `stopExperiments` | `Normal` | `ExperimentStopped` | "Experiment stopped: reason=… pods=[…]" |

Events appear in `kubectl describe obzev0resource <name>` under the Events section and in `kubectl get events`.

---

## Feature 10 — Chaos Report Generation

**Problem before:** After an experiment finished there was no record of what ran, for how long, on which pods, or why it stopped. Debugging required manually correlating log timestamps.

### `controller/internal/controller/report.go` (new file)

Defines the report structure and a `publishReport` function:

```go
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
```

`publishReport` does two things:

```go
func publishReport(cr *v1.Obzev0Resource, report ExperimentReport) {
    data, _ := json.Marshal(report)

    // 1. Structured log line — always emitted
    klog.Infof("ExperimentReport: %s", string(data))

    // 2. Kubernetes ConfigMap — written if ctrlClient is available
    cmName := fmt.Sprintf("obzev0-report-%s-%d", report.ExperimentName, time.Now().Unix())
    cm := &corev1.ConfigMap{
        ObjectMeta: metav1.ObjectMeta{
            Name:      cmName,
            Namespace: report.Namespace,
            Labels: map[string]string{
                "obzev0.io/report":     "true",
                "obzev0.io/experiment": report.ExperimentName,
            },
        },
        Data: map[string]string{"report.json": string(data)},
    }
    ctrlClient.Create(context.Background(), cm, &client.CreateOptions{})
}
```

The ConfigMap name is `obzev0-report-{experimentName}-{unixTimestamp}` to keep multiple reports from the same CR distinct. The label `obzev0.io/report=true` makes them easy to query:

```
kubectl get configmaps -l obzev0.io/report=true
kubectl get configmaps -l obzev0.io/experiment=my-experiment
```

### Integration points

`publishReport` is called from two locations:

- **`stopExperiments`** — called after every rollback or manual stop; records `stopReason`, `affectedPods`, and wall-clock duration.
- **`executeWorkflow`** — called at the end of the last step; additionally records `stepCount`.

---

## Feature 5 — Network Chaos (bandwidth, DNS, TCP RST)

Three new chaos modes were added as a single gRPC service (`NetworkChaosService`) in the daemon, requiring no new OS-level eBPF — they use existing Linux facilities via libraries already in the daemon's `go.mod`.

### New proto — `common/proto/networkChaos/networkChaos.proto`

```protobuf
service NetworkChaosService {
    rpc StartBandwidthLimit(BandwidthRequest)  returns (BandwidthResponse)  {};
    rpc StopBandwidthLimit(StopRequest)        returns (StopResponse)       {};
    rpc StartDNSChaos(DNSChaosRequest)         returns (DNSChaosResponse)   {};
    rpc StopDNSChaos(StopRequest)              returns (StopResponse)       {};
    rpc StartTCPReset(TCPResetRequest)         returns (TCPResetResponse)   {};
    rpc StopTCPReset(StopRequest)              returns (StopResponse)       {};
}
```

Go code generated with `protoc --go_out --go-grpc_out`.

### Bandwidth throttle — `daemon/api/grpc/networkChaos/bandwidthSvc.go`

Uses `vishvananda/netlink` (already a daemon dependency) to install a **TBF (Token Bucket Filter) qdisc** on the target network interface:

```go
func startBandwidth(iface string, rateKbps uint64, burstKb uint32) error {
    link, _ := netlink.LinkByName(iface)
    rateBytes := rateKbps * 128  // kbps → bytes/s
    qdisc := &netlink.Tbf{
        QdiscAttrs: netlink.QdiscAttrs{
            LinkIndex: link.Attrs().Index,
            Handle:    netlink.MakeHandle(1, 0),
            Parent:    netlink.HANDLE_ROOT,
        },
        Rate:   rateBytes,
        Limit:  uint32(rateBytes),
        Buffer: burstKb * 1024,
    }
    return netlink.QdiscAdd(qdisc)
}
```

`stopBandwidth` lists qdiscs on the interface and removes the one at `HANDLE_ROOT`. This restores full-speed networking without a node reboot.

### DNS chaos — `daemon/api/grpc/networkChaos/dnsChaos.go`

Starts a **raw UDP listener** on a configurable address (default `127.0.0.1:5353`) and handles each DNS query according to `mode`:

- `"nxdomain"` — responds with rcode 3 (NXDOMAIN) immediately, without forwarding
- `"servfail"` — responds with rcode 2 (SERVFAIL) immediately
- `"delay"` — sleeps `delayMs` milliseconds, then forwards to the upstream resolver

The DNS packet manipulation is done by directly setting bytes 2–3 of the wire-format DNS message (transaction ID is preserved, QR bit set, rcode injected). No external DNS library dependency.

The listener goroutine is stopped by cancelling a `context.Context` stored in an `atomic.Value` on the service struct.

### TCP RST injection — `daemon/api/grpc/networkChaos/tcpResetSvc.go`

Opens a **TCP listener** on `listenAddr`. For each accepted connection, a random float is compared to `resetRate`:

```go
func resetOrClose(conn net.Conn, resetRate float32) {
    tc := conn.(*net.TCPConn)
    if rand.Float32() < resetRate {
        tc.SetLinger(0)   // SO_LINGER=0 → RST on close instead of FIN
    }
    tc.Close()
}
```

Setting `SO_LINGER` to 0 before `Close()` causes the OS to send a TCP RST segment rather than the normal FIN/FIN-ACK handshake. The client sees `connection reset by peer`.

### Service struct — `daemon/api/grpc/networkChaos/service.go`

`NetworkChaosService` holds three independent `atomic.Value` cancel slots (one per mode) so all three modes can run concurrently on the same daemon pod:

```go
type NetworkChaosService struct {
    proto.UnimplementedNetworkChaosServiceServer
    bandwidthCancel atomic.Value  // stores interface name string
    dnsCancel       atomic.Value  // stores context.CancelFunc
    tcpResetCancel  atomic.Value  // stores context.CancelFunc
}
```

### CRD addition — `controller/api/v1/obzev0resource_types.go`

```go
type NetworkChaosConfig struct {
    Enabled            bool    `json:"enabled,omitempty"`
    BandwidthRateKbps  uint64  `json:"bandwidthRateKbps,omitempty"`
    Interface          string  `json:"interface,omitempty"`
    DNSChaosMode       string  `json:"dnsChaosMode,omitempty"`
    DNSDelayMs         int32   `json:"dnsDelayMs,omitempty"`
    DNSListenAddr      string  `json:"dnsListenAddr,omitempty"`
    DNSUpstream        string  `json:"dnsUpstream,omitempty"`
    TCPResetListenAddr string  `json:"tcpResetListenAddr,omitempty"`
    TCPResetRate       float32 `json:"tcpResetRate,omitempty"`
}
```

Added to both `Obzev0ResourceSpec` (top-level) and `ExperimentStep` (per-step in a workflow).

### Agent registration — `daemon/api/grpc/agent.go`

```go
nc := networkchaos.NetworkChaosService{}
netchaoproto.RegisterNetworkChaosServiceServer(grpcServer, &nc)
```

---

## Feature 6 — HTTP Fault Injection

Injects faults into HTTP traffic by running a **reverse proxy** on the daemon pod. The proxy sits between a client and a real upstream service, introducing configurable errors, delays, and connection aborts.

### New proto — `common/proto/httpFault/httpFault.proto`

```protobuf
service HTTPFaultService {
    rpc StartHTTPFault(HTTPFaultRequest)  returns (HTTPFaultResponse)  {};
    rpc StopHTTPFault(StopRequest)        returns (StopResponse)       {};
}

message HTTPFaultRequest {
    string listen_addr = 1;   // proxy listen, e.g. ":8880"
    string target_url  = 2;   // upstream, e.g. "http://real-service:9090"
    float  error_rate  = 3;   // fraction of requests to return error_code
    int32  error_code  = 4;   // HTTP status (default 500)
    int32  delay_ms    = 5;   // artificial delay on every request
    float  abort_rate  = 6;   // fraction of requests to drop without response
}
```

### Fault transport — `daemon/api/grpc/httpFault/proxy.go`

`httputil.ReverseProxy` is wrapped with a custom `http.RoundTripper`:

```go
type faultTransport struct {
    wrapped   http.RoundTripper
    errorRate float32
    errorCode int
    delayMs   int32
    abortRate float32
}

func (t *faultTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    if t.abortRate > 0 && rand.Float32() < t.abortRate {
        return nil, errors.New("injected abort")  // client receives a TCP error
    }
    if t.delayMs > 0 {
        time.Sleep(time.Duration(t.delayMs) * time.Millisecond)
    }
    if t.errorRate > 0 && rand.Float32() < t.errorRate {
        return &http.Response{StatusCode: t.errorCode, ...}, nil  // synthetic error response
    }
    return t.wrapped.RoundTrip(req)  // normal upstream call
}
```

The three fault modes are independent and additive:
- **error_rate** — a percentage of requests get back a synthetic HTTP error without ever reaching upstream.
- **delay_ms** — every request is slowed by this many milliseconds before forwarding.
- **abort_rate** — a percentage of requests receive a connection-level error (not an HTTP response at all).

### Service struct — `daemon/api/grpc/httpFault/service.go`

`StartHTTPFault` builds the proxy and starts a `http.Server` in a goroutine. The server pointer is stored in an `atomic.Value`. `StopHTTPFault` loads it and calls `Shutdown`:

```go
func (s *HTTPFaultService) StartHTTPFault(_ context.Context, req *proto.HTTPFaultRequest) (*proto.HTTPFaultResponse, error) {
    s.stopServer()   // replace any already-running proxy

    proxy, _ := newFaultProxy(req.TargetUrl, req.ErrorRate, req.ErrorCode, req.DelayMs, req.AbortRate)
    srv := &http.Server{Addr: req.ListenAddr, Handler: proxy}
    s.activeServer.Store(srv)
    go srv.ListenAndServe()
    // ...
}

func (s *HTTPFaultService) stopServer() {
    if v := s.activeServer.Swap(nil); v != nil {
        v.(*http.Server).Shutdown(context.Background())
    }
}
```

### CRD addition

```go
type HTTPFaultConfig struct {
    Enabled    bool    `json:"enabled,omitempty"`
    ListenAddr string  `json:"listenAddr,omitempty"`
    TargetURL  string  `json:"targetURL,omitempty"`
    ErrorRate  float32 `json:"errorRate,omitempty"`
    ErrorCode  int32   `json:"errorCode,omitempty"`
    DelayMs    int32   `json:"delayMs,omitempty"`
    AbortRate  float32 `json:"abortRate,omitempty"`
}
```

Added to both `Obzev0ResourceSpec` and `ExperimentStep`.

---

## Feature 7 — Web UI Dashboard

An embedded React single-page application served directly by the controller process. No separate deployment needed — the dashboard is accessible via `kubectl port-forward` on the same port as the REST API.

### REST API — `controller/internal/api/server.go`

`StartAPIServer` registers five routes on a plain `http.ServeMux`:

| Route | Method | Handler |
|---|---|---|
| `/api/experiments` | GET | Lists all `Obzev0Resource` CRs across all namespaces |
| `/api/pods` | GET | Lists pods with label `app=grpc-server` (daemon pods) |
| `/api/experiments/{ns}/{name}` | DELETE | Deletes the CR, stopping the experiment |
| `/api/reports` | GET | Lists ConfigMaps with label `obzev0.io/report=true` |
| `/api/healthz` | GET | Returns `200 ok` |

All responses are JSON with `Access-Control-Allow-Origin: *` so the Vite dev server (`npm run dev`) can call the real controller during development.

```go
func StartAPIServer(mgr ctrl.Manager, port int) {
    mux := http.NewServeMux()
    c := mgr.GetClient()
    mux.HandleFunc("/api/experiments", func(w http.ResponseWriter, r *http.Request) {
        listExperiments(w, r, c)
    })
    // ... other routes ...
    mux.Handle("/", buildUIHandler())
    go http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}
```

### Static file embedding — `controller/internal/api/embed.go`

```go
//go:embed ui/dist
var uiFS embed.FS
```

The `ui/dist` path is **relative to this file** (inside the `api` package directory). Vite's `outDir` is configured to write there so that `go build` picks up the compiled assets automatically.

`buildUIHandler` uses `fs.Sub` to strip the `ui/dist` prefix, then wraps a standard `http.FileServer`. Unknown paths fall back to `index.html` so React Router's client-side routing works:

```go
func buildUIHandler() http.Handler {
    sub, err := fs.Sub(uiFS, "dist")
    if err != nil {
        // No dist built yet — serve a plain status page with API links.
        return http.HandlerFunc(placeholderPage)
    }
    fileServer := http.FileServer(http.FS(sub))
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if _, err := sub.Open(strings.TrimPrefix(r.URL.Path, "/")); err != nil {
            r.URL.Path = "/"   // SPA fallback
        }
        fileServer.ServeHTTP(w, r)
    })
}
```

If `ui/dist` does not exist at build time (e.g. fresh clone without running `npm run build`), the embed directive falls back gracefully and the placeholder HTML page is served instead.

### React app — `ui/`

Built with **Vite + React + TypeScript**. `vite.config.ts` sets `outDir` to `../controller/internal/api/ui/dist` so the build output lands exactly where the Go embed directive expects it.

Three components:

- **`ExperimentCard`** — renders name, namespace, schedule, step count, blast radius %, rollback indicator, and a **Stop** button that calls `DELETE /api/experiments/{ns}/{name}`.
- **`PodList`** — renders a table of daemon pods with phase colour-coding (green = Running, yellow = Pending, red = Failed).
- **`ReportList`** — renders parsed experiment reports, showing duration, stop reason, and affected pods.

`App.tsx` polls all three API endpoints every 5 seconds using `setInterval`. On first mount it fetches immediately, then repeats. A tab bar switches between the three views.

```
cd ui
npm install
npm run build   # writes to controller/internal/api/ui/dist
# or for development:
npm run dev     # hot-reload proxy, API calls go to the real controller
```

### Controller wiring — `controller/cmd/main.go`

```go
flag.IntVar(&uiPort, "ui-port", 8090, "Port for the obzev0 web dashboard and REST API")
// ...
ctrapi.StartAPIServer(mgr, uiPort)
```

Access the dashboard:

```
kubectl port-forward deployment/obzev0-controller 8090:8090
# open http://localhost:8090
```

---

## Feature 8 — obzevMini TUI

An interactive terminal UI for launching and stopping chaos experiments directly from the CLI, without writing YAML. Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

### State machine — `cmd/cli/tui.go`

The TUI has four states:

```
stateMenu → (enter) → stateConfigure → (enter on last field) → stateRunning → (s/q) → stateStopped → (any key) → stateMenu
```

| State | What the user sees |
|---|---|
| `stateMenu` | Arrow-key menu of 6 chaos types |
| `stateConfigure` | Text input fields for the selected type |
| `stateRunning` | Spinner with "Started …" message; `s` or `q` stops |
| `stateStopped` | "✓ Stopped …" confirmation; any key returns to menu |

### Chaos types and their config fields

| Type | Fields |
|---|---|
| Latency Injection | gRPC address, request delay (ms), response delay (ms) |
| Packet Drop / Corrupt | gRPC address, drop rate, corrupt rate, duration |
| TCP RST Reset | gRPC address, listen address, reset rate |
| DNS Chaos | gRPC address, mode, DNS listen address, upstream |
| Bandwidth Throttle | gRPC address, interface, rate (kbps) |
| HTTP Fault Injection | gRPC address, listen address, target URL, error rate, error code, delay |

### Bubbletea model

```go
type tuiModel struct {
    state     tuiState
    cursor    int              // menu cursor
    chaosIdx  int              // selected chaos type
    inputs    []textinput.Model
    focusIdx  int
    spinner   spinner.Model
    statusMsg string
}
```

`Init()` returns `nil` (no initial command). `Update()` dispatches on `state` and delegates to per-state handlers. `View()` renders the current state using `lipgloss` styles (Catppuccin-inspired dark palette).

### Running the TUI

```
cd cmd/cli
go run . tui
```

Or after `go install`:

```
obzev0 tui
```
