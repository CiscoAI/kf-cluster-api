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

	cluster "github.com/CiscoAI/kf-cluster-api/api/v1alpha1"
	"github.com/CiscoAI/kf-cluster-api/pkg/kubernetes"
	"github.com/go-logr/logr"
	reconcilehelper "github.com/kubeflow/kubeflow/components/common/reconcilehelper"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// KfClusterReconciler reconciles a KfCluster object
type KfClusterReconciler struct {
	client.Client
	Log    logr.Logger     `json:"log,omitempty"`
	Scheme *runtime.Scheme `json:"scheme,omitempty"`
}

// +kubebuilder:rbac:groups=cluster.kubeflow.org,resources=kfclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cluster.kubeflow.org,resources=kfclusters/status,verbs=get;update;patch

// Reconcile - reconciles the KfCluster object
func (r *KfClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("kfcluster", req.NamespacedName)

	kfCluster := &cluster.KfCluster{}
	if err := r.Client.Get(ctx, req.NamespacedName, kfCluster); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Info("Error getting KfCluster resource")
		return ctrl.Result{}, err
	}

	// Provision cluster resources and get a kubernetes cluster
	// Store kubeconfig in the associated PV
	if kfCluster.Spec.Platform == cluster.KfGcp {
		err := r.reconcileGcp(ctx, kfCluster, log)
		if err != nil {
			log.Info("Error reconciling KfCluster on GCP")
			return ctrl.Result{}, err
		}
	} else if kfCluster.Spec.Platform == cluster.KfGeneric {
		err := r.reconcileGeneric(ctx, kfCluster, log)
		if err != nil {
			log.Info("Error reconciling KfCluster on k8s")
			return ctrl.Result{}, err
		}
	}

	// Use the kubeconfig stored in the associated PV
	// Use the platform resources to install Kubeflow
	err := r.reconcileKubeflow(ctx, kfCluster, log)
	if err != nil {
		log.Error(err, "error installing Kubeflow")
	}
	return ctrl.Result{}, nil
}

func (r *KfClusterReconciler) reconcileKubeflow(ctx context.Context, kfCluster *cluster.KfCluster, log logr.Logger) error {

	return nil
}

func (r *KfClusterReconciler) reconcileGcp(ctx context.Context, kfCluster *cluster.KfCluster, log logr.Logger) error {
	log.Info("Reconciling KfCluster on GCP")
	deployment, kfVolumeClaim := kubernetes.CreateDeployment(kfCluster)
	if deployment == nil {
		log.Info("Deploymeny spec wasn't generated")
		return fmt.Errorf("error generating deployment spec")
	}
	if kfVolumeClaim == nil {
		log.Info("VolumeClaim spec wasn't generated")
		return fmt.Errorf("error generating volumeclaim spec")
	}
	if err := ctrl.SetControllerReference(kfCluster, deployment, r.Scheme); err != nil {
		log.Info("unable to set controllereference for created deployment")
		return err
	}
	if err := ctrl.SetControllerReference(kfCluster, kfVolumeClaim, r.Scheme); err != nil {
		log.Info("unable to set controllereference for created volumeclaim")
		return err
	}
	justCreatedVolumeClaim := false
	if err := r.Get(ctx, types.NamespacedName{Name: kfVolumeClaim.Name, Namespace: kfVolumeClaim.Namespace}, kfVolumeClaim); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Creating volume claim for KfCluster")
			if err := r.Create(ctx, kfVolumeClaim); err != nil {
				log.Error(err, "Failed to create volume claim")
				return err
			}
			justCreatedVolumeClaim = true
		} else {
			log.Error(err, "error getting volume claim")
			return err
		}
	}
	if !justCreatedVolumeClaim {
		log.Info("volume claim already exists, updating volume claim")
		if err := r.Update(ctx, kfVolumeClaim); err != nil {
			log.Error(err, "error upadting volume claim")
			return err
		}
	}
	// TODO(swiftdiaries): fix deployment reconcile logic; reconcile loop re-creates deployment
	justCreatedDeployment := false
	existingDeployment := &appsv1.Deployment{}
	if err := r.Get(ctx, types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, existingDeployment); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Creating deployment for KfCluster")
			if err := r.Create(ctx, deployment); err != nil {
				log.Error(err, "error creating deployment")
				return err
			}
			kfCluster.Status.KubeconfigPath = "/mnt/volume/" + kfCluster.Name + "/kubeconfig"
			justCreatedDeployment = true
		} else {
			log.Error(err, "error getting deployment")
			return err
		}
	}
	if !justCreatedDeployment && reconcilehelper.CopyDeploymentSetFields(deployment, existingDeployment) {
		log.Info("deployment already exists, update deployment")
		if err := r.Update(ctx, existingDeployment); err != nil {
			log.Error(err, "error upadting deployment")
			return err
		}
	}

	if len(deployment.Status.Conditions) > 0 {
		clusterCondition := cluster.KfClusterCondition{
			State: &deployment.Status.Conditions[0].Type,
			Ready: deployment.Status.Conditions[0].Type == appsv1.DeploymentAvailable,
		}
		clusterConditionLen := len(kfCluster.Status.Conditions)
		if clusterConditionLen == 0 || kfCluster.Status.Conditions[clusterConditionLen-1].State != clusterCondition.State {
			kfCluster.Status.Conditions = append(kfCluster.Status.Conditions, clusterCondition)
		}
		log.Info("cluster condition", clusterCondition)
		err := r.Update(ctx, kfCluster)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *KfClusterReconciler) reconcileGeneric(ctx context.Context, kfCLuster *cluster.KfCluster, log logr.Logger) error {
	log.Info("Reconciling KfCluster on k8s")
	return nil
}

// SetupWithManager registers the controller reconciler logic with the manager binary
func (r *KfClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	err := ctrl.NewControllerManagedBy(mgr).
		For(&cluster.KfCluster{}).
		Complete(r)
	if err != nil {
		return err
	}
	return nil
}
