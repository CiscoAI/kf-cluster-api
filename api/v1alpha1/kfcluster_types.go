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
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IMPORTANT
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KfPlatform defines the platform for bootstrapping a KfCluster
type KfPlatform string

// KfClusterFinalizer - finalizer for all resources created by the KfCluster CRD
const (
	KfClusterFinalizer            = "kfcluster.kubeflow.org"
	KfGcp              KfPlatform = "gcp"
	KfGeneric          KfPlatform = "generic"
)

// KfClusterSpec defines the desired state of KfCluster
type KfClusterSpec struct {
	// Important: Run "make" to regenerate code after modifying this file
	Platform      KfPlatform `json:"platform,omitempty"`
	KfVersion     string     `json:"kf_version,omitempty"`
	ConfigMapName string     `json:"config_map_name,omitempty"`
	Apps          []string   `json:"apps,omitempty"`
	Secrets       []string   `json:"secrets,omitempty"`
}

// KfClusterStatus defines the observed state of KfCluster
type KfClusterStatus struct {
	Conditions     []KfClusterCondition `json:"conditions,omitempty"`
	KubeconfigPath string               `json:"kubeconfig_path,omitempty"`
}

// KfClusterCondition defines the possible states for the KfCluster
type KfClusterCondition struct {
	// Important: Run "make" to regenerate code after modifying this file
	State *appsv1.DeploymentConditionType `json:"state,omitempty"`
	Ready bool                            `json:"ready,omitempty"`
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
