package cluster

import (
	redisv1alpha1 "github.com/zdq0394/redis-cluster-operator/pkg/apis/redis/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	bootContainerName = "redis-cluster-boot"
	bootContainerCmd  = "/boot/boot.sh"
)

func generateRedisBootPod(bootImage string, clusterDomain string, rc *redisv1alpha1.RedisCluster,
	labels map[string]string, ownerRefs []metav1.OwnerReference) *corev1.Pod {
	name := generateName(bootPodNamePrefix, rc.Name)
	namespace := rc.Namespace
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Namespace:       namespace,
			Labels:          labels,
			OwnerReferences: ownerRefs,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  bootContainerName,
					Image: bootImage,
					Env: []corev1.EnvVar{
						{
							Name:  "CLUSTER_DOMAIN",
							Value: clusterDomain,
						},
						{
							Name:  "HEADLESS_SERVICE_NAME",
							Value: generateName(svcHeadlessNamePrefix, rc.Name),
						},
						{
							Name:  "NAMESPACE",
							Value: rc.Namespace,
						},
						{
							Name:  "STATEFULSET_NAME",
							Value: generateName(statefulsetNamePrefix, rc.Name),
						},
					},
					Command: []string{
						bootContainerCmd,
					},
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}
}
