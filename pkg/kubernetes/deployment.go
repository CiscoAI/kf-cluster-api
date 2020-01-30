package kubernetes

import (
	cluster "github.com/CiscoAI/kf-cluster-api/api/v1alpha1"
	"github.com/CiscoAI/kf-cluster-api/pkg/version"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateDeployment bootstraps k8s resources needed for a Kubeflow install
func CreateDeployment(kfCluster *cluster.KfCluster) (*v1.Deployment, *corev1.PersistentVolumeClaim) {
	labels := map[string]string{"kfcluster": kfCluster.Name}
	labelSelector := &metav1.LabelSelector{MatchLabels: labels}
	replicas := int32(1)
	kfPodSpec, kfVolumeClaim := createPodSpecAndVolumeClaim(kfCluster)
	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kfCluster.Name,
			Namespace: kfCluster.Namespace,
		},
		Spec: v1.DeploymentSpec{
			Replicas: &replicas,
			Selector: labelSelector,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec:       *kfPodSpec,
			},
		},
	}
	return deployment, kfVolumeClaim
}

func createPodSpecAndVolumeClaim(kfCluster *cluster.KfCluster) (*corev1.PodSpec, *corev1.PersistentVolumeClaim) {
	volumes := []corev1.Volume{}
	volumeMounts := []corev1.VolumeMount{}
	readOnlyMode := int32(444)
	requiredConfigMap := false
	var entrypointScript string
	if len(kfCluster.Spec.Secrets) > 0 {
		for _, secret := range kfCluster.Spec.Secrets {
			volumeSecret := &corev1.SecretVolumeSource{
				SecretName:  secret,
				DefaultMode: &readOnlyMode,
			}
			volumes = append(volumes, corev1.Volume{
				Name: secret,
				VolumeSource: corev1.VolumeSource{
					Secret: volumeSecret,
				},
			})
			volumeMounts = append(volumeMounts, corev1.VolumeMount{
				Name:      secret,
				ReadOnly:  true,
				MountPath: "/etc/" + secret,
			})
		}
	}
	// TODO(swiftdiaries): Programmatically get default StorageClass instead of hard-coding
	defaultStorageClass := "standard"
	resourceReq := make(map[corev1.ResourceName]resource.Quantity)
	resourceSize := int64(10)
	resourceReq[corev1.ResourceStorage] = *resource.NewQuantity(resourceSize, "Gi")
	defaultVolumeClaim := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kfCluster.Name,
			Namespace: kfCluster.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &defaultStorageClass,
			AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteMany},
			Resources: corev1.ResourceRequirements{
				Requests: resourceReq,
			},
		},
	}
	defaultVolume := corev1.Volume{
		Name: kfCluster.Name,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: defaultVolumeClaim.Name,
				ReadOnly:  false,
			},
		},
	}
	defaultVolumeMount := corev1.VolumeMount{
		Name:      kfCluster.Name,
		ReadOnly:  false,
		MountPath: "/mnt/volume/",
	}
	volumes = append(volumes, defaultVolume)
	volumeMounts = append(volumeMounts, defaultVolumeMount)
	if kfCluster.Spec.Platform == cluster.KfGcp {
		entrypointScript = "/gcp_entrypoint.sh"
	} else if kfCluster.Spec.Platform == cluster.KfGeneric {
		entrypointScript = "/generic_entrypoint.sh"
	}
	containers := []corev1.Container{
		corev1.Container{
			Name:            kfCluster.Name,
			Image:           "ciscoai/kf-clusterctl:" + version.Version,
			ImagePullPolicy: "Always",
			Command:         []string{"sh"},
			Args:            []string{entrypointScript},
			EnvFrom: []corev1.EnvFromSource{
				corev1.EnvFromSource{
					ConfigMapRef: &corev1.ConfigMapEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: kfCluster.Spec.ConfigMapName,
						},
						Optional: &requiredConfigMap,
					},
				},
			},
			VolumeMounts: volumeMounts,
		},
	}
	podSpec := &corev1.PodSpec{
		Containers: containers,
		Volumes:    volumes,
	}
	return podSpec, defaultVolumeClaim
}
