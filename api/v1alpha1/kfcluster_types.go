/*

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

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KfClusterSpec defines the desired state of KfCluster
type KfClusterSpec struct {
	// Important: Run "make" to regenerate code after modifying this file
	ClusterProvider KfClusterProvider `json:"clusterprovider,omitempty"`
	KfConfig        string            `json:"kf_config,omitempty"`
	Version         string            `json:"version,omitempty"`
	BuildKfctl      bool              `json:"build_kfctl,omitempty"`
}

// KfClusterProvider defines the desired cluster provider from where we source the Kubernetes nodes
type KfClusterProvider struct {
	// Important: Run "make" to regenerate code after modifying this file
	Existing ExistingK8s `json:"existing,omitempty"`
	Gce      GCE         `json:"gce,omitempty"`
}

// KfClusterStatus defines the observed state of KfCluster
type KfClusterStatus struct {
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// KfCluster is the Schema for the kfclusters API
type KfCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KfClusterSpec   `json:"spec,omitempty"`
	Status KfClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KfClusterList contains a list of KfCluster
type KfClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KfCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KfCluster{}, &KfClusterList{})
}
