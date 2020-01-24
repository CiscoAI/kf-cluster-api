package kubernetes

import (
	cluster "github.com/CiscoAI/kf-cluster-api/api/v1alpha1"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateStatefulSet(kfCluster *cluster.KfCluster) *v1.StatefulSet {
	secretVolumes := []corev1.Volume{}
	secretVolumeMounts := []corev1.VolumeMount{}
	readOnlyMode := int32(444)
	if len(kfCluster.Spec.Secrets) > 1 {
		for _, secret := range kfCluster.Spec.Secrets {
			volumeSecret := &corev1.SecretVolumeSource{
				SecretName:  secret,
				DefaultMode: &readOnlyMode,
			}
			secretVolumes = append(secretVolumes, corev1.Volume{
				Name: secret,
				VolumeSource: corev1.VolumeSource{
					Secret: volumeSecret,
				},
			})
			secretVolumeMounts = append(secretVolumeMounts, corev1.VolumeMount{
				Name:      secret,
				ReadOnly:  true,
				MountPath: "/etc/" + secret,
			})
		}
	}
	labels := map[string]string{"statefulset": kfCluster.Name}
	labelSelector := &metav1.LabelSelector{MatchLabels: labels}
	containers := []corev1.Container{
		corev1.Container{
			Image:           "ciscoai/kf-clusterctl:v0.1.0",
			ImagePullPolicy: "IfNotPresent",
			Args: []string{
				"kf-clusterctl",
				"create",
			},
			VolumeMounts: secretVolumeMounts,
		},
	}
	replicas := int32(1)
	statefulSet := &v1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kfCluster.Name,
			Namespace: kfCluster.Namespace,
		},
		Spec: v1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: labelSelector,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Containers: containers,
					Volumes:    secretVolumes,
				},
			},
		},
	}
	return statefulSet
}
