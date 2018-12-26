package service

import (
	redisv1alpha1 "github.com/zdq0394/redis-cluster-operator/pkg/apis/redis/v1alpha1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func generateRedisStatefulset(rc *redisv1alpha1.RedisCluster,
	labels map[string]string, ownerRefs []metav1.OwnerReference) *appsv1beta2.StatefulSet {
	name := generateName(statefulsetNamePrefix, rc.Name)
	serviceName := generateName(svcHeadlessNamePrefix, rc.Name)
	namespace := rc.Namespace

	spec := rc.Spec
	redisImage := spec.Redis.Image
	replicas := spec.Redis.Replicas

	ss := &appsv1beta2.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Namespace:       namespace,
			Labels:          labels,
			OwnerReferences: ownerRefs,
		},
		Spec: appsv1beta2.StatefulSetSpec{
			ServiceName: serviceName,
			Replicas:    &replicas,
			UpdateStrategy: appsv1beta2.StatefulSetUpdateStrategy{
				Type: "RollingUpdate",
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
					Containers: []corev1.Container{
						{
							Name:            "redis",
							Image:           redisImage,
							ImagePullPolicy: "Always",
							Command: []string{
								"redis-server",
							},
							Args: []string{
								"/etc/redis/redis.conf",
								"--protected-mode",
								"no",
							},
							Ports:        getContainerPorts(rc),
							VolumeMounts: getVolumeMounts(rc),
						},
					},
					Volumes: getVolumes(rc),
				},
			},
			VolumeClaimTemplates: getVolumeClaimTemplates(rc, labels, ownerRefs),
		},
	}
	return ss
}

func getVolumes(rc *redisv1alpha1.RedisCluster) []corev1.Volume {
	return []corev1.Volume{
		{
			Name: "redis-conf",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "redis-conf",
					},
				},
			},
		},
	}
}

func getVolumeMounts(rc *redisv1alpha1.RedisCluster) []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      "redis-conf",
			MountPath: "/etc/redis",
		},
		{
			Name:      "redis-data",
			MountPath: "/var/lib/redis",
		},
	}
}

func getContainerPorts(rc *redisv1alpha1.RedisCluster) []corev1.ContainerPort {
	return []corev1.ContainerPort{
		{
			Name:          "redis",
			ContainerPort: 6379,
			Protocol:      corev1.ProtocolTCP,
		},
		{
			Name:          "cluster",
			ContainerPort: 16379,
			Protocol:      corev1.ProtocolTCP,
		},
	}
}

func getVolumeClaimTemplates(rc *redisv1alpha1.RedisCluster,
	labels map[string]string, ownerRefs []metav1.OwnerReference) []corev1.PersistentVolumeClaim {
	storageSize := rc.Spec.Redis.Storage.Size
	storageClassName := rc.Spec.Redis.Storage.StorageClassName
	return []corev1.PersistentVolumeClaim{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:            "redis-data",
				Namespace:       rc.Namespace,
				Labels:          labels,
				OwnerReferences: ownerRefs,
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					"ReadWriteOnce",
				},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse(storageSize),
					},
				},
				StorageClassName: &storageClassName,
			},
		},
	}
}
