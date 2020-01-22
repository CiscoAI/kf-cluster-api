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

// ExistingK8s defines an existing k8s cluster on which we install Kubeflow
type ExistingK8s struct {
	Kubeconfig string `json:"kubeconfig,omitempty"`
}

// Kind defines a Kubernetes-in-Docker cluster upon which we bootstrap Kubeflow
type Kind struct {
	Kubeconfig string `json:"kubeconfig,omitempty"`
}

// GCE is the Google Compute Engine provider to provision VMs for nodes
type GCE struct {
}
