/*
Copyright 2022.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ServiceResourceMapSpec defines the desired state of ServiceResourceMap
type ServiceResourceMapSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	ServiceKindReference ServiceKindReference `json:"service_kind_reference"`
	ServiceMap           map[string]string    `json:"service_map"`
}

// ServiceResourceMapStatus defines the observed state of ServiceResourceMap
type ServiceResourceMapStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// ServiceResourceMap is the Schema for the serviceresourcemaps API
type ServiceResourceMap struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServiceResourceMapSpec   `json:"spec,omitempty"`
	Status ServiceResourceMapStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ServiceResourceMapList contains a list of ServiceResourceMap
type ServiceResourceMapList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceResourceMap `json:"items"`
}

type ServiceKindReference struct {
	ApiGroup string `json:"api_group"`
	Kind     string `json:"kind"`
}

func init() {
	SchemeBuilder.Register(&ServiceResourceMap{}, &ServiceResourceMapList{})
}
