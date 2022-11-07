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

// ProxableServiceSpec defines the desired state of ProxableService
type ProxableServiceSpec struct {
	ServiceResourceMapRef NamespacedName `json:"service_resource_map"`
	ServiceInstance       NamespacedName `json:"service_instance"`
}

type NamespacedName struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// ProxableServiceStatus defines the observed state of ProxableService
type ProxableServiceStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// ProxableService is the Schema for the proxableservices API
type ProxableService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProxableServiceSpec   `json:"spec,omitempty"`
	Status ProxableServiceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ProxableServiceList contains a list of ProxableService
type ProxableServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProxableService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProxableService{}, &ProxableServiceList{})
}
