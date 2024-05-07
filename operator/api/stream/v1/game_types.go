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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GameSpec defines the desired state of Game
type GameSpec struct {
	// Name of the game
	Name string `json:"name"`
	// ExecutableURL is the URL of the game executable, which will be downloaded and mounted into the container
	ExecutableURL string `json:"executableURL"`
}

// GameStatus defines the observed state of Game

type GameStatus struct {
	URL string `json:"url"`
}

//+kubebuilder:printcolumn:name="URL",type="string",JSONPath=`.status.url`
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Game is the Schema for the games API
type Game struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GameSpec   `json:"spec,omitempty"`
	Status GameStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GameList contains a list of Game
type GameList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Game `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Game{}, &GameList{})
}
