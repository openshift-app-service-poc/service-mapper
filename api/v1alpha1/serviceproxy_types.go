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

// ServiceProxySpec defines the desired state of ServiceProxy
type ServiceProxySpec struct {
	ServiceResourceMapRef string         `json:"service_resource_map"`
	ServiceInstance       NamespacedName `json:"service_instance"`
}

type NamespacedName struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// ServiceProxyStatus defines the observed state of ServiceProxy
type ServiceProxyStatus struct {
	Binding ServiceProxyStatusBinding `json:"binding"`
}

type ServiceProxyStatusBinding struct {
	Name string `json:"name"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ServiceProxy is the Schema for the serviceproxies API
type ServiceProxy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServiceProxySpec   `json:"spec,omitempty"`
	Status ServiceProxyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ServiceProxyList contains a list of ServiceProxy
type ServiceProxyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceProxy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ServiceProxy{}, &ServiceProxyList{})
}
