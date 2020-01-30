package kubeflow

import (
	cluster "github.com/CiscoAI/kf-cluster-api/api/v1alpha1"
)

// InstallKubeflow - takes in the kubernetes cluster created by the infra provider
// and installs kubeflow on it
func InstallKubeflow(kfCluster *cluster.KfCluster) error {
	return nil
}

// DeleteKubeflow - takes in a KfCluster and deletes KF components on it
func DeleteKubeflow(kfCluster *cluster.KfCluster) error {

	return nil
}

// UpgradeKubeflow - takes in a Kubeflow cluster and upgrades Kubeflow on it
func UpgradeKubeflow(kfCluster *cluster.KfCluster) error {

	return nil
}
