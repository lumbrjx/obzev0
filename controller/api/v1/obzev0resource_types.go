/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// +groupName=batch.github.com
// +kubebuilder:object:generate=true
// +kubebuilder:resource:scope=Namespaced,shortName=obz
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
type TcpConfig struct {
	Enabled  bool   `json:"enabled,omitempty"`
	ReqDelay int32  `json:"reqDelay,omitempty"`
	ResDelay int32  `json:"resDelay,omitempty"`
	Server   string `json:"server,omitempty"`
	Client   string `json:"client,omitempty"`
}

type TcAnalyserConfig struct {
	Enabled  bool   `json:"enabled,omitempty"`
	NetIFace string `json:"netIFace,omitempty"`
}

type PacketManipulationConfig struct {
	Enabled         bool   `json:"enabled,omitempty"`
	Server          string `json:"server,omitempty"`
	Client          string `json:"client,omitempty"`
	DurationSeconds int32  `json:"durationSeconds,omitempty"`
	DropRate        string `json:"dropRate,omitempty"`
	CorruptRate     string `json:"corruptRate,omitempty"`
}

// BlastRadius limits which daemon pods receive the chaos experiment.
// Empty fields mean "all" (no filtering applied for that dimension).
type BlastRadius struct {
	// Namespaces restricts chaos to pods running in these namespaces.
	Namespaces []string `json:"namespaces,omitempty"`
	// LabelSelector restricts chaos to pods whose labels match all key/value pairs.
	LabelSelector map[string]string `json:"labelSelector,omitempty"`
	// Percentage is the fraction of matching nodes to target (1-100). 0 means all.
	Percentage int32 `json:"percentage,omitempty"`
}

// RollbackPolicy defines conditions under which a running experiment is automatically stopped.
type RollbackPolicy struct {
	// MaxDurationSeconds stops the experiment after this many seconds. 0 = no limit.
	MaxDurationSeconds int32 `json:"maxDurationSeconds,omitempty"`
	// PrometheusURL is the base URL of a Prometheus instance (e.g. http://prometheus:9090).
	PrometheusURL string `json:"prometheusURL,omitempty"`
	// MetricQuery is a PromQL instant query. If its scalar result exceeds Threshold the experiment stops.
	MetricQuery string `json:"metricQuery,omitempty"`
	// Threshold is the metric value that triggers rollback.
	Threshold float64 `json:"threshold,omitempty"`
	// PollIntervalSeconds controls how often the metric is checked. Defaults to 10.
	PollIntervalSeconds int32 `json:"pollIntervalSeconds,omitempty"`
}

// NetworkChaosConfig configures bandwidth throttling, DNS failure injection, or TCP RST injection.
type NetworkChaosConfig struct {
	Enabled bool `json:"enabled,omitempty"`
	// Bandwidth throttling
	BandwidthRateKbps uint64 `json:"bandwidthRateKbps,omitempty"`
	Interface         string `json:"interface,omitempty"`
	// DNS chaos: mode is "nxdomain", "servfail", or "delay"
	DNSChaosMode  string `json:"dnsChaosMode,omitempty"`
	DNSDelayMs    int32  `json:"dnsDelayMs,omitempty"`
	DNSListenAddr string `json:"dnsListenAddr,omitempty"`
	DNSUpstream   string `json:"dnsUpstream,omitempty"`
	// TCP RST injection
	TCPResetListenAddr string  `json:"tcpResetListenAddr,omitempty"`
	TCPResetRate       float32 `json:"tcpResetRate,omitempty"`
}

// HTTPFaultConfig configures an HTTP fault-injection reverse proxy.
type HTTPFaultConfig struct {
	Enabled    bool    `json:"enabled,omitempty"`
	ListenAddr string  `json:"listenAddr,omitempty"`
	TargetURL  string  `json:"targetURL,omitempty"`
	ErrorRate  float32 `json:"errorRate,omitempty"`
	ErrorCode  int32   `json:"errorCode,omitempty"`
	DelayMs    int32   `json:"delayMs,omitempty"`
	AbortRate  float32 `json:"abortRate,omitempty"`
}

// ExperimentStep is a single phase inside a multi-step workflow.
type ExperimentStep struct {
	// Name is a human-readable label for this step.
	Name string `json:"name"`
	// DurationSeconds is how long this step runs before the next one starts. 0 = fire-and-forget.
	DurationSeconds                 int32                    `json:"durationSeconds,omitempty"`
	LatencyServiceConfig            TcpConfig                `json:"latencySvcConfig,omitempty"`
	TcAnalyserServiceConfig         TcAnalyserConfig         `json:"tcAnalyserSvcConfig,omitempty"`
	PacketManipulationServiceConfig PacketManipulationConfig `json:"packetManipulationSvcConfig,omitempty"`
	NetworkChaosConfig              NetworkChaosConfig       `json:"networkChaosSvcConfig,omitempty"`
	HTTPFaultConfig                 HTTPFaultConfig          `json:"httpFaultSvcConfig,omitempty"`
}

// Obzev0ResourceSpec defines the desired state of Obzev0Resource
type Obzev0ResourceSpec struct {
	LatencyServiceConfig            TcpConfig                `json:"latencySvcConfig,omitempty"`
	TcAnalyserServiceConfig         TcAnalyserConfig         `json:"tcAnalyserSvcConfig,omitempty"`
	PacketManipulationServiceConfig PacketManipulationConfig `json:"packetManipulationSvcConfig,omitempty"`
	NetworkChaosConfig              NetworkChaosConfig       `json:"networkChaosSvcConfig,omitempty"`
	HTTPFaultConfig                 HTTPFaultConfig          `json:"httpFaultSvcConfig,omitempty"`

	// BlastRadius limits the scope of this experiment. Empty = target all daemon pods.
	BlastRadius BlastRadius `json:"blastRadius,omitempty"`
	// RollbackPolicy defines automatic stop conditions for this experiment.
	RollbackPolicy RollbackPolicy `json:"rollbackPolicy,omitempty"`
	// Schedule is a standard cron expression (e.g. "*/5 * * * *"). When set, the experiment
	// runs on this schedule instead of immediately on CR creation.
	Schedule string `json:"schedule,omitempty"`
	// Steps defines a sequential workflow. When set, the top-level service configs are ignored.
	Steps []ExperimentStep `json:"steps,omitempty"`
}

// Obzev0ResourceStatus defines the observed state of Obzev0Resource
type Obzev0ResourceStatus struct {
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Obzev0Resource is the Schema for the obzev0resources API
type Obzev0Resource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   Obzev0ResourceSpec   `json:"spec,omitempty"`
	Status Obzev0ResourceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// Obzev0ResourceList contains a list of Obzev0Resource
type Obzev0ResourceList struct {
	metav1.TypeMeta `                 json:",inline"`
	metav1.ListMeta `                 json:"metadata,omitempty"`
	Items           []Obzev0Resource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Obzev0Resource{}, &Obzev0ResourceList{})
}
