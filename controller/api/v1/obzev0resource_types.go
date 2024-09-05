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

// Obzev0ResourceSpec defines the desired state of Obzev0Resource
type Obzev0ResourceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Obzev0Resource. Edit obzev0resource_types.go to remove/update
	LatencyServiceConfig            TcpConfig                `json:"latencySvcConfig,omitempty"`
	TcAnalyserServiceConfig         TcAnalyserConfig         `json:"tcAnalyserSvcConfig,omitempty"`
	PacketManipulationServiceConfig PacketManipulationConfig `json:"packetManipulationSvcConfig,omitempty"`
}

// Obzev0ResourceStatus defines the observed state of Obzev0Resource
type Obzev0ResourceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
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
