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
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var kfclusterlog = logf.Log.WithName("kfcluster-resource")

func (r *KfCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-cluster-kubeflow-org-v1alpha1-kfcluster,mutating=true,failurePolicy=fail,groups=cluster.kubeflow.org,resources=kfclusters,verbs=create;update,versions=v1alpha1,name=mkfcluster.kb.io

var _ webhook.Defaulter = &KfCluster{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *KfCluster) Default() {
	kfclusterlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-cluster-kubeflow-org-v1alpha1-kfcluster,mutating=false,failurePolicy=fail,groups=cluster.kubeflow.org,resources=kfclusters,versions=v1alpha1,name=vkfcluster.kb.io

var _ webhook.Validator = &KfCluster{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *KfCluster) ValidateCreate() error {
	kfclusterlog.Info("validate create", "name", r.Name)
	if r.Spec.Platform == "gcp" || r.Spec.Platform == "generic" {
		return nil
	}
	return fmt.Errorf("Invalid platform type. Please enter one of 'gcp', 'ccp' or 'generic'")
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *KfCluster) ValidateUpdate(old runtime.Object) error {
	kfclusterlog.Info("validate update", "name", r.Name)
	if r.Spec.Platform == "gcp" || r.Spec.Platform == "metal" {
		return nil
	}
	return fmt.Errorf("Invalid platform type. Please enter one of 'gcp' or 'metal")
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *KfCluster) ValidateDelete() error {
	kfclusterlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
