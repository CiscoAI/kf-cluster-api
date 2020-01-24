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

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cluster "github.com/CiscoAI/kf-cluster-api/api/v1alpha1"
	k8s "github.com/CiscoAI/kf-cluster-api/pkg/kubernetes"
)

// KfClusterReconciler reconciles a KfCluster object
type KfClusterReconciler struct {
	client.Client
	Log    logr.Logger     `json:"log,omitempty"`
	Scheme *runtime.Scheme `json:"scheme,omitempty"`
}

// +kubebuilder:rbac:groups=cluster.kubeflow.org,resources=kfclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cluster.kubeflow.org,resources=kfclusters/status,verbs=get;update;patch

func (r *KfClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("kfcluster", req.NamespacedName)

	kfCluster := &cluster.KfCluster{}
	if err := r.Client.Get(ctx, req.NamespacedName, kfCluster); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		//logger.Error(err, "Unable to find KfCluster resource")
		return ctrl.Result{}, err
	}
	if kfCluster.Spec.Platform == cluster.KfGcp {
		err := reconcileGcp(ctx, kfCluster)
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	if kfCluster.Spec.Platform == cluster.KfMetal {
		err := reconcileMetal(ctx, kfCluster)
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

func reconcileGcp(ctx context.Context, kfCluster *cluster.KfCluster) error {
	statefulSet := k8s.CreateStatefulSet(kfCluster)
	if statefulSet == nil {
		return fmt.Errorf("stateful set wasn't created")
	}
	return nil
}

func reconcileMetal(ctx context.Context, kfCLuster *cluster.KfCluster) error {
	return nil
}

func (r *KfClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	err := ctrl.NewControllerManagedBy(mgr).
		For(&cluster.KfCluster{}).
		Complete(r)
	if err != nil {
		return err
	}
	return nil
}
