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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ParentRefSpec struct {
	// Name of the UDPRoute
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}
type BackendRefSpec struct {
	// Name of the UDPRoute
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type RulesSpec struct {
	BackendRefs []BackendRefSpec `json:"backendRefs"`
}

// UDPRouteSpec defines the desired state of UDPRoute
type UDPRouteSpec struct {
	ParentRefs []ParentRefSpec `json:"parentRefs"`
	Rules      []RulesSpec     `json:"rules"`
}

// UDPRouteStatus defines the observed state of UDPRoute
type UDPRouteStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// UDPRoute is the Schema for the udproutes API
type UDPRoute struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UDPRouteSpec   `json:"spec,omitempty"`
	Status UDPRouteStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// UDPRouteList contains a list of UDPRoute
type UDPRouteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UDPRoute `json:"items"`
}

func init() {
	SchemeBuilder.Register(&UDPRoute{}, &UDPRouteList{})
}
